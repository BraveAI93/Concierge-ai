# Brave by Bruno / The Concierge
## Core Architecture, Brain Spine & Fleet Integration Decision Pack — v0.2
### Consolidated plan for Claude Code, ChatGPT, Gemini, and Perplexity

**Date:** 8 July 2026  
**Prepared by:** ChatGPT, consolidating Bruno’s question, Claude’s latest response, Gemini’s caution, and Perplexity’s Fleet Architecture direction.  
**Purpose:** Decide *when* and *how* to integrate the future Brain Spine / Fleet Architecture into The Concierge without destabilising V1.

---

## 0. Current State

### 0.1 AUTH GREEN

`AUTH GREEN` is now accepted as passed based on the latest Claude Code report and screenshots.

Confirmed evidence reported by Claude Code:

- `8fa6229` — middleware/cookie fix: middleware now reads `cai_token` instead of `ownerToken`.
- `f0673a1` — `_directToken` bypass fix: `/api/auth/login` now requires an existing valid `cai_token` before honouring `_directToken`.
- Both commits are pushed to `origin/main`.
- Working tree is clean.
- Real backend login returned `Set-Cookie: cai_token`.
- Cookie was `HttpOnly`.
- `cai_token` matched returned JSON token.
- Anonymous `_directToken` injection is blocked with `401`.
- Backend build is clean.
- Production is live on `bravebybruno.com`.
- `/theconcierge/dashboard` without cookie redirects to `/theconcierge/owner-auth`.

Loose ends noted but not blocking:

- Empty/missing backend token scenario still unverified.
- Local Supabase service-role key is broken for cleanup.
- Touch-scroll IDE corruption remains an operational risk.

**Decision:** Auth is closed enough to move on. Do not keep working on auth unless a real regression appears.

---

## 1. Why This Document Exists

Bruno asked a critical strategic question:

> If the future product core will be a Brain Spine / Fleet Architecture, does it make more sense to build/fix that core now and then rebuild functions on top of it, instead of polishing features that may need to be reworked later?

This is the correct question.

The answer is not binary. The correct path is to evaluate three levels separately:

1. **Full Fleet Brain** — long-term target architecture.
2. **Minimal Brain Spine** — possible near-term core abstraction.
3. **Current V1 features** — product and visual experience that still need to launch.

The goal is to avoid two opposite mistakes:

- **Mistake A:** Keep building features on a weak or fragmented AI core, then rebuild everything later.
- **Mistake B:** Over-engineer a full multi-provider brain now and destroy V1 momentum.

---

## 2. Distinction Between Three Concepts

### 2.1 Multi-AI Round Table — already active

This is the *team workflow* around Bruno:

- ChatGPT prepares strategy, packets, QA, prompts.
- Gemini reviews visual/UX direction.
- Perplexity researches external evidence.
- Claude / Claude Code handles repo verification and execution.
- Bruno decides.

This already works and is governed by the Product Operating System / AI Technical Round Table Protocol v0.5.

This does **not** mean the product runtime itself already has multiple AI engines.

---

### 2.2 Minimal Brain Spine — possible near-term architecture

This is the lightweight internal structure that future AI features would plug into.

It does **not** require adding new providers immediately.

It may include:

- common request envelope;
- model adapter interface;
- memory/provenance interface;
- verification/audit interface;
- tool/action-gate interface;
- final response synthesis interface;
- the current existing AI model still used underneath.

The purpose is to create a stable spine so future Concierge, Brave PA, Manager Agent, Trust Dot, memory, and verification features are not hardcoded separately.

---

### 2.3 Full Fleet Brain — V2/V3 target

This is the long-term premium architecture:

- one Concierge identity in front;
- one policy/control plane;
- one platform-owned memory spine;
- specialist role-based lanes;
- MCP/tool fabric;
- provider adapters;
- observability and evaluations;
- cost/latency routing;
- verification modes;
- action permissions.

This should **not** be implemented now.

---

## 3. Layered Architecture Proposal

### Layer 0 — Existing Stable Core

Current stack:

