package kernel

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func v05Now() time.Time { return time.Date(2026, time.August, 22, 10, 0, 0, 0, time.UTC) }
func v05Temporal(at time.Time) TemporalState {
	return TemporalState{EventAt: at, RecordedAt: at, EffectiveAt: at, AttentionAt: at}
}
func v05Fresh(at time.Time) FreshnessState {
	return FreshnessState{LastValidatedAt: at, Status: FreshnessFresh}
}
func v05Priority(stakes float64) PriorityFactors {
	return PriorityFactors{SubjectiveImportance: 0.4, ObjectiveStakes: stakes, ExpectedImpact: stakes, Reversibility: 0.8, Uncertainty: 0.2, OpportunityCost: 0.2, EffortAttentionCost: 0.2}
}
func v05Thread(person, id, entity string, at time.Time) Thread {
	return Thread{ID: id, PersonID: person, Anchors: []ThreadAnchor{{Kind: "entity", ID: entity, Label: "AI laptop"}}, Aliases: []string{"computer for AI work", "laptop"}, Status: ThreadActive, Importance: v05Priority(0.5), FirstRelevantAt: at.Add(-time.Hour), MostRecentRelevantAt: at, Confidence: 0.8, Freshness: v05Fresh(at), CreatedAt: at}
}
func v05Block(person, id string, at time.Time) InteractionBlock {
	return InteractionBlock{ID: id, PersonID: person, StartTemporal: v05Temporal(at), EvaluatedAt: at, IngestedAt: at, SourceType: "text_chat", Importance: v05Priority(0.5), SemanticState: "synthetic interaction", Confidence: 0.8, Freshness: v05Fresh(at), Provenance: []Provenance{{SourceType: "synthetic", SourceRef: id, CapturedAt: at}}}
}

func TestV05SemanticContinuityAcrossChatsAndParaphrase(t *testing.T) {
	now := v05Now()
	thread := v05Thread("person-a", "thread-laptop", "entity-laptop", now)
	resolver := DeterministicThreadResolver{Policy: ThreadResolverPolicy{SelectionThreshold: 0.50, AmbiguityMargin: 0.15}}
	// Different chat/session sources are provenance only; the entity anchor and
	// paraphrase cue resolve the same source-independent thread.
	resolution, err := resolver.Resolve("person-a", []SemanticTrigger{{Kind: TriggerParaphrase, Value: "that computer I needed for the AI work", EntityIDs: []string{"entity-laptop"}, Confidence: 1}}, []Thread{thread}, EvaluationMoment{WallClockAt: now})
	if err != nil || resolution.SelectedID != thread.ID || resolution.RequiresReview {
		t.Fatalf("paraphrase must resolve canonical thread: %+v err=%v", resolution, err)
	}
	first := v05Block("person-a", "block-chat-a", now.Add(-time.Hour))
	second := v05Block("person-a", "block-chat-b", now)
	if first.SourceType == second.SourceType && first.ID == second.ID {
		t.Fatal("fixture must represent distinct physical source containers")
	}
	link := ContinuityLink{ID: "link-b", PersonID: "person-a", Source: ContinuityRef{Kind: ContinuityBlock, ID: second.ID}, Target: ContinuityRef{Kind: ContinuityThread, ID: thread.ID}, Relation: ContinuitySameSubject, Why: "canonical entity anchor", Confidence: resolution.Candidates[0].Confidence, Temporal: v05Temporal(now), Freshness: v05Fresh(now), CreatedAt: now}
	if err := link.Validate(); err != nil {
		t.Fatalf("typed continuity link: %v", err)
	}
}

func TestV05AmbiguousThreadDoesNotSilentlyMerge(t *testing.T) {
	now := v05Now()
	a := v05Thread("person-a", "thread-a", "entity-shared", now)
	b := v05Thread("person-a", "thread-b", "entity-shared", now)
	resolver := DeterministicThreadResolver{Policy: ThreadResolverPolicy{SelectionThreshold: .5, AmbiguityMargin: .20}}
	resolution, err := resolver.Resolve("person-a", []SemanticTrigger{{Kind: TriggerEntity, EntityIDs: []string{"entity-shared"}, Confidence: 1}}, []Thread{a, b}, EvaluationMoment{WallClockAt: now})
	if err != nil || resolution.SelectedID != "" || !resolution.RequiresReview || len(resolution.Candidates) != 2 {
		t.Fatalf("ambiguous candidate must remain unresolved: %+v err=%v", resolution, err)
	}
}

