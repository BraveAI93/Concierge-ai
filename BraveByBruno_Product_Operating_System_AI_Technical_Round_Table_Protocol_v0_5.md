# BRAVE BY BRUNO / THE CONCIERGE
## Product Operating System + AI Technical Round Table Protocol - v0.5
### Build, Launch, and Post-Launch Governance for a Premium AI Product

**Status:** Adapted and repo-context-integrated by Claude (chat) at Bruno's request; ready for Claude Code to save/confirm inside the repo.  
**Date:** 7 July 2026  
**Supersedes:** v0.4 draft (ChatGPT + Gemini + Bruno)  
**Adapted by:** Claude (chat/mobile), using full session history, the Master Document, the Concierge Codex Audited v1, and tonight's live terminal work, per Bruno's explicit request to make v0.4 precise and complete.

**Current project state (corrected):** Auth/login is **not yet confirmed resolved** - it is mid-verification. Tonight's session fixed two real bugs (a cookie-name mismatch, and an anonymous `_directToken` cookie-injection bypass) and reverted one accidental file corruption (`concierge-backend/db/supabase.go`, caused by the known touch-scroll input bug - see 7.0). The backend builds cleanly (`BUILD_OK` confirmed). Still outstanding: a real email+password login producing a genuine `Set-Cookie: cai_token=...` header with a real, non-placeholder value. **AUTH GREEN cannot be declared until that specific evidence exists** - see 7.1 for the exact test.

