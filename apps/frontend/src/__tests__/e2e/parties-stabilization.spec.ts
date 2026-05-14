import { test, expect, Page } from '@playwright/test'

const BASE_URL = 'http://localhost:5173'

const createMockJwt = (payload: Record<string, unknown>) => {
  const header = Buffer.from(JSON.stringify({ alg: 'HS256', typ: 'JWT' })).toString('base64')
  const body = Buffer.from(JSON.stringify(payload)).toString('base64')
  return `${header}.${body}.signature`
}

test.describe('Parties Stabilization Phase 2: Header Actions', () => {
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

    // Mock API for Parties
    await page.route('**/api/parties/**', async route => {
      const url = route.request().url()
      if (route.request().method() === 'GET') {
        if (url.includes('/addresses')) {
          await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify([]) })
        } else if (url.includes('/contacts')) {
          await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify([]) })
        } else {
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({ 
              id: 'party-123', 
              name: 'Entidad Existente', 
              type: 'ORGANIZATION',
              status: 'ACTIVE',
              role: 'CLIENT',
              organization_profile: { name: 'Entidad Existente' }
            })
          })
        }
      } else {
        await route.continue()
      }
    })
    
    await page.route('**/api/parties', async route => {
      await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ data: [], total: 0 }) })
    })
  })

  test('Verify "Guardar" button is in the header on New Party page', async () => {
    await page.goto(`${BASE_URL}/parties/new`)

    // Wait for header to be visible
    const header = page.locator('.base-page-header')
    await expect(header).toBeVisible()

    // Check for "Guardar" button in the header actions area
    const saveButton = header.locator('button:has-text("Crear Entidad")')
    await expect(saveButton).toBeVisible()
  })

  test('Verify "Guardar" button is in the header on Edit Party page', async () => {
    await page.goto(`${BASE_URL}/parties/party-123`)
    
    // Enter Edit Mode
    const editButton = page.locator('button:has-text("Editar Datos")')
    await editButton.click()

    // Check for "Guardar" button in the header
    const saveButton = page.locator('.base-page-header button:has-text("Guardar Cambios")')
    await expect(saveButton).toBeVisible()
  })
})
