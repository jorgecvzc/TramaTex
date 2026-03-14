import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import salesApi from '../../services/salesApi'

// Mock fetch globally
globalThis.fetch = vi.fn()

describe('SalesApi Service', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    localStorage.setItem('tramatex_auth_token', 'test-token')
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  // ============================================================================
  // QUOTES
  // ============================================================================

  describe('createQuote', () => {
    it('should create a new quote', async () => {
      const mockQuote = {
        id: 'quote-001',
        quoteNumber: 'Q-001',
        party_id: 'client-001',
        totalAmount: { amount: 1000, currency: 'EUR' },
        status: 'DRAFT',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockQuote,
      })

      const result = await salesApi.createQuote({
        party_id: 'client-001',
        valid_until: '2026-03-18',
        line_items: [],
      })

      expect(result.id).toBe('quote-001')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/quotes'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('client-001'),
        })
      )
    })

    it('should handle error when creating quote', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Invalid client ID' }),
      })

      await expect(
        salesApi.createQuote({
          party_id: 'invalid',
            valid_until: '2026-03-18',
          line_items: [],
        })
      ).rejects.toThrow('Invalid client ID')
    })
  })

  describe('getQuote', () => {
    it('should get quote by ID', async () => {
      const mockQuote = {
        id: 'quote-001',
        quoteNumber: 'Q-001',
        status: 'DRAFT',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockQuote,
      })

      const result = await salesApi.getQuote('quote-001')

      expect(result.id).toBe('quote-001')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/quotes/quote-001'),
        expect.objectContaining({ method: 'GET' })
      )
    })

    it('should handle error when quote not found', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Quote not found' }),
      })

      await expect(salesApi.getQuote('invalid')).rejects.toThrow('Quote not found')
    })
  })

  describe('listQuotes', () => {
    it('should list quotes successfully', async () => {
      const mockQuotes = [
        { id: 'quote-001', quote_number: 'Q-001', status: 'DRAFT' },
        { id: 'quote-002', quote_number: 'Q-002', status: 'ISSUED' },
      ]

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockQuotes,
      })

      const result = await salesApi.listQuotes()

      expect(result).toHaveLength(2)
      expect(result[0].quote_number).toBe('Q-001')
    })

    it('should apply filters correctly', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })

      await salesApi.listQuotes({
        searchText: 'acme',
        partyId: 'party-001',
        status: 'DRAFT',
        fromDate: '2026-01-01',
        toDate: '2026-12-31',
      })

      const fetchUrl = (globalThis.fetch as any).mock.calls[0][0]
      expect(fetchUrl).toContain('search=acme')
      expect(fetchUrl).toContain('partyId=party-001')
      expect(fetchUrl).toContain('status=DRAFT')
      expect(fetchUrl).toContain('fromDate=2026-01-01')
      expect(fetchUrl).toContain('toDate=2026-12-31')
    })

    it('should handle error when listing quotes', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({}),
      })

      await expect(salesApi.listQuotes()).rejects.toThrow()
    })
  })

  describe('updateQuote', () => {
    it('should update quote', async () => {
      const mockQuote = {
        id: 'quote-001',
        quoteNumber: 'Q-001',
        valid_until: '2026-04-01',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockQuote,
      })

      const result = await salesApi.updateQuote('quote-001', {
        valid_until: '2026-04-01',
      })

      expect(result.valid_until).toBe('2026-04-01')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/quotes/quote-001'),
        expect.objectContaining({
          method: 'PUT',
          body: expect.stringContaining('2026-04-01'),
        })
      )
    })

    it('should handle error when updating quote', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Quote is not in DRAFT status' }),
      })

      await expect(salesApi.updateQuote('quote-001', {})).rejects.toThrow()
    })
  })

  describe('changeQuoteStatus', () => {
    it('should change quote status', async () => {
      const mockQuote = {
        id: 'quote-001',
        quoteNumber: 'Q-001',
        status: 'ISSUED',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockQuote,
      })

      const result = await salesApi.changeQuoteStatus('quote-001', 'ISSUED')

      expect(result.status).toBe('ISSUED')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/quotes/quote-001/status'),
        expect.objectContaining({
          method: 'PATCH',
          body: JSON.stringify({ newStatus: 'ISSUED' }),
        })
      )
    })

    it('should handle error when changing quote status', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Invalid status transition' }),
      })

      await expect(salesApi.changeQuoteStatus('quote-001', 'ACCEPTED')).rejects.toThrow()
    })
  })

  describe('convertQuoteToOrder', () => {
    it('should convert quote to order', async () => {
      const mockOrder = {
        id: 'order-001',
        order_number: 'O-001',
        quote_id: 'quote-001',
        status: 'PENDING',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockOrder,
      })

      const result = await salesApi.convertQuoteToOrder('quote-001', '2026-03-01')

      expect(result.id).toBe('order-001')
      expect(result.quote_id).toBe('quote-001')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/quotes/quote-001/convert'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('2026-03-01'),
        })
      )
    })

    it('should handle error when converting quote', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Quote is not in ACCEPTED status' }),
      })

      await expect(salesApi.convertQuoteToOrder('quote-001', '2026-03-01')).rejects.toThrow()
    })
  })

  // ============================================================================
  // ORDERS
  // ============================================================================

  describe('createOrder', () => {
    it('should create a new order', async () => {
      const mockOrder = {
        id: 'order-001',
        order_number: 'O-001',
        party_id: 'client-001',
        totalAmount: { amount: 1500, currency: 'EUR' },
        status: 'PENDING',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockOrder,
      })

      const result = await salesApi.createOrder({
        party_id: 'client-001',
        delivery_date: '2026-03-01',
        line_items: [],
      })

      expect(result.id).toBe('order-001')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('client-001'),
        })
      )
    })

    it('should handle error when creating order', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Invalid delivery date' }),
      })

      await expect(
        salesApi.createOrder({
          party_id: 'client-001',
            delivery_date: '2026-01-01',
          line_items: [],
        })
      ).rejects.toThrow('Invalid delivery date')
    })
  })

  describe('getOrder', () => {
    it('should get order by ID', async () => {
      const mockOrder = {
        id: 'order-001',
        order_number: 'O-001',
        status: 'PENDING',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockOrder,
      })

      const result = await salesApi.getOrder('order-001')

      expect(result.id).toBe('order-001')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders/order-001'),
        expect.objectContaining({ method: 'GET' })
      )
    })

    it('should handle error when order not found', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Order not found' }),
      })

      await expect(salesApi.getOrder('invalid')).rejects.toThrow('Order not found')
    })
  })

  describe('listOrders', () => {
    it('should list orders successfully', async () => {
      const mockOrders = [
        { id: 'order-001', order_number: 'O-001', status: 'PENDING' },
        { id: 'order-002', order_number: 'O-002', status: 'CONFIRMED' },
      ]

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockOrders,
      })

      const result = await salesApi.listOrders()

      expect(result).toHaveLength(2)
      expect(result[0].order_number).toBe('O-001')
    })

    it('should apply filters correctly', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })

      await salesApi.listOrders({
        searchText: 'acme',
        partyId: 'party-001',
        status: 'CONFIRMED',
        fromDate: '2026-01-01',
        toDate: '2026-12-31',
      })

      const fetchUrl = (globalThis.fetch as any).mock.calls[0][0]
      expect(fetchUrl).toContain('search=acme')
      expect(fetchUrl).toContain('partyId=party-001')
      expect(fetchUrl).toContain('status=CONFIRMED')
    })

    it('should handle error when listing orders', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({}),
      })

      await expect(salesApi.listOrders()).rejects.toThrow()
    })
  })

  describe('updateOrder', () => {
    it('should update order', async () => {
      const mockOrder = {
        id: 'order-001',
        order_number: 'O-001',
        delivery_date: '2026-03-15',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockOrder,
      })

      const result = await salesApi.updateOrder('order-001', {
        delivery_date: '2026-03-15',
      })

      expect(result.delivery_date).toBe('2026-03-15')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders/order-001'),
        expect.objectContaining({
          method: 'PUT',
          body: expect.stringContaining('2026-03-15'),
        })
      )
    })

    it('should handle error when updating order', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Order is not editable' }),
      })

      await expect(salesApi.updateOrder('order-001', {})).rejects.toThrow()
    })
  })

  describe('changeOrderStatus', () => {
    it('should change order status', async () => {
      const mockOrder = {
        id: 'order-001',
        order_number: 'O-001',
        status: 'CONFIRMED',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockOrder,
      })

      const result = await salesApi.changeOrderStatus('order-001', 'CONFIRMED')

      expect(result.status).toBe('CONFIRMED')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders/order-001/status'),
        expect.objectContaining({
          method: 'PATCH',
          body: JSON.stringify({ newStatus: 'CONFIRMED' }),
        })
      )
    })

    it('should handle error when changing order status', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Invalid status transition' }),
      })

      await expect(salesApi.changeOrderStatus('order-001', 'CANCELLED')).rejects.toThrow()
    })
  })

  describe('addOrderLineItem', () => {
    it('should add line item to order', async () => {
      const mockOrder = {
        id: 'order-001',
        order_number: 'O-001',
        line_items: [{ id: 'line-001', quantity: 10 }],
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockOrder,
      })

      const result = await salesApi.addOrderLineItem('order-001', {
        productVariantId: 'variant-001',
        quantity: 10,
        unitPrice: { amount: 100, currency: 'EUR' },
      })

      expect(result.line_items).toHaveLength(1)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders/order-001/line-items'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('variant-001'),
        })
      )
    })

    it('should handle error when adding line item', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Product variant not found' }),
      })

      await expect(
        salesApi.addOrderLineItem('order-001', {
          productVariantId: 'invalid',
          quantity: 1,
          unitPrice: { amount: 100, currency: 'EUR' },
        })
      ).rejects.toThrow()
    })
  })

  describe('updateOrderLineItem', () => {
    it('should update order line item', async () => {
      const mockOrder = {
        id: 'order-001',
        order_number: 'O-001',
        line_items: [{ id: 'line-001', quantity: 20 }],
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockOrder,
      })

      const result = await salesApi.updateOrderLineItem('order-001', 'line-001', {
        quantity: 20,
      })

      expect(result.line_items[0].quantity).toBe(20)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders/order-001/line-items/line-001'),
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify({ quantity: 20 }),
        })
      )
    })

    it('should handle error when updating line item', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Line item not found' }),
      })

      await expect(salesApi.updateOrderLineItem('order-001', 'invalid', {})).rejects.toThrow()
    })
  })

  describe('removeOrderLineItem', () => {
    it('should remove line item from order', async () => {
      const mockOrder = {
        id: 'order-001',
        order_number: 'O-001',
        line_items: [],
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockOrder,
      })

      const result = await salesApi.removeOrderLineItem('order-001', 'line-001')

      expect(result.line_items).toHaveLength(0)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders/order-001/line-items/line-001'),
        expect.objectContaining({ method: 'DELETE' })
      )
    })

    it('should handle error when removing line item', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Cannot remove line item from confirmed order' }),
      })

      await expect(salesApi.removeOrderLineItem('order-001', 'line-001')).rejects.toThrow()
    })
  })

  // ============================================================================
  // DELIVERY NOTES
  // ============================================================================

  describe('createDeliveryNote', () => {
    it('should create a new delivery note', async () => {
      const mockNote = {
        id: 'note-001',
        delivery_note_number: 'DN-001',
        party_id: 'order-001',
        dispatch_date: '2026-03-01',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockNote,
      })

      const result = await salesApi.createDeliveryNote({
        order_id: 'order-001',
        dispatch_date: '2026-03-01',
        line_items: [],
      })

      expect(result.id).toBe('note-001')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/delivery-notes'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('order-001'),
        })
      )
    })

    it('should handle error when creating delivery note', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Order not found' }),
      })

      await expect(
        salesApi.createDeliveryNote({
          order_id: 'invalid',
          dispatch_date: '2026-03-01',
          line_items: [],
        })
      ).rejects.toThrow('Order not found')
    })
  })

  describe('getDeliveryNote', () => {
    it('should get delivery note by ID', async () => {
      const mockNote = {
        id: 'note-001',
        delivery_note_number: 'DN-001',
        party_id: 'order-001',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockNote,
      })

      const result = await salesApi.getDeliveryNote('note-001')

      expect(result.id).toBe('note-001')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/delivery-notes/note-001'),
        expect.objectContaining({ method: 'GET' })
      )
    })

    it('should handle error when delivery note not found', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Delivery note not found' }),
      })

      await expect(salesApi.getDeliveryNote('invalid')).rejects.toThrow()
    })
  })

  describe('listDeliveryNotes', () => {
    it('should list delivery notes successfully', async () => {
      const mockNotes = [
        { id: 'note-001', delivery_note_number: 'DN-001', party_id: 'order-001' },
        { id: 'note-002', delivery_note_number: 'DN-002', party_id: 'order-002' },
      ]

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockNotes,
      })

      const result = await salesApi.listDeliveryNotes()

      expect(result).toHaveLength(2)
      expect(result[0].delivery_note_number).toBe('DN-001')
    })

    it('should filter by orderId and search', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })

      await salesApi.listDeliveryNotes({ orderId: 'order-001', searchText: 'acme' })

      const fetchUrl = (globalThis.fetch as any).mock.calls[0][0]
      expect(fetchUrl).toContain('salesOrderId=order-001')
      expect(fetchUrl).toContain('search=acme')
    })

    it('should handle error when listing delivery notes', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({}),
      })

      await expect(salesApi.listDeliveryNotes()).rejects.toThrow()
    })
  })

  // ============================================================================
  // INVOICES
  // ============================================================================

  describe('createInvoice', () => {
    it('should create a new invoice', async () => {
      const mockInvoice = {
        id: 'invoice-001',
        invoice_number: 'INV-001',
        party_id: 'order-001',
        totalAmount: { amount: 1500, currency: 'EUR' },
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockInvoice,
      })

      const result = await salesApi.createInvoice({
        order_id: 'order-001',
        due_date: '2026-03-31',
        line_items: [],
      })

      expect(result.id).toBe('invoice-001')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/invoices'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('order-001'),
        })
      )
    })

    it('should handle error when creating invoice', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Order not delivered' }),
      })

      await expect(
        salesApi.createInvoice({
          order_id: 'order-001',
          due_date: '2025-01-01',
          line_items: [],
        })
      ).rejects.toThrow('Order not delivered')
    })
  })

  describe('createSimplifiedInvoice', () => {
    it('should create a simplified invoice', async () => {
      const mockInvoice = {
        id: 'invoice-002',
        invoice_number: 'INV-002',
        party_id: 'order-002',
        totalAmount: { amount: 500, currency: 'EUR' },
        isSimplified: true,
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockInvoice,
      })

      const result = await salesApi.createSimplifiedInvoice({
        party_id: 'order-002',
        line_items: [],
        payment_method: 'CASH' as const,
      })

      expect(result.id).toBe('invoice-002')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/invoices/simplified'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('order-002'),
        })
      )
    })

    it('should handle error when creating simplified invoice', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Amount exceeds simplified invoice limit' }),
      })

      await expect(
        salesApi.createSimplifiedInvoice({
          party_id: 'order-001',
          line_items: [],
          payment_method: 'CASH' as const,
        })
      ).rejects.toThrow()
    })
  })

  describe('getInvoice', () => {
    it('should get invoice by ID', async () => {
      const mockInvoice = {
        id: 'invoice-001',
        invoice_number: 'INV-001',
        party_id: 'order-001',
      }

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockInvoice,
      })

      const result = await salesApi.getInvoice('invoice-001')

      expect(result.id).toBe('invoice-001')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/invoices/invoice-001'),
        expect.objectContaining({ method: 'GET' })
      )
    })

    it('should handle error when invoice not found', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Invoice not found' }),
      })

      await expect(salesApi.getInvoice('invalid')).rejects.toThrow('Invoice not found')
    })
  })

  describe('listInvoices', () => {
    it('should list invoices successfully', async () => {
      const mockInvoices = [
        { id: 'invoice-001', invoice_number: 'INV-001', party_id: 'order-001' },
        { id: 'invoice-002', invoice_number: 'INV-002', party_id: 'order-002' },
      ]

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockInvoices,
      })

      const result = await salesApi.listInvoices()

      expect(result).toHaveLength(2)
      expect(result[0].invoice_number).toBe('INV-001')
    })

    it('should filter by orderId and search', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })

      await salesApi.listInvoices({ orderId: 'order-001', searchText: 'acme' })

      const fetchUrl = (globalThis.fetch as any).mock.calls[0][0]
      expect(fetchUrl).toContain('orderId=order-001')
      expect(fetchUrl).toContain('search=acme')
    })

    it('should filter by partyId', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })

      await salesApi.listInvoices({ partyId: 'party-001' })

      const fetchUrl = (globalThis.fetch as any).mock.calls[0][0]
      expect(fetchUrl).toContain('partyId=party-001')
    })

    it('should handle error when listing invoices', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        json: async () => ({}),
      })

      await expect(salesApi.listInvoices()).rejects.toThrow()
    })
  })

  // ============================================================================
  // AUTHENTICATION & ERROR HANDLING
  // ============================================================================

  describe('Authentication and headers', () => {
    it('should include auth token in headers', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })

      await salesApi.listQuotes()

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: 'Bearer test-token',
          }),
        })
      )
    })

    it('should work without auth token', async () => {
      localStorage.removeItem('tramatex_auth_token')

      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })

      await salesApi.listQuotes()

      const headers = (globalThis.fetch as any).mock.calls[0][1].headers
      expect(headers.Authorization).toBeUndefined()
    })

    it('should handle network errors gracefully', async () => {
      ;(globalThis.fetch as any).mockRejectedValueOnce(new Error('Network error'))

      await expect(salesApi.listQuotes()).rejects.toThrow(/No se pudo conectar/)
    })
  })
})





