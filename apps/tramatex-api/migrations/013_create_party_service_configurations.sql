-- Migration: 013_create_party_service_configurations.sql
-- Created at: 2026-02-07
-- Description: Creates party_service_configurations table for Product module.

CREATE TABLE IF NOT EXISTS "party_service_configurations" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "party_id" UUID NOT NULL,
    "service_id" VARCHAR(255) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "configuration_details" JSONB,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "idx_party_service_configurations_party_id" ON "party_service_configurations" ("party_id");

DROP TRIGGER IF EXISTS update_party_service_configurations_updated_at ON "party_service_configurations";
CREATE TRIGGER update_party_service_configurations_updated_at
BEFORE UPDATE ON "party_service_configurations"
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
