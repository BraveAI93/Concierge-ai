package kernel

import (
	"time"
)

// OpportunityPolicyInput is a canonical observation prepared by domain
// invariants. It contains no product-specific weights or cut-offs.
type OpportunityPolicyInput struct {
	Factors             PriorityFactors
	Deadline            DeadlineFeasibility
	SoftConstraintCount int
	TemporalActive      bool
}

// PolicyAssessment is the policy-layer result consumed by the canonical
// evaluator. MustSurface remains explicit and never changes user goals.
type PolicyAssessment struct {
	Utility       float64
	Mismatch      PriorityMismatch
	DecisionBasis string
}

// EvaluationPolicy is injected into v0.2 evaluation. It owns arbitrary product
// scoring and decision policy; canonical records never own these coefficients.
type EvaluationPolicy interface {
	Validate() error
	AssessOpportunity(input OpportunityPolicyInput) (PolicyAssessment, error)
	Decide(assessment PolicyAssessment, hardBlocked bool, deadline DeadlineFeasibility) (DecisionKind, string, error)
	ScoreAttention(candidate AttentionCandidate, contextMatched bool) (float64, string, error)
}

// WeightedPolicyConfig is one replaceable policy implementation. Every weight,
// threshold, and penalty is supplied by configuration rather than canonical
// types or primitives.
type WeightedPolicyConfig struct {
	TemporalImportanceWeight     float64
	TemporalDeadlineWeight       float64
	TemporalRecencyWeight        float64
	TemporalAttentionDueWeight   float64
	SubjectiveImportanceWeight   float64
	ObjectiveStakesWeight        float64
	ExpectedImpactWeight         float64
	ReversibilityWeight          float64
	UncertaintyPenaltyWeight     float64
	OpportunityCostPenaltyWeight float64
	EffortAttentionPenaltyWeight float64
	SoftConstraintPenalty        float64
	FeasibleUrgencyWeight        float64
	AttentionContextBonus        float64
	AttentionInterruptionPenalty float64
	MustSurfaceSubjectiveAtMost  float64
	MustSurfaceObjectiveAtLeast  float64
	RecommendAt                  float64
	DeferAt                      float64
}

func (p WeightedPolicyConfig) Validate() error {
	for _, value := range []float64{
		p.TemporalImportanceWeight,
		p.TemporalDeadlineWeight,
		p.TemporalRecencyWeight,
		p.TemporalAttentionDueWeight,
		p.SubjectiveImportanceWeight,
		p.ObjectiveStakesWeight,
		p.ExpectedImpactWeight,
		p.ReversibilityWeight,
		p.UncertaintyPenaltyWeight,
		p.OpportunityCostPenaltyWeight,
		p.EffortAttentionPenaltyWeight,
		p.SoftConstraintPenalty,
		p.FeasibleUrgencyWeight,
		p.AttentionContextBonus,
		p.AttentionInterruptionPenalty,
		p.MustSurfaceSubjectiveAtMost,
		p.MustSurfaceObjectiveAtLeast,
		p.RecommendAt,
		p.DeferAt,
	} {
		if !isZeroOrUnit(value) {
			return ErrInvalidPolicy
		}
	}
	if p.DeferAt > p.RecommendAt {
		return ErrInvalidPolicy
	}
	return nil
}

func (p WeightedPolicyConfig) AssessOpportunity(input OpportunityPolicyInput) (PolicyAssessment, error) {
	if err := p.Validate(); err != nil || input.Factors.Validate() != nil {
		return PolicyAssessment{}, ErrInvalidPolicy
	}
	if !input.TemporalActive {
		return PolicyAssessment{DecisionBasis: "opportunity is not currently effective"}, nil
	}
	factors := input.Factors
	assessment := PolicyAssessment{
		Utility: clampUnit(
			p.SubjectiveImportanceWeight*factors.SubjectiveImportance +
				p.ObjectiveStakesWeight*factors.ObjectiveStakes +
				p.ExpectedImpactWeight*factors.ExpectedImpact +
				p.ReversibilityWeight*factors.Reversibility -
				p.UncertaintyPenaltyWeight*factors.Uncertainty -
				p.OpportunityCostPenaltyWeight*factors.OpportunityCost -
				p.EffortAttentionPenaltyWeight*factors.EffortAttentionCost -
				p.SoftConstraintPenalty*float64(input.SoftConstraintCount) +
				p.FeasibleUrgencyWeight*deadlineUrgency(input.Deadline),
		),
		DecisionBasis: "policy-scored independent priority factors and constraint inputs",
	}
	if factors.SubjectiveImportance <= p.MustSurfaceSubjectiveAtMost && factors.ObjectiveStakes >= p.MustSurfaceObjectiveAtLeast {
		assessment.Mismatch = PriorityMismatch{
			MustSurface:            true,
			PreservesGoalAuthority: true,
			SubjectiveImportance:   factors.SubjectiveImportance,
			ObjectiveStakes:        factors.ObjectiveStakes,
			Reason:                 "objective stakes are high while subjective importance is low; surface the conflict without overriding the person's goals",
		}
	}
	return assessment, nil
}

