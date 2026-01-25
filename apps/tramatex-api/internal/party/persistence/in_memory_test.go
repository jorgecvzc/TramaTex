package persistence

import (
	"context"
	"testing"

	"github.com/joran-cortez/tramatex/internal/party/domain"
)

// Test InMemoryOrganizationRepository
func TestInMemoryOrganizationRepository_Save_And_FindByID(t *testing.T) {
	repo := NewInMemoryOrganizationRepository()
	ctx := context.Background()

	id, _ := domain.NewOrganizationID("org-001")
	org, _ := domain.NewOrganization(id, "Test Corp", domain.OrganizationRoleClient, nil, "user-1")

	err := repo.Save(ctx, org)
	if err != nil {
		t.Errorf("Save should not error, got: %v", err)
	}

	retrieved, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Errorf("FindByID should not error, got: %v", err)
	}
	if retrieved.Name() != "Test Corp" {
		t.Errorf("Retrieved organization should have name 'Test Corp', got: %s", retrieved.Name())
	}
}

func TestInMemoryOrganizationRepository_FindByID_NotFound(t *testing.T) {
	repo := NewInMemoryOrganizationRepository()
	ctx := context.Background()

	id, _ := domain.NewOrganizationID("nonexistent")
	_, err := repo.FindByID(ctx, id)
	if err == nil {
		t.Error("FindByID should error for nonexistent organization")
	}
}

func TestInMemoryOrganizationRepository_FindByRole(t *testing.T) {
	repo := NewInMemoryOrganizationRepository()
	ctx := context.Background()

	// Add organizations with different roles
	id1, _ := domain.NewOrganizationID("org-001")
	org1, _ := domain.NewOrganization(id1, "Client Corp", domain.OrganizationRoleClient, nil, "user-1")
	repo.Save(ctx, org1)

	id2, _ := domain.NewOrganizationID("org-002")
	org2, _ := domain.NewOrganization(id2, "Supplier Corp", domain.OrganizationRoleSupplier, nil, "user-1")
	repo.Save(ctx, org2)

	clients, _ := repo.FindByRole(ctx, domain.OrganizationRoleClient)
	if len(clients) != 1 {
		t.Errorf("Should find 1 client, got: %d", len(clients))
	}

	suppliers, _ := repo.FindByRole(ctx, domain.OrganizationRoleSupplier)
	if len(suppliers) != 1 {
		t.Errorf("Should find 1 supplier, got: %d", len(suppliers))
	}
}

func TestInMemoryOrganizationRepository_Count(t *testing.T) {
	repo := NewInMemoryOrganizationRepository()
	ctx := context.Background()

	count, _ := repo.Count(ctx)
	if count != 0 {
		t.Errorf("Initial count should be 0, got: %d", count)
	}

	id, _ := domain.NewOrganizationID("org-001")
	org, _ := domain.NewOrganization(id, "Test Corp", domain.OrganizationRoleClient, nil, "user-1")
	repo.Save(ctx, org)

	count, _ = repo.Count(ctx)
	if count != 1 {
		t.Errorf("Count after save should be 1, got: %d", count)
	}
}

func TestInMemoryOrganizationRepository_Delete(t *testing.T) {
	repo := NewInMemoryOrganizationRepository()
	ctx := context.Background()

	id, _ := domain.NewOrganizationID("org-001")
	org, _ := domain.NewOrganization(id, "Test Corp", domain.OrganizationRoleClient, nil, "user-1")
	repo.Save(ctx, org)

	err := repo.Delete(ctx, id)
	if err != nil {
		t.Errorf("Delete should not error, got: %v", err)
	}

	exists, _ := repo.Exists(ctx, id)
	if exists {
		t.Error("Organization should not exist after delete")
	}
}

// Test InMemoryPersonRepository
func TestInMemoryPersonRepository_Save_And_FindByID(t *testing.T) {
	repo := NewInMemoryPersonRepository()
	ctx := context.Background()

	orgID, _ := domain.NewOrganizationID("org-001")
	personID, _ := domain.NewPersonID("person-001")
	email, _ := domain.NewEmail("john@example.com")
	person := domain.NewPerson(personID, orgID, "John", "Doe", email, "user-1")

	err := repo.Save(ctx, person)
	if err != nil {
		t.Errorf("Save should not error, got: %v", err)
	}

	retrieved, err := repo.FindByID(ctx, personID)
	if err != nil {
		t.Errorf("FindByID should not error, got: %v", err)
	}
	if retrieved.FullName() != "John Doe" {
		t.Errorf("Retrieved person should have full name 'John Doe', got: %s", retrieved.FullName())
	}
}

