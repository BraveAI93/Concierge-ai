-- CENTRAL INTELLIGENCE RUNTIME v0.4 — DISPOSABLE STAGING UP MIGRATION
--
-- This file is for an isolated disposable staging/local PostgreSQL target only.
-- Never apply it to production Supabase without an approved change plan.
-- No application route runs this migration. The local v0.4 test target is
-- ci_kernel_v04_test through a sandbox-only Unix socket on port 55432.

BEGIN;

CREATE SCHEMA IF NOT EXISTS ci_kernel_v04;

-- The no-login runtime role is a server-side database role, not a browser role.
-- The current migration executor is granted membership only for disposable
-- staging tests. Production must grant this role solely to a dedicated backend
-- principal after a separate security review.
DO $$
BEGIN
  CREATE ROLE ci_kernel_runtime NOLOGIN;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;
GRANT ci_kernel_runtime TO CURRENT_USER;

CREATE OR REPLACE FUNCTION ci_kernel_v04.current_person_id()
RETURNS text LANGUAGE sql STABLE AS $$
  SELECT NULLIF(current_setting('ci_kernel.person_id', true), '')
$$;

CREATE OR REPLACE FUNCTION ci_kernel_v04.current_stable_subject()
RETURNS text LANGUAGE sql STABLE AS $$
  SELECT NULLIF(current_setting('ci_kernel.stable_subject', true), '')
$$;

