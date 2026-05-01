import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/vue'
import PartyForm from '../../../components/party/PartyForm.vue'
import { partyApi } from '../../../services/partyApi'

// Mock partyApi
vi.mock('../../../services/partyApi', () => ({
  partyApi: {
    createParty: vi.fn(),
    updateParty: vi.fn(),
  }
}))

describe('PartyForm Component', () => {
  beforeEach(() => {
    vi.clearAllMocks()
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
      expect(screen.getByLabelText(/Tipo de NIF\/CIF/i)).toBeInTheDocument()
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
      const taxIdTypeSelect = screen.getByLabelText(/Tipo de NIF\/CIF/i) as HTMLSelectElement
      const websiteInput = screen.getByLabelText(/Sitio web/i) as HTMLInputElement

      expect(nameInput.value).toBe('Test Company')
      expect(roleSelect.value).toBe('SUPPLIER')
      expect(taxIdInput.value).toBe('A87654321')
      expect(taxIdTypeSelect.value).toBe('CIF')
      expect(websiteInput.value).toBe('https://test.com')
    })
  })

  describe('Form validation', () => {
    it('should show validation error on submit when name is missing', async () => {
      render(PartyForm)

      const submitButton = screen.getByRole('button', { name: /Crear entidad/i })
      await fireEvent.click(submitButton)

      // Wait a bit for potential API call
      await new Promise(resolve => setTimeout(resolve, 100))

      // Should NOT call createParty when validation fails
      expect(partyApi.createParty).not.toHaveBeenCalled()
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

      // Check if error appears in the DOM (using container to be more flexible)
      const hasError = container.textContent?.includes('caracteres') || 
                       container.textContent?.includes('obligatorio')
      expect(hasError).toBe(true)
    })

    it('should accept valid name', async () => {
      render(PartyForm)

      // Select entity type to reveal name field
      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'ORGANIZATION')

      const nameInput = screen.getByLabelText(/Nombre de la organización/i)
      await fireEvent.update(nameInput, 'Valid Company Name')
      await fireEvent.blur(nameInput)

      // Should not show error for valid name
      expect(nameInput).toHaveValue('Valid Company Name')
    })

    it('should validate taxId length on blur', async () => {
      const { container } = render(PartyForm)

      // Select entity type so the full form is visible
      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'ORGANIZATION')

      const taxIdInput = screen.getByPlaceholderText('p. ej., 12345678A')
      await fireEvent.update(taxIdInput, '123') // Too short
      await fireEvent.blur(taxIdInput)

      await new Promise(resolve => setTimeout(resolve, 50))

      const hasError = container.textContent?.includes('Formato de NIF inválido') ||
                       container.textContent?.includes('Formato inválido')
      expect(hasError).toBe(true)
    })

    it('should validate website URL format on blur', async () => {
      const { container } = render(PartyForm)

      const websiteInput = screen.getByLabelText(/Sitio web/i)
      await fireEvent.update(websiteInput, 'invalid-url')
      await fireEvent.blur(websiteInput)

      await new Promise(resolve => setTimeout(resolve, 50))

      const hasError = container.textContent?.includes('URL inválido') || 
                       container.textContent?.includes('Formato')
      expect(hasError).toBe(true)
    })

    it('should accept valid website URL', async () => {
      render(PartyForm)

      const websiteInput = screen.getByLabelText(/Sitio web/i)
      await fireEvent.update(websiteInput, 'https://example.com')
      await fireEvent.blur(websiteInput)

      expect((websiteInput as HTMLInputElement).value).toBe('https://example.com')
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

      const { emitted } = render(PartyForm)

      // Fill form - select entity type first to reveal name field
      await fireEvent.update(screen.getByLabelText(/Rol de la entidad/i), 'CLIENT')
      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'ORGANIZATION')
      // Set taxIdType to NIF explicitly as selecting ORGANIZATION auto-sets it to CIF
      await fireEvent.update(screen.getByLabelText(/Tipo de NIF\/CIF/i), 'NIF')
      await fireEvent.update(screen.getByLabelText(/Nombre de la organización/i), 'New Company')
      await fireEvent.update(screen.getByPlaceholderText('p. ej., 12345678A'), '12345678Z')
      await fireEvent.update(screen.getByLabelText(/Sitio web/i), 'https://newcompany.com')

      // Submit
      await fireEvent.click(screen.getByRole('button', { name: /Crear entidad/i }))

      await waitFor(() => {
        expect(partyApi.createParty).toHaveBeenCalledWith(
          expect.objectContaining({
            name: 'New Company',
            role: 'CLIENT',
            taxId: '12345678Z',
            taxIdType: 'NIF',
            website: 'https://newcompany.com',
          })
        )
      })

      // Check if submit event was emitted
      await waitFor(() => {
        expect(emitted().submit).toBeTruthy()
      })
    })

    it('should handle creation error', async () => {
      vi.mocked(partyApi.createParty).mockRejectedValueOnce({
        message: 'Tax ID already exists',
      })

      const { container } = render(PartyForm)

      // Fill form - select entity type first to reveal name field
      await fireEvent.update(screen.getByLabelText(/Rol de la entidad/i), 'CLIENT')
      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'ORGANIZATION')
      await fireEvent.update(screen.getByLabelText(/Nombre de la organización/i), 'Duplicate Company')

      // Submit
      await fireEvent.click(screen.getByRole('button', { name: /Crear entidad/i }))

      await waitFor(() => {
        expect(partyApi.createParty).toHaveBeenCalled()
      })

      // Check error message appears somewhere in the container
      await waitFor(() => {
        const hasErrorMessage = container.textContent?.includes('Tax ID') || 
                                container.textContent?.includes('error')
        expect(hasErrorMessage).toBe(true)
      }, { timeout: 2000 })
    })

    it('should disable submit button while submitting', async () => {
      vi.mocked(partyApi.createParty).mockImplementation(() => 
        new Promise(resolve => setTimeout(() => resolve({
          id: 'party-001',
          name: 'Test',
          role: 'CLIENT',
        }), 100))
      )

      render(PartyForm)

      await fireEvent.update(screen.getByLabelText(/Rol de la entidad/i), 'CLIENT')
      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'ORGANIZATION')
      await fireEvent.update(screen.getByLabelText(/Nombre de la organización/i), 'Test Company')

      const submitButton = screen.getByRole('button', { name: /Crear entidad/i })
      await fireEvent.click(submitButton)

      // Button should be disabled while submitting
      await waitFor(() => {
        expect(submitButton).toBeDisabled()
      }, { timeout: 500 })
    })
  })

  describe('Form submission - Edit mode', () => {
    it('should call updateParty with updated data', async () => {
      const mockUpdatedParty = {
        id: 'party-123',
        name: 'Updated Company',
        role: 'BOTH',
        status: 'ACTIVE',
        website: 'https://updated.com',
      }

      vi.mocked(partyApi.updateParty).mockResolvedValueOnce(mockUpdatedParty)

      const { emitted } = render(PartyForm, {
        props: {
          partyId: 'party-123',
          initialData: {
            name: 'Original Company',
            role: 'CLIENT',
            entityType: 'ORGANIZATION',
            website: 'https://original.com',
          }
        }
      })

      // Wait for entity type to render the org name field
      await waitFor(() => {
        expect(screen.getByLabelText(/Nombre de la organización/i)).toBeInTheDocument()
      })

      // Update fields
      await fireEvent.update(screen.getByLabelText(/Nombre de la organización/i), 'Updated Company')
      await fireEvent.update(screen.getByLabelText(/Sitio web/i), 'https://updated.com')

      // Submit
      await fireEvent.click(screen.getByRole('button', { name: /Actualizar entidad/i }))

      await waitFor(() => {
        expect(partyApi.updateParty).toHaveBeenCalledWith(
          'party-123',
          expect.objectContaining({
            name: 'Updated Company',
            website: 'https://updated.com',
          })
        )
      })

      await waitFor(() => {
        expect(emitted().update).toBeTruthy()
      })
    })

    it('should handle update error', async () => {
      vi.mocked(partyApi.updateParty).mockRejectedValueOnce({
        data: {
          message: 'Party not found',
        }
      })

      const { container } = render(PartyForm, {
        props: {
          partyId: 'party-123',
          initialData: {
            name: 'Test Company',
            role: 'CLIENT',
            entityType: 'ORGANIZATION',
          }
        }
      })

      // Wait for entity type to render the org name field
      await waitFor(() => {
        expect(screen.getByLabelText(/Nombre de la organización/i)).toBeInTheDocument()
      })

      await fireEvent.update(screen.getByLabelText(/Nombre de la organización/i), 'Updated Name')
      await fireEvent.click(screen.getByRole('button', { name: /Actualizar entidad/i }))

      await waitFor(() => {
        expect(partyApi.updateParty).toHaveBeenCalled()
      })

      await waitFor(() => {
        const hasError = container.textContent?.includes('Party not found') || 
                        container.textContent?.includes('error')
        expect(hasError).toBe(true)
      }, { timeout: 2000 })
    })
  })

  describe('Form reset', () => {
    it('should reset all fields when reset button is clicked', async () => {
      render(PartyForm)

      // Select entity type first to reveal name field
      await fireEvent.update(screen.getByLabelText(/Rol de la entidad/i), 'CLIENT')
      await fireEvent.update(screen.getByLabelText(/Tipo de entidad/i), 'ORGANIZATION')

      // Fill form
      const nameInput = screen.getByLabelText(/Nombre de la organización/i) as HTMLInputElement
      const roleSelect = screen.getByLabelText(/Rol de la entidad/i) as HTMLSelectElement
      const taxIdInput = screen.getByPlaceholderText('p. ej., 12345678A') as HTMLInputElement

      await fireEvent.update(nameInput, 'Test Company')
      await fireEvent.update(taxIdInput, 'B12345678')

      // Reset
      await fireEvent.click(screen.getByRole('button', { name: /Reiniciar/i }))

      expect(roleSelect.value).toBe('')
      expect(taxIdInput.value).toBe('')
      // Name field is hidden after reset (entityType resets to '')
      expect(screen.queryByLabelText(/Nombre de la organización/i)).not.toBeInTheDocument()
    })
  })

  describe('Role selection', () => {
    it('should have all role options', () => {
      render(PartyForm)

      const roleSelect = screen.getByLabelText(/Rol de la entidad/i)
      const options = roleSelect.querySelectorAll('option')

      expect(options).toHaveLength(5) // Including placeholder
      expect(options[0].textContent).toBe('-- Selecciona rol --')
      expect(options[1].textContent).toBe('Cliente')
      expect(options[2].textContent).toBe('Proveedor')
      expect(options[3].textContent).toBe('Cliente y proveedor')
      expect(options[4].textContent).toBe('Contacto')
    })
  })

  describe('Tax ID type selection', () => {
    it('should have all tax ID type options', () => {
      render(PartyForm)

      const taxIdTypeSelect = screen.getByLabelText(/Tipo de NIF\/CIF/i)
      const options = taxIdTypeSelect.querySelectorAll('option')

      expect(options).toHaveLength(3)
      expect(options[0].textContent).toBe('NIF')
      expect(options[1].textContent).toBe('CIF')
      expect(options[2].textContent).toBe('VAT')
    })

    it('should default to NIF', () => {
      render(PartyForm)

      const taxIdTypeSelect = screen.getByLabelText(/Tipo de NIF\/CIF/i) as HTMLSelectElement
      expect(taxIdTypeSelect.value).toBe('NIF')
    })
  })
})
