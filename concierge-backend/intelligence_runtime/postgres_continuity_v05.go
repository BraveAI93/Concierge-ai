package intelligence_runtime

import (
	"context"
	"sort"

	"github.com/BraveAI93/concierge-backend/kernel"
)

func (tx *postgresTransaction) SaveInteractionBlock(value kernel.InteractionBlock) error {
	if value.Validate() != nil {
		return ErrInvalidRuntimeConfig
	}
	return tx.insertRecord("interaction_block", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveThread(value kernel.Thread) error {
	if value.Validate() != nil {
		return ErrInvalidRuntimeConfig
	}
	return tx.insertRecord("thread", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveContinuityLink(value kernel.ContinuityLink) error {
	if err := tx.ensurePerson(value.PersonID); err != nil || value.Validate() != nil {
		return ErrCrossPersonAccess
	}
	for _, ref := range []kernel.ContinuityRef{value.Source, value.Target} {
		ok, err := tx.recordExists(string(ref.Kind), ref.ID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrCrossPersonAccess
		}
	}
	payload, err := jsonPayload(value)
	if err != nil {
		return err
	}
	_, err = tx.tx.ExecContext(tx.ctx, `INSERT INTO ci_kernel_v04.continuity_links
		(person_id,id,source_kind,source_id,target_kind,target_id,relation,payload,created_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8::jsonb,$9)`, value.PersonID, value.ID, value.Source.Kind, value.Source.ID, value.Target.Kind, value.Target.ID, value.Relation, string(payload), value.CreatedAt)
	if isUniqueViolation(err) {
		return ErrDuplicateRuntimeRecord
	}
	return err
}
func (tx *postgresTransaction) SaveThreadDelta(value kernel.ThreadDelta) error {
	if err := tx.ensurePerson(value.PersonID); err != nil || value.Validate() != nil {
		return ErrCrossPersonAccess
	}
	threadOK, err := tx.recordExists("thread", value.TargetThreadID)
	if err != nil {
		return err
	}
	originOK, err := tx.recordExists(string(value.Originating.Kind), value.Originating.ID)
	if err != nil {
		return err
	}
	if !threadOK || !originOK {
		return ErrCrossPersonAccess
	}
	if err := tx.insertRecord("thread_delta", value.ID, value.PersonID, value); err != nil {
		return err
	}
	payload, err := jsonPayload(value)
	if err != nil {
		return err
	}
	_, err = tx.tx.ExecContext(tx.ctx, `INSERT INTO ci_kernel_v04.thread_deltas
		(person_id,id,target_thread_id,origin_kind,origin_id,payload,event_at,evaluated_at)
		VALUES($1,$2,$3,$4,$5,$6::jsonb,$7,$8)`, value.PersonID, value.ID, value.TargetThreadID, value.Originating.Kind, value.Originating.ID, string(payload), value.EventAt, value.EvaluatedAt)
	if isUniqueViolation(err) {
		return ErrDuplicateRuntimeRecord
	}
	return err
}
func (tx *postgresTransaction) SaveCurrentThreadState(value kernel.CurrentThreadState) error {
	if err := tx.ensurePerson(value.PersonID); err != nil || value.Validate() != nil {
		return ErrCrossPersonAccess
	}
	ok, err := tx.recordExists("thread", value.ThreadID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCrossPersonAccess
	}
	return tx.insertRecord("current_thread_state", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveObservedInteractionSignal(value kernel.ObservedInteractionSignal) error {
	if err := tx.ensurePerson(value.PersonID); err != nil || value.Validate() != nil {
		return ErrCrossPersonAccess
	}
	if value.BlockID != "" {
		ok, err := tx.recordExists("interaction_block", value.BlockID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrCrossPersonAccess
		}
	}
	return tx.insertRecord("observed_interaction_signal", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SavePersonalInteractionBaseline(value kernel.PersonalInteractionBaseline) error {
	if err := tx.ensurePerson(value.PersonID); err != nil || value.Validate() != nil {
		return ErrCrossPersonAccess
	}
	return tx.insertRecord("personal_interaction_baseline", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveInferredInteractionState(value kernel.InferredInteractionState) error {
	if err := tx.ensurePerson(value.PersonID); err != nil || value.Validate() != nil {
		return ErrCrossPersonAccess
	}
	return tx.insertRecord("inferred_interaction_state", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveInteractionAdaptationDecision(value kernel.InteractionAdaptationDecision) error {
	if err := tx.ensurePerson(value.PersonID); err != nil || kernel.ValidateAttunementDecision(value) != nil {
		return ErrCrossPersonAccess
	}
	return tx.insertRecord("interaction_adaptation_decision", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveAttunementEpisode(value kernel.AttunementEpisode) error {
	if err := tx.ensurePerson(value.PersonID); err != nil || value.ID == "" || value.BlockID == "" {
		return ErrCrossPersonAccess
	}
	blockOK, err := tx.recordExists("interaction_block", value.BlockID)
	if err != nil {
		return err
	}
	if !blockOK {
		return ErrCrossPersonAccess
	}
	if value.AdaptationDecisionID != "" {
		ok, err := tx.recordExists("interaction_adaptation_decision", value.AdaptationDecisionID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrCrossPersonAccess
		}
	}
	return tx.insertRecord("attunement_episode", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveInteractionIntervention(value kernel.InteractionIntervention) error {
	if err := tx.ensurePerson(value.PersonID); err != nil || value.ID == "" || value.EpisodeID == "" {
		return ErrCrossPersonAccess
	}
	ok, err := tx.recordExists("attunement_episode", value.EpisodeID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCrossPersonAccess
	}
	return tx.insertRecord("interaction_intervention", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SaveInteractionOutcome(value kernel.InteractionOutcome) error {
	if err := tx.ensurePerson(value.PersonID); err != nil || value.ID == "" || value.EpisodeID == "" || value.InterventionID == "" || value.Privacy != kernel.PrivacyOutcomeEvidence {
		return ErrCrossPersonAccess
	}
	ok, err := tx.recordExists("attunement_episode", value.EpisodeID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCrossPersonAccess
	}
	var interventionPayload []byte
	if err := tx.tx.QueryRowContext(tx.ctx, "SELECT payload FROM ci_kernel_v04.records WHERE record_kind='interaction_intervention' AND person_id=$1 AND id=$2", value.PersonID, value.InterventionID).Scan(&interventionPayload); err != nil { return ErrCrossPersonAccess }
	var intervention kernel.InteractionIntervention
	if decodePayload(interventionPayload, &intervention) != nil || intervention.EpisodeID != value.EpisodeID { return ErrCrossPersonAccess }
	return tx.insertRecord("interaction_outcome", value.ID, value.PersonID, value)
}
func (tx *postgresTransaction) SavePersonalAttunementPattern(value kernel.PersonalAttunementPattern) error {
	if err := tx.ensurePerson(value.PersonID); err != nil || value.Validate() != nil {
		return ErrCrossPersonAccess
	}
	return tx.insertRecord("personal_attunement_pattern", value.ID, value.PersonID, value)
}

func (tx *postgresTransaction) listRecords(kind string, decode func([]byte) bool) {
	rows, err := tx.tx.QueryContext(tx.ctx, "SELECT payload FROM ci_kernel_v04.records WHERE person_id=$1 AND record_kind=$2 ORDER BY id", tx.personID, kind)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var payload []byte
		if rows.Scan(&payload) == nil {
			decode(payload)
		}
	}
}
func (tx *postgresTransaction) ListThreads() []kernel.Thread {
	values := make([]kernel.Thread, 0)
	tx.listRecords("thread", func(payload []byte) bool {
		var value kernel.Thread
		if decodePayload(payload, &value) == nil {
			values = append(values, value)
		}
		return true
	})
	return values
}
func (tx *postgresTransaction) ListInteractionBlocks(threadID string) []kernel.InteractionBlock {
	values := make([]kernel.InteractionBlock, 0)
	tx.listRecords("interaction_block", func(payload []byte) bool {
		var value kernel.InteractionBlock
		if decodePayload(payload, &value) == nil && containsID(value.ThreadIDs, threadID) {
			values = append(values, value)
		}
		return true
	})
	sort.Slice(values, func(i, j int) bool { return values[i].StartTemporal.EventAt.Before(values[j].StartTemporal.EventAt) })
	return values
}
func (tx *postgresTransaction) ListThreadDeltas(threadID string) []kernel.ThreadDelta {
	values := make([]kernel.ThreadDelta, 0)
	tx.listRecords("thread_delta", func(payload []byte) bool {
		var value kernel.ThreadDelta
		if decodePayload(payload, &value) == nil && value.TargetThreadID == threadID {
			values = append(values, value)
		}
		return true
	})
	sort.Slice(values, func(i, j int) bool { return values[i].ID < values[j].ID })
	return values
}
func (tx *postgresTransaction) ListCurrentThreadStates(threadID string) []kernel.CurrentThreadState {
	values := make([]kernel.CurrentThreadState, 0)
	tx.listRecords("current_thread_state", func(payload []byte) bool {
		var value kernel.CurrentThreadState
		if decodePayload(payload, &value) == nil && value.ThreadID == threadID {
			values = append(values, value)
		}
		return true
	})
	sort.Slice(values, func(i, j int) bool { return values[i].ReconstructedAt.After(values[j].ReconstructedAt) })
	return values
}
func (tx *postgresTransaction) ListAttunementPatterns(contextSignature string) []kernel.PersonalAttunementPattern {
	values := make([]kernel.PersonalAttunementPattern, 0)
	tx.listRecords("personal_attunement_pattern", func(payload []byte) bool {
		var value kernel.PersonalAttunementPattern
		if decodePayload(payload, &value) == nil && value.ContextSignature == contextSignature {
			values = append(values, value)
		}
		return true
	})
	return values
}

// Kernel repository wrappers retain the same person-scoped transaction boundary.
func (r *PostgresRuntimeRepository) SaveInteractionBlock(ctx context.Context, value kernel.InteractionBlock) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveInteractionBlock(value) })
}
func (r *PostgresRuntimeRepository) SaveThread(ctx context.Context, value kernel.Thread) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveThread(value) })
}
func (r *PostgresRuntimeRepository) SaveContinuityLink(ctx context.Context, value kernel.ContinuityLink) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveContinuityLink(value) })
}
func (r *PostgresRuntimeRepository) SaveThreadDelta(ctx context.Context, value kernel.ThreadDelta) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveThreadDelta(value) })
}
func (r *PostgresRuntimeRepository) SaveCurrentThreadState(ctx context.Context, value kernel.CurrentThreadState) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveCurrentThreadState(value) })
}
func (r *PostgresRuntimeRepository) SaveObservedInteractionSignal(ctx context.Context, value kernel.ObservedInteractionSignal) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveObservedInteractionSignal(value) })
}
func (r *PostgresRuntimeRepository) SavePersonalInteractionBaseline(ctx context.Context, value kernel.PersonalInteractionBaseline) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SavePersonalInteractionBaseline(value) })
}
func (r *PostgresRuntimeRepository) SaveInferredInteractionState(ctx context.Context, value kernel.InferredInteractionState) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveInferredInteractionState(value) })
}
func (r *PostgresRuntimeRepository) SaveInteractionAdaptationDecision(ctx context.Context, value kernel.InteractionAdaptationDecision) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveInteractionAdaptationDecision(value) })
}
func (r *PostgresRuntimeRepository) SaveAttunementEpisode(ctx context.Context, value kernel.AttunementEpisode) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveAttunementEpisode(value) })
}
func (r *PostgresRuntimeRepository) SaveInteractionIntervention(ctx context.Context, value kernel.InteractionIntervention) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveInteractionIntervention(value) })
}
func (r *PostgresRuntimeRepository) SaveInteractionOutcome(ctx context.Context, value kernel.InteractionOutcome) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SaveInteractionOutcome(value) })
}
func (r *PostgresRuntimeRepository) SavePersonalAttunementPattern(ctx context.Context, value kernel.PersonalAttunementPattern) error {
	return r.withAutoTransaction(ctx, value.PersonID, func(tx *postgresTransaction) error { return tx.SavePersonalAttunementPattern(value) })
}
