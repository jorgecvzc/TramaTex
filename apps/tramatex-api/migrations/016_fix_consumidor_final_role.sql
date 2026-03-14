-- ============================================================================
-- Migration: 016_fix_consumidor_final_role.sql
-- Description: Fix CONSUMIDOR FINAL party role from CUSTOMER to CLIENT.
--              Migration 012 incorrectly seeded role as 'CUSTOMER' but the
--              domain (party_types.go) defines the valid role as 'CLIENT'.
-- Date: 2026-03-11
-- ============================================================================

BEGIN;

-- Fix the role for CONSUMIDOR FINAL party
UPDATE party_roles
SET role = 'CLIENT'
WHERE party_id = '00000000-0000-0000-0000-000000000001'
  AND role = 'CUSTOMER';

COMMIT;
