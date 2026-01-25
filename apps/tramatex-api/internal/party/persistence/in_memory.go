package persistence

import (
	"context"
	"fmt"

	"github.com/joran-cortez/tramatex/internal/party/domain"
)

// InMemoryOrganizationRepository is an in-memory implementation of OrganizationRepository
type InMemoryOrganizationRepository struct {
	organizations map[string]*domain.Organization
}

// NewInMemoryOrganizationRepository creates a new in-memory organization repository
func NewInMemoryOrganizationRepository() *InMemoryOrganizationRepository {
	return &InMemoryOrganizationRepository{
		organizations: make(map[string]*domain.Organization),
	}
}

// Save persists an organization
func (r *InMemoryOrganizationRepository) Save(ctx context.Context, org *domain.Organization) error {
	if org == nil {
		return fmt.Errorf("organization cannot be nil")
	}
	r.organizations[org.ID().String()] = org
	return nil
}

// FindByID retrieves an organization by ID
func (r *InMemoryOrganizationRepository) FindByID(ctx context.Context, id domain.OrganizationID) (*domain.Organization, error) {
	org, exists := r.organizations[id.String()]
	if !exists {
		return nil, fmt.Errorf("organization not found: %s", id)
	}
	return org, nil
}

// FindAll retrieves all organizations
func (r *InMemoryOrganizationRepository) FindAll(ctx context.Context, filters *OrganizationFilters) ([]*domain.Organization, error) {
	var result []*domain.Organization

	for _, org := range r.organizations {
		// Apply filters
		if filters != nil {
			if filters.Status != nil && org.Status() != *filters.Status {
				continue
			}
			if filters.Role != nil && org.Role() != *filters.Role {
				continue
			}
			if filters.Name != "" && org.Name() != filters.Name {
				continue
			}
		}
		result = append(result, org)
	}

	return result, nil
}

// FindByRole retrieves organizations by role
func (r *InMemoryOrganizationRepository) FindByRole(ctx context.Context, role domain.OrganizationRole) ([]*domain.Organization, error) {
	var result []*domain.Organization
	for _, org := range r.organizations {
		if org.Role() == role {
			result = append(result, org)
		}
	}
	return result, nil
}

// Delete removes an organization
func (r *InMemoryOrganizationRepository) Delete(ctx context.Context, id domain.OrganizationID) error {
	if _, exists := r.organizations[id.String()]; !exists {
		return fmt.Errorf("organization not found: %s", id)
	}
	delete(r.organizations, id.String())
	return nil
}

// Exists checks if an organization exists
func (r *InMemoryOrganizationRepository) Exists(ctx context.Context, id domain.OrganizationID) (bool, error) {
	_, exists := r.organizations[id.String()]
	return exists, nil
}

// Count returns the total number of organizations
func (r *InMemoryOrganizationRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(r.organizations)), nil
}

// InMemoryPersonRepository is an in-memory implementation of PersonRepository
type InMemoryPersonRepository struct {
	persons map[string]*domain.Person
}

// NewInMemoryPersonRepository creates a new in-memory person repository
func NewInMemoryPersonRepository() *InMemoryPersonRepository {
	return &InMemoryPersonRepository{
		persons: make(map[string]*domain.Person),
	}
}

// Save persists a person
func (r *InMemoryPersonRepository) Save(ctx context.Context, person *domain.Person) error {
	if person == nil {
		return fmt.Errorf("person cannot be nil")
	}
	r.persons[person.ID().String()] = person
	return nil
}

// FindByID retrieves a person by ID
func (r *InMemoryPersonRepository) FindByID(ctx context.Context, id domain.PersonID) (*domain.Person, error) {
	person, exists := r.persons[id.String()]
	if !exists {
		return nil, fmt.Errorf("person not found: %s", id)
	}
	return person, nil
}

// FindByOrganization retrieves all persons for an organization
func (r *InMemoryPersonRepository) FindByOrganization(ctx context.Context, orgID domain.OrganizationID) ([]*domain.Person, error) {
	var result []*domain.Person
	for _, person := range r.persons {
		if person.OrganizationID() == orgID {
			result = append(result, person)
		}
	}
	return result, nil
}