func TestInMemoryPersonRepository_FindByOrganization(t *testing.T) {
	repo := NewInMemoryPersonRepository()
	ctx := context.Background()

	orgID, _ := domain.NewOrganizationID("org-001")

	// Add multiple persons
	for i := 1; i <= 3; i++ {
		personID, _ := domain.NewPersonID("person-00" + string(rune('0'+i)))
		email, _ := domain.NewEmail("person" + string(rune('0'+i)) + "@example.com")
		person := domain.NewPerson(personID, orgID, "Person", ""+string(rune('0'+i)), email, "user-1")
		repo.Save(ctx, person)
	}

	persons, _ := repo.FindByOrganization(ctx, orgID)
	if len(persons) != 3 {
		t.Errorf("Should find 3 persons, got: %d", len(persons))
	}
}

func TestInMemoryPersonRepository_FindByEmail(t *testing.T) {
	repo := NewInMemoryPersonRepository()
	ctx := context.Background()

	orgID, _ := domain.NewOrganizationID("org-001")
	personID, _ := domain.NewPersonID("person-001")
	email, _ := domain.NewEmail("john@example.com")
	person := domain.NewPerson(personID, orgID, "John", "Doe", email, "user-1")
	repo.Save(ctx, person)

	found, err := repo.FindByEmail(ctx, "john@example.com")
	if err != nil {
		t.Errorf("FindByEmail should not error, got: %v", err)
	}
	if found.FullName() != "John Doe" {
		t.Errorf("Found person should be John Doe, got: %s", found.FullName())
	}
}

func TestInMemoryPersonRepository_FindPrimaryContact(t *testing.T) {
	repo := NewInMemoryPersonRepository()
	ctx := context.Background()

	orgID, _ := domain.NewOrganizationID("org-001")

	// Add non-primary contact
	personID1, _ := domain.NewPersonID("person-001")
	email1, _ := domain.NewEmail("john@example.com")
	person1 := domain.NewPerson(personID1, orgID, "John", "Doe", email1, "user-1")
	repo.Save(ctx, person1)

	// Add primary contact
	personID2, _ := domain.NewPersonID("person-002")
	email2, _ := domain.NewEmail("jane@example.com")
	person2 := domain.NewPerson(personID2, orgID, "Jane", "Smith", email2, "user-1")
	person2.SetPrimaryContact(true)
	repo.Save(ctx, person2)

	primary, err := repo.FindPrimaryContact(ctx, orgID)
	if err != nil {
		t.Errorf("FindPrimaryContact should not error, got: %v", err)
	}
	if primary.FullName() != "Jane Smith" {
		t.Errorf("Primary contact should be Jane Smith, got: %s", primary.FullName())
	}
}

// Test InMemoryAddressRepository
func TestInMemoryAddressRepository_Save_And_FindByID(t *testing.T) {
	repo := NewInMemoryAddressRepository()
	ctx := context.Background()

	orgID, _ := domain.NewOrganizationID("org-001")
	addressID, _ := domain.NewAddressID("addr-001")
	addr, _ := domain.NewAddress("Calle 123", "Madrid", "Madrid", "28001", "Spain")

	err := repo.Save(ctx, addr, addressID, orgID)
	if err != nil {
		t.Errorf("Save should not error, got: %v", err)
	}

	retrieved, err := repo.FindByID(ctx, addressID)
	if err != nil {
		t.Errorf("FindByID should not error, got: %v", err)
	}
	if retrieved.City() != "Madrid" {
		t.Errorf("Retrieved address should be in Madrid, got: %s", retrieved.City())
	}
}

func TestInMemoryAddressRepository_FindByOrganization(t *testing.T) {
	repo := NewInMemoryAddressRepository()
	ctx := context.Background()

	orgID, _ := domain.NewOrganizationID("org-001")

	// Add multiple addresses
	for i := 1; i <= 2; i++ {
		addressID, _ := domain.NewAddressID("addr-00" + string(rune('0'+i)))
		addr, _ := domain.NewAddress("Calle "+string(rune('0'+i))+" 123", "Madrid", "Madrid", "2800"+string(rune('0'+i)), "Spain")
		repo.Save(ctx, addr, addressID, orgID)
	}

	addresses, _ := repo.FindByOrganization(ctx, orgID)
	if len(addresses) != 2 {
		t.Errorf("Should find 2 addresses, got: %d", len(addresses))
	}
}

func TestInMemoryAddressRepository_Delete(t *testing.T) {
	repo := NewInMemoryAddressRepository()
	ctx := context.Background()

	orgID, _ := domain.NewOrganizationID("org-001")
	addressID, _ := domain.NewAddressID("addr-001")
	addr, _ := domain.NewAddress("Calle 123", "Madrid", "Madrid", "28001", "Spain")
	repo.Save(ctx, addr, addressID, orgID)

	err := repo.Delete(ctx, addressID)
	if err != nil {
		t.Errorf("Delete should not error, got: %v", err)
	}

	exists, _ := repo.Exists(ctx, addressID)
	if exists {
		t.Error("Address should not exist after delete")
	}
}