func TestV05ThreadDeltaUpdatesCurrentStatePreservesHistory(t *testing.T) {
	now := v05Now()
	thread := v05Thread("person-a", "thread-laptop", "entity-laptop", now)
	baseline := CurrentThreadState{ID: "state-baseline", PersonID: "person-a", ThreadID: thread.ID, BaselineSummary: "budget 500", CurrentSummary: "budget 500", FieldValues: map[string]string{"budget": "500"}, ReconstructedAt: now.Add(-time.Hour), Freshness: v05Fresh(now)}
	old := ThreadDelta{ID: "delta-old", PersonID: "person-a", TargetThreadID: thread.ID, Originating: ContinuityRef{Kind: ContinuityBlock, ID: "block-laptop"}, SemanticChange: "budget is 500", AffectedConcept: "budget", FieldChanges: map[string]string{"budget": "500"}, Effects: []ThreadDeltaEffect{DeltaUpdatesState}, Confidence: .8, Importance: v05Priority(.6), EventAt: now.Add(-30 * time.Minute), EvaluatedAt: now.Add(-29 * time.Minute), CreatedAt: now.Add(-29 * time.Minute)}
	newer := ThreadDelta{ID: "delta-new", PersonID: "person-a", TargetThreadID: thread.ID, Originating: ContinuityRef{Kind: ContinuityBlock, ID: "block-contract"}, SemanticChange: "budget raised to 1200 by changed AI workload", AffectedConcept: "budget", FieldChanges: map[string]string{"budget": "1200"}, Effects: []ThreadDeltaEffect{DeltaUpdatesConstraint}, Confidence: .9, Importance: v05Priority(.9), EventAt: now, EvaluatedAt: now, SupersedesDeltaID: old.ID, CreatedAt: now}
	state, err := BuildCurrentThreadState(thread, baseline, []ThreadDelta{old, newer}, EvaluationMoment{WallClockAt: now})
	if err != nil || state.FieldValues["budget"] != "1200" || len(state.HistoricalDeltaIDs) != 1 || state.HistoricalDeltaIDs[0] != old.ID {
		t.Fatalf("current state must show new value and retain old provenance: %+v err=%v", state, err)
	}
}

func TestV05ImportanceSensitiveBoundedRetrievalAndNoFullHistoryScan(t *testing.T) {
	now := v05Now()
	thread := v05Thread("person-a", "thread-low", "entity-low", now)
	state := CurrentThreadState{ID: "state-low", PersonID: "person-a", ThreadID: thread.ID, CurrentSummary: "current", ReconstructedAt: now, Freshness: v05Fresh(now)}
	blocks := make([]InteractionBlock, 0, 200)
	for i := 0; i < 200; i++ {
		block := v05Block("person-a", fmt.Sprintf("block-%03d", i), now.Add(-time.Duration(i)*time.Hour))
		block.ThreadIDs = []string{thread.ID}
		blocks = append(blocks, block)
	}
	policy := DeterministicRetrievalDepthPolicy{KeyContinuityThreshold: .6, ReconstructionThreshold: .8, DeepAuditThreshold: .9}
	low, err := policy.Plan(RetrievalRequest{PersonID: "person-a", ThreadID: thread.ID, Priority: v05Priority(.1), RequestedAt: now}, thread, state, nil, blocks)
	if err != nil || low.Depth != RetrievalCurrentState || low.SelectedObjectCount != 1 || len(low.BlockIDs) != 0 {
		t.Fatalf("low stakes must avoid full history: %+v err=%v", low, err)
	}
	high, err := policy.Plan(RetrievalRequest{PersonID: "person-a", ThreadID: thread.ID, Priority: v05Priority(.95), EvidenceSensitive: true, RequestedAt: now}, thread, state, nil, blocks)
	if err != nil || high.Depth != RetrievalDeepAudit || len(high.BlockIDs) != 200 {
		t.Fatalf("high stakes audit can select full relevant blocks: %+v err=%v", high, err)
	}
}

