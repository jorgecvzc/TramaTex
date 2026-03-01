-- Migration: 008_add_party_discount
-- Description: Add default_discount_percentage to parties table for client-specific discounts

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'parties') THEN
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'parties' AND column_name = 'default_discount_percentage') THEN
            ALTER TABLE parties ADD COLUMN default_discount_percentage NUMERIC(5,2) NOT NULL DEFAULT 0;
            COMMENT ON COLUMN parties.default_discount_percentage IS 'Default discount percentage applied when this party acts as a client (0-100)';
        END IF;
    END IF;
END $$;
