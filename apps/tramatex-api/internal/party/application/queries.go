package application

import (
	"context"
	"fmt"

	"github.com/joran-cortez/tramatex/internal/party/domain"
	"github.com/joran-cortez/tramatex/internal/party/persistence"
)

// GetOrganizationQuery represents a query to get an organization by ID
type GetOrganizationQuery struct {
	ID string
}

// GetOrganizationHandler handles fetching a single organization
type GetOrganizationHandler struct {
	orgRepo persistence.OrganizationRepository
}

// NewGetOrganizationHandler creates a new handler
func NewGetOrganizationHandler(orgRepo persistence.OrganizationRepository) *GetOrganizationHandler {
	return &GetOrganizationHandler{orgRepo: orgRepo}
}

// Handle executes the get organization query
func (h *GetOrganizationHandler) Handle(ctx context.Context, query *GetOrganizationQuery) (*domain.Organization, error) {
	if query.ID == "" {
		return nil, fmt.Errorf("organization ID cannot be empty")
	}

	orgID, err := domain.NewOrganizationID(query.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	org, err := h.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch organization: %w", err)
	}

	return org, nil
}

// ListOrganizationsQuery represents a query to list organizations
type ListOrganizationsQuery struct {
	Status     string // ACTIVE, INACTIVE, or empty for all
	Role       string // CLIENT, SUPPLIER, BOTH, or empty for all
	Name       string // Partial match
	PageSize   int
	PageNumber int
}

// ListOrganizationsHandler handles listing organizations
type ListOrganizationsHandler struct {
	orgRepo persistence.OrganizationRepository
}

// NewListOrganizationsHandler creates a new handler
func NewListOrganizationsHandler(orgRepo persistence.OrganizationRepository) *ListOrganizationsHandler {
	return &ListOrganizationsHandler{orgRepo: orgRepo}
}

// Handle executes the list organizations query
func (h *ListOrganizationsHandler) Handle(ctx context.Context, query *ListOrganizationsQuery) ([]*domain.Organization, error) {
	// Build filters
	filters := &persistence.OrganizationFilters{
		Name:       query.Name,
		PageNumber: query.PageNumber,
		PageSize:   query.PageSize,
	}

	// Set default page size if not provided
	if filters.PageSize <= 0 {
		filters.PageSize = 10
	}
	if filters.PageNumber <= 0 {
		filters.PageNumber = 1
	}

	// Parse status if provided
	if query.Status != "" {
		status := domain.ParseOrganizationStatus(query.Status)
		if !status.IsValid() {
			return nil, fmt.Errorf("invalid status: %s", query.Status)
		}
		filters.Status = &status
	}

	// Parse role if provided
	if query.Role != "" {
		role := domain.ParseOrganizationRole(query.Role)
		if !role.IsValid() {
			return nil, fmt.Errorf("invalid role: %s", query.Role)
		}
		filters.Role = &role
	}

	orgs, err := h.orgRepo.FindAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	return orgs, nil
}

// ListOrganizationsByRoleQuery represents a query to list organizations by role
type ListOrganizationsByRoleQuery struct {
	Role string // CLIENT, SUPPLIER, BOTH
}

// ListOrganizationsByRoleHandler handles listing by role
type ListOrganizationsByRoleHandler struct {
	orgRepo persistence.OrganizationRepository
}

// NewListOrganizationsByRoleHandler creates a new handler
func NewListOrganizationsByRoleHandler(orgRepo persistence.OrganizationRepository) *ListOrganizationsByRoleHandler {
	return &ListOrganizationsByRoleHandler{orgRepo: orgRepo}
}

