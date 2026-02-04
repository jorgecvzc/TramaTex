/**
 * Party Module API Service
 * Handles communication with the backend Party module endpoints
 */

import { getApiBase } from './apiBase'

const API_BASE = getApiBase();

class PartyApiService {
  constructor() {
    this.baseUrl = `${API_BASE}/parties`;
  }

  /**
   * Get authorization header with user token
   */
  getHeaders(additionalHeaders = {}) {
    return {
      'Content-Type': 'application/json',
      'X-User-ID': this.getCurrentUserId(),
      ...additionalHeaders,
    };
  }

  /**
   * Get current user ID from auth context
   */
  getCurrentUserId() {
    return sessionStorage.getItem('userId') || 'anonymous';
  }

  /**
   * Handle API errors
   */
  async handleError(response, message) {
    let errorData;
    try {
      errorData = await response.json();
    } catch {
      errorData = { message: 'Ocurrió un error inesperado' };
    }
    
    const error = new Error(errorData.message || message);
    error.status = response.status;
    error.data = errorData;
    throw error;
  }

  async safeFetch(url, options, fallbackMessage) {
    try {
      return await fetch(url, options);
    } catch (error) {
      const message =
        fallbackMessage ||
        `No se pudo conectar con el servidor. Verifica tu conexión o que la API esté activa. (URL: ${url})`;
      const err = new Error(message);
      err.cause = error;
      throw err;
    }
  }

  // ============================================================================
  // PARTY ENDPOINTS
  // ============================================================================

  mapPartyToParty(party) {
    if (!party || !party.organization_profile) {
      return null;
    }

    const roles = party.roles || [];
    let role = '';
    if (roles.includes('CLIENT') && roles.includes('SUPPLIER')) {
      role = 'BOTH';
    } else if (roles.includes('CLIENT')) {
      role = 'CLIENT';
    } else if (roles.includes('SUPPLIER')) {
      role = 'SUPPLIER';
    }

    return {
      id: party.id,
      name: party.organization_profile.name,
      role,
      status: party.status,
      tax_id: party.organization_profile.tax_id,
      tax_id_type: party.organization_profile.tax_id_type,
      website: party.organization_profile.website,
      created_at: party.created_at,
      modified_at: party.modified_at,
    };
  }

  mapRoleToPartyRoles(role) {
    switch (role) {
      case 'CLIENT':
        return ['CLIENT'];
      case 'SUPPLIER':
        return ['SUPPLIER'];
      case 'BOTH':
        return ['CLIENT', 'SUPPLIER'];
      default:
        return [];
    }
  }

