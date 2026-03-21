/**
 * Sales Module Type Definitions
 * Defines TypeScript interfaces for Quote, Order, DeliveryNote, Invoice entities
 */

// ============================================================================
// ENUMS & LITERALS
// ============================================================================

export type QuoteStatus = 'DRAFT' | 'ISSUED' | 'ACCEPTED' | 'REJECTED' | 'EXPIRED' | 'CONVERTED'
export type OrderStatus = 'PENDING' | 'CONFIRMED' | 'PARTIALLY_DELIVERED' | 'DELIVERED' | 'CANCELLED' | 'PARTIALLY_INVOICED' | 'INVOICED'
export type DeliveryNoteStatus = 'PENDING' | 'DELIVERED' | 'CANCELLED'
export type InvoiceStatus = 'DRAFT' | 'ISSUED' | 'PAID' | 'OVERDUE' | 'VOID'
export type InvoiceType = 'STANDARD' | 'SIMPLIFIED'
export type PaymentMethod = 'CASH' | 'TRANSFER' | 'CREDIT_CARD' | 'CHECK' | 'OTHER'

export interface MesWorkRefItem {
  id: string
  workSetupId?: string | null
  workOrderId?: string | null
  description: string
}

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
  mesWorkRefs?: MesWorkRefItem[]
  generated_order_id?: string | null
  generated_order_number?: string | null
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
  sourceQuoteNumber?: string
  status: OrderStatus
  order_date: string
  delivery_date: string | null
  line_items: OrderLineItem[]
  subtotal: number
  discount_total: number
  tax_total: number
  total: number
  notes: string | null
  mesWorkRefs?: MesWorkRefItem[]
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
  series_code: string
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
  partyId: string
  expirationDate: string
  items: CreateQuoteLineItemRequest[]
  mesWorkRefs?: { workSetupId: string | null; description: string }[]
  notes?: string
}

export interface CreateQuoteLineItemRequest {
  productVariantId: string
  quantity: number
  unitPrice?: { amount: number; currency: string }
  discountPerUnit?: { amount: number; currency: string }
}

export interface UpdateQuoteRequest {
  expirationDate?: string
  mesWorkRefs?: { workSetupId: string | null; description: string }[]
  notes?: string
  items?: CreateQuoteLineItemRequest[]
}

export interface ConvertQuoteToOrderRequest {
  deliveryDate: string
}

export interface CreateOrderRequest {
  partyId: string
  quoteId?: string
  deliveryDate?: string
  items: CreateOrderLineItemRequest[]
  mesWorkRefs?: { workSetupId: string | null; description: string }[]
  notes?: string
}

export interface CreateOrderLineItemRequest {
  productVariantId: string
  quantity: number
  unitPrice?: { amount: number; currency: string }
  discountPerUnit?: { amount: number; currency: string }
}

export interface UpdateOrderRequest {
  partyId?: string
  deliveryDate?: string | null
  notes?: string
}

export interface UpdateOrderLineItemRequest {
  quantity?: number
  unitPrice?: { amount: number; currency: string }
  discountPerUnit?: { amount: number; currency: string }
}

export interface CreateDeliveryNoteRequest {
  salesOrderId: string
  deliveryDate: string
  items: CreateDeliveryLineItemRequest[]
  notes?: string
}

export interface CreateDeliveryLineItemRequest {
  salesOrderLineItemId: string
  deliveredQuantity: number
}

export interface CreateInvoiceRequest {
  partyId: string
  salesOrderIds: string[]
  deliveryNoteIds: string[]
  invoiceDate: string
  dueDate: string
  paymentTerms?: string
}

export interface CreateSimplifiedInvoiceRequest {
  partyId: string
  invoiceDate: string
  items: Array<{
    productVariantId: string
    quantity: number
    discountPercent?: number
  }>
}

// CreateInvoiceLineItemRequest removed — invoices are created from orders/delivery notes, not line-by-line

// ============================================================================
// FILTERS & PAGINATION
// ============================================================================

export interface ListQuotesFilters {
  searchText?: string
  partyId?: string
  status?: QuoteStatus
  fromDate?: string
  toDate?: string
  limit?: number
  pageNumber?: number
  pageSize?: number
}

export interface ListOrdersFilters {
  searchText?: string
  partyId?: string
  status?: OrderStatus
  fromDate?: string
  toDate?: string
  limit?: number
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
  limit?: number
  pageNumber?: number
  pageSize?: number
}

export interface ListInvoicesFilters {
  searchText?: string
  partyId?: string
  orderId?: string
  deliveryNoteId?: string
  status?: InvoiceStatus
  invoiceType?: InvoiceType
  fromDate?: string
  toDate?: string
  limit?: number
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
