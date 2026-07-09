# Product Reality Matrix

**AI:** Claude Code
**Mode:** Repo Execution Mode / Product Reality Audit Only
**Date:** 2026-07-09
**Status:** Audit only. No code changed.

Standard applied: no feature is treated as "done" because UI or code exists. Every feature is classified by structure, functionality, user/owner perception, and real tested behaviour.

Two findings drove this audit more than anything else: the public profile page's "Chat with me" widget is **not** the real AI chat, and the share-link/QR domain baked into the dashboard is stale.

---

## 1. Public landing / root (`/`)
- **UI:** `app/page.tsx` — static "Brave by Bruno / More to come" splash with one link to `/theconcierge`.
- **Perception:** minimal — reads as a holding page, not a feature claim.
- **Actual:** exactly what it looks like. Fully static, no data, no calls.
- **Frontend:** `app/page.tsx`. **Backend:** none. **Storage:** none. **External services:** none.
- **Real delivery or UI-only:** UI-only (by design — it's a placeholder splash).
- **Failure states:** N/A. **Owner visibility:** N/A. **Build-tested:** yes (part of `npm run build`). **Live-tested:** unknown.
- **Perception risk:** LOW. **Status:** REAL (as a placeholder — it doesn't overclaim).
- **To make real:** nothing required unless product direction wants a real ecosystem landing page.
- **Blocks launch:** NO.

## 2. `/theconcierge` home
- **UI:** `app/theconcierge/page.jsx` — hero, 5 demo buttons, "Build My Concierge" CTA, "Owner Login".
- **Perception:** a working product homepage.
- **Actual:** fully functional as a router — demo buttons go to real demo pages, onboarding/login buttons go to real flows. One latent bug: its auto-redirect-if-logged-in check reads `localStorage.getItem('ownerToken')`/`'cai_owner_token'`, both legacy keys from before the httpOnly-cookie auth migration (per the `cai_token` cookie-name fix) — these will never be set by the current login flow, so that branch of the check is dead weight (harmless, since the cookie check still works via `/api/auth/session`, but it's leftover code implying a state that can't occur).
- **Frontend:** `app/theconcierge/page.jsx`, `app/theconcierge/layout.jsx`. **Backend:** `/api/auth/session` (Next route). **Storage:** none directly. **External services:** none.
- **Real delivery or UI-only:** real navigation.
- **Failure states:** none visible (no error UI, but nothing here can meaningfully fail).
- **Owner visibility:** N/A. **Build-tested:** yes. **Live-tested:** unknown.
- **Perception risk:** LOW. **Status:** REAL.
- **To make real:** remove the dead `localStorage` check (cosmetic cleanup, not urgent).
- **Blocks launch:** NO.

## 3. Public profile pages (`/[slug]`) — HIGH RISK
- **UI:** `app/[slug]/page.jsx` → `PublicProfileClient.jsx` → `lib/generateBusinessPage.js`, injected via `document.write`. This is the page every real onboarded professional shares with clients.
- **Perception:** "Your 24/7 AI concierge" — a visitor expects to tap "Chat with me" and talk to a real AI that knows this professional's services, prices, and availability (this is literally the product's core pitch, restated in the bubble's own greeting text).
- **Actual:** the "Chat with me" bubble is entirely static, scripted HTML/JS — a hardcoded greeting (`"Hi 👋 I'm {name}'s AI assistant..."`), a canned fake-typing animation, and 4 hardcoded quick-reply pills. **No AI call happens on this page at all.** Every pill and the "Ask me anything" button call `openConcierge(msg)`, which does `window.open(conciergeURL + '?q=' + msg)` in a new tab — and `conciergeURL` is a hardcoded, stale URL: `https://concierge-ai-gamma.vercel.app/${slug}`. That URL routes back to this exact same page, which ignores the `?q=` param entirely (only `components/Chat.jsx`, mounted solely on `/demo/[id]`, reads `?q=`). Result: a real visitor clicking any chat action opens a new tab that renders the identical fake bubble again — a dead end, not a chat. Confirmed via repo-wide grep that `<Chat` (the real AI component) is used only in `app/demo/[id]/page.jsx` — nowhere in the real public-profile path.
- **Frontend:** `app/[slug]/page.jsx`, `app/[slug]/PublicProfileClient.jsx`, `lib/generateBusinessPage.js`.
- **Backend:** `GET /profile/:slug` (fetches the real profile data used to render the page — this part is real).
- **Storage:** `profiles` table (read only here).
- **External services:** none (the "AI" in the bubble is fabricated client-side text).
- **Real delivery or UI-only:** UI-only theater for the chat function; real data-driven rendering for everything else on the page (services, links, contact buttons).
- **Failure states:** none — it never even attempts a real call, so it can't visibly fail; it just silently isn't real.
- **Owner visibility:** owner has no way to see this is broken from the dashboard.
- **Build-tested:** yes (renders without error). **Live-tested:** unknown — but the logic is verifiable by reading the code; this isn't a guess.
- **Perception risk:** HIGH — this is the single biggest gap between what the product claims and what it does, on the page real clients actually see.
- **Status:** GHOST (looks fully real, is entirely scripted).
- **To make real:** mount the real `<Chat/>` component (or an equivalent) reachable from a real public route for non-demo slugs, and either remove the fake bubble or make it a genuine entry point into that route.
- **Blocks launch:** YES.

## 4. Demo pages (`/demo/[id]`)
- **UI:** `app/demo/[id]/page.jsx`, one of 5 hardcoded personas (`lib/constants.js` `DEMOS`).
- **Perception:** "try the AI concierge" — and this is where it's actually true.
- **Actual:** mounts the real `<Chat/>` component, real backend `/chat` call, real Anthropic response. Consent screen is auto-bypassed (`sessionStorage.setItem('cai_consent', ...)` set unconditionally on mount) — reasonable for a demo, but worth knowing it never exercises the real consent UI.
- **Frontend:** `app/demo/[id]/page.jsx`, `components/Chat.jsx`. **Backend:** `POST /chat`. **Storage:** `conversations`/`messages` (real rows written, tagged with the demo's `profileId`, e.g. `"bruno"`, mixed in with any real profile that happened to choose the same slug historically — low collision risk but not namespaced). **External services:** Anthropic.
- **Real delivery or UI-only:** real.
- **Failure states:** yes — `Chat.jsx`'s `catch` sets a visible error message.
- **Owner visibility:** demo conversations would appear in a real dashboard only if a profile with a matching slug/profileId exists — otherwise orphaned data with no owner to see it.
- **Build-tested:** yes. **Live-tested:** yes, per prior AUTH GREEN work and prior sessions describing live chat testing.
- **Perception risk:** LOW (does what it claims).
- **Status:** REAL.
- **To make real:** nothing required.
- **Blocks launch:** NO.

## 5. Chat widget / real AI concierge chat (the component itself)
- **UI:** `components/Chat.jsx` — consent screen, message thread, sensitive-topic alert banner, review modal, lead capture.
- **Perception:** a full AI concierge with human-review and owner-notification safety nets.
- **Actual:** the AI call itself is real (`POST /chat` → Anthropic). Sensitive-topic detection is real (keyword match). Owner notification is now real as of commits `412312e`/`d69daf0` — creates a durable record and attempts email, with truthful UI copy reflecting actual outcome. "Human review" (approve/edit/block) remains client-side only — the AI reply has already left the server before the visitor's browser decides whether to show/block it; there is no backend concept of a held-for-review message.
- **Frontend:** `components/Chat.jsx`, `lib/buildPrompt.js`. **Backend:** `POST /chat`, `POST /alert`, `POST /booking-request`, `POST /lead`. **Storage:** `conversations`, `messages`, `notes` (alert records), `leads`, `booking_requests`. **External services:** Anthropic, Resend.
- **Real delivery or UI-only:** mixed — chat and notification are real; review/block is UI-only.
- **Failure states:** yes, for chat call and for alert delivery.
- **Owner visibility:** conversations, leads, bookings, and alerts are all visible in the dashboard.
- **Build-tested:** yes. **Live-tested:** login/auth confirmed live; chat itself not confirmed live-tested in this audit session.
- **Perception risk:** MEDIUM (the "human review" framing implies a safety gate that doesn't structurally exist — a determined visitor's browser could just not run the JS, and nothing on the backend would have held the reply back).
- **Status:** PARTIAL-REAL.
- **To make real:** either implement real server-side holding of flagged replies, or reword the UI so "review" isn't implied as a pre-delivery gate.
- **Blocks launch:** borderline — not by itself, but combined with #3, the "review" claim on a page nobody can actually reach for real profiles is moot until #3 is fixed.

## 6. Owner dashboard (general)
- **UI:** `components/OwnerDashboard.jsx`, tabs: Overview, Notifications, Bookings, Notes, Forms, Insights, Profiles, Settings.
- **Perception:** command center for the business.
- **Actual:** genuinely functional CRUD across leads, bookings, notes, notifications, form submissions, multi-profile switching — all backed by real endpoints and real Supabase tables. One concrete bug found: the "share your link" text, copy button, and cross-profile "View" links are hardcoded to `https://concierge-ai-gamma.vercel.app/{slug}` (lines ~482, 484, 592) instead of the current domain — a real owner today would copy/share the wrong URL. (The `origin` fallback at line 176 is lower-risk — it's an SSR-only fallback that gets overridden by `window.location.origin` in the browser.)
- **Frontend:** `components/OwnerDashboard.jsx`. **Backend:** ~12 `/owner/*` endpoints. **Storage:** most tables. **External services:** none directly (delegates to Resend via backend for emails).
- **Real delivery or UI-only:** real.
- **Failure states:** inconsistent — most fetches silently fall back to empty arrays on failure (`.catch(() => setX([]))`) rather than surfacing an error to the owner.
- **Owner visibility:** this is the owner-visibility layer.
- **Build-tested:** yes. **Live-tested:** dashboard auth-gating live-verified previously; data flows not confirmed live-tested in this audit session.
- **Perception risk:** MEDIUM (silent failures could make the owner think "no leads today" when it's actually "the fetch failed").
- **Status:** PARTIAL-REAL (real functionality, wrong domain baked into share links).
- **To make real:** replace the 3 hardcoded `concierge-ai-gamma.vercel.app` occurrences with the actual current domain or `window.location.origin`; add visible (not silent) failure states to the initial data load.
- **Blocks launch:** YES — the share-link bug directly undermines the core distribution mechanism (every owner's "share your page" action).

## 7. Brave PA
- **UI:** `components/BravePAv2.jsx`, `components/BravePASettings.jsx`, mounted in owner dashboard shell.
- **Perception:** a full personal-assistant AI with web search, calendar read/write, drafting.
- **Actual:** the live component calls the real `/chat` endpoint with a Brave-PA-flavored prompt — genuinely real conversational AI. But most of the capability list in its own system prompt (`lib/buildPrompt.js` `buildBravePAPrompt`: "WEB SEARCH... CALENDAR: Read and add events to Google Calendar... BOOKING: Find links for tickets") is aspirational text told to the model, not backed by actual tool integrations — there is no web-search tool wired server-side (that only exists in the orphaned root `brave_pa_v2.js`, confirmed unused), no Google Calendar read/write endpoint anywhere in `main.go` despite `Profile.GoogleRefreshToken` existing as a schema field, and no ticket/booking API. The AI will talk about doing these things because it's instructed to, but cannot actually do most of them.
- **Frontend:** `components/BravePAv2.jsx`, `lib/buildPrompt.js`. **Backend:** `POST /chat` (shared, undifferentiated from Concierge chat). **Storage:** `conversations`/`messages` (untagged as Brave PA vs. Concierge). **External services:** Anthropic only.
- **Real delivery or UI-only:** the conversation itself is real; most named capabilities are not backed by real integrations.
- **Failure states:** generic fallback message on error.
- **Owner visibility:** it's the owner's own tool — they experience it directly.
- **Build-tested:** yes. **Live-tested:** unknown.
- **Perception risk:** HIGH — an owner asking "add this to my calendar" will get a plausible-sounding response from an LLM that has no calendar access, which is a convincing-but-false interaction, worse than an obvious placeholder.
- **Status:** PARTIAL-REAL.
- **To make real:** either wire real tool calls (calendar, search) or rewrite the system prompt to stop claiming capabilities that don't exist server-side.
- **Blocks launch:** depends on whether Brave PA is marketed pre-launch; if it's positioned as "does your calendar," yes.

## 8. Notifications (owner notification events)
- **UI:** `components/Chat.jsx` banner, `components/OwnerDashboard.jsx` Notifications tab.
- **Perception:** as of the last two commits, this now matches reality — "owner notified" only shows after real backend confirmation, sub-states shown truthfully.
- **Actual:** durable record (`notes` table, `note_type="alert"`), synchronous email attempt via Resend with explicit status, dashboard visibility, truthful frontend copy. This is the most recently hardened feature in the repo.
- **Frontend:** `components/Chat.jsx`, `components/OwnerDashboard.jsx`. **Backend:** `POST /alert`, `GET /owner/notifications`. **Storage:** `notes` table (reused, no new schema). **External services:** Resend (email only — no push).
- **Real delivery or UI-only:** real for record + email; not real for push (documented, not implemented) and not real for clustering/digest (documented, not implemented).
- **Failure states:** yes, explicit and surfaced.
- **Owner visibility:** yes, dedicated tab.
- **Build-tested:** yes (`go build`, `go vet`, `npm run build` all clean). **Live-tested:** NO — not exercised against a running server or real Resend account; whether `RESEND_API_KEY`/`OWNER_EMAIL` are actually set on Render is UNKNOWN.
- **Perception risk:** LOW now (down from HIGH before the fix).
- **Status:** PARTIAL-REAL (record + email real; push/digest explicitly absent, not faked).
- **To make real (fully):** confirm Render env vars; build push (VAPID + subscription + service worker registration) and a real scheduler for digest.
- **Blocks launch:** NO, in its current truthful form — it no longer overclaims.

## 9. News / daily-news service
- **UI:** OwnerDashboard "Insights" tab, "Generate now" button.
- **Perception:** an AI market-intelligence feature.
- **Actual:** real — reuses `POST /chat` with a market-intelligence system prompt, parses JSON, stores to `daily_news` table, displayed on dashboard. No scheduling — it only runs when the owner manually clicks "Generate now" or on first dashboard load if no news exists yet for today. There is no actual "daily" automation (no cron), so the "daily" framing is aspirational — it's really "on-demand, once cached per day."
- **Frontend:** `components/OwnerDashboard.jsx`. **Backend:** `POST /chat`, `GET/POST /owner/news`. **Storage:** `daily_news`. **External services:** Anthropic.
- **Real delivery or UI-only:** real generation, but not real automation.
- **Failure states:** yes (`catch(e) { setNews([]); }`).
- **Owner visibility:** direct.
- **Build-tested:** yes. **Live-tested:** unknown.
- **Perception risk:** MEDIUM (the word "daily" implies automation that doesn't exist).
- **Status:** PARTIAL-REAL.
- **To make real:** either rename to reflect on-demand generation, or add a real scheduler (same infra gap as notification digest — see #8).
- **Blocks launch:** NO (low-stakes feature, not core to trust).

## 10. Onboarding
- **UI:** `app/theconcierge/onboarding/page.jsx` → `components/Onboarding.jsx`, multi-step wizard.
- **Perception:** builds a real, working concierge.
- **Actual:** genuinely functional — saves a real profile (`POST /profile`), creates a real account (`POST /auth/signup`), real media upload, real AI-vision service-import from a screenshot (`POST /ai/import-services`). Shares the same stale-domain display bug as #6 at the handle-selection step (shows `concierge-ai-gamma.vercel.app/` as the URL prefix while the user picks their slug) — cosmetically wrong but doesn't break the actual save.
- **Frontend:** `components/Onboarding.jsx`. **Backend:** `POST /profile`, `POST /auth/signup`, `POST /media/upload`, `POST /ai/import-services`, `GET /check-slug/:slug`. **Storage:** `profiles`, `sessions`. **External services:** Anthropic (vision), Supabase Storage.
- **Real delivery or UI-only:** real.
- **Failure states:** present per-step (visible error text on failed saves).
- **Owner visibility:** it's the owner's own action.
- **Build-tested:** yes. **Live-tested:** signup/account creation live-verified as part of prior AUTH GREEN work.
- **Perception risk:** LOW-MEDIUM (the domain display bug is confusing but not destructive — the actual slug/profile saved is correct, only the displayed preview URL is wrong).
- **Status:** REAL (core flow) with one PARTIAL-REAL cosmetic sub-issue (same fix as #6).
- **To make real:** fix the hardcoded domain display at the handle step.
- **Blocks launch:** NO on its own, but reinforces #6's severity — it's the same bug appearing at signup time too.

## 11. Consent / legal forms
- **UI:** `Chat.jsx`'s "Before we begin" AI-processing consent screen; `components/LegalFormModal.jsx` for booking-specific legal forms; `components/FormPage.jsx` for shareable form links.
- **Perception:** GDPR-compliant consent capture.
- **Actual:** two separate consent mechanisms with different truthfulness. The general "AI processing" consent banner in `Chat.jsx` only writes to `sessionStorage` — never persisted server-side, despite explicitly citing UK/EU GDPR compliance in its own copy. The `LegalFormModal`/`FormPage` flow (health disclosures, injury declarations, image-rights releases, etc.) does persist real submissions via `POST /consent` and `POST /forms/:slug/:formType` to the `consents`/`form_submissions` tables — this part is real.
- **Frontend:** `components/Chat.jsx`, `components/LegalFormModal.jsx`, `components/FormPage.jsx`. **Backend:** `POST /consent`, `GET/POST /forms/:slug/:formType`. **Storage:** `consents`, `form_submissions`. **External services:** none.
- **Real delivery or UI-only:** split — legal forms real; general AI-consent banner UI-only.
- **Failure states:** present for form submission.
- **Owner visibility:** `GET /owner/form-submissions` — yes, for legal forms; no visibility into who saw/accepted the general AI-consent banner (there's nothing to see — it's not stored).
- **Build-tested:** yes. **Live-tested:** unknown.
- **Perception risk:** HIGH — this is a direct, specific GDPR compliance claim ("Data processing compliant with UK GDPR / EU GDPR") sitting on top of a consent record that doesn't exist.
- **Status:** PARTIAL-REAL (legal forms REAL, AI-processing consent claim GHOST).
- **To make real:** persist the general AI-processing consent server-side (already flagged as a pre-launch trust requirement in `docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md` §5 and `docs/decisions/AI_REVIEW_LEDGER_TRUST_AND_CINEMATIC_SHELL_v0.2.md`).
- **Blocks launch:** YES (already flagged in the decision ledger as a pre-launch blocker, independent of this audit).

## 12. Booking requests
- **UI:** proposed inline by the AI in chat, captured via regex in `Chat.jsx`, managed in OwnerDashboard "Bookings" tab.
- **Perception:** structured booking pipeline with owner accept/decline/counter.
- **Actual:** real, end-to-end — `POST /booking-request`, `GET /owner/bookings`, `PATCH /owner/bookings/:id` all functional against a real table. The one soft spot: detection of "a booking happened" is a regex match on the AI's free-text reply (`replyLower.includes('sent your request')`), not a structured tool call — so it's only as reliable as the model's phrasing consistency, not deterministic.
- **Frontend:** `components/Chat.jsx`, `components/OwnerDashboard.jsx`. **Backend:** the 3 endpoints above. **Storage:** `booking_requests`. **External services:** none.
- **Real delivery or UI-only:** real.
- **Failure states:** present.
- **Owner visibility:** yes, dedicated tab with accept/decline/counter actions.
- **Build-tested:** yes. **Live-tested:** unknown.
- **Perception risk:** MEDIUM (silent misses possible if the AI phrases confirmation differently than expected).
- **Status:** REAL, with a known fragility.
- **To make real (fully robust):** replace regex detection with a structured tool-call/JSON contract from the model — a natural `ToolGate` candidate.
- **Blocks launch:** NO.

## 13. Leads
- **UI:** captured via `Chat.jsx`'s lead-capture form after 3 turns; scored server-side; visible in dashboard.
- **Perception:** CRM-lite lead tracking with hot-lead alerting.
- **Actual:** real. `scoreLead()` is a deterministic keyword scorer (not AI, but functions correctly and consistently), `sendHotLeadEmail` fires a real Resend email on hot-lead capture (same infra as the notification email, proven pattern). One limitation: score is computed once at lead-creation time and never re-evaluated as the conversation continues.
- **Frontend:** `components/Chat.jsx`, `components/OwnerDashboard.jsx`. **Backend:** `POST /lead`, `GET /owner/leads`. **Storage:** `leads`. **External services:** Resend.
- **Real delivery or UI-only:** real.
- **Failure states:** present (`fmt.Printf("Warning:...")`-logged, not surfaced to the visitor, which is correct for a background op).
- **Owner visibility:** yes.
- **Build-tested:** yes. **Live-tested:** unknown.
- **Perception risk:** LOW.
- **Status:** REAL.
- **To make real:** nothing required for launch; re-scoring is a nice-to-have.
- **Blocks launch:** NO.

## 14. Media upload
- **UI:** photo/media upload in Onboarding and `OwnerEditProfile.jsx`.
- **Perception:** upload photos to your profile.
- **Actual:** real — `POST /media/upload` streams to Supabase Storage, returns a public URL, used in both onboarding and profile-edit flows.
- **Frontend:** `components/Onboarding.jsx`, `components/OwnerEditProfile.jsx`. **Backend:** `POST /media/upload`. **Storage:** Supabase Storage `media` bucket. **External services:** Supabase Storage.
- **Real delivery or UI-only:** real.
- **Failure states:** present (size limit, upload-failure JSON errors).
- **Owner visibility:** direct (their own upload).
- **Build-tested:** yes. **Live-tested:** unknown — cannot confirm the `media` bucket exists/is configured correctly on the live Supabase project from this environment.
- **Perception risk:** LOW.
- **Status:** REAL (pending live-bucket confirmation).
- **To make real:** confirm the bucket exists and has correct public-read policy in the live Supabase project.
- **Blocks launch:** NO, assuming the bucket is already configured (it's referenced consistently across two independent upload flows, suggesting it's an established, working piece).

## 15. Payments / Stripe
- **UI:** OwnerDashboard settings (onboarding CTA), booking-deposit flow (implied by `BookingPayment` schema).
- **Perception:** real payment processing via Stripe Connect.
- **Actual:** `concierge-backend/stripe_integration.go` is a real, properly structured Stripe Connect implementation — onboarding link creation, account status check, checkout session creation, webhook with real signature verification (`verifyStripeSignature`), platform-fee calculation, payments listing. Gracefully degrades with an explicit error if `STRIPE_SECRET_KEY` isn't set, rather than silently pretending to work.
- **Frontend:** not deeply inspected in this pass (out of primary scope — Stripe frontend wiring wasn't read line-by-line). **Backend:** `POST /stripe/onboard`, `GET /stripe/status`, `POST /stripe/checkout`, `POST /stripe/webhook`, `GET /owner/payments`. **Storage:** `booking_payments`. **External services:** Stripe, Resend (for a payment-related email, `sendResendEmail`).
- **Real delivery or UI-only:** appears real at the code level.
- **Failure states:** explicit (`"STRIPE_SECRET_KEY not set"`).
- **Owner visibility:** `GET /owner/payments`.
- **Build-tested:** yes (compiles as part of the same Go binary). **Live-tested:** UNKNOWN — cannot verify `STRIPE_SECRET_KEY`/`STRIPE_WEBHOOK_SECRET` are set in the live Render environment, or that a real charge has ever been processed, from this environment.
- **Perception risk:** MEDIUM, purely due to unknown live-configuration status, not code quality.
- **Status:** UNKNOWN (code looks REAL; live status unverifiable from here).
- **To make real (confirm):** verify Stripe env vars are set in Render and run one real test charge end-to-end.
- **Blocks launch:** only if payments are part of the launch claim — flag for explicit confirmation before marketing "book and pay" as working.

## 16. Auth / session
- **UI:** `app/theconcierge/owner-auth/page.jsx`, `middleware.ts` gate.
- **Perception:** secure login.
- **Actual:** real, and — per the two hardening commits this week (`8fa6229`, `f0673a1`) — specifically verified: `HttpOnly` cookie, correct cookie name, anonymous `_directToken` injection blocked with 401, production login confirmed live with real `Set-Cookie` evidence (AUTH GREEN).
- **Frontend:** `app/theconcierge/owner-auth/page.jsx`, `middleware.ts`, `lib/auth.js`, `app/api/auth/*`. **Backend:** `POST /auth/signup`, `POST /auth/login`. **Storage:** `sessions`, `profiles` (password hash+salt). **External services:** none.
- **Real delivery or UI-only:** real.
- **Failure states:** present (401s, invalid-credentials messaging).
- **Owner visibility:** N/A (it's the gate itself).
- **Build-tested:** yes. **Live-tested:** YES — the one feature in the entire matrix with concrete, first-party live-test evidence (real `Set-Cookie` header, matched token, blocked bypass, production redirect behavior), per AUTH GREEN.
- **Perception risk:** LOW.
- **Status:** REAL.
- **To make real:** nothing required.
- **Blocks launch:** NO — this is the one piece already cleared.

---

# Top 10 "ghost/placeholder" risks

1. **Public profile chat bubble (`/[slug]`)** — fully scripted, no real AI call, dead-end links to a stale domain. GHOST. (#3)
2. **Hardcoded `concierge-ai-gamma.vercel.app` share links in OwnerDashboard** — real owners get the wrong URL to share. (#6)
3. **General AI-processing consent banner** — claims GDPR compliance, persists nothing server-side. (#11)
4. **Brave PA's claimed capabilities** (calendar read/write, web search, ticket booking) — told to the model as instructions, not backed by real integrations. (#7)
5. **"Owner notified" claim** — RESOLVED as of `412312e`/`d69daf0`, listed here only as a before/after reference point; no longer a ghost.
6. **"Human review" of flagged AI replies** — implies a pre-delivery safety gate; is actually a post-delivery client-side UI toggle. (#5)
7. **"Daily" news generation** — implies automation; is actually on-demand-only, no scheduler exists. (#9)
8. **Push notifications** — service worker exists and has a `push` listener, but is unregistered, unreachable, and has zero delivery infrastructure behind it. (#8, documented not faked)
9. **Notification digest/clustering** — no scheduler exists anywhere in the codebase; not implemented, correctly not faked. (#8)
10. **Onboarding handle-picker domain display** — shows the same stale domain as #2, at the moment a new owner is choosing their public identity. (#10)

# Top 10 "partial-real but promising" features

1. **Notifications (alert record + email)** — genuinely hardened this week; only push/digest remain, and those are honestly documented as absent rather than faked.
2. **Brave PA conversational core** — the actual chat/reasoning loop is real Anthropic; only the tool-claims are aspirational.
3. **Booking requests** — fully real pipeline; only the AI-reply regex detection is fragile.
4. **Owner dashboard (minus the domain bug)** — extensive, genuinely functional CRUD across leads/notes/forms/bookings.
5. **News/insights generation** — real generation, just mislabeled as "daily."
6. **Onboarding core save flow** — real profile/account creation; only the URL preview is wrong.
7. **Legal/consent forms (`LegalFormModal`/`FormPage`)** — fully real and persisted, unlike the general chat-consent banner.
8. **Chat sensitive-topic detection** — real, deterministic keyword matching; just needs the review-gate story clarified.
9. **Stripe integration (code-level)** — looks properly built; needs live-env confirmation to move from UNKNOWN to REAL.
10. **Media upload** — real, working pattern reused cleanly across two flows; only needs bucket-config confirmation.

# Features safe to market now
- Auth/session (login, dashboard gating) — proven live.
- Demo pages — proven real AI chat experience.
- Leads capture and hot-lead email alerts.
- Booking-request pipeline (owner accept/decline/counter).
- Legal/consent form collection (distinct from the general AI-consent banner).
- Owner notifications for sensitive-topic alerts (as of this week's fix).
- Onboarding (core account/profile creation — not the preview URL).

# Features that must not be marketed yet
- **"Chat with your AI concierge on your public page"** — until #3 is fixed, this is the product's central promise and it doesn't work for real users.
- **"Share your link"** — until the domain bug (#6/#10) is fixed, shared links point to a stale/wrong domain.
- **"GDPR-compliant AI processing consent"** — until general consent is persisted server-side (#11).
- **Any specific Brave PA capability claim** ("manages your calendar," "searches the web for you") — until real tool integrations exist (#7).
- **"Daily" news/insights** — until it's either scheduled or relabeled (#9).
- **Push notifications** in any form — doesn't exist yet (#8).
- **Payments/Stripe** — until live env vars and a real test charge are confirmed (#15).

# Pre-launch blockers
1. Public profile page chat is non-functional for real users (#3) — the core product claim doesn't work outside demos.
2. Hardcoded stale domain in share links across dashboard + onboarding (#6, #10) — breaks the primary distribution mechanism.
3. General AI-processing consent claim is unbacked (#11) — already flagged independently in `docs/decisions/AI_REVIEW_LEDGER_TRUST_AND_CINEMATIC_SHELL_v0.2.md` as a blocker; this audit confirms it from the code side.
4. Stripe live-configuration status unverified (#15) — must be confirmed either way before any payment claim is made.

# Recommended next 3 implementation packets
1. **"Real Chat on Real Pages"** — replace or supplement the static bubble in `PublicProfileClient.jsx`/`generateBusinessPage.js` with a genuine route to the real `<Chat/>` component for non-demo slugs, and fix every hardcoded `concierge-ai-gamma.vercel.app` reference (4 files) to use the actual current domain. These two fixes are tightly coupled — the dead-end chat links and the wrong-domain share links are the same underlying "stale domain / stale architecture" issue and should be fixed together.
2. **"Consent Truth Path"** — same pattern already applied to notifications: persist the general AI-processing consent banner server-side (extend `POST /consent` or add a lightweight equivalent), matching the trust-hardening precedent already set this week.
3. **"Capability Claims Audit for Brave PA"** — either scope down `buildBravePAPrompt()`'s stated capabilities to what's real today, or scope up the backend (starting with the smallest one: Google Calendar, since `Profile.GoogleRefreshToken` already exists in the schema, suggesting it was planned but never finished).
