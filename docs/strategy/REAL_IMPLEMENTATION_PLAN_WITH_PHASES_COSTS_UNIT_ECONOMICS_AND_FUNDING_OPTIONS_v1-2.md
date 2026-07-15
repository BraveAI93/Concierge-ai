# Real Implementation Plan with Phases, Costs, Unit Economics, and Funding Options v1

**Project:** Brave by Bruno / The Concierge  
**Owner:** Bruno Aversa  
**Date:** 2026-07-13  
**Status:** Strategy / planning document only. No implementation has been started by this document.  
**Mode:** Vision-to-execution reconciliation, cost planning, funding readiness.

---

## 0. Executive decision

The next build should not be a generic calendar assistant, a cosmetic 3D shell, or a loose collection of standalone ideas.

The product should be planned as one shared **Platform Master** that can power multiple product shells:

1. **The Concierge** — business/client-facing product for independent professionals.
2. **Owner PA / Brave PA** — chief-of-staff layer for the owner.
3. **Standalone Brave PA** — future consumer/personal product.
4. **Creator PA / REELAI** — future creator workflow shell.
5. **Food Finder / Event Finder / Travel modules** — context modules that can live inside PA first, then become standalone if usage proves demand.

The correct immediate commercial strategy is:

> Build a sellable The Concierge + Owner PA + Growth Engine version first, while architecting it so approximately 80% of the PA intelligence can later become Standalone Brave PA.

Do not build the full standalone PA, full Food Finder, full Event Finder, full REELAI, full Living Star, or native apps before revenue.

---

## 1. Pricing decisions already locked

This document does **not** redesign pricing from scratch.

Use the pricing and percentages already defined in `BraveByBruno_COMPLETE_VISION_v1.md` and Master Document v5.0:

| Plan | Locked price | Notes for this implementation plan |
|---|---:|---|
| Free Trial | £0 / 14 days | Full Starter+ experience during trial. |
| Starter | £19/mo | Entry product: public page, Concierge, dashboard, lead capture, 2D Basic. |
| Starter+ | £24/mo | Adds 3D Signature fixed theme + 4-5 World Credit taster. |
| Pro | £29/mo | Main value plan: all modes, Growth Engine, reminders, package builder, 3D Signature, 1 bespoke World once at onboarding, 15 World Credits/mo, 1 profile. |
| Pro Multi | £53/mo | One person with several identities, 3 profiles, no agency/client-management features. |
| Agency | From £79/mo | Businesses managing other people's presences; up to 10 profiles; pricing above 3 profiles remains unresolved and needs volume-discount curve. |
| Enterprise | Custom | White label, API access, SLA, dedicated support. |

3D economics already locked as planning estimates:

- 1 World Credit approximately £0.30 raw cost.
- Full bespoke World, 6-8 elements, approximately £3.20 raw cost.
- Top-up pack: +25 World Credits for £16.
- Extra Profile 3D Creation: £9 flat, separate from edit credits.
- These are June 2026 estimates and must be validated with 5-10 real generations before launch pricing locks.

Stripe/platform fee logic already locked:

- Month 1-3: 0% platform fee.
- Month 4-6: 0.75% platform fee.
- Month 7-12: 0.50% platform fee.
- Above £5,000/month volume: 0.30% platform fee.
- Stripe processing fee remains separate and charged to the professional.

Therefore the work now is not pricing invention. It is:

1. Adjust feature access per plan.
2. Estimate build cost before revenue.
3. Estimate monthly burn before and after first users.
4. Estimate variable cost per user/action.
5. Define what must be funded by credits, customers, loan, or pre-sales.

---

## 2. Updated product structure

### 2.1 Shared Platform Master

The long-term architecture should be one shared intelligence layer, not separate bots.

```text
Shared Platform Master
├─ Personal & Professional Maximisation Engine
├─ Adaptive Chief of Staff / Executive Function Engine
├─ Situational Life & Opportunity Engine
├─ Growth & Market Intelligence Engine
├─ Intent Router / Capability Router
├─ Memory / Learning / Outcome Ledger
├─ Calendar / Task / Notification Core
├─ Professional Context Engine
├─ ActionGate
├─ Structured Booking / Payment Engine
├─ Food / Event / Travel Modules
├─ Relevant Visual Engine
├─ Living Star / Generative Universe
├─ AI Fleet / Provider Router
└─ Product shells
   ├─ The Concierge
   ├─ Standalone Brave PA
   ├─ Creator PA / REELAI
   ├─ Food Finder
   ├─ Event Finder
   ├─ Performers
   └─ Brave Social
```

