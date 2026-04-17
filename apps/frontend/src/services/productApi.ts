/**
 * Product Module API Service
 * Handles communication with the backend Product module endpoints
 */

import { api } from './api'
import type {
  CalculatedOptionSet,
  ListProductsFilters,
  ListVariantsFilters,
  ListBrandsFilters,
  ListProductGroupsFilters,
  ListAttributesFilters,
  PaginatedResponse,
  VariantGenerationOptions,
  ProductGroupType,
} from '../types/product'

// Frontend models with snake_case (UI layer compatibility)
interface ProductUI {
  id: string
  sku: string
  name: string
  long_name: string
  description: string
  product_type: string
  base_price?: number
  tax_rate?: number
  brand_id: string | null
  group_ids: string[]
  direct_attribute_ids: string[]
  is_active: boolean
  variants_count: number
  calculated_option_sets?: CalculatedOptionSet[]
}

interface VariantUI {
  id: string
  sku: string
  product_id: string
  barcode?: string
  base_cost?: number
  price?: number
  option_configuration: Record<string, string>
  status: string
  is_active: boolean
}

interface BrandUI {
  id: string
  name: string
  defaultMarkupPercentage: number
  is_active: boolean
  logo_url: string | null
}

interface ProductGroupUI {
  id: string
  name: string
  type: ProductGroupType
  is_active: boolean
  parent_group_id: string | null
  description: string | null
}

/**
 * Smart search result from backend.
 * type: "exact_variant" | "exact_product" | "partial_match" | "product_list" | "no_match"
 */
export interface SmartSearchResult {
  type: 'exact_variant' | 'exact_product' | 'partial_match' | 'product_list' | 'no_match'
  product?: ProductUI
  variant?: VariantUI
  products?: ProductUI[]
  optionSets?: Array<{
    attributeId: string
    attributeName: string
    attributeCode: string
    values: Array<{
      id: string
      value: string
      code: string
    }>
  }>
  selectedAttributes?: Record<string, string>
  matchingVariants?: VariantUI[]
}

class ProductApiService {
  private readonly baseUrl = '/products'
  private readonly brandsUrl = '/brands'
  private readonly groupsUrl = '/product-groups'
  private readonly attributesUrl = '/attributes'
  private readonly variantsUrl = '/variants'

  private async handleError(error: any, defaultMessage: string): Promise<never> {
    const errorData = error.response?.data
    const message = errorData?.error || errorData?.message || error.message || defaultMessage
    throw new Error(message)
  }

  /**
   * Transform variant from backend camelCase to frontend snake_case
   */
  private transformVariantResponse(v: any): VariantUI {
    return {
      id: v.id,
      sku: v.sku,
      product_id: v.productId || v.product_id,
      barcode: v.barcode,
      base_cost: v.baseCost !== undefined ? v.baseCost : v.base_cost,
      price: v.price,
      option_configuration: v.optionConfiguration || v.option_configuration || {},
      status: v.status,
      is_active: v.isActive !== undefined ? v.isActive : (v.is_active !== undefined ? v.is_active : true),
    }
  }

  // ============================================================================
  // PRODUCT ENDPOINTS
  // ============================================================================

  /**
   * List products with filters and pagination
   */
  async listProducts(filters: ListProductsFilters = {}): Promise<PaginatedResponse<ProductUI>> {
    try {
      const params: Record<string, any> = {}
      if (filters.search) params.search = filters.search
      if (filters.brandId) params.brandId = filters.brandId
      if (filters.groupId) params.groupId = filters.groupId
      if (filters.isActive !== undefined && filters.isActive !== '') params.isActive = String(filters.isActive)
      if (filters.type) params.productType = filters.type
      if (filters.productType) params.productType = filters.productType
      if (filters.pageNumber) params.page = filters.pageNumber
      if (filters.pageSize) params.page_size = filters.pageSize

      const response = await api.get(this.baseUrl, { params })
      const payload = response.data
      const rawProducts = Array.isArray(payload) ? payload : (payload.data || [])
      const products: ProductUI[] = rawProducts.map((p: any) => ({
      id: p.id,
      sku: p.sku,
      name: p.name,
      long_name: p.longName,
      description: p.description,
      product_type: p.productType,
      base_price: p.basePrice,
      tax_rate: p.taxRate !== undefined ? p.taxRate : 21,
      brand_id: p.brandId,
      group_ids: p.groupIds || [],
      direct_attribute_ids: p.directAttributeIds || [],
      is_active: p.isActive,
      variants_count: p.variantsCount || 0,
    }))
      return {
        data: products,
        total: products.length,
        page: filters.pageNumber || 1,
        pageSize: filters.pageSize || 10,
        totalPages: Math.ceil(products.length / (filters.pageSize || 10)),
      }
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar los productos')
    }
  }

