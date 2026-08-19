# Central Intelligence Runtime v0.3

> **Status:** local vertical slice only. This package is server-side composition around the independent `kernel` package. It has no registered HTTP route, no UI, no active Supabase adapter, no SQL execution, and defaults to disabled.

## Composition boundary

```text
Authenticated server principal
        │  (stable server subject; never a browser-selected target person)
        ▼
IdentityResolver ──► canonical Person + source-profile binding
        │
        ▼
ConversationMessageAdapter ──► Event / Evidence / Claim / Memory / Intent / OpenLoop
        │                                  │ source IDs + timestamps + provenance preserved
        ▼                                  ▼
RuntimeRepository transaction ──► temporal + interaction gap + finite attention allocation
        │
        ▼
Policy-injected Opportunity / Decision ──► ActionProposal / ActionGate(awaiting approval)
        │
        ▼
Person-scoped RuntimeState retrieval and deterministic idempotent replay
```

The canonical `kernel` package imports none of this package. This runtime package may reference the legacy `db.Conversation` and `db.Message` types only through a one-way adapter; it does not change or invoke legacy DB/Supabase functions.

## Stable identity mapping

A runtime request accepts an `AuthenticatedPrincipal` containing a **server-resolved stable subject ID**. `IdentityResolver` maps that subject to `PersonBinding`, whose immutable fields are a canonical `Person` and the internal source profile ID permitted to ingest conversations for that person.

The runtime request has **no target `PersonID`, world ID, public slug, or profile selector**. A source conversation whose internal `ProfileID` does not match the resolved binding is rejected. Production composition must resolve a session token server-side to a stable profile/account identity and then bind that internal identity to a canonical person. It must not treat a browser-provided slug as a canonical identity or authorization target.

## First source adapter

`ConversationMessageAdapter` accepts a legacy Concierge `db.Conversation` and `db.Message`, preserving their original IDs, profile ID, conversation/session IDs, message role, content, and timestamps in the runtime source record and provenance. It emits canonical records only when a conservative scheduling-request predicate is supported by the message text. A latest-safe-action timestamp can only be supplied alongside a verbatim deadline-evidence substring present in the source message; otherwise no deadline is inferred.

## Repository contract

`RuntimeRepository` provides person-scoped transaction, idempotency, state retrieval, attention-budget reads, and append-preserving runtime persistence. `InMemoryRuntimeRepository` is the local implementation. It implements the v0.2 `kernel.KernelRepository` save contract as well as the runtime transaction contract, with same-person checks and deterministic idempotency keys.

The in-memory repository is intentionally **not** a persistence activation. Its interfaces define the shape a future Postgres/Supabase adapter must implement. The review-only migration package in `../db/proposed_migrations/ci_kernel_v03/` is not executed.

## Feature boundary

`Feature` is an explicit server-side composition switch. `DisabledFeature()` is the constructor default. The runtime does nothing until a server composition supplies `EnabledFeature()`; no route currently does so. This is not a UI feature flag and no public caller can enable it.

## Vertical slice semantics

The runtime path is:

1. resolve server-authenticated principal to `PersonBinding`;
2. reject source-profile mismatch before any write;
3. map a source message to Event/Evidence/Claim/Memory plus PendingIntent/OpenLoop if the message is an unresolved scheduling request;
4. atomically persist source-derived records and an idempotency record;
5. evaluate interaction gap and finite attention budget;
6. create opportunity and use injected v0.2 policy for decision;
7. create ActionProposal and `ActionGate` in `awaiting_approval` when scheduling-draft permission requires approval;
8. retrieve the complete person-scoped state or replay the same source idempotently.

No action executor runs in this slice. An action gate can only be approved by the existing canonical primitive, and external action execution remains absent.
