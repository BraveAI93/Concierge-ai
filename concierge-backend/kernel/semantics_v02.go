package kernel

import (
	"errors"
	"time"
)

var (
	ErrInvalidInteractionGap  = errors.New("kernel: invalid interaction gap state")
	ErrInvalidAttentionBudget = errors.New("kernel: invalid attention budget")
	ErrInvalidAttentionItem   = errors.New("kernel: invalid attention item")
	ErrInvalidPriorityFactors = errors.New("kernel: invalid priority factors")
	ErrInvalidPolicy          = errors.New("kernel: invalid evaluation policy")
	ErrInfeasibleActionWindow = errors.New("kernel: action window is infeasible")
	ErrInvalidFreshness       = errors.New("kernel: invalid freshness state")
	ErrSupersessionDenied     = errors.New("kernel: supersession policy was not satisfied")
)

// EvaluationMoment is the independent wall-clock time at which the kernel
// evaluates state. It is intentionally supplied to deterministic primitives
// rather than stored as if it were semantic event time.
type EvaluationMoment struct {
	WallClockAt time.Time
}

func (m EvaluationMoment) Validate() error {
	if m.WallClockAt.IsZero() {
		return ErrInvalidTemporalState
	}
	return nil
}

// InteractionGapState describes time since the last relevant interaction. It
// is distinct from EventAt: an old event may have a recent interaction, and a
// recent event may follow a long gap.
type InteractionGapState struct {
	PersonID          string
	EntityID          string
	ContextID         string
	LastInteractionAt time.Time
	ObservedAt        time.Time
}

func (g InteractionGapState) Validate() error {
	if normalizeText(g.PersonID) == "" || g.LastInteractionAt.IsZero() || g.ObservedAt.IsZero() {
		return ErrInvalidInteractionGap
	}
	return nil
}

func (g InteractionGapState) Elapsed(at time.Time) (time.Duration, error) {
	if err := g.Validate(); err != nil || at.IsZero() || at.Before(g.LastInteractionAt) {
		return 0, ErrInvalidInteractionGap
	}
	return at.Sub(g.LastInteractionAt), nil
}

// EffortAttention separates total completion effort from the attention time
// that must be allocated. Both values are durations, not unitless scores.
type EffortAttention struct {
	EstimatedEffort    time.Duration
	EstimatedAttention time.Duration
	InterruptionCost   float64
	ContextSwitchCost  float64
}

func (e EffortAttention) Validate() error {
	if e.EstimatedEffort < 0 || e.EstimatedAttention < 0 || !isZeroOrUnit(e.InterruptionCost) || !isZeroOrUnit(e.ContextSwitchCost) {
		return ErrInvalidAttentionItem
	}
	return nil
}

func (e EffortAttention) RequiredAttention() time.Duration {
	if e.EstimatedAttention > 0 {
		return e.EstimatedAttention
	}
	return e.EstimatedEffort
}

// PriorityFactors preserves independent input dimensions. All values use a
// [0,1] scale but are not collapsed into generic importance, value, or risk.
// Reversibility is the ease of reversal: zero means irreversible, one means
// easily reversed. The policy chooses how every factor influences a decision.
type PriorityFactors struct {
	SubjectiveImportance float64
	ObjectiveStakes      float64
	ExpectedImpact       float64
	Reversibility        float64
	Uncertainty          float64
	OpportunityCost      float64
	EffortAttentionCost  float64
}

func (f PriorityFactors) Validate() error {
	for _, value := range []float64{
		f.SubjectiveImportance,
		f.ObjectiveStakes,
		f.ExpectedImpact,
		f.Reversibility,
		f.Uncertainty,
		f.OpportunityCost,
		f.EffortAttentionCost,
	} {
		if !isZeroOrUnit(value) {
			return ErrInvalidPriorityFactors
		}
	}
	return nil
}

// ActionWindow distinguishes a latest safe action time from a record expiry.
// LatestSafeActionAt answers when work must begin or complete; ExpiresAt in a
// TemporalState continues to describe applicability/validity only.
type ActionWindow struct {
	LatestSafeActionAt time.Time
	EstimatedEffort    time.Duration
}

func (w ActionWindow) HasDeadline() bool {
	return !w.LatestSafeActionAt.IsZero()
}

func (w ActionWindow) Validate() error {
	if w.EstimatedEffort < 0 {
		return ErrInfeasibleActionWindow
	}
	return nil
}

type DeadlineFeasibilityState string

const (
	DeadlineNone       DeadlineFeasibilityState = "no_deadline"
	DeadlineFeasible   DeadlineFeasibilityState = "feasible"
	DeadlineInfeasible DeadlineFeasibilityState = "infeasible"
)