  /**
   * Get product by ID
   */
  async getProduct(id: string): Promise<ProductUI> {
    try {
      const data = (await api.get(`${this.baseUrl}/${id}`)).data
      let calculatedOptionSets: CalculatedOptionSet[] = []
      try {
        const optionsData = await this.getCalculatedOptionSets(id)
        calculatedOptionSets = optionsData.attributes || []
      } catch (err) {
        console.warn('[productApi] Could not load calculated option sets:', err)
      }
      return {
        id: data.id,
        sku: data.sku,
        name: data.name,
        long_name: data.longName,
        description: data.description,
        product_type: data.productType,
        base_price: data.basePrice,
        tax_rate: data.taxRate !== undefined ? data.taxRate : 21,
        brand_id: data.brandId,
        group_ids: data.groupIds || [],
        direct_attribute_ids: data.directAttributeIds || [],
        is_active: data.isActive,
        variants_count: data.variantsCount || 0,
        calculated_option_sets: calculatedOptionSets,
      }
    } catch (e) {
      await this.handleError(e, 'Producto no encontrado')
    }
  }

  /**
   * Create a new product
   */
  async createProduct(data: {
    id?: string
    sku: string
    name?: string
    longName?: string
    long_name?: string
    description?: string
    productType?: string
    product_type?: string
    basePrice?: number
    base_price?: number
    taxRate?: number
    tax_rate?: number
    brandId?: string
    brand_id?: string | null
    groupIds?: string[]
    group_ids?: string[]
    directAttributeIds?: string[]
    direct_attribute_ids?: string[]
    attribute_ids?: string[]
    isActive?: boolean
    is_active?: boolean
  }): Promise<any> {
    try {
      const brandId = data.brandId ?? data.brand_id ?? null
      const normalizedBrandID = brandId === '' ? null : brandId
      const response = await api.post(this.baseUrl, {
        id: data.id,
        sku: data.sku,
        name: data.name,
        long_name: data.longName ?? data.long_name ?? '',
        description: data.description,
        product_type: data.productType ?? data.product_type ?? 'TANGIBLE',
        base_price: data.basePrice ?? data.base_price ?? 0,
        tax_rate: data.taxRate ?? data.tax_rate ?? 21.0,
        brand_id: normalizedBrandID,
        group_ids: data.groupIds ?? data.group_ids ?? [],
        direct_attribute_ids: data.directAttributeIds ?? data.direct_attribute_ids ?? data.attribute_ids ?? [],
        is_active: data.isActive ?? data.is_active ?? true,
      })
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo crear el producto')
    }
  }

