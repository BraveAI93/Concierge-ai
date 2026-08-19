# Central Intelligence Runtime v0.4 — Staging Persistence and Identity Design

> **Status:** staging-only architecture. It uses a disposable local PostgreSQL test target and contains no production Supabase configuration, route activation, UI, or executor.

## Postgres repository contract

`PostgresRuntimeRepository` implements `RuntimeRepository` and the v0.2 `kernel.KernelRepository` persistence ports. It stores full canonical records as append-preserving JSONB payloads, while retaining relational person IDs, record IDs, source idempotency keys, runtime replay IDs, profile bindings, and link references needed for isolation and query constraints. This preserves the exact Go-domain contracts while keeping schema reads/writes reviewable.

Every `RunInTransaction` call begins a database transaction, verifies that the canonical person exists, and binds all runtime transaction methods to that person ID. Writes check same-person parent references before inserting. Source lineage and replay keys are unique. Replay lookup locks the replay row with `FOR UPDATE`; competing ingestions serialize through a transaction-scoped advisory lock derived from the idempotency key. A failure rolls back the entire transaction.

The repository is tested only against `ci_kernel_v04_test` on a local PostgreSQL 16 process bound to a sandbox-only Unix socket and port `55432`. It never reads environment variables for Supabase URLs or service keys.

## Durable server identity boundary

The runtime receives a canonical `AuthenticatedPrincipal`, not a public profile slug. `ServerSessionIdentityAdapter` is the new production-shaped composition boundary:

```text
server-validated session token
  → ServerSessionSubjectLookup (must return immutable stable subject)
  → AuthenticatedPrincipal
  → PostgresIdentityResolver
  → PersonBinding + allowed internal SourceProfileID(s)
```

The existing legacy session table stores a token-to-slug relationship, not a demonstrated immutable account subject. v0.4 therefore does **not** invent a production resolver from that slug. The adapter boundary is implemented and tested with a server-side test lookup, but production wiring requires a reviewed migration/backfill from authenticated account/session identity to `person_binding_subjects` and `person_profile_links`.

A canonical person can have a primary source profile and explicit linked internal source profile IDs. This models profile switching/account linking without accepting a caller-selected profile. Every legacy source profile is checked against the binding’s server-loaded allowlist.

## Consent boundary

`DerivedMemoryConsentVerifier` is a narrow runtime port. Conversation-derived memory persistence invokes it before any transaction. `FailClosedConsentVerifier` rejects the operation when no verified result is supplied. The parallel P0 Trust/Data Boundary work remains authoritative: v0.4 neither creates a consent database nor declares a consent policy. Tests use an explicit allow-only verifier to prove the dependency injection path.

## Activation and kill switch

`RuntimeActivation` is an explicit server-side interface independent of the existing `feature_flags` table. `DisabledFeature()` remains the default. When a runtime service is constructed with an activation provider, both the static feature boundary and provider must allow execution. A false provider is a kill switch. No HTTP route instantiates the enabled composition.

## Migration/RLS alignment

The `ci_kernel_v04` package contains a locally applied staging schema, a rollback script, and a security/operations README. It aligns the real repository’s tables and constraints. RLS is enabled and forced on every kernel table; the local repository connects as a dedicated staging role granted controlled access. There is no anonymous, public, or broad authenticated policy. Production role policy, subject claim propagation, backup/restore, retention, and deletion require separately approved staging work.
