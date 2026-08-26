# Central Intelligence Recovered Lineage Attestation

**Recovered branch:** `recovery/central-intelligence-authoritative-v1`
**Repository:** `BraveAI93/Concierge-ai`
**Scope:** Central Intelligence Kernel/Runtime v0.1 through v0.5 source lineage only.
**Date:** 2026-08-25

## Attestation

The genuine repository anchor is `964ce2e246d544cf5cd9931c30c2819e42a0e3c7`. The verified v0.1–v0.4 recovery mbox was replayed from that anchor with recovery-specific committer metadata, generating new Git objects. The preserved v0.4.1 source snapshot was then restored in a separate recovery commit, and the verified v0.5 patch was applied and committed. The resulting v0.5 Central Intelligence source snapshot is byte-for-byte equivalent to the preserved frozen 37-file source manifest. [1] [2] [3]

The historical v0.4.1 commit object `d2b406f4729e3e14e121b49bf44ff1cd6724e469` remains **UNRECOVERED AFTER EXHAUSTIVE SEARCH**, as documented in the prior blocker record. It is neither present in this recovered lineage nor represented as a replacement object. [4]

> "The recovered lineage preserves verified Central Intelligence source content and behaviour. Where historical Git objects were unavailable, no replacement commit is represented as having the same historical object identity."

## Identity classification

| Classification | Identifier | Meaning |
|---|---|---|
| Genuine anchor | `964ce2e246d544cf5cd9931c30c2819e42a0e3c7` | Existing remote repository object used as the parent of the recovered lineage. |
| Historical reference | `ed9da46c20dcce5502f4e627fa94055e2a9d3b50` | Original v0.1 identifier; not reused. |
| Historical reference | `3ebbdc76a6ff2f356d722e98db666adb86d7ad1a` | Original v0.2 identifier; not reused. |
| Historical reference | `15b4de5a2addfcada7fee4391d23250578b4855f` | Original v0.3 identifier; not reused. |
| Historical reference | `ed64a8f5aea61fef8611f57362878e7927a7da25` | Original v0.4 identifier; not reused. |
| Historical and unrecoverable | `d2b406f4729e3e14e121b49bf44ff1cd6724e469` | Original v0.4.1 object; source-equivalent state recovered, object identity unavailable. |
| Historical reconstructed reference | `15a14f3cee1654f97db59f1d2b628d7a19a8f32e` | Earlier reconstructed v0.5 reference; not reused and not treated as canonical lineage. |

## Historical-to-recovered mapping

| Historical stage / reference | Recovered commit | Parent | Classification |
|---|---|---|---|
| Genuine anchor `964ce2e…` | `964ce2e246d544cf5cd9931c30c2819e42a0e3c7` | Existing repository ancestry | Genuine anchor; same existing object. |
| v0.1 `ed9da46…` | `f9c272af0c4aca33332e71427c1cea02b069457a` | `964ce2e246d544cf5cd9931c30c2819e42a0e3c7` | Replayed source change; new object identity. |
| v0.2 `3ebbdc7…` | `4a07952421720e2e21ce476d3f0144d82e446e11` | `f9c272af0c4aca33332e71427c1cea02b069457a` | Replayed source change; new object identity. |
| v0.3 `15b4de5…` | `bb291acd2cca49cc184f26d27116919d64d3d5d4` | `4a07952421720e2e21ce476d3f0144d82e446e11` | Replayed source change; new object identity. |
| v0.4 `ed64a8f…` | `8193bbbeeb0afb895c8d44202ea41257fd553002` | `bb291acd2cca49cc184f26d27116919d64d3d5d4` | Replayed source change; new object identity. |
| v0.4.1 `d2b406f…` | `5afc16bed4e53539241d67d614d9b12d0576d2c9` | `8193bbbeeb0afb895c8d44202ea41257fd553002` | Source-equivalent verified v0.4.1 state; **not** the unavailable historical object. |
| v0.5 reconstructed reference `15a14f3…` | `343795938581bf5e6e4e641503c97245bd92e29c` | `5afc16bed4e53539241d67d614d9b12d0576d2c9` | Verified patch application and source-equivalent v0.5 state; new recovered object identity. |

## Verified input and source-equivalence evidence

| Artifact | Verified SHA-256 / result |
|---|---|
| Final pre-lineage recovery package | `f8c2083b0318221615af28197f852160506808449977825a2530ad8749f01667` |
| Pre-re-anchor forensic archive | `cd733a52688cac9b4f37d698dbc7407f4a9843de4c762a116ec10c6443a57fd0` |
| v0.1–v0.4 recovery mbox | `3df240c86e37f5bdd01dbf5f90d4ed869e935419e418e7630c99cdf1244b3ed5` |
| v0.5 source patch | `3168d46f74bfd606ac8ab6b81463b0238c1d48a0dfc424b18e52aae03dd2871e`; 16 paths, 8 added, 8 modified |
| Frozen v0.5 source manifest | `6570d1d3e143bee7c5db8282797476a0c986373a27d5730a269e9538dda520a6`; 37 files |
| Recovered v0.5 comparison | Empty recursive diff; all 37 manifest entries verified |

The v0.5 restored content includes persistent conversational continuity contracts (`InteractionBlock`, `Thread`, `ContinuityLink`, `ThreadDelta`, `CurrentThreadState`, bounded retrieval, retrieval-depth policy, and `MustSurface` separation) and adaptive attunement contracts (baseline, observed signals, uncertain inferred state, reversible adaptation, episode/intervention/outcome lineage, context-scoped correlation patterns, decay, override, and anti-manipulation policy). These contracts are present in the frozen v0.5 source, associated deterministic primitives, runtime service/persistence adapters, and preserved v0.5 test files. [2] [5]

## Validation attestation

The recovered branch passed uncached full Go tests, the runtime race detector, static analysis, and `git diff --check`. Kernel coverage was 69.6%; runtime coverage was 64.6%. A fresh disposable PostgreSQL 16 database completed migration UP, PostgreSQL-backed tests, migration DOWN, and migration UP again. The final inventory was 11 tables, 23 RLS policies, 11 FORCE-RLS relations, and zero `PUBLIC` function-execution grants. The runtime and identity-provisioner privilege inventories were retained, and the explicit RLS, restart, rollback/replay, cross-person isolation, duplicate/concurrency, and continuity tests passed. The disposable instance was then stopped and verified inactive. [6]

## References

[1]: file:///home/ubuntu/CI_LINEAGE_INPUT_VERIFICATION_SUMMARY.txt "Verified recovery inputs"
[2]: file:///home/ubuntu/CI_LINEAGE_V05_SOURCE_EQUIVALENCE.txt "Recovered v0.5 source-equivalence proof"
[3]: file:///home/ubuntu/CI_LINEAGE_MBOX_REPLAY.txt "Recovered v0.1–v0.4 replay mapping"
[4]: file:///home/ubuntu/CI_V05_CANONICAL_REANCHOR_BLOCKER.md "Prior canonical-object blocker record"
[5]: file:///home/ubuntu/central-intelligence-recovered-lineage/concierge-backend/kernel/V05_CONTINUITY_ATTUNEMENT_DESIGN.md "v0.5 continuity and attunement design"
[6]: file:///home/ubuntu/CI_LINEAGE_PG_VALIDATION_SUMMARY.txt "Disposable PostgreSQL validation evidence"
