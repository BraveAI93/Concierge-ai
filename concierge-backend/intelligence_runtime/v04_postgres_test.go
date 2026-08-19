package intelligence_runtime

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/BraveAI93/concierge-backend/kernel"
)

func TestV04SharedVerticalSliceContractForMemoryAndPostgres(t *testing.T) {
	now := time.Date(2026, time.August, 21, 9, 0, 0, 0, time.UTC)
	t.Run("in_memory_reference", func(t *testing.T) {
		service, _, principal := runtimeFixture(t, now)
		exerciseVerticalSliceContract(t, service, principal, schedulingSource(now, "profile-internal-a", "conversation-shared-memory", "message-shared-memory"))
	})
	t.Run("disposable_postgres", func(t *testing.T) {
		service, _, principal := postgresRuntimeFixture(t, now)
		exerciseVerticalSliceContract(t, service, principal, schedulingSource(now, "profile-internal-a", "conversation-shared-postgres", "message-shared-postgres"))
	})
}

func exerciseVerticalSliceContract(t *testing.T, service RuntimeService, principal AuthenticatedPrincipal, source ConversationMessage) {
	t.Helper()
	first, err := service.IngestConversationMessage(context.Background(), principal, source)
	if err != nil {
		t.Fatalf("first vertical-slice ingestion: %v", err)
	}
	if first.Replayed || first.EventID == "" || first.MemoryID == "" || first.OpenLoopID == "" || first.ActionGateID == "" {
		t.Fatalf("expected canonical lineage, got %+v", first)
	}
	state, err := service.State(context.Background(), principal)
	if err != nil {
		t.Fatalf("state retrieval: %v", err)
	}
	if len(state.Sources) != 1 || len(state.Events) != 1 || len(state.Evidence) != 1 || len(state.Memories) != 1 || len(state.Claims) != 1 || len(state.OpenLoops) != 1 || len(state.Allocations) != 1 || len(state.Opportunities) != 1 || len(state.Decisions) != 1 || len(state.ActionProposals) != 1 || len(state.ActionGates) != 1 {
		t.Fatalf("unexpected persisted state: %+v", state)
	}
	if state.ActionGates[0].State != kernel.GateAwaitingApproval {
		t.Fatalf("expected approval gate, got %+v", state.ActionGates[0])
	}
	second, err := service.IngestConversationMessage(context.Background(), principal, source)
	if err != nil {
		t.Fatalf("idempotent replay: %v", err)
	}
	if !second.Replayed || second.EventID != first.EventID || second.ActionGateID != first.ActionGateID {
		t.Fatalf("expected deterministic replay: first=%+v second=%+v", first, second)
	}
	state, err = service.State(context.Background(), principal)
	if err != nil {
		t.Fatalf("state after replay: %v", err)
	}
	if len(state.Events) != 1 || len(state.OpenLoops) != 1 || len(state.ActionGates) != 1 || len(state.Replays) != 1 {
		t.Fatalf("replay duplicated state: %+v", state)
	}
}

func TestV04PostgresStagingRestartRetrievalAndSessionIdentity(t *testing.T) {
	now := time.Date(2026, time.August, 21, 9, 0, 0, 0, time.UTC)
	service, repo, principal := postgresRuntimeFixture(t, now)
	binding, err := (PostgresIdentityResolver{Repository: repo}).Resolve(context.Background(), principal)
	if err != nil || binding.Person.ID != "person-a" {
		t.Fatalf("stable subject must resolve person A: binding=%+v err=%v", binding, err)
	}
	sessions := NewStaticServerSessionSubjectLookup(map[string]string{"session-A-1": "owner-subject-a", "session-A-2": "owner-subject-a"})
	sessionAdapter := ServerSessionIdentityAdapter{Sessions: sessions, Bindings: PostgresIdentityResolver{Repository: repo}}
	firstBinding, err := sessionAdapter.ResolveSession(context.Background(), "session-A-1")
	if err != nil {
		t.Fatalf("resolve first stable session: %v", err)
	}
	secondBinding, err := sessionAdapter.ResolveSession(context.Background(), "session-A-2")
	if err != nil || secondBinding.Person.ID != firstBinding.Person.ID {
		t.Fatalf("sessions must map stably: first=%+v second=%+v err=%v", firstBinding, secondBinding, err)
	}
	if _, err := sessionAdapter.ResolveSession(context.Background(), "unknown-session"); !errors.Is(err, ErrUnknownIdentity) {
		t.Fatalf("unknown session must reject: %v", err)
	}

	source := schedulingSource(now, "profile-internal-a", "conversation-restart", "message-restart")
	first, err := service.IngestConversationMessage(context.Background(), principal, source)
	if err != nil {
		t.Fatalf("ingest before restart: %v", err)
	}
	if err := repo.Close(); err != nil {
		t.Fatalf("close first repository: %v", err)
	}
	secondRepo := openPostgresRepository(t)
	service.Repo = secondRepo
	service.Identity = PostgresIdentityResolver{Repository: secondRepo}
	state, err := service.State(context.Background(), principal)
	if err != nil {
		t.Fatalf("retrieve state after new repository instance: %v", err)
	}
	if len(state.Events) != 1 || len(state.Memories) != 1 || len(state.OpenLoops) != 1 || len(state.ActionGates) != 1 || state.Sources[0].MessageID != "message-restart" {
		t.Fatalf("durable restart state missing: %+v", state)
	}
	replay, err := service.IngestConversationMessage(context.Background(), principal, source)
	if err != nil || !replay.Replayed || replay.ActionGateID != first.ActionGateID {
		t.Fatalf("restart replay must be deterministic: replay=%+v err=%v", replay, err)
	}
}

