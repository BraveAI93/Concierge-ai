package intelligence_runtime

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/BraveAI93/concierge-backend/kernel"
)

// IngestContinuity records a local/synthetic interaction as a bounded block and
// source-independent thread relation. It deliberately does not invoke any HTTP
// handler, external account, external action executor, or raw-audio collector.
func (s RuntimeService) IngestContinuity(ctx context.Context, principal AuthenticatedPrincipal, input ContinuityInput) (ContinuityRuntimeResult, error) {
	if err := s.validateV05Composition(); err != nil {
		return ContinuityRuntimeResult{}, err
	}
	if !s.Feature.Enabled {
		return ContinuityRuntimeResult{}, ErrRuntimeDisabled
	}
	if s.Activation != nil {
		enabled, err := s.Activation.Enabled(ctx)
		if err != nil || !enabled {
			return ContinuityRuntimeResult{}, ErrRuntimeDisabled
		}
	}
	binding, err := s.Identity.Resolve(ctx, principal)
	if err != nil {
		return ContinuityRuntimeResult{}, err
	}
	if !binding.AllowsSourceProfile(input.Source.Conversation.ProfileID) {
		return ContinuityRuntimeResult{}, ErrSourceUnauthorized
	}
	if s.Consent == nil {
		return ContinuityRuntimeResult{}, ErrConsentNotVerified
	}
	if err := s.Consent.VerifyConversationDerivedMemory(ctx, binding, input.Source); err != nil {
		return ContinuityRuntimeResult{}, err
	}
	now := s.Clock.Now().UTC()
	if input.IdempotencyKey == "" || input.Block.PersonID != binding.Person.ID || input.Block.Validate() != nil {
		return ContinuityRuntimeResult{}, ErrInvalidRuntimeConfig
	}
	if input.ProposedThread != nil && (input.ProposedThread.PersonID != binding.Person.ID || input.ProposedThread.Validate() != nil) {
		return ContinuityRuntimeResult{}, ErrSourceUnauthorized
	}
	for _, delta := range input.Deltas {
		if delta.PersonID != binding.Person.ID || delta.Validate() != nil {
			return ContinuityRuntimeResult{}, ErrSourceUnauthorized
		}
	}
	for _, signal := range input.Signals {
		if signal.PersonID != binding.Person.ID || signal.Validate() != nil {
			return ContinuityRuntimeResult{}, ErrSourceUnauthorized
		}
	}
	if input.Baseline != nil && (input.Baseline.PersonID != binding.Person.ID || input.Baseline.Validate() != nil) {
		return ContinuityRuntimeResult{}, ErrSourceUnauthorized
	}

	result := ContinuityRuntimeResult{}
	err = s.Repo.RunInTransaction(ctx, binding.Person.ID, func(tx RuntimeTransaction) error {
		if replay, ok := tx.FindReplay(input.IdempotencyKey); ok {
			result = continuityResultFromReplay(replay)
			result.Replayed = true
			return nil
		}
		sourceRecord, sourceErr := continuitySourceRecord(binding, input.Source, now)
		if sourceErr != nil {
			return sourceErr
		}
		if err := tx.StoreSource(sourceRecord); err != nil {
			return err
		}
		thread, created, resolution, err := resolveOrCreateThread(tx, binding.Person.ID, input, s.ThreadResolver, now)
		if err != nil {
			return err
		}
		block := input.Block
		block.ThreadIDs = appendUniqueThreadID(block.ThreadIDs, thread.ID)
		priorBlocks := tx.ListInteractionBlocks(thread.ID)
		if len(priorBlocks) > 0 {
			previous := priorBlocks[len(priorBlocks)-1]
			lastAt := previous.StartTemporal.EventAt
			if previous.EndEventAt != nil {
				lastAt = *previous.EndEventAt
			}
			gap := kernel.InteractionGapState{PersonID: binding.Person.ID, ContextID: firstContextID(block.ContextIDs), LastInteractionAt: lastAt, ObservedAt: now}
			assessment, boundaryErr := s.BoundaryPolicy.AssessBoundary(previous, block, gap, kernel.EvaluationMoment{WallClockAt: now})
			if boundaryErr != nil {
				return boundaryErr
			}
			block.BoundaryEvidence = append(block.BoundaryEvidence, assessment.Evidence...)
		}
		if created {
			if err := tx.SaveThread(thread); err != nil {
				return err
			}
		}
		if err := tx.SaveInteractionBlock(block); err != nil {
			return err
		}
		link := kernel.ContinuityLink{ID: "continuity:block:" + block.ID + ":thread:" + thread.ID, PersonID: binding.Person.ID, Source: kernel.ContinuityRef{Kind: kernel.ContinuityBlock, ID: block.ID}, Target: kernel.ContinuityRef{Kind: kernel.ContinuityThread, ID: thread.ID}, Relation: kernel.ContinuitySameSubject, Why: continuityWhy(resolution, created), EvidenceIDs: block.EvidenceIDs, Provenance: block.Provenance, Confidence: maxResolutionConfidence(resolution, created), Temporal: block.StartTemporal, Freshness: block.Freshness, CreatedAt: now}
		if err := tx.SaveContinuityLink(link); err != nil {
			return err
		}

		allDeltas := tx.ListThreadDeltas(thread.ID)
		for _, delta := range input.Deltas {
			if delta.TargetThreadID != thread.ID {
				continue
			}
			if err := tx.SaveThreadDelta(delta); err != nil {
				return err
			}
			allDeltas = append(allDeltas, delta)
		}
		state := latestThreadState(tx.ListCurrentThreadStates(thread.ID), binding.Person.ID, thread.ID, now)
		state.ID = "thread_state:" + thread.ID + ":" + now.Format("20060102T150405.000000000")
		state.PersonID, state.ThreadID = binding.Person.ID, thread.ID
		state, err = kernel.BuildCurrentThreadState(thread, state, allDeltas, kernel.EvaluationMoment{WallClockAt: now})
		if err != nil {
			return err
		}
		if err := tx.SaveCurrentThreadState(state); err != nil {
			return err
		}

		adaptationID, episodeID, interventionID, err := s.persistAttunement(tx, binding.Person.ID, block, thread, input, now)
		if err != nil {
			return err
		}
		replay := RuntimeResult{PersonID: binding.Person.ID, IdempotencyKey: input.IdempotencyKey, InteractionBlockID: block.ID, ThreadID: thread.ID, ThreadStateID: state.ID, ContinuityLinkID: link.ID, AttunementEpisodeID: episodeID, InterventionID: interventionID, AdaptationID: adaptationID}
		if err := tx.StoreReplay(replay); err != nil {
			return err
		}
		result = continuityResultFromReplay(replay)
		return nil
	})
	if err != nil {
		return ContinuityRuntimeResult{}, err
	}
	return result, nil
}

