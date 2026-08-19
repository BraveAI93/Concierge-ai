-- CENTRAL INTELLIGENCE RUNTIME v0.3 — REVIEW-ONLY UP MIGRATION
--
-- DO NOT EXECUTE. This file is intentionally outside active migration tooling.
-- It has not been applied to Supabase or any database. See README.md for RLS,
-- identity, retention/deletion, rollback, and staging-activation requirements.
--
-- All UUIDs are application generated; this design assumes no database UUID
-- extension. The future activation must execute the complete file in one
-- reviewed transaction against a disposable staging project first.

BEGIN;

CREATE SCHEMA IF NOT EXISTS ci_kernel;

-- Canonical identity boundary. stable_subject and source_profile_id are
-- internal server-resolved values, never browser-selected public slugs.
CREATE TABLE ci_kernel.persons (
  id uuid PRIMARY KEY,
  display_name text NOT NULL,
  created_at timestamptz NOT NULL
);

CREATE TABLE ci_kernel.person_bindings (
  person_id uuid PRIMARY KEY REFERENCES ci_kernel.persons(id),
  stable_subject text NOT NULL UNIQUE,
  source_profile_id text NOT NULL UNIQUE,
  created_at timestamptz NOT NULL,
  retired_at timestamptz
);

-- Source idempotency and immutable source references. Source content is
-- personal data and must follow retention/deletion policy before activation.
CREATE TABLE ci_kernel.runtime_sources (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  source_profile_id text NOT NULL,
  conversation_id text NOT NULL,
  session_id text,
  message_id text NOT NULL,
  message_role text NOT NULL,
  content text NOT NULL,
  conversation_at timestamptz,
  message_at timestamptz NOT NULL,
  stored_at timestamptz NOT NULL,
  UNIQUE (person_id, message_id)
);

-- Four temporal clocks remain explicit. event_at is semantic event time;
-- recorded_at is knowledge-ingestion time and may precede a future event_at.
CREATE TABLE ci_kernel.events (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  kind text NOT NULL,
  summary text NOT NULL,
  context_ids jsonb NOT NULL DEFAULT '[]'::jsonb,
  evidence_ids jsonb NOT NULL DEFAULT '[]'::jsonb,
  event_at timestamptz NOT NULL,
  recorded_at timestamptz NOT NULL,
  effective_at timestamptz NOT NULL,
  attention_at timestamptz NOT NULL,
  expires_at timestamptz,
  created_at timestamptz NOT NULL,
  UNIQUE (person_id, id),
  CHECK (expires_at IS NULL OR expires_at >= effective_at)
);

CREATE TABLE ci_kernel.memories (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  kind text NOT NULL,
  summary text NOT NULL,
  context_ids jsonb NOT NULL DEFAULT '[]'::jsonb,
  event_at timestamptz NOT NULL,
  recorded_at timestamptz NOT NULL,
  effective_at timestamptz NOT NULL,
  attention_at timestamptz NOT NULL,
  expires_at timestamptz,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  UNIQUE (person_id, id),
  CHECK (expires_at IS NULL OR expires_at >= effective_at)
);

CREATE TABLE ci_kernel.claims (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  memory_id text NOT NULL REFERENCES ci_kernel.memories(id),
  statement text NOT NULL,
  subject_id text,
  predicate text,
  object_value text,
  confidence_score numeric(4,3) NOT NULL DEFAULT 0.500 CHECK (confidence_score BETWEEN 0 AND 1),
  stale_after_seconds bigint NOT NULL DEFAULT 0 CHECK (stale_after_seconds >= 0),
  last_validated_at timestamptz,
  last_revalidated_at timestamptz,
  freshness_status text NOT NULL DEFAULT 'fresh' CHECK (freshness_status IN ('fresh','stale','historical','superseded')),
  supersedes_claim_id text REFERENCES ci_kernel.claims(id),
  event_at timestamptz NOT NULL,
  recorded_at timestamptz NOT NULL,
  effective_at timestamptz NOT NULL,
  attention_at timestamptz NOT NULL,
  expires_at timestamptz,
  created_at timestamptz NOT NULL,
  UNIQUE (person_id, id),
  CHECK (expires_at IS NULL OR expires_at >= effective_at)
);

