package domain

import "fmt"

// OrganizationRole defines the role of an organization
type OrganizationRole string

const (
	OrganizationRoleClient   OrganizationRole = "CLIENT"
	OrganizationRoleSupplier OrganizationRole = "SUPPLIER"
	OrganizationRoleBoth     OrganizationRole = "BOTH"
)

// IsValid checks if the OrganizationRole is valid
func (r OrganizationRole) IsValid() bool {
	switch r {
	case OrganizationRoleClient, OrganizationRoleSupplier, OrganizationRoleBoth:
		return true
	default:
		return false
	}
}

// String returns the string representation of OrganizationRole
func (r OrganizationRole) String() string {
	return string(r)
}

// ParseOrganizationRole parses a string to OrganizationRole
func ParseOrganizationRole(s string) OrganizationRole {
	switch s {
	case "CLIENT":
		return OrganizationRoleClient
	case "SUPPLIER":
		return OrganizationRoleSupplier
	case "BOTH":
		return OrganizationRoleBoth
	default:
		return OrganizationRoleClient // Default value
	}
}

// OrganizationStatus defines the status of an organization
type OrganizationStatus string

const (
	OrganizationStatusActive   OrganizationStatus = "ACTIVE"
	OrganizationStatusInactive OrganizationStatus = "INACTIVE"
)

// IsValid checks if the OrganizationStatus is valid
func (s OrganizationStatus) IsValid() bool {
	switch s {
	case OrganizationStatusActive, OrganizationStatusInactive:
		return true
	default:
		return false
	}
}

// String returns the string representation of OrganizationStatus
func (s OrganizationStatus) String() string {
	return string(s)
}

// ParseOrganizationStatus parses a string to OrganizationStatus
func ParseOrganizationStatus(s string) OrganizationStatus {
	switch s {
	case "ACTIVE":
		return OrganizationStatusActive
	case "INACTIVE":
		return OrganizationStatusInactive
	default:
		return OrganizationStatusActive // Default value
	}
}

// OrganizationID represents an organization identifier
type OrganizationID string

// NewOrganizationID creates a new organization ID
func NewOrganizationID(id string) (OrganizationID, error) {
	if id == "" {
		return "", fmt.Errorf("organization ID cannot be empty")
	}
	return OrganizationID(id), nil
}

// String returns the string representation of OrganizationID
func (id OrganizationID) String() string {
	return string(id)
}

// Value returns the value of OrganizationID for database driver
func (id OrganizationID) Value() string {
	return string(id)
}

// PersonID represents a person identifier
type PersonID string

// NewPersonID creates a new person ID
func NewPersonID(id string) (PersonID, error) {
	if id == "" {
		return "", fmt.Errorf("person ID cannot be empty")
	}
	return PersonID(id), nil
}

// String returns the string representation of PersonID
func (id PersonID) String() string {
	return string(id)
}

// Value returns the value of PersonID for database driver
func (id PersonID) Value() string {
	return string(id)
}

// ContactID represents a contact identifier
type ContactID string

// NewContactID creates a new contact ID
func NewContactID(id string) (ContactID, error) {
	if id == "" {
		return "", fmt.Errorf("contact ID cannot be empty")
	}
	return ContactID(id), nil
}

// String returns the string representation of ContactID
func (id ContactID) String() string {
	return string(id)
}

// Value returns the value of ContactID for database driver
func (id ContactID) Value() string {
	return string(id)
}

// AddressID represents an address identifier
type AddressID string

// NewAddressID creates a new address ID
func NewAddressID(id string) (AddressID, error) {
	if id == "" {
		return "", fmt.Errorf("address ID cannot be empty")
	}
	return AddressID(id), nil
}

// String returns the string representation of AddressID
func (id AddressID) String() string {
	return string(id)
}

// Value returns the value of AddressID for database driver
func (id AddressID) Value() string {
	return string(id)
}
