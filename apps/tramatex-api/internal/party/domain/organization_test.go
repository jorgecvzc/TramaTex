package domain

import (
	"testing"
	"time"
)

// Test Organization Creation
func TestNewOrganization_Valid(t *testing.T) {
	id, _ := NewOrganizationID("org-001")
	taxID, _ := NewTaxID("A12345678", "CIF")

	org, err := NewOrganization(id, "Acme Corp", OrganizationRoleClient, taxID, "user-123")
	if err != nil {
		t.Errorf("NewOrganization should not error for valid input, got: %v", err)
	}
	if org == nil {
		t.Error("Organization should not be nil")
	}
	if org.Name() != "Acme Corp" {
		t.Errorf("Name should be 'Acme Corp', got: %s", org.Name())
	}
	if org.Status() != OrganizationStatusActive {
		t.Errorf("Status should be ACTIVE by default, got: %s", org.Status())
	}
	if org.CreatedBy() != "user-123" {
		t.Errorf("CreatedBy should be 'user-123', got: %s", org.CreatedBy())
	}
}

func TestNewOrganization_InvalidInputs(t *testing.T) {
	id, _ := NewOrganizationID("org-001")
	taxID, _ := NewTaxID("A12345678", "CIF")

	tests := []struct {
		id          OrganizationID
		name        string
		role        OrganizationRole
		createdBy   string
		description string
	}{
		{OrganizationID(""), "Acme", OrganizationRoleClient, "user-1", "empty ID"},
		{id, "", OrganizationRoleClient, "user-1", "empty name"},
		{id, "Acme", OrganizationRole("INVALID"), "user-1", "invalid role"},
		{id, "Acme", OrganizationRoleClient, "", "empty createdBy"},
	}

	for _, test := range tests {
		_, err := NewOrganization(test.id, test.name, test.role, taxID, test.createdBy)
		if err == nil {
			t.Errorf("NewOrganization should error for %s", test.description)
		}
	}
}

func TestOrganization_UpdateName(t *testing.T) {
	id, _ := NewOrganizationID("org-001")
	org, _ := NewOrganization(id, "Old Name", OrganizationRoleClient, nil, "user-1")

	err := org.UpdateName("New Name", "user-2")
	if err != nil {
		t.Errorf("UpdateName should not error, got: %v", err)
	}
	if org.Name() != "New Name" {
		t.Errorf("Name should be 'New Name', got: %s", org.Name())
	}
	if org.ModifiedBy() != "user-2" {
		t.Errorf("ModifiedBy should be 'user-2', got: %s", org.ModifiedBy())
	}
}

func TestOrganization_UpdateName_Invalid(t *testing.T) {
	id, _ := NewOrganizationID("org-001")
	org, _ := NewOrganization(id, "Name", OrganizationRoleClient, nil, "user-1")

	err := org.UpdateName("", "user-2")
	if err == nil {
		t.Error("UpdateName should error for empty name")
	}

	err = org.UpdateName("New", "")
	if err == nil {
		t.Error("UpdateName should error for empty modifiedBy")
	}
}

func TestOrganization_ActivateDeactivate(t *testing.T) {
	id, _ := NewOrganizationID("org-001")
	org, _ := NewOrganization(id, "Name", OrganizationRoleClient, nil, "user-1")

	// Should start as ACTIVE
	if org.Status() != OrganizationStatusActive {
		t.Errorf("Initial status should be ACTIVE, got: %s", org.Status())
	}

	// Deactivate
	err := org.Deactivate("user-2")
	if err != nil {
		t.Errorf("Deactivate should not error, got: %v", err)
	}
	if org.Status() != OrganizationStatusInactive {
		t.Errorf("Status should be INACTIVE, got: %s", org.Status())
	}

	// Try to deactivate again
	err = org.Deactivate("user-2")
	if err == nil {
		t.Error("Deactivate should error when already inactive")
	}

	// Activate
	err = org.Activate("user-3")
	if err != nil {
		t.Errorf("Activate should not error, got: %v", err)
	}
	if org.Status() != OrganizationStatusActive {
		t.Errorf("Status should be ACTIVE, got: %s", org.Status())
	}
}

