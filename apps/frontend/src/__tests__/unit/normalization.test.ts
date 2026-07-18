import { describe, it, expect } from 'vitest'
// We need to import the function, but it's not exported.
// For testing purposes, we'll test it via the public getQuote method with a mock.
import salesApi from '../../services/salesApi'
import { api } from '../../services/api'
import { vi } from 'vitest'

vi.mock('../../services/api', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
  }
}))

describe('Normalization Logic via SalesApi', () => {
  it('should normalize snake_case and PascalCase to camelCase', async () => {
    const snakeResponse = {
      id: 'q-1',
      quote_number: 'PRE-001',
      party_id: 'p-1',
      line_items: [
        {
          id: 'li-1',
          product_variant_id: 'v-1',
          unit_price: { amount: 10, currency: 'EUR' }
        }
      ]
    }

    vi.mocked(api.get).mockResolvedValueOnce({ data: snakeResponse })

    const result = await salesApi.getQuote('q-1')

    expect(result.quoteNumber).toBe('PRE-001')
    expect(result.partyId).toBe('p-1')
    expect(result.lineItems).toBeDefined()
    expect(result.lineItems[0].productVariantId).toBe('v-1')
    expect(result.lineItems[0].unitPrice.amount).toBe(10)
  })

  it('should handle PascalCase and ID variations', async () => {
    const pascalResponse = {
      ID: 'q-2',
      QuoteNumber: 'PRE-002',
      LineItems: [
        {
          ID: 'li-2',
          ProductVariantID: 'v-2',
          UnitPrice: { amount: 20, currency: 'EUR' }
        }
      ]
    }

    vi.mocked(api.get).mockResolvedValueOnce({ data: pascalResponse })

    const result = await salesApi.getQuote('q-2')

    expect(result.id).toBe('q-2')
    expect(result.quoteNumber).toBe('PRE-002')
    expect(result.lineItems).toBeDefined()
    expect(result.lineItems[0].id).toBe('li-2')
    expect(result.lineItems[0].productVariantId).toBe('v-2')
  })
  
  it('should handle "items" alias for lineItems', async () => {
    const aliasResponse = {
      id: 'q-3',
      items: [
        { id: 'li-3', productName: 'Test' }
      ]
    }

    vi.mocked(api.get).mockResolvedValueOnce({ data: aliasResponse })

    const result = await salesApi.getQuote('q-3')

    expect(result.id).toBe('q-3')
    expect(result.lineItems[0].productName).toBe('Test')
  })

  it('should auto-unwrap { data: ... } wrapper', async () => {
    const wrappedResponse = {
      data: {
        id: 'q-4',
        quoteNumber: 'PRE-004',
        lineItems: []
      }
    }

    vi.mocked(api.get).mockResolvedValueOnce({ data: wrappedResponse })

    const result = await salesApi.getQuote('q-4')

    expect(result.id).toBe('q-4')
    expect(result.quoteNumber).toBe('PRE-004')
  })

  it('should handle complex production-like response with nested items', async () => {
    const complexResponse = {
      ID: 'q-5',
      QuoteNumber: 'PRE-005',
      LineItems: [
        {
          ID: 'li-5',
          ProductVariantID: 'v-5',
          product_name: 'Textile Shirt',
          UnitPrice: { amount: 15, currency: 'EUR' },
          discount_percent: 5
        }
      ],
      mes_work_refs: [
        {
          id: 'ref-1',
          work_setup_id: 'setup-1'
        }
      ]
    }

    vi.mocked(api.get).mockResolvedValueOnce({ data: complexResponse })

    const result = await salesApi.getQuote('q-5')

    expect(result.id).toBe('q-5')
    expect(result.quoteNumber).toBe('PRE-005')
    expect(result.lineItems).toHaveLength(1)
    expect(result.lineItems[0].id).toBe('li-5')
    expect(result.lineItems[0].productVariantId).toBe('v-5')
    expect(result.lineItems[0].productName).toBe('Textile Shirt')
    expect(result.lineItems[0].unitPrice.amount).toBe(15)
    expect(result.lineItems[0].discountPercent).toBe(5)
    
    expect(result.mesWorkRefs).toHaveLength(1)
    expect(result.mesWorkRefs[0].workSetupId).toBe('setup-1')
  })
})