- Next.js 14 on Vercel.
- Go backend on Render.
- Supabase database.
- Claude/Anthropic server-side AI path.
- httpOnly cookie auth via `cai_token`.
- Product path under `/theconcierge`.
- Public vanity slugs at root.

This layer must remain stable.

---

### Layer 1 — Operating Governance

Already defined by v0.5:

- AI Technical Round Table.
- Product Operating System.
- Build, launch, and post-launch pipeline.
- Anti-chaos rules.
- Mode headers.
- One active integration at a time.
- Production verification over summaries.

This layer governs all future work.

---

### Layer 2 — Minimal Brain Spine

Potential near-term addition.

Suggested internal modules, subject to Claude Code audit:

```text
ai/
  orchestrator/
  model_adapters/
  memory/
  verification/
  audit/
  tools/
```

or repo-specific equivalent.

The spine should initially wrap the current AI behaviour without changing it.

Minimum interfaces:

1. `RequestEnvelope`
   - user ID;
   - session ID;
   - profile ID;
   - mode;
   - permissions;
   - risk level;
   - source route/component;
   - conversation context.

2. `ModelAdapter`
   - standard call interface;
   - provider-agnostic input/output;
   - error handling;
   - latency/cost metadata.

3. `MemoryInterface`
   - profile data;
   - preferences;
   - conversation history;
   - future vector/memory retrieval hook.

4. `VerificationInterface`
   - deterministic checks;
   - unsupported claim flag;
   - action safety flag;
   - future Trust Dot compatibility.

5. `ActionGate`
   - calendar/payment/message/tool permissions;
   - confirm-before-action logic;
   - high-risk escalation.

6. `AuditLogInterface`
   - request ID;
   - model used;
   - evidence used;
   - action taken;
   - verification status;
   - failure reason.

7. `ResponseSynthesisInterface`
   - final Concierge voice;
   - consistent persona;
   - user-facing formatting.

---

### Layer 3 — Lean Brain

Possible after Minimal Brain Spine.

Features:

- one main model still active;
- memory retrieval from Supabase;
- simple audit logs;
- deterministic verification rules;
- basic request modes: Fast / Smart / Verified;
- no full multi-provider routing yet.

This is the first real runtime upgrade.

---

### Layer 4 — Smart Concierge

Later phase.

Features:

- planner mode;
- selective external research;
- preference learning;
- first-pass verification;
- richer proactive suggestions;
- tool/action gating.

---

### Layer 5 — Full Fleet Brain

V2/V3 target.

Features:

- Claude/OpenAI/Gemini/Perplexity role-based lanes;
- OCR/document extraction lane;
- transcription lane;
- cheap classifier/router lane;
- retrieval/reranking lane;
- safety/moderation lane;
- MCP/tool gateway;
- evaluation dashboards;
- provider registry;
- cost-quality routing.

Rejected for immediate implementation.

---

## 4. Roadmap Options for Claude Code To Evaluate

### Option A — Visual First

**Sequence:**

1. AUTH GREEN.
2. Save v0.5 protocol and architecture docs.
3. Cinematic Shell.
4. Landing restoration.
5. Core Flow QA.
6. Dashboard visual integration.
7. Trust/security hardening.
8. Brain Spine audit and implementation later.

**Pros:**

- Faster visible progress.
- Product looks premium sooner.
- Lower backend architecture risk today.
- Good if current AI code is small and not about to be heavily expanded.

**Cons:**

- AI functions may continue growing on old structure.
- Possible future rework.
- Could delay the foundation for Trust Dot, Manager Agent, memory, and verification.

**Choose this if:**

- Claude finds current AI flow too tangled for a safe spine refactor now.
- Minimal Brain Spine would require schema changes or runtime behaviour changes.
- Cinematic Shell is urgently needed for presentation/marketing.

---

### Option B — Audit First, Then Decide

**Recommended default.**

**Sequence:**

1. AUTH GREEN.
2. Save v0.5 protocol and architecture docs.
3. Claude Code performs Brain Spine Readiness Audit.
4. Bruno + ChatGPT/Gemini/Perplexity review audit.
5. Decide:
   - implement minimal spine now;
   - delay spine until after Cinematic Shell;
   - document only.
6. Proceed with chosen route.

