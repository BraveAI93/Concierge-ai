package kernel

import (
	"sort"
	"strings"
	"time"
)

// BoundaryPolicyConfig is an injected local policy. MinimumGap is deliberately
// a composition setting, not a universal canonical chat timeout.
type BoundaryPolicyConfig struct {
	MinimumGap        time.Duration
	MinimumConfidence float64
}

type DeterministicBoundaryPolicy struct{ Config BoundaryPolicyConfig }

func (p DeterministicBoundaryPolicy) AssessBoundary(previous InteractionBlock, incoming InteractionBlock, gap InteractionGapState, at EvaluationMoment) (InteractionBoundaryAssessment, error) {
	if previous.Validate() != nil || incoming.Validate() != nil || previous.PersonID != incoming.PersonID || at.Validate() != nil || gap.Validate() != nil || gap.PersonID != previous.PersonID || p.Config.MinimumGap < 0 || !isZeroOrUnit(p.Config.MinimumConfidence) {
		return InteractionBoundaryAssessment{}, ErrInvalidInteractionBlock
	}
	elapsed, err := gap.Elapsed(at.WallClockAt)
	if err != nil {
		return InteractionBoundaryAssessment{}, err
	}
	assessment := InteractionBoundaryAssessment{InteractionGap: elapsed, EvaluatedAt: at.WallClockAt}
	for _, evidence := range incoming.BoundaryEvidence {
		if evidence.Confidence >= p.Config.MinimumConfidence {
			assessment.Reasons = append(assessment.Reasons, evidence.Reason)
			assessment.Evidence = append(assessment.Evidence, evidence)
		}
	}
	if p.Config.MinimumGap > 0 && elapsed >= p.Config.MinimumGap {
		assessment.Reasons = append(assessment.Reasons, BoundaryInactivity)
		assessment.Evidence = append(assessment.Evidence, BoundaryEvidence{Reason: BoundaryInactivity, Summary: "configured meaningful interaction gap exceeded", Confidence: p.Config.MinimumConfidence, ObservedAt: at.WallClockAt})
	}
	if previous.SourceType != "" && incoming.SourceType != "" && previous.SourceType != incoming.SourceType {
		assessment.Reasons = append(assessment.Reasons, BoundaryContextChange)
		assessment.Evidence = append(assessment.Evidence, BoundaryEvidence{Reason: BoundaryContextChange, Summary: "source type changed", Confidence: 1, ObservedAt: at.WallClockAt})
	}
	assessment.CloseBlock = len(assessment.Reasons) > 0
	assessment.StartNewBlock = assessment.CloseBlock
	if assessment.CloseBlock {
		for _, evidence := range assessment.Evidence {
			if evidence.Confidence > assessment.Confidence {
				assessment.Confidence = evidence.Confidence
			}
		}
	}
	return assessment, nil
}

// ThreadResolverPolicy supplies deterministic matching thresholds outside the
// Thread model. Candidate Margin prevents accidental merges on close scores.
type ThreadResolverPolicy struct {
	SelectionThreshold float64
	AmbiguityMargin    float64
}

type DeterministicThreadResolver struct{ Policy ThreadResolverPolicy }

