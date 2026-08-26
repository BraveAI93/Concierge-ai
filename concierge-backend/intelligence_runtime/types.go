// Package intelligence_runtime composes legacy Concierge records around the
// independent kernel package. It is deliberately not registered with HTTP or
// production persistence in v0.3.
package intelligence_runtime

import (
	"context"
	"errors"
	"time"

	"github.com/BraveAI93/concierge-backend/db"
	"github.com/BraveAI93/concierge-backend/kernel"
)

var (
	ErrRuntimeDisabled              = errors.New("intelligence runtime: feature is disabled")
	ErrUnknownIdentity              = errors.New("intelligence runtime: authenticated identity is unknown")
	ErrSourceUnauthorized           = errors.New("intelligence runtime: source profile is not bound to authenticated person")
	ErrCrossPersonAccess            = errors.New("intelligence runtime: cross-person access denied")
	ErrUnsupportedSource            = errors.New("intelligence runtime: source does not support canonical ingestion")
	ErrNoActiveGoal                 = errors.New("intelligence runtime: no active goal available for opportunity")
	ErrNoAttentionBudget            = errors.New("intelligence runtime: no active attention budget")
	ErrNoPermission                 = errors.New("intelligence runtime: no matching action permission")
	ErrInvalidRuntimeConfig         = errors.New("intelligence runtime: invalid runtime configuration")
	ErrDuplicateRuntimeRecord       = errors.New("intelligence runtime: duplicate runtime record")
	ErrUnsafePersistenceTarget      = errors.New("intelligence runtime: persistence target is not an approved local staging database")
	ErrIdentityProvisioningRequired = errors.New("intelligence runtime: identity provisioning requires the separate provisioner boundary")
)

// Feature is the server-side activation boundary. The default constructor is
// disabled and no route currently receives or changes this object.
type Feature struct {
	Enabled bool
}

func DisabledFeature() Feature { return Feature{} }
func EnabledFeature() Feature  { return Feature{Enabled: true} }

// AuthenticatedPrincipal must be created by server-side session/account
// authentication. It contains no public slug, profile selector, or target
// person ID, so a browser cannot redirect a request to another personal world.
type AuthenticatedPrincipal struct {
	StableSubjectID string
}

// PersonBinding is the stable server-side bridge between an authenticated
// Concierge account subject, an internal source profile ID, and a canonical
// Person. SourceProfileID is not a public routing slug.
type PersonBinding struct {
	StableSubjectID         string
	SourceProfileID         string
	AllowedSourceProfileIDs []string
	Person                  kernel.Person
	World                   kernel.PersonalWorld
}

// AllowsSourceProfile permits the primary internal profile and explicit linked
// internal profile IDs. It never treats a public browser slug as an identity.
func (b PersonBinding) AllowsSourceProfile(sourceProfileID string) bool {
	if sourceProfileID == "" {
		return false
	}
	if sourceProfileID == b.SourceProfileID {
		return true
	}
	for _, linked := range b.AllowedSourceProfileIDs {
		if linked == sourceProfileID {
			return true
		}
	}
	return false
}

// IdentityResolver is a runtime boundary. A production implementation must
// resolve an authenticated session to a stable account/profile subject first,
// then look up the binding without trusting caller-supplied profile identifiers.
type IdentityResolver interface {
	Resolve(ctx context.Context, principal AuthenticatedPrincipal) (PersonBinding, error)
}

// ConversationMessage is the first supported legacy source. The runtime
// retains both original records and does not query the legacy database itself.
type ConversationMessage struct {
	Conversation db.Conversation
	Message      db.Message
}

// SourceRecord preserves source IDs and timestamps for deterministic replay and
// provenance. The canonical evidence carries a SourceRef back to these values.
type SourceRecord struct {
	ID             string
	PersonID       string
	ProfileID      string
	ConversationID string
	SessionID      string
	MessageID      string
	MessageRole    string
	Content        string
	ConversationAt time.Time
	MessageAt      time.Time
	StoredAt       time.Time
}

// IngestionBundle is the conservative canonical output of the first source
// adapter. Nil intent/open loop means no unresolved work was justified.
type IngestionBundle struct {
	IdempotencyKey string
	Source         SourceRecord
	Event          kernel.Event
	Evidence       kernel.Evidence
	Memory         kernel.Memory
	Claim          kernel.Claim
	Intent         *kernel.PendingIntent
	OpenLoop       *kernel.OpenLoop
	Deadline       kernel.ActionWindow
}

