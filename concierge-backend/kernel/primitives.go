package kernel

import (
	"math"
	"time"
)

// LinkMemoryToEvent creates an idempotent reciprocal association after
// verifying that both records belong to the same person.
func LinkMemoryToEvent(memory *Memory, event *Event, reason string, linkedAt time.Time) (MemoryEventLink, error) {
	if memory == nil || event == nil || normalizeText(memory.ID) == "" || normalizeText(event.ID) == "" {
		return MemoryEventLink{}, ErrMissingIdentifier
	}
	if memory.PersonID == "" || event.PersonID == "" || memory.PersonID != event.PersonID {
		return MemoryEventLink{}, ErrPersonBoundary
	}
	if linkedAt.IsZero() {
		return MemoryEventLink{}, ErrInvalidTemporalState
	}
	memory.EventIDs = appendUnique(memory.EventIDs, event.ID)
	event.MemoryIDs = appendUnique(event.MemoryIDs, memory.ID)
	event.ContextIDs = uniqueValues(event.ContextIDs)
	return MemoryEventLink{
		MemoryID: memory.ID,
		EventID:  event.ID,
		Reason:   normalizeText(reason),
		LinkedAt: linkedAt,
	}, nil
}

// LinkEventToMemory is provided for adapters that hold an event reference and
// need an explicit reciprocity check. It delegates to LinkMemoryToEvent.
func LinkEventToMemory(event *Event, memory *Memory, reason string, linkedAt time.Time) (MemoryEventLink, error) {
	return LinkMemoryToEvent(memory, event, reason, linkedAt)
}

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func uniqueValues(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if normalizeText(value) != "" {
			result = appendUnique(result, value)
		}
	}
	return result
}

// TransitionPendingIntent validates the lifecycle graph and appends an
// immutable transition record. Terminal states cannot transition further.
func TransitionPendingIntent(intent *PendingIntent, next PendingIntentState, at time.Time, actor, reason string) error {
	if intent == nil || normalizeText(intent.ID) == "" || normalizeText(intent.PersonID) == "" {
		return ErrMissingIdentifier
	}
	if at.IsZero() || !isAllowedIntentTransition(intent.State, next) {
		return ErrInvalidIntentTransition
	}
	intent.Transitions = append(intent.Transitions, IntentTransition{
		From:   intent.State,
		To:     next,
		At:     at,
		Actor:  normalizeText(actor),
		Reason: normalizeText(reason),
	})
	intent.State = next
	intent.UpdatedAt = at
	return nil
}

func isAllowedIntentTransition(from, to PendingIntentState) bool {
	allowed := map[PendingIntentState]map[PendingIntentState]bool{
		IntentCaptured: {
			IntentClarifying: true, IntentReady: true, IntentCancelled: true, IntentExpired: true,
		},
		IntentClarifying: {
			IntentReady: true, IntentCancelled: true, IntentExpired: true,
		},
		IntentReady: {
			IntentProposed: true, IntentCancelled: true, IntentExpired: true,
		},
		IntentProposed: {
			IntentInProgress: true, IntentCancelled: true, IntentExpired: true,
		},
		IntentInProgress: {
			IntentCompleted: true, IntentCancelled: true, IntentExpired: true,
		},
	}
	return allowed[from][to]
}

