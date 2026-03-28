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
import { getApiBase } from './apiBase'

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

function normalizeEntity<T extends Record<string, any>>(obj: T): T {
  if (!obj) return obj;
  if (obj.status) obj.status = normalizeStatus(obj.status);
  
  if ('invoiceType' in obj) {
    obj.type = obj.invoiceType;
    obj.issueDate = obj.invoiceDate;
    obj.salesOrderIds = obj.relatedOrderIds || [];
    obj.deliveryNoteIds = obj.relatedDeliveryNoteIds || [];
  }
  return obj;
}

interface SalesApiError extends Error {
  cause?: Error
}

class SalesApi {
  private baseUrl: string

  constructor() {
    this.baseUrl = getApiBase() + '/sales'
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
      const message = fallbackMessage || `No se pudo conectar con el servidor. (URL: ${url})`
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

  // --- Quotes ---
  async getQuote(id: string): Promise<Quote> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes/${id}`, { method: 'GET', headers: this.getHeaders() })
    if (!response.ok) await this.handleError(response, 'No se pudo obtener el presupuesto')
    return normalizeEntity(await response.json())
  }

  async listQuotes(filters: ListQuotesFilters = {}): Promise<{ data: Quote[], total: number }> {
    const params = new URLSearchParams()
    if (filters.partyId) params.append('partyId', filters.partyId)
    if (filters.status) params.append('status', filters.status)
    if (filters.limit) params.append('limit', filters.limit.toString())

    const response = await this.safeFetch(`${this.baseUrl}/quotes?${params}`, { method: 'GET', headers: this.getHeaders() })
    if (!response.ok) await this.handleError(response, 'No se pudieron cargar los presupuestos')

    const res = await response.json()
    // Normalización idéntica a la de Pedidos (soporta array directo o envuelto)
    const rawData = Array.isArray(res) ? res : (res.data || [])
    return { data: rawData.map(normalizeEntity), total: res.total ?? rawData.length }
  }

  async createQuote(data: CreateQuoteRequest): Promise<Quote> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes`, { method: 'POST', headers: this.getHeaders(), body: JSON.stringify(data) })
    if (!response.ok) await this.handleError(response, 'No se pudo crear el presupuesto')
    return normalizeEntity(await response.json())
  }

  async updateQuote(id: string, data: UpdateQuoteRequest): Promise<Quote> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes/${id}`, { method: 'PUT', headers: this.getHeaders(), body: JSON.stringify(data) })
    if (!response.ok) await this.handleError(response, 'No se pudo actualizar el presupuesto')
    return normalizeEntity(await response.json())
  }

  async changeQuoteStatus(id: string, status: string): Promise<Quote> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes/${id}/status`, { 
      method: 'PATCH', 
      headers: this.getHeaders(), 
      body: JSON.stringify({ newStatus: frontendToBackendStatus[status] || status }) 
    })
    if (!response.ok) await this.handleError(response, 'No se pudo cambiar el estado del presupuesto')
    return normalizeEntity(await response.json())
  }

  async createOrderFromQuote(id: string): Promise<Order> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes/${id}/convert`, { 
      method: 'POST', 
      headers: this.getHeaders(),
      body: JSON.stringify({}) // Enviar body vacío para evitar errores de parseo en el backend
    })
    if (!response.ok) await this.handleError(response, 'No se pudo convertir el presupuesto en pedido')
    return normalizeEntity(await response.json())
  }

  async previewQuoteCalculation(partyId: string, items: any[]): Promise<any> {
    const response = await this.safeFetch(`${this.baseUrl}/quotes/preview-calculation`, { method: 'POST', headers: this.getHeaders(), body: JSON.stringify({ partyId, items }) })
    if (!response.ok) return null
    return await response.json()
  }

  // --- Orders ---
  async getOrder(id: string): Promise<Order> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/${id}`, { method: 'GET', headers: this.getHeaders() })
    if (!response.ok) await this.handleError(response, 'No se pudo obtener el pedido')
    return normalizeEntity(await response.json())
  }

  async listOrders(filters: ListOrdersFilters = {}): Promise<{ data: Order[], total: number }> {
    const params = new URLSearchParams()
    if (filters.partyId) params.append('partyId', filters.partyId)
    if (filters.status) params.append('status', filters.status)
    if (filters.fromDate) params.append('fromDate', filters.fromDate)
    if (filters.toDate) params.append('toDate', filters.toDate)
    if (filters.limit) params.append('limit', filters.limit.toString())
    
    const response = await this.safeFetch(`${this.baseUrl}/orders?${params}`, { method: 'GET', headers: this.getHeaders() })
    if (!response.ok) await this.handleError(response, 'No se pudieron cargar los pedidos')
    const res = await response.json()
    const rawData = Array.isArray(res) ? res : (res.data || [])
    return { data: rawData.map(normalizeEntity), total: res.total ?? rawData.length }
  }

  async createOrder(data: CreateOrderRequest): Promise<Order> {
    const response = await this.safeFetch(`${this.baseUrl}/orders`, { 
      method: 'POST', 
      headers: this.getHeaders(), 
      body: JSON.stringify({
        partyId: data.partyId,
        quoteId: data.quoteId || undefined,
        deliveryDate: data.deliveryDate,
        notes: data.notes || '',
        mesWorkRefs: data.mesWorkRefs || [],
        items: data.items || []
      }) 
    })
    if (!response.ok) await this.handleError(response, 'No se pudo crear el pedido')
    return normalizeEntity(await response.json())
  }

  async updateOrder(id: string, data: UpdateOrderRequest): Promise<Order> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/${id}`, { method: 'PUT', headers: this.getHeaders(), body: JSON.stringify(data) })
    if (!response.ok) await this.handleError(response, 'No se pudo actualizar el pedido')
    return normalizeEntity(await response.json())
  }

  async changeOrderStatus(id: string, status: string): Promise<void> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/${id}/status`, { method: 'PATCH', headers: this.getHeaders(), body: JSON.stringify({ newStatus: frontendToBackendStatus[status] || status }) })
    if (!response.ok) await this.handleError(response, 'No se pudo cambiar el estado del pedido')
  }

  async addOrderLineItem(orderId: string, item: any): Promise<void> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/${orderId}/line-items`, { method: 'POST', headers: this.getHeaders(), body: JSON.stringify({ item }) })
    if (!response.ok) await this.handleError(response, 'No se pudo añadir la línea')
  }

  async updateOrderLineItem(orderId: string, itemId: string, item: any): Promise<void> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/${orderId}/line-items/${itemId}`, { method: 'PUT', headers: this.getHeaders(), body: JSON.stringify(item) })
    if (!response.ok) await this.handleError(response, 'No se pudo actualizar la línea')
  }

  async removeOrderLineItem(orderId: string, itemId: string): Promise<void> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/${orderId}/line-items/${itemId}`, { method: 'DELETE', headers: this.getHeaders() })
    if (!response.ok) await this.handleError(response, 'No se pudo eliminar la línea')
  }

  // --- Delivery Notes ---
  async getDeliveryNote(id: string): Promise<DeliveryNote> {
    const response = await this.safeFetch(`${this.baseUrl}/delivery-notes/${id}`, { method: 'GET', headers: this.getHeaders() })
    if (!response.ok) await this.handleError(response, 'No se pudo obtener el albarán')
    return normalizeEntity(await response.json())
  }

  async listDeliveryNotes(filters: ListDeliveryNotesFilters = {}): Promise<{ data: DeliveryNote[], total: number }> {
    const params = new URLSearchParams()
    if (filters.orderId) params.append('salesOrderId', filters.orderId)
    if (filters.searchText) params.append('searchText', filters.searchText)
    const response = await this.safeFetch(`${this.baseUrl}/delivery-notes?${params}`, { method: 'GET', headers: this.getHeaders() })
    if (!response.ok) await this.handleError(response, 'No se pudieron cargar los albaranes')
    const res = await response.json()
    const rawDN = Array.isArray(res) ? res : (res.data || [])
    return { data: rawDN.map(normalizeEntity), total: res.total ?? rawDN.length }
  }

  async createDeliveryNote(data: CreateDeliveryNoteRequest): Promise<DeliveryNote> {
    const response = await this.safeFetch(`${this.baseUrl}/delivery-notes`, { method: 'POST', headers: this.getHeaders(), body: JSON.stringify(data) })
    if (!response.ok) await this.handleError(response, 'No se pudo crear el albarán')
    return normalizeEntity(await response.json())
  }

  async changeDeliveryNoteStatus(id: string, status: string): Promise<DeliveryNote> {
    const response = await this.safeFetch(`${this.baseUrl}/delivery-notes/${id}/status`, { 
      method: 'PATCH', 
      headers: this.getHeaders(), 
      body: JSON.stringify({ newStatus: frontendToBackendStatus[status] || status }) 
    })
    if (!response.ok) await this.handleError(response, 'No se pudo cambiar el estado del albarán')
    return normalizeEntity(await response.json())
  }

  // --- Invoices ---
  async getInvoice(id: string): Promise<Invoice> {
    const response = await this.safeFetch(`${this.baseUrl}/invoices/${id}`, { method: 'GET', headers: this.getHeaders() })
    if (!response.ok) await this.handleError(response, 'No se pudo obtener la factura')
    return normalizeEntity(await response.json())
  }

  async listInvoices(filters: ListInvoicesFilters = {}): Promise<{ data: Invoice[], total: number }> {
    const params = new URLSearchParams()
    if (filters.orderId) params.append('orderId', filters.orderId)
    if (filters.deliveryNoteId) params.append('deliveryNoteId', filters.deliveryNoteId)
    if (filters.searchText) params.append('searchText', filters.searchText)
    if (filters.status) params.append('status', filters.status)
    if (filters.type) params.append('type', filters.type)
    
    const response = await this.safeFetch(`${this.baseUrl}/invoices?${params}`, { method: 'GET', headers: this.getHeaders() })
    if (!response.ok) await this.handleError(response, 'No se pudieron cargar las facturas')
    const res = await response.json()
    const rawInv = Array.isArray(res) ? res : (res.data || [])
    return { data: rawInv.map(normalizeEntity), total: res.total ?? rawInv.length }
  }

  async createInvoice(data: CreateInvoiceRequest): Promise<Invoice> {
    const response = await this.safeFetch(`${this.baseUrl}/invoices`, { method: 'POST', headers: this.getHeaders(), body: JSON.stringify(data) })
    if (!response.ok) await this.handleError(response, 'No se pudo crear la factura')
    return normalizeEntity(await response.json())
  }

  async changeInvoiceStatus(id: string, status: string): Promise<Invoice> {
    const response = await this.safeFetch(`${this.baseUrl}/invoices/${id}/status`, { 
      method: 'PATCH', 
      headers: this.getHeaders(), 
      body: JSON.stringify({ newStatus: frontendToBackendStatus[status] || status }) 
    })
    if (!response.ok) await this.handleError(response, 'No se pudo cambiar el estado de la factura')
    return normalizeEntity(await response.json())
  }

  // --- Utils ---
  formatMoney(amount: any): string {
    if (!amount) return '—'
    const val = typeof amount === 'object' ? amount.amount : amount
    const currency = typeof amount === 'object' ? amount.currency : 'EUR'
    return new Intl.NumberFormat('es-ES', { style: 'currency', currency }).format(val)
  }

  formatUnitPrice(amount: any): string {
    if (!amount) return '—'
    const val = typeof amount === 'object' ? amount.amount : amount
    return new Intl.NumberFormat('es-ES', { style: 'decimal', minimumFractionDigits: 2, maximumFractionDigits: 4 }).format(val) + ' €'
  }

  async previewOrderCalculation(partyId: string, items: any[]): Promise<any> {
    const response = await this.safeFetch(`${this.baseUrl}/orders/preview-calculation`, { method: 'POST', headers: this.getHeaders(), body: JSON.stringify({ partyId, items }) })
    if (!response.ok) return null
    return await response.json()
  }

  getStatusClass(status: string): string {
    const statusMap: Record<string, string> = {
      // Common & Quotes
      'DRAFT': 'secondary', 'BORRADOR': 'secondary',
      'ISSUED': 'info', 'EMITIDA': 'info',
      'ACCEPTED': 'success', 'APPROVED': 'success', 'APROBADA': 'success',
      'REJECTED': 'danger', 'RECHAZADA': 'danger',
      'CONVERTED': 'primary', 'CONVERTIDA_A_PEDIDO': 'primary',
      'EXPIRED': 'secondary', 'EXPIRADA': 'secondary',
      // Orders & Delivery Notes
      'PENDING': 'warning', 'PENDIENTE': 'warning',
      'CONFIRMED': 'info', 'CONFIRMADO': 'info', 'EN_PREPARACION': 'info',
      'PROCESSING': 'primary', 'PROCESANDO': 'primary',
      'READY': 'success', 'LISTO_PARA_ENTREGA': 'success',
      'DELIVERED': 'success', 'ENTREGADO': 'success',
      'CANCELLED': 'danger', 'CANCELADO': 'danger',
      // Invoices
      'PAID': 'success', 'PAGADA': 'success',
      'OVERDUE': 'danger', 'VENCIDA': 'danger',
      'VOIDED': 'secondary', 'ANULADA': 'secondary'
    }
    return statusMap[status] || 'secondary'
  }

  getStatusLabel(status: string): string {
    const labels: Record<string, string> = {
      // Quotes
      'DRAFT': 'Borrador', 'BORRADOR': 'Borrador',
      'ISSUED': 'Emitido', 'EMITIDA': 'Emitido',
      'ACCEPTED': 'Aprobado', 'APPROVED': 'Aprobado', 'APROBADA': 'Aprobado',
      'REJECTED': 'Rechazado', 'RECHAZADA': 'Rechazado',
      'CONVERTED': 'Convertido', 'CONVERTIDA_A_PEDIDO': 'Convertido',
      'EXPIRED': 'Expirado', 'EXPIRADA': 'Expirado',
      // Orders & Delivery Notes
      'PENDING': 'Pendiente', 'PENDIENTE': 'Pendiente',
      'CONFIRMED': 'Confirmado', 'CONFIRMADO': 'Confirmado', 'EN_PREPARACION': 'En Preparación',
      'PROCESSING': 'En Taller', 'PROCESANDO': 'En Taller',
      'READY': 'Listo', 'LISTO_PARA_ENTREGA': 'Listo',
      'DELIVERED': 'Entregado', 'ENTREGADO': 'Entregado',
      'CANCELLED': 'Anulado', 'CANCELADO': 'Anulado',
      'PARTIALLY_DELIVERED': 'Entregado Parcial', 'ENTREGADO_PARCIALMENTE': 'Entregado Parcial',
      // Invoices
      'PAID': 'Pagada', 'PAGADA': 'Pagada',
      'OVERDUE': 'Vencida', 'VENCIDA': 'Vencida',
      'VOIDED': 'Anulada', 'ANULADA': 'Anulada',
      'INVOICED': 'Facturado', 'FACTURADO_COMPLETAMENTE': 'Facturado'
    }
    return labels[status] || status
  }

  getQuoteStatusLabel(status: string): string { return this.getStatusLabel(status); }
  getQuoteStatusClass(status: string): string { return this.getStatusClass(status); }

  formatDateForAPI(date: Date): string {
    return date.toISOString().split('T')[0]
  }
}

export default new SalesApi()