func TestV04PostgresConcurrentIngestionProducesSingleCanonicalLineage(t *testing.T) {
	now := time.Date(2026, time.August, 21, 9, 0, 0, 0, time.UTC)
	service, _, principal := postgresRuntimeFixture(t, now)
	source := schedulingSource(now, "profile-internal-a", "conversation-concurrent", "message-concurrent")
	start := make(chan struct{})
	results := make(chan RuntimeResult, 2)
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			result, err := service.IngestConversationMessage(context.Background(), principal, source)
			results <- result
			errs <- err
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	close(errs)
	var outcomes []RuntimeResult
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent ingestion error: %v", err)
		}
	}
	for result := range results {
		outcomes = append(outcomes, result)
	}
	if len(outcomes) != 2 || outcomes[0].EventID == "" || outcomes[0].EventID != outcomes[1].EventID {
		t.Fatalf("concurrent results must share lineage: %+v", outcomes)
	}
	replays := 0
	for _, result := range outcomes {
		if result.Replayed {
			replays++
		}
	}
	if replays != 1 {
		t.Fatalf("exactly one concurrent call must return replay, got %+v", outcomes)
	}
	state, err := service.State(context.Background(), principal)
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Events) != 1 || len(state.Memories) != 1 || len(state.OpenLoops) != 1 || len(state.ActionGates) != 1 || len(state.Replays) != 1 {
		t.Fatalf("concurrency created duplicates: %+v", state)
	}
}

func TestV04PostgresFailureRollbackThenDeterministicReplay(t *testing.T) {
	now := time.Date(2026, time.August, 21, 9, 0, 0, 0, time.UTC)
	service, repo, principal := postgresRuntimeFixture(t, now)
	failure := errors.New("force transaction rollback")
	err := repo.RunInTransaction(context.Background(), "person-a", func(tx RuntimeTransaction) error {
		if err := tx.StoreSource(SourceRecord{ID: "failed-source", PersonID: "person-a", ProfileID: "profile-internal-a", ConversationID: "failed-conversation", MessageID: "failed-message", MessageRole: "user", Content: "ignored", MessageAt: now, StoredAt: now}); err != nil {
			return err
		}
		if err := tx.SaveEvent(kernel.Event{ID: "failed-event", PersonID: "person-a", Kind: "test", Summary: "must roll back", Temporal: testTemporal(now), CreatedAt: now}); err != nil {
			return err
		}
		return failure
	})
	if !errors.Is(err, failure) {
		t.Fatalf("expected forced rollback, got %v", err)
	}
	state, err := service.State(context.Background(), principal)
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Sources) != 0 || len(state.Events) != 0 || len(state.OpenLoops) != 0 || len(state.ActionGates) != 0 {
		t.Fatalf("partial state survived failed transaction: %+v", state)
	}
	source := schedulingSource(now, "profile-internal-a", "conversation-after-failure", "message-after-failure")
	first, err := service.IngestConversationMessage(context.Background(), principal, source)
	if err != nil {
		t.Fatalf("ingest after failure: %v", err)
	}
	replay, err := service.IngestConversationMessage(context.Background(), principal, source)
	if err != nil || !replay.Replayed || replay.EventID != first.EventID {
		t.Fatalf("replay after failure must remain deterministic: %+v %v", replay, err)
	}
}

