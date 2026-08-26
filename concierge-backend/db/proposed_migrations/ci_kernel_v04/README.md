# Central Intelligence Staging Persistence Package — v0.5

> **Validation target:** `ci_kernel_v04_test` on a disposable local PostgreSQL 16.15 cluster bound to `/tmp/ci-v05-pgsocket` on port `55433`. This target is sandbox-local, has no Supabase endpoint, and is stopped after validation. It is **not** a production migration package.

## Package Contents

| File | Purpose |
|---|---|
| `001_staging_up.sql` | Creates the isolated `ci_kernel_v04` schema, identity roles, tenant-safe persistence keys, v0.5 continuity edges, forced RLS, and least-privilege grants. |
| `001_staging_down.sql` | Destructive disposable-only schema rollback after an export and dependency review. |
| `README.md` | Documents target constraints, identity authority, typed graph integrity, privacy treatment, RLS, grants, and non-activation boundaries. |

## Tenant-Safe Storage Model

`people`, `worlds`, `person_binding_subjects`, and `person_profile_links` maintain the reviewed v0.4.1 durable identity boundary. `runtime_sources`, `records`, `memory_event_links`, `attention_allocations`, and `runtime_replays` use person-scoped operational keys. v0.5 adds relational `continuity_links` and `thread_deltas` tables because their typed graph references must be tenant-safe at the SQL boundary rather than only inside JSONB.

| Table | Key and integrity property |
|---|---|
| `runtime_sources` | `(person_id, id)` primary key and `(person_id, message_id)` source identity. |
| `records` | `(record_kind, person_id, id)` primary key for canonical block, thread, state, signal, inference, decision, episode, intervention, outcome, and pattern payloads. |
| `continuity_links` | `(person_id, id)` key plus composite source and target FKs into the same person’s canonical records. |
| `thread_deltas` | `(person_id, id)` key plus same-person FKs to target `thread`, originating record, and its own `thread_delta` record payload. |
| `memory_event_links`, `attention_allocations`, `runtime_replays` | Retain v0.4.1 person-scoped keys and same-person integrity. |

A raw source ID, canonical ID, delta ID, replay key, or allocation ID in one Personal World cannot reserve, reveal, or collide with the same logical value in another Person’s world.

## v0.5 Continuity and Attunement Persistence

`records` stores complete canonical JSONB payloads for `InteractionBlock`, `Thread`, `CurrentThreadState`, `ObservedInteractionSignal`, `PersonalInteractionBaseline`, `InferredInteractionState`, `InteractionAdaptationDecision`, `AttunementEpisode`, `InteractionIntervention`, `InteractionOutcome`, and `PersonalAttunementPattern`, in addition to all retained v0.2–v0.4.1 domain objects.

`continuity_links` records typed, evidence-bearing block/thread or other canonical-object relationships. `thread_deltas` records material change to thread state with explicit source, target thread, event/evaluation times, and tenant-safe foreign keys. Original source records and canonical evidence remain authoritative; current state is reconstructable from append-preserving delta lineage.

Raw communication signal is intentionally not a default persisted payload in this package. The runtime accepts only abstract synthetic/local measurement values in the v0.5 vertical slice. Derived baselines, uncertain non-diagnostic interaction hypotheses, reversible adaptation decisions, observed outcomes, and decaying correlation-only patterns use their respective canonical privacy classifications. No medical or psychiatric diagnosis, population stereotype, causal claim, raw audio, or external account collection is introduced.

## Identity Authority, RLS, and Grants

`ci_kernel_runtime` can resolve existing stable subject and internal profile bindings, read person-scoped runtime data, and append ordinary runtime records. It cannot provision, rebind, attach, update, or delete high-trust identity records. `ci_kernel_identity_provisioner` retains the v0.4.1 INSERT-only initial identity boundary and has no UPDATE or DELETE authority.

RLS is both enabled and **forced** on all eleven staging tables, including `continuity_links` and `thread_deltas`. Runtime policies permit only person-context SELECT/INSERT operations and the existing `worlds` compatibility UPDATE. There is no runtime DELETE grant. `PUBLIC` is revoked from the schema, all tables, and all current functions. Default function privileges in the schema also revoke `PUBLIC EXECUTE`, preventing a future helper function from silently widening the surface. The runtime receives execution only for the two context functions.

The local migration executor is a disposable test convenience. A non-disposable environment requires separate reviewed runtime/provisioner principals, secret distribution, audited role membership, and independent RLS/adversarial validation.

## Backup, Retention, and Production Exclusions

Before any non-disposable use, operators must export the target schema, verify database identity, and rehearse restoration separately. Conversation content, provenance, derived continuity, interaction measurements, outcomes, and pattern records may be personal data. Retention, deletion, access, legal hold, encryption/key management, consent, and backup-restoration requirements remain subject to the parallel Trust/Data Boundary approval.

Do **not** apply this package to production Supabase. v0.5 does not integrate with P0, activate HTTP routes, add a user interface, execute an external action, use live accounts, push, merge, or deploy. Remaining prerequisites include production identity backfill, P0 consent and recording-path integration, privacy/retention approvals, production-grade feature/kill-switch ownership, observability, load testing, backup/restore rehearsal, controlled activation, and a separately reviewed protocol boundary for rich semantic or audio analysis.
