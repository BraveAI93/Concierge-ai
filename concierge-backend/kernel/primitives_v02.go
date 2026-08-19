package kernel

import "sort"

// EvaluateDeadlineFeasibility compares the latest safe action time with the
// estimated effort. It does not use record expiry as a substitute for a task
// deadline and explicitly preserves the infeasible state.
func EvaluateDeadlineFeasibility(window ActionWindow, moment EvaluationMoment) (DeadlineFeasibility, error) {
	if err := window.Validate(); err != nil || moment.Validate() != nil {
		return DeadlineFeasibility{}, ErrInfeasibleActionWindow
	}
	if !window.HasDeadline() {
		return DeadlineFeasibility{State: DeadlineNone, EstimatedEffort: window.EstimatedEffort, Reason: "no latest safe action time supplied"}, nil
	}
	remaining := window.LatestSafeActionAt.Sub(moment.WallClockAt)
	result := DeadlineFeasibility{
		LatestSafeAt:    window.LatestSafeActionAt,
		Remaining:       remaining,
		EstimatedEffort: window.EstimatedEffort,
	}
	if remaining < window.EstimatedEffort {
		result.State = DeadlineInfeasible
		result.Reason = "remaining action window is shorter than estimated effort"
		return result, nil
	}
	result.State = DeadlineFeasible
	result.Reason = "estimated effort fits within the remaining action window"
	return result, nil
}

// EvaluateOpportunityWithPolicy enforces canonical invariants and delegates
// all arbitrary score weights, soft penalties, mismatch thresholds, and
// recommend/defer cut-offs to the injected EvaluationPolicy.
func EvaluateOpportunityWithPolicy(opportunity *Opportunity, goals []Goal, constraints []Constraint, policy EvaluationPolicy, moment EvaluationMoment) (OpportunityEvaluation, error) {
	if opportunity == nil || normalizeText(opportunity.ID) == "" || normalizeText(opportunity.PersonID) == "" || policy == nil || moment.Validate() != nil {
		return OpportunityEvaluation{}, ErrMissingIdentifier
	}
	if err := opportunity.Temporal.Validate(); err != nil {
		return OpportunityEvaluation{}, err
	}
	if err := opportunity.Priority.Validate(); err != nil {
		return OpportunityEvaluation{}, err
	}
	if err := opportunity.AttentionNeed.Validate(); err != nil {
		return OpportunityEvaluation{}, err
	}
	if err := policy.Validate(); err != nil {
		return OpportunityEvaluation{}, err
	}

	at := moment.WallClockAt
	evaluation := OpportunityEvaluation{EvaluatedAt: at}
	if !opportunity.Temporal.IsActive(at) {
		evaluation.DecisionBasis = "opportunity is not currently effective"
		opportunity.Evaluation = evaluation
		return evaluation, nil
	}
	if _, activeGoalCount := activeGoalImportance(opportunity, goals, at); activeGoalCount == 0 {
		evaluation.DecisionBasis = "no referenced active goal"
		opportunity.Evaluation = evaluation
		return evaluation, nil
	}

	softConstraintCount := 0
	for _, constraint := range constraints {
		if constraint.PersonID != opportunity.PersonID {
			return OpportunityEvaluation{}, ErrPersonBoundary
		}
		if !contains(opportunity.ConstraintIDs, constraint.ID) || !constraint.Active || !constraint.Temporal.IsActive(at) {
			continue
		}
		if constraint.Kind == ConstraintHard {
			evaluation.HardBlocked = true
			evaluation.DecisionBasis = "blocked by hard constraint: " + constraint.Title
			opportunity.Evaluation = evaluation
			return evaluation, nil
		}
		if constraint.Kind == ConstraintSoft {
			softConstraintCount++
		}
	}

	deadline, err := EvaluateDeadlineFeasibility(opportunity.ActionWindow, moment)
	if err != nil {
		return OpportunityEvaluation{}, err
	}
	evaluation.Deadline = deadline
	assessment, err := policy.AssessOpportunity(OpportunityPolicyInput{
		Factors:             opportunity.Priority,
		Deadline:            deadline,
		SoftConstraintCount: softConstraintCount,
		TemporalActive:      true,
	})
	if err != nil {
		return OpportunityEvaluation{}, err
	}
	evaluation.Utility = assessment.Utility
	evaluation.Mismatch = assessment.Mismatch
	evaluation.DecisionBasis = assessment.DecisionBasis
	opportunity.Evaluation = evaluation
	return evaluation, nil
}