### 2.2 First-class capability families to add to the roadmap

The earlier strategy documents covered many technical capabilities, but the current vision requires these additional first-class capability families.

#### A. Personal & Professional Maximisation Engine

Purpose:

> Help the owner become the strongest, most sustainable version of themselves personally and professionally.

It should maximise:

- quality of life;
- calm;
- income;
- consistency;
- personal brand;
- professional growth;
- client conversion;
- creative output;
- executive function;
- realistic time use;
- recovery and sustainability.

Core principle:

> The assistant should know when to adapt itself around the owner, and when to gently suggest changes to the owner's routine, behaviour, structure, offers, timing, or priorities.

#### B. Adaptive Chief of Staff / Executive Function Engine

Purpose:

> Help the owner use their real brain as a superpower while reducing friction from executive-function gaps, time blindness, overcommitment, scattered attention, and non-linear routines.

Must include:

- Time Reality Engine;
- personalised buffer/reminder logic;
- Transition Support;
- Adaptive Routine Planning;
- Energy-Aware Scheduling;
- Anti-Overcommitment Checks;
- Personal + Work Calendar Fusion;
- Opportunity Fit Suggestions;
- Strength Amplification;
- Friction Reduction;
- End-of-Day Reflection;
- learning from which plans work or fail.

Legal boundary:

- This is productivity, planning, routine, and organisational support.
- It must not diagnose, treat ADHD, replace therapy, or provide medical advice.
- ADHD-style support is framed as user-controlled personalisation: reminders, buffers, routine adaptation, and friction reduction.

#### C. Situational Life & Opportunity Engine

Purpose:

> Convert messy human needs into realistic scenarios that balance energy, calendar, budget, income goals, clients, personal preferences, and business opportunities.

Example owner request:

> "I am tired. I need a holiday."

Expected behaviour:

1. Understand recovery need.
2. Ask only missing constraints: budget, dates, duration, desired intensity.
3. Check calendar and commitments.
4. Check income target and cash gap.
5. Identify movable clients or flexible obligations.
6. Suggest income sprint through packages/offers/follow-ups/content.
7. Generate 3-5 scenario cards.
8. Execute no real action without approval.
9. Learn which scenario worked.

Example outputs:

- Leave now, cheaper and simpler.
- Leave in 3 weeks after a £450 income sprint.
- Move Client A to their second-choice date and free 5 days.
- Take a 2-day recovery block instead of full travel.
- Choose a city/event that also creates content/networking value.

#### D. Growth & Market Intelligence Engine

Purpose:

> Understand the business, offer, market, audience, pricing, content capacity, acquisition paths, and conversion gaps; then turn those into growth actions.

Submodules:

- Product Understanding Engine;
- Market Access / Research Engine;
- Competitor / Market Comparison Engine;
- Pricing Gap and Price Elasticity Suggester;
- Offer / Package Builder;
- Growth Roadmap Generator;
- Content and Campaign Generator;
- Lifestyle Quiz / Acquisition Funnel;
- Abandoned Enquiry Recovery;
- Post-Session Follow-up;
- Automated Review Engine;
- Referral Engine;
- Appointment Gap Filler;
- Contextual Upsell Engine;
- Partnership / Cross-referral Network;
- Outcome Learning;
- Cohort Benchmark Intelligence, legal-gated.

#### E. Learning Onboarding / Personal Universe Setup

The existing onboarding is not wrong. It is incomplete.

Current onboarding becomes:

> Business & Public Profile Setup.

The new onboarding becomes:

```text
Personal Universe Setup
├─ Business & Public Profile Setup  (existing onboarding)
├─ Owner Operating Profile          (routine, energy, executive function)
├─ Growth & Market Profile          (income, offers, market, goals)
├─ PA Behaviour Setup               (style, challenge level, proactivity)
├─ Memory & Learning Consent        (what may be remembered)
└─ Visual Universe Setup            (theme, 3D world, Living Star later)
```

Existing onboarding fields that must be preserved:

- professional modes;
- name, business name, city, tagline;
- tone;
- writing example;
- services/acts/packages;
- performer details: showreel, equipment, fee, availability;
- creator details: platforms, niche, collab packages, media kit;
- B&B / property details;
- sensitive topics;
- contact links: booking, calendar, Instagram, WhatsApp, Telegram, phone, email, agent;
- location and coverage;
- photos/media/gallery/video;
- legal forms;
- handle/share link;
- account creation.

New onboarding fields required:

- real preparation time;
- reminder style;
- time blindness / buffer preferences;
- energy patterns;
- routine style;
- transition support;
- overcommitment risk;
- dietary preferences and restrictions;
- food/event/travel preferences;
- budget comfort zones;
- recovery preferences;
- personal goals;
- professional goals;
- income targets;
- services to sell more;
- services that drain energy;
- ideal client and wrong-fit client;
- content capacity;
- growth style;
- PA style: direct, calm, challenging, gentle, strategic, creative;
- proactive permission levels;
- memory permissions;
- connected-tool permissions;
- what the PA should never suggest.

---

## 3. Revised tier access

### 3.1 Starter — £19/mo

Purpose:

> Basic client-facing Concierge and public business profile.

Included:

- 1 mode;
- public page;
- real AI client Concierge;
- dashboard;
- lead capture;
- basic notifications;
- CSV export;
- services/prices/media/contact links;
- 2D Basic theme;
- basic onboarding;
- limited Growth hints, not full Growth Engine;
- no advanced Owner PA;
- no bespoke 3D World except trial/frozen logic.

Cost strategy:

- Must remain low-AI-cost.
- Do not include heavy search, scenario planning, or deep PA.
- Use this as entry and upgrade path.

### 3.2 Starter+ — £24/mo

Purpose:

> Starter plus visual hook.

Included:

- Starter;
- 3D Signature fixed theme;
- 4-5 World Credit taster/month;
- no bespoke ongoing full-generation loop;
- no full Owner PA;
- no full Growth Engine.

Cost strategy:

- Taster should be enough for one meaningful edit, not enough to cannibalise Pro.
- 3D generation must remain credit-limited.

### 3.3 Pro — £29/mo

Purpose:

> Main plan where the product becomes indispensable for the owner.

Included:

- all modes;
- unlimited conversations subject to fair-use/cost guardrails;
- Owner PA / Chief of Staff light;
- Growth Engine MVP;
- package builder;
- reminders;
- email campaigns;
- scenario planning light;
- Learning Onboarding / Personal Universe Setup;
- proactive suggestions within limits;
- 3D Signature theme;
- 1 bespoke World once at onboarding;
- 15 World Credits/month;
- 1 profile.

This is the plan that should carry the real product value.

### 3.4 Pro Multi — £53/mo

Purpose:

> One person with several identities.

Included:

- Pro;
- 3 profiles;
- one shared owner account;
- profile switching;
- shared PA understanding owner identity across professional identities;
- shared World Credit pool;
- no agency/team/client-management features.

Cost strategy:

- Margin is thinner due to multiple profile creations.
- Still cheaper and clearer than buying three Pro accounts.

### 3.5 Agency — from £79/mo

Purpose:

> Businesses managing other people's presences.

Included:

- up to 10 profiles in initial concept;
- team permissions;
- client management;
- white label;
- custom domain;
- priority support;
- referral system;
- higher Growth Engine and Market Intelligence access;
- volume usage controls.

Open pricing issue:

- Pricing above 3 profiles needs a dedicated volume-discount curve.
- Do not guess this in implementation.

### 3.6 Standalone Brave PA — future product, not pre-revenue build

Purpose:

> Consumer/personal product based on the same Platform Master.

Positioning:

- The standalone PA is not a new brain.
- It is the consumer shell around the same PA intelligence used by The Concierge.
- Approximately 80% of the intelligence should be shared with Owner PA.

Possible future tiers:

| Future tier | Possible position | Notes |
|---|---|---|
| Brave PA Free | Basic planning and manual inputs | Not for immediate build. |
| Brave PA Plus | Calendar, reminders, preferences, food/events suggestions | Requires context providers. |
| Brave PA Pro | scenario planning, travel, budget, routine, proactive alerts | Requires memory, calendar, ActionGate. |
| Creator PA / REELAI add-on | content/reels/collabs | Later shell. |

