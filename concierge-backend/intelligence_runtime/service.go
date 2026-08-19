package intelligence_runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/BraveAI93/concierge-backend/kernel"
)

// IngestConversationMessage is the first real vertical slice. It accepts an
// already-authenticated server principal and legacy source record, resolves the
// canonical person, then writes a complete source-to-ActionGate lineage in one
// transaction. It does not call an external action executor.
func (s RuntimeService) IngestConversationMessage(ctx context.Context, principal AuthenticatedPrincipal, source ConversationMessage) (RuntimeResult, error) {
	if !s.Feature.Enabled {
		return RuntimeResult{}, ErrRuntimeDisabled
	}
	if s.Identity == nil || s.Adapter == nil || s.Repo == nil || s.Clock == nil || s.Policy == nil || s.Config.Validate() != nil {
		return RuntimeResult{}, ErrInvalidRuntimeConfig
	}
	binding, err := s.Identity.Resolve(ctx, principal)
	if err != nil {
		return RuntimeResult{}, err
	}
	// The conversation's internal profile ID is checked only against the
	// server-resolved binding. No public caller may nominate a target person.
	if source.Conversation.ProfileID != binding.SourceProfileID {
		return RuntimeResult{}, ErrSourceUnauthorized
	}
	now := s.Clock.Now().UTC()
	bundle, err := s.Adapter.Map(binding, source, now)
	if err != nil {
		return RuntimeResult{}, err
	}
	if bundle.Source.PersonID != binding.Person.ID || bundle.Event.PersonID != binding.Person.ID || bundle.Evidence.PersonID != binding.Person.ID || bundle.Memory.PersonID != binding.Person.ID || bundle.Claim.PersonID != binding.Person.ID {
		return RuntimeResult{}, ErrSourceUnauthorized
	}

	result := RuntimeResult{}
	err = s.Repo.RunInTransaction(ctx, binding.Person.ID, func(tx RuntimeTransaction) error {
		if replay, ok := tx.FindReplay(bundle.IdempotencyKey); ok {
			result = replay
			result.Replayed = true
			return nil
		}
		if err := tx.StoreSource(bundle.Source); err != nil {
			return err
		}
		link, err := kernel.LinkMemoryToEvent(&bundle.Memory, &bundle.Event, "source message supports the memory", now)
		if err != nil {
			return err
		}
		if err := tx.SaveEvent(bundle.Event); err != nil {
			return err
		}
		if err := tx.SaveEvidence(bundle.Evidence); err != nil {
			return err
		}
		if err := tx.SaveMemory(bundle.Memory); err != nil {
			return err
		}
		if err := tx.SaveClaim(bundle.Claim); err != nil {
			return err
		}
		if err := tx.SaveMemoryEventLink(link); err != nil {
			return err
		}

		var intentID, loopID string
		if bundle.Intent != nil && bundle.OpenLoop != nil {
			if err := kernel.TransitionPendingIntent(bundle.Intent, kernel.IntentReady, now, "runtime", "source message classified as supported unresolved scheduling work"); err != nil {
				return err
			}
			if err := tx.SavePendingIntent(*bundle.Intent); err != nil {
				return err
			}
			if err := tx.SaveOpenLoop(*bundle.OpenLoop); err != nil {
				return err
			}
			intentID, loopID = bundle.Intent.ID, bundle.OpenLoop.ID
		}
		if bundle.OpenLoop == nil {
			return ErrUnsupportedSource
		}

		budget, ok := tx.CurrentAttentionBudget(now)
		if !ok {
			return ErrNoAttentionBudget
		}
		return s.composePlan(tx, binding, bundle, budget, now, intentID, loopID, &result)
	})
	if err != nil {
		return RuntimeResult{}, err
	}
	return result, nil
}

