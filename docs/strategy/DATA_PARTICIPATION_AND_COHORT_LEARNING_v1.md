# Data Participation and Cohort Learning v1

**AI:** Claude Code
**Mode:** Repo Execution Mode / Strategy Documentation Only
**Date:** 2026-07-09
**Status:** Strategy document. No code implemented, no code changed. No participation mode described below exists in the repo today — this is new design, cross-referenced against `docs/strategy/COHORT_LEARNING_AND_ANONYMOUS_PRODUCT_INTELLIGENCE_v1.md`, which remains the fuller treatment of the cohort-learning pipeline specifically. This document adds the **Data Participation Mode** framing on top — a per-user/per-owner tier system that determines which data behaviors are even allowed to apply to them.

## Purpose

Every claim the assistant makes about memory, personalization, or learning must be gated by which Data Participation Mode is actually active for that specific user or owner. This document defines those modes precisely, so that `docs/strategy/ASSISTANT_VERSION_LADDER_AND_KNOWLEDGE_PACKS_v1.md`'s V4 (Memory/Personalization) and beyond can never silently apply to someone who hasn't chosen that mode.

---

## Data Participation Modes

### Essential / No Memory
- **What it means:** the assistant answers using only the business's static profile data (services, pricing, tone) already provided for that conversation — nothing about this specific visitor is stored or reused beyond the current message exchange needed to generate a reply.
- **Applies to:** the default mode for every anonymous public-page visitor, unless and until they explicitly opt into something richer.
- **Data stored:** none beyond what's operationally required to serve the current reply (and the already-real, separately-consented lead/booking-request data if the visitor explicitly submits it).
- **Activation state:** `ACTIVE_PUBLIC` — this should be the default, safe, always-available mode.
- **Assistant constraint:** must never say "I remember," "as we discussed before," or imply continuity across visits.

### Session Only
- **What it means:** within a single browser session (the existing `cai_session` identifier already used for rate-limiting/tracking a conversation thread), the assistant can reference earlier messages in *that same session* — this is not new memory, it's the same behavior `Chat.jsx` already exhibits within one open conversation today, named explicitly here so it's not confused with cross-session memory.
- **Applies to:** every visitor, automatically, for the duration of one open conversation — this is not a new opt-in, it's naming existing real behavior.
- **Data stored:** the existing `conversations`/`messages` rows already written today — no new storage.
- **Activation state:** `ACTIVE_PUBLIC` — already real.
- **Assistant constraint:** must not reference anything from a *different* session, even for the same visitor, unless Personal Memory mode (below) is active.

### Personal Memory
- **What it means:** a specific, identified user (most realistically: an owner using Brave PA, or a returning client who has created some form of account/identity) explicitly opts in to the assistant recalling their own prior conversations, preferences, or history across sessions.
- **Applies to:** only users who have taken an explicit, separate consent action for this specific mode — never a default, never inferred from usage.
- **Data stored:** re-read access to that specific user's own `conversations`/`messages`/preference data (this is `docs/strategy/ASSISTANT_VERSION_LADDER_AND_KNOWLEDGE_PACKS_v1.md`'s V4 in practice).
- **Legal/privacy requirements:** explicit consent, a real way to view/export/delete the stored memory (currently does not exist anywhere in the product for any data type — see `docs/strategy/EXTERNAL_DEPENDENCIES_COSTS_AND_LEGAL_READINESS_v1.md` §6 for the same gap already flagged for GDPR rights generally).
- **Activation state:** `GHOST_FORBIDDEN` today (Brave PA's prompt already implies this mode is active for everyone, which it isn't) → `DORMANT_BUILT` once V4 is built and unmarketed → `ACTIVE_PRIVATE` once proven and genuinely opt-in.
- **Assistant constraint:** must only ever say "I remember" to a user who has actually opted into this mode — this is the single most important rule in this document, and directly supersedes anything a system prompt might otherwise imply.

### Connected Assistant
- **What it means:** the data-participation dimension of `docs/strategy/ASSISTANT_VERSION_LADDER_AND_KNOWLEDGE_PACKS_v1.md`'s V2/V3 — the owner has explicitly connected an external service (Calendar, email, etc.) and consented to the assistant reading/acting on that connected data.
- **Applies to:** owners only, per-connector, opt-in.
- **Data stored:** OAuth tokens and whatever data the connector reads (e.g., calendar events) — scoped to what that specific connector needs, not a general data grant.
- **Legal/privacy requirements:** standard OAuth consent screen per connector; privacy notice disclosure of each processor involved.
- **Activation state:** `SPEC_ONLY` today → `DORMANT_BUILT` (beta) → `ACTIVE_PRIVATE` per connector, exactly as scoped in the Version Ladder document.
- **Assistant constraint:** must never claim access to a connector the owner hasn't actually connected — this is the direct fix for the currently-false calendar/search claims in Brave PA's system prompt.

