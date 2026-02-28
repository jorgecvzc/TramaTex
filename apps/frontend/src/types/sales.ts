/**
 * Sales Module Type Definitions
 * Defines TypeScript interfaces for Quote, Order, DeliveryNote, Invoice entities
 */

// ============================================================================
// ENUMS & LITERALS
// ============================================================================

export type QuoteStatus = 'DRAFT' | 'SENT' | 'ACCEPTED' | 'REJECTED' | 'EXPIRED'
export type OrderStatus = 'PENDING' | 'CONFIRMED' | 'IN_PRODUCTION' | 'READY_TO_SHIP' | 'SHIPPED' | 'DELIVERED' | 'CANCELLED'
export type DeliveryNoteStatus = 'DRAFT' | 'DISPATCHED' | 'DELIVERED' | 'CANCELLED'
export type InvoiceStatus = 'DRAFT' | 'ISSUED' | 'PAID' | 'OVERDUE' | 'CANCELLED'
export type InvoiceType = 'STANDARD' | 'SIMPLIFIED'
export type PaymentMethod = 'CASH' | 'TRANSFER' | 'CREDIT_CARD' | 'CHECK' | 'OTHER'

// ============================================================================
// LINE ITEM ENTITIES
// ============================================================================

export interface QuoteLineItem {
  id: string
  mes_work_id?: string | null
  product_variant_id: string
  product_name: string
  quantity: number
  unit_price: number
  discount_percentage: number
  subtotal: number
  total: number
}

export interface OrderLineItem {
  id: string
  mes_work_id?: string | null
  product_variant_id: string
  product_name: string
  quantity: number
  unit_price: number
  discount_percentage: number
  subtotal: number
  total: number
  production_status: string
  notes: string | null
}

export interface DeliveryLineItem {
  id: string
  order_line_item_id: string
  product_variant_id: string
  product_name: string
  quantity_delivered: number
  notes: string | null
}

export interface InvoiceLineItem {
  id: string
  order_line_item_id: string
  product_variant_id: string
  product_name: string
  quantity: number
  unit_price: number
  discount_percentage: number
  subtotal: number
  tax_amount: number
  total: number
}

// ============================================================================
// QUOTE ENTITIES
// ============================================================================

export interface Quote {
  id: string
  quote_number: string
  party_id: string
  party_name: string
  status: QuoteStatus
  issue_date: string
  valid_until: string
  line_items: QuoteLineItem[]
  subtotal: number
  discount_total: number
  tax_total: number
  total: number
  notes: string | null
  created_at: string
  updated_at: string
}

// ============================================================================
// ORDER ENTITIES
// ============================================================================

export interface Order {
  id: string
  order_number: string
  party_id: string
  party_name: string
  quote_id: string | null
  status: OrderStatus
  order_date: string
  delivery_date: string | null
  line_items: OrderLineItem[]
  subtotal: number
  discount_total: number
  tax_total: number
  total: number
  notes: string | null
  created_at: string
  updated_at: string
}

// ============================================================================
// DELIVERY NOTE ENTITIES
// ============================================================================

export interface DeliveryNote {
  id: string
  delivery_note_number: string
  order_id: string
  order_number: string
  party_id: string
  party_name: string
  status: DeliveryNoteStatus
  dispatch_date: string
  delivery_date: string | null
  line_items: DeliveryLineItem[]
  notes: string | null
  created_at: string
  updated_at: string
}

// ============================================================================
// INVOICE ENTITIES
// ============================================================================

export interface Invoice {
  id: string
  invoice_number: string
  invoice_type: InvoiceType
  order_id: string | null
  order_number: string | null
  party_id: string
  party_name: string
  status: InvoiceStatus
  issue_date: string
  due_date: string
  payment_date: string | null
  payment_method: PaymentMethod | null
  line_items: InvoiceLineItem[]
  subtotal: number
  discount_total: number
  tax_total: number
  total: number
  notes: string | null
  created_at: string
  updated_at: string
}

// ============================================================================
// REQUEST DTOs
// ============================================================================

export interface CreateQuoteRequest {
  party_id: string
  valid_until: string
  line_items: CreateQuoteLineItemRequest[]
  notes?: string
}

export interface CreateQuoteLineItemRequest {
  mesWorkId?: string
  product_variant_id: string
  quantity: number
  unit_price: number
  discount_percentage?: number
}

export interface UpdateQuoteRequest {
  status?: QuoteStatus
  valid_until?: string
  line_items?: CreateQuoteLineItemRequest[]
  notes?: string
}

export interface ConvertQuoteToOrderRequest {
  delivery_date: string
}

export interface CreateOrderRequest {
  party_id: string
  quote_id?: string
  delivery_date?: string
  line_items: CreateOrderLineItemRequest[]
  notes?: string
}

export interface CreateOrderLineItemRequest {
  mesWorkId?: string
  product_variant_id: string
  quantity: number
  unit_price: number
  discount_percentage?: number
}

export interface UpdateOrderRequest {
  status?: OrderStatus
  delivery_date?: string | null
  line_items?: CreateOrderLineItemRequest[]
  notes?: string
}

export interface UpdateOrderLineItemRequest {
  quantity?: number
  unit_price?: number
  discount_percentage?: number
  production_status?: string
  notes?: string
}

export interface CreateDeliveryNoteRequest {
  order_id: string
  dispatch_date: string
  line_items: CreateDeliveryLineItemRequest[]
  notes?: string
}

export interface CreateDeliveryLineItemRequest {
  order_line_item_id: string
  quantity_delivered: number
  notes?: string
}

export interface CreateInvoiceRequest {
  order_id: string
  due_date: string
  line_items: CreateInvoiceLineItemRequest[]
  notes?: string
}

export interface CreateSimplifiedInvoiceRequest {
  partyId: string
  invoiceDate: string
  items: Array<{
    productVariantId: string
    quantity: number
  }>
}

export interface CreateInvoiceLineItemRequest {
  order_line_item_id: string
  quantity: number
  unit_price: number
  discount_percentage?: number
  tax_percentage?: number
}

// ============================================================================
// FILTERS & PAGINATION
// ============================================================================

export interface ListQuotesFilters {
  searchText?: string
  partyId?: string
  status?: QuoteStatus
  fromDate?: string
  toDate?: string
  pageNumber?: number
  pageSize?: number
}

export interface ListOrdersFilters {
  searchText?: string
  partyId?: string
  status?: OrderStatus
  fromDate?: string
  toDate?: string
  pageNumber?: number
  pageSize?: number
}

export interface ListDeliveryNotesFilters {
  searchText?: string
  orderId?: string
  partyId?: string
  status?: DeliveryNoteStatus
  fromDate?: string
  toDate?: string
  pageNumber?: number
  pageSize?: number
}

export interface ListInvoicesFilters {
  searchText?: string
  partyId?: string
  orderId?: string
  status?: InvoiceStatus
  invoiceType?: InvoiceType
  fromDate?: string
  toDate?: string
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

export interface SalesError {
  error: string
  message?: string
  cause?: Error
}
