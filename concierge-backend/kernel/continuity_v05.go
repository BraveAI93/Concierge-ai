package kernel

import (
	"errors"
	"time"
)

var (
	ErrInvalidInteractionBlock    = errors.New("kernel: invalid interaction block")
	ErrInvalidThread              = errors.New("kernel: invalid conversational thread")
	ErrInvalidContinuityLink      = errors.New("kernel: invalid continuity link")
	ErrInvalidThreadDelta         = errors.New("kernel: invalid thread delta")
	ErrInvalidRetrievalRequest    = errors.New("kernel: invalid retrieval request")
	ErrInvalidInteractionSignal   = errors.New("kernel: invalid observed interaction signal")
	ErrInvalidInteractionBaseline = errors.New("kernel: invalid personal interaction baseline")
	ErrInvalidInteractionInference = errors.New("kernel: invalid inferred interaction state")
	ErrUnsafeAdaptation           = errors.New("kernel: unsafe or disallowed interaction adaptation")
	ErrAttunementDisabled         = errors.New("kernel: attunement is disabled by user control")
	ErrInvalidAttunementPattern   = errors.New("kernel: invalid attunement pattern")
)

// PrivacyClass categorizes v0.5 records for future consent, retention, and
// access-control policy. It records no entitlement by itself.
type PrivacyClass string

const (
	PrivacyRawCommunicationSignal PrivacyClass = "raw_communication_signal"
	PrivacyDerivedBaseline        PrivacyClass = "derived_baseline"
	PrivacyInferredState          PrivacyClass = "inferred_interaction_state"
	PrivacyAdaptationDecision     PrivacyClass = "adaptation_decision"
	PrivacyOutcomeEvidence        PrivacyClass = "outcome_evidence"
	PrivacyLearnedPattern         PrivacyClass = "learned_pattern"
)

// InteractionSourceRef is provenance only. UI conversation/session identifiers
// must not become a canonical Block or Thread identity.
type InteractionSourceRef struct {
	SourceType     string
	ConversationID string
	SessionID      string
	MessageID      string
	DeviceClass    string
	CapturedAt     time.Time
}

type BlockBoundaryReason string

const (
	BoundaryInactivity          BlockBoundaryReason = "meaningful_inactivity"
	BoundaryExplicitEnding      BlockBoundaryReason = "explicit_ending"
	BoundaryTopicTransition     BlockBoundaryReason = "major_topic_transition"
	BoundaryContextChange       BlockBoundaryReason = "context_source_change"
	BoundaryTaskCompletion      BlockBoundaryReason = "bounded_task_completion"
	BoundarySemanticDiscontinuity BlockBoundaryReason = "semantic_discontinuity"
)

type BoundaryEvidence struct {
	Reason     BlockBoundaryReason
	Summary    string
	EvidenceID string
	Confidence float64
	ObservedAt time.Time
}

// InteractionBlock is the canonical semantic/time-bounded interaction segment.
// It is deliberately not a physical chat or session.
type InteractionBlock struct {
	ID                string
	PersonID          string
	SourceRefs        []InteractionSourceRef
	StartTemporal     TemporalState
	EndEventAt        *time.Time
	EvaluatedAt       time.Time
	IngestedAt        time.Time
	SourceType        string
	TopicLabels       []string
	EntityIDs         []string
	ContextIDs        []string
	GoalIDs           []string
	DecisionIDs       []string
	OpenLoopIDs       []string
	PendingIntentIDs  []string
	EvidenceIDs       []string
	OutcomeIDs        []string
	Importance        PriorityFactors
	SemanticState     string
	Provenance        []Provenance
	Confidence        float64
	Freshness         FreshnessState
	BoundaryEvidence  []BoundaryEvidence
	SupersedesBlockID string
	ThreadIDs         []string
}

func (b InteractionBlock) Validate() error {
	if normalizeText(b.ID) == "" || normalizeText(b.PersonID) == "" || b.StartTemporal.Validate() != nil || b.EvaluatedAt.IsZero() || b.IngestedAt.IsZero() || !isZeroOrUnit(b.Confidence) || b.Importance.Validate() != nil {
		return ErrInvalidInteractionBlock
	}
	if b.EndEventAt != nil && b.EndEventAt.Before(b.StartTemporal.EventAt) {
		return ErrInvalidInteractionBlock
	}
	if b.Freshness.Validate() != nil {
		return ErrInvalidInteractionBlock
	}
	for _, evidence := range b.BoundaryEvidence {
		if evidence.Reason == "" || evidence.ObservedAt.IsZero() || !isZeroOrUnit(evidence.Confidence) {
			return ErrInvalidInteractionBlock
		}
	}
	return nil
}

