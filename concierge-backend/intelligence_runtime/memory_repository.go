package intelligence_runtime

import (
	"context"
	"maps"
	"sort"
	"sync"
	"time"

	"github.com/BraveAI93/concierge-backend/kernel"
)

// InMemoryRuntimeRepository is a production-shaped local implementation. A
// transaction clones the complete state, applies writes to the clone, and swaps
// it into place only on success. It opens no network connection and executes no
// SQL. Its methods intentionally satisfy kernel.KernelRepository as well as the
// runtime-specific transaction contract.
type InMemoryRuntimeRepository struct {
	mu    sync.RWMutex
	state *memoryState
}

type memoryState struct {
	persons       map[string]kernel.Person
	worlds        map[string]kernel.PersonalWorld
	entities      map[string]kernel.Entity
	contexts      map[string]kernel.Context
	sources       map[string]SourceRecord
	events        map[string]kernel.Event
	memories      map[string]kernel.Memory
	claims        map[string]kernel.Claim
	evidence      map[string]kernel.Evidence
	links         map[string]kernel.MemoryEventLink
	goals         map[string]kernel.Goal
	constraints   map[string]kernel.Constraint
	intents       map[string]kernel.PendingIntent
	loops         map[string]kernel.OpenLoop
	budgets       map[string]kernel.AttentionBudget
	allocations   map[string]kernel.AttentionAllocation
	opportunities map[string]kernel.Opportunity
	decisions     map[string]kernel.Decision
	permissions   map[string]kernel.Permission
	proposals     map[string]kernel.ActionProposal
	gates         map[string]kernel.ActionGate
	outcomes      map[string]kernel.Outcome
	audits        map[string]kernel.SelfAudit
	lineages      map[string]kernel.ClaimLineage
	replays       map[string]RuntimeResult
}

func newMemoryState() *memoryState {
	return &memoryState{
		persons: make(map[string]kernel.Person), worlds: make(map[string]kernel.PersonalWorld), entities: make(map[string]kernel.Entity), contexts: make(map[string]kernel.Context),
		sources: make(map[string]SourceRecord), events: make(map[string]kernel.Event), memories: make(map[string]kernel.Memory), claims: make(map[string]kernel.Claim), evidence: make(map[string]kernel.Evidence), links: make(map[string]kernel.MemoryEventLink),
		goals: make(map[string]kernel.Goal), constraints: make(map[string]kernel.Constraint), intents: make(map[string]kernel.PendingIntent), loops: make(map[string]kernel.OpenLoop), budgets: make(map[string]kernel.AttentionBudget), allocations: make(map[string]kernel.AttentionAllocation),
		opportunities: make(map[string]kernel.Opportunity), decisions: make(map[string]kernel.Decision), permissions: make(map[string]kernel.Permission), proposals: make(map[string]kernel.ActionProposal), gates: make(map[string]kernel.ActionGate), outcomes: make(map[string]kernel.Outcome), audits: make(map[string]kernel.SelfAudit), lineages: make(map[string]kernel.ClaimLineage), replays: make(map[string]RuntimeResult),
	}
}

func (s *memoryState) clone() *memoryState {
	return &memoryState{
		persons: maps.Clone(s.persons), worlds: maps.Clone(s.worlds), entities: maps.Clone(s.entities), contexts: maps.Clone(s.contexts),
		sources: maps.Clone(s.sources), events: maps.Clone(s.events), memories: maps.Clone(s.memories), claims: maps.Clone(s.claims), evidence: maps.Clone(s.evidence), links: maps.Clone(s.links),
		goals: maps.Clone(s.goals), constraints: maps.Clone(s.constraints), intents: maps.Clone(s.intents), loops: maps.Clone(s.loops), budgets: maps.Clone(s.budgets), allocations: maps.Clone(s.allocations),
		opportunities: maps.Clone(s.opportunities), decisions: maps.Clone(s.decisions), permissions: maps.Clone(s.permissions), proposals: maps.Clone(s.proposals), gates: maps.Clone(s.gates), outcomes: maps.Clone(s.outcomes), audits: maps.Clone(s.audits), lineages: maps.Clone(s.lineages), replays: maps.Clone(s.replays),
	}
}

func NewInMemoryRuntimeRepository() *InMemoryRuntimeRepository {
	return &InMemoryRuntimeRepository{state: newMemoryState()}
}

