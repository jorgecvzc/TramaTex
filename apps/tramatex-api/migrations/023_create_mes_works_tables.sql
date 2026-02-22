-- Migration: 023_create_mes_works_tables.sql
-- Created at: 2026-02-21
-- Description: Creates MES work tables (mes_works, mes_work_service_groups, mes_work_tasks).

BEGIN;

CREATE TABLE IF NOT EXISTS "mes_works" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "work_number" VARCHAR(50) NOT NULL UNIQUE,
    "work_name" VARCHAR(200) NOT NULL,
    "party_id" VARCHAR(36) NOT NULL,
    "tangible_group_id" UUID NOT NULL,
    "garment_notes" TEXT,
    "status" VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    "priority" VARCHAR(20) NOT NULL DEFAULT 'NORMAL',
    "start_date" TIMESTAMP WITH TIME ZONE,
    "due_date" TIMESTAMP WITH TIME ZONE,
    "completed_date" TIMESTAMP WITH TIME ZONE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT "fk_mes_works_party" FOREIGN KEY ("party_id") REFERENCES "parties" ("id") ON DELETE RESTRICT,
    CONSTRAINT "fk_mes_works_tangible_group" FOREIGN KEY ("tangible_group_id") REFERENCES "product_groups" ("id") ON DELETE RESTRICT,
    CONSTRAINT "chk_mes_works_status" CHECK ("status" IN ('DRAFT', 'PENDING', 'IN_PROGRESS', 'ON_HOLD', 'COMPLETED', 'CANCELLED')),
    CONSTRAINT "chk_mes_works_priority" CHECK ("priority" IN ('LOW', 'NORMAL', 'HIGH', 'URGENT'))
);

CREATE INDEX IF NOT EXISTS "idx_mes_works_party_id" ON "mes_works" ("party_id");
CREATE INDEX IF NOT EXISTS "idx_mes_works_status" ON "mes_works" ("status");
CREATE INDEX IF NOT EXISTS "idx_mes_works_priority" ON "mes_works" ("priority");

CREATE TABLE IF NOT EXISTS "mes_work_service_groups" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "mes_work_id" UUID NOT NULL,
    "service_group_id" UUID NOT NULL,
    "position_id" UUID NOT NULL,
    "design_file_path" VARCHAR(500),
    "notes" TEXT,
    "sequence" INT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT "fk_mes_work_service_groups_work" FOREIGN KEY ("mes_work_id") REFERENCES "mes_works" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_mes_work_service_groups_service_group" FOREIGN KEY ("service_group_id") REFERENCES "service_groups" ("id") ON DELETE RESTRICT,
    CONSTRAINT "fk_mes_work_service_groups_position" FOREIGN KEY ("position_id") REFERENCES "positions" ("id") ON DELETE RESTRICT,
    CONSTRAINT "chk_mes_work_service_groups_sequence_positive" CHECK ("sequence" > 0)
);

CREATE INDEX IF NOT EXISTS "idx_mes_work_service_groups_work" ON "mes_work_service_groups" ("mes_work_id");
CREATE INDEX IF NOT EXISTS "idx_mes_work_service_groups_sequence" ON "mes_work_service_groups" ("mes_work_id", "sequence");

CREATE TABLE IF NOT EXISTS "mes_work_tasks" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "mes_work_service_group_id" UUID NOT NULL,
    "task_id" UUID NOT NULL,
    "sequence" INT NOT NULL,
    "status" VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    "assigned_to" UUID,
    "started_at" TIMESTAMP WITH TIME ZONE,
    "completed_at" TIMESTAMP WITH TIME ZONE,
    "notes" TEXT,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT "fk_mes_work_tasks_group" FOREIGN KEY ("mes_work_service_group_id") REFERENCES "mes_work_service_groups" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_mes_work_tasks_task" FOREIGN KEY ("task_id") REFERENCES "tasks" ("id") ON DELETE RESTRICT,
    CONSTRAINT "chk_mes_work_tasks_sequence_positive" CHECK ("sequence" > 0),
    CONSTRAINT "chk_mes_work_tasks_status" CHECK ("status" IN ('PENDING', 'IN_PROGRESS', 'COMPLETED', 'BLOCKED', 'SKIPPED'))
);

CREATE INDEX IF NOT EXISTS "idx_mes_work_tasks_group" ON "mes_work_tasks" ("mes_work_service_group_id");
CREATE INDEX IF NOT EXISTS "idx_mes_work_tasks_status" ON "mes_work_tasks" ("status");

DROP TRIGGER IF EXISTS update_mes_works_updated_at ON "mes_works";
DROP TRIGGER IF EXISTS update_mes_work_service_groups_updated_at ON "mes_work_service_groups";
DROP TRIGGER IF EXISTS update_mes_work_tasks_updated_at ON "mes_work_tasks";

CREATE TRIGGER update_mes_works_updated_at BEFORE UPDATE ON "mes_works" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_mes_work_service_groups_updated_at BEFORE UPDATE ON "mes_work_service_groups" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_mes_work_tasks_updated_at BEFORE UPDATE ON "mes_work_tasks" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

COMMIT;
