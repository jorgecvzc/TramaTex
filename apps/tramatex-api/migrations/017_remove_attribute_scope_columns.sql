-- Migration: Remove scope columns from attributes table
-- Date: 2026-02-12
-- Reason: Simplify attribute system for MVP - defer scope logic to post-MVP

BEGIN;

-- Drop scope columns from attributes table
ALTER TABLE attributes 
DROP COLUMN IF EXISTS scope_brand_id,
DROP COLUMN IF EXISTS scope_group_id;

-- Note: No indexes to drop as these columns didn't have explicit indexes

COMMIT;