CREATE TABLE ci_kernel.evidence (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  claim_id text NOT NULL REFERENCES ci_kernel.claims(id),
  stance text NOT NULL CHECK (stance IN ('supports','contradicts','neutral')),
  summary text NOT NULL,
  quality numeric(4,3) NOT NULL CHECK (quality BETWEEN 0 AND 1),
  relevance numeric(4,3) NOT NULL CHECK (relevance BETWEEN 0 AND 1),
  authority numeric(4,3) NOT NULL DEFAULT 0 CHECK (authority BETWEEN 0 AND 1),
  source_type text NOT NULL,
  source_ref text NOT NULL,
  actor text,
  captured_at timestamptz NOT NULL,
  checksum text,
  event_at timestamptz NOT NULL,
  recorded_at timestamptz NOT NULL,
  effective_at timestamptz NOT NULL,
  attention_at timestamptz NOT NULL,
  expires_at timestamptz,
  created_at timestamptz NOT NULL,
  UNIQUE (person_id, id),
  CHECK (expires_at IS NULL OR expires_at >= effective_at)
);

CREATE TABLE ci_kernel.memory_event_links (
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  memory_id text NOT NULL REFERENCES ci_kernel.memories(id),
  event_id text NOT NULL REFERENCES ci_kernel.events(id),
  reason text NOT NULL,
  linked_at timestamptz NOT NULL,
  PRIMARY KEY (person_id, memory_id, event_id)
);

CREATE TABLE ci_kernel.claim_lineage (
  claim_id text PRIMARY KEY REFERENCES ci_kernel.claims(id),
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  supersedes_claim_id text REFERENCES ci_kernel.claims(id),
  evidence_ids jsonb NOT NULL DEFAULT '[]'::jsonb,
  preserves_history boolean NOT NULL DEFAULT true,
  recorded_at timestamptz NOT NULL
);

CREATE TABLE ci_kernel.goals (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  title text NOT NULL,
  description text NOT NULL,
  subjective_importance numeric(4,3) NOT NULL CHECK (subjective_importance BETWEEN 0 AND 1),
  status text NOT NULL CHECK (status IN ('active','paused','achieved','abandoned')),
  event_at timestamptz NOT NULL,
  recorded_at timestamptz NOT NULL,
  effective_at timestamptz NOT NULL,
  attention_at timestamptz NOT NULL,
  expires_at timestamptz,
  created_at timestamptz NOT NULL,
  UNIQUE (person_id, id),
  CHECK (expires_at IS NULL OR expires_at >= effective_at)
);

CREATE TABLE ci_kernel.constraints (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  kind text NOT NULL CHECK (kind IN ('hard','soft')),
  title text NOT NULL,
  description text NOT NULL,
  active boolean NOT NULL,
  event_at timestamptz NOT NULL,
  recorded_at timestamptz NOT NULL,
  effective_at timestamptz NOT NULL,
  attention_at timestamptz NOT NULL,
  expires_at timestamptz,
  created_at timestamptz NOT NULL,
  UNIQUE (person_id, id),
  CHECK (expires_at IS NULL OR expires_at >= effective_at)
);

CREATE TABLE ci_kernel.pending_intents (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  summary text NOT NULL,
  state text NOT NULL CHECK (state IN ('captured','clarifying','ready','proposed','in_progress','completed','cancelled','expired')),
  goal_id text REFERENCES ci_kernel.goals(id),
  memory_id text REFERENCES ci_kernel.memories(id),
  context_ids jsonb NOT NULL DEFAULT '[]'::jsonb,
  event_at timestamptz NOT NULL,
  recorded_at timestamptz NOT NULL,
  effective_at timestamptz NOT NULL,
  attention_at timestamptz NOT NULL,
  expires_at timestamptz,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  UNIQUE (person_id, id),
  CHECK (expires_at IS NULL OR expires_at >= effective_at)
);

