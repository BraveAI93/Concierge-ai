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
	persons             map[string]kernel.Person
	worlds              map[string]kernel.PersonalWorld
	entities            map[string]kernel.Entity
	contexts            map[string]kernel.Context
	sources             map[string]SourceRecord
	events              map[string]kernel.Event
	memories            map[string]kernel.Memory
	claims              map[string]kernel.Claim
	evidence            map[string]kernel.Evidence
	links               map[string]kernel.MemoryEventLink
	goals               map[string]kernel.Goal
	constraints         map[string]kernel.Constraint
	intents             map[string]kernel.PendingIntent
	loops               map[string]kernel.OpenLoop
	budgets             map[string]kernel.AttentionBudget
	allocations         map[string]kernel.AttentionAllocation
	opportunities       map[string]kernel.Opportunity
	decisions           map[string]kernel.Decision
	permissions         map[string]kernel.Permission
	proposals           map[string]kernel.ActionProposal
	gates               map[string]kernel.ActionGate
	outcomes            map[string]kernel.Outcome
	audits              map[string]kernel.SelfAudit
	lineages            map[string]kernel.ClaimLineage
	blocks              map[string]kernel.InteractionBlock
	threads             map[string]kernel.Thread
	continuityLinks     map[string]kernel.ContinuityLink
	threadDeltas        map[string]kernel.ThreadDelta
	threadStates        map[string]kernel.CurrentThreadState
	signals             map[string]kernel.ObservedInteractionSignal
	baselines           map[string]kernel.PersonalInteractionBaseline
	inferredStates      map[string]kernel.InferredInteractionState
	adaptations         map[string]kernel.InteractionAdaptationDecision
	attunementEpisodes  map[string]kernel.AttunementEpisode
	interventions       map[string]kernel.InteractionIntervention
	interactionOutcomes map[string]kernel.InteractionOutcome
	attunementPatterns  map[string]kernel.PersonalAttunementPattern
	replays             map[string]RuntimeResult
}

func newMemoryState() *memoryState {
	return &memoryState{
		persons: make(map[string]kernel.Person), worlds: make(map[string]kernel.PersonalWorld), entities: make(map[string]kernel.Entity), contexts: make(map[string]kernel.Context),
		sources: make(map[string]SourceRecord), events: make(map[string]kernel.Event), memories: make(map[string]kernel.Memory), claims: make(map[string]kernel.Claim), evidence: make(map[string]kernel.Evidence), links: make(map[string]kernel.MemoryEventLink),
		goals: make(map[string]kernel.Goal), constraints: make(map[string]kernel.Constraint), intents: make(map[string]kernel.PendingIntent), loops: make(map[string]kernel.OpenLoop), budgets: make(map[string]kernel.AttentionBudget), allocations: make(map[string]kernel.AttentionAllocation),
		opportunities: make(map[string]kernel.Opportunity), decisions: make(map[string]kernel.Decision), permissions: make(map[string]kernel.Permission), proposals: make(map[string]kernel.ActionProposal), gates: make(map[string]kernel.ActionGate), outcomes: make(map[string]kernel.Outcome), audits: make(map[string]kernel.SelfAudit), lineages: make(map[string]kernel.ClaimLineage), blocks: make(map[string]kernel.InteractionBlock), threads: make(map[string]kernel.Thread), continuityLinks: make(map[string]kernel.ContinuityLink), threadDeltas: make(map[string]kernel.ThreadDelta), threadStates: make(map[string]kernel.CurrentThreadState), signals: make(map[string]kernel.ObservedInteractionSignal), baselines: make(map[string]kernel.PersonalInteractionBaseline), inferredStates: make(map[string]kernel.InferredInteractionState), adaptations: make(map[string]kernel.InteractionAdaptationDecision), attunementEpisodes: make(map[string]kernel.AttunementEpisode), interventions: make(map[string]kernel.InteractionIntervention), interactionOutcomes: make(map[string]kernel.InteractionOutcome), attunementPatterns: make(map[string]kernel.PersonalAttunementPattern), replays: make(map[string]RuntimeResult),
	}
}

