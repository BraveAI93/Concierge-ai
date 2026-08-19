package kernel

import (
	"errors"
	"testing"
	"time"
)

func TestV02FutureEventRecordedBeforeSemanticEventTimeIsValid(t *testing.T) {
	now := fixedV02Now()
	appointment := TemporalState{
		EventAt:     now.Add(24 * time.Hour),
		RecordedAt:  now,
		EffectiveAt: now,
		AttentionAt: now.Add(12 * time.Hour),
	}
	if err := appointment.Validate(); err != nil {
		t.Fatalf("future appointment recorded today must be valid: %v", err)
	}
}

func TestV02EventAgeIsIndependentFromInteractionGap(t *testing.T) {
	now := fixedV02Now()
	temporal := TemporalState{
		EventAt:     now.Add(-30 * 24 * time.Hour),
		RecordedAt:  now.Add(-29 * 24 * time.Hour),
		EffectiveAt: now.Add(-30 * 24 * time.Hour),
		AttentionAt: now,
	}
	gap := InteractionGapState{
		PersonID:          "person-1",
		EntityID:          "entity-client",
		LastInteractionAt: now.Add(-time.Hour),
		ObservedAt:        now,
	}
	evaluation, err := DefaultV02Policy().EvaluateTemporal(TemporalPriorityInput{
		Moment:               EvaluationMoment{WallClockAt: now},
		Temporal:             temporal,
		InteractionGap:       gap,
		Effort:               EffortAttention{EstimatedEffort: 15 * time.Minute, EstimatedAttention: 10 * time.Minute},
		SubjectiveImportance: 0.8,
	})
	if err != nil {
		t.Fatalf("evaluate independent temporal dimensions: %v", err)
	}
	if evaluation.InteractionGap != time.Hour {
		t.Fatalf("expected a one-hour interaction gap, got %s", evaluation.InteractionGap)
	}
	if evaluation.Recency >= 0.5 {
		t.Fatalf("expected a thirty-day semantic event to have low event recency, got %f", evaluation.Recency)
	}
}

func TestV02AttentionBudgetConstrainsAndContextuallyResurfacesOpenLoops(t *testing.T) {
	now := fixedV02Now()
	budget := AttentionBudget{
		ID:                "budget-morning",
		PersonID:          "person-1",
		WindowStart:       now,
		WindowEnd:         now.Add(time.Hour),
		AttentionCapacity: 30 * time.Minute,
		MaxCompetingItems: 1,
		CurrentContext:    AttentionContext{ContextIDs: []string{"context-client-work"}, EntityIDs: []string{"entity-sam"}},
	}
	candidates := []AttentionCandidate{
		{
			OpenLoopID:   "loop-context-match",
			PersonID:     "person-1",
			Label:        "Respond to Sam",
			ContextIDs:   []string{"context-client-work"},
			EntityIDs:    []string{"entity-sam"},
			BasePriority: 0.70,
			Effort:       EffortAttention{EstimatedAttention: 20 * time.Minute},
		},
		{
			OpenLoopID:   "loop-nonmatch-high",
			PersonID:     "person-1",
			Label:        "Review unrelated receipt",
			ContextIDs:   []string{"context-finance"},
			BasePriority: 0.80,
			Effort:       EffortAttention{EstimatedAttention: 20 * time.Minute},
		},
		{
			OpenLoopID:   "loop-nonmatch-second",
			PersonID:     "person-1",
			Label:        "Plan travel",
			ContextIDs:   []string{"context-travel"},
			BasePriority: 0.75,
			Effort:       EffortAttention{EstimatedAttention: 20 * time.Minute},
		},
	}
	allocation, err := AllocateAttention(budget, candidates, DefaultV02Policy())
	if err != nil {
		t.Fatalf("allocate attention: %v", err)
	}
	if len(allocation.Surfaced) != 1 || len(allocation.Deferred) != 2 {
		t.Fatalf("budget must constrain simultaneous resurfacing, got %+v", allocation)
	}
	if allocation.Surfaced[0].OpenLoopID != "loop-context-match" || !allocation.Surfaced[0].ContextMatched {
		t.Fatalf("expected matching context to resurface first, got %+v", allocation.Surfaced[0])
	}
	if allocation.UsedAttention != 20*time.Minute || allocation.RemainingCapacity != 10*time.Minute {
		t.Fatalf("expected finite capacity accounting, got %+v", allocation)
	}
}