CREATE TABLE ci_kernel.pending_intent_transitions (
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  pending_intent_id text NOT NULL REFERENCES ci_kernel.pending_intents(id),
  sequence integer NOT NULL,
  from_state text NOT NULL,
  to_state text NOT NULL,
  occurred_at timestamptz NOT NULL,
  actor text,
  reason text NOT NULL,
  PRIMARY KEY (person_id, pending_intent_id, sequence)
);

CREATE TABLE ci_kernel.open_loops (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  pending_intent_id text NOT NULL REFERENCES ci_kernel.pending_intents(id),
  label text NOT NULL,
  context_ids jsonb NOT NULL DEFAULT '[]'::jsonb,
  entity_ids jsonb NOT NULL DEFAULT '[]'::jsonb,
  last_interaction_at timestamptz,
  interaction_observed_at timestamptz,
  estimated_effort_seconds bigint NOT NULL DEFAULT 0 CHECK (estimated_effort_seconds >= 0),
  estimated_attention_seconds bigint NOT NULL DEFAULT 0 CHECK (estimated_attention_seconds >= 0),
  interruption_cost numeric(4,3) NOT NULL DEFAULT 0 CHECK (interruption_cost BETWEEN 0 AND 1),
  context_switch_cost numeric(4,3) NOT NULL DEFAULT 0 CHECK (context_switch_cost BETWEEN 0 AND 1),
  event_at timestamptz NOT NULL,
  recorded_at timestamptz NOT NULL,
  effective_at timestamptz NOT NULL,
  attention_at timestamptz NOT NULL,
  expires_at timestamptz,
  resolved_at timestamptz,
  created_at timestamptz NOT NULL,
  UNIQUE (person_id, id),
  CHECK (expires_at IS NULL OR expires_at >= effective_at)
);

CREATE TABLE ci_kernel.attention_budgets (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  window_start timestamptz NOT NULL,
  window_end timestamptz NOT NULL CHECK (window_end > window_start),
  attention_capacity_seconds bigint NOT NULL CHECK (attention_capacity_seconds >= 0),
  max_competing_items integer NOT NULL CHECK (max_competing_items > 0),
  interruption_cost numeric(4,3) NOT NULL DEFAULT 0 CHECK (interruption_cost BETWEEN 0 AND 1),
  current_context_ids jsonb NOT NULL DEFAULT '[]'::jsonb,
  current_entity_ids jsonb NOT NULL DEFAULT '[]'::jsonb,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  UNIQUE (person_id, id),
  CHECK (attention_capacity_seconds <= EXTRACT(EPOCH FROM (window_end-window_start)))
);

CREATE TABLE ci_kernel.attention_allocations (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  attention_budget_id text NOT NULL REFERENCES ci_kernel.attention_budgets(id),
  open_loop_id text NOT NULL REFERENCES ci_kernel.open_loops(id),
  selection_state text NOT NULL CHECK (selection_state IN ('surfaced','deferred')),
  policy_score numeric(4,3),
  context_matched boolean NOT NULL DEFAULT false,
  attention_reserved_seconds bigint NOT NULL DEFAULT 0 CHECK (attention_reserved_seconds >= 0),
  reason text NOT NULL,
  evaluated_at timestamptz NOT NULL,
  UNIQUE (person_id, attention_budget_id, open_loop_id)
);

