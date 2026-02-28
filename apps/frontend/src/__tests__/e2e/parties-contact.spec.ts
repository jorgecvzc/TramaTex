import { test, expect, type Page } from '@playwright/test'

const BASE_URL = 'http://localhost:5173'

type PartyPayload = {
  id: string
  roles: string[]
  status: 'ACTIVE' | 'INACTIVE'
  organization_profile: {
    name: string
    tax_id: string
    tax_id_type: 'NIF' | 'CIF' | 'VAT'
    website: string | null
  } | null
  person_profile: {
    first_name: string
    last_name: string
  } | null
  created_at: string
  modified_at: string
}

const clientParty: PartyPayload = {
  id: 'party-001',
  roles: ['CLIENT'],
  status: 'ACTIVE',
  organization_profile: {
    name: 'Acme Textil',
    tax_id: 'B12345678',
    tax_id_type: 'CIF',
    website: null,
  },
  person_profile: null,
  created_at: '2026-02-10T08:00:00Z',
  modified_at: '2026-02-10T08:00:00Z',
}

const employeeParty: PartyPayload = {
  id: 'party-001',
  roles: ['EMPLOYEE'],
  status: 'ACTIVE',
  organization_profile: {
    name: 'Acme Textil',
    tax_id: 'B12345678',
    tax_id_type: 'CIF',
    website: null,
  },
  person_profile: null,
  created_at: '2026-02-10T08:00:00Z',
  modified_at: '2026-02-23T09:00:00Z',
}

async function mockAuthState(page: Page) {
  await page.addInitScript(() => {
    localStorage.setItem('tramatex_auth_token', 'test-token')
    localStorage.setItem('tramatex_refresh_token', 'refresh-token')
    localStorage.setItem('tramatex_user', JSON.stringify({ id: 'user-123', email: 'test@example.com' }))
  })
}

test.describe('Parties CONTACT compatibility flow', () => {
  test.beforeEach(async ({ page }) => {
    await mockAuthState(page)

    let listCallCount = 0
    let partyState: PartyPayload = { ...clientParty }

    await page.route('**/api/parties**', async route => {
      const request = route.request()
      const url = new URL(request.url())

      if (request.method() === 'GET' && url.pathname.endsWith('/api/parties')) {
        const roleFilter = url.searchParams.get('role')

        if (roleFilter === 'CONTACT') {
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({
              data: [employeeParty],
              total: 1,
              page: 1,
              limit: 10,
            }),
          })
          return
        }

        listCallCount += 1
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: [listCallCount > 1 ? employeeParty : clientParty],
            total: 1,
            page: 1,
            limit: 10,
          }),
        })
        return
      }

      if (request.method() === 'GET' && url.pathname.endsWith('/api/parties/party-001')) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify(partyState),
        })
        return
      }

      await route.fallback()
    })

    await page.route('**/api/parties/party-001', async route => {
      if (route.request().method() !== 'PUT') {
        await route.fallback()
        return
      }

      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(partyState),
      })
    })

    let rolesPostCount = 0

    await page.route('**/api/parties/party-001/roles', async route => {
      if (route.request().method() !== 'POST') {
        await route.fallback()
        return
      }

      rolesPostCount += 1

      if (rolesPostCount === 1) {
        await route.fulfill({
          status: 400,
          contentType: 'application/json',
          body: JSON.stringify({ message: 'invalid party role: CONTACT' }),
        })
        return
      }

      partyState = { ...employeeParty }

      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({}),
      })
    })

    await page.route('**/api/parties/party-001/roles/CLIENT', async route => {
      if (route.request().method() !== 'DELETE') {
        await route.fallback()
        return
      }

      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({}),
      })
    })

    await page.route('**/api/parties/*/relationships**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: [] }),
      })
    })

    await page.route('**/api/parties/*/contact-details**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: [] }),
      })
    })

    await page.route('**/api/parties/*/addresses**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: [], total: 0 }),
      })
    })
  })

  test('updates to CONTACT and filters Contactos in list', async ({ page }) => {
    await page.goto(`${BASE_URL}/parties/party-001`)

    await expect(page.getByRole('heading', { name: 'Detalle de entidad' })).toBeVisible()

    await page.getByRole('button', { name: /Editar entidad/i }).click()
    await page.locator('#editRole').selectOption('CONTACT')

    await Promise.all([
      page.waitForResponse(response =>
        response.url().includes('/api/parties/party-001/roles') && response.request().method() === 'POST'
      ),
      page.getByRole('button', { name: /Guardar cambios/i }).click(),
    ])

    await expect(page.locator('.badge.role-contact')).toContainText('Contacto')

    await page.goto(`${BASE_URL}/parties`)

    await expect(page.getByRole('heading', { name: 'Gestión de entidades' })).toBeVisible()

    await page.locator('label:has-text("Filtrar por rol") + select').selectOption('CONTACT')

    await expect(page.locator('tbody tr .role-pill')).toContainText('Contacto')
    await expect(page.locator('tbody tr .party-link')).toContainText('Acme Textil')
  })
})