func (s RuntimeService) composePlan(tx RuntimeTransaction, binding PersonBinding, bundle IngestionBundle, budget kernel.AttentionBudget, now time.Time, intentID, loopID string, result *RuntimeResult) error {
	temporalEvaluation, err := s.Policy.EvaluateTemporal(kernel.TemporalPriorityInput{
		Moment:               kernel.EvaluationMoment{WallClockAt: now},
		Temporal:             bundle.OpenLoop.Attention,
		InteractionGap:       bundle.OpenLoop.InteractionGap,
		Effort:               bundle.OpenLoop.AttentionNeed,
		SubjectiveImportance: s.Config.SchedulingPriority.SubjectiveImportance,
	})
	if err != nil {
		return err
	}
	candidate, err := kernel.AttentionCandidateFromOpenLoop(*bundle.OpenLoop, temporalEvaluation.Utility)
	if err != nil {
		return err
	}
	allocation, err := kernel.AllocateAttention(budget, []kernel.AttentionCandidate{candidate}, s.Policy)
	if err != nil {
		return err
	}
	if err := tx.SaveAttentionAllocation(allocation); err != nil {
		return err
	}
	if len(allocation.Surfaced) == 0 {
		return ErrNoAttentionBudget
	}
	goals := tx.ListActiveGoals(now)
	if len(goals) == 0 {
		return ErrNoActiveGoal
	}
	goal := goals[0]

	opportunity := kernel.Opportunity{
		ID:            "opportunity:" + bundle.Source.ID,
		PersonID:      binding.Person.ID,
		Title:         "Prepare a scheduling response from source message " + bundle.Source.MessageID,
		Summary:       "Source-supported unresolved scheduling request; no calendar availability is assumed.",
		GoalIDs:       []string{goal.ID},
		EvidenceIDs:   []string{bundle.Evidence.ID},
		Temporal:      bundle.Event.Temporal,
		Priority:      s.Config.SchedulingPriority,
		ActionWindow:  bundle.Deadline,
		AttentionNeed: s.Config.SchedulingEffort,
		CreatedAt:     now,
	}
	evaluation, err := kernel.EvaluateOpportunityWithPolicy(&opportunity, goals, nil, s.Policy, kernel.EvaluationMoment{WallClockAt: now})
	if err != nil {
		return err
	}
	if err := tx.SaveOpportunity(opportunity); err != nil {
		return err
	}
	decision, err := kernel.DecideOpportunityWithPolicy(opportunity, "decision:"+bundle.Source.ID, s.Policy, kernel.EvaluationMoment{WallClockAt: now})
	if err != nil {
		return err
	}
	if err := tx.SaveDecision(decision); err != nil {
		return err
	}

	permission, ok := matchingPermission(tx.ListActivePermissions(now), binding.Person.ID, s.Config.RequestedScope, now)
	if !ok {
		return ErrNoPermission
	}
	proposal := kernel.ActionProposal{
		ID:            "proposal:" + bundle.Source.ID,
		PersonID:      binding.Person.ID,
		OpportunityID: opportunity.ID,
		DecisionID:    decision.ID,
		Title:         "Create an unsent scheduling-response draft for review",
		Requested:     s.Config.RequestedScope,
		PermissionID:  permission.ID,
		Parameters: map[string]string{
			"source_message_id": bundle.Source.MessageID,
			"source_ref":        bundle.Evidence.Provenance.SourceRef,
			"intent_id":         intentID,
			"open_loop_id":      loopID,
		},
		CreatedAt: now,
	}
	gate := kernel.ActionGate{
		ID:               "gate:" + bundle.Source.ID,
		PersonID:         binding.Person.ID,
		ActionProposalID: proposal.ID,
		State:            kernel.GateDraft,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := tx.SaveActionProposal(proposal); err != nil {
		return err
	}
	if err := kernel.PrepareActionGate(&gate, proposal, permission, now, "runtime"); err != nil {
		return err
	}
	if err := tx.SaveActionGate(gate); err != nil {
		return err
	}

	*result = RuntimeResult{
		PersonID:       binding.Person.ID,
		IdempotencyKey: bundle.IdempotencyKey,
		EventID:        bundle.Event.ID,
		EvidenceID:     bundle.Evidence.ID,
		MemoryID:       bundle.Memory.ID,
		ClaimID:        bundle.Claim.ID,
		IntentID:       intentID,
		OpenLoopID:     loopID,
		OpportunityID:  opportunity.ID,
		DecisionID:     decision.ID,
		ProposalID:     proposal.ID,
		ActionGateID:   gate.ID,
	}
	_ = evaluation
	return tx.StoreReplay(*result)
}

func matchingPermission(permissions []kernel.Permission, personID string, requested kernel.Scope, now time.Time) (kernel.Permission, bool) {
	for _, permission := range permissions {
		if permission.PersonID != personID || !permission.Temporal.IsActive(now) {
			continue
		}
		if kernel.ValidateProposalPermission(kernel.ActionProposal{ID: "permission-check", PersonID: personID, PermissionID: permission.ID, Requested: requested}, permission, now) == nil {
			return permission, true
		}
	}
	return kernel.Permission{}, false
}

// State exposes only the authenticated principal's resolved personal world.
// There is deliberately no target-person parameter.
func (s RuntimeService) State(ctx context.Context, principal AuthenticatedPrincipal) (RuntimeState, error) {
	if !s.Feature.Enabled {
		return RuntimeState{}, ErrRuntimeDisabled
	}
	if s.Identity == nil || s.Repo == nil {
		return RuntimeState{}, ErrInvalidRuntimeConfig
	}
	binding, err := s.Identity.Resolve(ctx, principal)
	if err != nil {
		return RuntimeState{}, err
	}
	return s.Repo.ReadState(ctx, binding.Person.ID, binding.Person.ID)
}

func (s RuntimeService) String() string {
	return fmt.Sprintf("RuntimeService{enabled:%t}", s.Feature.Enabled)
}
