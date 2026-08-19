package kernel

import (
	stdcontext "context"
	"time"
)

// PersonalWorldRepository is the future persistence boundary for the
// person-scoped aggregate. Implementations must reject cross-person reads and
// writes rather than relying on callers to enforce tenancy.
type PersonalWorldRepository interface {
	FindPerson(ctx stdcontext.Context, personID string) (Person, error)
	FindWorld(ctx stdcontext.Context, personID string) (PersonalWorld, error)
	SaveWorld(ctx stdcontext.Context, world PersonalWorld) error
	SaveEntity(ctx stdcontext.Context, entity Entity) error
	SaveContext(ctx stdcontext.Context, context Context) error
}

// KnowledgeRepository stores the canonical event, memory, claim, evidence, and
// relationship records. It deliberately accepts canonical values, never legacy
// profile maps or raw provider payloads.
type KnowledgeRepository interface {
	SaveEvent(ctx stdcontext.Context, event Event) error
	SaveMemory(ctx stdcontext.Context, memory Memory) error
	SaveClaim(ctx stdcontext.Context, claim Claim) error
	SaveEvidence(ctx stdcontext.Context, evidence Evidence) error
	SaveMemoryEventLink(ctx stdcontext.Context, link MemoryEventLink) error
}

// IntentionRepository stores attention, goals, constraints, and pending work.
type IntentionRepository interface {
	SaveGoal(ctx stdcontext.Context, goal Goal) error
	SaveConstraint(ctx stdcontext.Context, constraint Constraint) error
	SavePendingIntent(ctx stdcontext.Context, intent PendingIntent) error
	SaveOpenLoop(ctx stdcontext.Context, loop OpenLoop) error
}

// DeliberationRepository stores a reproducible evaluation trail. Evaluation
// inputs and outputs are persisted together to support later self-audit.
type DeliberationRepository interface {
	SaveOpportunity(ctx stdcontext.Context, opportunity Opportunity) error
	SaveDecision(ctx stdcontext.Context, decision Decision) error
	SavePermission(ctx stdcontext.Context, permission Permission) error
	SaveActionProposal(ctx stdcontext.Context, proposal ActionProposal) error
	SaveActionGate(ctx stdcontext.Context, gate ActionGate) error
}

// LearningRepository records actual results and the self-audit derived from
// them. It is separate so an executor cannot mark its own result as successful.
type LearningRepository interface {
	SaveOutcome(ctx stdcontext.Context, outcome Outcome) error
	SaveSelfAudit(ctx stdcontext.Context, audit SelfAudit) error
}

// AttentionRepository is a future persistence port for finite attention
// budgets, their transparent selection/deferral records, and claim lineage.
type AttentionRepository interface {
	SaveAttentionBudget(ctx stdcontext.Context, budget AttentionBudget) error
	SaveAttentionAllocation(ctx stdcontext.Context, allocation AttentionAllocation) error
	SaveClaimLineage(ctx stdcontext.Context, lineage ClaimLineage) error
}

// KernelRepository groups storage ports for a composition root. This package
// supplies no concrete implementation in v0.2.
type KernelRepository interface {
	PersonalWorldRepository
	KnowledgeRepository
	IntentionRepository
	DeliberationRepository
	LearningRepository
	AttentionRepository
}

// LegacyEventInput is intentionally minimal. Existing application records must
// be mapped into it by a compatibility adapter rather than imported directly by
// the domain package.
type LegacyEventInput struct {
	PersonID    string
	ExternalID  string
	Kind        string
	Summary     string
	OccurredAt  time.Time
	ObservedAt  time.Time
	EffectiveAt time.Time
	AttentionAt time.Time
	SourceType  string
	SourceRef   string
}

// EventMapper is the boundary for mapping any legacy source into an Event.
type EventMapper interface {
	ToKernelEvent(input LegacyEventInput) (Event, Evidence, error)
}

// ActionExecutor is a side-effect port. An implementation must execute only a
// gate already in GateApproved state and return a receipt that can be attached
// as outcome provenance. It must not mutate domain records itself.
type ActionExecutor interface {
	Execute(ctx stdcontext.Context, proposal ActionProposal, gate ActionGate) (ActionReceipt, error)
}

type ActionReceipt struct {
	ExternalID string
	OccurredAt time.Time
	Provenance Provenance
}

// Clock makes current-time dependency explicit for service orchestration. All
// deterministic primitives receive their reference time as an argument.
type Clock interface {
	Now() time.Time
}
