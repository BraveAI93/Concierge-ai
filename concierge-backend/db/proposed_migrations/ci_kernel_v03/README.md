# CI Kernel v0.3 Runtime Migration Package — Review Only

> **Not executed.** This package is a reviewable migration design for a future isolated `ci_kernel` schema. It is not registered with a migration runner, has not been sent to Supabase, and must not be applied without a separate staging approval.

## Purpose

The package turns the earlier commented storage proposal into a reviewable up/down migration pair for the first runtime vertical slice. It includes canonical person bindings, source-message idempotency, the v0.2 temporal/attention/freshness model, planning and action-gate lineage, append-preserving audit records, indexes, and RLS design notes.

## Identity mapping

`ci_kernel.person_bindings` maps a canonical `person_id` to a **server-resolved stable subject** and an internal Concierge `source_profile_id`. It does not accept a public profile slug and does not use browser-selected profile identifiers as an authorization target. A production importer must obtain `stable_subject` only after server-side session authentication, then use the binding to constrain source conversation ingestion.

Existing `profiles` and `sessions` are not joined from an anonymous/public client context. Any historical backfill must execute through a reviewed server-side job using service credentials with explicit profile-to-subject mapping. The production security findings around `profiles`, `feature_flags`, and `audit_events` are therefore treated as isolation constraints: this schema does not depend on public reads or writes to those tables.

## RLS design

The SQL package enables RLS on every `ci_kernel` table. It intentionally creates **no anonymous policy** and no broad `authenticated` policy. The expected controlled-staging pattern is a server-only service role or a dedicated internal database role that sets a verified canonical-person context inside a transaction. Public clients never query this schema directly.

Before activation, security review must decide whether to use a dedicated backend role, a security-definer function with strict subject checks, or a Supabase JWT claim mapping. Every child table must be protected by a person-scoped policy and direct writes to append/history tables must remain backend-only.

## Transactions and integrity

The up migration defines same-person foreign-reference strategy through person IDs on primary records, composite/trigger review markers, and a required transaction procedure for ingestion. The runtime must atomically record source idempotency, Event/Evidence/Claim/Memory, links, intent/open loop, allocation, opportunity/decision, proposal/gate, and replay result. The initial vertical slice simulates this contract in memory; it does not execute the SQL.

## Retention and deletion

Source message content, provenance references, evidence, and audit records may contain personal data. Retention periods, legal holds, data-subject deletion, derived-memory deletion, and aggregate/audit exceptions require privacy review before activation. The current design preserves history; a production deletion policy must define whether records are cryptographically erased, tombstoned, or retained under a legal basis. Cascading deletes are deliberately not enabled blindly.

## Rollback

`001_runtime_rollback.sql` drops only the isolated `ci_kernel` schema. It is appropriate only before data activation or after an approved backup/export and dependency audit. Rollback does not restore deleted source data and must never be run against production without a change-management plan.

## Activation blockers

The package remains blocked on a staging Postgres/Supabase project, reviewed stable-subject mapping, complete RLS implementation validation, retention/deletion policy, transaction function review, operational metrics, key management, controlled rollout procedures, and security sign-off for all profile/session/audit/feature-flag touchpoints.
