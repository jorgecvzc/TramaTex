/**
 * Product Module API Service
 * Handles communication with the backend Product module endpoints
 */

import { getApiBase } from './apiBase'
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

const API_BASE = getApiBase()

interface ProductApiError extends Error {
  status?: number
  data?: unknown
  cause?: Error
}

// Frontend models with snake_case (UI layer compatibility)
interface ProductUI {
  id: string
  sku: string
  name: string
  long_name: string
  description: string
  product_type: string
  brand_id: string | null
  group_ids: string[]
  direct_attribute_ids: string[]
  is_active: boolean
  variants_count: number
}

interface VariantUI {
  id: string
  sku: string
  product_id: string
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

class ProductApiService {
  private baseUrl: string
  private brandsUrl: string
  private groupsUrl: string
  private attributesUrl: string

  constructor() {
    this.baseUrl = `${API_BASE}/products`
    this.brandsUrl = `${API_BASE}/brands`
    this.groupsUrl = `${API_BASE}/product-groups`
    this.attributesUrl = `${API_BASE}/attributes`
  }

  /**
   * Get authorization header with user token
   */
  private getHeaders(additionalHeaders: Record<string, string> = {}): Record<string, string> {
    const token = localStorage.getItem('tramatex_auth_token')
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'X-User-ID': this.getCurrentUserId(),
      ...additionalHeaders,
    }
    
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }
    
