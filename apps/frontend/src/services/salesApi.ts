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
import { api } from './api'

// ============================================================================
// STATUS & TYPE HELPERS
// ============================================================================

function normalizeEntity<T extends Record<string, any>>(obj: T): T {
  if (!obj) return obj;
  
  // Base normalization (snake_case to camelCase)
  const snakeToCamelMap: Record<string, string> = {
    'party_id': 'partyId',
    'party_name': 'partyName',
    'party_reference': 'partyReference',
    'order_number': 'orderNumber',
    'quote_number': 'quoteNumber',
    'order_date': 'orderDate',
    'quote_date': 'quoteDate',
    'delivery_date': 'deliveryDate',
    'expiration_date': 'expirationDate',
    'line_items': 'lineItems',
    'mes_work_refs': 'mesWorkRefs',
    'tax_id': 'taxId',
    'invoice_number': 'invoiceNumber',
    'invoice_date': 'invoiceDate',
    'invoice_type': 'invoiceType',
    'related_order_ids': 'relatedOrderIds',
    'related_delivery_note_ids': 'relatedDeliveryNoteIds',
    'delivered_quantity': 'deliveredQuantity',
    'unit_price': 'unitPrice',
    'discount_percent': 'discountPercent',
    'product_variant_id': 'productVariantId',
    'variant_sku': 'variantSku',
    'product_name': 'productName',
    'work_setup_id': 'workSetupId',
    'work_order_id': 'workOrderId',
    'subtotal_amount': 'subtotal', // Handle variations
    'tax_amount': 'taxAmount',
    'total_amount': 'total'
  };

  Object.entries(snakeToCamelMap).forEach(([snake, camel]) => {
    if (snake in obj && !(camel in obj)) {
      (obj as any)[camel] = obj[snake];
    }
  });

  // Ensure status is normalized to uppercase for consistent UI logic
  if (obj.status && typeof obj.status === 'string') {
    obj.status = normalizeSalesStatus(obj.status);
  }

  // Deep normalization for items/refs
  if (Array.isArray(obj.lineItems)) {
    obj.lineItems = obj.lineItems.map(normalizeEntity);
  }
  if (Array.isArray(obj.mesWorkRefs)) {
    obj.mesWorkRefs = obj.mesWorkRefs.map(normalizeEntity);
  }
  if (Array.isArray(obj.items)) {
    obj.items = obj.items.map(normalizeEntity);
  }
  
  // Specific legacy fixes
  if ('invoiceType' in obj || 'invoice_type' in obj || 'type' in obj) {
    obj.issueDate = obj.invoiceDate || obj.issueDate;
    obj.salesOrderIds = obj.relatedOrderIds || obj.salesOrderIds || [];
    obj.deliveryNoteIds = obj.relatedDeliveryNoteIds || obj.deliveryNoteIds || [];
    obj.type = obj.type || obj.invoiceType;
  }

  return obj;
}

function normalizeSalesStatus(status: string): string {
  const key = String(status || "").trim().toUpperCase();
  const statusMap: Record<string, string> = {
    PENDIENTE: "PENDING",
    BORRADOR: "DRAFT",
    EMITIDA: "ISSUED",
    EMITIDO: "ISSUED",
    SENT: "ISSUED",
    ENVIADO: "ISSUED",
    ENVIADA: "ISSUED",
    APROBADO: "APPROVED",
    APROBADA: "APPROVED",
    ACEPTADO: "APPROVED",
    ACEPTADA: "APPROVED",
    ACCEPTED: "APPROVED",
    RECHAZADO: "REJECTED",
    RECHAZADA: "REJECTED",
    PAGADA: "PAID",
    PAGADO: "PAID",
    CANCELADO: "CANCELLED",
    CANCELADA: "CANCELLED",
    ENTREGADO: "DELIVERED",
    ENTREGADA: "DELIVERED",
    FACTURADO: "INVOICED",
    FACTURADA: "INVOICED",
    PREPARACION: "IN_PREPARATION",
    EN_PREPARACION: "IN_PREPARATION",
    EN_PREPARACIÓN: "IN_PREPARATION",
  };

  return statusMap[key] || key;
}

function resolveTotal(res: any, rawData: any[]): number {
  if (typeof res?.total === 'number') return res.total
  if (typeof res?.count === 'number') return res.count
  if (typeof res?.total_count === 'number') return res.total_count
  return rawData.length
}