Do not build this as a separate product before The Concierge has revenue or pilot customers.

---

## 4. Build phases, timing, costs, dependencies, blockers

### Phase 0 — Reality baseline and planning lock

**Timing:** 2-4 days  
**Pre-revenue cash:** £0-£50  
**Monthly burn impact:** none  
**Build type:** docs / audit / planning only

Deliverables:

- Save this document.
- Save cross-reference document for missing engines.
- Update feature registry planning notes.
- Confirm Day 6 PA web search is represented.
- Confirm old onboarding fields are preserved in the new onboarding plan.
- Confirm pricing is locked and only feature access/unit economics are changing.

Blockers:

- None technical.

Do not touch:

- production code;
- schema;
- auth;
- Stripe;
- provider keys;
- live feature flags.

### Phase 1 — Launch-truth and infrastructure hardening

**Timing:** 1-2 weeks  
**Pre-revenue cash:** £100-£500  
**Monthly burn impact:** £30-£150 depending hosting plan choices

Deliverables:

- Server-owned prompt templates.
- Live notification email verification.
- Notification preferences.
- Sound settings.
- Basic push plan, not necessarily full push implementation.
- Stripe live/test config verification.
- Supabase media bucket verification.
- Render free/starter/standard decision.
- Cost guardrails for AI/search.

Dependencies:

- existing feature flags;
- existing audit events;
- existing notification work;
- existing auth green status.

Blockers:

- If server still trusts client-supplied system prompts, PA/Growth personalisation remains risky.
- If notification email is not live-verified, no owner-notified claim should be shown.
- If Stripe is not live/test verified, no end-to-end payment claim.

External cost notes:

- Vercel Hobby is $0/mo; Vercel Pro is $20/mo and includes $20 usage credit.
- Render Web Services include Free, Starter at $7/mo, Standard at $25/mo.
- Render Cron Jobs start from $1/mo.
- Supabase Pro is $25/mo and includes 100,000 monthly active users, 8 GB disk, 250 GB bandwidth, and $10/mo compute credits.

### Phase 2 — Learning Onboarding / Personal Universe Setup v1

**Timing:** 2-3 weeks  
**Pre-revenue cash:** £200-£900  
**Monthly burn impact:** mostly AI usage; no heavy fixed cost if done carefully

Deliverables:

- Keep existing onboarding as Business & Public Profile Setup.
- Add Personal Universe Setup structure.
- Add Quick Setup vs Deep Setup flow.
- Add PA behaviour preferences.
- Add owner operating profile questions.
- Add Growth profile questions.
- Add diet/food/event/travel preferences as optional fields.
- Add memory/learning permission scaffolding, even if memory itself is dormant.
- Store data in structured profile_data sections first, unless a dedicated schema is approved later.

Implementation principle:

> Do not ask everything in one long form. Use progressive conversational onboarding.

Suggested data structure inside profile_data first:

```text
profile_data.business
profile_data.public_profile
profile_data.services
profile_data.media
profile_data.legal
profile_data.owner_operating_profile
profile_data.growth_profile
profile_data.pa_preferences
profile_data.life_preferences
profile_data.memory_permissions
profile_data.visual_world
```

Blockers:

- Needs server-owned prompt templates if the PA is going to use the data safely.
- Needs privacy copy update before claiming memory/learning.
- Needs view/edit/delete roadmap before active personal memory.

### Phase 3 — Owner PA / Adaptive Chief of Staff v1

**Timing:** 2-4 weeks  
**Pre-revenue cash:** £300-£1,200  
**Monthly burn impact:** £50-£250 depending usage

Deliverables:

- Daily Planning Engine v1.
- Time Reality Engine v1.
- Reminder style and buffer logic.
- Anti-overcommitment warnings.
- Energy-aware scheduling suggestions.
- Morning plan and evening reflection.
- Work/personal commitment distinction.
- Manual calendar import or read-only calendar connection depending readiness.
- Owner-only responses, not public client surface.

MVP examples:

- "You have worked a lot this week; choose low-friction admin tonight."
- "This plan is theoretically possible but not realistic with your usual prep/travel time."
- "You have an 8pm fixed appointment; a 10pm show is realistic only if you eat nearby and do not go home first."

Blockers:

- Calendar connector if live calendar reading is expected.
- Personal memory controls if PA claims it remembers patterns across sessions.

### Phase 4 — Growth & Market Intelligence Engine MVP

**Timing:** 3-5 weeks  
**Pre-revenue cash:** £400-£1,500  
**Monthly burn impact:** £50-£300 depending search/research volume

Deliverables:

- Product Understanding Engine.
- Offer/package builder.
- Income target and cash-gap planner.
- Services-to-push logic.
- Follow-up suggestions.
- Content/campaign generator.
- Market/search-backed research with citations or source links.
- Growth Roadmap Generator.
- Basic outcome logging.

MVP examples:

- "To make £450 before your trip, sell 3 x 90-minute packages or 6 x 60-minute sessions."
- "This service has high margin but low visibility; push it this week."
- "These 5 leads should be followed up before creating more content."

Blockers:

- Needs clear separation between known business data, inferred estimate, and live research.
- Needs rate limits for web search.
- Needs no fake claims about market certainty.

### Phase 5 — Situational Life & Opportunity Engine v1

**Timing:** 3-6 weeks  
**Pre-revenue cash:** £500-£2,000  
**Monthly burn impact:** £100-£400 depending external providers

Deliverables:

- Scenario Planner.
- Calendar window logic.
- Budget/income gap planner.
- Workload/recovery suggestions.
- Food/event/travel preference matching, initially lightweight/search-based.
- Visual scenario cards.
- Action proposals but no autonomous execution.

MVP examples:

- "I need a holiday" -> 3 travel/recovery/income scenarios.
- "I finished work and have an 8pm appointment" -> food/time/diet options.
- "Can I go to a 10pm show after this?" -> route/time/buffer answer.

Dependencies:

- Owner Operating Profile;
- Growth Engine;
- Calendar/commitment model;
- Context/search provider;
- ActionGate for real actions.

Blockers:

- Food/menu/restaurant data can be costly or unreliable if done deeply.
- Do not build full Food Finder here. Build only PA-context food recommendations first.

### Phase 6 — Calendar, reminders, tasks, proactive notifications

**Timing:** 3-5 weeks  
**Pre-revenue cash:** £300-£1,500  
**Monthly burn impact:** £20-£150 initially

Deliverables:

- Google Calendar read-only.
- Calendar write only after approve/confirm.
- Task/reminder system.
- Morning briefing.
- Evening reflection.
- Follow-up reminders.
- Push notifications after proof on real device.
- Digest/clustering scheduler.

Dependencies:

- Google OAuth credentials;
- Render cron or equivalent scheduler;
- push VAPID keys;
- notification preferences.

Blockers:

- OAuth privacy disclosure.
- Connected Assistant consent per connector.
- No calendar claims until live connector works.

### Phase 7 — ActionGate + booking/payment engine v1

**Timing:** 4-8 weeks  
**Pre-revenue cash:** £500-£2,500  
**Monthly burn impact:** mostly Stripe/payment and AI usage

Deliverables:

- Action proposal model.
- Approve/deny/revise flow.
- Calendar write after approval.
- Booking request -> deposit/payment link -> confirmation flow.
- Deposit logic.
- Refund/cancellation policy display.
- Audit events for every consequential action.

Blockers:

- Stripe live/test config.
- Payment terms.
- Calendar write.
- Clear rollback/repair path.

### Phase 8 — Relevant Visual Engine and visual cards

**Timing:** 2-4 weeks for cards; 6-12 weeks for full Living Star  
**Pre-revenue cash:** £200-£1,000 for cards; £3,000-£12,000+ for full Living Star path  
**Monthly burn impact:** low for cards; variable/high for 3D generation

Deliverables for pre-revenue:

- Scenario cards.
- Budget gap bar.
- Calendar window card.
- Growth action checklist.
- Food/event option cards.
- 3D credit validation with 5-10 real generations.

Do not build pre-revenue unless funded:

- full Living Star;
- full conversational 3D world editing;
- native app visual system;
- advanced generative core.

Blockers:

- 3D API/provider cost validation.
- Gemini visual review.
- Quality benchmark.

### Phase 9 — Standalone PA shell planning, not build

**Timing:** 1 week planning  
**Pre-revenue cash:** £0-£200  
**Monthly burn impact:** none if planning only