  /**
   * Update product
   */
  async updateProduct(id: string, data: {
    name?: string
    longName?: string
    long_name?: string
    sku?: string
    barcode?: string
    basePrice?: number
    base_price?: number
    taxRate?: number
    tax_rate?: number
    productType?: string
    product_type?: string
    description?: string
    brandId?: string
    brand_id?: string | null
    groupIds?: string[]
    group_ids?: string[]
    directAttributeIds?: string[]
    direct_attribute_ids?: string[]
    attribute_ids?: string[]
    isActive?: boolean
    is_active?: boolean
  }): Promise<ProductUI> {
    try {
      const brandId = data.brandId ?? data.brand_id
      const normalizedBrandID = brandId === '' ? null : brandId
      const response = await api.put(`${this.baseUrl}/${id}`, {
        name: data.name,
        long_name: data.longName ?? data.long_name,
        sku: data.sku,
        barcode: data.barcode,
        base_price: data.basePrice ?? data.base_price,
        tax_rate: data.taxRate ?? data.tax_rate,
        product_type: data.productType ?? data.product_type,
        description: data.description,
        brand_id: normalizedBrandID,
        group_ids: data.groupIds ?? data.group_ids,
        direct_attribute_ids: data.directAttributeIds ?? data.direct_attribute_ids ?? data.attribute_ids,
        is_active: data.isActive ?? data.is_active,
      })
      const updated = response.data
      return {
        id: updated.id,
        sku: updated.sku,
        name: updated.name,
        long_name: updated.longName,
        base_price: updated.basePrice,
        tax_rate: updated.taxRate,
        description: updated.description,
        product_type: updated.productType,
        brand_id: updated.brandId,
        group_ids: updated.groupIds || [],
        direct_attribute_ids: updated.directAttributeIds || [],
        is_active: updated.isActive,
        variants_count: updated.variantsCount || 0,
      }
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar el producto')
    }
  }

  /**
   * Change product status (activate/deactivate)
   */
  async changeProductStatus(id: string, isActive: boolean): Promise<any> {
    try {
      const response = await api.put(`${this.baseUrl}/${id}`, { is_active: isActive })
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar el estado')
    }
  }

  /**
   * Get calculated option sets for a product (inherited + direct)
   */
  async getCalculatedOptionSets(productId: string): Promise<{ attributes: CalculatedOptionSet[] }> {
    try {
      const response = await api.get(`${this.baseUrl}/${productId}/calculated-option-sets`)
      const data = response.data
      const attributes = Array.isArray(data) ? data : (data.attributes || [])
      return { attributes }
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar los atributos del producto')
    }
  }

  /**
   * Assign option set directly to product
   */
  async assignOptionSetToProduct(productId: string, optionSetId: string): Promise<any> {
    try {
      const response = await api.post(`${this.baseUrl}/${productId}/direct-option-sets`, {
        option_set_id: optionSetId,
      })
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo asignar el atributo al producto')
    }
  }

  /**
   * Get calculated attributes for a product (alias for getCalculatedOptionSets)
   */
  async getCalculatedAttributes(productId: string): Promise<{ attributes: CalculatedOptionSet[] }> {
    return this.getCalculatedOptionSets(productId)
  }

  // ============================================================================
  // PRODUCT VARIANT ENDPOINTS
  // ============================================================================

  /**
   * List variants for a product
   */
  async listProductVariants(productId: string, filters: ListVariantsFilters = {}): Promise<{
    data: VariantUI[]
    variants: VariantUI[]
    total: number
    page: number
    page_size: number
  }> {
    try {
      const params: Record<string, any> = {}
      if (filters.isActive !== undefined) params.is_active = String(filters.isActive)
      if (filters.pageNumber) params.page = filters.pageNumber
      if (filters.pageSize) params.page_size = filters.pageSize

      const response = await api.get(`${this.baseUrl}/${productId}/variants`, { params })
      const payload = response.data
      const rawVariants = Array.isArray(payload) ? payload : (payload.data || [])
      const variants: VariantUI[] = rawVariants.map((v: any) => this.transformVariantResponse(v))
      return {
        data: variants,
        variants: variants,
        total: variants.length,
        page: filters.pageNumber || 1,
        page_size: filters.pageSize || 10,
      }
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar las variantes')
    }
  }

  /**
   * Get variant by ID
   */
  async getVariant(variantId: string): Promise<VariantUI> {
    try {
      const response = await api.get(`${this.variantsUrl}/${variantId}`)
      return this.transformVariantResponse(response.data)
    } catch (e) {
      await this.handleError(e, 'Variante no encontrada')
    }
  }

