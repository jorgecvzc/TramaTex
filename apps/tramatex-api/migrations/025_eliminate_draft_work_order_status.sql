-- Migration 025: Eliminate DRAFT status from MES work orders
-- DRAFT is semantically identical to PENDING (no restrictions, no guard logic).
-- All existing DRAFT work orders are migrated to PENDING.
UPDATE mes_works SET status = 'PENDING' WHERE status = 'DRAFT';
