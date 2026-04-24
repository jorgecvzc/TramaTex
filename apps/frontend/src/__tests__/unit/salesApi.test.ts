import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import salesApi from '../../services/salesApi'
import { api } from '../../services/api'

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

      vi.mocked(api.post).mockResolvedValueOnce({
        data: mockQuote,
        status: 201,
      })

      const result = await salesApi.createQuote({
        party_id: 'client-001',
        valid_until: '2026-03-18',
        line_items: [],
      })

      expect(result.id).toBe('quote-001')
      expect(api.post).toHaveBeenCalledWith(
        expect.stringContaining('/sales/quotes'),
        expect.objectContaining({
          party_id: 'client-001',
          valid_until: '2026-03-18'
        })
      )
    })

    it('should handle error when creating quote', async () => {
      vi.mocked(api.post).mockRejectedValueOnce({
        response: {
          data: { error: 'Invalid client ID' },
          status: 400
        }
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

      vi.mocked(api.get).mockResolvedValueOnce({
        data: mockQuote,
        status: 200,
      })

      const result = await salesApi.getQuote('quote-001')

      expect(result.id).toBe('quote-001')
      expect(api.get).toHaveBeenCalledWith(
        expect.stringContaining('/sales/quotes/quote-001')
      )
    })

    it('should handle error when quote not found', async () => {
      vi.mocked(api.get).mockRejectedValueOnce({
        response: {
          data: { error: 'Quote not found' },
          status: 404
        }
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

      vi.mocked(api.get).mockResolvedValueOnce({
        data: mockQuotes,
        status: 200,
      })

      const result = await salesApi.listQuotes()

      expect(result.data).toHaveLength(2)
      expect(result.data[0].quote_number).toBe('Q-001')
    })

    it('should apply filters correctly', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: [],
        status: 200,
      })

      await salesApi.listQuotes({
        searchText: 'acme',
        partyId: 'party-001',
        status: 'DRAFT',
        fromDate: '2026-01-01',
        toDate: '2026-12-31',
      })

      expect(api.get).toHaveBeenCalledWith(
        expect.stringContaining('/sales/quotes'),
        expect.objectContaining({
          params: expect.objectContaining({
            search: 'acme',
            partyId: 'party-001',
            status: 'DRAFT'
          })
        })
      )
    })

    it('should handle error when listing quotes', async () => {
      vi.mocked(api.get).mockRejectedValueOnce(new Error('API Error'))

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

      vi.mocked(api.put).mockResolvedValueOnce({
        data: mockQuote,
        status: 200,
      })

      const result = await salesApi.updateQuote('quote-001', {
        valid_until: '2026-04-01',
      })

      expect(result.valid_until).toBe('2026-04-01')
      expect(api.put).toHaveBeenCalledWith(
        expect.stringContaining('/sales/quotes/quote-001'),
        expect.objectContaining({
          valid_until: '2026-04-01',
        })
      )
    })

    it('should handle error when updating quote', async () => {
      vi.mocked(api.put).mockRejectedValueOnce({
        response: {
          data: { error: 'Quote is not in DRAFT status' },
          status: 400
        }
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

      vi.mocked(api.patch).mockResolvedValueOnce({
        data: mockQuote,
        status: 200,
      })

      const result = await salesApi.changeQuoteStatus('quote-001', 'ISSUED')

      expect(result.status).toBe('ISSUED')
      expect(api.patch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/quotes/quote-001/status'),
        expect.objectContaining({ newStatus: 'ISSUED' })
      )
    })

    it('should handle error when changing quote status', async () => {
      vi.mocked(api.patch).mockRejectedValueOnce({
        response: {
          data: { error: 'Invalid status transition' },
          status: 400
        }
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

      vi.mocked(api.post).mockResolvedValueOnce({
        data: mockOrder,
        status: 201,
      })

      const result = await salesApi.convertQuoteToOrder('quote-001', '2026-03-01')

      expect(result.id).toBe('order-001')
      expect(result.quote_id).toBe('quote-001')
      expect(api.post).toHaveBeenCalledWith(
        expect.stringContaining('/sales/quotes/quote-001/convert'),
        expect.objectContaining({
          deliveryDate: expect.stringContaining('2026-03-01'),
        })
      )
    })

    it('should handle error when converting quote', async () => {
      vi.mocked(api.post).mockRejectedValueOnce({
        response: {
          data: { error: 'Quote is not in ACCEPTED status' },
          status: 400
        }
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

      vi.mocked(api.post).mockResolvedValueOnce({
        data: mockOrder,
        status: 201,
      })

      const result = await salesApi.createOrder({
        party_id: 'client-001',
        delivery_date: '2026-03-01',
        line_items: [],
      })

      expect(result.id).toBe('order-001')
      expect(api.post).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders'),
        expect.objectContaining({
          partyId: 'client-001',
        })
      )
    })

    it('should handle error when creating order', async () => {
      vi.mocked(api.post).mockRejectedValueOnce({
        response: {
          data: { error: 'Invalid delivery date' },
          status: 400
        }
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

      vi.mocked(api.get).mockResolvedValueOnce({
        data: mockOrder,
        status: 200,
      })

      const result = await salesApi.getOrder('order-001')

      expect(result.id).toBe('order-001')
      expect(api.get).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders/order-001')
      )
    })

    it('should handle error when order not found', async () => {
      vi.mocked(api.get).mockRejectedValueOnce({
        response: {
          data: { error: 'Order not found' },
          status: 404
        }
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

      vi.mocked(api.get).mockResolvedValueOnce({
        data: mockOrders,
        status: 200,
      })

      const result = await salesApi.listOrders()

      expect(result.data).toHaveLength(2)
      expect(result.data[0].order_number).toBe('O-001')
    })

    it('should apply filters correctly', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: [],
        status: 200,
      })

      await salesApi.listOrders({
        searchText: 'acme',
        partyId: 'party-001',
        status: 'CONFIRMED',
        fromDate: '2026-01-01',
        toDate: '2026-12-31',
      })

      expect(api.get).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders'),
        expect.objectContaining({
          params: expect.objectContaining({
            search: 'acme',
            partyId: 'party-001',
            status: 'CONFIRMED'
          })
        })
      )
    })

    it('should handle error when listing orders', async () => {
      vi.mocked(api.get).mockRejectedValueOnce(new Error('API Error'))

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

      vi.mocked(api.put).mockResolvedValueOnce({
        data: mockOrder,
        status: 200,
      })

      const result = await salesApi.updateOrder('order-001', {
        delivery_date: '2026-03-15',
      })

      expect(result.delivery_date).toBe('2026-03-15')
      expect(api.put).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders/order-001'),
        expect.objectContaining({
          deliveryDate: '2026-03-15',
        })
      )
    })

    it('should handle error when updating order', async () => {
      vi.mocked(api.put).mockRejectedValueOnce({
        response: {
          data: { error: 'Order is not editable' },
          status: 400
        }
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

      vi.mocked(api.patch).mockResolvedValueOnce({
        data: mockOrder,
        status: 200,
      })

      const result = await salesApi.changeOrderStatus('order-001', 'CONFIRMED')

      expect(result.status).toBe('CONFIRMED')
      expect(api.patch).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders/order-001/status'),
        expect.objectContaining({ newStatus: 'CONFIRMED' })
      )
    })

    it('should handle error when changing order status', async () => {
      vi.mocked(api.patch).mockRejectedValueOnce({
        response: {
          data: { error: 'Invalid status transition' },
          status: 400
        }
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

      vi.mocked(api.post).mockResolvedValueOnce({
        data: mockOrder,
        status: 201,
      })

      const result = await salesApi.addOrderLineItem('order-001', {
        productVariantId: 'variant-001',
        quantity: 10,
        unitPrice: { amount: 100, currency: 'EUR' },
      })

      expect(result.line_items).toHaveLength(1)
      expect(api.post).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders/order-001/line-items'),
        expect.objectContaining({
          item: expect.objectContaining({
             productVariantId: 'variant-001',
          })
        })
      )
    })

    it('should handle error when adding line item', async () => {
      vi.mocked(api.post).mockRejectedValueOnce(new Error('API Error'))

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

      vi.mocked(api.put).mockResolvedValueOnce({
        data: mockOrder,
        status: 200,
      })

      const result = await salesApi.updateOrderLineItem('order-001', 'line-001', {
        quantity: 20,
      })

      expect(result.line_items[0].quantity).toBe(20)
      expect(api.put).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders/order-001/line-items/line-001'),
        expect.objectContaining({ quantity: 20 })
      )
    })

    it('should handle error when updating line item', async () => {
      vi.mocked(api.put).mockRejectedValueOnce(new Error('API Error'))

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

      vi.mocked(api.delete).mockResolvedValueOnce({
        data: mockOrder,
        status: 200,
      })

      const result = await salesApi.removeOrderLineItem('order-001', 'line-001')

      expect(result.line_items).toHaveLength(0)
      expect(api.delete).toHaveBeenCalledWith(
        expect.stringContaining('/sales/orders/order-001/line-items/line-001')
      )
    })

    it('should handle error when removing line item', async () => {
      vi.mocked(api.delete).mockRejectedValueOnce({
        response: {
          data: { error: 'Cannot remove line item from confirmed order' },
          status: 400
        }
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

      vi.mocked(api.post).mockResolvedValueOnce({
        data: mockNote,
        status: 201,
      })

      const result = await salesApi.createDeliveryNote({
        order_id: 'order-001',
        dispatch_date: '2026-03-01',
        line_items: [],
      })

      expect(result.id).toBe('note-001')
      expect(api.post).toHaveBeenCalledWith(
        expect.stringContaining('/sales/delivery-notes'),
        expect.objectContaining({
          order_id: 'order-001',
        })
      )
    })

    it('should handle error when creating delivery note', async () => {
      vi.mocked(api.post).mockRejectedValueOnce({
        response: {
          data: { error: 'Order not found' },
          status: 404
        }
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

      vi.mocked(api.get).mockResolvedValueOnce({
        data: mockNote,
        status: 200,
      })

      const result = await salesApi.getDeliveryNote('note-001')

      expect(result.id).toBe('note-001')
      expect(api.get).toHaveBeenCalledWith(
        expect.stringContaining('/sales/delivery-notes/note-001')
      )
    })

    it('should handle error when delivery note not found', async () => {
      vi.mocked(api.get).mockRejectedValueOnce(new Error('API Error'))

      await expect(salesApi.getDeliveryNote('invalid')).rejects.toThrow()
    })
  })

  describe('listDeliveryNotes', () => {
    it('should list delivery notes successfully', async () => {
      const mockNotes = [
        { id: 'note-001', delivery_note_number: 'DN-001', party_id: 'order-001' },
        { id: 'note-002', delivery_note_number: 'DN-002', party_id: 'order-002' },
      ]

      vi.mocked(api.get).mockResolvedValueOnce({
        data: mockNotes,
        status: 200,
      })

      const result = await salesApi.listDeliveryNotes()

      expect(result.data).toHaveLength(2)
      expect(result.data[0].delivery_note_number).toBe('DN-001')
    })

    it('should filter by orderId and search', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: [],
        status: 200,
      })

      await salesApi.listDeliveryNotes({ orderId: 'order-001', searchText: 'acme' })

      expect(api.get).toHaveBeenCalledWith(
        expect.stringContaining('/sales/delivery-notes'),
        expect.objectContaining({
          params: expect.objectContaining({
            salesOrderId: 'order-001',
            search: 'acme'
          })
        })
      )
    })

    it('should handle error when listing delivery notes', async () => {
      vi.mocked(api.get).mockRejectedValueOnce(new Error('API Error'))

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

      vi.mocked(api.post).mockResolvedValueOnce({
        data: mockInvoice,
        status: 201,
      })

      const result = await salesApi.createInvoice({
        order_id: 'order-001',
        due_date: '2026-03-31',
        line_items: [],
      })

      expect(result.id).toBe('invoice-001')
      expect(api.post).toHaveBeenCalledWith(
        expect.stringContaining('/sales/invoices'),
        expect.objectContaining({
          order_id: 'order-001',
        })
      )
    })

    it('should handle error when creating invoice', async () => {
      vi.mocked(api.post).mockRejectedValueOnce({
        response: {
          data: { error: 'Order not delivered' },
          status: 400
        }
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

      vi.mocked(api.post).mockResolvedValueOnce({
        data: mockInvoice,
        status: 201,
      })

      const result = await salesApi.createSimplifiedInvoice({
        party_id: 'order-002',
        line_items: [],
        payment_method: 'CASH' as const,
      })

      expect(result.id).toBe('invoice-002')
      expect(api.post).toHaveBeenCalledWith(
        expect.stringContaining('/sales/invoices/simplified'),
        expect.objectContaining({
          party_id: 'order-002',
        })
      )
    })

    it('should handle error when creating simplified invoice', async () => {
      vi.mocked(api.post).mockRejectedValueOnce({
        response: {
          data: { error: 'Amount exceeds simplified invoice limit' },
          status: 400
        }
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

      vi.mocked(api.get).mockResolvedValueOnce({
        data: mockInvoice,
        status: 200,
      })

      const result = await salesApi.getInvoice('invoice-001')

      expect(result.id).toBe('invoice-001')
      expect(api.get).toHaveBeenCalledWith(
        expect.stringContaining('/sales/invoices/invoice-001')
      )
    })

    it('should handle error when invoice not found', async () => {
      vi.mocked(api.get).mockRejectedValueOnce({
        response: {
          data: { error: 'Invoice not found' },
          status: 404
        }
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

      vi.mocked(api.get).mockResolvedValueOnce({
        data: mockInvoices,
        status: 200,
      })

      const result = await salesApi.listInvoices()

      expect(result.data).toHaveLength(2)
      expect(result.data[0].invoice_number).toBe('INV-001')
    })

    it('should filter by orderId and search', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: [],
        status: 200,
      })

      await salesApi.listInvoices({ orderId: 'order-001', searchText: 'acme' })

      expect(api.get).toHaveBeenCalledWith(
        expect.stringContaining('/sales/invoices'),
        expect.objectContaining({
          params: expect.objectContaining({
            orderId: 'order-001',
            search: 'acme'
          })
        })
      )
    })

    it('should filter by partyId', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: [],
        status: 200,
      })

      await salesApi.listInvoices({ partyId: 'party-001' })

      expect(api.get).toHaveBeenCalledWith(
        expect.stringContaining('/sales/invoices'),
        expect.objectContaining({
          params: expect.objectContaining({
            partyId: 'party-001'
          })
        })
      )
    })

    it('should handle error when listing invoices', async () => {
      vi.mocked(api.get).mockRejectedValueOnce(new Error('API Error'))

      await expect(salesApi.listInvoices()).rejects.toThrow()
    })
  })

  // ============================================================================
  // AUTHENTICATION & ERROR HANDLING
  // ============================================================================

  describe('Authentication and headers', () => {
    it('should work with auth token', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: [],
        status: 200,
      })

      await salesApi.listQuotes()

      expect(api.get).toHaveBeenCalled()
    })

    it('should work without auth token', async () => {
      localStorage.removeItem('tramatex_auth_token')

      vi.mocked(api.get).mockResolvedValueOnce({
        data: [],
        status: 200,
      })

      await salesApi.listQuotes()

      expect(api.get).toHaveBeenCalled()
    })

    it('should handle network errors gracefully', async () => {
      vi.mocked(api.get).mockRejectedValueOnce(new Error('Network error'))

      await expect(salesApi.listQuotes()).rejects.toThrow()
    })
  })
})
