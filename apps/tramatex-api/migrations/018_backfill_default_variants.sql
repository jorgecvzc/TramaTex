-- Migration 018: Backfill default variants for products without attributes
-- Products that have no variants get a default variant (empty attribute_values)
-- so they can be used in sales documents.

INSERT INTO product_variants (id, product_id, sku, status, attribute_values, is_active, created_by, created_at, updated_at)
SELECT
    gen_random_uuid(),
    p.id,
    p.sku,
    'CONFIRMED',
    '{}',
    true,
    'system-migration',
    NOW(),
    NOW()
FROM products p
WHERE p.is_active = true
  AND p.deleted_at IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM product_variants pv
      WHERE pv.product_id = p.id
        AND pv.deleted_at IS NULL
  );

COMMENT ON TABLE product_variants IS 'Product variants - concrete products with specific attributes. Products without attributes get a default variant with empty attribute_values.';
