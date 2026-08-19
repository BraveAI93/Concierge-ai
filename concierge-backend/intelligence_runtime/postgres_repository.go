package intelligence_runtime

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/BraveAI93/concierge-backend/kernel"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const defaultRuntimeRole = "ci_kernel_runtime"

// PostgresRuntimeRepository is the staging-only, real PostgreSQL implementation
// of the runtime and kernel persistence contracts. It has no Supabase-specific
// code and never reads environment variables: callers must pass an explicit
// local/test DSN. Production composition is intentionally absent.
type PostgresRuntimeRepository struct {
	db          *sql.DB
	role        string
	beginTxHook func()
}

func OpenPostgresRuntimeRepository(ctx context.Context, dsn string) (*PostgresRuntimeRepository, error) {
	if dsn == "" {
		return nil, ErrInvalidRuntimeConfig
	}
	config, err := pgx.ParseConfig(dsn)
	if err != nil || !isApprovedLocalStagingTarget(config.Host, config.Database) {
		return nil, ErrUnsafePersistenceTarget
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	var databaseName string
	var socketOnly bool
	if err := db.QueryRowContext(ctx, "SELECT current_database(), inet_server_addr() IS NULL").Scan(&databaseName, &socketOnly); err != nil || databaseName != "ci_kernel_v04_test" || !socketOnly {
		db.Close()
		return nil, ErrUnsafePersistenceTarget
	}
	return &PostgresRuntimeRepository{db: db, role: defaultRuntimeRole}, nil
}

func isApprovedLocalStagingTarget(host, database string) bool {
	if database != "ci_kernel_v04_test" {
		return false
	}
	if strings.HasPrefix(host, "/") {
		return true
	}
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}

func (r *PostgresRuntimeRepository) Close() error {
	if r == nil || r.db == nil {
		return nil
	}
	return r.db.Close()
}

func (r *PostgresRuntimeRepository) begin(ctx context.Context, personID, stableSubject string) (*sql.Tx, error) {
	if r == nil || r.db == nil || personID == "" {
		return nil, ErrInvalidRuntimeConfig
	}
	if r.beginTxHook != nil {
		r.beginTxHook()
	}
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	if _, err = tx.ExecContext(ctx, "SET LOCAL ROLE "+r.role); err != nil {
		tx.Rollback()
		return nil, err
	}
	if _, err = tx.ExecContext(ctx, "SELECT set_config('ci_kernel.person_id', $1, true)", personID); err != nil {
		tx.Rollback()
		return nil, err
	}
	if _, err = tx.ExecContext(ctx, "SELECT set_config('ci_kernel.stable_subject', $1, true)", stableSubject); err != nil {
		tx.Rollback()
		return nil, err
	}
	return tx, nil
}

func setTxContext(ctx context.Context, tx *sql.Tx, personID, stableSubject string) error {
	if _, err := tx.ExecContext(ctx, "SELECT set_config('ci_kernel.person_id', $1, true)", personID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "SELECT set_config('ci_kernel.stable_subject', $1, true)", stableSubject); err != nil {
		return err
	}
	return nil
}

// RunInTransaction checks the canonical person inside the same database
// transaction before exposing a person-bound transaction to the service.
func (r *PostgresRuntimeRepository) RunInTransaction(ctx context.Context, personID string, fn func(RuntimeTransaction) error) error {
	if fn == nil {
		return ErrInvalidRuntimeConfig
	}
	tx, err := r.begin(ctx, personID, "")
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var exists bool
	if err := tx.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM ci_kernel_v04.people WHERE id=$1)", personID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return ErrUnknownIdentity
	}
	if err := fn(&postgresTransaction{ctx: ctx, tx: tx, personID: personID}); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *PostgresRuntimeRepository) withAutoTransaction(ctx context.Context, personID string, fn func(*postgresTransaction) error) error {
	return r.RunInTransaction(ctx, personID, func(tx RuntimeTransaction) error { return fn(tx.(*postgresTransaction)) })
}

type postgresTransaction struct {
	ctx      context.Context
	tx       *sql.Tx
	personID string
}

func (tx *postgresTransaction) ensurePerson(personID string) error {
	if personID == "" || personID != tx.personID {
		return ErrCrossPersonAccess
	}
	return nil
}

func jsonPayload(value any) ([]byte, error)           { return json.Marshal(value) }
func decodePayload(raw []byte, destination any) error { return json.Unmarshal(raw, destination) }

