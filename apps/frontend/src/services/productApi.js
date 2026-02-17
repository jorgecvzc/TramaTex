/**
 * Product Module API Service
 * Handles communication with the backend Product module endpoints
 */

import { getApiBase } from "./apiBase.ts";

const API_BASE = getApiBase();

class ProductApiService {
  constructor() {
    this.baseUrl = `${API_BASE}/products`;
    this.brandsUrl = `${API_BASE}/brands`;
    this.groupsUrl = `${API_BASE}/product-groups`;
    this.attributesUrl = `${API_BASE}/attributes`;
  }

  /**
   * Get authorization header with user token
   */
  getHeaders(additionalHeaders = {}) {
    const token = localStorage.getItem('tramatex_auth_token');
    const headers = {
      "Content-Type": "application/json",
      "X-User-ID": this.getCurrentUserId(),
      ...additionalHeaders,
    };
    
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
    
    return headers;
  }

  /**
   * Get current user ID from auth context
   */
  getCurrentUserId() {
    try {
      const userStr = localStorage.getItem('tramatex_user');
      if (userStr) {
        const user = JSON.parse(userStr);
        return user.id || 'anonymous';
      }
    } catch (error) {
      console.error('[ProductApi] Error parsing user:', error);
    }
    return 'anonymous';
  }

  /**
   * Handle API errors
   */
  async handleError(response, message) {
    let errorData;
    try {
      errorData = await response.json();
    } catch {
      errorData = { message: "Ocurrió un error inesperado" };
    }

    // El backend puede enviar el error en 'error' o 'message'
    const errorMessage = errorData.error || errorData.message || message;
    const error = new Error(errorMessage);
    error.status = response.status;
    error.data = errorData;
    throw error;
  }

  async safeFetch(url, options, fallbackMessage) {
    try {
      return await fetch(url, options);
    } catch (error) {
      const message =
        fallbackMessage ||
        `No se pudo conectar con el servidor. Verifica tu conexión o que la API esté activa. (URL: ${url})`;
      const err = new Error(message);
      err.cause = error;
      throw err;
    }
  }

  // ============================================================================
  // PRODUCT ENDPOINTS
  // ============================================================================

  /**
   * List products with filters and pagination
   */
  async listProducts(filters = {}) {
    const params = new URLSearchParams();

    if (filters.search) {
      // Backend should search in both name and SKU
      params.append("search", filters.search);
    }
    if (filters.brandId) params.append("brandId", filters.brandId);
    if (filters.groupId) params.append("groupId", filters.groupId);
    if (filters.isActive !== undefined && filters.isActive !== "") {
      params.append("isActive", filters.isActive);
    }
    if (filters.productType) params.append("productType", filters.productType);
    if (filters.pageNumber) params.append("page", filters.pageNumber);
    if (filters.pageSize) params.append("page_size", filters.pageSize);

    const url = params.toString() ? `${this.baseUrl}?${params}` : this.baseUrl;

    const response = await this.safeFetch(url, {
      method: "GET",
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, "No se pudieron cargar los productos");
    }

    const payload = await response.json();
    
    // Backend devuelve un array directo, no una estructura con data/total
    // Adaptamos la respuesta al formato esperado por el componente
    const rawProducts = Array.isArray(payload) ? payload : (payload.data || []);
    
    // Transformar de camelCase (backend) a snake_case (frontend)
    const products = rawProducts.map(p => ({
      id: p.id,
      sku: p.sku,
      name: p.name,
      long_name: p.longName,
      description: p.description,
      product_type: p.productType,
      brand_id: p.brandId,
      group_ids: p.groupIds || [],
      direct_attribute_ids: p.directAttributeIds || [],
      is_active: p.isActive,
      variants_count: p.variantsCount || 0,
    }));
    
    return {
      data: products,
      total: products.length,
      page: filters.pageNumber || 1,
      page_size: filters.pageSize || 10,
    };
  }

  /**
   * Get product by ID
   */
  async getProduct(id) {
    const response = await this.safeFetch(`${this.baseUrl}/${id}`, {
      method: "GET",
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, "Producto no encontrado");
    }

    const data = await response.json();
    // Transformar de camelCase (backend) a snake_case (frontend)
    return {
      id: data.id,
      sku: data.sku,
      name: data.name,
      long_name: data.longName,
      description: data.description,
      product_type: data.productType,
      brand_id: data.brandId,
      group_ids: data.groupIds || [],
      direct_attribute_ids: data.directAttributeIds || [],
      is_active: data.isActive,
      variants_count: data.variantsCount || 0,
    };
  }

