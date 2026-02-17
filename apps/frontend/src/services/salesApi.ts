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
  UpdateOrderLineItemRequest,
} from '../types/sales'

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

    return await response.json()
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

    return await response.json()
  }

  /**
   * List quotes with filters
   */
  async listQuotes(filters: ListQuotesFilters = {}): Promise<Quote[]> {
    const params = new URLSearchParams()
    
    if (filters.partyId) params.append('partyId', filters.partyId)
    if (filters.status) params.append('status', filters.status)
    if (filters.fromDate) params.append('fromDate', filters.fromDate)
    if (filters.toDate) params.append('toDate', filters.toDate)

    const url = params.toString() ? `${this.baseUrl}/quotes?${params}` : `${this.baseUrl}/quotes`
    
    const response = await this.safeFetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las cotizaciones')
    }

    return await response.json()
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

    return await response.json()
  }

  /**
   * Change quote status
   */
  async changeQuoteStatus(id: string, newStatus: QuoteStatus): Promise<Quote> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes/${id}/status`, {
      method: 'PATCH',
      headers: this.getHeaders(),
      body: JSON.stringify({ newStatus }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo cambiar el estado de la cotización')
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

    return await response.json()
  }

  // ============================================================================
  // ORDERS ENDPOINTS
  // ============================================================================

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

    return await response.json()
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

    return await response.json()
  }

  /**
   * List orders with filters
   */
  async listOrders(filters: ListOrdersFilters = {}): Promise<Order[]> {
    const params = new URLSearchParams()
    
    if (filters.partyId) params.append('partyId', filters.partyId)
    if (filters.status) params.append('status', filters.status)
    if (filters.fromDate) params.append('fromDate', filters.fromDate)
    if (filters.toDate) params.append('toDate', filters.toDate)

    const url = params.toString() ? `${this.baseUrl}/orders?${params}` : `${this.baseUrl}/orders`
    
    const response = await this.safeFetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar los pedidos')
    }

    return await response.json()
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

    return await response.json()
  }

  /**
   * Change order status
   */
  async changeOrderStatus(id: string, newStatus: OrderStatus): Promise<Order> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/${id}/status`, {
      method: 'PATCH',
      headers: this.getHeaders(),
      body: JSON.stringify({ newStatus }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo cambiar el estado del pedido')
    }

    return await response.json()
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

    return await response.json()
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

    return await response.json()
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

    return await response.json()
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

    return await response.json()
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

    return await response.json()
  }

  /**
   * List delivery notes with filters
   */
  async listDeliveryNotes(filters: ListDeliveryNotesFilters = {}): Promise<DeliveryNote[]> {
    const params = new URLSearchParams()
    
    if (filters.orderId) params.append('salesOrderId', filters.orderId)
    if (filters.partyId) params.append('partyId', filters.partyId)
    if (filters.status) params.append('status', filters.status)
    if (filters.fromDate) params.append('fromDate', filters.fromDate)
    if (filters.toDate) params.append('toDate', filters.toDate)

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

    return await response.json()
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

    return await response.json()
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

    return await response.json()
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

    return await response.json()
  }

  /**
   * List invoices with filters
   */
  async listInvoices(filters: ListInvoicesFilters = {}): Promise<Invoice[]> {
    const params = new URLSearchParams()
    
    if (filters.partyId) params.append('partyId', filters.partyId)
    if (filters.orderId) params.append('orderId', filters.orderId)
    if (filters.invoiceType) params.append('type', filters.invoiceType)
    if (filters.status) params.append('status', filters.status)
    if (filters.fromDate) params.append('fromDate', filters.fromDate)
    if (filters.toDate) params.append('toDate', filters.toDate)

    const url = params.toString() ? `${this.baseUrl}/invoices?${params}` : `${this.baseUrl}/invoices`
    
    const response = await this.safeFetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las facturas')
    }

    return await response.json()
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
      'PENDING': 'warning',
      'CONFIRMED': 'info',
      'IN_PROGRESS': 'primary',
      'COMPLETED': 'success',
      'CANCELLED': 'danger',
      'DRAFT': 'secondary',
      'SENT': 'info',
      'PAID': 'success',
      'OVERDUE': 'danger',
    }
    return statusMap[status] || 'secondary'
  }

  /**
   * Get status label in Spanish
   */
  getStatusLabel(status: string): string {
    const statusLabels: Record<string, string> = {
      'PENDING': 'Pendiente',
      'CONFIRMED': 'Confirmado',
      'IN_PROGRESS': 'En Progreso',
      'COMPLETED': 'Completado',
      'CANCELLED': 'Cancelado',
      'DRAFT': 'Borrador',
      'SENT': 'Enviado',
      'PAID': 'Pagado',
      'OVERDUE': 'Vencido',
    }
    return statusLabels[status] || status
  }
}

export default new SalesApi()
