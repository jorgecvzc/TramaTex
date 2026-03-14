/**
 * Pricing Module Type Definitions
 * Defines TypeScript interfaces for price calculations, pricing rules, and client pricing
 */

// ============================================================================
// PRICE CALCULATION ENTITIES
// ============================================================================

export interface PriceCalculationRequest {
  product_variant_id: string
  client_id: string
  quantity: number
}

export interface BaseSalesPriceRequest {
  productId: string
  variantId: string
}

export interface FinalSalePriceRequest {
  saleItems: SaleItem[]
  clientId: string
  saleDate: string  // ISO 8601 date string
}

export interface SaleItem {
  productVariantId: string
  quantity: number
}

export interface PriceCalculationResult {
  product_variant_id: string
  client_id: string
  quantity: number
  base_price: number
  discount_percentage: number
  discount_amount: number
  final_price: number
  applied_rules: AppliedPricingRule[]
  calculated_at: string
}

export interface MoneyResult {
  amount: number
  currency: string
}

export interface BaseSalesPriceResult {
  variantId: string
  baseCost: MoneyResult
  baseSalesPrice: MoneyResult
  taxRate: number
}

export interface CalculatedSaleItem {
  productVariantId: string
  quantity: number
  baseCost: MoneyResult
  baseSalesPrice: MoneyResult
  finalPrice: MoneyResult
  taxRate: number
  finalPriceWithTax: MoneyResult
}

export interface FinalSalePriceResult {
  calculatedItems: CalculatedSaleItem[]
  saleTotal: MoneyResult
  saleTotalWithTax: MoneyResult
}

// ============================================================================
// PRICING RULE ENTITIES
// ============================================================================

export interface PricingRule {
  id: string
  name: string
  description: string | null
  rule_type: PricingRuleType
  priority: number
  is_active: boolean
  start_date: string | null
  end_date: string | null
  conditions: PricingRuleCondition[]
  discount_percentage: number
  created_at: string
  updated_at: string
}

export type PricingRuleType = 'VOLUME_DISCOUNT' | 'CLIENT_SPECIFIC' | 'SEASONAL' | 'PRODUCT_CATEGORY'

export interface PricingRuleCondition {
  field: string
  operator: ConditionOperator
  value: string | number
}

export type ConditionOperator = 'EQUALS' | 'GREATER_THAN' | 'LESS_THAN' | 'IN_RANGE' | 'CONTAINS'

export interface AppliedPricingRule {
  rule_id: string
  rule_name: string
  discount_percentage: number
  priority: number
}

// ============================================================================
// CLIENT PRICING ENTITIES (OVERRIDES)
// ============================================================================

export interface ClientPricing {
  id: string
  client_id: string
  product_variant_id: string
  custom_price: number
  discount_percentage: number
  is_active: boolean
  valid_from: string | null
  valid_to: string | null
  created_at: string
  updated_at: string
}

// ============================================================================
// REQUEST DTOs
// ============================================================================

export interface CreatePricingRuleRequest {
  name: string
  description?: string
  rule_type: PricingRuleType
  priority: number
  discount_percentage: number
  start_date?: string
  end_date?: string
  conditions?: PricingRuleCondition[]
}

export interface UpdatePricingRuleRequest {
  name?: string
  description?: string
  priority?: number
  discount_percentage?: number
  is_active?: boolean
  start_date?: string | null
  end_date?: string | null
  conditions?: PricingRuleCondition[]
}

export interface CreateClientPricingRequest {
  client_id: string
  product_variant_id: string
  custom_price?: number
  discount_percentage?: number
  valid_from?: string
  valid_to?: string
}

export interface UpdateClientPricingRequest {
  custom_price?: number
  discount_percentage?: number
  is_active?: boolean
  valid_from?: string | null
  valid_to?: string | null
}

// ============================================================================
// FILTERS & PAGINATION
// ============================================================================

export interface ListPricingRulesFilters {
  ruleType?: PricingRuleType
  isActive?: boolean
  pageNumber?: number
  pageSize?: number
}

export interface ListClientPricingFilters {
  clientId?: string
  productVariantId?: string
  isActive?: boolean
  pageNumber?: number
  pageSize?: number
}

// ============================================================================
// API ERROR HANDLING
// ============================================================================

export interface PricingError {
  error: string
  message?: string
  cause?: Error
}
