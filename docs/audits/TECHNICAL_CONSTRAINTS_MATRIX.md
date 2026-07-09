# Technical Constraints Matrix

**AI:** Claude Code
**Mode:** Repo Execution Mode / Technical Constraints Revision Audit Only
**Date:** 2026-07-09
**Status:** Audit only. No code changed.

This extends `docs/audits/MASTER_VISION_VS_REPO_REALITY.md` and `docs/audits/MASTER_VISION_VS_REPO_REALITY_REVISION.md` with a technical-constraints pass covering infrastructure, external dependencies, APIs, subscriptions, env vars, cost exposure, and build feasibility for 26 concepts (A-Z).

## Cross-cutting infra note

The repo's own doc confirms the Go backend runs on **Render's free tier** today — `PRODUCT_OPERATING_SYSTEM_v0.5.md` §7.7: "The Render upgrade specifically must happen before any real client is sent a link, per existing Master Doc pricing rules." Free-tier Render services spin down on inactivity and incur cold-start latency (typically 30-60s). This affects the cost/risk rating of every feature below that depends on fast/reliable backend response — especially Stripe webhooks (Stripe expects fast acks), notifications, and anything cron-adjacent. Flagged once here rather than repeated 26 times, but it materially affects G, H, K, R, and U.

---

## A. Real chat on real public profile pages
1. Status: **GHOST**
2. Frontend: mount real `<Chat/>` (or equivalent) on a route reachable for non-demo slugs; remove/replace the static bubble in `generateBusinessPage.js`
3. Backend: none new — reuses existing `POST /chat`
4. DB/schema: none
5. External APIs: none new — reuses existing Anthropic path
6. Provider accounts: none new
7. Env vars/secrets: none new (existing `ANTHROPIC_API_KEY` on Render)
8. OAuth: no
9. Webhooks: no
10. Cron/worker: no
11. PWA/service-worker: no
12. Storage/bucket: no
13. Legal/privacy work: no additional beyond existing chat consent
14. DPIA before activation: no
15. Dormant/private build first: not applicable — this is a pure bug-fix/completion of an already-public feature, no reason to gate
16. Runs on current stack: **YES**, fully — Vercel + existing backend, no new infra
17. New paid subscription/provider: **NO**
18. Cost/risk level: **FREE/LOW**
19. Complexity: **SMALL-MEDIUM** (routing + component wiring, not new systems)
20. Main technical blocker: none — it's a completion task, not a new build
21. Main operational blocker: none
22. Main legal blocker: none
23. Minimum safe V1: mount `<Chat/>` behind the real slug route with the same consent gate already proven on `/demo/*`
24. Dormant/private build strategy: N/A
25. Activation condition: ship directly — no gate needed
26. Live test required before marketing: **YES** — confirm a real (non-demo) slug produces a real AI reply end-to-end in production
27. Include in next packet: **YES — highest priority**

