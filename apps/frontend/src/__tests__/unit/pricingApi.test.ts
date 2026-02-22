import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { pricingApi } from '../../services/pricingApi'
import type { PriceCalculationResult, BaseSalesPriceResult, FinalSalePriceResult } from '../../types/pricing'

// Mock fetch globally
globalThis.fetch = vi.fn()

describe('PricingApi Service', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    localStorage.setItem('tramatex_auth_token', 'test-token')
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  describe('calculatePrice', () => {
    it('should calculate price for product variant successfully', async () => {
      const mockResult: PriceCalculationResult = {
        product_variant_id: 'variant-001',
        client_id: 'client-001',
        quantity: 10,
        base_price: 100.00,
        discount_percentage: 15,
        discount_amount: 15.00,
        final_price: 85.00,
        applied_rules: [],
        calculated_at: '2026-02-18T10:00:00Z',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResult,
      })

      const result = await pricingApi.calculatePrice('variant-001', 'client-001', 10)

      expect(result).toEqual(mockResult)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/pricing/calculate'),
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            'Authorization': 'Bearer test-token',
          }),
          body: JSON.stringify({
            product_variant_id: 'variant-001',
            client_id: 'client-001',
            quantity: 10,
          }),
        })
      )
    })

    it('should handle error when calculating price', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Product variant not found' }),
      })

      await expect(pricingApi.calculatePrice('invalid', 'client-001', 10))
        .rejects.toThrow('Product variant not found')
    })

    it('should include auth token in headers', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({}),
      })

      await pricingApi.calculatePrice('var-1', 'client-1', 1)

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            'Authorization': 'Bearer test-token',
          }),
        })
      )
    })
  })

  describe('calculateBaseSalesPrice', () => {
    it('should calculate base sales price successfully', async () => {
      const mockResult: BaseSalesPriceResult = {
        productId: 'prod-001',
        variantId: 'variant-001',
        basePrice: 120.00,
        cost: 80.00,
        margin: 40.00,
        calculatedAt: '2026-02-18T10:00:00Z',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResult,
      })

      const result = await pricingApi.calculateBaseSalesPrice('prod-001', 'variant-001')

      expect(result).toEqual(mockResult)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/pricing/base-sales-price/calculate'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            productId: 'prod-001',
            variantId: 'variant-001',
          }),
        })
      )
    })

    it('should handle error when calculating base sales price', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Invalid product ID' }),
      })

      await expect(pricingApi.calculateBaseSalesPrice('invalid', 'variant-001'))
        .rejects.toThrow('Invalid product ID')
    })
  })

  describe('calculateFinalSalePrice', () => {
    it('should calculate final sale price with modifications', async () => {
      const saleItems = [
        { productVariantId: 'var-001', quantity: 5 },
        { productVariantId: 'var-002', quantity: 3 },
      ]

      const mockResult: FinalSalePriceResult = {
        subtotal: 650.00,
        discounts: [],
        total_discount: 50.00,
        final_total: 600.00,
        items: [
          { productVariantId: 'var-001', quantity: 5, unit_price: 100.00, subtotal: 500.00, discount: 30.00, total: 470.00 },
          { productVariantId: 'var-002', quantity: 3, unit_price: 50.00, subtotal: 150.00, discount: 20.00, total: 130.00 },
        ],
        client_id: 'client-001',
        sale_date: '2026-02-18',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResult,
      })

      const saleDate = new Date('2026-02-18')
      const result = await pricingApi.calculateFinalSalePrice(saleItems, 'client-001', saleDate)

      expect(result).toEqual(mockResult)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/pricing/final-sale-price/calculate'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('2026-02-18'),
        })
      )
    })

    it('should handle error when calculating final price', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Client not found' }),
      })

      await expect(pricingApi.calculateFinalSalePrice([], 'invalid', new Date()))
        .rejects.toThrow('Client not found')
    })
  })

  describe('listPricingRules', () => {
    it('should list all pricing rules successfully', async () => {
      const mockRules = [
        {
          id: 'rule-001',
          name: 'Volume Discount',
          type: 'DISCOUNT',
          value: 10,
          active: true,
        },
        {
          id: 'rule-002',
          name: 'Early Payment',
          type: 'DISCOUNT',
          value: 5,
          active: true,
        },
      ]

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockRules,
      })

      const result = await pricingApi.listPricingRules()

      expect(result).toEqual(mockRules)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/pricing/rules'),
        expect.objectContaining({
          method: 'GET',
        })
      )
    })

    it('should handle error when listing rules', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({}),
      })

      await expect(pricingApi.listPricingRules()).rejects.toThrow('Error obteniendo reglas de precio')
    })
  })

  describe('createPricingRule', () => {
    it('should create a new pricing rule successfully', async () => {
      const newRule = {
        name: 'New Rule',
        rule_type: 'VOLUME_DISCOUNT' as const,
        priority: 10,
        discount_percentage: 15,
      }

      const mockCreatedRule = {
        id: 'rule-003',
        ...newRule,
        active: true,
        created_at: '2026-02-18T10:00:00Z',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockCreatedRule,
      })

      const result = await pricingApi.createPricingRule(newRule)

      expect(result).toEqual(mockCreatedRule)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/pricing/rules'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(newRule),
        })
      )
    })

    it('should handle error when creating rule', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Invalid rule data' }),
      })

      await expect(pricingApi.createPricingRule({ name: '', rule_type: 'VOLUME_DISCOUNT', priority: 1, discount_percentage: 0 }))
        .rejects.toThrow('Invalid rule data')
    })
  })

  describe('getPricingHistory', () => {
    it('should fetch pricing history for variant', async () => {
      const mockHistory = [
        {
          product_variant_id: 'variant-001',
          client_id: 'client-001',
          quantity: 5,
          base_price: 100.00,
          discount_percentage: 10,
          discount_amount: 10.00,
          final_price: 90.00,
          applied_rules: [],
          calculated_at: '2026-02-15T10:00:00Z',
        },
        {
          product_variant_id: 'variant-001',
          client_id: 'client-001',
          quantity: 10,
          base_price: 100.00,
          discount_percentage: 15,
          discount_amount: 15.00,
          final_price: 85.00,
          applied_rules: [],
          calculated_at: '2026-02-16T10:00:00Z',
        },
      ]

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockHistory,
      })

      const result = await pricingApi.getPricingHistory('variant-001')

      expect(result).toEqual(mockHistory)
      expect(result).toHaveLength(2)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/pricing/history/variant-001'),
        expect.objectContaining({
          method: 'GET',
        })
      )
    })

    it('should handle error when fetching history', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({}),
      })

      await expect(pricingApi.getPricingHistory('invalid'))
        .rejects.toThrow('Error obteniendo historial de precios')
    })
  })

  describe('createClientPricingOverride', () => {
    it('should create client-specific pricing override', async () => {
      const mockOverride = {
        id: 'override-001',
        client_id: 'client-001',
        product_variant_id: 'variant-001',
        fixed_price: 95.00,
        currency: 'EUR',
        effective_from: '2026-02-18T00:00:00Z',
        effective_to: null,
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockOverride,
      })

      const effectiveFrom = new Date('2026-02-18')
      const result = await pricingApi.createClientPricingOverride(
        'client-001',
        'variant-001',
        95.00,
        'EUR',
        effectiveFrom
      )

      expect(result).toEqual(mockOverride)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/pricing/client-overrides'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('client-001'),
        })
      )
    })

    it('should create override with end date', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({}),
      })

      const effectiveFrom = new Date('2026-02-18')
      const effectiveTo = new Date('2026-03-18')
      
      await pricingApi.createClientPricingOverride(
        'client-001',
        'variant-001',
        95.00,
        'EUR',
        effectiveFrom,
        effectiveTo
      )

      const callBody = JSON.parse((globalThis.fetch as any).mock.calls[0][1].body)
      expect(callBody.effective_to).not.toBeNull()
    })

    it('should handle error when creating override', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Duplicate override' }),
      })

      await expect(pricingApi.createClientPricingOverride('client-001', 'var-001', 95.00))
        .rejects.toThrow('Duplicate override')
    })
  })

  describe('listBaseSalesPriceRules', () => {
    it('should list base sales price rules', async () => {
      const mockRules = [
        { id: 'rule-001', name: 'Standard Markup', markup: 1.5 },
        { id: 'rule-002', name: 'Premium Markup', markup: 2.0 },
      ]

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockRules,
      })

      const result = await pricingApi.listBaseSalesPriceRules()

      expect(result).toEqual(mockRules)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/pricing/base-sales-rules'),
        expect.objectContaining({
          method: 'GET',
        })
      )
    })

    it('should handle error when listing base sales rules', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({}),
      })

      await expect(pricingApi.listBaseSalesPriceRules())
        .rejects.toThrow('Error obteniendo reglas de precio base de venta')
    })
  })

  describe('createBaseSalesPriceRule', () => {
    it('should create base sales price rule', async () => {
      const ruleData = { name: 'New Markup', markup: 1.75 }
      const mockCreatedRule = { id: 'rule-003', ...ruleData }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockCreatedRule,
      })

      const result = await pricingApi.createBaseSalesPriceRule(ruleData)

      expect(result).toEqual(mockCreatedRule)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/pricing/base-sales-rules'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(ruleData),
        })
      )
    })

    it('should handle error when creating base sales rule', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Invalid markup value' }),
      })

      await expect(pricingApi.createBaseSalesPriceRule({ markup: -1 }))
        .rejects.toThrow('Invalid markup value')
    })
  })

  describe('updateBaseSalesPriceRule', () => {
    it('should update base sales price rule', async () => {
      const ruleData = { name: 'Updated Markup', markup: 1.8 }
      const mockUpdatedRule = { id: 'rule-001', ...ruleData }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockUpdatedRule,
      })

      const result = await pricingApi.updateBaseSalesPriceRule('rule-001', ruleData)

      expect(result).toEqual(mockUpdatedRule)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/pricing/base-sales-rules/rule-001'),
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(ruleData),
        })
      )
    })

    it('should handle error when updating base sales rule', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Rule not found' }),
      })

      await expect(pricingApi.updateBaseSalesPriceRule('invalid', {}))
        .rejects.toThrow('Rule not found')
    })
  })

  describe('createSaleModificationRule', () => {
    it('should create sale modification rule', async () => {
      const ruleData = { type: 'DISCOUNT', value: 10, condition: 'quantity > 100' }
      const mockCreatedRule = { id: 'mod-rule-001', ...ruleData }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockCreatedRule,
      })

      const result = await pricingApi.createSaleModificationRule(ruleData)

      expect(result).toEqual(mockCreatedRule)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/pricing/sale-modification-rules'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(ruleData),
        })
      )
    })

    it('should handle error when creating modification rule', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Invalid condition syntax' }),
      })

      await expect(pricingApi.createSaleModificationRule({ condition: 'invalid' }))
        .rejects.toThrow('Invalid condition syntax')
    })
  })

  describe('updateSaleModificationRule', () => {
    it('should update sale modification rule', async () => {
      const ruleData = { type: 'DISCOUNT', value: 15 }
      const mockUpdatedRule = { id: 'mod-rule-001', ...ruleData }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockUpdatedRule,
      })

      const result = await pricingApi.updateSaleModificationRule('mod-rule-001', ruleData)

      expect(result).toEqual(mockUpdatedRule)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/pricing/sale-modification-rules/mod-rule-001'),
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(ruleData),
        })
      )
    })

    it('should handle error when updating modification rule', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Modification rule not found' }),
      })

      await expect(pricingApi.updateSaleModificationRule('invalid', {}))
        .rejects.toThrow('Modification rule not found')
    })
  })

  describe('Authentication and headers', () => {
    it('should work without auth token', async () => {
      localStorage.removeItem('tramatex_auth_token')

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ([]),
      })

      await pricingApi.listPricingRules()

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.not.objectContaining({
            'Authorization': expect.anything(),
          }),
        })
      )
    })

    it('should include Content-Type header in POST requests', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({}),
      })

      await pricingApi.createPricingRule({ name: 'Test Rule', rule_type: 'VOLUME_DISCOUNT', priority: 10, discount_percentage: 10 })

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
        })
      )
    })
  })
})
