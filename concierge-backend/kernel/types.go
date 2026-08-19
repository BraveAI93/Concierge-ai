// Package kernel contains the pure, person-centric domain model for the
// Central Intelligence Kernel. It intentionally has no HTTP, database, or
// external-service dependencies.
package kernel

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrMissingIdentifier       = errors.New("kernel: identifier is required")
	ErrPersonBoundary          = errors.New("kernel: records belong to different people")
	ErrInvalidTemporalState    = errors.New("kernel: invalid temporal state")
	ErrInvalidIntentTransition = errors.New("kernel: invalid pending intent transition")
	ErrInvalidGateTransition   = errors.New("kernel: invalid action gate transition")
	ErrMissingApproval         = errors.New("kernel: action requires approval")
	ErrPermissionDenied        = errors.New("kernel: permission scope does not allow action")
	ErrInactiveTemporalState   = errors.New("kernel: temporal state is inactive")
	ErrMissingOutcomeEvidence  = errors.New("kernel: outcome requires evidence")
)

// Person is the owner of a PersonalWorld. PersonID is required on all
// person-scoped records so adapters can enforce isolation before persistence.
type Person struct {
	ID          string
	DisplayName string
	WorldID     string
	CreatedAt   time.Time
}

// PersonalWorld is the canonical, person-owned aggregate boundary. Reference
// slices deliberately avoid embedding mutable child records in the root.
type PersonalWorld struct {
	ID         string
	PersonID   string
	EntityIDs  []string
	ContextIDs []string
	UpdatedAt  time.Time
}

// EntityKind distinguishes people, organisations, places, and user-defined
// concepts without requiring a legacy profile schema.
type EntityKind string

const (
	EntityPerson       EntityKind = "person"
	EntityOrganization EntityKind = "organization"
	EntityPlace        EntityKind = "place"
	EntityConcept      EntityKind = "concept"
	EntityObject       EntityKind = "object"
)

// Entity represents a named item in a personal world. Aliases resolve
// alternate spellings and references into the same canonical entity.
type Entity struct {
	ID        string
	PersonID  string
	Kind      EntityKind
	Name      string
	Aliases   []Alias
	Metadata  map[string]string
	CreatedAt time.Time
}

// Alias binds a normalized user-facing label to an Entity. The adapter owns
// normalization policy; the domain model records the resulting value.
type Alias struct {
	ID         string
	EntityID   string
	Value      string
	Normalized string
	Source     string
	CreatedAt  time.Time
}

// Context is the current or historical situational frame for a person.
type Context struct {
	ID         string
	PersonID   string
	Kind       string
	Label      string
	EntityIDs  []string
	Attributes map[string]string
	Temporal   TemporalState
	CreatedAt  time.Time
}

// Event is an observed occurrence. EventAt is held in Temporal rather than
// duplicated so all time policy uses the same four-dimensional representation.
type Event struct {
	ID          string
	PersonID    string
	Kind        string
	Summary     string
	EntityIDs   []string
	ContextIDs  []string
	MemoryIDs   []string
	Temporal    TemporalState
	EvidenceIDs []string
	CreatedAt   time.Time
}

// MemoryKind distinguishes factual, procedural, preference, and episodic
// records, while Claim states an assessable proposition held by a memory.
type MemoryKind string

const (
	MemoryEpisodic   MemoryKind = "episodic"
	MemorySemantic   MemoryKind = "semantic"
	MemoryPreference MemoryKind = "preference"
	MemoryProcedure  MemoryKind = "procedure"
)

