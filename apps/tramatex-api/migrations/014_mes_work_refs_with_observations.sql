-- Migration 014: Refactor MES from UUID array to JSONB with observations
--
-- Changes document-level MES references from a simple UUID[] to a JSONB array
-- of objects with {mes_work_id, observations} to allow per-MES notes.
-- Also removes the unused per-line-item mes_work_id columns.

-- Step 1: Add new JSONB column to quotes
ALTER TABLE quotes ADD COLUMN IF NOT EXISTS mes_work_refs JSONB DEFAULT NULL;
COMMENT ON COLUMN quotes.mes_work_refs IS 'Document-level MES work references with observations: [{mes_work_id, observations}]';

-- Migrate existing data from UUID array to JSONB
UPDATE quotes
SET mes_work_refs = (
    SELECT jsonb_agg(jsonb_build_object('mes_work_id', elem::text, 'observations', ''))
    FROM unnest(mes_work_ids) AS elem
)
WHERE mes_work_ids IS NOT NULL AND array_length(mes_work_ids, 1) > 0;

-- Drop old UUID array column
ALTER TABLE quotes DROP COLUMN IF EXISTS mes_work_ids;

-- Step 2: Add new JSONB column to sales_orders
ALTER TABLE sales_orders ADD COLUMN IF NOT EXISTS mes_work_refs JSONB DEFAULT NULL;
COMMENT ON COLUMN sales_orders.mes_work_refs IS 'Document-level MES work references with observations: [{mes_work_id, observations}]';

-- Migrate existing data from UUID array to JSONB
UPDATE sales_orders
SET mes_work_refs = (
    SELECT jsonb_agg(jsonb_build_object('mes_work_id', elem::text, 'observations', ''))
    FROM unnest(mes_work_ids) AS elem
)
WHERE mes_work_ids IS NOT NULL AND array_length(mes_work_ids, 1) > 0;

-- Drop old UUID array column
ALTER TABLE sales_orders DROP COLUMN IF EXISTS mes_work_ids;

-- Step 3: Remove per-line-item MES columns (always NULL in MVP)
-- Drop FK constraints first
ALTER TABLE quote_line_items DROP CONSTRAINT IF EXISTS fk_quote_line_items_mes_work;
ALTER TABLE order_line_items DROP CONSTRAINT IF EXISTS fk_order_line_items_mes_work;

-- Drop indexes
DROP INDEX IF EXISTS idx_quote_line_items_mes_work_id;
DROP INDEX IF EXISTS idx_order_line_items_mes_work_id;

-- Drop columns
ALTER TABLE quote_line_items DROP COLUMN IF EXISTS mes_work_id;
ALTER TABLE order_line_items DROP COLUMN IF EXISTS mes_work_id;
