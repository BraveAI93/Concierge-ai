-- CENTRAL INTELLIGENCE RUNTIME v0.4 — DISPOSABLE STAGING ROLLBACK
--
-- NEVER RUN AGAINST PRODUCTION. Before running in staging, confirm the target
-- is disposable, all ci_kernel_v04 data is backed up/exported if needed, no
-- legal hold or retention obligation exists, and no running server depends on
-- the schema. This rollback is destructive.

BEGIN;

DROP SCHEMA IF EXISTS ci_kernel_v04 CASCADE;

DO $$
BEGIN
  REVOKE ci_kernel_runtime FROM CURRENT_USER;
EXCEPTION WHEN undefined_object THEN NULL;
END $$;
DO $$
BEGIN
  REVOKE ci_kernel_identity_provisioner FROM CURRENT_USER;
EXCEPTION WHEN undefined_object THEN NULL;
END $$;
DROP ROLE IF EXISTS ci_kernel_identity_provisioner;
DROP ROLE IF EXISTS ci_kernel_runtime;

COMMIT;