## B. Stale domain/share-link correction
1. Status: **PARTIAL-REAL** (functions, but violates the repo's own Anti-Chaos Rule #17)
2. Frontend: replace 4 hardcoded `concierge-ai-gamma.vercel.app` occurrences (`generateBusinessPage.js`, `BusinessPagePreview.jsx` fallback, `OwnerDashboard.jsx` x3, `Onboarding.jsx`) with the real domain or `window.location.origin`
3. Backend: none
4. DB/schema: none
5. External APIs: none
6. Provider accounts: none
7. Env vars/secrets: optionally introduce a `NEXT_PUBLIC_APP_URL` env var to avoid hardcoding again
8. OAuth: no. 9. Webhooks: no. 10. Cron: no. 11. PWA: no. 12. Storage: no.
13. Legal/privacy work: no
14. DPIA: no
15. Dormant/private build first: N/A — trivial fix, ship directly
16. Runs on current stack: **YES**
17. New paid subscription: **NO**
18. Cost/risk: **FREE/LOW**
19. Complexity: **SMALL**
20. Main technical blocker: none
21. Main operational blocker: none
22. Main legal blocker: none — but note it's a live violation of the repo's own binding Anti-Chaos Rule #17, so treat as urgent despite low technical difficulty
23. Minimum safe V1: string replacement, no new logic
24. Dormant/private build strategy: N/A
25. Activation condition: ship directly
26. Live test required: **YES** — copy the share link from a real dashboard session and confirm it resolves to `bravebybruno.com`
27. Include in next packet: **YES**, bundle with A (same root-cause fix)

## C. Consent Truth Path
1. Status: **GHOST**
2. Frontend: persist `giveConsent()` server-side instead of `sessionStorage`-only (small change to `Chat.jsx`); serve `privacy.html` at a real reachable route
3. Backend: reuse existing `POST /consent`, or add a lightweight equivalent scoped to general AI-processing consent (distinct from legal-form consent)
4. DB/schema: none required — the existing `consents` table already has the needed columns (`profile_id`, `session_id`, timestamps); this can reuse it with a `forms_agreed` value like `"ai_processing"` rather than a real form name
5. External APIs: none
6. Provider accounts: none
7. Env vars/secrets: none
8. OAuth: no. 9. Webhooks: no. 10. Cron: no. 11. PWA: no.
12. Storage/bucket: for the privacy-page fix specifically — either add a `public/` folder (simplest) or convert `privacy.html` to a real Next.js page under `app/privacy/page.jsx` (more idiomatic, recommended since it also gets it into the same routing/build pipeline as everything else)
13. Legal/privacy work: recommend a real legal review of `privacy.html`'s content once it's reachable again — separate from the technical fix
14. DPIA before activation: no (the technical fix doesn't need one; the broader consent-content review is good practice, not a hard DPIA trigger by itself)
15. Dormant/private build first: N/A — this is fixing an active production defect, ship directly
16. Runs on current stack: **YES**
17. New paid subscription: **NO**
18. Cost/risk: **FREE/LOW** technically; reputational/legal risk is currently HIGH while unfixed (a live 404 on a GDPR disclosure)
19. Complexity: **SMALL** (routing fix) + **SMALL** (consent persistence, same pattern as legal forms)
20. Main technical blocker: none
21. Main operational blocker: none
22. Main legal blocker: none for the fix itself; the content of the policy once reachable should still get a real review
23. Minimum safe V1: privacy page reachable + consent persisted, exactly mirroring the already-proven legal-form-consent pattern
24. Dormant/private build strategy: N/A
25. Activation condition: ship directly — this is not a feature to gate, it's a defect to close
26. Live test required: **YES** — `curl -L https://www.bravebybruno.com/privacy` must return 200, and a real consent click must produce a real DB row
27. Include in next packet: **YES — most urgent single item across both audits**

## D. Server-owned prompt templates
1. Status: **REAL** (chat works) but with an unmanaged trust gap — the server currently trusts a client-supplied `system_prompt` string on every `/chat` call with no rebuild/validation (flagged in `docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md` §2)
2. Frontend: `Chat.jsx`/`BravePAv2.jsx` would stop sending a full prompt string and instead send a `profile_id` + `mode`
3. Backend: `handleChat` would look up/rebuild the system prompt server-side from stored profile data instead of trusting the client's string
4. DB/schema: none required — `profiles.profile_data` already has what's needed to rebuild the prompt server-side
5. External APIs: none
6. Provider accounts: none
7. Env vars/secrets: none
8. OAuth: no. 9. Webhooks: no. 10. Cron: no. 11. PWA: no. 12. Storage: no.
13. Legal/privacy work: no
14. DPIA: no
15. Dormant/private build first: **YES, recommended** — this changes a security-relevant code path on the busiest endpoint in the product; worth a careful, isolated rollout even though it's conceptually simple
16. Runs on current stack: **YES**
17. New paid subscription: **NO**
18. Cost/risk: **LOW** technically, but touches the most business-critical endpoint — treat with the same care as the auth-hardening work
19. Complexity: **MEDIUM** (requires moving `lib/buildPrompt.js`'s logic — or an equivalent — server-side in Go, a real port, not a trivial change)
20. Main technical blocker: porting JS prompt-building logic to Go (or calling back into a shared service) — a real but bounded engineering task
21. Main operational blocker: must not break the currently-working chat behavior; needs the golden-file/snapshot testing approach already recommended in the Brain Spine audit
22. Main legal blocker: none
23. Minimum safe V1: server rebuilds the prompt from `profile_id` alone for the Concierge-chat case first (highest-value, least complex); leave Brave PA's richer context-dependent prompt for a follow-up
24. Dormant/private build strategy: ship behind a flag so it can be A/B-verified against the old client-supplied path before fully cutting over
25. Activation condition: golden-file test showing identical replies for identical inputs before/after
26. Live test required: **YES**
27. Include in next packet: recommended as its own packet, not bundled with A/B/C — it's a different risk category (security hardening of a live endpoint, not a UI/routing fix)

## E. Brave PA capability truth / Google Calendar connector
1. Status: **PARTIAL-REAL** (conversation) / **GHOST** (claimed calendar/search capability)
2. Frontend: `BravePASettings.jsx`/`BravePAv2.jsx` would need a "Connect Google Calendar" UI once the backend exists
3. Backend: new OAuth callback handler + calendar-read endpoint in `main.go`
4. DB/schema: none new — `Profile.GoogleRefreshToken` already exists as a column, unused; this is genuinely ready to receive a real value
5. External APIs: Google Calendar API
6. Provider accounts: Google Cloud Console project (per the Master Doc excerpt, "CLIENT_ID/CLIENT_SECRET retrieval from Google Cloud Console was the last confirmed step" — i.e., this may already partially exist from the interrupted Blocco A3 work; needs confirming with Bruno, not assumed)
7. Env vars/secrets needed (names only): `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_OAUTH_REDIRECT_URI`
8. OAuth: **YES** — this is the core of the feature
9. Webhooks: optional (Google Calendar push notifications channel) — not required for a read-only V1
10. Cron/worker: only if building proactive calendar-based reminders later; not required for basic read
11. PWA: no
12. Storage/bucket: no
13. Legal/privacy work: update `privacy.html` (once reachable) to disclose Google as a processor and the calendar scope requested
14. DPIA before activation: no — standard OAuth scope consent is sufficient, this isn't special-category data
15. Dormant/private build first: **YES** — natural to ship owner-opt-in/beta before claiming broadly
16. Runs on current stack: **YES** — no new infra category, just a new external API integration from the existing Go backend
17. New paid subscription: **NO** — Google Calendar API has a generous free tier for this use case
18. Cost/risk: **LOW**
19. Complexity: **MEDIUM** (OAuth flow + token refresh handling + calendar-read parsing)
20. Main technical blocker: none structural — genuinely the closest-to-done item in this whole audit, per the Master Doc's own "last confirmed step" note
21. Main operational blocker: confirming whether the Google Cloud Console project/credentials from the interrupted work still exist or need recreating
22. Main legal blocker: none beyond the standard OAuth consent screen
23. Minimum safe V1: read-only calendar view surfaced to Brave PA ("what's on my calendar today")
24. Dormant/private build strategy: ship as an opt-in "Connect Calendar" button, owner-only, not defaulted on
25. Activation condition: real OAuth token successfully refreshes and a real event list is returned
26. Live test required: **YES**
27. Include in next packet: **YES** — recommended alongside the Brave PA prompt truth-fix, since together they close the gap between claim and capability

## F. Manager Agent connector architecture
1. Status: **SPEC_ONLY**
2. Frontend: a connector-management UI (list of available/connected services) — doesn't exist
3. Backend: a generalized connector registry/interface — doesn't exist; E above would be its first concrete instance
4. DB/schema: a `connectors` or `integrations` table (new) to track which owner has which service connected, with tokens/status — new schema required
5. External APIs: depends on which connectors are prioritized (Calendar first, per Codex sequencing)
6. Provider accounts: one per connector, added incrementally
7. Env vars/secrets: one set per connector (see E for the first)
8. OAuth: yes, per-connector
9. Webhooks: optional, per-connector
10. Cron/worker: eventually useful for sync/refresh jobs, not required for V1
11. PWA: no
12. Storage: no
13. Legal/privacy work: privacy notice updates per connector added
14. DPIA: no, unless a future connector touches special-category data
15. Dormant/private build first: **YES** — this is explicitly how the Codex itself wants it built ("incremental, one connector at a time")
16. Runs on current stack: **YES** for the architecture itself; each connector may need its own provider account
17. New paid subscription: depends on connector — none required just to build the registry itself
18. Cost/risk: **LOW-MEDIUM** depending on how many connectors are added and their individual pricing
19. Complexity: **LARGE** for the generalized architecture; **MEDIUM** per individual connector once the pattern exists
20. Main technical blocker: needs a real schema design for `connectors`/`integrations` before the first connector (E) can be built cleanly — SPEC_MISSING for what exact fields/permissions model the Codex intends
21. Main operational blocker: Invyted is explicitly named as excluded "until it publishes a public API" — a dependency outside the team's control
22. Main legal blocker: none structurally; per-connector as they're added
23. Minimum safe V1: build the Calendar connector (E) first as a concrete instance, then extract the generalized pattern from it rather than over-designing upfront
24. Dormant/private build strategy: exactly as Codex specifies — one connector at a time, owner-opt-in
25. Activation condition: first real connector (Calendar) proven reliable in production
26. Live test required: **YES**, per connector
27. Include in next packet: not as its own packet yet — fold into E; revisit as its own packet once 2+ connectors exist and a real pattern needs extracting

## G. Owner notifications
1. Status: **ACTIVE_PRIVATE** (record + email real; live-config for email unconfirmed)
2. Frontend: already built (`Chat.jsx` alert flow, `OwnerDashboard.jsx` Notifications tab)
3. Backend: already built (`POST /alert`, `GET /owner/notifications`)
4. DB/schema: none — reuses `notes` table, no migration needed (already verified: no migration pattern exists in this repo)
5. External APIs: Resend (already integrated)
6. Provider accounts: Resend (already exists, per `sendHotLeadEmail` precedent)
7. Env vars/secrets (names only, already referenced in code): `RESEND_API_KEY`, `OWNER_EMAIL`, `RESEND_FROM_EMAIL`
8. OAuth: no. 9. Webhooks: no (Resend is a simple POST API, not webhook-driven for sending). 10. Cron: no, for this alert path specifically.
11. PWA/service-worker: not required for the email path (only for push, see H)
12. Storage: no
13. Legal/privacy work: no
14. DPIA: no
15. Dormant/private: already effectively private (owner-only visibility) — no further gating needed
16. Runs on current stack: **YES**, fully
17. New paid subscription: **NO** — already using existing Resend account
18. Cost/risk: **LOW** — but genuinely UNKNOWN whether Render free-tier cold-starts could delay the synchronous email-send-before-responding pattern built this week; worth a live timing check
19. Complexity: already built — **SMALL** remaining work (live-config verification only)
20. Main technical blocker: none remaining
21. Main operational blocker: confirm `RESEND_API_KEY`/`OWNER_EMAIL` are actually set on Render — this is the one open item
22. Main legal blocker: none
23. Minimum safe V1: already shipped
24. Dormant/private build strategy: N/A, already live
25. Activation condition: already active
26. Live test required: **YES** — trigger a real sensitive-topic message in production and confirm `email_status: "sent"` in the response, not just `"disabled_missing_env"`
27. Include in next packet: only as a verification task, not a build task

## H. Push notifications
1. Status: **GHOST** (client listener exists, zero delivery pipeline)
2. Frontend: register the existing `service_worker.js` (currently orphaned — no `public/` folder, no `serviceWorker.register()` call anywhere); build a permission-request + subscribe flow
3. Backend: new `POST /push/subscribe` endpoint; call a web-push send function from `handleCreateAlert`
4. DB/schema: new — a subscription-storage mechanism is needed; could reuse the `notes` pattern (`note_type: "push_subscription"`) to avoid a migration, or a dedicated table if Bruno wants push subscriptions modeled more explicitly — SPEC_MISSING on which Bruno prefers
5. External APIs: Web Push protocol (no third-party API beyond the browser's own push service — Google/Mozilla/Apple push services, which are free and don't require an account, only VAPID keys)
6. Provider accounts: none required — Web Push doesn't need a paid provider, just a self-generated VAPID key pair
7. Env vars/secrets (names only): `VAPID_PUBLIC_KEY`, `VAPID_PRIVATE_KEY`, `VAPID_SUBJECT` — not created, must be generated fresh, never hardcoded
8. OAuth: no
9. Webhooks: no
10. Cron/worker: no, for basic push (only needed for digest, see K)
11. PWA/service-worker: **YES — this is the core requirement**, and currently the biggest structural gap (no `public/` folder exists at all in this Next.js app)
12. Storage/bucket: no
13. Legal/privacy work: update privacy notice to disclose push notification data (subscription endpoint, which is technically an identifier)
14. DPIA before activation: no — push subscriptions are standard, low-risk technical data, not special-category
15. Dormant/private build first: **YES** — build the pipeline, keep it owner-opt-in and unpublicized until proven reliable
16. Runs on current stack: **YES**, once a `public/` folder is added — no new infrastructure category, just missing pieces
17. New paid subscription: **NO**
18. Cost/risk: **LOW** cost; **MEDIUM** complexity risk (getting service-worker registration right in a Next.js App Router project the first time)
19. Complexity: **MEDIUM**
20. Main technical blocker: no `public/` folder exists in this Next.js project at all — needs creating (also affects `privacy.html`'s eventual home, could be solved together with C)
21. Main operational blocker: VAPID keys don't exist yet and must be generated and stored in Render env
22. Main legal blocker: none
23. Minimum safe V1: subscribe + a single test push ("you have a new sensitive-topic alert") reusing the alert path already built
24. Dormant/private build strategy: owner-only opt-in toggle, not defaulted on, not marketed until reliability is proven over real usage
25. Activation condition: a real push notification successfully received on a real device
26. Live test required: **YES**, on a real browser/device — this cannot be meaningfully verified by build/vet alone
27. Include in next packet: recommended as its own packet after A/B/C/D are done — real infrastructure work, not a quick fix

## I. Notification preferences
1. Status: **GHOST** (`notifPrefs` exists, stored in `localStorage`, consumed by nothing)
2. Frontend: `OwnerDashboard.jsx` settings tab already has the UI; needs to actually gate behavior and persist server-side
3. Backend: extend the profile/settings save path to store preferences, or a small new `POST/GET /owner/notification-prefs` pair
4. DB/schema: could reuse `profiles.profile_data` JSON blob (no migration) or add real columns — SPEC_MISSING on which Bruno prefers; reusing the blob is lower-risk given no migration pattern exists
5. External APIs: none
6. Provider accounts: none
7. Env vars: none
8. OAuth: no. 9. Webhooks: no. 10. Cron: no. 11. PWA: no. 12. Storage: no.
13. Legal/privacy work: no
14. DPIA: no
15. Dormant/private build first: N/A — small fix, ship directly
16. Runs on current stack: **YES**
17. New paid subscription: **NO**
18. Cost/risk: **FREE/LOW**
19. Complexity: **SMALL**
20. Main technical blocker: none
21. Main operational blocker: none
22. Main legal blocker: none
23. Minimum safe V1: persist prefs server-side; gate the actual sound-playing (J) and future push (H) on the relevant toggle values
24. Dormant/private build strategy: N/A
25. Activation condition: ship directly once J exists to actually connect to (no point wiring prefs with nothing real to gate yet, beyond H/J)
26. Live test required: **YES** — toggle off, confirm no sound/notification; toggle on, confirm one does happen
27. Include in next packet: bundle with J, low priority relative to A/B/C

## J. Sound settings
1. Status: **GHOST** (`soundEnabled`/`soundStyle` toggles exist, zero code ever plays a sound — confirmed via repo-wide search for `Audio`/`.play()`)
2. Frontend: add a real `new Audio(...).play()` call triggered on real alert arrival, gated by `notifPrefs.soundEnabled`
3. Backend: none — this can be entirely client-side once an alert is confirmed via the already-real `/alert` response
4. DB/schema: none
5. External APIs: none
6. Provider accounts: none
7. Env vars: none
8. OAuth: no. 9. Webhooks: no. 10. Cron: no. 11. PWA: no.
12. Storage/bucket: a small audio asset file needs to be added to the (currently nonexistent) `public/` folder
13. Legal/privacy work: no
14. DPIA: no
15. Dormant/private build first: N/A
16. Runs on current stack: **YES**
17. New paid subscription: **NO**
18. Cost/risk: **FREE/LOW**
19. Complexity: **SMALL**
20. Main technical blocker: needs the same `public/` folder creation flagged under H
21. Main operational blocker: none
22. Main legal blocker: none
23. Minimum safe V1: one sound, plays on real alert confirmation, respects the mute toggle
24. Dormant/private build strategy: N/A
25. Activation condition: ship directly
26. Live test required: **YES**
27. Include in next packet: bundle with I, low priority relative to A/B/C

## K. Digest/clustering
1. Status: **SPEC_ONLY** (no scheduler of any kind exists anywhere in the codebase — confirmed via `go.mod`, no cron/ticker/queue library present)
2. Frontend: a "digest frequency" preference selector (partially exists already — `OwnerDashboard.jsx` has `digestFreq` state and calls `PUT /owner/digest-prefs`, but that route doesn't exist in `main.go` — a pre-existing broken call, confirmed by grep, worth noting even though not asked to fix here)
3. Backend: a new scheduled job/endpoint that batches pending notifications and sends one digest email
4. DB/schema: none required if reusing `profiles.DigestFrequency`/`DigestLastSent` (already exist as unused columns) — genuinely ready to receive real values, similar to the Calendar `GoogleRefreshToken` situation
5. External APIs: none new — reuses Resend
6. Provider accounts: none new
7. Env vars: none new
8. OAuth: no. 9. Webhooks: no.
10. Cron/worker: **YES — this is the core missing piece.** Options, none chosen yet (UNKNOWN, needs provider selection): (a) Render Cron Job (native to the existing hosting provider, simplest), (b) GitHub Actions scheduled workflow hitting a new endpoint, (c) a third-party scheduler like EasyCron or cron-job.org hitting a new endpoint. Given the backend is already on Render, option (a) is the natural fit.
11. PWA: no. 12. Storage: no.
13. Legal/privacy work: no
14. DPIA: no
15. Dormant/private build first: **YES** — build and test against a small owner group before claiming "digest" broadly
16. Runs on current stack: **NO** — requires adding a scheduler, which the current stack (Vercel serverless + Render web dyno) has no native mechanism for beyond Render's own Cron Job product (a paid feature on some Render plans — needs confirming against current plan tier)
17. New paid subscription: possibly — depends on which Render plan tier is active and whether Cron Jobs are included; UNKNOWN, needs confirming
18. Cost/risk: **LOW-MEDIUM**, pending the plan-tier question above
19. Complexity: **MEDIUM**
20. Main technical blocker: no scheduler infrastructure exists; also should first fix the pre-existing broken `/owner/digest-prefs` call while in this area
21. Main operational blocker: confirming Render plan/Cron Job availability
22. Main legal blocker: none
23. Minimum safe V1: a single daily digest email of unread alerts, manually triggerable first, scheduled second
24. Dormant/private build strategy: manual-trigger-only V1, scheduled version gated behind confirmed reliability
25. Activation condition: a real scheduled job fires and a real digest email is confirmed received
26. Live test required: **YES**
27. Include in next packet: not urgent — bundle with H as a "notification maturity" packet after A/B/C ship

## L. Voice input/output without learning
1. Status: **LEGAL_LOCKED** (treated conservatively — see the Revision Audit's discussion of the unresolved scope ambiguity between this and M)
2. Frontend: Web Speech API (browser-native, no account needed) or a hosted STT/TTS provider
3. Backend: if using a hosted provider, a proxy endpoint to avoid exposing provider keys client-side
4. DB/schema: none, if genuinely no session/transcript is retained
5. External APIs (provider not chosen, UNKNOWN — 3 viable options): (a) browser-native Web Speech API — free, no account, but inconsistent cross-browser support and lower quality; (b) a hosted STT/TTS provider (e.g., a general-purpose speech API) — paid, higher quality, requires an account; (c) reuse Anthropic's existing relationship if/when a voice-capable model endpoint becomes relevant — speculative, not confirmed available for this use case
6. Provider accounts: depends on option chosen above
7. Env vars: depends on provider (names TBD once selected)
8. OAuth: no, typically simple API-key auth for most STT/TTS providers
9. Webhooks: no. 10. Cron: no.
11. PWA/service-worker: no additional requirement beyond microphone permission handling
12. Storage/bucket: only if any audio is temporarily buffered — should be avoided/minimized given the legal sensitivity
13. Legal/privacy work: required — precisely the scope-boundary question flagged in the Revision Audit (does "no learning" genuinely mean nothing is retained even transiently?)
14. DPIA before activation: conditionally yes — depends entirely on Bruno's legal counsel confirming that a genuinely transient, non-learning voice I/O implementation falls outside the Codex's binding block; this determination itself should come from the legal review, not be assumed by engineering
15. Dormant/private: **DO_NOT_BUILD_YET** until the scope question is resolved with Bruno and legal — building even a "just I/O" version without that resolution risks building something that turns out to need the same DPIA as M anyway
16. Runs on current stack: likely yes for browser-native option; UNKNOWN for hosted providers pending selection
17. New paid subscription: depends on provider chosen
18. Cost/risk: UNKNOWN pending legal scope resolution — this is the actual blocker, not cost
19. Complexity: **MEDIUM**, once unblocked
20. Main technical blocker: none significant once legally cleared
21. Main operational blocker: provider selection
22. Main legal blocker: the unresolved scope question itself — this is the single most important open question this audit surfaces for C/L
23. Minimum safe V1: N/A until legal clears scope
24. Dormant/private build strategy: N/A — do not build yet
25. Activation condition: explicit legal sign-off that transient voice I/O (no learning/storage) is out of scope for the Codex's binding block, or a completed DPIA if it isn't
26. Live test required: N/A yet
27. Include in next packet: **NO** — flag for Bruno/legal decision first, not an engineering packet

## M. Voice/video/session learning
1. Status: **LEGAL_LOCKED** — confirmed against every one of the 8 required safeguards in the classification rule, 0 present (full detail in the Revision Audit)
2. Frontend/Backend/DB/APIs: **DO_NOT_BUILD_YET** — no speculative technical scoping performed here, per the repo's own explicit instruction not to schedule this into any sprint
3. Provider accounts/env vars: **DO_NOT_BUILD_YET**
4. Legal/privacy work: required — DPIA, explicit consent flow, privacy notice update, retention policy, deletion/export controls, biometric/health-data safeguards, processor list, activation gating — all 8, per this task's own classification rule
5. DPIA before activation: **YES, explicitly and non-negotiably**, per `PRODUCT_OPERATING_SYSTEM_v0.5.md` line 723
6. Dormant/private: **DO_NOT_BUILD_YET**
7. Runs on current stack: not evaluated — premature until legally cleared
8. Cost/risk: **HIGH** (both financial — a real speech/video processing provider is a real ongoing cost — and legal/reputational)
9. Complexity: **VERY LARGE**, whenever it is eventually cleared
10. Main technical blocker: N/A — legal is the blocker, not technology
11. Main operational blocker: N/A
12. Main legal blocker: complete absence of DPIA/legal review — the sole and total blocker
13. Minimum safe V1: N/A
14. Include in next packet: **NO — explicitly excluded, per the Codex's own binding instruction**

## N. Brave Star 3 behavioural states
1. Status: **SPEC_ONLY**
2. Frontend: a visual state-indicator component (replacing "the single decibel-line concept") — no other spec detail available in this repo
3. Backend: SPEC_MISSING — unclear whether the 3 states are driven by real backend signals (e.g., conversation activity, alert status) or purely client-side animation; this determines whether backend work is needed at all
4. DB/schema: none apparent, pending spec
5. External APIs: none apparent
6. Provider accounts: none
7. Env vars: none
8. OAuth: no. 9. Webhooks: no. 10. Cron: no. 11. PWA: no. 12. Storage: no.
13. Legal/privacy work: no
14. DPIA: no
15. Dormant/private build first: DORMANT_BUILD_OK once the spec exists — visual-only, low-risk to build behind a flag
16. Runs on current stack: likely yes, pending spec
17. New paid subscription: **NO**, most likely
18. Cost/risk: **FREE/LOW**, pending spec
19. Complexity: UNKNOWN — depends entirely on the missing spec; could be SMALL (pure CSS/animation) or MEDIUM (backend-state-driven)
20. Main technical blocker: SPEC_MISSING — cannot be built accurately without the Codex's fuller definition
21. Main operational blocker: none beyond the spec gap
22. Main legal blocker: none
23. Minimum safe V1: N/A until spec exists
24. Dormant/private build strategy: build behind the feature-flag layer (X) once both exist, preview to Bruno only first
25. Activation condition: spec obtained + flag layer exists
26. Live test required: visual QA, not a functional live test
27. Include in next packet: **NO** — blocked on spec acquisition, not an engineering task yet

## O. Trust Dot
1. Status: **SPEC_ONLY**
2. Frontend: a visual trust indicator — no spec detail beyond the name exists in this repo
3. Backend: SPEC_MISSING — likely depends on the audit-events system (Y) existing first, to have something real to indicate trust about
4. DB/schema: depends on Y
5. External APIs: none apparent
6. Provider accounts: none
7. Env vars: none
8. OAuth: no. 9. Webhooks: no. 10. Cron: no. 11. PWA: no. 12. Storage: no.
13. Legal/privacy work: no
14. DPIA: no
15. Dormant/private build first: DORMANT_BUILD_OK once Y exists to feed it real data — building it before Y would risk it becoming exactly the kind of decorative-only ghost UI already flagged elsewhere in this audit
16. Runs on current stack: likely yes
17. New paid subscription: **NO**
18. Cost/risk: **FREE/LOW**
19. Complexity: SMALL (UI) + dependent on Y's complexity (MEDIUM-LARGE) for it to mean anything real
20. Main technical blocker: SPEC_MISSING, and a real dependency on Y not yet existing
21. Main operational blocker: none
22. Main legal blocker: none
23. Minimum safe V1: N/A until spec + Y exist
24. Dormant/private build strategy: build only after Y, to avoid recreating the "real-looking but wired to nothing" pattern already found repeatedly in this audit
25. Activation condition: Y exists and produces real events to visualize
26. Live test required: yes, once built
27. Include in next packet: **NO** — sequenced after Y, not this round

## P. Location-aware background / city / weather / day-night theme
1. Status: **SPEC_ONLY**
2. Frontend: a rendering layer (skyline/gradient shifting by time-of-day + weather condition)
3. Backend: a thin proxy endpoint to the chosen weather API (to avoid exposing a weather-API key client-side)
4. DB/schema: none required — could cache the profile's lat/lng already captured during onboarding
5. External APIs (provider not chosen, UNKNOWN — 3 viable options): (a) OpenWeatherMap — widely used, has a free tier; (b) WeatherAPI.com — free tier available; (c) Open-Meteo — genuinely free, no API key required at all, which would avoid the secret-management question entirely and is worth strong consideration for exactly that reason
6. Provider accounts: depends on choice above; Open-Meteo would need none
7. Env vars: depends on choice; `WEATHER_API_KEY` if a keyed provider is chosen, none if Open-Meteo
8. OAuth: no. 9. Webhooks: no. 10. Cron: no (fetch-on-load is sufficient, no need to poll continuously).
11. PWA: no. 12. Storage: no.
13. Legal/privacy work: minor — disclose the weather lookup uses the profile's stored location, not live visitor geolocation (a meaningfully lower privacy bar)
14. DPIA: no
15. Dormant/private build first: DORMANT_BUILD_OK — purely visual, low-risk, natural Cinematic Shell companion piece
16. Runs on current stack: **YES**
17. New paid subscription: **NO**, if Open-Meteo is chosen; otherwise likely still free-tier feasible
18. Cost/risk: **FREE/LOW**
19. Complexity: **SMALL-MEDIUM**
20. Main technical blocker: none significant
21. Main operational blocker: provider selection (recommend Open-Meteo specifically for its no-key simplicity)
22. Main legal blocker: none
23. Minimum safe V1: static skyline + time-of-day only, weather as a fast-follow
24. Dormant/private build strategy: build alongside Cinematic Shell, preview internally first
25. Activation condition: Cinematic Shell milestone reached
26. Live test required: visual QA
27. Include in next packet: bundle with Q (Cinematic Shell), not before

## Q. Cinematic Shell / Three.js
1. Status: **SPEC_ONLY** (`three` installed, zero usage anywhere in the app)
2. Frontend: the entire feature — a new landing/intro experience
3. Backend: none apparent
4. DB/schema: none
5. External APIs: none apparent (client-side rendering)
6. Provider accounts: none
7. Env vars: none apparent
8. OAuth: no. 9. Webhooks: no. 10. Cron: no. 11. PWA: no. 12. Storage: possibly for 3D assets (models/textures) — needs the Integration Packet to specify
13. Legal/privacy work: no
14. DPIA: no
15. Dormant/private build first: N/A — this is explicitly meant to be public-facing by design; the only gate is "not before B ships" (per the Revision Audit's own recommendation, to avoid wrapping a non-functional core chat in a polished shell)
16. Runs on current stack: **YES** for the library itself; asset hosting (3D models, textures) may need a CDN/storage decision not yet made
17. New paid subscription: UNKNOWN — depends entirely on the not-yet-produced Integration Packet
18. Cost/risk: UNKNOWN, pending the packet
19. Complexity: **LARGE**, typical for a real Three.js production feature
20. Main technical blocker: the Integration Packet itself doesn't exist yet — this is the actual blocker, not engineering capacity
21. Main operational blocker: none beyond the packet
22. Main legal blocker: none
23. Minimum safe V1: N/A until the packet exists
24. Dormant/private build strategy: N/A — public by nature once built
25. Activation condition: Integration Packet delivered + B (real chat on real pages) already shipped
26. Live test required: cross-browser/device visual + performance QA
27. Include in next packet: **NO** — sequenced after A/B/C and after the Packet is produced, not an immediate engineering task

## R. Generative Core / Build Your Universe onboarding
1. Status: mechanism = **REAL** (conventional wizard); branding = **GHOST**
2. Frontend: `components/Onboarding.jsx` already exists and works
3. Backend: `POST /profile`, `POST /auth/signup` already exist and work
4. DB/schema: none new
5. External APIs: already uses Anthropic vision for the one generative step
6. Provider accounts: none new
7. Env vars: none new
8. OAuth: no. 9. Webhooks: no. 10. Cron: no. 11. PWA: no. 12. Storage: already uses Supabase Storage for media.
13. Legal/privacy work: no
14. DPIA: no
15. Dormant/private build first: N/A — already public
16. Runs on current stack: **YES**, fully
17. New paid subscription: **NO** for the current mechanism; a genuinely more "generative" experience (more AI-driven steps) would increase Anthropic usage cost proportionally, not require a new provider
18. Cost/risk: **FREE/LOW** as-is
19. Complexity: **SMALL** to relabel copy honestly; **MEDIUM-LARGE** if genuinely rebuilding it to be more generative
20. Main technical blocker: SPEC_MISSING — the actual "Generative Core" definition isn't in this repo, so there's nothing concrete to build toward beyond the existing wizard
21. Main operational blocker: none for the current mechanism
22. Main legal blocker: none
23. Minimum safe V1: honest copy over the existing wizard (cheapest fix, immediately available)
24. Dormant/private build strategy: N/A for the copy fix; a genuinely generative rebuild could be prototyped dormant/owner-preview first
25. Activation condition: N/A for copy; spec acquisition for any real rebuild
26. Live test required: no, for copy-only fix
27. Include in next packet: copy fix is cheap enough to bundle with V (product copy) work; real rebuild is not ready to scope

## S. Memory / personalization
1. Status: storage = DORMANT_BUILT (real tables, unused for retrieval); "personalization" as a claim = **GHOST**
2. Frontend: none required for basic wiring (server-side prompt augmentation)
3. Backend: a `MemoryInterface` reading `conversations`/`messages`/`profiles` back into future prompts — already scoped conceptually in `docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md` §7/§8C
4. DB/schema: none required for basic version (re-reading existing tables); a real vector/embeddings store would be new schema/infra, explicitly V2/V3 per the Brain Spine audit
5. External APIs: none for basic version; an embeddings provider (could be Anthropic, could be a dedicated embeddings service) for real semantic memory later
6. Provider accounts: none for basic version
7. Env vars: none for basic version
8. OAuth: no. 9. Webhooks: no. 10. Cron: no. 11. PWA: no. 12. Storage: a vector store (e.g., pgvector on Supabase, or a dedicated vector DB) only for the advanced version — UNKNOWN, provider not chosen
13. Legal/privacy work: yes, proportional to scope — genuine behavioral/preference learning is closer to profiling under GDPR than simple transactional storage
14. DPIA before activation: likely yes for real behavioral learning/profiling; not for the current inert storage
15. Dormant/private build first: **YES**, strongly recommended — this is Brain Spine territory, and the Brain Spine audit itself recommends deferral
16. Runs on current stack: **YES** for a basic re-read version (Supabase already stores the data); **NO** for real semantic/vector memory without adding pgvector or an external vector store
17. New paid subscription: **NO** for basic; possibly for a dedicated vector store, depending on choice
18. Cost/risk: **LOW** for basic; **MEDIUM-HIGH** for real semantic memory (ongoing embeddings cost, storage cost, and the legal/DPIA question above)
19. Complexity: **MEDIUM** for basic re-read; **LARGE-VERY LARGE** for real semantic memory
20. Main technical blocker: none for basic version; vector-store selection for advanced
21. Main operational blocker: none
22. Main legal blocker: the profiling/DPIA question, for the advanced version specifically
23. Minimum safe V1: re-read the last N messages/leads into the prompt server-side — no new infra, immediately available
24. Dormant/private build strategy: ship basic version silently (no "personalized" claim yet), gather internal confidence before marketing it as such
25. Activation condition: basic version proven to improve response quality without legal exposure; advanced version gated behind DPIA if pursued
26. Live test required: **YES** for basic version — confirm the AI actually references prior context
27. Include in next packet: **NO** — explicitly Brain Spine/V2 territory per the existing Brain Spine audit's own sequencing; do not bundle into near-term packets

## T. News / insights / daily briefing automation
1. Status: **PARTIAL-REAL** (generation real, automation absent)
2. Frontend: already exists
3. Backend: already exists (`POST /chat` reuse, `GET/POST /owner/news`)
4. DB/schema: none new — reuses `daily_news`
5. External APIs: none new
6. Provider accounts: none new
7. Env vars: none new
8. OAuth: no. 9. Webhooks: no.
10. Cron/worker: same gap as K — no scheduler exists; this and K should share whatever scheduling solution is chosen rather than building two separate ones
11. PWA: no. 12. Storage: no.
13. Legal/privacy work: no
14. DPIA: no
15. Dormant/private build first: N/A — already public, just mislabeled
16. Runs on current stack: **YES** for on-demand (already does); **NO** for real automation without the same scheduler gap as K
17. New paid subscription: **NO**, or shares whatever K resolves
18. Cost/risk: **FREE/LOW** (each generation is one Anthropic call, already happening)
19. Complexity: **SMALL** (copy fix) or **MEDIUM** (real scheduling, shared with K)
20. Main technical blocker: same as K
21. Main operational blocker: none beyond K
22. Main legal blocker: none
23. Minimum safe V1: relabel as "on-demand insights" immediately (cheapest fix); real automation as a fast-follow once K's scheduler exists
24. Dormant/private build strategy: N/A for the label fix
25. Activation condition: label fix ships directly; automation gated on K
26. Live test required: no, for label fix
27. Include in next packet: label fix bundles cheaply with V (product copy); real automation bundles with K

## U. Stripe payments
1. Status: **UNKNOWN** (code REAL, live-config unverifiable from this environment)
2. Frontend: not deeply audited this pass (out of scope per instructions not to touch Stripe)
3. Backend: already exists, real signature verification, real fee calc
4. DB/schema: none new — `booking_payments` already exists
5. External APIs: Stripe (already integrated)
6. Provider accounts: Stripe (already exists, per the code's own env-var checks)
7. Env vars (names only, already referenced in code): `STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET`
8. OAuth: Stripe Connect uses its own onboarding-link flow (already implemented), not classic OAuth
9. Webhooks: **YES, already implemented** with real signature verification (`verifyStripeSignature`)
10. Cron: no. 11. PWA: no. 12. Storage: no.
13. Legal/privacy work: confirm payment terms/refund policy exist and are accurate (not verified this session)
14. DPIA: no — card data is Stripe's PCI-DSS responsibility, not collected directly
15. Dormant/private build first: N/A — already built; the open question is live-config, not build strategy
16. Runs on current stack: **YES**, at the code level
17. New paid subscription: **NO** new — Stripe itself takes a transaction fee, standard for the category, not a new subscription cost
18. Cost/risk: UNKNOWN until live-verified; the free-tier Render cold-start risk flagged at the top of this document is a real concern here specifically, since Stripe webhooks expect timely acknowledgment
19. Complexity: already built — **SMALL** remaining work (verification only)
20. Main technical blocker: Render free-tier cold-start could cause Stripe webhook timeouts — worth checking explicitly, since this is the one feature where backend response latency has real financial consequences
21. Main operational blocker: confirm `STRIPE_SECRET_KEY`/`STRIPE_WEBHOOK_SECRET` are set in Render; run one real test charge
22. Main legal blocker: confirm payment terms/refund policy currency
23. Minimum safe V1: already shipped at the code level
24. Dormant/private build strategy: N/A
25. Activation condition: one real end-to-end test charge succeeds and the webhook is acknowledged before Stripe's timeout
26. Live test required: **YES — the single most consequential live test in this entire matrix**, given real money is involved
27. Include in next packet: verification-only task, high priority given financial exposure, but not a build task

## V. Media upload
1. Status: **ACTIVE_PUBLIC**, pending live-bucket confirmation
2. Frontend/Backend: already exist and work in two independent flows
3. DB/schema: none new
4. External APIs: Supabase Storage (already integrated)
5. Provider accounts: Supabase (already exists)
6. Env vars (already referenced): `SUPABASE_URL`, `SUPABASE_KEY`
7. OAuth: no. 8. Webhooks: no. 9. Cron: no. 10. PWA: no.
11. Storage/bucket: the `media` bucket must exist with correct public-read policy — unverifiable from this environment
12. Legal/privacy work: no
13. DPIA: no
14. Dormant/private build first: N/A — already public
15. Runs on current stack: **YES**
16. New paid subscription: **NO** — Supabase Storage already in use
17. Cost/risk: **LOW**, contingent on the bucket-config confirmation
18. Complexity: already built — **SMALL** remaining work (verification only)
19. Main technical blocker: none apparent
20. Main operational blocker: confirm bucket exists/policy is correct
21. Main legal blocker: none
22. Minimum safe V1: already shipped
23. Live test required: **YES** — one real upload through each of the two flows
24. Include in next packet: verification-only, low priority given consistent cross-flow usage already suggests it works

## W. Legal forms / health forms
1. Status: **ACTIVE_PUBLIC** as a mechanism
2. Frontend/Backend/DB: already exist and work (`LegalFormModal.jsx`, `FormPage.jsx`, `POST /consent`, `POST /forms/:slug/:formType`)
3. External APIs: none
4. Provider accounts: none
5. Env vars: none new
6. OAuth: no. 7. Webhooks: no. 8. Cron: no. 9. PWA: no. 10. Storage: no.
11. Legal/privacy work: the specific health/medical disclosure wording should get a real legal sanity check
12. DPIA before activation: arguably yes for the health-category form types specifically (health disclosure, pregnancy disclosure) — this is the narrowest, most achievable DPIA scope identified across both audits
13. Dormant/private build first: N/A — already live; the DPIA question doesn't require un-shipping it, just formalizing the review
14. Runs on current stack: **YES**
15. New paid subscription: **NO**
16. Cost/risk: **LOW** technical risk; the DPIA gap is the real exposure, proportional to how much health data is actually being collected in production today
17. Complexity: already built
18. Main technical blocker: none
19. Main operational blocker: none
20. Main legal blocker: the missing DPIA-trigger review for health-category forms specifically
21. Minimum safe V1: already shipped; the DPIA is a parallel-track legal deliverable, not a code change
22. Live test required: no additional beyond what's already presumably in use
23. Include in next packet: not an engineering packet — flag as the "Health-Data DPIA Scoping" legal deliverable already recommended in the Revision Audit

## X. Feature registry / feature flags / activation states
1. Status: **GHOST as a system** (doesn't exist at all)
2. Frontend: a way to read a feature's state and conditionally render (dormant/coming-soon/live)
3. Backend: a way to serve those states — could be as simple as a config file/JSON deployed with the app, or a real `feature_flags` table
4. DB/schema: new, if built as a real table (recommended over a static file, so states can change without a redeploy) — genuinely new schema, and this repo has no existing migration pattern, so the same "no safe migration pattern" consideration flagged in the notification work applies here too
5. External APIs: UNKNOWN, provider not chosen — 3 viable options: (a) build it in-repo (a new Supabase table + simple admin toggle) — cheapest, most aligned with the existing stack, no new vendor; (b) GrowthBook — open-source-friendly, self-hostable or cloud, has a generous free tier; (c) LaunchDarkly — industry-standard, but likely overkill in cost/complexity for this product's current size
6. Provider accounts: none, if built in-repo (recommended); one, if (b) or (c) chosen
7. Env vars: none for in-repo option; an SDK key for (b)/(c)
8. OAuth: no. 9. Webhooks: no. 10. Cron: no. 11. PWA: no. 12. Storage: no.
13. Legal/privacy work: no
14. DPIA: no
15. Dormant/private build first: this feature is itself the mechanism for building things dormant — recommend building it first, simply, in-repo
16. Runs on current stack: **YES**, for the recommended in-repo approach
17. New paid subscription: **NO**, if built in-repo (recommended)
18. Cost/risk: **FREE/LOW**
19. Complexity: **MEDIUM** for a real, reusable version; **SMALL** for a minimal first cut (a single new table + one backend check + one frontend hook)
20. Main technical blocker: no schema/migration pattern exists in this repo yet for any new table — this would be the first, and needs the same "explain exactly what's needed, don't improvise" discipline already applied to the notification work
21. Main operational blocker: none
22. Main legal blocker: none
23. Minimum safe V1: a single `feature_flags` table (name, state, updated_at), one backend endpoint to read it, one frontend hook to consume it
24. Dormant/private build strategy: N/A — this is the enabling infrastructure, not a feature to gate itself
25. Activation condition: ship directly, it unlocks everything else
26. Live test required: **YES** — confirm toggling a flag actually changes visible behavior
27. Include in next packet: **YES — high priority**, recommended as packet 3 in the sequencing below

## Y. Audit events / provenance / trust logs
1. Status: **SPEC_ONLY** as a general system; the notification work this week is a real, working instance of the pattern, just not generalized
2. Frontend: minimal — mostly a backend/data concern; a dashboard view of events would be a nice-to-have, not core
3. Backend: a shared logging function/table that every "we did X" claim writes to, generalizing `dashboard_record_created`/`email_status` from the notification work
4. DB/schema: new — an `audit_events` (or similarly named) table; same "no migration pattern exists" consideration as X
5. External APIs: none required
6. Provider accounts: none
7. Env vars: none
8. OAuth: no. 9. Webhooks: no. 10. Cron: no. 11. PWA: no. 12. Storage: no.
13. Legal/privacy work: no
14. DPIA: no
15. Dormant/private build first: N/A — this is internal infrastructure, not a user-facing claim to gate
16. Runs on current stack: **YES**
17. New paid subscription: **NO**
18. Cost/risk: **FREE/LOW**
19. Complexity: **MEDIUM** (designing a genuinely reusable event schema, not just copy-pasting the notification pattern once)
20. Main technical blocker: none significant — mostly a design/schema question
21. Main operational blocker: none
22. Main legal blocker: none
23. Minimum safe V1: generalize the exact pattern already proven this week (real record + explicit status + never claim success without a backend signal) into one reusable table/function
24. Dormant/private build strategy: N/A
25. Activation condition: ship directly once designed
26. Live test required: **YES** — confirm a real event is written for at least one more claim beyond notifications (e.g., booking-request creation)
27. Include in next packet: reasonable to bundle with X (both are foundational infrastructure), not urgent enough to precede A/B/C

## Z. Owner dashboard / assistant settings
1. Status: **PARTIAL-REAL** (extensive real CRUD; sound/notification prefs inert; personality settings real)
2. Frontend: already exists, mostly functional
3. Backend: already exists for the real parts (~12 `/owner/*` endpoints)
4. DB/schema: none new required for the fixes identified (I, J bundle here)
5. External APIs: none new
6. Provider accounts: none new
7. Env vars: none new
8. OAuth: no. 9. Webhooks: no. 10. Cron: no. 11. PWA: no. 12. Storage: no.
13. Legal/privacy work: no
14. DPIA: no
15. Dormant/private build first: N/A — already public
16. Runs on current stack: **YES**
17. New paid subscription: **NO**
18. Cost/risk: **FREE/LOW**
19. Complexity: **SMALL** for the identified fixes (silent-failure surfacing, domain correction already covered in B, sound/prefs wiring already covered in I/J)
20. Main technical blocker: none
21. Main operational blocker: none
22. Main legal blocker: none
23. Minimum safe V1: already mostly shipped; add visible (not silent) error states on the initial data-load `Promise.all` failure path
24. Dormant/private build strategy: N/A
25. Activation condition: ship directly
26. Live test required: **YES** — force a fetch failure and confirm the owner sees an error, not a silent empty state
27. Include in next packet: bundle the error-surfacing fix with B, since both touch `OwnerDashboard.jsx`'s data-loading logic

---

# External Services / Accounts Required
- **Already integrated (live-config partially unverified):** Anthropic (chat — confirmed live), Supabase (DB/Storage/Auth — confirmed live), Resend (email — integration real, live env vars unconfirmed), Stripe (payments — integration real, live env vars unconfirmed), Render (backend hosting — confirmed live, free tier confirmed via repo's own doc).
- **Needed, not yet integrated:** Google Cloud Console project (Calendar OAuth, E/F — may partially already exist per the interrupted Blocco A3 work), a weather API (P — recommend Open-Meteo, no key required), a scheduler (K/T — recommend Render Cron Job, pending plan-tier confirmation), a speech/TTS provider (C/L — not selected, blocked on legal scope first).
- **Explicitly not needed:** any new paid feature-flag SaaS (recommend building X in-repo); any new push-notification provider beyond self-generated VAPID keys (H).

# Required Env Vars / Secrets Checklist (names only, no values)
- Already in code, live status to confirm: `ANTHROPIC_API_KEY`, `SUPABASE_URL`, `SUPABASE_KEY`, `RESEND_API_KEY`, `OWNER_EMAIL`, `RESEND_FROM_EMAIL`, `STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET`, `AUTH_SALT_SECRET`, `ADMIN_KEY`, `OWNER_KEY`.
- Net-new, not yet created: `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_OAUTH_REDIRECT_URI` (E), `VAPID_PUBLIC_KEY`, `VAPID_PRIVATE_KEY`, `VAPID_SUBJECT` (H), a weather-provider key only if a keyed provider is chosen over Open-Meteo (P), a speech-provider key set (C/L, pending selection and legal clearance).

# Paid Subscription / Cost Exposure List
- **Confirmed real ongoing costs today:** Anthropic API usage, Render hosting (free tier — upgrade explicitly required before real client links go out, per the repo's own Master Doc pricing rule), Resend (has a free tier, volume-dependent), Stripe (transaction-fee based, no subscription).
- **Potential new costs:** Render Cron Job or plan upgrade (K/T, pending tier confirmation), a weather API if a keyed provider is chosen over the free Open-Meteo option (P — avoidable), a speech/TTS provider (C/L, cost UNKNOWN pending provider selection and legal clearance), Google Calendar API (E — free tier almost certainly sufficient at this product's scale).
- **Explicitly avoidable costs:** a paid feature-flag SaaS (X — recommend in-repo instead), a paid weather API (P — recommend Open-Meteo).

# Legal / DPIA Locked Features
- **M (Voice/video/session learning):** LEGAL_LOCKED, 0/8 required safeguards present, do not schedule.
- **C/L (Voice I/O without learning):** LEGAL_LOCKED pending scope resolution — DO_NOT_BUILD_YET until Bruno/legal explicitly separate it from M.
- **S (Memory/personalization, advanced/semantic version):** conditionally DPIA-relevant if genuine behavioral profiling is built — the current inert storage is not locked, but the advanced version should be scoped with legal input before building.
- **W (Health/legal forms):** not locked, but carries an unaddressed DPIA-trigger gap for the health-category form types specifically — the narrowest, most achievable legal deliverable in this whole audit.

# Dormant Build Candidates
N (Brave Star, pending spec), O (Trust Dot, pending Y), P (Location-aware background), H (Push, owner-opt-in), F (Manager Agent connectors beyond the first), S (basic memory re-read, unmarketed).

# Do-Not-Build-Yet Features
M (Voice/video/session learning) — explicit, binding. C/L (Voice I/O) — pending legal scope resolution, treat with the same caution until resolved.

# Quick Wins with No External Dependency
B (domain correction), C-consent (persist AI-processing consent, reuses existing table), I/J (wire notification/sound prefs to already-real alert data), R-copy (relabel "Generative Core" honestly), T-copy (relabel "daily" honestly), Z (surface silent dashboard failures).

# Features needing new schema/migrations
X (feature_flags table), Y (audit_events table), F (connectors/integrations table, if built as a real generalized system rather than starting with E alone) — all three share the same open question: this repo has no established migration pattern for any table, so the first one built should set a deliberate precedent (in-repo, versioned, reviewed) rather than improvising ad hoc, exactly as flagged for the notification work.

# Features needing live config verification
G (Resend env vars), U (Stripe env vars + one real test charge — highest financial stakes in this list), V (Supabase `media` bucket policy), A (real slug -> real AI reply in production), C (privacy page reachability + consent persistence).

# Features needing provider selection
K/T (scheduler — recommend Render Cron Job, pending plan confirmation), P (weather — recommend Open-Meteo), C/L (speech/TTS — blocked on legal scope first, selection secondary), X (feature-flag mechanism — recommend in-repo).

# Recommended next 5 implementation packets

1. **"Consent, Domain & Real Chat"** (bundles A + B + C) — the three items sharing the same "stale architecture / live production defect" root cause, all SMALL-MEDIUM complexity, zero new infrastructure, highest product-trust value, and the only items with a live-confirmed active defect (`/privacy` 404).
2. **"Feature Flag Foundation"** (X, with Y as a closely-related follow-on) — SMALL-MEDIUM complexity, zero new external cost, and structurally unlocks safely shipping N/O/P/H as dormant rather than either hiding or overclaiming them — the single highest-leverage infrastructure investment identified across both audits.
3. **"Brave PA Capability Truth + Calendar Connector"** (bundles D + E) — closes an active misrepresentation (D) and completes the closest-to-done unfinished feature in the repo (E, per the Master Doc's own "last confirmed step" note), MEDIUM complexity, LOW cost (Google Calendar API free tier).
4. **"Notification Maturity"** (bundles H + I + J + K, verify G) — MEDIUM-LARGE complexity, LOW-MEDIUM cost pending scheduler/plan confirmation, but every piece is well-scoped with no legal blockers, and directly extends the truthfulness work already proven this week.
5. **"Live Config & Financial Verification"** (U + V + G, cross-cutting) — not new building, but the highest-stakes unverified items (real money via Stripe, real email delivery, real file storage) — should happen in parallel with the above, not after, given the financial exposure of an unverified Stripe webhook path on a free-tier backend with cold-start risk.
