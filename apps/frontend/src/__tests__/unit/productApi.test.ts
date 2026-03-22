import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { productApi } from '../../services/productApi'

// Mock fetch globally
globalThis.fetch = vi.fn()

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
          longName: 'Full Product Name 1',
          description: 'Description 1',
          productType: 'CONFIGURABLE',
          brandId: 'brand-001',
          groupIds: ['group-001'],
          directAttributeIds: ['attr-001'],
          isActive: true,
          variantsCount: 5,
        },
      ]

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockProducts,
      })

      const result = await productApi.listProducts()

      expect(result.data).toHaveLength(1)
      expect(result.data[0].name).toBe('Product 1')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/products'),
        expect.objectContaining({ method: 'GET' })
      )
    })

    it('should apply filters correctly', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })

      await productApi.listProducts({
        search: 'test',
        brandId: 'brand-001',
        isActive: true,
        pageNumber: 2,
        pageSize: 20,
      })

      const fetchUrl = (globalThis.fetch as any).mock.calls[0][0]
      expect(fetchUrl).toContain('search=test')
      expect(fetchUrl).toContain('brandId=brand-001')
      expect(fetchUrl).toContain('isActive=true')
      expect(fetchUrl).toContain('page=2')
      expect(fetchUrl).toContain('page_size=20')
    })

    it('should handle error when listing products', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Database error' }),
      })

      await expect(productApi.listProducts()).rejects.toThrow('Database error')
    })
  })

  describe('getProduct', () => {
    it('should get product by ID', async () => {
      const mockProduct = {
        id: 'prod-001',
        sku: 'P001',
        name: 'Product 1',
        longName: 'Full Product Name 1',
        description: 'Description',
        productType: 'SIMPLE',
        brandId: 'brand-001',
        groupIds: [],
        directAttributeIds: [],
        isActive: true,
        variantsCount: 0,
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockProduct,
      })

      const result = await productApi.getProduct('prod-001')

      expect(result.id).toBe('prod-001')
      expect(result.name).toBe('Product 1')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/products/prod-001'),
        expect.objectContaining({ method: 'GET' })
      )
    })

    it('should handle error when product not found', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Product not found' }),
      })

      await expect(productApi.getProduct('invalid')).rejects.toThrow()
    })
  })

  describe('createProduct', () => {
    it('should create a new product', async () => {
      const newProduct = {
        sku: 'P002',
        name: 'New Product',
        longName: 'New Product Full Name',
        description: 'Description',
        productType: 'SIMPLE',
        brandId: 'brand-001',
        isActive: true,
      }

      const mockResponse = { id: 'prod-002', ...newProduct }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await productApi.createProduct(newProduct)

      expect(result.id).toBe('prod-002')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/products'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('P002'),
        })
      )
    })

    it('should handle error when creating product', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'SKU already exists' }),
      })

      await expect(
        productApi.createProduct({
          sku: 'DUP',
          name: 'Duplicate',
          longName: 'Duplicate',
          description: '',
          productType: 'SIMPLE',
        })
      ).rejects.toThrow('SKU already exists')
    })
  })

  describe('updateProduct', () => {
    it('should update product', async () => {
      const updateData = {
        name: 'Updated Name',
        longName: 'Updated Long Name',
        description: 'Updated description',
        isActive: false,
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: 'prod-001', ...updateData }),
      })

      const result = await productApi.updateProduct('prod-001', updateData)

      expect(result.name).toBe('Updated Name')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/products/prod-001'),
        expect.objectContaining({
          method: 'PUT',
          body: expect.stringContaining('Updated Name'),
        })
      )
    })

    it('should handle error when updating product', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Product not found' }),
      })

      await expect(productApi.updateProduct('invalid', { name: 'Test' })).rejects.toThrow()
    })
  })

  describe('changeProductStatus', () => {
    it('should change product status', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ success: true }),
      })

      await productApi.changeProductStatus('prod-001', false)

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/products/prod-001'),
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify({ is_active: false }),
        })
      )
    })

    it('should handle error when changing status', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Status change failed' }),
      })

      await expect(productApi.changeProductStatus('invalid', true)).rejects.toThrow()
    })
  })

  // ============================================================================
  // OPTION SETS & ATTRIBUTES
  // ============================================================================

  describe('getCalculatedOptionSets', () => {
    it('should get calculated option sets for product', async () => {
      const mockSets = {
        attributes: [
          { attribute_id: 'opt-001', attribute_name: 'Color', options: ['Red', 'Blue'] },
          { attribute_id: 'opt-002', attribute_name: 'Size', options: ['S', 'M', 'L'] },
        ],
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockSets,
      })

      const result = await productApi.getCalculatedOptionSets('prod-001')

      expect(result.attributes).toHaveLength(2)
      expect(result.attributes[0].attribute_name).toBe('Color')
    })

    it('should handle error when getting option sets', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({}),
      })

      await expect(productApi.getCalculatedOptionSets('invalid')).rejects.toThrow()
    })
  })

  describe('assignOptionSetToProduct', () => {
    it('should assign option set to product', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ success: true }),
      })

      await productApi.assignOptionSetToProduct('prod-001', 'optset-001')

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/products/prod-001/direct-option-sets'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ option_set_id: 'optset-001' }),
        })
      )
    })

    it('should handle error when assigning option set', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Option set not found' }),
      })

      await expect(productApi.assignOptionSetToProduct('prod-001', 'invalid')).rejects.toThrow()
    })
  })

  describe('getCalculatedAttributes', () => {
    it('should get calculated attributes for product', async () => {
      const mockAttrs = {
        attributes: [{ id: 'attr-001', name: 'Material', value: 'Cotton' }],
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockAttrs,
      })

      const result = await productApi.getCalculatedAttributes('prod-001')

      expect(result.attributes).toHaveLength(1)
    })
  })

  // ============================================================================
  // VARIANTS
  // ============================================================================

  describe('listProductVariants', () => {
    it('should list variants for product', async () => {
      const mockVariants = [
        {
          id: 'var-001',
          sku: 'P001-RED-S',
          productId: 'prod-001',
          optionConfiguration: { color: 'Red', size: 'S' },
          status: 'ACTIVE',
          isActive: true,
        },
      ]

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockVariants,
      })

      const result = await productApi.listProductVariants('prod-001')

      expect(result.variants).toHaveLength(1)
      expect(result.variants[0].sku).toBe('P001-RED-S')
    })

    it('should handle error when listing variants', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Product not found' }),
      })

      await expect(productApi.listProductVariants('invalid')).rejects.toThrow()
    })
  })

  describe('getVariant', () => {
    it('should get variant by ID', async () => {
      const mockVariant = {
        id: 'var-001',
        sku: 'VAR001',
        productId: 'prod-001',
        status: 'ACTIVE',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockVariant,
      })

      const result = await productApi.getVariant('var-001')

      expect(result.id).toBe('var-001')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/variants/var-001'),
        expect.objectContaining({ method: 'GET' })
      )
    })
  })

  describe('getVariantBySku', () => {
    it('should get variant by SKU', async () => {
      const mockVariant = { id: 'var-001', sku: 'VAR001' }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockVariant,
      })

      const result = await productApi.getVariantBySku('VAR001')

      expect(result.sku).toBe('VAR001')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/variants?sku=VAR001'),
        expect.objectContaining({ method: 'GET' })
      )
    })

    it('should handle error when variant SKU not found', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Variant not found' }),
      })

      await expect(productApi.getVariantBySku('INVALID')).rejects.toThrow()
    })
  })

  describe('generateVariants', () => {
    it('should generate variants for product', async () => {
      const mockResponse = {
        created: 6,
        variants: [{ id: 'var-001' }, { id: 'var-002' }],
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await productApi.generateVariants('prod-001', {
        inheritFromProduct: true,
      })

      expect(result.created).toBe(6)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/products/prod-001/variants/generate'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('inheritFromProduct'),
        })
      )
    })

    it('should handle error when generating variants', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'No option sets defined' }),
      })

      await expect(productApi.generateVariants('prod-001')).rejects.toThrow()
    })
  })

  describe('updateVariant', () => {
    it('should update variant', async () => {
      const updateData = { sku: 'NEW-SKU', isActive: false }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: 'var-001', ...updateData }),
      })

      await productApi.updateVariant('var-001', updateData)

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/variants/var-001'),
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(updateData),
        })
      )
    })
  })

  // ============================================================================
  // BRANDS
  // ============================================================================

  describe('listBrands', () => {
    it('should list brands successfully', async () => {
      const mockBrands = [
        { id: 'brand-001', name: 'Brand A', isActive: true, logoUrl: null },
        { id: 'brand-002', name: 'Brand B', isActive: true, logoUrl: 'http://logo.url' },
      ]

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockBrands,
      })

      const result = await productApi.listBrands()

      expect(result.data).toHaveLength(2)
      expect(result.data[0].name).toBe('Brand A')
    })

    it('should handle error when listing brands', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({}),
      })

      await expect(productApi.listBrands()).rejects.toThrow()
    })
  })

  describe('getBrand', () => {
    it('should get brand by ID', async () => {
      const mockBrand = {
        id: 'brand-001',
        name: 'Brand A',
        isActive: true,
        logoUrl: null,
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockBrand,
      })

      const result = await productApi.getBrand('brand-001')

      expect(result.id).toBe('brand-001')
      expect(result.name).toBe('Brand A')
    })
  })

  describe('createBrand', () => {
    it('should create a new brand', async () => {
      const newBrand = { name: 'New Brand', isActive: true }
      const mockResponse = { id: 'brand-003', ...newBrand, logoUrl: null }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await productApi.createBrand(newBrand)

      expect(result.id).toBe('brand-003')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/brands'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('New Brand'),
        })
      )
    })

    it('should handle error when creating brand', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Brand name already exists' }),
      })

      await expect(productApi.createBrand({ name: 'Duplicate' })).rejects.toThrow()
    })
  })

  describe('updateBrand', () => {
    it('should update brand', async () => {
      const updateData = { name: 'Updated Brand', isActive: false }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: 'brand-001', ...updateData }),
      })

      await productApi.updateBrand('brand-001', updateData)

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/brands/brand-001'),
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(updateData),
        })
      )
    })
  })

  // ============================================================================
  // PRODUCT GROUPS
  // ============================================================================

  describe('listProductGroups', () => {
    it('should list product groups successfully', async () => {
      const mockGroups = [
        {
          id: 'group-001',
          name: 'Electronics',
          type: 'TANGIBLE',
          isActive: true,
          parentGroupId: null,
          description: 'Electronic products',
        },
      ]

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockGroups,
      })

      const result = await productApi.listProductGroups()

      expect(result.data).toHaveLength(1)
      expect(result.data[0].name).toBe('Electronics')
      expect(result.data[0].type).toBe('TANGIBLE')
    })

    it('should handle error when listing groups', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({}),
      })

      await expect(productApi.listProductGroups()).rejects.toThrow()
    })
  })

  describe('getProductGroup', () => {
    it('should get product group by ID', async () => {
      const mockGroup = {
        id: 'group-001',
        name: 'Electronics',
        type: 'TANGIBLE',
        isActive: true,
        parentGroupId: null,
        description: 'Electronic products',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockGroup,
      })

      const result = await productApi.getProductGroup('group-001')

      expect(result.id).toBe('group-001')
      expect(result.name).toBe('Electronics')
      expect(result.type).toBe('TANGIBLE')
    })
  })

  describe('createProductGroup', () => {
    it('should create a new product group', async () => {
      const newGroup = {
        name: 'Furniture',
        type: 'TANGIBLE',
        isActive: true,
        description: 'Furniture products',
      }

      const mockResponse = { id: 'group-002', ...newGroup, parentGroupId: null }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await productApi.createProductGroup(newGroup)

      expect(result.id).toBe('group-002')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/product-groups'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('Furniture'),
        })
      )
    })

    it('should create a SERVICE type product group', async () => {
      const newGroup = {
        name: 'Consulting Services',
        type: 'SERVICE',
        isActive: true,
      }

      const mockResponse = { id: 'group-003', ...newGroup, parentGroupId: null }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await productApi.createProductGroup(newGroup)

      expect(result.id).toBe('group-003')
      expect(result.name).toBe('Consulting Services')
      const fetchCall = (globalThis.fetch as any).mock.calls[(globalThis.fetch as any).mock.calls.length - 1]
      const body = JSON.parse(fetchCall[1].body)
      expect(body.type).toBe('SERVICE')
    })

    it('should handle error when creating group', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Group name already exists' }),
      })

      await expect(productApi.createProductGroup({ name: 'Duplicate' })).rejects.toThrow()
    })
  })

  describe('updateProductGroup', () => {
    it('should update product group', async () => {
      const updateData = { name: 'Updated Group', description: 'New description' }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: 'group-001', ...updateData }),
      })

      await productApi.updateProductGroup('group-001', updateData)

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/product-groups/group-001'),
        expect.objectContaining({
          method: 'PUT',
          body: expect.stringContaining('Updated Group'),
        })
      )
    })
  })

  // ============================================================================
  // ATTRIBUTES
  // ============================================================================

  describe('listAttributes', () => {
    it('should list attributes successfully', async () => {
      const mockAttributes = {
        data: [
          { id: 'attr-001', name: 'Color', type: 'OPTION_SET', isActive: true },
          { id: 'attr-002', name: 'Material', type: 'TEXT', isActive: true },
        ],
        total: 2,
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockAttributes,
      })

      const result = await productApi.listAttributes()

      expect(result.data).toHaveLength(2)
      expect(result.data[0].name).toBe('Color')
    })

    it('should handle error when listing attributes', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({}),
      })

      await expect(productApi.listAttributes()).rejects.toThrow()
    })
  })

  describe('getAttribute', () => {
    it('should get attribute by ID', async () => {
      const mockAttr = {
        id: 'attr-001',
        name: 'Color',
        type: 'OPTION_SET',
        isActive: true,
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockAttr,
      })

      const result = await productApi.getAttribute('attr-001')

      expect(result.id).toBe('attr-001')
      expect(result.name).toBe('Color')
    })
  })

  describe('createAttribute', () => {
    it('should create a new attribute', async () => {
      const newAttr = {
        name: 'Weight',
        code: 'weight',
        type: 'NUMBER' as const,
        scope: 'PRODUCT' as const,
      }

      const mockResponse = { id: 'attr-003', ...newAttr }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await productApi.createAttribute(newAttr)

      expect(result.id).toBe('attr-003')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/attributes'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('Weight'),
        })
      )
    })

    it('should handle error when creating attribute', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Attribute name already exists' }),
      })

      await expect(
        productApi.createAttribute({
          name: 'Duplicate',
          code: 'duplicate',
        })
      ).rejects.toThrow()
    })
  })

  describe('updateAttribute', () => {
    it('should update attribute', async () => {
      const updateData = { name: 'Updated Attribute', isActive: false }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: 'attr-001', ...updateData }),
      })

      await productApi.updateAttribute('attr-001', updateData)

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/attributes/attr-001'),
        expect.objectContaining({
          method: 'PUT',
          body: expect.stringContaining('Updated Attribute'),
        })
      )
    })
  })

  // ============================================================================
  // ERROR HANDLING & AUTHENTICATION
  // ============================================================================

  describe('Authentication and headers', () => {
    it('should include auth token in headers', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })

      await productApi.listProducts()

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: 'Bearer test-token',
          }),
        })
      )
    })

    it('should work without auth token', async () => {
      localStorage.removeItem('tramatex_auth_token')

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })

      await productApi.listProducts()

      const headers = (globalThis.fetch as any).mock.calls[0][1].headers
      expect(headers.Authorization).toBeUndefined()
    })

    it('should include user ID in headers', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })

      await productApi.listProducts()

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            'X-User-ID': 'user-123',
          }),
        })
      )
    })

    it('should use "anonymous" when user ID not available', async () => {
      localStorage.removeItem('tramatex_user')

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })

      await productApi.listProducts()

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            'X-User-ID': 'anonymous',
          }),
        })
      )
    })

    it('should handle network errors gracefully', async () => {
      ;(globalThis.fetch as any).mockRejectedValueOnce(new Error('Network error'))

      await expect(productApi.listProducts()).rejects.toThrow(/No se pudo conectar/)
    })
  })
})
