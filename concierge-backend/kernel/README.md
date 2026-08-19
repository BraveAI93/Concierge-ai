# Central Intelligence Kernel v0.2

> **Status:** isolated, dependency-free, person-centric domain kernel. The package is not wired into HTTP handlers, UI, Supabase, the legacy database, production, or external action providers.

The Central Intelligence Kernel models a person’s world, evidence, unfinished work, decisions, permissions, outcomes, and learning. v0.2 corrects the semantic model introduced in v0.1. It keeps useful ingestion and validity timestamps, but it no longer treats them as substitutes for interaction gap or attention-time capacity. It also moves arbitrary scoring and decision policy out of canonical primitives.

## Architectural boundaries

| Layer | Responsibility | v0.2 implementation |
|---|---|---|
| Canonical domain | Person-scoped facts, time, attention needs, evidence lineage, permissions, outcomes, and invariants | Pure Go records and deterministic primitives |
| Policy layer | Weights, cut-offs, penalties, attention ranking, mismatch thresholds, and decision behavior | Injected `EvaluationPolicy`; `WeightedPolicyConfig` is one replaceable implementation |
| Adapter boundary | Maps legacy/external data to canonical records and persists/executed approved actions | Interfaces only; no concrete adapter is enabled |
| Existing application shell | Current HTTP, UI, legacy DB/Supabase, and deployment behavior | Intentionally unchanged |
| Proposed persistence | Future normalized storage and RLS review target | Commented-only design artifact; never executed |

The package never imports the existing `db` package. No package function opens a network connection, issues SQL, or performs an external side effect.

## Canonical model

| Concern | Canonical objects |
|---|---|
| Personal world | `Person`, `PersonalWorld`, `Entity`, `Alias`, `Context` |
| Records and truth | `Event`, `Memory`, `Claim`, `Evidence`, `Provenance`, `MemoryEventLink` |
| Intent and attention | `Goal`, `Constraint`, `PendingIntent`, `OpenLoop`, `InteractionGapState`, `EffortAttention`, `AttentionBudget`, `AttentionCandidate`, `AttentionAllocation` |
| Decision | `PriorityFactors`, `ActionWindow`, `DeadlineFeasibility`, `PriorityMismatch`, `Opportunity`, `Decision` |
| Action and learning | `Permission`, `Scope`, `ActionProposal`, `ActionGate`, `Outcome`, `SelfAudit` |
| Freshness and lineage | `FreshnessState`, `ClaimLineage`, `ClaimSelection`, `SupersessionPolicy` |

All canonical records retain a `PersonID` boundary. Cross-person links, opportunities, constraints, permissions, attention candidates, and supersession transitions are rejected before policy evaluation.

## Independent temporal cognition

v0.2 makes the required temporal concepts distinct. `RecordedAt`, `EffectiveAt`, and `ExpiresAt` remain useful auxiliary clocks, but none replaces wall-clock evaluation, semantic event time, interaction gap, or effort/attention duration.

| Independent concept | Canonical representation | Meaning |
|---|---|---|
| Wall-clock/current evaluation time | `EvaluationMoment.WallClockAt` | The moment a deterministic primitive evaluates the world; supplied by caller, not stored as semantic fact time |
| Semantic event time | `TemporalState.EventAt` | When the occurrence is intended or happened |
| Elapsed interaction/gap time | `InteractionGapState.LastInteractionAt` and `Elapsed(at)` | Time since the last relevant interaction for a person/entity/context, independent of event age |
| Effort/attention time | `EffortAttention.EstimatedEffort` and `EstimatedAttention` | Duration to complete work and duration to allocate from finite attention capacity |
| Ingestion time | `TemporalState.RecordedAt` | When the kernel learned or recorded the fact; it may precede a future `EventAt` |
| Validity interval | `TemporalState.EffectiveAt` and optional `ExpiresAt` | When a record is applicable; expiry is not a task deadline |
| Reconsideration time | `TemporalState.AttentionAt` | When an item should be reconsidered, not a capacity reservation |

`TemporalState.Validate` intentionally permits `RecordedAt < EventAt`. A scheduled appointment learned today and occurring tomorrow is valid.

## Priority and decision semantics

`PriorityFactors` distinguishes subjective importance, objective stakes, expected impact, reversibility, uncertainty, opportunity cost, and effort/attention cost. `ActionWindow.LatestSafeActionAt` and `EstimatedEffort` separately model deadline feasibility. Generic v0.1 fields remain only as clearly marked compatibility inputs; v0.2 callers must not use `ExpectedValue`, `Risk`, or `Importance` as silent substitutes for the independent factors.

| Factor | Representation | Notes |
|---|---|---|
| Subjective importance | `PriorityFactors.SubjectiveImportance`, `Goal.SubjectiveImportance` | Preserves the person’s stated preference and goal authority |
| Objective stakes | `PriorityFactors.ObjectiveStakes` | Separate from subjective importance and expected impact |
| Deadline/time-to-action | `ActionWindow.LatestSafeActionAt`, `DeadlineFeasibility` | Compared with estimated effort; never inferred from generic expiry |
| Expected impact | `PriorityFactors.ExpectedImpact` | Separate from stakes and preference |
| Reversibility | `PriorityFactors.Reversibility` | Ease of reversing the action; not risk |
| Uncertainty | `PriorityFactors.Uncertainty` plus claim confidence/evidence | Separate from impact and stakes |
| Opportunity cost | `PriorityFactors.OpportunityCost` | A separate cost input, not a soft constraint alias |
| Effort/attention cost | `PriorityFactors.EffortAttentionCost` and `EffortAttention` durations | Scalar policy cost plus schedulable duration/capacity |

