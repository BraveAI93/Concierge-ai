-- CENTRAL INTELLIGENCE RUNTIME v0.4.1 — DISPOSABLE STAGING ONLY
--
-- This file is for an isolated disposable local/staging PostgreSQL target only.
-- Never apply it to production Supabase without an approved change plan. No
-- application route runs this migration. The local v0.4.1 test target is
-- ci_kernel_v04_test through a sandbox-only Unix socket on port 55432.

BEGIN;

CREATE SCHEMA IF NOT EXISTS ci_kernel_v04;

-- These no-login roles are server-side database roles, never browser roles.
-- In local tests only, the migration executor receives membership to exercise
-- both boundaries. Production must map each role to separate reviewed backend
-- principals; ordinary runtime credentials must never receive the provisioner.
DO $$
BEGIN
  CREATE ROLE ci_kernel_runtime NOLOGIN;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;
DO $$
BEGIN
  CREATE ROLE ci_kernel_identity_provisioner NOLOGIN;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;
GRANT ci_kernel_runtime TO CURRENT_USER;
GRANT ci_kernel_identity_provisioner TO CURRENT_USER;

CREATE OR REPLACE FUNCTION ci_kernel_v04.current_person_id()
RETURNS text LANGUAGE sql STABLE AS $$
  SELECT NULLIF(current_setting('ci_kernel.person_id', true), '')
$$;

CREATE OR REPLACE FUNCTION ci_kernel_v04.current_stable_subject()
RETURNS text LANGUAGE sql STABLE AS $$
  SELECT NULLIF(current_setting('ci_kernel.stable_subject', true), '')
$$;

-- Canonical Person and PersonalWorld payloads remain complete v0.2 objects.
-- Person ID is a canonical globally unique identity. All operational/logical
-- identifiers below are tenant-scoped at the relational boundary.
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

