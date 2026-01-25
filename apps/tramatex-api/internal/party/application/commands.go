package application

import (
	"context"
	"fmt"

	"github.com/joran-cortez/tramatex/internal/party/domain"
	"github.com/joran-cortez/tramatex/internal/party/persistence"
)

// CreateOrganizationCommand represents a command to create an organization
type CreateOrganizationCommand struct {
	ID        string
	Name      string
	Role      string // CLIENT, SUPPLIER, BOTH
	TaxID     string
	TaxIDType string
	Website   string
	CreatedBy string
}

// CreateOrganizationHandler handles organization creation
type CreateOrganizationHandler struct {
	orgRepo persistence.OrganizationRepository
}

// NewCreateOrganizationHandler creates a new handler
func NewCreateOrganizationHandler(orgRepo persistence.OrganizationRepository) *CreateOrganizationHandler {
	return &CreateOrganizationHandler{orgRepo: orgRepo}
}

// Handle executes the create organization command
func (h *CreateOrganizationHandler) Handle(ctx context.Context, cmd *CreateOrganizationCommand) (*domain.Organization, error) {
	// Validate inputs
	if cmd.ID == "" {
		return nil, fmt.Errorf("organization ID cannot be empty")
	}
	if cmd.Name == "" {
		return nil, fmt.Errorf("organization name cannot be empty")
	}
	if cmd.CreatedBy == "" {
		return nil, fmt.Errorf("createdBy user ID cannot be empty")
	}

	// Create ID
	orgID, err := domain.NewOrganizationID(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	// Parse role
	role := domain.ParseOrganizationRole(cmd.Role)
	if !role.IsValid() {
		return nil, fmt.Errorf("invalid organization role: %s", cmd.Role)
	}

	// Create Tax ID if provided
	var taxID *domain.TaxID
	if cmd.TaxID != "" {
		taxIDType := cmd.TaxIDType
		if taxIDType == "" {
			taxIDType = "NIF" // Default type
		}
		taxID, err = domain.NewTaxID(cmd.TaxID, taxIDType)
		if err != nil {
			return nil, fmt.Errorf("invalid tax ID: %w", err)
		}
	}

	// Create organization aggregate
	org, err := domain.NewOrganization(orgID, cmd.Name, role, taxID, cmd.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Set website if provided
	if cmd.Website != "" {
		if err := org.UpdateWebsite(cmd.Website, cmd.CreatedBy); err != nil {
			return nil, fmt.Errorf("failed to set website: %w", err)
		}
	}

	// Save to repository
	if err := h.orgRepo.Save(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to save organization: %w", err)
	}

	return org, nil
}

// UpdateOrganizationCommand represents a command to update an organization
type UpdateOrganizationCommand struct {
	ID         string
	Name       string
	Website    string
	Notes      string
	ModifiedBy string
}

// UpdateOrganizationHandler handles organization updates
type UpdateOrganizationHandler struct {
	orgRepo persistence.OrganizationRepository
}

// NewUpdateOrganizationHandler creates a new handler
func NewUpdateOrganizationHandler(orgRepo persistence.OrganizationRepository) *UpdateOrganizationHandler {
	return &UpdateOrganizationHandler{orgRepo: orgRepo}
}

// Handle executes the update organization command
func (h *UpdateOrganizationHandler) Handle(ctx context.Context, cmd *UpdateOrganizationCommand) (*domain.Organization, error) {
	// Validate inputs
	if cmd.ID == "" {
		return nil, fmt.Errorf("organization ID cannot be empty")
	}
	if cmd.ModifiedBy == "" {
		return nil, fmt.Errorf("modifiedBy user ID cannot be empty")
	}

	// Find organization
	orgID, err := domain.NewOrganizationID(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	org, err := h.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	// Update fields
	if cmd.Name != "" {
		if err := org.UpdateName(cmd.Name, cmd.ModifiedBy); err != nil {
			return nil, fmt.Errorf("failed to update name: %w", err)
		}
	}

	if cmd.Website != "" {
		if err := org.UpdateWebsite(cmd.Website, cmd.ModifiedBy); err != nil {
			return nil, fmt.Errorf("failed to update website: %w", err)
		}
	}

	if cmd.Notes != "" {
		if err := org.UpdateNotes(cmd.Notes, cmd.ModifiedBy); err != nil {
			return nil, fmt.Errorf("failed to update notes: %w", err)
		}
	}

	// Save changes
	if err := h.orgRepo.Save(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to save changes: %w", err)
	}

	return org, nil
}

// ChangeOrganizationStatusCommand represents a command to change organization status
type ChangeOrganizationStatusCommand struct {
	ID         string
	Status     string // ACTIVE, INACTIVE
	ModifiedBy string
}

// ChangeOrganizationStatusHandler handles status changes
type ChangeOrganizationStatusHandler struct {
	orgRepo persistence.OrganizationRepository
}

// NewChangeOrganizationStatusHandler creates a new handler
func NewChangeOrganizationStatusHandler(orgRepo persistence.OrganizationRepository) *ChangeOrganizationStatusHandler {
	return &ChangeOrganizationStatusHandler{orgRepo: orgRepo}
}

// Handle executes the change status command
func (h *ChangeOrganizationStatusHandler) Handle(ctx context.Context, cmd *ChangeOrganizationStatusCommand) (*domain.Organization, error) {
	// Validate inputs
	if cmd.ID == "" {
		return nil, fmt.Errorf("organization ID cannot be empty")
	}
	if cmd.ModifiedBy == "" {
		return nil, fmt.Errorf("modifiedBy user ID cannot be empty")
	}

	// Find organization
	orgID, err := domain.NewOrganizationID(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	org, err := h.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	// Parse status
	status := domain.ParseOrganizationStatus(cmd.Status)
	if !status.IsValid() {
		return nil, fmt.Errorf("invalid status: %s", cmd.Status)
	}

	// Change status
	if status == domain.OrganizationStatusActive {
		if err := org.Activate(cmd.ModifiedBy); err != nil {
			return nil, fmt.Errorf("failed to activate: %w", err)
		}
	} else {
		if err := org.Deactivate(cmd.ModifiedBy); err != nil {
			return nil, fmt.Errorf("failed to deactivate: %w", err)
		}
	}

	// Save changes
	if err := h.orgRepo.Save(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to save changes: %w", err)
	}

	return org, nil
}

// AddPersonCommand represents a command to add a person to an organization
type AddPersonCommand struct {
	ID             string
	OrganizationID string
	FirstName      string
	LastName       string
	Email          string
	Phone          string
	JobTitle       string
	IsPrimary      bool
	CreatedBy      string
}

// AddPersonHandler handles adding a person
type AddPersonHandler struct {
	orgRepo    persistence.OrganizationRepository
	personRepo persistence.PersonRepository
}

// NewAddPersonHandler creates a new handler
func NewAddPersonHandler(
	orgRepo persistence.OrganizationRepository,
	personRepo persistence.PersonRepository,
) *AddPersonHandler {
	return &AddPersonHandler{
		orgRepo:    orgRepo,
		personRepo: personRepo,
	}
}

// Handle executes the add person command
func (h *AddPersonHandler) Handle(ctx context.Context, cmd *AddPersonCommand) (*domain.Person, error) {
	// Validate inputs
	if cmd.ID == "" {
		return nil, fmt.Errorf("person ID cannot be empty")
	}
	if cmd.OrganizationID == "" {
		return nil, fmt.Errorf("organization ID cannot be empty")
	}
	if cmd.FirstName == "" || cmd.LastName == "" {
		return nil, fmt.Errorf("first name and last name cannot be empty")
	}
	if cmd.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if cmd.CreatedBy == "" {
		return nil, fmt.Errorf("createdBy user ID cannot be empty")
	}

	// Verify organization exists
	orgID, err := domain.NewOrganizationID(cmd.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	exists, err := h.orgRepo.Exists(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to check organization: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("organization not found")
	}

	// Create person ID
	personID, err := domain.NewPersonID(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid person ID: %w", err)
	}

	// Create email
	email, err := domain.NewEmail(cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	// Create person
	person := domain.NewPerson(personID, orgID, cmd.FirstName, cmd.LastName, email, cmd.CreatedBy)

	// Set phone if provided
	if cmd.Phone != "" {
		phone, err := domain.NewPhone(cmd.Phone)
		if err != nil {
			return nil, fmt.Errorf("invalid phone: %w", err)
		}
		person.SetPhone(phone)
	}

	// Set job title if provided
	if cmd.JobTitle != "" {
		person.SetJobTitle(cmd.JobTitle)
	}

	// Set primary contact flag
	if cmd.IsPrimary {
		person.SetPrimaryContact(true)
	}

	// Save person
	if err := h.personRepo.Save(ctx, person); err != nil {
		return nil, fmt.Errorf("failed to save person: %w", err)
	}

	return person, nil
}

// AddAddressCommand represents a command to add an address to an organization
type AddAddressCommand struct {
	ID             string
	OrganizationID string
	Street         string
	City           string
	Province       string
	PostalCode     string
	Country        string
	IsPrimary      bool
	CreatedBy      string
}

// AddAddressHandler handles adding an address
type AddAddressHandler struct {
	orgRepo     persistence.OrganizationRepository
	addressRepo persistence.AddressRepository
}

// NewAddAddressHandler creates a new handler
func NewAddAddressHandler(
	orgRepo persistence.OrganizationRepository,
	addressRepo persistence.AddressRepository,
) *AddAddressHandler {
	return &AddAddressHandler{
		orgRepo:     orgRepo,
		addressRepo: addressRepo,
	}
}

// Handle executes the add address command
func (h *AddAddressHandler) Handle(ctx context.Context, cmd *AddAddressCommand) (*domain.Address, error) {
	// Validate inputs
	if cmd.ID == "" {
		return nil, fmt.Errorf("address ID cannot be empty")
	}
	if cmd.OrganizationID == "" {
		return nil, fmt.Errorf("organization ID cannot be empty")
	}
	if cmd.Street == "" || cmd.City == "" {
		return nil, fmt.Errorf("street and city cannot be empty")
	}
	if cmd.CreatedBy == "" {
		return nil, fmt.Errorf("createdBy user ID cannot be empty")
	}

	// Verify organization exists
	orgID, err := domain.NewOrganizationID(cmd.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	exists, err := h.orgRepo.Exists(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to check organization: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("organization not found")
	}

	// Create address ID
	addressID, err := domain.NewAddressID(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid address ID: %w", err)
	}

	// Create address
	addr, err := domain.NewAddress(cmd.Street, cmd.City, cmd.Province, cmd.PostalCode, cmd.Country)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	// Save address
	if err := h.addressRepo.Save(ctx, addr, addressID, orgID); err != nil {
		return nil, fmt.Errorf("failed to save address: %w", err)
	}

	return addr, nil
}
