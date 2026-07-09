# Brain Spine Readiness Audit

**AI:** Claude Code
**Mode:** Repo Execution Mode / Documentation + Audit Only
**Date:** 2026-07-09
**Status:** Audit only. No implementation. No behaviour changed.

Context: AUTH GREEN accepted (evidence: pushed commits `8fa6229`, `f0673a1`; real `Set-Cookie: cai_token` with `HttpOnly`; token match; anonymous `_directToken` blocked with 401; production live on `bravebybruno.com`). This audit evaluates whether to introduce a Minimal Brain Spine before Cinematic Shell, per `docs/architecture/BRAIN_SPINE_FLEET_INTEGRATION_DECISION_PACK_v0.2.md`. The full Fleet Brain (`docs/architecture/FLEET_BRAIN_ARCHITECTURE_BRIEF.md`) is a future target, not an active task.

---

## 1. Documents saved / moved

Three uploaded planning documents were located (they arrived via GitHub web-upload commits to `origin/main`, not the Codespace filesystem — found via `git fetch` after an initial on-disk search came up empty) and moved into structured documentation folders:

| Original (repo root) | Saved path |
|---|---|
| `BraveByBruno_Product_Operating_System_AI_Technical_Round_Table_Protocol_v0_5.md` | `docs/protocol/PRODUCT_OPERATING_SYSTEM_v0.5.md` |
| `05cfdf97.md` (identified by content as "The Concierge Fleet Architecture Brief") | `docs/architecture/FLEET_BRAIN_ARCHITECTURE_BRIEF.md` |
| `BraveByBruno_Brain_Spine_Fleet_Integration_Decision_Pack_v0.2.md` | `docs/architecture/BRAIN_SPINE_FLEET_INTEGRATION_DECISION_PACK_v0.2.md` |

`docs/protocol/` and `docs/architecture/` were created since neither existed.

## 2. Documentation commit hash

`b9fc477` — "Organize uploaded strategy docs into docs/protocol and docs/architecture" (pure renames, 3 files, 0 insertions/deletions). **Local only at time of that commit — push status was not re-verified in this pass; confirm with `git status`/`git log origin/main..HEAD` before assuming it is on the remote.**

---

## 3. Current AI-related repo map

| Path | Purpose | Layer | Status |
|---|---|---|---|
| `components/Chat.jsx` | Public Concierge chat widget (consent gate, alert keywords, review modal, lead capture, booking-intent parsing) | Frontend | Active |
| `components/BravePAv2.jsx` | Owner-facing Brave PA assistant UI | Frontend | Active |
| `components/BravePASettings.jsx` | Brave PA config (name/personality) — no direct AI call | Frontend | Active |
| `components/OwnerDashboard.jsx` | Leads/conversations/bookings/notes dashboard; also triggers news generation | Frontend | Active |
| `components/AdminDashboard.jsx` | Bruno-only stats incl. estimated Claude cost (display only) | Frontend | Active |
| `components/Onboarding.jsx` | Onboarding wizard; calls vision AI to import services from a screenshot | Frontend | Active |
| `components/LegalFormModal.jsx` | Legal/consent forms; posts to `/consent` | Frontend | Active |
| `components/FormPage.jsx`, `BusinessPagePreview.jsx`, `OwnerEditProfile.jsx` | Supporting UI, no AI | Frontend | Active |
| `lib/buildPrompt.js` | `buildPrompt()` (Concierge system prompt), `buildBravePAPrompt()`, `generateProactiveMessage()`, `getQuickActions()` — the only prompt-construction logic in the repo | Shared | Active |
| `lib/constants.js` | `BACKEND_URL`, `GUARDRAILS` text, `MAX_TURNS`, 5 hardcoded `DEMOS` system prompts | Shared | Active |
| `lib/auth.js` | httpOnly cookie helpers (`cai_token`/`cai_slug`) | Backend (Next route handlers) | Active |
| `lib/generateBusinessPage.js` | Static public-page generation — no AI | Frontend | Active |
| `middleware.ts` | Gates `/theconcierge/dashboard` on `cai_token` cookie | Backend (edge) | Active |
| `app/api/auth/login/route.js`, `signup`, `logout`, `session` | Proxy to Go backend auth, set httpOnly cookies | Backend | Active |
| `api/chat.js` | Vercel serverless fn: rate limiting, turn cap, in-memory analytics, calls Anthropic (`claude-haiku-4-5`) directly | Backend | **Likely unused** — no frontend calls this path |
| `brave_pa_v2.js` (repo root) | Standalone legacy Brave PA prototype; calls `api.anthropic.com` **directly from the browser** | Frontend (orphaned) | **Unused** — not imported/scripted anywhere |
| `concierge-backend/main.go` | Go/Gin monolith: all 30+ routes, `callAnthropic()`, `callAnthropicVision()`, deterministic `scoreLead()` (keyword-based, not AI) | Backend | Active — this is the real runtime |
| `concierge-backend/db/supabase.go` | Data access layer over Supabase REST (all tables) | Backend | Active |
| `concierge-backend/stripe_integration.go` | Payments, not AI-related | Backend | Active |
| `bruno.html`, `marco.html`, `nour.html`, `sofia.html`, `alex.html`, `index.html`, `dashboard.html` (repo root) | Pre-Next.js static prototypes, duplicate the `DEMOS` content | Frontend | **Uncertain** — not confirmed still served/linked; flagged as repo-hygiene risk, not traced further |

