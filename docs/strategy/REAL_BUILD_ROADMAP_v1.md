# Real Build Roadmap v1

**AI:** Claude Code
**Mode:** Repo Execution Mode / Strategy Documentation Only
**Date:** 2026-07-09
**Status:** Strategy document. No code implemented, no code changed. Every task below requires Bruno's explicit go-ahead before any implementation begins, per the existing Diamond Protocol (`docs/protocol/PRODUCT_OPERATING_SYSTEM_v0.5.md` §3).

This roadmap operationalizes `docs/strategy/MASTER_PRODUCT_BLUEPRINT_REALITY_ALIGNED_v1.md` into ordered, dependency-aware tasks. Every task is grounded in one or more of the four audits — nothing here invents a capability the audits didn't already establish. Activation states use: `ACTIVE_PUBLIC`, `ACTIVE_PRIVATE`, `DORMANT_BUILT`, `SPEC_ONLY`, `LEGAL_LOCKED`, `GHOST_FORBIDDEN`, `UNKNOWN`.

---

## Phase 0 — Documentation and reality baseline

**Status: already complete.** This phase's deliverable is the four audits plus this five-document strategy set.

### Task 0.1 — Reality baseline (retrospective entry)
1. **Objective:** establish a truthful, evidence-based picture of repo reality before any further building.
2. **Reason for order:** nothing downstream can be sequenced safely without it — every later phase cites this baseline.
3. **Dependencies:** none.
4. **Technical requirements:** none (read-only audit work).
5. **External services/accounts:** none.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** N/A (documentation, not a feature).
8. **Definition of done:** the four audits + five strategy documents exist and are pushed. **Met.**
9. **Live test required:** N/A.
10. **What not to touch:** N/A — already done, do not re-litigate without new evidence.

---

## Phase 1 — Product Operating Spine / Feature Flags / Activation States

Goal: build the one piece of infrastructure that lets every later phase ship features as honestly dormant instead of either hidden or overclaimed.

### Task 1.1 — Feature flags table + read/consume mechanism
1. **Objective:** a single source of truth for each feature's activation state, readable by frontend and backend.
2. **Reason for order:** every later "build it dormant" recommendation in this roadmap depends on this existing first (`docs/audits/TECHNICAL_CONSTRAINTS_MATRIX.md` §X).
3. **Dependencies:** none.
4. **Technical requirements:** new `feature_flags` table (name, state, updated_at) — this repo has no established migration pattern for any table, so this is the first; the exact SQL DDL must be reviewed with Bruno and run manually against Supabase, not improvised by Claude Code, exactly as the notification work's own precedent requires. One backend read endpoint. One frontend consumption hook.
5. **External services/accounts:** none — recommend building in-repo rather than adopting GrowthBook/LaunchDarkly (cost/complexity not justified at current scale).
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** N/A — this is infrastructure, not a user-facing feature to gate.
8. **Definition of done:** toggling a flag in the table visibly changes behavior in a real deployed environment.
9. **Live test required:** yes — toggle a test flag, confirm the change is observable.
10. **What not to touch:** auth, existing tables, any feature this system will later gate (build the mechanism first, migrate features to it after).

### Task 1.2 — Audit events / trust-log table
1. **Objective:** generalize the "real record + explicit status, never claim success without a backend signal" pattern already proven in this week's notification work into a reusable event log.
2. **Reason for order:** natural follow-on to 1.1; both are foundational infrastructure and share the "first new table in this repo" precedent-setting concern.
3. **Dependencies:** Task 1.1 (shares the schema-precedent discussion, can be reviewed together).
4. **Technical requirements:** new `audit_events` table; one shared logging function called from at least the existing notification path plus one additional claim (e.g., booking-request creation) to prove genuine reusability, not a one-off.
5. **External services/accounts:** none.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** N/A — infrastructure.
8. **Definition of done:** two independent product claims ("owner notified," "booking request sent") both write through the same reusable event function.
9. **Live test required:** yes.
10. **What not to touch:** the existing notification code's already-proven behavior — wrap/extend, don't rewrite.

### Task 1.3 — Seed the Feature Registry
1. **Objective:** populate the new `feature_flags` table with every feature from `docs/strategy/FEATURE_REGISTRY_AND_ACTIVATION_STATES_v1.md` at its currently audited state.
2. **Reason for order:** makes the registry document operationally real, not just a reference file.
3. **Dependencies:** Task 1.1.
4. **Technical requirements:** a one-time seed script/migration reflecting document 4's table.
5. **External services/accounts:** none.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** N/A.
8. **Definition of done:** every row in the registry document has a corresponding row in the live table.
9. **Live test required:** spot-check 3-5 rows against live behavior.
10. **What not to touch:** nothing beyond the new table.

---

## Phase 2 — Truth-Critical Public Core

Goal: close the two live, active truth-gaps identified across all four audits before anything else public-facing proceeds.