// FindByEmail retrieves a person by email
func (r *InMemoryPersonRepository) FindByEmail(ctx context.Context, email string) (*domain.Person, error) {
	for _, person := range r.persons {
		if person.Email() != nil && person.Email().String() == email {
			return person, nil
		}
	}
	return nil, fmt.Errorf("person not found with email: %s", email)
}

// FindPrimaryContact retrieves the primary contact for an organization
func (r *InMemoryPersonRepository) FindPrimaryContact(ctx context.Context, orgID domain.OrganizationID) (*domain.Person, error) {
	for _, person := range r.persons {
		if person.OrganizationID() == orgID && person.IsPrimaryContact() {
			return person, nil
		}
	}
	return nil, fmt.Errorf("primary contact not found for organization: %s", orgID)
}

// Delete removes a person
func (r *InMemoryPersonRepository) Delete(ctx context.Context, id domain.PersonID) error {
	if _, exists := r.persons[id.String()]; !exists {
		return fmt.Errorf("person not found: %s", id)
	}
	delete(r.persons, id.String())
	return nil
}

// Exists checks if a person exists
func (r *InMemoryPersonRepository) Exists(ctx context.Context, id domain.PersonID) (bool, error) {
	_, exists := r.persons[id.String()]
	return exists, nil
}

// InMemoryAddressRepository is an in-memory implementation of AddressRepository
type InMemoryAddressRepository struct {
	addresses map[string]struct {
		address *domain.Address
		orgID   domain.OrganizationID
	}
}

// NewInMemoryAddressRepository creates a new in-memory address repository
func NewInMemoryAddressRepository() *InMemoryAddressRepository {
	return &InMemoryAddressRepository{
		addresses: make(map[string]struct {
			address *domain.Address
			orgID   domain.OrganizationID
		}),
	}
}

// Save persists an address
func (r *InMemoryAddressRepository) Save(ctx context.Context, address *domain.Address, addressID domain.AddressID, orgID domain.OrganizationID) error {
	if address == nil {
		return fmt.Errorf("address cannot be nil")
	}
	r.addresses[addressID.String()] = struct {
		address *domain.Address
		orgID   domain.OrganizationID
	}{
		address: address,
		orgID:   orgID,
	}
	return nil
}

// FindByID retrieves an address by ID
func (r *InMemoryAddressRepository) FindByID(ctx context.Context, id domain.AddressID) (*domain.Address, error) {
	entry, exists := r.addresses[id.String()]
	if !exists {
		return nil, fmt.Errorf("address not found: %s", id)
	}
	return entry.address, nil
}

// FindByOrganization retrieves all addresses for an organization
func (r *InMemoryAddressRepository) FindByOrganization(ctx context.Context, orgID domain.OrganizationID) ([]*domain.Address, error) {
	var result []*domain.Address
	for _, entry := range r.addresses {
		if entry.orgID == orgID {
			result = append(result, entry.address)
		}
	}
	return result, nil
}

// FindPrimary retrieves the primary address for an organization
func (r *InMemoryAddressRepository) FindPrimary(ctx context.Context, orgID domain.OrganizationID) (*domain.Address, error) {
	for _, entry := range r.addresses {
		if entry.orgID == orgID {
			// Since Address is a value object without ID, we can't mark primary here
			// This would need additional state tracking
			return entry.address, nil
		}
	}
	return nil, fmt.Errorf("no address found for organization: %s", orgID)
}

// Delete removes an address
func (r *InMemoryAddressRepository) Delete(ctx context.Context, id domain.AddressID) error {
	if _, exists := r.addresses[id.String()]; !exists {
		return fmt.Errorf("address not found: %s", id)
	}
	delete(r.addresses, id.String())
	return nil
}

// Exists checks if an address exists
func (r *InMemoryAddressRepository) Exists(ctx context.Context, id domain.AddressID) (bool, error) {
	_, exists := r.addresses[id.String()]
	return exists, nil
}
