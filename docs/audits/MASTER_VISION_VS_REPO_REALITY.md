# Master Vision vs Repo Reality Matrix

**AI:** Claude Code
**Mode:** Repo Execution Mode / Master Vision vs Repo Reality Audit Only
**Date:** 2026-07-09
**Status:** Audit only. No code changed.

Context: Bruno clarified that The Concierge is not meant to be a simple dashboard/chat product. The intended product is an interactive, voice-first, adaptive assistant layer that becomes the master interface of the Brave by Bruno ecosystem. This audit compares the current repo against that vision across 12 source-of-truth concepts.

**Sourcing note:** the Master Document and Concierge Codex Audited v1 — ranked #2 and #4 in this repo's own Source of Truth Hierarchy (`docs/protocol/PRODUCT_OPERATING_SYSTEM_v0.5.md` §2) — are not present anywhere in this repository. Only the Product Operating System v0.5 doc exists, which excerpts fragments of both (§7.6 "Codex-Locked Features Queue", §7.7 "Master Doc Blocco A3-A5"). For items covered by those excerpts, this audit uses real, quoted repo text. For items not covered by any repo document, the vision column is marked **UNDOCUMENTED IN REPO** and repo reality is audited against the task-prompt phrasing instead, flagged as such. No content is fabricated or attributed to the Master Document/Codex beyond what is actually quoted from this repo.

---

## 1. The Concierge as operating engine of the ecosystem
- **Intended vision:** UNDOCUMENTED IN REPO. No repo doc frames Concierge as an "engine" other ecosystem products plug into. Closest repo evidence: commit history shows the product was "relocated to `/theconcierge`, root becomes ecosystem container" (root `/` now says "Brave by Bruno — More to come → The Concierge").
- **Current repo implementation:** Concierge is one self-contained product living under `/theconcierge`; the root `/` is an empty shell with a single link.
- **Structure:** Frontend-only distinction (`app/page.tsx` vs `app/theconcierge/*`); no backend concept of "ecosystem" exists — one Go backend serves Concierge alone.
- **Function:** Concierge works as a standalone product; nothing calls into it as infrastructure from elsewhere, because nothing else exists yet.
- **Perception:** the root page's "More to come" is honest — it doesn't overclaim.
- **Status:** FUTURE (root is a real placeholder for a real future container; the "engine" framing itself has no implementation to evaluate).
- **What's needed:** requires the Master Document itself to define what "operating engine" means structurally before this can be scoped.
- **Blocks launch:** NO.
- **Build timing:** after Core Flow QA — this is an ecosystem-level architecture decision, not a pre-launch item.

## 2. Brave PA as omnipresent guide
- **Intended vision:** UNDOCUMENTED IN REPO (task-prompt phrasing only — "omnipresent guide").
- **Current repo implementation:** `BravePAv2.jsx` is mounted in exactly 2 places — `DashboardClient.jsx` (main owner dashboard) and `OwnerEditClient.jsx` (profile-edit page). Both are authenticated, owner-only screens.
- **Structure:** Frontend component + shared `/chat` backend endpoint (undifferentiated from Concierge chat).
- **Function:** works as a real chat assistant on those 2 screens; absent from onboarding, owner-auth/login, and the entire public/visitor-facing side.
- **Perception:** an owner would experience Brave PA as present on their main working screens, not as truly omnipresent across the whole product.
- **Status:** PARTIAL-REAL (real where it exists, but "omnipresent" overstates its actual footprint — it covers 2 of ~7 owner-facing routes).
- **What's needed:** either mount it globally (a persistent floating layer across all authenticated routes) or adjust the "omnipresent" framing to "available on your main dashboard screens."
- **Blocks launch:** NO.
- **Build timing:** during Cinematic Shell or after — expanding its footprint is a UI/shell-layer decision, not core-flow-blocking.

## 3. Voice-first assistant experience
- **Intended vision:** the repo explicitly and repeatedly groups "voice" with biometric/video/session-learning as requiring legal review before any build (`PRODUCT_OPERATING_SYSTEM_v0.5.md` lines 303, 699, 723). It does not separately define a lower-risk "voice-first interaction" (e.g., speech-to-text input/TTS output with no session learning/storage) as distinct from the DPIA-blocked "Voice & Video Session Learning." This is a real ambiguity worth resolving with Bruno, not something resolved unilaterally here.
- **Current repo implementation:** zero. No speech recognition, no text-to-speech, no microphone/audio-input code anywhere in the codebase (confirmed by repo-wide grep).
- **Structure:** none exists.
- **Function:** none.
- **Perception:** N/A — nothing user-facing claims this today.
- **Status:** BLOCKED-LEGAL, treated conservatively as covered by the existing voice/video legal block until Bruno explicitly separates "voice interaction" from "voice session learning" in the Master Document.
- **What's needed:** a legal/DPIA review (per the Codex's own binding block) covering whichever scope Bruno intends — full session learning, or lighter-weight voice I/O.
- **Blocks launch:** NO (not built, so nothing to block; but must not be marketed).
- **Build timing:** after Core Flow QA, and only after the DPIA/legal review the repo itself already mandates.