func TestV02LowSubjectivePriorityHighObjectiveStakesMustSurfaceWithoutGoalOverride(t *testing.T) {
	now := fixedV02Now()
	opportunity := v02Opportunity(now, "opportunity-mismatch", PriorityFactors{
		SubjectiveImportance: 0.10,
		ObjectiveStakes:      0.95,
		ExpectedImpact:       0.80,
		Reversibility:        0.70,
		Uncertainty:          0.10,
		OpportunityCost:      0.10,
		EffortAttentionCost:  0.10,
	})
	goal := v02Goal(now, "goal-1")
	evaluation, err := EvaluateOpportunityWithPolicy(&opportunity, []Goal{goal}, nil, DefaultV02Policy(), EvaluationMoment{WallClockAt: now})
	if err != nil {
		t.Fatalf("evaluate mismatch: %v", err)
	}
	if !evaluation.Mismatch.MustSurface || !evaluation.Mismatch.PreservesGoalAuthority {
		t.Fatalf("expected explicit must-surface mismatch that preserves goal authority, got %+v", evaluation.Mismatch)
	}
	decision, err := DecideOpportunityWithPolicy(opportunity, "decision-mismatch", DefaultV02Policy(), EvaluationMoment{WallClockAt: now})
	if err != nil {
		t.Fatalf("decide mismatch: %v", err)
	}
	if decision.Kind != DecisionSurface {
		t.Fatalf("expected must-surface decision, got %+v", decision)
	}
}

func TestV02DeadlineFeasibilityDistinguishesUrgentFromImpossible(t *testing.T) {
	now := fixedV02Now()
	feasible, err := EvaluateDeadlineFeasibility(ActionWindow{LatestSafeActionAt: now.Add(45 * time.Minute), EstimatedEffort: 30 * time.Minute}, EvaluationMoment{WallClockAt: now})
	if err != nil || feasible.State != DeadlineFeasible {
		t.Fatalf("expected feasible urgent action window, got %+v err=%v", feasible, err)
	}
	infeasible, err := EvaluateDeadlineFeasibility(ActionWindow{LatestSafeActionAt: now.Add(10 * time.Minute), EstimatedEffort: 30 * time.Minute}, EvaluationMoment{WallClockAt: now})
	if err != nil || infeasible.State != DeadlineInfeasible {
		t.Fatalf("expected impossible action window, got %+v err=%v", infeasible, err)
	}

	opportunity := v02Opportunity(now, "opportunity-infeasible", fullPriority())
	opportunity.ActionWindow = ActionWindow{LatestSafeActionAt: now.Add(10 * time.Minute), EstimatedEffort: 30 * time.Minute}
	goal := v02Goal(now, "goal-1")
	if _, err := EvaluateOpportunityWithPolicy(&opportunity, []Goal{goal}, nil, DefaultV02Policy(), EvaluationMoment{WallClockAt: now}); err != nil {
		t.Fatalf("evaluate infeasible opportunity: %v", err)
	}
	decision, err := DecideOpportunityWithPolicy(opportunity, "decision-infeasible", DefaultV02Policy(), EvaluationMoment{WallClockAt: now})
	if err != nil {
		t.Fatalf("decide infeasible opportunity: %v", err)
	}
	if decision.Kind != DecisionDefer || opportunity.Evaluation.Deadline.State != DeadlineInfeasible {
		t.Fatalf("infeasible task must differ from feasible urgency, got decision=%+v evaluation=%+v", decision, opportunity.Evaluation)
	}
}