CREATE TABLE ci_kernel.opportunities (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  title text NOT NULL,
  summary text NOT NULL,
  goal_ids jsonb NOT NULL DEFAULT '[]'::jsonb,
  constraint_ids jsonb NOT NULL DEFAULT '[]'::jsonb,
  evidence_ids jsonb NOT NULL DEFAULT '[]'::jsonb,
  subjective_importance numeric(4,3) NOT NULL CHECK (subjective_importance BETWEEN 0 AND 1),
  objective_stakes numeric(4,3) NOT NULL CHECK (objective_stakes BETWEEN 0 AND 1),
  expected_impact numeric(4,3) NOT NULL CHECK (expected_impact BETWEEN 0 AND 1),
  reversibility numeric(4,3) NOT NULL CHECK (reversibility BETWEEN 0 AND 1),
  uncertainty numeric(4,3) NOT NULL CHECK (uncertainty BETWEEN 0 AND 1),
  opportunity_cost numeric(4,3) NOT NULL CHECK (opportunity_cost BETWEEN 0 AND 1),
  effort_attention_cost numeric(4,3) NOT NULL CHECK (effort_attention_cost BETWEEN 0 AND 1),
  latest_safe_action_at timestamptz,
  estimated_effort_seconds bigint NOT NULL DEFAULT 0 CHECK (estimated_effort_seconds >= 0),
  estimated_attention_seconds bigint NOT NULL DEFAULT 0 CHECK (estimated_attention_seconds >= 0),
  utility numeric(4,3),
  hard_blocked boolean NOT NULL DEFAULT false,
  priority_mismatch_must_surface boolean NOT NULL DEFAULT false,
  priority_mismatch_reason text,
  deadline_feasibility text CHECK (deadline_feasibility IN ('no_deadline','feasible','infeasible')),
  decision_basis text,
  evaluated_at timestamptz,
  event_at timestamptz NOT NULL,
  recorded_at timestamptz NOT NULL,
  effective_at timestamptz NOT NULL,
  attention_at timestamptz NOT NULL,
  expires_at timestamptz,
  created_at timestamptz NOT NULL,
  UNIQUE (person_id, id),
  CHECK (expires_at IS NULL OR expires_at >= effective_at)
);

CREATE TABLE ci_kernel.decisions (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  opportunity_id text NOT NULL REFERENCES ci_kernel.opportunities(id),
  kind text NOT NULL CHECK (kind IN ('recommend','defer','decline','must_surface')),
  utility numeric(4,3) NOT NULL CHECK (utility BETWEEN 0 AND 1),
  reason text NOT NULL,
  created_at timestamptz NOT NULL,
  UNIQUE (person_id, id)
);

CREATE TABLE ci_kernel.permissions (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  requires_approval boolean NOT NULL,
  can_auto_approve boolean NOT NULL,
  granted_by text NOT NULL,
  effective_at timestamptz NOT NULL,
  expires_at timestamptz,
  created_at timestamptz NOT NULL,
  UNIQUE (person_id, id),
  CHECK (expires_at IS NULL OR expires_at >= effective_at)
);

CREATE TABLE ci_kernel.permission_scopes (
  permission_id text NOT NULL REFERENCES ci_kernel.permissions(id),
  capability text NOT NULL,
  resource text NOT NULL,
  PRIMARY KEY (permission_id, capability, resource)
);

CREATE TABLE ci_kernel.action_proposals (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  opportunity_id text NOT NULL REFERENCES ci_kernel.opportunities(id),
  decision_id text NOT NULL REFERENCES ci_kernel.decisions(id),
  permission_id text NOT NULL REFERENCES ci_kernel.permissions(id),
  title text NOT NULL,
  requested_capability text NOT NULL,
  requested_resource text NOT NULL,
  parameters jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL,
  UNIQUE (person_id, id)
);

CREATE TABLE ci_kernel.action_gates (
  id text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  action_proposal_id text NOT NULL UNIQUE REFERENCES ci_kernel.action_proposals(id),
  state text NOT NULL CHECK (state IN ('draft','awaiting_approval','approved','rejected','expired','executed')),
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  UNIQUE (person_id, id)
);

