package intelligence_runtime

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/BraveAI93/concierge-backend/kernel"
)

func v05RuntimeInput(now time.Time, conversationID, messageID, blockID, threadID string, proposed bool, includeAttunement bool) ContinuityInput {
	blockAt := now.Add(-5 * time.Minute)
	block := kernel.InteractionBlock{ID: blockID, PersonID: "person-a", SourceRefs: []kernel.InteractionSourceRef{{SourceType: "synthetic_conversation", ConversationID: conversationID, MessageID: messageID, CapturedAt: now}}, StartTemporal: kernel.TemporalState{EventAt: blockAt, RecordedAt: now, EffectiveAt: blockAt, AttentionAt: now}, EvaluatedAt: now, IngestedAt: now, SourceType: "text_chat", TopicLabels: []string{"AI laptop budget"}, EntityIDs: []string{"entity-laptop"}, ContextIDs: []string{"conversation:" + conversationID}, Importance: kernel.PriorityFactors{SubjectiveImportance: .5, ObjectiveStakes: .7, ExpectedImpact: .7, Reversibility: .8, Uncertainty: .2, OpportunityCost: .2, EffortAttentionCost: .2}, SemanticState: "user discussed budget for computer used in AI work", Provenance: []kernel.Provenance{{SourceType: "synthetic", SourceRef: messageID, CapturedAt: now}}, Confidence: .8, Freshness: kernel.FreshnessState{LastValidatedAt: now, Status: kernel.FreshnessFresh}}
	input := ContinuityInput{IdempotencyKey: "continuity:" + messageID, Source: schedulingSource(now, "profile-internal-a", conversationID, messageID), Block: block, Triggers: []kernel.SemanticTrigger{{Kind: kernel.TriggerParaphrase, Value: "computer needed for AI work", EntityIDs: []string{"entity-laptop"}, Confidence: 1}}, AttunementControl: kernel.AttunementNormal, ContextSignature: "work"}
	if proposed {
		input.ProposedThread = &kernel.Thread{ID: threadID, PersonID: "person-a", Anchors: []kernel.ThreadAnchor{{Kind: "entity", ID: "entity-laptop", Label: "AI laptop"}}, Aliases: []string{"computer for AI work"}, Status: kernel.ThreadActive, Importance: block.Importance, FirstRelevantAt: blockAt, MostRecentRelevantAt: now, Confidence: .8, Freshness: kernel.FreshnessState{LastValidatedAt: now, Status: kernel.FreshnessFresh}, CreatedAt: now}
		input.Deltas = []kernel.ThreadDelta{{ID: "delta:" + messageID, PersonID: "person-a", TargetThreadID: threadID, Originating: kernel.ContinuityRef{Kind: kernel.ContinuityBlock, ID: blockID}, SemanticChange: "budget initially set to 500", AffectedConcept: "budget", FieldChanges: map[string]string{"budget": "500"}, Effects: []kernel.ThreadDeltaEffect{kernel.DeltaUpdatesConstraint}, Confidence: .8, Importance: block.Importance, EventAt: blockAt, EvaluatedAt: now, CreatedAt: now}}
	}
	if includeAttunement {
		input.Baseline = &kernel.PersonalInteractionBaseline{ID: "baseline:person-a:work", PersonID: "person-a", ContextSignature: "work", Metrics: []kernel.BaselineMetric{{Kind: kernel.SignalResponseLatency, Mean: 1, Tolerance: 1, ObservationCount: 6}}, ObservationCount: 6, Confidence: .7, LastValidatedAt: now, DecayAfter: 24 * time.Hour, Freshness: kernel.FreshnessState{LastValidatedAt: now, Status: kernel.FreshnessFresh}, Privacy: kernel.PrivacyDerivedBaseline}
		input.Signals = []kernel.ObservedInteractionSignal{{ID: "signal:" + messageID, PersonID: "person-a", BlockID: blockID, Kind: kernel.SignalResponseLatency, Value: 4, Unit: "seconds", ContextIDs: []string{"work"}, ObservedAt: now, Provenance: kernel.Provenance{SourceType: "synthetic", SourceRef: messageID, CapturedAt: now}, Privacy: kernel.PrivacyRawCommunicationSignal, Confidence: 1}}
	}
	return input
}