func TestV02ReversibilityAndOpportunityCostChangeEvaluationIndependently(t *testing.T) {
	now := fixedV02Now()
	policy := DefaultV02Policy()
	policy.SubjectiveImportanceWeight = 0
	policy.ObjectiveStakesWeight = 0
	policy.ExpectedImpactWeight = 0
	policy.ReversibilityWeight = 0.50
	policy.UncertaintyPenaltyWeight = 0
	policy.OpportunityCostPenaltyWeight = 0.50
	policy.EffortAttentionPenaltyWeight = 0
	policy.SoftConstraintPenalty = 0
	policy.FeasibleUrgencyWeight = 0
	policy.RecommendAt = 0.90
	policy.DeferAt = 0.10
	goal := v02Goal(now, "goal-1")
	lowCostReversible := v02Opportunity(now, "opportunity-reversible", PriorityFactors{Reversibility: 1, OpportunityCost: 0})
	irreversibleHighCost := v02Opportunity(now, "opportunity-irreversible", PriorityFactors{Reversibility: 0, OpportunityCost: 1})
	first, err := EvaluateOpportunityWithPolicy(&lowCostReversible, []Goal{goal}, nil, policy, EvaluationMoment{WallClockAt: now})
	if err != nil {
		t.Fatalf("evaluate reversible opportunity: %v", err)
	}
	second, err := EvaluateOpportunityWithPolicy(&irreversibleHighCost, []Goal{goal}, nil, policy, EvaluationMoment{WallClockAt: now})
	if err != nil {
		t.Fatalf("evaluate irreversible opportunity: %v", err)
	}
	if first.Utility <= second.Utility || first.Utility != 0.5 || second.Utility != 0 {
		t.Fatalf("expected independent reversibility and opportunity-cost effect, first=%+v second=%+v", first, second)
	}
}

func TestV02StaleClaimCanBeExplicitlyRevalidated(t *testing.T) {
	now := fixedV02Now()
	claim := Claim{
		ID:       "claim-stale",
		PersonID: "person-1",
		Freshness: FreshnessState{
			StaleAfter:      24 * time.Hour,
			LastValidatedAt: now.Add(-48 * time.Hour),
			Status:          FreshnessFresh,
		},
	}
	status, err := EvaluateClaimFreshness(claim, EvaluationMoment{WallClockAt: now})
	if err != nil || status != FreshnessStale {
		t.Fatalf("expected stale claim, got %s err=%v", status, err)
	}
	if err := RevalidateClaim(&claim, EvaluationMoment{WallClockAt: now}); err != nil {
		t.Fatalf("revalidate claim: %v", err)
	}
	status, err = EvaluateClaimFreshness(claim, EvaluationMoment{WallClockAt: now})
	if err != nil || status != FreshnessFresh || !claim.Freshness.LastRevalidatedAt.Equal(now) {
		t.Fatalf("expected explicitly revalidated fresh claim, got %+v status=%s err=%v", claim.Freshness, status, err)
	}
}