// Handle executes the list by role query
func (h *ListOrganizationsByRoleHandler) Handle(ctx context.Context, query *ListOrganizationsByRoleQuery) ([]*domain.Organization, error) {
	if query.Role == "" {
		return nil, fmt.Errorf("role cannot be empty")
	}

	role := domain.ParseOrganizationRole(query.Role)
	if !role.IsValid() {
		return nil, fmt.Errorf("invalid role: %s", query.Role)
	}

	orgs, err := h.orgRepo.FindByRole(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	return orgs, nil
}

// GetPersonQuery represents a query to get a person by ID
type GetPersonQuery struct {
	ID string
}

// GetPersonHandler handles fetching a single person
type GetPersonHandler struct {
	personRepo persistence.PersonRepository
}

// NewGetPersonHandler creates a new handler
func NewGetPersonHandler(personRepo persistence.PersonRepository) *GetPersonHandler {
	return &GetPersonHandler{personRepo: personRepo}
}

// Handle executes the get person query
func (h *GetPersonHandler) Handle(ctx context.Context, query *GetPersonQuery) (*domain.Person, error) {
	if query.ID == "" {
		return nil, fmt.Errorf("person ID cannot be empty")
	}

	personID, err := domain.NewPersonID(query.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid person ID: %w", err)
	}

	person, err := h.personRepo.FindByID(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch person: %w", err)
	}

	return person, nil
}

// ListPersonsByOrganizationQuery represents a query to list persons in an organization
type ListPersonsByOrganizationQuery struct {
	OrganizationID string
}

// ListPersonsByOrganizationHandler handles listing persons
type ListPersonsByOrganizationHandler struct {
	personRepo persistence.PersonRepository
}

// NewListPersonsByOrganizationHandler creates a new handler
func NewListPersonsByOrganizationHandler(personRepo persistence.PersonRepository) *ListPersonsByOrganizationHandler {
	return &ListPersonsByOrganizationHandler{personRepo: personRepo}
}

// Handle executes the list persons query
func (h *ListPersonsByOrganizationHandler) Handle(ctx context.Context, query *ListPersonsByOrganizationQuery) ([]*domain.Person, error) {
	if query.OrganizationID == "" {
		return nil, fmt.Errorf("organization ID cannot be empty")
	}

	orgID, err := domain.NewOrganizationID(query.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	persons, err := h.personRepo.FindByOrganization(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list persons: %w", err)
	}

	return persons, nil
}

// GetPersonByEmailQuery represents a query to get a person by email
type GetPersonByEmailQuery struct {
	Email string
}

// GetPersonByEmailHandler handles fetching by email
type GetPersonByEmailHandler struct {
	personRepo persistence.PersonRepository
}

// NewGetPersonByEmailHandler creates a new handler
func NewGetPersonByEmailHandler(personRepo persistence.PersonRepository) *GetPersonByEmailHandler {
	return &GetPersonByEmailHandler{personRepo: personRepo}
}

// Handle executes the get by email query
func (h *GetPersonByEmailHandler) Handle(ctx context.Context, query *GetPersonByEmailQuery) (*domain.Person, error) {
	if query.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	person, err := h.personRepo.FindByEmail(ctx, query.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch person: %w", err)
	}

	return person, nil
}

// GetPrimaryContactQuery represents a query to get primary contact
type GetPrimaryContactQuery struct {
	OrganizationID string
}

// GetPrimaryContactHandler handles fetching primary contact
type GetPrimaryContactHandler struct {
	personRepo persistence.PersonRepository
}

// NewGetPrimaryContactHandler creates a new handler
func NewGetPrimaryContactHandler(personRepo persistence.PersonRepository) *GetPrimaryContactHandler {
	return &GetPrimaryContactHandler{personRepo: personRepo}
}

// Handle executes the get primary contact query
func (h *GetPrimaryContactHandler) Handle(ctx context.Context, query *GetPrimaryContactQuery) (*domain.Person, error) {
	if query.OrganizationID == "" {
		return nil, fmt.Errorf("organization ID cannot be empty")
	}

	orgID, err := domain.NewOrganizationID(query.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	person, err := h.personRepo.FindPrimaryContact(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch primary contact: %w", err)
	}

	return person, nil
}

// ListAddressesByOrganizationQuery represents a query to list addresses
type ListAddressesByOrganizationQuery struct {
	OrganizationID string
}

// ListAddressesByOrganizationHandler handles listing addresses
type ListAddressesByOrganizationHandler struct {
	addressRepo persistence.AddressRepository
}

// NewListAddressesByOrganizationHandler creates a new handler
func NewListAddressesByOrganizationHandler(addressRepo persistence.AddressRepository) *ListAddressesByOrganizationHandler {
	return &ListAddressesByOrganizationHandler{addressRepo: addressRepo}
}

// Handle executes the list addresses query
func (h *ListAddressesByOrganizationHandler) Handle(ctx context.Context, query *ListAddressesByOrganizationQuery) ([]*domain.Address, error) {
	if query.OrganizationID == "" {
		return nil, fmt.Errorf("organization ID cannot be empty")
	}

	orgID, err := domain.NewOrganizationID(query.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	addresses, err := h.addressRepo.FindByOrganization(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list addresses: %w", err)
	}

	return addresses, nil
}

// GetPrimaryAddressQuery represents a query to get primary address
type GetPrimaryAddressQuery struct {
	OrganizationID string
}

// GetPrimaryAddressHandler handles fetching primary address
type GetPrimaryAddressHandler struct {
	addressRepo persistence.AddressRepository
}

// NewGetPrimaryAddressHandler creates a new handler
func NewGetPrimaryAddressHandler(addressRepo persistence.AddressRepository) *GetPrimaryAddressHandler {
	return &GetPrimaryAddressHandler{addressRepo: addressRepo}
}

// Handle executes the get primary address query
func (h *GetPrimaryAddressHandler) Handle(ctx context.Context, query *GetPrimaryAddressQuery) (*domain.Address, error) {
	if query.OrganizationID == "" {
		return nil, fmt.Errorf("organization ID cannot be empty")
	}

	orgID, err := domain.NewOrganizationID(query.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	address, err := h.addressRepo.FindPrimary(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch primary address: %w", err)
	}

	return address, nil
}