### Task 2.1 — Consent & Disclosure Truth Path
1. **Objective:** make `/privacy` reachable in production and persist the general AI-processing consent server-side.
2. **Reason for order:** **highest priority in the entire roadmap** — a live-confirmed 404 on a GDPR disclosure (`curl -L https://www.bravebybruno.com/privacy` → 404, confirmed in the Revision Audit) plus an unpersisted compliance claim is an active defect, not a future risk.
3. **Dependencies:** none — can start immediately.
4. **Technical requirements:** add a `public/` folder or convert `privacy.html` to a real Next.js page (`app/privacy/page.jsx` recommended, gets it into the standard build/routing pipeline); persist `Chat.jsx`'s `giveConsent()` via the existing `POST /consent` pattern (no new schema — reuse `consents` table with a value like `"ai_processing"`).
5. **External services/accounts:** none.
6. **Legal/privacy requirements:** recommend a real legal review of `privacy.html`'s *content* once reachable (separate from this technical fix).
7. **Activation state after build:** `GHOST_FORBIDDEN` → `ACTIVE_PUBLIC`.
8. **Definition of done:** `/privacy` returns 200 in production; a real consent click produces a real DB row.
9. **Live test required:** yes — `curl -L` the production URL; trigger a real consent click and confirm the DB row.
10. **What not to touch:** legal-form consent (`LegalFormModal.jsx`/`FormPage.jsx`) — already real, do not modify.

### Task 2.2 — Real Chat on Real Pages
1. **Objective:** replace the scripted "Chat with me" bubble on `/[slug]` with the real `<Chat/>` component (or equivalent) reachable for non-demo profiles.
2. **Reason for order:** the product's core promise doesn't function outside `/demo/*` today — this is the single most important functional gap across all four audits.
3. **Dependencies:** none technically, but should ship in the same pass as 2.3 (same root cause).
4. **Technical requirements:** mount `<Chat/>` on a real route for non-demo slugs; remove or repurpose the static bubble in `lib/generateBusinessPage.js`; reuses the existing `POST /chat` endpoint, no new backend work.
5. **External services/accounts:** none new — reuses existing Anthropic integration.
6. **Legal/privacy requirements:** none beyond the existing chat consent (now real per 2.1).
7. **Activation state after build:** `GHOST_FORBIDDEN` → `ACTIVE_PUBLIC`.
8. **Definition of done:** a real (non-demo) slug produces a real AI reply end-to-end in production.
9. **Live test required:** yes — visit a real profile URL, send a message, confirm a genuine Anthropic-backed reply.
10. **What not to touch:** the demo chat flow (`/demo/*`) — already real, do not modify; auth/routing beyond the new route itself.

### Task 2.3 — Stale domain/share-link correction
1. **Objective:** replace all 4 hardcoded `concierge-ai-gamma.vercel.app` occurrences with the real production domain.
2. **Reason for order:** shares the exact root cause as 2.2 (stale architecture assumptions); live-confirmed to violate the repo's own binding Anti-Chaos Rule #17.
3. **Dependencies:** none — can ship alongside or immediately after 2.2.
4. **Technical requirements:** string replacement in `generateBusinessPage.js`, `BusinessPagePreview.jsx`, `OwnerDashboard.jsx` (x3), `Onboarding.jsx`; optionally introduce a `NEXT_PUBLIC_APP_URL` env var to prevent recurrence.
5. **External services/accounts:** none.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** `PARTIAL-REAL` → `ACTIVE_PUBLIC`.
8. **Definition of done:** every share link/URL preview in the product resolves to `bravebybruno.com`.
9. **Live test required:** yes — copy a real share link from a real dashboard session.
10. **What not to touch:** the underlying profile/slug data — this is a display/link-generation fix only.

