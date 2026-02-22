-- Migration 026: Add default_markup_percentage column to brands table
-- Purpose: Add support for brand-level markup percentage used in pricing calculations
-- Related: Product pricing enhancement - brands can now define a default markup percentage
--          that is applied when calculating base sales prices from variant base costs

ALTER TABLE brands 
ADD COLUMN default_markup_percentage NUMERIC(5,2) NOT NULL DEFAULT 0.00;

COMMENT ON COLUMN brands.default_markup_percentage IS 
'Default markup percentage applied to product variants of this brand in pricing calculations. 
Value represents percentage (e.g., 30.0 = 30%). Must be >= 0. 
Used by BaseSalesPriceRule: BaseSalesPrice = VariantBaseCost * (1 + (DefaultMarkupPercentage / 100))';

-- Update existing brands to have 0% markup (neutral for existing prices)
-- Operations team should update markup percentages as needed
UPDATE brands SET default_markup_percentage = 0.00 WHERE default_markup_percentage IS NULL;