---

## 4. Current AI call locations

| Location | Function/endpoint | Provider/model | Input | Output | User-facing | Risk |
|---|---|---|---|---|---|---|
| `concierge-backend/main.go` `callAnthropic()` | `POST /chat` → `handleChat` | Anthropic `claude-sonnet-4-6` | `messages[]`, **client-supplied** `system_prompt` | reply text, saved to `messages` table, returned as JSON | **Yes** — live path for Concierge chat, Brave PA, *and* dashboard news generation (all reuse this one endpoint) | **Medium-High** — public, unauthenticated for the Concierge-chat case; server trusts a client-supplied system prompt with no server-side rebuild/validation (prompt-injection/tampering surface); no rate limiting in the Go backend |
| `concierge-backend/main.go` `callAnthropicVision()` | `POST /ai/import-services` → `handleImportServices` | Anthropic `claude-haiku-4-5-20251001` | base64 image + fixed extraction prompt | parsed JSON services array | Owner-only (auth required), non-conversational | **Low** |
| `api/chat.js` | Vercel serverless `/api/chat` | Anthropic `claude-haiku-4-5-20251001` | `messages[]`, `systemPrompt`, `turnCount`, `profileId` | raw Anthropic response | **Not currently reachable** — `components/Chat.jsx` calls `BACKEND_URL/chat` (Go backend), not this route | **Medium** — dead code holding its own `ANTHROPIC_API_KEY`; diverges from the live path (different model, different max_tokens, its own analytics/rate-limit logic) |
| `brave_pa_v2.js` (root) | Inline `fetch('https://api.anthropic.com/v1/messages')` | Anthropic `claude-sonnet-4-6` | conversation messages, system prompt, `web_search` tool enabled | reply text | **Orphaned — not wired into any page**, and the fetch has no `x-api-key` header, so it would fail even if loaded | **High if ever reactivated** (client-side secret-exposure pattern); **currently inert** |

---

## 5. Current chat / Concierge / Brave PA flow map