func deadlineUrgency(feasibility DeadlineFeasibility) float64 {
	if feasibility.State != DeadlineFeasible || feasibility.Remaining <= 0 {
		return 0
	}
	if feasibility.EstimatedEffort <= 0 {
		return 1
	}
	// The ratio remains a policy input transformed into an urgency signal. The
	// caller-selected policy config controls its weight in the final score.
	ratio := float64(feasibility.EstimatedEffort) / float64(feasibility.Remaining)
	return clampUnit(ratio)
}

func (p WeightedPolicyConfig) Decide(assessment PolicyAssessment, hardBlocked bool, deadline DeadlineFeasibility) (DecisionKind, string, error) {
	if err := p.Validate(); err != nil {
		return "", "", err
	}
	if hardBlocked {
		return DecisionDecline, "blocked by hard constraint", nil
	}
	if deadline.State == DeadlineInfeasible {
		return DecisionDefer, "valuable but infeasible within the latest safe action window", nil
	}
	if assessment.Mismatch.MustSurface {
		return DecisionSurface, assessment.Mismatch.Reason, nil
	}
	if assessment.Utility >= p.RecommendAt {
		return DecisionRecommend, assessment.DecisionBasis, nil
	}
	if assessment.Utility >= p.DeferAt {
		return DecisionDefer, assessment.DecisionBasis, nil
	}
	return DecisionDecline, assessment.DecisionBasis, nil
}

func (p WeightedPolicyConfig) ScoreAttention(candidate AttentionCandidate, contextMatched bool) (float64, string, error) {
	if err := p.Validate(); err != nil || candidate.Validate() != nil {
		return 0, "", ErrInvalidPolicy
	}
	score := candidate.BasePriority -
		p.AttentionInterruptionPenalty*candidate.Effort.InterruptionCost -
		p.EffortAttentionPenaltyWeight*candidate.Effort.ContextSwitchCost
	reason := "selected by base priority and attention costs"
	if contextMatched {
		score += p.AttentionContextBonus
		reason = "selected with matching context or entity"
	}
	return clampUnit(score), reason, nil
}

// DefaultV02Policy is a selectable product policy, not a canonical contract.
// Applications may inject another configuration without changing domain types.
func DefaultV02Policy() WeightedPolicyConfig {
	return WeightedPolicyConfig{
		TemporalImportanceWeight:     0.25,
		TemporalDeadlineWeight:       0.30,
		TemporalRecencyWeight:        0.20,
		TemporalAttentionDueWeight:   0.25,
		SubjectiveImportanceWeight:   0.20,
		ObjectiveStakesWeight:        0.30,
		ExpectedImpactWeight:         0.20,
		ReversibilityWeight:          0.05,
		UncertaintyPenaltyWeight:     0.10,
		OpportunityCostPenaltyWeight: 0.05,
		EffortAttentionPenaltyWeight: 0.05,
		SoftConstraintPenalty:        0.05,
		FeasibleUrgencyWeight:        0.05,
		AttentionContextBonus:        0.15,
		AttentionInterruptionPenalty: 0.10,
		MustSurfaceSubjectiveAtMost:  0.30,
		MustSurfaceObjectiveAtLeast:  0.80,
		RecommendAt:                  0.65,
		DeferAt:                      0.35,
	}
}

