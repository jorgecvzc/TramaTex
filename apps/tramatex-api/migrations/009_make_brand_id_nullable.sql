-- Migration to make brand_id nullable in products.
-- The check keeps this migration safe on legacy databases where this table/column
-- can differ from current consolidated schema.

DO $$
BEGIN
	IF EXISTS (
		SELECT 1
		FROM information_schema.columns
		WHERE table_schema = 'public'
		  AND table_name = 'products'
		  AND column_name = 'brand_id'
	) THEN
		EXECUTE 'ALTER TABLE products ALTER COLUMN brand_id DROP NOT NULL';
	END IF;
END $$;