// DeadlineFeasibility preserves the difference between an urgent task that is
// feasible and one whose action window is already impossible.
type DeadlineFeasibility struct {
	State           DeadlineFeasibilityState
	LatestSafeAt    time.Time
	Remaining       time.Duration
	EstimatedEffort time.Duration
	Reason          string
}

// PriorityMismatch makes a high-objective-stakes/low-subjective-importance
// conflict explicit. It requests surfacing without silently changing goals.
type PriorityMismatch struct {
	MustSurface            bool
	PreservesGoalAuthority bool
	SubjectiveImportance   float64
	ObjectiveStakes        float64
	Reason                 string
}

// AttentionContext is the current context against which open loops are
// selectively resurfaced. Entity and context matches are evaluated separately.
type AttentionContext struct {
	ContextIDs []string
	EntityIDs  []string
}

// AttentionBudget is a finite, person-scoped allocation of attention for a
// window. Capacity and MaxCompetingItems are both hard constraints.
type AttentionBudget struct {
	ID                string
	PersonID          string
	WindowStart       time.Time
	WindowEnd         time.Time
	AttentionCapacity time.Duration
	MaxCompetingItems int
	InterruptionCost  float64
	CurrentContext    AttentionContext
}

func (b AttentionBudget) Validate() error {
	if normalizeText(b.ID) == "" || normalizeText(b.PersonID) == "" || b.WindowStart.IsZero() || b.WindowEnd.IsZero() || !b.WindowEnd.After(b.WindowStart) || b.AttentionCapacity < 0 || b.AttentionCapacity > b.WindowEnd.Sub(b.WindowStart) || b.MaxCompetingItems <= 0 || !isZeroOrUnit(b.InterruptionCost) {
		return ErrInvalidAttentionBudget
	}
	return nil
}

// AttentionCandidate is an unresolved OpenLoop plus the policy inputs needed
// to select or defer it within a finite budget.
type AttentionCandidate struct {
	OpenLoopID   string
	PersonID     string
	Label        string
	ResolvedAt   *time.Time
	ContextIDs   []string
	EntityIDs    []string
	BasePriority float64
	Effort       EffortAttention
}

func (c AttentionCandidate) Validate() error {
	if normalizeText(c.OpenLoopID) == "" || normalizeText(c.PersonID) == "" || !isZeroOrUnit(c.BasePriority) || c.Effort.Validate() != nil {
		return ErrInvalidAttentionItem
	}
	return nil
}

type AttentionSelection struct {
	OpenLoopID        string
	Score             float64
	ContextMatched    bool
	AttentionReserved time.Duration
	Reason            string
}

type AttentionDeferral struct {
	OpenLoopID string
	Reason     string
}

// AttentionAllocation reports both surfaced and deferred loops so the
// scheduler cannot hide budget trade-offs behind a score alone.
type AttentionAllocation struct {
	BudgetID          string
	PersonID          string
	Surfaced          []AttentionSelection
	Deferred          []AttentionDeferral
	UsedAttention     time.Duration
	RemainingCapacity time.Duration
}

// FreshnessStatus separates current, stale, historical, and superseded claims.
type FreshnessStatus string

const (
	FreshnessFresh      FreshnessStatus = "fresh"
	FreshnessStale      FreshnessStatus = "stale"
	FreshnessHistorical FreshnessStatus = "historical"
	FreshnessSuperseded FreshnessStatus = "superseded"
)

// FreshnessState records staleness and revalidation without deleting prior
// observations. StaleAfter is a policy-independent record property; policy
// decides how stale claims may be used.
type FreshnessState struct {
	StaleAfter        time.Duration
	LastValidatedAt   time.Time
	LastRevalidatedAt time.Time
	Status            FreshnessStatus
}

func (f FreshnessState) Validate() error {
	if f.StaleAfter < 0 || f.LastValidatedAt.IsZero() {
		return ErrInvalidFreshness
	}
	if !f.LastRevalidatedAt.IsZero() && f.LastRevalidatedAt.Before(f.LastValidatedAt) {
		return ErrInvalidFreshness
	}
	return nil
}

type ClaimLineage struct {
	ClaimID           string
	SupersedesClaimID string
	EvidenceIDs       []string
	PreservesHistory  bool
	RecordedAt        time.Time
}

type ClaimSelection struct {
	CurrentClaimID     string
	HistoricalClaimIDs []string
	Reason             string
}

// SupersessionPolicy is injected into the primitive that may promote a newer
// contradictory claim. Newness alone is deliberately never sufficient.
type SupersessionPolicy struct {
	MinimumAuthority     float64
	MinimumRelevance     float64
	RequireContradiction bool
}

func (p SupersessionPolicy) Validate() error {
	if !isZeroOrUnit(p.MinimumAuthority) || !isZeroOrUnit(p.MinimumRelevance) {
		return ErrInvalidPolicy
	}
	return nil
}