func TestV05RetrievalIsNotInterruptionAndMustSurfaceIsExplicit(t *testing.T) {
	now := v05Now()
	link := ContinuityLink{ID: "link", PersonID: "person-a", Source: ContinuityRef{Kind: ContinuityBlock, ID: "block"}, Target: ContinuityRef{Kind: ContinuityThread, ID: "thread"}, Relation: ContinuityMateriallyImpacts, Why: "related context", Confidence: .7, Temporal: v05Temporal(now), Freshness: v05Fresh(now), CreatedAt: now}
	quiet, err := DecideContinuitySurface(link, v05Priority(.2), DeadlineFeasibility{State: DeadlineNone}, ContinuitySurfacePolicy{MustSurfaceStakes: .8, MustSurfaceUrgency: .9}, EvaluationMoment{WallClockAt: now})
	if err != nil || quiet.Mode == SurfaceMust {
		t.Fatalf("low-value relation must remain non-interrupting: %+v err=%v", quiet, err)
	}
	critical, err := DecideContinuitySurface(link, v05Priority(.95), DeadlineFeasibility{State: DeadlineInfeasible}, ContinuitySurfacePolicy{MustSurfaceStakes: .8, MustSurfaceUrgency: .9}, EvaluationMoment{WallClockAt: now})
	if err != nil || critical.Mode != SurfaceMust {
		t.Fatalf("critical deadline/stakes must become must-surface: %+v err=%v", critical, err)
	}
}

func v05Baseline(person string, tolerance float64, at time.Time) PersonalInteractionBaseline {
	return PersonalInteractionBaseline{ID: "baseline:" + person, PersonID: person, ContextSignature: "work", Metrics: []BaselineMetric{{Kind: SignalResponseLatency, Mean: 1, Tolerance: tolerance, ObservationCount: 6}}, ObservationCount: 6, Confidence: .7, LastValidatedAt: at, DecayAfter: 24 * time.Hour, Freshness: v05Fresh(at), Privacy: PrivacyDerivedBaseline}
}
func v05Signal(person string, value float64, at time.Time) ObservedInteractionSignal {
	return ObservedInteractionSignal{ID: "signal:" + person, PersonID: person, BlockID: "block:" + person, Kind: SignalResponseLatency, Value: value, Unit: "seconds", ContextIDs: []string{"work"}, ObservedAt: at, Provenance: Provenance{SourceType: "synthetic", SourceRef: person, CapturedAt: at}, Privacy: PrivacyRawCommunicationSignal, Confidence: 1}
}

func TestV05PersonalBaselineAndNoEmotionalCertainty(t *testing.T) {
	now := v05Now()
	statesA, err := InferInteractionStates([]ObservedInteractionSignal{v05Signal("person-a", 4, now)}, v05Baseline("person-a", 1, now), BaselineInferencePolicy{MinimumObservations: 3, DeviationMultiplier: 1.5}, EvaluationMoment{WallClockAt: now})
	if err != nil || len(statesA) != 1 || statesA[0].Hypothesis != HypothesisPossibleLowEnergy || statesA[0].Confidence >= 1 || len(statesA[0].AlternativeExplanations) == 0 {
		t.Fatalf("same signal should create only uncertain personal hypothesis: %+v err=%v", statesA, err)
	}
	statesB, err := InferInteractionStates([]ObservedInteractionSignal{v05Signal("person-b", 4, now)}, v05Baseline("person-b", 10, now), BaselineInferencePolicy{MinimumObservations: 3, DeviationMultiplier: 1.5}, EvaluationMoment{WallClockAt: now})
	if err != nil || len(statesB) != 0 {
		t.Fatalf("same signal must not be strong for a different personal baseline: %+v err=%v", statesB, err)
	}
}

