# The Concierge Fleet Architecture Brief

## Purpose

This brief proposes the target fleet architecture for **The Concierge**: a premium assistant-manager platform designed to feel like one unified private chief of staff while using a governed fleet of frontier models, specialist services, tools, and memory systems behind the scenes.[cite:98][cite:12][cite:101]

The core recommendation is to evolve the current system from a single-model assistant into a **hierarchical, policy-governed, multi-model fleet** with one user-facing identity, one deterministic control plane, one platform-owned memory spine, one MCP-centric integration fabric, and multiple specialist lanes invoked only when appropriate.[cite:76][cite:12][cite:98]

## Executive Summary

The strongest architecture for The Concierge is **not** a flat swarm of agents and **not** a simple “add more APIs” approach. The best architecture is a layered hierarchy in which orchestration, memory, governance, and tool policy belong to the platform, while models act as replaceable cognitive engines selected according to role, risk, cost, and context.[cite:12][cite:98][cite:101]

This architecture is the strongest option for five reasons:[cite:12][cite:76][cite:101]

- It keeps the **user experience unified**, so the product feels like one Concierge rather than many disconnected assistants.[cite:98]
- It improves **reliability**, because routing, fallback, and workflow control are explicit rather than improvised.[cite:76][cite:109]
- It improves **cost discipline**, because premium models are reserved for judgment-heavy tasks while narrow subtasks are pushed to cheaper specialist lanes.[cite:83][cite:41][cite:6]
- It improves **personalization**, because memory and provenance stay inside the platform rather than inside one provider’s context window.[cite:98][cite:12]
- It improves **evolvability**, because new lanes, tools, and providers can be added through a registry-governed framework rather than by reworking the whole product each time.[cite:101][cite:77]

## Strategic Design Principle

The governing principle should be:

**One identity in front. One control plane in charge. Many engines underneath. One final judgment layer before the user sees the answer.**[cite:98][cite:76]

This principle matters because a concierge product is judged at the experience layer, not at the model roster layer. The user should perceive continuity, taste, memory, initiative, and trustworthiness, while the backend remains modular and replaceable.[cite:98][cite:12]

## Why a Hierarchical Fleet Is Best

Microsoft’s multi-agent reference architecture presents modular, governed systems with an orchestrator, registry, memory/RAG, tool mediation, and support for multiple workflow patterns rather than one flat universal pattern.[cite:98][cite:12] Microsoft’s workflow-oriented guidance also emphasizes deterministic control for sequencing, compliance, and auditability in multi-agent workflows.[cite:109]

That makes a **hierarchical fleet** the best objective design for The Concierge because different requests require different execution patterns:[cite:12][cite:109]

- A simple reminder should take the shortest, cheapest route.[cite:12]
- A complex planning task may require decomposition and synthesis.[cite:12]
- A high-stakes recommendation may need verification or multi-lane review before being shown to the user.[cite:101]
- A tool-driven workflow may need deterministic gating, permissions, and approval checkpoints.[cite:76][cite:109]

A flat agent swarm makes these flows harder to predict, harder to govern, and harder to cost-control. A hierarchical fleet makes them easier to reason about and easier to evolve.[cite:76][cite:101]

## Target Architecture

### 1. User Identity Layer

This is the layer the user experiences directly. It should contain:[cite:98]

- The Concierge persona.
- Tone and stylistic rules.
- Response presentation logic.
- Session continuity.
- User-visible assurance modes such as fast, deep, and triple-checked.[cite:98]

This layer should never expose internal model switching. The product must feel like one coherent premium assistant at all times.[cite:98]

### 2. Policy and Control Plane

This is the true operating core of the platform. It should own:[cite:76][cite:101]

- Routing policy.
- Risk classification.
- Latency and cost budgets.
- Permission and compliance rules.
- Fallback and escalation logic.
- Approval requirements for sensitive actions.
- Response mode policy, including when to use fast, deep, or verified workflows.[cite:76][cite:109]

This layer should be deterministic first. Conductor’s workflow model is relevant because it frames orchestration as explicit logic rather than another LLM trying to invent the entire control flow at runtime.[cite:76]

### 3. Supervisor / Planner Layer

This layer decides workflow shape:[cite:12][cite:109]

- Single-lane response.
- Sequential pipeline.
- Parallel multi-lane review.
- Verification flow.
- Human-gated action flow.[cite:12]