CREATE TABLE ci_kernel.action_gate_transitions (
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  action_gate_id text NOT NULL REFERENCES ci_kernel.action_gates(id),
  sequence integer NOT NULL,
  from_state text NOT NULL,
  to_state text NOT NULL,
  occurred_at timestamptz NOT NULL,
  actor text,
  reason text NOT NULL,
  PRIMARY KEY (person_id, action_gate_id, sequence)
);

CREATE TABLE ci_kernel.runtime_replays (
  idempotency_key text PRIMARY KEY,
  person_id uuid NOT NULL REFERENCES ci_kernel.persons(id),
  event_id text NOT NULL REFERENCES ci_kernel.events(id),
  evidence_id text NOT NULL REFERENCES ci_kernel.evidence(id),
  memory_id text NOT NULL REFERENCES ci_kernel.memories(id),
  claim_id text NOT NULL REFERENCES ci_kernel.claims(id),
  intent_id text REFERENCES ci_kernel.pending_intents(id),
  open_loop_id text REFERENCES ci_kernel.open_loops(id),
  opportunity_id text REFERENCES ci_kernel.opportunities(id),
  decision_id text REFERENCES ci_kernel.decisions(id),
  proposal_id text REFERENCES ci_kernel.action_proposals(id),
  action_gate_id text REFERENCES ci_kernel.action_gates(id),
  created_at timestamptz NOT NULL
);

-- Performance and deterministic replay paths.
CREATE INDEX ci_runtime_sources_person_message_idx ON ci_kernel.runtime_sources(person_id, message_id);
CREATE INDEX ci_events_person_attention_idx ON ci_kernel.events(person_id, attention_at);
CREATE INDEX ci_open_loops_person_unresolved_idx ON ci_kernel.open_loops(person_id, resolved_at, attention_at);
CREATE INDEX ci_attention_budgets_person_window_idx ON ci_kernel.attention_budgets(person_id, window_start, window_end);
CREATE INDEX ci_claims_person_freshness_idx ON ci_kernel.claims(person_id, freshness_status, last_validated_at);
CREATE INDEX ci_evidence_claim_authority_idx ON ci_kernel.evidence(claim_id, authority DESC);
CREATE INDEX ci_gates_person_state_idx ON ci_kernel.action_gates(person_id, state);

-- Same-person integrity must be enforced for cross-table IDs in the future
-- transaction function or a reviewed trigger because simple FKs do not compare
-- each referenced row's person_id. The ingestion transaction must check every
-- source/event/memory/claim/intent/loop/opportunity/decision/proposal/gate ID
-- belongs to the transaction's canonical person before commit.
--
-- CREATE FUNCTION ci_kernel.ingest_runtime_slice(...) RETURNS ... AS $$
--   -- assert server-supplied canonical person, binding, and all same-person FKs
--   -- insert source with ON CONFLICT(idempotency_key) return prior replay
--   -- append all history records and action gate transition atomically
-- $$ LANGUAGE plpgsql SECURITY DEFINER;
--
-- The function must have a dedicated internal owner, a fixed search_path, no
-- dynamic SQL, explicit stable-subject authorization, and staging security tests.

-- RLS is enabled everywhere. No INSERT/SELECT/UPDATE/DELETE policy is granted
-- to anon or public roles. A future server-only role or audited security-definer
-- function may receive narrow person-scoped access after review.
ALTER TABLE ci_kernel.persons ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.person_bindings ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.runtime_sources ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.events ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.memories ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.claims ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.evidence ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.memory_event_links ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.claim_lineage ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.goals ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.constraints ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.pending_intents ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.pending_intent_transitions ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.open_loops ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.attention_budgets ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.attention_allocations ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.opportunities ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.decisions ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.permissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.permission_scopes ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.action_proposals ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.action_gates ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.action_gate_transitions ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel.runtime_replays ENABLE ROW LEVEL SECURITY;

-- No policies are created in this proposal. Absence of a policy denies direct
-- client access. Do not add a public policy to work around existing profiles,
-- feature_flags, or audit_events security findings.

COMMIT;