func TestV05AttunementSafetyControlsAndOutcomeLearning(t *testing.T) {
	now := v05Now()
	state := InferredInteractionState{ID: "inference", PersonID: "person-a", BlockID: "block-a", Hypothesis: HypothesisPossibleOverload, EvidenceSignalIDs: []string{"signal"}, AlternativeExplanations: []string{"context"}, Confidence: .5, Uncertainty: .5, EventAt: now, EvaluatedAt: now, Freshness: v05Fresh(now), Privacy: PrivacyInferredState}
	policy := DeterministicAttunementSafetyPolicy{DefaultMaxChoices: 3}
	decision, err := policy.Decide(AttunementNormal, []InferredInteractionState{state}, nil, EvaluationMoment{WallClockAt: now})
	if err != nil || decision.MaxChoices != 1 || decision.FactualContentChanged || decision.Objective != ObjectiveReduceCognitiveFriction {
		t.Fatalf("safe adaptation must reduce cognitive load without changing facts: %+v err=%v", decision, err)
	}
	if _, err := policy.Decide(AttunementDisabled, []InferredInteractionState{state}, nil, EvaluationMoment{WallClockAt: now}); !errors.Is(err, ErrAttunementDisabled) {
		t.Fatalf("user disabled control must take effect: %v", err)
	}
	unsafe := decision
	unsafe.Objective = ObjectiveEngagementMaximization
	if err := ValidateAttunementDecision(unsafe); !errors.Is(err, ErrUnsafeAdaptation) {
		t.Fatalf("engagement maximization must be rejected: %v", err)
	}

	pattern := PersonalAttunementPattern{ID: "pattern", PersonID: "person-a", ContextSignature: "work", Hypothesis: state.Hypothesis, StrategyFingerprint: "short-direct-one-choice", ObservationCount: 0, Confidence: 0, Freshness: v05Fresh(now), LastOutcomeAt: now, UserOverridable: true, CorrelationOnly: true, Privacy: PrivacyLearnedPattern}
	beneficial := InteractionOutcome{ID: "outcome-1", PersonID: "person-a", EpisodeID: "episode", Status: InteractionOutcomeBeneficial, ContextSignature: "work", OccurredAt: now, RecordedAt: now, Privacy: PrivacyOutcomeEvidence}
	pattern, err = UpdateAttunementPattern(pattern, beneficial, EvaluationMoment{WallClockAt: now})
	if err != nil || pattern.ObservationCount != 1 || pattern.Confidence >= .5 {
		t.Fatalf("one outcome must create only weak correlation: %+v err=%v", pattern, err)
	}
	for i := 2; i <= 4; i++ {
		beneficial.ID = fmt.Sprintf("outcome-%d", i)
		beneficial.OccurredAt = now.Add(time.Duration(i) * time.Minute)
		pattern, err = UpdateAttunementPattern(pattern, beneficial, EvaluationMoment{WallClockAt: beneficial.OccurredAt})
		if err != nil {
			t.Fatal(err)
		}
	}
	beforeFailure := pattern.Confidence
	failure := beneficial
	failure.ID = "outcome-failure"
	failure.Status = InteractionOutcomeUnhelpful
	failure.OccurredAt = now.Add(time.Hour)
	pattern, err = UpdateAttunementPattern(pattern, failure, EvaluationMoment{WallClockAt: failure.OccurredAt})
	if err != nil || pattern.Confidence >= beforeFailure {
		t.Fatalf("contradictory evidence must lower confidence: before=%f after=%f err=%v", beforeFailure, pattern.Confidence, err)
	}
	foreign := failure
	foreign.ContextSignature = "health"
	if _, err := UpdateAttunementPattern(pattern, foreign, EvaluationMoment{WallClockAt: now}); !errors.Is(err, ErrInvalidAttunementPattern) {
		t.Fatalf("cross-context result must not generalize automatically: %v", err)
	}
	pattern.Freshness.StaleAfter = time.Hour
	if weight := PatternDecisionWeight(pattern, EvaluationMoment{WallClockAt: now.Add(48 * time.Hour)}); weight >= pattern.Confidence {
		t.Fatalf("stale pattern must decay: weight=%f confidence=%f", weight, pattern.Confidence)
	}
}
