package intelligence_runtime

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/BraveAI93/concierge-backend/db"
	"github.com/BraveAI93/concierge-backend/kernel"
)

func TestRuntimeVerticalSliceConversationToApprovalGateAndIdempotentReplay(t *testing.T) {
	now := time.Date(2026, time.August, 20, 9, 0, 0, 0, time.UTC)
	service, repo, principal := runtimeFixture(t, now)
	source := schedulingSource(now, "profile-internal-a", "conversation-1", "message-1")

	first, err := service.IngestConversationMessage(context.Background(), principal, source)
	if err != nil {
		t.Fatalf("run vertical slice: %v", err)
	}
	if first.Replayed || first.EventID == "" || first.OpenLoopID == "" || first.ActionGateID == "" {
		t.Fatalf("expected durable source-to-gate identifiers, got %+v", first)
	}

	state, err := service.State(context.Background(), principal)
	if err != nil {
		t.Fatalf("retrieve runtime state: %v", err)
	}
	if len(state.Sources) != 1 || len(state.Events) != 1 || len(state.Evidence) != 1 || len(state.Memories) != 1 || len(state.Claims) != 1 || len(state.PendingIntents) != 1 || len(state.OpenLoops) != 1 || len(state.Allocations) != 1 || len(state.Opportunities) != 1 || len(state.Decisions) != 1 || len(state.ActionProposals) != 1 || len(state.ActionGates) != 1 {
		t.Fatalf("expected complete persisted runtime lineage, got %+v", state)
	}
	if state.ActionGates[0].State != kernel.GateAwaitingApproval {
		t.Fatalf("action proposal must remain approval-gated, got %+v", state.ActionGates[0])
	}
	if state.Evidence[0].Provenance.SourceRef != "concierge://conversations/conversation-1/messages/message-1" || state.Sources[0].MessageID != "message-1" || state.Events[0].PersonID != first.PersonID {
		t.Fatalf("source provenance or person lineage was not preserved, state=%+v", state)
	}
	if !containsID(state.Events[0].MemoryIDs, state.Memories[0].ID) || !containsID(state.Memories[0].EventIDs, state.Events[0].ID) {
		t.Fatalf("expected reciprocal persisted memory/event link, event=%+v memory=%+v", state.Events[0], state.Memories[0])
	}

	second, err := service.IngestConversationMessage(context.Background(), principal, source)
	if err != nil {
		t.Fatalf("replay source: %v", err)
	}
	if !second.Replayed || second.EventID != first.EventID || second.ActionGateID != first.ActionGateID {
		t.Fatalf("expected deterministic idempotent replay, first=%+v second=%+v", first, second)
	}
	state, err = service.State(context.Background(), principal)
	if err != nil {
		t.Fatalf("retrieve replayed state: %v", err)
	}
	if len(state.Events) != 1 || len(state.OpenLoops) != 1 || len(state.ActionGates) != 1 || len(state.Replays) != 1 {
		t.Fatalf("idempotent replay must not duplicate canonical state, got %+v", state)
	}

	if err := repo.RunInTransaction(context.Background(), first.PersonID, func(tx RuntimeTransaction) error {
		if replay, ok := tx.FindReplay(first.IdempotencyKey); !ok || replay.EventID != first.EventID {
			t.Fatalf("replay record must be retrievable inside repository transaction")
		}
		return nil
	}); err != nil {
		t.Fatalf("replay repository contract: %v", err)
	}
}

