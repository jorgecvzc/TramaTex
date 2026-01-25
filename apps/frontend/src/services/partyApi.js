/**
 * Party Module API Service
 * Handles communication with the backend Party module endpoints
 */

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

class PartyApiService {
  constructor() {
    this.baseUrl = `${API_BASE}/organizations`;
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
   * TODO: Integrate with auth module
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
      errorData = { message: 'Unknown error occurred' };
    }
    
    const error = new Error(errorData.message || message);
    error.status = response.status;
    error.data = errorData;
    throw error;
  }

  // ============================================================================
  // ORGANIZATION ENDPOINTS
  // ============================================================================

  /**
   * Create a new organization
   */
  async createOrganization(data) {
    const response = await fetch(this.baseUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      await this.handleError(response, 'Failed to create organization');
    }

    return response.json();
  }

  /**
   * Get organization by ID
   */
  async getOrganization(id) {
    const response = await fetch(`${this.baseUrl}/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, 'Organization not found');
    }

    return response.json();
  }

  /**
   * List organizations with filters and pagination
   */
  async listOrganizations(filters = {}) {
    const params = new URLSearchParams();
    
    if (filters.name) params.append('name', filters.name);
    if (filters.role) params.append('role', filters.role);
    if (filters.status) params.append('status', filters.status);
    if (filters.pageNumber) params.append('page', filters.pageNumber);
    if (filters.pageSize) params.append('page_size', filters.pageSize);

    const url = params.toString() ? `${this.baseUrl}?${params}` : this.baseUrl;
    
    const response = await fetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, 'Failed to fetch organizations');
    }

    return response.json();
  }

  /**
   * Update organization
   */
  async updateOrganization(id, data) {
    const response = await fetch(`${this.baseUrl}/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      await this.handleError(response, 'Failed to update organization');
    }

    return response.json();
  }

  /**
   * Change organization status (activate/deactivate)
   */
  async changeOrganizationStatus(id, status) {
    const response = await fetch(`${this.baseUrl}/${id}/status`, {
      method: 'PATCH',
      headers: this.getHeaders(),
      body: JSON.stringify({ status }),
    });

    if (!response.ok) {
      await this.handleError(response, 'Failed to update organization status');
    }

    return response.json();
  }

  // ============================================================================
  // PERSON ENDPOINTS
  // ============================================================================

  /**
   * Add person (contact) to organization
   */
  async addPerson(organizationId, data) {
    const response = await fetch(`${this.baseUrl}/${organizationId}/persons`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      await this.handleError(response, 'Failed to add person');
    }

    return response.json();
  }

  /**
   * Get person by ID
   */
  async getPerson(id) {
    const response = await fetch(`${API_BASE}/persons/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, 'Person not found');
    }

    return response.json();
  }

  /**
   * List persons in organization
   */
  async listPersons(organizationId) {
    const response = await fetch(`${this.baseUrl}/${organizationId}/persons`, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, 'Failed to fetch persons');
    }

    return response.json();
  }

  /**
   * Get primary contact for organization
   */
  async getPrimaryContact(organizationId) {
    const response = await fetch(`${this.baseUrl}/${organizationId}/primary-contact`, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      // 404 is acceptable if no primary contact exists
      if (response.status === 404) {
        return null;
      }
      await this.handleError(response, 'Failed to fetch primary contact');
    }

    return response.json();
  }

  // ============================================================================
  // ADDRESS ENDPOINTS
  // ============================================================================

  /**
   * Add address to organization
   */
  async addAddress(organizationId, data) {
    const response = await fetch(`${this.baseUrl}/${organizationId}/addresses`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      await this.handleError(response, 'Failed to add address');
    }

    return response.json();
  }

  /**
   * List addresses for organization
   */
  async listAddresses(organizationId) {
    const response = await fetch(`${this.baseUrl}/${organizationId}/addresses`, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      await this.handleError(response, 'Failed to fetch addresses');
    }

    return response.json();
  }

  /**
   * Get primary address for organization
   */
  async getPrimaryAddress(organizationId) {
    const response = await fetch(`${this.baseUrl}/${organizationId}/primary-address`, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      // 404 is acceptable if no primary address exists
      if (response.status === 404) {
        return null;
      }
      await this.handleError(response, 'Failed to fetch primary address');
    }

    return response.json();
  }
}

// Export singleton instance
export const partyApi = new PartyApiService();

// Also export class for testing
export default PartyApiService;
