import { test, expect, Page } from '@playwright/test'

const BASE_URL = 'http://localhost:5173'

const createMockJwt = (payload: Record<string, unknown>) => {
  const header = Buffer.from(JSON.stringify({ alg: 'HS256', typ: 'JWT' })).toString('base64')
  const body = Buffer.from(JSON.stringify(payload)).toString('base64')
  return `${header}.${body}.signature`
}

test.describe('Sales Stabilization Phase 1: Quotes and Invoices', () => {
  let page: Page

  test.beforeEach(async ({ browser }) => {
    const context = await browser.newContext()
    page = await context.newPage()

    // Setup Auth State
    const token = createMockJwt({
      sub: 'admin-123',
      email: 'admin@tramatex.com',
      role: 'admin',
      exp: Math.floor(Date.now() / 1000) + 3600
    })

    await page.addInitScript((tokenValue) => {
      window.localStorage.setItem('tramatex_auth_token', tokenValue)
      window.localStorage.setItem('tramatex_user', JSON.stringify({
        id: 'admin-123',
        email: 'admin@tramatex.com',
        role: 'admin'
      }))
    }, token)

    // Mock API responses for common data
    await page.route('**/api/parties/*', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ id: 'party-123', name: 'Cliente de Prueba E2E', organization_profile: { name: 'Cliente de Prueba E2E' } })
      })
    })

    await page.route('**/api/mes/work-types', async route => {
      await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify([]) })
    })

    await page.route('**/api/mes/positions', async route => {
      await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify([]) })
    })
  })

  test('Verify "Convertir a Pedido" button appears for ISSUED quotes', async () => {
    const quoteId = 'quote-issued-123'
    
    // Mock ISSUED quote
    await page.route(`**/api/sales/quotes/${quoteId}`, async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: quoteId,
          quoteNumber: 'PRE-2026-001',
          status: 'ISSUED', // Important for the test
          partyId: 'party-123',
          quoteDate: '2026-05-10T10:00:00Z',
          expirationDate: '2026-06-10T10:00:00Z',
          subtotal: { amount: 1000, currency: 'EUR' },
          taxAmount: { amount: 210, currency: 'EUR' },
          total: { amount: 1210, currency: 'EUR' },
          lineItems: [],
          mesWorkRefs: []
        })
      })
    })

    await page.goto(`${BASE_URL}/sales/quotes/${quoteId}`)

    // Wait for the page to load
    await expect(page.locator('h1.page-title')).toContainText('PRE-2026-001')
    
    // Check for the conversion button
    const convertButton = page.locator('button:has-text("Convertir a Pedido")')
    await expect(convertButton).toBeVisible()
    
    // Optional: Click and verify modal
    await convertButton.click()
    await expect(page.locator('text=¿Está seguro de que desea convertir este presupuesto')).toBeVisible()
  })

  test('Verify "Convertir a Pedido" button appears for EMITIDO (masculine) quotes', async () => {
    const quoteId = 'quote-emitido-123'
    
    // Mock quote with masculine Spanish status
    await page.route(`**/api/sales/quotes/${quoteId}`, async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: quoteId,
          quoteNumber: 'PRE-2026-002',
          status: 'EMITIDO', 
          partyId: 'party-123',
          quoteDate: '2026-05-10T10:00:00Z',
          expirationDate: '2026-06-10T10:00:00Z',
          subtotal: { amount: 1000, currency: 'EUR' },
          taxAmount: { amount: 210, currency: 'EUR' },
          total: { amount: 1210, currency: 'EUR' },
          lineItems: [],
          mesWorkRefs: []
        })
      })
    })

    await page.goto(`${BASE_URL}/sales/quotes/${quoteId}`)

    // Check for the conversion button - should be visible due to normalization
    const convertButton = page.locator('button:has-text("Convertir a Pedido")')
    await expect(convertButton).toBeVisible()
  })

  test('Verify Invoice payment transition flow', async () => {
    const invoiceId = 'invoice-pay-123'
    
    await page.route(`**/api/sales/invoices/${invoiceId}`, async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: invoiceId,
          invoiceNumber: 'FAC-PAY-001',
          status: 'ISSUED',
          partyId: 'party-123',
          total: { amount: 100, currency: 'EUR' },
          lineItems: []
        })
      })
    })

    // Mock the status update API call
    let statusUpdated = false
    await page.route(`**/api/sales/invoices/${invoiceId}/status`, async route => {
      const postData = route.request().postDataJSON()
      if (postData.newStatus === 'PAID') {
        statusUpdated = true
      }
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ id: invoiceId, status: 'PAID' })
      })
    })

    await page.goto(`${BASE_URL}/sales/invoices/${invoiceId}`)

    const payButton = page.locator('button:has-text("Registrar Cobro")')
    await expect(payButton).toBeVisible()
    
    await payButton.click()
    
    // Confirm in the prompt
    await expect(page.locator('text=¿Desea marcar esta factura como PAGADA?')).toBeVisible()
    await page.locator('button.btn-success:has-text("Confirmar Cobro")').click()
    
    expect(statusUpdated).toBe(true)
  })

  test('Verify InvoiceDetail.vue loads correctly with camelCase fields', async () => {
    const invoiceId = 'invoice-123'
    
    // Mock Invoice with snake_case fields (backend standard) to test frontend normalization
    await page.route(`**/api/sales/invoices/${invoiceId}`, async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: invoiceId,
          invoice_number: 'FAC-2026-001',
          status: 'ISSUED',
          party_id: 'party-123',
          invoice_date: '2026-05-12T08:00:00Z',
          due_date: '2026-06-12T08:00:00Z',
          subtotal_amount: { amount: 500, currency: 'EUR' },
          tax_amount: { amount: 105, currency: 'EUR' },
          total_amount: { amount: 605, currency: 'EUR' },
          line_items: [
            {
              product_variant_id: 'var-1',
              variant_sku: 'SKU-001',
              product_name: 'Producto de Prueba',
              quantity: 2,
              unit_price: { amount: 250, currency: 'EUR' }
            }
          ]
        })
      })
    })

    await page.goto(`${BASE_URL}/sales/invoices/${invoiceId}`)

    // Verify header title (uses invoiceNumber from normalization)
    await expect(page.locator('h1.page-title')).toContainText('FAC-2026-001')
    
    // Verify party name is loaded
    await expect(page.locator('text=Cliente de Prueba E2E').first()).toBeVisible()
    
    // Verify economic summary (uses total from normalization)
    await expect(page.locator('.amount')).toContainText('605,00')
    
    // Verify table lines (uses lineItems from normalization)
    await expect(page.locator('text=Producto de Prueba').first()).toBeVisible()
    await expect(page.locator('text=SKU-001').first()).toBeVisible()
  })
})
