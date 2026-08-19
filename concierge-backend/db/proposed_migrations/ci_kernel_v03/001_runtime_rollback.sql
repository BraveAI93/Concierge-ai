-- CENTRAL INTELLIGENCE RUNTIME v0.3 — REVIEW-ONLY ROLLBACK
--
-- DO NOT EXECUTE. This rollback is valid only before data activation or after
-- approved export/backup, dependency, retention, and change-management review.
-- It must never be run against production merely to undo an application deploy.

BEGIN;

-- Verify no external dependency, legal hold, retention requirement, or active
-- runtime deployment references ci_kernel before any destructive operation.
DROP SCHEMA IF EXISTS ci_kernel CASCADE;

COMMIT;