// EvaluateTemporalUtility measures attention priority at a supplied reference
// time. It always returns zero for a state that is not currently applicable.
func EvaluateTemporalUtility(temporal TemporalState, importance float64, at time.Time) (TemporalEvaluation, error) {
	if temporal.Validate() != nil || at.IsZero() || !isZeroOrUnit(importance) {
		return TemporalEvaluation{}, ErrInvalidTemporalState
	}
	evaluation := TemporalEvaluation{
		Importance:  importance,
		EvaluatedAt: at,
		Active:      temporal.IsActive(at),
	}
	if !evaluation.Active {
		return evaluation, nil
	}

	ageHours := at.Sub(temporal.EventAt).Hours()
	if ageHours < 0 {
		ageHours = 0
	}
	evaluation.Recency = 1 / (1 + ageHours/(24*7))

	if temporal.ExpiresAt != nil {
		hoursToExpiry := temporal.ExpiresAt.Sub(at).Hours()
		if hoursToExpiry <= 0 {
			evaluation.DeadlineUrgency = 1
		} else {
			evaluation.DeadlineUrgency = 1 / (1 + hoursToExpiry/24)
		}
	}

	hoursToAttention := temporal.AttentionAt.Sub(at).Hours()
	if hoursToAttention <= 0 {
		evaluation.AttentionDue = 1
	} else {
		evaluation.AttentionDue = 1 / (1 + hoursToAttention/24)
	}

	// Retained v0.1 callers use an explicit compatibility policy. Coefficients
	// are no longer embedded in this canonical primitive.
	utility, err := LegacyV01Policy().LegacyTemporalEvaluation(evaluation.Importance, evaluation.DeadlineUrgency, evaluation.Recency, evaluation.AttentionDue)
	if err != nil {
		return TemporalEvaluation{}, err
	}
	evaluation.Utility = utility
	return evaluation, nil
}

// AssessClaimConfidence produces a conflict-aware assessment from explicit
// evidence quality, relevance, stance, and provenance completeness.
func AssessClaimConfidence(claim *Claim, evidence []Evidence, at time.Time) (ConfidenceAssessment, error) {
	if claim == nil || normalizeText(claim.ID) == "" || normalizeText(claim.PersonID) == "" || at.IsZero() {
		return ConfidenceAssessment{}, ErrMissingIdentifier
	}
	assessment := ConfidenceAssessment{Score: 0.5, AssessedAt: at}
	for _, item := range evidence {
		if item.PersonID != claim.PersonID || item.ClaimID != claim.ID {
			return ConfidenceAssessment{}, ErrPersonBoundary
		}
		if !isZeroOrUnit(item.Quality) || !isZeroOrUnit(item.Relevance) {
			return ConfidenceAssessment{}, ErrInvalidTemporalState
		}
		weight := item.Quality * item.Relevance
		assessment.EvidenceCount++
		if hasProvenance(item.Provenance) {
			assessment.ProvenanceCount++
		}
		switch item.Stance {
		case EvidenceSupports:
			assessment.SupportingWeight += weight
		case EvidenceContradicts:
			assessment.ConflictingWeight += weight
		}
	}

	totalWeight := assessment.SupportingWeight + assessment.ConflictingWeight
	if totalWeight == 0 {
		claim.Confidence = assessment
		return assessment, nil
	}

	signedSupport := (assessment.SupportingWeight - assessment.ConflictingWeight) / totalWeight
	coverage := 1 - math.Exp(-totalWeight)
	provenanceCompleteness := float64(assessment.ProvenanceCount) / float64(assessment.EvidenceCount)
	assessment.Score = clampUnit(0.5 + 0.5*signedSupport*coverage*provenanceCompleteness)
	claim.Confidence = assessment
	return assessment, nil
}

func hasProvenance(provenance Provenance) bool {
	return normalizeText(provenance.SourceType) != "" &&
		normalizeText(provenance.SourceRef) != "" &&
		!provenance.CapturedAt.IsZero()
}