**Pros:**

- No blind architecture decision.
- Uses Claude Code’s real repo visibility.
- Avoids wasting time on theoretical debate.
- Does not mutate runtime.
- Allows better sequencing.

**Cons:**

- Uses some Claude Code time.
- Delays Cinematic Shell slightly.
- Requires discipline not to turn audit into implementation.

**Choose this if:**

- The team wants evidence before deciding.
- Brain Spine may be important, but the risk is unknown.
- Claude Code has enough availability for inspection.

---

### Option C — Minimal Brain Spine Now

**Sequence:**

1. AUTH GREEN.
2. Save v0.5 protocol and architecture docs.
3. Brain Spine Readiness Audit.
4. If audit says low-risk, implement Minimal Brain Spine as no-behaviour-change architecture layer.
5. Migrate existing AI calls into the wrapper.
6. Build/test.
7. Then Cinematic Shell.

**Pros:**

- Prevents later rework.
- Creates a better foundation before major AI features.
- Supports future Trust Dot, Manager Agent, memory, verification.
- Aligns with the long-term one-AI-layer ecosystem.

**Cons:**

- Could destabilise working AI flows.
- May be invisible progress for Bruno/marketing.
- Could become too broad if not strictly scoped.
- May require more Claude credits and more careful testing.

**Choose this only if Claude confirms:**

- no user-facing behaviour changes;
- no new providers;
- no new API keys;
- no schema changes, or schema changes are strictly optional/deferred;
- limited file changes;
- easy rollback;
- current AI behaviour remains identical;
- build and core flows can be tested quickly.

---

### Option D — Full Fleet Now

**Rejected for now.**

**Sequence would be:**

1. Add provider registry.
2. Add Claude/OpenAI/Gemini/Perplexity lanes.
3. Add MCP gateway.
4. Add memory/provenance.
5. Add verification and dashboards.

**Why rejected:**

- Too large for V1.
- Too many providers and keys.
- Too much latency/cost complexity.
- Too much testing burden.
- Risks delaying the launch.
- Contradicts the rule: sequence over parallel.

**Revisit after:**

- Core Flow QA passes.
- Launch readiness improves.
- There is real usage/traction.
- Lean Brain proves useful.

---

## 5. Recommended Decision Right Now

Recommended path:

```text
AUTH GREEN accepted
→ Save documents
→ Brain Spine Readiness Audit only
→ Review audit here
→ Decide Minimal Spine now vs Cinematic Shell first
```

Do **not** implement before the audit.

Do **not** build the full Fleet Brain now.

Do **not** start Cinematic Shell before the Brain Spine timing decision is made, unless Claude’s audit is unavailable and visual progress becomes the highest-priority fallback.

---

## 6. Decision Criteria After Claude’s Audit

### Implement Minimal Brain Spine now if:

- it can wrap existing AI calls without behaviour change;
- no new provider keys are needed;
- no production data migration is needed;
- schema changes are not needed, or can be deferred;
- rollback is simple;
- affected files are limited;
- it makes future AI features safer to build;
- Claude can describe tests clearly.

### Delay Minimal Brain Spine if:

- current AI flow is fragile;
- it touches too many unrelated files;
- it requires schema changes;
- it changes runtime behaviour;
- it delays urgent public visual restoration;
- it risks auth/session stability;
- Claude cannot give a precise rollback plan.

### Document-only if:

- the architecture is valid but too large now;
- current V1 features can safely continue without it;
- no major AI feature is being built immediately;
- Claude availability is limited.

---

## 7. Brain Spine Readiness Audit — Claude Code Prompt

Use after docs are uploaded/saved and AUTH GREEN is accepted.