## 4. Brave Star — 3 behavioural states
- **Intended vision:** real, quoted repo text — `PRODUCT_OPERATING_SYSTEM_v0.5.md` §7.6: "Brave Star - 3 behavioural states (replaces the single decibel-line concept) - no external blocker, ready now." No further specification of what the 3 states are exists in this repo.
- **Current repo implementation:** none. Zero references to "Brave Star" or "decibel" anywhere in actual code (confirmed by grep).
- **Structure:** none exists. **Function:** none. **Perception:** N/A.
- **Status:** FUTURE (cleared to build, per the Codex's own words, but not started).
- **What's needed:** the Codex's fuller spec (not present in this repo) defining what the 3 behavioural states actually are and display; this repo only has the name and clearance status.
- **Blocks launch:** NO.
- **Build timing:** the Codex explicitly says "ready now" (no external blocker) — candidate for before or during Cinematic Shell if Bruno wants it as part of the visual identity, since it's cleared independent of the visual-shell work.

## 5. Generative Core / Build Your Universe onboarding
- **Intended vision:** UNDOCUMENTED IN REPO. No repo doc uses "Generative Core" or "Build Your Universe" anywhere.
- **Current repo implementation:** `components/Onboarding.jsx`, a structured multi-step form wizard (`STEPS` array in `lib/constants.js`) — "What describes you?", "The basics", "Your tone", etc. It is not framed as generative in any repo copy; the only AI-generative moment in the whole flow is the one-shot vision-based service-import from a screenshot (`POST /ai/import-services`).
- **Structure:** real frontend wizard + real backend save (`POST /profile`, `POST /auth/signup`).
- **Function:** produces a real, working profile — confirmed real in the Product Reality Matrix (`docs/audits/PRODUCT_REALITY_MATRIX.md` #10).
- **Perception:** a user experiences this as filling out a form, not as "building a universe."
- **Status:** the underlying mechanism is REAL; the "Generative Core / Build Your Universe" framing/branding is GHOST — it doesn't exist anywhere in the product, so there's a naming/positioning gap between whatever vision document uses that phrase and what a user actually sees.
- **What's needed:** either the Master Document's actual spec for what "Generative Core" adds beyond the current wizard (more AI-driven generation throughout, not just one screenshot-import step), or confirmation that "Build Your Universe" is meant as marketing language over the existing wizard rather than a functional change.
- **Blocks launch:** NO — the underlying onboarding works; only the aspirational framing is absent.
- **Build timing:** during Cinematic Shell if it's a presentation/copy layer; after Core Flow QA if it implies real new generative mechanics.

## 6. Cinematic Shell as public visual shell only
- **Intended vision:** repo docs describe Cinematic Shell as a "Three.js Cosmic Intro" integration packet for the public landing experience (`PRODUCT_OPERATING_SYSTEM_v0.5.md` §1007, §1290; `BRAIN_SPINE_FLEET_INTEGRATION_DECISION_PACK_v0.2.md` §9). None of the repo docs explicitly state the boundary "public visual shell only" (i.e., that it must not touch auth/routing/backend) — that specific framing is closest to Bruno's own restriction lists in recent tasks, not a quoted Master Document constraint.
- **Current repo implementation:** none. `three` is listed in `package.json` as a dependency but is not imported anywhere in `app/`, `components/`, or `lib/` (confirmed by repo-wide search) — the library is installed, nothing built with it yet.
- **Structure:** dependency present, zero implementation. **Function:** none. **Perception:** N/A.
- **Status:** FUTURE.
- **What's needed:** the actual Integration Packet (referenced repeatedly across docs as the next deliverable from ChatGPT/Gemini, not yet produced or saved to this repo).
- **Blocks launch:** NO (not started, so nothing broken yet — but per prior instructions, it must not touch the AI/chat/auth path when it lands).
- **Build timing:** this is the "during Cinematic Shell" phase by definition — it's the next major visual milestone per `docs/decisions/AI_REVIEW_LEDGER_TRUST_AND_CINEMATIC_SHELL_v0.2.md`.

## 7. Owner notifications and proactive alerts
- **Intended vision:** partially documented — Bruno's own product decisions across this session (captured in `docs/decisions/AI_REVIEW_LEDGER_TRUST_AND_CINEMATIC_SHELL_v0.2.md`) establish that notification claims must be backend-real; "proactive alerts" beyond that isn't separately specified in any repo doc.
- **Current repo implementation:** sensitive-topic alerts are now REAL end-to-end (record + email, per commits `412312e`/`d69daf0`, confirmed in `docs/audits/PRODUCT_REALITY_MATRIX.md` #8). Separately, `generateProactiveMessage()` in `lib/buildPrompt.js` produces canned, time-of-day-based greeting text ("Good morning! You have N hot leads...") — this is client-side only, computed in the browser from data already fetched, not a push/email/real-time proactive system.
- **Structure:** alert path — real frontend + real backend + real storage + real email. Proactive-message path — frontend-only, no backend trigger.
- **Function:** alert path works as audited; proactive-message path only "activates" if the owner happens to have the dashboard open at the right moment — it can't reach them otherwise.
- **Perception:** "proactive" implies the assistant reaches out to the owner; today it only speaks up if they're already looking at the screen.
- **Status:** alerts = REAL. Proactive messaging = PARTIAL-REAL (real logic, not actually proactive in delivery).
- **What's needed:** push notification infrastructure (already scoped as a blocker in the notification work — VAPID keys, subscription storage, service worker registration) to make "proactive" actually reach the owner outside the open tab.
- **Blocks launch:** NO for alerts (already real); proactive messaging should not be marketed as "reaches you anywhere" until push exists.
- **Build timing:** alerts are done; proactive/push is after Core Flow QA, per the blockers already documented in the notification work.

## 8. Assistant settings: voice, sound, notifications, personality
- **Intended vision:** UNDOCUMENTED IN REPO as a unified "Assistant Settings" concept; closest repo reality is scattered across two components.
- **Current repo implementation:** `components/BravePASettings.jsx` — real, persisted personality/name settings (`paName`, `personality` — used in `buildBravePAPrompt()`, genuinely affects the AI's behavior). `components/OwnerDashboard.jsx`'s `notifPrefs` — `newLead`, `hotLead`, `newBooking`, `newMessage`, `dailyBriefing`, `soundEnabled`, `soundStyle` toggles, but stored only in `localStorage`, never sent to the backend, and — critically — nothing in the codebase ever reads `soundEnabled`/`soundStyle` to actually play a sound (confirmed: zero `Audio`/`.play()` calls anywhere in the repo). These toggles currently do nothing beyond persisting their own state locally.
- **Structure:** personality settings = real frontend + real backend persistence (via profile save). Notification/sound prefs = frontend-only, no backend, no consuming logic.
- **Function:** personality settings genuinely change AI behavior. Sound/notification toggles are inert switches.
- **Perception:** an owner toggling "sound on" would reasonably expect to hear a sound on a new lead — this will never happen with the current code.
- **Status:** personality = REAL. Voice settings = FUTURE (nothing exists to configure). Sound settings = GHOST (a real-looking, interactive UI control wired to nothing). Notification channel prefs = GHOST for the same reason, and now additionally inconsistent with the actually real notification system built this session, which doesn't check these preferences at all.
- **What's needed:** either wire `notifPrefs` into the real notification path (check `soundEnabled` before playing a sound client-side on alert arrival; persist prefs server-side so they mean something across devices) or remove/relabel the toggles as "coming soon."
- **Blocks launch:** the sound/notification toggles being non-functional isn't launch-blocking by itself, but it's a direct perception-risk match to the ghost/placeholder pattern already flagged in the Product Reality Matrix — worth fixing alongside that work rather than separately.
- **Build timing:** during Core Flow QA (it's a small, contained fix, not a visual-shell or architecture item).

## 9. Manager Agent / connector direction
- **Intended vision:** real, quoted repo text — `PRODUCT_OPERATING_SYSTEM_v0.5.md` §7.6: "Manager Agent (incremental, one connector at a time) - ready now; Invyted stays excluded until it publishes a public API." Also §7.7 Blocco A3: Google Calendar OAuth (read-only) is named as "mid-implementation — CLIENT_ID/CLIENT_SECRET retrieval from Google Cloud Console was the last confirmed step."
- **Current repo implementation:** no "Manager Agent" code exists anywhere (confirmed by grep). However, `db.Profile.GoogleRefreshToken` does exist as a schema field with zero code reading or writing it, and Brave PA's system prompt (`buildBravePAPrompt()`) tells the model it can "Read and add events to Google Calendar" — a capability with no backend behind it, already flagged as a false-capability-claim risk in the Product Reality Matrix (#7, Brave PA).
- **Structure:** one orphaned schema field is the only trace of Blocco A3's prior progress; no Manager Agent architecture exists.
- **Function:** none — the field is unused, and the AI's claimed calendar capability is fiction.
- **Perception:** an owner asking Brave PA "add this to my calendar" gets a plausible-sounding but false response.
- **Status:** FUTURE for Manager Agent itself; GHOST for the specific "I can manage your calendar" claim already live in the Brave PA system prompt today.
- **What's needed:** resume Blocco A3 (Google Calendar OAuth) as the first Manager Agent connector, per the Master Doc's own stated sequencing — this repo confirms it was genuinely in-progress before being interrupted by the auth-hardening work, not abandoned by decision.
- **Blocks launch:** the missing calendar feature doesn't block launch; the false claim about it already flagged in the Product Reality Matrix does, if Brave PA is marketed pre-launch.
- **Build timing:** Blocco A3 resumption is explicitly scoped in the Master Doc excerpt as happening "before or alongside Phase B1" (Public Experience Restoration) — i.e., before Cinematic Shell, per the repo's own stated sequencing, not after.

## 10. Trust Dot
- **Intended vision:** real, quoted repo text — `PRODUCT_OPERATING_SYSTEM_v0.5.md` §7.6: "Trust Dot - no external blocker, ready now." No further specification of what it visually is or does exists in this repo (also referenced only as a future conceptual seam in `docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md` §8F — "trust/verification metadata could be attached later").
- **Current repo implementation:** none. Zero code references anywhere (confirmed by grep).
- **Structure:** none. **Function:** none. **Perception:** N/A.
- **Status:** FUTURE (cleared, not started).
- **What's needed:** the Codex's fuller spec (not present in this repo) for what the Trust Dot actually indicates and where it appears.
- **Blocks launch:** NO.
- **Build timing:** cleared "ready now" per the Codex — candidate for before/during Cinematic Shell if it's meant to be part of the chat UI's visual trust signaling (which would tie it naturally to the review/notification truthfulness work already done this session).

## 11. Location-aware background
- **Intended vision:** real, quoted repo text — `PRODUCT_OPERATING_SYSTEM_v0.5.md` §7.6: "Sfondo Location-Aware (skyline + weather + time of day) - no external blocker, ready now."
- **Current repo implementation:** none as a visual feature. The only location-related code is a one-time reverse-geocoding call in `Onboarding.jsx` (`nominatim.openstreetmap.org/reverse`) used to populate the profile's location text field during signup — unrelated to a dynamic visual background. No weather API integration exists anywhere; "weather" only appears as text the AI is instructed to talk about in Brave PA's system prompt, never a real data call.
- **Structure:** none for the background feature itself. **Function:** none. **Perception:** N/A — nothing claims this today.
- **Status:** FUTURE (cleared, not started).
- **What's needed:** a weather API integration (new external service/key — out of scope for this audit to select) plus a skyline/time-of-day rendering layer, most naturally built alongside Cinematic Shell since both are visual/presentation-layer work.
- **Blocks launch:** NO.
- **Build timing:** during Cinematic Shell — it's visual-shell work by nature and is explicitly cleared with no dependency blocking it.

## 12. Voice & Video Session Learning — BLOCKED-LEGAL
- **Intended vision:** the most explicitly and unambiguously documented item in the entire matrix. `PRODUCT_OPERATING_SYSTEM_v0.5.md` line 723: "Voice & Video Session Learning stays in spec only - DPIA and GDPR/UK GDPR legal review required before any build, per the Codex's binding block. Do not schedule it into a sprint." Reinforced at lines 303 and 699.
- **Current repo implementation:** none whatsoever — confirmed by repo-wide search for any voice/video capture or session-learning code. This is a case where repo reality and the documented legal block are in full agreement.
- **Structure:** none. **Function:** none. **Perception:** N/A — correctly, nothing markets this.
- **Status:** BLOCKED-LEGAL, exactly as documented — this is the one item where vision and reality already match perfectly.
- **What's needed:** a completed DPIA and GDPR/UK GDPR legal review before any implementation work — explicitly not a Claude Code task.
- **Blocks launch:** N/A (correctly not built, correctly not planned).
- **Build timing:** not to be scheduled into any sprint, per the Codex's own binding instruction, until the legal review exists — this applies regardless of Cinematic Shell or Core Flow QA timing.

---

# Top 10 gaps between vision and repo reality

1. **Voice-first / voice interaction** — entirely unbuilt, and the repo doesn't even clearly distinguish "voice I/O" from the DPIA-blocked "Voice & Video Session Learning," risking either over-blocking a legitimate feature or under-blocking a sensitive one.
2. **Brave PA's claimed calendar/search capabilities** — the AI is told it can do things (Google Calendar, web search) that have zero backend implementation; this is a live, active misrepresentation today, not just a missing future feature.
3. **Sound/notification preference toggles** — a fully interactive-looking settings UI wired to nothing; the same ghost pattern already flagged in the Product Reality Matrix, now confirmed to extend into "Assistant Settings."
4. **Generative Core / Build Your Universe** — the onboarding that exists is a conventional form wizard; whatever generative/universe-building experience the vision implies doesn't exist in any form, not even partially.
5. **Cinematic Shell** — `three` is installed as a dependency and literally nothing else exists; zero implementation despite being the most-referenced upcoming milestone across every doc in this repo.
6. **Trust Dot** — named and cleared in the Codex, zero implementation, and no spec detailed enough in this repo to build from.
7. **Brave Star (3 behavioural states)** — same pattern as Trust Dot: named, cleared, zero spec, zero code.
8. **Location-aware background** — cleared, zero implementation, and would require a net-new external weather service not yet decided.
9. **Manager Agent / Google Calendar (Blocco A3)** — the Master Doc excerpt says this was "mid-implementation" before being interrupted; today only an orphaned schema field remains, and the interruption was never resumed.
10. **"Omnipresent" Brave PA** — real on 2 of ~7 owner-facing screens, absent from the entire public/visitor side and from onboarding/login — a meaningful gap between "omnipresent" and its actual footprint.

# Features safe to market now
- Brave PA's core conversational ability (distinct from its false capability claims) — already listed as safe in the Product Reality Matrix.
- Onboarding's actual mechanics (distinct from any "Generative Core" branding) — already listed as safe in the Product Reality Matrix.
- Everything already cleared in `docs/audits/PRODUCT_REALITY_MATRIX.md` "Features safe to market now" still holds; this audit doesn't change that list, only adds vision-alignment context.
- Brave Star / Trust Dot / Sfondo — not applicable, none exist yet, so there's nothing to market either way.

# Features that must not be marketed yet
- Any "voice" claim, in any framing, until Bruno and legal resolve the voice-I/O-vs-session-learning ambiguity.
- Brave PA's calendar/search/booking capabilities as stated in its own system prompt.
- Sound/notification preference controls as functional.
- "Omnipresent" as a description of Brave PA's availability.
- "Generative"/"Build Your Universe" framing of onboarding, unless and until it's actually generative.
- Trust Dot, Brave Star, Location-aware background — none exist to market.

# Features that should be visually represented but labelled "coming soon"
- **Trust Dot** — cleared, no blocker, natural to preview as "coming soon" alongside the notification-truthfulness work already shipped.
- **Brave Star (3 behavioural states)** — cleared, no blocker, a visible placeholder could set expectations without overclaiming.
- **Location-aware background** — cleared, visual-only, low-risk to tease.
- **Manager Agent / Calendar connector** — reasonable to show as "Connect your calendar (coming soon)" rather than pretending Brave PA already does it.
- **Voice** — explicitly should NOT be teased as "coming soon" without the DPIA/legal review existing first, per the Codex's own binding language ("do not schedule it into a sprint") — even a "coming soon" label could be read as a premature commitment.

# Recommended next 3 implementation packets
1. **"Capability Truth Pass on Brave PA"** — same category as items already fixed this session: either strip the calendar/search/booking claims from `buildBravePAPrompt()` or resume Blocco A3 (Google Calendar OAuth) to back them, since the Master Doc excerpt confirms this was already mid-build and only interrupted, not cancelled. This directly extends the "Capability Claims Audit for Brave PA" packet already recommended in the Product Reality Matrix.
2. **"Assistant Settings Wiring"** — connect `notifPrefs` (sound/notification toggles) to real behavior: play a real sound client-side when a real alert arrives (the alert system is already real as of this session), and persist prefs server-side instead of `localStorage`-only. Small, contained, closes a concrete ghost-UI gap.
3. **"Codex Visual Trio Spec Request"** — before any code, this needs the actual Codex spec (not present in this repo) for Trust Dot, Brave Star's 3 states, and Sfondo Location-Aware pulled from the Master Document/Codex and saved into this repo the same way the other vision docs were — otherwise any implementation would be guessing at undocumented requirements.