**Concierge (public chat):**
- Entry: `/[slug]` (`app/[slug]/page.jsx` → `PublicProfileClient.jsx`) or `/demo/[id]`, both render `<Chat/>` with `profile`, `systemPrompt`, `profileId`.
- Frontend: `components/Chat.jsx` — consent screen (sessionStorage only, see gap below) → `callAPI()` → `fetch(BACKEND_URL + '/chat')`.
- Backend route: `POST /chat` on the Go backend (Render) → `handleChat`.
- Model call: `callAnthropic()`, system prompt built **client-side** via `lib/buildPrompt.js` and sent as plain text in the request body.
- Storage/logging: `conversations` + `messages` rows saved to Supabase. No moderation step server-side.
- Final response path: JSON `{reply, conversation_id, score}` → displayed in `Chat.jsx` after client-side currency conversion.
- Human review/moderation step: **client-side only.** `Chat.jsx`'s `alertMap` keyword-matches the *user's* message; if matched, the already-generated AI reply is routed into a `ReviewModal` (approve/edit/block) before being added to the visible message list. This is a UX gate in the browser, not a backend workflow — the reply already left the server in the fetch response regardless of the review outcome.
- Session/auth dependency: none — anonymous visitor, `sessionStorage`-based `cai_session` id only.
- **Gap found:** the "AI processing" consent banner in `Chat.jsx` (`giveConsent()`) only writes to `sessionStorage` — it never calls the backend `/consent` endpoint. Only the separate legal-form flow (`LegalFormModal.jsx` → `POST /consent`) is persisted server-side. So the GDPR-relevant "user agreed to AI processing" claim currently has no server-side record.

**Brave PA (owner-facing):**
- Entry: `components/BravePAv2.jsx`, mounted in the owner dashboard shell (distinct from the orphaned root `brave_pa_v2.js`).
- Model call: same `POST /chat` on the Go backend, system prompt built via `buildBravePAPrompt()` (also in `lib/buildPrompt.js`), no `web_search` tool wired server-side (the tool-enabled version only exists in the orphaned root file).
- Storage: same `conversations`/`messages` tables — Brave PA sessions are not structurally distinguished from Concierge client sessions except by `profile_id`/`session_id` string convention.
- Human review/moderation: none — Brave PA talks directly to the owner, no review gate applies.

**News/insights (owner dashboard):**
- `OwnerDashboard.jsx`'s `generateNewsNow()` also calls `POST /chat` (the same endpoint), with a market-intelligence system prompt, then separately `POST`s the parsed result to `/owner/news` for storage. UNCERTAIN whether this reuse was deliberate architecture or coincidental convenience.

**UNCERTAIN:** the static HTML pages (`bruno.html` etc.) — not confirmed whether they're still linked/served or fully superseded by the Next.js app; evidence missing is Vercel routing/rewrites and any inbound links to them.

---

## 6. Booking intent / alerts / human review flow map

| Flow | Starts | Stores data | Who sees it | Uses AI? | Would benefit from Brain Spine later? |
|---|---|---|---|---|---|
| Booking intent detection | `Chat.jsx` regex-scans the AI reply text for phrases like "sent your request" | `POST /booking-request` → `booking_requests` table | Owner, via `OwnerDashboard.jsx` | No — pure string matching on the AI's free-text output, not a structured tool call | Yes — a `ToolGate`/structured-output contract would replace this |
| Lead scoring | `scoreLead()` in `main.go`, runs on every `/chat` call | `leads` table (via separate `POST /lead` call from `Chat.jsx`) | Owner dashboard | No — deterministic keyword scoring | Mild — could stay deterministic |
| Sensitive-topic alerts | `Chat.jsx`'s `alertMap`, client-side keyword match on user message | Not stored server-side — UI-only banner | End-user sees "🔔 {name} notified"; **no actual notification is sent to the owner** | No | Yes — currently a UX illusion; needs a real pipeline |
| Human review (approve/edit/block) | `Chat.jsx` `ReviewModal`, triggered by the same `alertMap` | Not stored — client-side state only | End-user only (owner never sees blocked/edited replies) | No | Yes — clearest first candidate to wrap |
| Hot-lead email | `sendHotLeadEmail()` in `main.go`, fires when `handleLead` receives `score:"hot"` | N/A (transactional email via Resend) | Owner's email | No | No — already server-side and reliable |
| News/insights | `OwnerDashboard.jsx` `generateNewsNow()` | `daily_news` table | Owner dashboard | Yes — reuses `/chat` | Yes — should get its own `ModelAdapter` call shape distinct from conversational chat |