func (r DeterministicThreadResolver) Resolve(personID string, triggers []SemanticTrigger, threads []Thread, at EvaluationMoment) (ThreadResolution, error) {
	if normalizeText(personID) == "" || at.Validate() != nil || !isZeroOrUnit(r.Policy.SelectionThreshold) || !isZeroOrUnit(r.Policy.AmbiguityMargin) {
		return ThreadResolution{}, ErrInvalidThread
	}
	resolution := ThreadResolution{PersonID: personID, ResolvedAt: at.WallClockAt}
	for _, thread := range threads {
		if thread.PersonID != personID || thread.Validate() != nil {
			continue
		}
		candidate := scoreThreadCandidate(thread, triggers)
		if candidate.Score > 0 {
			resolution.Candidates = append(resolution.Candidates, candidate)
		}
	}
	sort.Slice(resolution.Candidates, func(i, j int) bool {
		if resolution.Candidates[i].Score == resolution.Candidates[j].Score {
			return resolution.Candidates[i].ThreadID < resolution.Candidates[j].ThreadID
		}
		return resolution.Candidates[i].Score > resolution.Candidates[j].Score
	})
	if len(resolution.Candidates) == 0 {
		return resolution, nil
	}
	top := resolution.Candidates[0]
	if top.Score < r.Policy.SelectionThreshold {
		top.Unresolved = true
		resolution.Candidates[0] = top
		resolution.RequiresReview = true
		return resolution, nil
	}
	if len(resolution.Candidates) > 1 && top.Score-resolution.Candidates[1].Score < r.Policy.AmbiguityMargin {
		top.Unresolved = true
		resolution.Candidates[0] = top
		resolution.RequiresReview = true
		return resolution, nil
	}
	resolution.SelectedID = top.ThreadID
	return resolution, nil
}

func scoreThreadCandidate(thread Thread, triggers []SemanticTrigger) ThreadCandidate {
	candidate := ThreadCandidate{ThreadID: thread.ID}
	seenKinds := map[SemanticTriggerKind]bool{}
	for _, trigger := range triggers {
		if !isZeroOrUnit(trigger.Confidence) {
			continue
		}
		matched := false
		weight := 0.0
		for _, anchor := range thread.Anchors {
			if containsString(trigger.EntityIDs, anchor.ID) || (trigger.Normalized != "" && normalizedEqual(trigger.Normalized, anchor.Label)) || (trigger.Value != "" && normalizedEqual(trigger.Value, anchor.Label)) {
				matched, weight = true, 0.70
				break
			}
		}
		if !matched {
			for _, alias := range thread.Aliases {
				if normalizedEqual(trigger.Normalized, alias) || normalizedEqual(trigger.Value, alias) {
					matched, weight = true, 0.55
					break
				}
			}
		}
		if matched {
			candidate.Score += weight * trigger.Confidence
			if !seenKinds[trigger.Kind] {
				candidate.MatchedBy = append(candidate.MatchedBy, trigger.Kind)
				seenKinds[trigger.Kind] = true
			}
		}
	}
	if candidate.Score > 1 {
		candidate.Score = 1
	}
	candidate.Confidence = candidate.Score
	candidate.Reason = "matched canonical thread anchors or aliases"
	return candidate
}

func normalizedEqual(left, right string) bool {
	return normalizeText(left) != "" && normalizeText(left) == normalizeText(right)
}