// ConversationMessageAdapter maps legacy source records one way into canonical
// objects. It must preserve source-supported facts and reject unsupported text.
type ConversationMessageAdapter interface {
	Map(binding PersonBinding, source ConversationMessage, now time.Time) (IngestionBundle, error)
}

// RuntimeConfig owns product-specific source heuristics and planning defaults.
// It is intentionally outside the canonical kernel package.
type RuntimeConfig struct {
	SchedulingPriority kernel.PriorityFactors
	SchedulingEffort   kernel.EffortAttention
	RequestedScope     kernel.Scope
}

func DefaultRuntimeConfig() RuntimeConfig {
	return RuntimeConfig{
		SchedulingPriority: kernel.PriorityFactors{
			SubjectiveImportance: 0.50,
			ObjectiveStakes:      0.70,
			ExpectedImpact:       0.70,
			Reversibility:        0.80,
			Uncertainty:          0.30,
			OpportunityCost:      0.20,
			EffortAttentionCost:  0.20,
		},
		SchedulingEffort: kernel.EffortAttention{
			EstimatedEffort:    15 * time.Minute,
			EstimatedAttention: 10 * time.Minute,
			InterruptionCost:   0.10,
			ContextSwitchCost:  0.10,
		},
		RequestedScope: kernel.Scope{Capability: "create_draft", Resource: "email:client"},
	}
}

func (c RuntimeConfig) Validate() error {
	if c.SchedulingPriority.Validate() != nil || c.SchedulingEffort.Validate() != nil || c.RequestedScope.Capability == "" || c.RequestedScope.Resource == "" {
		return ErrInvalidRuntimeConfig
	}
	return nil
}

// RuntimeResult contains durable identifiers from a completed vertical slice
// and is returned unchanged by a deterministic idempotent replay.
type RuntimeResult struct {
	PersonID            string
	IdempotencyKey      string
	EventID             string
	EvidenceID          string
	MemoryID            string
	ClaimID             string
	IntentID            string
	OpenLoopID          string
	OpportunityID       string
	DecisionID          string
	ProposalID          string
	ActionGateID        string
	InteractionBlockID  string
	ThreadID            string
	ThreadStateID       string
	ContinuityLinkID    string
	AttunementEpisodeID string
	InterventionID      string
	AdaptationID        string
	Replayed            bool
}

// RuntimeState is a person-scoped retrieval view. It preserves the complete
// source-to-action lineage used by the local vertical-slice proof.
type RuntimeState struct {
	Person                    kernel.Person
	World                     kernel.PersonalWorld
	Sources                   []SourceRecord
	Events                    []kernel.Event
	Evidence                  []kernel.Evidence
	Memories                  []kernel.Memory
	Claims                    []kernel.Claim
	Goals                     []kernel.Goal
	Constraints               []kernel.Constraint
	PendingIntents            []kernel.PendingIntent
	OpenLoops                 []kernel.OpenLoop
	AttentionBudgets          []kernel.AttentionBudget
	Allocations               []kernel.AttentionAllocation
	Opportunities             []kernel.Opportunity
	Decisions                 []kernel.Decision
	Permissions               []kernel.Permission
	ActionProposals           []kernel.ActionProposal
	ActionGates               []kernel.ActionGate
	Outcomes                  []kernel.Outcome
	Audits                    []kernel.SelfAudit
	InteractionBlocks         []kernel.InteractionBlock
	Threads                   []kernel.Thread
	ContinuityLinks           []kernel.ContinuityLink
	ThreadDeltas              []kernel.ThreadDelta
	CurrentThreadStates       []kernel.CurrentThreadState
	ObservedSignals           []kernel.ObservedInteractionSignal
	InteractionBaselines      []kernel.PersonalInteractionBaseline
	InferredInteractionStates []kernel.InferredInteractionState
	AdaptationDecisions       []kernel.InteractionAdaptationDecision
	AttunementEpisodes        []kernel.AttunementEpisode
	InteractionInterventions  []kernel.InteractionIntervention
	InteractionOutcomes       []kernel.InteractionOutcome
	AttunementPatterns        []kernel.PersonalAttunementPattern
	Replays                   []RuntimeResult
}

// RuntimeRepository defines transactional, person-scoped persistence required
// by this slice. The in-memory implementation has identical contracts to a
// future Postgres adapter but is not a production persistence activation.
type RuntimeRepository interface {
	RunInTransaction(ctx context.Context, personID string, fn func(RuntimeTransaction) error) error
	ReadState(ctx context.Context, requesterPersonID, targetPersonID string) (RuntimeState, error)
}

