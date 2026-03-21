-- 029_work_order_setup_id_nullable.sql
-- Makes work_setup_id optional in work_orders.
-- WorkOrders can now be created from Sales on order confirmation without
-- requiring a backing WorkSetup configuration.

ALTER TABLE work_orders ALTER COLUMN work_setup_id DROP NOT NULL;
