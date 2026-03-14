-- Migration 019: Add invoice traceability to delivery note line items
-- Adds invoice_line_item_id FK to delivery_note_line_items for direct DN↔Invoice linking.
-- MVP: 1:1 (one DN line → one invoice line). Post-MVP: N:1 (multiple DN lines → one invoice line).

ALTER TABLE delivery_note_line_items
    ADD COLUMN invoice_line_item_id UUID REFERENCES invoice_line_items(id);

CREATE INDEX idx_dn_line_items_invoice_line_item_id
    ON delivery_note_line_items(invoice_line_item_id)
    WHERE invoice_line_item_id IS NOT NULL;
