package persistence

import (
	"context"
	"testing"

	"github.com/joran-cortez/tramatex/internal/party/domain"
)

// TestPostgreSQLOrganizationRepository_Integration tests PostgreSQL organization repository
// Note: These tests require a running PostgreSQL instance with test database
func TestPostgreSQLOrganizationRepository_Save_And_FindByID_Integration(t *testing.T) {
	// Skip if no PostgreSQL available
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUp(); err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer func() {
		if err := tdb.TearDown(); err != nil {
			t.Logf("Failed to tear down test database: %v", err)
		}
	}()

	repo := NewPostgreSQLOrganizationRepository(tdb.DB)
	ctx := context.Background()

	// Create and save organization
	id, _ := domain.NewOrganizationID("org-integration-001")
	org, _ := domain.NewOrganization(id, "Test Integration Corp", domain.OrganizationRoleClient, nil, "test-user")

	err := repo.Save(ctx, org)
	if err != nil {
		t.Errorf("Save should not error, got: %v", err)
	}

	// Retrieve and verify
	retrieved, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Errorf("FindByID should not error, got: %v", err)
	}
	if retrieved.Name() != "Test Integration Corp" {
		t.Errorf("Organization name mismatch, got: %s", retrieved.Name())
	}
}

func TestPostgreSQLOrganizationRepository_FindByRole_Integration(t *testing.T) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUp(); err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer func() {
		if err := tdb.TearDown(); err != nil {
			t.Logf("Failed to tear down test database: %v", err)
		}
	}()

	repo := NewPostgreSQLOrganizationRepository(tdb.DB)
	ctx := context.Background()

	// Create organizations with different roles
	id1, _ := domain.NewOrganizationID("org-client-001")
	org1, _ := domain.NewOrganization(id1, "Client Corp", domain.OrganizationRoleClient, nil, "test-user")
	if err := repo.Save(ctx, org1); err != nil {
		t.Fatalf("Failed to save org1: %v", err)
	}

	id2, _ := domain.NewOrganizationID("org-supplier-001")
	org2, _ := domain.NewOrganization(id2, "Supplier Corp", domain.OrganizationRoleSupplier, nil, "test-user")
	if err := repo.Save(ctx, org2); err != nil {
		t.Fatalf("Failed to save org2: %v", err)
	}

	// Query by role
	clients, _ := repo.FindByRole(ctx, domain.OrganizationRoleClient)
	if len(clients) != 1 {
		t.Errorf("Should find 1 client, got: %d", len(clients))
	}

	suppliers, _ := repo.FindByRole(ctx, domain.OrganizationRoleSupplier)
	if len(suppliers) != 1 {
		t.Errorf("Should find 1 supplier, got: %d", len(suppliers))
	}
}

func TestPostgreSQLOrganizationRepository_Count_Integration(t *testing.T) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUp(); err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer func() {
		if err := tdb.TearDown(); err != nil {
			t.Logf("Failed to tear down test database: %v", err)
		}
	}()

	repo := NewPostgreSQLOrganizationRepository(tdb.DB)
	ctx := context.Background()

	count, _ := repo.Count(ctx)
	if count != 0 {
		t.Errorf("Initial count should be 0, got: %d", count)
	}

	// Save organizations
	for i := 1; i <= 3; i++ {
		id, _ := domain.NewOrganizationID("org-" + string(rune('0'+i)))
		org, _ := domain.NewOrganization(id, "Org "+string(rune('0'+i)), domain.OrganizationRoleClient, nil, "test-user")
		if err := repo.Save(ctx, org); err != nil {
			t.Fatalf("Failed to save org: %v", err)
		}
	}

	count, _ = repo.Count(ctx)
	if count != 3 {
		t.Errorf("Count after saves should be 3, got: %d", count)
	}
}

