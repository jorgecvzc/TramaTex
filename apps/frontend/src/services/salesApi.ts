/**
 * Sales API Service
 * 
 * Provides methods for managing quotes, orders, delivery notes and invoices.
 * Handles all communication with the Sales module backend.
 * 
 * @module services/salesApi
 */

import type {
  Quote,
  Order,
  DeliveryNote,
  Invoice,
  CreateQuoteRequest,
  UpdateQuoteRequest,
  CreateOrderRequest,
  UpdateOrderRequest,
  CreateDeliveryNoteRequest,
  CreateInvoiceRequest,
  CreateSimplifiedInvoiceRequest,
  ListQuotesFilters,
  ListOrdersFilters,
  ListDeliveryNotesFilters,
  ListInvoicesFilters,
  QuoteStatus,
  OrderStatus,
  DeliveryNoteStatus,
  InvoiceStatus,
  UpdateOrderLineItemRequest,
} from '../types/sales'

// ============================================================================
// STATUS MAPPING: Backend (Spanish) ↔ Frontend (English)
// ============================================================================

const backendToFrontendStatus: Record<string, string> = {
  // Quote statuses
  'BORRADOR': 'DRAFT',
  'EMITIDA': 'ISSUED',
  'APROBADA': 'ACCEPTED',
  'RECHAZADA': 'REJECTED',
  'EXPIRADA': 'EXPIRED',
  'CONVERTIDA_A_PEDIDO': 'CONVERTED',
  // Order statuses
  'PENDIENTE': 'PENDING',
  'EN_PREPARACION': 'CONFIRMED',
  'ENTREGADO_PARCIALMENTE': 'PARTIALLY_DELIVERED',
  'ENTREGADO': 'DELIVERED',
  'CANCELADO': 'CANCELLED',
  'FACTURADO_PARCIALMENTE': 'PARTIALLY_INVOICED',
  'FACTURADO_COMPLETAMENTE': 'INVOICED',
  // Delivery note statuses
  // (PENDIENTE already mapped above)
  // (ENTREGADO already mapped above)
  // (CANCELADO already mapped above)
  // Invoice statuses
  // (BORRADOR already mapped above)
  // (EMITIDA already mapped above)
  'PAGADA': 'PAID',
  'VENCIDA': 'OVERDUE',
  'ANULADA': 'VOID',
}

const frontendToBackendStatus: Record<string, string> = Object.fromEntries(
  Object.entries(backendToFrontendStatus).map(([k, v]) => [v, k])
)

function normalizeStatus(status: string): string {
  return backendToFrontendStatus[status] || status
}

function toBackendStatus(status: string): string {
  return frontendToBackendStatus[status] || status
}

function normalizeEntity<T extends Record<string, any>>(obj: T): T {
  // Normalize invoice-specific field names from backend camelCase
  if (obj && 'invoiceType' in obj) {
    obj.type = obj.invoiceType
    obj.issueDate = obj.invoiceDate
    obj.salesOrderIds = obj.relatedOrderIds || []
    obj.deliveryNoteIds = obj.relatedDeliveryNoteIds || []
  }
  return obj
}

interface SalesApiError extends Error {
  cause?: Error
}

interface MoneyAmount {
  amount: number
  currency: string
}

class SalesApi {
  private baseUrl: string

  constructor() {
    this.baseUrl = '/api/sales'
  }

  private getHeaders(additionalHeaders: Record<string, string> = {}): Record<string, string> {
    const token = localStorage.getItem('tramatex_auth_token')
    return {
      'Content-Type': 'application/json',
      ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
      ...additionalHeaders,
    }
  }

  private async safeFetch(url: string, options: RequestInit, fallbackMessage: string | null = null): Promise<Response> {
    try {
      return await fetch(url, options)
    } catch (error) {
      const message =
        fallbackMessage ||
        `No se pudo conectar con el servidor. Verifica tu conexión o que la API esté activa. (URL: ${url})`
      const err = new Error(message) as SalesApiError
      err.cause = error as Error
      throw err
    }
  }

  private async handleError(response: Response, defaultMessage: string): Promise<never> {
    let errorData: { error?: string; message?: string } | undefined
    try {
      errorData = await response.json()
    } catch {
      throw new Error(defaultMessage)
    }
    throw new Error(errorData?.error || errorData?.message || defaultMessage)
  }

  // ============================================================================
  // QUOTES ENDPOINTS
  // ============================================================================