-- Canonical person and personal-world records retain JSON payloads for the
-- complete v0.2 canonical object while relational IDs enforce tenancy.
CREATE TABLE ci_kernel_v04.people (
  id text PRIMARY KEY,
  payload jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE ci_kernel_v04.worlds (
  person_id text PRIMARY KEY REFERENCES ci_kernel_v04.people(id) ON DELETE RESTRICT,
  payload jsonb NOT NULL,
  updated_at timestamptz NOT NULL DEFAULT now()
);

-- Stable subject and source-profile mappings are server-resolved. The source
-- profile is an internal legacy profile ID, not a caller-supplied public slug.
CREATE TABLE ci_kernel_v04.person_binding_subjects (
  stable_subject text PRIMARY KEY,
  person_id text NOT NULL REFERENCES ci_kernel_v04.people(id) ON DELETE RESTRICT,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE ci_kernel_v04.person_profile_links (
  person_id text NOT NULL REFERENCES ci_kernel_v04.people(id) ON DELETE RESTRICT,
  source_profile_id text NOT NULL UNIQUE,
  is_primary boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (person_id, source_profile_id)
);
CREATE UNIQUE INDEX person_profile_links_primary_idx
  ON ci_kernel_v04.person_profile_links(person_id) WHERE is_primary;

-- Source lineage is separate from generic canonical record persistence so the
-- idempotency key and original message/conversation timestamps are explicit.
CREATE TABLE ci_kernel_v04.runtime_sources (
  id text PRIMARY KEY,
  person_id text NOT NULL REFERENCES ci_kernel_v04.people(id) ON DELETE RESTRICT,
  source_profile_id text NOT NULL,
  conversation_id text NOT NULL,
  session_id text,
  message_id text NOT NULL,
  message_role text NOT NULL,
  content text NOT NULL,
  conversation_at timestamptz,
  message_at timestamptz NOT NULL,
  stored_at timestamptz NOT NULL,
  payload jsonb NOT NULL,
  UNIQUE (person_id, message_id),
  FOREIGN KEY (person_id, source_profile_id)
    REFERENCES ci_kernel_v04.person_profile_links(person_id, source_profile_id)
    ON DELETE RESTRICT
);

-- All v0.2/v0.3 canonical records use an append-preserving payload table. The
-- record_kind column is indexed for durable querying and direct same-person
-- parent checks in the repository transaction. payload retains temporal state,
-- freshness/supersession, provenance, claims, attention data, outcomes, etc.
CREATE TABLE ci_kernel_v04.records (
  record_kind text NOT NULL,
  id text NOT NULL,
  person_id text NOT NULL REFERENCES ci_kernel_v04.people(id) ON DELETE RESTRICT,
  payload jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (record_kind, id),
  UNIQUE (record_kind, person_id, id)
);
CREATE INDEX records_person_kind_idx ON ci_kernel_v04.records(person_id, record_kind, id);
CREATE INDEX records_payload_gin_idx ON ci_kernel_v04.records USING gin(payload);

-- A single memory/event link is immutable for its person-scoped pair.
CREATE TABLE ci_kernel_v04.memory_event_links (
  person_id text NOT NULL REFERENCES ci_kernel_v04.people(id) ON DELETE RESTRICT,
  memory_kind text NOT NULL DEFAULT 'memory' CHECK (memory_kind = 'memory'),
  memory_id text NOT NULL,
  event_kind text NOT NULL DEFAULT 'event' CHECK (event_kind = 'event'),
  event_id text NOT NULL,
  payload jsonb NOT NULL,
  linked_at timestamptz NOT NULL,
  PRIMARY KEY (person_id, memory_id, event_id),
  FOREIGN KEY (memory_kind, person_id, memory_id)
    REFERENCES ci_kernel_v04.records(record_kind, person_id, id) DEFERRABLE INITIALLY IMMEDIATE,
  FOREIGN KEY (event_kind, person_id, event_id)
    REFERENCES ci_kernel_v04.records(record_kind, person_id, id) DEFERRABLE INITIALLY IMMEDIATE
);

-- Attention allocation is preserved with a stable allocation key and full
-- v0.2 allocation payload. The budget remains a canonical record.
CREATE TABLE ci_kernel_v04.attention_allocations (
  allocation_id text PRIMARY KEY,
  person_id text NOT NULL REFERENCES ci_kernel_v04.people(id) ON DELETE RESTRICT,
  budget_kind text NOT NULL DEFAULT 'attention_budget' CHECK (budget_kind = 'attention_budget'),
  budget_id text NOT NULL,
  payload jsonb NOT NULL,
  evaluated_at timestamptz NOT NULL,
  FOREIGN KEY (budget_kind, person_id, budget_id)
    REFERENCES ci_kernel_v04.records(record_kind, person_id, id) DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX attention_allocations_person_budget_idx
  ON ci_kernel_v04.attention_allocations(person_id, budget_id, evaluated_at);

-- The replay is the durable idempotency serialization point. Its result holds
-- the complete Event/Evidence/Memory/OpenLoop/Decision/Gate identifier lineage.
CREATE TABLE ci_kernel_v04.runtime_replays (
  idempotency_key text PRIMARY KEY,
  person_id text NOT NULL REFERENCES ci_kernel_v04.people(id) ON DELETE RESTRICT,
  payload jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX runtime_replays_person_idx ON ci_kernel_v04.runtime_replays(person_id, idempotency_key);

-- RLS is enabled and forced on every kernel table. There are no anon, public,
-- or broad authenticated policies. The only policy principal is the internal
-- ci_kernel_runtime role, and each policy binds access to SET LOCAL runtime
-- session context. The repository sets person_id/subject within each SQL tx.
ALTER TABLE ci_kernel_v04.people ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.worlds ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.person_binding_subjects ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.person_profile_links ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.runtime_sources ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.records ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.memory_event_links ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.attention_allocations ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.runtime_replays ENABLE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.people FORCE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.worlds FORCE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.person_binding_subjects FORCE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.person_profile_links FORCE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.runtime_sources FORCE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.records FORCE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.memory_event_links FORCE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.attention_allocations FORCE ROW LEVEL SECURITY;
ALTER TABLE ci_kernel_v04.runtime_replays FORCE ROW LEVEL SECURITY;

GRANT USAGE ON SCHEMA ci_kernel_v04 TO ci_kernel_runtime;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA ci_kernel_v04 TO ci_kernel_runtime;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA ci_kernel_v04 TO ci_kernel_runtime;

CREATE POLICY people_runtime_person ON ci_kernel_v04.people
  FOR ALL TO ci_kernel_runtime
  USING (id = ci_kernel_v04.current_person_id())
  WITH CHECK (id = ci_kernel_v04.current_person_id());
CREATE POLICY worlds_runtime_person ON ci_kernel_v04.worlds
  FOR ALL TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id())
  WITH CHECK (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY bindings_runtime_subject ON ci_kernel_v04.person_binding_subjects
  FOR ALL TO ci_kernel_runtime
  USING (stable_subject = ci_kernel_v04.current_stable_subject())
  WITH CHECK (stable_subject = ci_kernel_v04.current_stable_subject());
CREATE POLICY profile_links_runtime_person ON ci_kernel_v04.person_profile_links
  FOR ALL TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id())
  WITH CHECK (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY sources_runtime_person ON ci_kernel_v04.runtime_sources
  FOR ALL TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id())
  WITH CHECK (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY records_runtime_person ON ci_kernel_v04.records
  FOR ALL TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id())
  WITH CHECK (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY links_runtime_person ON ci_kernel_v04.memory_event_links
  FOR ALL TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id())
  WITH CHECK (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY allocations_runtime_person ON ci_kernel_v04.attention_allocations
  FOR ALL TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id())
  WITH CHECK (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY replays_runtime_person ON ci_kernel_v04.runtime_replays
  FOR ALL TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id())
  WITH CHECK (person_id = ci_kernel_v04.current_person_id());

-- No GRANT is made to PUBLIC. Do not add an anon/public shortcut. Existing
-- production profiles, feature_flags, and audit_events security issues are not
-- a dependency of this isolated schema.
REVOKE ALL ON SCHEMA ci_kernel_v04 FROM PUBLIC;
REVOKE ALL ON ALL TABLES IN SCHEMA ci_kernel_v04 FROM PUBLIC;

COMMIT;
