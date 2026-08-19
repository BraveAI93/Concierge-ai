# Central Intelligence Kernel v0.1

> **Status:** isolated, dependency-free domain kernel. This package is not wired into HTTP handlers, does not issue database calls, and does not perform external actions.

The Central Intelligence Kernel is a person-centric server-side domain package. It establishes canonical objects and deterministic policy primitives while preserving the existing Concierge application as a compatibility shell. The package is intentionally pure Go: all persistence and side effects are represented by narrow interfaces at the boundary.

## Architectural boundaries

| Layer | Responsibility | v0.1 implementation |
|---|---|---|
| `kernel` domain | Canonical facts, temporal state, linking, evaluation, consent and audit policy | Implemented as pure Go types and functions |
| Adapter boundary | Maps legacy/external data into the canonical model and persists canonical records | Interfaces only; no Supabase adapter is enabled |
| Existing backend shell | Current HTTP routes, existing profile flow, Supabase utilities, and UI compatibility | Intentionally unchanged |
| Proposed persistence | Future normalized storage shape and tenancy/RLS review point | SQL proposal only; not executed |

The package never imports the existing `db` package. This prevents implicit production persistence, makes all policy behavior testable in memory, and avoids changing the current service contract.

## Canonical model

| Concern | Canonical objects |
|---|---|
| Personal world | `Person`, `PersonalWorld`, `Entity`, `Alias`, `Context` |
| Records and truth | `Event`, `Memory`, `Claim`, `Evidence`, `Provenance`, `MemoryEventLink` |
| Intent and attention | `Goal`, `Constraint`, `PendingIntent`, `OpenLoop`, `TemporalState`, `TemporalEvaluation` |
| Deliberation and action | `Opportunity`, `Decision`, `ActionProposal`, `ActionGate`, `Permission`, `Scope` |
| Learning | `Outcome`, `SelfAudit` |

Every canonical record carries a `PersonID` so person boundaries can be verified before data are linked or evaluated. Evidence holds provenance directly rather than treating confidence as an unsupported scalar.

## Deterministic primitives

| Primitive | Contract |
|---|---|
| Memory/event linking | `LinkMemoryToEvent` rejects empty identifiers and cross-person links, then writes a reciprocal, idempotent relationship. |
| Pending-intent lifecycle | `TransitionPendingIntent` accepts only the declared lifecycle graph and appends an immutable transition record. |
| Four-dimensional time | `TemporalState` separates `EventAt`, `RecordedAt`, `EffectiveAt`, and `AttentionAt`; optional expiry bounds current applicability. |
| Temporal utility | `EvaluateTemporalUtility` combines importance, deadline urgency, recency, and scheduled attention. It returns zero for inactive/expired state. |
| Evidence confidence | `AssessClaimConfidence` uses signed evidence support, quality, relevance, and provenance presence. Contradictory evidence raises conflict and moves the score toward uncertainty. |
| Opportunity and decision | `EvaluateOpportunity` evaluates active goals, hard/soft constraints, value, effort, risk, confidence, and temporal priority. `DecideOpportunity` returns `recommend`, `defer`, or `decline` from fixed thresholds. |
| Action approval | Permission scope and the `ActionGate` are separate. Permission determines whether a proposal is allowed; the gate records the human approval lifecycle before execution. |
| Outcome/self-audit | `AuditOutcome` compares expected and observed utility, requires outcome evidence, and emits a repeatable learning status. |

## Time semantics

`TemporalState` has four distinct clocks. `EventAt` is when something happened; `RecordedAt` is when the kernel learned it; `EffectiveAt` is when the state became applicable; and `AttentionAt` is when it should be reconsidered. `ExpiresAt`, when present, closes the effective interval rather than representing a fifth clock.

## Lifecycle invariants

| Lifecycle | Allowed transitions |
|---|---|
| Pending intent | `captured → clarifying → ready → proposed → in_progress → completed`; `captured`, `clarifying`, `ready`, `proposed`, and `in_progress` may also move to `cancelled`; non-terminal states may move to `expired`. |
| Action gate | `draft → awaiting_approval → approved → executed`; `awaiting_approval → rejected` or `expired`; `approved → expired`. A proposal without human approval may move `draft → approved` only after permission validation. |

## Integration contract

A future integration should instantiate repository and adapter implementations in a new service composition root, map legacy inputs explicitly, read permissions at evaluation time, and execute side effects only after an approved `ActionProposal`. It must not bypass the `ActionGate` or mutate evidence provenance.

The SQL proposal in `../db/proposed_migrations/20260819_central_intelligence_kernel_v0_1.sql` is a design artifact only. It is intentionally not registered with deployment tooling and has not been applied to Supabase.

## Verification

`scenario_test.go` contains a synthetic end-to-end proof of the full path: event → personal world/context → memory/open loop → temporal attention evaluation → opportunity/decision → action proposal/approval → outcome/provenance. Unit tests cover cross-person linking, invalid intent transitions, temporal validity, conflict-aware confidence, scope enforcement, and outcome invariants.

Run the package tests from `concierge-backend`:

```bash
/usr/local/go/bin/go test ./kernel
```

Run the complete existing backend suite from the same directory:

```bash
/usr/local/go/bin/go test ./...
```

No environment variables, network services, database tables, HTTP routes, or UI behavior are required by these tests.

## Non-goals for v0.1

This release does not enable persistence, background processing, LLM reasoning, automation execution, cross-person inference, UI changes, or production migrations. These are integration decisions that must be reviewed separately after the domain contracts have been accepted.

## Sources of truth

The package code and tests are the executable source of truth. This document explains the public architecture and must be updated alongside any future behavioral change.
