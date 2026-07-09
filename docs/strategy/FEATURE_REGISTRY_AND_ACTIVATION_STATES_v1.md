# Feature Registry and Activation States v1

**AI:** Claude Code
**Mode:** Repo Execution Mode / Strategy Documentation Only
**Date:** 2026-07-09
**Status:** Strategy document. No code implemented, no code changed.

This is the master feature ledger for Brave by Bruno / The Concierge, compiled from all four audits. It is the intended seed data for the real `feature_flags` table to be built in Roadmap Phase 1, Task 1.1. Every row is grounded in specific audit evidence — nothing here is asserted without a source.

**Activation states used:** `ACTIVE_PUBLIC`, `ACTIVE_PRIVATE`, `DORMANT_BUILT`, `SPEC_ONLY`, `LEGAL_LOCKED`, `GHOST_FORBIDDEN`, `UNKNOWN`.

Note on vocabulary: prior audits in this repo used a plain `GHOST` label. This document adopts the stricter `GHOST_FORBIDDEN` term specified for this document set — it carries an explicit instruction, not just a description: **features in this state must not be marketed, and their misleading UI/copy should be treated as an active defect, not a backlog item.**

---

## Registry

| ID | Feature | Current state | Why | Evidence |
|---|---|---|---|---|
| REG-01 | Auth / session (login, dashboard gating) | `ACTIVE_PUBLIC` | Live-verified, real `Set-Cookie`, real token match, real bypass block | AUTH GREEN evidence; `PRODUCT_REALITY_MATRIX.md` #16 |
| REG-02 | Demo AI chat (`/demo/*`) | `ACTIVE_PUBLIC` | Real `<Chat/>`, real backend call, live-verified this session (`/demo/bruno` → 200) | `PRODUCT_REALITY_MATRIX.md` #4; `MASTER_VISION_VS_REPO_REALITY_REVISION.md` live evidence |
| REG-03 | Real AI chat on real (non-demo) public profile pages | **`GHOST_FORBIDDEN`** | Scripted bubble, no live AI call, dead-end links to a stale domain | `PRODUCT_REALITY_MATRIX.md` #3; confirmed with live production evidence in the Revision Audit |
| REG-04 | Public profile page data rendering (services, contact links) | `ACTIVE_PUBLIC` | Real `GET /profile/:slug` data, correctly rendered | `PRODUCT_REALITY_MATRIX.md` #3 |
| REG-05 | Leads / CRM-lite + hot-lead email | `ACTIVE_PUBLIC` | Real deterministic scoring, real Resend email | `PRODUCT_REALITY_MATRIX.md` #13 |
| REG-06 | Booking-request pipeline | `ACTIVE_PUBLIC` | Real end-to-end, known regex-detection fragility (see REG-06a) | `PRODUCT_REALITY_MATRIX.md` #12 |
| REG-06a | Booking-intent detection (regex-based) | `ACTIVE_PUBLIC` with known fragility | Detects via string match on AI reply text, not a structured tool call | `TECHNICAL_CONSTRAINTS_MATRIX.md` O |
| REG-07 | Legal/consent forms (health, injury, image-rights, etc.) | `ACTIVE_PUBLIC` | Real submission, real persistence, real owner visibility | `PRODUCT_REALITY_MATRIX.md` #11 |
| REG-07a | Health-category form DPIA coverage | `UNKNOWN` (gap) | Special-category data collected without a documented DPIA-trigger review | `TECHNICAL_CONSTRAINTS_MATRIX.md` W |
| REG-08 | General AI-processing consent (chat consent screen) | **`GHOST_FORBIDDEN`** | `sessionStorage`-only, never persisted server-side, cites GDPR compliance falsely | `PRODUCT_REALITY_MATRIX.md` #11 |
| REG-09 | Privacy policy / Terms of Service page (`/privacy`) | **`GHOST_FORBIDDEN`** | Live-confirmed 404 in production; both legal disclosures unreachable | `MASTER_VISION_VS_REPO_REALITY_REVISION.md` I, live evidence |
| REG-10 | Owner dashboard (leads, bookings, notes, forms) | `ACTIVE_PUBLIC` | Extensive real CRUD | `PRODUCT_REALITY_MATRIX.md` #6 |
| REG-10a | Owner dashboard share-link / QR domain | `PARTIAL-REAL` (violates repo's own Anti-Chaos Rule #17) | Hardcoded to a Vercel alias, not the production domain; live-confirmed as a functioning duplicate, not dead, but still a rule violation | `MASTER_VISION_VS_REPO_REALITY_REVISION.md` live evidence |
| REG-10b | Owner dashboard silent data-load failures | `PARTIAL-REAL` | Fetch failures silently fall back to empty state instead of visible error | `PRODUCT_REALITY_MATRIX.md` #6 |
| REG-11 | Onboarding (profile/account creation) | `ACTIVE_PUBLIC` | Real save flow, live-verified account creation | `PRODUCT_REALITY_MATRIX.md` #10 |
| REG-11a | Onboarding "Generative Core / Build Your Universe" branding | **`GHOST_FORBIDDEN`** | Underlying mechanism is a conventional wizard; the branding language doesn't exist anywhere in the product | `MASTER_VISION_VS_REPO_REALITY_REVISION.md` F |
| REG-12 | Media upload | `ACTIVE_PUBLIC` (pending bucket confirmation) | Real, two independent flows; live Supabase bucket policy unverified | `PRODUCT_REALITY_MATRIX.md` #14 |
| REG-13 | Owner notifications — sensitive-topic alert record + email | `ACTIVE_PRIVATE` | Real durable record + real email attempt with explicit status, built this week; live-config unconfirmed | `TECHNICAL_CONSTRAINTS_MATRIX.md` G |
| REG-14 | Push notifications | **`GHOST_FORBIDDEN`** | Client-side listener exists, zero delivery pipeline, nothing can trigger it | `TECHNICAL_CONSTRAINTS_MATRIX.md` H |
| REG-15 | Notification preferences (`notifPrefs`) | **`GHOST_FORBIDDEN`** | `localStorage`-only, consumed by nothing | `TECHNICAL_CONSTRAINTS_MATRIX.md` I |
| REG-16 | Sound settings | **`GHOST_FORBIDDEN`** | Toggle exists, zero code plays a sound anywhere in the repo | `TECHNICAL_CONSTRAINTS_MATRIX.md` J |
| REG-17 | News/insights digest ("daily" framing) | `PARTIAL-REAL` | Real generation, no automation exists; "daily" is a false framing over an on-demand mechanism | `TECHNICAL_CONSTRAINTS_MATRIX.md` T |
| REG-18 | Notification digest/clustering | `SPEC_ONLY` | No scheduler exists anywhere in the codebase | `TECHNICAL_CONSTRAINTS_MATRIX.md` K |
| REG-19 | Brave PA conversation | `ACTIVE_PRIVATE` | Real AI conversation, owner-only, shares the undifferentiated `/chat` endpoint | `PRODUCT_REALITY_MATRIX.md` #7 |
| REG-19a | Brave PA claimed capabilities (calendar, web search, ticket booking) | **`GHOST_FORBIDDEN`** | System prompt claims capabilities with zero backend integration | `MASTER_VISION_VS_REPO_REALITY_REVISION.md` D |
| REG-19b | Brave PA "omnipresent" framing | `PARTIAL-REAL` | Mounted on 2 of ~7 owner-facing routes; absent from public/visitor side entirely | `MASTER_VISION_VS_REPO_REALITY_REVISION.md` D |
| REG-20 | Server-owned prompt templates | `ACTIVE_PUBLIC` (chat works) with an unmanaged trust gap | Server trusts client-supplied `system_prompt` with no rebuild/validation | `BRAIN_SPINE_READINESS_AUDIT.md` §2; `TECHNICAL_CONSTRAINTS_MATRIX.md` D |
| REG-21 | Manager Agent connector architecture | `SPEC_ONLY` | No architecture exists; one orphaned schema field is the only trace of prior work | `TECHNICAL_CONSTRAINTS_MATRIX.md` F |
| REG-22 | Google Calendar connector (Blocco A3) | `SPEC_ONLY` (genuinely closest-to-done) | Mid-implementation before being interrupted; schema field ready, OAuth not built | `TECHNICAL_CONSTRAINTS_MATRIX.md` E |
| REG-23 | Trust Dot | `SPEC_ONLY` | Cleared by Codex, zero implementation, no spec beyond the name in this repo | `TECHNICAL_CONSTRAINTS_MATRIX.md` O |
| REG-24 | Brave Star (3 behavioural states) | `SPEC_ONLY` | Cleared by Codex, zero implementation, no spec beyond the name | `TECHNICAL_CONSTRAINTS_MATRIX.md` N |
| REG-25 | Location-aware background / weather / city ambience | `SPEC_ONLY` | Cleared by Codex, zero implementation | `TECHNICAL_CONSTRAINTS_MATRIX.md` P |
| REG-26 | Cinematic Shell | `SPEC_ONLY` | `three` dependency installed, zero usage anywhere in the app; Integration Packet not yet produced | `TECHNICAL_CONSTRAINTS_MATRIX.md` Q |
| REG-27 | The Concierge as ecosystem operating engine | `SPEC_ONLY` (effectively **SPEC_MISSING**) | No repo doc defines this framing structurally | `MASTER_VISION_VS_REPO_REALITY_REVISION.md` A |
| REG-28 | Voice-first interaction (no learning) | `LEGAL_LOCKED` | Treated conservatively pending scope resolution vs. REG-29 | `MASTER_VISION_VS_REPO_REALITY_REVISION.md` C |
| REG-29 | Voice & Video Session Learning | `LEGAL_LOCKED` | Explicit, binding block; 0/8 required legal safeguards present | `PRODUCT_OPERATING_SYSTEM_v0.5.md` line 723; classification confirmed in Revision Audit |
| REG-30 | Memory / conversation history storage | `DORMANT_BUILT` | Real tables (`conversations`/`messages`), never re-read into future prompts | `BRAIN_SPINE_READINESS_AUDIT.md` §8C |
| REG-30a | "Personalization" / "remembers you" claim | **`GHOST_FORBIDDEN`** | Brave PA's own prompt implies persistent understanding that doesn't exist | `TECHNICAL_CONSTRAINTS_MATRIX.md` S |
| REG-31 | Advanced/semantic memory (vector-based) | `SPEC_ONLY` | Explicitly V2/V3, no infrastructure exists | `BRAIN_SPINE_READINESS_AUDIT.md` §11 |
| REG-32 | Stripe payments (Connect, checkout, webhook) | `UNKNOWN` | Code-level REAL (real signature verification, real fee calc); live env-var status unverifiable from this environment | `TECHNICAL_CONSTRAINTS_MATRIX.md` U |
| REG-33 | Minimal Brain Spine (architecture layer) | `SPEC_ONLY` | Fully scoped in its own audit, explicitly deferred until after Cinematic Shell + Core Flow QA, pending Bruno approval | `BRAIN_SPINE_READINESS_AUDIT.md` §15/§17 |
| REG-34 | Feature registry / feature flags / activation states (this system) | **`GHOST_FORBIDDEN`** (as a system — doesn't exist) | No flag/toggle mechanism exists anywhere in the codebase | `TECHNICAL_CONSTRAINTS_MATRIX.md` X |
| REG-35 | Audit events / provenance / trust logs | `SPEC_ONLY` | The notification work is a real, working instance of the pattern, not yet generalized | `TECHNICAL_CONSTRAINTS_MATRIX.md` Y |
| REG-36 | Root landing page (`/`) | `ACTIVE_PUBLIC` | Honest placeholder, doesn't overclaim | `MASTER_VISION_VS_REPO_REALITY_REVISION.md` A |
| REG-37 | `/theconcierge` home | `ACTIVE_PUBLIC` | Fully functional router; one dead legacy `localStorage` check, harmless | `MASTER_VISION_VS_REPO_REALITY_REVISION.md` context |

---

## Counts by state (as of this document's creation)

| State | Count |
|---|---|
| `ACTIVE_PUBLIC` | 14 |
| `ACTIVE_PRIVATE` | 2 |
| `DORMANT_BUILT` | 1 |
| `SPEC_ONLY` | 12 |
| `LEGAL_LOCKED` | 2 |
| `GHOST_FORBIDDEN` | 9 |
| `UNKNOWN` | 2 |

Nine `GHOST_FORBIDDEN` entries is the single most important number in this document — each one represents a live or recently-live claim the product makes that isn't structurally backed. Closing REG-03, REG-08, and REG-09 (Roadmap Phase 2) resolves three of the nine and the three with the highest live-production stakes.

## How this feeds Roadmap Phase 1

Task 1.3 in `docs/strategy/REAL_BUILD_ROADMAP_v1.md` seeds the real `feature_flags` table directly from this table's `Current state` column. As each roadmap task ships, this document (and the live table it seeds) should be updated in the same commit that closes the corresponding task — this registry is meant to stay a living, accurate ledger, not a one-time snapshot.
