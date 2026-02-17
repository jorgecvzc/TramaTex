-- Migration: 008_add_creation_identifier_to_party_roles
-- Description: Add creation_identifier column to party_roles table
-- Created: 2026-02-12

BEGIN;

ALTER TABLE party_roles
ADD COLUMN creation_identifier VARCHAR(255) NULL;

COMMIT;