class SalesApi {
  private readonly moduleBase = '/sales'

  private async handleError(error: any, defaultMessage: string): Promise<never> {
    const errorData = error.response?.data
    const message = errorData?.error || errorData?.message || error.message || defaultMessage
    throw new Error(message)
  }

  // --- Quotes ---
  async getQuote(id: string): Promise<Quote> {
    try {
      const response = await api.get(`${this.moduleBase}/quotes/${id}`)
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo obtener el presupuesto')
    }
  }

  async listQuotes(filters: ListQuotesFilters = {}): Promise<{ data: Quote[], total: number }> {
    const params: any = {}
    if (filters.partyId) params.partyId = filters.partyId
    if ((filters as any).searchText) params.search = (filters as any).searchText
    if ((filters as any).fromDate) params.fromDate = (filters as any).fromDate
    if ((filters as any).toDate) params.toDate = (filters as any).toDate
    if (filters.status) params.status = normalizeSalesStatus(filters.status)
    if (filters.limit) params.limit = filters.limit

    try {
      const response = await api.get(`${this.moduleBase}/quotes`, { params })
      const res = response.data
      const rawData = Array.isArray(res) ? res : (res.data || [])
      return { data: rawData.map(normalizeEntity), total: resolveTotal(res, rawData) }
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar los presupuestos')
    }
  }

  async createQuote(data: CreateQuoteRequest): Promise<Quote> {
    try {
      const response = await api.post(`${this.moduleBase}/quotes`, data)
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo crear el presupuesto')
    }
  }

  async updateQuote(id: string, data: UpdateQuoteRequest): Promise<Quote> {
    try {
      const response = await api.put(`${this.moduleBase}/quotes/${id}`, data)
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar el presupuesto')
    }
  }

  async changeQuoteStatus(id: string, status: string): Promise<Quote> {
    try {
      const response = await api.patch(`${this.moduleBase}/quotes/${id}/status`, { newStatus: status })
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo cambiar el estado del presupuesto')
    }
  }

  async createOrderFromQuote(id: string, deliveryDate?: string): Promise<Order> {
    const body: Record<string, string> = {}
    if (deliveryDate) {
      body.deliveryDate = deliveryDate.includes('T') ? deliveryDate : `${deliveryDate}T00:00:00Z`
    }
    try {
      const response = await api.post(`${this.moduleBase}/quotes/${id}/convert`, body)
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo convertir el presupuesto en pedido')
    }
  }

  async convertQuoteToOrder(id: string, deliveryDate?: string): Promise<Order> {
    return this.createOrderFromQuote(id, deliveryDate)
  }

  async previewQuoteCalculation(partyId: string, items: any[]): Promise<any> {
    try {
      const response = await api.post(`${this.moduleBase}/quotes/preview`, { partyId, items })
      return response.data
    } catch (e) {
      return null
    }
  }

  // --- Orders ---
  async getOrder(id: string): Promise<Order> {
    try {
      const response = await api.get(`${this.moduleBase}/orders/${id}`)
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo obtener el pedido')
    }
  }

  async listOrders(filters: ListOrdersFilters = {}): Promise<{ data: Order[], total: number }> {
    const params: any = {}
    if (filters.partyId) params.partyId = filters.partyId
    if ((filters as any).searchText) params.search = (filters as any).searchText
    if (filters.status) params.status = normalizeSalesStatus(filters.status)
    if (filters.fromDate) params.fromDate = filters.fromDate
    if (filters.toDate) params.toDate = filters.toDate
    if (filters.limit) params.limit = filters.limit
    
    try {
      const response = await api.get(`${this.moduleBase}/orders`, { params })
      const res = response.data
      const rawData = Array.isArray(res) ? res : (res.data || [])
      return { data: rawData.map(normalizeEntity), total: resolveTotal(res, rawData) }
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar los pedidos')
    }
  }

  async createOrder(data: CreateOrderRequest): Promise<Order> {
    const payload: any = data as any
    try {
      const response = await api.post(`${this.moduleBase}/orders`, {
        partyId: payload.partyId ?? payload.party_id,
        quoteId: payload.quoteId ?? payload.quote_id ?? undefined,
        deliveryDate: payload.deliveryDate ?? payload.delivery_date,
        notes: payload.notes || '',
        mesWorkRefs: payload.mesWorkRefs ?? payload.mes_work_refs ?? [],
        items: payload.items ?? payload.line_items ?? []
      })
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo crear el pedido')
    }
  }