Deliverables:

- Standalone Brave PA product shell document.
- Shared intelligence map showing 80% reuse from Owner PA.
- Food/Event/Travel modules as PA modules first.
- No separate app/product build yet.

Trigger to build:

- At least 10-20 paying Concierge users or external funding.
- Evidence that owner-side PA usage is high.

### Phase 10 — First paid beta / founding customers

**Timing:** starts as soon as Phase 2-4 are credible  
**Pre-revenue cash:** minimal  
**Revenue goal:** £1,000-£3,000 cash through founding beta

Recommended beta offer:

- £99-£199 setup for first selected customers.
- £29/mo Pro after setup or after trial.
- Limit beta to 10-20 people.
- Manual high-touch onboarding by Bruno allowed; do not pretend everything is fully automated yet.

Expected revenue:

- 10 beta users x £149 setup = £1,490 upfront.
- 10 beta users x £29/mo = £290/mo recurring.
- 20 beta users x £149 setup = £2,980 upfront.
- 20 beta users x £29/mo = £580/mo recurring.

This is faster and safer than waiting for grants.

---

## 5. Pre-revenue budget scenarios

### Scenario A — Lean sellable version

**Goal:** first paying beta with real value, not full vision.  
**Timing:** 4-8 weeks.  
**Cash needed before selling:** £750-£2,000.  
**Monthly burn:** £100-£300.

Includes:

- trust fixes;
- improved onboarding;
- Owner PA light;
- Growth Engine MVP;
- web search;
- dashboard/basic notification verification;
- Stripe/Resend/Supabase/Render verification;
- no full Living Star;
- no standalone PA;
- no full Food/Event modules;
- no native apps.

This is the recommended path.

### Scenario B — Strong beta / investor-demo version

**Goal:** product feels serious, personalised, and differentiated.  
**Timing:** 8-12 weeks.  
**Cash needed before selling:** £2,000-£6,000.  
**Monthly burn:** £200-£600.

Adds:

- scenario planning;
- read-only calendar;
- proactive reminders;
- visual cards;
- first food/event/travel contextual suggestions;
- stronger Growth Engine;
- stronger PA onboarding.

This is ideal if funded by credits/pre-sales.

### Scenario C — Wow demo / early funding showcase

**Goal:** very impressive demo with Living Star direction and PA standalone vision visible.  
**Timing:** 3-6 months.  
**Cash needed before selling/funding:** £6,000-£15,000+ if mostly built with AI/Claude Code and minimal external help.  
**External team version:** £30,000-£100,000+.

Adds:

- deeper visual/generative layer;
- early Living Star;
- advanced scenario planning;
- stronger memory controls;
- broader provider integrations;
- polished standalone PA shell.

Do not choose this path without funding, credits, or pre-sales.

---

## 6. Monthly burn model

### 6.1 Before customers

| Cost line | Lean target | Strong beta target | Notes |
|---|---:|---:|---|
| Vercel | $0-$20/mo | $20/mo | Pro only if needed for commercial/team/limits. |
| Render backend | $7-$25/mo | $25-$85/mo | Starter or Standard likely enough early. |
| Render cron | $1+/mo | $1-$20/mo | Only when proactive/digest exists. |
| Supabase | $0-$25/mo | $25/mo | Free until real scale/security requires Pro. |
| Anthropic/Claude API | £20-£150/mo | £100-£400/mo | Depends heavily on usage and model selection. |
| Search/research | £10-£100/mo | £50-£250/mo | Use Brave/Perplexity carefully. |
| Email/Resend | low | low/medium | Verify live env. |
| Maps/food/event | £0-£50/mo | £50-£300/mo | Avoid deep provider dependency pre-revenue. |
| 3D generation tests | £20-£200 one-off | £100-£500 one-off | Validate unit economics. |
| Legal/privacy light review | optional £500-£2,000 | recommended £500-£2,000 | Especially before sensitive/memory claims. |

### 6.2 After first users

| Stage | Approx recurring revenue | Target monthly infra/API burn | Comment |
|---|---:|---:|---|
| 5 Pro users | £145/mo | <£150/mo | Nearly covers lean API/hosting if controlled. |
| 10 Pro users | £290/mo | <£300/mo | Break-even target from master doc. |
| 20 Pro users | £580/mo | £300-£600/mo | Funds continued build. |
| 50 Pro users | £1,450/mo | £500-£1,200/mo | Can fund maps, better memory, visual tests. |
| 100 mixed users | £2,000-£4,000+/mo | £1,000-£2,500/mo | Start considering more providers/human help. |

