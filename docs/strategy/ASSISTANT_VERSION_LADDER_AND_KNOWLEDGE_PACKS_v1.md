# Assistant Version Ladder and Knowledge Packs v1

**AI:** Claude Code
**Mode:** Repo Execution Mode / Strategy Documentation Only
**Date:** 2026-07-09
**Status:** Strategy document. No code implemented, no code changed.

## Purpose

Defines a maturity ladder (V0-V6) for what "the AI assistant" in The Concierge actually is at each stage, mapped directly onto `docs/strategy/REAL_BUILD_ROADMAP_v1.md`'s phases and `docs/strategy/FEATURE_REGISTRY_AND_ACTIVATION_STATES_v1.md`'s activation states — so that at any point, it is possible to say precisely which version is live, which is next, and which claims are currently allowed. Also defines Knowledge Packs — modular domain-content units that inform an assistant's system prompt per professional type. The "Knowledge Pack" concept itself is **SPEC_MISSING** as a formally named idea anywhere in the repo prior to this document; each pack below is mapped to what already exists (or doesn't) in the codebase rather than invented fresh.

---

## Assistant Version Ladder

### V0 — Scripted Assistant
- **What it can do:** display a pre-written greeting and a fixed set of quick-reply buttons on a public profile page.
- **What it must not claim:** anything implying a live AI is responding — "AI assistant," "ask me anything," "available 24/7" are all currently false claims at this version.
- **Data needed:** none — fully static.
- **External services:** none.
- **Legal/privacy requirements:** none technically, but the *misrepresentation itself* is the concern, not data handling.
- **Activation state:** **`GHOST_FORBIDDEN`** — this is not a future version to build toward, it is **the current state of the public profile page chat bubble today** (`docs/strategy/FEATURE_REGISTRY_AND_ACTIVATION_STATES_v1.md` REG-03), and it must be replaced, not iterated on. Listed here only to mark the floor this ladder starts above.
- **Live test required:** N/A — this version should not exist in production once Roadmap Task 2.2 ships.
- **Build timing:** N/A — this is the defect being fixed, not a target.

### V1 — Real Chat Assistant
- **What it can do:** a real, live AI conversation about the specific business's services, pricing, and availability, reachable from both demo pages (already true) and real public profile pages (once Roadmap Task 2.2 ships).
- **What it must not claim:** memory across sessions, capabilities beyond answering from the business's own profile data, any tool action (booking, calendar, payment) beyond proposing next steps in conversation.
- **Data needed:** the business's `profile_data` (services, tone, sensitive topics) — already real and collected at onboarding.
- **External services:** Anthropic (already integrated).
- **Legal/privacy requirements:** the AI-processing consent must be real (Roadmap Task 2.1) before this version can honestly claim to be live anywhere.
- **Activation state:** `ACTIVE_PUBLIC` on demo pages today; target state for real pages after Roadmap Phase 2.
- **Live test required:** yes — a real (non-demo) slug must produce a real, business-specific AI reply in production.
- **Build timing:** Roadmap Phase 2 (before Cinematic Shell).