// InteractionBoundaryAssessment records an evaluated decision, not an implicit
// hard-coded time rule.
type InteractionBoundaryAssessment struct {
	CloseBlock      bool
	StartNewBlock   bool
	Reasons         []BlockBoundaryReason
	Evidence        []BoundaryEvidence
	Confidence      float64
	InteractionGap  time.Duration
	EvaluatedAt     time.Time
}

// InteractionBoundaryPolicy is injected by composition. Implementations may
// use richer local semantic analysis later without changing canonical records.
type InteractionBoundaryPolicy interface {
	AssessBoundary(previous InteractionBlock, incoming InteractionBlock, gap InteractionGapState, at EvaluationMoment) (InteractionBoundaryAssessment, error)
}

type ThreadStatus string

const (
	ThreadActive   ThreadStatus = "active"
	ThreadDormant  ThreadStatus = "dormant"
	ThreadResolved ThreadStatus = "resolved"
)

type ThreadAnchor struct {
	Kind  string
	ID    string
	Label string
}

// Thread organizes source-independent continuity. It never replaces the
// source Event/Memory/Evidence records that support it.
type Thread struct {
	ID                   string
	PersonID             string
	Anchors              []ThreadAnchor
	Aliases              []string
	CurrentStateID       string
	Status               ThreadStatus
	Importance           PriorityFactors
	OpenLoopIDs          []string
	FirstRelevantAt      time.Time
	MostRecentRelevantAt time.Time
	LastMaterialUpdateAt time.Time
	EvidenceIDs          []string
	LinkedThreadIDs      []string
	Confidence           float64
	Freshness            FreshnessState
	RetrievalPolicy      string
	ContinuitySummary    string
	CreatedAt            time.Time
}

func (t Thread) Validate() error {
	if normalizeText(t.ID) == "" || normalizeText(t.PersonID) == "" || t.Status == "" || t.FirstRelevantAt.IsZero() || t.MostRecentRelevantAt.IsZero() || t.CreatedAt.IsZero() || !isZeroOrUnit(t.Confidence) || t.Importance.Validate() != nil || t.Freshness.Validate() != nil {
		return ErrInvalidThread
	}
	if t.MostRecentRelevantAt.Before(t.FirstRelevantAt) || (!t.LastMaterialUpdateAt.IsZero() && t.LastMaterialUpdateAt.Before(t.FirstRelevantAt)) {
		return ErrInvalidThread
	}
	return nil
}

type ContinuityObjectKind string

const (
	ContinuityBlock       ContinuityObjectKind = "interaction_block"
	ContinuityThread      ContinuityObjectKind = "thread"
	ContinuityEvent       ContinuityObjectKind = "event"
	ContinuityMemory      ContinuityObjectKind = "memory"
	ContinuityOpenLoop    ContinuityObjectKind = "open_loop"
	ContinuityIntent      ContinuityObjectKind = "pending_intent"
	ContinuityDecision    ContinuityObjectKind = "decision"
	ContinuityOutcome     ContinuityObjectKind = "outcome"
	ContinuityAttunement  ContinuityObjectKind = "attunement_episode"
)

type ContinuityRef struct {
	Kind ContinuityObjectKind
	ID   string
}

type ContinuityRelation string

const (
	ContinuitySameSubject        ContinuityRelation = "same_subject"
	ContinuityUpdates            ContinuityRelation = "updates"
	ContinuityDependsOn          ContinuityRelation = "depends_on"
	ContinuityContradicts        ContinuityRelation = "contradicts"
	ContinuityResolves           ContinuityRelation = "resolves"
	ContinuityCausesFollowup     ContinuityRelation = "causes_followup_for"
	ContinuityRelatedEntity      ContinuityRelation = "related_entity"
	ContinuityContextualReference ContinuityRelation = "contextual_reference"
	ContinuitySupersedes         ContinuityRelation = "supersedes"
	ContinuityMateriallyImpacts  ContinuityRelation = "materially_impacts"
)

