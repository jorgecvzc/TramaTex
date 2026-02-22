-- Migration 027: Add tax_rate column to products table
-- Purpose: Add support for product-specific tax rates (VAT/IVA)
-- Related: Each product can have a different tax rate (e.g., 21%, 10%, 4% in Spain)
--          This affects pricing calculations and sales invoicing

ALTER TABLE products 
ADD COLUMN tax_rate NUMERIC(5,2) NOT NULL DEFAULT 21.00;

COMMENT ON COLUMN products.tax_rate IS 
'Tax rate (VAT/IVA) applicable to this product as a percentage. 
Value represents percentage (e.g., 21.0 = 21% IVA).
Common Spanish VAT rates: 21% (general), 10% (reduced), 4% (super-reduced).
Used in pricing calculations and sales invoicing to calculate final price with tax.
Formula: FinalPriceWithTax = BasePrice * (1 + (TaxRate / 100))';

-- Set default 21% IVA for existing products (general VAT rate in Spain)
UPDATE products SET tax_rate = 21.00 WHERE tax_rate IS NULL;