type Memory struct {
	ID         string
	PersonID   string
	Kind       MemoryKind
	Summary    string
	ClaimIDs   []string
	EventIDs   []string
	ContextIDs []string
	Temporal   TemporalState
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Claim struct {
	ID           string
	PersonID     string
	MemoryID     string
	Statement    string
	SubjectID    string
	Predicate    string
	Object       string
	EvidenceIDs  []string
	Confidence   ConfidenceAssessment
	Temporal     TemporalState
	Freshness    FreshnessState
	Lineage      ClaimLineage
	SupersedesID string
	CreatedAt    time.Time
}

// EvidenceStance declares whether a piece of evidence supports or contradicts
// a claim. Ambiguous evidence is intentionally not silently treated as support.
type EvidenceStance string

const (
	EvidenceSupports    EvidenceStance = "supports"
	EvidenceContradicts EvidenceStance = "contradicts"
	EvidenceNeutral     EvidenceStance = "neutral"
)

type Evidence struct {
	ID         string
	PersonID   string
	ClaimID    string
	Stance     EvidenceStance
	Summary    string
	Quality    float64
	Relevance  float64
	Authority  float64
	Provenance Provenance
	Temporal   TemporalState
	CreatedAt  time.Time
}

// Provenance is immutable attribution metadata. It is attached to every
// Evidence record, and it is repeated as evidence IDs in outcomes for audit.
type Provenance struct {
	SourceType string
	SourceRef  string
	Actor      string
	CapturedAt time.Time
	Checksum   string
}

// Goal expresses an active or retired desired state.
type GoalStatus string

const (
	GoalActive    GoalStatus = "active"
	GoalPaused    GoalStatus = "paused"
	GoalAchieved  GoalStatus = "achieved"
	GoalAbandoned GoalStatus = "abandoned"
)

type Goal struct {
	ID                   string
	PersonID             string
	Title                string
	Description          string
	Importance           float64 // v0.1 compatibility field; policy must not infer objective stakes from it.
	SubjectiveImportance float64
	Status               GoalStatus
	Temporal             TemporalState
	CreatedAt            time.Time
}

// Constraint may prevent action altogether (hard) or reduce its utility (soft).
type ConstraintKind string

const (
	ConstraintHard ConstraintKind = "hard"
	ConstraintSoft ConstraintKind = "soft"
)

type Constraint struct {
	ID          string
	PersonID    string
	Kind        ConstraintKind
	Title       string
	Description string
	Active      bool
	Temporal    TemporalState
	CreatedAt   time.Time
}

// PendingIntentState is the deterministic lifecycle of an unfinished intention.
type PendingIntentState string

const (
	IntentCaptured   PendingIntentState = "captured"
	IntentClarifying PendingIntentState = "clarifying"
	IntentReady      PendingIntentState = "ready"
	IntentProposed   PendingIntentState = "proposed"
	IntentInProgress PendingIntentState = "in_progress"
	IntentCompleted  PendingIntentState = "completed"
	IntentCancelled  PendingIntentState = "cancelled"
	IntentExpired    PendingIntentState = "expired"
)

type PendingIntent struct {
	ID          string
	PersonID    string
	Summary     string
	State       PendingIntentState
	GoalID      string
	MemoryID    string
	ContextIDs  []string
	Temporal    TemporalState
	Transitions []IntentTransition
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type IntentTransition struct {
	From   PendingIntentState
	To     PendingIntentState
	At     time.Time
	Actor  string
	Reason string
}

// OpenLoop represents the attention-facing aspect of a pending intention.
type OpenLoop struct {
	ID              string
	PersonID        string
	PendingIntentID string
	Label           string
	Attention       TemporalState
	InteractionGap  InteractionGapState
	AttentionNeed   EffortAttention
	ContextIDs      []string
	EntityIDs       []string
	ResolvedAt      *time.Time
	CreatedAt       time.Time
}

// TemporalState intentionally separates four timelines: occurrence,
// recording, applicability, and attention. Expiry limits applicability.
type TemporalState struct {
	EventAt     time.Time
	RecordedAt  time.Time
	EffectiveAt time.Time
	AttentionAt time.Time
	ExpiresAt   *time.Time
}

// Opportunity is a candidate intervention. It contains scored inputs, not an
// execution instruction. The evaluation output is captured in Evaluation.
type Opportunity struct {
	ID                 string
	PersonID           string
	Title              string
	Summary            string
	GoalIDs            []string
	ConstraintIDs      []string
	EvidenceIDs        []string
	Temporal           TemporalState
	Priority           PriorityFactors
	ActionWindow       ActionWindow
	AttentionNeed      EffortAttention
	GoalAlignment      float64 // v0.1 compatibility field; v0.2 policy uses Priority factors.
	ExpectedValue      float64 // v0.1 compatibility field; not an objective-stakes substitute.
	Effort             float64 // v0.1 compatibility field; not an attention-duration substitute.
	Risk               float64 // v0.1 compatibility field; not a reversibility substitute.
	EvidenceConfidence float64 // v0.1 compatibility field; v0.2 uses explicit uncertainty.
	TemporalPriority   float64 // v0.1 compatibility field; policy may use deadline feasibility instead.
	Evaluation         OpportunityEvaluation
	CreatedAt          time.Time
}

type OpportunityEvaluation struct {
	Utility       float64
	HardBlocked   bool
	SoftPenalty   float64 // v0.1 compatibility output; v0.2 policy owns any penalty.
	Deadline      DeadlineFeasibility
	Mismatch      PriorityMismatch
	DecisionBasis string
	EvaluatedAt   time.Time
}

// Decision records a deliberative result independently of the action proposal.
type DecisionKind string

const (
	DecisionRecommend DecisionKind = "recommend"
	DecisionDefer     DecisionKind = "defer"
	DecisionDecline   DecisionKind = "decline"
	DecisionSurface   DecisionKind = "must_surface"
)

type Decision struct {
	ID            string
	PersonID      string
	OpportunityID string
	Kind          DecisionKind
	Utility       float64
	Reason        string
	CreatedAt     time.Time
}

// Scope is the least-privilege operation boundary for a Permission and an
// ActionProposal. A star is allowed only as an explicit literal wildcard.
type Scope struct {
	Capability string
	Resource   string
}

type Permission struct {
	ID               string
	PersonID         string
	Scopes           []Scope
	RequiresApproval bool
	CanAutoApprove   bool
	Temporal         TemporalState
	GrantedBy        string
	CreatedAt        time.Time
}

// ActionProposal must be approved by its ActionGate before an adapter performs
// the requested capability against its requested resource.
type ActionProposal struct {
	ID            string
	PersonID      string
	OpportunityID string
	DecisionID    string
	Title         string
	Requested     Scope
	PermissionID  string
	Parameters    map[string]string
	CreatedAt     time.Time
}

type ActionGateState string

const (
	GateDraft            ActionGateState = "draft"
	GateAwaitingApproval ActionGateState = "awaiting_approval"
	GateApproved         ActionGateState = "approved"
	GateRejected         ActionGateState = "rejected"
	GateExpired          ActionGateState = "expired"
	GateExecuted         ActionGateState = "executed"
)

type ActionGate struct {
	ID               string
	PersonID         string
	ActionProposalID string
	State            ActionGateState
	Transitions      []GateTransition
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type GateTransition struct {
	From   ActionGateState
	To     ActionGateState
	At     time.Time
	Actor  string
	Reason string
}

// Outcome records the observed result of an approved proposal; it is never an
// implicit success inferred from an executed action.
type OutcomeStatus string

const (
	OutcomeSucceeded OutcomeStatus = "succeeded"
	OutcomePartial   OutcomeStatus = "partial"
	OutcomeFailed    OutcomeStatus = "failed"
	OutcomeCancelled OutcomeStatus = "cancelled"
)

type Outcome struct {
	ID               string
	PersonID         string
	ActionProposalID string
	Status           OutcomeStatus
	Summary          string
	ExpectedUtility  float64
	ObservedUtility  float64
	EvidenceIDs      []string
	Provenance       []Provenance
	OccurredAt       time.Time
	CreatedAt        time.Time
}

type AuditStatus string

const (
	AuditConfirmed AuditStatus = "confirmed"
	AuditAdjusted  AuditStatus = "adjusted"
	AuditLearned   AuditStatus = "learned"
	AuditEscalate  AuditStatus = "escalate"
)

type SelfAudit struct {
	ID                 string
	PersonID           string
	OutcomeID          string
	Status             AuditStatus
	ExpectedUtility    float64
	ObservedUtility    float64
	UtilityDelta       float64
	EvidenceSufficient bool
	Summary            string
	CreatedAt          time.Time
}

// MemoryEventLink is a reciprocal association with explicit linking semantics.
type MemoryEventLink struct {
	MemoryID string
	EventID  string
	Reason   string
	LinkedAt time.Time
}

// ConfidenceAssessment preserves the components used to derive a claim score.
type ConfidenceAssessment struct {
	Score             float64
	SupportingWeight  float64
	ConflictingWeight float64
	EvidenceCount     int
	ProvenanceCount   int
	AssessedAt        time.Time
}

// TemporalEvaluation explains why something should receive attention now.
type TemporalEvaluation struct {
	Active          bool
	Importance      float64
	DeadlineUrgency float64
	Recency         float64
	InteractionGap  time.Duration
	AttentionDue    float64
	Utility         float64
	EvaluatedAt     time.Time
}

func (t TemporalState) Validate() error {
	if t.EventAt.IsZero() || t.RecordedAt.IsZero() || t.EffectiveAt.IsZero() || t.AttentionAt.IsZero() {
		return ErrInvalidTemporalState
	}
	// RecordedAt is knowledge/ingestion time, not semantic event time. A future
	// appointment may therefore be validly recorded before EventAt occurs.
	if t.ExpiresAt != nil && t.ExpiresAt.Before(t.EffectiveAt) {
		return ErrInvalidTemporalState
	}
	return nil
}

func (t TemporalState) IsActive(at time.Time) bool {
	if t.Validate() != nil || at.Before(t.EffectiveAt) {
		return false
	}
	return t.ExpiresAt == nil || !at.After(*t.ExpiresAt)
}

func (s Scope) Matches(request Scope) bool {
	capabilityMatches := s.Capability == "*" || s.Capability == request.Capability
	resourceMatches := s.Resource == "*" || s.Resource == request.Resource
	return capabilityMatches && resourceMatches
}

func normalizeText(value string) string {
	return strings.TrimSpace(value)
}

func isZeroOrUnit(value float64) bool {
	return value >= 0 && value <= 1
}

func clampUnit(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}