func containsString(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

// BuildCurrentThreadState applies only current deltas to a baseline. Superseded
// deltas remain explicitly historical and retain independent provenance.
func BuildCurrentThreadState(thread Thread, baseline CurrentThreadState, deltas []ThreadDelta, at EvaluationMoment) (CurrentThreadState, error) {
	if thread.Validate() != nil || baseline.PersonID != thread.PersonID || baseline.ThreadID != thread.ID || at.Validate() != nil {
		return CurrentThreadState{}, ErrInvalidThread
	}
	state := baseline
	if state.FieldValues == nil {
		state.FieldValues = map[string]string{}
	}
	superseded := map[string]bool{}
	for _, delta := range deltas {
		if delta.PersonID == thread.PersonID && delta.TargetThreadID == thread.ID && delta.SupersedesDeltaID != "" {
			superseded[delta.SupersedesDeltaID] = true
		}
	}
	applicable := make([]ThreadDelta, 0, len(deltas))
	for _, delta := range deltas {
		if delta.PersonID != thread.PersonID || delta.TargetThreadID != thread.ID || delta.Validate() != nil {
			continue
		}
		applicable = append(applicable, delta)
	}
	sort.Slice(applicable, func(i, j int) bool {
		if applicable[i].EventAt.Equal(applicable[j].EventAt) {
			return applicable[i].ID < applicable[j].ID
		}
		return applicable[i].EventAt.Before(applicable[j].EventAt)
	})
	state.IncludedDeltaIDs = nil
	state.HistoricalDeltaIDs = nil
	for _, delta := range applicable {
		if superseded[delta.ID] {
			state.HistoricalDeltaIDs = append(state.HistoricalDeltaIDs, delta.ID)
			continue
		}
		for field, value := range delta.FieldChanges {
			state.FieldValues[field] = value
		}
		state.CurrentSummary = delta.SemanticChange
		state.IncludedDeltaIDs = append(state.IncludedDeltaIDs, delta.ID)
		state.EvidenceIDs = appendContinuityUnique(state.EvidenceIDs, delta.EvidenceIDs...)
	}
	state.ReconstructedAt = at.WallClockAt
	return state, state.Validate()
}

func appendContinuityUnique(existing []string, additions ...string) []string {
	seen := make(map[string]bool, len(existing))
	for _, item := range existing {
		seen[item] = true
	}
	for _, item := range additions {
		if item != "" && !seen[item] {
			existing = append(existing, item)
			seen[item] = true
		}
	}
	return existing
}

// DeterministicRetrievalDepthPolicy bounds retrieval by selected canonical
// records instead of reconstructing a lifetime transcript.
type DeterministicRetrievalDepthPolicy struct {
	KeyContinuityThreshold  float64
	ReconstructionThreshold float64
	DeepAuditThreshold      float64
}

func (p DeterministicRetrievalDepthPolicy) Plan(request RetrievalRequest, thread Thread, state CurrentThreadState, deltas []ThreadDelta, blocks []InteractionBlock) (RetrievalPlan, error) {
	if normalizeText(request.PersonID) == "" || request.PersonID != thread.PersonID || request.ThreadID != thread.ID || request.RequestedAt.IsZero() || request.Priority.Validate() != nil || thread.Validate() != nil || state.Validate() != nil || !isZeroOrUnit(p.KeyContinuityThreshold) || !isZeroOrUnit(p.ReconstructionThreshold) || !isZeroOrUnit(p.DeepAuditThreshold) {
		return RetrievalPlan{}, ErrInvalidRetrievalRequest
	}
	importance := maxFloat(request.Priority.SubjectiveImportance, request.Priority.ObjectiveStakes, request.Priority.ExpectedImpact)
	plan := RetrievalPlan{Depth: RetrievalCurrentState, ThreadStateIDs: []string{state.ID}, Reason: "current materialized state only"}
	if request.EvidenceSensitive || request.Priority.Uncertainty >= p.DeepAuditThreshold || importance >= p.DeepAuditThreshold {
		plan.Depth = RetrievalDeepAudit
		plan.Reason = "high evidence sensitivity, uncertainty, or importance justifies audit depth"
	} else if importance >= p.ReconstructionThreshold || request.Priority.OpportunityCost >= p.ReconstructionThreshold {
		plan.Depth = RetrievalReconstructed
		plan.Reason = "importance or opportunity cost justifies chronological reconstruction"
	} else if importance >= p.KeyContinuityThreshold || len(thread.OpenLoopIDs) > 0 {
		plan.Depth = RetrievalKeyContinuity
		plan.Reason = "material continuity or unresolved work justifies key deltas"
	}
	if plan.Depth == RetrievalAwareness {
		return plan, nil
	}
	if plan.Depth == RetrievalKeyContinuity || plan.Depth == RetrievalReconstructed || plan.Depth == RetrievalDeepAudit {
		for _, delta := range deltas {
			if delta.PersonID == request.PersonID && delta.TargetThreadID == thread.ID {
				plan.DeltaIDs = append(plan.DeltaIDs, delta.ID)
			}
		}
		plan.OpenLoopIDs = append(plan.OpenLoopIDs, thread.OpenLoopIDs...)
	}
	if plan.Depth == RetrievalReconstructed || plan.Depth == RetrievalDeepAudit {
		for _, block := range blocks {
			if block.PersonID == request.PersonID && containsString(block.ThreadIDs, thread.ID) {
				plan.BlockIDs = append(plan.BlockIDs, block.ID)
			}
		}
	}
	if plan.Depth == RetrievalDeepAudit {
		plan.EvidenceIDs = append(plan.EvidenceIDs, state.EvidenceIDs...)
	}
	plan.SelectedObjectCount = len(plan.ThreadStateIDs) + len(plan.DeltaIDs) + len(plan.OpenLoopIDs) + len(plan.BlockIDs) + len(plan.EvidenceIDs)
	return plan, nil
}

func maxFloat(values ...float64) float64 {
	var maximum float64
	for _, value := range values {
		if value > maximum {
			maximum = value
		}
	}
	return maximum
}

// DecideContinuitySurface separates contextual knowledge from interruption.
type ContinuitySurfacePolicy struct {
	MustSurfaceStakes  float64
	MustSurfaceUrgency float64
}

func DecideContinuitySurface(link ContinuityLink, priority PriorityFactors, deadline DeadlineFeasibility, policy ContinuitySurfacePolicy, at EvaluationMoment) (ContinuitySurfaceDecision, error) {
	if link.Validate() != nil || priority.Validate() != nil || at.Validate() != nil || !isZeroOrUnit(policy.MustSurfaceStakes) || !isZeroOrUnit(policy.MustSurfaceUrgency) {
		return ContinuitySurfaceDecision{}, ErrInvalidContinuityLink
	}
	decision := ContinuitySurfaceDecision{LinkID: link.ID, Mode: SurfaceContextOnly, Reason: "relation retained for contextual use without interruption", EvaluatedAt: at.WallClockAt}
	if link.Target.Kind == ContinuityThread {
		decision.ThreadID = link.Target.ID
	}
	if priority.ObjectiveStakes >= policy.MustSurfaceStakes || deadline.State == DeadlineInfeasible || (deadline.State == DeadlineFeasible && priority.ExpectedImpact >= policy.MustSurfaceUrgency) {
		decision.Mode, decision.Reason = SurfaceMust, "high objective stakes or deadline feasibility requires attention policy review"
	} else if priority.SubjectiveImportance > 0 || link.Confidence >= 0.7 {
		decision.Mode, decision.Reason = SurfaceContextOnly, "related continuity available for response adaptation"
	}
	return decision, nil
}

// BaselineInferencePolicy defines sufficient local evidence before a signal may
// create an uncertain interaction hypothesis.
type BaselineInferencePolicy struct {
	MinimumObservations int
	DeviationMultiplier float64
}

func InferInteractionStates(signals []ObservedInteractionSignal, baseline PersonalInteractionBaseline, policy BaselineInferencePolicy, at EvaluationMoment) ([]InferredInteractionState, error) {
	if baseline.Validate() != nil || at.Validate() != nil || policy.MinimumObservations < 1 || policy.DeviationMultiplier <= 0 || policy.DeviationMultiplier > 10 {
		return nil, ErrInvalidInteractionBaseline
	}
	if baseline.ObservationCount < policy.MinimumObservations {
		return nil, nil
	}
	metrics := map[InteractionSignalKind]BaselineMetric{}
	for _, metric := range baseline.Metrics {
		metrics[metric.Kind] = metric
	}
	states := make([]InferredInteractionState, 0)
	for _, signal := range signals {
		if signal.PersonID != baseline.PersonID || signal.Validate() != nil {
			continue
		}
		metric, ok := metrics[signal.Kind]
		if !ok || metric.ObservationCount < policy.MinimumObservations || metric.Tolerance == 0 {
			continue
		}
		deviation := (signal.Value - metric.Mean) / metric.Tolerance
		if deviation < 0 {
			deviation = -deviation
		}
		if deviation < policy.DeviationMultiplier {
			continue
		}
		hypothesis := hypothesisForSignal(signal.Kind, signal.Value-metric.Mean)
		if hypothesis == "" {
			continue
		}
		confidence := deviation / (policy.DeviationMultiplier * 3)
		if confidence > 0.75 {
			confidence = 0.75
		}
		states = append(states, InferredInteractionState{ID: "inference:" + signal.ID, PersonID: signal.PersonID, BlockID: signal.BlockID, Hypothesis: hypothesis, EvidenceSignalIDs: []string{signal.ID}, AlternativeExplanations: []string{"contextual variation", "measurement uncertainty"}, Confidence: confidence, Uncertainty: 1 - confidence, ContextIDs: signal.ContextIDs, EventAt: signal.ObservedAt, EvaluatedAt: at.WallClockAt, Freshness: FreshnessState{LastValidatedAt: at.WallClockAt, Status: FreshnessFresh}, Privacy: PrivacyInferredState})
	}
	return states, nil
}

func hypothesisForSignal(kind InteractionSignalKind, signedDeviation float64) InteractionHypothesis {
	switch kind {
	case SignalPace, SignalInteractionDensity, SignalTopicSwitching:
		if signedDeviation > 0 {
			return HypothesisPossibleUrgency
		}
	case SignalResponseLatency, SignalPauseDensity:
		if signedDeviation > 0 {
			return HypothesisPossibleLowEnergy
		}
	case SignalDirectness:
		if signedDeviation < 0 {
			return HypothesisPossibleUncertainty
		}
	}
	return ""
}

// DeterministicAttunementSafetyPolicy produces only low-risk presentation
// choices. No objective can request manipulation, dependence, impersonation,
// or engagement maximization.
type DeterministicAttunementSafetyPolicy struct {
	DefaultMaxChoices int
}

func (p DeterministicAttunementSafetyPolicy) Decide(control AttunementControlMode, states []InferredInteractionState, patterns []PersonalAttunementPattern, at EvaluationMoment) (InteractionAdaptationDecision, error) {
	if at.Validate() != nil || p.DefaultMaxChoices < 1 {
		return InteractionAdaptationDecision{}, ErrUnsafeAdaptation
	}
	if control == AttunementDisabled {
		return InteractionAdaptationDecision{}, ErrAttunementDisabled
	}
	decision := InteractionAdaptationDecision{ID: "adaptation:" + at.WallClockAt.UTC().Format("20060102T150405.000000000"), Control: control, Objective: ObjectiveComprehension, ResponsePace: 0.5, ResponseVerbosity: 0.5, SentenceLength: 0.5, Directness: 0.5, Warmth: 0.5, MaxChoices: p.DefaultMaxChoices, ResurfaceOpenLoops: false, AvoidUnrelatedSurface: true, FactualContentChanged: false, Reason: "stable low-risk interaction presentation", Confidence: 0.2, EvaluatedAt: at.WallClockAt, Reversible: true, Privacy: PrivacyAdaptationDecision}
	for _, state := range states {
		if state.Validate() != nil {
			continue
		}
		decision.PersonID, decision.BlockID = state.PersonID, state.BlockID
		if state.Hypothesis == HypothesisPossibleUrgency || state.Hypothesis == HypothesisPossibleOverload || state.Hypothesis == HypothesisPossibleLowEnergy {
			decision.Objective = ObjectiveReduceCognitiveFriction
			decision.ResponseVerbosity, decision.SentenceLength, decision.MaxChoices = 0.3, 0.3, 1
			decision.Confidence = state.Confidence
			decision.Reason = "uncertain self-baseline deviation supports reduced cognitive load"
		}
	}
	for _, pattern := range patterns {
		if pattern.Validate() != nil || pattern.Hypothesis == "" || PatternDecisionWeight(pattern, at) == 0 {
			continue
		}
		if pattern.UnhelpfulCount > pattern.BeneficialCount {
			decision.ResponseVerbosity, decision.SentenceLength, decision.MaxChoices = 0.5, 0.5, p.DefaultMaxChoices
			decision.Reason = "contradictory context-specific outcomes keep the adjustment minimal; correlation is not treated as causation"
			continue
		}
		if pattern.BeneficialCount > pattern.UnhelpfulCount && PatternDecisionWeight(pattern, at) >= 0.4 {
			decision.Confidence = minFloat(0.75, decision.Confidence+PatternDecisionWeight(pattern, at)*0.2)
			decision.Reason = "repeated context-specific outcome correlation supports a conservative reversible adjustment"
		}
	}
	if control == AttunementReduced || control == AttunementTemporaryOverride {
		decision.ResponseVerbosity = 0.5
		decision.MaxChoices = p.DefaultMaxChoices
		decision.ResurfaceOpenLoops = false
		decision.Reason = "user control limits attunement strength"
	}
	if decision.PersonID == "" {
		return InteractionAdaptationDecision{}, ErrUnsafeAdaptation
	}
	return decision, ValidateAttunementDecision(decision)
}

func ValidateAttunementDecision(decision InteractionAdaptationDecision) error {
	if decision.Validate() != nil || decision.FactualContentChanged || !decision.Reversible {
		return ErrUnsafeAdaptation
	}
	switch decision.Objective {
	case ObjectiveCovertPersuasion, ObjectiveEngagementMaximization, ObjectiveDependency:
		return ErrUnsafeAdaptation
	}
	return nil
}

// UpdateAttunementPattern records a correlation-only association. Contradictory
// outcomes reduce confidence; one anecdote remains weak by construction.
func UpdateAttunementPattern(pattern PersonalAttunementPattern, outcome InteractionOutcome, at EvaluationMoment) (PersonalAttunementPattern, error) {
	if pattern.Validate() != nil || outcome.PersonID != pattern.PersonID || outcome.ContextSignature != pattern.ContextSignature || at.Validate() != nil || outcome.OccurredAt.IsZero() || outcome.Privacy != PrivacyOutcomeEvidence {
		return PersonalAttunementPattern{}, ErrInvalidAttunementPattern
	}
	updated := pattern
	updated.ObservationCount++
	switch outcome.Status {
	case InteractionOutcomeBeneficial:
		updated.BeneficialCount++
	case InteractionOutcomeUnhelpful:
		updated.UnhelpfulCount++
	case InteractionOutcomeMixed:
		updated.MixedCount++
	}
	// Laplace-like conservative confidence: success evidence is discounted by
	// contradictory/mixed evidence and cannot become strong from one episode.
	denominator := float64(updated.ObservationCount + 4)
	numerator := float64(updated.BeneficialCount) + 1
	confidence := numerator / denominator
	penalty := float64(updated.UnhelpfulCount)*0.15 + float64(updated.MixedCount)*0.05
	updated.Confidence = confidence - penalty
	if updated.Confidence < 0 {
		updated.Confidence = 0
	}
	updated.LastOutcomeAt = outcome.OccurredAt
	updated.Freshness.LastRevalidatedAt = at.WallClockAt
	updated.Freshness.Status = FreshnessFresh
	return updated, updated.Validate()
}

func PatternDecisionWeight(pattern PersonalAttunementPattern, at EvaluationMoment) float64 {
	if pattern.Validate() != nil || at.Validate() != nil || at.WallClockAt.Before(pattern.LastOutcomeAt) {
		return 0
	}
	if pattern.Freshness.StaleAfter == 0 {
		return pattern.Confidence
	}
	age := at.WallClockAt.Sub(pattern.LastOutcomeAt)
	if age <= pattern.Freshness.StaleAfter {
		return pattern.Confidence
	}
	return pattern.Confidence / (1 + float64(age/pattern.Freshness.StaleAfter))
}

func minFloat(left, right float64) float64 {
	if left < right { return left }
	return right
}

func normalizeContextSignature(values []string) string {
	copyValues := append([]string(nil), values...)
	sort.Strings(copyValues)
	return strings.Join(copyValues, "|")
}