### Sensitive / Coach / Body / Health
- **What it means:** any mode where the assistant would track, reference, or act on health, body, fitness, or other special-category-adjacent data over time (the Version Ladder document's Workout/Body and Health/Benefits Knowledge Packs, made real).
- **Applies to:** no one today — this mode does not exist in any active form.
- **Data stored:** N/A until legally scoped.
- **Legal/privacy requirements:** requires its own DPIA, separate from and in addition to the general Health-Data DPIA Scoping already planned for the existing legal forms (`docs/strategy/REAL_BUILD_ROADMAP_v1.md` Task 7.2) — ongoing tracking/learning from health-adjacent data over time is a materially different, higher-risk activity than one-time disclosure collection, and should not be assumed to be covered by the same review.
- **Activation state:** **`LEGAL_LOCKED`** by default — same tier as Voice/Video below, not merely "not built yet."
- **Assistant constraint:** must never imply this mode is active, and the existing `GUARDRAILS` text (`lib/constants.js`) already correctly forbids medical advice regardless — this document reinforces that boundary rather than loosening it.

### Voice/Video Session Learning
- **What it means:** exactly `docs/strategy/ASSISTANT_VERSION_LADDER_AND_KNOWLEDGE_PACKS_v1.md`'s V6 — genuine learning from recorded voice/video sessions over time.
- **Applies to:** no one — not built, not to be scheduled into any sprint.
- **Legal/privacy requirements:** full DPIA, all 8 safeguards, per `docs/protocol/PRODUCT_OPERATING_SYSTEM_v0.5.md` line 723's binding instruction.
- **Activation state:** **`LEGAL_LOCKED`**, explicit, non-negotiable.
- **Assistant constraint:** N/A — this mode must never be reachable by any user under any circumstance until the DPIA (`docs/strategy/REAL_BUILD_ROADMAP_v1.md` Task 10.2) is complete.

---

## Anonymous / Cohort Learning — separate from personal memory

Cohort learning is **not a Data Participation Mode a user opts into individually** — it is a separate, product-level system that operates on aggregated, anonymous statistics across many users, and must remain structurally isolated from every mode above. A user in "Essential / No Memory" mode and a user in "Personal Memory" mode can both, in principle, contribute anonymous operational signals to cohort learning — participation in cohort learning is not the same axis as, and must never be conflated with, which personal-memory mode they've chosen.

**This section summarizes; `docs/strategy/COHORT_LEARNING_AND_ANONYMOUS_PRODUCT_INTELLIGENCE_v1.md` is the authoritative, fuller treatment — consult it for the full pipeline design, retention rules, audit events, and the open legal-basis question.**

### Allowed anonymous operational signals
Business category, coarse visitor intent (pricing/availability/general/booking-ready), interface language, tone variant used, coarse funnel stage reached, booking outcome, coarse conversation length (a bucket — e.g., "short/medium/long," never an exact turn count tied to a specific session), and error type (e.g., "AI service error," "rate limit hit" — operational health signals, not content). All must be segment-labeled and stripped of anything identifying, per the full detail in the companion document's §3 and §7.

### Forbidden / default-avoid signals
Full or partial transcripts, names, emails, phone numbers, health/body data, any sensitive/protected trait (age, sex/gender, health, sexuality, ethnicity, disability, religion/politics) as a segmentation axis or inference target — none of these may be used for demographic targeting or cohort segmentation unless a specific, freestanding future need is explicitly justified and separately consented, which is not authorized by this document and would require its own legal review, not an extension of this framework. Full detail in the companion document's §4 and §8.

### Minimum cohort thresholds
No cohort statistic may be computed, stored, or exposed for any combined segment representing fewer than **k=20 distinct underlying businesses/profiles** — set stricter than the common k≥5-10 baseline specifically because of this product's current small scale, where naive cohorts could trivially resolve to one or two real, identifiable businesses. Full detail and rationale in the companion document's §9.

### Rule: pseudonymised data is not anonymous
Data that is stripped of direct identifiers but still traceable back to a specific conversation, session, or business (even indirectly, even by combining fields) remains **personal data under GDPR**, not anonymous data. Only the final, k-anonymous aggregated statistics — never the pre-aggregation raw signals — may be treated as genuinely anonymous. This distinction is the reason the companion document flags an open legal-basis question (§14 there) for the raw-signal collection stage specifically, rather than assuming the whole pipeline is consent-free by default.

### Rule: the assistant must not claim to remember or personalize unless the relevant mode is active
This is the single rule that ties both documents together. Cohort learning may change the *product* over time (better default prompts, better routing, better copy) — it must never cause the assistant to say or imply anything to an individual user that depends on data about *them specifically*, unless their own Personal Memory (or Connected Assistant, or Sensitive/Coach) mode is genuinely active. A user in Essential/No Memory mode who happens to be statistically typical of a cohort must never hear "people like you usually..." — that would leak cohort-derived inference into an individual interaction in a way this framework does not permit.

---

## Build order and activation states

| Step | Deliverable | Activation state target | Depends on |
|---|---|---|---|
| 1 | Data Participation Mode field on user/session/owner records | Infrastructure (no state itself) | `docs/strategy/REAL_BUILD_ROADMAP_v1.md` Phase 1 (feature flags) |
| 2 | Essential / No Memory enforced as true default | `ACTIVE_PUBLIC` | Step 1 |
| 3 | Session Only formally named/documented (no new code — already real) | `ACTIVE_PUBLIC` | None |
| 4 | Cohort learning Minimum Safe V1 (per companion document) | `ACTIVE_PRIVATE` | §14 legal basis determination (companion doc), Roadmap Phase 1 |
| 5 | Connected Assistant mode (per-connector consent gating) | `SPEC_ONLY` → `DORMANT_BUILT` → `ACTIVE_PRIVATE` | Roadmap Phase 3/8 (connectors) |
| 6 | Personal Memory mode (real opt-in, real view/export/delete controls) | `GHOST_FORBIDDEN` → `DORMANT_BUILT` → `ACTIVE_PRIVATE` | Roadmap Phase 9 (Brain Spine, V4) |
| 7 | Sensitive / Coach / Body / Health mode | `LEGAL_LOCKED` | Dedicated DPIA, separate from Task 7.2 |
| 8 | Voice/Video Session Learning mode | `LEGAL_LOCKED` | Roadmap Task 10.2 (full DPIA) |

No step above may be reordered to move a `LEGAL_LOCKED` item earlier than its stated legal dependency, and no step may claim a higher activation state than its actual build/verification status — per the same "no `GHOST_FORBIDDEN` capability may ship" discipline established in the companion cohort-learning document.

Stop.