// ContinuityLink provides a typed, evidence-backed graph edge. It is never a
// generic related=true flag.
type ContinuityLink struct {
	ID            string
	PersonID      string
	Source        ContinuityRef
	Target        ContinuityRef
	Relation      ContinuityRelation
	Why           string
	EvidenceIDs   []string
	Provenance    []Provenance
	Confidence    float64
	Temporal      TemporalState
	Freshness     FreshnessState
	SupersedesID  string
	CreatedAt     time.Time
}

func (l ContinuityLink) Validate() error {
	if normalizeText(l.ID) == "" || normalizeText(l.PersonID) == "" || l.Source.Kind == "" || normalizeText(l.Source.ID) == "" || l.Target.Kind == "" || normalizeText(l.Target.ID) == "" || l.Relation == "" || normalizeText(l.Why) == "" || !isZeroOrUnit(l.Confidence) || l.Temporal.Validate() != nil || l.Freshness.Validate() != nil || l.CreatedAt.IsZero() {
		return ErrInvalidContinuityLink
	}
	return nil
}

type ThreadDeltaEffect string

const (
	DeltaUpdatesOpenLoop    ThreadDeltaEffect = "updates_open_loop"
	DeltaUpdatesDecision    ThreadDeltaEffect = "updates_decision"
	DeltaUpdatesDeadline    ThreadDeltaEffect = "updates_deadline"
	DeltaUpdatesConstraint  ThreadDeltaEffect = "updates_constraint"
	DeltaUpdatesOpportunity ThreadDeltaEffect = "updates_opportunity"
	DeltaUpdatesOutcome     ThreadDeltaEffect = "updates_outcome"
	DeltaUpdatesState       ThreadDeltaEffect = "updates_state"
)

// ThreadDelta is an append-preserving material update. FieldChanges is a
// compact typed state projection; original source evidence remains available.
type ThreadDelta struct {
	ID                  string
	PersonID            string
	TargetThreadID      string
	Originating         ContinuityRef
	SemanticChange      string
	AffectedConcept     string
	FieldChanges        map[string]string
	Effects             []ThreadDeltaEffect
	EvidenceIDs         []string
	Provenance          []Provenance
	Confidence          float64
	Importance          PriorityFactors
	EventAt             time.Time
	EvaluatedAt         time.Time
	SupersedesDeltaID   string
	CreatedAt           time.Time
}

func (d ThreadDelta) Validate() error {
	if normalizeText(d.ID) == "" || normalizeText(d.PersonID) == "" || normalizeText(d.TargetThreadID) == "" || d.Originating.Kind == "" || normalizeText(d.Originating.ID) == "" || normalizeText(d.SemanticChange) == "" || d.EventAt.IsZero() || d.EvaluatedAt.IsZero() || d.CreatedAt.IsZero() || !isZeroOrUnit(d.Confidence) || d.Importance.Validate() != nil {
		return ErrInvalidThreadDelta
	}
	return nil
}

// CurrentThreadState is a materialized/reconstructable view. It records the
// baseline and delta lineage rather than discarding older reasoning.
type CurrentThreadState struct {
	ID                 string
	PersonID           string
	ThreadID           string
	BaselineSummary    string
	CurrentSummary     string
	FieldValues        map[string]string
	IncludedDeltaIDs   []string
	HistoricalDeltaIDs []string
	EvidenceIDs        []string
	ReconstructedAt    time.Time
	Freshness          FreshnessState
}

func (s CurrentThreadState) Validate() error {
	if normalizeText(s.ID) == "" || normalizeText(s.PersonID) == "" || normalizeText(s.ThreadID) == "" || s.ReconstructedAt.IsZero() || s.Freshness.Validate() != nil {
		return ErrInvalidThread
	}
	return nil
}

type SemanticTriggerKind string

const (
	TriggerExactWord  SemanticTriggerKind = "exact_word"
	TriggerAlias      SemanticTriggerKind = "alias"
	TriggerEntity     SemanticTriggerKind = "entity"
	TriggerProject    SemanticTriggerKind = "project"
	TriggerPlace      SemanticTriggerKind = "place"
	TriggerEvent      SemanticTriggerKind = "event"
	TriggerGoal       SemanticTriggerKind = "goal"
	TriggerDeadline   SemanticTriggerKind = "deadline"
	TriggerDecision   SemanticTriggerKind = "decision"
	TriggerOpenLoop   SemanticTriggerKind = "open_loop"
	TriggerParaphrase SemanticTriggerKind = "paraphrase"
	TriggerRelation   SemanticTriggerKind = "relation"
)

