/**
 * Pricing API Service
 * 
 * Service layer for Pricing module integration
 * Handles price calculations, pricing rules, and client-specific overrides
 */

import { fetchWithAuth } from './apiBase'

const BASE_URL = '/api/pricing'

/**
 * Calculate price for a specific product variant, client, and quantity
 * @param {string} productVariantId - UUID of the product variant
 * @param {string} clientId - UUID of the client
 * @param {number} quantity - Quantity to calculate price for
 * @returns {Promise<Object>} Price calculation result with breakdown
 */
export async function calculatePrice(productVariantId, clientId, quantity) {
  const response = await fetchWithAuth(`${BASE_URL}/calculate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
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
 * @param {string} productId - UUID of the product
 * @param {string} variantId - UUID of the variant
 * @returns {Promise<Object>} Base sales price calculation
 */
export async function calculateBaseSalesPrice(productId, variantId) {
  const response = await fetchWithAuth(`${BASE_URL}/base-sales-price/calculate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
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
 * @param {Array<Object>} saleItems - Array of {productVariantId, quantity}
 * @param {string} clientId - UUID of the client
 * @param {Date} saleDate - Sale date for calculations
 * @returns {Promise<Object>} Final sale price with all modifications applied
 */
export async function calculateFinalSalePrice(saleItems, clientId, saleDate) {
  const response = await fetchWithAuth(`${BASE_URL}/final-sale-price/calculate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
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
 * @returns {Promise<Array>} Array of pricing rules
 */
export async function listPricingRules() {
  const response = await fetchWithAuth(`${BASE_URL}/rules`, {
    method: 'GET',
  })

  if (!response.ok) {
    throw new Error('Error obteniendo reglas de precio')
  }

  return response.json()
}

/**
 * Create a new pricing rule
 * @param {Object} ruleData - Pricing rule data
 * @returns {Promise<Object>} Created pricing rule
 */
export async function createPricingRule(ruleData) {
  const response = await fetchWithAuth(`${BASE_URL}/rules`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
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
 * @param {string} variantId - UUID of the product variant
 * @returns {Promise<Array>} Array of price calculations
 */
export async function getPricingHistory(variantId) {
  const response = await fetchWithAuth(`${BASE_URL}/history/${variantId}`, {
    method: 'GET',
  })

  if (!response.ok) {
    throw new Error('Error obteniendo historial de precios')
  }

  return response.json()
}

/**
 * Create client-specific pricing override
 * @param {string} clientId - UUID of the client
 * @param {string} productVariantId - UUID of the product variant
 * @param {number} fixedPrice - Fixed price for this client
 * @param {string} currency - Currency code (e.g., 'EUR')
 * @param {Date} effectiveFrom - Start date for this override
 * @param {Date|null} effectiveTo - Optional end date
 * @returns {Promise<Object>} Created client pricing override
 */
export async function createClientPricingOverride(
  clientId,
  productVariantId,
  fixedPrice,
  currency = 'EUR',
  effectiveFrom = new Date(),
  effectiveTo = null
) {
  const response = await fetchWithAuth(`${BASE_URL}/client-overrides`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
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
 * @returns {Promise<Array>} Array of base sales price rules
 */
export async function listBaseSalesPriceRules() {
  const response = await fetchWithAuth(`${BASE_URL}/base-sales-rules`, {
    method: 'GET',
  })

  if (!response.ok) {
    throw new Error('Error obteniendo reglas de precio base de venta')
  }

  return response.json()
}

/**
 * Create base sales price rule (ADR-015)
 * @param {Object} ruleData - Base sales price rule data
 * @returns {Promise<Object>} Created rule
 */
export async function createBaseSalesPriceRule(ruleData) {
  const response = await fetchWithAuth(`${BASE_URL}/base-sales-rules`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
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
 * @param {string} ruleId - UUID of the rule
 * @param {Object} ruleData - Updated rule data
 * @returns {Promise<Object>} Updated rule
 */
export async function updateBaseSalesPriceRule(ruleId, ruleData) {
  const response = await fetchWithAuth(`${BASE_URL}/base-sales-rules/${ruleId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
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
 * @param {Object} ruleData - Sale modification rule data
 * @returns {Promise<Object>} Created rule
 */
export async function createSaleModificationRule(ruleData) {
  const response = await fetchWithAuth(`${BASE_URL}/sale-modification-rules`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
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
 * @param {string} ruleId - UUID of the rule
 * @param {Object} ruleData - Updated rule data
 * @returns {Promise<Object>} Updated rule
 */
export async function updateSaleModificationRule(ruleId, ruleData) {
  const response = await fetchWithAuth(`${BASE_URL}/sale-modification-rules/${ruleId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
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
