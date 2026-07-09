# Master Vision vs Reality — Revision Audit

**AI:** Claude Code
**Mode:** Repo Execution Mode / Master Vision vs Reality Revision Audit Only
**Date:** 2026-07-09
**Status:** Audit only. No code changed.

This is the first complete revision pass of `docs/audits/MASTER_VISION_VS_REPO_REALITY.md`, extended with live production verification (direct HTTP checks against `bravebybruno.com` and the backend), 22 concepts (A-V), and activation-state classification.

**MASTER DOCUMENT NOT FOUND IN REPO.** Confirmed again by direct search — no Master Document, Codex, or ecosystem spec file exists anywhere in this repository. This audit uses only real repo text (quoted where used) plus fresh live production evidence gathered via direct HTTP checks, which upgrades several findings from code-inferred to confirmed-live.

## Live evidence gathered this pass

- `https://www.bravebybruno.com/privacy` → **404** (confirmed live, `curl -L`). Both "Privacy Policy" and "Terms of Service" links in `Chat.jsx`'s consent screen point here — Terms of Use is a section inside `privacy.html`, not a separate file, so this single 404 breaks both legal disclosures simultaneously.
- `/demo/bruno` → 200 (live). `/theconcierge/owner-auth` → 200 (live). Backend `/health` → real live response (`{"status":"ok",...}`).
- `https://concierge-ai-gamma.vercel.app/` → 200, and serving the **same current build** as `bravebybruno.com` (not a frozen/stale deployment — it's Vercel's auto-assigned project alias for the same deployment). This revises the domain-bug finding from the Product Reality Matrix: it's not "broken," it's a live, functioning duplicate URL — but the repo's own **Anti-Chaos Rule #17** explicitly forbids exactly this pattern: "Verify only against the real production domain - a Vercel preview/alias URL can silently diverge from `bravebybruno.com` and must never be used to declare something PASS" (`PRODUCT_OPERATING_SYSTEM_v0.5.md` §14.17). Hardcoding that alias into every owner's share link is a direct, live-confirmed violation of the project's own binding rule, independent of whether it currently "works."

---

## A. The Concierge as ecosystem operating engine
1. **Vision:** UNDOCUMENTED IN REPO.
2. **Repo implementation:** one standalone product under `/theconcierge`; root `/` is an empty container.
3. **Perception:** honest — root says "More to come."
4. **Structure:** frontend routing split only; no backend "ecosystem" concept; one Go backend serves Concierge alone.
5. **End-to-end function:** N/A — nothing to function as an engine yet.
6. **Live/build-tested:** build-tested only.
7. **External service/API needed:** unknown, undefined.
8. **Legal/privacy review needed:** no.
9. **DPIA needed:** no.
10. **Activation state:** SPEC_ONLY (not even that — no spec exists in-repo).
11. **To make real:** Master Document must define what "operating engine" means structurally.
12. **To make premium:** N/A until defined.
13. **Must not market yet:** "master interface of the ecosystem" — nothing backs this claim today.
14. **Build timing:** after Core Flow QA — ecosystem architecture, not launch-critical.

## B. Public profile pages / real chat on real pages
1. **Vision:** every professional's public page should let a visitor talk to a real AI concierge (the product's stated core pitch, per its own on-page copy).
2. **Repo implementation:** `app/[slug]/page.jsx` → `PublicProfileClient.jsx` → `lib/generateBusinessPage.js`. The "Chat with me" bubble is fully static/scripted HTML — hardcoded greeting, fake typing animation, hardcoded quick-reply pills. Confirmed via repo-wide grep: the real `<Chat/>` component is used only in `app/demo/[id]/page.jsx`, nowhere in the real public-profile path. Every chat action opens a new tab to a hardcoded `conciergeURL` (the Vercel alias, see above) which routes back to this same static page and ignores the `?q=` param entirely.
3. **Perception:** a real visitor believes they're talking to a live AI; they are talking to canned JavaScript.
4. **Structure:** frontend-only fiction for the chat function; real backend (`GET /profile/:slug`) for the rest of the page's data (services, links).
5. **End-to-end function:** NO — the chat path structurally cannot reach the AI backend from this page.
6. **Live/build-tested:** the page itself is live (build-tested + reachable), but the chat function was never live-testable because it doesn't exist as a real call.
7. **External service needed:** none additional — it would just need to reuse the existing Anthropic path already live elsewhere.
8. **Legal/privacy review needed:** no more than the existing chat consent requirements.
9. **DPIA needed:** no.
10. **Activation state:** **GHOST**.
11. **To make real:** mount the real `<Chat/>` component (or equivalent) on a reachable route for non-demo slugs.
12. **To make premium:** an in-page chat panel (not a new-tab redirect) with the same polish as the demo experience.
13. **Must not market yet:** "chat with your AI concierge" on real business pages — the single most important thing not to claim publicly right now.
14. **Build timing:** before Cinematic Shell — this is core product function, not visual polish; shipping a visual shell around a non-functional core chat would compound the misrepresentation.

