-- Migration 008: Add SUSPENDED to work_orders status check constraint
-- The Go domain defines ProductionStatusSuspended = "SUSPENDED" but the original
-- constraint in 006_init_mes.sql did not include it, causing a DB violation when
-- cancelling a confirmed sales order (which tries to suspend its MES work orders).

ALTER TABLE work_orders
    DROP CONSTRAINT IF EXISTS chk_work_orders_status;

ALTER TABLE work_orders
    ADD CONSTRAINT chk_work_orders_status
        CHECK (status IN ('PENDING', 'IN_PROGRESS', 'ON_HOLD', 'COMPLETED', 'CANCELLED', 'SUSPENDED'));