  /**
   * Create a new product
   */
  async createProduct(data) {
    const response = await this.safeFetch(this.baseUrl, {
      method: "POST",
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: data.id,
        sku: data.sku,
        name: data.name,
        long_name: data.longName,
        description: data.description,
        product_type: data.productType,
        brand_id: data.brandId,
        group_ids: data.groupIds || [],
        direct_attribute_ids: data.directAttributeIds || [],
        is_active: data.isActive !== undefined ? data.isActive : true,
      }),
    });

    if (!response.ok) {
      await this.handleError(response, "No se pudo crear el producto");
    }

    return response.json();
  }

  /**
   * Update product
   */
  async updateProduct(id, data) {
    const response = await this.safeFetch(`${this.baseUrl}/${id}`, {
      method: "PUT",
      headers: this.getHeaders(),
      body: JSON.stringify({
        name: data.name,
        long_name: data.longName,
        description: data.description,
        brand_id: data.brandId,
        group_ids: data.groupIds,
        direct_attribute_ids: data.directAttributeIds,
        is_active: data.isActive,
      }),
    });

    if (!response.ok) {
      await this.handleError(response, "No se pudo actualizar el producto");
    }

    const updated = await response.json();
    // Transformar de camelCase (backend) a snake_case (frontend)
    return {
      id: updated.id,
      sku: updated.sku,
      name: updated.name,
      long_name: updated.longName,
      description: updated.description,
      product_type: updated.productType,
      brand_id: updated.brandId,
      group_ids: updated.groupIds || [],
      direct_attribute_ids: updated.directAttributeIds || [],
      is_active: updated.isActive,
      variants_count: updated.variantsCount || 0,
    };
  }

  /**
   * Change product status (activate/deactivate)
   */
  async changeProductStatus(id, isActive) {
    const response = await this.safeFetch(`${this.baseUrl}/${id}/status`, {
      method: "PATCH",
      headers: this.getHeaders(),
      body: JSON.stringify({ is_active: isActive }),
    });

    if (!response.ok) {
      await this.handleError(response, "No se pudo actualizar el estado");
    }

    return response.json();
  }

  /**
   * Get calculated option sets for a product (inherited + direct)
   */
  async getCalculatedOptionSets(productId) {
    const response = await this.safeFetch(
      `${this.baseUrl}/${productId}/calculated-option-sets`,
      {
        method: "GET",
        headers: this.getHeaders(),
      },
    );

    if (!response.ok) {
      await this.handleError(
        response,
        "No se pudieron cargar los atributos del producto",
      );
    }

    const data = await response.json();
    // Backend puede devolver array directo o objeto con attributes
    const attributes = Array.isArray(data) ? data : (data.attributes || []);
    return { attributes };
  }

  /**
   * Assign option set directly to product
   */
  async assignOptionSetToProduct(productId, optionSetId) {
    const response = await this.safeFetch(
      `${this.baseUrl}/${productId}/direct-option-sets`,
      {
        method: "POST",
        headers: this.getHeaders(),
        body: JSON.stringify({ option_set_id: optionSetId }),
      },
    );

    if (!response.ok) {
      await this.handleError(
        response,
        "No se pudo asignar el atributo al producto",
      );
    }

    return response.json();
  }

  /**
   * Get calculated attributes for a product (alias for getCalculatedOptionSets)
   * Returns inherited and direct attributes with their hierarchy
   */
  async getCalculatedAttributes(productId) {
    return this.getCalculatedOptionSets(productId);
  }

  // ============================================================================
  // PRODUCT VARIANT ENDPOINTS
  // ============================================================================

  /**
   * List variants for a product
   */
  async listProductVariants(productId, filters = {}) {
    const params = new URLSearchParams();

    if (filters.status) params.append("status", filters.status);
    if (filters.isActive !== undefined && filters.isActive !== "") {
      params.append("is_active", filters.isActive);
    }
    if (filters.pageNumber) params.append("page", filters.pageNumber);
    if (filters.pageSize) params.append("page_size", filters.pageSize);

    const url = params.toString()
      ? `${this.baseUrl}/${productId}/variants?${params}`
      : `${this.baseUrl}/${productId}/variants`;

    const response = await this.safeFetch(url, {
      method: "GET",
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, "No se pudieron cargar las variantes");
    }

    const payload = await response.json();
    // Backend puede devolver array directo o estructura con data
    const rawVariants = Array.isArray(payload) ? payload : (payload.data || []);
    
    // Transformar variantes de camelCase a snake_case si es necesario
    const variants = rawVariants.map(v => ({
      id: v.id,
      sku: v.sku,
      product_id: v.productId,
      option_configuration: v.optionConfiguration,
      status: v.status,
      is_active: v.isActive,
    }));
    
    return {
      data: variants,
      variants: variants, // Alias para compatibilidad
      total: variants.length,
      page: filters.pageNumber || 1,
      page_size: filters.pageSize || 10,
    };
  }

  /**
   * Get variant by ID
   */
  async getVariant(variantId) {
    const response = await this.safeFetch(
      `${this.baseUrl}/variants/${variantId}`,
      {
        method: "GET",
        headers: this.getHeaders(),
      },
    );

    if (!response.ok) {
      await this.handleError(response, "Variante no encontrada");
    }

    return response.json();
  }

  /**
   * Get variant by SKU
   */
  async getVariantBySku(sku) {
    const response = await this.safeFetch(
      `${this.baseUrl}/variants?sku=${sku}`,
      {
        method: "GET",
        headers: this.getHeaders(),
      },
    );

    if (!response.ok) {
      await this.handleError(response, "Variante no encontrada");
    }

    return response.json();
  }

  /**
   * Find or create variant (JIT creation)
   */
  async findOrCreateVariant(productId, optionConfiguration) {
    const response = await this.safeFetch(
      `${this.baseUrl}/${productId}/variants/find-or-create`,
      {
        method: "POST",
        headers: this.getHeaders(),
        body: JSON.stringify({ option_configuration: optionConfiguration }),
      },
    );

    if (!response.ok) {
      await this.handleError(response, "No se pudo crear/obtener la variante");
    }

    return response.json();
  }

  /**
   * Generate variants for product (bulk creation)
   */
  async generateVariants(productId, options = {}) {
    const response = await this.safeFetch(
      `${this.baseUrl}/${productId}/variants/generate`,
      {
        method: "POST",
        headers: this.getHeaders(),
        body: JSON.stringify(options),
      },
    );

    if (!response.ok) {
      await this.handleError(
        response,
        "No se pudo iniciar la generación de variantes",
      );
    }

    return response.json();
  }

  /**
   * Update variant
   */
  async updateVariant(variantId, data) {
    const response = await this.safeFetch(
      `${this.baseUrl}/variants/${variantId}`,
      {
        method: "PUT",
        headers: this.getHeaders(),
        body: JSON.stringify(data),
      },
    );

    if (!response.ok) {
      await this.handleError(response, "No se pudo actualizar la variante");
    }

    return response.json();
  }

  // ============================================================================
  // BRAND ENDPOINTS
  // ============================================================================

  /**
   * List all brands
   */
  async listBrands(filters = {}) {
    const params = new URLSearchParams();

    if (filters.isActive !== undefined && filters.isActive !== "") {
      params.append("isActive", filters.isActive);
    }

    const url = params.toString()
      ? `${this.brandsUrl}?${params}`
      : this.brandsUrl;

    const response = await this.safeFetch(url, {
      method: "GET",
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, "No se pudieron cargar las marcas");
    }

    const payload = await response.json();
    const rawBrands = Array.isArray(payload) ? payload : (payload.data || []);
    
    // Transformar de camelCase (backend) a snake_case (frontend)
    const brands = rawBrands.map(b => ({
      id: b.id,
      name: b.name,
      is_active: b.isActive,
      logo_url: b.logoUrl,
    }));
    
    return {
      data: brands,
      total: brands.length,
    };
  }

  /**
   * Get brand by ID
   */
  async getBrand(id) {
    const response = await this.safeFetch(`${this.brandsUrl}/${id}`, {
      method: "GET",
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, "Marca no encontrada");
    }

    const data = await response.json();
    // Transformar de camelCase (backend) a snake_case (frontend)
    return {
      id: data.id,
      name: data.name,
      is_active: data.isActive,
      logo_url: data.logoUrl,
    };
  }

  /**
   * Create brand
   */
  async createBrand(data) {
    const response = await this.safeFetch(this.brandsUrl, {
      method: "POST",
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: data.id,
        name: data.name,
        is_active: data.isActive !== undefined ? data.isActive : true,
      }),
    });

    if (!response.ok) {
      await this.handleError(response, "No se pudo crear la marca");
    }

    return response.json();
  }

  /**
   * Update brand
   */
  async updateBrand(id, data) {
    const response = await this.safeFetch(`${this.brandsUrl}/${id}`, {
      method: "PUT",
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      await this.handleError(response, "No se pudo actualizar la marca");
    }

    return response.json();
  }

  // ============================================================================
  // PRODUCT GROUP ENDPOINTS
  // ============================================================================

  /**
   * List all product groups (categories)
   */
  async listProductGroups(filters = {}) {
    const params = new URLSearchParams();

    if (filters.isActive !== undefined && filters.isActive !== "") {
      params.append("isActive", filters.isActive);
    }
    if (filters.parentGroupId) {
      params.append("parentGroupId", filters.parentGroupId);
    }

    const url = params.toString()
      ? `${this.groupsUrl}?${params}`
      : this.groupsUrl;

    const response = await this.safeFetch(url, {
      method: "GET",
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, "No se pudieron cargar las categorías");
    }

    const payload = await response.json();
    const rawGroups = Array.isArray(payload) ? payload : (payload.data || []);
    
    // Transformar de camelCase (backend) a snake_case (frontend)
    const groups = rawGroups.map(g => ({
      id: g.id,
      name: g.name,
      is_active: g.isActive,
      parent_group_id: g.parentGroupId,
      description: g.description,
    }));
    
    return {
      data: groups,
      total: groups.length,
    };
  }

  /**
   * Get product group by ID
   */
  async getProductGroup(id) {
    const response = await this.safeFetch(`${this.groupsUrl}/${id}`, {
      method: "GET",
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, "Categoría no encontrada");
    }

    const data = await response.json();
    // Transformar de camelCase (backend) a snake_case (frontend)
    return {
      id: data.id,
      name: data.name,
      is_active: data.isActive,
      parent_group_id: data.parentGroupId,
      description: data.description,
    };
  }

  /**
   * Create product group
   */
  async createProductGroup(data) {
    const response = await this.safeFetch(this.groupsUrl, {
      method: "POST",
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: data.id,
        name: data.name,
        parent_group_id: data.parentGroupId || null,
        is_active: data.isActive !== undefined ? data.isActive : true,
      }),
    });

    if (!response.ok) {
      await this.handleError(response, "No se pudo crear la categoría");
    }

    return response.json();
  }

  /**
   * Update product group
   */
  async updateProductGroup(id, data) {
    const response = await this.safeFetch(`${this.groupsUrl}/${id}`, {
      method: "PUT",
      headers: this.getHeaders(),
      body: JSON.stringify({
        name: data.name,
        parent_group_id: data.parentGroupId || null,
        is_active: data.isActive !== undefined ? data.isActive : true,
      }),
    });

    if (!response.ok) {
      await this.handleError(response, "No se pudo actualizar la categoría");
    }

    return response.json();
  }

  // ============================================================================
  // ATTRIBUTE ENDPOINTS
  // ============================================================================

  /**
   * List all attributes
   */
  async listAttributes(filters = {}) {
    const params = new URLSearchParams();

    if (filters.scopeBrandId)
      params.append("scope_brand_id", filters.scopeBrandId);
    if (filters.scopeGroupId)
      params.append("scope_group_id", filters.scopeGroupId);

    const url = params.toString()
      ? `${this.attributesUrl}?${params}`
      : this.attributesUrl;

    const response = await this.safeFetch(url, {
      method: "GET",
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, "No se pudieron cargar los atributos");
    }

    const payload = await response.json();
    return {
      data: payload.data || [],
      total: payload.total || 0,
    };
  }

  /**
   * Get attribute by ID
   */
  async getAttribute(id) {
    const response = await this.safeFetch(`${this.attributesUrl}/${id}`, {
      method: "GET",
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, "Atributo no encontrado");
    }

    return response.json();
  }

  /**
   * Create attribute
   */
  async createAttribute(data) {
    const response = await this.safeFetch(this.attributesUrl, {
      method: "POST",
      headers: this.getHeaders(),
      body: JSON.stringify({
        name: data.name,
        code: data.code,
        sortOrder: data.order || 0,
        values: data.values || [],
      }),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      const errorMessage = errorData.error || "No se pudo crear el atributo";
      throw new Error(errorMessage);
    }

    return response.json();
  }

  /**
   * Update attribute
   */
  async updateAttribute(id, data) {
    const response = await this.safeFetch(`${this.attributesUrl}/${id}`, {
      method: "PUT",
      headers: this.getHeaders(),
      body: JSON.stringify({
        name: data.name,
        code: data.code,
        sortOrder: data.order || 0,
        values: data.values || [],
      }),
    });

    if (!response.ok) {
      await this.handleError(response, "No se pudo actualizar el atributo");
    }

    return response.json();
  }
}

// Export singleton instance
export const productApi = new ProductApiService();

// Also export class for testing
export default ProductApiService;