---

## 7. Current memory / profile / session / conversation state

| Data | Table/file | Access path | Structured enough for future memory retrieval? | Gaps/risks |
|---|---|---|---|---|
| User (owner) profile | `profiles` table, `profile_data` column is a JSON blob | `db.GetProfile`/`UpdateProfile` via `/owner/profile` | Partially — top-level fields structured; services/tone/sensitive topics/legal forms live in one opaque `profile_data` JSON string | No schema/typing on `profile_data`; future memory layer needs a parser, not a query |
| Conversation history | `conversations` + `messages` tables | `db.SaveConversation`/`SaveMessage`, read via `/owner/conversations` | Yes, structurally — role/content/timestamp per message | No retention/expiry policy visible; no distinction between Concierge-client and Brave-PA-owner conversations beyond `profile_id` |
| Lead/contact data | `leads` table | `/owner/leads`, `/lead` | Yes | Score is a static string set once at creation, never re-scored |
| Session/auth state | `sessions` table (token→slug), `cai_token`/`cai_slug` httpOnly cookies | `middleware.ts`, `lib/auth.js`, `authenticateToken()` in Go | Yes | Token expiry/rotation not confirmed either way — not this audit's scope |
| Consent/GDPR state | `consents` table (legal forms only) | `LegalFormModal.jsx` → `/consent` | Partially | See gap in §5 — general "AI processing" consent is not in this table, only sessionStorage |
| Business page data | `profiles.profile_data` (same blob) | `lib/generateBusinessPage.js` reads it | Same as profile — unstructured blob | Same as above |
| Dashboard state | Mostly ephemeral React state in `OwnerDashboard.jsx`; persisted pieces (`booking_prefs`, `notes`, `daily_news`) are proper tables | Various `/owner/*` endpoints | Yes for the persisted pieces | None found |

---

## 8. Existing abstraction layer check

**Reusable pieces that already exist:**
- `callAnthropic()` / `callAnthropicVision()` in `main.go` — a de facto (but not formalized) model-adapter, already the single choke point for every real Anthropic call.
- `lib/buildPrompt.js` — a de facto prompt-builder layer, already separated from UI and from the network call.
- `db/supabase.go` — a clean data-access layer, already separated from route handlers.

**Missing pieces:**
- No request envelope — `handleChat` takes raw `messages[]` + a client-trusted `system_prompt` string, nothing carrying user/session/risk/permissions as first-class fields.
- No verification/audit interface — nothing inspects a reply before it's returned; "review" is UI-only.
- No tool/action gate — booking intent is regex-sniffed off free text, not a structured action.
- No memory/provenance interface — conversation history is stored but never re-read into a later prompt; each `/chat` call is stateless from the model's perspective beyond whatever `messages[]` the client happens to resend.
- No audit log — no request ID, no record of which model/version produced which reply.

**Duplication found:**
- Two independent Anthropic call sites with different models/token limits (`api/chat.js` at `claude-haiku-4-5` vs `main.go` at `claude-sonnet-4-6`), one of them dead.
- Two independent Brave-PA prompt-building/call implementations (`components/BravePAv2.jsx`+`lib/buildPrompt.js`, live, vs the orphaned root `brave_pa_v2.js`, which duplicates `PA_PERSONALITIES` and `buildBravePAPrompt` verbatim inline instead of importing them).

**Where abstraction would reduce future rebuild:** most value is in wrapping `handleChat` — it's the single busiest, most-reused code path (Concierge, Brave PA, and news generation all funnel through it today with no differentiation), so it's both the highest-leverage and lowest-risk place to introduce a spine.

---

## 9. Minimal Brain Spine proposed module structure

All paths are proposals only — nothing created.