func TestRuntimeIdentityResolutionAndCrossPersonIsolation(t *testing.T) {
	now := time.Date(2026, time.August, 20, 9, 0, 0, 0, time.UTC)
	service, repo, principalA := runtimeFixture(t, now)
	principalB := AuthenticatedPrincipal{StableSubjectID: "owner-subject-b"}
	if _, err := service.Identity.Resolve(context.Background(), principalA); err != nil {
		t.Fatalf("correct owner must resolve canonical person: %v", err)
	}
	if _, err := service.Identity.Resolve(context.Background(), AuthenticatedPrincipal{StableSubjectID: "unknown-owner"}); !errors.Is(err, ErrUnknownIdentity) {
		t.Fatalf("unknown identity must be rejected, got %v", err)
	}

	wrongProfile := schedulingSource(now, "profile-internal-b", "conversation-wrong", "message-wrong")
	if _, err := service.IngestConversationMessage(context.Background(), principalA, wrongProfile); !errors.Is(err, ErrSourceUnauthorized) {
		t.Fatalf("public caller must not redirect owner A to owner B profile/world, got %v", err)
	}
	stateA, err := service.State(context.Background(), principalA)
	if err != nil || len(stateA.Events) != 0 {
		t.Fatalf("rejected profile redirect must not write owner A state, state=%+v err=%v", stateA, err)
	}

	// Seed the second person only after proving the requested source cannot use
	// their profile via owner A's authenticated principal.
	bindingB := testBinding("owner-subject-b", "profile-internal-b", "person-b", now)
	if err := repo.SeedBinding(bindingB); err != nil {
		t.Fatalf("seed person B: %v", err)
	}
	if _, err := repo.ReadState(context.Background(), "person-b", "person-a"); !errors.Is(err, ErrCrossPersonAccess) {
		t.Fatalf("person B must not read person A world, got %v", err)
	}
	foreignEvent := kernel.Event{ID: "foreign-event", PersonID: "person-a", Temporal: testTemporal(now)}
	if err := repo.RunInTransaction(context.Background(), "person-b", func(tx RuntimeTransaction) error { return tx.SaveEvent(foreignEvent) }); !errors.Is(err, ErrCrossPersonAccess) {
		t.Fatalf("person B must not write person A event, got %v", err)
	}
	if _, err := service.State(context.Background(), principalB); err != nil {
		t.Fatalf("resolved person B should retrieve only their own empty world: %v", err)
	}
}

func TestRuntimeRepositoryTransactionRollsBackAndRejectsCrossPersonLinks(t *testing.T) {
	now := time.Date(2026, time.August, 20, 9, 0, 0, 0, time.UTC)
	_, repo, _ := runtimeFixture(t, now)
	rollback := errors.New("force rollback")
	err := repo.RunInTransaction(context.Background(), "person-a", func(tx RuntimeTransaction) error {
		if err := tx.StoreSource(SourceRecord{ID: "rollback-source", PersonID: "person-a", ProfileID: "profile-internal-a", MessageID: "rollback-message", MessageAt: now, StoredAt: now}); err != nil {
			return err
		}
		return rollback
	})
	if !errors.Is(err, rollback) {
		t.Fatalf("expected transaction failure, got %v", err)
	}
	state, err := repo.ReadState(context.Background(), "person-a", "person-a")
	if err != nil {
		t.Fatalf("read post-rollback state: %v", err)
	}
	if len(state.Sources) != 0 {
		t.Fatalf("transaction must discard writes on error, got %+v", state.Sources)
	}

	if err := repo.RunInTransaction(context.Background(), "person-a", func(tx RuntimeTransaction) error {
		return tx.SaveMemoryEventLink(kernel.MemoryEventLink{MemoryID: "missing-memory", EventID: "missing-event", LinkedAt: now})
	}); !errors.Is(err, ErrCrossPersonAccess) {
		t.Fatalf("link must reject missing or foreign person records, got %v", err)
	}
}

func TestConservativeAdapterDoesNotInventUnsupportedFacts(t *testing.T) {
	now := time.Date(2026, time.August, 20, 9, 0, 0, 0, time.UTC)
	binding := testBinding("owner-subject-a", "profile-internal-a", "person-a", now)
	unsupported := ConversationMessage{
		Conversation: db.Conversation{ID: "conversation-smalltalk", ProfileID: "profile-internal-a", StartedAt: now},
		Message:      db.Message{ID: "message-smalltalk", ConversationID: "conversation-smalltalk", Role: "user", Content: "Thanks for your help today.", CreatedAt: now},
	}
	if _, err := (ConservativeConversationAdapter{}).Map(binding, unsupported, now); !errors.Is(err, ErrUnsupportedSource) {
		t.Fatalf("unsupported source must not invent a claim or open loop, got %v", err)
	}
}

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