func TestV02AuthoritativeContradictionSupersedesWithoutDeletingHistory(t *testing.T) {
	now := fixedV02Now()
	oldClaim := Claim{
		ID:        "claim-old",
		PersonID:  "person-1",
		Temporal:  v02Temporal(now.Add(-48 * time.Hour)),
		Freshness: FreshnessState{LastValidatedAt: now.Add(-24 * time.Hour), Status: FreshnessFresh},
	}
	newClaim := Claim{
		ID:        "claim-new",
		PersonID:  "person-1",
		Temporal:  v02Temporal(now.Add(-24 * time.Hour)),
		Freshness: FreshnessState{LastValidatedAt: now, Status: FreshnessFresh},
	}
	before, err := SelectCurrentClaim([]Claim{oldClaim, newClaim}, EvaluationMoment{WallClockAt: now})
	if err != nil {
		t.Fatalf("select before explicit supersession: %v", err)
	}
	if before.CurrentClaimID != "" {
		t.Fatalf("newness alone must not select a current claim, got %+v", before)
	}
	evidence := Evidence{
		ID:         "evidence-authoritative-new",
		PersonID:   "person-1",
		ClaimID:    newClaim.ID,
		Stance:     EvidenceContradicts,
		Quality:    0.95,
		Relevance:  0.95,
		Authority:  0.95,
		Provenance: Provenance{SourceType: "signed-record", SourceRef: "record://v2", CapturedAt: now, Checksum: "authoritative-v2"},
	}
	policy := SupersessionPolicy{MinimumAuthority: 0.80, MinimumRelevance: 0.80, RequireContradiction: true}
	if err := SupersedeClaim(&oldClaim, &newClaim, evidence, policy, EvaluationMoment{WallClockAt: now}); err != nil {
		t.Fatalf("supersede with newer authoritative contradiction: %v", err)
	}
	selection, err := SelectCurrentClaim([]Claim{oldClaim, newClaim}, EvaluationMoment{WallClockAt: now})
	if err != nil {
		t.Fatalf("select after supersession: %v", err)
	}
	if selection.CurrentClaimID != newClaim.ID || !contains(selection.HistoricalClaimIDs, oldClaim.ID) {
		t.Fatalf("expected new current claim while older claim remains historical, got %+v", selection)
	}
	if oldClaim.Freshness.Status != FreshnessSuperseded || newClaim.SupersedesID != oldClaim.ID || !newClaim.Lineage.PreservesHistory || !contains(newClaim.Lineage.EvidenceIDs, evidence.ID) {
		t.Fatalf("expected preserved supersession lineage and provenance, old=%+v new=%+v", oldClaim, newClaim)
	}

	lowAuthority := evidence
	lowAuthority.ID = "evidence-low-authority"
	lowAuthority.Authority = 0.10
	anotherClaim := newClaim
	anotherClaim.ID = "claim-no-auto-supersession"
	lowAuthority.ClaimID = anotherClaim.ID
	if err := SupersedeClaim(&newClaim, &anotherClaim, lowAuthority, policy, EvaluationMoment{WallClockAt: now}); !errors.Is(err, ErrSupersessionDenied) {
		t.Fatalf("newer evidence without authority must not supersede automatically, got %v", err)
	}
}

func TestV02InjectedPolicyConfigurationsChangeDecisionWithoutChangingDomainData(t *testing.T) {
	now := fixedV02Now()
	goal := v02Goal(now, "goal-1")
	original := v02Opportunity(now, "opportunity-policy", PriorityFactors{
		SubjectiveImportance: 0.50,
		ObjectiveStakes:      0.50,
		ExpectedImpact:       0.50,
		Reversibility:        0.50,
		Uncertainty:          0.10,
		OpportunityCost:      0.10,
		EffortAttentionCost:  0.10,
	})

	permissive := DefaultV02Policy()
	permissive.RecommendAt = 0.30
	permissive.DeferAt = 0.20
	strict := DefaultV02Policy()
	strict.RecommendAt = 0.90
	strict.DeferAt = 0.80

	withPermissive := original
	if _, err := EvaluateOpportunityWithPolicy(&withPermissive, []Goal{goal}, nil, permissive, EvaluationMoment{WallClockAt: now}); err != nil {
		t.Fatalf("evaluate permissive policy: %v", err)
	}
	permissiveDecision, err := DecideOpportunityWithPolicy(withPermissive, "decision-permissive", permissive, EvaluationMoment{WallClockAt: now})
	if err != nil {
		t.Fatalf("decide permissive policy: %v", err)
	}

	withStrict := original
	if _, err := EvaluateOpportunityWithPolicy(&withStrict, []Goal{goal}, nil, strict, EvaluationMoment{WallClockAt: now}); err != nil {
		t.Fatalf("evaluate strict policy: %v", err)
	}
	strictDecision, err := DecideOpportunityWithPolicy(withStrict, "decision-strict", strict, EvaluationMoment{WallClockAt: now})
	if err != nil {
		t.Fatalf("decide strict policy: %v", err)
	}
	if permissiveDecision.Kind != DecisionRecommend || strictDecision.Kind != DecisionDecline || original.Priority != withPermissive.Priority || original.Priority != withStrict.Priority {
		t.Fatalf("same canonical data must support distinct injected policy decisions, permissive=%+v strict=%+v", permissiveDecision, strictDecision)
	}
}

