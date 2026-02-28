-- ============================================================================
-- Migration: 007_add_party_contact_fields.sql
-- Description: Add phone and email fields to person_profiles and organization_profiles
-- Date: 2026-02-27
-- Purpose: Allow all party entities to have direct contact information
-- ============================================================================

BEGIN;

-- ============================================================================
-- ADD PHONE AND EMAIL TO PERSON PROFILES
-- ============================================================================
ALTER TABLE person_profiles 
ADD COLUMN IF NOT EXISTS phone VARCHAR(30),
ADD COLUMN IF NOT EXISTS email VARCHAR(255);

COMMENT ON COLUMN person_profiles.phone IS 'Contact phone number for person';
COMMENT ON COLUMN person_profiles.email IS 'Contact email address for person';

-- ============================================================================
-- ADD PHONE AND EMAIL TO ORGANIZATION PROFILES
-- ============================================================================
ALTER TABLE organization_profiles 
ADD COLUMN IF NOT EXISTS phone VARCHAR(30),
ADD COLUMN IF NOT EXISTS email VARCHAR(255);

COMMENT ON COLUMN organization_profiles.phone IS 'Primary contact phone number for organization';
COMMENT ON COLUMN organization_profiles.email IS 'Primary contact email address for organization';

COMMIT;