---

## 7. Unit economics assumptions

These are planning estimates, not audited financials.

### 7.1 AI conversation cost

Legacy Master Document estimated approximately £0.02 per conversation for the early V1. Keep this as a target, not a guarantee.

Current provider reality:

- Claude Sonnet 4.6 / 4.5 pricing is materially higher than Haiku.
- Haiku should be used for simple classification, extraction, and low-risk routing where quality is sufficient.
- Sonnet should be reserved for owner PA reasoning, Growth Engine, scenario planning, and higher-value outputs.

Planning target:

| Use case | Suggested model tier | Target cost logic |
|---|---|---|
| client FAQ / simple service questions | Haiku or cheaper fast model | very low cost |
| lead scoring / intent classification | Haiku / rules first | very low cost |
| Growth plan / scenario planner | Sonnet or stronger reasoning | higher cost, Pro-only |
| market research | search + summariser | Pro/Agency only or capped |
| 3D prompt generation | specialised model/provider | credit-gated |

### 7.2 Search/research cost

Search should be treated as a metered capability.

Rules:

- Starter: no unrestricted live research.
- Starter+: no unrestricted live research.
- Pro: capped research/search for Growth and PA.
- Agency: larger allowance.
- Deep research: manual or add-on, not included unlimited.

### 7.3 3D cost

3D must remain credit-based.

Rules:

- Trial and Starter+ taster remain small.
- Pro includes enough credits to feel meaningful, not unlimited.
- Extra Profile 3D Creation remains separate fee.
- Agent-Tailored Worlds are quoted from £59, not automated flat work.
- Validate with 5-10 real generations before locking production margins.

### 7.4 Maps/Food/Event cost

Food/Event/Travel recommendations can become expensive if based on paid Places/Routes APIs at scale.

Rules:

- Pre-revenue: lightweight search + manual preferences + simple maps only.
- Pro beta: limited context suggestions.
- Standalone PA/Food/Event products: only after demand is proven or external funding is secured.

---

## 8. Funding options and immediate order

### 8.1 Best immediate non-dilutive path: startup/cloud credits

Apply first because these reduce API/cloud pressure without debt.

#### Google for Startups Cloud Program

Why relevant:

- AI-first startup angle.
- Credits can support cloud, AI, data, and infrastructure experiments.

Current public offer snapshot:

- Pre-funded startups: $2,000 credits to build MVP.
- Seed-Series A: up to $200,000, or up to $350,000 for AI-first startups.

Action:

- Prepare one-page pitch.
- Emphasise AI concierge, personal/professional operating system, Growth Engine, AI-first startup, London founder, live MVP.

#### AWS Activate / AWS for Startups

Why relevant:

- AWS credits can offset infrastructure, AI/ML, and Bedrock model costs.

Current public offer snapshot:

- Up to $200,000 AWS Activate Credits.
- AI startups may be eligible for additional credits.

Action:

- Apply after/alongside Google.
- If using Bedrock later, credits can matter.

#### NVIDIA Inception

Why relevant:

- AI/3D/generative direction.
- Partner offers, cloud credits, training, investor exposure.

Current public offer snapshot:

- Free program.
- No membership fees, application fees, or equity requirements.
- Requires incorporated company, working website, and at least one developer.

Action:

- Apply after incorporation decision, if eligible.
- Strong fit for Living Star / generative world pitch.

### 8.2 Cash path: UK Start Up Loan

Why relevant:

- Cash, not credits.
- Can fund the stronger beta version if Bruno chooses debt.

Current public terms snapshot:

- Borrow up to £25,000.
- Fixed interest rate 7.5% per year.
- Repayment over 1-5 years.
- 12 months free mentoring.
- No arrangement/early-repayment fees.

Caution:

- This is debt.
- Do not take it before checking personal affordability, benefits/admin implications, and realistic repayment plan.

### 8.3 Fastest practical funding: pre-sales / founding beta

Recommended.

Offer:

