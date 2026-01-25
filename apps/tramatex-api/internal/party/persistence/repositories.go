package persistence

import (
	"context"

	"github.com/joran-cortez/tramatex/internal/party/domain"
)

// OrganizationRepository defines methods for organization persistence
type OrganizationRepository interface {
	// Save persists a new organization or updates an existing one
	Save(ctx context.Context, org *domain.Organization) error

	// FindByID retrieves an organization by its ID
	FindByID(ctx context.Context, id domain.OrganizationID) (*domain.Organization, error)

	// FindAll retrieves all organizations with optional filtering
	FindAll(ctx context.Context, filters *OrganizationFilters) ([]*domain.Organization, error)

	// FindByRole retrieves organizations by role
	FindByRole(ctx context.Context, role domain.OrganizationRole) ([]*domain.Organization, error)

	// Delete marks an organization as deleted (soft delete or removes it)
	Delete(ctx context.Context, id domain.OrganizationID) error

	// Exists checks if an organization exists
	Exists(ctx context.Context, id domain.OrganizationID) (bool, error)

	// Count returns the total number of organizations
	Count(ctx context.Context) (int64, error)
}

// PersonRepository defines methods for person persistence
type PersonRepository interface {
	// Save persists a new person or updates an existing one
	Save(ctx context.Context, person *domain.Person) error

	// FindByID retrieves a person by ID
	FindByID(ctx context.Context, id domain.PersonID) (*domain.Person, error)

	// FindByOrganization retrieves all persons for an organization
	FindByOrganization(ctx context.Context, orgID domain.OrganizationID) ([]*domain.Person, error)

	// FindByEmail retrieves a person by email
	FindByEmail(ctx context.Context, email string) (*domain.Person, error)

	// FindPrimaryContact retrieves the primary contact for an organization
	FindPrimaryContact(ctx context.Context, orgID domain.OrganizationID) (*domain.Person, error)

	// Delete removes a person
	Delete(ctx context.Context, id domain.PersonID) error

	// Exists checks if a person exists
	Exists(ctx context.Context, id domain.PersonID) (bool, error)
}

// AddressRepository defines methods for address persistence
type AddressRepository interface {
	// Save persists a new address or updates an existing one
	Save(ctx context.Context, address *domain.Address, addressID domain.AddressID, orgID domain.OrganizationID) error

	// FindByID retrieves an address by ID
	FindByID(ctx context.Context, id domain.AddressID) (*domain.Address, error)

	// FindByOrganization retrieves all addresses for an organization
	FindByOrganization(ctx context.Context, orgID domain.OrganizationID) ([]*domain.Address, error)

	// FindPrimary retrieves the primary address for an organization
	FindPrimary(ctx context.Context, orgID domain.OrganizationID) (*domain.Address, error)

	// Delete removes an address
	Delete(ctx context.Context, id domain.AddressID) error

	// Exists checks if an address exists
	Exists(ctx context.Context, id domain.AddressID) (bool, error)
}

// OrganizationFilters defines filtering options for organizations
type OrganizationFilters struct {
	Status     *domain.OrganizationStatus
	Role       *domain.OrganizationRole
	Name       string
	PageSize   int
	PageNumber int
}