func TestPostgreSQLPersonRepository_Save_And_FindByEmail_Integration(t *testing.T) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUp(); err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer func() {
		if err := tdb.TearDown(); err != nil {
			t.Logf("Failed to tear down test database: %v", err)
		}
	}()

	// Setup: create organization first
	orgRepo := NewPostgreSQLOrganizationRepository(tdb.DB)
	personRepo := NewPostgreSQLPersonRepository(tdb.DB)
	ctx := context.Background()

	orgID, _ := domain.NewOrganizationID("org-integration-001")
	org, _ := domain.NewOrganization(orgID, "Test Corp", domain.OrganizationRoleClient, nil, "test-user")
	if err := orgRepo.Save(ctx, org); err != nil {
		t.Fatalf("Failed to save org: %v", err)
	}

	// Create and save person
	personID, _ := domain.NewPersonID("person-integration-001")
	email, _ := domain.NewEmail("john.integration@example.com")
	person := domain.NewPerson(personID, orgID, "John", "Integration", email, "test-user")

	err := personRepo.Save(ctx, person)
	if err != nil {
		t.Errorf("Save person should not error, got: %v", err)
	}

	// Retrieve by email
	retrieved, err := personRepo.FindByEmail(ctx, "john.integration@example.com")
	if err != nil {
		t.Errorf("FindByEmail should not error, got: %v", err)
	}
	if retrieved.FirstName() != "John" {
		t.Errorf("Person first name mismatch, got: %s", retrieved.FirstName())
	}
}

func TestPostgreSQLPersonRepository_FindByOrganization_Integration(t *testing.T) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUp(); err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer func() {
		if err := tdb.TearDown(); err != nil {
			t.Logf("Failed to tear down test database: %v", err)
		}
	}()

	orgRepo := NewPostgreSQLOrganizationRepository(tdb.DB)
	personRepo := NewPostgreSQLPersonRepository(tdb.DB)
	ctx := context.Background()

	// Create organization
	orgID, _ := domain.NewOrganizationID("org-integration-002")
	org, _ := domain.NewOrganization(orgID, "Multi-Person Corp", domain.OrganizationRoleClient, nil, "test-user")
	if err := orgRepo.Save(ctx, org); err != nil {
		t.Fatalf("Failed to save org: %v", err)
	}

	// Create multiple persons
	for i := 1; i <= 3; i++ {
		personID, _ := domain.NewPersonID("person-integration-00" + string(rune('0'+i)))
		email, _ := domain.NewEmail("person" + string(rune('0'+i)) + "@integration.com")
		person := domain.NewPerson(personID, orgID, "Person", ""+string(rune('0'+i)), email, "test-user")
		if err := personRepo.Save(ctx, person); err != nil {
			t.Fatalf("Failed to save person: %v", err)
		}
	}

	// Query by organization
	persons, _ := personRepo.FindByOrganization(ctx, orgID)
	if len(persons) != 3 {
		t.Errorf("Should find 3 persons, got: %d", len(persons))
	}
}

func TestPostgreSQLAddressRepository_Save_And_FindByID_Integration(t *testing.T) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUp(); err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer func() {
		if err := tdb.TearDown(); err != nil {
			t.Logf("Failed to tear down test database: %v", err)
		}
	}()

	orgRepo := NewPostgreSQLOrganizationRepository(tdb.DB)
	addressRepo := NewPostgreSQLAddressRepository(tdb.DB)
	ctx := context.Background()

	// Create organization
	orgID, _ := domain.NewOrganizationID("org-integration-003")
	org, _ := domain.NewOrganization(orgID, "Address Test Corp", domain.OrganizationRoleClient, nil, "test-user")
	if err := orgRepo.Save(ctx, org); err != nil {
		t.Fatalf("Failed to save org: %v", err)
	}

	// Create and save address
	addressID, _ := domain.NewAddressID("addr-integration-001")
	addr, _ := domain.NewAddress("Calle Principal 123", "Madrid", "Madrid", "28001", "Spain")
	err := addressRepo.Save(ctx, addr, addressID, orgID)
	if err != nil {
		t.Errorf("Save address should not error, got: %v", err)
	}

	// Retrieve and verify
	retrieved, err := addressRepo.FindByID(ctx, addressID)
	if err != nil {
		t.Errorf("FindByID should not error, got: %v", err)
	}
	if retrieved.City() != "Madrid" {
		t.Errorf("Address city mismatch, got: %s", retrieved.City())
	}
}