func isUniqueViolation(err error) bool {
	var sqlErr interface{ SQLState() string }
	if errors.As(err, &sqlErr) {
		return sqlErr.SQLState() == "23505"
	}
	return false
}

func (tx *postgresTransaction) insertRecord(kind, id, personID string, value any) error {
	if err := tx.ensurePerson(personID); err != nil {
		return err
	}
	payload, err := jsonPayload(value)
	if err != nil {
		return err
	}
	_, err = tx.tx.ExecContext(tx.ctx, "INSERT INTO ci_kernel_v04.records(record_kind,id,person_id,payload) VALUES($1,$2,$3,$4::jsonb)", kind, id, personID, string(payload))
	if isUniqueViolation(err) {
		return ErrDuplicateRuntimeRecord
	}
	return err
}

func (tx *postgresTransaction) recordExists(kind, id string) (bool, error) {
	var exists bool
	err := tx.tx.QueryRowContext(tx.ctx, "SELECT EXISTS (SELECT 1 FROM ci_kernel_v04.records WHERE record_kind=$1 AND id=$2 AND person_id=$3)", kind, id, tx.personID).Scan(&exists)
	return exists, err
}

func (tx *postgresTransaction) FindReplay(key string) (RuntimeResult, bool) {
	if key == "" {
		return RuntimeResult{}, false
	}
	// A transaction-scoped advisory lock serializes same-source ingestions even
	// before a replay row exists. It is released automatically on rollback/commit.
	if _, err := tx.tx.ExecContext(tx.ctx, "SELECT pg_advisory_xact_lock(hashtext($1))", key); err != nil {
		return RuntimeResult{}, false
	}
	var payload []byte
	err := tx.tx.QueryRowContext(tx.ctx, "SELECT payload FROM ci_kernel_v04.runtime_replays WHERE idempotency_key=$1 AND person_id=$2 FOR UPDATE", key, tx.personID).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return RuntimeResult{}, false
	}
	if err != nil || decodePayload(payload, &RuntimeResult{}) != nil {
		return RuntimeResult{}, false
	}
	var result RuntimeResult
	if err := decodePayload(payload, &result); err != nil {
		return RuntimeResult{}, false
	}
	return result, true
}

func (tx *postgresTransaction) StoreReplay(result RuntimeResult) error {
	if err := tx.ensurePerson(result.PersonID); err != nil {
		return err
	}
	payload, err := jsonPayload(result)
	if err != nil {
		return err
	}
	_, err = tx.tx.ExecContext(tx.ctx, "INSERT INTO ci_kernel_v04.runtime_replays(idempotency_key,person_id,payload) VALUES($1,$2,$3::jsonb)", result.IdempotencyKey, result.PersonID, string(payload))
	if isUniqueViolation(err) {
		return ErrDuplicateRuntimeRecord
	}
	return err
}
func (tx *postgresTransaction) StoreSource(source SourceRecord) error {
	if err := tx.ensurePerson(source.PersonID); err != nil {
		return err
	}
	payload, err := jsonPayload(source)
	if err != nil {
		return err
	}
	_, err = tx.tx.ExecContext(tx.ctx, `INSERT INTO ci_kernel_v04.runtime_sources
		(id,person_id,source_profile_id,conversation_id,session_id,message_id,message_role,content,conversation_at,message_at,stored_at,payload)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12::jsonb)`, source.ID, source.PersonID, source.ProfileID, source.ConversationID, source.SessionID, source.MessageID, source.MessageRole, source.Content, nullableTime(source.ConversationAt), source.MessageAt, source.StoredAt, string(payload))
	if isUniqueViolation(err) {
		return ErrDuplicateRuntimeRecord
	}
	return err
}

func nullableTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}
func (tx *postgresTransaction) SaveEvent(value kernel.Event) error {
	return tx.insertRecord("event", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveEvidence(value kernel.Evidence) error {
	return tx.insertRecord("evidence", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveMemory(value kernel.Memory) error {
	return tx.insertRecord("memory", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveClaim(value kernel.Claim) error {
	return tx.insertRecord("claim", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SavePendingIntent(value kernel.PendingIntent) error {
	return tx.insertRecord("pending_intent", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveOpenLoop(value kernel.OpenLoop) error {
	if err := tx.ensurePerson(value.PersonID); err != nil {
		return err
	}
	ok, err := tx.recordExists("pending_intent", value.PendingIntentID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCrossPersonAccess
	}
	return tx.insertRecord("open_loop", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveOpportunity(value kernel.Opportunity) error {
	return tx.insertRecord("opportunity", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveDecision(value kernel.Decision) error {
	if err := tx.ensurePerson(value.PersonID); err != nil {
		return err
	}
	ok, err := tx.recordExists("opportunity", value.OpportunityID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCrossPersonAccess
	}
	return tx.insertRecord("decision", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveActionProposal(value kernel.ActionProposal) error {
	if err := tx.ensurePerson(value.PersonID); err != nil {
		return err
	}
	for _, parent := range []struct{ kind, id string }{{"opportunity", value.OpportunityID}, {"decision", value.DecisionID}, {"permission", value.PermissionID}} {
		ok, err := tx.recordExists(parent.kind, parent.id)
		if err != nil {
			return err
		}
		if !ok {
			return ErrCrossPersonAccess
		}
	}
	return tx.insertRecord("action_proposal", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveActionGate(value kernel.ActionGate) error {
	if err := tx.ensurePerson(value.PersonID); err != nil {
		return err
	}
	ok, err := tx.recordExists("action_proposal", value.ActionProposalID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCrossPersonAccess
	}
	return tx.insertRecord("action_gate", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveGoal(value kernel.Goal) error {
	return tx.insertRecord("goal", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveConstraint(value kernel.Constraint) error {
	return tx.insertRecord("constraint", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SavePermission(value kernel.Permission) error {
	return tx.insertRecord("permission", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveMemoryEventLink(link kernel.MemoryEventLink) error {
	if err := tx.ensurePerson(link.PersonID); err != nil {
		return err
	}
	memoryOK, err := tx.recordExists("memory", link.MemoryID)
	if err != nil {
		return err
	}
	eventOK, err := tx.recordExists("event", link.EventID)
	if err != nil {
		return err
	}
	if !memoryOK || !eventOK {
		return ErrCrossPersonAccess
	}
	payload, err := jsonPayload(link)
	if err != nil {
		return err
	}
	_, err = tx.tx.ExecContext(tx.ctx, `INSERT INTO ci_kernel_v04.memory_event_links(person_id,memory_id,event_id,payload,linked_at)
		VALUES($1,$2,$3,$4::jsonb,$5)`, tx.personID, link.MemoryID, link.EventID, string(payload), link.LinkedAt)
	if isUniqueViolation(err) {
		return ErrDuplicateRuntimeRecord
	}
	return err
}
func (tx *postgresTransaction) SaveAttentionBudget(value kernel.AttentionBudget) error {
	if err := value.Validate(); err != nil {
		return err
	}
	return tx.insertRecord("attention_budget", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveAttentionAllocation(value kernel.AttentionAllocation) error {
	if err := tx.ensurePerson(value.PersonID); err != nil {
		return err
	}
	budgetOK, err := tx.recordExists("attention_budget", value.BudgetID)
	if err != nil {
		return err
	}
	if !budgetOK {
		return ErrCrossPersonAccess
	}
	payload, err := jsonPayload(value)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(payload)
	allocationID := "allocation:" + value.BudgetID + ":" + hex.EncodeToString(sum[:8])
	_, err = tx.tx.ExecContext(tx.ctx, "INSERT INTO ci_kernel_v04.attention_allocations(allocation_id,person_id,budget_id,payload,evaluated_at) VALUES($1,$2,$3,$4::jsonb,now())", allocationID, value.PersonID, value.BudgetID, string(payload))
	if isUniqueViolation(err) {
		return ErrDuplicateRuntimeRecord
	}
	return err
}

func (tx *postgresTransaction) ListActiveGoals(at time.Time) []kernel.Goal {
	rows, err := tx.tx.QueryContext(tx.ctx, "SELECT payload FROM ci_kernel_v04.records WHERE person_id=$1 AND record_kind='goal' ORDER BY id", tx.personID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	goals := make([]kernel.Goal, 0)
	for rows.Next() {
		var payload []byte
		if rows.Scan(&payload) == nil {
			var goal kernel.Goal
			if decodePayload(payload, &goal) == nil && goal.Status == kernel.GoalActive && goal.Temporal.IsActive(at) {
				goals = append(goals, goal)
			}
		}
	}
	return goals
}
func (tx *postgresTransaction) ListActivePermissions(at time.Time) []kernel.Permission {
	rows, err := tx.tx.QueryContext(tx.ctx, "SELECT payload FROM ci_kernel_v04.records WHERE person_id=$1 AND record_kind='permission' ORDER BY id", tx.personID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	permissions := make([]kernel.Permission, 0)
	for rows.Next() {
		var payload []byte
		if rows.Scan(&payload) == nil {
			var permission kernel.Permission
			if decodePayload(payload, &permission) == nil && permission.Temporal.IsActive(at) {
				permissions = append(permissions, permission)
			}
		}
	}
	return permissions
}
func (tx *postgresTransaction) CurrentAttentionBudget(at time.Time) (kernel.AttentionBudget, bool) {
	rows, err := tx.tx.QueryContext(tx.ctx, "SELECT payload FROM ci_kernel_v04.records WHERE person_id=$1 AND record_kind='attention_budget' ORDER BY id", tx.personID)
	if err != nil {
		return kernel.AttentionBudget{}, false
	}
	defer rows.Close()
	for rows.Next() {
		var payload []byte
		if rows.Scan(&payload) == nil {
			var budget kernel.AttentionBudget
			if decodePayload(payload, &budget) == nil && !at.Before(budget.WindowStart) && !at.After(budget.WindowEnd) {
				return budget, true
			}
		}
	}
	return kernel.AttentionBudget{}, false
}

// SeedBinding and companion seed methods are staging-only test setup helpers.
// They use the same RLS-protected role and are not wired to legacy persistence.
func (r *PostgresRuntimeRepository) SeedBinding(ctx context.Context, binding PersonBinding) error {
	if r == nil || binding.Person.ID == "" || binding.World.ID == "" || binding.World.PersonID != binding.Person.ID || binding.StableSubjectID == "" || binding.SourceProfileID == "" {
		return ErrInvalidRuntimeConfig
	}
	tx, err := r.begin(ctx, binding.Person.ID, binding.StableSubjectID)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	personPayload, err := jsonPayload(binding.Person)
	if err != nil {
		return err
	}
	worldPayload, err := jsonPayload(binding.World)
	if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO ci_kernel_v04.people(id,payload) VALUES($1,$2::jsonb)", binding.Person.ID, string(personPayload)); isUniqueViolation(err) {
		return ErrDuplicateRuntimeRecord
	} else if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO ci_kernel_v04.worlds(person_id,payload) VALUES($1,$2::jsonb)", binding.Person.ID, string(worldPayload)); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO ci_kernel_v04.person_binding_subjects(stable_subject,person_id) VALUES($1,$2)", binding.StableSubjectID, binding.Person.ID); err != nil {
		return err
	}
	profiles := append([]string{binding.SourceProfileID}, binding.AllowedSourceProfileIDs...)
	seen := make(map[string]bool)
	for index, profile := range profiles {
		if profile == "" || seen[profile] {
			continue
		}
		seen[profile] = true
		if _, err = tx.ExecContext(ctx, "INSERT INTO ci_kernel_v04.person_profile_links(person_id,source_profile_id,is_primary) VALUES($1,$2,$3)", binding.Person.ID, profile, index == 0); err != nil {
			return err
		}
	}
	return tx.Commit()
}
func (r *PostgresRuntimeRepository) SeedGoal(ctx context.Context, value kernel.Goal) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveGoal(value) })
}
func (r *PostgresRuntimeRepository) SeedPermission(ctx context.Context, value kernel.Permission) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SavePermission(value) })
}
func (r *PostgresRuntimeRepository) SeedAttentionBudget(ctx context.Context, value kernel.AttentionBudget) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveAttentionBudget(value) })
}

// ResolveBinding loads the binding and all explicitly linked internal source
// profiles under the stable-subject RLS context, then loads person/world under
// the resulting canonical-person RLS context.
func (r *PostgresRuntimeRepository) ResolveBinding(ctx context.Context, stableSubject string) (PersonBinding, error) {
	if r == nil || r.db == nil || stableSubject == "" {
		return PersonBinding{}, ErrUnknownIdentity
	}
	tx, err := r.begin(ctx, "identity-bootstrap", stableSubject)
	if err != nil {
		return PersonBinding{}, err
	}
	defer tx.Rollback()
	// The identity-bootstrap person is intentionally not an actual person; RLS
	// policy on bindings uses only stable subject for this first lookup.
	var personID string
	if err := tx.QueryRowContext(ctx, "SELECT person_id FROM ci_kernel_v04.person_binding_subjects WHERE stable_subject=$1", stableSubject).Scan(&personID); errors.Is(err, sql.ErrNoRows) {
		return PersonBinding{}, ErrUnknownIdentity
	} else if err != nil {
		return PersonBinding{}, err
	}
	if err := setTxContext(ctx, tx, personID, stableSubject); err != nil {
		return PersonBinding{}, err
	}
	var personRaw, worldRaw []byte
	if err := tx.QueryRowContext(ctx, "SELECT payload FROM ci_kernel_v04.people WHERE id=$1", personID).Scan(&personRaw); err != nil {
		return PersonBinding{}, err
	}
	if err := tx.QueryRowContext(ctx, "SELECT payload FROM ci_kernel_v04.worlds WHERE person_id=$1", personID).Scan(&worldRaw); err != nil {
		return PersonBinding{}, err
	}
	var person kernel.Person
	var world kernel.PersonalWorld
	if err := decodePayload(personRaw, &person); err != nil {
		return PersonBinding{}, err
	}
	if err := decodePayload(worldRaw, &world); err != nil {
		return PersonBinding{}, err
	}
	rows, err := tx.QueryContext(ctx, "SELECT source_profile_id,is_primary FROM ci_kernel_v04.person_profile_links WHERE person_id=$1 ORDER BY is_primary DESC,source_profile_id", personID)
	if err != nil {
		return PersonBinding{}, err
	}
	defer rows.Close()
	profiles := make([]string, 0)
	var primary string
	for rows.Next() {
		var profile string
		var isPrimary bool
		if err := rows.Scan(&profile, &isPrimary); err != nil {
			return PersonBinding{}, err
		}
		if isPrimary {
			primary = profile
		} else {
			profiles = append(profiles, profile)
		}
	}
	if err := rows.Err(); err != nil {
		return PersonBinding{}, err
	}
	if primary == "" {
		return PersonBinding{}, ErrUnknownIdentity
	}
	if err := tx.Commit(); err != nil {
		return PersonBinding{}, err
	}
	return PersonBinding{StableSubjectID: stableSubject, SourceProfileID: primary, AllowedSourceProfileIDs: profiles, Person: person, World: world}, nil
}

func (r *PostgresRuntimeRepository) ReadState(ctx context.Context, requesterPersonID, targetPersonID string) (RuntimeState, error) {
	if requesterPersonID == "" || requesterPersonID != targetPersonID {
		return RuntimeState{}, ErrCrossPersonAccess
	}
	tx, err := r.begin(ctx, requesterPersonID, "")
	if err != nil {
		return RuntimeState{}, err
	}
	defer tx.Rollback()
	state, err := readPostgresState(ctx, tx, requesterPersonID)
	if err != nil {
		return RuntimeState{}, err
	}
	if err := tx.Commit(); err != nil {
		return RuntimeState{}, err
	}
	return state, nil
}

func readPostgresState(ctx context.Context, tx *sql.Tx, personID string) (RuntimeState, error) {
	var personRaw, worldRaw []byte
	if err := tx.QueryRowContext(ctx, "SELECT payload FROM ci_kernel_v04.people WHERE id=$1", personID).Scan(&personRaw); errors.Is(err, sql.ErrNoRows) {
		return RuntimeState{}, ErrUnknownIdentity
	} else if err != nil {
		return RuntimeState{}, err
	}
	if err := tx.QueryRowContext(ctx, "SELECT payload FROM ci_kernel_v04.worlds WHERE person_id=$1", personID).Scan(&worldRaw); err != nil {
		return RuntimeState{}, err
	}
	var state RuntimeState
	if err := decodePayload(personRaw, &state.Person); err != nil {
		return RuntimeState{}, err
	}
	if err := decodePayload(worldRaw, &state.World); err != nil {
		return RuntimeState{}, err
	}
	rows, err := tx.QueryContext(ctx, "SELECT record_kind,payload FROM ci_kernel_v04.records WHERE person_id=$1 ORDER BY record_kind,id", personID)
	if err != nil {
		return RuntimeState{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var kind string
		var payload []byte
		if err := rows.Scan(&kind, &payload); err != nil {
			return RuntimeState{}, err
		}
		if err := appendRecord(&state, kind, payload); err != nil {
			return RuntimeState{}, err
		}
	}
	if err := rows.Err(); err != nil {
		return RuntimeState{}, err
	}
	rows, err = tx.QueryContext(ctx, "SELECT payload FROM ci_kernel_v04.runtime_sources WHERE person_id=$1 ORDER BY id", personID)
	if err != nil {
		return RuntimeState{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var payload []byte
		if err := rows.Scan(&payload); err != nil {
			return RuntimeState{}, err
		}
		var source SourceRecord
		if err := decodePayload(payload, &source); err != nil {
			return RuntimeState{}, err
		}
		state.Sources = append(state.Sources, source)
	}
	rows.Close()
	rows, err = tx.QueryContext(ctx, "SELECT payload FROM ci_kernel_v04.attention_allocations WHERE person_id=$1 ORDER BY allocation_id", personID)
	if err != nil {
		return RuntimeState{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var payload []byte
		if err := rows.Scan(&payload); err != nil {
			return RuntimeState{}, err
		}
		var item kernel.AttentionAllocation
		if err := decodePayload(payload, &item); err != nil {
			return RuntimeState{}, err
		}
		state.Allocations = append(state.Allocations, item)
	}
	rows.Close()
	rows, err = tx.QueryContext(ctx, "SELECT payload FROM ci_kernel_v04.runtime_replays WHERE person_id=$1 ORDER BY idempotency_key", personID)
	if err != nil {
		return RuntimeState{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var payload []byte
		if err := rows.Scan(&payload); err != nil {
			return RuntimeState{}, err
		}
		var replay RuntimeResult
		if err := decodePayload(payload, &replay); err != nil {
			return RuntimeState{}, err
		}
		state.Replays = append(state.Replays, replay)
	}
	sortRuntimeState(&state)
	return state, rows.Err()
}

func appendRecord(state *RuntimeState, kind string, payload []byte) error {
	switch kind {
	case "event":
		var value kernel.Event
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.Events = append(state.Events, value)
	case "evidence":
		var value kernel.Evidence
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.Evidence = append(state.Evidence, value)
	case "memory":
		var value kernel.Memory
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.Memories = append(state.Memories, value)
	case "claim":
		var value kernel.Claim
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.Claims = append(state.Claims, value)
	case "goal":
		var value kernel.Goal
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.Goals = append(state.Goals, value)
	case "constraint":
		var value kernel.Constraint
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.Constraints = append(state.Constraints, value)
	case "pending_intent":
		var value kernel.PendingIntent
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.PendingIntents = append(state.PendingIntents, value)
	case "open_loop":
		var value kernel.OpenLoop
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.OpenLoops = append(state.OpenLoops, value)
	case "attention_budget":
		var value kernel.AttentionBudget
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.AttentionBudgets = append(state.AttentionBudgets, value)
	case "opportunity":
		var value kernel.Opportunity
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.Opportunities = append(state.Opportunities, value)
	case "decision":
		var value kernel.Decision
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.Decisions = append(state.Decisions, value)
	case "permission":
		var value kernel.Permission
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.Permissions = append(state.Permissions, value)
	case "action_proposal":
		var value kernel.ActionProposal
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.ActionProposals = append(state.ActionProposals, value)
	case "action_gate":
		var value kernel.ActionGate
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.ActionGates = append(state.ActionGates, value)
	case "outcome":
		var value kernel.Outcome
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.Outcomes = append(state.Outcomes, value)
	case "self_audit":
		var value kernel.SelfAudit
		if err := decodePayload(payload, &value); err != nil {
			return err
		}
		state.Audits = append(state.Audits, value)
	}
	return nil
}

// kernel.KernelRepository compatibility wrappers. They preserve the same
// person-bound transaction model used by runtime writes.
func (r *PostgresRuntimeRepository) FindPerson(ctx context.Context, personID string) (kernel.Person, error) {
	state, err := r.ReadState(ctx, personID, personID)
	return state.Person, err
}
func (r *PostgresRuntimeRepository) FindWorld(ctx context.Context, personID string) (kernel.PersonalWorld, error) {
	state, err := r.ReadState(ctx, personID, personID)
	return state.World, err
}
func (r *PostgresRuntimeRepository) SaveWorld(ctx context.Context, value kernel.PersonalWorld) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error {
		payload, err := jsonPayload(value)
		if err != nil {
			return err
		}
		_, err = tx.tx.ExecContext(tx.ctx, "UPDATE ci_kernel_v04.worlds SET payload=$1::jsonb,updated_at=now() WHERE person_id=$2", string(payload), value.PersonID)
		return err
	})
}
func (r *PostgresRuntimeRepository) SaveEntity(ctx context.Context, value kernel.Entity) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.insertRecord("entity", value.ID, value.PersonID, value) })
}
func (r *PostgresRuntimeRepository) SaveContext(ctx context.Context, value kernel.Context) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error {
		return tx.insertRecord("context", value.ID, value.PersonID, value)
	})
}
func (r *PostgresRuntimeRepository) SaveEvent(ctx context.Context, value kernel.Event) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveEvent(value) })
}
func (r *PostgresRuntimeRepository) SaveMemory(ctx context.Context, value kernel.Memory) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveMemory(value) })
}
func (r *PostgresRuntimeRepository) SaveClaim(ctx context.Context, value kernel.Claim) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveClaim(value) })
}
func (r *PostgresRuntimeRepository) SaveEvidence(ctx context.Context, value kernel.Evidence) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveEvidence(value) })
}

func (r *PostgresRuntimeRepository) SaveMemoryEventLink(ctx context.Context, link kernel.MemoryEventLink) error {
	if link.PersonID == "" {
		return ErrCrossPersonAccess
	}
	return r.withAutoTransaction(ctx, link.PersonID, func(tx *postgresTransaction) error { return tx.SaveMemoryEventLink(link) })
}
func (r *PostgresRuntimeRepository) SaveGoal(ctx context.Context, value kernel.Goal) error {
	return r.SeedGoal(ctx, value)
}
func (r *PostgresRuntimeRepository) SaveConstraint(ctx context.Context, value kernel.Constraint) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveConstraint(value) })
}
func (r *PostgresRuntimeRepository) SavePendingIntent(ctx context.Context, value kernel.PendingIntent) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SavePendingIntent(value) })
}
func (r *PostgresRuntimeRepository) SaveOpenLoop(ctx context.Context, value kernel.OpenLoop) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveOpenLoop(value) })
}
func (r *PostgresRuntimeRepository) SaveOpportunity(ctx context.Context, value kernel.Opportunity) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveOpportunity(value) })
}
func (r *PostgresRuntimeRepository) SaveDecision(ctx context.Context, value kernel.Decision) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveDecision(value) })
}
func (r *PostgresRuntimeRepository) SavePermission(ctx context.Context, value kernel.Permission) error {
	return r.SeedPermission(ctx, value)
}
func (r *PostgresRuntimeRepository) SaveActionProposal(ctx context.Context, value kernel.ActionProposal) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveActionProposal(value) })
}
func (r *PostgresRuntimeRepository) SaveActionGate(ctx context.Context, value kernel.ActionGate) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveActionGate(value) })
}
func (r *PostgresRuntimeRepository) SaveOutcome(ctx context.Context, value kernel.Outcome) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error {
		return tx.insertRecord("outcome", value.ID, value.PersonID, value)
	})
}
func (r *PostgresRuntimeRepository) SaveSelfAudit(ctx context.Context, value kernel.SelfAudit) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error {
		return tx.insertRecord("self_audit", value.ID, value.PersonID, value)
	})
}
func (r *PostgresRuntimeRepository) SaveAttentionBudget(ctx context.Context, value kernel.AttentionBudget) error {
	return r.SeedAttentionBudget(ctx, value)
}
func (r *PostgresRuntimeRepository) SaveAttentionAllocation(ctx context.Context, value kernel.AttentionAllocation) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveAttentionAllocation(value) })
}

func (r *PostgresRuntimeRepository) SaveClaimLineage(ctx context.Context, value kernel.ClaimLineage) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error {
		claimOK, err := tx.recordExists("claim", value.ClaimID)
		if err != nil {
			return err
		}
		if !claimOK {
			return ErrCrossPersonAccess
		}
		return tx.insertRecord("claim_lineage", value.ClaimID, value.PersonID, value)
	})
}

// Explicitly retain interface checks for the staging adapter.
var _ RuntimeRepository = (*PostgresRuntimeRepository)(nil)
var _ kernel.KernelRepository = (*PostgresRuntimeRepository)(nil)