// RuntimeTransaction exposes only append-preserving writes and bounded reads.
// Implementations must commit all writes atomically or discard them on error.
type RuntimeTransaction interface {
	FindReplay(idempotencyKey string) (RuntimeResult, bool)
	StoreSource(source SourceRecord) error
	SaveEvent(event kernel.Event) error
	SaveEvidence(evidence kernel.Evidence) error
	SaveMemory(memory kernel.Memory) error
	SaveClaim(claim kernel.Claim) error
	SaveMemoryEventLink(link kernel.MemoryEventLink) error
	SavePendingIntent(intent kernel.PendingIntent) error
	SaveOpenLoop(loop kernel.OpenLoop) error
	SaveOpportunity(opportunity kernel.Opportunity) error
	SaveDecision(decision kernel.Decision) error
	SaveActionProposal(proposal kernel.ActionProposal) error
	SaveActionGate(gate kernel.ActionGate) error
	SaveAttentionAllocation(allocation kernel.AttentionAllocation) error
	SaveInteractionBlock(block kernel.InteractionBlock) error
	SaveThread(thread kernel.Thread) error
	SaveContinuityLink(link kernel.ContinuityLink) error
	SaveThreadDelta(delta kernel.ThreadDelta) error
	SaveCurrentThreadState(state kernel.CurrentThreadState) error
	SaveObservedInteractionSignal(signal kernel.ObservedInteractionSignal) error
	SavePersonalInteractionBaseline(baseline kernel.PersonalInteractionBaseline) error
	SaveInferredInteractionState(state kernel.InferredInteractionState) error
	SaveInteractionAdaptationDecision(decision kernel.InteractionAdaptationDecision) error
	SaveAttunementEpisode(episode kernel.AttunementEpisode) error
	SaveInteractionIntervention(intervention kernel.InteractionIntervention) error
	SaveInteractionOutcome(outcome kernel.InteractionOutcome) error
	SavePersonalAttunementPattern(pattern kernel.PersonalAttunementPattern) error
	ListThreads() []kernel.Thread
	ListInteractionBlocks(threadID string) []kernel.InteractionBlock
	ListThreadDeltas(threadID string) []kernel.ThreadDelta
	ListCurrentThreadStates(threadID string) []kernel.CurrentThreadState
	ListAttunementPatterns(contextSignature string) []kernel.PersonalAttunementPattern
	ListActiveGoals(at time.Time) []kernel.Goal
	ListActivePermissions(at time.Time) []kernel.Permission
	CurrentAttentionBudget(at time.Time) (kernel.AttentionBudget, bool)
	StoreReplay(result RuntimeResult) error
}

// TemporalEvaluationPolicy adds the independent v0.2 temporal cognition
// contract to the policy required by this vertical slice.
type TemporalEvaluationPolicy interface {
	kernel.EvaluationPolicy
	EvaluateTemporal(input kernel.TemporalPriorityInput) (kernel.TemporalEvaluation, error)
}

// RuntimeService is the explicit composition root. It is not an HTTP handler.
// ContinuityInput is a synthetic/local source-to-continuity request. The
// supplied source remains provenance; it cannot select a Person or bypass the
// server-resolved profile binding. No raw audio or external account read occurs.
type ContinuityInput struct {
	IdempotencyKey    string
	Source            ConversationMessage
	Block             kernel.InteractionBlock
	Triggers          []kernel.SemanticTrigger
	ProposedThread    *kernel.Thread
	Deltas            []kernel.ThreadDelta
	Signals           []kernel.ObservedInteractionSignal
	Baseline          *kernel.PersonalInteractionBaseline
	AttunementControl kernel.AttunementControlMode
	ContextSignature  string
}

type ContinuityRuntimeResult struct {
	PersonID            string
	IdempotencyKey      string
	BlockID             string
	ThreadID            string
	ThreadStateID       string
	ContinuityLinkID    string
	AttunementEpisodeID string
	InterventionID      string
	AdaptationID        string
	Replayed            bool
}

// RuntimeService is the existing v0.3/v0.4 scheduling slice composition root.
type RuntimeService struct {
	Feature              Feature
	Activation           RuntimeActivation
	Consent              DerivedMemoryConsentVerifier
	Identity             IdentityResolver
	Adapter              ConversationMessageAdapter
	Repo                 RuntimeRepository
	Clock                kernel.Clock
	Policy               TemporalEvaluationPolicy
	Config               RuntimeConfig
	BoundaryPolicy       kernel.InteractionBoundaryPolicy
	ThreadResolver       kernel.SemanticThreadResolver
	RetrievalDepthPolicy kernel.RetrievalDepthPolicy
	AttunementPolicy     kernel.AttunementSafetyPolicy
}