func (s RuntimeService) persistAttunement(tx RuntimeTransaction, personID string, block kernel.InteractionBlock, thread kernel.Thread, input ContinuityInput, now time.Time) (string, string, string, error) {
	if input.Baseline == nil || len(input.Signals) == 0 {
		return "", "", "", nil
	}
	if err := tx.SavePersonalInteractionBaseline(*input.Baseline); err != nil {
		return "", "", "", err
	}
	for _, signal := range input.Signals {
		if err := tx.SaveObservedInteractionSignal(signal); err != nil {
			return "", "", "", err
		}
	}
	states, err := kernel.InferInteractionStates(input.Signals, *input.Baseline, kernel.BaselineInferencePolicy{MinimumObservations: 3, DeviationMultiplier: 1.5}, kernel.EvaluationMoment{WallClockAt: now})
	if err != nil {
		return "", "", "", err
	}
	for _, state := range states {
		if err := tx.SaveInferredInteractionState(state); err != nil {
			return "", "", "", err
		}
	}
	if len(states) == 0 || input.AttunementControl == kernel.AttunementDisabled {
		return "", "", "", nil
	}
	patterns := tx.ListAttunementPatterns(input.ContextSignature)
	decision, err := s.AttunementPolicy.Decide(input.AttunementControl, states, patterns, kernel.EvaluationMoment{WallClockAt: now})
	if errors.Is(err, kernel.ErrAttunementDisabled) {
		return "", "", "", nil
	}
	if err != nil {
		return "", "", "", err
	}
	decision.ID = "adaptation:" + block.ID
	decision.PersonID, decision.BlockID = personID, block.ID
	if err := kernel.ValidateAttunementDecision(decision); err != nil {
		return "", "", "", err
	}
	if err := tx.SaveInteractionAdaptationDecision(decision); err != nil {
		return "", "", "", err
	}
	episode := kernel.AttunementEpisode{ID: "attunement:" + block.ID, PersonID: personID, BlockID: block.ID, ThreadIDs: []string{thread.ID}, AdaptationDecisionID: decision.ID, CreatedAt: now}
	for _, signal := range input.Signals {
		episode.SignalIDs = append(episode.SignalIDs, signal.ID)
	}
	for _, state := range states {
		episode.InferredStateIDs = append(episode.InferredStateIDs, state.ID)
	}
	if err := tx.SaveAttunementEpisode(episode); err != nil {
		return "", "", "", err
	}
	intervention := kernel.InteractionIntervention{ID: "intervention:" + block.ID, PersonID: personID, EpisodeID: episode.ID, DecisionID: decision.ID, Summary: decision.Reason, OccurredAt: now, Reversible: decision.Reversible, Provenance: []kernel.Provenance{{SourceType: "local_attunement_policy", SourceRef: decision.ID, CapturedAt: now}}}
	if err := tx.SaveInteractionIntervention(intervention); err != nil {
		return "", "", "", err
	}
	return decision.ID, episode.ID, intervention.ID, nil
}