func TestV05RuntimePersistentContinuityAndAttunementVerticalSlice(t *testing.T) {
	now := time.Date(2026, time.August, 22, 12, 0, 0, 0, time.UTC)
	service, repo, principal := runtimeFixture(t, now)
	firstInput := v05RuntimeInput(now, "continuity-chat-a", "continuity-message-a", "block-a", "thread-laptop", true, true)
	first, err := service.IngestContinuity(context.Background(), principal, firstInput)
	if err != nil || first.Replayed || first.ThreadID != "thread-laptop" || first.BlockID == "" || first.ThreadStateID == "" || first.AttunementEpisodeID == "" || first.AdaptationID == "" || first.InterventionID == "" {
		t.Fatalf("v0.5 first ingestion must persist full local lineage: %+v err=%v", first, err)
	}
	secondInput := v05RuntimeInput(now.Add(40*time.Minute), "continuity-chat-b", "continuity-message-b", "block-b", "thread-laptop", false, false)
	// The second chat uses different physical source identifiers but a semantic
	// entity trigger. It must continue the existing source-independent Thread.
	service.Clock = fixedClock{now: now.Add(40 * time.Minute)}
	second, err := service.IngestContinuity(context.Background(), principal, secondInput)
	if err != nil || second.ThreadID != first.ThreadID || second.Replayed {
		t.Fatalf("same semantic subject across chats must continue thread: first=%+v second=%+v err=%v", first, second, err)
	}
	replay, err := service.IngestContinuity(context.Background(), principal, secondInput)
	if err != nil || !replay.Replayed || replay.BlockID != second.BlockID {
		t.Fatalf("continuity replay must be deterministic: replay=%+v err=%v", replay, err)
	}
	state, err := service.State(context.Background(), principal)
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Sources) != 2 || len(state.InteractionBlocks) != 2 || len(state.Threads) != 1 || len(state.ContinuityLinks) != 2 || len(state.ThreadDeltas) != 1 || len(state.CurrentThreadStates) != 2 || len(state.ObservedSignals) != 1 || len(state.InferredInteractionStates) != 1 || len(state.AdaptationDecisions) != 1 || len(state.AttunementEpisodes) != 1 || len(state.InteractionInterventions) != 1 {
		t.Fatalf("unexpected persisted v0.5 state: %+v", state)
	}
	if len(state.InteractionBlocks[1].BoundaryEvidence) == 0 {
		t.Fatalf("continued block must preserve evaluated segmentation evidence")
	}

	resolution, plan, err := service.ResolveContinuity(context.Background(), principal, kernel.RetrievalRequest{PersonID: "person-a", Priority: kernel.PriorityFactors{SubjectiveImportance: .2, ObjectiveStakes: .2, ExpectedImpact: .2, Reversibility: .8, Uncertainty: .2, OpportunityCost: .2, EffortAttentionCost: .2}, Triggers: firstInput.Triggers, RequestedAt: now.Add(40 * time.Minute)})
	if err != nil || resolution.SelectedID != first.ThreadID || plan.Depth != kernel.RetrievalCurrentState || len(plan.BlockIDs) != 0 {
		t.Fatalf("low-stakes retrieval must select bounded continuity without full history: resolution=%+v plan=%+v err=%v", resolution, plan, err)
	}

	outcome := kernel.InteractionOutcome{ID: "attunement-outcome-a", PersonID: "person-a", EpisodeID: first.AttunementEpisodeID, InterventionID: first.InterventionID, Status: kernel.InteractionOutcomeBeneficial, Summary: "synthetic user-completion evidence", EvidenceIDs: []string{"synthetic-evidence"}, ExplicitFeedback: "helpful", ContextSignature: "work", TimeToOutcome: 2 * time.Minute, OccurredAt: now.Add(45 * time.Minute), RecordedAt: now.Add(45 * time.Minute), Privacy: kernel.PrivacyOutcomeEvidence}
	seed := kernel.PersonalAttunementPattern{ID: "pattern:work:overload", PersonID: "person-a", ContextSignature: "work", Hypothesis: kernel.HypothesisPossibleLowEnergy, StrategyFingerprint: "reduced-cognitive-load", ObservationCount: 0, Confidence: 0, Freshness: kernel.FreshnessState{LastValidatedAt: now, Status: kernel.FreshnessFresh}, LastOutcomeAt: now, UserOverridable: true, CorrelationOnly: true, Privacy: kernel.PrivacyLearnedPattern}
	pattern, err := service.RecordAttunementOutcome(context.Background(), principal, outcome, seed)
	if err != nil || pattern.ObservationCount != 1 || !pattern.CorrelationOnly || pattern.Confidence >= .5 {
		t.Fatalf("closed loop must produce only a weak reversible correlation: %+v err=%v", pattern, err)
	}
	state, err = service.State(context.Background(), principal)
	if err != nil || len(state.InteractionOutcomes) != 1 || len(state.AttunementPatterns) != 1 {
		t.Fatalf("outcome learning persistence missing: state=%+v err=%v", state, err)
	}

	if err := repo.RunInTransaction(context.Background(), "person-a", func(tx RuntimeTransaction) error {
		return tx.SaveContinuityLink(kernel.ContinuityLink{ID: "foreign-link", PersonID: "person-b", Source: kernel.ContinuityRef{Kind: kernel.ContinuityBlock, ID: "block-a"}, Target: kernel.ContinuityRef{Kind: kernel.ContinuityThread, ID: "thread-laptop"}, Relation: kernel.ContinuitySameSubject, Why: "cross-person", Confidence: .5, Temporal: firstInput.Block.StartTemporal, Freshness: kernel.FreshnessState{LastValidatedAt: now, Status: kernel.FreshnessFresh}, CreatedAt: now})
	}); !errors.Is(err, ErrCrossPersonAccess) {
		t.Fatalf("cross-person continuity write must reject: %v", err)
	}
}

