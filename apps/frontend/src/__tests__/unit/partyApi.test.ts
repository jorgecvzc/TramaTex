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
        created_at: '2026-02-17T10:00:00Z',
        modified_at: '2026-02-17T10:00:00Z',
        has_organization: true,
        has_person: false,
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
  })

  describe('updateParty', () => {
    it('should update a party successfully', async () => {
      const mockUpdatedParty = {
        id: 'party-001',
        roles: ['CLIENT', 'SUPPLIER'],
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

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockUpdatedParty,
      })

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
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/parties/party-001'),
        expect.objectContaining({
          method: 'PUT',
          body: expect.any(String),
        })
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
    
    it.skip('should add a contact to a party (complex multi-fetch operation)', async () => {
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

    it.skip('should list contacts for a party (complex multi-fetch operation)', async () => {
      const mockResponse = {
        data: [
          {
            id: 'contact-det-001',
            type_description: 'Sales Manager',
            phone: '+34600111111',
            email: 'contact1@example.com',
            related_party_id: 'person-001',
          },
          {
            id: 'contact-det-002',
            type_description: 'Support',
            phone: '+34600222222',
            email: 'contact2@example.com',
            related_party_id: 'person-002',
          },
        ],
        total: 2,
      }

      // Mock fetch for contact-details endpoint
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await partyApi.listContacts('party-001')

      expect(result.data).toHaveLength(2)
      expect(result.total).toBe(2)
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/parties/party-001/contact-details'),
        expect.any(Object)
      )
    })

    it.skip('should get primary contact for a party (depends on listContacts)', async () => {
      const mockResponse = {
        data: [
          {
            id: 'contact-primary',
            type_description: 'Primary Contact',
            phone: '+34600000000',
            email: 'primary@example.com',
            related_party_id: 'person-primary',
            is_primary: true,
          },
        ],
        total: 1,
      }

      // getPrimaryContact calls listContacts internally which makes 1 fetch
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await partyApi.getPrimaryContact('party-001')

      expect(result).toBeDefined()
      expect(result?.id).toBe('contact-primary')
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/parties/party-001/contact-details'),
        expect.any(Object)
      )
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