### V1.5 — Owner PA Truthful
- **What it can do:** everything Brave PA's real conversational core already does today (real Anthropic-backed conversation, owner-only) — but with every capability claim in its system prompt matching a real, working backend.
- **What it must not claim:** calendar management, web search, or ticket booking unless the specific backend for that claim exists (see V2 below) — this version is explicitly about **removing false claims**, not adding new capability.
- **Data needed:** the owner's business data already fetched into `businessData` (leads, services, conversation counts) — already real.
- **External services:** Anthropic (already integrated); none new for the truth-pass itself.
- **Legal/privacy requirements:** none new.
- **Activation state:** current state is `ACTIVE_PRIVATE` (conversation) blended with `GHOST_FORBIDDEN` (claimed capabilities) — this version's entire purpose is collapsing that into a single honest `ACTIVE_PRIVATE`.
- **Live test required:** yes — ask the assistant to perform each previously-claimed action and confirm the response now honestly reflects real capability (either performs it, if V2's Calendar connector has landed, or clearly says it cannot).
- **Build timing:** Roadmap Phase 3, Task 3.1 — before Cinematic Shell.

### V2 — Connected Assistant
- **What it can do:** real actions through real external connectors — starting with read-only Google Calendar (Roadmap Task 3.3), then incrementally more (email, search) per Roadmap Phase 8's "one connector at a time" sequencing.
- **What it must not claim:** any connector not yet built; "manages your whole calendar" when only read access exists; any connector for a service not yet integrated (e.g., must not claim ticket-booking ability without a real booking-API connector).
- **Data needed:** OAuth tokens per connector (e.g., `Profile.GoogleRefreshToken`, already an existing unused schema field).
- **External services:** Google Calendar API (first connector); others as added.
- **Legal/privacy requirements:** standard OAuth scope disclosure per connector; privacy notice updates (once `/privacy` is reachable, per Roadmap Task 2.1).
- **Activation state:** `SPEC_ONLY` today → `DORMANT_BUILT` (owner-opt-in beta) → `ACTIVE_PRIVATE` per connector as each is proven.
- **Live test required:** yes, per connector — a real OAuth token refresh and a real data read/write.
- **Build timing:** Roadmap Phase 3 (Calendar first) and Phase 8 (subsequent connectors) — Calendar before Cinematic Shell, later connectors after Core Flow QA.

### V3 — Operating Assistant
- **What it can do:** takes real, gated actions on the owner's behalf across connected services — not just reading a calendar but proposing and (with explicit approval) creating events, not just answering booking questions but executing a confirmed booking end-to-end through a structured tool-call contract (replacing the current regex-based detection, Roadmap Task 7.3) rather than free-text inference.
- **What it must not claim:** full autonomy — every consequential action must remain gated by an explicit approval step (this is the `ActionGate` concept already scoped in `docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md` §7-8E, not a new idea invented here).
- **Data needed:** everything from V2, plus a structured action-request/approval log.
- **External services:** whichever connectors from V2 are live at this point.
- **Legal/privacy requirements:** none new beyond V2's, but action logging should feed the audit-events system (`docs/strategy/REAL_BUILD_ROADMAP_v1.md` Task 1.2).
- **Activation state:** `SPEC_ONLY` — this version does not exist even in partial form today; it depends on the Minimal Brain Spine's `ActionGate` interface (Roadmap Phase 9, Task 9.1) and at least one mature connector (V2/Phase 8).
- **Live test required:** yes — a real gated action, proposed, approved, and executed, with a real audit-event trail.
- **Build timing:** after Roadmap Phase 8 (connectors) and Phase 9's Minimal Brain Spine (`ActionGate`) — not before.
- **Framing note:** this is the closest version to Bruno's "operating engine" framing from `docs/strategy/MASTER_PRODUCT_BLUEPRINT_REALITY_ALIGNED_v1.md`, but that framing itself remains **SPEC_MISSING** as a formal architecture beyond this description — V3 is this document's best-grounded interpretation, not a confirmed spec.

### V4 — Memory / Personalization Assistant
- **What it can do:** references a specific user's or business's own prior conversations/preferences across sessions — the actual, structural meaning of "personalization," as opposed to the false claim currently implied by Brave PA's system prompt (`docs/strategy/FEATURE_REGISTRY_AND_ACTIVATION_STATES_v1.md` REG-30a).
- **What it must not claim:** anything beyond what `docs/strategy/DATA_PARTICIPATION_AND_COHORT_LEARNING_v1.md`'s **Personal Memory mode** actually stores for that specific user — memory must never be implied for a user who has not opted into that data participation mode.
- **Data needed:** re-read access to `conversations`/`messages`/`leads` (already stored, never re-read today — `docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md` §8C), gated per-user by their chosen Data Participation Mode.
- **External services:** none for basic re-read; a vector/embeddings store only for advanced semantic memory (later sub-stage of this version).
- **Legal/privacy requirements:** proportional to depth — basic re-read of a business's own operational data is lower-risk; genuine behavioral/preference learning approaches profiling and likely needs a DPIA (`docs/strategy/DATA_PARTICIPATION_AND_COHORT_LEARNING_v1.md` covers this fully).
- **Activation state:** `GHOST_FORBIDDEN` (claim) / `DORMANT_BUILT` (storage) today → `DORMANT_BUILT` (real, unmarketed) → `ACTIVE_PRIVATE` once proven, per Roadmap Task 9.2.
- **Live test required:** yes — confirm the assistant demonstrably references real prior context for an opted-in user.
- **Build timing:** Roadmap Phase 9, explicitly after Cinematic Shell and Core Flow QA.

### V5 — Voice I/O Assistant
- **What it can do:** transient speech input/output with no session learning or long-term retention — **only if** Roadmap Task 10.1 (legal scope resolution) concludes this is genuinely separable from V6 below.
- **What it must not claim:** any voice capability at all until Task 10.1 clears, and even then, must not claim any learning/memory from voice sessions specifically (that would make it V6, not V5).
- **Data needed:** transient audio only, no persistence beyond the single exchange, per the "no learning" boundary this version depends on.
- **External services:** provider not chosen — options and constraints fully detailed in `docs/strategy/EXTERNAL_DEPENDENCIES_COSTS_AND_LEGAL_READINESS_v1.md` §2.
- **Legal/privacy requirements:** conditional DPIA depending on Task 10.1's outcome — this version is **`LEGAL_LOCKED`** by default until that resolution exists, same as V6, pending the scope-separation determination.
- **Activation state:** `LEGAL_LOCKED` (conservative default) → `DORMANT_BUILT` only if and after Task 10.1 clears it.
- **Live test required:** yes, once built.
- **Build timing:** Roadmap Phase 10, Task 10.3 — **do not build any part of this before Task 10.1 completes.**

### V6 — Session Learning / Coach Assistant
- **What it can do:** nothing today, and nothing until a full DPIA exists — this version represents genuine voice/video session learning (e.g., a fitness-coaching assistant that observes and learns from recorded sessions over time).
- **What it must not claim:** anything, under any circumstance, until Roadmap Task 10.2 (full DPIA) is complete.
- **Data needed:** not scoped — premature before legal clearance, per this document's own discipline of not speculating on technical constraints ahead of a legal gate (matching `docs/audits/TECHNICAL_CONSTRAINTS_MATRIX.md`'s own treatment of this exact feature).
- **External services:** a real speech/video processing provider — real ongoing cost, not selected.
- **Legal/privacy requirements:** full DPIA, all 8 safeguards specified in `docs/audits/MASTER_VISION_VS_REPO_REALITY_REVISION.md`'s classification (DPIA, explicit consent flow, privacy notice, retention policy, deletion/export controls, biometric/health-data safeguards, processor list, activation gating) — currently 0 of 8 present anywhere in this repo.
- **Activation state:** **`LEGAL_LOCKED`**, explicit and non-negotiable, per `docs/protocol/PRODUCT_OPERATING_SYSTEM_v0.5.md` line 723: "Voice & Video Session Learning stays in spec only... Do not schedule it into a sprint."
- **Live test required:** N/A until legally cleared.
- **Build timing:** last in the entire roadmap (Roadmap Task 10.4) — **do not build any part of this feature under any circumstance until the DPIA is complete.**

