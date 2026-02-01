package domain

import (
	"fmt"
	"time"
)

// Organization represents a party (client or supplier) in the TramaTex system
type Organization struct {
	id         OrganizationID
	name       string
	role       OrganizationRole
	status     OrganizationStatus
	taxID      *TaxID
	website    string
	notes      string
	createdBy  string // User ID who created
	createdAt  time.Time
	modifiedBy string // User ID who last modified
	modifiedAt time.Time
	persons    []*Person
	addresses  []*Address
}

// NewOrganization creates a new Organization aggregate root
func NewOrganization(
	id OrganizationID,
	name string,
	role OrganizationRole,
	taxID *TaxID,
	createdBy string,
) (*Organization, error) {
	if id.String() == "" {
		return nil, fmt.Errorf("organization ID cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("organization name cannot be empty")
	}
	if !role.IsValid() {
		return nil, fmt.Errorf("invalid organization role: %s", role)
	}
	if createdBy == "" {
		return nil, fmt.Errorf("createdBy user ID cannot be empty")
	}

	now := time.Now()
	return &Organization{
		id:         id,
		name:       name,
		role:       role,
		status:     OrganizationStatusActive,
		taxID:      taxID,
		createdBy:  createdBy,
		createdAt:  now,
		modifiedBy: createdBy,
		modifiedAt: now,
		persons:    make([]*Person, 0),
		addresses:  make([]*Address, 0),
	}, nil
}

// ID returns the organization ID
func (o *Organization) ID() OrganizationID {
	return o.id
}

// Name returns the organization name
func (o *Organization) Name() string {
	return o.name
}

// UpdateName updates the organization name
func (o *Organization) UpdateName(name string, modifiedBy string) error {
	if name == "" {
		return fmt.Errorf("organization name cannot be empty")
	}
	if modifiedBy == "" {
		return fmt.Errorf("modifiedBy user ID cannot be empty")
	}
	o.name = name
	o.modifiedBy = modifiedBy
	o.modifiedAt = time.Now()
	return nil
}

// Role returns the organization role
func (o *Organization) Role() OrganizationRole {
	return o.role
}

// Status returns the organization status
func (o *Organization) Status() OrganizationStatus {
	return o.status
}

// Activate activates the organization
func (o *Organization) Activate(modifiedBy string) error {
	if modifiedBy == "" {
		return fmt.Errorf("modifiedBy user ID cannot be empty")
	}
	if o.status == OrganizationStatusActive {
		return fmt.Errorf("organization is already active")
	}
	o.status = OrganizationStatusActive
	o.modifiedBy = modifiedBy
	o.modifiedAt = time.Now()
	return nil
}

// Deactivate deactivates the organization
func (o *Organization) Deactivate(modifiedBy string) error {
	if modifiedBy == "" {
		return fmt.Errorf("modifiedBy user ID cannot be empty")
	}
	if o.status == OrganizationStatusInactive {
		return fmt.Errorf("organization is already inactive")
	}
	o.status = OrganizationStatusInactive
	o.modifiedBy = modifiedBy
	o.modifiedAt = time.Now()
	return nil
}

// TaxID returns the tax ID
func (o *Organization) TaxID() *TaxID {
	return o.taxID
}

// UpdateTaxID updates the tax ID
func (o *Organization) UpdateTaxID(taxID *TaxID, modifiedBy string) error {
	if modifiedBy == "" {
		return fmt.Errorf("modifiedBy user ID cannot be empty")
	}
	o.taxID = taxID
	o.modifiedBy = modifiedBy
	o.modifiedAt = time.Now()
	return nil
}

// Website returns the website
func (o *Organization) Website() string {
	return o.website
}

// UpdateWebsite updates the website
func (o *Organization) UpdateWebsite(website string, modifiedBy string) error {
	if modifiedBy == "" {
		return fmt.Errorf("modifiedBy user ID cannot be empty")
	}
	o.website = website
	o.modifiedBy = modifiedBy
	o.modifiedAt = time.Now()
	return nil
}

// Notes returns the notes
func (o *Organization) Notes() string {
	return o.notes
}

// UpdateNotes updates the notes
func (o *Organization) UpdateNotes(notes string, modifiedBy string) error {
	if modifiedBy == "" {
		return fmt.Errorf("modifiedBy user ID cannot be empty")
	}
	o.notes = notes
	o.modifiedBy = modifiedBy
	o.modifiedAt = time.Now()
	return nil
}

// CreatedBy returns who created the organization
func (o *Organization) CreatedBy() string {
	return o.createdBy
}

// CreatedAt returns when the organization was created
func (o *Organization) CreatedAt() time.Time {
	return o.createdAt
}

// ModifiedBy returns who last modified the organization
func (o *Organization) ModifiedBy() string {
	return o.modifiedBy
}

// ModifiedAt returns when the organization was last modified
func (o *Organization) ModifiedAt() time.Time {
	return o.modifiedAt
}

// AddPerson adds a person to the organization
func (o *Organization) AddPerson(person *Person) error {
	if person == nil {
		return fmt.Errorf("person cannot be nil")
	}
	// Check if person already exists
	for _, p := range o.persons {
		if p.ID() == person.ID() {
			return fmt.Errorf("person already exists in organization")
		}
	}
	o.persons = append(o.persons, person)
	return nil
}

// Persons returns all persons in the organization
func (o *Organization) Persons() []*Person {
	return o.persons
}

// AddAddress adds an address to the organization
func (o *Organization) AddAddress(address *Address) error {
	if address == nil {
		return fmt.Errorf("address cannot be nil")
	}
	// Check if address already exists (basic check)
	for _, a := range o.addresses {
		if a.Equals(address) {
			return fmt.Errorf("address already exists in organization")
		}
	}
	o.addresses = append(o.addresses, address)
	return nil
}

// Addresses returns all addresses for the organization
func (o *Organization) Addresses() []*Address {
	return o.addresses
}