```text
AI: Claude Code
Mode: Repo Execution Mode / Audit Only
Objective: Produce a Brain Spine Readiness Audit and Minimal Spine timing recommendation. Do not implement anything.

Output expected:
- current AI flow map;
- current AI-related files;
- current memory/session/profile state map;
- proposed Minimal Brain Spine structure;
- risk analysis;
- timing recommendation;
- whether to implement before or after Cinematic Shell.

Do not touch / do not decide:
implementation, new provider APIs, Supabase schema changes, MCP gateway, billing, auth, routing, Stripe, production behaviour, full Fleet implementation, Cinematic Shell.

Context:
AUTH GREEN is accepted. We are deciding whether to build a Minimal Brain Spine before continuing with Cinematic Shell and landing restoration.

We have a future Fleet Brain Architecture target:
- one user-facing Concierge identity;
- one platform-owned control plane;
- platform-owned memory/provenance;
- role-based specialist lanes;
- tool/MCP fabric;
- observability and governance.

We are NOT implementing the full fleet now.

Please inspect the repo and return:

1. Current AI-related files, routes, endpoints, backend functions, and frontend components.
2. Current chat / Concierge / Brave PA / booking intent / alert / review flows.
3. Where AI calls currently happen.
4. Where memory, user profile, session, and conversation state are currently stored.
5. Whether there is already any partial abstraction layer.
6. Proposed minimal Brain Spine module structure, using the existing stack.
7. Interfaces needed:
   - request envelope;
   - model adapter;
   - memory/provenance;
   - verification/audit;
   - tool/action gate;
   - final response synthesis.
8. What can be introduced with no user-facing behaviour change.
9. What would require schema changes.
10. What would require new provider keys.
11. What should be delayed to V2/V3.
12. Risks of doing this before Cinematic Shell.
13. Risks of delaying it until after Cinematic Shell.
14. Recommended next action:
   - implement minimal spine now;
   - document only;
   - delay until after Cinematic Shell.

Important:
- Do not edit code.
- Do not create implementation files.
- Do not install dependencies.
- Do not change runtime behaviour.
- Do not touch auth.
- Do not start Cinematic Shell.
- Return only the audit and proposal.
```

---

## 8. If Claude Recommends Implementing Minimal Spine Now

Use only after review and explicit Bruno approval.

```text
AI: Claude Code
Mode: Repo Execution Mode
Objective: Implement the approved Minimal Brain Spine as a no-behaviour-change architecture layer.

Output expected:
- files changed;
- build result;
- tests run;
- confirmation existing AI behaviour is unchanged;
- rollback plan;
- next recommended migration step.

Do not touch / do not decide:
new providers, MCP, Supabase schema unless explicitly approved, auth, routing, Stripe, Cinematic Shell, full Fleet implementation.

Rules:
- Existing AI behaviour must remain unchanged.
- Existing model provider remains the only active runtime provider.
- No Gemini, Perplexity, OpenAI, OCR, transcription, MCP, or extra provider API integration yet.
- No user-visible changes.
- No database migration unless separately approved.
- Add only interfaces/wrappers needed so future AI features can use the same spine.

Return:
1. files changed;
2. build result;
3. tests run;
4. what behaviour remains unchanged;
5. rollback plan;
6. next recommended migration step.
```

---

## 9. If Claude Recommends Delaying Spine

Proceed to Cinematic Shell.

Next ChatGPT task:

```text
AI: ChatGPT
Mode: Execution / Computer Mode
Objective: Prepare the Integration Packet for Cinematic Shell — Three.js Cosmic Intro.
Output expected:
- complete Integration Packet;
- Next.js 14 safety;
- WebGL fallback;
- reduced motion;
- asset media requirements;
- manual tests;
- Claude Code prompt.

Do not touch / do not decide:
repo integration, auth, routing, middleware, backend, Stripe, Supabase, env vars, global CSS.
```

---

## 10. Files To Save In Repo

Suggested documentation paths:

```text
docs/protocol/PRODUCT_OPERATING_SYSTEM_v0.5.md
docs/architecture/FLEET_BRAIN_ARCHITECTURE_BRIEF.md
docs/architecture/BRAIN_SPINE_FLEET_INTEGRATION_DECISION_PACK_v0.2.md
```

If these folders do not exist, Claude Code may create them for documentation only.

---

## 11. Final Recommendation

The best strategy is:

> Full Fleet Brain later.  
> Minimal Brain Spine only after repo audit.  
> Cinematic Shell next if the spine is risky.  
> No implementation without Claude Code’s evidence-based recommendation.

This keeps the project ambitious without becoming chaotic.
