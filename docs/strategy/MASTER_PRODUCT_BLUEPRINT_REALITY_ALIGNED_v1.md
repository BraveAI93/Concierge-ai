# Master Product Blueprint — Reality-Aligned v1

**AI:** Claude Code
**Mode:** Repo Execution Mode / Strategy Documentation Only
**Date:** 2026-07-09
**Status:** Strategy document. No code implemented, no code changed.

## Purpose

This is the "front door" strategy document for Brave by Bruno / The Concierge. It synthesizes four completed audits and two protocol documents into one reality-aligned picture of what the product actually is today, what it is meant to become, and what stands between the two. It does not replace the Master Document or Concierge Codex Audited v1 — both remain **SPEC_MISSING from this repository** (confirmed by repeated direct search across every audit in this series) — it works around that gap by being explicit about which claims are grounded in real repo text versus which are undocumented.

**Source documents (all already in this repo):**
- `docs/audits/PRODUCT_REALITY_MATRIX.md`
- `docs/audits/MASTER_VISION_VS_REPO_REALITY_REVISION.md` (supersedes `MASTER_VISION_VS_REPO_REALITY.md` for content purposes; both are kept for history)
- `docs/audits/TECHNICAL_CONSTRAINTS_MATRIX.md`
- `docs/decisions/AI_REVIEW_LEDGER_TRUST_AND_CINEMATIC_SHELL_v0.2.md`
- `docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md`
- `docs/protocol/PRODUCT_OPERATING_SYSTEM_v0.5.md`

## What the product actually is today (grounded)

Brave by Bruno's Concierge is, in its currently real form: a Next.js/Vercel frontend talking to a single Go backend on Render, backed by Supabase, that lets an independent professional (personal trainer, chef, photographer, etc.) create a public profile, receive AI-answered client enquiries **only in the 5 demo personas** (`/demo/*`), capture leads and booking requests, collect legal/consent forms, and manage all of it from an authenticated owner dashboard that also hosts a real conversational assistant (Brave PA). Auth is proven live (`AUTH GREEN`). Payments (Stripe) are built to a real standard at the code level but unverified live. Notifications for sensitive-topic chat alerts are real as of this week (durable record + email, truthful UI copy).

