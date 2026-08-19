package kernel

import (
	"errors"
	"testing"
	"time"
)

func TestSyntheticEndToEndPersonCentricKernelScenario(t *testing.T) {
	// A fixed timestamp makes every stage, score, and state transition
	// deterministic and independent of the execution environment.
	now := time.Date(2026, time.August, 19, 9, 0, 0, 0, time.UTC)
	expiresAt := now.Add(2 * time.Hour)
	temporal := TemporalState{
		EventAt:     now.Add(-time.Hour),
		RecordedAt:  now.Add(-55 * time.Minute),
		EffectiveAt: now.Add(-time.Hour),
		AttentionAt: now.Add(-10 * time.Minute),
		ExpiresAt:   &expiresAt,
	}

	// 1. A personal world contains a canonical entity and current context.
	person := Person{ID: "person-ava", DisplayName: "Ava", WorldID: "world-ava", CreatedAt: now.Add(-24 * time.Hour)}
	world := PersonalWorld{ID: person.WorldID, PersonID: person.ID, EntityIDs: []string{"entity-client"}, ContextIDs: []string{"context-meeting"}, UpdatedAt: now}
	client := Entity{
		ID:        "entity-client",
		PersonID:  person.ID,
		Kind:      EntityPerson,
		Name:      "Sam Rivera",
		Aliases:   []Alias{{ID: "alias-sam", EntityID: "entity-client", Value: "Sam", Normalized: "sam", Source: "calendar", CreatedAt: now}},
		CreatedAt: now,
	}
	meetingContext := Context{ID: "context-meeting", PersonID: person.ID, Kind: "meeting", Label: "Client scheduling", EntityIDs: []string{client.ID}, Temporal: temporal, CreatedAt: now}
	if world.PersonID != person.ID || meetingContext.PersonID != person.ID || client.PersonID != person.ID {
		t.Fatal("personal world objects must share the person boundary")
	}

	// 2. An observed event is linked to a memory and its claim is supported by
	// source-attributed evidence.
	event := Event{
		ID:         "event-reschedule-request",
		PersonID:   person.ID,
		Kind:       "calendar_message",
		Summary:    "Sam asked to reschedule tomorrow's meeting.",
		EntityIDs:  []string{client.ID},
		ContextIDs: []string{meetingContext.ID},
		Temporal:   temporal,
		CreatedAt:  now,
	}
	memory := Memory{
		ID:         "memory-reschedule",
		PersonID:   person.ID,
		Kind:       MemoryEpisodic,
		Summary:    "Sam needs a meeting reschedule.",
		ContextIDs: []string{meetingContext.ID},
		Temporal:   temporal,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	link, err := LinkMemoryToEvent(&memory, &event, "event establishes the unfinished scheduling request", now)
	if err != nil {
		t.Fatalf("link memory to event: %v", err)
	}
	if link.MemoryID != memory.ID || link.EventID != event.ID || !contains(memory.EventIDs, event.ID) || !contains(event.MemoryIDs, memory.ID) {
		t.Fatalf("expected reciprocal memory/event link, got memory=%+v event=%+v link=%+v", memory, event, link)
	}

	claim := Claim{
		ID:        "claim-reschedule-needed",
		PersonID:  person.ID,
		MemoryID:  memory.ID,
		Statement: "Sam needs a new meeting time.",
		SubjectID: client.ID,
		Predicate: "needs",
		Object:    "meeting reschedule",
		Temporal:  temporal,
		CreatedAt: now,
	}
	evidence := Evidence{
		ID:        "evidence-calendar-message",
		PersonID:  person.ID,
		ClaimID:   claim.ID,
		Stance:    EvidenceSupports,
		Summary:   "Calendar message explicitly asks for a different meeting time.",
		Quality:   0.95,
		Relevance: 0.95,
		Provenance: Provenance{
			SourceType: "calendar_message",
			SourceRef:  "calendar://messages/evt-42",
			Actor:      "Sam Rivera",
			CapturedAt: now,
			Checksum:   "synthetic-calendar-message-v1",
		},
		Temporal:  temporal,
		CreatedAt: now,
	}
	confidence, err := AssessClaimConfidence(&claim, []Evidence{evidence}, now)
	if err != nil {
		t.Fatalf("assess claim confidence: %v", err)
	}
	if confidence.Score <= 0.5 || confidence.ProvenanceCount != 1 || claim.Confidence.Score != confidence.Score {
		t.Fatalf("expected provenance-backed positive confidence, got %+v", confidence)
	}

	// 3. The event creates both an open loop and a stateful pending intention.
	intent := PendingIntent{
		ID:         "intent-reschedule",
		PersonID:   person.ID,
		Summary:    "Offer Sam alternative meeting times.",
		State:      IntentCaptured,
		GoalID:     "goal-client-relationship",
		MemoryID:   memory.ID,
		ContextIDs: []string{meetingContext.ID},
		Temporal:   temporal,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	for _, next := range []PendingIntentState{IntentClarifying, IntentReady, IntentProposed} {
		if err := TransitionPendingIntent(&intent, next, now, "kernel-test", "synthetic lifecycle progression"); err != nil {
			t.Fatalf("transition pending intent to %s: %v", next, err)
		}
	}
	openLoop := OpenLoop{
		ID:              "loop-reschedule",
		PersonID:        person.ID,
		PendingIntentID: intent.ID,
		Label:           "Respond to Sam's rescheduling request",
		Attention:       temporal,
		CreatedAt:       now,
	}
	if intent.State != IntentProposed || openLoop.PendingIntentID != intent.ID {
		t.Fatalf("expected a proposed pending intent with matching open loop: intent=%+v loop=%+v", intent, openLoop)
	}

	// 4. Four-dimensional time produces attention priority before an opportunity
	// is evaluated. The evaluation sees event, record, effective, and attention
	// time separately in TemporalState.
	temporalPriority, err := EvaluateTemporalUtility(openLoop.Attention, 0.90, now)
	if err != nil {
		t.Fatalf("evaluate temporal utility: %v", err)
	}
	if !temporalPriority.Active || temporalPriority.AttentionDue != 1 || temporalPriority.Utility < 0.7 {
		t.Fatalf("expected high active attention priority, got %+v", temporalPriority)
	}

	goal := Goal{
		ID:         intent.GoalID,
		PersonID:   person.ID,
		Title:      "Maintain responsive client relationships",
		Importance: 0.90,
		Status:     GoalActive,
		Temporal:   temporal,
		CreatedAt:  now,
	}
	opportunity := Opportunity{
		ID:                 "opportunity-draft-reschedule",
		PersonID:           person.ID,
		Title:              "Draft a rescheduling reply for Sam",
		Summary:            "Propose three available meeting times based on the scheduling request.",
		GoalIDs:            []string{goal.ID},
		EvidenceIDs:        []string{evidence.ID},
		Temporal:           temporal,
		GoalAlignment:      0.95,
		ExpectedValue:      0.94,
		Effort:             0.20,
		Risk:               0.10,
		EvidenceConfidence: confidence.Score,
		TemporalPriority:   temporalPriority.Utility,
		CreatedAt:          now,
	}
	evaluation, err := EvaluateOpportunity(&opportunity, []Goal{goal}, nil, now)
	if err != nil {
		t.Fatalf("evaluate opportunity: %v", err)
	}
	if evaluation.HardBlocked || evaluation.Utility < 0.65 {
		t.Fatalf("expected a viable high-utility opportunity, got %+v", evaluation)
	}
	decision, err := DecideOpportunity(opportunity, "decision-draft-reschedule", now)
	if err != nil {
		t.Fatalf("decide opportunity: %v", err)
	}
	if decision.Kind != DecisionRecommend {
		t.Fatalf("expected recommendation, got %+v", decision)
	}

	// 5. An action proposal is scope-checked, then receives an explicit human
	// approval before the kernel marks it executed. No external action occurs.
	permission := Permission{
		ID:               "permission-create-draft",
		PersonID:         person.ID,
		Scopes:           []Scope{{Capability: "create_draft", Resource: "email:client"}},
		RequiresApproval: true,
		CanAutoApprove:   false,
		Temporal:         temporal,
		GrantedBy:        person.ID,
		CreatedAt:        now,
	}
	proposal := ActionProposal{
		ID:            "proposal-draft-reschedule",
		PersonID:      person.ID,
		OpportunityID: opportunity.ID,
		DecisionID:    decision.ID,
		Title:         "Create, but do not send, a rescheduling email draft",
		Requested:     Scope{Capability: "create_draft", Resource: "email:client"},
		PermissionID:  permission.ID,
		Parameters:    map[string]string{"recipient_entity_id": client.ID},
		CreatedAt:     now,
	}
	gate := ActionGate{
		ID:               "gate-draft-reschedule",
		PersonID:         person.ID,
		ActionProposalID: proposal.ID,
		State:            GateDraft,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := PrepareActionGate(&gate, proposal, permission, now, "kernel-test"); err != nil {
		t.Fatalf("prepare action gate: %v", err)
	}
	if gate.State != GateAwaitingApproval {
		t.Fatalf("expected approval gate, got %s", gate.State)
	}
	if err := ApproveAction(&gate, proposal, permission, now, person.ID); err != nil {
		t.Fatalf("approve action: %v", err)
	}
	if err := ExecuteApprovedAction(&gate, now, "synthetic-executor"); err != nil {
		t.Fatalf("execute approved action: %v", err)
	}
	if gate.State != GateExecuted || len(gate.Transitions) != 3 {
		t.Fatalf("expected draft→awaiting_approval→approved→executed, got %+v", gate)
	}

	// 6. An observed outcome carries evidence and provenance, which supports an
	// auditable comparison between expected and realized utility.
	outcome := Outcome{
		ID:               "outcome-draft-created",
		PersonID:         person.ID,
		ActionProposalID: proposal.ID,
		Status:           OutcomeSucceeded,
		Summary:          "A draft containing three alternate slots was created.",
		ExpectedUtility:  decision.Utility,
		ObservedUtility:  0.82,
		EvidenceIDs:      []string{"evidence-draft-receipt"},
		Provenance: []Provenance{{
			SourceType: "synthetic_email_adapter",
			SourceRef:  "draft://msg-123",
			Actor:      "synthetic-executor",
			CapturedAt: now,
			Checksum:   "synthetic-receipt-v1",
		}},
		OccurredAt: now,
		CreatedAt:  now,
	}
	audit, err := AuditOutcome(outcome, "audit-draft-created", now)
	if err != nil {
		t.Fatalf("audit outcome: %v", err)
	}
	if !audit.EvidenceSufficient || audit.Status == AuditEscalate || audit.OutcomeID != outcome.ID {
		t.Fatalf("expected evidence-backed non-escalation audit, got %+v", audit)
	}
}

func TestLinkMemoryToEventRejectsCrossPersonAndIsIdempotent(t *testing.T) {
	now := time.Date(2026, time.August, 19, 9, 0, 0, 0, time.UTC)
	temporal := validTemporal(now)
	memory := Memory{ID: "memory-1", PersonID: "person-a", Temporal: temporal}
	foreignEvent := Event{ID: "event-foreign", PersonID: "person-b", Temporal: temporal}
	if _, err := LinkMemoryToEvent(&memory, &foreignEvent, "invalid", now); !errors.Is(err, ErrPersonBoundary) {
		t.Fatalf("expected person boundary error, got %v", err)
	}

	event := Event{ID: "event-1", PersonID: "person-a", Temporal: temporal}
	if _, err := LinkMemoryToEvent(&memory, &event, "first", now); err != nil {
		t.Fatalf("first link: %v", err)
	}
	if _, err := LinkMemoryToEvent(&memory, &event, "second", now); err != nil {
		t.Fatalf("idempotent link: %v", err)
	}
	if len(memory.EventIDs) != 1 || len(event.MemoryIDs) != 1 {
		t.Fatalf("expected idempotent reciprocal references, got memory=%+v event=%+v", memory, event)
	}
}

func TestPendingIntentOnlyAllowsDeclaredStateTransitions(t *testing.T) {
	now := time.Date(2026, time.August, 19, 9, 0, 0, 0, time.UTC)
	intent := PendingIntent{ID: "intent-1", PersonID: "person-1", State: IntentCaptured, Temporal: validTemporal(now)}
	if err := TransitionPendingIntent(&intent, IntentCompleted, now, "person-1", "skip"); !errors.Is(err, ErrInvalidIntentTransition) {
		t.Fatalf("expected rejected state skip, got %v", err)
	}
	if err := TransitionPendingIntent(&intent, IntentReady, now, "person-1", "already clear"); err != nil {
		t.Fatalf("captured to ready should be allowed: %v", err)
	}
	if err := TransitionPendingIntent(&intent, IntentProposed, now, "person-1", "recommendation created"); err != nil {
		t.Fatalf("ready to proposed should be allowed: %v", err)
	}
	if err := TransitionPendingIntent(&intent, IntentCompleted, now, "person-1", "skip execution"); !errors.Is(err, ErrInvalidIntentTransition) {
		t.Fatalf("proposed to completed should be rejected, got %v", err)
	}
}

func TestFourDimensionalTimeAndTemporalUtility(t *testing.T) {
	now := time.Date(2026, time.August, 19, 9, 0, 0, 0, time.UTC)
	futureEventRecordedNow := TemporalState{
		EventAt:     now.Add(24 * time.Hour),
		RecordedAt:  now,
		EffectiveAt: now,
		AttentionAt: now,
	}
	if err := futureEventRecordedNow.Validate(); err != nil {
		t.Fatalf("future event recorded before its semantic event time must be valid, got %v", err)
	}

	expiredAt := now.Add(-time.Minute)
	expired := validTemporal(now)
	expired.ExpiresAt = &expiredAt
	evaluation, err := EvaluateTemporalUtility(expired, 0.8, now)
	if err != nil {
		t.Fatalf("evaluate expired temporal state: %v", err)
	}
	if evaluation.Active || evaluation.Utility != 0 {
		t.Fatalf("expired state must have zero current utility, got %+v", evaluation)
	}
}

func TestConfidenceLowersWhenReliableEvidenceContradictsClaim(t *testing.T) {
	now := time.Date(2026, time.August, 19, 9, 0, 0, 0, time.UTC)
	claim := Claim{ID: "claim-1", PersonID: "person-1"}
	provenance := Provenance{SourceType: "test", SourceRef: "test://source", CapturedAt: now}
	support := Evidence{ID: "support", PersonID: claim.PersonID, ClaimID: claim.ID, Stance: EvidenceSupports, Quality: 0.9, Relevance: 0.9, Provenance: provenance}
	contradiction := Evidence{ID: "contradiction", PersonID: claim.PersonID, ClaimID: claim.ID, Stance: EvidenceContradicts, Quality: 0.9, Relevance: 0.9, Provenance: provenance}
	positive, err := AssessClaimConfidence(&claim, []Evidence{support}, now)
	if err != nil {
		t.Fatalf("positive confidence: %v", err)
	}
	conflicted, err := AssessClaimConfidence(&claim, []Evidence{support, contradiction}, now)
	if err != nil {
		t.Fatalf("conflicted confidence: %v", err)
	}
	if conflicted.Score >= positive.Score || conflicted.Score != 0.5 || conflicted.ConflictingWeight == 0 {
		t.Fatalf("expected balanced contradiction to move score toward uncertainty, positive=%+v conflicted=%+v", positive, conflicted)
	}
}

func TestHardConstraintAndScopeEnforcementPreventAction(t *testing.T) {
	now := time.Date(2026, time.August, 19, 9, 0, 0, 0, time.UTC)
	temporal := validTemporal(now)
	goal := Goal{ID: "goal-1", PersonID: "person-1", Importance: 0.9, Status: GoalActive, Temporal: temporal}
	hardConstraint := Constraint{ID: "constraint-1", PersonID: "person-1", Kind: ConstraintHard, Title: "Do not contact during leave", Active: true, Temporal: temporal}
	opportunity := Opportunity{
		ID: "opportunity-1", PersonID: "person-1", GoalIDs: []string{goal.ID}, ConstraintIDs: []string{hardConstraint.ID}, Temporal: temporal,
		GoalAlignment: 0.9, ExpectedValue: 0.9, Effort: 0.1, Risk: 0.1, EvidenceConfidence: 0.9, TemporalPriority: 0.9,
	}
	evaluation, err := EvaluateOpportunity(&opportunity, []Goal{goal}, []Constraint{hardConstraint}, now)
	if err != nil {
		t.Fatalf("evaluate hard constraint: %v", err)
	}
	if !evaluation.HardBlocked || evaluation.Utility != 0 {
		t.Fatalf("hard constraint must block opportunity, got %+v", evaluation)
	}

	proposal := ActionProposal{ID: "proposal-1", PersonID: "person-1", PermissionID: "permission-1", Requested: Scope{Capability: "send_email", Resource: "email:client"}}
	permission := Permission{ID: "permission-1", PersonID: "person-1", Scopes: []Scope{{Capability: "create_draft", Resource: "email:client"}}, Temporal: temporal}
	if err := ValidateProposalPermission(proposal, permission, now); !errors.Is(err, ErrPermissionDenied) {
		t.Fatalf("expected disallowed capability to fail scope validation, got %v", err)
	}
}

func TestOutcomeRequiresEvidenceAndProvenance(t *testing.T) {
	now := time.Date(2026, time.August, 19, 9, 0, 0, 0, time.UTC)
	outcome := Outcome{ID: "outcome-1", PersonID: "person-1", ExpectedUtility: 0.7, ObservedUtility: 0.7}
	if _, err := AuditOutcome(outcome, "audit-1", now); !errors.Is(err, ErrMissingOutcomeEvidence) {
		t.Fatalf("expected missing outcome evidence rejection, got %v", err)
	}
}

func validTemporal(now time.Time) TemporalState {
	expiresAt := now.Add(24 * time.Hour)
	return TemporalState{
		EventAt:     now.Add(-time.Hour),
		RecordedAt:  now.Add(-30 * time.Minute),
		EffectiveAt: now.Add(-time.Hour),
		AttentionAt: now.Add(-time.Minute),
		ExpiresAt:   &expiresAt,
	}
}