  async updateOrder(id: string, data: UpdateOrderRequest): Promise<Order> {
    const payload: any = data as any
    try {
      const response = await api.put(`${this.moduleBase}/orders/${id}`, {
        partyId: payload.partyId ?? payload.party_id,
        deliveryDate: payload.deliveryDate ?? payload.delivery_date,
        notes: payload.notes || '',
        mesWorkRefs: payload.mesWorkRefs ?? payload.mes_work_refs ?? []
      })
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar el pedido')
    }
  }

  async changeOrderStatus(id: string, status: string): Promise<any> {
    try {
      const response = await api.patch(`${this.moduleBase}/orders/${id}/status`, { newStatus: normalizeSalesStatus(status) })
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo cambiar el estado del pedido')
    }
  }

  async confirmOrder(id: string): Promise<any> {
    return this.changeOrderStatus(id, 'IN_PREPARATION')
  }

  async cancelOrder(id: string): Promise<any> {
    return this.changeOrderStatus(id, 'CANCELLED')
  }

  async reactivateOrder(id: string): Promise<any> {
    return this.changeOrderStatus(id, 'PENDING')
  }

  async addOrderLineItem(orderId: string, item: any): Promise<any> {
    try {
      const response = await api.post(`${this.moduleBase}/orders/${orderId}/line-items`, { item })
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo añadir la línea')
    }
  }

  async updateOrderLineItem(orderId: string, itemId: string, item: any): Promise<any> {
    try {
      const response = await api.put(`${this.moduleBase}/orders/${orderId}/line-items/${itemId}`, item)
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar la línea')
    }
  }

  async removeOrderLineItem(orderId: string, itemId: string): Promise<any> {
    try {
      const response = await api.delete(`${this.moduleBase}/orders/${orderId}/line-items/${itemId}`)
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo eliminar la línea')
    }
  }

  // --- Delivery Notes ---
  async getDeliveryNote(id: string): Promise<DeliveryNote> {
    try {
      const response = await api.get(`${this.moduleBase}/delivery-notes/${id}`)
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo obtener el albarán')
    }
  }

  async listDeliveryNotes(filters: ListDeliveryNotesFilters = {}): Promise<{ data: DeliveryNote[], total: number }> {
    const params: any = {}
    if (filters.orderId) params.salesOrderId = filters.orderId
    if (filters.searchText) params.search = filters.searchText
    if (filters.status) params.status = normalizeSalesStatus(filters.status)
    
    try {
      const response = await api.get(`${this.moduleBase}/delivery-notes`, { params })
      const res = response.data
      const rawDN = Array.isArray(res) ? res : (res.data || [])
      return { data: rawDN.map(normalizeEntity), total: resolveTotal(res, rawDN) }
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar los albaranes')
    }
  }

  async createDeliveryNote(data: CreateDeliveryNoteRequest): Promise<DeliveryNote> {
    try {
      const response = await api.post(`${this.moduleBase}/delivery-notes`, data)
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo crear el albarán')
    }
  }

  async changeDeliveryNoteStatus(id: string, status: string): Promise<DeliveryNote> {
    try {
      const response = await api.patch(`${this.moduleBase}/delivery-notes/${id}/status`, { newStatus: normalizeSalesStatus(status) })
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo cambiar el estado del albarán')
    }
  }

  async deleteDeliveryNote(id: string): Promise<void> {
    try {
      await api.delete(`${this.moduleBase}/delivery-notes/${id}`)
    } catch (e) {
      await this.handleError(e, 'No se pudo eliminar el albarán')
    }
  }

  // --- Invoices ---
  async getInvoice(id: string): Promise<Invoice> {
    try {
      const response = await api.get(`${this.moduleBase}/invoices/${id}`)
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo obtener la factura')
    }
  }

  async listInvoices(filters: ListInvoicesFilters = {}): Promise<{ data: Invoice[], total: number }> {
    const params: any = {}
    if (filters.orderId) params.orderId = filters.orderId
    if ((filters as any).partyId) params.partyId = (filters as any).partyId
    if (filters.deliveryNoteId) params.deliveryNoteId = filters.deliveryNoteId
    if (filters.searchText) params.search = filters.searchText
    if (filters.status) params.status = normalizeSalesStatus(filters.status)
    if (filters.type) params.type = filters.type
    
    try {
      const response = await api.get(`${this.moduleBase}/invoices`, { params })
      const res = response.data
      const rawInv = Array.isArray(res) ? res : (res.data || [])
      return { data: rawInv.map(normalizeEntity), total: resolveTotal(res, rawInv) }
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar las facturas')
    }
  }

