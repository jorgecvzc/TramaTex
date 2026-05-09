import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/vue'
import { createPinia, setActivePinia } from 'pinia'
import PartyForm from '../../../components/party/PartyForm.vue'
import { partyApi } from '../../../services/partyApi'
import { useToastStore } from '../../../stores/toast'

// Mock partyApi
vi.mock('../../../services/partyApi', () => ({
  partyApi: {
    createParty: vi.fn(),
    updateParty: vi.fn(),
  }
}))

describe('PartyForm Component', () => {
  let toastStore: any

  beforeEach(() => {
    setActivePinia(createPinia())
    toastStore = useToastStore()
    vi.clearAllMocks()
    vi.spyOn(toastStore, 'error')
    vi.spyOn(toastStore, 'success')
    vi.spyOn(toastStore, 'warning')
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  describe('Initial rendering', () => {
    it('should render form with all fields in create mode', async () => {
      render(PartyForm)

      expect(screen.getByRole('heading', { name: /Crear entidad/i })).toBeInTheDocument()
      expect(screen.getByLabelText(/Rol de la entidad/i)).toBeInTheDocument()

      // Select entity type to reveal name field
      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'ORGANIZATION')
      expect(screen.getByLabelText(/Nombre de la organización/i)).toBeInTheDocument()
      expect(screen.getByPlaceholderText('p. ej., 12345678A')).toBeInTheDocument()
      expect(screen.getByLabelText(/Tipo identificación/i)).toBeInTheDocument()
      expect(screen.getByLabelText(/Sitio web/i)).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /Crear entidad/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /Reiniciar/i })).toBeInTheDocument()
    })

    it('should render form in edit mode when partyId is provided', () => {
      render(PartyForm, {
        props: {
          partyId: 'party-123',
          initialData: {
            name: 'Acme Corporation',
            role: 'CLIENT',
            entityType: 'ORGANIZATION',
            taxId: 'B12345678',
          }
        }
      })

      expect(screen.getByRole('heading', { name: /Editar entidad/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /Actualizar entidad/i })).toBeInTheDocument()
      expect(screen.getByLabelText(/Notas/i)).toBeInTheDocument()
    })

    it('should populate fields with initialData', async () => {
      render(PartyForm, {
        props: {
          initialData: {
            name: 'Test Company',
            role: 'SUPPLIER',
            entityType: 'ORGANIZATION',
            taxId: 'A87654321',
            taxIdType: 'CIF',
            website: 'https://test.com',
          }
        }
      })

      // Wait for entity type to be set so org name field appears
      await waitFor(() => {
        expect(screen.getByLabelText(/Nombre de la organización/i)).toBeInTheDocument()
      })

      // Check input values directly
      const nameInput = screen.getByLabelText(/Nombre de la organización/i) as HTMLInputElement
      const roleSelect = screen.getByLabelText(/Rol de la entidad/i) as HTMLSelectElement
      const taxIdInput = screen.getByPlaceholderText('p. ej., 12345678A') as HTMLInputElement
      const taxIdTypeSelect = screen.getByLabelText(/Tipo identificación/i) as HTMLSelectElement
      const websiteInput = screen.getByLabelText(/Sitio web/i) as HTMLInputElement

      expect(nameInput.value).toBe('Test Company')
      expect(roleSelect.value).toBe('SUPPLIER')
      expect(taxIdInput.value).toBe('A87654321')
      expect(taxIdTypeSelect.value).toBe('CIF')
      expect(websiteInput.value).toBe('https://test.com')
    })
  })

  describe('PERSON entity specific fields', () => {
    it('should show first and last name fields for PERSON', async () => {
      render(PartyForm)
      
      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'PERSON')
      
      expect(screen.getByLabelText(/Nombre \*/i)).toBeInTheDocument()
      expect(screen.getByLabelText(/Apellidos \*/i)).toBeInTheDocument()
      expect(screen.queryByLabelText(/Nombre de la organización/i)).not.toBeInTheDocument()
    })

    it('should populate first and last name from name if not provided', async () => {
      render(PartyForm, {
        props: {
          initialData: {
            name: 'Juan Perez Nadal',
            entityType: 'PERSON',
            role: 'CLIENT'
          }
        }
      })

      await waitFor(() => {
        expect(screen.getByLabelText(/Nombre \*/i)).toHaveValue('Juan')
        expect(screen.getByLabelText(/Apellidos \*/i)).toHaveValue('Perez Nadal')
      })
    })
  })

  describe('Form validation', () => {
    it('should show validation error on submit when name is missing', async () => {
      const { container } = render(PartyForm)

      // Need to select role and entityType so validateForm can be reached
      await fireEvent.update(screen.getByLabelText(/Rol de la entidad/i), 'CLIENT')
      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'ORGANIZATION')

      // Submit the form
      const form = container.querySelector('form')
      expect(form).toBeTruthy()
      await fireEvent.submit(form!)

      // Should NOT call createParty when validation fails
      expect(partyApi.createParty).not.toHaveBeenCalled()
      expect(toastStore.error).toHaveBeenCalledWith('Corrige los errores antes de continuar')
    })

    it('should validate name field on blur', async () => {
      const { container } = render(PartyForm)

      // Select entity type to reveal name field
      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'ORGANIZATION')

      const nameInput = screen.getByLabelText(/Nombre de la organización/i)
      await fireEvent.update(nameInput, 'AB') // Too short
      await fireEvent.blur(nameInput)

      // Wait for Vue to update DOM
      await new Promise(resolve => setTimeout(resolve, 50))

      // Check if error appears in the DOM
      expect(container.textContent).toContain('caracteres')
    })

    it('should validate taxId length on blur', async () => {
      const { container } = render(PartyForm)

      // Select entity type so the full form is visible
      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'ORGANIZATION')

      const taxIdInput = screen.getByPlaceholderText('p. ej., 12345678A')
      await fireEvent.update(taxIdInput, '123') // Too short
      await fireEvent.blur(taxIdInput)

      await new Promise(resolve => setTimeout(resolve, 50))

      expect(container.textContent).toContain('Formato inválido')
    })
  })

  describe('Form submission - Create mode', () => {
    it('should call createParty with valid data', async () => {
      const mockParty = {
        id: 'party-001',
        name: 'New Company',
        role: 'CLIENT',
        status: 'ACTIVE',
        tax_id: 'B12345678',
        tax_id_type: 'CIF',
        website: 'https://newcompany.com',
        created_at: '2026-02-17T10:00:00Z',
        modified_at: '2026-02-17T10:00:00Z',
        has_organization: true,
        has_person: false,
      }

      vi.mocked(partyApi.createParty).mockResolvedValueOnce(mockParty)

      const { container, emitted } = render(PartyForm)

      // Fill form
      await fireEvent.update(screen.getByLabelText(/Rol de la entidad/i), 'CLIENT')
      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'ORGANIZATION')
      
      // Wait for dynamic fields
      await waitFor(() => {
        expect(screen.getByLabelText(/Nombre de la organización/i)).toBeInTheDocument()
      })

      await fireEvent.update(screen.getByLabelText(/Nombre de la organización/i), 'New Company')
      await fireEvent.update(screen.getByPlaceholderText('p. ej., 12345678A'), '12345678Z')
      await fireEvent.update(screen.getByLabelText(/Sitio web/i), 'https://newcompany.com')

      // Submit
      const form = container.querySelector('form')
      await fireEvent.submit(form!)

      await waitFor(() => {
        expect(partyApi.createParty).toHaveBeenCalledWith(
          expect.objectContaining({
            name: 'New Company',
            role: 'CLIENT',
            taxId: '12345678Z',
          })
        )
      })

      expect(toastStore.success).toHaveBeenCalledWith('Entidad creada correctamente')
      expect(emitted().submit).toBeTruthy()
    })
  })

  describe('Form reset', () => {
    it('should reset all fields when reset button is clicked', async () => {
      render(PartyForm)

      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'ORGANIZATION')
      await fireEvent.update(screen.getByLabelText(/Rol de la entidad/i), 'CLIENT')

      await waitFor(() => {
        expect(screen.getByLabelText(/Nombre de la organización/i)).toBeInTheDocument()
      })

      const nameInput = screen.getByLabelText(/Nombre de la organización/i) as HTMLInputElement
      const taxIdInput = screen.getByPlaceholderText('p. ej., 12345678A') as HTMLInputElement

      await fireEvent.update(nameInput, 'Test Company')
      await fireEvent.update(taxIdInput, 'B12345678')

      // Reset
      await fireEvent.click(screen.getByRole('button', { name: /Reiniciar/i }))

      expect((screen.getByLabelText(/Rol de la entidad/i) as HTMLSelectElement).value).toBe('')
      // Name field is hidden after reset
      expect(screen.queryByLabelText(/Nombre de la organización/i)).not.toBeInTheDocument()
    })
  })

  describe('Tax ID type selection', () => {
    it('should default to CIF for ORGANIZATION', async () => {
      render(PartyForm)
      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'ORGANIZATION')
      
      const taxIdTypeSelect = screen.getByLabelText(/Tipo identificación/i) as HTMLSelectElement
      expect(taxIdTypeSelect.value).toBe('CIF')
    })

    it('should default to NIF for PERSON', async () => {
      render(PartyForm)
      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'PERSON')
      
      const taxIdTypeSelect = screen.getByLabelText(/Tipo identificación/i) as HTMLSelectElement
      expect(taxIdTypeSelect.value).toBe('NIF')
    })
  })
})
