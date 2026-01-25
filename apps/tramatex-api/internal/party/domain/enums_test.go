package domain

import "testing"

// Test OrganizationRole
func TestOrganizationRole_IsValid(t *testing.T) {
	tests := []struct {
		role    OrganizationRole
		isValid bool
	}{
		{OrganizationRoleClient, true},
		{OrganizationRoleSupplier, true},
		{OrganizationRoleBoth, true},
		{OrganizationRole("INVALID"), false},
		{OrganizationRole(""), false},
	}

	for _, test := range tests {
		if test.role.IsValid() != test.isValid {
			t.Errorf("OrganizationRole(%s).IsValid() should be %v", test.role, test.isValid)
		}
	}
}

func TestOrganizationRole_String(t *testing.T) {
	if OrganizationRoleClient.String() != "CLIENT" {
		t.Errorf("String() should return 'CLIENT', got: %s", OrganizationRoleClient.String())
	}
	if OrganizationRoleSupplier.String() != "SUPPLIER" {
		t.Errorf("String() should return 'SUPPLIER', got: %s", OrganizationRoleSupplier.String())
	}
	if OrganizationRoleBoth.String() != "BOTH" {
		t.Errorf("String() should return 'BOTH', got: %s", OrganizationRoleBoth.String())
	}
}

// Test OrganizationStatus
func TestOrganizationStatus_IsValid(t *testing.T) {
	tests := []struct {
		status  OrganizationStatus
		isValid bool
	}{
		{OrganizationStatusActive, true},
		{OrganizationStatusInactive, true},
		{OrganizationStatus("INVALID"), false},
		{OrganizationStatus(""), false},
	}

	for _, test := range tests {
		if test.status.IsValid() != test.isValid {
			t.Errorf("OrganizationStatus(%s).IsValid() should be %v", test.status, test.isValid)
		}
	}
}

func TestOrganizationStatus_String(t *testing.T) {
	if OrganizationStatusActive.String() != "ACTIVE" {
		t.Errorf("String() should return 'ACTIVE', got: %s", OrganizationStatusActive.String())
	}
	if OrganizationStatusInactive.String() != "INACTIVE" {
		t.Errorf("String() should return 'INACTIVE', got: %s", OrganizationStatusInactive.String())
	}
}

// Test OrganizationID
func TestNewOrganizationID_Valid(t *testing.T) {
	id, err := NewOrganizationID("org-123")
	if err != nil {
		t.Errorf("NewOrganizationID should not error for valid ID, got: %v", err)
	}
	if id.String() != "org-123" {
		t.Errorf("OrganizationID.String() should return 'org-123', got: %s", id.String())
	}
}

func TestNewOrganizationID_Invalid(t *testing.T) {
	_, err := NewOrganizationID("")
	if err == nil {
		t.Error("NewOrganizationID should error for empty ID")
	}
}

// Test PersonID
func TestNewPersonID_Valid(t *testing.T) {
	id, err := NewPersonID("person-123")
	if err != nil {
		t.Errorf("NewPersonID should not error for valid ID, got: %v", err)
	}
	if id.String() != "person-123" {
		t.Errorf("PersonID.String() should return 'person-123', got: %s", id.String())
	}
}

func TestNewPersonID_Invalid(t *testing.T) {
	_, err := NewPersonID("")
	if err == nil {
		t.Error("NewPersonID should error for empty ID")
	}
}

// Test ContactID
func TestNewContactID_Valid(t *testing.T) {
	id, err := NewContactID("contact-123")
	if err != nil {
		t.Errorf("NewContactID should not error for valid ID, got: %v", err)
	}
	if id.String() != "contact-123" {
		t.Errorf("ContactID.String() should return 'contact-123', got: %s", id.String())
	}
}

func TestNewContactID_Invalid(t *testing.T) {
	_, err := NewContactID("")
	if err == nil {
		t.Error("NewContactID should error for empty ID")
	}
}

// Test AddressID
func TestNewAddressID_Valid(t *testing.T) {
	id, err := NewAddressID("address-123")
	if err != nil {
		t.Errorf("NewAddressID should not error for valid ID, got: %v", err)
	}
	if id.String() != "address-123" {
		t.Errorf("AddressID.String() should return 'address-123', got: %s", id.String())
	}
}

func TestNewAddressID_Invalid(t *testing.T) {
	_, err := NewAddressID("")
	if err == nil {
		t.Error("NewAddressID should error for empty ID")
	}
}