  async createInvoice(data: CreateInvoiceRequest): Promise<Invoice> {
    try {
      const response = await api.post(`${this.moduleBase}/invoices`, data)
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo crear la factura')
    }
  }

  async createSimplifiedInvoice(data: CreateSimplifiedInvoiceRequest): Promise<Invoice> {
    try {
      const response = await api.post(`${this.moduleBase}/invoices/simplified`, data)
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo crear la factura simplificada')
    }
  }

  async updateInvoice(id: string, data: any): Promise<Invoice> {
    try {
      const response = await api.put(`${this.moduleBase}/invoices/${id}`, data)
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar la factura')
    }
  }

  async changeInvoiceStatus(id: string, status: string): Promise<Invoice> {
    try {
      const response = await api.patch(`${this.moduleBase}/invoices/${id}/status`, { newStatus: normalizeSalesStatus(status) })
      return normalizeEntity(response.data)
    } catch (e) {
      await this.handleError(e, 'No se pudo cambiar el estado de la factura')
    }
  }

  // --- Utils ---
  formatMoney(amount: any): string {
    if (amount === undefined || amount === null) return '—'
    const val = typeof amount === 'object' ? (amount.amount ?? 0) : (amount ?? 0)
    const currency = typeof amount === 'object' ? (amount.currency ?? 'EUR') : 'EUR'
    return new Intl.NumberFormat('es-ES', { style: 'currency', currency }).format(val)
  }

  formatUnitPrice(amount: any): string {
    if (!amount) return '—'
    const val = typeof amount === 'object' ? amount.amount : amount
    return new Intl.NumberFormat('es-ES', { style: 'decimal', minimumFractionDigits: 2, maximumFractionDigits: 4 }).format(val) + ' €'
  }

  async previewOrderCalculation(partyId: string, items: any[]): Promise<any> {
    try {
      const response = await api.post(`${this.moduleBase}/orders/preview`, { partyId, items })
      return response.data
    } catch (e) {
      return null
    }
  }

  getStatusClass(status: string): string {
    const statusMap: Record<string, string> = {
      'DRAFT': 'secondary',
      'ISSUED': 'info',
      'ACCEPTED': 'success',
      'APPROVED': 'success',
      'REJECTED': 'danger',
      'CONVERTED': 'success',
      'CONVERTED_TO_ORDER': 'success',
      'EXPIRED': 'secondary',
      'PENDING': 'warning',
      'IN_PREPARATION': 'info',
      'READY_FOR_PRODUCTION': 'primary',
      'PROCESSING': 'primary',
      'READY': 'success',
      'DELIVERED': 'success',
      'CANCELLED': 'danger',
      'PAID': 'success',
      'OVERDUE': 'danger',
      'VOID': 'secondary'
    }
    return statusMap[status] || 'secondary'
  }

  getStatusLabel(status: string): string {
    const labels: Record<string, string> = {
      'DRAFT': 'Borrador',
      'ISSUED': 'Emitido',
      'ACCEPTED': 'Aprobado',
      'APPROVED': 'Aprobado',
      'REJECTED': 'Rechazado',
      'CONVERTED': 'Aprobado',
      'CONVERTED_TO_ORDER': 'Aprobado (en Pedido)',
      'EXPIRED': 'Expirado',
      'PENDING': 'Pendiente',
      'IN_PREPARATION': 'En Preparación',
      'READY_FOR_PRODUCTION': 'Lanzado a Producción',
      'PROCESSING': 'En Taller',
      'READY': 'Listo',
      'DELIVERED': 'Entregado',
      'CANCELLED': 'Anulado',
      'PARTIALLY_DELIVERED': 'Entregado Parcial',
      'PAID': 'Pagada',
      'OVERDUE': 'Vencida',
      'VOID': 'Anulada',
      'PARTIALLY_INVOICED': 'Facturado Parcial',
      'INVOICED': 'Facturado'
    }
    return labels[status] || status
  }

  getQuoteStatusLabel(status: string): string { return this.getStatusLabel(status); }
  getQuoteStatusClass(status: string): string { return this.getStatusClass(status); }

  formatDateForAPI(date: Date): string {
    return date.toISOString()
  }
}

export default new SalesApi()