## C. Voice-first assistant interface
1. **Vision:** stated by Bruno as the product's north star; not separately defined from "Voice & Video Session Learning" anywhere in the repo.
2. **Repo implementation:** none. Zero speech/audio code anywhere (confirmed by repo-wide grep, repeated this session).
3. **Perception:** N/A.
4. **Structure:** none.
5. **End-to-end function:** none.
6. **Live/build-tested:** neither.
7. **External service needed:** yes — a speech-to-text and/or TTS provider, not yet selected anywhere in this repo.
8. **Legal/privacy review needed:** yes, per the repo's own binding text (see L below) if any session data is retained/learned from; unclear if simple transient voice I/O (no storage) would need the same review — this ambiguity is unresolved.
9. **DPIA needed:** conditionally yes, pending Bruno's clarification of scope.
10. **Activation state:** **LEGAL_LOCKED** (treated conservatively, per repo's own grouping of voice with biometric/video/session-learning).
11. **To make real:** resolve the scope ambiguity with Bruno, then the legal review the repo already mandates.
12. **To make premium:** N/A until unblocked.
13. **Must not market yet:** any voice claim, in any framing.
14. **Build timing:** after legal/privacy review — explicitly, per the Codex's own binding language, "do not schedule it into a sprint."

## D. Brave PA / owner assistant
1. **Vision:** adaptive assistant/guide (task framing); Codex excerpt separately confirms Google Calendar OAuth (Blocco A3) was mid-implementation.
2. **Repo implementation:** `components/BravePAv2.jsx`, real conversational AI via shared `/chat` endpoint. System prompt (`buildBravePAPrompt()`) claims web search, Google Calendar read/write, ticket booking — none of these are backed by real integrations. No web-search tool wired server-side, no calendar endpoint anywhere in `main.go` despite `Profile.GoogleRefreshToken` existing as an unused schema field.
3. **Perception:** an owner asking "add this to my calendar" gets a plausible, false response.
4. **Structure:** frontend + shared `/chat` backend (real); calendar/search claims = fiction, no structure behind them.
5. **End-to-end function:** conversation itself, yes; claimed tool capabilities, no.
6. **Live/build-tested:** build-tested; conversational core not confirmed live-tested this session (demo chat was; Brave PA specifically was not).
7. **External service needed:** Google Calendar API (OAuth), a web-search provider — neither present.
8. **Legal review needed:** only if/when calendar OAuth is resumed (standard OAuth scope disclosure, not DPIA-level).
9. **DPIA needed:** no.
10. **Activation state:** PARTIAL-REAL (conversation core) / **GHOST** (claimed tool capabilities).
11. **To make real:** resume Blocco A3 (Google Calendar OAuth — genuinely closer to done than most features, per the Master Doc excerpt) or strip the false claims from the system prompt.
12. **To make premium:** real calendar read/write would be a genuinely differentiating feature once built.
13. **Must not market yet:** "manages your calendar," "searches the web for you," or any specific tool claim.
14. **Build timing:** Blocco A3 resumption is explicitly scoped in the Master Doc excerpt as before or alongside Phase B1 — i.e., before Cinematic Shell, per the repo's own stated sequencing.

## E. Brave Star 3 behavioural states
1. **Vision:** real, quoted — `PRODUCT_OPERATING_SYSTEM_v0.5.md` §7.6: "Brave Star - 3 behavioural states (replaces the single decibel-line concept) - no external blocker, ready now." No further spec exists in-repo.
2. **Repo implementation:** none (confirmed by grep).
3. **Perception:** N/A.
4. **Structure:** none. **Function:** none. **Tested:** neither.
7. **External service needed:** none apparent from the one-line description.
8. **Legal review needed:** no. 9. **DPIA needed:** no.
10. **Activation state:** SPEC_ONLY (cleared but unspecified beyond a name).
11. **To make real:** the Codex's fuller spec (not in this repo) for what the 3 states actually are/display.
12. **To make premium:** this is explicitly a premium-identity feature by its own description ("replaces the single decibel-line concept").
13. **Must not market yet:** anything specific about it — there's no spec to build to or claim yet.
14. **Build timing:** cleared "ready now," independent of other blockers — candidate for before or during Cinematic Shell since it's part of visual identity.

## F. Generative Core / Build Your Universe onboarding
1. **Vision:** UNDOCUMENTED IN REPO.
2. **Repo implementation:** `components/Onboarding.jsx`, a conventional multi-step form wizard. Only generative moment: one-shot vision-based service-import from a screenshot.
3. **Perception:** feels like filling out a form, not "building a universe."
4. **Structure:** real frontend wizard + real backend save (`POST /profile`, `POST /auth/signup`).
5. **Function:** produces a real, working profile — genuinely REAL as a mechanism.
6. **Live/build-tested:** signup/account creation live-verified per AUTH GREEN work.
7. **External service needed:** none beyond what's already used (Anthropic vision, Supabase Storage).
8. **Legal review needed:** no. 9. **DPIA needed:** no.
10. **Activation state:** the underlying wizard = ACTIVE_PUBLIC; the "Generative Core / Build Your Universe" branding = GHOST (doesn't exist in any form users see).
11. **To make real:** the actual Master Document spec for what "Generative Core" adds beyond the current wizard.
12. **To make premium:** more AI-driven generation throughout the flow, not just one screenshot-import step.
13. **Must not market yet:** "Generative Core," "Build Your Universe" — until the mechanism matches the name.
14. **Build timing:** during Cinematic Shell if it's copy/presentation; after Core Flow QA if it implies real new mechanics.

## G. Cinematic Shell
1. **Vision:** "Three.js Cosmic Intro" for the public landing experience (`PRODUCT_OPERATING_SYSTEM_v0.5.md` §1007/§1290; Decision Pack §9).
2. **Repo implementation:** `three` is a listed dependency; zero imports anywhere in `app/`, `components/`, or `lib/` (confirmed by repo-wide search, again this session).
3. **Perception:** N/A — most-referenced upcoming milestone across every doc, zero visible progress.
4. **Structure:** dependency only. **Function:** none. **Tested:** neither.
7. **External service needed:** none apparent (client-side rendering library).
8. **Legal review:** no. 9. **DPIA:** no.
10. **Activation state:** SPEC_ONLY (an Integration Packet is repeatedly referenced as the next deliverable but has not been produced/saved to this repo).
11. **To make real:** the actual Integration Packet from ChatGPT/Gemini.
12. **To make premium:** this is inherently the premium-presentation feature.
13. **Must not market yet:** nothing — it's not built, so nothing to overclaim yet.
14. **Build timing:** this is the "during Cinematic Shell" milestone by definition — but should not start, per B above, until real chat exists on real public pages, or the shell will wrap a non-functional core.

## H. Notifications / email / push / digest / preferences
1. **Vision:** Bruno's explicit product decision this session — notification claims must be backend-real; push is "essential to the long-term product"; preferences are part of product direction.
2. **Repo implementation:** sensitive-topic alerts are REAL end-to-end (durable record + synchronous email attempt + explicit status, commits `412312e`/`d69daf0`). Push: `service_worker.js` exists with a working `push` listener but is unregistered (no `public/` folder, no `serviceWorker.register()` call anywhere), no VAPID keys, no subscription storage, no Go web-push library. Digest/clustering: no scheduler exists anywhere in the codebase. Preferences: `notifPrefs` in `OwnerDashboard.jsx` is `localStorage`-only, never sent to the backend, and nothing reads `soundEnabled`/`soundStyle` to actually play a sound (zero `Audio`/`.play()` calls anywhere).
3. **Perception:** now accurate for the alert-record/email path (UI copy reflects real outcome as of this week); sound/push toggles look interactive but do nothing.
4. **Structure:** alert path = real frontend + real backend + real storage (`notes` table, reused) + real Resend email. Push/digest = no structure. Prefs = frontend-only, unconsumed.
5. **End-to-end function:** alert record + email, yes. Push, digest, and preference-driven behavior, no.
6. **Live/build-tested:** build-tested (`go build`, `go vet`, `npm run build` all clean this week); not live-tested — whether `RESEND_API_KEY`/`OWNER_EMAIL` are actually set on Render is unconfirmed from this environment.
7. **External service needed:** Resend (already integrated, live-config unknown); for push: a VAPID key pair (new secret, not created); for digest: a real scheduler (Render Cron Job or equivalent, doesn't exist).
8. **Legal review needed:** no (transactional notification, not sensitive-category data).
9. **DPIA needed:** no.
10. **Activation state:** alert record + email = ACTIVE_PRIVATE (real but live-config unverified); push = DORMANT_BUILT at best (only a client-side listener exists, no delivery pipeline) — more accurately GHOST since nothing can actually trigger it; digest = SPEC_ONLY; sound/notification preferences = GHOST.
11. **To make real:** confirm Render env vars for email; build VAPID+subscription+registration for push; build a real scheduler for digest; wire `notifPrefs` into actual behavior.
12. **To make premium:** real push + a genuinely adaptive digest would be a strong differentiator once built.
13. **Must not market yet:** push notifications in any form; "daily digest"; sound/notification preferences as functional.
14. **Build timing:** push and digest are explicitly documented as blocked on new infrastructure — after Core Flow QA. Preference-wiring is small and contained — fine during Core Flow QA.

## I. Consent Truth Path / AI-processing consent
1. **Vision:** the product's own consent screen claims UK/EU GDPR compliance and links to a Privacy Policy and Terms of Service.
2. **Repo implementation:** `Chat.jsx`'s general AI-processing consent (`giveConsent()`) writes only to `sessionStorage`, never persisted server-side, despite explicitly citing GDPR compliance in its own copy. Live-confirmed this session: the Privacy Policy/Terms links both point to `/privacy`, which returns a real 404 in production. A genuinely thorough `privacy.html` document exists in the repo (GDPR rights sections, Anthropic-as-processor disclosure, retention policy, terms of use) — but it is not reachable, because there is no `public/` folder and no rewrite rule for it (confirmed by reading `next.config.js`/`vercel.json`).
3. **Perception:** a visitor believes they've read and can access a compliant privacy policy; the link is dead.
4. **Structure:** consent capture = frontend-only, unpersisted. Privacy document = exists in repo, unreachable in production.
5. **End-to-end function:** NO, in two independent ways — consent isn't stored, and the disclosure it references can't be read.
6. **Live/build-tested:** live-tested this session — `curl -L https://www.bravebybruno.com/privacy` → 404, confirmed directly, not inferred.
7. **External service needed:** none — this is a wiring problem, not a missing capability.
8. **Legal review needed:** the content of `privacy.html` should get a real legal review before being relied upon, but the immediate problem is technical (it's unreachable), not legal.
9. **DPIA needed:** no, but this is exactly the kind of gap a DPIA process would have caught.
10. **Activation state:** GHOST for both the consent record and the linked disclosure — the product actively claims compliance mechanisms that don't function.
11. **To make real:** (a) move/serve `privacy.html` at a real route (either a `public/` folder + rewrite, or convert it to a proper Next.js page), (b) persist the general AI-processing consent server-side, matching the pattern already used for legal-form consent (`POST /consent`).
12. **To make premium:** a real, visible "your data" control panel (view/export/delete) tied to the same consent record.
13. **Must not market yet:** "GDPR-compliant AI processing consent" — this is a direct, specific, currently false claim.
14. **Build timing:** this is a live production bug affecting a legal disclosure, already independently flagged in `docs/decisions/AI_REVIEW_LEDGER_TRUST_AND_CINEMATIC_SHELL_v0.2.md` as a blocker — before Cinematic Shell, and arguably the most urgent single item in this entire audit given it's a live 404 on a legal document, not just a design gap.

## J. Legal forms / health forms / booking forms
1. **Vision:** structured pre-booking consent (health disclosures, injury/allergy declarations, image-rights releases, etc.).
2. **Repo implementation:** `components/LegalFormModal.jsx` and `components/FormPage.jsx`, real submission via `POST /consent` and `POST /forms/:slug/:formType`, persisted to `consents`/`form_submissions` tables, visible to the owner via `GET /owner/form-submissions`.
3. **Perception:** matches reality — this is the one consent mechanism in the product that actually works as claimed.
4. **Structure:** fully real, frontend + backend + storage.
5. **End-to-end function:** yes.
6. **Live/build-tested:** build-tested; live-test status unknown (not exercised this session).
7. **External service needed:** none.
8. **Legal review needed:** the specific form language (health/medical disclosure wording) should get a real legal sanity check, but the mechanism itself is sound.
9. **DPIA needed:** arguably yes if health-category data is being collected via the "Health & medical disclosure" / "Pregnancy disclosure" form types — this touches special-category data under GDPR, which the general "input validation" checklist in the repo's own Phase B4 doesn't specifically flag.
10. **Activation state:** ACTIVE_PUBLIC as a mechanism; the health-category forms specifically warrant a DPIA-trigger review before being called "compliant."
11. **To make real:** already real; add the DPIA-trigger review for the health-data-collecting form types specifically.
12. **To make premium:** N/A — this is a compliance feature, not a differentiation feature.
13. **Must not market yet:** nothing new — this one's already safe to claim as-is.
14. **Build timing:** DPIA-trigger review for health forms specifically — after legal/privacy review, can run in parallel with other work.

## K. Memory / learning / personalization
1. **Vision:** UNDOCUMENTED IN REPO explicitly, but implied throughout (adaptive assistant, personalization).
2. **Repo implementation:** conversation history is stored (`conversations`/`messages` tables) but never re-read into a later prompt — each `/chat` call is stateless from the model's perspective beyond whatever `messages[]` the client happens to resend in that single request (confirmed in `docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md` §8C). No cross-session memory, no preference learning, no retrieval layer of any kind.
3. **Perception:** Brave PA's own system prompt claims it "notices things" and "cares about [owner]'s success" — implying persistent understanding that doesn't exist; each conversation starts cold except for whatever business data (leads, services) is freshly fetched into that turn's prompt.
4. **Structure:** storage exists (real tables); retrieval/re-use layer does not.
5. **End-to-end function:** no — nothing "remembers" across sessions in the sense a user would expect from "personalization."
6. **Live/build-tested:** the storage writes are real and presumably live; the memory concept has nothing to test because it isn't built.
7. **External service needed:** for real semantic memory, an embeddings/vector store — explicitly scoped as V2/V3 in `docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md` §11.
8. **Legal review needed:** yes, proportional to scope — persistent behavioral/preference learning about a person raises different GDPR considerations (profiling) than simple transactional storage.
9. **DPIA needed:** likely yes if genuine behavioral learning/profiling is built, not for the current inert storage.
10. **Activation state:** GHOST for "personalization" as a claim; the underlying storage is real but inert (closer to DORMANT_BUILT for the storage layer alone).
11. **To make real:** the Minimal Brain Spine's `MemoryInterface` (already scoped in the Brain Spine audit) is the natural next step — read stored history back into future prompts before claiming any personalization.
12. **To make premium:** this is precisely where the product's long-term moat is supposed to live, per the Fleet Architecture Brief's own language.
13. **Must not market yet:** "personalized," "remembers you," "learns your preferences" — none of this is true today.
14. **Build timing:** V2/V3, per the Brain Spine audit's own sequencing — this is squarely Brain Spine territory, explicitly deferred.

## L. Voice/video/session learning
1. **Vision:** the most explicitly documented item in the entire repo. `PRODUCT_OPERATING_SYSTEM_v0.5.md` line 723: "Voice & Video Session Learning stays in spec only - DPIA and GDPR/UK GDPR legal review required before any build, per the Codex's binding block. Do not schedule it into a sprint." Reinforced at lines 303, 699.
2. **Repo implementation:** none whatsoever (confirmed by repeated repo-wide search).
3. **Perception:** N/A — correctly, nothing markets this.
4. **Structure/function:** none, correctly.
5. **Live/build-tested:** neither, correctly.
6. **External service needed:** would require a speech/video processing provider — not selected, not relevant until unblocked.
7. **Legal review needed:** yes, explicitly required before any build.
8. **DPIA needed:** yes, explicitly required.
9. **Special legal/privacy rule check:** LEGAL_LOCKED requires the repo to demonstrate a DPIA, explicit consent flow, privacy notice, retention policy, deletion/export controls, biometric/health-data safeguards, processor list, and activation gating. Checking each: DPIA — not in repo. Explicit consent flow for voice/video specifically — not in repo. Privacy notice — exists (`privacy.html`) but is currently unreachable in production and doesn't mention voice/video at all. Retention policy — exists for text chat only. Deletion/export controls — described but not implementable by a user (no in-product mechanism). Biometric/health-data safeguards — none. Processor list — Anthropic named for text chat; none named for voice/video since nothing is integrated. Activation gating — none, but moot since nothing is built.
10. **Activation state:** LEGAL_LOCKED, exactly as documented — the one item where the repo's own stated policy and repo reality are in full agreement, and the strictest reading of the classification rule confirms it decisively (0 of 8 required safeguards met).
11. **To make real:** a completed DPIA and GDPR/UK GDPR legal review — explicitly not a Claude Code task.
12. **To make premium:** N/A until unblocked.
13. **Must not market yet:** anything.
14. **Build timing:** not to be scheduled into any sprint, per the Codex's own binding instruction, regardless of Cinematic Shell or Core Flow QA timing.

## M. Location-aware visual themes / weather / city ambience
1. **Vision:** real, quoted — `PRODUCT_OPERATING_SYSTEM_v0.5.md` §7.6: "Sfondo Location-Aware (skyline + weather + time of day) - no external blocker, ready now."
2. **Repo implementation:** none as a visual feature. Only location-related code: a one-time reverse-geocoding call in `Onboarding.jsx` for the profile's location text field. No weather API integration anywhere.
3. **Perception:** N/A — nothing claims this today.
4. **Structure:** none. **Function:** none. **Tested:** neither.
7. **External service needed:** a weather API (net-new, not yet selected).
8. **Legal review:** no. 9. **DPIA:** no.
10. **Activation state:** SPEC_ONLY (cleared, unbuilt).
11. **To make real:** weather API integration + skyline/time-of-day rendering layer.
12. **To make premium:** this is inherently a premium-presentation feature by its own description.
13. **Must not market yet:** nothing to overclaim — not built.
14. **Build timing:** during Cinematic Shell — visual-shell work by nature, explicitly cleared with no dependency blocking it.

## N. News / insights / daily briefing
1. **Vision:** AI market-intelligence feature, framed as "daily."
2. **Repo implementation:** real — reuses `POST /chat` with a market-intelligence prompt, stores to `daily_news`, shown in dashboard. No scheduling exists — runs only on manual "Generate now" click or first-load-if-empty-today.
3. **Perception:** "daily" implies automation that doesn't exist.
4. **Structure:** real generation + storage; no automation layer.
5. **End-to-end function:** yes, for generation; no, for the "daily" automation claim.
6. **Live/build-tested:** build-tested; live-test unknown.
7. **External service needed:** none beyond existing Anthropic integration.
8. **Legal review:** no. 9. **DPIA:** no.
10. **Activation state:** PARTIAL-REAL (real mechanism, false automation framing).
11. **To make real:** rename to reflect on-demand generation, or add a real scheduler.
12. **To make premium:** genuine daily automation once a scheduler exists.
13. **Must not market yet:** "daily" as an automated claim.
14. **Build timing:** low-stakes — after Core Flow QA, bundled naturally with the digest/scheduler work in H.

## O. Booking requests
1. **Vision:** structured booking pipeline with owner accept/decline/counter.
2. **Repo implementation:** real, end-to-end — `POST /booking-request`, `GET /owner/bookings`, `PATCH /owner/bookings/:id`, real table. Detection that "a booking happened" is a regex match on the AI's free-text reply, not a structured tool call.
3. **Perception:** matches reality for the pipeline itself; silent-miss risk if the AI phrases confirmation unexpectedly.
4. **Structure:** fully real. 5. **Function:** yes, with the fragility noted.
6. **Live/build-tested:** build-tested; live-test unknown.
7. **External service needed:** none.
8. **Legal review:** no. 9. **DPIA:** no.
10. **Activation state:** ACTIVE_PUBLIC, with a known fragility.
11. **To make real (fully robust):** replace regex detection with a structured tool-call/JSON contract.
12. **To make premium:** deterministic booking confirmation, calendar sync (ties to Manager Agent).
13. **Must not market yet:** nothing — safe as-is.
14. **Build timing:** robustness fix during Core Flow QA.

## P. Leads / CRM-lite
1. **Vision:** lead tracking with hot-lead alerting.
2. **Repo implementation:** real — deterministic keyword `scoreLead()`, real Resend hot-lead email, dashboard visibility. Score computed once at creation, never re-evaluated.
3. **Perception:** matches reality.
4. **Structure:** fully real. 5. **Function:** yes.
6. **Live/build-tested:** build-tested; live-test unknown.
7. **External service needed:** Resend (already integrated).
8. **Legal review:** no. 9. **DPIA:** no.
10. **Activation state:** ACTIVE_PUBLIC.
11. **To make real:** already real; re-scoring is a nice-to-have.
12. **To make premium:** re-scoring as conversation continues, lead-quality trends over time.
13. **Must not market yet:** nothing.
14. **Build timing:** no action needed for launch.

## Q. Media upload
1. **Vision:** upload photos/media to your profile.
2. **Repo implementation:** real — `POST /media/upload` to Supabase Storage, used in both Onboarding and `OwnerEditProfile.jsx`.
3. **Perception:** matches reality.
4. **Structure:** fully real.
5. **End-to-end function:** yes, assuming the Supabase `media` bucket is correctly configured (unverifiable from this environment).
6. **Live/build-tested:** build-tested; live-bucket status unknown.
7. **External service needed:** Supabase Storage (already integrated).
8. **Legal review:** no. 9. **DPIA:** no.
10. **Activation state:** ACTIVE_PUBLIC, pending live-bucket confirmation.
11. **To make real:** confirm bucket exists/public-read policy in live Supabase project.
12. **To make premium:** N/A — utility feature.
13. **Must not market yet:** nothing.
14. **Build timing:** confirm bucket config, low urgency.

## R. Stripe / payments
1. **Vision:** real payment processing via Stripe Connect for bookings/deposits.
2. **Repo implementation:** `concierge-backend/stripe_integration.go` — real Connect onboarding, checkout, webhook with real signature verification, platform-fee calc, payments listing. Gracefully degrades if `STRIPE_SECRET_KEY` unset.
3. **Perception:** appears real at the code level; live-config status is the only open question.
4. **Structure:** fully real at the code level.
5. **End-to-end function:** appears yes, code-level; unconfirmed live.
6. **Live/build-tested:** build-tested (compiles); live-tested: UNKNOWN.
7. **External service needed:** Stripe (integration present, live-config unconfirmed).
8. **Legal review:** standard payment terms/refund policy should be confirmed present and correct (not verified this session).
9. **DPIA needed:** no.
10. **Activation state:** UNKNOWN (code REAL, live status unverifiable from here).
11. **To make real (confirm):** verify Stripe env vars in Render, run one real test charge end-to-end.
12. **To make premium:** N/A — this is infrastructure, not a differentiator.
13. **Must not market yet:** "book and pay" until live-verified.
14. **Build timing:** verification (not building) — can happen any time.

## S. Manager Agent / connectors / calendar / email / search
1. **Vision:** real, quoted — "Manager Agent (incremental, one connector at a time) - ready now; Invyted stays excluded until it publishes a public API." Blocco A3: Google Calendar OAuth (read-only), "mid-implementation."
2. **Repo implementation:** no Manager Agent architecture exists. One orphaned schema field (`Profile.GoogleRefreshToken`) is the only trace of Blocco A3's prior progress.
3. **Perception:** Brave PA claims calendar access it doesn't have.
4. **Structure:** one dead schema field; no connector framework.
5. **End-to-end function:** none for calendar/search/Manager Agent.
6. **Live/build-tested:** neither, for anything Manager-Agent-specific.
7. **External service needed:** Google Calendar API (OAuth) is the immediately actionable one.
8. **Legal review needed:** standard OAuth scope/consent disclosure, not DPIA-level.
9. **DPIA needed:** no.
10. **Activation state:** SPEC_ONLY generally; GHOST for the specific "I manage your calendar" claim already live in Brave PA's prompt.
11. **To make real:** resume Blocco A3.
12. **To make premium:** incremental connectors exactly as the Codex excerpt specifies.
13. **Must not market yet:** any specific connector capability until built.
14. **Build timing:** before or alongside Phase B1 — before Cinematic Shell, per the repo's own stated sequencing.

## T. Trust Dot / verification / audit events
1. **Vision:** real, quoted — "Trust Dot - no external blocker, ready now."
2. **Repo implementation:** none. No general audit-event system exists either — the closest thing is the notification system's explicit `email_status`/`dashboard_record_created` fields built this week.
3. **Perception:** N/A — not built.
4. **Structure:** none. **Function:** none. **Tested:** neither.
7. **External service needed:** none apparent.
8. **Legal review:** no. 9. **DPIA:** no.
10. **Activation state:** SPEC_ONLY.
11. **To make real:** the Codex's fuller spec (not in this repo); structurally, this week's notification-status pattern is the right template to generalize.
12. **To make premium:** exactly what a "Trust Dot" implies — a visible, real trust signal tied to actually-verified events.
13. **Must not market yet:** anything specific.
14. **Build timing:** cleared "ready now" — natural to build during Cinematic Shell.

## U. Feature flags / dormant features / activation states
1. **Vision:** UNDOCUMENTED IN REPO — introduced by this audit's own vocabulary.
2. **Repo implementation:** no feature-flag system exists anywhere.
3. **Perception:** N/A.
4. **Structure:** none. 5. **Function:** N/A.
6. **Live/build-tested:** N/A.
7. **External service needed:** optional — a flag service or a simple in-repo config table would both work.
8. **Legal review:** no. 9. **DPIA:** no.
10. **Activation state:** GHOST as a system — the concept of activation states doesn't exist in the codebase at all right now.
11. **To make real:** see Recommended product architecture spine.
12. **To make premium:** enables safely shipping Trust Dot/Brave Star/Location-aware as DORMANT_BUILT ahead of full launch.
13. **Must not market yet:** N/A — infrastructure, not a claim.
14. **Build timing:** before Cinematic Shell — every "build it dormant" recommendation depends on this existing first.

## V. Product copy / launch claims
1. **Vision:** Bruno's own Anti-Chaos Rule #11: "Marketing does not outrun product truth."
2. **Repo implementation:** at least three concrete violations confirmed this session: (a) AI-processing GDPR consent claim with a dead-linked privacy policy; (b) public-page chat bubble presenting scripted text as live AI; (c) Brave PA claiming calendar/search capabilities that don't exist.
3. **Perception:** confident, specific copy in all three cases with no way for a user to know otherwise.
4. **Structure:** N/A — cross-cutting.
5. **End-to-end function:** N/A.
6. **Live/build-tested:** live-verified this session for (a); code-verified for (b) and (c).
7. **External service needed:** none.
8. **Legal review:** yes — (a) specifically.
9. **DPIA needed:** no.
10. **Activation state:** N/A (cross-cutting).
11. **To make real:** fix the three instances; adopt the "truthful sub-state copy" pattern proven this week as house style.
12. **To make premium:** consistent truthful copy is itself a trust signal.
13. **Must not market yet:** see B, D, I above.
14. **Build timing:** before Cinematic Shell, ideally alongside the structural fixes they depend on.

---

# Top 15 gaps between vision and repo

1. Public-page chat is entirely scripted (B).
2. `/privacy` 404s in live production (I).
3. AI-processing consent is never persisted server-side (I).
4. Brave PA claims calendar/search capabilities with zero backend (D, S).
5. No feature-flag/activation-state system exists (U).
6. Voice-first has no scope boundary from DPIA-blocked session learning (C, L).
7. Sound/notification preference toggles are fully inert (H).
8. Trust Dot, Brave Star, Sfondo Location-Aware — all cleared, all unbuilt, all lacking spec.
9. No memory/personalization exists despite Brave PA's own prompt implying it (K).
10. Google Calendar OAuth (Blocco A3) was mid-implementation and silently abandoned (S).
11. Vercel alias URL hardcoded into share links, violating Anti-Chaos Rule #17.
12. "Daily" news generation has no scheduler (N).
13. Push notifications have a client-side listener with zero delivery pipeline (H).
14. No general audit-event/trust system exists (H, T).
15. "Omnipresent" Brave PA is real on 2 of ~7 owner-facing routes (D).

# Features safe to market now
- Auth/session, demo pages, leads/hot-lead email, booking-request pipeline, legal/consent form collection, owner notifications (alert path), onboarding core mechanism.

# Features that must not be marketed yet
- "Chat with your AI concierge" on real pages; "share your link"; "GDPR-compliant AI processing consent"; any Brave PA specific capability claim; "daily" insights; push notifications; payments (until live-verified).

# Features that should be visually represented but labelled "coming soon"
- Trust Dot, Brave Star, Location-aware background, calendar connector — once a feature-flag layer exists.
- Voice should NOT be teased as "coming soon" without the DPIA/legal review existing first.

# Recommended next 3 implementation packets
1. **Consent & Disclosure Truth Path** — serve `privacy.html` at a real route and persist AI-processing consent server-side.
2. **Real Chat on Real Pages** — mount real `<Chat/>` on public profile pages; fix Vercel-alias hardcoding in the same pass.
3. **Capability Claims Audit for Brave PA** — strip or back the calendar/search/booking claims.

# Final one-page executive verdict for Bruno

The repo is in better shape than "prototype patchwork" in the places that matter most operationally — auth, booking, leads, legal-form consent, and (as of this week) honest notifications are all genuinely real and load-bearing. But the product's headline promise — a real AI concierge a client can talk to on a real business's real page — currently does not work outside the 5 demo personas, and the legal disclosure it should be able to point to is returning a live 404 right now, not hypothetically. Those two facts alone mean the product is not yet safe to launch publicly as described, regardless of how polished Cinematic Shell makes it look. The good news: both are narrow, well-understood fixes, not architectural rewrites, and they share a root cause (stale assumptions about domain/routing left over from an earlier version of the app) — fixing them together is efficient, not just necessary. Everything genuinely aspirational in the vision — Trust Dot, Brave Star, voice, location-aware ambience, real memory, a real Manager Agent — is honestly absent, mostly by the Codex's own admission ("ready now" but unbuilt, or explicitly legal-locked), which is a much healthier state than if the repo were pretending otherwise. The single highest-leverage structural investment available right now is a feature-flag/activation-state layer: it doesn't exist today, and its absence is why so much of this product is binary "real or ghost" instead of safely, honestly "dormant and coming soon." Build that first, fix the two live truth-gaps second, and the rest of this roadmap becomes a sequencing exercise rather than a trust problem.