func TestV05RuntimeDisabledOrUserAttunementDisabledProducesNoAdaptation(t *testing.T) {
	now := time.Date(2026, time.August, 22, 13, 0, 0, 0, time.UTC)
	service, repo, principal := runtimeFixture(t, now)
	service.Feature = DisabledFeature()
	if _, err := service.IngestContinuity(context.Background(), principal, v05RuntimeInput(now, "disabled", "disabled-message", "disabled-block", "disabled-thread", true, true)); !errors.Is(err, ErrRuntimeDisabled) {
		t.Fatalf("default-off feature must reject v0.5 path: %v", err)
	}
	service.Feature = EnabledFeature()
	input := v05RuntimeInput(now, "attunement-off", "attunement-off-message", "attunement-off-block", "attunement-off-thread", true, true)
	input.AttunementControl = kernel.AttunementDisabled
	result, err := service.IngestContinuity(context.Background(), principal, input)
	if err != nil || result.AdaptationID != "" || result.AttunementEpisodeID != "" {
		t.Fatalf("user control disabled must persist no adaptation: %+v err=%v", result, err)
	}
	state, err := repo.ReadState(context.Background(), "person-a", "person-a")
	if err != nil || len(state.AdaptationDecisions) != 0 || len(state.AttunementEpisodes) != 0 {
		t.Fatalf("disabled attunement must not write adaptive records: %+v err=%v", state, err)
	}
}

func TestV05PostgresContinuityRestartRoundTripAndRLS(t *testing.T) {
	now := time.Date(2026, time.August, 22, 14, 0, 0, 0, time.UTC)
	service, repo, principal := postgresRuntimeFixture(t, now)
	input := v05RuntimeInput(now, "postgres-continuity", "postgres-continuity-message", "postgres-block", "postgres-thread", true, true)
	result, err := service.IngestContinuity(context.Background(), principal, input)
	if err != nil {
		t.Fatalf("persist v0.5 block/thread to local postgres: %v", err)
	}
	if err := repo.Close(); err != nil {
		t.Fatal(err)
	}
	secondRepo := openPostgresRepository(t)
	defer secondRepo.Close()
	service.Repo = secondRepo
	service.Identity = PostgresIdentityResolver{Repository: secondRepo}
	state, err := service.State(context.Background(), principal)
	if err != nil || len(state.InteractionBlocks) != 1 || len(state.Threads) != 1 || len(state.ContinuityLinks) != 1 || len(state.CurrentThreadStates) != 1 || len(state.AttunementEpisodes) != 1 || len(state.InteractionInterventions) != 1 {
		t.Fatalf("restart must reconstruct v0.5 records: %+v err=%v", state, err)
	}
	if state.InteractionBlocks[0].ID != result.BlockID || state.Threads[0].ID != result.ThreadID {
		t.Fatalf("restart IDs must preserve lineage: result=%+v state=%+v", result, state)
	}
	if _, err := secondRepo.ReadState(context.Background(), "person-b", "person-a"); !errors.Is(err, ErrCrossPersonAccess) {
		t.Fatalf("cross-person state read must still fail: %v", err)
	}
}

