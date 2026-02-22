-- Migration: 022_create_mes_master_data_tables.sql
-- Created at: 2026-02-20
-- Description: Creates MES foundation master data tables (tasks, positions, service_groups, service_group_tasks).

BEGIN;

-- Tasks
CREATE TABLE IF NOT EXISTS "tasks" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(100) NOT NULL,
    "description" TEXT,
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS "idx_tasks_is_active" ON "tasks" ("is_active");

-- Positions
CREATE TABLE IF NOT EXISTS "positions" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(100) NOT NULL,
    "code" VARCHAR(50) NOT NULL,
    "description" TEXT,
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_positions_code" ON "positions" ("code");
CREATE INDEX IF NOT EXISTS "idx_positions_is_active" ON "positions" ("is_active");

-- Service Groups
CREATE TABLE IF NOT EXISTS "service_groups" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(100) NOT NULL,
    "description" TEXT,
    "product_group_id" UUID,
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT "fk_service_groups_product_group" FOREIGN KEY ("product_group_id") REFERENCES "product_groups" ("id") ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS "idx_service_groups_is_active" ON "service_groups" ("is_active");
CREATE INDEX IF NOT EXISTS "idx_service_groups_product_group_id" ON "service_groups" ("product_group_id");

-- Service Group Tasks (many-to-many with order)
CREATE TABLE IF NOT EXISTS "service_group_tasks" (
    "service_group_id" UUID NOT NULL,
    "task_id" UUID NOT NULL,
    "sequence" INT NOT NULL,
    PRIMARY KEY ("service_group_id", "task_id"),
    CONSTRAINT "fk_service_group_tasks_service_group" FOREIGN KEY ("service_group_id") REFERENCES "service_groups" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_service_group_tasks_task" FOREIGN KEY ("task_id") REFERENCES "tasks" ("id") ON DELETE CASCADE,
    CONSTRAINT "chk_service_group_tasks_sequence_positive" CHECK ("sequence" > 0)
);
CREATE INDEX IF NOT EXISTS "idx_service_group_tasks_service_group_id" ON "service_group_tasks" ("service_group_id");
CREATE INDEX IF NOT EXISTS "idx_service_group_tasks_sequence" ON "service_group_tasks" ("service_group_id", "sequence");

-- Update triggers
DROP TRIGGER IF EXISTS update_tasks_updated_at ON "tasks";
DROP TRIGGER IF EXISTS update_positions_updated_at ON "positions";
DROP TRIGGER IF EXISTS update_service_groups_updated_at ON "service_groups";

CREATE TRIGGER update_tasks_updated_at BEFORE UPDATE ON "tasks" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_positions_updated_at BEFORE UPDATE ON "positions" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_service_groups_updated_at BEFORE UPDATE ON "service_groups" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

COMMIT;