type SemanticTrigger struct {
	Kind        SemanticTriggerKind
	Value       string
	Normalized  string
	EntityIDs   []string
	ContextIDs  []string
	EvidenceIDs []string
	Temporal    *TemporalState
	Confidence  float64
}

type ThreadCandidate struct {
	ThreadID      string
	Score         float64
	MatchedBy     []SemanticTriggerKind
	Reason        string
	Confidence    float64
	Unresolved    bool
}

type ThreadResolution struct {
	PersonID      string
	SelectedID    string
	Candidates    []ThreadCandidate
	RequiresReview bool
	ResolvedAt    time.Time
}

// SemanticThreadResolver is a local/policy boundary for deterministic v0.5
// matching and future richer semantic similarity providers.
type SemanticThreadResolver interface {
	Resolve(personID string, triggers []SemanticTrigger, threads []Thread, at EvaluationMoment) (ThreadResolution, error)
}

type RetrievalDepth string

const (
	RetrievalAwareness     RetrievalDepth = "awareness"
	RetrievalCurrentState  RetrievalDepth = "current_state"
	RetrievalKeyContinuity RetrievalDepth = "key_continuity"
	RetrievalReconstructed RetrievalDepth = "reconstructed_thread"
	RetrievalDeepAudit     RetrievalDepth = "evidence_deep_audit"
)

type RetrievalRequest struct {
	PersonID        string
	ThreadID        string
	TemporalStart   *time.Time
	TemporalEnd     *time.Time
	Triggers        []SemanticTrigger
	Priority        PriorityFactors
	EvidenceSensitive bool
	AttentionBudget *AttentionBudget
	RequestedAt     time.Time
}

type RetrievalPlan struct {
	Depth              RetrievalDepth
	ThreadStateIDs     []string
	DeltaIDs           []string
	OpenLoopIDs        []string
	BlockIDs           []string
	EvidenceIDs        []string
	SelectedObjectCount int
	Reason             string
}

type RetrievalDepthPolicy interface {
	Plan(request RetrievalRequest, thread Thread, state CurrentThreadState, deltas []ThreadDelta, blocks []InteractionBlock) (RetrievalPlan, error)
}

type ContinuitySurfaceMode string

const (
	SurfaceSilent      ContinuitySurfaceMode = "silent_delta"
	SurfaceContextOnly ContinuitySurfaceMode = "context_only"
	SurfaceBrief       ContinuitySurfaceMode = "brief_mention"
	SurfaceMust        ContinuitySurfaceMode = "must_surface"
)

type ContinuitySurfaceDecision struct {
	ThreadID  string
	LinkID    string
	Mode      ContinuitySurfaceMode
	Reason    string
	EvaluatedAt time.Time
}

// --- Adaptive attunement ---------------------------------------------------

type InteractionSignalKind string

const (
	SignalPace              InteractionSignalKind = "pace"
	SignalPauseDensity      InteractionSignalKind = "pause_density"
	SignalResponseLatency   InteractionSignalKind = "response_latency"
	SignalSentenceLength    InteractionSignalKind = "sentence_length"
	SignalDirectness        InteractionSignalKind = "directness"
	SignalLexicalComplexity InteractionSignalKind = "lexical_complexity"
	SignalTopicSwitching    InteractionSignalKind = "topic_switching"
	SignalInteractionDensity InteractionSignalKind = "interaction_density"
)

// ObservedInteractionSignal accepts abstract synthetic measurements only. It
// deliberately omits raw audio or raw transcript persistence.
type ObservedInteractionSignal struct {
	ID             string
	PersonID       string
	BlockID        string
	Kind           InteractionSignalKind
	Value          float64
	Unit           string
	ContextIDs     []string
	ObservedAt     time.Time
	Provenance     Provenance
	Privacy        PrivacyClass
	Confidence     float64
}