**The one fact that most needs to be said plainly:** the product's headline promise — a visitor can chat with a real AI concierge on a real professional's real public page — **does not work today**. The public-page chat widget is scripted HTML with no live AI call behind it (`docs/audits/PRODUCT_REALITY_MATRIX.md` #3, confirmed again with live production evidence in the Revision Audit). This is the single highest-priority gap in this entire blueprint.

## What the product is meant to become (per available specs — flagged where ungrounded)

Per `docs/protocol/PRODUCT_OPERATING_SYSTEM_v0.5.md` §2 (Source of Truth Hierarchy) and §7.6 (Codex-Locked Features Queue), the intended product includes, **as named and cleared by the Concierge Codex Audited v1** (full spec **SPEC_MISSING** from this repo — only these names and one-line clearances exist):
- **Trust Dot** — no external blocker, ready now, no further spec in-repo.
- **Brave Star** — 3 behavioural states, replaces a single decibel-line concept, no further spec in-repo.
- **Sfondo Location-Aware** — skyline + weather + time of day, no further spec in-repo.
- **Manager Agent** — incremental connectors, Google Calendar (Blocco A3) genuinely mid-implementation before being interrupted.
- **Offline-First Sync + Priority Reminder Fallback** — build after Manager Agent core.
- **Cinematic Shell** — Three.js Cosmic Intro, public visual layer; Integration Packet not yet produced.
- **Voice & Video Session Learning** — explicitly **LEGAL_LOCKED**, "stays in spec only," DPIA and GDPR/UK GDPR legal review required before any build, per the Codex's own binding block.

Bruno's own framing (introduced directly in this task series, **not found elsewhere in repo docs**, so marked **SPEC_MISSING** as a formal vision statement even though it is his stated direction): The Concierge as the operating engine of the ecosystem; Brave PA as an omnipresent adaptive guide; a voice-first interaction layer; a Generative Core / "Build Your Universe" onboarding experience. None of these framings currently exist in the product beyond a conventional form wizard and a dashboard-scoped assistant — see the Feature Registry (document 4) for the full accounting.

## Architecture pillars — activation-state snapshot

| Pillar | Current activation state | Evidence |
|---|---|---|
| Auth / session | ACTIVE_PUBLIC | Live-verified, AUTH GREEN |
| Demo AI chat | ACTIVE_PUBLIC | Live-verified this session (`/demo/bruno` → 200, real backend call) |
| Real-page AI chat | **GHOST_FORBIDDEN** | Confirmed scripted, no live call — see Feature Registry |
| Leads / CRM-lite | ACTIVE_PUBLIC | Real, tested |
| Booking requests | ACTIVE_PUBLIC | Real, with a known regex-detection fragility |
| Legal/consent forms | ACTIVE_PUBLIC | Real, persisted |
| General AI-processing consent | **GHOST_FORBIDDEN** | `sessionStorage`-only; linked privacy policy live-404s |
| Owner notifications (alert + email) | ACTIVE_PRIVATE | Real, built this week; live-config unconfirmed |
| Push notifications | GHOST_FORBIDDEN (candidate DORMANT_BUILT once built) | Listener exists, zero pipeline |
| Brave PA conversation | ACTIVE_PRIVATE | Real, owner-only |
| Brave PA claimed capabilities (calendar/search) | **GHOST_FORBIDDEN** | Prompt claims capability with zero backend |
| Media upload | ACTIVE_PUBLIC (pending bucket confirmation) | Real, UNKNOWN live-bucket status |
| Stripe payments | UNKNOWN | Real code, live-config unconfirmed |
| Trust Dot / Brave Star / Location-aware | SPEC_ONLY | Named, cleared, zero implementation |
| Cinematic Shell | SPEC_ONLY | Dependency installed, zero usage |
| Manager Agent / connectors | SPEC_ONLY (Calendar: closest to real) | One orphaned schema field |
| Memory / personalization | GHOST_FORBIDDEN (claim) / DORMANT_BUILT (storage) | Storage real, never re-read |
| Voice (any form) | LEGAL_LOCKED | Explicit, binding, 0/8 required safeguards present |
| Feature flags / activation states | GHOST_FORBIDDEN (as a system — doesn't exist) | No mechanism exists anywhere in repo |

Full per-feature detail lives in `docs/strategy/FEATURE_REGISTRY_AND_ACTIVATION_STATES_v1.md`.

## The core strategic move this blueprint makes

Two things are true at once, and this blueprint treats both as true rather than picking one: (1) the product's operational plumbing — auth, leads, bookings, legal forms, and now notifications — is genuinely solid, tested, and load-bearing; (2) the product's headline promise and several of its most vision-aligned features (Trust Dot, Brave Star, voice, real memory, a real Manager Agent) are honestly absent, most of them correctly so per the Codex's own gating. The strategic move is **not** "rush to build the vision" and **not** "just polish what exists" — it is to close the two live truth-gaps first (real chat, consent/privacy), then build the one piece of infrastructure (feature flags/activation states) that lets everything else be shipped honestly as dormant/coming-soon rather than either hidden or overclaimed, and only then move toward the more ambitious vision items in a sequence that never lets marketing outrun product truth (Anti-Chaos Rule #11, already binding in this repo).

## How to use this document set

- **This document** — orientation, "what is real, what is meant, what's the gap."
- `REAL_BUILD_ROADMAP_v1.md` — the phased, ordered task list to close the gap, Phase 0 through Phase 10.
- `IMPLEMENTATION_PACKETS_INDEX_v1.md` — a flat, quick-reference index of every packet, cross-referenced to its phase and source audit.
- `FEATURE_REGISTRY_AND_ACTIVATION_STATES_v1.md` — the full feature-by-feature ledger, the seed data for the real feature-flags system once built in Phase 1.
- `EXTERNAL_DEPENDENCIES_COSTS_AND_LEGAL_READINESS_v1.md` — every external account, env var name, cost exposure, and legal/DPIA requirement in one place.

None of these documents authorize implementation by themselves. Each roadmap task in document 2 still requires Bruno's explicit go-ahead before Claude Code touches code, per the existing Product Operating System's own Diamond Protocol.
