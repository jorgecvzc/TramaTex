-- 030_work_setup_id_nullable_in_join_tables.sql
-- Makes work_setup_id optional in sales join tables.
-- Allows order/quote work references to exist without a backing MES WorkSetup
-- (e.g. "personalizada" references entered by the user without selecting a setup).
-- The FK constraint is preserved; only NOT NULL is dropped so NULL is a valid value.

BEGIN;
ALTER TABLE quote_work_setups ALTER COLUMN work_setup_id DROP NOT NULL;
ALTER TABLE order_work_setups ALTER COLUMN work_setup_id DROP NOT NULL;
COMMIT;