// SeedBinding, SeedGoal, SeedPermission, and SeedAttentionBudget represent a
// controlled local composition setup. They do not connect to legacy storage.
func (r *InMemoryRuntimeRepository) SeedBinding(binding PersonBinding) error {
	if r == nil || binding.Person.ID == "" || binding.World.PersonID != binding.Person.ID || binding.World.ID == "" {
		return ErrInvalidRuntimeConfig
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.state.persons[binding.Person.ID]; exists {
		return ErrDuplicateRuntimeRecord
	}
	r.state.persons[binding.Person.ID] = binding.Person
	r.state.worlds[binding.Person.ID] = binding.World
	return nil
}

func (r *InMemoryRuntimeRepository) SeedGoal(goal kernel.Goal) error {
	return r.withAutoTransaction(goal.PersonID, func(tx *memoryTransaction) error { return tx.SaveGoal(goal) })
}

func (r *InMemoryRuntimeRepository) SeedPermission(permission kernel.Permission) error {
	return r.withAutoTransaction(permission.PersonID, func(tx *memoryTransaction) error { return tx.SavePermission(permission) })
}

func (r *InMemoryRuntimeRepository) SeedAttentionBudget(budget kernel.AttentionBudget) error {
	return r.withAutoTransaction(budget.PersonID, func(tx *memoryTransaction) error { return tx.SaveAttentionBudget(budget) })
}

func (r *InMemoryRuntimeRepository) RunInTransaction(_ context.Context, personID string, fn func(RuntimeTransaction) error) error {
	if r == nil || personID == "" || fn == nil {
		return ErrInvalidRuntimeConfig
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.state.persons[personID]; !ok {
		return ErrUnknownIdentity
	}
	candidate := r.state.clone()
	tx := &memoryTransaction{personID: personID, state: candidate}
	if err := fn(tx); err != nil {
		return err
	}
	r.state = candidate
	return nil
}

func (r *InMemoryRuntimeRepository) withAutoTransaction(personID string, fn func(*memoryTransaction) error) error {
	return r.RunInTransaction(context.Background(), personID, func(tx RuntimeTransaction) error {
		return fn(tx.(*memoryTransaction))
	})
}

type memoryTransaction struct {
	personID string
	state    *memoryState
}

func (tx *memoryTransaction) ensurePerson(personID string) error {
	if personID == "" || personID != tx.personID {
		return ErrCrossPersonAccess
	}
	if _, ok := tx.state.persons[personID]; !ok {
		return ErrUnknownIdentity
	}
	return nil
}

func saveRecord[T any](destination map[string]T, id string, record T) error {
	if id == "" {
		return ErrInvalidRuntimeConfig
	}
	if _, exists := destination[id]; exists {
		return ErrDuplicateRuntimeRecord
	}
	destination[id] = record
	return nil
}

func (tx *memoryTransaction) FindReplay(key string) (RuntimeResult, bool) {
	result, ok := tx.state.replays[key]
	return result, ok
}
func (tx *memoryTransaction) StoreReplay(result RuntimeResult) error {
	if err := tx.ensurePerson(result.PersonID); err != nil {
		return err
	}
	return saveRecord(tx.state.replays, result.IdempotencyKey, result)
}
func (tx *memoryTransaction) StoreSource(source SourceRecord) error {
	if err := tx.ensurePerson(source.PersonID); err != nil {
		return err
	}
	return saveRecord(tx.state.sources, source.ID, source)
}
func (tx *memoryTransaction) SaveEvent(event kernel.Event) error {
	if err := tx.ensurePerson(event.PersonID); err != nil {
		return err
	}
	return saveRecord(tx.state.events, event.ID, event)
}
func (tx *memoryTransaction) SaveEvidence(evidence kernel.Evidence) error {
	if err := tx.ensurePerson(evidence.PersonID); err != nil {
		return err
	}
	return saveRecord(tx.state.evidence, evidence.ID, evidence)
}
func (tx *memoryTransaction) SaveMemory(memory kernel.Memory) error {
	if err := tx.ensurePerson(memory.PersonID); err != nil {
		return err
	}
	return saveRecord(tx.state.memories, memory.ID, memory)
}
func (tx *memoryTransaction) SaveClaim(claim kernel.Claim) error {
	if err := tx.ensurePerson(claim.PersonID); err != nil {
		return err
	}
	return saveRecord(tx.state.claims, claim.ID, claim)
}
func (tx *memoryTransaction) SaveMemoryEventLink(link kernel.MemoryEventLink) error {
	memory, memoryOK := tx.state.memories[link.MemoryID]
	event, eventOK := tx.state.events[link.EventID]
	if link.PersonID != tx.personID || !memoryOK || !eventOK || memory.PersonID != tx.personID || event.PersonID != tx.personID {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.links, link.MemoryID+":"+link.EventID, link)
}
func (tx *memoryTransaction) SavePendingIntent(intent kernel.PendingIntent) error {
	if err := tx.ensurePerson(intent.PersonID); err != nil {
		return err
	}
	return saveRecord(tx.state.intents, intent.ID, intent)
}
func (tx *memoryTransaction) SaveOpenLoop(loop kernel.OpenLoop) error {
	if err := tx.ensurePerson(loop.PersonID); err != nil {
		return err
	}
	if intent, ok := tx.state.intents[loop.PendingIntentID]; !ok || intent.PersonID != loop.PersonID {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.loops, loop.ID, loop)
}
func (tx *memoryTransaction) SaveOpportunity(opportunity kernel.Opportunity) error {
	if err := tx.ensurePerson(opportunity.PersonID); err != nil {
		return err
	}
	return saveRecord(tx.state.opportunities, opportunity.ID, opportunity)
}
func (tx *memoryTransaction) SaveDecision(decision kernel.Decision) error {
	if err := tx.ensurePerson(decision.PersonID); err != nil {
		return err
	}
	if opportunity, ok := tx.state.opportunities[decision.OpportunityID]; !ok || opportunity.PersonID != decision.PersonID {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.decisions, decision.ID, decision)
}
func (tx *memoryTransaction) SaveActionProposal(proposal kernel.ActionProposal) error {
	if err := tx.ensurePerson(proposal.PersonID); err != nil {
		return err
	}
	decision, decisionOK := tx.state.decisions[proposal.DecisionID]
	opportunity, opportunityOK := tx.state.opportunities[proposal.OpportunityID]
	if !decisionOK || !opportunityOK || decision.PersonID != proposal.PersonID || opportunity.PersonID != proposal.PersonID {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.proposals, proposal.ID, proposal)
}
func (tx *memoryTransaction) SaveActionGate(gate kernel.ActionGate) error {
	if err := tx.ensurePerson(gate.PersonID); err != nil {
		return err
	}
	proposal, ok := tx.state.proposals[gate.ActionProposalID]
	if !ok || proposal.PersonID != gate.PersonID {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.gates, gate.ID, gate)
}
func (tx *memoryTransaction) SaveAttentionAllocation(allocation kernel.AttentionAllocation) error {
	if err := tx.ensurePerson(allocation.PersonID); err != nil {
		return err
	}
	if _, ok := tx.state.budgets[allocation.BudgetID]; !ok {
		return ErrInvalidRuntimeConfig
	}
	return saveRecord(tx.state.allocations, allocation.BudgetID, allocation)
}
func (tx *memoryTransaction) SaveGoal(goal kernel.Goal) error {
	if err := tx.ensurePerson(goal.PersonID); err != nil {
		return err
	}
	return saveRecord(tx.state.goals, goal.ID, goal)
}
func (tx *memoryTransaction) SaveConstraint(constraint kernel.Constraint) error {
	if err := tx.ensurePerson(constraint.PersonID); err != nil {
		return err
	}
	return saveRecord(tx.state.constraints, constraint.ID, constraint)
}
func (tx *memoryTransaction) SavePermission(permission kernel.Permission) error {
	if err := tx.ensurePerson(permission.PersonID); err != nil {
		return err
	}
	return saveRecord(tx.state.permissions, permission.ID, permission)
}
func (tx *memoryTransaction) SaveAttentionBudget(budget kernel.AttentionBudget) error {
	if err := tx.ensurePerson(budget.PersonID); err != nil {
		return err
	}
	if err := budget.Validate(); err != nil {
		return err
	}
	return saveRecord(tx.state.budgets, budget.ID, budget)
}

func (tx *memoryTransaction) ListActiveGoals(at time.Time) []kernel.Goal {
	goals := make([]kernel.Goal, 0)
	for _, goal := range tx.state.goals {
		if goal.PersonID == tx.personID && goal.Status == kernel.GoalActive && goal.Temporal.IsActive(at) {
			goals = append(goals, goal)
		}
	}
	sort.Slice(goals, func(i, j int) bool { return goals[i].ID < goals[j].ID })
	return goals
}
func (tx *memoryTransaction) ListActivePermissions(at time.Time) []kernel.Permission {
	permissions := make([]kernel.Permission, 0)
	for _, permission := range tx.state.permissions {
		if permission.PersonID == tx.personID && permission.Temporal.IsActive(at) {
			permissions = append(permissions, permission)
		}
	}
	sort.Slice(permissions, func(i, j int) bool { return permissions[i].ID < permissions[j].ID })
	return permissions
}
func (tx *memoryTransaction) CurrentAttentionBudget(at time.Time) (kernel.AttentionBudget, bool) {
	var selected kernel.AttentionBudget
	found := false
	for _, budget := range tx.state.budgets {
		if budget.PersonID != tx.personID || at.Before(budget.WindowStart) || at.After(budget.WindowEnd) {
			continue
		}
		if !found || budget.ID < selected.ID {
			selected, found = budget, true
		}
	}
	return selected, found
}

// ReadState rejects a target person other than the authenticated requester's
// resolved canonical person. This is the repository-level second isolation
// check in addition to source-profile validation at ingestion.
func (r *InMemoryRuntimeRepository) ReadState(_ context.Context, requesterPersonID, targetPersonID string) (RuntimeState, error) {
	if r == nil || requesterPersonID == "" || requesterPersonID != targetPersonID {
		return RuntimeState{}, ErrCrossPersonAccess
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	person, ok := r.state.persons[targetPersonID]
	if !ok {
		return RuntimeState{}, ErrUnknownIdentity
	}
	world := r.state.worlds[targetPersonID]
	return buildRuntimeState(r.state, person, world), nil
}

func buildRuntimeState(s *memoryState, person kernel.Person, world kernel.PersonalWorld) RuntimeState {
	state := RuntimeState{Person: person, World: world}
	for _, source := range s.sources {
		if source.PersonID == person.ID {
			state.Sources = append(state.Sources, source)
		}
	}
	for _, event := range s.events {
		if event.PersonID == person.ID {
			state.Events = append(state.Events, event)
		}
	}
	for _, evidence := range s.evidence {
		if evidence.PersonID == person.ID {
			state.Evidence = append(state.Evidence, evidence)
		}
	}
	for _, memory := range s.memories {
		if memory.PersonID == person.ID {
			state.Memories = append(state.Memories, memory)
		}
	}
	for _, claim := range s.claims {
		if claim.PersonID == person.ID {
			state.Claims = append(state.Claims, claim)
		}
	}
	for _, goal := range s.goals {
		if goal.PersonID == person.ID {
			state.Goals = append(state.Goals, goal)
		}
	}
	for _, constraint := range s.constraints {
		if constraint.PersonID == person.ID {
			state.Constraints = append(state.Constraints, constraint)
		}
	}
	for _, intent := range s.intents {
		if intent.PersonID == person.ID {
			state.PendingIntents = append(state.PendingIntents, intent)
		}
	}
	for _, loop := range s.loops {
		if loop.PersonID == person.ID {
			state.OpenLoops = append(state.OpenLoops, loop)
		}
	}
	for _, budget := range s.budgets {
		if budget.PersonID == person.ID {
			state.AttentionBudgets = append(state.AttentionBudgets, budget)
		}
	}
	for _, allocation := range s.allocations {
		if allocation.PersonID == person.ID {
			state.Allocations = append(state.Allocations, allocation)
		}
	}
	for _, opportunity := range s.opportunities {
		if opportunity.PersonID == person.ID {
			state.Opportunities = append(state.Opportunities, opportunity)
		}
	}
	for _, decision := range s.decisions {
		if decision.PersonID == person.ID {
			state.Decisions = append(state.Decisions, decision)
		}
	}
	for _, permission := range s.permissions {
		if permission.PersonID == person.ID {
			state.Permissions = append(state.Permissions, permission)
		}
	}
	for _, proposal := range s.proposals {
		if proposal.PersonID == person.ID {
			state.ActionProposals = append(state.ActionProposals, proposal)
		}
	}
	for _, gate := range s.gates {
		if gate.PersonID == person.ID {
			state.ActionGates = append(state.ActionGates, gate)
		}
	}
	for _, outcome := range s.outcomes {
		if outcome.PersonID == person.ID {
			state.Outcomes = append(state.Outcomes, outcome)
		}
	}
	for _, audit := range s.audits {
		if audit.PersonID == person.ID {
			state.Audits = append(state.Audits, audit)
		}
	}
	for _, replay := range s.replays {
		if replay.PersonID == person.ID {
			state.Replays = append(state.Replays, replay)
		}
	}
	sortRuntimeState(&state)
	return state
}

func sortRuntimeState(state *RuntimeState) {
	sort.Slice(state.Sources, func(i, j int) bool { return state.Sources[i].ID < state.Sources[j].ID })
	sort.Slice(state.Events, func(i, j int) bool { return state.Events[i].ID < state.Events[j].ID })
	sort.Slice(state.Evidence, func(i, j int) bool { return state.Evidence[i].ID < state.Evidence[j].ID })
	sort.Slice(state.Memories, func(i, j int) bool { return state.Memories[i].ID < state.Memories[j].ID })
	sort.Slice(state.Claims, func(i, j int) bool { return state.Claims[i].ID < state.Claims[j].ID })
	sort.Slice(state.PendingIntents, func(i, j int) bool { return state.PendingIntents[i].ID < state.PendingIntents[j].ID })
	sort.Slice(state.OpenLoops, func(i, j int) bool { return state.OpenLoops[i].ID < state.OpenLoops[j].ID })
	sort.Slice(state.Opportunities, func(i, j int) bool { return state.Opportunities[i].ID < state.Opportunities[j].ID })
	sort.Slice(state.Decisions, func(i, j int) bool { return state.Decisions[i].ID < state.Decisions[j].ID })
	sort.Slice(state.ActionProposals, func(i, j int) bool { return state.ActionProposals[i].ID < state.ActionProposals[j].ID })
	sort.Slice(state.ActionGates, func(i, j int) bool { return state.ActionGates[i].ID < state.ActionGates[j].ID })
}

// KernelRepository compatibility methods. These route writes through the same
// transaction mechanism and preserve PersonID checks used by runtime writes.
func (r *InMemoryRuntimeRepository) FindPerson(_ context.Context, personID string) (kernel.Person, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	person, ok := r.state.persons[personID]
	if !ok {
		return kernel.Person{}, ErrUnknownIdentity
	}
	return person, nil
}
func (r *InMemoryRuntimeRepository) FindWorld(_ context.Context, personID string) (kernel.PersonalWorld, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	world, ok := r.state.worlds[personID]
	if !ok {
		return kernel.PersonalWorld{}, ErrUnknownIdentity
	}
	return world, nil
}
func (r *InMemoryRuntimeRepository) SaveWorld(_ context.Context, world kernel.PersonalWorld) error {
	return r.withAutoTransaction(world.PersonID, func(tx *memoryTransaction) error {
		if err := tx.ensurePerson(world.PersonID); err != nil {
			return err
		}
		tx.state.worlds[world.PersonID] = world
		return nil
	})
}
func (r *InMemoryRuntimeRepository) SaveEntity(_ context.Context, entity kernel.Entity) error {
	return r.withAutoTransaction(entity.PersonID, func(tx *memoryTransaction) error {
		if err := tx.ensurePerson(entity.PersonID); err != nil {
			return err
		}
		return saveRecord(tx.state.entities, entity.ID, entity)
	})
}
func (r *InMemoryRuntimeRepository) SaveContext(_ context.Context, item kernel.Context) error {
	return r.withAutoTransaction(item.PersonID, func(tx *memoryTransaction) error {
		if err := tx.ensurePerson(item.PersonID); err != nil {
			return err
		}
		return saveRecord(tx.state.contexts, item.ID, item)
	})
}
func (r *InMemoryRuntimeRepository) SaveEvent(_ context.Context, event kernel.Event) error {
	return r.withAutoTransaction(event.PersonID, func(tx *memoryTransaction) error { return tx.SaveEvent(event) })
}
func (r *InMemoryRuntimeRepository) SaveMemory(_ context.Context, memory kernel.Memory) error {
	return r.withAutoTransaction(memory.PersonID, func(tx *memoryTransaction) error { return tx.SaveMemory(memory) })
}
func (r *InMemoryRuntimeRepository) SaveClaim(_ context.Context, claim kernel.Claim) error {
	return r.withAutoTransaction(claim.PersonID, func(tx *memoryTransaction) error { return tx.SaveClaim(claim) })
}
func (r *InMemoryRuntimeRepository) SaveEvidence(_ context.Context, evidence kernel.Evidence) error {
	return r.withAutoTransaction(evidence.PersonID, func(tx *memoryTransaction) error { return tx.SaveEvidence(evidence) })
}
func (r *InMemoryRuntimeRepository) SaveMemoryEventLink(_ context.Context, link kernel.MemoryEventLink) error {
	if link.PersonID == "" {
		return ErrCrossPersonAccess
	}
	return r.withAutoTransaction(link.PersonID, func(tx *memoryTransaction) error { return tx.SaveMemoryEventLink(link) })
}
func (r *InMemoryRuntimeRepository) SaveGoal(_ context.Context, goal kernel.Goal) error {
	return r.SeedGoal(goal)
}
func (r *InMemoryRuntimeRepository) SaveConstraint(_ context.Context, constraint kernel.Constraint) error {
	return r.withAutoTransaction(constraint.PersonID, func(tx *memoryTransaction) error { return tx.SaveConstraint(constraint) })
}
func (r *InMemoryRuntimeRepository) SavePendingIntent(_ context.Context, intent kernel.PendingIntent) error {
	return r.withAutoTransaction(intent.PersonID, func(tx *memoryTransaction) error { return tx.SavePendingIntent(intent) })
}
func (r *InMemoryRuntimeRepository) SaveOpenLoop(_ context.Context, loop kernel.OpenLoop) error {
	return r.withAutoTransaction(loop.PersonID, func(tx *memoryTransaction) error { return tx.SaveOpenLoop(loop) })
}
func (r *InMemoryRuntimeRepository) SaveOpportunity(_ context.Context, opportunity kernel.Opportunity) error {
	return r.withAutoTransaction(opportunity.PersonID, func(tx *memoryTransaction) error { return tx.SaveOpportunity(opportunity) })
}
func (r *InMemoryRuntimeRepository) SaveDecision(_ context.Context, decision kernel.Decision) error {
	return r.withAutoTransaction(decision.PersonID, func(tx *memoryTransaction) error { return tx.SaveDecision(decision) })
}
func (r *InMemoryRuntimeRepository) SavePermission(_ context.Context, permission kernel.Permission) error {
	return r.SeedPermission(permission)
}
func (r *InMemoryRuntimeRepository) SaveActionProposal(_ context.Context, proposal kernel.ActionProposal) error {
	return r.withAutoTransaction(proposal.PersonID, func(tx *memoryTransaction) error { return tx.SaveActionProposal(proposal) })
}
func (r *InMemoryRuntimeRepository) SaveActionGate(_ context.Context, gate kernel.ActionGate) error {
	return r.withAutoTransaction(gate.PersonID, func(tx *memoryTransaction) error { return tx.SaveActionGate(gate) })
}
func (r *InMemoryRuntimeRepository) SaveOutcome(_ context.Context, outcome kernel.Outcome) error {
	return r.withAutoTransaction(outcome.PersonID, func(tx *memoryTransaction) error {
		if err := tx.ensurePerson(outcome.PersonID); err != nil {
			return err
		}
		return saveRecord(tx.state.outcomes, outcome.ID, outcome)
	})
}
func (r *InMemoryRuntimeRepository) SaveSelfAudit(_ context.Context, audit kernel.SelfAudit) error {
	return r.withAutoTransaction(audit.PersonID, func(tx *memoryTransaction) error {
		if err := tx.ensurePerson(audit.PersonID); err != nil {
			return err
		}
		return saveRecord(tx.state.audits, audit.ID, audit)
	})
}
func (r *InMemoryRuntimeRepository) SaveAttentionBudget(_ context.Context, budget kernel.AttentionBudget) error {
	return r.SeedAttentionBudget(budget)
}
func (r *InMemoryRuntimeRepository) SaveAttentionAllocation(_ context.Context, allocation kernel.AttentionAllocation) error {
	return r.withAutoTransaction(allocation.PersonID, func(tx *memoryTransaction) error { return tx.SaveAttentionAllocation(allocation) })
}
func (r *InMemoryRuntimeRepository) SaveClaimLineage(_ context.Context, lineage kernel.ClaimLineage) error {
	if lineage.PersonID == "" { return ErrCrossPersonAccess }
	return r.withAutoTransaction(lineage.PersonID, func(tx *memoryTransaction) error {
		claim, ok := tx.state.claims[lineage.ClaimID]
		if !ok || claim.PersonID != lineage.PersonID { return ErrCrossPersonAccess }
		return saveRecord(tx.state.lineages, lineage.ClaimID, lineage)
	})
}

var _ RuntimeRepository = (*InMemoryRuntimeRepository)(nil)
var _ kernel.KernelRepository = (*InMemoryRuntimeRepository)(nil)
