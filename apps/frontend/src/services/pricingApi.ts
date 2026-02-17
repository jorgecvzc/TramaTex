/**
 * Pricing API Service
 * 
 * Service layer for Pricing module integration
 * Handles price calculations, pricing rules, and client-specific overrides
 */

import type {
  PriceCalculationResult,
  BaseSalesPriceResult,
  FinalSalePriceResult,
  SaleItem,
  PricingRule,
  CreatePricingRuleRequest,
} from '../types/pricing'

const BASE_URL = '/api/pricing'

/**
 * Get authorization header with user token
 */
function getHeaders(): Record<string, string> {
  const token = localStorage.getItem('tramatex_auth_token')
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }
  return headers
}

/**
 * Calculate price for a specific product variant, client, and quantity
 */
export async function calculatePrice(
  productVariantId: string,
  clientId: string,
  quantity: number
): Promise<PriceCalculationResult> {
  const response = await fetch(`${BASE_URL}/calculate`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({
      product_variant_id: productVariantId,
      client_id: clientId,
      quantity: quantity,
    }),
  })

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Error calculando precio')
  }

  return response.json()
}

/**
 * Calculate base sales price for a product variant (ADR-015)
 */
export async function calculateBaseSalesPrice(
  productId: string,
  variantId: string
): Promise<BaseSalesPriceResult> {
  const response = await fetch(`${BASE_URL}/base-sales-price/calculate`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({
      productId: productId,
      variantId: variantId,
    }),
  })

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Error calculando precio base de venta')
  }

  return response.json()
}

/**
 * Calculate final sale price with modifications (ADR-015)
 */
export async function calculateFinalSalePrice(
  saleItems: SaleItem[],
  clientId: string,
  saleDate: Date
): Promise<FinalSalePriceResult> {
  const response = await fetch(`${BASE_URL}/final-sale-price/calculate`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({
      saleItems: saleItems,
      clientId: clientId,
      saleDate: saleDate.toISOString(),
    }),
  })

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Error calculando precio final de venta')
  }

  return response.json()
}

/**
 * List all pricing rules
 */
export async function listPricingRules(): Promise<PricingRule[]> {
  const response = await fetch(`${BASE_URL}/rules`, {
    method: 'GET',
    headers: getHeaders(),
  })

  if (!response.ok) {
    throw new Error('Error obteniendo reglas de precio')
  }

  return response.json()
}

/**
 * Create a new pricing rule
 */
export async function createPricingRule(ruleData: CreatePricingRuleRequest): Promise<PricingRule> {
  const response = await fetch(`${BASE_URL}/rules`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(ruleData),
  })

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Error creando regla de precio')
  }

  return response.json()
}

/**
 * Get pricing history for a specific variant
 */
export async function getPricingHistory(variantId: string): Promise<PriceCalculationResult[]> {
  const response = await fetch(`${BASE_URL}/history/${variantId}`, {
    method: 'GET',
    headers: getHeaders(),
  })

  if (!response.ok) {
    throw new Error('Error obteniendo historial de precios')
  }

  return response.json()
}

/**
 * Create client-specific pricing override
 */
export async function createClientPricingOverride(
  clientId: string,
  productVariantId: string,
  fixedPrice: number,
  currency: string = 'EUR',
  effectiveFrom: Date = new Date(),
  effectiveTo: Date | null = null
): Promise<any> {
  const response = await fetch(`${BASE_URL}/client-overrides`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({
      client_id: clientId,
      product_variant_id: productVariantId,
      fixed_price: fixedPrice,
      currency: currency,
      effective_from: effectiveFrom.toISOString(),
      effective_to: effectiveTo ? effectiveTo.toISOString() : null,
    }),
  })

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Error creando precio específico de cliente')
  }

  return response.json()
}

/**
 * List base sales price rules (ADR-015)
 */
export async function listBaseSalesPriceRules(): Promise<any[]> {
  const response = await fetch(`${BASE_URL}/base-sales-rules`, {
    method: 'GET',
    headers: getHeaders(),
  })

  if (!response.ok) {
    throw new Error('Error obteniendo reglas de precio base de venta')
  }

  return response.json()
}

/**
 * Create base sales price rule (ADR-015)
 */
export async function createBaseSalesPriceRule(ruleData: any): Promise<any> {
  const response = await fetch(`${BASE_URL}/base-sales-rules`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(ruleData),
  })

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Error creando regla de precio base de venta')
  }

  return response.json()
}

/**
 * Update base sales price rule (ADR-015)
 */
export async function updateBaseSalesPriceRule(ruleId: string, ruleData: any): Promise<any> {
  const response = await fetch(`${BASE_URL}/base-sales-rules/${ruleId}`, {
    method: 'PUT',
    headers: getHeaders(),
    body: JSON.stringify(ruleData),
  })

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Error actualizando regla de precio base de venta')
  }

  return response.json()
}

/**
 * Create sale modification rule (ADR-015)
 */
export async function createSaleModificationRule(ruleData: any): Promise<any> {
  const response = await fetch(`${BASE_URL}/sale-modification-rules`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(ruleData),
  })

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Error creando regla de modificación de venta')
  }

  return response.json()
}

/**
 * Update sale modification rule (ADR-015)
 */
export async function updateSaleModificationRule(ruleId: string, ruleData: any): Promise<any> {
  const response = await fetch(`${BASE_URL}/sale-modification-rules/${ruleId}`, {
    method: 'PUT',
    headers: getHeaders(),
    body: JSON.stringify(ruleData),
  })

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Error actualizando regla de modificación de venta')
  }

  return response.json()
}

export const pricingApi = {
  calculatePrice,
  calculateBaseSalesPrice,
  calculateFinalSalePrice,
  listPricingRules,
  createPricingRule,
  getPricingHistory,
  createClientPricingOverride,
  listBaseSalesPriceRules,
  createBaseSalesPriceRule,
  updateBaseSalesPriceRule,
  createSaleModificationRule,
  updateSaleModificationRule,
}