It should **not** decompose every request by default. Its purpose is to decide when decomposition adds value and when it would only add latency and inconsistency.[cite:12][cite:109]

### 4. Specialist Fleet Layer

The fleet should be organized by **role**, not by vendor. Vendors can change over time; roles remain stable.[cite:98][cite:101]

#### Premium reasoning lane

This lane should use frontier models for deep cognition, nuanced drafting, planning, and final synthesis.[cite:98]

- **Claude**: deep reasoning, structured plans, long-form refinement, nuanced writing.[cite:98]
- **OpenAI**: structured outputs, orchestration-friendly reasoning, tool-enabled agent flows, and final judgment when needed. OpenAI’s Responses API supports remote MCP servers and integrated tools, making it strong for production agent workflows.[cite:26][cite:38]

#### Live research lane

This lane should use **Perplexity / Computer** for current, sourced, web-grounded research and evidence gathering. Market coverage in 2026 describes Perplexity as operating a multi-model research system for high-quality research synthesis.[cite:28][cite:40]

#### Multimodal / visual lane

This lane should use **Gemini** for screenshot analysis, visual review, document interpretation, and interface-aware reasoning where seeing matters, not just reading.[cite:98]

#### Fast classifier / router lane

This lane should use a low-cost model for intent classification, route selection, low-risk transformation, and guardrail filtering. Current cost-routing guidance shows that many production systems reduce spend by routing routine work to smaller models and keeping premium reasoning for complex tasks.[cite:41][cite:83]

#### OCR / document lane

This lane should use specialist document extraction services rather than premium reasoning models for every PDF or receipt. Mistral’s API lineup, for example, includes OCR 4 at $4 per 1,000 pages and Document AI at $5 per 1,000 pages, showing the economic logic of a separate extraction lane.[cite:6]

#### Voice / transcription lane

This lane should handle voice notes, spoken commands, and meeting capture before escalating to higher reasoning layers. Mistral’s API pricing lists transcription at roughly $0.003 to $0.006 per minute depending on model, which is much more appropriate than sending raw audio into premium orchestration paths.[cite:6]

#### Retrieval lane

This lane should own embeddings, search, reranking, and memory recall. Microsoft’s reference architecture treats knowledge and RAG as platform components rather than side features, which is especially important for a personalized concierge system.[cite:98][cite:12]

Open multimodal retrieval tooling is becoming more capable as well. For example, current reporting on Qwen3-VL embedding and reranking positions these models for unified search across text, screenshots, images, and video workflows.[cite:46][cite:54]

#### Safety / moderation lane

This lane should perform moderation, policy checks, and basic risk filtering using dedicated classifiers. This keeps safety checks cheap, consistent, and inspectable.[cite:6]

#### Code / integration lane

This lane should support connectors, scripts, schemas, integration helpers, and MCP server work. Mistral’s public API lineup includes coding-focused models such as Devstral 2, Devstral Small 2, and Codestral, illustrating the value of a dedicated coding lane instead of sending engineering tasks through user-facing reasoning paths.[cite:6]

### 5. Tool and Integration Layer

This layer should be **MCP-first**. Microsoft’s architecture describes MCP as a decoupled way for agents to invoke tools with validated parameters, policy mediation, and structured returns rather than raw ad hoc tool calls.[cite:12]

This layer should include:

- Calendar, reminders, tasks.
- Email, messaging, communications.
- Search, browser, and file tools.
- Travel, booking, and commerce partners.
- CRM and business systems.
- Internal operational APIs.[cite:12][cite:38]

A centralized tool fabric is important because tool sprawl is otherwise one of the main sources of prompt bloat, risk, and failure in agent systems.[cite:12][cite:71]

### 6. Memory and Knowledge Layer

This layer should be platform-owned and split into separate stores:[cite:98][cite:12]

- **Profile memory**: preferences, habits, tone, priorities, recurring relationships.
- **Operational memory**: active plans, pending actions, commitments, workflow state.
- **Knowledge memory**: notes, uploads, approved references, structured entities.
- **Evidence and provenance memory**: source links, verification logs, policy outcomes, action history.[cite:98][cite:101]

This is strategically critical. The long-term moat of The Concierge will come less from access to any one model and more from memory quality, continuity, trust structure, and how those memories are turned into action.[cite:98][cite:12]

### 7. Observability, Evaluation, and Governance Layer

This layer should trace and govern every lane, workflow, and model decision. The registry-governed agent lifecycle paper argues that agent operations should be governed through evaluation-driven promotion, retirement, and registry-based discovery rather than one-time benchmark selection.[cite:101][cite:77]