func TestV04PostgresCrossPersonAndProfileBoundary(t *testing.T) {
	now := time.Date(2026, time.August, 21, 9, 0, 0, 0, time.UTC)
	service, repo, principalA := postgresRuntimeFixture(t, now)
	bindingB := testBinding("owner-subject-b", "profile-internal-b", "person-b", now)
	if err := repo.SeedBinding(context.Background(), bindingB); err != nil {
		t.Fatalf("seed person B: %v", err)
	}
	if _, err := repo.ReadState(context.Background(), "person-b", "person-a"); !errors.Is(err, ErrCrossPersonAccess) {
		t.Fatalf("cross-person read must fail: %v", err)
	}
	foreign := kernel.Event{ID: "foreign-event", PersonID: "person-a", Kind: "test", Summary: "foreign", Temporal: testTemporal(now), CreatedAt: now}
	if err := repo.RunInTransaction(context.Background(), "person-b", func(tx RuntimeTransaction) error { return tx.SaveEvent(foreign) }); !errors.Is(err, ErrCrossPersonAccess) {
		t.Fatalf("cross-person write must fail: %v", err)
	}
	if _, err := service.IngestConversationMessage(context.Background(), principalA, schedulingSource(now, "profile-internal-b", "conversation-redirect", "message-redirect")); !errors.Is(err, ErrSourceUnauthorized) {
		t.Fatalf("caller-selected profile redirect must fail: %v", err)
	}
	// Exercise database RLS directly under the internal role with person B's
	// transaction context. Even though the test connection owns the schema, the
	// forced policy must hide person A's runtime records from person B.
	if _, err := service.IngestConversationMessage(context.Background(), principalA, schedulingSource(now, "profile-internal-a", "conversation-rls", "message-rls")); err != nil {
		t.Fatalf("seed person A RLS record: %v", err)
	}
	rlsTx, err := repo.db.Begin()
	if err != nil {
		t.Fatalf("begin RLS adversarial read: %v", err)
	}
	defer rlsTx.Rollback()
	if _, err := rlsTx.Exec("SET LOCAL ROLE ci_kernel_runtime"); err != nil {
		t.Fatalf("set runtime role: %v", err)
	}
	if _, err := rlsTx.Exec("SELECT set_config('ci_kernel.person_id', 'person-b', true)"); err != nil {
		t.Fatalf("set person B context: %v", err)
	}
	if _, err := rlsTx.Exec("SELECT set_config('ci_kernel.stable_subject', 'owner-subject-b', true)"); err != nil {
		t.Fatalf("set person B subject: %v", err)
	}
	var visible int
	if err := rlsTx.QueryRow("SELECT count(*) FROM ci_kernel_v04.records WHERE person_id='person-a'").Scan(&visible); err != nil || visible != 0 {
		t.Fatalf("RLS must hide person A records from person B, visible=%d err=%v", visible, err)
	}
	if err := repo.SeedBinding(context.Background(), PersonBinding{StableSubjectID: "owner-subject-linked", SourceProfileID: "profile-primary-linked", AllowedSourceProfileIDs: []string{"profile-linked-alt"}, Person: kernel.Person{ID: "person-linked", DisplayName: "Linked", WorldID: "world-linked", CreatedAt: now}, World: kernel.PersonalWorld{ID: "world-linked", PersonID: "person-linked", UpdatedAt: now}}); err != nil {
		t.Fatalf("seed linked profile binding: %v", err)
	}
	linked, err := (PostgresIdentityResolver{Repository: repo}).Resolve(context.Background(), AuthenticatedPrincipal{StableSubjectID: "owner-subject-linked"})
	if err != nil || !linked.AllowsSourceProfile("profile-linked-alt") {
		t.Fatalf("explicit linked internal profile must resolve: binding=%+v err=%v", linked, err)
	}
}

