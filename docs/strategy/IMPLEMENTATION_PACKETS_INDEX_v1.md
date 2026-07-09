# Implementation Packets Index v1

**AI:** Claude Code
**Mode:** Repo Execution Mode / Strategy Documentation Only
**Date:** 2026-07-09
**Status:** Strategy document. No code implemented, no code changed.

Flat, quick-reference index of every implementation packet in `docs/strategy/REAL_BUILD_ROADMAP_v1.md`. Use this to find a packet's phase, status, and source evidence without reading the full roadmap. Status is `NOT_STARTED` for every packet as of this document's creation — nothing in this document set has been implemented.

| # | Packet | Phase | Tasks bundled | Activation state target | Status | Source audit |
|---|---|---|---|---|---|---|
| P1 | Feature Flags & Audit Events | 1 | 1.1, 1.2, 1.3 | Infrastructure (no user-facing state) | NOT_STARTED | `TECHNICAL_CONSTRAINTS_MATRIX.md` §X, §Y |
| P2 | Consent, Domain & Real Chat | 2 | 2.1, 2.2, 2.3 | `GHOST_FORBIDDEN` → `ACTIVE_PUBLIC` | NOT_STARTED | `PRODUCT_REALITY_MATRIX.md` #3, #6; `MASTER_VISION_VS_REPO_REALITY_REVISION.md` B, I; `TECHNICAL_CONSTRAINTS_MATRIX.md` A, B, C |
| P3 | Dashboard Reliability Surfacing | 2 | 2.4 | `PARTIAL-REAL` → `ACTIVE_PUBLIC` | NOT_STARTED | `PRODUCT_REALITY_MATRIX.md` #6 |
| P4 | Brave PA Capability Truth + Calendar Connector | 3 | 3.1, 3.2, 3.3 | `GHOST_FORBIDDEN` → `ACTIVE_PRIVATE`/`DORMANT_BUILT` | NOT_STARTED | `MASTER_VISION_VS_REPO_REALITY_REVISION.md` D, S; `TECHNICAL_CONSTRAINTS_MATRIX.md` D, E |
| P5 | Notification Maturity | 4 | 4.1, 4.2, 4.3, 4.4, 4.5 | `ACTIVE_PRIVATE` (confirmed) / `GHOST_FORBIDDEN` → `ACTIVE_PUBLIC`/`DORMANT_BUILT` | NOT_STARTED | `TECHNICAL_CONSTRAINTS_MATRIX.md` G, H, I, J, K |
| P6 | Codex Visual Trio Spec Request | 5 | 5.1 | N/A (produces docs) | NOT_STARTED | `MASTER_VISION_VS_REPO_REALITY_REVISION.md` E, N, O |
| P7 | Trust Dot / Brave Star / Location Ambient (dormant build) | 5 | 5.2, 5.3, 5.4 | `SPEC_ONLY` → `DORMANT_BUILT` | NOT_STARTED (blocked on P6) | `TECHNICAL_CONSTRAINTS_MATRIX.md` N, O, P |
| P8 | Cinematic Shell Integration Packet | 6 | 6.1 | N/A (produces docs) | NOT_STARTED (blocked on P2) | `PRODUCT_OPERATING_SYSTEM_v0.5.md` §1007/§1290 |
| P9 | Cinematic Shell Implementation | 6 | 6.2 | `SPEC_ONLY` → `ACTIVE_PUBLIC` | NOT_STARTED (blocked on P8, P2) | `TECHNICAL_CONSTRAINTS_MATRIX.md` Q |
| P10 | Live Config & Financial Verification | 7 | 7.1 | `UNKNOWN` → `ACTIVE_PUBLIC` | NOT_STARTED | `TECHNICAL_CONSTRAINTS_MATRIX.md` U, V; §10 |
| P11 | Health-Data DPIA Scoping | 7 | 7.2 | N/A (legal deliverable) | NOT_STARTED | `TECHNICAL_CONSTRAINTS_MATRIX.md` W; §5 |
| P12 | Booking-Intent Structured Detection | 7 | 7.3 | `ACTIVE_PUBLIC` (robustness only) | NOT_STARTED | `PRODUCT_REALITY_MATRIX.md` #12 |
| P13 | Pre-Launch QA Day | 7 | 7.4 | N/A (gate) | NOT_STARTED (blocked on P2, P5, P10) | `PRODUCT_OPERATING_SYSTEM_v0.5.md` §8 Phase L3 |
| P14 | Connector Architecture Generalization | 8 | 8.1 | `SPEC_ONLY` → `DORMANT_BUILT` | NOT_STARTED (blocked on P4) | `TECHNICAL_CONSTRAINTS_MATRIX.md` F |
| P15 | Incremental Connectors (email, search, ...) | 8 | 8.2 | `DORMANT_BUILT` per connector | NOT_STARTED (blocked on P14) | `TECHNICAL_CONSTRAINTS_MATRIX.md` F |
| P16 | Minimal Brain Spine | 9 | 9.1 | Internal architecture (no state shift) | NOT_STARTED (blocked on P9, P13, Bruno approval) | `BRAIN_SPINE_READINESS_AUDIT.md` §7-11 |
| P17 | Basic Memory Re-read | 9 | 9.2 | `GHOST_FORBIDDEN` → `DORMANT_BUILT` | NOT_STARTED (blocked on P16) | `TECHNICAL_CONSTRAINTS_MATRIX.md` S |
| P18 | Advanced Semantic Memory | 9 | 9.3 | `GHOST_FORBIDDEN` → legal-gated | NOT_STARTED (blocked on P17, DPIA) | `TECHNICAL_CONSTRAINTS_MATRIX.md` S |
| P19 | Voice Legal Scope Resolution | 10 | 10.1 | N/A (decision) | NOT_STARTED | `MASTER_VISION_VS_REPO_REALITY_REVISION.md` C, L |
| P20 | Voice/Video DPIA | 10 | 10.2 | N/A (legal deliverable) | NOT_STARTED (blocked on P19) | `MASTER_VISION_VS_REPO_REALITY_REVISION.md` L |
| P21 | Voice I/O Build (conditional) | 10 | 10.3 | `LEGAL_LOCKED` → conditional | DO_NOT_BUILD_YET | `TECHNICAL_CONSTRAINTS_MATRIX.md` L |
| P22 | Voice & Video Session Learning | 10 | 10.4 | `LEGAL_LOCKED` | DO_NOT_BUILD_YET | `PRODUCT_OPERATING_SYSTEM_v0.5.md` line 723 |

## Reading this table

- **Status** reflects implementation reality as of this document's creation, not intent — everything is `NOT_STARTED` because this document set is strategy-only.
- **"Blocked on"** notes are dependency reminders, not permission gates by themselves — every packet still requires Bruno's explicit go-ahead per the Diamond Protocol, regardless of whether its technical dependencies are satisfied.
- **P21 and P22** are marked `DO_NOT_BUILD_YET` rather than merely `NOT_STARTED` — this is a deliberate distinction. No engineering time should be spent scoping these further until their respective legal gates (P19, P20) complete.
- Packets P6, P8, P11, P19, P20 are **not engineering tasks** — they produce documents or decisions, and should not be assigned to Claude Code as implementation work.

## Recommended near-term sequencing (first 5, unchanged from the Technical Constraints Matrix's own recommendation)

1. P2 — Consent, Domain & Real Chat
2. P1 — Feature Flags & Audit Events
3. P4 — Brave PA Capability Truth + Calendar Connector
4. P5 — Notification Maturity
5. P10 — Live Config & Financial Verification (can run in parallel with 2-4, not strictly after)
