/**
 * Party Module Type Definitions
 * Defines TypeScript interfaces for Party, Contact, Address entities
 */

// ============================================================================
// ENUMS & LITERALS
// ============================================================================

export type PartyRole = 'CLIENT' | 'SUPPLIER' | 'BOTH' | 'CONTACT'
export type PartyStatus = 'ACTIVE' | 'INACTIVE'
export type TaxIdType = 'RUT' | 'DNI' | 'CUIT' | 'CUIL' | 'OTHER'
export type ContactType = 'EMPLOYEE' | 'EXTERNAL'

// ============================================================================
// PROFILES
// ============================================================================

export interface OrganizationProfile {
  name: string
  tax_id: string | null
  tax_id_type: TaxIdType | null
  website: string | null
}

export interface PersonProfile {
  first_name: string
  last_name: string
}

// ============================================================================
// PARTY ENTITIES
// ============================================================================

/**
 * Party Domain Entity (Backend Response)
 */
export interface Party {
  id: string
  status: PartyStatus
  roles: string[]  // Array of roles like ['CLIENT'], ['SUPPLIER'], or ['CLIENT', 'SUPPLIER']
  organization_profile?: OrganizationProfile
  person_profile?: PersonProfile
  created_at: string
  modified_at: string
}

/**
 * Party UI Model (Mapped for Frontend)
 */
export interface PartyUI {
  id: string
  name: string
  role: PartyRole
  status: PartyStatus
  tax_id: string | null
  tax_id_type: TaxIdType | null
  website: string | null
  phone: string | null
  email: string | null
  created_at: string
  modified_at: string
  has_organization: boolean
  has_person: boolean
}

// ============================================================================
// CONTACT ENTITIES
// ============================================================================

export interface Contact {
  id: string
  first_name: string
  last_name: string
  email: string
  phone: string
  job_title: string
  is_primary: boolean
  created_at: string
}

export interface ContactDetails {
  email: string
  phone: string
  type_description: string
  related_party_id?: string  // Added for backend responses
}

// ============================================================================
// ADDRESS ENTITIES
// ============================================================================

export interface Address {
  id: string
  street: string
  city: string
  province: string
  country: string
  postal_code: string
  is_primary: boolean
  created_at: string
}

// ============================================================================
// REQUEST DTOs
// ============================================================================

export type EntityType = 'PERSON' | 'ORGANIZATION'

export interface CreatePartyRequest {
  id?: string
  name: string
  role: PartyRole
  taxId?: string | null
  taxIdType?: TaxIdType | null
  website?: string | null
  phone?: string | null
  email?: string | null
  entityType?: EntityType  // New field to specify entity type
  firstName?: string  // For person entities
  lastName?: string   // For person entities
}

export interface UpdatePartyRequest {
  name: string
  role?: PartyRole
  taxId?: string | null
  taxIdType?: TaxIdType | null
  website?: string | null
  phone?: string | null
  email?: string | null
}

export interface ChangePartyStatusRequest {
  status: PartyStatus
}

export interface CreateContactRequest {
  id?: string
  existingContactId?: string
  firstName: string
  lastName: string
  email: string
  phone: string
  jobTitle: string
  isPrimary?: boolean
}

export interface UpdateContactRequest {
  email: string
  phone: string
  typeDescription: string
  isPrimary: boolean
}

// ============================================================================
// FILTERS & PAGINATION
// ============================================================================

export interface ListPartiesFilters {
  name?: string
  role?: PartyRole
  status?: PartyStatus
  type?: 'organization' | 'person'
  pageNumber?: number
  pageSize?: number
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// ============================================================================
// BATCH OPERATIONS
// ============================================================================

export interface PartyBatchItem {
  id: string
  name: string
  reference?: string
  tax_id: string | null
  tax_id_type: TaxIdType | null
}

export type PartyBatchMap = Record<string, PartyBatchItem>

// ============================================================================
// API ERROR HANDLING
// ============================================================================

export interface PartyError {
  message: string
  status?: number
  data?: unknown
  cause?: Error
}
