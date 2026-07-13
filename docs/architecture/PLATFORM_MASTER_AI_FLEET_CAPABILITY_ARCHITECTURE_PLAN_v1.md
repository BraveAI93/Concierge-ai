# Platform Master / AI Fleet / Capability Architecture Plan v1

**AI:** Claude Code
**Mode:** Repo Review Only / Day 6.5 Platform Master + AI Fleet + API Capability Architecture
**Date:** 2026-07-13
**Status:** Architecture review and decision-ready plan only. No code implemented, no schema changed, no feature flags flipped, no env vars added.

Context: this document synthesizes the current, live state of Sprint 01 (Days 3-6) against the pre-existing strategy set — `docs/architecture/FLEET_BRAIN_ARCHITECTURE_BRIEF.md`, `docs/architecture/BRAIN_SPINE_FLEET_INTEGRATION_DECISION_PACK_v0.2.md`, `docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md`, and `docs/strategy/FEATURE_REGISTRY_AND_ACTIVATION_STATES_v1.md` — to answer: how does The Concierge become one Platform Master with two access layers, and how do real capability APIs (web search, weather, calendar, maps, tickets, booking actions) plug into that same architecture without becoming random bolt-ons.

---

## Conclusion (preserve exactly)

- Owner PA and Public Concierge are two access layers of one Platform Master.
- Shared `feature_flags` and `audit_events` govern truth for both layers.
- Web search is the first activation target.
- No further external capability should be coded until `BRAVE_SEARCH_API_KEY` exists and one real web-search end-to-end test passes.
- A multi-AI model router / provider router is future architecture, not immediate implementation.

---

## 1. Current state

**AI/model logic today:** One live model call site, `callAnthropic()` (`claude-sonnet-4-6`), shared undifferentiated by `POST /chat` across Concierge public chat, Brave PA, and dashboard news-gen. `callAnthropicVision()` (`claude-haiku-4-5-20251001`) handles onboarding service-import only. Prompt construction is cleanly split into two functions in `lib/buildPrompt.js` — `buildPrompt()` (public) and `buildBravePAPrompt()` (owner) — that split is the existing two-layer seam. No request envelope, no verification/audit-of-reply step, no memory re-read (conversations/messages are stored but never read back into a later prompt), no structured tool/action gate (booking intent is still regex-sniffed off free text).

**Provider/API logic today:** Anthropic only, live. Resend (email) is real and separate from AI. Stripe is real code with unknown live-env status. No Brave Search, OpenAI, Gemini, Perplexity, Google Calendar, or maps code runs in production. A Brave Search adapter exists (`braveWebSearch()`, Day 6) but has never made a real network call.

**Already real:** `feature_flags` + `GET /flags` (Day 3); `audit_events` write-only store + `logAuditEvent()` (Day 4); Brave PA capability-truth prompt logic reading live flags (Day 5); owner-only `/pa/websearch` with flag→auth→provider gate ordering and its own three-event audit trail (Day 6).

**Built-but-locked:** the entire Day 6 web-search path — code complete, `pa_web_search` still `GHOST_FORBIDDEN` live, no `BRAVE_SEARCH_API_KEY`.

**Still ghost/forbidden (live-confirmed 2026-07-13):** `pa_weather`, `pa_calendar`, `pa_ticket_search`, `pa_booking_actions`, `memory_personalization`, `cohort_learning`, `notification_preferences`, `notification_sound`, `push_notifications` all `GHOST_FORBIDDEN`; `voice_io`/`voice_video_session_learning` `LEGAL_LOCKED`. REG-19a in the registry doc is now partly stale — it predates Day 5/6's truth-wiring and should be reconciled in a future docs-only pass.

---

## 2. Platform Master model

Owner PA and Public Concierge are already, structurally, two front doors onto one spine — not two products. They share: `profiles.profile_data`, `feature_flags` (one row governs what either layer may honestly claim), the `callAnthropic` adapter/credentials, and `audit_events` as a single provenance store.

**Must never cross owner → public:** session tokens, any other client's lead/contact data, internal capability-state language (a visitor must never see "GHOST_FORBIDDEN"/"DORMANT_BUILT"-style internal framing — only the honest, plain-English fallback copy Day 5 wrote), raw `audit_events` rows (already enforced at the DB grant level — no anon SELECT).

**Must never cross public → owner:** nothing meaningful crosses this direction today; visitors don't write anything owner-private.

**Formalization, not rebuild:** the two-function prompt split already existing in `lib/buildPrompt.js` is the correct shape. The Platform Master model just names it: one shared spine (`model_router` + `provider_adapters` + `feature_flags`/`audit_events`) behind two personas that never speak in each other's voice.

---

## 3. AI Fleet architecture

