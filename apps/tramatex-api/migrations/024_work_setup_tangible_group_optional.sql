-- Migration 024: Make tangible_group_id optional in work_setups
-- Allows creating WorkSetup records from Sales without requiring a product group.
ALTER TABLE work_setups ALTER COLUMN tangible_group_id DROP NOT NULL;