---

## Knowledge Packs

Modular domain-content units that could inform an assistant's system prompt per professional type. Each pack is mapped to what already exists in the repo (mostly the `modes`/onboarding-step fields in `lib/constants.js`) rather than treated as a new system to build from scratch — most of the *data collection* already exists; what doesn't exist is a formal "pack" abstraction layering behavior on top of it.

| Pack | What it can do | What it must not claim | Data needed | External services | Legal/privacy | Activation state | Live test | Build timing |
|---|---|---|---|---|---|---|---|---|
| **Brand** | Inform tone, voice, and identity framing in the system prompt | A "brand strategy" capability beyond reflecting what the owner entered | `name`, `tag`, `tone` fields — already real, collected at onboarding | None | None | `ACTIVE_PUBLIC` (as existing onboarding data) / `SPEC_ONLY` (as a formal "pack" abstraction) | N/A — already flows into prompts today | Already live; formalizing as a named pack is a documentation task, not a build task |
| **Business** | Answer questions about services, pricing, booking | Availability it doesn't actually have live access to (ties to V2/V3's calendar connector) | `services` array — already real | None | None | `ACTIVE_PUBLIC` | Already tested via demo chat | Already live |
| **Creator** | Answer questions about platforms, content niche, collaboration packages | Follower counts or engagement metrics it cannot verify live | `platforms`, `niche`, `collab`, `mediakit` fields — already exist in onboarding `STEPS` (`lib/constants.js`) | None | None | `ACTIVE_PUBLIC` (data exists) / `UNKNOWN` (whether it's meaningfully exercised in practice — not verified this session) | Should be tested with a real creator-mode profile | Already live if a profile uses "Creator" mode; formal pack framing is documentation only |
| **Performance/Casting** | Answer questions about showreel, equipment, fee range, availability for events | Booking confirmation without a real booking mechanism behind it | `showreel`, `equip_own`, `equip_need`, `fee`, `availability` fields — already exist in onboarding `STEPS` | None | None | `ACTIVE_PUBLIC` (data exists) | Should be tested with a real performer-mode profile | Already live if used; formal pack framing is documentation only |
| **Admin/Life** | Task/reminder tracking, briefings — the Brave PA capabilities already partially real in the owner dashboard | Cross-device sync or persistence beyond what's actually stored (task list persistence was not independently re-verified in this session — mark **UNKNOWN**, not assumed real) | Whatever task/reminder state exists in `BravePAv2.jsx` — **status not independently re-confirmed in this session, flag as UNKNOWN rather than asserting** | None currently; a real scheduler (Roadmap Phase 4/8) for anything time-based | None | `UNKNOWN` — requires verification before claiming real | Required before any claim is made about this pack | Verification task first, not a build task |
| **Travel/Collab** | Matching businesses/creators for collaboration, travel-related booking assistance | Any live capability — this does not exist in any form today | None currently collected | A travel/booking API — not selected, not scoped | None beyond standard data handling once built | **`GHOST_FORBIDDEN`** if ever implied today — must not be claimed; `SPEC_ONLY` as a future pack | N/A | Not scoped in the current roadmap — would need its own packet definition first |
| **Workout/Body** | Track workout history, body metrics, progress over time (for personal-training-style professionals) | Any current implementation — this touches health-adjacent data and does not exist today | Would require structured health/fitness data collection — none exists | None currently | **Requires DPIA-level review before any build** — body/fitness metrics are health-adjacent data under GDPR, same caution class as the existing health-disclosure legal forms (`docs/strategy/EXTERNAL_DEPENDENCIES_COSTS_AND_LEGAL_READINESS_v1.md` §6-7) | **`LEGAL_LOCKED`** by default, same caution tier as Health/Benefits below, until a real DPIA scoping exists specifically for this pack | N/A until legally scoped | Not scoped in the current roadmap — should be evaluated alongside Roadmap Task 7.2's DPIA work, not built ahead of it |
| **Health/Benefits** | Reflect health/medical disclosures already collected via legal forms (injury declarations, pregnancy disclosure, allergy declarations) back into the assistant's awareness for a given booking | Any AI-generated health advice, diagnosis, or recommendation — the existing `GUARDRAILS` text in `lib/constants.js` already explicitly forbids this, and this pack must not weaken that | Existing `consents`/`form_submissions` data — already real, already collected | None | **This is exactly the data covered by Roadmap Task 7.2's Health-Data DPIA Scoping** — this pack must not be built or expanded ahead of that review completing | `ACTIVE_PUBLIC` (underlying form collection) / **`LEGAL_LOCKED`** (any new AI-facing use of that data beyond the existing guardrail-constrained disclosure-awareness) | Required before any expansion beyond current form collection | Gated entirely behind Roadmap Task 7.2 |

**Cross-cutting rule for all packs:** a pack may only inform what the assistant says if the underlying data was collected with a legal basis that covers that use — this is the same discipline `docs/strategy/DATA_PARTICIPATION_AND_COHORT_LEARNING_v1.md` applies to personal memory and cohort learning, extended to domain-content packs. No pack should silently expand what data is used for beyond what the owner understood when they entered it during onboarding.

## How this ladder relates to the roadmap

| Version | Roadmap phase |
|---|---|
| V0 → V1 | Phase 2 (Truth-Critical Public Core) |
| V1.5 | Phase 3 (Assistant Capability Truth) |
| V2 | Phase 3 (Calendar) + Phase 8 (subsequent connectors) |
| V3 | After Phase 8 + Phase 9's `ActionGate` |
| V4 | Phase 9 (Brain Spine / Memory) |
| V5 | Phase 10, conditional on Task 10.1 |
| V6 | Phase 10, Task 10.4, DPIA-gated, last |
