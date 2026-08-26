# Central Intelligence Runtime v0.5 — Local Continuity and Attunement Slice

## Activated Surface

There is **no route, UI, P0 integration, production database target, external account reader, audio collector, action executor, push, merge, or deployment** in this slice. `RuntimeService.IngestContinuity`, `RecordAttunementOutcome`, and `ResolveContinuity` are server-side composition methods only. `Feature.Enabled`, the server activation port, authenticated stable-subject resolution, internal source-profile ownership, and verified conversation-derived-memory consent are checked before continuity writes.

The only exercised source is a synthetic/local `ConversationMessage` fixture. The runtime preserves its identifiers as a `SourceRecord` and as `InteractionSourceRef` provenance but never makes a physical conversation or session the canonical memory boundary.

## Vertical Slice

```text
synthetic ConversationMessage
  → server stable subject → PersonBinding → consent/profile validation
  → SourceRecord + InteractionBlock
  → semantic trigger resolver → selected Thread or new unresolved-safe Thread
  → typed ContinuityLink + append-only ThreadDelta
  → reconstructable CurrentThreadState
  → abstract ObservedInteractionSignal + personal baseline
  → uncertain InferredInteractionState
  → reversible, user-controlled InteractionAdaptationDecision
  → AttunementEpisode → InteractionIntervention
  → observed InteractionOutcome → decaying PersonalAttunementPattern
```

The block-to-thread relationship is persisted in `continuity_links`, and material cross-thread updates are persisted in `thread_deltas`; both have person-scoped composite foreign keys. All other canonical v0.5 records persist as typed `records` JSONB payloads with person-scoped primary keys. This preserves append history while avoiding a premature wide relational schema for every evolving canonical field.

## Policy Injection and Attention

The runtime injects `InteractionBoundaryPolicy`, `SemanticThreadResolver`, `RetrievalDepthPolicy`, and `AttunementSafetyPolicy`. The local deterministic policies are test substitutes, not embedded product truth. `ResolveContinuity` returns a bounded retrieval plan for the authenticated person only. Low-value continuity remains current state or contextual knowledge; it does not create a notification, OpenLoop, or task.

A separate helper compiles a typed relation into a `ContinuitySurfaceDecision` only through caller-provided policy. It creates no direct attention or action side effect. Existing finite `AttentionBudget` and ActionGate policy remains the only route to attention competition or action approval.

## Attunement Safety and Privacy

The slice stores only abstract synthetic measurement values. It intentionally does not persist raw audio, raw transcript analysis, medical labels, psychiatric diagnoses, demographic assumptions, population comparisons, or claims of causality. Inference records include uncertainty and alternative explanations. Adaptation decisions are reversible, cannot change factual content, and reject persuasion, engagement maximization, dependency creation, impersonation, and hidden certainty.

Outcomes update a person/context-specific `PersonalAttunementPattern` as a correlation-only, user-overridable, decaying record. One observed outcome remains weak; repeated contradictory outcomes lower confidence and reduce adaptation. A disabled user-control mode writes no adaptation decision, episode, or intervention.

## Staging Persistence

The v0.5 extension was applied and tested only against a disposable local PostgreSQL cluster. `ci_kernel_runtime` has SELECT/INSERT access to the continuity tables and no DELETE grant. RLS is enabled and forced. `PUBLIC` has no schema, table, current-function, or future-schema-function execution privilege. The separate v0.4.1 identity provisioner remains the only initial identity-writing boundary.

Before any controlled staging activation, remaining work includes P0 consent/recording composition, durable production subject migration, reviewed principals/secrets, retention/deletion governance, observability, backup/restore rehearsal, rich semantic/audio provider review, load tests, and explicit activation ownership.