func (s ObservedInteractionSignal) Validate() error {
	if normalizeText(s.ID) == "" || normalizeText(s.PersonID) == "" || s.Kind == "" || s.ObservedAt.IsZero() || s.Privacy != PrivacyRawCommunicationSignal || !isZeroOrUnit(s.Confidence) {
		return ErrInvalidInteractionSignal
	}
	return nil
}

type BaselineMetric struct {
	Kind             InteractionSignalKind
	Mean             float64
	Tolerance        float64
	ObservationCount int
}

// PersonalInteractionBaseline is person/context scoped. It is a derived
// aggregate, not a demographic or population stereotype.
type PersonalInteractionBaseline struct {
	ID                  string
	PersonID            string
	ContextSignature    string
	Metrics             []BaselineMetric
	ObservationCount    int
	Confidence          float64
	LastValidatedAt     time.Time
	DecayAfter          time.Duration
	Freshness           FreshnessState
	Privacy             PrivacyClass
}

func (b PersonalInteractionBaseline) Validate() error {
	if normalizeText(b.ID) == "" || normalizeText(b.PersonID) == "" || b.ObservationCount < 0 || !isZeroOrUnit(b.Confidence) || b.LastValidatedAt.IsZero() || b.DecayAfter < 0 || b.Freshness.Validate() != nil || b.Privacy != PrivacyDerivedBaseline {
		return ErrInvalidInteractionBaseline
	}
	for _, metric := range b.Metrics {
		if metric.Kind == "" || metric.Tolerance < 0 || metric.ObservationCount < 0 {
			return ErrInvalidInteractionBaseline
		}
	}
	return nil
}

type InteractionHypothesis string

const (
	HypothesisPossibleUrgency    InteractionHypothesis = "possible_urgency"
	HypothesisPossibleOverload   InteractionHypothesis = "possible_overload"
	HypothesisPossibleEngagement InteractionHypothesis = "possible_high_engagement"
	HypothesisPossibleLowEnergy  InteractionHypothesis = "possible_low_energy"
	HypothesisPossibleUncertainty InteractionHypothesis = "possible_uncertainty"
	HypothesisPossibleFrustration InteractionHypothesis = "possible_frustration"
	HypothesisPossibleExcitement InteractionHypothesis = "possible_excitement"
)

// InferredInteractionState is a non-diagnostic, uncertain interaction
// hypothesis. Alternatives preserve ambiguity rather than implying certainty.
type InferredInteractionState struct {
	ID                     string
	PersonID               string
	BlockID                string
	Hypothesis             InteractionHypothesis
	EvidenceSignalIDs      []string
	AlternativeExplanations []string
	Confidence             float64
	Uncertainty            float64
	ContextIDs             []string
	EventAt                time.Time
	EvaluatedAt            time.Time
	Freshness              FreshnessState
	Privacy                PrivacyClass
}

func (s InferredInteractionState) Validate() error {
	if normalizeText(s.ID) == "" || normalizeText(s.PersonID) == "" || s.Hypothesis == "" || s.EventAt.IsZero() || s.EvaluatedAt.IsZero() || !isZeroOrUnit(s.Confidence) || !isZeroOrUnit(s.Uncertainty) || s.Freshness.Validate() != nil || s.Privacy != PrivacyInferredState {
		return ErrInvalidInteractionInference
	}
	return nil
}

type AttunementControlMode string

const (
	AttunementNormal            AttunementControlMode = "normal"
	AttunementReduced           AttunementControlMode = "reduced"
	AttunementDisabled          AttunementControlMode = "disabled"
	AttunementTemporaryOverride AttunementControlMode = "temporary_override"
)

type AdaptationObjective string

const (
	ObjectiveComprehension          AdaptationObjective = "improve_comprehension"
	ObjectiveReduceCognitiveFriction AdaptationObjective = "reduce_cognitive_friction"
	ObjectiveQuietSupport           AdaptationObjective = "quiet_support"
	ObjectiveCovertPersuasion       AdaptationObjective = "covert_persuasion"
	ObjectiveEngagementMaximization AdaptationObjective = "engagement_maximization"
	ObjectiveDependency             AdaptationObjective = "dependency_creation"
)

