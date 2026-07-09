# External Dependencies, Costs, and Legal Readiness v1

**AI:** Claude Code
**Mode:** Repo Execution Mode / Strategy Documentation Only
**Date:** 2026-07-09
**Status:** Strategy document. No code implemented, no code changed.

Consolidates every external service, provider account, environment variable name, cost exposure, and legal/DPIA requirement identified across all four audits into one reference. No secret values appear anywhere in this document — names only, per the standing rule already applied throughout this audit series.

---

## 1. Already-integrated external services (live-config status noted)

| Service | Used by | Live-confirmed this session? | Notes |
|---|---|---|---|
| Anthropic | Chat (`/chat`), vision import (`/ai/import-services`) | **YES** — real replies confirmed via `/demo/bruno` | Core AI provider, no change recommended |
| Supabase (DB) | All persistent data | **YES** — backend `/health` responded live | No migration pattern exists in this repo — first new table (Roadmap Phase 1) should establish a deliberate, reviewed precedent |
| Supabase (Storage) | Media upload | Bucket policy **NOT verified** | Two independent upload flows both depend on it; consistent usage suggests it works, but unconfirmed live |
| Render (backend hosting) | Go backend | **YES** — `/health` responded live | **Confirmed on free tier** per the repo's own doc (`PRODUCT_OPERATING_SYSTEM_v0.5.md` §7.7: "Render upgrade... must happen before any real client is sent a link"). Cold-start latency is a real risk for Stripe webhooks and synchronous notification email specifically. |
| Resend | Hot-lead email, sensitive-topic alert email | `RESEND_API_KEY`/`OWNER_EMAIL`/`RESEND_FROM_EMAIL` **NOT verified live** | Code path is real and correctly degrades to `disabled_missing_env` if unset — verification task, not a build task |
| Stripe | Connect, checkout, webhook | `STRIPE_SECRET_KEY`/`STRIPE_WEBHOOK_SECRET` **NOT verified live**; no real charge confirmed processed | Highest financial stakes of any unverified item in this document set |
| Google (partial) | Reverse-geocoding via Nominatim/OpenStreetMap (onboarding location capture) | Not independently verified, but a standard public API, low risk | Not the same as the Google Calendar OAuth integration below, which does not yet exist |

## 2. Needed, not yet integrated

| Need | For | Provider status | Recommendation |
|---|---|---|---|
| Google Calendar OAuth | Manager Agent's first connector (Blocco A3) | May **partially already exist** — the Master Doc excerpt states CLIENT_ID/CLIENT_SECRET retrieval from Google Cloud Console was "the last confirmed step" before interruption | **Confirm with Bruno before assuming a fresh setup is needed** |
| Weather API | Location-aware background | **UNKNOWN, provider not chosen** — 3 viable options: (a) OpenWeatherMap (free tier), (b) WeatherAPI.com (free tier), (c) Open-Meteo (genuinely free, **no API key required at all**) | Recommend Open-Meteo specifically — avoids secret management entirely |
| A scheduler | Digest/clustering, real "daily" news automation | **UNKNOWN** — needs a scheduler mechanism the current stack has no native equivalent for beyond Render's own Cron Job product | Recommend Render Cron Job (native to existing hosting); **COST_VERIFY_REQUIRED against current Render plan tier** |
| Speech/TTS provider | Voice I/O (if legally cleared) | **UNKNOWN, provider not chosen, and selection is secondary to the legal gate** — 3 viable options: (a) browser-native Web Speech API (free, inconsistent quality), (b) a hosted STT/TTS provider (paid), (c) a future Anthropic voice-capable endpoint if relevant (speculative) | Do not select until Roadmap Task 10.1 (legal scope resolution) completes |
| Vector store / embeddings provider | Advanced semantic memory | **UNKNOWN, provider not chosen** — explicitly V2/V3, not urgent | pgvector on Supabase is the lowest-new-infrastructure option if/when pursued |

## 3. Explicitly not needed

- **A paid feature-flag SaaS** (GrowthBook, LaunchDarkly) — recommend building the feature-flag mechanism in-repo (a new Supabase table + simple read/write), cheaper and more aligned with the existing stack at current scale.
- **A new push-notification provider** — Web Push requires no paid account, only a self-generated VAPID key pair.

---

## 4. Required env vars / secrets checklist (names only, no values)

### Already in code, live status to confirm
`ANTHROPIC_API_KEY`, `SUPABASE_URL`, `SUPABASE_KEY`, `RESEND_API_KEY`, `OWNER_EMAIL`, `RESEND_FROM_EMAIL`, `STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET`, `AUTH_SALT_SECRET`, `ADMIN_KEY`, `OWNER_KEY`.