// DecideOpportunityWithPolicy is the policy-injected companion to evaluation.
// It retains the exact evaluated utility and emits a decision without storing
// cut-offs in canonical types.
func DecideOpportunityWithPolicy(opportunity Opportunity, decisionID string, policy EvaluationPolicy, moment EvaluationMoment) (Decision, error) {
	if normalizeText(decisionID) == "" || policy == nil || moment.Validate() != nil || normalizeText(opportunity.ID) == "" || normalizeText(opportunity.PersonID) == "" {
		return Decision{}, ErrMissingIdentifier
	}
	if err := policy.Validate(); err != nil {
		return Decision{}, err
	}
	kind, reason, err := policy.Decide(PolicyAssessment{
		Utility:       opportunity.Evaluation.Utility,
		Mismatch:      opportunity.Evaluation.Mismatch,
		DecisionBasis: opportunity.Evaluation.DecisionBasis,
	}, opportunity.Evaluation.HardBlocked, opportunity.Evaluation.Deadline)
	if err != nil {
		return Decision{}, err
	}
	return Decision{
		ID:            decisionID,
		PersonID:      opportunity.PersonID,
		OpportunityID: opportunity.ID,
		Kind:          kind,
		Utility:       opportunity.Evaluation.Utility,
		Reason:        reason,
		CreatedAt:     moment.WallClockAt,
	}, nil
}

// AttentionCandidateFromOpenLoop exposes an unresolved canonical OpenLoop to
// the finite allocator without an adapter inventing duplicate fields. Base
// priority remains an explicit policy input and is not stored as a hidden loop
// score.
func AttentionCandidateFromOpenLoop(loop OpenLoop, basePriority float64) (AttentionCandidate, error) {
	if normalizeText(loop.ID) == "" || normalizeText(loop.PersonID) == "" || !isZeroOrUnit(basePriority) || loop.AttentionNeed.Validate() != nil {
		return AttentionCandidate{}, ErrInvalidAttentionItem
	}
	return AttentionCandidate{
		OpenLoopID:   loop.ID,
		PersonID:     loop.PersonID,
		Label:        loop.Label,
		ResolvedAt:   loop.ResolvedAt,
		ContextIDs:   append([]string(nil), loop.ContextIDs...),
		EntityIDs:    append([]string(nil), loop.EntityIDs...),
		BasePriority: basePriority,
		Effort:       loop.AttentionNeed,
	}, nil
}

// AllocateAttention deterministically surfaces only unresolved, person-scoped
// loops that fit inside a finite AttentionBudget. Candidates are ordered by
// policy score, then by open-loop ID to make tie handling stable.
func AllocateAttention(budget AttentionBudget, candidates []AttentionCandidate, policy EvaluationPolicy) (AttentionAllocation, error) {
	if err := budget.Validate(); err != nil || policy == nil || policy.Validate() != nil {
		return AttentionAllocation{}, ErrInvalidAttentionBudget
	}
	allocation := AttentionAllocation{
		BudgetID:          budget.ID,
		PersonID:          budget.PersonID,
		RemainingCapacity: budget.AttentionCapacity,
	}
	type scoredCandidate struct {
		candidate AttentionCandidate
		score     float64
		reason    string
		matched   bool
	}
	eligible := make([]scoredCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if err := candidate.Validate(); err != nil {
			return AttentionAllocation{}, err
		}
		if candidate.PersonID != budget.PersonID {
			return AttentionAllocation{}, ErrPersonBoundary
		}
		if candidate.ResolvedAt != nil {
			allocation.Deferred = append(allocation.Deferred, AttentionDeferral{OpenLoopID: candidate.OpenLoopID, Reason: "open loop is already resolved"})
			continue
		}
		matched := overlaps(candidate.ContextIDs, budget.CurrentContext.ContextIDs) || overlaps(candidate.EntityIDs, budget.CurrentContext.EntityIDs)
		score, reason, err := policy.ScoreAttention(candidate, matched)
		if err == nil {
			score = clampUnit(score - budget.InterruptionCost*candidate.Effort.InterruptionCost)
		}
		if err != nil {
			return AttentionAllocation{}, err
		}
		eligible = append(eligible, scoredCandidate{candidate: candidate, score: score, reason: reason, matched: matched})
	}
	sort.SliceStable(eligible, func(i, j int) bool {
		if eligible[i].score == eligible[j].score {
			return eligible[i].candidate.OpenLoopID < eligible[j].candidate.OpenLoopID
		}
		return eligible[i].score > eligible[j].score
	})
	for _, candidate := range eligible {
		required := candidate.candidate.Effort.RequiredAttention()
		if len(allocation.Surfaced) >= budget.MaxCompetingItems {
			allocation.Deferred = append(allocation.Deferred, AttentionDeferral{OpenLoopID: candidate.candidate.OpenLoopID, Reason: "maximum competing items reached"})
			continue
		}
		if required > allocation.RemainingCapacity {
			allocation.Deferred = append(allocation.Deferred, AttentionDeferral{OpenLoopID: candidate.candidate.OpenLoopID, Reason: "insufficient remaining attention capacity"})
			continue
		}
		allocation.Surfaced = append(allocation.Surfaced, AttentionSelection{
			OpenLoopID:        candidate.candidate.OpenLoopID,
			Score:             candidate.score,
			ContextMatched:    candidate.matched,
			AttentionReserved: required,
			Reason:            candidate.reason,
		})
		allocation.UsedAttention += required
		allocation.RemainingCapacity -= required
	}
	return allocation, nil
}

func overlaps(left, right []string) bool {
	for _, first := range left {
		for _, second := range right {
			if first != "" && first == second {
				return true
			}
		}
	}
	return false
}