| Module | Purpose | Proposed path | Wraps | Changes behaviour? | Needs schema changes? |
|---|---|---|---|---|---|
| Request envelope | Carries user/session/profile/mode/risk alongside the existing `messages[]`/`system_prompt` | `concierge-backend/ai/envelope.go` | The struct `handleChat` already binds (`ChatRequest`) | NO | NO |
| Model adapter interface | Formalizes `callAnthropic`/`callAnthropicVision` behind a provider-agnostic signature | `concierge-backend/ai/adapter.go` | Existing `callAnthropic()`/`callAnthropicVision()` (body moves, logic unchanged) | NO | NO |
| Current provider adapter | The concrete Anthropic implementation of the adapter interface | `concierge-backend/ai/adapter_anthropic.go` | Same HTTP call, same models, same env vars | NO | NO |
| Memory/provenance interface | A read-only accessor over existing `conversations`/`messages`/`profiles` tables, so future features stop hand-rolling queries | `concierge-backend/ai/memory.go` | `db.GetConversationsByProfile` etc. (already exist) | NO | NO |
| Verification/audit interface | A pass-through no-op today; the seam future Trust Dot/verification work will implement against | `concierge-backend/ai/verify.go` | Nothing yet — new seam, always returns "pass" | NO | NO |
| Action/tool gate interface | Replaces the current regex-sniffing of booking intent with an explicit, still-deterministic check | `concierge-backend/ai/actiongate.go` | The booking-intent regex currently inline in `Chat.jsx` — could move server-side or stay client-side initially | NO (if kept as same regex, just relocated/named) | NO |
| Response synthesis interface | Ensures persona/formatting rules (currently baked into each `buildPrompt`) are applied consistently before a reply returns | `concierge-backend/ai/synthesize.go` | Nothing new — today's flow already returns raw model text; this would be a thin pass-through initially | NO | NO |

The important property: every module above is proposed to be a **wrapper around code that already exists and already runs**, not new logic.

---

## 10. Required interfaces

**A. Request envelope**
`user_id` (nullable — anonymous chat has none today), `profile_id`, `mode` (`concierge` | `brave_pa` | `news_gen` — currently undistinguished), `permissions` (n/a today, placeholder), `route/source` (which frontend surface called in), `risk_level` (n/a today, placeholder), `conversation/session state` (`session_id`, existing `conversation_id`).

**B. Model adapter**
Input: envelope + `messages[]` + `system_prompt`. Output: reply text + raw provider response. Errors: today's `handleChat` just returns a flat 500 — the adapter should preserve that behaviour but make the error typed. Provider metadata: model name/version (already known, just not surfaced). Cost/latency placeholder: not currently measured anywhere — would be new instrumentation, additive only.

**C. Memory/provenance**
What's retrieved: none today — each `/chat` call is stateless beyond client-resent `messages[]`. Where it came from: N/A yet. Timestamp/confidence/source: N/A yet — this interface today would just describe the *existing* conversation-history tables, not add retrieval behaviour.

**D. Verification/audit**
What's checked: nothing today. Evidence: N/A. Risk result: N/A. Pass/fail/warn: today would always be "pass" (no-op) — a seam for later.

**E. Tool/action gate**
Requested action: today only "booking request" exists, detected by regex on the reply text in `Chat.jsx`. Permission needed: none today (unauthenticated visitor can trigger it). Approval status: none — it fires automatically. Side effects: one `POST /booking-request` write. Rollback/undo: none exists (owner can only change `status` via `PATCH /owner/bookings/:id`).

**F. Final response synthesis**
Concierge voice today is entirely encoded in the system prompt text (`lib/buildPrompt.js`); there is no separate synthesis step. Trust/verification metadata: nowhere to attach it yet — this interface would be the future seam for a Trust Dot–style annotation, not something to build now.

---

## 11. No-behaviour-change feasibility

**YES, but only in stages.**

