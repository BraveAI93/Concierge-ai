# Cohort Learning and Anonymous Product Intelligence v1

**AI:** Claude Code
**Mode:** Repo Execution Mode / Strategy Documentation Only
**Date:** 2026-07-09
**Status:** Strategy document. No code implemented, no code changed. Nothing described here exists in the repo today — this entire document is `SPEC_ONLY`/new design, not an audit of existing behavior.

## Purpose

Bruno wants The Concierge to learn faster from aggregate usage patterns without building personal profiles of individual users or a central memory of who anyone is. This document defines a strict boundary between **personal memory** (which this repo does not yet build — see `docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md` §K, memory/learning/personalization) and **cohort learning** (aggregate, anonymous, operational-pattern learning across many users/businesses, which also does not yet exist). It is written so that a future Claude Code implementation task can build the minimum safe version without guessing at the privacy boundary.

---

## 1. Personal Memory vs Cohort Learning

| | Personal Memory | Cohort Learning |
|---|---|---|
| Subject | One identifiable user, visitor, or business owner | An aggregate pattern shared by many, no individual traceable |
| Storage unit | A record tied to a person/session/account | A statistic/count/rate, not a record |
| Purpose | Personalize *this* user's experience ("remembers you") | Improve the *product* for everyone (better default copy, better routing, better prompts) |
| Consent basis | Requires the user's own consent to that specific processing | Can potentially rely on legitimate interest **if and only if** the aggregation is genuinely irreversible — see §14 |
| Current repo state | `GHOST_FORBIDDEN` as a claim; `DORMANT_BUILT` as inert storage (per `docs/strategy/FEATURE_REGISTRY_AND_ACTIVATION_STATES_v1.md` REG-30/REG-30a) | Does not exist in any form — this document proposes it for the first time |
| Relationship to this doc | Explicitly out of scope — governed separately by the Brain Spine audit and Roadmap Phase 9 | In scope |

**The hard rule this document exists to enforce:** cohort learning must never become a backdoor to personal memory. A cohort statistic ("62% of personal-trainer visitors ask about pricing before availability") is allowed. A record that lets anyone reconstruct "this specific visitor asked about pricing" is not cohort learning — it is personal memory, and must go through the personal-memory consent path, not this one.

---

## 2. Anonymous / aggregated operational signals

An "operational signal" in this framework is a small, structured, non-narrative fact extracted from an interaction — never the interaction itself. Examples: "a pricing question occurred," "the conversation reached the booking-proposal stage," "the visitor's browser language was Italian." A signal is only usable for cohort learning once it has been stripped of anything that could tie it back to one visitor, one business, or one conversation, **directly or indirectly** (see §14 on why pseudonymous data does not qualify).

The distinction that matters most for this product specifically: **operational signals describe what happened in a category of interaction, never what was said.** No message content, no paraphrase, no summary of a specific exchange qualifies as an operational signal — those are transcript-adjacent and belong to personal memory's consent boundary, not this one.

---

## 3. What can be collected in minimal mode

