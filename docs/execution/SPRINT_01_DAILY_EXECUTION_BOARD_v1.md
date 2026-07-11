# Sprint 01 — Daily Execution Board v1

**AI:** Claude Code
**Mode:** Repo Execution Mode / Execution Planning Documentation Only
**Date:** 2026-07-09
**Status:** Execution planning document. No code implemented, no code changed, no schema changed, no env vars touched. This is not a strategy document — it is the practical day-by-day board built on top of the strategy set already in `docs/strategy/` and `docs/audits/`. Open this document each morning; it should tell you exactly what today is for.

Every day below still requires Bruno's explicit go-ahead before Claude Code touches any code, per the existing Diamond Protocol (`docs/protocol/PRODUCT_OPERATING_SYSTEM_v0.5.md` §3). This board sequences the work; it does not pre-authorize it.

---

## Daily Status Template

Close out every day in this sprint by filling in this template with real evidence, not a summary — per Anti-Chaos Rule #16 ("Concrete evidence beats a confident summary, always"). A blank copy is appended to each day's card below; fill it in at end of day and leave it in place as the permanent record of what actually happened, don't delete it once complete.

```
**Task:**
**Status:** PASS / FAIL / PARTIAL / BLOCKED
**Files changed:**
**Commit:**
**Build/test:**
**Production proof:**
**Remaining risk:**
**Next safest action:**
```

- **Status** is never "done" or "looks good" — it is one of the four words above, chosen honestly. `PARTIAL` and `BLOCKED` are valid, useful outcomes, not failures to hide.
- **Production proof** means a real artifact — a `curl` result, a screenshot, a `Set-Cookie` header, a real DB row, a real push notification received — never a description of what should have happened.
- **Next safest action** is filled in even on a `PASS` day — it's the bridge to tomorrow's card, and on a `BLOCKED`/`FAIL` day it's the actual next step, not "investigate further."

---

## 1. Working rhythm

Each working day follows the same five-beat rhythm:

1. **Morning decision (Bruno).** Read today's card below. Confirm the day's single main packet is still the right priority (things may have shifted overnight — a live incident, a blocked dependency). Give Claude Code the explicit go-ahead for today's task, or adjust.
2. **ChatGPT packet (if needed).** For days that need a prepared Integration Packet, prompt, or QA checklist before Claude Code starts, ChatGPT produces it first. Not every day needs this — most days reuse the packet definitions already in `docs/strategy/IMPLEMENTATION_PACKETS_INDEX_v1.md`.
3. **Claude Code implementation.** Claude Code executes today's task only — one main packet per day, maximum. If today's task naturally splits (e.g., Days 7-8, 12-13), Claude Code still stops at a clean, verifiable checkpoint each day rather than running two days of work together.
4. **Bruno live verification.** Every day ends with a real, concrete check against production or a real test environment — never a summary alone (Anti-Chaos Rule #16: "Concrete evidence beats a confident summary, always"). Each day card below states exactly what to check.
5. **Evening commit/status.** Claude Code commits and pushes (documentation, or code once implementation phases begin — always with Bruno's explicit approval per commit, per standing instructions). Bruno notes the day's outcome before closing out — a one-line status is enough: done / partial / blocked, and why.

**One main packet per day, maximum.** Some days are decision-only or verification-only and involve no code — those still count as "today's one thing." Do not pull tomorrow's packet forward just because today finished early; use spare time for the live-test/verification step, not for starting the next packet early.

---

## 2. Assistant responsibility matrix

| Assistant | Role in this sprint |
|---|---|
| **Bruno** | Final decision-maker on every day's go-ahead; performs all live verification steps personally (per Anti-Chaos Rule #16 — Claude's own summary is not sufficient); makes the Day 6 (Calendar credentials) and Day 10 (Render upgrade) decisions specifically. |
| **ChatGPT** | Prepares Integration Packets, prompts, and QA checklists ahead of days that need them; turns this board's outcomes into the next sprint's planning input. |
| **Claude Code** | Repo executor — implements each day's packet, runs build/vet checks, commits and pushes only what's explicitly approved, reports real diffs and real test results, never a summary alone. |
| **Perplexity** | On call for the parallel legal track (privacy policy language research, DPIA best-practice research) and the parallel cost/provider track (comparing scheduler/weather-provider options) — not scheduled for a specific day unless a day card calls it out. |
| **Gemini** | On call for visual/UX review once Phase 5-6 (visual identity, Cinematic Shell) work approaches — not active in this 14-day sprint's scope, which stays in Roadmap Phases 1-4 and part of 7. |
| **Legal/GDPR consultant** | Owns the parallel legal track (§5 below) — privacy policy content review, Health-Data DPIA, voice legal scope resolution — runs alongside the engineering days, not blocking them except where explicitly noted. |

---

## 3-4. Day-by-day plan (First Operational Order)

### Day 1 — P2: Consent, Domain & Real Chat (start)
- **Main task:** fix the live `/privacy` 404 and persist AI-processing consent server-side.
- **Packet/task reference:** P2 / Roadmap Task 2.1 (Consent Truth Path).
- **Why this day:** the single most urgent item across every audit — a live-confirmed 404 on a GDPR disclosure, active right now, not a future risk.
- **Assistant responsible:** Claude Code (implementation), Bruno (go-ahead + verification).
- **Bruno decision/action needed:** approve serving `privacy.html` at a real route (recommend converting to `app/privacy/page.jsx`) and persisting `Chat.jsx`'s consent via the existing `POST /consent` pattern.
- **Technical constraint closed:** REG-09 (`GHOST_FORBIDDEN` → `ACTIVE_PUBLIC`) — privacy page reachable.
- **Legal/privacy constraint closed or deferred:** technical reachability closed today; full legal content review of the policy text is **deferred** to the parallel legal track (§5).
- **External service/account needed:** none.
- **Expected cost category:** £0.
- **Definition of done:** `/privacy` reachable at a real Next.js route; `giveConsent()` writes a real row to the `consents` table.
- **Live test required:** `curl -L https://www.bravebybruno.com/privacy` → 200; a real consent click in a real browser produces a real DB row Bruno can see.
- **What not to touch:** legal-form consent (`LegalFormModal.jsx`/`FormPage.jsx`) — already real, do not modify; auth; routing beyond the new privacy route.
- **Stop condition:** stop once both live tests pass. If the `POST /consent` reuse turns out to need a schema change beyond what's already reviewed, stop and report — do not improvise a new table today.
- **Daily status:**
```
**Task:** Fix live /privacy 404; persist AI-processing consent server-side; align privacy-policy text with actual data handling.
**Status:** PASS
**Files changed:**
  - app/privacy/page.jsx (new, 303 lines) — real Next.js route at /privacy
  - components/Chat.jsx (+39/-5) — giveConsent() now POSTs to BACKEND_URL/consent instead of sessionStorage-only
  - app/privacy/page.jsx (follow-up, +6/-7) — corrected policy text that had shipped saying "not stored" when handleChat (concierge-backend/main.go:191-216) already persists every message via db.SaveConversation/db.SaveMessage
**Commit:**
  - 52451f7 "fix: persist chat consent and serve privacy page" (privacy route + consent POST)
  - ec6e2c4 "fix: correct privacy policy to reflect actual message storage" (policy-text correction)
**Push result:** both pushed to origin/main — `git push` output: `52451f7..ec6e2c4  main -> main`. origin/main HEAD confirmed at ec6e2c4.
**Build/test:**
  - `npm run build` → " ✓ Compiled successfully", /privacy listed as a static (○) prerendered route in the output route table.
  - `go build ./...` (concierge-backend) → exits clean, no errors.
**Production proof:**
  - `curl -s -o /dev/null -w "HTTP %{http_code}" -L https://www.bravebybruno.com/privacy` → `HTTP 200`.
**Consent DB-row proof:**
  - Real POST to live backend: `curl -X POST https://concierge-backend-80rb.onrender.com/consent` with a disposable `session_id=qa-test-1783732146`, `profile_id=qa-test-profile` → `{"id":"927cc69a-2b69-4c02-9a3f-c5108accafaa","status":"ok"}`, HTTP 200.
  - Confirmed via direct Supabase REST read (`GET /rest/v1/consents?id=eq.927cc69a-...`) → row returned with matching session_id, forms_agreed=["ai_processing"], real created_at timestamp `2026-07-11T01:09:07.63Z`. This is a real row in the real `consents` table, not an inference from code.
**Remaining risk:**
  - The disposable QA test consent row (id `927cc69a-...`) could not be deleted for cleanup — anon/publishable Supabase key has no DELETE grant on `consents` (`42501 permission denied`, same known gap noted for the `profiles` table in prior auth-hardening work). Row contains no real PII (empty name/email, synthetic profile_id) so low sensitivity, but a working service-role key for admin cleanup is still missing locally — worth fixing before more QA rows accumulate.
  - Full legal content review of the privacy policy text (beyond the technical accuracy fix made today) is still deferred to the parallel legal track per this day's own scope — today closed technical reachability + factual accuracy, not a lawyer-reviewed policy.
**Next safest action:** Day 2 — mount real `<Chat/>` on real (non-demo) public profile pages and fix the hardcoded Vercel-alias share-link domain.
```

### Day 2 — Finish P2 + dashboard reliability if needed
- **Main task:** mount real `<Chat/>` on real public profile pages; fix the hardcoded Vercel-alias share-link domain; add dashboard error surfacing if time allows.
- **Packet/task reference:** P2 / Roadmap Tasks 2.2 (Real Chat on Real Pages), 2.3 (domain correction), 2.4 (dashboard reliability, stretch goal).
- **Why this day:** the product's core promise doesn't function outside `/demo/*` today — this is the highest-priority functional gap in the entire audit series, and shares its root cause with the domain fix.
- **Assistant responsible:** Claude Code (implementation), Bruno (go-ahead + verification).
- **Bruno decision/action needed:** confirm which real profile slug to use for the live test.
- **Technical constraint closed:** REG-03 (`GHOST_FORBIDDEN` → `ACTIVE_PUBLIC`), REG-10a (domain violation of Anti-Chaos Rule #17 fixed).
- **Legal/privacy constraint closed or deferred:** none new — depends on Day 1's consent fix already being live.
- **External service/account needed:** none — reuses existing Anthropic integration.
- **Expected cost category:** £0 (marginal increase in Anthropic API calls from real usage, not a new cost line).
- **Definition of done:** a real (non-demo) profile URL produces a real AI reply end-to-end; every share link in the dashboard/onboarding resolves to `bravebybruno.com`.
- **Live test required:** visit a real profile URL, send a message, confirm a genuine Anthropic-backed reply; copy a share link from a real dashboard session and confirm the domain.
- **What not to touch:** the demo chat flow (already real); the underlying profile/slug data.
- **Stop condition:** stop once the real-slug chat test and the share-link domain test both pass. Dashboard error-surfacing (2.4) is a stretch goal for this day only — if it doesn't fit, defer it, do not let it delay Day 1-2's core fix.
- **Daily status:**
```
**Task:**
**Status:** PASS / FAIL / PARTIAL / BLOCKED
**Files changed:**
**Commit:**
**Build/test:**
**Production proof:**
**Remaining risk:**
**Next safest action:**
```

### Day 3 — P1: Feature Flags & Audit Events (table + mechanism)
- **Main task:** build the `feature_flags` table and read/consume mechanism.
- **Packet/task reference:** P1 / Roadmap Task 1.1.
- **Why this day:** this is the enabling infrastructure every later "ship it dormant" day in this sprint (Days 7-8, 12-13) depends on — must exist before those days arrive.
- **Assistant responsible:** Claude Code (implementation), Bruno (schema review + go-ahead + verification).
- **Bruno decision/action needed:** **review and approve the exact `feature_flags` table DDL before it's run against Supabase** — this repo has no established migration pattern, so this is the first, and per standing practice, Claude Code must explain exactly what's needed rather than improvise.
- **Technical constraint closed:** REG-34 (`GHOST_FORBIDDEN` as a system → real infrastructure exists).
- **Legal/privacy constraint closed or deferred:** none.
- **External service/account needed:** none — building in-repo, not adopting a paid flag SaaS.
- **Expected cost category:** £0.
- **Definition of done:** toggling a test flag in the table visibly changes behavior in a real deployed environment.
- **Live test required:** toggle a flag, confirm the observable change.
- **What not to touch:** any existing table; any feature this system will later gate (the migration happens on later days, not today).
- **Stop condition:** stop once the table exists and one real toggle is proven live. Do not begin migrating existing features onto the flag system today — that's Day 4 and beyond.
- **Daily status:**
```
**Task:**
**Status:** PASS / FAIL / PARTIAL / BLOCKED
**Files changed:**
**Commit:**
**Build/test:**
**Production proof:**
**Remaining risk:**
**Next safest action:**
```

### Day 4 — Seed feature registry + audit events
- **Main task:** build the `audit_events` table and generalize this week's notification-status pattern into it; seed `feature_flags` from `docs/strategy/FEATURE_REGISTRY_AND_ACTIVATION_STATES_v1.md`.
- **Packet/task reference:** P1 / Roadmap Tasks 1.2 (audit events), 1.3 (seed registry).
- **Why this day:** natural follow-on to Day 3, makes the registry document operationally real rather than a static reference file.
- **Assistant responsible:** Claude Code (implementation), Bruno (go-ahead + verification).
- **Bruno decision/action needed:** approve the `audit_events` DDL (same schema-precedent discipline as Day 3).
- **Technical constraint closed:** REG-35 (`SPEC_ONLY` → real, generalized instance).
- **Legal/privacy constraint closed or deferred:** none.
- **External service/account needed:** none.
- **Expected cost category:** £0.
- **Definition of done:** every row in the Feature Registry document has a corresponding live row in `feature_flags`; two independent product claims (notification + at least one more, e.g., booking-request creation) both write through the same reusable `audit_events` function.
- **Live test required:** spot-check 3-5 registry rows against live table state; confirm both claim types produce real audit-event rows.
- **What not to touch:** the existing, already-proven notification code's behavior — wrap/extend, don't rewrite.
- **Stop condition:** stop once the registry seed and the two-claim-type audit-event proof both pass.
- **Daily status:**
```
**Task:**
**Status:** PASS / FAIL / PARTIAL / BLOCKED
**Files changed:**
**Commit:**
**Build/test:**
**Production proof:**
**Remaining risk:**
**Next safest action:**
```

### Day 5 — P4: Brave PA Capability Truth
- **Main task:** remove or gate the calendar/search/booking claims in `buildBravePAPrompt()` that have no backend behind them.
- **Packet/task reference:** P4 / Roadmap Task 3.1.
- **Why this day:** an active, ongoing misrepresentation — Brave PA is telling real users it can do things it cannot, today, right now.
- **Assistant responsible:** Claude Code (implementation), Bruno (go-ahead + verification).
- **Bruno decision/action needed:** confirm the replacement copy is acceptable (i.e., Brave PA should now clearly and honestly say what it cannot yet do, rather than silently going quiet on those topics).
- **Technical constraint closed:** REG-19a (`GHOST_FORBIDDEN` → removed or flagged pending Day 7-8's connector).
- **Legal/privacy constraint closed or deferred:** none.
- **External service/account needed:** none.
- **Expected cost category:** £0.
- **Definition of done:** every capability claim in the system prompt is either backed by real code or explicitly and honestly absent.
- **Live test required:** ask Brave PA to do each previously-claimed action and confirm the response now matches reality.
- **What not to touch:** Brave PA's real conversational core — do not degrade it while fixing the claims.
- **Stop condition:** stop once the live test passes for every previously-false claim.
- **Daily status:**
```
**Task:**
**Status:** PASS / FAIL / PARTIAL / BLOCKED
**Files changed:**
**Commit:**
**Build/test:**
**Production proof:**
**Remaining risk:**
**Next safest action:**
```

### Day 6 — Google Calendar OAuth setup decision
- **Main task:** decide whether the Google Cloud Console project/credentials from the previously-interrupted Blocco A3 work still exist, or need recreating.
- **Packet/task reference:** precursor to P4 / Roadmap Task 3.3.
- **Why this day:** the Master Doc excerpt states this work was "mid-implementation" — this decision determines whether Days 7-8 can proceed as planned or need to start further back.
- **Assistant responsible:** **Bruno** (this is a decision/lookup day, not an implementation day) — Claude Code can assist by searching for any existing local record of the credentials setup, but cannot access Google Cloud Console itself.
- **Bruno decision/action needed:** check Google Cloud Console directly for an existing project/credentials; decide whether to reuse or recreate.
- **Technical constraint closed:** none yet — this is the decision that unblocks Days 7-8.
- **Legal/privacy constraint closed or deferred:** none.
- **External service/account needed:** Google Cloud Console (existing or new project).
- **Expected cost category:** £0 (Google Calendar API free tier is expected to be sufficient at this scale).
- **Definition of done:** a clear yes/no on whether existing credentials exist, and the `GOOGLE_CLIENT_ID`/`GOOGLE_CLIENT_SECRET` values are in hand (not shared with Claude Code — only confirmation they exist and are ready to be set as Render env vars).
- **Live test required:** N/A — decision day.
- **What not to touch:** N/A.
- **Stop condition:** stop once the decision is made and documented (a one-line note is enough: "reusing existing project" or "creating new project, ETA X"). If credentials cannot be resolved today, Days 7-8 should be pushed later in the sprint rather than started without them.
- **Daily status:**
```
**Task:**
**Status:** PASS / FAIL / PARTIAL / BLOCKED
**Files changed:**
**Commit:**
**Build/test:**
**Production proof:**
**Remaining risk:**
**Next safest action:**
```

### Day 7-8 — Calendar Connector V1 (if credentials are ready)
- **Main task:** build the OAuth callback handler + read-only calendar-read endpoint.
- **Packet/task reference:** P4 / Roadmap Task 3.3.
- **Why this day:** the genuinely closest-to-done unfinished feature in the repo, per the Master Doc's own "last confirmed step" note — and directly resolves Day 5's remaining calendar-claim gap.
- **Assistant responsible:** Claude Code (implementation), Bruno (env var entry + go-ahead + verification).
- **Bruno decision/action needed:** enter `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_OAUTH_REDIRECT_URI` into Render's environment directly (Claude Code never handles secret values); approve shipping this as owner-opt-in/beta, not defaulted on.
- **Technical constraint closed:** REG-22 (`SPEC_ONLY` → `DORMANT_BUILT`, owner-opt-in).
- **Legal/privacy constraint closed or deferred:** standard OAuth consent screen closed today; full privacy-notice disclosure of Google as a processor should land alongside or shortly after (coordinate with the parallel legal track, §5).
- **External service/account needed:** Google Cloud Console project (per Day 6's decision), Google Calendar API.
- **Expected cost category:** £0/LOW (free tier expected sufficient).
- **Definition of done:** a real OAuth token successfully refreshes and a real event list is returned to Brave PA for an opted-in owner.
- **Live test required:** yes — a real "Connect Calendar" flow, a real token refresh, a real "what's on my calendar today" answer.
- **What not to touch:** the product's own login/auth system — this is a separate OAuth flow; existing profile data.
- **Stop condition:** if Day 6 did not resolve credentials, **do not proceed** — use these two days for a different, unblocked task instead (recommend pulling forward Day 11's notification-preferences work) and revisit Calendar once credentials are confirmed. If credentials are ready, stop once the live test passes.
- **Daily status:**
```
**Task:**
**Status:** PASS / FAIL / PARTIAL / BLOCKED
**Files changed:**
**Commit:**
**Build/test:**
**Production proof:**
**Remaining risk:**
**Next safest action:**
```

### Day 9 — P10: Live Config & Financial Verification
- **Main task:** confirm Stripe, Resend, and Supabase media-bucket live configuration; run one real Stripe test charge.
- **Packet/task reference:** P10 / Roadmap Task 7.1.
- **Why this day:** the highest financial and trust stakes in this entire sprint — real money via Stripe — must be confirmed before any launch claim, not discovered late.
- **Assistant responsible:** Claude Code (verification scripting/checks), Bruno (performs the real test charge, confirms results).
- **Bruno decision/action needed:** authorize and perform one real Stripe test transaction; confirm `RESEND_API_KEY`/`OWNER_EMAIL`/`RESEND_FROM_EMAIL` and `STRIPE_SECRET_KEY`/`STRIPE_WEBHOOK_SECRET` are genuinely set in Render (Claude Code checks for presence, never values).
- **Technical constraint closed:** REG-32 (`UNKNOWN` → `ACTIVE_PUBLIC`, if verification passes); REG-13's live-config confirmation (email).
- **Legal/privacy constraint closed or deferred:** confirm payment terms/refund policy text exists and is accurate — **flag to the legal track if not**, don't attempt to write it today.
- **External service/account needed:** Stripe, Resend, Supabase Storage (all already integrated).
- **Expected cost category:** £0 for the verification itself; **COST_VERIFY_REQUIRED** is resolved, not introduced, today.
- **Definition of done:** one real Stripe test charge succeeds with the webhook acknowledged before timeout; a real sensitive-topic alert returns `email_status: "sent"`; one real media upload succeeds through both existing flows.
- **Live test required:** yes — this entire day is the live test.
- **What not to touch:** Stripe integration code itself, unless verification reveals a real defect — if so, stop and scope that as its own separate task, do not silently patch it today.
- **Stop condition:** stop once all three verifications pass or fail with a clear, documented reason. A failure here is a valid, useful outcome — do not force a pass.
- **Daily status:**
```
**Task:**
**Status:** PASS / FAIL / PARTIAL / BLOCKED
**Files changed:**
**Commit:**
**Build/test:**
**Production proof:**
**Remaining risk:**
**Next safest action:**
```

### Day 10 — Render production decision
- **Main task:** decide on the Render plan upgrade.
- **Packet/task reference:** cross-cutting infrastructure decision, referenced throughout `docs/audits/TECHNICAL_CONSTRAINTS_MATRIX.md`.
- **Why this day:** the repo's own doc states this upgrade "must happen before any real client is sent a link" — and Day 9 just demonstrated why (Stripe webhook timing, cold-start risk); this is the natural decision point right after that evidence lands.
- **Assistant responsible:** **Bruno** (billing decision) — Claude Code can summarize the cold-start risk evidence from Day 9 but does not make spending decisions.
- **Bruno decision/action needed:** decide whether to upgrade the Render plan now (recommended, given Day 9's evidence and the repo's own pricing rule) or defer, and check whether the current/target plan includes Cron Job support (needed for Days 12-13 and Roadmap Phase 4's digest work later).
- **Technical constraint closed:** resolves the free-tier cold-start risk flagged repeatedly across the audits.
- **Legal/privacy constraint closed or deferred:** none.
- **External service/account needed:** Render (existing account, plan change only).
- **Expected cost category:** **COST_VERIFY_REQUIRED** — exact plan pricing not available to Claude Code from this environment.
- **Definition of done:** a documented decision (upgrade now / defer with a named reason) and, if upgrading, confirmation of Cron Job availability on the new plan.
- **Live test required:** N/A — decision day, though if upgraded, worth re-running Day 9's Stripe webhook timing check post-upgrade.
- **What not to touch:** N/A.
- **Stop condition:** stop once the decision is documented, regardless of which way it goes.
- **Daily status:**
```
**Task:**
**Status:** PASS / FAIL / PARTIAL / BLOCKED
**Files changed:**
**Commit:**
**Build/test:**
**Production proof:**
**Remaining risk:**
**Next safest action:**
```

### Day 11 — Notification preferences + sound
- **Main task:** persist `notifPrefs` server-side and make them actually gate behavior; wire a real sound to real alert arrival.
- **Packet/task reference:** P5 / Roadmap Tasks 4.2 (preferences), 4.3 (sound).
- **Why this day:** cheap, contained, closes two concrete ghost-UI gaps flagged in every audit; natural to do together since sound depends on the preference existing.
- **Assistant responsible:** Claude Code (implementation), Bruno (go-ahead + verification).
- **Bruno decision/action needed:** confirm whether preferences reuse `profiles.profile_data` (no migration) or need real columns — recommend the blob reuse, lower risk.
- **Technical constraint closed:** REG-15 and REG-16 (`GHOST_FORBIDDEN` → `ACTIVE_PUBLIC`).
- **Legal/privacy constraint closed or deferred:** none.
- **External service/account needed:** none for preferences; an audio asset file needs a `public/` folder (may already exist from earlier sprint work if Days 7-8 or elsewhere required it — check before recreating).
- **Expected cost category:** £0.
- **Definition of done:** toggling a preference off suppresses the corresponding behavior; toggling it on enables it; a real alert plays a real sound when enabled.
- **Live test required:** yes — toggle off, confirm silence; toggle on, confirm sound plays on a real alert.
- **What not to touch:** the underlying alert-creation logic — this only gates display/delivery, not detection.
- **Stop condition:** stop once both toggles are proven live.
- **Daily status:**
```
**Task:**
**Status:** PASS / FAIL / PARTIAL / BLOCKED
**Files changed:**
**Commit:**
**Build/test:**
**Production proof:**
**Remaining risk:**
**Next safest action:**
```

### Day 12-13 — Push notification V1
- **Main task:** register the existing orphaned `service_worker.js`, build a subscribe flow, wire real Web Push delivery.
- **Packet/task reference:** P5 / Roadmap Task 4.4.
- **Why this day:** real infrastructure work, appropriately sequenced after the cheaper wins (Day 11) and after Days 3-4's flag system exists to ship this dormant/opt-in.
- **Assistant responsible:** Claude Code (implementation), Bruno (VAPID key generation/entry + go-ahead + verification on a real device).
- **Bruno decision/action needed:** generate and enter `VAPID_PUBLIC_KEY`, `VAPID_PRIVATE_KEY`, `VAPID_SUBJECT` into Render directly; decide on the subscription-storage approach (reuse `notes` pattern vs. dedicated table — recommend `notes` reuse to avoid another new-table decision this sprint); approve shipping as owner-opt-in only.
- **Technical constraint closed:** REG-14 (`GHOST_FORBIDDEN` → `DORMANT_BUILT`, owner-opt-in).
- **Legal/privacy constraint closed or deferred:** update privacy notice to disclose push-subscription data — coordinate with the legal track, don't block the technical build on it.
- **External service/account needed:** none required (Web Push needs no paid provider) — only the self-generated VAPID keys.
- **Expected cost category:** £0.
- **Definition of done:** a real push notification is received on a real device for a real, opted-in test subscription.
- **Live test required:** yes, on a real browser/device — this cannot be verified by build/vet alone.
- **What not to touch:** the existing alert-record/email path — push is additive, not a replacement.
- **Stop condition:** if no `public/` folder exists yet by this point in the sprint, creating it is part of Day 12's work, not a blocker to defer past. Stop once one real push notification is confirmed received.
- **Daily status:**
```
**Task:**
**Status:** PASS / FAIL / PARTIAL / BLOCKED
**Files changed:**
**Commit:**
**Build/test:**
**Production proof:**
**Remaining risk:**
**Next safest action:**
```

### Day 14 — Core QA mini-gate
- **Main task:** a lightweight, sprint-scoped version of the full Pre-Launch QA Day (Roadmap Task 7.4), covering only what this sprint actually built.
- **Packet/task reference:** scoped subset of P13 / Roadmap Task 7.4.
- **Why this day:** closes the sprint with a real, evidence-based checkpoint rather than assuming everything landed correctly.
- **Assistant responsible:** Bruno (performs verification), Claude Code (assembles the checklist and runs build/vet one final time), ChatGPT (helps turn today's outcome into Sprint 02's planning input).
- **Bruno decision/action needed:** walk through every live test listed on Days 1-13 one more time, end to end, in a single sitting, and confirm nothing has regressed.
- **Technical constraint closed:** confirms Weeks 1-2's milestone (see §9) is genuinely met, not just individually-tested-and-forgotten.
- **Legal/privacy constraint closed or deferred:** confirm the parallel legal track's status (§5) — note what's closed vs. still in progress, do not treat silence as completion.
- **External service/account needed:** none new.
- **Expected cost category:** £0.
- **Definition of done:** every Day 1-13 live test re-passes in one continuous session; a documented list of what's still open going into Sprint 02.
- **Live test required:** yes — this entire day is a live test session.
- **What not to touch:** nothing new — this is verification only, no new building.
- **Stop condition:** stop at the end of the day regardless of outcome; document results honestly, including any regressions found, rather than extending the day to "fix and re-test" — that becomes Sprint 02's first task if needed.
- **Daily status:**
```
**Task:**
**Status:** PASS / FAIL / PARTIAL / BLOCKED
**Files changed:**
**Commit:**
**Build/test:**
**Production proof:**
**Remaining risk:**
**Next safest action:**
```

---

## 5. Parallel legal track

Runs alongside the engineering days above, owned by Bruno + the Legal/GDPR consultant, not blocking most engineering days except where explicitly noted (Day 1's technical fix does not wait for this track; the *content* review does).

- **Privacy policy review** — real legal review of `privacy.html`'s content, once it's technically reachable (Day 1). Target: complete within this sprint, does not block any engineering day.
- **Health-data DPIA scoping** — the narrowest, most achievable legal deliverable identified across all audits (Roadmap Task 7.2). Should be commissioned this sprint; does not block any of the 14 engineering days above, since it concerns already-live legal-form data, not new functionality.
- **Voice legal scope resolution** — Roadmap Task 10.1, determining whether transient voice I/O is separable from the full session-learning block. Not scheduled in this 14-day sprint's engineering days (Phase 10 is late-roadmap) — but can be commissioned in parallel now if Bruno wants it resolved before it becomes a blocker later.
- **Voice/video DPIA** — Roadmap Task 10.2, explicitly later, depends on Task 10.1 completing first. Not a Sprint 01 deliverable.
- **Sensitive / Coach / Body / Health remains `LEGAL_LOCKED`** for the entire sprint — no engineering day above touches this, and none should.
- **Voice/Video Session Learning remains `LEGAL_LOCKED`** for the entire sprint — same.

---

## 6. Parallel cost/provider track

| Decision | Owner | Sprint day it's needed by | Status target by end of sprint |
|---|---|---|---|
| Render upgrade decision | Bruno | Day 10 (but informed by Day 9's evidence) | Decided either way |
| Resend live env verification | Bruno + Claude Code | Day 9 | Confirmed |
| Stripe live test | Bruno | Day 9 | Confirmed |
| Supabase media bucket verification | Bruno + Claude Code | Day 9 | Confirmed |
| Google Calendar OAuth setup | Bruno | Day 6 | Decided (reuse or recreate) |
| Weather provider decision | Bruno | Not needed this sprint (Roadmap Phase 5) | No action needed — defer |
| Scheduler/cron decision | Bruno | Informed by Day 10's Render plan check | Can defer full decision to Sprint 02 (Roadmap Phase 4's digest work, Task 4.5) unless Day 12-13's push work reveals an earlier need |

---

## 7. Assistant maturity track

Tracked against `docs/strategy/ASSISTANT_VERSION_LADDER_AND_KNOWLEDGE_PACKS_v1.md`:

- **V0 (Scripted Assistant) must be removed/replaced** — closes on Day 2.
- **V1 (Real Chat Assistant) becomes real** on real pages — closes on Day 2, alongside V0's removal (they are the same change, viewed from two sides).
- **V1.5 (Owner PA Truthful)** — closes on Day 5.
- **V2 (Connected Assistant) begins with Calendar** — starts Days 7-8, reaches `DORMANT_BUILT` (owner-opt-in beta) by end of sprint if credentials are ready; otherwise carries into Sprint 02.
- **V3 (Operating Assistant)** — not reached this sprint; depends on Roadmap Phase 8 (more connectors) and Phase 9 (`ActionGate`), both beyond this sprint's scope.
- **V4 (Memory/Personalization)** — not reached this sprint; Roadmap Phase 9, after Cinematic Shell and Core Flow QA.
- **V5/V6 remain legal-gated** — no engineering movement this sprint, consistent with §5.

---

## 8. Data participation track

Tracked against `docs/strategy/DATA_PARTICIPATION_AND_COHORT_LEARNING_v1.md`:

- **Essential / No Memory** remains the true default throughout — reinforced by Day 5's Brave PA truth-pass (no mode should ever be implied that isn't active).
- **Session Only** is already real — no sprint action needed, just confirmed as unchanged.
- **Personal Memory is not active this sprint** — no engineering day above builds toward it; it remains `GHOST_FORBIDDEN` (claim) until Roadmap Phase 9.
- **Connected Assistant** — begins per-connector with Calendar (Days 7-8), exactly as scoped.
- **Cohort Learning is not build-now** — this sprint does not schedule `docs/strategy/COHORT_LEARNING_AND_ANONYMOUS_PRODUCT_INTELLIGENCE_v1.md`'s Minimum Safe V1. It depends on: the §14 legal-basis determination (not commissioned this sprint — flag for Sprint 02 planning if Bruno wants to start it), the feature-flag system (built Day 3, so technically available after), audit events (built Day 4, available after), and a scheduler decision (deferred per §6 above). All four prerequisites should exist by Sprint 02, but this sprint does not build the cohort pipeline itself.
- **Sensitive/Coach/Voice** remain `LEGAL_LOCKED` — no sprint action, consistent with §5.

---

## 9. Weekly milestones

- **Week 1 (Days 1-7):** public core truth fixed — real chat on real pages, consent/privacy reachable, share-links correct, feature-flag/audit-event infrastructure live, Brave PA's false claims removed, Calendar connector underway.
- **Week 2 (Days 8-14):** flags fully seeded, Brave PA truth complete, Calendar connector proven (if credentials allowed), financial/live-config verification complete, Render decision made, notification maturity (preferences, sound, push) delivered, sprint closes on a real QA checkpoint.
- **Week 3 (beyond this sprint):** visual identity spec acquisition (Trust Dot, Brave Star, Location-Ambient — Roadmap Phase 5) and Cinematic Shell Integration Packet preparation (Roadmap Phase 6, Task 6.1) — not started in Sprint 01, this is the natural Sprint 02 focus once Phase 2-4 truth-and-infrastructure work is proven stable.
- **Week 4 (beyond this sprint):** full Core Flow QA and Launch Readiness (Roadmap Phase 7, Task 7.4) at full scope, once visual/Cinematic Shell work (Week 3) has landed — Day 14 of this sprint is a mini-gate covering only Sprint 01's own scope, not the full launch gate.

---

## 10. Launch-safe definition

**What must pass before sharing real client links** (i.e., before Bruno sends `bravebybruno.com/<real-slug>` to an actual prospective client, not a test):

- Day 1-2's real-chat and consent/privacy fixes are live and confirmed (REG-03, REG-08, REG-09 all out of `GHOST_FORBIDDEN`).
- Day 9's Stripe/Resend/media verification has passed, or Stripe-dependent features are explicitly not mentioned in what's shared with that client.
- Day 10's Render decision has been made — per the repo's own binding pricing rule, the free-tier backend should not be what a real client's traffic hits.
- Day 5's Brave PA truth-pass is live, if Brave PA is part of what's shown to that client.

**What can remain dormant** (i.e., does not need to be finished before sharing real links, as long as it isn't claimed):
- Calendar connector (Days 7-8) — fine to ship after real links go out, as long as no false claim exists in the meantime (Day 5 already ensures this).
- Push notifications (Days 12-13) — fine to remain `DORMANT_BUILT`/owner-opt-in, not required for a client-facing launch.
- Notification preferences/sound (Day 11) — cosmetic-adjacent, not launch-blocking either way.
- Everything in Roadmap Phases 5, 6, 8, 9, 10 — none of this sprint's scope, none required for a launch-safe state.

**What must not be marketed**, regardless of sprint progress:
- Any voice claim, in any framing (§5, §7, §8 — consistently `LEGAL_LOCKED`).
- "Personalized" / "remembers you" (Personal Memory mode is not active this sprint).
- Any Brave PA capability beyond what Day 5 confirms is real, and beyond Calendar once Days 7-8 land (search, ticket-booking, etc. remain unbuilt and unclaimed).
- "Daily" automated insights (Roadmap Task T, not in this sprint's scope — the on-demand mechanism remains, the automation does not exist yet).
- Any cohort-learning or "we learn from usage" product claim — explicitly not built this sprint (§8).

Stop.