  /**
   * Get variant by SKU
   */
  async getVariantBySku(sku: string): Promise<VariantUI> {
    try {
      const response = await api.get(this.variantsUrl, { params: { sku } })
      return this.transformVariantResponse(response.data)
    } catch (e) {
      await this.handleError(e, 'Variante no encontrada')
    }
  }

  /**
   * Smart search: searches by SKU, barcode, or partial reference.
   * Returns typed result for auto-resolution in sales line items.
   */
  async smartSearch(query: string): Promise<SmartSearchResult> {
    try {
      const response = await api.get(`${this.baseUrl}/smart-search`, { params: { q: query } })
      return this.transformSmartSearchResponse(response.data)
    } catch (e) {
      await this.handleError(e, 'Error en búsqueda inteligente')
    }
  }

  /**
   * Transform smart search response from backend camelCase to frontend format
   */
  private transformSmartSearchResponse(data: any): SmartSearchResult {
    const result: SmartSearchResult = {
      type: data.type,
    }

    if (data.product) {
      result.product = this.transformProductResponse(data.product)
    }
    if (data.variant) {
      result.variant = this.transformVariantResponse(data.variant)
    }
    if (data.products) {
      result.products = data.products.map((p: any) => this.transformProductResponse(p))
    }
    if (data.optionSets) {
      result.optionSets = data.optionSets.map((os: any) => ({
        attributeId: os.id,
        attributeName: os.name,
        attributeCode: os.code,
        values: (os.values || []).map((v: any) => ({
          id: v.id,
          value: v.value,
          code: v.code,
        })),
      }))
    }
    if (data.selectedAttributes) {
      result.selectedAttributes = data.selectedAttributes
    }
    if (data.matchingVariants) {
      result.matchingVariants = data.matchingVariants.map((v: any) => this.transformVariantResponse(v))
    }

    return result
  }

  /**
   * Transform product from backend camelCase to frontend snake_case
   */
  private transformProductResponse(p: any): ProductUI {
    return {
      id: p.id,
      sku: p.sku,
      name: p.name,
      long_name: p.longName || '',
      description: p.description || '',
      product_type: p.productType || '',
      base_price: p.basePrice,
      tax_rate: p.taxRate,
      brand_id: p.brandId || null,
      group_ids: p.groupIds || [],
      direct_attribute_ids: p.directAttributeIds || [],
      is_active: p.isActive,
      variants_count: p.variantsCount || 0,
      calculated_option_sets: p.calculatedOptionSets,
    }
  }

  /**
   * Find or create variant (JIT creation)
   */
  async findOrCreateVariant(productId: string, optionConfiguration: Record<string, string>): Promise<{ variant: VariantUI }> {
    try {
      const response = await api.post(`${this.baseUrl}/${productId}/variants/find-or-create`, {
        optionConfiguration,
      })
      const payload = response.data
      const variantData = payload?.variant || payload
      return { variant: this.transformVariantResponse(variantData) }
    } catch (e) {
      await this.handleError(e, 'No se pudo crear/obtener la variante')
    }
  }

  /**
   * Generate variants for product (bulk creation)
   */
  async generateVariants(productId: string, options: VariantGenerationOptions = {}): Promise<any> {
    try {
      const response = await api.post(`${this.baseUrl}/${productId}/variants/generate`, options)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo iniciar la generación de variantes')
    }
  }

