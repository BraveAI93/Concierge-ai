# AI Review Ledger — Trust Hardening and Cinematic Shell

**Project:** Brave by Bruno / The Concierge  
**Version:** v0.2  
**Date:** 2026-07-09  
**Owner:** Bruno Aversa  
**Status:** Decision log / documentation only. Not an implementation packet.
**Update in v0.2:** Adds Claude Chat final review and upgrades the fake “owner notified” UI issue to an immediate small blocker.

---

## 1. Current Decision

The current agreed direction is:

1. **Do not implement Brain Spine now.**
2. **Proceed toward Cinematic Shell next.**
3. **Treat Trust Hardening as a pre-launch blocker.**
4. **Treat the fake “owner notified” UI claim as an immediate small blocker:** either wire a real backend notification/audit record or remove/neutralise the claim before public/demo use.
5. **Verify AI-processing consent wording before finalising Cinematic Shell/landing copy.**
6. **Reopen Minimal Brain Spine after Cinematic Shell and Core Flow QA.**
7. **Do not make public trust/security/moderation/notification claims that the backend cannot currently prove.**

This ledger exists to prevent the AI assistants from reopening already-settled decisions or mixing roles.

---

## 2. Role Map

### Bruno
Final decision-maker and product owner. Approves or rejects changes.

### ChatGPT
Integrator and operational strategist. Turns Claude/Perplexity/Gemini outputs into decision packets, prompts, checklists, and roadmaps.

### Claude Code
Repo executor. Saves documentation, audits repo reality, commits/pushes, and implements only explicitly approved scopes.

### Perplexity
Compliance/security/evidence researcher. Supplies external research and best-practice context, not repo implementation decisions.

### Gemini
Visual/product-trust reviewer. Reviews visual direction, copy boundaries, UX honesty, and launch presentation implications.

---

## 3. Claude Code — Brain Spine Readiness Audit Summary

Claude Code produced `BRAIN_SPINE_READINESS_AUDIT.md`.

Key technical conclusion:

- The current live AI runtime is small and concentrated mostly around the Go backend `POST /chat` / `handleChat` flow.
- Minimal Brain Spine is feasible later, but should not be implemented before Cinematic Shell.
- Recommended roadmap: **Option B — audit/document now, Cinematic Shell next, Brain Spine later.**

Key trust findings:

1. The backend currently trusts a client-supplied `system_prompt`.
2. AI-processing consent in the chat banner is currently stored only in browser `sessionStorage`, not as a durable backend consent record.
3. Sensitive-topic alerts / “owner notified” / human-review UX are currently client-side/cosmetic and do not create a real backend notification or owner review record.

Operational conclusion:

- Brain Spine is not the next implementation task.
- Trust Hardening must be tracked separately before launch.

---

## 4. Perplexity — Compliance / Security Evidence

Perplexity reviewed best practices for AI concierge apps and confirmed that the safest and most defensible pattern is:

1. Keep system-prompt authority on the server.
2. Persist AI-processing consent as a real backend record when consent is the lawful basis.
3. Never display “owner notified” unless a real backend event happened and was logged.

Practical checklist from Perplexity:

- Build effective system prompts server-side from trusted templates, route type, user role, permissions, product mode, and policy version.
- Do not treat client `system_prompt` as authoritative.
- Validate any client-sent mode hints against an allowlist.
- Keep secrets, internal instructions, and guardrails out of browser/client-editable payloads wherever possible.
- Add server-side logging for prompt/template version, model, request ID, tool calls, moderation outcome, and high-risk action decisions.
- Persist AI-processing consent server-side with subject/session identifier, timestamp, consent text shown, consent version, purpose, capture method, and withdrawal state.
- Do not show “owner notified,” “sent to owner,” or “escalated to owner” unless backend notification/queue/audit record exists.
- Use truthful interim wording such as “flagged for review” or “ready to send” until backend confirmation exists.

Pre-launch fixes identified by Perplexity:

1. Server-owned prompts.
2. Persisted AI-processing consent.
3. Truthful notification claims.
4. Sensitive review logic moved off the client if intended as a genuine safeguard.

Items that can wait:

- Advanced indirect prompt-injection defenses for retrieved content/attachments/tools.
- Richer observability and red-team coverage.
- Consent-management UX history and renewal prompts.
- Full provenance separation and audit analytics.

---

## 5. Gemini — Visual / Product Trust Review

Gemini reviewed the audit from a visual/product-trust perspective.

Key conclusion:

- The Cinematic Shell can proceed visually because it is a brand/experience layer and does not itself worsen backend trust gaps.
- However, the copy inside the shell must not imply security, moderation, compliance, or notification systems that are not currently backed by the backend.

Landing-page copy implications:

- Do not claim “Full GDPR / UK-GDPR compliant AI processing” yet for the chat flow.
- Do not claim “fully moderated AI.”
- Do not claim “instant owner intervention.”
- Do not imply that the owner reviews every blocked/edited message.
- Do not imply system prompts are locked server-side until that is actually true.

Allowed visual direction:

- Premium.
- Cinematic.
- Cosmic / Universe metaphor.
- Creator-led.
- Intelligent conversational experience.
- Early/private version.
- Trust-aware roadmap.

Avoid public copy such as:

- “Enterprise-grade security.”
- “Tamper-proof.”
- “Fully GDPR compliant AI processing.”
- “Owner instantly notified.”
- “Fully moderated AI.”
- “Emergency alerts.”
- “Every blocked message is reviewed.”
- “Secure prompt vault.”
- “Permanent consent logging.”

Recommended honest alternatives:

- “AI-powered concierge interface.”
- “Premium conversational experience.”
- “Designed for creator-led businesses.”
- “Owner-controlled experience.”
- “Trust-aware roadmap.”
- “Flagged for review” instead of “owner notified,” unless real backend notification exists.

---

## 5A. Claude Chat — Final Review of Claude Code Audit

Claude Chat reviewed the Claude Code Brain Spine Readiness Audit as a separate assistant, not as repo executor.

Verdict:

1. **Agrees with Option B:** Cinematic Shell next, Brain Spine later. The audit evidence shows the AI surface is concentrated around `handleChat`, so delaying Brain Spine carries lower rework risk than feared.
2. **Agrees Trust Hardening is a pre-launch blocker**, but strengthens the priority on one item: the current UI claim that the owner was notified, when no real backend notification is triggered, is not merely a later engineering gap. It is an active false claim shown to users.
3. **Recommends an immediate small blocker:** either connect a real backend notification/audit event or remove/neutralise the “owner notified” language now.
4. **Adds a copy-safety dependency:** before finalising Cinematic Shell/landing copy, verify the GDPR/AI-processing consent gap so the new public page does not promise privacy or data-handling behaviour that is not currently implemented.

Operational update from Claude Chat:

- Brain Spine remains later.
- Cinematic Shell can proceed.
- The false notification claim should be fixed or neutralised independently, as a small trust patch.
- The AI-processing consent gap must constrain any public privacy/compliance copy.

---

## 6. Integrated Decision — ChatGPT

The integrated decision is:

```text
Docs push
→ AI Review Ledger push
→ Immediate trust-copy patch: remove/neutralise fake “owner notified” claim unless backend notification exists
→ Cinematic Shell Integration Packet with strict copy boundaries
→ Trust Hardening Mini-Packet before launch
→ Core Flow QA
→ Launch readiness
→ Minimal Brain Spine later
```

Brain Spine remains a future architecture layer. It is not the immediate task.

Trust Hardening is separate from Brain Spine. The app does not need the full Brain Spine to stop making false trust claims or to persist consent properly.

Immediate rule: public/demo UI must not say “owner notified,” “sent to owner,” or similar unless a real backend notification/audit event exists. If the backend event is not implemented yet, the copy should be changed to neutral language such as “flagged for review” or the toast should be disabled.

---

## 7. Open Decisions

These must be decided before launch:

1. **Resolved for now:** the chat must remove/neutralise “owner notified” immediately unless backend notification/audit exists.
2. Should server-owned prompt templates be fixed before or after Cinematic Shell?
3. Should AI-processing consent be persisted before any public demo with real users?
4. What exact public copy is allowed in the first live version?
5. Should Trust Hardening be its own Claude Code packet before Core Flow QA?

---

## 8. Next Packets

### Packet 0 — Immediate Trust-Copy Patch

Purpose:
Remove or neutralise any UI text that falsely says the owner was notified/escalated when no backend notification/audit record exists.

Likely scope:

- Search for text such as “owner notified”, “sent to owner”, “escalated”, or similar.
- Replace with truthful copy such as “flagged for review” only if accurate, or disable the toast.
- Do not implement the full backend notification system in this tiny patch unless explicitly approved.
- Do not start Brain Spine.
- Do not touch auth, Stripe, provider integrations, or schema unless separately approved.

### Packet 1 — Cinematic Shell Integration Packet

Purpose:
Restore/build the premium cinematic visual shell while avoiding backend trust/security claims.

Hard boundaries:

- Do not touch auth.
- Do not touch `concierge-backend/main.go` unless explicitly approved later.
- Do not touch `lib/buildPrompt.js`.
- Do not modify `components/Chat.jsx` trust logic inside the visual packet.
- Do not start Brain Spine.
- Do not add new providers.
- Do not add weather/GPS/API keys in V1.
- Do not make privacy, GDPR, moderation, security, owner-notification, or consent claims unless the backend already supports them.

### Packet 2 — Trust Hardening Mini-Packet

Purpose:
Fix or neutralise trust-relevant gaps before launch.

Likely scope:

1. Server-owned prompts.
2. Persist AI-processing consent.
3. Replace fake notification language or add backend notification record.
4. Clarify review/alert UX honestly.

### Packet 3 — Minimal Brain Spine

Purpose:
Introduce clean AI architecture later, after Cinematic Shell and Core Flow QA.

Not now.

---

## 9. New Visual Feature Idea — Location-Aware Ambient Theme

Feature name proposal:

- `Location-Aware Ambient Theme`
- `City Weather Visual Mode`

Classification:

- Cinematic Shell / dashboard visual layer.
- Not Brain Spine.
- Not Trust Hardening.

V1 boundaries:

- Use city selected in profile/business page only.
- No GPS.
- No real-time user geolocation.
- Day/night state can be static or browser-time based.
- Weather widget can be placeholder/mock only.
- No live weather API.
- No new provider keys.
- No privacy-sensitive location logic.

V2 possibilities:

- Live weather API.
- User/location consent.
- Rich city-specific visual widgets.
- Local time and atmospheric background.

Example:

```text
Bruno — London
Night mode
Cloudy visual ambience
Cosmic city background + weather-style widget + local time
```

---

## 10. Update Protocol

This ledger should be updated append-only.

Rules:

1. Do not rewrite previous decisions unless explicitly superseded.
2. Add dated sections for new AI reviews.
3. Label source clearly: Claude Code / Perplexity / Gemini / ChatGPT / Bruno.
4. Keep final decisions separate from raw assistant opinions.
5. Claude Code may save and commit this document, but should not invent missing assistant content.
6. If a source output was not provided to Claude Code, mark it as `PENDING`, not fabricated.