### Net-new, not yet created
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_OAUTH_REDIRECT_URI` — Google Calendar connector.
- `VAPID_PUBLIC_KEY`, `VAPID_PRIVATE_KEY`, `VAPID_SUBJECT` — push notifications. Must be generated fresh; never hardcoded or reused from any example/tutorial.
- `WEATHER_API_KEY` — only if a keyed provider is chosen over Open-Meteo.
- A speech-provider key set — names TBD once a provider is selected, which itself is gated behind the Roadmap's Phase 10 legal resolution.
- Optionally, `NEXT_PUBLIC_APP_URL` — to prevent recurrence of the hardcoded-domain issue (REG-10a in the Feature Registry).

---

## 5. Paid subscription / cost exposure list

| Item | Status | Cost category |
|---|---|---|
| Anthropic API usage | Real, ongoing | Confirmed real cost |
| Render hosting | Real, ongoing (free tier today) | **Confirmed real cost — upgrade explicitly required before real client links go out, per the repo's own Master Doc pricing rule** |
| Resend | Has a free tier, volume-dependent | Likely low cost at current scale, **COST_VERIFY_REQUIRED** at growth |
| Stripe | Transaction-fee based, no subscription | Confirmed real cost structure, no fixed fee |
| Render Cron Job / plan upgrade | Potential new cost | **COST_VERIFY_REQUIRED** — depends on current plan tier |
| Weather API | Avoidable | Choose Open-Meteo to keep this at zero |
| Speech/TTS provider | Potential new cost, if voice is ever cleared | **COST_VERIFY_REQUIRED**, gated behind legal review anyway |
| Google Calendar API | Likely free at this product's scale | Standard free tier almost certainly sufficient |
| Feature-flag SaaS | Avoidable | Build in-repo instead |
| Vector store / embeddings | Potential new cost, V2/V3 only | **COST_VERIFY_REQUIRED** if/when pursued |

---

## 6. Legal / privacy consultation needs

- **What kind of professional:** a data-protection/GDPR specialist covering both UK and EU GDPR (the product's own privacy policy claims dual compliance) — not a generalist lawyer; ideally with DPIA drafting experience.
- **Why:** three concrete gaps, in priority order:
  1. **`/privacy` returns a live 404 in production** — an active compliance defect right now, not a future risk (Roadmap Task 2.1, technical fix; legal review of content is a parallel follow-on).
  2. **Health-category data collection via legal forms has no documented DPIA-trigger review** (Roadmap Task 7.2) — the narrowest, most achievable legal deliverable identified across all four audits, since this data is already being collected today.
  3. **Voice/video work is entirely pre-conditioned on a DPIA that doesn't exist** (Roadmap Phase 10) — already correctly not built, but the review should exist well before any engineering work starts, not scoped reactively.
- **Before which activation:** hard gate before any voice/video engineering work (Roadmap Task 10.2); soft gate (should happen soon, not launch-blocking by itself) before confidently calling the health-disclosure forms "compliant."
- **Rough urgency level:** **HIGH** for the dead privacy-policy link and the health-form DPIA-trigger question — both concern data already live in production today. **LOW** for voice/video — nothing is built, so nothing is at risk yet, but the review should be commissioned well ahead of any engineering interest in the feature, not after.
- **What documents to prepare first:**
  1. A real DPIA covering current health-data collection via legal forms — achievable now, narrower scope than the voice/video question.
  2. A confirmed, legally-reviewed refresh of `privacy.html`'s content once it's technically reachable again (Roadmap Task 2.1).
  3. A data retention and deletion/export mechanism to actually back the GDPR rights `privacy.html` already claims to offer — currently described in the policy text but not implementable by a real user (no in-product delete-my-data control exists anywhere in the repo).

---

## 7. Legal / DPIA locked features summary

| Feature | State | Safeguards present | Gate |
|---|---|---|---|
| Voice & Video Session Learning | `LEGAL_LOCKED` | 0 of 8 required (DPIA, consent flow, privacy notice, retention policy, deletion/export controls, biometric/health-data safeguards, processor list, activation gating) | Full DPIA (Roadmap Task 10.2) |
| Voice I/O without learning | `LEGAL_LOCKED` (conservative default) | Scope not yet separated from the above | Legal scope resolution (Roadmap Task 10.1) |
| Advanced/semantic memory (profiling-adjacent) | Conditionally locked | N/A — not yet built | DPIA if genuine behavioral learning is pursued (Roadmap Task 9.3) |
| Health-category legal forms | Not locked, but gap flagged | Mechanism real; DPIA-trigger review missing | Health-Data DPIA Scoping (Roadmap Task 7.2) — does not require un-shipping the feature, only formalizing the review |

---

## 8. Cross-reference

This document consolidates but does not replace the full detail in `docs/audits/TECHNICAL_CONSTRAINTS_MATRIX.md` (per-feature technical constraints, items A-Z) and `docs/audits/MASTER_VISION_VS_REPO_REALITY_REVISION.md` (live production evidence). Consult those documents for the reasoning behind each line item above.
