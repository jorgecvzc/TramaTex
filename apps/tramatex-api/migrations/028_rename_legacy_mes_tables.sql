-- 028_rename_legacy_mes_tables.sql
-- Renames MES tables from legacy names to canonical domain names:
--   service_groups          → work_types
--   service_group_tasks     → work_type_tasks
--   mes_works               → work_orders
--   mes_work_service_groups → work_order_lines
--   mes_work_tasks          → work_order_tasks
--
-- Also renames relevant FK columns to match new table conventions.

BEGIN;

-- 1. Rename tables
ALTER TABLE service_groups          RENAME TO work_types;
ALTER TABLE service_group_tasks     RENAME TO work_type_tasks;
ALTER TABLE mes_works               RENAME TO work_orders;
ALTER TABLE mes_work_service_groups RENAME TO work_order_lines;
ALTER TABLE mes_work_tasks          RENAME TO work_order_tasks;

-- 2. Rename FK columns that used legacy names
ALTER TABLE work_type_tasks   RENAME COLUMN service_group_id          TO work_type_id;
ALTER TABLE work_order_lines  RENAME COLUMN mes_work_id               TO work_order_id;
ALTER TABLE work_order_lines  RENAME COLUMN service_group_id          TO work_type_id;
ALTER TABLE work_order_tasks  RENAME COLUMN mes_work_service_group_id TO work_order_line_id;

-- 3. Update order_work_setups FK to reference renamed table
-- (The FK constraint from migration 026 references mes_works; update it)
ALTER TABLE order_work_setups DROP CONSTRAINT IF EXISTS order_work_setups_work_order_id_fkey;
ALTER TABLE order_work_setups
    ADD CONSTRAINT order_work_setups_work_order_id_fkey
    FOREIGN KEY (work_order_id) REFERENCES work_orders(id);

COMMIT;