func (s RuntimeService) validateV05Composition() error {
	if s.Identity == nil || s.Repo == nil || s.Clock == nil || s.ThreadResolver == nil || s.AttunementPolicy == nil || s.RetrievalDepthPolicy == nil {
		return ErrInvalidRuntimeConfig
	}
	return nil
}

func resolveOrCreateThread(tx RuntimeTransaction, personID string, input ContinuityInput, resolver kernel.SemanticThreadResolver, now time.Time) (kernel.Thread, bool, kernel.ThreadResolution, error) {
	resolution, err := resolver.Resolve(personID, input.Triggers, tx.ListThreads(), kernel.EvaluationMoment{WallClockAt: now})
	if err != nil {
		return kernel.Thread{}, false, resolution, err
	}
	if resolution.SelectedID != "" {
		for _, thread := range tx.ListThreads() {
			if thread.ID == resolution.SelectedID {
				return thread, false, resolution, nil
			}
		}
	}
	if input.ProposedThread == nil {
		return kernel.Thread{}, false, resolution, ErrUnsupportedSource
	}
	// Ambiguous matches never overwrite/merge a candidate. The new thread is
	// source-independent and can be reviewed/linked later through typed edges.
	return *input.ProposedThread, true, resolution, nil
}

func latestThreadState(states []kernel.CurrentThreadState, personID, threadID string, now time.Time) kernel.CurrentThreadState {
	if len(states) > 0 {
		sort.Slice(states, func(i, j int) bool { return states[i].ReconstructedAt.After(states[j].ReconstructedAt) })
		return states[0]
	}
	return kernel.CurrentThreadState{ID: "thread_state:" + threadID + ":baseline", PersonID: personID, ThreadID: threadID, BaselineSummary: "initial continuity state", CurrentSummary: "initial continuity state", FieldValues: map[string]string{}, ReconstructedAt: now, Freshness: kernel.FreshnessState{LastValidatedAt: now, Status: kernel.FreshnessFresh}}
}

func continuityWhy(resolution kernel.ThreadResolution, created bool) string {
	if created {
		return "new source-independent thread proposed because no safe existing match was selected"
	}
	if resolution.RequiresReview {
		return "candidate relation retained without automatic merge"
	}
	return "semantic triggers matched canonical anchors or aliases"
}

func maxResolutionConfidence(resolution kernel.ThreadResolution, created bool) float64 {
	if created {
		return 0.5
	}
	if len(resolution.Candidates) == 0 {
		return 0.5
	}
	return resolution.Candidates[0].Confidence
}