  /**
   * Update variant
   */
  async updateVariant(variantId: string, data: any): Promise<VariantUI> {
    try {
      const response = await api.put(`${this.variantsUrl}/${variantId}`, data)
      return this.transformVariantResponse(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar la variante')
    }
  }

  // ============================================================================
  // BRAND ENDPOINTS
  // ============================================================================

  /**
   * List all brands
   */
  async listBrands(filters: ListBrandsFilters = {}): Promise<{ data: BrandUI[]; total: number }> {
    try {
      const params: Record<string, any> = {}
      if (filters.isActive !== undefined && filters.isActive !== '') params.isActive = String(filters.isActive)
      const response = await api.get(this.brandsUrl, { params })
      const payload = response.data
      const rawBrands = Array.isArray(payload) ? payload : (payload.data || [])
      const brands: BrandUI[] = rawBrands.map((b: any) => ({
      id: b.id,
      name: b.name,
      defaultMarkupPercentage: b.default_markup_percentage ?? b.defaultMarkupPercentage ?? 0,
      is_active: b.is_active ?? b.isActive,
      logo_url: b.logo_url ?? b.logoUrl,
    }))
      return { data: brands, total: brands.length }
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar las marcas')
    }
  }

  /**
   * Get brand by ID
   */
  async getBrand(id: string): Promise<BrandUI> {
    try {
      const b = (await api.get(`${this.brandsUrl}/${id}`)).data
      return {
        id: b.id,
        name: b.name,
        defaultMarkupPercentage: b.default_markup_percentage ?? b.defaultMarkupPercentage ?? 0,
        is_active: b.is_active ?? b.isActive,
        logo_url: b.logo_url ?? b.logoUrl,
      }
    } catch (e) {
      await this.handleError(e, 'Marca no encontrada')
    }
  }

  /**
   * Create brand
   */
  async createBrand(data: { id?: string; name: string; defaultMarkupPercentage?: number; isActive?: boolean }): Promise<any> {
    try {
      const response = await api.post(this.brandsUrl, {
        id: data.id,
        name: data.name,
        defaultMarkupPercentage: data.defaultMarkupPercentage ?? 0,
        isActive: data.isActive !== undefined ? data.isActive : true,
      })
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo crear la marca')
    }
  }

  /**
   * Update brand
   */
  async updateBrand(id: string, data: any): Promise<any> {
    try {
      const response = await api.put(`${this.brandsUrl}/${id}`, data)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar la marca')
    }
  }

  /**
   * Delete brand
   */
  async deleteBrand(id: string): Promise<void> {
    try {
      await api.delete(`${this.brandsUrl}/${id}`)
    } catch (e) {
      await this.handleError(e, 'No se pudo eliminar la marca')
    }
  }

  // ============================================================================
  // PRODUCT GROUP ENDPOINTS
  // ============================================================================

  /**
   * List all product groups (categories)
   */
  async listProductGroups(filters: ListProductGroupsFilters = {}): Promise<{ data: ProductGroupUI[]; total: number }> {
    try {
      const params: Record<string, any> = {}
      if (filters.isActive !== undefined && filters.isActive !== '') params.isActive = String(filters.isActive)
      if (filters.parentId) params.parentGroupId = filters.parentId
      const response = await api.get(this.groupsUrl, { params })
      const payload = response.data
      const rawGroups = Array.isArray(payload) ? payload : (payload.data || [])
      const groups: ProductGroupUI[] = rawGroups.map((g: any) => ({
      id: g.id,
      name: g.name,
      type: g.type || 'TANGIBLE',
      is_active: g.isActive,
      parent_group_id: g.parent_group_id,
      description: g.description,
    }))
      return { data: groups, total: groups.length }
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar las categorías')
    }
  }

  /**
   * Get product group by ID
   */
  async getProductGroup(id: string): Promise<ProductGroupUI> {
    try {
      const data = (await api.get(`${this.groupsUrl}/${id}`)).data
      return {
        id: data.id,
        name: data.name,
        type: data.type || 'TANGIBLE',
        is_active: data.isActive,
        parent_group_id: data.parent_group_id,
        description: data.description,
      }
    } catch (e) {
      await this.handleError(e, 'Categoría no encontrada')
    }
  }

  /**
   * Create product group
   */
  async createProductGroup(data: {
    id?: string
    name: string
    type?: string // TANGIBLE or SERVICE
    parentGroupId?: string
    isActive?: boolean
  }): Promise<any> {
    try {
      const response = await api.post(this.groupsUrl, {
        id: data.id,
        name: data.name,
        type: data.type || 'TANGIBLE',
        parent_id: data.parentGroupId || null,
        isActive: data.isActive !== undefined ? data.isActive : true,
      })
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo crear la categoría')
    }
  }

  /**
   * Update product group
   */
  async updateProductGroup(id: string, data: {
    name?: string
    type?: string // TANGIBLE or SERVICE
    parentGroupId?: string | null
    isActive?: boolean
  }): Promise<any> {
    try {
      const response = await api.put(`${this.groupsUrl}/${id}`, {
        name: data.name,
        type: data.type,
        parent_id: data.parentGroupId || undefined,
        clear_parent: data.parentGroupId === null ? true : undefined,
        isActive: data.isActive !== undefined ? data.isActive : undefined,
      })
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar la categoría')
    }
  }

  /**
   * Delete product group
   */
  async deleteProductGroup(id: string): Promise<void> {
    try {
      await api.delete(`${this.groupsUrl}/${id}`)
    } catch (e) {
      await this.handleError(e, 'No se pudo eliminar la categoría')
    }
  }

  // ============================================================================
  // ATTRIBUTE ENDPOINTS
  // ============================================================================

  /**
   * List all attributes
   */
  async listAttributes(filters: ListAttributesFilters = {}): Promise<{ data: any[]; total: number }> {
    try {
      const params: Record<string, any> = {}
      if (filters.scope) params.scope = filters.scope
      if (filters.isConfigurable !== undefined) params.isConfigurable = String(filters.isConfigurable)
      const response = await api.get(this.attributesUrl, { params })
      const payload = response.data
      return { data: payload.data || [], total: payload.total || 0 }
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar los atributos')
    }
  }

  /**
   * Get attribute by ID
   */
  async getAttribute(id: string): Promise<any> {
    try {
      const response = await api.get(`${this.attributesUrl}/${id}`)
      return response.data
    } catch (e) {
      await this.handleError(e, 'Atributo no encontrado')
    }
  }

  /**
   * Create attribute
   */
  async createAttribute(data: {
    name: string
    code: string
    values?: Array<{
      value: string
      code: string
      hasPriceModifier?: boolean
      modifierType?: string
      modifierAmount?: number
    }>
  }): Promise<any> {
    try {
      const response = await api.post(this.attributesUrl, {
        name: data.name,
        code: data.code,
        values: (data.values || []).map(v => {
          const hasPriceModifier = v.hasPriceModifier || false
          return {
            value: v.value,
            code: v.code,
            hasPriceModifier,
            modifierType: hasPriceModifier && v.modifierType ? v.modifierType : null,
            modifierAmount: hasPriceModifier ? (v.modifierAmount || 0) : 0,
          }
        }),
      })
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo crear el atributo')
    }
  }

  /**
   * Update attribute
   */
  async updateAttribute(id: string, data: {
    name?: string
    code?: string
    values?: Array<{
      id?: string
      value: string
      code: string
      hasPriceModifier?: boolean
      modifierType?: string
      modifierAmount?: number
    }>
  }): Promise<any> {
    try {
      const response = await api.put(`${this.attributesUrl}/${id}`, {
        name: data.name,
        code: data.code,
        values: (data.values || []).map(v => {
          const hasPriceModifier = v.hasPriceModifier || false
          return {
            id: v.id || null,
            value: v.value,
            code: v.code,
            hasPriceModifier,
            modifierType: hasPriceModifier && v.modifierType ? v.modifierType : null,
            modifierAmount: hasPriceModifier ? (v.modifierAmount || 0) : 0,
          }
        }),
      })
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar el atributo')
    }
  }

  /**
   * Delete attribute
   */
  async deleteAttribute(id: string): Promise<void> {
    try {
      await api.delete(`${this.attributesUrl}/${id}`)
    } catch (e) {
      await this.handleError(e, 'No se pudo eliminar el atributo')
    }
  }
}

// Export singleton instance
export const productApi = new ProductApiService()

// Also export class for testing
export default ProductApiService