func TestV04ConsentAndActivationFailClosed(t *testing.T) {
	now := time.Date(2026, time.August, 21, 9, 0, 0, 0, time.UTC)
	service, repo, principal := postgresRuntimeFixture(t, now)
	source := schedulingSource(now, "profile-internal-a", "conversation-consent", "message-consent")
	service.Consent = FailClosedConsentVerifier{}
	if _, err := service.IngestConversationMessage(context.Background(), principal, source); !errors.Is(err, ErrConsentNotVerified) {
		t.Fatalf("unverified consent must fail closed: %v", err)
	}
	state, err := repo.ReadState(context.Background(), "person-a", "person-a")
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Events) != 0 {
		t.Fatalf("consent failure must not persist memory lineage: %+v", state)
	}
	service.Consent = NewStaticConsentVerifier([]PersonBinding{testBinding("owner-subject-a", "profile-internal-a", "person-a", now)})
	service.Activation = StaticRuntimeActivation{Allowed: false}
	if _, err := service.IngestConversationMessage(context.Background(), principal, source); !errors.Is(err, ErrRuntimeDisabled) {
		t.Fatalf("kill switch must disable runtime: %v", err)
	}
}

func postgresRuntimeFixture(t *testing.T, now time.Time) (RuntimeService, *PostgresRuntimeRepository, AuthenticatedPrincipal) {
	t.Helper()
	repo := openPostgresRepository(t)
	resetPostgresSchema(t, repo)
	binding := testBinding("owner-subject-a", "profile-internal-a", "person-a", now)
	if err := repo.SeedBinding(context.Background(), binding); err != nil {
		t.Fatalf("seed binding: %v", err)
	}
	goal := kernel.Goal{ID: "goal-relationship", PersonID: "person-a", Title: "Maintain responsive client relationships", SubjectiveImportance: 0.6, Status: kernel.GoalActive, Temporal: testTemporal(now), CreatedAt: now}
	if err := repo.SeedGoal(context.Background(), goal); err != nil {
		t.Fatalf("seed goal: %v", err)
	}
	permission := kernel.Permission{ID: "permission-draft", PersonID: "person-a", Scopes: []kernel.Scope{{Capability: "create_draft", Resource: "email:client"}}, RequiresApproval: true, CanAutoApprove: false, Temporal: testTemporal(now), GrantedBy: "person-a", CreatedAt: now}
	if err := repo.SeedPermission(context.Background(), permission); err != nil {
		t.Fatalf("seed permission: %v", err)
	}
	budget := kernel.AttentionBudget{ID: "budget-morning", PersonID: "person-a", WindowStart: now.Add(-time.Hour), WindowEnd: now.Add(time.Hour), AttentionCapacity: 30 * time.Minute, MaxCompetingItems: 1, CurrentContext: kernel.AttentionContext{ContextIDs: []string{"conversation:conversation-shared-postgres", "conversation:conversation-restart", "conversation:conversation-concurrent", "conversation:conversation-after-failure"}}}
	if err := repo.SeedAttentionBudget(context.Background(), budget); err != nil {
		t.Fatalf("seed attention budget: %v", err)
	}
	service := RuntimeService{Feature: EnabledFeature(), Activation: StaticRuntimeActivation{Allowed: true}, Consent: NewStaticConsentVerifier([]PersonBinding{binding}), Identity: PostgresIdentityResolver{Repository: repo}, Adapter: ConservativeConversationAdapter{}, Repo: repo, Clock: fixedClock{now: now}, Policy: kernel.DefaultV02Policy(), Config: DefaultRuntimeConfig()}
	t.Cleanup(func() { repo.Close() })
	return service, repo, AuthenticatedPrincipal{StableSubjectID: "owner-subject-a"}
}

func openPostgresRepository(t *testing.T) *PostgresRuntimeRepository {
	t.Helper()
	dsn := os.Getenv("CI_V04_TEST_DSN")
	if dsn == "" {
		t.Skip("CI_V04_TEST_DSN is not configured for disposable local PostgreSQL")
	}
	repo, err := OpenPostgresRuntimeRepository(context.Background(), dsn)
	if err != nil {
		t.Fatalf("open local Postgres repository: %v", err)
	}
	return repo
}

func resetPostgresSchema(t *testing.T, repo *PostgresRuntimeRepository) {
	t.Helper()
	migration, err := os.ReadFile("../db/proposed_migrations/ci_kernel_v04/001_staging_up.sql")
	if err != nil {
		t.Fatalf("read staging migration: %v", err)
	}
	if _, err := repo.db.Exec("DROP SCHEMA IF EXISTS ci_kernel_v04 CASCADE"); err != nil {
		t.Fatalf("drop disposable schema: %v", err)
	}
	if _, err := repo.db.Exec(string(migration)); err != nil {
		t.Fatalf("apply disposable staging schema: %v", err)
	}
	var dbName string
	if err := repo.db.QueryRow("SELECT current_database()").Scan(&dbName); err != nil || dbName != "ci_kernel_v04_test" {
		t.Fatalf("expected disposable target ci_kernel_v04_test, name=%q err=%v", dbName, err)
	}
}

