/**
 * Product Module Type Definitions
 * Defines TypeScript interfaces for Product, Variant, Brand, Group, Attribute entities
 */

// ============================================================================
// ENUMS & LITERALS
// ============================================================================

export type ProductStatus = 'ACTIVE' | 'INACTIVE'
export type ProductType = 'SIMPLE' | 'CONFIGURABLE'
export type AttributeType = 'DROPDOWN' | 'TEXT' | 'NUMBER' | 'COLOR' | 'SIZE'
export type AttributeScope = 'PRODUCT' | 'VARIANT'

// ============================================================================
// PRODUCT ENTITIES
// ============================================================================

export interface Product {
  id: string
  sku: string
  name: string
  description: string | null
  type: ProductType
  is_active: boolean
  brand_id: string | null
  group_id: string | null
  base_cost: number
  base_price: number
  base_weight: number | null
  stock_quantity: number
  min_stock_level: number
  configurable_attributes: ConfigurableAttribute[]
  created_at: string
  updated_at: string
}

export interface ConfigurableAttribute {
  attribute_id: string
  attribute_name: string
  values: string[]
}

// ============================================================================
// VARIANT ENTITIES
// ============================================================================

export interface ProductVariant {
  id: string
  product_id: string
  sku: string
  variant_attributes: VariantAttribute[]
  cost: number
  price: number
  weight: number | null
  stock_quantity: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface VariantAttribute {
  attribute_id: string
  attribute_name: string
  value: string
}

export interface VariantGenerationOptions {
  inheritFromProduct?: boolean
  baseCost?: number
  basePrice?: number
  baseWeight?: number
}

// ============================================================================
// BRAND ENTITIES
// ============================================================================

export interface Brand {
  id: string
  name: string
  code: string
  description: string | null
  logo_url: string | null
  is_active: boolean
  created_at: string
  updated_at: string
}

// ============================================================================
// PRODUCT GROUP ENTITIES
// ============================================================================

export interface ProductGroup {
  id: string
  name: string
  code: string
  description: string | null
  parent_id: string | null
  is_active: boolean
  created_at: string
  updated_at: string
}

// ============================================================================
// ATTRIBUTE ENTITIES
// ============================================================================

export interface Attribute {
  id: string
  name: string
  code: string
  type: AttributeType
  scope: AttributeScope
  is_required: boolean
  is_configurable: boolean
  options: AttributeOption[]
  default_value: string | null
  created_at: string
  updated_at: string
}

export interface AttributeOption {
  value: string
  label: string
  sort_order: number
}

export interface CalculatedOptionSet {
  attribute_id: string
  attribute_name: string
  options: string[]
}

// ============================================================================
// REQUEST DTOs
// ============================================================================

export interface CreateProductRequest {
  sku: string
  name: string
  description?: string
  type: ProductType
  brand_id?: string
  group_id?: string
  base_cost: number
  base_price: number
  base_weight?: number
  min_stock_level?: number
  configurable_attributes?: ConfigurableAttribute[]
}

export interface UpdateProductRequest {
  name?: string
  description?: string
  brand_id?: string | null
  group_id?: string | null
  base_cost?: number
  base_price?: number
  base_weight?: number | null
  min_stock_level?: number
  is_active?: boolean
}

export interface UpdateVariantRequest {
  cost?: number
  price?: number
  weight?: number | null
  stock_quantity?: number
  is_active?: boolean
}

export interface CreateBrandRequest {
  name: string
  code: string
  description?: string
  logo_url?: string
}

export interface UpdateBrandRequest {
  name?: string
  code?: string
  description?: string
  logo_url?: string
  is_active?: boolean
}

export interface CreateProductGroupRequest {
  name: string
  code: string
  description?: string
  parent_id?: string
}

export interface UpdateProductGroupRequest {
  name?: string
  code?: string
  description?: string
  parent_id?: string | null
  is_active?: boolean
}

export interface CreateAttributeRequest {
  name: string
  code: string
  type: AttributeType
  scope: AttributeScope
  is_required?: boolean
  is_configurable?: boolean
  options?: AttributeOption[]
  default_value?: string
}

export interface UpdateAttributeRequest {
  name?: string
  code?: string
  is_required?: boolean
  is_configurable?: boolean
  options?: AttributeOption[]
  default_value?: string
}

// ============================================================================
// FILTERS & PAGINATION
// ============================================================================

export interface ListProductsFilters {
  search?: string
  brandId?: string
  groupId?: string
  isActive?: boolean | ''
  type?: ProductType
  pageNumber?: number
  pageSize?: number
}

export interface ListVariantsFilters {
  isActive?: boolean
  minStock?: number
  pageNumber?: number
  pageSize?: number
}

export interface ListBrandsFilters {
  search?: string
  isActive?: boolean | ''
  pageNumber?: number
  pageSize?: number
}

export interface ListProductGroupsFilters {
  search?: string
  parentId?: string
  isActive?: boolean | ''
  pageNumber?: number
  pageSize?: number
}

export interface ListAttributesFilters {
  scope?: AttributeScope
  isConfigurable?: boolean
  pageNumber?: number
  pageSize?: number
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

// ============================================================================
// API ERROR HANDLING
// ============================================================================

export interface ProductError {
  message: string
  status?: number
  data?: unknown
  cause?: Error
}