This layer should capture:

- Trace IDs per workflow.
- Cost per route, lane, and completed outcome.
- Latency per step.
- Fallback path taken.
- Evaluation results.
- Human override points.
- Failure patterns and replayable diagnostics.[cite:71][cite:101]

This is not optional in a serious agent fleet. Current 2026 guidance on agent observability consistently frames cost control and quality control as observability problems.[cite:71][cite:82][cite:84]

## Recommended Execution Hierarchy

The default request flow should be:

1. **Ingress**
2. **Mode and risk classification**
3. **Deterministic shortcut check**
4. **Supervisor decision**
5. **Specialist execution**
6. **Premium synthesis**
7. **Optional validator / verifier**
8. **Action gate for external side effects**
9. **Response delivery + memory update + trace log**[cite:109][cite:76][cite:101]

This structure gives the platform three major strengths:[cite:12][cite:76][cite:101]

- Easy tasks stay cheap and fast.[cite:12]
- Complex tasks get the right orchestration depth.[cite:109]
- High-stakes tasks can be verified before any answer or action is surfaced.[cite:101]

## Provider and Role Map

| Layer / role | Recommended owner | Why |
|---|---|---|
| User-facing Concierge | The platform | Must remain vendor-independent and consistent.[cite:98] |
| Policy/control plane | The platform | Core moat and governance logic should be owned, not outsourced.[cite:101] |
| Deterministic orchestration | Workflow engine + rules | More auditability and lower token waste than improvised orchestration.[cite:76][cite:109] |
| Deep reasoning / drafting | Claude | Strong structured reasoning, refinement, and long-form synthesis.[cite:98] |
| Tool-aware reasoning / structured execution | OpenAI | Strong tool support, remote MCP, structured outputs, and workflow integration.[cite:26][cite:38] |
| Current live research | Perplexity / Computer | Best fit for sourced, current, market-aware research.[cite:28][cite:40] |
| Visual / multimodal inspection | Gemini | Strong multimodal analysis role in the fleet.[cite:98] |
| Classification / moderation / first-pass transforms | Small or specialist models | Better unit economics for frequent low-risk tasks.[cite:41][cite:83] |
| OCR / transcription / embeddings / reranking | Specialist services | Better fit and cost for narrow technical subtasks.[cite:6][cite:46] |
| Memory / provenance | The platform | Source of personalization, trust, and continuity.[cite:98][cite:12] |
| Registry / evaluation / observability | The platform | Needed for evidence-driven lifecycle management.[cite:101][cite:71] |

## Why Organize by Role, Not by Provider

A provider-first architecture becomes brittle over time because providers change pricing, rate limits, quality, governance features, and strengths. A role-first architecture lets the platform preserve its own internal logic while swapping the best engine into a given role as the market evolves.[cite:101][cite:83]

This matters especially in a fast-moving 2026 market, where multi-model routing and runtime selection are already being presented as ways to reduce cost while maintaining or improving outcome quality.[cite:83][cite:42][cite:41]

## Open-Source and Specialist Lanes

Open-source and niche specialist lanes should be added under the same governance framework, not bolted on as separate mini-products. The most useful categories for The Concierge are:[cite:6][cite:46][cite:49]

- Cheap general text models for classification and low-risk transformations.[cite:6]
- OCR and document extraction lanes for forms, invoices, receipts, and PDFs.[cite:6]
- Transcription lanes for voice workflows.[cite:6]
- Retrieval lanes for embeddings and reranking, including multimodal retrieval support.[cite:46][cite:54]
- Moderation lanes for policy and content filtering.[cite:6]
- Coding lanes for scripts, connectors, and automation helpers.[cite:6]

Self-hosting becomes attractive later when volumes justify it. Current self-hosted inference guidance emphasizes hybrid routing and gateway patterns, especially for stable low-cost workloads such as embeddings, classification, and background automation.[cite:49]

## Cost Logic

The platform’s cost strategy should follow a simple rule: **use the most expensive cognition only where it produces meaningful user value**.[cite:83][cite:41]

