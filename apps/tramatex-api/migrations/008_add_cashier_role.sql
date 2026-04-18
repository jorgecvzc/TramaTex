-- ============================================================================
-- Migration: 008_add_cashier_role.sql
-- Purpose: Add cashier role to IAM role constraint for TPV access
-- Date: 2026-04-17
-- ============================================================================

ALTER TABLE users
  DROP CONSTRAINT IF EXISTS chk_role;

ALTER TABLE users
  ADD CONSTRAINT chk_role CHECK (role IN ('admin', 'commercial', 'cashier', 'designer', 'workshop'));
