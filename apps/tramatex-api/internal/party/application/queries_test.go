package application

import (
	"context"
	"testing"

	"github.com/joran-cortez/tramatex/internal/party/domain"
	"github.com/joran-cortez/tramatex/internal/party/persistence"
)

func TestGetOrganizationHandler_Success(t *testing.T) {
	orgRepo := persistence.NewInMemoryOrganizationRepository()
	ctx := context.Background()

	// Create organization
	orgID, _ := domain.NewOrganizationID("org-001")
	org, _ := domain.NewOrganization(orgID, "Test Corp", domain.OrganizationRoleClient, nil, "user-1")
	orgRepo.Save(ctx, org)

	// Query it
	handler := NewGetOrganizationHandler(orgRepo)
	query := &GetOrganizationQuery{ID: "org-001"}

	result, err := handler.Handle(ctx, query)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if result.Name() != "Test Corp" {
		t.Errorf("Organization name mismatch, got: %s", result.Name())
	}
}

func TestGetOrganizationHandler_NotFound(t *testing.T) {
	orgRepo := persistence.NewInMemoryOrganizationRepository()
	handler := NewGetOrganizationHandler(orgRepo)
	ctx := context.Background()

	query := &GetOrganizationQuery{ID: "nonexistent"}
	_, err := handler.Handle(ctx, query)
	if err == nil {
		t.Error("Handle should error for nonexistent organization")
	}
}

func TestListOrganizationsHandler_Success(t *testing.T) {
	orgRepo := persistence.NewInMemoryOrganizationRepository()
	ctx := context.Background()

	// Create organizations
	for i := 1; i <= 3; i++ {
		orgID, _ := domain.NewOrganizationID("org-00" + string(rune('0'+i)))
		org, _ := domain.NewOrganization(orgID, "Org "+string(rune('0'+i)), domain.OrganizationRoleClient, nil, "user-1")
		orgRepo.Save(ctx, org)
	}

	handler := NewListOrganizationsHandler(orgRepo)
	query := &ListOrganizationsQuery{PageSize: 10, PageNumber: 1}

	orgs, err := handler.Handle(ctx, query)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if len(orgs) != 3 {
		t.Errorf("Should return 3 organizations, got: %d", len(orgs))
	}
}

func TestListOrganizationsHandler_FilterByRole(t *testing.T) {
	orgRepo := persistence.NewInMemoryOrganizationRepository()
	ctx := context.Background()

	// Create organizations with different roles
	orgID1, _ := domain.NewOrganizationID("org-001")
	org1, _ := domain.NewOrganization(orgID1, "Client Corp", domain.OrganizationRoleClient, nil, "user-1")
	orgRepo.Save(ctx, org1)

	orgID2, _ := domain.NewOrganizationID("org-002")
	org2, _ := domain.NewOrganization(orgID2, "Supplier Corp", domain.OrganizationRoleSupplier, nil, "user-1")
	orgRepo.Save(ctx, org2)

	handler := NewListOrganizationsHandler(orgRepo)
	query := &ListOrganizationsQuery{Role: "CLIENT", PageSize: 10, PageNumber: 1}

	orgs, err := handler.Handle(ctx, query)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if len(orgs) != 1 {
		t.Errorf("Should return 1 CLIENT organization, got: %d", len(orgs))
	}
}

func TestListOrganizationsByRoleHandler_Success(t *testing.T) {
	orgRepo := persistence.NewInMemoryOrganizationRepository()
	ctx := context.Background()

	// Create suppliers
	for i := 1; i <= 2; i++ {
		orgID, _ := domain.NewOrganizationID("supp-00" + string(rune('0'+i)))
		org, _ := domain.NewOrganization(orgID, "Supplier "+string(rune('0'+i)), domain.OrganizationRoleSupplier, nil, "user-1")
		orgRepo.Save(ctx, org)
	}

	handler := NewListOrganizationsByRoleHandler(orgRepo)
	query := &ListOrganizationsByRoleQuery{Role: "SUPPLIER"}

	orgs, err := handler.Handle(ctx, query)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if len(orgs) != 2 {
		t.Errorf("Should return 2 suppliers, got: %d", len(orgs))
	}
}

func TestGetPersonHandler_Success(t *testing.T) {
	personRepo := persistence.NewInMemoryPersonRepository()
	ctx := context.Background()

	// Create person
	orgID, _ := domain.NewOrganizationID("org-001")
	personID, _ := domain.NewPersonID("person-001")
	email, _ := domain.NewEmail("john@example.com")
	person := domain.NewPerson(personID, orgID, "John", "Doe", email, "user-1")
	personRepo.Save(ctx, person)

	// Query it
	handler := NewGetPersonHandler(personRepo)
	query := &GetPersonQuery{ID: "person-001"}

	result, err := handler.Handle(ctx, query)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if result.FirstName() != "John" {
		t.Errorf("Person first name mismatch, got: %s", result.FirstName())
	}
}