// EvaluateClaimFreshness reports whether a claim is fresh, stale, historical,
// or superseded at an evaluation moment. It never deletes prior claim history.
func EvaluateClaimFreshness(claim Claim, moment EvaluationMoment) (FreshnessStatus, error) {
	if normalizeText(claim.ID) == "" || normalizeText(claim.PersonID) == "" || moment.Validate() != nil || claim.Freshness.Validate() != nil {
		return "", ErrInvalidFreshness
	}
	if claim.Freshness.Status == FreshnessSuperseded || normalizeText(claim.SupersedesID) != "" && claim.Freshness.Status == FreshnessHistorical {
		return FreshnessHistorical, nil
	}
	if claim.Freshness.Status == FreshnessHistorical {
		return FreshnessHistorical, nil
	}
	if claim.Freshness.StaleAfter > 0 && moment.WallClockAt.After(claim.Freshness.LastValidatedAt.Add(claim.Freshness.StaleAfter)) {
		return FreshnessStale, nil
	}
	return FreshnessFresh, nil
}

// RevalidateClaim creates a new state on the same historical claim. It never
// replaces evidence or provenance, and resets freshness only from an explicit
// validation moment.
func RevalidateClaim(claim *Claim, moment EvaluationMoment) error {
	if claim == nil || moment.Validate() != nil || claim.Freshness.Validate() != nil {
		return ErrInvalidFreshness
	}
	if moment.WallClockAt.Before(claim.Freshness.LastValidatedAt) {
		return ErrInvalidFreshness
	}
	claim.Freshness.LastValidatedAt = moment.WallClockAt
	claim.Freshness.LastRevalidatedAt = moment.WallClockAt
	claim.Freshness.Status = FreshnessFresh
	return nil
}

// SupersedeClaim requires explicit contradictory, authoritative, relevant, and
// provenance-backed evidence. Recency alone cannot invoke this transition.
func SupersedeClaim(previous *Claim, replacement *Claim, evidence Evidence, policy SupersessionPolicy, moment EvaluationMoment) error {
	if previous == nil || replacement == nil || moment.Validate() != nil || policy.Validate() != nil {
		return ErrSupersessionDenied
	}
	if normalizeText(previous.ID) == "" || normalizeText(replacement.ID) == "" || previous.PersonID == "" || previous.PersonID != replacement.PersonID || evidence.PersonID != replacement.PersonID || evidence.ClaimID != replacement.ID {
		return ErrPersonBoundary
	}
	if evidence.Stance != EvidenceContradicts && policy.RequireContradiction {
		return ErrSupersessionDenied
	}
	if evidence.Authority < policy.MinimumAuthority || evidence.Relevance < policy.MinimumRelevance || !hasProvenance(evidence.Provenance) {
		return ErrSupersessionDenied
	}
	if replacement.Temporal.EventAt.Before(previous.Temporal.EventAt) {
		return ErrSupersessionDenied
	}
	previous.Freshness.Status = FreshnessSuperseded
	previous.Lineage.ClaimID = previous.ID
	previous.Lineage.PreservesHistory = true
	replacement.SupersedesID = previous.ID
	replacement.Lineage = ClaimLineage{
		ClaimID:           replacement.ID,
		SupersedesClaimID: previous.ID,
		EvidenceIDs:       appendUnique(replacement.Lineage.EvidenceIDs, evidence.ID),
		PreservesHistory:  true,
		RecordedAt:        moment.WallClockAt,
	}
	replacement.EvidenceIDs = appendUnique(replacement.EvidenceIDs, evidence.ID)
	replacement.Freshness.Status = FreshnessFresh
	if replacement.Freshness.LastValidatedAt.IsZero() {
		replacement.Freshness.LastValidatedAt = moment.WallClockAt
	}
	return nil
}

// SelectCurrentClaim reports a current claim only when it is explicitly
// unambiguous after freshness and supersession status are applied. It does not
// choose the newest candidate merely by timestamp.
func SelectCurrentClaim(claims []Claim, moment EvaluationMoment) (ClaimSelection, error) {
	if moment.Validate() != nil {
		return ClaimSelection{}, ErrInvalidFreshness
	}
	selection := ClaimSelection{}
	current := make([]Claim, 0, len(claims))
	for _, claim := range claims {
		status, err := EvaluateClaimFreshness(claim, moment)
		if err != nil {
			return ClaimSelection{}, err
		}
		if status == FreshnessFresh {
			current = append(current, claim)
		} else {
			selection.HistoricalClaimIDs = append(selection.HistoricalClaimIDs, claim.ID)
		}
	}
	if len(current) == 1 {
		selection.CurrentClaimID = current[0].ID
		selection.Reason = "single fresh, non-superseded claim"
		return selection, nil
	}
	if len(current) == 0 {
		selection.Reason = "no fresh current claim; retain historical records"
		return selection, nil
	}
	for _, claim := range current {
		selection.HistoricalClaimIDs = append(selection.HistoricalClaimIDs, claim.ID)
	}
	selection.Reason = "multiple fresh claims require explicit policy selection; newest is not auto-selected"
	return selection, nil
}