  /**
   * Create a new party
   */
  async createParty(data) {
    const roles = this.mapRoleToPartyRoles(data.role);
    const response = await this.safeFetch(this.baseUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: data.id,
        status: 'ACTIVE',
        roles,
        organization_profile: {
          name: data.name,
          tax_id: data.taxId,
          tax_id_type: data.taxIdType,
          website: data.website,
        },
      }),
    });

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear la entidad');
    }

    const party = await response.json();
    return this.mapPartyToParty(party);
  }

  /**
   * Get party by ID
   */
  async getParty(id) {
    const response = await this.safeFetch(`${this.baseUrl}/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, 'Entidad no encontrada');
    }

    const party = await response.json();
    return this.mapPartyToParty(party);
  }

  /**
   * List parties with filters and pagination
   */
  async listParties(filters = {}) {
    const params = new URLSearchParams();
    
    if (filters.name) params.append('name', filters.name);
    if (filters.role && filters.role !== 'BOTH') params.append('role', filters.role);
    if (filters.status) params.append('status', filters.status);
    if (filters.pageNumber) params.append('page', filters.pageNumber);
    if (filters.pageSize) params.append('page_size', filters.pageSize);
    params.append('type', 'ORGANIZATION');

    const url = params.toString() ? `${this.baseUrl}?${params}` : this.baseUrl;
    
    const response = await this.safeFetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las entidades');
    }

    const payload = await response.json();
    const data = (payload.data || [])
      .map((party) => this.mapPartyToParty(party))
      .filter(Boolean);
    return {
      ...payload,
      data,
    };
  }

  /**
   * Update party
   */
  async updateParty(id, data) {
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
    });

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar la entidad');
    }

    const party = await response.json();
    return this.mapPartyToParty(party);
  }

  /**
   * Change party status (activate/deactivate)
   */
  async changePartyStatus(id, status) {
    const response = await this.safeFetch(`${this.baseUrl}/${id}/status`, {
      method: 'PATCH',
      headers: this.getHeaders(),
      body: JSON.stringify({ status }),
    });

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar el estado');
    }

    const party = await response.json();
    return this.mapPartyToParty(party);
  }

  // ============================================================================
  // CONTACT ENDPOINTS
  // ============================================================================

  mapPartyToContact(party, contactDetails, isPrimary = false) {
    if (!party || !party.person_profile) {
      return null;
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
    };
  }

  /**
   * Add contact to party
   */
  async addContact(partyId, data) {
    const createPersonResponse = await this.safeFetch(this.baseUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: data.id,
        status: 'ACTIVE',
        roles: ['EMPLOYEE'],
        person_profile: {
          first_name: data.firstName,
          last_name: data.lastName,
        },
      }),
    });

    if (!createPersonResponse.ok) {
      await this.handleError(createPersonResponse, 'No se pudo agregar el contacto');
    }

    const personParty = await createPersonResponse.json();

    const relationshipResponse = await this.safeFetch(`${this.baseUrl}/${personParty.id}/relationships`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: `rel-${personParty.id}-${partyId}`,
        to_party_id: partyId,
        type: 'IS_EMPLOYEE_OF',
      }),
    });

    if (!relationshipResponse.ok) {
      await this.handleError(relationshipResponse, 'No se pudo vincular el contacto');
    }

    const contactResponse = await this.safeFetch(`${this.baseUrl}/${partyId}/contact-details`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        id: `contact-${personParty.id}`,
        type_description: data.jobTitle || 'Contacto',
        phone: data.phone,
        email: data.email,
        related_party_id: personParty.id,
      }),
    });

    if (!contactResponse.ok) {
      await this.handleError(contactResponse, 'No se pudieron guardar los datos de contacto');
    }

    const contactDetails = await contactResponse.json();
    return this.mapPartyToContact(personParty, contactDetails, data.isPrimary);
  }

  /**
   * Get contact by ID
   */
  async getContact(id) {
    const response = await this.safeFetch(`${this.baseUrl}/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, 'Contacto no encontrado');
    }

    const personParty = await response.json();
    const relationshipsResponse = await this.safeFetch(`${this.baseUrl}/${id}/relationships`, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!relationshipsResponse.ok) {
      await this.handleError(relationshipsResponse, 'No se pudieron cargar las relaciones');
    }

    const relationshipsPayload = await relationshipsResponse.json();
    const relationships = relationshipsPayload.data || [];
    const employeeRel = relationships.find((rel) => rel.type === 'IS_EMPLOYEE_OF' && rel.to_party_id);

    if (!employeeRel) {
      return this.mapPartyToContact(personParty, null, false);
    }

    const contactsResponse = await this.safeFetch(`${this.baseUrl}/${employeeRel.to_party_id}/contact-details`, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!contactsResponse.ok) {
      await this.handleError(contactsResponse, 'No se pudieron cargar los datos de contacto');
    }

    const contactsPayload = await contactsResponse.json();
    const contact = (contactsPayload.data || []).find((item) => item.related_party_id === id);

    return this.mapPartyToContact(personParty, contact, false);
  }

  /**
   * List contacts for party
   */
  async listContacts(partyId) {
    const response = await this.safeFetch(`${this.baseUrl}/${partyId}/relationships`, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar los contactos');
    }

    const relationshipsPayload = await response.json();
    const relationships = relationshipsPayload.data || [];
    const personIds = relationships
      .filter((rel) => rel.type === 'IS_EMPLOYEE_OF' && rel.to_party_id === partyId)
      .map((rel) => rel.from_party_id);

    const contactsResponse = await this.safeFetch(`${this.baseUrl}/${partyId}/contact-details`, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!contactsResponse.ok) {
      await this.handleError(contactsResponse, 'No se pudieron cargar los datos de contacto');
    }

    const contactsPayload = await contactsResponse.json();
    const contacts = contactsPayload.data || [];

    const persons = await Promise.all(
      personIds.map(async (personId) => {
        const personResponse = await this.safeFetch(`${this.baseUrl}/${personId}`, {
          method: 'GET',
          headers: this.getHeaders(),
        });

        if (!personResponse.ok) {
          return null;
        }

        const personParty = await personResponse.json();
        const contactDetails = contacts.find((item) => item.related_party_id === personId);
        return this.mapPartyToContact(personParty, contactDetails, false);
      })
    );

    return {
      data: persons.filter(Boolean),
      total: persons.filter(Boolean).length,
    };
  }

  /**
   * Get primary contact for party
   */
  async getPrimaryContact(partyId) {
    const response = await this.listContacts(partyId);
    const persons = response.data || [];
    if (persons.length === 0) {
      return null;
    }
    return { ...persons[0], is_primary: true };
  }

  // ============================================================================
  // ADDRESS ENDPOINTS
  // ============================================================================

  /**
   * Add address to party
   */
  async addPartyAddress(partyId, data) {
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
    });

    if (!response.ok) {
      await this.handleError(response, 'No se pudo agregar la dirección');
    }

    return response.json();
  }

  /**
   * List addresses for party
   */
  async listPartyAddresses(partyId) {
    const response = await this.safeFetch(`${this.baseUrl}/${partyId}/addresses`, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las direcciones');
    }

    const payload = await response.json();
    const data = (payload.data || []).map((address) => ({
      id: `${address.street}-${address.postal_code}-${address.city}`,
      street: address.street,
      city: address.city,
      province: address.province,
      postal_code: address.postal_code,
      country: address.country,
      is_primary: false,
      created_at: null,
    }));

    return {
      ...payload,
      data,
    };
  }

  /**
   * Get primary address for party
   */
  async getPrimaryPartyAddress(partyId) {
    const payload = await this.listPartyAddresses(partyId);
    const addresses = payload.data || [];
    if (addresses.length === 0) {
      return null;
    }
    return { ...addresses[0], is_primary: true };
  }

}

// Export singleton instance
export const partyApi = new PartyApiService();

// Also export class for testing
export default PartyApiService;
