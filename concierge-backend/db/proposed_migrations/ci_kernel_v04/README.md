# Central Intelligence Runtime v0.4.1 — Disposable Staging Migration Package

> **Target used for validation:** `ci_kernel_v04_test` on a local PostgreSQL 16.14 server bound to the sandbox-only Unix socket `/home/ubuntu/ci-v04-pgsocket` and local port `55432`. This is not Supabase and has no production network endpoint.

## Package contents

| File | Purpose |
|---|---|
| `001_staging_up.sql` | Creates the isolated `ci_kernel_v04` schema, separate runtime/provisioner roles, canonical persistence tables, person-scoped keys, forced RLS, and table-specific grants. |
| `001_staging_down.sql` | Destructive staging-only rollback after backup/export and dependency review. |
| `README.md` | Target, identity authority, tenancy, RLS, least-privilege, backup/restore, retention, and activation guidance. |

## Repository Tables and Tenant-Safe Keys

The repository uses `people`, `worlds`, `person_binding_subjects`, `person_profile_links`, `runtime_sources`, `records`, `memory_event_links`, `attention_allocations`, and `runtime_replays`. `records` contains complete JSONB payloads for v0.2/v0.3 canonical Event, Evidence, Memory, Claim, Goal, Constraint, PendingIntent, OpenLoop, AttentionBudget, Opportunity, Decision, Permission, ActionProposal, ActionGate, Outcome, SelfAudit, and ClaimLineage values.

| Table | Identifier boundary in v0.4.1 |
|---|---|
| `person_binding_subjects` | `stable_subject` remains globally unique because it is a canonical immutable server identity. |
| `person_profile_links` | `source_profile_id` remains globally unique because an internal source profile maps to exactly one Person. |
| `runtime_sources` | Primary key is `(person_id, id)`; unique source message identity is `(person_id, message_id)`. |
| `records` | Primary key is `(record_kind, person_id, id)`. |
| `memory_event_links` | Primary key and both foreign references include `person_id`. |
| `attention_allocations` | Primary key is `(person_id, allocation_id)`. |
| `runtime_replays` | Primary key is `(person_id, idempotency_key)`. |

A logical source, canonical record, allocation, or replay identifier belonging to one Person therefore does not reserve, reveal, or collide with the same logical identifier in another Person's world.

## Identity Authority and Provisioning

`ci_kernel_runtime` is the ordinary ingestion/read role. It can resolve an existing stable-subject binding and read linked internal profile IDs, but it has no INSERT, UPDATE, or DELETE grant or RLS policy on `people`, `person_binding_subjects`, or `person_profile_links`.

`ci_kernel_identity_provisioner` is the separate high-trust identity-management role. `PostgresIdentityProvisioner` uses this boundary for **initial INSERT-only** creation of Person, PersonalWorld, stable-subject binding, and internal source-profile links. It has no UPDATE or DELETE authority. Consequently, a stable subject cannot be rebound, and an existing profile cannot be attached to a different Person, through normal runtime or provisioner operations. Any future rebinding/deletion workflow must be separately designed, reviewed, audited, and authorized.

The local migration executor receives membership to both no-login roles only to exercise the boundary in disposable tests. A real staging or production deployment must use distinct reviewed backend principals and must never grant the provisioner role to ordinary runtime credentials.

## Transaction, Integrity, and Replay Model

The runtime starts a transaction, sets canonical person and stable-subject context locally, acquires an advisory lock on the person-scoped idempotency key, checks same-person parent references, writes append-preserving records, and commits only when all writes succeed. The advisory lock serializes duplicate work; replay records need only read/insert authority and do not use `FOR UPDATE`.

Composite foreign keys tie memory-event links and attention allocation budget references to the same `person_id` in `records`. The repository performs corresponding person-bound checks before writes. Future production hardening should add a reviewed stored-procedure or trigger layer for every graph edge represented in JSONB.

## RLS, Functions, and Least Privilege

RLS is enabled and forced for every `ci_kernel_v04` table. There are no policies for `anon`, `PUBLIC`, or any broad client-facing role. Runtime policies are split into the required SELECT, INSERT, and World compatibility UPDATE operations; provisioner policies are INSERT-only for high-trust identity tables.

`PUBLIC` is explicitly revoked from the schema, every table, and every schema function. The runtime receives schema usage, execution of only the two context functions, SELECT on required runtime tables, INSERT on normal runtime-write tables, and UPDATE only on `worlds`. It receives no DELETE privilege. The provisioner receives only schema usage and INSERT on the four initial identity tables.

The migration does not read, join, or trust existing `profiles`, `feature_flags`, or `audit_events` tables. Public browser slugs are not identity inputs. Feature activation remains an application-side server activation/kill-switch port independent of legacy feature flags.

## Backup, Restore, Retention, and Deletion

Before any staging apply, take a database/schema export, record the target database identity, and test restore separately. Conversation content, evidence provenance, derived memory, and audit data may be personal data; retention, deletion, subject-access, legal-hold, encryption/key management, and backup-restoration requirements require the parallel Trust/Data Boundary approval before non-disposable use.

`001_staging_down.sql` drops the complete schema and is only appropriate for disposable staging after the above checks. It is not a production rollback plan.

## Production Exclusions

Do not use this package against production Supabase. Remaining prerequisites include durable server-session-to-stable-subject migration/backfill, separate reviewed database principals and secret distribution, production RLS/adversarial validation, P0 consent integration, privacy/retention governance, observability, backup/restore rehearsal, load testing, controlled feature ownership, and approved HTTP composition. v0.4.1 does not activate an HTTP route, UI, action executor, push, merge, or deployment.