func TestV05PostgresContinuityGraphRLSAndLeastPrivilege(t *testing.T) {
	now := time.Date(2026, time.August, 22, 15, 0, 0, 0, time.UTC)
	_, repo, _ := postgresRuntimeFixture(t, now)
	defer repo.Close()
	for _, table := range []string{"continuity_links", "thread_deltas"} {
		var forced bool
		if err := repo.db.QueryRow("SELECT relforcerowsecurity FROM pg_class c JOIN pg_namespace n ON n.oid=c.relnamespace WHERE n.nspname='ci_kernel_v04' AND c.relname=$1", table).Scan(&forced); err != nil || !forced {
			t.Fatalf("v0.5 table %s must retain forced RLS, forced=%t err=%v", table, forced, err)
		}
		var canDelete bool
		if err := repo.db.QueryRow("SELECT has_table_privilege('ci_kernel_runtime', 'ci_kernel_v04.' || $1, 'DELETE')", table).Scan(&canDelete); err != nil || canDelete {
			t.Fatalf("runtime role must not delete %s, canDelete=%t err=%v", table, canDelete, err)
		}
	}
	var publicExecute bool
	if err := repo.db.QueryRow("SELECT has_function_privilege('public', 'ci_kernel_v04.current_person_id()', 'EXECUTE')").Scan(&publicExecute); err != nil || publicExecute {
		t.Fatalf("PUBLIC must not execute staging identity function, publicExecute=%t err=%v", publicExecute, err)
	}
	bindingB := testBinding("owner-subject-b", "profile-internal-b", "person-b", now)
	if err := (PostgresIdentityProvisioner{Repository: repo}).Provision(context.Background(), bindingB); err != nil { t.Fatal(err) }
	if err := repo.RunInTransaction(context.Background(), "person-a", func(tx RuntimeTransaction) error {
		return tx.SaveContinuityLink(kernel.ContinuityLink{ID: "cross-person-v05-link", PersonID: "person-b", Source: kernel.ContinuityRef{Kind: kernel.ContinuityBlock, ID: "other"}, Target: kernel.ContinuityRef{Kind: kernel.ContinuityThread, ID: "other"}, Relation: kernel.ContinuitySameSubject, Why: "attempt", Confidence: .5, Temporal: testTemporal(now), Freshness: kernel.FreshnessState{LastValidatedAt: now, Status: kernel.FreshnessFresh}, CreatedAt: now})
	}); !errors.Is(err, ErrCrossPersonAccess) { t.Fatalf("runtime transaction must reject cross-person v0.5 edge: %v", err) }
}

func TestV05PostgresOutcomeRequiresSameEpisodeIntervention(t *testing.T) {
	now := time.Date(2026, time.August, 22, 16, 0, 0, 0, time.UTC)
	service, repo, principal := postgresRuntimeFixture(t, now)
	result, err := service.IngestContinuity(context.Background(), principal, v05RuntimeInput(now, "outcome-parent", "outcome-parent-message", "outcome-parent-block", "outcome-parent-thread", true, true))
	if err != nil { t.Fatal(err) }
	wrong := kernel.InteractionOutcome{ID: "wrong-intervention", PersonID: "person-a", EpisodeID: result.AttunementEpisodeID, InterventionID: "not-an-intervention", Status: kernel.InteractionOutcomeUnknown, ContextSignature: "work", OccurredAt: now, RecordedAt: now, Privacy: kernel.PrivacyOutcomeEvidence}
	seed := kernel.PersonalAttunementPattern{ID: "wrong-pattern", PersonID: "person-a", ContextSignature: "work", Hypothesis: kernel.HypothesisPossibleLowEnergy, StrategyFingerprint: "test", Freshness: kernel.FreshnessState{LastValidatedAt: now, Status: kernel.FreshnessFresh}, LastOutcomeAt: now, UserOverridable: true, CorrelationOnly: true, Privacy: kernel.PrivacyLearnedPattern}
	if _, err := service.RecordAttunementOutcome(context.Background(), principal, wrong, seed); !errors.Is(err, ErrCrossPersonAccess) { t.Fatalf("outcome must require same-person same-episode intervention: %v", err) }
	state, err := repo.ReadState(context.Background(), "person-a", "person-a")
	if err != nil || len(state.InteractionOutcomes) != 0 { t.Fatalf("failed outcome transaction must roll back: %+v err=%v", state, err) }
}

func TestV05MigrationDeclaresTypedContinuityEdges(t *testing.T) {
	migration, err := os.ReadFile("../db/proposed_migrations/ci_kernel_v04/001_staging_up.sql")
	if err != nil { t.Fatal(err) }
	text := strings.ToLower(string(migration))
	for _, required := range []string{"create table ci_kernel_v04.continuity_links", "create table ci_kernel_v04.thread_deltas", "foreign key (source_kind, person_id, source_id)", "foreign key (target_kind, person_id, target_id)", "force row level security", "alter default privileges in schema ci_kernel_v04 revoke execute on functions from public"} {
		if !strings.Contains(text, required) { t.Fatalf("v0.5 migration missing required hardening fragment %q", required) }
	}
}