  /**
   * Create a new quote
   */
  async createQuote(data: CreateQuoteRequest): Promise<Quote> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear la cotización')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * Get quote by ID
   */
  async getQuote(id: string): Promise<Quote> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'Cotización no encontrada')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * List quotes with filters
   */
  async listQuotes(filters: ListQuotesFilters = {}): Promise<Quote[]> {
    const params = new URLSearchParams()
    
    if (filters.searchText) params.append('search', filters.searchText)
    if (filters.partyId) params.append('partyId', filters.partyId)
    if (filters.status) params.append('status', toBackendStatus(filters.status))
    if (filters.fromDate) params.append('fromDate', filters.fromDate)
    if (filters.toDate) params.append('toDate', filters.toDate)
    if (filters.limit) params.append('limit', String(filters.limit))

    const url = params.toString() ? `${this.baseUrl}/quotes?${params}` : `${this.baseUrl}/quotes`
    
    const response = await this.safeFetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las cotizaciones')
    }

    return (await response.json()).map(normalizeEntity)
  }

  /**
   * Update quote
   */
  async updateQuote(id: string, data: UpdateQuoteRequest): Promise<Quote> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar la cotización')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * Change quote status
   */
  async changeQuoteStatus(id: string, newStatus: QuoteStatus): Promise<Quote> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes/${id}/status`, {
      method: 'PATCH',
      headers: this.getHeaders(),
      body: JSON.stringify({ newStatus: toBackendStatus(newStatus) }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo cambiar el estado de la cotización')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * Preview quote calculation — returns backend-computed line subtotals and totals
   * without persisting anything.
   */
  async previewQuoteCalculation(partyId: string, items: Array<{
    productVariantId: string
    quantity: number
    unitPrice?: { amount: number; currency: string }
    discountPercent?: number
  }>): Promise<{
    lineItems: Array<{
      productVariantId: string
      productName: string
      variantSku: string
      quantity: number
      listUnitPrice: MoneyAmount
      unitPrice: MoneyAmount
      taxRate: number
      discountPercent: number
      discountPerUnit: MoneyAmount
      subtotal: MoneyAmount
      taxAmount: MoneyAmount
    }>
    subtotal: MoneyAmount
    taxAmount: MoneyAmount
    total: MoneyAmount
  }> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes/preview`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ partyId, items }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo calcular la previsualización')
    }

    return await response.json()
  }

  /**
   * Convert quote to order
   */
  async convertQuoteToOrder(id: string, deliveryDate: string): Promise<Order> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes/${id}/convert`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ quoteId: id, deliveryDate }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo convertir la cotización a pedido')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * Accept quote and convert to order in a single operation
   */
  async acceptAndConvertQuote(id: string, deliveryDate: string): Promise<Order> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes/${id}/accept-and-convert`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ deliveryDate }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo aceptar y convertir la cotización a pedido')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * Delete a draft quote
   */
  async deleteQuote(id: string): Promise<void> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes/${id}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo eliminar el presupuesto')
    }
  }

  // ============================================================================
  // ORDERS ENDPOINTS
  // ============================================================================

  /**
   * Preview order calculation — returns backend-computed line subtotals and totals
   * without persisting anything.
   */
  async previewOrderCalculation(partyId: string, items: Array<{
    productVariantId: string
    quantity: number
    unitPrice?: { amount: number; currency: string }
    discountPercent?: number
  }>): Promise<{
    lineItems: Array<{
      productVariantId: string
      productName: string
      variantSku: string
      quantity: number
      listUnitPrice: MoneyAmount
      unitPrice: MoneyAmount
      taxRate: number
      discountPercent: number
      discountPerUnit: MoneyAmount
      subtotal: MoneyAmount
      taxAmount: MoneyAmount
    }>
    subtotal: MoneyAmount
    taxAmount: MoneyAmount
    total: MoneyAmount
  }> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/preview`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ partyId, items }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo calcular la previsualización del pedido')
    }

    return await response.json()
  }

  /**
   * Create a new sales order
   */
  async createOrder(data: CreateOrderRequest): Promise<Order> {
    const response = await this.safeFetch(`${this.baseUrl}/orders`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear el pedido')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * Get order by ID
   */
  async getOrder(id: string): Promise<Order> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'Pedido no encontrado')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * List orders with filters
   */
  async listOrders(filters: ListOrdersFilters = {}): Promise<Order[]> {
    const params = new URLSearchParams()
    
    if (filters.searchText) params.append('search', filters.searchText)
    if (filters.partyId) params.append('partyId', filters.partyId)
    if (filters.status) params.append('status', toBackendStatus(filters.status))
    if (filters.fromDate) params.append('fromDate', filters.fromDate)
    if (filters.toDate) params.append('toDate', filters.toDate)
    if (filters.limit) params.append('limit', String(filters.limit))

    const url = params.toString() ? `${this.baseUrl}/orders?${params}` : `${this.baseUrl}/orders`
    
    const response = await this.safeFetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar los pedidos')
    }

    return (await response.json()).map(normalizeEntity)
  }

  /**
   * Update order details
   */
  async updateOrder(id: string, data: UpdateOrderRequest): Promise<Order> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar el pedido')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * Change order status
   */
  async changeOrderStatus(id: string, newStatus: OrderStatus): Promise<Order> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/${id}/status`, {
      method: 'PATCH',
      headers: this.getHeaders(),
      body: JSON.stringify({ newStatus: toBackendStatus(newStatus) }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo cambiar el estado del pedido')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * Add line item to order
   */
  async addOrderLineItem(orderId: string, item: any): Promise<Order> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/${orderId}/line-items`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ item }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo agregar la línea al pedido')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * Update order line item
   */
  async updateOrderLineItem(orderId: string, lineItemId: string, updates: UpdateOrderLineItemRequest): Promise<Order> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/${orderId}/line-items/${lineItemId}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify(updates),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar la línea del pedido')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * Remove line item from order
   */
  async removeOrderLineItem(orderId: string, lineItemId: string): Promise<Order> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/${orderId}/line-items/${lineItemId}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo eliminar la línea del pedido')
    }

    return normalizeEntity(await response.json())
  }

  // ============================================================================
  // DELIVERY NOTES ENDPOINTS
  // ============================================================================

  /**
   * Create a delivery note
   */
  async createDeliveryNote(data: CreateDeliveryNoteRequest): Promise<DeliveryNote> {
    const response = await this.safeFetch(`${this.baseUrl}/delivery-notes`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear el albarán')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * Get delivery note by ID
   */
  async getDeliveryNote(id: string): Promise<DeliveryNote> {
    const response = await this.safeFetch(`${this.baseUrl}/delivery-notes/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'Albarán no encontrado')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * List delivery notes with filters
   */
  async listDeliveryNotes(filters: ListDeliveryNotesFilters = {}): Promise<DeliveryNote[]> {
    const params = new URLSearchParams()
    
    if (filters.searchText) params.append('search', filters.searchText)
    if (filters.orderId) params.append('salesOrderId', filters.orderId)
    if (filters.partyId) params.append('partyId', filters.partyId)
    if (filters.status) params.append('status', toBackendStatus(filters.status))
    if (filters.fromDate) params.append('fromDate', filters.fromDate)
    if (filters.toDate) params.append('toDate', filters.toDate)
    if (filters.limit) params.append('limit', String(filters.limit))

    const url = params.toString() 
      ? `${this.baseUrl}/delivery-notes?${params}` 
      : `${this.baseUrl}/delivery-notes`
    
    const response = await this.safeFetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar los albaranes')
    }

    return (await response.json()).map(normalizeEntity)
  }

  /**
   * Change delivery note status
   */
  async changeDeliveryNoteStatus(id: string, newStatus: DeliveryNoteStatus): Promise<DeliveryNote> {
    const response = await this.safeFetch(`${this.baseUrl}/delivery-notes/${id}/status`, {
      method: 'PATCH',
      headers: this.getHeaders(),
      body: JSON.stringify({ newStatus: toBackendStatus(newStatus) }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo cambiar el estado del albarán')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * Change invoice status
   */
  async changeInvoiceStatus(id: string, newStatus: InvoiceStatus): Promise<Invoice> {
    const response = await this.safeFetch(`${this.baseUrl}/invoices/${id}/status`, {
      method: 'PATCH',
      headers: this.getHeaders(),
      body: JSON.stringify({ newStatus: toBackendStatus(newStatus) }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo cambiar el estado de la factura')
    }

    return normalizeEntity(await response.json())
  }

  // ============================================================================
  // INVOICES ENDPOINTS
  // ============================================================================

  /**
   * Create a standard invoice
   */
  async createInvoice(data: CreateInvoiceRequest): Promise<Invoice> {
    const response = await this.safeFetch(`${this.baseUrl}/invoices`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear la factura')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * Create a simplified invoice (ticket)
   */
  async createSimplifiedInvoice(data: CreateSimplifiedInvoiceRequest): Promise<Invoice> {
    const response = await this.safeFetch(`${this.baseUrl}/invoices/simplified`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear el ticket')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * Get invoice by ID
   */
  async getInvoice(id: string): Promise<Invoice> {
    const response = await this.safeFetch(`${this.baseUrl}/invoices/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'Factura no encontrada')
    }

    return normalizeEntity(await response.json())
  }

  /**
   * List invoices with filters
   */
  async listInvoices(filters: ListInvoicesFilters = {}): Promise<Invoice[]> {
    const params = new URLSearchParams()
    
    if (filters.searchText) params.append('search', filters.searchText)
    if (filters.partyId) params.append('partyId', filters.partyId)
    if (filters.orderId) params.append('orderId', filters.orderId)
    if (filters.deliveryNoteId) params.append('deliveryNoteId', filters.deliveryNoteId)
    if (filters.invoiceType) params.append('type', filters.invoiceType)
    if (filters.status) params.append('status', toBackendStatus(filters.status))
    if (filters.fromDate) params.append('fromDate', filters.fromDate)
    if (filters.toDate) params.append('toDate', filters.toDate)
    if (filters.limit) params.append('limit', String(filters.limit))

    const url = params.toString() ? `${this.baseUrl}/invoices?${params}` : `${this.baseUrl}/invoices`
    
    const response = await this.safeFetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las facturas')
    }

    return (await response.json()).map(normalizeEntity)
  }

  // ============================================================================
  // UTILITY METHODS
  // ============================================================================

  /**
   * Format money object for display
   */
  formatMoney(money: MoneyAmount | null | undefined): string {
    if (!money || typeof money.amount === 'undefined') return '—'
    const formatter = new Intl.NumberFormat('es-ES', {
      style: 'currency',
      currency: money.currency || 'EUR',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    })
    return formatter.format(money.amount)
  }

  /**
   * Format unit price (up to 3 decimals) — for line item prices and discounts per unit
   */
  formatUnitPrice(money: MoneyAmount | null | undefined): string {
    if (!money || typeof money.amount === 'undefined') return '—'
    const formatter = new Intl.NumberFormat('es-ES', {
      style: 'currency',
      currency: money.currency || 'EUR',
      minimumFractionDigits: 2,
      maximumFractionDigits: 3,
    })
    return formatter.format(money.amount)
  }

  /**
   * Format date for API
   */
  formatDateForAPI(date: Date | string | null | undefined): string | null {
    if (!date) return null
    if (typeof date === 'string') return date
    return date.toISOString()
  }

  /**
   * Parse date from API
   */
  parseDateFromAPI(dateString: string | null | undefined): Date | null {
    if (!dateString) return null
    return new Date(dateString)
  }

  /**
   * Get status badge class
   */
  getStatusClass(status: string): string {
    const statusMap: Record<string, string> = {
      'BORRADOR': 'secondary',
      'EMITIDA': 'info',
      'APROBADA': 'success',
      'RECHAZADA': 'danger',
      'EXPIRADA': 'secondary',
      'CONVERTIDA_A_PEDIDO': 'primary',
      'PENDIENTE': 'warning',
      'EN_PREPARACION': 'info',
      'ENTREGADO_PARCIALMENTE': 'primary',
      'ENTREGADO': 'success',
      'CANCELADO': 'danger',
      'FACTURADO_PARCIALMENTE': 'info',
      'FACTURADO_COMPLETAMENTE': 'success',
      'PAGADA': 'success',
      'VENCIDA': 'danger',
      'ANULADA': 'secondary',
    }
    return statusMap[status] || 'secondary'
  }

  /**
   * Get status label in Spanish
   */
  getStatusLabel(status: string): string {
    const statusLabels: Record<string, string> = {
      'BORRADOR': 'Borrador',
      'EMITIDA': 'Emitida',
      'APROBADA': 'Aprobada',
      'RECHAZADA': 'Rechazada',
      'EXPIRADA': 'Expirada',
      'CONVERTIDA_A_PEDIDO': 'Convertida a Pedido',
      'PENDIENTE': 'Pendiente',
      'EN_PREPARACION': 'En Preparación',
      'ENTREGADO_PARCIALMENTE': 'Entregado Parcialmente',
      'ENTREGADO': 'Entregado',
      'CANCELADO': 'Cancelado',
      'FACTURADO_PARCIALMENTE': 'Facturado Parcialmente',
      'FACTURADO_COMPLETAMENTE': 'Facturado Completamente',
      'PAGADA': 'Pagada',
      'VENCIDA': 'Vencida',
      'ANULADA': 'Anulada',
    }
    return statusLabels[status] || status
  }

  /**
   * List pending work setups from confirmed orders (for MES Dashboard)
   */
  async listPendingWorkSetups(): Promise<any[]> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/pending-work-setups`, {
      method: 'GET',
      headers: this.getHeaders(),
    })
    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las configuraciones pendientes')
    }
    return response.json()
  }
}

export default new SalesApi()