func firstContextID(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func appendUniqueThreadID(values []string, id string) []string {
	for _, existing := range values {
		if existing == id {
			return values
		}
	}
	if id != "" {
		return append(values, id)
	}
	return values
}

func continuitySourceRecord(binding PersonBinding, source ConversationMessage, now time.Time) (SourceRecord, error) {
	conversation, message := source.Conversation, source.Message
	if conversation.ID == "" || message.ID == "" || message.ConversationID != conversation.ID || !binding.AllowsSourceProfile(conversation.ProfileID) || message.CreatedAt.IsZero() {
		return SourceRecord{}, ErrUnsupportedSource
	}
	return SourceRecord{ID: "continuity-source:" + message.ID, PersonID: binding.Person.ID, ProfileID: conversation.ProfileID, ConversationID: conversation.ID, SessionID: conversation.SessionID, MessageID: message.ID, MessageRole: message.Role, Content: message.Content, ConversationAt: conversation.StartedAt, MessageAt: message.CreatedAt, StoredAt: now}, nil
}

func continuityResultFromReplay(replay RuntimeResult) ContinuityRuntimeResult {
	return ContinuityRuntimeResult{PersonID: replay.PersonID, IdempotencyKey: replay.IdempotencyKey, BlockID: replay.InteractionBlockID, ThreadID: replay.ThreadID, ThreadStateID: replay.ThreadStateID, ContinuityLinkID: replay.ContinuityLinkID, AttunementEpisodeID: replay.AttunementEpisodeID, InterventionID: replay.InterventionID, AdaptationID: replay.AdaptationID}
}

// RecordAttunementOutcome closes only the local evidence loop. It does not
// claim causation or execute an action. A caller must provide a synthetic/local
// observed outcome with evidence/provenance before a pattern is updated.
func (s RuntimeService) RecordAttunementOutcome(ctx context.Context, principal AuthenticatedPrincipal, outcome kernel.InteractionOutcome, seed kernel.PersonalAttunementPattern) (kernel.PersonalAttunementPattern, error) {
	if err := s.validateV05Composition(); err != nil {
		return kernel.PersonalAttunementPattern{}, err
	}
	if !s.Feature.Enabled {
		return kernel.PersonalAttunementPattern{}, ErrRuntimeDisabled
	}
	binding, err := s.Identity.Resolve(ctx, principal)
	if err != nil {
		return kernel.PersonalAttunementPattern{}, err
	}
	if outcome.PersonID != binding.Person.ID || seed.PersonID != binding.Person.ID || outcome.Privacy != kernel.PrivacyOutcomeEvidence {
		return kernel.PersonalAttunementPattern{}, ErrSourceUnauthorized
	}
	now := s.Clock.Now().UTC()
	var updated kernel.PersonalAttunementPattern
	err = s.Repo.RunInTransaction(ctx, binding.Person.ID, func(tx RuntimeTransaction) error {
		if err := tx.SaveInteractionOutcome(outcome); err != nil {
			return err
		}
		pattern := seed
		for _, existing := range tx.ListAttunementPatterns(outcome.ContextSignature) {
			if existing.StrategyFingerprint == seed.StrategyFingerprint && existing.Hypothesis == seed.Hypothesis {
				pattern = existing
				break
			}
		}
		updated, err = kernel.UpdateAttunementPattern(pattern, outcome, kernel.EvaluationMoment{WallClockAt: now})
		if err != nil {
			return err
		}
		return tx.SavePersonalAttunementPattern(updated)
	})
	return updated, err
}

// ResolveContinuity is a domain/API method rather than a route. It returns a
// bounded retrieval plan for the authenticated Person only.
func (s RuntimeService) ResolveContinuity(ctx context.Context, principal AuthenticatedPrincipal, request kernel.RetrievalRequest) (kernel.ThreadResolution, kernel.RetrievalPlan, error) {
	if err := s.validateV05Composition(); err != nil {
		return kernel.ThreadResolution{}, kernel.RetrievalPlan{}, err
	}
	binding, err := s.Identity.Resolve(ctx, principal)
	if err != nil {
		return kernel.ThreadResolution{}, kernel.RetrievalPlan{}, err
	}
	if request.PersonID != binding.Person.ID {
		return kernel.ThreadResolution{}, kernel.RetrievalPlan{}, ErrCrossPersonAccess
	}
	state, err := s.Repo.ReadState(ctx, binding.Person.ID, binding.Person.ID)
	if err != nil {
		return kernel.ThreadResolution{}, kernel.RetrievalPlan{}, err
	}
	resolution, err := s.ThreadResolver.Resolve(binding.Person.ID, request.Triggers, state.Threads, kernel.EvaluationMoment{WallClockAt: request.RequestedAt})
	if err != nil || resolution.SelectedID == "" {
		return resolution, kernel.RetrievalPlan{Depth: kernel.RetrievalAwareness, Reason: "no safely selected thread"}, err
	}
	for _, thread := range state.Threads {
		if thread.ID != resolution.SelectedID {
			continue
		}
		for _, threadState := range state.CurrentThreadStates {
			if threadState.ThreadID == thread.ID {
				request.ThreadID = thread.ID
				plan, planErr := s.RetrievalDepthPolicy.Plan(request, thread, threadState, state.ThreadDeltas, state.InteractionBlocks)
				return resolution, plan, planErr
			}
		}
	}
	return resolution, kernel.RetrievalPlan{}, ErrUnsupportedSource
}

// CompileContinuityAttention is intentionally narrow: support-graph links are
// not visible work until the existing attention policy accepts a MustSurface
// decision. It returns no OpenLoop/notification/action side effect.
func CompileContinuityAttention(link kernel.ContinuityLink, priority kernel.PriorityFactors, deadline kernel.DeadlineFeasibility, policy kernel.ContinuitySurfacePolicy, at time.Time) (kernel.ContinuitySurfaceDecision, error) {
	return kernel.DecideContinuitySurface(link, priority, deadline, policy, kernel.EvaluationMoment{WallClockAt: at})
}