    return headers
  }

  /**
   * Get current user ID from auth context
   */
  private getCurrentUserId(): string {
    try {
      const userStr = localStorage.getItem('tramatex_user')
      if (userStr) {
        const user = JSON.parse(userStr)
        return user.id || 'anonymous'
      }
    } catch (error) {
      console.error('[ProductApi] Error parsing user:', error)
    }
    return 'anonymous'
  }

  /**
   * Handle API errors
   */
  private async handleError(response: Response, message: string): Promise<never> {
    let errorData: { error?: string; message?: string } | undefined
    try {
      errorData = await response.json()
    } catch {
      errorData = { message: 'Ocurrió un error inesperado' }
    }

    // El backend puede enviar el error en 'error' o 'message'
    const errorMessage = errorData?.error || errorData?.message || message
    const error = new Error(errorMessage) as ProductApiError
    error.status = response.status
    error.data = errorData
    throw error
  }

  private async safeFetch(url: string, options: RequestInit, fallbackMessage?: string): Promise<Response> {
    try {
      return await fetch(url, options)
    } catch (error) {
      const message =
        fallbackMessage ||
        `No se pudo conectar con el servidor. Verifica tu conexión o que la API esté activa. (URL: ${url})`
      const err = new Error(message) as ProductApiError
      err.cause = error as Error
      throw err
    }
  }

  // ============================================================================
  // PRODUCT ENDPOINTS
  // ============================================================================

  /**
   * List products with filters and pagination
   */
  async listProducts(filters: ListProductsFilters = {}): Promise<PaginatedResponse<ProductUI>> {
    const params = new URLSearchParams()

    if (filters.search) {
      params.append('search', filters.search)
    }
    if (filters.brandId) params.append('brandId', filters.brandId)
    if (filters.groupId) params.append('groupId', filters.groupId)
    if (filters.isActive !== undefined && filters.isActive !== '') {
      params.append('isActive', String(filters.isActive))
    }
    if (filters.type) params.append('productType', filters.type)
    if (filters.pageNumber) params.append('page', filters.pageNumber.toString())
    if (filters.pageSize) params.append('page_size', filters.pageSize.toString())

    const url = params.toString() ? `${this.baseUrl}?${params}` : this.baseUrl

    const response = await this.safeFetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar los productos')
    }

    const payload = await response.json()
    
    // Backend devuelve un array directo, no una estructura con data/total
    const rawProducts = Array.isArray(payload) ? payload : (payload.data || [])
    
    // Transformar de camelCase (backend) a snake_case (frontend)
    const products: ProductUI[] = rawProducts.map((p: any) => ({
      id: p.id,
      sku: p.sku,
      name: p.name,
      long_name: p.longName,
      description: p.description,
      product_type: p.productType,
      base_price: p.basePrice,
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
  }

  /**
   * Get product by ID
   */
  async getProduct(id: string): Promise<ProductUI> {
    const response = await this.safeFetch(`${this.baseUrl}/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'Producto no encontrado')
    }

    const data = await response.json()
    return {
      id: data.id,
      sku: data.sku,
      name: data.name,
      long_name: data.longName,
      description: data.description,
      product_type: data.productType,
      base_price: data.basePrice,
      brand_id: data.brandId,
      group_ids: data.groupIds || [],
      direct_attribute_ids: data.directAttributeIds || [],
      is_active: data.isActive,
      variants_count: data.variantsCount || 0,
    }
  }

  /**
   * Create a new product
   */
  async createProduct(data: {
    id?: string
    sku: string
    name: string
    longName: string
    description: string
    productType: string
    basePrice: number
    brandId?: string
    groupIds?: string[]
    directAttributeIds?: string[]
    isActive?: boolean
  }): Promise<any> {
    const response = await this.safeFetch(this.baseUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: data.id,
        sku: data.sku,
        name: data.name,
        long_name: data.longName,
        description: data.description,
        product_type: data.productType,
        base_price: data.basePrice,
        brand_id: data.brandId,
        group_ids: data.groupIds || [],
        direct_attribute_ids: data.directAttributeIds || [],
        is_active: data.isActive !== undefined ? data.isActive : true,
      }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear el producto')
    }

    return response.json()
  }

  /**
   * Update product
   */
  async updateProduct(id: string, data: {
    name?: string
    longName?: string
    sku?: string
    barcode?: string
    basePrice?: number
    productType?: string
    description?: string
    brandId?: string
    groupIds?: string[]
    directAttributeIds?: string[]
    isActive?: boolean
  }): Promise<ProductUI> {
    const response = await this.safeFetch(`${this.baseUrl}/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify({
        name: data.name,
        long_name: data.longName,
        sku: data.sku,
        barcode: data.barcode,
        base_price: data.basePrice,
        product_type: data.productType,
        description: data.description,
        brand_id: data.brandId,
        group_ids: data.groupIds,
        direct_attribute_ids: data.directAttributeIds,
        is_active: data.isActive,
      }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar el producto')
    }

    const updated = await response.json()
    return {
      id: updated.id,
      sku: updated.sku,
      name: updated.name,
      long_name: updated.longName,
      barcode: updated.barcode,
      basePrice: updated.basePrice,
      description: updated.description,
      product_type: updated.productType,
      brand_id: updated.brandId,
      group_ids: updated.groupIds || [],
      direct_attribute_ids: updated.directAttributeIds || [],
      is_active: updated.isActive,
      variants_count: updated.variantsCount || 0,
    }
  }

  /**
   * Change product status (activate/deactivate)
   */
  async changeProductStatus(id: string, isActive: boolean): Promise<any> {
    const response = await this.safeFetch(`${this.baseUrl}/${id}/status`, {
      method: 'PATCH',
      headers: this.getHeaders(),
      body: JSON.stringify({ is_active: isActive }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar el estado')
    }

    return response.json()
  }

  /**
   * Get calculated option sets for a product (inherited + direct)
   */
  async getCalculatedOptionSets(productId: string): Promise<{ attributes: CalculatedOptionSet[] }> {
    const response = await this.safeFetch(
      `${this.baseUrl}/${productId}/calculated-option-sets`,
      {
        method: 'GET',
        headers: this.getHeaders(),
      },
    )

    if (!response.ok) {
      await this.handleError(
        response,
        'No se pudieron cargar los atributos del producto',
      )
    }

    const data = await response.json()
    const attributes = Array.isArray(data) ? data : (data.attributes || [])
    return { attributes }
  }

  /**
   * Assign option set directly to product
   */
  async assignOptionSetToProduct(productId: string, optionSetId: string): Promise<any> {
    const response = await this.safeFetch(
      `${this.baseUrl}/${productId}/direct-option-sets`,
      {
        method: 'POST',
        headers: this.getHeaders(),
        body: JSON.stringify({ option_set_id: optionSetId }),
      },
    )

    if (!response.ok) {
      await this.handleError(
        response,
        'No se pudo asignar el atributo al producto',
      )
    }

    return response.json()
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
    const params = new URLSearchParams()

    if (filters.isActive !== undefined) {
      params.append('is_active', String(filters.isActive))
    }
    if (filters.pageNumber) params.append('page', filters.pageNumber.toString())
    if (filters.pageSize) params.append('page_size', filters.pageSize.toString())

    const url = params.toString()
      ? `${this.baseUrl}/${productId}/variants?${params}`
      : `${this.baseUrl}/${productId}/variants`

    const response = await this.safeFetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las variantes')
    }

    const payload = await response.json()
    const rawVariants = Array.isArray(payload) ? payload : (payload.data || [])
    
    const variants: VariantUI[] = rawVariants.map((v: any) => ({
      id: v.id,
      sku: v.sku,
      product_id: v.productId,
      option_configuration: v.optionConfiguration,
      status: v.status,
      is_active: v.isActive,
    }))
    
    return {
      data: variants,
      variants: variants,
      total: variants.length,
      page: filters.pageNumber || 1,
      page_size: filters.pageSize || 10,
    }
  }

  /**
   * Get variant by ID
   */
  async getVariant(variantId: string): Promise<any> {
    const response = await this.safeFetch(
      `${this.baseUrl}/variants/${variantId}`,
      {
        method: 'GET',
        headers: this.getHeaders(),
      },
    )

    if (!response.ok) {
      await this.handleError(response, 'Variante no encontrada')
    }

    return response.json()
  }

  /**
   * Get variant by SKU
   */
  async getVariantBySku(sku: string): Promise<any> {
    const response = await this.safeFetch(
      `${this.baseUrl}/variants?sku=${sku}`,
      {
        method: 'GET',
        headers: this.getHeaders(),
      },
    )

    if (!response.ok) {
      await this.handleError(response, 'Variante no encontrada')
    }

    return response.json()
  }

  /**
   * Find or create variant (JIT creation)
   */
  async findOrCreateVariant(productId: string, optionConfiguration: Record<string, string>): Promise<any> {
    const response = await this.safeFetch(
      `${this.baseUrl}/${productId}/variants/find-or-create`,
      {
        method: 'POST',
        headers: this.getHeaders(),
        body: JSON.stringify({ option_configuration: optionConfiguration }),
      },
    )

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear/obtener la variante')
    }

    return response.json()
  }

  /**
   * Generate variants for product (bulk creation)
   */
  async generateVariants(productId: string, options: VariantGenerationOptions = {}): Promise<any> {
    const response = await this.safeFetch(
      `${this.baseUrl}/${productId}/variants/generate`,
      {
        method: 'POST',
        headers: this.getHeaders(),
        body: JSON.stringify(options),
      },
    )

    if (!response.ok) {
      await this.handleError(
        response,
        'No se pudo iniciar la generación de variantes',
      )
    }

    return response.json()
  }

  /**
   * Update variant
   */
  async updateVariant(variantId: string, data: any): Promise<any> {
    const response = await this.safeFetch(
      `${this.baseUrl}/variants/${variantId}`,
      {
        method: 'PUT',
        headers: this.getHeaders(),
        body: JSON.stringify(data),
      },
    )

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar la variante')
    }

    return response.json()
  }

  // ============================================================================
  // BRAND ENDPOINTS
  // ============================================================================

  /**
   * List all brands
   */
  async listBrands(filters: ListBrandsFilters = {}): Promise<{ data: BrandUI[]; total: number }> {
    const params = new URLSearchParams()

    if (filters.isActive !== undefined && filters.isActive !== '') {
      params.append('isActive', String(filters.isActive))
    }

    const url = params.toString()
      ? `${this.brandsUrl}?${params}`
      : this.brandsUrl

    const response = await this.safeFetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las marcas')
    }

    const payload = await response.json()
    const rawBrands = Array.isArray(payload) ? payload : (payload.data || [])
    
    const brands: BrandUI[] = rawBrands.map((b: any) => ({
      id: b.id,
      name: b.name,
      is_active: b.isActive,
      logo_url: b.logoUrl,
    }))
    
    return {
      data: brands,
      total: brands.length,
    }
  }

  /**
   * Get brand by ID
   */
  async getBrand(id: string): Promise<BrandUI> {
    const response = await this.safeFetch(`${this.brandsUrl}/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'Marca no encontrada')
    }

    const data = await response.json()
    return {
      id: data.id,
      name: data.name,
      is_active: data.isActive,
      logo_url: data.logoUrl,
    }
  }

  /**
   * Create brand
   */
  async createBrand(data: { id?: string; name: string; defaultMarkupPercentage?: number; isActive?: boolean }): Promise<any> {
    const response = await this.safeFetch(this.brandsUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: data.id,
        name: data.name,
        defaultMarkupPercentage: data.defaultMarkupPercentage ?? 0,
        is_active: data.isActive !== undefined ? data.isActive : true,
      }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear la marca')
    }

    return response.json()
  }

  /**
   * Update brand
   */
  async updateBrand(id: string, data: any): Promise<any> {
    const response = await this.safeFetch(`${this.brandsUrl}/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar la marca')
    }

    return response.json()
  }

  // ============================================================================
  // PRODUCT GROUP ENDPOINTS
  // ============================================================================

  /**
   * List all product groups (categories)
   */
  async listProductGroups(filters: ListProductGroupsFilters = {}): Promise<{ data: ProductGroupUI[]; total: number }> {
    const params = new URLSearchParams()

    if (filters.isActive !== undefined && filters.isActive !== '') {
      params.append('isActive', String(filters.isActive))
    }
    if (filters.parentId) {
      params.append('parentGroupId', filters.parentId)
    }

    const url = params.toString()
      ? `${this.groupsUrl}?${params}`
      : this.groupsUrl

    const response = await this.safeFetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las categorías')
    }

    const payload = await response.json()
    const rawGroups = Array.isArray(payload) ? payload : (payload.data || [])
    
    const groups: ProductGroupUI[] = rawGroups.map((g: any) => ({
      id: g.id,
      name: g.name,
      type: g.type || 'TANGIBLE',
      is_active: g.isActive,
      parent_group_id: g.parentGroupId,
      description: g.description,
    }))
    
    return {
      data: groups,
      total: groups.length,
    }
  }

  /**
   * Get product group by ID
   */
  async getProductGroup(id: string): Promise<ProductGroupUI> {
    const response = await this.safeFetch(`${this.groupsUrl}/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'Categoría no encontrada')
    }

    const data = await response.json()
    return {
      id: data.id,
      name: data.name,
      type: data.type || 'TANGIBLE',
      is_active: data.isActive,
      parent_group_id: data.parentGroupId,
      description: data.description,
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
    const response = await this.safeFetch(this.groupsUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: data.id,
        name: data.name,
        type: data.type || 'TANGIBLE',
        parent_group_id: data.parentGroupId || null,
        is_active: data.isActive !== undefined ? data.isActive : true,
      }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear la categoría')
    }

    return response.json()
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
    const response = await this.safeFetch(`${this.groupsUrl}/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify({
        name: data.name,
        type: data.type,
        parent_group_id: data.parentGroupId !== undefined ? data.parentGroupId : undefined,
        is_active: data.isActive !== undefined ? data.isActive : undefined,
      }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar la categoría')
    }

    return response.json()
  }

  // ============================================================================
  // ATTRIBUTE ENDPOINTS
  // ============================================================================

  /**
   * List all attributes
   */
  async listAttributes(filters: ListAttributesFilters = {}): Promise<{ data: any[]; total: number }> {
    const params = new URLSearchParams()

    if (filters.scope) params.append('scope', filters.scope)
    if (filters.isConfigurable !== undefined) {
      params.append('isConfigurable', String(filters.isConfigurable))
    }

    const url = params.toString()
      ? `${this.attributesUrl}?${params}`
      : this.attributesUrl

    const response = await this.safeFetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar los atributos')
    }

    const payload = await response.json()
    return {
      data: payload.data || [],
      total: payload.total || 0,
    }
  }

  /**
   * Get attribute by ID
   */
  async getAttribute(id: string): Promise<any> {
    const response = await this.safeFetch(`${this.attributesUrl}/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'Atributo no encontrado')
    }

    return response.json()
  }

  /**
   * Create attribute
   */
  async createAttribute(data: {
    name: string
    code: string
    order?: number
    values?: string[]
  }): Promise<any> {
    const response = await this.safeFetch(this.attributesUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        name: data.name,
        code: data.code,
        sortOrder: data.order || 0,
        values: data.values || [],
      }),
    })

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      const errorMessage = errorData.error || 'No se pudo crear el atributo'
      throw new Error(errorMessage)
    }

    return response.json()
  }

  /**
   * Update attribute
   */
  async updateAttribute(id: string, data: {
    name?: string
    code?: string
    order?: number
    values?: string[]
  }): Promise<any> {
    const response = await this.safeFetch(`${this.attributesUrl}/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify({
        name: data.name,
        code: data.code,
        sortOrder: data.order || 0,
        values: data.values || [],
      }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar el atributo')
    }

    return response.json()
  }
}

// Export singleton instance
export const productApi = new ProductApiService()

// Also export class for testing
export default ProductApiService