// EvaluateOpportunity evaluates candidate utility against active goals,
// person-scoped constraints, confidence, and temporal priority. A matching
// active hard constraint always blocks the opportunity.
// EvaluateOpportunity is a v0.1 compatibility wrapper. It explicitly maps
// legacy generic values into a policy adapter; v0.2 callers must populate the
// independent PriorityFactors and call EvaluateOpportunityWithPolicy.
func EvaluateOpportunity(opportunity *Opportunity, goals []Goal, constraints []Constraint, at time.Time) (OpportunityEvaluation, error) {
	if opportunity == nil || normalizeText(opportunity.ID) == "" || normalizeText(opportunity.PersonID) == "" || at.IsZero() {
		return OpportunityEvaluation{}, ErrMissingIdentifier
	}
	for _, value := range []float64{
		opportunity.GoalAlignment,
		opportunity.ExpectedValue,
		opportunity.Effort,
		opportunity.Risk,
		opportunity.EvidenceConfidence,
		opportunity.TemporalPriority,
	} {
		if !isZeroOrUnit(value) {
			return OpportunityEvaluation{}, ErrInvalidTemporalState
		}
	}
	activeImportance, _ := activeGoalImportance(opportunity, goals, at)
	legacy := *opportunity
	legacy.Priority = PriorityFactors{
		SubjectiveImportance: opportunity.GoalAlignment * activeImportance,
		ObjectiveStakes:      opportunity.ExpectedValue,
		ExpectedImpact:       opportunity.TemporalPriority,
		Reversibility:        1 - opportunity.Risk,
		Uncertainty:          1 - opportunity.EvidenceConfidence,
		OpportunityCost:      0,
		EffortAttentionCost:  opportunity.Effort,
	}
	legacy.AttentionNeed = EffortAttention{}
	evaluation, err := EvaluateOpportunityWithPolicy(&legacy, goals, constraints, LegacyV01Policy(), EvaluationMoment{WallClockAt: at})
	if err != nil {
		return OpportunityEvaluation{}, err
	}
	opportunity.Evaluation = evaluation
	return evaluation, nil
}

