/**
 * Party Module API Service
 * Handles communication with the backend Party module endpoints
 */

import { getApiBase } from './apiBase'
import type {
  Party,
  PartyUI,
  PartyRole,
  PartyStatus,
  Contact,
  ContactDetails,
  Address,
  CreatePartyRequest,
  UpdatePartyRequest,
  CreateContactRequest,
  ListPartiesFilters,
  PaginatedResponse,
  PartyBatchMap,
  PartyBatchItem,
} from '../types/party'

const API_BASE = getApiBase()

interface PartyApiError extends Error {
  status?: number
  data?: unknown
  cause?: Error
}

class PartyApiService {
  private baseUrl: string

  constructor() {
    this.baseUrl = `${API_BASE}/parties`
  }

  private generateId(): string {
    if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
      return crypto.randomUUID()
    }

    const random = Math.random().toString(16).slice(2, 10)
    const timestamp = Date.now().toString(16)
    return `${timestamp}${random}`.slice(0, 36)
  }

  /**
   * Get authorization header with user token
   */
  private getHeaders(additionalHeaders: Record<string, string> = {}): Record<string, string> {
    const token = localStorage.getItem('tramatex_auth_token')
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'X-User-ID': this.getCurrentUserId(),
      ...additionalHeaders,
    }
    
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }
    
    return headers
  }

  /**
   * Get current user ID from auth context
   */
  private getCurrentUserId(): string {
    try {
      const userStr = localStorage.getItem('tramatex_user')
      if (userStr) {
        const user = JSON.parse(userStr)
        return user.id || 'anonymous'
      }
    } catch (error) {
      console.error('[PartyApi] Error parsing user:', error)
    }
    return 'anonymous'
  }

  /**
   * Handle API errors
   */
  private async handleError(response: Response, message: string): Promise<never> {
    let errorData: { message?: string; error?: string } | undefined
    try {
      errorData = await response.json()
    } catch {
      errorData = { message: 'Ocurrió un error inesperado' }
    }
    
    const error = new Error(errorData?.message || errorData?.error || message) as PartyApiError
    error.status = response.status
    error.data = errorData
    throw error
  }

  private async safeFetch(url: string, options: RequestInit, fallbackMessage?: string): Promise<Response> {
    try {
      return await fetch(url, options)
    } catch (error) {
      const message =
        fallbackMessage ||
        `No se pudo conectar con el servidor. Verifica tu conexión o que la API esté activa. (URL: ${url})`
      const err = new Error(message) as PartyApiError
      err.cause = error as Error
      throw err
    }
  }

  // ============================================================================
  // PARTY ENDPOINTS
  // ============================================================================

  private mapPartyToParty(party: Party | null): PartyUI | null {
    if (!party) {
      return null
    }

    const roles = party.roles || []
    let role: PartyRole = 'CLIENT'
    if (roles.includes('CLIENT') && roles.includes('SUPPLIER')) {
      role = 'BOTH'
    } else if (roles.includes('CLIENT')) {
      role = 'CLIENT'
    } else if (roles.includes('SUPPLIER')) {
      role = 'SUPPLIER'
    } else if (roles.includes('CONTACT') || roles.includes('EMPLOYEE')) {
      role = 'CONTACT'
    }

    // Support both organization and person profiles (FIXED 2026-02-14)
    let name = '(Sin nombre)'
    let tax_id: string | null = null
    let tax_id_type = null
    let website: string | null = null

    if (party.organization_profile) {
      name = party.organization_profile.name
      tax_id = party.organization_profile.tax_id
      tax_id_type = party.organization_profile.tax_id_type
      website = party.organization_profile.website
    } else if (party.person_profile) {
      name = `${party.person_profile.first_name} ${party.person_profile.last_name}`
    }

    return {
      id: party.id,
      name,
      role,
      status: party.status,
      tax_id,
      tax_id_type,
      website,
      default_discount_percentage: party.default_discount_percentage ?? 0,
      created_at: party.created_at,
      modified_at: party.modified_at,
      has_organization: !!party.organization_profile,
      has_person: !!party.person_profile,
    }
  }

  private mapRoleToPartyRoles(role: PartyRole): string[] {
    switch (role) {
      case 'CLIENT':
        return ['CLIENT']
      case 'SUPPLIER':
        return ['SUPPLIER']
      case 'BOTH':
        return ['CLIENT', 'SUPPLIER']
      case 'CONTACT':
        return ['CONTACT']
      default:
        return []
    }
  }

  private async addRoleWithCompatibility(id: string, role: string): Promise<void> {
    const addResponse = await this.safeFetch(`${this.baseUrl}/${id}/roles`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ role }),
    })

    if (addResponse.ok) {
      return
    }

    if (role === 'CONTACT') {
      const fallbackResponse = await this.safeFetch(`${this.baseUrl}/${id}/roles`, {
        method: 'POST',
        headers: this.getHeaders(),
        body: JSON.stringify({ role: 'EMPLOYEE' }),
      })

      if (fallbackResponse.ok) {
        return
      }
    }

    await this.handleError(addResponse, `No se pudo agregar el rol ${role}`)
  }

  private async syncPartyRoles(id: string, currentRoles: string[], targetRole: PartyRole): Promise<void> {
    const normalizedCurrentRoles = (currentRoles || []).map((role) => String(role).toUpperCase())
    const targetRoles = this.mapRoleToPartyRoles(targetRole)

    const normalizedTargetRoles =
      targetRole === 'CONTACT' ? ['CONTACT', 'EMPLOYEE'] : targetRoles

    const rolesToAdd = targetRoles.filter((role) => !normalizedCurrentRoles.includes(role))
    const rolesToRemove = normalizedCurrentRoles.filter((role) => !normalizedTargetRoles.includes(role))

    for (const role of rolesToAdd) {
      await this.addRoleWithCompatibility(id, role)
    }

    for (const role of rolesToRemove) {
      const removeResponse = await this.safeFetch(`${this.baseUrl}/${id}/roles/${role}`, {
        method: 'DELETE',
        headers: this.getHeaders(),
      })

      if (!removeResponse.ok) {
        await this.handleError(removeResponse, `No se pudo quitar el rol ${role}`)
      }
    }
  }

  /**
   * Create a new party
   */
  async createParty(data: CreatePartyRequest): Promise<PartyUI | null> {
    const roles = this.mapRoleToPartyRoles(data.role)
    
    // Determine entity type and build appropriate profile
    const entityType = data.entityType || 'ORGANIZATION'
    const body: any = {
      id: data.id,
      status: 'ACTIVE',
      roles,
    }

    if (entityType === 'PERSON') {
      // Create person profile
      body.person_profile = {
        first_name: data.firstName || '',
        last_name: data.lastName || '',
      }
      // Person entities can still have tax_id (NIE/NIF)
      if (data.taxId) {
        body.organization_profile = {
          name: `${data.firstName} ${data.lastName}`,
          tax_id: data.taxId,
          tax_id_type: data.taxIdType,
          website: data.website,
        }
      }
    } else {
      // Create organization profile
      body.organization_profile = {
        name: data.name,
        tax_id: data.taxId,
        tax_id_type: data.taxIdType,
        website: data.website,
      }
    }

    const response = await this.safeFetch(this.baseUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(body),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear la entidad')
    }

    const party: Party = await response.json()
    return this.mapPartyToParty(party)
  }

  /**
   * Get party by ID
   */
  async getParty(id: string): Promise<PartyUI | null> {
    const response = await this.safeFetch(`${this.baseUrl}/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'Entidad no encontrada')
    }

    const party: Party = await response.json()
    return this.mapPartyToParty(party)
  }

  /**
   * Get multiple parties by IDs (batch operation)
   */
  async getPartiesBatch(ids: string[]): Promise<PartyBatchMap> {
    if (!ids || ids.length === 0) {
      return {}
    }

    // Remove duplicates
    const uniqueIds = [...new Set(ids)]
    
    const idsParam = uniqueIds.join(',')
    const response = await this.safeFetch(
      `${this.baseUrl}/batch?ids=${encodeURIComponent(idsParam)}`,
      {
        method: 'GET',
        headers: this.getHeaders(),
      }
    )

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las entidades')
    }

    const parties: PartyBatchItem[] = await response.json()
    
    // Convert array to map for easy lookup
    const partyMap: PartyBatchMap = {}
    for (const party of parties) {
      partyMap[party.id] = {
        id: party.id,
        name: party.name,
        reference: party.reference,
        tax_id: party.tax_id,
        tax_id_type: party.tax_id_type,
      }
    }
    
    return partyMap
  }

  /**
   * List parties with filters and pagination
   */
  async listParties(filters: ListPartiesFilters = {}): Promise<PaginatedResponse<PartyUI>> {
    const params = new URLSearchParams()
    
    if (filters.name) params.append('name', filters.name)
    if (filters.role && filters.role !== 'BOTH') params.append('role', filters.role)
    if (filters.status) params.append('status', filters.status)
    if (filters.pageNumber) params.append('page', filters.pageNumber.toString())
    if (filters.pageSize) params.append('page_size', filters.pageSize.toString())
    // Allow filtering by type if provided, otherwise show all
    if (filters.type) params.append('type', filters.type)

    const url = params.toString() ? `${this.baseUrl}?${params}` : this.baseUrl

    let response = await this.safeFetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok && filters.role === 'CONTACT') {
      let errorData: { message?: string; error?: string } | undefined
      try {
        errorData = await response.json()
      } catch {
        errorData = undefined
      }

      const errorMessage = (errorData?.message || errorData?.error || '').toLowerCase()
      if (response.status === 400 && errorMessage.includes('invalid role') && errorMessage.includes('contact')) {
        const fallbackParams = new URLSearchParams(params)
        fallbackParams.set('role', 'EMPLOYEE')
        const fallbackUrl = `${this.baseUrl}?${fallbackParams.toString()}`
        response = await this.safeFetch(fallbackUrl, {
          method: 'GET',
          headers: this.getHeaders(),
        })
      }
    }

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las entidades')
    }

    const payload: { data: Party[] } & Omit<PaginatedResponse<Party>, 'data'> = await response.json()
    const data = (payload.data || [])
      .map((party) => this.mapPartyToParty(party))
      .filter((party): party is PartyUI => party !== null)
    return {
      ...payload,
      data,
    }
  }

  /**
   * Update party
   */
  async updateParty(id: string, data: UpdatePartyRequest): Promise<PartyUI | null> {
    const response = await this.safeFetch(`${this.baseUrl}/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify({
        organization_profile: {
          name: data.name,
          tax_id: data.taxId,
          tax_id_type: data.taxIdType,
          website: data.website,
        },
      }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar la entidad')
    }

    const party: Party = await response.json()
    let mappedParty = this.mapPartyToParty(party)
    const currentRoles = party.roles || []
    const targetRoles = data.role ? this.mapRoleToPartyRoles(data.role) : []
    const requiresRoleSync =
      !!data.role &&
      (targetRoles.length !== currentRoles.length ||
        targetRoles.some((role) => !currentRoles.includes(role)))

    if (data.role && requiresRoleSync) {
      await this.syncPartyRoles(id, currentRoles, data.role)
      mappedParty = await this.getParty(id)
    }

    return mappedParty
  }

  /**
   * Change party status (activate/deactivate)
   */
  async changePartyStatus(id: string, status: PartyStatus): Promise<PartyUI | null> {
    const response = await this.safeFetch(`${this.baseUrl}/${id}/status`, {
      method: 'PATCH',
      headers: this.getHeaders(),
      body: JSON.stringify({ status }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar el estado')
    }

    const party: Party = await response.json()
    return this.mapPartyToParty(party)
  }

  async deleteParty(id: string): Promise<void> {
    const response = await this.safeFetch(`${this.baseUrl}/${id}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo eliminar la entidad')
    }
  }

  /**
   * List party relationships
   */
  async listRelationships(partyId: string): Promise<Array<{ id: string; type: string; to_party_id: string; from_party_id: string }>> {
    const response = await this.safeFetch(`${this.baseUrl}/${partyId}/relationships`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las relaciones')
    }

    const payload: { data: Array<{ id: string; type: string; to_party_id: string; from_party_id: string }> } = await response.json()
    return payload.data || []
  }

  // ============================================================================
  // CONTACT ENDPOINTS
  // ============================================================================

  private mapPartyToContact(party: Party | null, contactDetails: ContactDetails | null | undefined, isPrimary = false): Contact | null {
    if (!party || !party.person_profile) {
      return null
    }

    return {
      id: party.id,
      first_name: party.person_profile.first_name,
      last_name: party.person_profile.last_name,
      email: contactDetails?.email || '',
      phone: contactDetails?.phone || '',
      job_title: contactDetails?.type_description || '',
      is_primary: isPrimary,
      created_at: party.created_at,
    }
  }

  /**
   * Add contact to party
   */
  async addContact(partyId: string, data: CreateContactRequest): Promise<Contact | null> {
    const personId = (data.id || '').trim() || this.generateId()

    const createPersonResponse = await this.safeFetch(this.baseUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: personId,
        status: 'ACTIVE',
        roles: ['CONTACT'],
        person_profile: {
          first_name: data.firstName,
          last_name: data.lastName,
        },
      }),
    })

    if (!createPersonResponse.ok) {
      await this.handleError(createPersonResponse, 'No se pudo agregar el contacto')
    }

    const personParty: Party = await createPersonResponse.json()

    const relationshipResponse = await this.safeFetch(`${this.baseUrl}/${personParty.id}/relationships`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: this.generateId(),
        to_party_id: partyId,
        type: 'IS_EMPLOYEE_OF',
      }),
    })

    if (!relationshipResponse.ok) {
      await this.handleError(relationshipResponse, 'No se pudo vincular el contacto')
    }

    const contactResponse = await this.safeFetch(`${this.baseUrl}/${partyId}/contact-details`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: this.generateId(),
        type_description: data.jobTitle || 'Contacto',
        phone: data.phone,
        email: data.email,
        related_party_id: personParty.id,
      }),
    })

    if (!contactResponse.ok) {
      await this.handleError(contactResponse, 'No se pudieron guardar los datos de contacto')
    }

    const contactDetails: ContactDetails = await contactResponse.json()
    return this.mapPartyToContact(personParty, contactDetails, false)
  }

  async linkExistingContact(
    partyId: string,
    contactId: string,
    data: Pick<CreateContactRequest, 'jobTitle' | 'email' | 'phone' | 'isPrimary'>,
  ): Promise<Contact | null> {
    const relationshipResponse = await this.safeFetch(`${this.baseUrl}/${contactId}/relationships`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: this.generateId(),
        to_party_id: partyId,
        type: 'IS_EMPLOYEE_OF',
      }),
    })

    if (!relationshipResponse.ok) {
      await this.handleError(relationshipResponse, 'No se pudo vincular el contacto existente')
    }

    const contactResponse = await this.safeFetch(`${this.baseUrl}/${partyId}/contact-details`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: this.generateId(),
        type_description: data.jobTitle || 'Contacto',
        phone: data.phone,
        email: data.email,
        related_party_id: contactId,
      }),
    })

    if (!contactResponse.ok) {
      await this.handleError(contactResponse, 'No se pudieron guardar los datos del contacto existente')
    }

    return this.getContact(contactId)
  }

  async listAvailableContactsForParty(partyId: string): Promise<Contact[]> {
    // Get all entities without role filter to ensure we get all physical persons
    const listResponse = await this.listParties({ pageNumber: 1, pageSize: 500 })
    
    // Any physical person (persona física) can be a contact, regardless of role
    const candidateParties = (listResponse.data || []).filter((party) => {
      const hasPerson = party.has_person
      const isNotSelf = party.id !== partyId
      return hasPerson && isNotSelf
    })

    // Get current party's contacts to exclude them
    const currentContactsResponse = await this.safeFetch(`${this.baseUrl}/${partyId}/contact-details`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    let currentContactIds: string[] = []
    if (currentContactsResponse.ok) {
      const currentContactsPayload: { data: Array<{ related_party_id: string }> } = await currentContactsResponse.json()
      currentContactIds = (currentContactsPayload.data || []).map(c => c.related_party_id)
    }

    // Filter out contacts already linked to this party
    const availableParties = candidateParties.filter(party => !currentContactIds.includes(party.id))

    // Convert to Contact objects
    const contacts = await Promise.all(
      availableParties.map(async (party) => {
        try {
          return await this.getContact(party.id)
        } catch (error) {
          console.error(`Error al obtener contacto ${party.id}:`, error)
          return null
        }
      })
    )

    const validContacts = contacts.filter((contact): contact is Contact => contact !== null)
    return validContacts
  }

  /**
   * Get contact by ID
   */
  async getContact(id: string): Promise<Contact | null> {
    const response = await this.safeFetch(`${this.baseUrl}/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'Contacto no encontrado')
    }

    const personParty: Party = await response.json()
    const relationshipsResponse = await this.safeFetch(`${this.baseUrl}/${id}/relationships`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!relationshipsResponse.ok) {
      await this.handleError(relationshipsResponse, 'No se pudieron cargar las relaciones')
    }

    const relationshipsPayload: { data: Array<{ type: string; to_party_id: string; from_party_id: string }> } = await relationshipsResponse.json()
    const relationships = relationshipsPayload.data || []
    const employeeRel = relationships.find((rel) => rel.type === 'IS_EMPLOYEE_OF' && rel.to_party_id)

    if (!employeeRel) {
      return this.mapPartyToContact(personParty, null, false)
    }

    const contactsResponse = await this.safeFetch(`${this.baseUrl}/${employeeRel.to_party_id}/contact-details`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!contactsResponse.ok) {
      await this.handleError(contactsResponse, 'No se pudieron cargar los datos de contacto')
    }

    const contactsPayload: { data: Array<ContactDetails & { related_party_id: string }> } = await contactsResponse.json()
    const contact = (contactsPayload.data || []).find((item) => item.related_party_id === id)

    return this.mapPartyToContact(personParty, contact, false)
  }

  /**
   * List contacts for party
   */
  async listContacts(partyId: string): Promise<{ data: Contact[]; total: number }> {
    const response = await this.safeFetch(`${this.baseUrl}/${partyId}/relationships`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar los contactos')
    }

    const relationshipsPayload: { data: Array<{ type: string; to_party_id: string; from_party_id: string }> } = await response.json()
    const relationships = relationshipsPayload.data || []
    const personIds = relationships
      .filter((rel) => rel.type === 'IS_EMPLOYEE_OF' && rel.to_party_id === partyId)
      .map((rel) => rel.from_party_id)

    const contactsResponse = await this.safeFetch(`${this.baseUrl}/${partyId}/contact-details`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!contactsResponse.ok) {
      await this.handleError(contactsResponse, 'No se pudieron cargar los datos de contacto')
    }

    const contactsPayload: { data: Array<ContactDetails & { related_party_id: string }> } = await contactsResponse.json()
    const contacts = contactsPayload.data || []

    const persons = await Promise.all(
      personIds.map(async (personId) => {
        const personResponse = await this.safeFetch(`${this.baseUrl}/${personId}`, {
          method: 'GET',
          headers: this.getHeaders(),
        })

        if (!personResponse.ok) {
          return null
        }

        const personParty: Party = await personResponse.json()
        const contactDetails = contacts.find((item) => item.related_party_id === personId)
        return this.mapPartyToContact(personParty, contactDetails, false)
      })
    )

    const validPersons = persons.filter((person): person is Contact => person !== null)
    return {
      data: validPersons,
      total: validPersons.length,
    }
  }

  /**
   * Get primary contact for party
   */
  async getPrimaryContact(partyId: string): Promise<Contact | null> {
    const response = await this.listContacts(partyId)
    const persons = response.data || []
    if (persons.length === 0) {
      return null
    }
    return { ...persons[0], is_primary: true }
  }

  /**
   * Remove contact from party
   * @param deleteIfNoReferences If true, also deletes the contact party if it has no other references
   */
  async removeContact(partyId: string, contactId: string, deleteIfNoReferences: boolean = false): Promise<void> {
    const contactsResponse = await this.safeFetch(`${this.baseUrl}/${partyId}/contact-details`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!contactsResponse.ok) {
      await this.handleError(contactsResponse, 'No se pudieron cargar los datos de contacto')
    }

    const contactsPayload: { data: Array<ContactDetails & { id: string; related_party_id: string }> } = await contactsResponse.json()
    const contacts = contactsPayload.data || []
    const contactDetails = contacts.find((item) => item.related_party_id === contactId)

    if (contactDetails?.id) {
      const deleteUrl = `${this.baseUrl}/${partyId}/contact-details/${contactDetails.id}${deleteIfNoReferences ? '?deleteIfNoReferences=true' : ''}`
      const deleteContactResponse = await this.safeFetch(deleteUrl, {
        method: 'DELETE',
        headers: this.getHeaders(),
      })

      if (!deleteContactResponse.ok) {
        await this.handleError(deleteContactResponse, 'No se pudo eliminar el detalle del contacto')
      }
    }

    const relationshipsResponse = await this.safeFetch(`${this.baseUrl}/${contactId}/relationships`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!relationshipsResponse.ok) {
      await this.handleError(relationshipsResponse, 'No se pudieron cargar las relaciones del contacto')
    }

    const relationshipsPayload: {
      data: Array<{ id: string; type: string; to_party_id: string; from_party_id: string }>
    } = await relationshipsResponse.json()

    const relationship = (relationshipsPayload.data || []).find(
      (rel) => rel.type === 'IS_EMPLOYEE_OF' && rel.to_party_id === partyId
    )

    if (relationship?.id) {
      const deleteRelationshipResponse = await this.safeFetch(
        `${this.baseUrl}/${contactId}/relationships/${relationship.id}`,
        {
          method: 'DELETE',
          headers: this.getHeaders(),
        }
      )

      if (!deleteRelationshipResponse.ok) {
        await this.handleError(deleteRelationshipResponse, 'No se pudo desvincular el contacto')
      }
    }

    const remainingRelationshipsResponse = await this.safeFetch(`${this.baseUrl}/${contactId}/relationships`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!remainingRelationshipsResponse.ok) {
      await this.handleError(remainingRelationshipsResponse, 'No se pudo verificar el estado del contacto')
    }

    const remainingRelationshipsPayload: {
      data: Array<{ id: string; type: string; to_party_id: string; from_party_id: string }>
    } = await remainingRelationshipsResponse.json()

    const hasEmployeeRelationships = (remainingRelationshipsPayload.data || []).some(
      (rel) => rel.type === 'IS_EMPLOYEE_OF'
    )

    if (!hasEmployeeRelationships) {
      const contactParty = await this.getParty(contactId)
      if (contactParty && contactParty.role === 'CONTACT' && contactParty.has_person && !contactParty.has_organization) {
        await this.deleteParty(contactId)
      }
    }
  }

  // ============================================================================
  // ADDRESS ENDPOINTS
  // ============================================================================

  /**
   * Add address to party
   */
  async addPartyAddress(partyId: string, data: { id?: string; street: string; city: string; province: string; postalCode: string; country: string }): Promise<Address> {
    const response = await this.safeFetch(`${this.baseUrl}/${partyId}/addresses`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: data.id,
        street: data.street,
        city: data.city,
        province: data.province,
        postal_code: data.postalCode,
        country: data.country,
      }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo agregar la dirección')
    }

    return response.json()
  }

  /**
   * List addresses for party
   */
  async listPartyAddresses(partyId: string): Promise<{ data: Address[]; total: number }> {
    const response = await this.safeFetch(`${this.baseUrl}/${partyId}/addresses`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las direcciones')
    }

    const payload: { data: Array<{ id: string; street: string; city: string; province: string; postal_code: string; country: string; is_primary?: boolean; created_at?: string }> } = await response.json()
    const data: Address[] = (payload.data || []).map((address) => ({
      id: address.id,
      street: address.street,
      city: address.city,
      province: address.province,
      country: address.country,
      postal_code: address.postal_code,
      is_primary: address.is_primary || false,
      created_at: address.created_at || '',
    }))

    return {
      data,
      total: data.length,
    }
  }

  /**
   * Update address for party
   */
  async updatePartyAddress(partyId: string, addressId: string, data: { id?: string; street: string; city: string; province: string; postalCode: string; country: string }): Promise<Address> {
    const response = await this.safeFetch(`${this.baseUrl}/${partyId}/addresses/${addressId}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify({
        street: data.street,
        city: data.city,
        province: data.province,
        postal_code: data.postalCode,
        country: data.country,
      }),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar la dirección')
    }

    return response.json()
  }

  /**
   * Delete address from party
   */
  async deletePartyAddress(partyId: string, addressId: string): Promise<void> {
    const response = await this.safeFetch(`${this.baseUrl}/${partyId}/addresses/${addressId}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo eliminar la dirección')
    }
  }

  /**
   * Get primary address for party
   */
  async getPrimaryPartyAddress(partyId: string): Promise<Address | null> {
    const payload = await this.listPartyAddresses(partyId)
    const addresses = payload.data || []
    if (addresses.length === 0) {
      return null
    }
    return { ...addresses[0], is_primary: true }
  }

}

// Export singleton instance
export const partyApi = new PartyApiService()

// Also export class for testing
export default PartyApiService