A low subjective importance combined with high objective stakes produces an explicit `PriorityMismatch` with `MustSurface=true` and `PreservesGoalAuthority=true`. The policy returns `DecisionSurface` unless a hard constraint blocks it or the action window is infeasible. This surfaces the conflict without silently overriding the person’s goals or silently burying the item.

## Finite Attention Budget and contextual resurfacing

`AttentionBudget` is person-scoped and finite. It has a time window, an attention-duration capacity, a maximum number of competing items, interruption cost, and current context/entity IDs. `AllocateAttention` receives unresolved `AttentionCandidate` records and produces both selected and deferred items with reasons.

The allocator validates person boundaries, context/entity matching, unresolved state, required attention duration, total capacity, and maximum competing items. It deterministically ranks by injected policy score and then open-loop ID. A candidate may be deferred because it is resolved, exceeds remaining attention capacity, or exceeds the maximum competing-item limit. This is an actual capacity constraint rather than an independent score list.

## World drift, freshness, and provenance-preserving supersession

Historical records are not overwritten. `FreshnessState` tracks `StaleAfter`, `LastValidatedAt`, `LastRevalidatedAt`, and status. `EvaluateClaimFreshness` can return fresh, stale, historical, or superseded; `RevalidateClaim` is explicit. `ClaimLineage` records a superseding claim, the supporting evidence, timestamp, and historical preservation marker.

`SupersedeClaim` rejects a claim transition unless injected `SupersessionPolicy` requirements are met: the replacement evidence must have required authority and relevance, carry provenance, and (when configured) explicitly contradict the predecessor. Newness is an input to the caller’s policy flow, but it is not enough by itself. `SelectCurrentClaim` refuses to select among multiple fresh claims merely because one is newer; it returns an ambiguity result until an explicit supersession policy has resolved the lineage.

## Policy injection

Canonical records and invariants contain no recommendation thresholds, score weights, or soft-constraint penalties. `EvaluationPolicy` defines the policy boundary. `WeightedPolicyConfig` contains configurable temporal and opportunity weights, mismatch thresholds, attention ranking costs, and recommend/defer thresholds. `DefaultV02Policy` is a selectable configuration, not the domain contract. `LegacyV01Policy` is an explicit compatibility adapter used only by retained v0.1 wrappers.

| Primitive | Canonical responsibility | Policy responsibility |
|---|---|---|
| `EvaluateOpportunityWithPolicy` | Validate records, person boundaries, activity, goal references, hard constraints, and deadline feasibility | Compute utility, soft-constraint effect, mismatch state, and ranking |
| `DecideOpportunityWithPolicy` | Preserve evaluated result and hard-block/deadline state | Choose recommend, defer, decline, or must-surface outcome |
| `AllocateAttention` | Enforce person, unresolved-state, capacity, and count invariants | Rank candidates and define contextual bonus/cost |
| `EvaluateTemporal` | Preserve independent time inputs | Weight temporal urgency and attention factors |

## v0.2 deterministic primitives

| Primitive | Contract |
|---|---|
| `LinkMemoryToEvent` | Creates an idempotent reciprocal relationship and rejects cross-person links. |
| `TransitionPendingIntent` | Enforces the declared intent lifecycle and appends immutable transitions. |
| `EvaluateDeadlineFeasibility` | Classifies no deadline, feasible deadline, or infeasible action window from latest safe action time and effort duration. |
| `EvaluateOpportunityWithPolicy` | Produces canonical evaluation state using an injected policy; hard constraints remain domain invariants. |
| `AllocateAttention` | Finite, deterministic selection/deferral of unresolved loops with contextual resurfacing. |
| `EvaluateClaimFreshness` / `RevalidateClaim` | Detects staleness and performs explicit revalidation without deleting history. |
| `SupersedeClaim` / `SelectCurrentClaim` | Preserves lineage and provenance; authority, relevance, contradiction, and policy—not newness alone—govern supersession. |
| `PrepareActionGate` / `ApproveAction` | Keeps explicit permission and human approval checks before execution. |
| `AuditOutcome` | Requires outcome evidence and provenance before learning classification. |

## Verification

Run from `concierge-backend`:

```bash
/usr/local/go/bin/go test ./...
/usr/local/go/bin/go test -cover ./kernel
/usr/local/go/bin/go vet ./kernel
```

The test suite includes all retained v0.1 scenarios plus v0.2 regressions for future event recording, independent event age and interaction gap, finite attention capacity and contextual resurfacing, priority mismatch surfacing, deadline feasibility, reversibility/opportunity-cost effects, staleness and revalidation, authoritative provenance-preserving supersession, and injected policies producing different decisions over unchanged canonical data.

## Non-goals

v0.2 does not implement persistence, background execution, scheduling integration, LLM reasoning, cross-person inference, HTTP routes, UI, Supabase migration, legacy database mutation, production deployment, or external action execution. The commented SQL proposal is not a migration and must not be executed.