func runtimeFixture(t *testing.T, now time.Time) (RuntimeService, *InMemoryRuntimeRepository, AuthenticatedPrincipal) {
	t.Helper()
	bindingA := testBinding("owner-subject-a", "profile-internal-a", "person-a", now)
	resolver := NewStaticIdentityResolver([]PersonBinding{bindingA, testBinding("owner-subject-b", "profile-internal-b", "person-b", now)})
	repo := NewInMemoryRuntimeRepository()
	if err := repo.SeedBinding(bindingA); err != nil {
		t.Fatalf("seed person A: %v", err)
	}
	goal := kernel.Goal{ID: "goal-relationship", PersonID: "person-a", Title: "Maintain responsive client relationships", SubjectiveImportance: 0.6, Status: kernel.GoalActive, Temporal: testTemporal(now), CreatedAt: now}
	if err := repo.SeedGoal(goal); err != nil {
		t.Fatalf("seed goal: %v", err)
	}
	permission := kernel.Permission{
		ID:               "permission-draft",
		PersonID:         "person-a",
		Scopes:           []kernel.Scope{{Capability: "create_draft", Resource: "email:client"}},
		RequiresApproval: true,
		CanAutoApprove:   false,
		Temporal:         testTemporal(now),
		GrantedBy:        "person-a",
		CreatedAt:        now,
	}
	if err := repo.SeedPermission(permission); err != nil {
		t.Fatalf("seed permission: %v", err)
	}
	budget := kernel.AttentionBudget{
		ID:                "budget-morning",
		PersonID:          "person-a",
		WindowStart:       now.Add(-time.Hour),
		WindowEnd:         now.Add(time.Hour),
		AttentionCapacity: 30 * time.Minute,
		MaxCompetingItems: 1,
		CurrentContext:    kernel.AttentionContext{ContextIDs: []string{"conversation:conversation-1"}},
	}
	if err := repo.SeedAttentionBudget(budget); err != nil {
		t.Fatalf("seed attention budget: %v", err)
	}
	service := RuntimeService{
		Feature:    EnabledFeature(),
		Activation: StaticRuntimeActivation{Allowed: true},
		Consent:    NewStaticConsentVerifier([]PersonBinding{bindingA}),
		Identity:   resolver,
		Adapter:    ConservativeConversationAdapter{},
		Repo:       repo,
		Clock:      fixedClock{now: now},
		Policy:     kernel.DefaultV02Policy(),
		Config:     DefaultRuntimeConfig(),
		BoundaryPolicy: kernel.DeterministicBoundaryPolicy{Config: kernel.BoundaryPolicyConfig{MinimumGap: 30 * time.Minute, MinimumConfidence: 0.7}},
		ThreadResolver: kernel.DeterministicThreadResolver{Policy: kernel.ThreadResolverPolicy{SelectionThreshold: 0.5, AmbiguityMargin: 0.15}},
		RetrievalDepthPolicy: kernel.DeterministicRetrievalDepthPolicy{KeyContinuityThreshold: 0.6, ReconstructionThreshold: 0.8, DeepAuditThreshold: 0.9},
		AttunementPolicy: kernel.DeterministicAttunementSafetyPolicy{DefaultMaxChoices: 3},
	}
	return service, repo, AuthenticatedPrincipal{StableSubjectID: "owner-subject-a"}
}

func testBinding(subject, profileID, personID string, now time.Time) PersonBinding {
	person := kernel.Person{ID: personID, DisplayName: personID, WorldID: "world-" + personID, CreatedAt: now}
	return PersonBinding{StableSubjectID: subject, SourceProfileID: profileID, Person: person, World: kernel.PersonalWorld{ID: person.WorldID, PersonID: personID, UpdatedAt: now}}
}

func testTemporal(now time.Time) kernel.TemporalState {
	return kernel.TemporalState{EventAt: now.Add(-time.Minute), RecordedAt: now, EffectiveAt: now.Add(-time.Minute), AttentionAt: now}
}

func schedulingSource(now time.Time, profileID, conversationID, messageID string) ConversationMessage {
	return ConversationMessage{
		Conversation: db.Conversation{ID: conversationID, ProfileID: profileID, SessionID: "session-1", StartedAt: now.Add(-10 * time.Minute), MessageCount: 1},
		Message:      db.Message{ID: messageID, ConversationID: conversationID, Role: "user", Content: "Could you please reschedule my appointment? deadline: 2026-08-20T10:00:00Z", CreatedAt: now.Add(-5 * time.Minute)},
	}
}

func containsID(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func TestRuntimeFeatureIsDisabledByDefault(t *testing.T) {
	now := time.Date(2026, time.August, 20, 9, 0, 0, 0, time.UTC)
	service, repo, principal := runtimeFixture(t, now)
	service.Feature = DisabledFeature()
	if _, err := service.IngestConversationMessage(context.Background(), principal, schedulingSource(now, "profile-internal-a", "conversation-disabled", "message-disabled")); !errors.Is(err, ErrRuntimeDisabled) {
		t.Fatalf("disabled feature must reject runtime execution, got %v", err)
	}
	state, err := repo.ReadState(context.Background(), "person-a", "person-a")
	if err != nil || len(state.Events) != 0 || len(state.OpenLoops) != 0 {
		t.Fatalf("disabled feature must not persist state, state=%+v err=%v", state, err)
	}
}