- **model_router** — doesn't formally exist; today it's two hardcoded model strings (sonnet for chat, haiku for vision) picked implicitly by which function is called. A router would sit exactly where `handleChat` dispatches to `callAnthropic()`.
- **provider_adapters** — `callAnthropic()`/`callAnthropicVision()` are the de facto Anthropic adapters; `braveWebSearch()` (Day 6) is the first non-Anthropic adapter, and it's already correctly isolated behind a narrow `query → []WebSearchResult` interface.
- **task_classifier** — doesn't exist as a shared module yet. The Day 6 `searchTriggerKeywords` heuristic in `BravePAv2.jsx` is its seed, not its final form — if weather/calendar/tickets each grow their own inline keyword list, that's the exact duplication to avoid; it should graduate into one shared, testable classifier once a second capability lands.
- **permission_boundary** — already the right shape: flag-state check, then `authenticateToken()`, then provider call (Day 6's `/pa/websearch`). Generalize this ordering to every future capability rather than reimplementing per-endpoint.
- **audit_events usage** — already the fleet's provenance/observability primitive (Day 4-6): request/completed/failed triads, low-risk metadata only. Natural next step (not now): add a `layer` (owner/public) field so it can answer "what did the fleet do, for which layer, how often" as more capabilities land.
- **fallback behaviour** — already established and correct: missing/failed capability → honest, useful answer, never silence, never fabrication (`buildCapabilitySummary`). Standing rule for every future capability.
- **avoiding "multiple bots":** capability adapters must return structured data only; only the active persona's prompt-builder decides phrasing. Never let a specialist lane "speak" directly.

---

## 4. External capability API architecture

All six should share one shape — flag check → auth/permission check → provider call → structured+cited result → audit event → persona-voiced response. Day 6's `/pa/websearch` is the reference implementation.

| Capability | State | Note |
|---|---|---|
| Web search | Built, locked | Owner-only, cited, `GHOST_FORBIDDEN` pending key |
| Weather | Not started | Reuse the same 3-gate shape; `pa_weather` flag already seeded |
| Calendar (read-only) | Not started this sprint | REG-22 `SPEC_ONLY`, schema field reserved; deferred not cancelled (original Day 6 plan) |
| Maps/geocoding | Not started | Correctly later per the stated safest order |
| Ticket/event search | Not started | Likely same structured+cited shape as web search |
| Booking actions | Not started | Highest risk — real side effects; should imitate the existing "propose, don't auto-act" pattern already proven in `booking_requests` (REG-06), not bypass it |

---

## 5. Provider abstraction

Today's only capability-provider hardcoding (Brave Search) is already correctly isolated inside `braveWebSearch()` — prompt-building code only ever sees structured `WebSearchResult`, never Brave's request/response shape. That's the right instinct to keep formal later: a generic `CapabilityResult{Title, URL, Source, Snippet, Timestamp}` shape already happens to fit weather/tickets/maps too. What should stay provider-specific: the HTTP call, auth header, and response-parsing struct — exactly what Day 6 did. Model-provider abstraction (Claude vs. a future OpenAI/Gemini lane) is a separate, larger axis than capability-provider abstraction (Brave Search vs. an alternative search API) — nothing here needs the former yet since only one model provider is active.

---

## 6. Feature flags

- **AI Fleet control flags:** none exist by that name yet — not needed with a fleet of one provider.
- **Web/weather/calendar/tickets:** `pa_web_search`, `pa_weather`, `pa_calendar`, `pa_ticket_search` — all seeded, all `GHOST_FORBIDDEN` today.
- **Owner-only:** `owner_pa`, `dashboard_insights`, and all five `pa_*` capability flags — by design, per the stated owner-first order.
- **Public-safe:** `ai_processing_consent`, `public_real_chat`, `public_visual_shell`.
- **Must stay locked indefinitely (not an engineering gap):** `voice_io`/`voice_video_session_learning` (`LEGAL_LOCKED` — 0/8 legal safeguards, REG-29); `memory_personalization` (`GHOST_FORBIDDEN` — no retrieval infra); `cohort_learning` (`GHOST_FORBIDDEN` — no legal-basis review); `pa_booking_actions` (should stay locked longest of the six PA capabilities — highest real-world risk).

---

## 7. Env/API credentials

- **Eventually required:** `BRAVE_SEARCH_API_KEY` (+ optional `BRAVE_SEARCH_BASE_URL`) — already coded for; `GOOGLE_CLIENT_ID`/`GOOGLE_CLIENT_SECRET`/`GOOGLE_OAUTH_REDIRECT_URI` — calendar, per the original Day 6/7-8 plan; weather/maps/ticket-search provider keys — not yet chosen.
- **Missing today (confirmed 2026-07-13):** `BRAVE_SEARCH_API_KEY` absent from both `concierge-backend/.env` and process env. Google Calendar credential status is still the standing open decision from the original Day 6 plan.
- **Activate first:** web search — it's the only capability fully built and gated, needing only the key + one real end-to-end proof.

---

## 8. Review-gate system

- **Perplexity:** live/sourced research — e.g. choosing a weather/maps/ticket provider, pricing/ToS comparison. Matches its existing role in the Fleet Brief.
- **Gemini:** visual/UX and multimodal judgment only — not relevant to this capability work.
- **Claude Chat:** architecture sequencing and reconciling recommendations before Bruno decides.
- **Claude Code:** execution only, after Bruno approves a specific scoped packet — never unprompted architecture decisions.

---

## 9. Recommended next implementation packet

**Day 7 should build first:** nothing code-side until `BRAVE_SEARCH_API_KEY` exists. The real next packet is unchanged from Day 6's own conclusion: get the key → set it in Render + local `.env` → run one real owner-only search end-to-end (query → cited results → audit events written) → only then flip `pa_web_search` to `ACTIVE_PRIVATE`.

**What should wait:** weather/calendar/maps/tickets/booking-action code; any `model_router`/`task_classifier` formalization — premature with only one live capability-provider; any `audit_events` schema change (e.g. a `layer` field) — worth proposing later, not urgent now.

**Stop conditions:** stop the moment `BRAVE_SEARCH_API_KEY` is unavailable (document, don't fake); stop before touching any other capability's code without its own separate approved packet (the "one active integration at a time" rule); stop before generalizing an adapter interface until at least two real capability-providers exist to justify it.