func fixedV02Now() time.Time {
	return time.Date(2026, time.August, 20, 9, 0, 0, 0, time.UTC)
}

func v02Temporal(now time.Time) TemporalState {
	return TemporalState{
		EventAt:     now,
		RecordedAt:  now,
		EffectiveAt: now.Add(-time.Hour),
		AttentionAt: now,
	}
}

func v02Goal(now time.Time, id string) Goal {
	return Goal{
		ID:                   id,
		PersonID:             "person-1",
		SubjectiveImportance: 0.5,
		Status:               GoalActive,
		Temporal:             v02Temporal(now),
	}
}

func v02Opportunity(now time.Time, id string, factors PriorityFactors) Opportunity {
	return Opportunity{
		ID:            id,
		PersonID:      "person-1",
		GoalIDs:       []string{"goal-1"},
		Temporal:      v02Temporal(now),
		Priority:      factors,
		AttentionNeed: EffortAttention{EstimatedEffort: 10 * time.Minute, EstimatedAttention: 10 * time.Minute},
	}
}

func fullPriority() PriorityFactors {
	return PriorityFactors{
		SubjectiveImportance: 0.8,
		ObjectiveStakes:      0.8,
		ExpectedImpact:       0.8,
		Reversibility:        0.8,
		Uncertainty:          0.1,
		OpportunityCost:      0.1,
		EffortAttentionCost:  0.1,
	}
}

func TestV02OpenLoopFeedsFiniteContextualAttentionAllocation(t *testing.T) {
	now := fixedV02Now()
	loop := OpenLoop{
		ID:            "loop-from-domain",
		PersonID:      "person-1",
		Label:         "Follow up with Sam",
		ContextIDs:    []string{"context-client-work"},
		EntityIDs:     []string{"entity-sam"},
		AttentionNeed: EffortAttention{EstimatedAttention: 15 * time.Minute},
	}
	candidate, err := AttentionCandidateFromOpenLoop(loop, 0.60)
	if err != nil {
		t.Fatalf("derive candidate from unresolved open loop: %v", err)
	}
	budget := AttentionBudget{
		ID:                "budget-domain-loop",
		PersonID:          "person-1",
		WindowStart:       now,
		WindowEnd:         now.Add(30 * time.Minute),
		AttentionCapacity: 20 * time.Minute,
		MaxCompetingItems: 1,
		CurrentContext:    AttentionContext{ContextIDs: []string{"context-client-work"}},
	}
	allocation, err := AllocateAttention(budget, []AttentionCandidate{candidate}, DefaultV02Policy())
	if err != nil {
		t.Fatalf("allocate unresolved open loop: %v", err)
	}
	if len(allocation.Surfaced) != 1 || allocation.Surfaced[0].OpenLoopID != loop.ID || !allocation.Surfaced[0].ContextMatched {
		t.Fatalf("expected contextual resurfacing from canonical open loop, got %+v", allocation)
	}
	resolvedAt := now
	loop.ResolvedAt = &resolvedAt
	resolvedCandidate, err := AttentionCandidateFromOpenLoop(loop, 0.60)
	if err != nil {
		t.Fatalf("derive candidate from resolved loop: %v", err)
	}
	allocation, err = AllocateAttention(budget, []AttentionCandidate{resolvedCandidate}, DefaultV02Policy())
	if err != nil {
		t.Fatalf("allocate resolved loop: %v", err)
	}
	if len(allocation.Surfaced) != 0 || len(allocation.Deferred) != 1 || allocation.Deferred[0].Reason != "open loop is already resolved" {
		t.Fatalf("resolved loop must not compete for attention, got %+v", allocation)
	}
}
