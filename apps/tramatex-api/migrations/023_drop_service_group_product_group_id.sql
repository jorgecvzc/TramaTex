-- 023: Remove obsolete product_group_id from service_groups (WorkType)
-- This column was removed from the domain during the MES refactor.
-- The product/garment association now belongs to WorkSetup, not WorkType.
DROP INDEX IF EXISTS idx_service_groups_product_group_id;
ALTER TABLE service_groups DROP CONSTRAINT IF EXISTS fk_service_groups_product_group;
ALTER TABLE service_groups DROP COLUMN IF EXISTS product_group_id;
