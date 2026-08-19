# Central Intelligence Runtime v0.4 — Disposable Staging Migration Package

> **Target used for validation:** `ci_kernel_v04_test` on a local PostgreSQL 16.14 server bound to the sandbox-only Unix socket `/home/ubuntu/ci-v04-pgsocket` and local port `55432`. This is not Supabase and has no production network endpoint.

## Package contents

| File | Purpose |
|---|---|
| `001_staging_up.sql` | Creates the isolated `ci_kernel_v04` schema, runtime role, canonical persistence tables, indexes, same-person composite references, and RLS policies. |
| `001_staging_down.sql` | Destructive staging-only rollback after backup/export and dependency review. |
| `README.md` | Target, security, backup/restore, retention, and activation guidance. |

## Tables aligned to the real v0.4 repository

The repository uses the following schema objects: `people`, `worlds`, `person_binding_subjects`, `person_profile_links`, `runtime_sources`, `records`, `memory_event_links`, `attention_allocations`, and `runtime_replays`.

`records` contains full JSONB payloads for v0.2/v0.3 canonical Event, Evidence, Memory, Claim, Goal, Constraint, PendingIntent, OpenLoop, AttentionBudget, Opportunity, Decision, Permission, ActionProposal, ActionGate, Outcome, SelfAudit, and ClaimLineage values. The relational person ID, kind, record ID, source key, profile binding, replay key, and composite references preserve tenant isolation and deterministic access patterns while retaining full canonical Go payloads.

## Transaction and integrity model

The repository begins a database transaction for each runtime operation, sets a local canonical-person context, acquires an advisory lock on the idempotency key, checks person-scoped parent references, writes append-preserving records, and commits only when all writes succeed. A unique `runtime_replays.idempotency_key` is the durable replay serialization point. Composite foreign keys tie memory-event links and attention allocation budget references to the same `person_id` in `records`.

Future production hardening should add a reviewed stored procedure or trigger layer for every graph edge represented in JSONB. The v0.4 repository already performs the required checks at the persistence boundary; the schema enforces the most critical relational link edges.

## RLS and existing security boundaries

RLS is enabled and forced for every `ci_kernel_v04` table. No policy is created for `anon`, `PUBLIC`, or a broad client-facing authenticated role. Only the no-login internal `ci_kernel_runtime` role receives policies, and each policy requires transaction-local person and stable-subject context.

The migration does not read, join, or trust the existing `profiles`, `feature_flags`, or `audit_events` tables. Public browser slugs are not identity inputs. Feature activation uses an application-side server activation/kill-switch port rather than the existing insecure feature-flag table.

## Backup, restore, retention, and deletion

Before any staging apply, take a database/schema export and record the target database identity. Test restore separately before considering a controlled activation. Conversation content, evidence provenance, derived memory, and audit data may be personal data; retention, deletion, consent, subject-access, legal-hold, encryption/key management, and backup-restoration rules require approval from the parallel Trust/Data Boundary effort before any non-disposable use.

`001_staging_down.sql` drops the complete schema and is only appropriate for disposable staging after the above checks. It is not a production rollback plan.

## Production exclusions

Do not use this package against production Supabase. The remaining prerequisites are durable server-session-to-stable-subject migration/backfill, reviewed internal role management, production RLS policy validation, privacy/consent policy integration, observability, backup/restore rehearsal, concurrency/load testing, controlled feature ownership, and approved HTTP composition. v0.4 does not activate an HTTP route, UI, action executor, push, merge, or deployment.
