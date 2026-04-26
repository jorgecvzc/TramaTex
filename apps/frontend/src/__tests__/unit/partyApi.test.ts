import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { partyApi } from '../../services/partyApi'
import type { PartyUI, CreatePartyRequest, UpdatePartyRequest, Contact } from '../../types/party'

// Mock fetch globally
global.fetch = vi.fn()

describe('PartyApi Service', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    localStorage.setItem('tramatex_auth_token', 'test-token')
    localStorage.setItem('tramatex_user', JSON.stringify({ id: 'user-123' }))
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  describe('createParty', () => {
    it('should create a party successfully with organization profile', async () => {
      const mockParty: PartyUI = {
        id: 'party-001',
        name: 'Acme Corporation',
        role: 'CLIENT',
        status: 'ACTIVE',
        tax_id: 'B12345678',
        tax_id_type: 'CIF',
        website: 'https://acme.com',
        phone: null,
        email: null,
        default_discount_percentage: 0,
        notes: null,
        created_at: '2026-02-17T10:00:00Z',
        modified_at: '2026-02-17T10:00:00Z',
        has_organization: true,
        has_person: false,
        can_delete: true,
      }

      const mockBackendResponse = {
        id: 'party-001',
        roles: ['CLIENT'],
        status: 'ACTIVE',
        organization_profile: {
          name: 'Acme Corporation',
          tax_id: 'B12345678',
          tax_id_type: 'CIF',
          website: 'https://acme.com',
          notes: null,
        },
        person_profile: null,
        created_at: '2026-02-17T10:00:00Z',
        modified_at: '2026-02-17T10:00:00Z',
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockBackendResponse,
      })

      const createData: CreatePartyRequest = {
        name: 'Acme Corporation',
        role: 'CLIENT',
        taxId: 'B12345678',
        taxIdType: 'CIF',
        website: 'https://acme.com',
      }

      const result = await partyApi.createParty(createData)

      expect(result).toEqual(mockParty)
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/parties'),
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            'Authorization': 'Bearer test-token',
            'X-User-ID': 'user-123',
          }),
          body: expect.any(String),
        })
      )
    })

    it('should handle API error when creating party', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => ({ message: 'Tax ID already exists' }),
      })

      const createData: CreatePartyRequest = {
        name: 'Duplicate Corp',
        role: 'CLIENT',
        taxId: 'B12345678',
      }

      await expect(partyApi.createParty(createData)).rejects.toThrow('Tax ID already exists')
    })

    it('should handle network error when creating party', async () => {
      ;(global.fetch as any).mockRejectedValueOnce(new Error('Network error'))

      const createData: CreatePartyRequest = {
        name: 'Test Party',
        role: 'CLIENT',
      }

      await expect(partyApi.createParty(createData)).rejects.toThrow(/No se pudo conectar/)
    })
  })

  describe('getParty', () => {
    it('should fetch a party by ID successfully', async () => {
      const mockBackendParty = {
        id: 'party-002',
        roles: ['SUPPLIER'],
        status: 'ACTIVE',
        organization_profile: {
          name: 'Supplier Inc',
          tax_id: 'A87654321',
          tax_id_type: 'CIF',
          website: null,
        },
        person_profile: null,
        created_at: '2026-02-15T09:00:00Z',
        modified_at: '2026-02-15T09:00:00Z',
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockBackendParty,
      })

      const result = await partyApi.getParty('party-002')

      expect(result).toMatchObject({
        id: 'party-002',
        name: 'Supplier Inc',
        role: 'SUPPLIER',
        status: 'ACTIVE',
        tax_id: 'A87654321',
      })
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/parties/party-002'),
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Authorization': 'Bearer test-token',
          }),
        })
      )
    })

    it('should return null when party not found', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: async () => ({ message: 'Party not found' }),
      })

      await expect(partyApi.getParty('non-existent')).rejects.toThrow('Party not found')
    })
  })

  describe('listParties', () => {
    it('should list parties with filters', async () => {
      const mockBackendResponse = {
        data: [
          {
            id: 'party-001',
            roles: ['CLIENT'],
            status: 'ACTIVE',
            organization_profile: {
              name: 'Client One',
              tax_id: 'B11111111',
              tax_id_type: 'CIF',
              website: null,
            },
            person_profile: null,
            created_at: '2026-02-10T08:00:00Z',
            modified_at: '2026-02-10T08:00:00Z',
          },
          {
            id: 'party-002',
            roles: ['CLIENT'],
            status: 'ACTIVE',
            organization_profile: {
              name: 'Client Two',
              tax_id: 'B22222222',
              tax_id_type: 'CIF',
              website: null,
            },
            person_profile: null,
            created_at: '2026-02-11T08:00:00Z',
            modified_at: '2026-02-11T08:00:00Z',
          },
        ],
        total: 2,
        page: 1,
        limit: 10,
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockBackendResponse,
      })

      const filters = { role: 'CLIENT' as const, pageNumber: 1, pageSize: 10 }
      const result = await partyApi.listParties(filters)

      expect(result.data).toHaveLength(2)
      expect(result.total).toBe(2)
      expect(result.data[0].name).toBe('Client One')
      expect(result.data[1].name).toBe('Client Two')
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/parties?role=CLIENT'),
        expect.any(Object)
      )
    })

    it('should list parties without filters', async () => {
      const mockResponse = {
        data: [],
        total: 0,
        page: 1,
        limit: 10,
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await partyApi.listParties()

      expect(result.data).toHaveLength(0)
      expect(result.total).toBe(0)
    })

    it('should fallback CONTACT filter to EMPLOYEE when backend rejects CONTACT', async () => {
      const fallbackPayload = {
        data: [
          {
            id: 'party-contact-legacy',
            roles: ['EMPLOYEE'],
            status: 'ACTIVE',
            organization_profile: null,
            person_profile: {
              first_name: 'Pablo',
              last_name: 'Nadal',
            },
            created_at: '2026-02-11T08:00:00Z',
            modified_at: '2026-02-11T08:00:00Z',
          },
        ],
        total: 1,
        page: 1,
        limit: 10,
      }

      ;(global.fetch as any)
        .mockResolvedValueOnce({
          ok: false,
          status: 400,
          json: async () => ({ message: 'invalid role: CONTACT' }),
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => fallbackPayload,
        })

      const result = await partyApi.listParties({ role: 'CONTACT', pageNumber: 1, pageSize: 10 })

      expect(result.data).toHaveLength(1)
      expect(result.data[0].role).toBe('CONTACT')
      expect(global.fetch).toHaveBeenNthCalledWith(
        1,
        expect.stringContaining('/parties?role=CONTACT'),
        expect.any(Object)
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        2,
        expect.stringContaining('/parties?role=EMPLOYEE'),
        expect.any(Object)
      )
    })

    it('should map EMPLOYEE backend role as CONTACT in UI', async () => {
      const mockBackendResponse = {
        data: [
          {
            id: 'party-emp-001',
            roles: ['EMPLOYEE'],
            status: 'ACTIVE',
            organization_profile: null,
            person_profile: {
              first_name: 'Lucia',
              last_name: 'Martinez',
            },
            created_at: '2026-02-11T08:00:00Z',
            modified_at: '2026-02-11T08:00:00Z',
          },
        ],
        total: 1,
        page: 1,
        limit: 10,
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockBackendResponse,
      })

      const result = await partyApi.listParties()

      expect(result.data).toHaveLength(1)
      expect(result.data[0].role).toBe('CONTACT')
      expect(result.data[0].name).toBe('Lucia Martinez')
    })
  })

  describe('updateParty', () => {
    it('should update a party successfully', async () => {
      const mockCurrentParty = {
        id: 'party-001',
        roles: ['CLIENT', 'SUPPLIER'],
        status: 'ACTIVE',
        organization_profile: {
          name: 'Test Corp',
          tax_id: 'B12345678',
          tax_id_type: 'CIF',
          website: null,
        },
        person_profile: null,
        created_at: '2026-02-10T08:00:00Z',
        modified_at: '2026-02-17T11:00:00Z',
      }

      const mockUpdatedParty = {
        ...mockCurrentParty,
        organization_profile: {
          ...mockCurrentParty.organization_profile,
          name: 'Updated Corp',
          website: 'https://updated.com',
        },
      }

      // New flow: GET (fetch current roles) → no role sync needed (BOTH = CLIENT+SUPPLIER already) → PUT
      ;(global.fetch as any)
        .mockResolvedValueOnce({ ok: true, json: async () => mockCurrentParty }) // GET current
        .mockResolvedValueOnce({ ok: true, json: async () => mockUpdatedParty }) // PUT

      const updateData: UpdatePartyRequest = {
        name: 'Updated Corp',
        role: 'BOTH',
        website: 'https://updated.com',
      }

      const result = await partyApi.updateParty('party-001', updateData)

      expect(result).toMatchObject({
        id: 'party-001',
        name: 'Updated Corp',
        role: 'BOTH',
        website: 'https://updated.com',
      })
      expect(global.fetch).toHaveBeenCalledTimes(2)
      expect(global.fetch).toHaveBeenNthCalledWith(
        2,
        expect.stringContaining('/parties/party-001'),
        expect.objectContaining({
          method: 'PUT',
          body: expect.any(String),
        })
      )
    })

    it('should sync roles when updating from SUPPLIER to CLIENT', async () => {
      const currentPartyRaw = {
        id: 'party-001',
        roles: ['SUPPLIER'],
        status: 'ACTIVE',
        organization_profile: {
          name: 'Updated Corp',
          tax_id: 'B12345678',
          tax_id_type: 'CIF',
          website: 'https://updated.com',
        },
        person_profile: null,
        created_at: '2026-02-10T08:00:00Z',
        modified_at: '2026-02-17T11:00:00Z',
      }

      const updatedBackendParty = {
        ...currentPartyRaw,
        roles: ['CLIENT'],
      }

      // New flow: GET → POST add CLIENT → DELETE SUPPLIER → PUT
      ;(global.fetch as any)
        .mockResolvedValueOnce({ ok: true, json: async () => currentPartyRaw }) // GET current roles
        .mockResolvedValueOnce({ ok: true, json: async () => ({}) }) // POST add CLIENT
        .mockResolvedValueOnce({ ok: true, json: async () => ({}) }) // DELETE SUPPLIER
        .mockResolvedValueOnce({ ok: true, json: async () => updatedBackendParty }) // PUT

      const result = await partyApi.updateParty('party-001', {
        name: 'Updated Corp',
        role: 'CLIENT',
      })

      expect(result?.role).toBe('CLIENT')
      expect(global.fetch).toHaveBeenCalledTimes(4)
      expect(global.fetch).toHaveBeenNthCalledWith(
        1,
        expect.stringContaining('/parties/party-001'),
        expect.objectContaining({ method: 'GET' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        2,
        expect.stringContaining('/parties/party-001/roles'),
        expect.objectContaining({ method: 'POST' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        3,
        expect.stringContaining('/parties/party-001/roles/SUPPLIER'),
        expect.objectContaining({ method: 'DELETE' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        4,
        expect.stringContaining('/parties/party-001'),
        expect.objectContaining({ method: 'PUT' })
      )
    })

    it('should fallback to EMPLOYEE when CONTACT role is rejected', async () => {
      const updatedBackendParty = {
        id: 'party-contact-001',
        roles: ['CLIENT'],
        status: 'ACTIVE',
        organization_profile: {
          name: 'Contact Corp',
          tax_id: 'B33333333',
          tax_id_type: 'CIF',
          website: null,
        },
        person_profile: null,
        created_at: '2026-02-10T08:00:00Z',
        modified_at: '2026-02-17T11:00:00Z',
      }

      const finalPartyAfterSync = {
        ...updatedBackendParty,
        roles: ['EMPLOYEE'],
      }

      // New flow: GET current → POST CONTACT (fails) → POST EMPLOYEE (fallback) → DELETE CLIENT → PUT
      ;(global.fetch as any)
        .mockResolvedValueOnce({ ok: true, json: async () => updatedBackendParty }) // GET current roles
        .mockResolvedValueOnce({ ok: false, status: 400, json: async () => ({ message: 'invalid party role: CONTACT' }) }) // POST add CONTACT (fails)
        .mockResolvedValueOnce({ ok: true, json: async () => ({}) }) // POST add EMPLOYEE fallback
        .mockResolvedValueOnce({ ok: true, json: async () => ({}) }) // DELETE CLIENT
        .mockResolvedValueOnce({ ok: true, json: async () => finalPartyAfterSync }) // PUT

      const result = await partyApi.updateParty('party-contact-001', {
        name: 'Contact Corp',
        role: 'CONTACT',
      })

      expect(result?.role).toBe('CONTACT')
      expect(global.fetch).toHaveBeenCalledTimes(5)
      expect(global.fetch).toHaveBeenNthCalledWith(
        1,
        expect.stringContaining('/parties/party-contact-001'),
        expect.objectContaining({ method: 'GET' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        2,
        expect.stringContaining('/parties/party-contact-001/roles'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ role: 'CONTACT' }),
        })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        3,
        expect.stringContaining('/parties/party-contact-001/roles'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ role: 'EMPLOYEE' }),
        })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        4,
        expect.stringContaining('/parties/party-contact-001/roles/CLIENT'),
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })

  describe('changePartyStatus', () => {
    it('should change party status to INACTIVE', async () => {
      const mockParty = {
        id: 'party-001',
        roles: ['CLIENT'],
        status: 'INACTIVE',
        organization_profile: {
          name: 'Test Corp',
          tax_id: 'B12345678',
          tax_id_type: 'CIF',
          website: null,
        },
        person_profile: null,
        created_at: '2026-02-10T08:00:00Z',
        modified_at: '2026-02-17T11:30:00Z',
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockParty,
      })

      const result = await partyApi.changePartyStatus('party-001', 'INACTIVE')

      expect(result?.status).toBe('INACTIVE')
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/parties/party-001/status'),
        expect.objectContaining({
          method: 'PATCH',
          body: JSON.stringify({ status: 'INACTIVE' }),
        })
      )
    })
  })

  describe('Contact management', () => {
    // Note: Contact management functions make multiple nested fetch calls
    // which are complex to mock properly. These are integration-level tests.
    
    it('should add a contact to a party', async () => {
      // addContact makes 3 fetch calls: create person, create relationship, create contact-details
      // Mock all 3 responses
      const mockPersonParty: any = {
        id: 'person-001',
        roles: ['EMPLOYEE'],
        status: 'ACTIVE',
        person_profile: {
          first_name: 'John',
          last_name: 'Doe',
        },
        organization_profile: null,
        created_at: '2026-02-17T12:00:00Z',
        modified_at: '2026-02-17T12:00:00Z',
      }

      const mockRelationship = {
        id: 'rel-person-001-party-001',
        from_party_id: 'person-001',
        to_party_id: 'party-001',
        type: 'IS_EMPLOYEE_OF',
      }

      const mockContactDetails = {
        id: 'contact-person-001',
        type_description: 'Contacto',
        phone: '+34600000000',
        email: 'john@example.com',
        related_party_id: 'person-001',
      }

      // Mock 3 sequential fetch calls
      ;(global.fetch as any)
        .mockResolvedValueOnce({
          ok: true,
          json: async () => mockPersonParty,
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => mockRelationship,
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => mockContactDetails,
        })

      const contactData = {
        email: 'john@example.com',
        phone: '+34600000000',
        firstName: 'John',
        lastName: 'Doe',
        jobTitle: 'Contacto',
      }

      const result = await partyApi.addContact('party-001', contactData)

      expect(result).toMatchObject({
        id: 'person-001',
        first_name: 'John',
        last_name: 'Doe',
        email: 'john@example.com',
        phone: '+34600000000',
      })
      expect(global.fetch).toHaveBeenCalledTimes(3)
    })

    it('should list contacts for a party', async () => {
      const relationshipsResponse = {
        data: [
          {
            id: 'rel-001',
            type: 'IS_EMPLOYEE_OF',
            from_party_id: 'person-001',
            to_party_id: 'party-001',
          },
          {
            id: 'rel-002',
            type: 'IS_EMPLOYEE_OF',
            from_party_id: 'person-002',
            to_party_id: 'party-001',
          },
        ],
      }

      const contactDetailsResponse = {
        data: [
          {
            id: 'cd-001',
            type_description: 'Sales Manager',
            phone: '+34600111111',
            email: 'contact1@example.com',
            related_party_id: 'person-001',
          },
          {
            id: 'cd-002',
            type_description: 'Support',
            phone: '+34600222222',
            email: 'contact2@example.com',
            related_party_id: 'person-002',
          },
        ],
      }

      const personOne = {
        id: 'person-001',
        roles: ['CONTACT'],
        status: 'ACTIVE',
        person_profile: { first_name: 'John', last_name: 'Doe' },
        organization_profile: null,
        created_at: '2026-02-10T08:00:00Z',
        modified_at: '2026-02-10T08:00:00Z',
      }

      const personTwo = {
        id: 'person-002',
        roles: ['CONTACT'],
        status: 'ACTIVE',
        person_profile: { first_name: 'Jane', last_name: 'Smith' },
        organization_profile: null,
        created_at: '2026-02-11T08:00:00Z',
        modified_at: '2026-02-11T08:00:00Z',
      }

      ;(global.fetch as any)
        .mockResolvedValueOnce({ ok: true, json: async () => relationshipsResponse })
        .mockResolvedValueOnce({ ok: true, json: async () => contactDetailsResponse })
        .mockResolvedValueOnce({ ok: true, json: async () => personOne })
        .mockResolvedValueOnce({ ok: true, json: async () => personTwo })

      const result = await partyApi.listContacts('party-001')

      expect(result.data).toHaveLength(2)
      expect(result.total).toBe(2)
      expect(result.data[0]).toMatchObject({
        first_name: 'John',
        email: 'contact1@example.com',
      })
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/parties/party-001/relationships'),
        expect.any(Object)
      )
    })

    it('should get primary contact for a party', async () => {
      const relationshipsResponse = {
        data: [
          {
            id: 'rel-primary',
            type: 'IS_EMPLOYEE_OF',
            from_party_id: 'person-primary',
            to_party_id: 'party-001',
          },
        ],
      }

      const contactDetailsResponse = {
        data: [
          {
            id: 'cd-primary',
            type_description: 'Primary Contact',
            phone: '+34600000000',
            email: 'primary@example.com',
            related_party_id: 'person-primary',
          },
        ],
      }

      const personPrimary = {
        id: 'person-primary',
        roles: ['CONTACT'],
        status: 'ACTIVE',
        person_profile: { first_name: 'Primary', last_name: 'User' },
        organization_profile: null,
        created_at: '2026-02-12T08:00:00Z',
        modified_at: '2026-02-12T08:00:00Z',
      }

      ;(global.fetch as any)
        .mockResolvedValueOnce({ ok: true, json: async () => relationshipsResponse })
        .mockResolvedValueOnce({ ok: true, json: async () => contactDetailsResponse })
        .mockResolvedValueOnce({ ok: true, json: async () => personPrimary })

      const result = await partyApi.getPrimaryContact('party-001')

      expect(result).toBeDefined()
      expect(result?.id).toBe('person-primary')
      expect(result?.is_primary).toBe(true)
    })

    it('should get a contact by ID', async () => {
      const personParty = {
        id: 'person-777',
        roles: ['CONTACT'],
        status: 'ACTIVE',
        person_profile: {
          first_name: 'Luis',
          last_name: 'Gómez',
        },
        organization_profile: null,
        created_at: '2026-02-15T09:00:00Z',
        modified_at: '2026-02-15T09:00:00Z',
      }

      const relationshipPayload = {
        data: [
          {
            id: 'rel-777',
            type: 'IS_EMPLOYEE_OF',
            from_party_id: 'person-777',
            to_party_id: 'party-500',
          },
        ],
      }

      const contactDetailsPayload = {
        data: [
          {
            id: 'cd-777',
            type_description: 'Comercial',
            phone: '+34612345678',
            email: 'luis@example.com',
            related_party_id: 'person-777',
          },
        ],
      }

      ;(global.fetch as any)
        .mockResolvedValueOnce({ ok: true, json: async () => personParty })
        .mockResolvedValueOnce({ ok: true, json: async () => relationshipPayload })
        .mockResolvedValueOnce({ ok: true, json: async () => contactDetailsPayload })

      const result = await partyApi.getContact('person-777')

      expect(result).toMatchObject({
        id: 'person-777',
        first_name: 'Luis',
        last_name: 'Gómez',
        email: 'luis@example.com',
        job_title: 'Comercial',
      })
    })

    it('should get contact without employee relationship', async () => {
      const personParty = {
        id: 'person-888',
        roles: ['CONTACT'],
        status: 'ACTIVE',
        person_profile: {
          first_name: 'Ana',
          last_name: 'Pérez',
        },
        organization_profile: null,
        created_at: '2026-02-15T09:00:00Z',
        modified_at: '2026-02-15T09:00:00Z',
      }

      ;(global.fetch as any)
        .mockResolvedValueOnce({ ok: true, json: async () => personParty })
        .mockResolvedValueOnce({ ok: true, json: async () => ({ data: [] }) })

      const result = await partyApi.getContact('person-888')

      expect(result).toMatchObject({
        id: 'person-888',
        first_name: 'Ana',
        last_name: 'Pérez',
        email: '',
      })
    })

    it('should fail when deleting relationship during removeContact', async () => {
      ;(global.fetch as any)
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({
            data: [
              {
                id: 'cd-001',
                type_description: 'Contacto',
                phone: '+34600111222',
                email: 'contact@example.com',
                related_party_id: 'person-001',
              },
            ],
          }),
        })
        .mockResolvedValueOnce({ ok: true, json: async () => ({}) })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({
            data: [
              {
                id: 'rel-001',
                type: 'IS_EMPLOYEE_OF',
                from_party_id: 'person-001',
                to_party_id: 'party-001',
              },
            ],
          }),
        })
        .mockResolvedValueOnce({
          ok: false,
          status: 400,
          json: async () => ({ message: 'cannot delete relationship' }),
        })

      await expect(partyApi.removeContact('party-001', 'person-001')).rejects.toThrow('cannot delete relationship')
    })

    it('should link an existing contact to a party', async () => {
      const linkedContactParty = {
        id: 'person-500',
        roles: ['CONTACT'],
        status: 'ACTIVE',
        person_profile: {
          first_name: 'Mario',
          last_name: 'Lopez',
        },
        organization_profile: null,
        created_at: '2026-02-17T12:00:00Z',
        modified_at: '2026-02-17T12:00:00Z',
      }

      const relationshipsPayload = {
        data: [
          {
            id: 'rel-500',
            type: 'IS_EMPLOYEE_OF',
            from_party_id: 'person-500',
            to_party_id: 'party-001',
          },
        ],
      }

      const contactDetailsPayload = {
        data: [
          {
            id: 'cd-500',
            type_description: 'Comercial',
            phone: '+34655555555',
            email: 'mario@example.com',
            related_party_id: 'person-500',
          },
        ],
      }

      ;(global.fetch as any)
        .mockResolvedValueOnce({ ok: true, json: async () => ({}) })
        .mockResolvedValueOnce({ ok: true, json: async () => ({}) })
        .mockResolvedValueOnce({ ok: true, json: async () => linkedContactParty })
        .mockResolvedValueOnce({ ok: true, json: async () => relationshipsPayload })
        .mockResolvedValueOnce({ ok: true, json: async () => contactDetailsPayload })

      const result = await partyApi.linkExistingContact('party-001', 'person-500', {
        jobTitle: 'Comercial',
        email: 'mario@example.com',
        phone: '+34655555555',
        isPrimary: false,
      })

      expect(result).toMatchObject({
        id: 'person-500',
        first_name: 'Mario',
        last_name: 'Lopez',
      })
      expect(global.fetch).toHaveBeenNthCalledWith(
        1,
        expect.stringContaining('/parties/person-500/relationships'),
        expect.objectContaining({ method: 'POST' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        2,
        expect.stringContaining('/parties/party-001/contact-details'),
        expect.objectContaining({ method: 'POST' })
      )
    })

    it('should list orphan contacts available to link', async () => {
      const listPayload = {
        data: [
          {
            id: 'person-501',
            roles: ['CONTACT'],
            status: 'ACTIVE',
            organization_profile: null,
            person_profile: { first_name: 'Ana', last_name: 'Ruiz' },
            created_at: '2026-02-10T08:00:00Z',
            modified_at: '2026-02-10T08:00:00Z',
          },
          {
            id: 'person-502',
            roles: ['EMPLOYEE'],
            status: 'ACTIVE',
            organization_profile: null,
            person_profile: { first_name: 'Luis', last_name: 'Diaz' },
            created_at: '2026-02-10T08:00:00Z',
            modified_at: '2026-02-10T08:00:00Z',
          },
        ],
        total: 2,
        page: 1,
        limit: 200,
      }

      const orphanRelationships = { data: [] }
      const linkedRelationships = {
        data: [
          {
            id: 'rel-502',
            type: 'IS_EMPLOYEE_OF',
            from_party_id: 'person-502',
            to_party_id: 'party-777',
          },
        ],
      }

      const orphanParty = {
        id: 'person-501',
        roles: ['CONTACT'],
        status: 'ACTIVE',
        person_profile: { first_name: 'Ana', last_name: 'Ruiz' },
        organization_profile: null,
        created_at: '2026-02-10T08:00:00Z',
        modified_at: '2026-02-10T08:00:00Z',
      }

      const linkedParty = {
        id: 'person-502',
        roles: ['EMPLOYEE'],
        status: 'ACTIVE',
        person_profile: { first_name: 'Luis', last_name: 'Diaz' },
        organization_profile: null,
        created_at: '2026-02-10T08:00:00Z',
        modified_at: '2026-02-10T08:00:00Z',
      }

      ;(global.fetch as any)
        .mockResolvedValueOnce({ ok: true, json: async () => listPayload })
        .mockResolvedValueOnce({ ok: true, json: async () => ({ data: [] }) })
        .mockResolvedValueOnce({ ok: true, json: async () => orphanParty })
        .mockResolvedValueOnce({ ok: true, json: async () => orphanRelationships })
        .mockResolvedValueOnce({ ok: true, json: async () => linkedParty })
        .mockResolvedValueOnce({ ok: true, json: async () => linkedRelationships })
        .mockResolvedValueOnce({ ok: true, json: async () => ({ data: [] }) })

      const result = await partyApi.listAvailableContactsForParty('party-001')

      expect(result).toHaveLength(1)
      expect(result[0].id).toBe('person-501')
      expect(result[0].first_name).toBe('Ana')
    })

    it('should remove contact details and relationship successfully', async () => {
      ;(global.fetch as any)
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({
            data: [
              {
                id: 'cd-001',
                type_description: 'Contacto',
                phone: '+34600111222',
                email: 'contact@example.com',
                related_party_id: 'person-001',
              },
            ],
          }),
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({}),
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({
            data: [
              {
                id: 'rel-001',
                type: 'IS_EMPLOYEE_OF',
                from_party_id: 'person-001',
                to_party_id: 'party-001',
              },
            ],
          }),
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({}),
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({ data: [] }),
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({
            id: 'person-001',
            roles: ['CONTACT'],
            status: 'ACTIVE',
            person_profile: { first_name: 'John', last_name: 'Doe' },
            organization_profile: null,
            created_at: '2026-02-10T08:00:00Z',
            modified_at: '2026-02-10T08:00:00Z',
          }),
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({}),
        })

      await expect(partyApi.removeContact('party-001', 'person-001')).resolves.toBeUndefined()

      expect(global.fetch).toHaveBeenCalledTimes(7)
      expect(global.fetch).toHaveBeenNthCalledWith(
        1,
        expect.stringContaining('/parties/party-001/contact-details'),
        expect.objectContaining({ method: 'GET' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        2,
        expect.stringContaining('/parties/party-001/contact-details/cd-001'),
        expect.objectContaining({ method: 'DELETE' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        3,
        expect.stringContaining('/parties/person-001/relationships'),
        expect.objectContaining({ method: 'GET' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        4,
        expect.stringContaining('/parties/person-001/relationships/rel-001'),
        expect.objectContaining({ method: 'DELETE' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        5,
        expect.stringContaining('/parties/person-001/relationships'),
        expect.objectContaining({ method: 'GET' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        6,
        expect.stringContaining('/parties/person-001'),
        expect.objectContaining({ method: 'GET' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        7,
        expect.stringContaining('/parties/person-001'),
        expect.objectContaining({ method: 'DELETE' })
      )
    })

    it('should handle remove contact when there is no details or relationship', async () => {
      ;(global.fetch as any)
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({ data: [] }),
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({ data: [] }),
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({ data: [] }),
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({
            id: 'person-404',
            roles: ['CONTACT'],
            status: 'ACTIVE',
            person_profile: { first_name: 'Ghost', last_name: 'User' },
            organization_profile: null,
            created_at: '2026-02-10T08:00:00Z',
            modified_at: '2026-02-10T08:00:00Z',
          }),
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({}),
        })

      await expect(partyApi.removeContact('party-001', 'person-404')).resolves.toBeUndefined()

      expect(global.fetch).toHaveBeenCalledTimes(5)
      expect(global.fetch).toHaveBeenNthCalledWith(
        1,
        expect.stringContaining('/parties/party-001/contact-details'),
        expect.objectContaining({ method: 'GET' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        2,
        expect.stringContaining('/parties/person-404/relationships'),
        expect.objectContaining({ method: 'GET' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        3,
        expect.stringContaining('/parties/person-404/relationships'),
        expect.objectContaining({ method: 'GET' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        4,
        expect.stringContaining('/parties/person-404'),
        expect.objectContaining({ method: 'GET' })
      )
      expect(global.fetch).toHaveBeenNthCalledWith(
        5,
        expect.stringContaining('/parties/person-404'),
        expect.objectContaining({ method: 'DELETE' })
      )
    })

    it('should not delete contact entity if it still has employee relationships', async () => {
      ;(global.fetch as any)
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({ data: [] }),
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({
            data: [
              {
                id: 'rel-keep',
                type: 'IS_EMPLOYEE_OF',
                from_party_id: 'person-900',
                to_party_id: 'party-010',
              },
            ],
          }),
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({
            data: [
              {
                id: 'rel-keep',
                type: 'IS_EMPLOYEE_OF',
                from_party_id: 'person-900',
                to_party_id: 'party-010',
              },
            ],
          }),
        })

      await expect(partyApi.removeContact('party-001', 'person-900')).resolves.toBeUndefined()

      expect(global.fetch).toHaveBeenCalledTimes(3)
      expect(global.fetch).not.toHaveBeenCalledWith(
        expect.stringContaining('/parties/person-900'),
        expect.objectContaining({ method: 'DELETE' })
      )
    })

    it('should list party relationships', async () => {
      const mockRelationships = [
        {
          id: 'rel-001',
          type: 'IS_EMPLOYEE_OF',
          from_party_id: 'person-001',
          to_party_id: 'party-001',
        },
        {
          id: 'rel-002',
          type: 'IS_EMPLOYEE_OF',
          from_party_id: 'person-001',
          to_party_id: 'party-002',
        },
      ]

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: mockRelationships }),
      })

      const result = await partyApi.listRelationships('person-001')

      expect(result).toHaveLength(2)
      expect(result[0]).toMatchObject({
        id: 'rel-001',
        type: 'IS_EMPLOYEE_OF',
        from_party_id: 'person-001',
        to_party_id: 'party-001',
      })
      expect(global.fetch).toHaveBeenCalledWith(
        '/api/parties/person-001/relationships',
        expect.objectContaining({ method: 'GET' })
      )
    })

    it('should fail when listing relationships with error', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => ({ message: 'server error' }),
      })

      await expect(partyApi.listRelationships('person-001')).rejects.toThrow()
    })

    it('should remove contact without deleting party when deleteIfNoReferences is false', async () => {
      ;(global.fetch as any)
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({
            data: [
              {
                id: 'cd-001',
                type_description: 'Contacto',
                phone: '+34600111222',
                email: 'contact@example.com',
                related_party_id: 'person-001',
              },
            ],
          }),
        })
        .mockResolvedValueOnce({ ok: true, json: async () => ({}) })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({
            data: [
              {
                id: 'rel-001',
                type: 'IS_EMPLOYEE_OF',
                from_party_id: 'person-001',
                to_party_id: 'party-001',
              },
            ],
          }),
        })
        .mockResolvedValueOnce({ ok: true, json: async () => ({}) })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({
            data: [
              {
                id: 'rel-002',
                type: 'IS_EMPLOYEE_OF',
                from_party_id: 'person-001',
                to_party_id: 'party-002',
              },
            ],
          }),
        })

      await expect(partyApi.removeContact('party-001', 'person-001', false)).resolves.toBeUndefined()

      expect(global.fetch).toHaveBeenCalledTimes(5)
      // Verify the DELETE URL does NOT contain deleteIfNoReferences query param
      const deleteCall = (global.fetch as any).mock.calls[1]
      expect(deleteCall[0]).toContain('/contact-details/cd-001')
      expect(deleteCall[0]).not.toContain('deleteIfNoReferences=true')
    })

    it('should remove contact AND delete party when deleteIfNoReferences is true', async () => {
      ;(global.fetch as any)
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({
            data: [
              {
                id: 'cd-001',
                type_description: 'Contacto',
                phone: '+34600111222',
                email: 'contact@example.com',
                related_party_id: 'person-001',
              },
            ],
          }),
        })
        .mockResolvedValueOnce({ ok: true, json: async () => ({}) })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({
            data: [
              {
                id: 'rel-001',
                type: 'IS_EMPLOYEE_OF',
                from_party_id: 'person-001',
                to_party_id: 'party-001',
              },
            ],
          }),
        })
        .mockResolvedValueOnce({ ok: true, json: async () => ({}) })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({ data: [] }),
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => ({
            id: 'person-001',
            roles: ['CONTACT'],
            status: 'ACTIVE',
            person_profile: { first_name: 'John', last_name: 'Doe', has_person: true, has_organization: false },
            organization_profile: null,
            created_at: '2026-02-10T08:00:00Z',
            modified_at: '2026-02-10T08:00:00Z',
          }),
        })
        .mockResolvedValueOnce({ ok: true, json: async () => ({}) })

      await expect(partyApi.removeContact('party-001', 'person-001', true)).resolves.toBeUndefined()

      expect(global.fetch).toHaveBeenCalledTimes(7)
      // Verify the DELETE URL DOES contain deleteIfNoReferences query param
      const deleteCall = (global.fetch as any).mock.calls[1]
      expect(deleteCall[0]).toContain('/contact-details/cd-001?deleteIfNoReferences=true')
    })
  })

  describe('Address management', () => {
    it('should add address to a party', async () => {
      const mockAddress = {
        id: 'addr-001',
        party_id: 'party-001',
        street: 'Calle Principal 123',
        city: 'Madrid',
        province: 'Madrid',
        postal_code: '28001',
        country: 'España',
        is_primary: true,
        created_at: '2026-02-17T13:00:00Z',
        modified_at: '2026-02-17T13:00:00Z',
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockAddress,
      })

      const addressData = {
        street: 'Calle Principal 123',
        city: 'Madrid',
        province: 'Madrid',
        postalCode: '28001',
        country: 'España',
      }

      const result = await partyApi.addPartyAddress('party-001', addressData)

      expect(result).toEqual(mockAddress)
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/parties/party-001/addresses'),
        expect.objectContaining({
          method: 'POST',
        })
      )
    })

    it('should list addresses for a party', async () => {
      const mockAddresses = {
        data: [
          {
            id: 'addr-001',
            party_id: 'party-001',
            street: 'Calle Principal 123',
            city: 'Madrid',
            province: 'Madrid',
            postal_code: '28001',
            country: 'España',
            is_primary: true,
            created_at: '2026-02-15T08:00:00Z',
            modified_at: '2026-02-15T08:00:00Z',
          },
        ],
        total: 1,
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockAddresses,
      })

      const result = await partyApi.listPartyAddresses('party-001')

      expect(result.data).toHaveLength(1)
      expect(result.total).toBe(1)
      expect(result.data[0].city).toBe('Madrid')
    })

    it('should handle API error when adding address', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => ({ message: 'invalid address' }),
      })

      await expect(
        partyApi.addPartyAddress('party-001', {
          street: 'Calle Falsa',
          city: 'Madrid',
          province: 'Madrid',
          postalCode: '28001',
          country: 'España',
        })
      ).rejects.toThrow('invalid address')
    })

    it('should handle API error when listing addresses', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => ({ message: 'address service down' }),
      })

      await expect(partyApi.listPartyAddresses('party-001')).rejects.toThrow('address service down')
    })

    it('should return null when no primary address exists', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [], total: 0 }),
      })

      const result = await partyApi.getPrimaryPartyAddress('party-001')
      expect(result).toBeNull()
    })

    it('should return first address as primary when available', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          data: [
            {
              street: 'Calle Mayor 1',
              city: 'Madrid',
              province: 'Madrid',
              postal_code: '28013',
              country: 'España',
            },
          ],
          total: 1,
        }),
      })

      const result = await partyApi.getPrimaryPartyAddress('party-001')
      expect(result?.is_primary).toBe(true)
      expect(result?.city).toBe('Madrid')
    })
  })

  describe('getPartiesBatch', () => {
    it('should fetch multiple parties by IDs', async () => {
      // Backend returns simplified PartyBatchItem[] structure
      const mockBatchResponse = [
        {
          id: 'party-001',
          name: 'Party One',
          tax_id: 'B11111111',
          tax_id_type: 'CIF',
        },
        {
          id: 'party-002',
          name: 'Party Two',
          tax_id: 'B22222222',
          tax_id_type: 'CIF',
        },
      ]

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockBatchResponse,
      })

      const result = await partyApi.getPartiesBatch(['party-001', 'party-002'])

      expect(Object.keys(result)).toHaveLength(2)
      expect(result['party-001']?.name).toBe('Party One')
      expect(result['party-002']?.name).toBe('Party Two')
      expect(result['party-001']?.tax_id).toBe('B11111111')
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/parties/batch?ids=party-001%2Cparty-002'),
        expect.any(Object)
      )
    })

    it('should return empty map when no IDs provided', async () => {
      const result = await partyApi.getPartiesBatch([])

      expect(Object.keys(result)).toHaveLength(0)
      expect(global.fetch).not.toHaveBeenCalled()
    })
  })

  describe('Authentication and headers', () => {
    it('should include auth token in requests', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [], total: 0, page: 1, limit: 10 }),
      })

      await partyApi.listParties()

      expect(global.fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: 'Bearer test-token',
          }),
        })
      )
    })

    it('should include user ID in requests', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [], total: 0, page: 1, limit: 10 }),
      })

      await partyApi.listParties()

      expect(global.fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            'X-User-ID': 'user-123',
          }),
        })
      )
    })

    it('should handle missing auth token gracefully', async () => {
      localStorage.removeItem('tramatex_auth_token')

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [], total: 0, page: 1, limit: 10 }),
      })

      await partyApi.listParties()

      expect(global.fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.not.objectContaining({
            Authorization: expect.any(String),
          }),
        })
      )
    })

    it('should fallback to anonymous when user not in localStorage', async () => {
      localStorage.removeItem('tramatex_user')

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: [], total: 0, page: 1, limit: 10 }),
      })

      await partyApi.listParties()

      expect(global.fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            'X-User-ID': 'anonymous',
          }),
        })
      )
    })
  })
})