### Task 2.4 — Owner dashboard silent-failure surfacing
1. **Objective:** replace silent `.catch(() => setX([]))` fallbacks in `OwnerDashboard.jsx`'s initial data load with visible error states.
2. **Reason for order:** low cost, bundles naturally with 2.3 since both touch the same file's data-loading logic.
3. **Dependencies:** none.
4. **Technical requirements:** small frontend change only.
5. **External services/accounts:** none.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** `PARTIAL-REAL` → `ACTIVE_PUBLIC` (for the dashboard's reliability signaling specifically).
8. **Definition of done:** a forced fetch failure produces a visible error to the owner, not a silent empty state.
9. **Live test required:** yes.
10. **What not to touch:** the underlying `/owner/*` endpoints themselves — this is UI-only.

---

## Phase 3 — Assistant Capability Truth

Goal: stop Brave PA from claiming capabilities it doesn't have, and close the highest-leverage genuinely-in-progress gap (calendar).

### Task 3.1 — Brave PA capability truth pass
1. **Objective:** remove or gate the calendar/search/booking claims in `buildBravePAPrompt()` that have no backend behind them.
2. **Reason for order:** an active, ongoing misrepresentation (not a missing future feature) — Brave PA currently tells users it can do things it cannot.
3. **Dependencies:** ideally sequenced with 3.3 (Calendar connector), since that resolves the calendar claim by building it rather than removing it.
4. **Technical requirements:** edit `lib/buildPrompt.js`'s `buildBravePAPrompt()` to remove claims not backed by 3.3, or make them conditional on the flag system (Phase 1) once the real connector exists.
5. **External services/accounts:** none for the truth-pass itself.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** claimed capability `GHOST_FORBIDDEN` → either removed or `ACTIVE_PRIVATE` once backed by 3.3.
8. **Definition of done:** every capability claim in the system prompt is backed by real, working code.
9. **Live test required:** yes — ask Brave PA to do each previously-claimed action and confirm the response matches real capability.
10. **What not to touch:** Brave PA's core conversational behavior — already real, do not degrade.

### Task 3.2 — Server-owned prompt templates
1. **Objective:** stop trusting a client-supplied `system_prompt` string on every `/chat` call; rebuild it server-side from stored profile data.
2. **Reason for order:** a real trust/security gap on the busiest endpoint in the product (`docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md` §2), best fixed once the flag system (Phase 1) exists to allow a safe staged rollout.
3. **Dependencies:** Phase 1 (flag system, for staged rollout); should not block Phase 2.
4. **Technical requirements:** port `lib/buildPrompt.js`'s Concierge-chat prompt-building logic (or an equivalent) into the Go backend; `handleChat` accepts `profile_id` + `mode` instead of a full prompt string.
5. **External services/accounts:** none.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** `ACTIVE_PUBLIC` (unchanged from the user's perspective — this is a trust hardening, not a new feature).
8. **Definition of done:** golden-file test shows byte-identical replies for identical inputs before/after the change.
9. **Live test required:** yes — the same golden-file comparison run against production.
10. **What not to touch:** auth, the currently-working chat response shape, Brave PA's richer context-dependent prompt (leave for a follow-up).

### Task 3.3 — Google Calendar connector V1 (Manager Agent's first connector)
1. **Objective:** resume Blocco A3 — read-only Google Calendar OAuth — the genuinely closest-to-done unfinished feature in the repo.
2. **Reason for order:** per the Master Doc excerpt itself, this was scoped as "before or alongside Phase B1" — i.e., early; also the natural fix for 3.1's calendar claim.
3. **Dependencies:** none technically; pairs naturally with 3.1.
4. **Technical requirements:** OAuth callback handler + calendar-read endpoint in `main.go`; `Profile.GoogleRefreshToken` schema field already exists, unused — no migration needed.
5. **External services/accounts:** Google Cloud Console project — **confirm with Bruno whether the CLIENT_ID/CLIENT_SECRET from the interrupted prior attempt still exist before assuming a fresh setup is needed.**
6. **Env vars (names only):** `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_OAUTH_REDIRECT_URI`.
7. **Legal/privacy requirements:** update the (now-reachable, per 2.1) privacy notice to disclose Google as a processor and the calendar scope requested; standard OAuth consent screen, not DPIA-level.
8. **Activation state after build:** `SPEC_ONLY`/`GHOST_FORBIDDEN` (claim) → `DORMANT_BUILT` first (owner-opt-in beta), then `ACTIVE_PRIVATE` once proven reliable.
9. **Definition of done:** a real OAuth token successfully refreshes and a real event list is returned to Brave PA.
10. **Live test required:** yes.
11. **What not to touch:** auth (this is a separate OAuth flow, not the product's own login), existing profile data.

---

## Phase 4 — Notification Center foundation

Goal: mature the notification system already hardened this week into a genuinely complete, verified center.

### Task 4.1 — Live-config verification of notification email
1. **Objective:** confirm `RESEND_API_KEY`/`OWNER_EMAIL`/`RESEND_FROM_EMAIL` are actually set on Render.
2. **Reason for order:** the notification system was built and code-verified this week but never live-verified; this is the cheapest, highest-confidence task to close that gap.
3. **Dependencies:** none.
4. **Technical requirements:** none — verification only, no code change.
5. **External services/accounts:** Resend (already exists).
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** `ACTIVE_PRIVATE` (confirmed) — no state change, just confidence.
8. **Definition of done:** a real sensitive-topic message in production returns `email_status: "sent"`.
9. **Live test required:** yes, this task **is** the live test.
10. **What not to touch:** the notification code itself — this is a config check, not a code change.

### Task 4.2 — Notification preferences wiring
1. **Objective:** persist `notifPrefs` server-side and make it actually gate behavior.
2. **Reason for order:** cheap, contained, closes a concrete ghost-UI gap flagged in every audit in this series.
3. **Dependencies:** Phase 1 recommended (reuse the same schema-precedent discussion) but not strictly required — could reuse `profiles.profile_data` blob instead.
4. **Technical requirements:** small backend save path + frontend wiring.
5. **External services/accounts:** none.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** `GHOST_FORBIDDEN` → `ACTIVE_PUBLIC`.
8. **Definition of done:** toggling a preference measurably changes behavior (see 4.3).
9. **Live test required:** yes.
10. **What not to touch:** the underlying alert-creation logic — this only gates *display/delivery*, not detection.

### Task 4.3 — Sound settings wiring
1. **Objective:** make `soundEnabled`/`soundStyle` actually play a sound on real alert arrival.
2. **Reason for order:** bundles with 4.2, same root gap.
3. **Dependencies:** 4.2 (for the preference to check), Task 1.1's `public/` folder work (needed for an audio asset).
4. **Technical requirements:** one `new Audio(...).play()` call, gated by the preference.
5. **External services/accounts:** none.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** `GHOST_FORBIDDEN` → `ACTIVE_PUBLIC`.
8. **Definition of done:** a real alert triggers a real sound when enabled, no sound when disabled.
9. **Live test required:** yes.
10. **What not to touch:** the alert-creation logic itself.

### Task 4.4 — Push notification pipeline
1. **Objective:** register the existing (orphaned) `service_worker.js`, build a subscribe flow, and wire real Web Push delivery.
2. **Reason for order:** real infrastructure work, sequenced after the cheaper wins in this phase; Bruno has stated push is "essential to the long-term product," so it should not be indefinitely deferred.
3. **Dependencies:** a `public/` folder (also needed by 4.3); Phase 1 (ship dormant/owner-opt-in via flags).
4. **Technical requirements:** `POST /push/subscribe` endpoint; a subscription-storage mechanism (reuse the `notes` pattern with a new `note_type`, or a dedicated table — **needs a decision with Bruno, SPEC_MISSING on preference**); call a web-push send function from `handleCreateAlert`.
5. **External services/accounts:** none required — Web Push needs no paid provider, only self-generated VAPID keys.
6. **Env vars (names only):** `VAPID_PUBLIC_KEY`, `VAPID_PRIVATE_KEY`, `VAPID_SUBJECT` — must be generated fresh, never hardcoded.
7. **Legal/privacy requirements:** update privacy notice to disclose push-subscription data.
8. **Activation state after build:** `GHOST_FORBIDDEN` → `DORMANT_BUILT` (owner-opt-in, not marketed until proven reliable) → `ACTIVE_PRIVATE` once confirmed.
9. **Definition of done:** a real push notification is received on a real device.
10. **Live test required:** yes, on a real browser/device — cannot be verified by build/vet alone.
11. **What not to touch:** the existing alert-record/email path — push is additive, not a replacement.

### Task 4.5 — Digest/clustering scheduler
1. **Objective:** add a real scheduler and a batched digest email; also fixes the pre-existing broken `PUT /owner/digest-prefs` frontend call (confirmed to have no backend route).
2. **Reason for order:** last in this phase — depends on a scheduler decision that also benefits Phase 8 (daily news automation).
3. **Dependencies:** a scheduler choice — recommend Render Cron Job (native to existing hosting), pending confirmation of current Render plan tier's Cron Job availability.
4. **Technical requirements:** reuse `profiles.DigestFrequency`/`DigestLastSent` (already exist, unused — no migration needed); a new scheduled endpoint.
5. **External services/accounts:** possibly a Render plan upgrade, **COST_VERIFY_REQUIRED**.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** `SPEC_ONLY` → `DORMANT_BUILT` (manual-trigger V1) → `ACTIVE_PRIVATE` (scheduled, proven).
8. **Definition of done:** a real scheduled job fires and a real digest email is received.
9. **Live test required:** yes.
10. **What not to touch:** the individual alert-notification path — digest is additive/batched, not a replacement.

---

## Phase 5 — Visual Identity / Trust Dot / Brave Star / Location Ambient

Goal: build the three Codex-cleared visual-identity features — but only once real specs exist.

### Task 5.1 — Acquire Codex spec for Trust Dot, Brave Star, Sfondo
1. **Objective:** obtain the actual Concierge Codex Audited v1 spec for these three features — **this repo currently only has their names and clearance status, nothing else** (SPEC_MISSING).
2. **Reason for order:** must happen before any code in this phase — building from a one-line name risks recreating the "real-looking but wired to nothing" pattern flagged repeatedly across the audits.
3. **Dependencies:** none — this is a Bruno/ChatGPT/documentation task, not an engineering task.
4. **Technical requirements:** none.
5. **External services/accounts:** none.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** N/A — this task produces a document, not a feature.
8. **Definition of done:** the three specs are saved into this repo's `docs/` tree, following the same pattern as the existing architecture docs.
9. **Live test required:** N/A.
10. **What not to touch:** N/A.

### Task 5.2 — Trust Dot (dormant)
1. **Objective:** build a real, backend-driven trust indicator.
2. **Reason for order:** depends on Phase 1's audit-events system existing to have something genuine to indicate trust *about* — building it earlier risks decorative-only UI.
3. **Dependencies:** Task 1.2 (audit events), Task 5.1 (spec).
4. **Technical requirements:** per spec, once obtained.
5. **External services/accounts:** none apparent from current naming.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** `SPEC_ONLY` → `DORMANT_BUILT` (owner/internal preview) → `ACTIVE_PUBLIC` once validated.
8. **Definition of done:** per spec.
9. **Live test required:** yes.
10. **What not to touch:** N/A pending spec.

### Task 5.3 — Brave Star 3 behavioural states (dormant)
1. **Objective:** build the visual state-indicator replacing "the single decibel-line concept."
2. **Reason for order:** cleared "ready now" per Codex, independent of other Phase 5 items, but still blocked on 5.1's spec.
3. **Dependencies:** Task 5.1 (spec — unclear if this needs real backend signals or is purely client-side animation; the spec determines scope).
4. **Technical requirements:** per spec, once obtained.
5. **External services/accounts:** none apparent.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** `SPEC_ONLY` → `DORMANT_BUILT` → `ACTIVE_PUBLIC`.
8. **Definition of done:** per spec.
9. **Live test required:** visual QA.
10. **What not to touch:** N/A pending spec.

### Task 5.4 — Location-aware background (dormant)
1. **Objective:** skyline + weather + time-of-day rendering layer.
2. **Reason for order:** natural Cinematic Shell companion, can build in parallel with 5.2/5.3.
3. **Dependencies:** Task 5.1 (spec, lighter need here since the one-line description is more actionable than Trust Dot/Brave Star's).
4. **Technical requirements:** a thin backend proxy to the chosen weather API (avoid exposing a key client-side, if a keyed provider is chosen); reuse the profile's stored lat/lng from onboarding.
5. **External services/accounts:** weather provider — **provider not chosen, recommend Open-Meteo (genuinely free, no API key required, avoids secret management entirely)**; alternatives: OpenWeatherMap, WeatherAPI.com.
6. **Env vars (names only):** `WEATHER_API_KEY` only if a keyed provider is chosen; none if Open-Meteo.
7. **Legal/privacy requirements:** minor — disclose the weather lookup uses stored profile location, not live visitor geolocation.
8. **Activation state after build:** `SPEC_ONLY` → `DORMANT_BUILT` → `ACTIVE_PUBLIC`.
9. **Definition of done:** a real skyline/weather/time-of-day background renders based on real profile location data.
10. **Live test required:** visual QA.
11. **What not to touch:** onboarding's existing location-capture flow — reuse, don't rebuild.

---

## Phase 6 — Cinematic Shell

Goal: build the public visual shell — but never before real chat exists on real pages.

### Task 6.1 — Obtain Cinematic Shell Integration Packet
1. **Objective:** get the actual Three.js Cosmic Intro Integration Packet from ChatGPT/Gemini, per the existing Product Operating System workflow.
2. **Reason for order:** repeatedly referenced across every doc in this repo as the next deliverable, but never produced or saved here — this is the actual blocker, not engineering capacity.
3. **Dependencies:** **Phase 2 must be complete first** — shipping a polished visual shell around a non-functional core chat would compound the misrepresentation already flagged in every audit.
4. **Technical requirements:** none — documentation/planning task.
5. **External services/accounts:** none.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** N/A.
8. **Definition of done:** packet saved to this repo's docs tree.
9. **Live test required:** N/A.
10. **What not to touch:** N/A.

### Task 6.2 — Implement Cinematic Shell per packet
1. **Objective:** build the actual public landing/intro experience.
2. **Reason for order:** last in this phase, gated on 6.1 and Phase 2.
3. **Dependencies:** Task 6.1, Phase 2 complete.
4. **Technical requirements:** per packet; `three` is already an installed dependency, zero current usage.
5. **External services/accounts:** UNKNOWN, depends entirely on the packet; possible 3D asset hosting/CDN decision.
6. **Legal/privacy requirements:** none apparent.
7. **Activation state after build:** `SPEC_ONLY` → `ACTIVE_PUBLIC` (this is inherently public-facing by design, not a candidate for dormant staging).
8. **Definition of done:** per packet's own acceptance criteria.
9. **Live test required:** cross-browser/device visual + performance QA.
10. **What not to touch:** per the packet's own explicit restrictions (historically: auth, routing, backend, Stripe, env vars, global CSS) plus, per this roadmap: do not touch the AI/chat path (`main.go`, `lib/buildPrompt.js`, `Chat.jsx`, `BravePAv2.jsx`) — confirmed already as a standing instruction in `docs/decisions/AI_REVIEW_LEDGER_TRUST_AND_CINEMATIC_SHELL_v0.2.md`.

---

## Phase 7 — Core Flow QA and Launch Readiness

Goal: verify everything that was built-but-unverified, and close the remaining narrow legal gap.

### Task 7.1 — Live verification sweep
1. **Objective:** confirm live-config for every "code real, live status unknown" item: Stripe (`STRIPE_SECRET_KEY`/`STRIPE_WEBHOOK_SECRET` + one real test charge), Supabase `media` bucket policy, and a final end-to-end confirmation of 2.1/2.2/2.3/4.1.
2. **Reason for order:** the highest financial and trust stakes in the whole roadmap (real money via Stripe) — must happen before any launch claim, not after.
3. **Dependencies:** Phase 2 and Task 4.1 complete.
4. **Technical requirements:** none — verification only.
5. **External services/accounts:** Stripe, Supabase (both already integrated).
6. **Legal/privacy requirements:** confirm payment terms/refund policy exist and are accurate.
7. **Activation state after build:** `UNKNOWN` → `ACTIVE_PUBLIC` (Stripe, Media) once confirmed.
8. **Definition of done:** one real Stripe test charge succeeds with the webhook acknowledged before timeout (flagging the known Render free-tier cold-start risk); one real media upload succeeds through both existing flows.
9. **Live test required:** yes — this task **is** the live test suite.
10. **What not to touch:** Stripe integration code itself unless the verification reveals a real defect, in which case treat as its own separately-scoped fix, not a silent change here.

### Task 7.2 — Health-Data DPIA Scoping
1. **Objective:** the narrowest, most achievable legal deliverable identified across all four audits — a real DPIA covering the existing health/injury/pregnancy disclosure form types specifically (not the much larger voice/video question).
2. **Reason for order:** these forms are already live and already collecting special-category data today; this is real, present exposure, not speculative.
3. **Dependencies:** none — a legal-track task, can run in parallel with engineering phases.
4. **Technical requirements:** none for the DPIA itself; may produce follow-up engineering tasks (e.g., retention/deletion controls) not yet scoped.
5. **External services/accounts:** a data-protection/GDPR specialist (external, not Claude Code).
6. **Legal/privacy requirements:** this task **is** the legal/privacy requirement.
7. **Activation state after build:** the legal-forms feature remains `ACTIVE_PUBLIC` throughout — this task closes a compliance gap, not a functional one.
8. **Definition of done:** a completed DPIA document exists and is saved to this repo's docs tree (e.g., `docs/legal/`).
9. **Live test required:** N/A — legal deliverable.
10. **What not to touch:** the legal-forms code itself, unless the DPIA specifically recommends a change.

### Task 7.3 — Booking-intent structured detection
1. **Objective:** replace the regex-based "a booking happened" detection with a structured tool-call/JSON contract from the model.
2. **Reason for order:** a known fragility (silent-miss risk if AI phrasing varies), worth closing during QA rather than discovering it live post-launch.
3. **Dependencies:** none technically; benefits from Task 3.2's server-owned prompt work being done first.
4. **Technical requirements:** a structured output contract for booking confirmations, replacing the current `replyLower.includes('sent your request')` pattern in `Chat.jsx`.
5. **External services/accounts:** none new.
6. **Legal/privacy requirements:** none.
7. **Activation state after build:** `ACTIVE_PUBLIC` (unchanged externally — this is a robustness fix).
8. **Definition of done:** booking detection no longer depends on exact AI phrasing.
9. **Live test required:** yes — vary the AI's confirmation phrasing and confirm detection still fires.
10. **What not to touch:** the booking-request storage/dashboard flow itself — already real, do not modify.

### Task 7.4 — Full Pre-Launch QA Day
1. **Objective:** execute the existing Product Operating System's Phase L3 (Pre-Launch QA Day) now that the above gaps are closed.
2. **Reason for order:** last in this phase by definition — it's the gate to the Launch Pipeline.
3. **Dependencies:** all prior tasks in Phases 2, 4, 7.
4. **Technical requirements:** per `docs/protocol/PRODUCT_OPERATING_SYSTEM_v0.5.md` §8 Phase L3 (not re-specified here — that document remains authoritative for QA-day mechanics).
5. **External services/accounts:** none new.
6. **Legal/privacy requirements:** confirm Task 7.2 is either complete or explicitly accepted as an open risk by Bruno.
7. **Activation state after build:** N/A — this is a gate, not a feature.
8. **Definition of done:** per the existing Phase L3 definition in the Product Operating System doc.
9. **Live test required:** yes, comprehensively.
10. **What not to touch:** N/A — this phase is itself the verification of everything else.

---

## Phase 8 — Manager Agent and Connectors

Goal: generalize the pattern proven by the Calendar connector (Phase 3) into a real, incremental connector architecture.

### Task 8.1 — Generalize connector architecture
1. **Objective:** extract a real `connectors`/`integrations` pattern from the concrete Calendar connector built in Phase 3, rather than over-designing upfront.
2. **Reason for order:** deliberately sequenced *after* a first real connector exists — building the abstraction before a concrete instance risks guessing wrong.
3. **Dependencies:** Task 3.3 (Calendar connector) complete and proven in production.
4. **Technical requirements:** a new `connectors`/`integrations` table (name, owner, status, tokens) — new schema, same "no established migration pattern" precedent-setting concern as Phase 1.
5. **External services/accounts:** none new for the architecture itself.
6. **Legal/privacy requirements:** privacy notice updates per connector added.
7. **Activation state after build:** `SPEC_ONLY` → `DORMANT_BUILT`.
8. **Definition of done:** the Calendar connector is refactored to use the new generalized pattern without behavior change.
9. **Live test required:** yes — confirm Calendar connector still works identically post-refactor.
10. **What not to touch:** Invyted integration — explicitly excluded per the Codex "until it publishes a public API," not a Claude Code decision to revisit.

### Task 8.2 — Additional connectors, incrementally
1. **Objective:** add the next connector (candidates: email, search) one at a time, per the Codex's own explicit sequencing instruction.
2. **Reason for order:** last, and deliberately paced — "incremental, one connector at a time" is the Codex's own stated approach, not a Claude Code simplification.
3. **Dependencies:** Task 8.1.
4. **Technical requirements:** per connector, to be scoped individually when reached — not pre-specified here to avoid guessing ahead of real need.
5. **External services/accounts:** per connector.
6. **Legal/privacy requirements:** per connector.
7. **Activation state after build:** `DORMANT_BUILT` per connector, owner-opt-in.
8. **Definition of done:** per connector.
9. **Live test required:** yes, per connector.
10. **What not to touch:** connectors not yet reached in the sequence.

---

## Phase 9 — Brain Spine / Memory / Personalization

Goal: only after Cinematic Shell and Core Flow QA, per the existing Brain Spine Readiness Audit's own recommendation.

### Task 9.1 — Minimal Brain Spine implementation
1. **Objective:** implement the Minimal Brain Spine exactly as scoped in `docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md` §7-11, if and only if Bruno explicitly approves after reviewing that audit.
2. **Reason for order:** the Brain Spine audit's own recommendation was "document only, then Cinematic Shell, then implement later" — this roadmap honors that sequencing rather than overriding it.
3. **Dependencies:** Phase 6 (Cinematic Shell) and Phase 7 (Core Flow QA) complete; Bruno's explicit approval.
4. **Technical requirements:** exactly as specified in the Brain Spine audit — request envelope, model adapter, memory/provenance interface, verification/audit interface, action/tool gate, response synthesis — all as pass-through wrappers around existing behavior, no new provider, no schema change for the minimal version.
5. **External services/accounts:** none — explicitly uses the existing Anthropic path only.
6. **Legal/privacy requirements:** none for the minimal (pass-through) version.
7. **Activation state after build:** internal architecture change — no user-facing activation state shifts.
8. **Definition of done:** golden-file tests confirm zero behavior change versus pre-Spine responses.
9. **Live test required:** yes, per the Brain Spine audit's own §9 recommendation.
10. **What not to touch:** exactly the exclusion list already specified in the Brain Spine audit (auth, routing, Supabase schema beyond what's explicitly approved, Stripe, new providers, MCP).

### Task 9.2 — Basic memory re-read (dormant, unmarketed)
1. **Objective:** read stored conversation history back into future prompts server-side — the smallest real step toward "personalization."
2. **Reason for order:** depends on 9.1's `MemoryInterface` existing as the clean integration point.
3. **Dependencies:** Task 9.1.
4. **Technical requirements:** no new schema — reuses existing `conversations`/`messages` tables.
5. **External services/accounts:** none for this basic version.
6. **Legal/privacy requirements:** proportional — basic re-read of a business's own operational data (leads, services) is lower-risk than behavioral profiling; still worth a light privacy-notice mention.
7. **Activation state after build:** `GHOST_FORBIDDEN` (claim) / `DORMANT_BUILT` (storage) → `DORMANT_BUILT` (real, but unmarketed) → `ACTIVE_PRIVATE` once confidence is built internally.
8. **Definition of done:** the AI demonstrably references prior context in a real conversation.
9. **Live test required:** yes.
10. **What not to touch:** do not market as "personalized" or "remembers you" until this has real production track record — see `docs/strategy/EXTERNAL_DEPENDENCIES_COSTS_AND_LEGAL_READINESS_v1.md` for the copy constraints.

### Task 9.3 — Advanced/semantic memory (legal-gated)
1. **Objective:** real embeddings-based semantic memory, if and when pursued.
2. **Reason for order:** last in this phase, explicitly V2/V3 per the Brain Spine audit's own sequencing.
3. **Dependencies:** Task 9.2 proven; a legal review of the profiling question.
4. **Technical requirements:** a vector store (e.g., pgvector on Supabase, or a dedicated vector DB — **provider not chosen, UNKNOWN**); an embeddings provider.
5. **External services/accounts:** UNKNOWN, pending provider selection; likely a new cost line — **COST_VERIFY_REQUIRED**.
6. **Legal/privacy requirements:** likely a DPIA — genuine behavioral/preference learning is closer to profiling under GDPR than simple transactional storage.
7. **Activation state after build:** `GHOST_FORBIDDEN` → gated behind DPIA outcome.
8. **Definition of done:** to be scoped once legal review completes.
9. **Live test required:** yes, once built.
10. **What not to touch:** N/A until this task is actually reached — do not pre-build.

---

## Phase 10 — Voice-first and Voice/Video Session Learning — legal-locked until cleared

Goal: resolve the scope ambiguity and legal gate before any engineering work, per the Codex's own binding instruction to not schedule this into any sprint.

### Task 10.1 — Legal scope resolution (Bruno + legal)
1. **Objective:** explicitly separate, or confirm as one and the same, (a) transient voice I/O with no learning/storage and (b) Voice & Video Session Learning — the repo currently conflates these by grouping "voice" with the DPIA-blocked category everywhere it appears.
2. **Reason for order:** must happen before any other Phase 10 task — engineering cannot safely make this determination.
3. **Dependencies:** none — a Bruno/legal-track decision, not an engineering task.
4. **Technical requirements:** none.
5. **External services/accounts:** a data-protection/GDPR specialist.
6. **Legal/privacy requirements:** this task **is** the legal/privacy requirement.
7. **Activation state after build:** N/A — this task produces a decision, not a feature.
8. **Definition of done:** a documented decision (saved to this repo) on whether voice I/O is in-scope for the existing legal block or requires its own, narrower review.
9. **Live test required:** N/A.
10. **What not to touch:** N/A.

### Task 10.2 — DPIA completion (voice/video, full scope)
1. **Objective:** complete the DPIA and GDPR/UK GDPR legal review explicitly required by `PRODUCT_OPERATING_SYSTEM_v0.5.md` line 723 before any Voice & Video Session Learning build.
2. **Reason for order:** explicit, binding precondition — "Do not schedule it into a sprint" until this exists.
3. **Dependencies:** Task 10.1.
4. **Technical requirements:** none — legal deliverable.
5. **External services/accounts:** data-protection/GDPR specialist.
6. **Legal/privacy requirements:** all 8 safeguards checked in the Revision Audit's classification (DPIA, explicit consent flow, privacy notice, retention policy, deletion/export controls, biometric/health-data safeguards, processor list, activation gating) — currently 0/8 present.
7. **Activation state after build:** N/A until complete.
8. **Definition of done:** a completed DPIA exists and is saved to this repo.
9. **Live test required:** N/A.
10. **What not to touch:** N/A.

### Task 10.3 — Voice I/O build (only if Task 10.1 clears it as separate/lower-risk)
1. **Objective:** build transient, non-learning voice input/output, if and only if Task 10.1 concludes this is genuinely out of scope for the full session-learning block.
2. **Reason for order:** conditional — this task may not proceed at all, depending on 10.1's outcome.
3. **Dependencies:** Task 10.1 explicitly clearing this scope.
4. **Technical requirements:** provider not chosen — 3 viable options: (a) browser-native Web Speech API (free, no account, inconsistent cross-browser quality); (b) a hosted STT/TTS provider (paid, higher quality); (c) a future Anthropic voice-capable endpoint if/when relevant (speculative, unconfirmed).
5. **External services/accounts:** depends on provider chosen.
6. **Legal/privacy requirements:** per Task 10.1's determination.
7. **Activation state after build:** `LEGAL_LOCKED` → `DORMANT_BUILT` (if cleared) → `ACTIVE_PRIVATE`.
8. **Definition of done:** to be scoped once 10.1 clears.
9. **Live test required:** yes, once built.
10. **What not to touch:** do not build any part of this until Task 10.1 is formally resolved — this is a **DO_NOT_BUILD_YET** item until then.

### Task 10.4 — Voice & Video Session Learning (last, full DPIA required)
1. **Objective:** the full vision-stated capability, if and when pursued.
2. **Reason for order:** explicitly last in the entire roadmap, per the Codex's own binding block.
3. **Dependencies:** Task 10.2 (completed DPIA) — non-negotiable.
4. **Technical requirements:** not scoped here — premature until legally cleared, per this task's own instruction not to speculate on technical constraints before the legal gate.
5. **External services/accounts:** a real speech/video processing provider — real ongoing cost, not yet selected.
6. **Legal/privacy requirements:** full DPIA, all 8 safeguards from Task 10.2.
7. **Activation state after build:** `LEGAL_LOCKED` until Task 10.2 completes, then re-evaluated.
8. **Definition of done:** N/A — not to be defined until legally cleared.
9. **Live test required:** N/A.
10. **What not to touch:** **do not build any part of this feature under any circumstance until Task 10.2 is complete.** This is the one absolute rule in this entire roadmap.

---

## Summary sequencing table

| Phase | Focus | Gated by |
|---|---|---|
| 0 | Reality baseline | — (complete) |
| 1 | Feature flags / audit events | — |
| 2 | Truth-critical public core | — |
| 3 | Assistant capability truth | Phase 1 recommended |
| 4 | Notification Center foundation | Phase 1 (partial) |
| 5 | Trust Dot / Brave Star / Location | Phase 1, spec acquisition |
| 6 | Cinematic Shell | Phase 2 complete |
| 7 | Core Flow QA / Launch Readiness | Phases 2, 4 |
| 8 | Manager Agent / connectors | Phase 3 (Calendar proven) |
| 9 | Brain Spine / Memory | Phases 6, 7, Bruno approval |
| 10 | Voice — legal-locked | Legal review, all phases independent of timing |
