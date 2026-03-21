-- 031_add_description_to_work_ref_tables.sql
-- Adds a user-entered description field to Sales work reference join tables.
-- Previously, description came only from the linked work_setup (via JOIN),
-- which meant custom-reference descriptions (workSetupId = NULL) were never persisted.

ALTER TABLE quote_work_setups ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '';
ALTER TABLE order_work_setups ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '';