func TestOrganization_AddPerson(t *testing.T) {
	orgID, _ := NewOrganizationID("org-001")
	org, _ := NewOrganization(orgID, "Name", OrganizationRoleClient, nil, "user-1")

	personID, _ := NewPersonID("person-001")
	email, _ := NewEmail("john@example.com")
	person := NewPerson(personID, orgID, "John", "Doe", email, "user-1")

	err := org.AddPerson(person)
	if err != nil {
		t.Errorf("AddPerson should not error, got: %v", err)
	}

	persons := org.Persons()
	if len(persons) != 1 {
		t.Errorf("Should have 1 person, got: %d", len(persons))
	}

	// Try to add same person again
	err = org.AddPerson(person)
	if err == nil {
		t.Error("AddPerson should error when person already exists")
	}

	// Try to add nil person
	err = org.AddPerson(nil)
	if err == nil {
		t.Error("AddPerson should error for nil person")
	}
}

func TestOrganization_AddAddress(t *testing.T) {
	orgID, _ := NewOrganizationID("org-001")
	org, _ := NewOrganization(orgID, "Name", OrganizationRoleClient, nil, "user-1")
	addressID, _ := NewAddressID("addr-001")

	addr, _ := NewAddress("Calle 123", "Madrid", "Madrid", "28001", "Spain")

	err := org.AddAddress(addr, addressID)
	if err != nil {
		t.Errorf("AddAddress should not error, got: %v", err)
	}

	addresses := org.Addresses()
	if len(addresses) != 1 {
		t.Errorf("Should have 1 address, got: %d", len(addresses))
	}

	// Try to add same address again
	err = org.AddAddress(addr, addressID)
	if err == nil {
		t.Error("AddAddress should error when address already exists")
	}

	// Try to add nil address
	addressID2, _ := NewAddressID("addr-002")
	err = org.AddAddress(nil, addressID2)
	if err == nil {
		t.Error("AddAddress should error for nil address")
	}
}

func TestOrganization_TimestampUpdates(t *testing.T) {
	id, _ := NewOrganizationID("org-001")
	org, _ := NewOrganization(id, "Name", OrganizationRoleClient, nil, "user-1")

	initialModifiedAt := org.ModifiedAt()
	time.Sleep(10 * time.Millisecond) // Ensure time difference

	org.UpdateName("New Name", "user-2")

	if org.ModifiedAt().Equal(initialModifiedAt) {
		t.Error("ModifiedAt should update when organization is modified")
	}
}

// Test Person Creation
func TestNewPerson_Valid(t *testing.T) {
	orgID, _ := NewOrganizationID("org-001")
	personID, _ := NewPersonID("person-001")
	email, _ := NewEmail("john@example.com")

	person := NewPerson(personID, orgID, "John", "Doe", email, "user-1")

	if person == nil {
		t.Error("Person should not be nil")
	}
	if person.FullName() != "John Doe" {
		t.Errorf("FullName should be 'John Doe', got: %s", person.FullName())
	}
	if person.IsPrimaryContact() {
		t.Error("Should not be primary contact by default")
	}
}

func TestPerson_SetPhone(t *testing.T) {
	orgID, _ := NewOrganizationID("org-001")
	personID, _ := NewPersonID("person-001")
	email, _ := NewEmail("john@example.com")
	person := NewPerson(personID, orgID, "John", "Doe", email, "user-1")

	phone, _ := NewPhone("666123456")
	person.SetPhone(phone)

	if person.Phone() == nil {
		t.Error("Phone should not be nil after setting")
	}
	if !person.Phone().Equals(phone) {
		t.Errorf("Phone should be set correctly")
	}
}

func TestPerson_SetJobTitle(t *testing.T) {
	orgID, _ := NewOrganizationID("org-001")
	personID, _ := NewPersonID("person-001")
	email, _ := NewEmail("john@example.com")
	person := NewPerson(personID, orgID, "John", "Doe", email, "user-1")

	person.SetJobTitle("Manager")
	if person.JobTitle() != "Manager" {
		t.Errorf("JobTitle should be 'Manager', got: %s", person.JobTitle())
	}
}

func TestPerson_SetPrimaryContact(t *testing.T) {
	orgID, _ := NewOrganizationID("org-001")
	personID, _ := NewPersonID("person-001")
	email, _ := NewEmail("john@example.com")
	person := NewPerson(personID, orgID, "John", "Doe", email, "user-1")

	if person.IsPrimaryContact() {
		t.Error("Should not be primary contact by default")
	}

	person.SetPrimaryContact(true)
	if !person.IsPrimaryContact() {
		t.Error("Should be primary contact after setting")
	}
}