**On "Claude" vs "Claude Code":** v0.4 used both names without distinguishing them, which risked real confusion. **Claude (chat)** is this mobile/desktop conversation - strategy continuity, reading terminal screenshots, live relay (see new 5.5). **Claude Code** is the terminal agent running inside the Codespace (`stunning-carnival`, repo `concierge-ai`) - the only one that actually edits files and pushes commits. **Correction (confirmed against Anthropic's own help center):** Claude (chat) and Claude Code are NOT independent on usage limits - they draw from the same shared pool (a 5-hour rolling window plus a separate weekly cap). Heavy use in one reduces headroom in the other. This is exactly why routing more coordination through Claude (chat) does not add capacity - see new section 5.6 for how to actually reduce load on this shared pool.

**Immediate use:** pass this document to Claude Code when available so it can save/adapt it inside the repo, then resume the AUTH GREEN verification exactly as specified in 7.1 - not restart it, and not declare it passed from a summary alone.

> **Changelog, v0.4 -> v0.5 (Claude chat adaptation pass):** corrected current-state claims (auth not yet green; rate-limit reset time was wrong); distinguished Claude (chat) from Claude Code as separate roles (2, 3, 4, 5.5, 6.2, 17); added the proven Live Relay & Verification role and its Safe-to-Self-Approve rule (5.5); added the concrete AUTH GREEN test procedure (7.1); added the touch-scroll input/file-corruption hazard as a named repo-specific risk (7.0); mapped the Concierge Codex Audited v1's six build-ready locked features into the Build Pipeline so none are silently dropped (7.6); mapped the Master Document's Blocco A3-A5 feature roadmap into the pipeline so it isn't orphaned by this protocol (7.7); added real repo paths and the exact save/verify prompts (16); added three Anti-Chaos Rules from tonight's proven practice (14.15-14.17). No locked decision - stack, URL architecture, pricing, roadmap order - was reopened or changed.

---

## 0. Why v0.4 Exists

v0.3 defined how Bruno, ChatGPT, Gemini, Perplexity Plus, and Claude Code work together without creating chaos.

v0.4 expands that protocol into a full **Product Operating System**: not only how to build features, but how to move from the current technical state to a polished launch, and then how to manage the product after launch.

The purpose is not to replace employees with magic. The purpose is to give Bruno a precise operating structure that lets a solo founder coordinate AI tools like a small, disciplined product team.

**Core promise:**

> The system must allow Bruno to operate the project without deep technical knowledge, while still producing an experience that feels premium, safe, intentional, trustworthy, and worthy of a luxury technology brand.

**Core principle:**

> Simple in front. Complex behind. Safe always.

---

## 1. The Product Operating System

The protocol now has three levels:

### 1.1 Build Pipeline

What must happen from now until the product is technically and visually ready:

- fix and confirm critical infrastructure;
- restore the premium public experience;
- QA the core flow;
- harden legal, security, privacy, and trust;
- prepare the product for launch.

### 1.2 Launch Pipeline

What must happen before and during the marketing launch:

- define positioning;
- prepare launch assets;
- prepare demo script;
- verify public pages;
- test conversion paths;
- prepare fallback plans;
- go live deliberately, not chaotically.

### 1.3 Management Pipeline

What must happen after launch:

- triage bugs;
- answer support;
- collect feedback;
- prioritise improvements;
- publish content;
- monitor metrics;
- protect legal/security integrity;
- keep the product evolving without fragmenting.

---

## 2. Source of Truth Hierarchy

The round table is not a debate where the loudest model wins. The hierarchy is fixed:

1. **Bruno** - Founder, Vision Holder, Product Owner, final decision maker.
2. **Master Document** - strategic constitution: vision, URL architecture, product decisions, pricing, roadmap, build order.
3. **Session Handover** - live technical state and next immediate action.
4. **Concierge Codex Audited** - approved feature library, technical/legal/security status, external blockers.
5. **GitHub repository** - actual source code state.
6. **Vercel/Render/Supabase production state** - real deployed truth (production domain only - never a Vercel preview/alias URL, which can silently diverge - see 14.17).
7. **Claude (chat)** - live relay, verification discipline, and terminal supervision between Bruno and Claude Code (see 5.5). Does not edit the repo.
8. **Claude Code** - only repo integrator.
9. **ChatGPT** - strategy, Integration Packets, QA plans, prompts, documentation.
10. **Gemini** - visual/UX/reference reviewer.
11. **Perplexity Plus** - external research, technical docs, competitor research, compliance evidence.

No AI may reopen locked decisions about stack, routing, URL architecture, pricing, Universe Engine, roadmap, or legal guardrails unless Bruno explicitly asks for a strategic reconsideration and the Master Document is formally updated.

---

## 3. Diamond Protocol

**Bruno decides.  
Master Doc governs.  
Perplexity researches.  
ChatGPT prepares.  
Gemini observes.  
Claude (chat) relays and verifies in real time.  
Claude Code integrates.  
Production verifies.**

This is the required sequence.

The rule is not "multiple AIs build the product." The rule is:

> Multiple AIs prepare and review. One pipeline integrates. Production verifies.

---

## 4. Mode of Work - Required on Every Task

Every task must state the AI and mode before the content.

This avoids confusion between a strategic conversation and an operational handoff.

### 4.1 Mandatory Task Header Template

Every task should start like this:

```text
AI:
Mode:
Objective:
Output expected:
Do not touch / do not decide:
```

Example:

```text
AI: ChatGPT
Mode: Execution / Computer Mode
Objective: Prepare an Integration Packet for the Cinematic Shell.
Output expected: scope, component path, dependencies, fallback, tests, Claude prompt.
Do not touch / do not decide: auth, routing, middleware, backend, database, Stripe, env vars, global CSS.
```

### 4.2 Task Mode Must Be Pre-Selected Where Possible

Bruno should not have to choose the mode from zero every time. Where possible, the mode should be specified by the task itself.

Default task labels:

- **Strategic decision** -> ChatGPT, Strategy Mode.
- **Updated documentation / protocol** -> ChatGPT, Execution / Computer Mode; Claude Code only saves/adapts.
- **External facts / recent technical docs** -> Perplexity, Research Mode.
- **Legal/security evidence** -> Perplexity, Compliance / Evidence Mode.
- **Screenshot / luxury UI critique** -> Gemini, Visual Review Mode.
- **Micro CSS/UI notes** -> Gemini, Visual Execution Notes Mode.
- **Feature packet** -> ChatGPT, Execution / Computer Mode.
- **Repo integration** -> Claude Code, Repo Execution Mode only.
- **Bug triage** -> ChatGPT Strategy first, then Claude Code Repo Execution if a code fix is needed.
- **Launch planning** -> ChatGPT Strategy, then ChatGPT Execution for checklists/assets.
- **Post-launch incident** -> ChatGPT Strategy for triage, Claude Code for technical fix, Perplexity only if current policy/security docs are needed.

If Bruno has uncertainty, he asks strategic questions first. No operational packet is produced until the direction is clear.

### 4.3 Two Major Modes

#### Strategy / Normal Mode

Used for thinking, deciding, comparing, questioning, and prioritising.

Typical outputs:

- opinion;
- risk analysis;
- pros/cons;
- roadmap decision;
- product reasoning;
- strategic recommendation.

A Strategy Mode answer is **not ready to be given directly to Claude Code**.

#### Execution / Computer Mode

Used for turning a decision into a controlled operational deliverable.

Typical outputs:

- Integration Packet;
- Research Packet;
- Visual Review Packet;
- Claude Code prompt;
- QA checklist;
- legal/security checklist;
- asset list;
- launch checklist;
- post-launch triage plan;
- step-by-step execution plan.

An Execution Mode output can be prepared for handoff, but still requires Bruno approval before reaching Claude Code.

### 4.4 Simple Decision Rule for Bruno

Ask:

> Am I thinking, researching, looking, preparing, or building?

- **Thinking** = ChatGPT Strategy Mode.
- **Researching** = Perplexity Research or Compliance Evidence Mode.
- **Looking** = Gemini Visual Review Mode.
- **Preparing** = ChatGPT Execution / Computer Mode.
- **Building** = Claude Code Repo Execution Mode.
- **Watching the terminal / deciding a menu prompt** = Claude (chat) Live Relay Mode.

---

## 5. AI Usage Matrix

### 5.1 ChatGPT

#### Strategy Mode

Use ChatGPT Strategy Mode when Bruno needs to:

- decide priority;
- compare alternatives;
- identify risks;
- turn chaos into roadmap;
- test whether a feature belongs to V1, V2, or later;
- evaluate product, technical, legal, security, marketing, or monetisation implications.

Prompt pattern:

```text
AI: ChatGPT
Mode: Strategy Mode
Objective: Decide whether to do X or Y first.
Output expected: recommendation, dependencies, risks, next step.
Do not prepare an Integration Packet yet.
```

#### Execution / Computer Mode

Use ChatGPT Execution Mode when Bruno needs:

- Integration Packets;
- Claude prompts;
- QA test plans;
- structured handovers;
- copy blocks;
- feature specifications;
- legal/security checklists;
- launch checklists;
- post-launch management procedures;
- documentation updates.

Prompt pattern:

```text
AI: ChatGPT
Mode: Execution / Computer Mode
Objective: Prepare a complete Integration Packet for [feature].
Output expected: scope, allowed files, forbidden files, dependencies, fallback, tests, Claude prompt.
Do not assume permission to integrate into the repo.
```

### 5.2 Perplexity Plus

Perplexity is the external evidence layer. It is not the final decision-maker and it does not produce repo-ready code.

#### Research Mode

Use Perplexity Research Mode for:

- latest documentation;
- technical best practices;
- library compatibility;
- current frontend patterns;
- competitor analysis;
- design/technology benchmarks;
- pricing/API updates.

Prompt pattern:

```text
AI: Perplexity Plus
Mode: Research Mode
Objective: Find current best practices for [topic].
Output expected: reliable sources, date, risks, recommendation, links.
Do not write final implementation code.
```

#### Compliance / Evidence Mode

Use Perplexity Compliance Evidence Mode for:

- GDPR / UK GDPR;
- ICO guidance;
- EDPB guidance;
- OWASP;
- NCSC;
- payments/security;
- biometric, voice, video, health, legal, or sensitive-data features.

Prompt pattern:

```text
AI: Perplexity Plus
Mode: Compliance / Evidence Mode
Objective: Research current official guidance on [topic].
Output expected: official sources, compliance checklist, unresolved legal questions, professional-review triggers.
Do not give final legal advice. Flag where a lawyer or DPIA is required.
```

Perplexity output should be passed to ChatGPT as raw evidence before ChatGPT prepares the operational packet.

### 5.3 Gemini

Gemini is the visual and UX eye. It can also review reference boards and competitor screenshots.

#### Visual Review Mode

Use Gemini Visual Review Mode for:

- screenshots;
- reference boards;
- luxury competitor pages;
- UI hierarchy;
- glassmorphism quality;
- typography;
- spacing;
- palette;
- perceived luxury/trust/futurism.

Prompt pattern:

```text
AI: Gemini
Mode: Visual Review Mode
Objective: Review this screenshot/reference for premium UI and visual coherence.
Output expected: visual critique, UX issues, luxury perception, improvement notes.
Do not propose changes to auth, routing, middleware, backend, database, Stripe, env vars, or general architecture.
```

#### Visual Execution Notes Mode

Use Gemini Visual Execution Notes Mode for:

- localised CSS suggestions;
- before/after visual notes;
- asset/reference comments to pass back to ChatGPT;
- precise visual adjustments that can be converted into an Integration Packet.

Prompt pattern:

```text
AI: Gemini
Mode: Visual Execution Notes Mode
Objective: Convert the visual review into localised implementation notes.
Output expected: spacing, blur, typography, colour, motion, asset notes.
Do not write the final Integration Packet. Do not touch architecture.
```

### 5.4 Claude Code

Claude Code operates only in Repo Execution Mode during build sessions.

Use Claude Code for:

- file edits;
- repo integration;
- dependency installation;
- build/test;
- commit;
- push;
- deploy verification;
- route checks;
- production verification.

Prompt pattern:

```text
AI: Claude Code
Mode: Repo Execution Mode
Objective: Integrate only the approved packet below.
Output expected: files changed, build result, tests, commit, production PASS/FAIL.
Do not modify auth, routing, middleware, backend, database, Stripe, env vars, global CSS, or unrelated files unless explicitly authorised.
```

Claude should not brainstorm freely while integrating. If Claude identifies a strategic issue, it must pause and ask Bruno.

### 5.5 Claude (Chat/Mobile) - Live Relay & Verification Mode

This role was not present in v0.3/v0.4 but has been doing real, proven work every session: reading Claude Code's terminal output directly from screenshots, translating its numbered approval menus, deciding what to approve or pause on, and catching gaps between what Claude Code claims and what it has actually shown proof of.

Use Claude (chat) Live Relay Mode when Bruno is:

- inside an active Claude Code session and needs to know what to answer to a menu prompt;
- on a touchscreen device where terminal scrolling can inject garbage input (see 7.0);
- unsure whether a Claude Code summary is backed by real evidence (a diff, a header, a PID) or is just a confident sentence;
- coordinating a long back-and-forth where losing track of what's still unverified is the real risk.

**Core rule Claude (chat) enforces on every terminal step: a summary is a claim, not proof.** Concrete evidence only - a real `git diff`, an actual `Set-Cookie` header, a printed PID, a `BUILD_OK`. Never accept "fixed," "confirmed," or "done" as a substitute for showing the artifact.

**Safe-to-self-approve rule (proven across this entire session):** Bruno may approve numbered menu prompts himself, without waiting for Claude (chat), when the command is local-only and read-only, or writes only to `/tmp` or `.next` - for example: `curl` against `localhost`, `cat`, `grep`, `npm start`/`npm run dev`, `pkill`/`kill` on local dev processes, reading files. Claude (chat) should always be consulted before: `git push`, any `DROP`/`DELETE` on a database table, or any write outside `/tmp` and `.next` (including `git checkout --` on a tracked file, which is safe but still a working-tree mutation worth a second look).

Prompt pattern (mostly implicit - Bruno simply pastes a screenshot):

```text
AI: Claude (chat)
Mode: Live Relay Mode
Objective: Tell me what to answer to this terminal prompt, and why.
Output expected: which option is safe, what the command actually does, what to watch for in the result, and what remains unverified afterward.
Do not accept a Claude Code summary as proof without a concrete artifact.
```

Claude (chat) does not edit files or run commands itself in this mode - it only reads what Bruno shows it and tells him what to do next.

### 5.6 When To Consult More Than One Agent (And When Not To)

Consulting multiple assistants on the same question is not free - even when it doesn't touch Claude's pool, it still costs Bruno's time and attention to relay between them. Use this test before doing it:

> Is this decision expensive to reverse (money spent, code deeply integrated, a public-facing commitment), **and** does it genuinely need more than one kind of expertise (visual eye, legal/technical research, strategic tradeoff) rather than one clear skill?

**Both yes -> consult multiple agents**, in the cheapest pools first:
- Architecture or sequencing calls that trade off against each other (e.g., whether a feature belongs before or after another) - ChatGPT (strategy) + Perplexity (comparable research).
- Visual decisions with a technical constraint attached (e.g., a heavy animation's mobile performance) - Gemini (visual) + a technical feasibility check.
- Anything legal/compliance-adjacent (data collection, consent, retention) - Perplexity (research) informing ChatGPT's strategic write-up, before it ever reaches a build decision.
- Any request to reopen a locked decision - by 2's own rule this needs Bruno's explicit trigger anyway, and multiple perspectives are justified precisely because it's rare and high-stakes.

**Either no -> single agent, fastest available, do not multi-consult:**
- Anything already locked (re-litigating it wastes tokens across every pool involved, not just one).
- Routine repo work matching an existing pattern - Claude Code alone.
- The live terminal relay/verification loop itself - adding more assistants into that loop recreates the exact "which menu, which number" confusion this protocol exists to prevent.
- Mechanical or administrative steps (a push, a file save, a rename).
- Small stylistic tweaks with no architectural weight - default to whichever assistant is already looking at it.

**Cost discipline:** when a question does warrant multiple agents, let ChatGPT/Gemini/Perplexity do the back-and-forth first and bring Claude (chat) only the synthesized outcome - not the full debate - unless the decision specifically needs repo-context knowledge to weigh in. This keeps the shared Claude/Claude Code pool reserved for what only it can do.

---

## 6. Operating Lanes

Every task belongs to one lane. This helps Bruno understand what kind of work is happening.

### 6.1 Product Lane

Owns:

- user value;
- feature priority;
- pricing logic;
- onboarding promise;
- market positioning;
- launch readiness.

Default AI:

- ChatGPT Strategy Mode.
- Perplexity Research Mode when external evidence is needed.

### 6.2 Engineering Lane

Owns:

- repo integration;
- build;
- deploy;
- auth;
- backend;
- data;
- performance;
- technical QA.

Default AI:

- Claude Code Repo Execution Mode.
- Claude (chat) Live Relay Mode for real-time terminal supervision while Claude Code runs.
- ChatGPT Execution Mode for packets/prompts.

### 6.3 Visual Experience Lane

Owns:

- premium look;
- UX hierarchy;
- motion quality;
- spatial/cosmic identity;
- asset quality;
- responsive polish.

Default AI:

- Gemini Visual Review Mode.
- ChatGPT Execution Mode for final packet.

### 6.4 Research and Compliance Lane

Owns:

- current documentation;
- legal evidence;
- GDPR/UK GDPR;
- OWASP/NCSC references;
- competitor benchmark;
- unresolved external-risk flags.

Default AI:

- Perplexity Research or Compliance / Evidence Mode.
- ChatGPT Strategy Mode for interpretation.

### 6.5 Marketing and Launch Lane

Owns:

- messaging;
- demo script;
- waitlist/launch copy;
- social content;
- email sequences;
- press/investor narrative;
- visual launch assets.

Default AI:

- ChatGPT Strategy/Execution Mode.
- Gemini Visual Review Mode.
- Perplexity Research Mode for benchmarks/trends.

### 6.6 Support and Growth Lane

Owns:

- bug triage;
- user support;
- feedback loops;
- feature prioritisation;
- growth experiments;
- churn/retention learning.

Default AI:

- ChatGPT Strategy Mode for triage.
- Claude Code Repo Execution Mode for fixes.
- Perplexity only if fresh external evidence is needed.

---

## 7. Build Pipeline - From Now to Product Ready

### 7.0 Known Repo-Specific Hazard - Touch-Scroll Input Corruption

This is not hypothetical - it has already happened twice this project. Claude Code runs as a full-screen terminal UI with no native scrollback; on a touchscreen, scrolling directly over the terminal or an open editor tab can be misread as keystrokes. It has produced garbage in the command input line (`NaN;NaNMaN;...`) and, once, silently corrupted an open, uncommitted file (`concierge-backend/db/supabase.go` - a truncated string literal plus a dangling function, confirmed by the file no longer compiling).

Standing rule: before typing into an active Claude Code prompt, or before trusting any unstaged/uncommitted diff as intentional work, check for this signature - a string or line cut off mid-word, followed by content that doesn't belong. If found, treat it as corruption, not a real edit: confirm with `git diff` / a direct file read, then `git checkout -- <path>` to discard it and rebuild to confirm (`BUILD_OK`). Scroll using the panel's scrollbar edge, not a swipe over the content area, and check the input line is empty before typing.

### Phase B0 - Protocol and Auth Green

**Goal:** establish operating control before new feature work.

Tasks:

1. Save/adapt this v0.5 protocol inside the repo.
2. **Confirm Auth/login with concrete evidence, not a summary** - the exact test:
   - Restart the local dev server if not already running.
   - `curl -s -D - -o <body>.json -c <cookiejar>.txt -X POST http://localhost:3000/api/auth/login -H "Content-Type: application/json" -d '{"email":"<real test account email>","password":"<real password>"}'`
   - Confirm the printed response headers include `Set-Cookie: cai_token=<a real, non-placeholder value>`.
   - Confirm the same test account, without that cookie, is redirected to `/theconcierge/owner-auth` (not the marketing landing page) when hitting `/theconcierge/dashboard`.
   - Confirm the anonymous `_directToken` cookie-injection path (`/api/auth/login`) is closed or restricted to non-production - show the source line, not a description of it.
3. Confirm session cookie behaviour on the real production domain (`bravebybruno.com`), not the Vercel preview alias, which can silently diverge from production - see 14.17.
4. Confirm dashboard gating (logged out -> login page; logged in -> dashboard).
5. Confirm refresh persistence.
6. Confirm logout clears the cookie.
7. Confirm production domain behaviour end to end.

Gate to next phase:

```text
AUTH GREEN - PASS
Protocol saved:
Commit:
Set-Cookie header captured (real value, not placeholder): YES/NO
_directToken bypass closed, source confirmed: YES/NO
Deploy:
Domain tested (bravebybruno.com, not alias):
Known issues:
```

No visual feature begins until this is passed, and it is not passed on a Claude Code summary alone.

### Phase B1 - Public Experience Restoration

**Goal:** make the public site feel worthy of the vision again.

Primary tasks:

1. Cinematic Shell - Three.js Cosmic Intro.
2. Restore full landing content.
3. Restore Difference section.
4. Restore pricing cards.
5. Restore clear CTAs.
6. Confirm Light/Easy fallback plan.

Required AI flow:

1. Perplexity Research Mode for current Three.js + Next.js 14 best practices if needed.
2. ChatGPT Execution Mode for the Cinematic Shell Integration Packet.
3. Gemini Visual Review Mode for screenshot/reference review.
4. ChatGPT refines packet.
5. Claude Code integrates one packet only.
6. Production verifies.

Gate to next phase:

- public landing loads;
- no auth regression;
- no console-breaking errors;
- mobile works;
- fallback works;
- the page communicates premium identity in the first 5 seconds.

### Phase B2 - Core Flow QA

**Goal:** prove the product works as a product, not just as a beautiful landing page.

Test:

- signup;
- login;
- onboarding;
- dashboard;
- chat;
- business page creation;
- QR/contact buttons;
- welcome email;
- legal consent;
- profile switching if relevant;
- logout and re-login.

Gate to next phase:

```text
CORE FLOW - PASS
Critical blockers:
High issues:
Medium issues:
Low issues:
Launch-safe: YES/NO
```

### Phase B3 - Dashboard Visual Integration

**Goal:** make the dashboard elegant and usable without turning it into visual noise.

Rule:

- operational data = readable, flat, high contrast;
- brand/status elements = cosmic/premium aesthetic.

Tasks:

- glassmorphism consistency;
- card hierarchy;
- clickable insights;
- readable lead stats;
- mobile dashboard sanity;
- no animation overload.

Gate to next phase:

- dashboard supports decisions quickly;
- visuals feel premium but not distracting;
- no core flow regression.

### Phase B4 - Trust, Legal, and Security Hardening

**Goal:** make the product safe enough to show seriously.

Minimum checks:

- auth/session review;
- cookie flags;
- API key exposure check;
- input validation;
- output escaping;
- rate limiting;
- Supabase permissions/RLS/grants review;
- Stripe webhook signature verification if payment flow is active;
- privacy policy;
- terms;
- AI transparency;
- GDPR consent logs;
- retention/deletion policy;
- DPIA trigger list.

Important:

Voice/video/biometric/session-learning features must remain blocked until professional legal review and DPIA are completed.

Gate to next phase:

```text
TRUST GATE - PASS/FAIL
Security blockers:
Legal blockers:
Privacy blockers:
AI transparency complete:
Launch-safe: YES/NO
```

### 7.6 Codex-Locked Features Queue

The Concierge Codex Audited v1 already cleared six features for build (locked), independent of this protocol. This pipeline must not let them get lost behind the visual/auth work above. Once B0-B2 pass, these become the standing backlog for Product Lane + Engineering Lane, prioritised by Bruno, in this order unless he says otherwise (dependency order per the Codex itself):

1. **Trust Dot** - no external blocker, ready now.
2. **Brave Star - 3 behavioural states** (replaces the single decibel-line concept) - no external blocker, ready now.
3. **Sfondo Location-Aware** (skyline + weather + time of day) - no external blocker, ready now.
4. **Offline-First Sync + Priority Reminder Fallback** - build after Manager Agent core, shares its connector infrastructure.
5. **Manager Agent** (incremental, one connector at a time) - ready now; Invyted stays excluded until it publishes a public API.
6. **Dashboard Visual Integration** - already scheduled as Phase B3 above; no change.

**Voice & Video Session Learning stays in spec only** - DPIA and GDPR/UK GDPR legal review required before any build, per the Codex's binding block. Do not schedule it into a sprint.

### 7.7 Feature Roadmap Integration - Master Doc Blocco A3-A5

This protocol governs *how* work gets coordinated. It does not replace the Master Document's Blocco sequence, which remains the source of truth for *what* gets built and in what order (hierarchy position 2). To avoid the roadmap being silently orphaned by this new protocol:

- **Blocco A3 (in progress before tonight's auth detour):** Google Calendar OAuth (read-only) is mid-implementation - CLIENT_ID/CLIENT_SECRET retrieval from Google Cloud Console was the last confirmed step. This resumes as Engineering Lane work once AUTH GREEN passes, before or alongside Phase B1, since it is closer to completion than any new visual work.
- **Blocco A4 (B&B mode, installable PWA):** queued after Phase B2 (Core Flow QA), since B&B is an eighth professional mode that must not be added before the existing seven are proven stable end to end.
- **Blocco A5 (Stripe booking + deposit, Render $7/month upgrade, WhatsApp Business API via Twilio, email campaigns):** explicitly gated behind the Launch Pipeline (section 8) - several A5 items involve real client-facing money or messaging and should not go live before Phase L3 Pre-Launch QA passes. The Render upgrade specifically must happen before any real client is sent a link, per existing Master Doc pricing rules.

No Blocco A3-A5 item is cancelled, deprioritised, or reopened for debate by this protocol. This section only sequences them against the new B0-B4/L1-L5 structure.

---

## 8. Launch Pipeline - From Product Ready to Public Launch

### Phase L1 - Launch Positioning

**Goal:** define what the market must understand immediately.

Outputs:

- one-sentence positioning;
- hero headline;
- subheadline;
- target users;
- key pain points;
- key proof points;
- founder story angle;
- offer/CTA.

Default AI:

- ChatGPT Strategy Mode.
- Perplexity Research Mode for competitor positioning.
- Gemini Visual Review Mode for brand perception.

Gate:

- Bruno can explain the product in 15 seconds.

### Phase L2 - Launch Asset Preparation

**Goal:** prepare everything needed to present the product publicly.

Assets:

- landing screenshots;
- short demo video plan;
- founder demo script;
- social launch captions;
- LinkedIn announcement;
- Instagram/TikTok/Reel concept;
- email to warm contacts;
- investor/founder intro message;
- press/partnership pitch;
- FAQ answers;
- support email template.

Gate:

- all launch assets are in a folder/document;
- no missing CTA;
- no unsupported legal/security claim.

### Phase L3 - Pre-Launch QA Day

**Goal:** simulate a real visitor before launching.

Checklist:

- homepage loads on mobile and desktop;
- sign-up works;
- login works;
- demo path works;
- dashboard loads;
- contact/business page works;
- legal pages visible;
- support route/email visible;
- analytics active if used;
- no broken links;
- no obvious console errors;
- fallback plan ready if deployment fails.

Gate:

```text
PRE-LAUNCH QA - PASS/FAIL
Launch allowed: YES/NO
Blockers:
Fallback plan:
```

### Phase L4 - Controlled Launch

**Goal:** launch deliberately, not explosively.

Recommended order:

1. Private review by 2-5 trusted people.
2. Fix critical blockers only.
3. Soft launch to warm network.
4. Collect feedback.
5. Public launch content.
6. Founder/partner outreach.
7. Post-launch monitoring.

Rule:

No major feature work during launch window unless it fixes a critical blocker.

### Phase L5 - Launch Retrospective

Within 48-72 hours after launch, Bruno should collect:

- what worked;
- what confused users;
- what broke;
- what converted;
- what people asked for;
- what should not be changed yet;
- next highest-leverage improvement.

Output:

```text
LAUNCH RETRO
Traffic / interest:
Signups:
Feedback themes:
Bugs:
Conversion blockers:
Next 3 priorities:
Do not touch yet:
```

---

## 9. Post-Launch Management Pipeline

### 9.1 Weekly Operating Rhythm

A simple founder rhythm:

**Daily, 10 minutes:**

- check critical errors;
- check user messages;
- check signup/payment/contact events;
- note urgent blockers.

**Twice weekly, 30 minutes:**

- review feedback;
- classify bugs;
- choose one improvement;
- prepare one Claude task if needed.

**Weekly, 60 minutes:**

- review metrics;
- review marketing content;
- decide next sprint;
- update the handover/status doc.

### 9.2 Bug Severity System

Use this classification:

**S0 - Emergency**  
Security breach, data leak, payment failure affecting users, site inaccessible. Stop all normal work.

**S1 - Critical**  
Login broken, signup broken, dashboard inaccessible, core chat unusable, launch-blocking production issue.

**S2 - High**  
Important feature broken but workaround exists, major UI break on mobile, business page broken for some users.

**S3 - Medium**  
Visual bug, non-critical copy issue, minor dashboard inconsistency, UX friction.

**S4 - Low**  
Nice-to-have improvement, polish, future idea.

Rules:

- S0/S1 interrupt the roadmap.
- S2 is next sprint unless launch is blocked.
- S3/S4 go to backlog.
- No new feature is started while S0/S1 is unresolved.

### 9.3 Support Triage

Every user issue should be captured as:

```text
User issue:
User impact:
Severity:
Page/flow affected:
Screenshot/link:
Reproduction steps:
Owner:
Next action:
```

AI flow:

- ChatGPT Strategy Mode to classify.
- Claude Code Repo Execution Mode if code fix needed.
- Gemini Visual Review if UI issue.
- Perplexity Compliance Mode if legal/security/privacy issue.

### 9.4 Growth and Content Management

Post-launch growth should not create product chaos.

Monthly content lanes:

- founder story;
- product demo;
- creator/freelancer pain point;
- behind-the-scenes build;
- AI trust/security explanation;
- case study/test user story;
- premium visual teaser.

Default AI:

- ChatGPT Strategy/Execution for content plan and captions.
- Gemini Visual Review for creative direction.
- Perplexity Research for trends/market context.

Rule:

Marketing must never make claims the product, legal docs, or security posture cannot support.

### 9.5 Product Improvement Loop

Every improvement must answer:

1. What user pain does this solve?
2. Is it V1, V2, or later?
3. Does it increase trust, conversion, retention, or revenue?
4. Does it create legal/security risk?
5. Can it be built as a small packet?
6. What is the rollback plan?

If not clear, it remains in backlog.

---

## 10. Integration Packet Standard

Every UI/React/Three.js/Tailwind feature must reach Claude as an Integration Packet, not as loose code.

### 10.1 Required Sections

1. **Packet name**
2. **Mode header** - AI, mode, objective, output expected, forbidden scope.
3. **Purpose** - what the feature does and why it exists.
4. **Master Doc alignment** - which locked decision it supports.
5. **Priority** - Now / Soon / Later / Blocked.
6. **Proposed component path** - example: `components/cinematic/CinematicShell.jsx`.
7. **Proposed import path** - example: `app/theconcierge/page.jsx`.
8. **Allowed file changes** - exact list.
9. **Forbidden file changes** - auth, routing, middleware, backend, database, Stripe, env vars, global CSS unless explicitly approved.
10. **Dependencies** - required packages and alternatives.
11. **Asset media required** - `.glb`, `.webp`, textures, fallback image, icons, fonts, video, audio, size constraints.
12. **Next.js compatibility** - client component safety, no `window/document` before mount, SSR/hydration safety.
13. **Fallback strategy** - WebGL failure, reduced motion, slow device, Light Mode.
14. **Performance notes** - resize handling, animation cleanup, mobile FPS, bundle concerns.
15. **Accessibility notes** - reduced motion, contrast, no aggressive flashing, readable content.
16. **Code** - isolated component/code only.
17. **Manual tests** - PASS/FAIL checklist.
18. **Claude Code prompt** - exact instruction to integrate safely.
19. **Rollback plan** - files to revert if the feature breaks.

### 10.2 Maximum Standby Rule

ChatGPT may prepare at most two Integration Packets in standby. Prefer one.

Current standby target:

1. **Cinematic Shell - Three.js Cosmic Intro**

No second packet should be created until the Cinematic Shell packet is accepted or explicitly paused.

---

## 11. Work Allowed While Claude Code Is Blocked

When Claude Code is blocked, the repo is frozen. However, preparation continues.

### 11.1 Allowed Parallel Work

ChatGPT may prepare:

- Integration Packet drafts;
- Claude prompts;
- QA checklists;
- landing copy;
- pricing copy;
- security/legal checklist;
- launch plan;
- post-launch triage templates;
- technical handover notes;
- project documentation.

Gemini may prepare:

- screenshot review;
- reference board review;
- competitor luxury UI analysis;
- visual execution notes;
- palette/typography/spacing critique.

Perplexity may prepare:

- latest technical documentation research;
- Next.js/Three.js compatibility research;
- competitor research;
- legal/security evidence packets;
- GDPR/ICO/OWASP/NCSC research summaries.

Asset preparation may also happen:

- define texture requirements;
- define `.webp` fallback image requirements;
- define `.glb` model requirements if needed;
- define particle/motion visual rules;
- prepare image-generation prompts for fallback assets;
- compress/format asset specs;
- define naming conventions and target file sizes.

### 11.2 Forbidden While Claude Is Blocked

Do not change or prepare direct repo mutations involving:

- auth;
- login;
- cookies/session;
- middleware;
- routing;
- backend Go;
- Supabase schema;
- Stripe;
- env vars;
- global CSS;
- dependency upgrades;
- security architecture;
- legal/sensitive-data features;
- production deployment.

---

## 12. Decision Trees

### 12.1 If Bruno Has a New Idea

1. Ask ChatGPT in Strategy Mode.
2. Decide Now / Soon / Later / Blocked / Kill.
3. If Now, ask for an Execution Packet.
4. If visual, ask Gemini for visual review.
5. If current evidence is needed, ask Perplexity before the packet.
6. Claude integrates only after approval.

### 12.2 If Bruno Finds a Bug

1. Capture screenshot/error/route.
2. Ask ChatGPT Strategy Mode to classify severity.
3. If S0/S1/S2, create a Claude fix prompt.
4. Claude fixes one issue only.
5. Verify production.
6. Update status.

### 12.3 If Bruno Is Unsure What To Do Next

Ask ChatGPT:

```text
AI: ChatGPT
Mode: Strategy Mode
Objective: Decide the next highest-leverage task from the current state.
Output expected: one recommended next task, why, blockers, who to ask, and exact next prompt.
Do not prepare multiple parallel workstreams.
```

### 12.4 If Claude Is Blocked

1. Freeze repo work.
2. Prepare one standby packet only.
3. Prepare launch/QA/support docs if useful.
4. Do not open more than two active preparation threads.
5. When Claude returns, resume with the next saved task.

### 12.5 If A Task Feels Too Big

Split it into:

- research;
- strategy;
- packet;
- visual review;
- repo integration;
- QA;
- documentation.

Never send a huge multi-feature task to Claude.

---

## 13. Definition of Done

A feature is not done when code is written. It is done only when:

- it is integrated into the repo;
- it builds;
- it deploys;
- it works on the real domain;
- it does not break auth;
- it does not break mobile;
- it has fallback behavior;
- it respects accessibility where applicable;
- it has no obvious console errors;
- it does not introduce legal/security risk;
- its status is documented;
- it can be explained in one sentence.

A launch asset is done only when:

- the claim is accurate;
- the CTA is clear;
- the visual matches the brand;
- legal/security claims are not exaggerated;
- it points to a working page or action;
- Bruno can use it without editing under pressure.

A post-launch fix is done only when:

- the issue is reproduced or reasonably diagnosed;
- the fix is integrated;
- the affected flow is retested;
- the user impact is documented;
- the next prevention step is noted.

---

## 14. Anti-Chaos Rules

1. **One active integration at a time.**
2. **Maximum two standby packets; prefer one.**
3. **No global refactor during a critical fix.**
4. **No "while we are here" changes.**
5. **No AI changes architecture unless Bruno explicitly asks for strategic reconsideration.**
6. **Every task must specify AI + Mode.**
7. **Every Claude task must specify forbidden areas.**
8. **Production domain is truth.**
9. **Claude summaries are not verification.**
10. **If Claude menu options appear, read the current menu every time; never assume a number means the same thing as before.**
11. **Marketing does not outrun product truth.**
12. **A launch blocker beats a beautiful idea.**
13. **A legal/security blocker beats a launch deadline.**
14. **If Bruno is overwhelmed, reduce to one next task.**
15. **Never trust an uncommitted diff without reading it** - touch-scroll input corruption is a confirmed, repeat repo-specific hazard, not a hypothetical (see 7.0).
16. **Concrete evidence beats a confident summary, always** - a `Set-Cookie` header, a `git diff`, a printed PID, or a `BUILD_OK` is proof; "fixed," "confirmed," or "done" in prose is a claim.
17. **Verify only against the real production domain** - a Vercel preview/alias URL can silently diverge from `bravebybruno.com` and must never be used to declare something PASS.

---

## 15. Non-Technical Operator Guide for Bruno

### 15.1 The Five Questions

Before doing anything, ask:

1. What am I trying to achieve?
2. Am I thinking, researching, looking, preparing, or building?
3. Which AI owns that mode?
4. What must not be touched?
5. What would prove this worked?

### 15.2 The One-Task Rule

If the next action cannot be written in one sentence, it is too broad.

Good:

> Prepare the Integration Packet for Cinematic Shell.

Too broad:

> Make the site look finished and premium and fix anything that is wrong.

### 15.3 The Safety Phrase

If confused, use this prompt:

```text
AI: ChatGPT
Mode: Strategy Mode
Objective: I am overwhelmed and need the next safest action.
Output expected: one next step, who to ask, exact prompt, and what not to touch.
Do not create a full roadmap unless I ask.
```

---

## 16. Immediate Next Actions

### 16.1 When Claude Code Is Next Available

Send Claude Code:

1. This v0.5 protocol document.
2. Instruction to save/adapt it inside the repo (see 16.2).
3. Instruction not to start new code work.
4. Immediately after saving, resume the AUTH GREEN verification exactly as specified in 7.1 - do not restart from zero, and do not accept a summary in place of the `Set-Cookie` header and `git diff` evidence still outstanding from tonight's session.

### 16.2 Claude Code Handoff Prompt

```text
AI: Claude Code
Mode: Repo Execution Mode
Objective: Save/adapt the Product Operating System + AI Technical Round Table Protocol v0.5 as an official project document.
Output expected: saved path, modifications made, confirmation that no feature code was started.
Do not touch / do not decide: auth, routing, middleware, backend, database, Stripe, env vars, Universe Engine, pricing, roadmap, or any feature implementation.

Claude, save this document as the official operating protocol for the Brave by Bruno / The Concierge project.

Suggested filename:
docs/protocol/PRODUCT_OPERATING_SYSTEM_v0.5.md (confirm the actual path/convention with a quick look at the repo root and any existing docs/ folder before committing to this path - report back whatever path you actually use).

Context:
- Auth/login is NOT yet confirmed resolved - it is mid-verification. See section 7 (Phase B0) for the exact outstanding test (Set-Cookie header with a real cai_token value from a genuine email+password login).
- The Master Document remains the strategic source of truth; the Concierge Codex Audited v1 remains the approved-feature source of truth.
- This protocol coordinates Bruno, Claude (chat), ChatGPT, Gemini, Perplexity Plus, and Claude Code - six roles, not five.
- Do not reopen locked decisions.
- You may only add repo-specific operational details beyond what v0.5 already includes: exact confirmed paths, build/test commands, branch/commit policy.
- After saving/adapting, resume the AUTH GREEN test from Phase B0. Do not start the Cinematic Shell until AUTH GREEN passes AND Bruno explicitly approves the next task.

Return:
1. saved path;
2. changes/additions made;
3. confirmation that no feature code was started;
4. result of the resumed AUTH GREEN test (Set-Cookie header content, or exactly why it's still blocked).
```

### 16.3 After Protocol Is Saved

Claude Code resumes AUTH GREEN using the exact procedure in Phase B0:

```text
AI: Claude Code
Mode: Repo Execution Mode
Objective: Complete the AUTH GREEN test from protocol Phase B0 and report PASS/FAIL with concrete evidence.
Output expected: Set-Cookie header content (real value), git diff of the _directToken fix, BUILD_OK confirmation, PASS/FAIL on commit, deploy, domain, login, session, dashboard, refresh, logout.
Do not touch / do not decide: new visual features, Cinematic Shell, pricing cards, dashboard restyle, backend expansion, Stripe, Supabase schema, global CSS.
```

### 16.4 After AUTH GREEN

Next task becomes:

```text
AI: ChatGPT
Mode: Execution / Computer Mode
Objective: Prepare the Integration Packet for Cinematic Shell - Three.js Cosmic Intro.
Output expected: complete Integration Packet with Next.js safety, WebGL fallback, reduced motion, asset media required, manual tests, Claude prompt.
Do not touch / do not decide: repo integration, auth, routing, middleware, backend, Stripe, Supabase, env vars, global CSS.
```

---

## 17. Final Operating Sentence

> Think with ChatGPT. Research with Perplexity. Look with Gemini. Relay and verify with Claude. Build with Claude Code. Decide with Bruno. Verify in production.

Expanded v0.5 version:

> Build in sequence. Launch with proof. Manage with discipline. Keep the experience simple in front, complex behind, and safe always.

This is the working model until formally changed.