func TestV04MigrationContainsNoPublicRuntimePolicy(t *testing.T) {
	migration, err := os.ReadFile("../db/proposed_migrations/ci_kernel_v04/001_staging_up.sql")
	if err != nil {
		t.Fatal(err)
	}
	executable := make([]string, 0)
	for _, line := range strings.Split(string(migration), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "--") {
			executable = append(executable, strings.ToLower(line))
		}
	}
	text := strings.Join(executable, "\n")
	if strings.Contains(text, "to anon") || strings.Contains(text, "to public") || strings.Contains(text, "to authenticated") {
		t.Fatalf("migration must not grant broad public runtime policy")
	}
}

func TestV04PostgresPreservesClaimLineageAndCurrentHistoricalSelection(t *testing.T) {
	now := time.Date(2026, time.August, 21, 9, 0, 0, 0, time.UTC)
	service, repo, principal := postgresRuntimeFixture(t, now)
	older := kernel.Claim{
		ID:        "claim-historical",
		PersonID:  "person-a",
		Statement: "The client prefers a morning appointment.",
		Temporal:  testTemporal(now),
		Freshness: kernel.FreshnessState{LastValidatedAt: now, Status: kernel.FreshnessSuperseded},
		Lineage:   kernel.ClaimLineage{PersonID: "person-a", ClaimID: "claim-historical", PreservesHistory: true, RecordedAt: now},
		CreatedAt: now,
	}
	current := kernel.Claim{
		ID:           "claim-current",
		PersonID:     "person-a",
		Statement:    "The client now prefers an afternoon appointment.",
		SupersedesID: older.ID,
		Temporal:     testTemporal(now.Add(time.Minute)),
		Freshness:    kernel.FreshnessState{LastValidatedAt: now, Status: kernel.FreshnessFresh},
		Lineage:      kernel.ClaimLineage{PersonID: "person-a", ClaimID: "claim-current", SupersedesClaimID: older.ID, EvidenceIDs: []string{"evidence-lineage"}, PreservesHistory: true, RecordedAt: now},
		CreatedAt:    now,
	}
	if err := repo.SaveClaim(context.Background(), older); err != nil {
		t.Fatalf("persist historical claim: %v", err)
	}
	if err := repo.SaveClaim(context.Background(), current); err != nil {
		t.Fatalf("persist current claim: %v", err)
	}
	if err := repo.SaveClaimLineage(context.Background(), current.Lineage); err != nil {
		t.Fatalf("persist claim lineage: %v", err)
	}
	state, err := service.State(context.Background(), principal)
	if err != nil {
		t.Fatalf("retrieve claims: %v", err)
	}
	claims := make([]kernel.Claim, 0, 2)
	for _, claim := range state.Claims {
		if claim.ID == older.ID || claim.ID == current.ID {
			claims = append(claims, claim)
		}
	}
	selection, err := kernel.SelectCurrentClaim(claims, kernel.EvaluationMoment{WallClockAt: now})
	if err != nil || selection.CurrentClaimID != current.ID || !containsID(selection.HistoricalClaimIDs, older.ID) {
		t.Fatalf("expected durable current/historical selection, selection=%+v err=%v", selection, err)
	}
	tx, err := repo.begin(context.Background(), "person-a", "")
	if err != nil {
		t.Fatalf("open scoped lineage read: %v", err)
	}
	defer tx.Rollback()
	var count int
	if err := tx.QueryRow("SELECT count(*) FROM ci_kernel_v04.records WHERE person_id='person-a' AND record_kind='claim_lineage' AND id='claim-current'").Scan(&count); err != nil || count != 1 {
		t.Fatalf("expected persisted lineage record, count=%d err=%v", count, err)
	}
}

func TestV04PostgresRepositoryRejectsUnsafeTargetBeforeConnecting(t *testing.T) {
	_, err := OpenPostgresRuntimeRepository(context.Background(), "postgresql://example.invalid/production")
	if !errors.Is(err, ErrUnsafePersistenceTarget) {
		t.Fatalf("expected unsafe target rejection before connection, got %v", err)
	}
}
