import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { productApi } from '../../services/productApi'
import { api } from '../../services/api'

describe('ProductApi Service', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    localStorage.setItem('tramatex_auth_token', 'test-token')
    localStorage.setItem('tramatex_user', JSON.stringify({ id: 'user-123' }))
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  // ============================================================================
  // PRODUCTS
  // ============================================================================

  describe('listProducts', () => {
    it('should list products successfully', async () => {
      const mockProducts = [
        {
          id: 'prod-001',
          sku: 'P001',
          name: 'Product 1',
          isActive: true,
        },
      ]

      vi.mocked(api.get).mockResolvedValueOnce({
        data: mockProducts,
        status: 200,
      })

      const result = await productApi.listProducts()

      expect(result.data).toHaveLength(1)
      expect(result.data[0].name).toBe('Product 1')
      expect(api.get).toHaveBeenCalledWith(
        expect.stringContaining('/products'),
        expect.any(Object)
      )
    })

    it('should apply filters correctly', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: [],
        status: 200,
      })

      await productApi.listProducts({
        search: 'test',
        brandId: 'brand-001',
        isActive: true,
      })

      expect(api.get).toHaveBeenCalledWith(
        expect.stringContaining('/products'),
        expect.objectContaining({
          params: expect.objectContaining({
            search: 'test',
            brandId: 'brand-001'
          })
        })
      )
    })

    it('should handle error when listing products', async () => {
      vi.mocked(api.get).mockRejectedValueOnce({
        response: {
          data: { error: 'Database error' },
          status: 500
        }
      })

      await expect(productApi.listProducts()).rejects.toThrow('Database error')
    })
  })

  describe('getProduct', () => {
    it('should get product by ID', async () => {
      const mockProduct = {
        id: 'prod-001',
        name: 'Product 1',
      }

      vi.mocked(api.get).mockResolvedValueOnce({
        data: mockProduct,
        status: 200,
      })

      const result = await productApi.getProduct('prod-001')

      expect(result.id).toBe('prod-001')
      expect(api.get).toHaveBeenCalledWith(
        expect.stringContaining('/products/prod-001')
      )
    })
  })

  describe('createProduct', () => {
    it('should create a new product', async () => {
      const mockProduct = { id: 'prod-001', sku: 'P001' }

      vi.mocked(api.post).mockResolvedValueOnce({
        data: mockProduct,
        status: 201,
      })

      const result = await productApi.createProduct({
        sku: 'P001',
        name: 'New Product',
        productType: 'SIMPLE',
      })

      expect(result.id).toBe('prod-001')
      expect(api.post).toHaveBeenCalled()
    })
  })

  describe('Authentication and headers', () => {
    it('should work with auth token', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: [],
        status: 200,
      })

      await productApi.listProducts()
      expect(api.get).toHaveBeenCalled()
    })

    it('should handle network errors gracefully', async () => {
      vi.mocked(api.get).mockRejectedValueOnce(new Error('Network error'))
      await expect(productApi.listProducts()).rejects.toThrow()
    })
  })
})