func (s *memoryState) clone() *memoryState {
	return &memoryState{
		persons: maps.Clone(s.persons), worlds: maps.Clone(s.worlds), entities: maps.Clone(s.entities), contexts: maps.Clone(s.contexts),
		sources: maps.Clone(s.sources), events: maps.Clone(s.events), memories: maps.Clone(s.memories), claims: maps.Clone(s.claims), evidence: maps.Clone(s.evidence), links: maps.Clone(s.links),
		goals: maps.Clone(s.goals), constraints: maps.Clone(s.constraints), intents: maps.Clone(s.intents), loops: maps.Clone(s.loops), budgets: maps.Clone(s.budgets), allocations: maps.Clone(s.allocations),
		opportunities: maps.Clone(s.opportunities), decisions: maps.Clone(s.decisions), permissions: maps.Clone(s.permissions), proposals: maps.Clone(s.proposals), gates: maps.Clone(s.gates), outcomes: maps.Clone(s.outcomes), audits: maps.Clone(s.audits), lineages: maps.Clone(s.lineages), blocks: maps.Clone(s.blocks), threads: maps.Clone(s.threads), continuityLinks: maps.Clone(s.continuityLinks), threadDeltas: maps.Clone(s.threadDeltas), threadStates: maps.Clone(s.threadStates), signals: maps.Clone(s.signals), baselines: maps.Clone(s.baselines), inferredStates: maps.Clone(s.inferredStates), adaptations: maps.Clone(s.adaptations), attunementEpisodes: maps.Clone(s.attunementEpisodes), interventions: maps.Clone(s.interventions), interactionOutcomes: maps.Clone(s.interactionOutcomes), attunementPatterns: maps.Clone(s.attunementPatterns), replays: maps.Clone(s.replays),
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

func (tx *memoryTransaction) SaveInteractionBlock(block kernel.InteractionBlock) error {
	if err := tx.ensurePerson(block.PersonID); err != nil || block.Validate() != nil {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.blocks, block.ID, block)
}
func (tx *memoryTransaction) SaveThread(thread kernel.Thread) error {
	if err := tx.ensurePerson(thread.PersonID); err != nil || thread.Validate() != nil {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.threads, thread.ID, thread)
}
func (tx *memoryTransaction) SaveContinuityLink(link kernel.ContinuityLink) error {
	if err := tx.ensurePerson(link.PersonID); err != nil || link.Validate() != nil {
		return ErrCrossPersonAccess
	}
	if link.Source.Kind == kernel.ContinuityBlock {
		if block, ok := tx.state.blocks[link.Source.ID]; !ok || block.PersonID != link.PersonID {
			return ErrCrossPersonAccess
		}
	}
	if link.Target.Kind == kernel.ContinuityThread {
		if thread, ok := tx.state.threads[link.Target.ID]; !ok || thread.PersonID != link.PersonID {
			return ErrCrossPersonAccess
		}
	}
	return saveRecord(tx.state.continuityLinks, link.ID, link)
}
func (tx *memoryTransaction) SaveThreadDelta(delta kernel.ThreadDelta) error {
	if err := tx.ensurePerson(delta.PersonID); err != nil || delta.Validate() != nil {
		return ErrCrossPersonAccess
	}
	if thread, ok := tx.state.threads[delta.TargetThreadID]; !ok || thread.PersonID != delta.PersonID {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.threadDeltas, delta.ID, delta)
}
func (tx *memoryTransaction) SaveCurrentThreadState(state kernel.CurrentThreadState) error {
	if err := tx.ensurePerson(state.PersonID); err != nil || state.Validate() != nil {
		return ErrCrossPersonAccess
	}
	if thread, ok := tx.state.threads[state.ThreadID]; !ok || thread.PersonID != state.PersonID {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.threadStates, state.ID, state)
}
func (tx *memoryTransaction) SaveObservedInteractionSignal(signal kernel.ObservedInteractionSignal) error {
	if err := tx.ensurePerson(signal.PersonID); err != nil || signal.Validate() != nil {
		return ErrCrossPersonAccess
	}
	if signal.BlockID != "" {
		if block, ok := tx.state.blocks[signal.BlockID]; !ok || block.PersonID != signal.PersonID {
			return ErrCrossPersonAccess
		}
	}
	return saveRecord(tx.state.signals, signal.ID, signal)
}
func (tx *memoryTransaction) SavePersonalInteractionBaseline(baseline kernel.PersonalInteractionBaseline) error {
	if err := tx.ensurePerson(baseline.PersonID); err != nil || baseline.Validate() != nil {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.baselines, baseline.ID, baseline)
}
func (tx *memoryTransaction) SaveInferredInteractionState(state kernel.InferredInteractionState) error {
	if err := tx.ensurePerson(state.PersonID); err != nil || state.Validate() != nil {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.inferredStates, state.ID, state)
}
func (tx *memoryTransaction) SaveInteractionAdaptationDecision(decision kernel.InteractionAdaptationDecision) error {
	if err := tx.ensurePerson(decision.PersonID); err != nil || kernel.ValidateAttunementDecision(decision) != nil {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.adaptations, decision.ID, decision)
}
func (tx *memoryTransaction) SaveAttunementEpisode(episode kernel.AttunementEpisode) error {
	if err := tx.ensurePerson(episode.PersonID); err != nil || episode.ID == "" || episode.BlockID == "" || episode.CreatedAt.IsZero() {
		return ErrCrossPersonAccess
	}
	if block, ok := tx.state.blocks[episode.BlockID]; !ok || block.PersonID != episode.PersonID {
		return ErrCrossPersonAccess
	}
	if episode.AdaptationDecisionID != "" {
		if decision, ok := tx.state.adaptations[episode.AdaptationDecisionID]; !ok || decision.PersonID != episode.PersonID {
			return ErrCrossPersonAccess
		}
	}
	return saveRecord(tx.state.attunementEpisodes, episode.ID, episode)
}
func (tx *memoryTransaction) SaveInteractionIntervention(intervention kernel.InteractionIntervention) error {
	if err := tx.ensurePerson(intervention.PersonID); err != nil || intervention.ID == "" || intervention.EpisodeID == "" || intervention.OccurredAt.IsZero() {
		return ErrCrossPersonAccess
	}
	if episode, ok := tx.state.attunementEpisodes[intervention.EpisodeID]; !ok || episode.PersonID != intervention.PersonID {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.interventions, intervention.ID, intervention)
}
func (tx *memoryTransaction) SaveInteractionOutcome(outcome kernel.InteractionOutcome) error {
	if err := tx.ensurePerson(outcome.PersonID); err != nil || outcome.ID == "" || outcome.EpisodeID == "" || outcome.InterventionID == "" || outcome.OccurredAt.IsZero() || outcome.Privacy != kernel.PrivacyOutcomeEvidence {
		return ErrCrossPersonAccess
	}
	if episode, ok := tx.state.attunementEpisodes[outcome.EpisodeID]; !ok || episode.PersonID != outcome.PersonID {
		return ErrCrossPersonAccess
	}
	if intervention, ok := tx.state.interventions[outcome.InterventionID]; !ok || intervention.PersonID != outcome.PersonID || intervention.EpisodeID != outcome.EpisodeID {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.interactionOutcomes, outcome.ID, outcome)
}
func (tx *memoryTransaction) SavePersonalAttunementPattern(pattern kernel.PersonalAttunementPattern) error {
	if err := tx.ensurePerson(pattern.PersonID); err != nil || pattern.Validate() != nil {
		return ErrCrossPersonAccess
	}
	return saveRecord(tx.state.attunementPatterns, pattern.ID, pattern)
}
func (tx *memoryTransaction) ListThreads() []kernel.Thread {
	threads := make([]kernel.Thread, 0)
	for _, thread := range tx.state.threads {
		if thread.PersonID == tx.personID {
			threads = append(threads, thread)
		}
	}
	sort.Slice(threads, func(i, j int) bool { return threads[i].ID < threads[j].ID })
	return threads
}
func (tx *memoryTransaction) ListInteractionBlocks(threadID string) []kernel.InteractionBlock {
	blocks := make([]kernel.InteractionBlock, 0)
	for _, block := range tx.state.blocks {
		if block.PersonID == tx.personID && containsID(block.ThreadIDs, threadID) {
			blocks = append(blocks, block)
		}
	}
	sort.Slice(blocks, func(i, j int) bool { return blocks[i].StartTemporal.EventAt.Before(blocks[j].StartTemporal.EventAt) })
	return blocks
}
func (tx *memoryTransaction) ListThreadDeltas(threadID string) []kernel.ThreadDelta {
	deltas := make([]kernel.ThreadDelta, 0)
	for _, delta := range tx.state.threadDeltas {
		if delta.PersonID == tx.personID && delta.TargetThreadID == threadID {
			deltas = append(deltas, delta)
		}
	}
	sort.Slice(deltas, func(i, j int) bool { return deltas[i].ID < deltas[j].ID })
	return deltas
}
func (tx *memoryTransaction) ListCurrentThreadStates(threadID string) []kernel.CurrentThreadState {
	states := make([]kernel.CurrentThreadState, 0)
	for _, state := range tx.state.threadStates {
		if state.PersonID == tx.personID && state.ThreadID == threadID {
			states = append(states, state)
		}
	}
	sort.Slice(states, func(i, j int) bool { return states[i].ReconstructedAt.After(states[j].ReconstructedAt) })
	return states
}
func (tx *memoryTransaction) ListAttunementPatterns(contextSignature string) []kernel.PersonalAttunementPattern {
	patterns := make([]kernel.PersonalAttunementPattern, 0)
	for _, pattern := range tx.state.attunementPatterns {
		if pattern.PersonID == tx.personID && pattern.ContextSignature == contextSignature {
			patterns = append(patterns, pattern)
		}
	}
	sort.Slice(patterns, func(i, j int) bool { return patterns[i].ID < patterns[j].ID })
	return patterns
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
	for _, block := range s.blocks {
		if block.PersonID == person.ID {
			state.InteractionBlocks = append(state.InteractionBlocks, block)
		}
	}
	for _, thread := range s.threads {
		if thread.PersonID == person.ID {
			state.Threads = append(state.Threads, thread)
		}
	}
	for _, link := range s.continuityLinks {
		if link.PersonID == person.ID {
			state.ContinuityLinks = append(state.ContinuityLinks, link)
		}
	}
	for _, delta := range s.threadDeltas {
		if delta.PersonID == person.ID {
			state.ThreadDeltas = append(state.ThreadDeltas, delta)
		}
	}
	for _, threadState := range s.threadStates {
		if threadState.PersonID == person.ID {
			state.CurrentThreadStates = append(state.CurrentThreadStates, threadState)
		}
	}
	for _, signal := range s.signals {
		if signal.PersonID == person.ID {
			state.ObservedSignals = append(state.ObservedSignals, signal)
		}
	}
	for _, baseline := range s.baselines {
		if baseline.PersonID == person.ID {
			state.InteractionBaselines = append(state.InteractionBaselines, baseline)
		}
	}
	for _, inferred := range s.inferredStates {
		if inferred.PersonID == person.ID {
			state.InferredInteractionStates = append(state.InferredInteractionStates, inferred)
		}
	}
	for _, adaptation := range s.adaptations {
		if adaptation.PersonID == person.ID {
			state.AdaptationDecisions = append(state.AdaptationDecisions, adaptation)
		}
	}
	for _, episode := range s.attunementEpisodes {
		if episode.PersonID == person.ID {
			state.AttunementEpisodes = append(state.AttunementEpisodes, episode)
		}
	}
	for _, intervention := range s.interventions {
		if intervention.PersonID == person.ID {
			state.InteractionInterventions = append(state.InteractionInterventions, intervention)
		}
	}
	for _, outcome := range s.interactionOutcomes {
		if outcome.PersonID == person.ID {
			state.InteractionOutcomes = append(state.InteractionOutcomes, outcome)
		}
	}
	for _, pattern := range s.attunementPatterns {
		if pattern.PersonID == person.ID {
			state.AttunementPatterns = append(state.AttunementPatterns, pattern)
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
	sort.Slice(state.InteractionBlocks, func(i, j int) bool { return state.InteractionBlocks[i].ID < state.InteractionBlocks[j].ID })
	sort.Slice(state.Threads, func(i, j int) bool { return state.Threads[i].ID < state.Threads[j].ID })
	sort.Slice(state.ContinuityLinks, func(i, j int) bool { return state.ContinuityLinks[i].ID < state.ContinuityLinks[j].ID })
	sort.Slice(state.ThreadDeltas, func(i, j int) bool { return state.ThreadDeltas[i].ID < state.ThreadDeltas[j].ID })
	sort.Slice(state.CurrentThreadStates, func(i, j int) bool { return state.CurrentThreadStates[i].ID < state.CurrentThreadStates[j].ID })
	sort.Slice(state.ObservedSignals, func(i, j int) bool { return state.ObservedSignals[i].ID < state.ObservedSignals[j].ID })
	sort.Slice(state.InteractionBaselines, func(i, j int) bool { return state.InteractionBaselines[i].ID < state.InteractionBaselines[j].ID })
	sort.Slice(state.InferredInteractionStates, func(i, j int) bool {
		return state.InferredInteractionStates[i].ID < state.InferredInteractionStates[j].ID
	})
	sort.Slice(state.AdaptationDecisions, func(i, j int) bool { return state.AdaptationDecisions[i].ID < state.AdaptationDecisions[j].ID })
	sort.Slice(state.AttunementEpisodes, func(i, j int) bool { return state.AttunementEpisodes[i].ID < state.AttunementEpisodes[j].ID })
	sort.Slice(state.InteractionOutcomes, func(i, j int) bool { return state.InteractionOutcomes[i].ID < state.InteractionOutcomes[j].ID })
	sort.Slice(state.AttunementPatterns, func(i, j int) bool { return state.AttunementPatterns[i].ID < state.AttunementPatterns[j].ID })
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
	if lineage.PersonID == "" {
		return ErrCrossPersonAccess
	}
	return r.withAutoTransaction(lineage.PersonID, func(tx *memoryTransaction) error {
		claim, ok := tx.state.claims[lineage.ClaimID]
		if !ok || claim.PersonID != lineage.PersonID {
			return ErrCrossPersonAccess
		}
		return saveRecord(tx.state.lineages, lineage.ClaimID, lineage)
	})
}
func (r *InMemoryRuntimeRepository) SaveInteractionBlock(_ context.Context, block kernel.InteractionBlock) error {
	return r.withAutoTransaction(block.PersonID, func(tx *memoryTransaction) error { return tx.SaveInteractionBlock(block) })
}
func (r *InMemoryRuntimeRepository) SaveThread(_ context.Context, thread kernel.Thread) error {
	return r.withAutoTransaction(thread.PersonID, func(tx *memoryTransaction) error { return tx.SaveThread(thread) })
}
func (r *InMemoryRuntimeRepository) SaveContinuityLink(_ context.Context, link kernel.ContinuityLink) error {
	return r.withAutoTransaction(link.PersonID, func(tx *memoryTransaction) error { return tx.SaveContinuityLink(link) })
}
func (r *InMemoryRuntimeRepository) SaveThreadDelta(_ context.Context, delta kernel.ThreadDelta) error {
	return r.withAutoTransaction(delta.PersonID, func(tx *memoryTransaction) error { return tx.SaveThreadDelta(delta) })
}
func (r *InMemoryRuntimeRepository) SaveCurrentThreadState(_ context.Context, state kernel.CurrentThreadState) error {
	return r.withAutoTransaction(state.PersonID, func(tx *memoryTransaction) error { return tx.SaveCurrentThreadState(state) })
}
func (r *InMemoryRuntimeRepository) SaveObservedInteractionSignal(_ context.Context, signal kernel.ObservedInteractionSignal) error {
	return r.withAutoTransaction(signal.PersonID, func(tx *memoryTransaction) error { return tx.SaveObservedInteractionSignal(signal) })
}
func (r *InMemoryRuntimeRepository) SavePersonalInteractionBaseline(_ context.Context, baseline kernel.PersonalInteractionBaseline) error {
	return r.withAutoTransaction(baseline.PersonID, func(tx *memoryTransaction) error { return tx.SavePersonalInteractionBaseline(baseline) })
}
func (r *InMemoryRuntimeRepository) SaveInferredInteractionState(_ context.Context, state kernel.InferredInteractionState) error {
	return r.withAutoTransaction(state.PersonID, func(tx *memoryTransaction) error { return tx.SaveInferredInteractionState(state) })
}
func (r *InMemoryRuntimeRepository) SaveInteractionAdaptationDecision(_ context.Context, decision kernel.InteractionAdaptationDecision) error {
	return r.withAutoTransaction(decision.PersonID, func(tx *memoryTransaction) error { return tx.SaveInteractionAdaptationDecision(decision) })
}
func (r *InMemoryRuntimeRepository) SaveAttunementEpisode(_ context.Context, episode kernel.AttunementEpisode) error {
	return r.withAutoTransaction(episode.PersonID, func(tx *memoryTransaction) error { return tx.SaveAttunementEpisode(episode) })
}
func (r *InMemoryRuntimeRepository) SaveInteractionIntervention(_ context.Context, intervention kernel.InteractionIntervention) error {
	return r.withAutoTransaction(intervention.PersonID, func(tx *memoryTransaction) error { return tx.SaveInteractionIntervention(intervention) })
}
func (r *InMemoryRuntimeRepository) SaveInteractionOutcome(_ context.Context, outcome kernel.InteractionOutcome) error {
	return r.withAutoTransaction(outcome.PersonID, func(tx *memoryTransaction) error { return tx.SaveInteractionOutcome(outcome) })
}
func (r *InMemoryRuntimeRepository) SavePersonalAttunementPattern(_ context.Context, pattern kernel.PersonalAttunementPattern) error {
	return r.withAutoTransaction(pattern.PersonID, func(tx *memoryTransaction) error { return tx.SavePersonalAttunementPattern(pattern) })
}

var _ RuntimeRepository = (*InMemoryRuntimeRepository)(nil)
var _ kernel.KernelRepository = (*InMemoryRuntimeRepository)(nil)
