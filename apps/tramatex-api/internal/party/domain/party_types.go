package domain

import "fmt"

// PartyStatus defines the status of a party
type PartyStatus string

const (
	PartyStatusActive   PartyStatus = "ACTIVE"
	PartyStatusInactive PartyStatus = "INACTIVE"
)

func (s PartyStatus) IsValid() bool {
	switch s {
	case PartyStatusActive, PartyStatusInactive:
		return true
	default:
		return false
	}
}

// PartyRoleType defines the role type of a party
type PartyRoleType string

const (
	PartyRoleClient   PartyRoleType = "CLIENT"
	PartyRoleSupplier PartyRoleType = "SUPPLIER"
	PartyRoleEmployee PartyRoleType = "EMPLOYEE"
)

func (r PartyRoleType) IsValid() bool {
	switch r {
	case PartyRoleClient, PartyRoleSupplier, PartyRoleEmployee:
		return true
	default:
		return false
	}
}

// RelationshipType defines relationship types between parties
type RelationshipType string

const (
	RelationshipIsEmployeeOf   RelationshipType = "IS_EMPLOYEE_OF"
	RelationshipIsSubsidiaryOf RelationshipType = "IS_SUBSIDIARY_OF"
)

func (t RelationshipType) IsValid() bool {
	switch t {
	case RelationshipIsEmployeeOf, RelationshipIsSubsidiaryOf:
		return true
	default:
		return false
	}
}

// PartyID represents a party identifier
type PartyID string

func NewPartyID(id string) (PartyID, error) {
	if id == "" {
		return "", fmt.Errorf("party ID cannot be empty")
	}
	return PartyID(id), nil
}

func (id PartyID) String() string {
	return string(id)
}

func (id PartyID) Value() string {
	return string(id)
}

// PartyRelationshipID represents a relationship identifier
type PartyRelationshipID string

func NewPartyRelationshipID(id string) (PartyRelationshipID, error) {
	if id == "" {
		return "", fmt.Errorf("relationship ID cannot be empty")
	}
	return PartyRelationshipID(id), nil
}

func (id PartyRelationshipID) String() string {
	return string(id)
}

func (id PartyRelationshipID) Value() string {
	return string(id)
}

// ContactDetailsID represents a contact detail identifier
type ContactDetailsID string

func NewContactDetailsID(id string) (ContactDetailsID, error) {
	if id == "" {
		return "", fmt.Errorf("contact details ID cannot be empty")
	}
	return ContactDetailsID(id), nil
}

func (id ContactDetailsID) String() string {
	return string(id)
}

func (id ContactDetailsID) Value() string {
	return string(id)
}

// AddressID represents an address identifier
type AddressID string

func NewAddressID(id string) (AddressID, error) {
	if id == "" {
		return "", fmt.Errorf("address ID cannot be empty")
	}
	return AddressID(id), nil
}

func (id AddressID) String() string {
	return string(id)
}

func (id AddressID) Value() string {
	return string(id)
}