-- Stable subject and internal source-profile mappings are high-trust identity
-- data. A stable subject and source profile are globally unique identities by
-- canonical contract; only the dedicated provisioner can create mappings.
CREATE TABLE ci_kernel_v04.person_binding_subjects (
  stable_subject text PRIMARY KEY,
  person_id text NOT NULL REFERENCES ci_kernel_v04.people(id) ON DELETE RESTRICT,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE ci_kernel_v04.person_profile_links (
  person_id text NOT NULL REFERENCES ci_kernel_v04.people(id) ON DELETE RESTRICT,
  source_profile_id text NOT NULL,
  is_primary boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (person_id, source_profile_id),
  UNIQUE (source_profile_id)
);
CREATE UNIQUE INDEX person_profile_links_primary_idx
  ON ci_kernel_v04.person_profile_links(person_id) WHERE is_primary;

-- Source IDs are logical source-system identifiers, not global CI identities.
-- The same source ID/message ID may validly occur in two separate worlds.
CREATE TABLE ci_kernel_v04.runtime_sources (
  person_id text NOT NULL REFERENCES ci_kernel_v04.people(id) ON DELETE RESTRICT,
  id text NOT NULL,
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
  PRIMARY KEY (person_id, id),
  UNIQUE (person_id, message_id),
  FOREIGN KEY (person_id, source_profile_id)
    REFERENCES ci_kernel_v04.person_profile_links(person_id, source_profile_id)
    ON DELETE RESTRICT
);

-- Canonical record IDs are person-owned identifiers. `record_kind` lets a
-- Person use the same raw ID in distinct canonical kinds where valid, while
-- person_id prevents one tenant's logical ID from colliding with another's.
CREATE TABLE ci_kernel_v04.records (
  record_kind text NOT NULL,
  person_id text NOT NULL REFERENCES ci_kernel_v04.people(id) ON DELETE RESTRICT,
  id text NOT NULL,
  payload jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (record_kind, person_id, id)
);
CREATE INDEX records_person_kind_idx ON ci_kernel_v04.records(person_id, record_kind, id);
CREATE INDEX records_payload_gin_idx ON ci_kernel_v04.records USING gin(payload);

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

-- Allocation IDs are derived operational identifiers and therefore tenant-safe.
CREATE TABLE ci_kernel_v04.attention_allocations (
  person_id text NOT NULL REFERENCES ci_kernel_v04.people(id) ON DELETE RESTRICT,
  allocation_id text NOT NULL,
  budget_kind text NOT NULL DEFAULT 'attention_budget' CHECK (budget_kind = 'attention_budget'),
  budget_id text NOT NULL,
  payload jsonb NOT NULL,
  evaluated_at timestamptz NOT NULL,
  PRIMARY KEY (person_id, allocation_id),
  FOREIGN KEY (budget_kind, person_id, budget_id)
    REFERENCES ci_kernel_v04.records(record_kind, person_id, id) DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX attention_allocations_person_budget_idx
  ON ci_kernel_v04.attention_allocations(person_id, budget_id, evaluated_at);

-- Replay/idempotency is scoped to the canonical person. A source key in one
-- world cannot reveal or reserve the same logical key in another world.
CREATE TABLE ci_kernel_v04.runtime_replays (
  person_id text NOT NULL REFERENCES ci_kernel_v04.people(id) ON DELETE RESTRICT,
  idempotency_key text NOT NULL,
  payload jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (person_id, idempotency_key)
);
CREATE INDEX runtime_replays_person_idx ON ci_kernel_v04.runtime_replays(person_id, idempotency_key);

-- RLS is enabled and forced on every actual kernel table. There are no anon,
-- PUBLIC, or broad authenticated policies. The runtime role only has policies
-- needed for normal person-scoped reads/inserts (and World compatibility
-- updates). The provisioner only has INSERT policies for identity creation.
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

-- Defense in depth: never leave schema/function execution available to PUBLIC.
REVOKE ALL ON SCHEMA ci_kernel_v04 FROM PUBLIC;
REVOKE ALL ON ALL TABLES IN SCHEMA ci_kernel_v04 FROM PUBLIC;
REVOKE ALL ON ALL FUNCTIONS IN SCHEMA ci_kernel_v04 FROM PUBLIC;

-- Runtime least-privilege grants. It cannot create, update, or delete identity
-- bindings; it cannot delete any canonical/runtime record.
GRANT USAGE ON SCHEMA ci_kernel_v04 TO ci_kernel_runtime;
GRANT EXECUTE ON FUNCTION ci_kernel_v04.current_person_id() TO ci_kernel_runtime;
GRANT EXECUTE ON FUNCTION ci_kernel_v04.current_stable_subject() TO ci_kernel_runtime;
GRANT SELECT ON ci_kernel_v04.people, ci_kernel_v04.worlds,
  ci_kernel_v04.person_binding_subjects, ci_kernel_v04.person_profile_links,
  ci_kernel_v04.runtime_sources, ci_kernel_v04.records,
  ci_kernel_v04.memory_event_links, ci_kernel_v04.attention_allocations,
  ci_kernel_v04.runtime_replays TO ci_kernel_runtime;
GRANT INSERT ON ci_kernel_v04.runtime_sources, ci_kernel_v04.records,
  ci_kernel_v04.memory_event_links, ci_kernel_v04.attention_allocations,
  ci_kernel_v04.runtime_replays TO ci_kernel_runtime;
GRANT UPDATE ON ci_kernel_v04.worlds TO ci_kernel_runtime;

-- Identity provisioner is the only role that can initially create Person,
-- world, stable-subject, and internal profile bindings. It receives no update
-- or delete authority, so rebinding requires a separately reviewed future flow.
GRANT USAGE ON SCHEMA ci_kernel_v04 TO ci_kernel_identity_provisioner;
GRANT INSERT ON ci_kernel_v04.people, ci_kernel_v04.worlds,
  ci_kernel_v04.person_binding_subjects, ci_kernel_v04.person_profile_links
  TO ci_kernel_identity_provisioner;

CREATE POLICY people_runtime_read ON ci_kernel_v04.people
  FOR SELECT TO ci_kernel_runtime
  USING (id = ci_kernel_v04.current_person_id());
CREATE POLICY worlds_runtime_read ON ci_kernel_v04.worlds
  FOR SELECT TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY worlds_runtime_update ON ci_kernel_v04.worlds
  FOR UPDATE TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id())
  WITH CHECK (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY bindings_runtime_read ON ci_kernel_v04.person_binding_subjects
  FOR SELECT TO ci_kernel_runtime
  USING (stable_subject = ci_kernel_v04.current_stable_subject());
CREATE POLICY profile_links_runtime_read ON ci_kernel_v04.person_profile_links
  FOR SELECT TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY sources_runtime_read ON ci_kernel_v04.runtime_sources
  FOR SELECT TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY sources_runtime_insert ON ci_kernel_v04.runtime_sources
  FOR INSERT TO ci_kernel_runtime
  WITH CHECK (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY records_runtime_read ON ci_kernel_v04.records
  FOR SELECT TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY records_runtime_insert ON ci_kernel_v04.records
  FOR INSERT TO ci_kernel_runtime
  WITH CHECK (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY links_runtime_read ON ci_kernel_v04.memory_event_links
  FOR SELECT TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY links_runtime_insert ON ci_kernel_v04.memory_event_links
  FOR INSERT TO ci_kernel_runtime
  WITH CHECK (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY allocations_runtime_read ON ci_kernel_v04.attention_allocations
  FOR SELECT TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY allocations_runtime_insert ON ci_kernel_v04.attention_allocations
  FOR INSERT TO ci_kernel_runtime
  WITH CHECK (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY replays_runtime_read ON ci_kernel_v04.runtime_replays
  FOR SELECT TO ci_kernel_runtime
  USING (person_id = ci_kernel_v04.current_person_id());
CREATE POLICY replays_runtime_insert ON ci_kernel_v04.runtime_replays
  FOR INSERT TO ci_kernel_runtime
  WITH CHECK (person_id = ci_kernel_v04.current_person_id());

CREATE POLICY people_provisioner_insert ON ci_kernel_v04.people
  FOR INSERT TO ci_kernel_identity_provisioner WITH CHECK (true);
CREATE POLICY worlds_provisioner_insert ON ci_kernel_v04.worlds
  FOR INSERT TO ci_kernel_identity_provisioner WITH CHECK (true);
CREATE POLICY bindings_provisioner_insert ON ci_kernel_v04.person_binding_subjects
  FOR INSERT TO ci_kernel_identity_provisioner WITH CHECK (true);
CREATE POLICY profile_links_provisioner_insert ON ci_kernel_v04.person_profile_links
  FOR INSERT TO ci_kernel_identity_provisioner WITH CHECK (true);

COMMIT;