// InteractionAdaptationDecision changes interaction presentation, not facts or
// authority. Values are normalized strengths; MaxChoices is an explicit cap.
type InteractionAdaptationDecision struct {
	ID                   string
	PersonID             string
	BlockID              string
	InferredStateIDs     []string
	Control              AttunementControlMode
	Objective            AdaptationObjective
	ResponsePace         float64
	ResponseVerbosity    float64
	SentenceLength       float64
	Directness           float64
	Warmth               float64
	MaxChoices           int
	ResurfaceOpenLoops   bool
	AvoidUnrelatedSurface bool
	FactualContentChanged bool
	Reason               string
	Confidence           float64
	EvaluatedAt          time.Time
	Reversible           bool
	Privacy              PrivacyClass
}

func (d InteractionAdaptationDecision) Validate() error {
	if normalizeText(d.ID) == "" || normalizeText(d.PersonID) == "" || d.Control == "" || d.Objective == "" || d.MaxChoices < 0 || d.EvaluatedAt.IsZero() || !isZeroOrUnit(d.ResponsePace) || !isZeroOrUnit(d.ResponseVerbosity) || !isZeroOrUnit(d.SentenceLength) || !isZeroOrUnit(d.Directness) || !isZeroOrUnit(d.Warmth) || !isZeroOrUnit(d.Confidence) || d.Privacy != PrivacyAdaptationDecision {
		return ErrUnsafeAdaptation
	}
	return nil
}

type AttunementEpisode struct {
	ID                   string
	PersonID             string
	BlockID              string
	SignalIDs            []string
	InferredStateIDs     []string
	AdaptationDecisionID string
	ThreadIDs            []string
	CreatedAt            time.Time
}

type InteractionIntervention struct {
	ID          string
	PersonID    string
	EpisodeID   string
	DecisionID  string
	Summary     string
	OccurredAt  time.Time
	Reversible  bool
	Provenance  []Provenance
}

type InteractionOutcomeStatus string

const (
	InteractionOutcomeBeneficial InteractionOutcomeStatus = "beneficial"
	InteractionOutcomeMixed      InteractionOutcomeStatus = "mixed"
	InteractionOutcomeUnhelpful  InteractionOutcomeStatus = "unhelpful"
	InteractionOutcomeUnknown    InteractionOutcomeStatus = "unknown"
)

type InteractionOutcome struct {
	ID                  string
	PersonID            string
	EpisodeID           string
	InterventionID      string
	Status              InteractionOutcomeStatus
	Summary             string
	EvidenceIDs         []string
	ExplicitFeedback    string
	ContextSignature    string
	TimeToOutcome       time.Duration
	AlternativeExplanations []string
	OccurredAt          time.Time
	RecordedAt          time.Time
	Privacy             PrivacyClass
}

// PersonalAttunementPattern is a decaying correlation record. It cannot claim
// that an intervention caused an outcome.
type PersonalAttunementPattern struct {
	ID                  string
	PersonID            string
	ContextSignature    string
	Hypothesis          InteractionHypothesis
	StrategyFingerprint string
	ObservationCount    int
	BeneficialCount     int
	UnhelpfulCount      int
	MixedCount          int
	Confidence          float64
	Freshness           FreshnessState
	LastOutcomeAt       time.Time
	UserOverridable     bool
	CorrelationOnly     bool
	Privacy             PrivacyClass
}

func (p PersonalAttunementPattern) Validate() error {
	if normalizeText(p.ID) == "" || normalizeText(p.PersonID) == "" || p.Hypothesis == "" || normalizeText(p.StrategyFingerprint) == "" || p.ObservationCount < 0 || p.BeneficialCount < 0 || p.UnhelpfulCount < 0 || p.MixedCount < 0 || p.BeneficialCount+p.UnhelpfulCount+p.MixedCount > p.ObservationCount || !isZeroOrUnit(p.Confidence) || p.LastOutcomeAt.IsZero() || p.Freshness.Validate() != nil || !p.UserOverridable || !p.CorrelationOnly || p.Privacy != PrivacyLearnedPattern {
		return ErrInvalidAttunementPattern
	}
	return nil
}

// AttunementSafetyPolicy rejects unsafe objectives and produces reversible,
// user-controlled interaction choices. It never changes factual content.
type AttunementSafetyPolicy interface {
	Decide(control AttunementControlMode, states []InferredInteractionState, patterns []PersonalAttunementPattern, at EvaluationMoment) (InteractionAdaptationDecision, error)
}