- Founding beta setup: £99-£199.
- Monthly Pro: £29/mo.
- Limited to 10-20 professionals.
- High-touch onboarding by Bruno included.
- Honest beta language: not full automation, but a concierge/PA system being personalised.

Why best:

- Cash arrives faster than grants.
- Validates demand.
- Funds the next features.
- Creates real user data/outcomes.

### 8.4 Grants / accelerators

Use as parallel track, not immediate cash dependency.

Potential categories:

- AI startup accelerators;
- creator economy accelerators;
- accessibility/neurodiversity productivity tools;
- solo-professional / SME productivity;
- London startup support;
- Innovate UK-style grants if fit.

Do not block the build waiting for grants.

---

## 9. Recommended immediate funding plan

### Week 1

1. Save this plan.
2. Create 1-page funding/credits application brief.
3. Prepare screenshots/live demo proof.
4. Apply to Google for Startups Cloud.
5. Apply to AWS Activate.
6. Decide whether incorporation is needed before NVIDIA Inception.
7. Build founding beta waitlist/offer.

### Week 2

1. Contact 10-20 possible beta professionals.
2. Offer high-touch onboarding.
3. Charge setup fee only when the onboarding/product promise is honest.
4. Keep monthly burn below £300.
5. Begin Phase 1 and Phase 2 implementation only after docs are locked.

---

## 10. First implementation packet after this document

Do not start the full plan at once.

Next implementation packet should be:

```text
Learning Onboarding / Personal Universe Setup — Audit and Design Packet
```

Scope:

- Audit existing `components/Onboarding.jsx` and `lib/constants.js`.
- Map every existing field to its new role.
- Propose new data sections inside `profile_data`.
- No UI implementation yet.
- No schema changes yet.
- No provider changes.
- No auth/routing/Stripe changes.
- Output a field map and build packet.

Why this first:

- The PA, Growth Engine, Situational Engine, Food/Event/Travel suggestions, and future Standalone PA all depend on knowing the person and business.
- Onboarding is the root of personalisation.
- The old onboarding is valuable and must be preserved, not replaced blindly.

---

## 11. Blockers to track before selling

Hard blockers for paid launch:

1. Public Concierge must be real on real pages.
2. AI consent path must be real and truthful.
3. Owner notification claims must be backed by backend event/delivery.
4. Stripe must be verified before payment claims.
5. Server-owned prompt authority must be in place before deep PA personalisation is relied on.
6. Personal memory cannot be claimed unless opt-in + view/export/delete controls exist.
7. 3D credit costs must be validated with real generation tests.
8. Growth Engine outputs must distinguish facts, estimates, and live research.
9. Standalone PA must not be marketed as ready until modules exist beyond The Concierge shell.
10. Agency pricing above 3 profiles remains open.

Soft blockers:

- legal/privacy review for sensitive data and memory;
- maps/places provider choice;
- calendar OAuth readiness;
- push notification real-device proof;
- Render hosting upgrade decision;
- funding/credit applications.

---

## 12. Decision summary

Build before selling:

- trust fixes;
- onboarding expansion;
- Owner PA light;
- Growth Engine MVP;
- scenario planning light;
- search-backed research with caps;
- basic visual cards;
- Stripe/notification verification;
- 3D cost validation tests.

Do not build before selling unless funded:

- full standalone PA;
- full Food Finder;
- full Event Finder;
- full REELAI;
- full Living Star;
- native apps;
- wearable interface;
- voice/video learning;
- advanced cohort intelligence;
- full AI Fleet.

Main recommendation:

> Target a £750-£2,000 lean sellable version first, or £2,000-£6,000 strong beta if startup credits/pre-sales land. Use founding beta customers to fund expansion. Keep the same Platform Master underneath so The Concierge, Owner PA, Standalone PA, Food/Event/Travel, and Creator PA do not become disconnected products.

---

## 13. Next safest action

Create one documentation/audit packet for Claude Code:

```text
AI: Claude Code
Mode: Repo Documentation / Audit Only
Objective: Audit the current onboarding and produce the field map for Learning Onboarding / Personal Universe Setup.
Do not implement UI.
Do not change schema.
Do not touch auth, routing, Stripe, provider keys, backend runtime, or feature flags.
Output expected: existing fields, retained fields, new proposed fields, data structure, blockers, and next implementation packet.
```

Stop after that audit before implementing.