Why: every module proposed in §9 is a pass-through wrapper around code paths that already exist and are already exercised in production (`callAnthropic`, `db.SaveConversation`, the booking-intent regex). None require touching the Anthropic request body, the response shape returned to the frontend, or the DB schema. The risk is entirely in *sequencing and discipline*, not in the wrapper concept itself.

- Which existing flow would be first to wrap: `handleChat` / `POST /chat` — highest-leverage, best-understood, already battle-tested, benefits all three consumers (Concierge, Brave PA, news gen) at once.
- Which flow should be left untouched: **auth** (`/auth/login`, `/auth/signup`, `middleware.ts`, `lib/auth.js`) — just hardened (cookie-name fix + `_directToken` bypass fix, commits `f0673a1`/`8fa6229`) and out of scope. Also leave Stripe/payment path and `handleMediaUpload` untouched — no AI involvement.
- What tests would prove behaviour is unchanged: a before/after diff of raw HTTP responses from `POST /chat` given identical input — byte-for-byte reply text, identical `conversation_id`/`score` shape, identical DB rows written. Golden-file/snapshot test, not a new test category.

---

## 12. What would require database schema changes

**Not needed for minimal spine:** request envelope, model adapter, provider adapter, verification/audit no-op, action-gate wrapper, response-synthesis pass-through — all of §9 reads/writes exactly the tables that already exist.

**Needed later for real memory/provenance:** a way to distinguish `mode` (concierge vs brave_pa vs news_gen) on `conversations` rows (currently inferred, not stored); an audit-log table (request id, model used, verification result) — doesn't exist today; possibly a `profile_data` schema/typed-columns migration if the JSON blob ever needs to be queried rather than fetched-and-parsed.

**Should be delayed to V2/V3:** embeddings/vector columns, RLS policy work beyond what already exists, storage buckets beyond the existing `media` bucket, any indexes tied to retrieval-quality tuning.

---

## 13. What would require new provider keys or external services

**Not required for minimal spine:** everything — the minimal spine's entire point is to wrap the existing `ANTHROPIC_API_KEY` path with zero new providers.