func TestListPersonsByOrganizationHandler_Success(t *testing.T) {
	personRepo := persistence.NewInMemoryPersonRepository()
	ctx := context.Background()

	orgID, _ := domain.NewOrganizationID("org-001")

	// Create multiple persons
	for i := 1; i <= 3; i++ {
		personID, _ := domain.NewPersonID("person-00" + string(rune('0'+i)))
		email, _ := domain.NewEmail("person" + string(rune('0'+i)) + "@example.com")
		person := domain.NewPerson(personID, orgID, "Person", ""+string(rune('0'+i)), email, "user-1")
		personRepo.Save(ctx, person)
	}

	handler := NewListPersonsByOrganizationHandler(personRepo)
	query := &ListPersonsByOrganizationQuery{OrganizationID: "org-001"}

	persons, err := handler.Handle(ctx, query)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if len(persons) != 3 {
		t.Errorf("Should return 3 persons, got: %d", len(persons))
	}
}

func TestGetPersonByEmailHandler_Success(t *testing.T) {
	personRepo := persistence.NewInMemoryPersonRepository()
	ctx := context.Background()

	// Create person
	orgID, _ := domain.NewOrganizationID("org-001")
	personID, _ := domain.NewPersonID("person-001")
	email, _ := domain.NewEmail("john@example.com")
	person := domain.NewPerson(personID, orgID, "John", "Doe", email, "user-1")
	personRepo.Save(ctx, person)

	// Query by email
	handler := NewGetPersonByEmailHandler(personRepo)
	query := &GetPersonByEmailQuery{Email: "john@example.com"}

	result, err := handler.Handle(ctx, query)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if result.FullName() != "John Doe" {
		t.Errorf("Person full name mismatch, got: %s", result.FullName())
	}
}

func TestGetPrimaryContactHandler_Success(t *testing.T) {
	personRepo := persistence.NewInMemoryPersonRepository()
	ctx := context.Background()

	orgID, _ := domain.NewOrganizationID("org-001")

	// Create primary contact
	personID, _ := domain.NewPersonID("person-001")
	email, _ := domain.NewEmail("primary@example.com")
	person := domain.NewPerson(personID, orgID, "Primary", "Contact", email, "user-1")
	person.SetPrimaryContact(true)
	personRepo.Save(ctx, person)

	// Query primary contact
	handler := NewGetPrimaryContactHandler(personRepo)
	query := &GetPrimaryContactQuery{OrganizationID: "org-001"}

	result, err := handler.Handle(ctx, query)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if !result.IsPrimaryContact() {
		t.Error("Result should be marked as primary contact")
	}
}

func TestListAddressesByOrganizationHandler_Success(t *testing.T) {
	addressRepo := persistence.NewInMemoryAddressRepository()
	ctx := context.Background()

	orgID, _ := domain.NewOrganizationID("org-001")

	// Create multiple addresses
	for i := 1; i <= 2; i++ {
		addressID, _ := domain.NewAddressID("addr-00" + string(rune('0'+i)))
		addr, _ := domain.NewAddress("Calle "+string(rune('0'+i))+" 123", "Madrid", "Madrid", "2800"+string(rune('0'+i)), "Spain")
		addressRepo.Save(ctx, addr, addressID, orgID)
	}

	handler := NewListAddressesByOrganizationHandler(addressRepo)
	query := &ListAddressesByOrganizationQuery{OrganizationID: "org-001"}

	addresses, err := handler.Handle(ctx, query)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if len(addresses) != 2 {
		t.Errorf("Should return 2 addresses, got: %d", len(addresses))
	}
}

func TestGetPrimaryAddressHandler_Success(t *testing.T) {
	addressRepo := persistence.NewInMemoryAddressRepository()
	ctx := context.Background()

	orgID, _ := domain.NewOrganizationID("org-001")

	// Create address
	addressID, _ := domain.NewAddressID("addr-001")
	addr, _ := domain.NewAddress("Calle Principal 123", "Madrid", "Madrid", "28001", "Spain")
	addressRepo.Save(ctx, addr, addressID, orgID)

	// Query primary address
	handler := NewGetPrimaryAddressHandler(addressRepo)
	query := &GetPrimaryAddressQuery{OrganizationID: "org-001"}

	result, err := handler.Handle(ctx, query)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if result.City() != "Madrid" {
		t.Errorf("Address city mismatch, got: %s", result.City())
	}
}