func activeGoalImportance(opportunity *Opportunity, goals []Goal, at time.Time) (float64, int) {
	var total float64
	var count int
	for _, goal := range goals {
		if goal.PersonID != opportunity.PersonID || goal.Status != GoalActive || !goal.Temporal.IsActive(at) {
			continue
		}
		if contains(opportunity.GoalIDs, goal.ID) {
			total += goal.Importance
			count++
		}
	}
	if count == 0 {
		return 0, 0
	}
	return total / float64(count), count
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

// DecideOpportunity applies fixed v0.1 thresholds to an evaluation. Thresholds
// are intentionally local and visible until policy configuration is introduced.
// DecideOpportunity is a v0.1 compatibility wrapper. New callers must inject
// a policy through DecideOpportunityWithPolicy rather than use fixed cut-offs.
func DecideOpportunity(opportunity Opportunity, decisionID string, at time.Time) (Decision, error) {
	return DecideOpportunityWithPolicy(opportunity, decisionID, LegacyV01Policy(), EvaluationMoment{WallClockAt: at})
}

// ValidateProposalPermission ensures the permission is person-scoped, active,
// and explicitly grants the exact requested scope.
func ValidateProposalPermission(proposal ActionProposal, permission Permission, at time.Time) error {
	if normalizeText(proposal.ID) == "" || normalizeText(proposal.PersonID) == "" || normalizeText(permission.ID) == "" {
		return ErrMissingIdentifier
	}
	if proposal.PersonID != permission.PersonID {
		return ErrPersonBoundary
	}
	if proposal.PermissionID != permission.ID || !permission.Temporal.IsActive(at) {
		return ErrPermissionDenied
	}
	for _, scope := range permission.Scopes {
		if scope.Matches(proposal.Requested) {
			return nil
		}
	}
	return ErrPermissionDenied
}

// PrepareActionGate validates permission before changing a draft proposal into
// an approval state. It never executes an action.
func PrepareActionGate(gate *ActionGate, proposal ActionProposal, permission Permission, at time.Time, actor string) error {
	if gate == nil || gate.ActionProposalID != proposal.ID || gate.PersonID != proposal.PersonID {
		return ErrPersonBoundary
	}
	if err := ValidateProposalPermission(proposal, permission, at); err != nil {
		return err
	}
	if permission.RequiresApproval || !permission.CanAutoApprove {
		return transitionActionGate(gate, GateAwaitingApproval, at, actor, "permission validation completed; awaiting approval")
	}
	return transitionActionGate(gate, GateApproved, at, actor, "permission validation completed; auto-approval allowed")
}

// ApproveAction records a human approval only after permission validation.
func ApproveAction(gate *ActionGate, proposal ActionProposal, permission Permission, at time.Time, actor string) error {
	if normalizeText(actor) == "" {
		return ErrMissingIdentifier
	}
	if err := ValidateProposalPermission(proposal, permission, at); err != nil {
		return err
	}
	if gate == nil || gate.State != GateAwaitingApproval {
		return ErrMissingApproval
	}
	return transitionActionGate(gate, GateApproved, at, actor, "human approval")
}

// ExecuteApprovedAction marks only the gate state. The calling adapter remains
// responsible for an external side effect and must separately record an Outcome.
func ExecuteApprovedAction(gate *ActionGate, at time.Time, actor string) error {
	if gate == nil || gate.State != GateApproved {
		return ErrMissingApproval
	}
	return transitionActionGate(gate, GateExecuted, at, actor, "action dispatched by adapter")
}

func transitionActionGate(gate *ActionGate, next ActionGateState, at time.Time, actor, reason string) error {
	if gate == nil || normalizeText(gate.ID) == "" || at.IsZero() || !isAllowedGateTransition(gate.State, next) {
		return ErrInvalidGateTransition
	}
	gate.Transitions = append(gate.Transitions, GateTransition{
		From:   gate.State,
		To:     next,
		At:     at,
		Actor:  normalizeText(actor),
		Reason: normalizeText(reason),
	})
	gate.State = next
	gate.UpdatedAt = at
	return nil
}

func isAllowedGateTransition(from, to ActionGateState) bool {
	allowed := map[ActionGateState]map[ActionGateState]bool{
		GateDraft: {
			GateAwaitingApproval: true, GateApproved: true,
		},
		GateAwaitingApproval: {
			GateApproved: true, GateRejected: true, GateExpired: true,
		},
		GateApproved: {
			GateExecuted: true, GateExpired: true,
		},
	}
	return allowed[from][to]
}

// AuditOutcome requires evidence and provenance, then compares observed versus
// expected utility to produce a deterministic learning classification.
func AuditOutcome(outcome Outcome, auditID string, at time.Time) (SelfAudit, error) {
	if normalizeText(outcome.ID) == "" || normalizeText(outcome.PersonID) == "" || normalizeText(auditID) == "" || at.IsZero() {
		return SelfAudit{}, ErrMissingIdentifier
	}
	if !isZeroOrUnit(outcome.ExpectedUtility) || !isZeroOrUnit(outcome.ObservedUtility) {
		return SelfAudit{}, ErrInvalidTemporalState
	}
	if len(outcome.EvidenceIDs) == 0 || len(outcome.Provenance) == 0 {
		return SelfAudit{}, ErrMissingOutcomeEvidence
	}
	for _, provenance := range outcome.Provenance {
		if !hasProvenance(provenance) {
			return SelfAudit{}, ErrMissingOutcomeEvidence
		}
	}

	delta := outcome.ObservedUtility - outcome.ExpectedUtility
	audit := SelfAudit{
		ID:                 auditID,
		PersonID:           outcome.PersonID,
		OutcomeID:          outcome.ID,
		ExpectedUtility:    outcome.ExpectedUtility,
		ObservedUtility:    outcome.ObservedUtility,
		UtilityDelta:       delta,
		EvidenceSufficient: true,
		CreatedAt:          at,
	}

	switch {
	case outcome.Status == OutcomeFailed && outcome.ExpectedUtility >= 0.5:
		audit.Status = AuditEscalate
		audit.Summary = "high-expectation action failed; review assumptions and execution context"
	case math.Abs(delta) <= 0.15:
		audit.Status = AuditConfirmed
		audit.Summary = "observed utility matched the expected utility"
	case delta < 0:
		audit.Status = AuditLearned
		audit.Summary = "observed utility was lower than expected; reduce confidence in similar proposals"
	default:
		audit.Status = AuditAdjusted
		audit.Summary = "observed utility exceeded expectation; refine future opportunity estimates"
	}
	return audit, nil
}