OpenAI’s current pricing documentation shows large spreads between flagship and smaller models, with GPT-5.5 listed at $5 input and $30 output per 1M tokens, GPT-5.4 at $2.5 input and $15 output, GPT-5.4-mini at $0.75 input and $4.5 output, and GPT-5.4-nano at $0.20 input and $1.25 output.[cite:7] Mistral’s current API pricing shows much cheaper specialist and small-model options, such as Mistral Small 4 at $0.15 input and $0.6 output per 1M tokens, OCR 4 at $4 per 1,000 pages, and transcription as low as $0.003 per minute.[cite:6]

This suggests the best economic pattern for The Concierge:[cite:41][cite:83][cite:6]

- Reserve premium reasoning lanes for synthesis, judgment, nuance, and high-value user interactions.[cite:83]
- Route extraction, filtering, transcription, and retrieval to specialist lanes.[cite:6]
- Keep routing and moderation cheap and fast.[cite:41]
- Use observability to measure cost per resolved outcome, not just cost per token.[cite:71][cite:84]

## Recommended Rollout Phases

### Phase 1: Build the spine

Build first:[cite:76][cite:98][cite:101]

- Policy engine.
- Deterministic workflow orchestration.
- Unified tracing.
- Platform-owned memory and provenance.
- Initial registry for lanes and tools.[cite:76][cite:101]

This phase matters because adding more models before building the spine only increases chaos.[cite:76]

### Phase 2: Add the premium fleet

Add next:[cite:98][cite:26][cite:28]

- Claude lane.
- OpenAI lane.
- Perplexity research lane.
- Gemini multimodal lane.[cite:98]

This phase establishes the flagship cognition profile of the product.[cite:98]

### Phase 3: Add specialist utility lanes

Add after the premium fleet:[cite:6][cite:46]

- OCR.
- Transcription.
- Embeddings and reranking.
- Moderation.
- Cheap classifier / router.[cite:6]

This phase improves cost control, modality coverage, and operational quality.[cite:6][cite:41]

### Phase 4: Build the action fabric

Add next:[cite:12][cite:109]

- MCP gateway.
- Tool permission policies.
- External integrations.
- Approval gates for sensitive actions.[cite:12]

This phase turns cognitive quality into reliable execution.[cite:109]

### Phase 5: Governance maturity

Add after stable operation:[cite:101][cite:71]

- Evaluation harness.
- Promotion / retirement rules.
- A/B routing tests.
- Cost-quality dashboards.
- Lane-level scorecards.[cite:101]

This phase makes the architecture durable and self-improving.[cite:101]

## Non-Negotiable Principles

The following should be treated as architectural constants:[cite:98][cite:76][cite:101]

- One user-facing Concierge identity.[cite:98]
- Platform-owned memory and provenance.[cite:98][cite:12]
- Deterministic orchestration before free-form swarm behavior.[cite:76][cite:109]
- Role-based lanes, not provider chaos.[cite:101]
- Registry, observability, and evaluations from the start.[cite:101][cite:71]
- Premium models used for judgment, not every subtask.[cite:83][cite:6]
- MCP or equivalent governed tool mediation as the default integration path.[cite:12][cite:38]

## Message for Claude, Gemini, and ChatGPT

This proposal should be discussed with each system according to its likely role inside the fleet.[cite:98][cite:26]

### For Claude

Focus on the **premium reasoning lane**:[cite:98]

- Deep planning.
- Structured decomposition when needed.
- Nuanced drafting.
- Long-form synthesis.
- Final response quality in deep mode.[cite:98]

### For Gemini

Focus on the **multimodal / visual lane**:[cite:98]

- Screenshot interpretation.
- Interface analysis.
- Visual document review.
- Image-based context extraction.[cite:98]

### For ChatGPT / OpenAI

Focus on the **control-plane-adjacent lane**:[cite:26][cite:38]

- Tool orchestration.
- Structured outputs.
- MCP integration.
- Workflow-compatible reasoning.
- Operational robustness for agent execution.[cite:26][cite:38]

## Final Recommendation

The best way to add this to the roadmap is **not** as a list of APIs to integrate. It should be presented as a product operating architecture with a clear hierarchy:[cite:98][cite:76][cite:101]

1. Identity layer.
2. Policy/control plane.
3. Supervisor/planner.
4. Specialist fleet.
5. Tool/MCP fabric.
6. Memory and provenance spine.
7. Observability, evaluation, and governance.[cite:98][cite:101]

That structure gives The Concierge the best chance of becoming a premium flagship assistant-manager platform rather than a loose bundle of AI features. It creates a system that is stronger technically, clearer strategically, safer operationally, and more defensible commercially.[cite:98][cite:76][cite:101]