**Optional later:** OCR/document AI (only relevant if `handleImportServices`'s vision call needs to scale beyond Claude vision), embeddings/reranking (only relevant once real memory retrieval is built), moderation service (only relevant once the client-side "human review" UX gets a real backend counterpart).

**Definitely V2/V3:** OpenAI, Gemini, Perplexity, transcription, MCP servers — all explicitly named in the Fleet Architecture Brief as later-phase lanes, none touched by anything in this repo today.

**Confirmed explicitly:** the minimal spine, if built, should use the existing Anthropic path only (`claude-sonnet-4-6` for chat, `claude-haiku-4-5-20251001` for vision) — no new provider unless Bruno approves otherwise.

---

## 14. Risk comparison: Brain Spine before Cinematic Shell vs after

**Option A — Minimal Brain Spine before Cinematic Shell**
- Benefits: wraps `handleChat` (the one path all three AI features share) once instead of three times later; the client-supplied `system_prompt` trust gap and the cosmetic "human review" gap (both found in this audit) become visible/nameable inside a spine rather than staying implicit in `Chat.jsx`.
- Risks: any mistake touches the single most business-critical endpoint in the product (public chat, revenue-adjacent via booking intent); done right after an auth-hardening night, it adds a second live-wire change in close succession.
- Files likely affected: `concierge-backend/main.go` (`handleChat`, `handleImportServices`), new `concierge-backend/ai/*.go` files. Frontend untouched if done as backend-only wrappers.
- Delay risk: low — additive, doesn't block Cinematic Shell technically.
- Regression risk: low-to-medium, contingent entirely on discipline (pass-through only, golden-file tested per §11).
- Impact on launch timeline: small if scoped exactly as §9 describes; grows fast if scope creeps into fixing the trust-gap/review-gap findings at the same time.
- Effect on future rebuild avoidance: real — Brave PA, Concierge chat, and news-gen currently share one undifferentiated code path; a spine here is the natural place to formalize that sharing before more features pile onto it.

**Option B — Cinematic Shell first, Brain Spine later**
- Benefits: visible progress now; zero risk to the just-stabilized auth/chat path; lets AUTH GREEN and recent fixes sit and prove stable before another change lands.
- Risks: Brave PA / news-gen / Concierge chat keep growing on the undifferentiated `handleChat` path in the meantime — each new feature added the old way is one more thing the eventual spine has to migrate.
- Files likely affected: Three.js/visual layer only — no overlap with `main.go` or the AI path.
- Rework risk: moderate — not "wasted work," but every new AI touchpoint added between now and the spine is a small amount of extra migration surface later.
- Launch/presentation benefit: real and immediate — Cinematic Shell is what gets shown, not the backend architecture.
- Product stability impact: best of the two options — no change to the just-hardened, revenue-facing chat path.

---

## 15. Recommended roadmap option

**B — Do Brain Spine audit/documentation now (this document), then Cinematic Shell, then implement the spine later**, with one addition: the two concrete gaps this audit surfaced (client-trusted `system_prompt`, and "human review"/sensitive-topic alerts being UI-only cosmetic gates with no backend record or real owner notification) are worth a short, separate conversation with Bruno before Cinematic Shell starts — not because they block anything, but because they're real product-trust claims (the consent screen literally tells visitors "AI conversations... processed by Anthropic" and the UI shows "🔔 owner notified" when no notification is sent) that exist independently of Brain Spine timing.

Justification: the repo evidence supports Option B's own stated trigger conditions — current AI flow is small, well-understood, and centered on one function (`handleChat`), which is exactly the situation the decision pack says makes a later, better-scoped spine implementation low-risk; nothing found here suggests urgency to spine-wrap before Cinematic Shell, and nothing found here suggests the AI flow is too fragile to revisit later.

If delaying (as recommended), what must happen before reopening it:
1. Cinematic Shell ships and is stable.
2. Core Flow QA passes.
3. Bruno decides whether to address the `system_prompt`-trust and cosmetic-review-gate findings as their own small task, independent of spine timing.

---

## 16. Exact next Claude Code prompt

```text
AI: Claude Code
Mode: Repo Execution Mode / Documentation + Audit Only
Objective: Prepare the Cinematic Shell Integration Packet review checklist and confirm no
AI/chat-path files are touched by it, using the Brain Spine Readiness Audit
(docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md and
docs/architecture/BRAIN_SPINE_FLEET_INTEGRATION_DECISION_PACK_v0.2.md) as the reference for
what must stay untouched.

Output expected:
- confirmation that Cinematic Shell's integration packet (once provided by ChatGPT/Gemini)
  does not modify concierge-backend/main.go, lib/buildPrompt.js, components/Chat.jsx,
  components/BravePAv2.jsx, middleware.ts, or lib/auth.js;
- a short list of the two findings from this audit (client-trusted system_prompt on
  POST /chat; cosmetic-only human-review/alert UX with no backend record or owner
  notification) presented to Bruno as a separate decision, not bundled into Cinematic Shell;
- nothing else.

Do not touch / do not decide:
Brain Spine implementation, new provider APIs, Supabase schema, MCP, billing, auth, routing,
Stripe, production behaviour, the two findings above (flag only, do not fix).

Exact scope: read-only review of the Cinematic Shell integration packet against the AI-path
file list above.

Tests: none — this is a read-only compatibility check, not a code change.

Stop condition: stop immediately after returning the three outputs above. Do not proceed to
implementation of Cinematic Shell or Brain Spine without Bruno's explicit next instruction.
```

---

## 17. Final recommendation

Document only for now (Option B) — the AI path is small and centered on one function, so a better-scoped spine later is lower-risk than spine work squeezed in before Cinematic Shell; but raise the two trust-relevant findings (client-controlled system prompt, cosmetic-only review/alerts) with Bruno as their own decision, separate from spine timing. Next: use the §16 prompt. Stop.