- Coarse **event types**, not content: `pricing_question_asked`, `booking_proposed`, `booking_confirmed`, `sensitive_topic_flagged` (type only, e.g. "Injuries & pain" — already an existing category in `lib/constants.js`'s `alertMap`, not new), `lead_captured`, `form_completed`.
- **Business category** (already a first-class field in this product — personal trainer, photographer, chef, etc., from the existing `profession`/mode fields in `profiles.profile_data`).
- **Coarse funnel stage** reached (browsed → asked question → lead captured → booking proposed → booking confirmed) — a stage label, not a transcript.
- **Interface language** (already collected — `lang` parameter, already a UI-level signal, not new).
- **Coarse timing** (hour-of-day bucket, day-of-week) — never a precise timestamp tied to a specific session in the aggregate store.
- **Aggregate counts and rates only** — "N conversations reached booking stage this week across all personal trainers," never "conversation X reached booking stage."

## 4. What must not be collected in minimal mode

- **Full or partial chat transcripts** — explicitly forbidden by this task's own rules, and structurally the highest re-identification risk in this product given how specific and personal service-business conversations are (an injury description, a specific date/venue, a specific price negotiation can each be uniquely identifying on their own).
- **Message content of any kind**, including "just a snippet" or "just the first message" — a single sentence can be enough to re-identify a small local business's client.
- **Session identifiers, IP addresses, device fingerprints, or any token that could be joined back to `conversations`/`messages`/`leads` tables** — cohort data must live in a separate store with no foreign key, not even an indirect one, back to identifiable tables.
- **Names, emails, phone numbers, or any contact detail** — even in aggregate-adjacent form (e.g., "the most common lead email domain" is still forbidden — domain-level email data is still personal data under GDPR).
- **Exact geolocation** — city-level at coarsest is defensible for a "location" signal; anything more precise (postcode, lat/lng) is forbidden in this store, even though the product already collects lat/lng for other real purposes (onboarding location display) — that data must never flow into the cohort store.
- **Any protected/sensitive characteristic** — see §8 below.

## 5. What can be learned from cohorts

- Which business categories see the highest booking-conversion rates, to improve default prompts/copy for that category.
- Which coarse funnel stage most conversations drop off at, per business category — to improve the product's own flow (e.g., "the lead-capture form appears too early for chef inquiries").
- Which sensitive-topic categories (by type, not content — e.g., "Injuries & pain" as a label) trigger most often, to prioritize the Notification Center work in `docs/strategy/REAL_BUILD_ROADMAP_v1.md` Phase 4 toward the categories that matter most in practice.
- Aggregate tone-variant performance (see §7) — which prompt tone framing correlates with higher booking-confirmation rates, in aggregate, across many businesses of the same category.
- Language-level patterns — e.g., whether Italian-language conversations convert differently from English, to inform localization priority.

## 6. What must never be inferred without explicit consent

- Anything about an individual visitor's identity, demographics, or behavior pattern across sessions.
- Anything that infers a protected characteristic (see §8) from behavioral signal, even indirectly (e.g., inferring likely health status from which service category someone repeatedly asks about — this is exactly the kind of "innocuous signals combine into a sensitive inference" risk that k-anonymity thresholds alone do not fully protect against, and is why §8 exists as a hard rule, not just a default-off setting).
- Anything about a specific business's performance relative to *named* competitors — cross-business comparison must stay at the category level, never "Bruno's Concierge converts worse than [specific other named professional]."
- Any consequential decision about a specific user or business based on cohort data (see §15's forbidden-use rule) — cohort data may change the *product* (default copy, prompt tuning, feature prioritization); it must never change what a specific user experiences based on which cohort they were placed in without their own separate, explicit consent for that personalization (which would then be personal memory, not cohort learning, and governed by that separate path).

---

## 7. Segment design

Segments are the axes cohort statistics are grouped by. Each is a coarse category, never a precise or continuous value.

| Segment | Values (examples) | Source (already exists in repo?) |
|---|---|---|
| Business category | personal trainer, chef, photographer, dance teacher, consultant, etc. | Yes — existing `profession`/onboarding mode fields |
| Visitor intent (coarse) | pricing inquiry, availability inquiry, general info, booking-ready | New — derived from existing coarse keyword categories already used in lead scoring (`scoreLead()`), not new invasive detection |
| Language | interface language selected | Yes — existing `lang` parameter |
| Service type | the service category within a profession (e.g., "deep tissue massage" bucketed as "wellness service," not stored verbatim) | Partially — needs bucketing logic, not raw service names, to avoid small-business re-identification via a distinctive service name |
| Funnel stage | browsed, asked question, lead captured, booking proposed, booking confirmed | New — derivable from existing event points already logged (leads, booking_requests) |
| Tone variant | which system-prompt tone framing was used (existing `tone` array in onboarding, e.g., "warm & personal," "direct & no-nonsense") | Yes — existing onboarding field |
| Outcome | booking confirmed / not confirmed / abandoned | New — derivable from existing `booking_requests.status` |

**Segment combination risk:** combining too many segments at once (e.g., business category + service type + language + tone variant + outcome, all at once) can shrink a cohort down to a handful of real businesses, defeating anonymity even though no single field is sensitive. §9's minimum thresholds must be enforced on the **combined** segment, not per-field.

## 8. Avoid using sensitive or demographic segmentation as default

The following must **never** be a default or automatic segmentation axis, and must never be inferred from behavioral signal even as a byproduct of another analysis: **age, sex/gender, health status, sexuality, ethnicity, disability, religion/political views.**

This is a hard rule, not a configurable default. If a future business need genuinely requires any of these (e.g., a health-and-wellness-specific product insight), that would require its own explicit, freestanding legal/consent review — it does not fall under this document's cohort-learning framework at all, and should be treated with the same caution as `docs/strategy/EXTERNAL_DEPENDENCIES_COSTS_AND_LEGAL_READINESS_v1.md`'s health-data DPIA discussion, not bolted onto general product analytics.

---

## 9. Minimum cohort thresholds to reduce re-identification risk

This product's current real-world scale makes this section unusually important: with a small number of onboarded businesses, a naive cohort ("personal trainers in London using the direct-tone variant") could easily resolve to one or two real businesses today, even without collecting anything individually sensitive.

- **Hard floor: no cohort statistic may be computed, stored, or exposed (even internally) for any combined segment with fewer than k=20 distinct underlying businesses/profiles contributing to it.** This is stricter than common k-anonymity minimums (often cited around k=5-10) specifically because this product's current scale makes smaller cohorts trivially re-identifiable by anyone who knows the local market (e.g., "the only Italian-speaking personal trainer in this city" is a real, current risk, not a hypothetical).
- Below k=20, the combined segment must be **generalized** (drop a segmentation axis, e.g., report at the country level instead of city level, or the broader "wellness" category instead of "personal trainer") until the threshold is met, or simply not reported at all for that period.
- This threshold should be revisited and can be relaxed somewhat as the number of onboarded businesses genuinely grows — but the review that changes it should be deliberate, not automatic, and should be documented (not silently adjusted in code).
- **Never report a raw count below the threshold, even as "fewer than 20."** Small-number disclosure (e.g., "1-4 businesses matched") can itself be re-identifying in a small local market — suppress the row entirely rather than showing a small-but-vague number.

## 10. Retention rules

- Cohort statistics (the aggregated counts/rates themselves) may be retained indefinitely, since — if correctly built — they contain no personal data by definition.
- The **underlying event-level signals** that feed into cohort aggregation (before aggregation, still segment-labeled but not yet summed into a k≥20 statistic) must have a **short, defined retention window** — recommend 90 days — after which raw event-level rows are deleted and only the already-aggregated statistics survive. This limits the blast radius if the pre-aggregation store were ever compromised or misconfigured.
- No cohort-learning data of any kind should ever be exported alongside, or joined with, personal data exports (e.g., a GDPR data-subject access request response) — the two systems must remain structurally and operationally separate, per the "keep personal memory separate from product analytics" rule.

## 11. Audit events needed

Extending the audit-events system already scoped in `docs/strategy/REAL_BUILD_ROADMAP_v1.md` Phase 1 (Task 1.2) rather than inventing a parallel system:

- `cohort_signal_recorded` — event type, segment labels, timestamp bucket (not exact timestamp) — no personal identifiers.
- `cohort_aggregation_run` — when a rollup from raw signals to k-anonymous statistics occurred, how many underlying signals were included, whether any segment combination was suppressed for failing the k≥20 threshold.
- `cohort_raw_signal_purged` — confirms the 90-day retention rule (§10) actually executed, with a count of rows purged, not their content.
- `cohort_data_access` — any time a human (Bruno, or a future team member) queries the aggregated cohort statistics, for the same "never claim something happened without a real backend record" discipline already proven in the notification work this week.

## 12. Feature flags / activation states needed

Using the `feature_flags` system scoped in `docs/strategy/REAL_BUILD_ROADMAP_v1.md` Phase 1:

- `cohort_learning_signal_collection` — gates whether raw signals are recorded at all. Should default **off** until Minimum Safe V1 (§ below) is built and reviewed.
- `cohort_learning_aggregate_display` — gates whether aggregated statistics are shown anywhere (even internally to Bruno) — should only flip on once §9's threshold enforcement is proven working, not assumed.
- `cohort_learning_public_analytics` — a separate, higher-bar flag for any future public-facing aggregate stat (e.g., a marketing claim like "businesses using The Concierge see X% more bookings") — this must stay off far longer than the internal flags, since public claims carry their own truthfulness bar per the existing Anti-Chaos Rule #11 ("Marketing does not outrun product truth").

## 13. Privacy notice copy needed

Once `/privacy` is reachable again (per `docs/strategy/REAL_BUILD_ROADMAP_v1.md` Task 2.1 — this document's recommendations are moot until that live defect is fixed), the privacy notice needs a new, clearly-labeled section, separate from the existing "AI Conversations" section, along these lines (**for legal review, not final copy — flagged as such**):

> **Product improvement analytics.** We collect anonymous, aggregated patterns about how our AI assistant is used — such as which types of questions are common for a given business category, or which stage of a conversation most often leads to a booking — to improve the product for all users. This analysis never includes the content of your conversations, your name, contact details, or any information that could identify you or the business you contacted. Statistics are only calculated and used when there are enough businesses in a category to prevent any one business or conversation from being identifiable.

This copy must not ship without: (a) the technical implementation actually matching every claim in it, and (b) real legal review — this document proposes structure and honesty constraints, not final legal language.

## 14. Consent or legitimate-interest decision needed

This is a genuine open legal question this document cannot resolve unilaterally, and is flagged as **SPEC_MISSING / requires Bruno + legal decision** before Minimum Safe V1 ships:

- Under GDPR, **genuinely anonymous data (irreversibly, with no reasonable means of re-identification) falls outside GDPR's scope entirely** (Recital 26) — if the aggregation design in this document is correctly implemented and the k≥20 threshold genuinely holds, the resulting *statistics* may not require a consent basis at all, because they are not personal data.
- **However, the pre-aggregation, segment-labeled raw signals (§10) are not yet anonymous — they are, at best, pseudonymous** (tied to a specific conversation/session even if not to a named person), and pseudonymous data **remains personal data under GDPR**. This raw-signal collection stage therefore needs its own legal basis — most plausibly **legitimate interest** (product improvement, with a documented legitimate-interest assessment weighing this against user privacy expectations) rather than requiring fresh explicit consent from every visitor, but this is a determination for Bruno's data-protection specialist to make formally, not an engineering assumption.
- **This document explicitly does not claim to have resolved this question.** It recommends: legitimate interest as the likely basis for the raw-signal collection stage (subject to a real legitimate-interest assessment), and no consent basis needed at all for the final aggregated statistics (subject to confirming the k≥20 design genuinely achieves anonymity, not mere pseudonymity) — but both of these are recommendations for legal review, not conclusions this document is authorized to finalize.

## 15. Activation state

- `ACTIVE_PUBLIC` — reserved **only** for genuinely anonymous aggregate analytics that have cleared the k≥20 threshold and the legal review in §14. Nothing in this system starts here.
- `ACTIVE_PRIVATE` — for early cohort tuning, internal-only (Bruno/team visibility of aggregated statistics), once Minimum Safe V1 is built and the raw-signal collection legal basis (§14) is confirmed. This is the realistic first real state for this system.
- `LEGAL_LOCKED` — for anything touching sensitive/health/voice/video behavioral learning, or any segmentation axis listed in §8 — these are not merely deprioritized, they are **forbidden as defaults** and would require their own freestanding legal review entirely separate from this document's framework, per the rule stated in §8.
- **No `GHOST_FORBIDDEN` state is acceptable for this system at any point** — per this task's own rule ("No GHOST capabilities"), if Minimum Safe V1 cannot be built with genuine k≥20 enforcement and genuine separation from personal-memory data, it should not ship at all rather than shipping a cohort-learning claim that isn't structurally true. This is the same discipline already proven in this week's notification-truthfulness work, applied to a new, higher-stakes domain (aggregate behavioral data across potentially many real small businesses and their real clients).

---

## Cohort Learning Model

A three-layer pipeline: (1) **coarse event signals** recorded per interaction, segment-labeled, no content, no direct identifiers, short retention (§3, §10); (2) **aggregation** into k≥20-thresholded statistics on a defined schedule, raw signals purged after 90 days (§9, §10); (3) **consumption** of only the aggregated, anonymous statistics by product decisions or (eventually, cautiously) public claims — the raw layer is never queried directly for product decisions, only the aggregate layer is.

## Allowed Signals

Coarse event types (pricing question, booking proposed, lead captured, sensitive-topic-flagged by category label), business category, coarse funnel stage, interface language, coarse time bucket, tone variant used, booking outcome — all segment-labeled, all content-free, all k≥20-thresholded before any use. See §3 and §7.

## Forbidden Signals

Full or partial transcripts, any message content or snippet, session/IP/device identifiers joined to identifiable tables, names/emails/phones/contact data, precise geolocation, any protected/sensitive characteristic (age, sex/gender, health, sexuality, ethnicity, disability, religion/politics) as a segmentation axis or inference target, any cohort below the k≥20 threshold. See §4 and §8.

## Minimum Safe V1

1. A `cohort_signals` table, separate from all identifiable tables, storing only: event type, segment labels (business category, coarse funnel stage, language, tone variant), a coarse time bucket — no foreign keys to `conversations`, `messages`, `leads`, or `profiles` beyond the business-category label itself resolved at write time (not stored as a joinable ID).
2. A scheduled aggregation job (can share the scheduler infrastructure being built for `docs/strategy/REAL_BUILD_ROADMAP_v1.md` Phase 4's digest work) that rolls raw signals into `cohort_statistics`, enforcing k≥20 and suppressing sub-threshold rows.
3. A 90-day purge job for the raw `cohort_signals` table, logged via `cohort_raw_signal_purged` audit events.
4. Internal-only (`ACTIVE_PRIVATE`) visibility of `cohort_statistics` to start — no public claim, no automated product behavior change yet, purely observational for Bruno.
5. The §14 legal basis determination completed **before** step 1 begins collecting any real data, not after.

## Technical Structure Needed

New tables: `cohort_signals` (raw, 90-day retention), `cohort_statistics` (aggregated, indefinite retention). New scheduled job (shares Phase 4's scheduler decision). New feature flags: `cohort_learning_signal_collection`, `cohort_learning_aggregate_display`, `cohort_learning_public_analytics`. New audit event types (§11). No new external provider or paid service required — this is buildable entirely on the existing Supabase/Go stack. Per this repo's established pattern, no migration mechanism currently exists — the same "explain exactly what's needed, don't improvise" discipline from the notification work applies to these two new tables.

## Legal/Privacy Notes

The central open question is §14 (legal basis for the raw-signal collection stage) — this must be resolved by Bruno's data-protection specialist before Minimum Safe V1 begins collecting any data, not treated as a formality. The privacy notice addition in §13 cannot ship until `/privacy` itself is reachable (`docs/strategy/REAL_BUILD_ROADMAP_v1.md` Task 2.1) and the copy has had real legal review. The k≥20 threshold (§9) is an engineering safeguard, not a legal one by itself — it supports but does not replace the legitimate-interest assessment.

## Recommended Build Order

1. §14 legal basis determination (Bruno + legal — not an engineering task).
2. `docs/strategy/REAL_BUILD_ROADMAP_v1.md` Phase 1 (feature flags, audit events) — this system depends on both.
3. Minimum Safe V1 (this document's own §, above) — raw signal collection + aggregation job + purge job, `ACTIVE_PRIVATE` only.
4. Internal review of real aggregated statistics over a meaningful period (recommend minimum 90 days of real data) before considering any `ACTIVE_PUBLIC` graduation.
5. Only then, a freestanding decision (with fresh legal review) on any public-facing aggregate claim — never assumed as an automatic next step.

Stop.