func TestPostgreSQLAddressRepository_FindByOrganization_Integration(t *testing.T) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUp(); err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer func() {
		if err := tdb.TearDown(); err != nil {
			t.Logf("Failed to tear down test database: %v", err)
		}
	}()

	orgRepo := NewPostgreSQLOrganizationRepository(tdb.DB)
	addressRepo := NewPostgreSQLAddressRepository(tdb.DB)
	ctx := context.Background()

	// Create organization
	orgID, _ := domain.NewOrganizationID("org-integration-004")
	org, _ := domain.NewOrganization(orgID, "Multi-Address Corp", domain.OrganizationRoleClient, nil, "test-user")
	if err := orgRepo.Save(ctx, org); err != nil {
		t.Fatalf("Failed to save org: %v", err)
	}

	// Create multiple addresses
	for i := 1; i <= 2; i++ {
		addressID, _ := domain.NewAddressID("addr-integration-00" + string(rune('0'+i)))
		addr, _ := domain.NewAddress("Calle "+string(rune('0'+i))+" 123", "Barcelona", "Barcelona", "0800"+string(rune('0'+i)), "Spain")
		if err := addressRepo.Save(ctx, addr, addressID, orgID); err != nil {
			t.Fatalf("Failed to save address: %v", err)
		}
	}

	// Query by organization
	addresses, _ := addressRepo.FindByOrganization(ctx, orgID)
	if len(addresses) != 2 {
		t.Errorf("Should find 2 addresses, got: %d", len(addresses))
	}
}

// Benchmark tests
func BenchmarkPostgreSQLOrganizationRepository_Save(b *testing.B) {
	tdb := NewTestDB(&testing.T{})
	if tdb.DB == nil {
		b.Skip("PostgreSQL not available for benchmarks")
	}

	if err := tdb.SetUp(); err != nil {
		b.Fatalf("Failed to set up test database: %v", err)
	}
	defer func() {
		if err := tdb.TearDown(); err != nil {
			b.Logf("Failed to tear down test database: %v", err)
		}
	}()

	repo := NewPostgreSQLOrganizationRepository(tdb.DB)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id, _ := domain.NewOrganizationID("org-bench-" + string(rune(i)))
		org, _ := domain.NewOrganization(id, "Benchmark Org", domain.OrganizationRoleClient, nil, "bench-user")
		if err := repo.Save(ctx, org); err != nil {
			b.Fatalf("Failed to save org: %v", err)
		}
	}
}

func BenchmarkPostgreSQLPersonRepository_FindByEmail(b *testing.B) {
	tdb := NewTestDB(&testing.T{})
	if tdb.DB == nil {
		b.Skip("PostgreSQL not available for benchmarks")
	}

	if err := tdb.SetUp(); err != nil {
		b.Fatalf("Failed to set up test database: %v", err)
	}
	defer func() {
		if err := tdb.TearDown(); err != nil {
			b.Logf("Failed to tear down test database: %v", err)
		}
	}()

	orgRepo := NewPostgreSQLOrganizationRepository(tdb.DB)
	personRepo := NewPostgreSQLPersonRepository(tdb.DB)
	ctx := context.Background()

	// Setup: create org and person
	orgID, _ := domain.NewOrganizationID("org-bench")
	org, _ := domain.NewOrganization(orgID, "Benchmark Org", domain.OrganizationRoleClient, nil, "bench-user")
	if err := orgRepo.Save(ctx, org); err != nil {
		b.Fatalf("Failed to save org: %v", err)
	}

	personID, _ := domain.NewPersonID("person-bench")
	email, _ := domain.NewEmail("bench@example.com")
	person := domain.NewPerson(personID, orgID, "Bench", "User", email, "bench-user")
	if err := personRepo.Save(ctx, person); err != nil {
		b.Fatalf("Failed to save person: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := personRepo.FindByEmail(ctx, "bench@example.com"); err != nil {
			b.Fatalf("FindByEmail failed: %v", err)
		}
	}
}