// LegacyV01Policy is an explicit compatibility adapter used only by deprecated
// v0.1 wrapper functions. It keeps existing tests and consumers stable without
// treating generic v0.1 fields as canonical v0.2 semantics.
func LegacyV01Policy() WeightedPolicyConfig {
	return WeightedPolicyConfig{
		TemporalImportanceWeight:     0.35,
		TemporalDeadlineWeight:       0.30,
		TemporalRecencyWeight:        0.20,
		TemporalAttentionDueWeight:   0.15,
		SubjectiveImportanceWeight:   0.30,
		ObjectiveStakesWeight:        0.25,
		ExpectedImpactWeight:         0.22,
		ReversibilityWeight:          0.00,
		UncertaintyPenaltyWeight:     0.12,
		OpportunityCostPenaltyWeight: 0.00,
		EffortAttentionPenaltyWeight: 0.08,
		SoftConstraintPenalty:        0.10,
		FeasibleUrgencyWeight:        0.15,
		AttentionContextBonus:        0.00,
		AttentionInterruptionPenalty: 0.00,
		MustSurfaceSubjectiveAtMost:  0.00,
		MustSurfaceObjectiveAtLeast:  1.00,
		RecommendAt:                  0.65,
		DeferAt:                      0.35,
	}
}

// LegacyTemporalEvaluation is a compatibility policy result for the v0.1
// wrapper. Its coefficients live in this policy file, not in domain logic.
func (p WeightedPolicyConfig) LegacyTemporalEvaluation(importance, deadline, recency, attentionDue float64) (float64, error) {
	if err := p.Validate(); err != nil {
		return 0, err
	}
	for _, value := range []float64{importance, deadline, recency, attentionDue} {
		if !isZeroOrUnit(value) {
			return 0, ErrInvalidPolicy
		}
	}
	// This method preserves v0.1 behavior only when called with LegacyV01Policy.
	// New callers can supply a different injected policy configuration.
	return clampUnit(
		p.TemporalImportanceWeight*importance +
			p.TemporalDeadlineWeight*deadline +
			p.TemporalRecencyWeight*recency +
			p.TemporalAttentionDueWeight*attentionDue,
	), nil
}

// NewTemporalPriorityInput ensures wall-clock evaluation, semantic event time,
// interaction gap, and effort/attention duration remain independently visible.
type TemporalPriorityInput struct {
	Moment               EvaluationMoment
	Temporal             TemporalState
	InteractionGap       InteractionGapState
	Effort               EffortAttention
	SubjectiveImportance float64
}

func (p WeightedPolicyConfig) EvaluateTemporal(input TemporalPriorityInput) (TemporalEvaluation, error) {
	if err := p.Validate(); err != nil || input.Moment.Validate() != nil || input.Temporal.Validate() != nil || input.InteractionGap.Validate() != nil || input.Effort.Validate() != nil || !isZeroOrUnit(input.SubjectiveImportance) {
		return TemporalEvaluation{}, ErrInvalidPolicy
	}
	at := input.Moment.WallClockAt
	result := TemporalEvaluation{Importance: input.SubjectiveImportance, EvaluatedAt: at, Active: input.Temporal.IsActive(at)}
	if !result.Active {
		return result, nil
	}
	gap, err := input.InteractionGap.Elapsed(at)
	if err != nil {
		return TemporalEvaluation{}, err
	}
	result.InteractionGap = gap
	result.Recency = 1 / (1 + input.Temporal.EventAt.Sub(at).Abs().Hours()/(24*7))
	result.AttentionDue = attentionDueAt(input.Temporal.AttentionAt, at)
	if input.Temporal.ExpiresAt != nil {
		result.DeadlineUrgency = attentionDueAt(*input.Temporal.ExpiresAt, at)
	}
	// Policy configuration, rather than the canonical evaluation, chooses score
	// weights. The generic scoring path remains documented as compatibility.
	score, err := p.LegacyTemporalEvaluation(result.Importance, result.DeadlineUrgency, result.Recency, result.AttentionDue)
	if err != nil {
		return TemporalEvaluation{}, err
	}
	result.Utility = score
	return result, nil
}

func attentionDueAt(target, at time.Time) float64 {
	hours := target.Sub(at).Hours()
	if hours <= 0 {
		return 1
	}
	return 1 / (1 + hours/24)
}
