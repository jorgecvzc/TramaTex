-- Migration 009: Move MES from per-line-item to document-level
--
-- MVP Decision: MES works are associated at the document level (quote/order),
-- not per individual line item. Multiple MES works can be linked to a single
-- document as an informational guide. Users identify which MES work corresponds
-- to which products manually.
--
-- The existing mes_work_id columns on line items are preserved for Post-MVP
-- when a more functional product-MES binding will be implemented.
-- See: docs/post-mvp/post-mvp-roadmap.md (Section 1)

-- Add document-level MES work IDs array to quotes
ALTER TABLE quotes ADD COLUMN IF NOT EXISTS mes_work_ids UUID[] DEFAULT NULL;
COMMENT ON COLUMN quotes.mes_work_ids IS 'Document-level MES work references (MVP: informational guide, multiple allowed)';

-- Add document-level MES work IDs array to sales_orders
ALTER TABLE sales_orders ADD COLUMN IF NOT EXISTS mes_work_ids UUID[] DEFAULT NULL;
COMMENT ON COLUMN sales_orders.mes_work_ids IS 'Document-level MES work references (MVP: informational guide, multiple allowed)';

-- Note: Existing mes_work_id on quote_line_items and order_line_items
-- columns are intentionally kept for backward compatibility and Post-MVP use.
