# Central Intelligence v0.5 — Persistent Continuity and Adaptive Attunement Design

## Purpose and Boundary

v0.5 extends the **canonical kernel** with person-scoped, append-preserving structures for conversational continuity and adaptive interaction. A physical chat is a source container, not a memory boundary. The canonical system therefore organizes semantic/time-bounded **InteractionBlocks**, long-lived **Threads**, typed **ContinuityLinks**, append-only **ThreadDeltas**, and reconstructable **CurrentThreadState** without replacing Events, Memories, Claims, Evidence, OpenLoops, Decisions, or Outcomes.

Adaptive attunement is likewise an evidence-backed interaction layer rather than an emotional diagnosis system. It preserves a distinction among observation, uncertain inference, reversible adaptation, observed outcome, and a weak person/context-specific association pattern. It never asserts medical or psychiatric conditions, causes, or emotional certainty.

## Canonical Choice: `InteractionBlock`

`InteractionBlock` is the v0.5 canonical name for a bounded interaction period. It contains source references only as provenance; it is not a chat/session identity. The block retains semantic start/end times where known, evaluation/ingestion time, labels, entity and context anchors, affected work, evidence, importance, confidence, freshness, supersession, and its thread membership.

The `InteractionBoundaryPolicy` is injected. Its deterministic test implementation evaluates explicit end markers, source/context change, task completion, semantic discontinuity, and a supplied interaction gap. It records the evidence and confidence leading to a split or continuation. No universal inactivity duration is embedded into the canonical contract.

## Thread and Continuity Graph

A `Thread` is a source-independent ongoing subject anchored by canonical entity/project/context IDs and aliases. It organizes continuity only; underlying Events, Memories, Evidence, Claims, and source blocks remain authoritative. `SemanticTrigger` values can come from aliases, entities, projects, contexts, goals, OpenLoops, decisions, temporal cues, or local deterministic paraphrase cues. The resolver returns ranked candidates. It only attaches a block to a Thread if injected policy confidence crosses its selection threshold; otherwise it returns an unresolved candidate relation rather than silently merging unrelated contexts.

`ContinuityLink` has a typed source, target, relation, evidence/provenance, confidence, freshness, evaluation time, and optional supersession. Links are distinct from visible attention. Support-graph relations must be compiled through the existing finite AttentionBudget and policy before they can interrupt the user.

`ThreadDelta` is the append-preserving mechanism for material cross-thread or later updates. `BuildCurrentThreadState` applies non-superseded deltas in deterministic temporal order to a baseline and records included versus historical delta lineage. This allows later changes to affect current continuity while preserving superseded historical provenance.

## Retrieval and Attention

`RetrievalDepthPolicy` chooses an explicit awareness/current/key-continuity/reconstructed/deep-audit depth from existing priority dimensions, deadline urgency, reversibility, uncertainty, opportunity cost, effort cost, evidence sensitivity, and available attention capacity. Retrieval returns selected object IDs and operation counts, not a lifetime transcript. The ordering is Thread state/index, latest material deltas, unresolved work, relevant blocks, then original evidence only at the justified depth.

`ContinuitySurfaceDecision` is distinct from retrieval. Low-value continuity can be used contextually or stored as a silent delta. A high-stakes/deadline-critical link can become `MustSurface` through policy and then participate in attention allocation; relation discovery alone does not create an OpenLoop, notification, or task.

## Adaptive Attunement and Safety

Raw future voice/acoustic payloads are deliberately not a v0.5 persistence default. `ObservedInteractionSignal` stores only abstract, local synthetic measurements and privacy classification. `PersonalInteractionBaseline` is person/context scoped, confidence-bearing, minimally observation-counted, and freshness/decay aware. `InferredInteractionState` is a hypothesis with alternative explanations and uncertainty; it cannot represent diagnostic certainty.

`InteractionAdaptationDecision` changes only interaction parameters such as verbosity, pace, choice count, directness, and resurfacing behavior. It includes a user-control mode (`normal`, `reduced`, `disabled`, `temporary_override`) and an explicit objective. `AttunementSafetyPolicy` rejects persuasion, engagement maximization, dependency, impersonation/parody, concealed uncertainty, and agitation amplification. The policy is fail-closed if user control disables attunement.

The closed loop is modeled as `AttunementEpisode` → `InteractionIntervention` → `InteractionOutcome` → `PersonalAttunementPattern`. Patterns record supporting, failing, and mixed outcomes, context signature, time-to-outcome, explicit feedback, alternative explanations, confidence, and decay. They are explicitly correlations, not causal claims. One outcome can only produce a weak pattern; repeated context-similar evidence may raise confidence, while contradictory outcomes lower it.

## Persistence and Integration

The runtime remains default-off and has no HTTP handler, UI, external action executor, production Supabase adapter, P0 dependency, or live account source. A local runtime service may accept synthetic interaction inputs only after existing identity, source-profile ownership, activation, and consent boundaries pass. New persistence uses canonical record payloads plus relational `continuity_links` and `thread_deltas` edges where person-scoped graph integrity materially matters. All operational keys are person scoped.

The v0.4.1 `ci_kernel_runtime` role retains no identity-provisioning capability. New tables receive `ENABLE` and `FORCE ROW LEVEL SECURITY`, person-bound policies, table-specific grants, no DELETE grant, and no PUBLIC function execution. JSONB retains canonical payloads, but only noncritical derived references remain JSONB-only; this limitation is documented rather than hidden.

## Privacy Classification

| Class | v0.5 treatment |
|---|---|
| `raw_communication_signal` | Abstract input only; no raw audio/text capture added by v0.5. |
| `derived_baseline` | Minimal aggregate/range, person/context scoped, decaying and revalidatable. |
| `inferred_interaction_state` | Uncertain non-diagnostic hypothesis with evidence and alternatives. |
| `adaptation_decision` | Reversible interaction parameters, objective, and user control state. |
| `outcome_evidence` | Person-scoped observed result and provenance; no causal assertion. |
| `learned_pattern` | Decaying, weak correlation only, context-bound, user-overridable. |

## Rationale for Local Deterministic Policies

The implementation intentionally uses typed resolver, boundary, retrieval-depth, attunement-safety, and learning-policy contracts with deterministic local implementations. This verifies the architecture without embedding provider-specific embeddings, audio analysis, population stereotypes, or hidden product thresholds. Richer semantic/audio implementations can be added only through these ports after privacy, consent, and quality review.
