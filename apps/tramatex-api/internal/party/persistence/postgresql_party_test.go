package persistence

import (
	"context"
	"testing"

	"github.com/joran-cortez/tramatex/internal/party/domain"
)

const testUserID = "00000000-0000-0000-0000-000000000001"

func TestPostgreSQLPartyRepository_Save_And_FindByID_Integration(t *testing.T) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUpParty(); err != nil {
		t.Fatalf("Failed to set up party schema: %v", err)
	}
	defer func() {
		if err := tdb.TearDownParty(); err != nil {
			t.Logf("Failed to tear down party schema: %v", err)
		}
	}()

	repo := NewPostgreSQLPartyRepository(tdb.DB)
	ctx := context.Background()

	partyID, _ := domain.NewPartyID("party-001")
	personProfile, _ := domain.NewPersonProfile("Ana", "Pérez")
	taxID, _ := domain.NewTaxID("B12345678", "CIF")
	orgProfile, _ := domain.NewOrganizationProfile("Textiles Pérez", taxID, "https://textiles.local")

	contactID, _ := domain.NewContactDetailsID("contact-001")
	phone, _ := domain.NewPhone("+34 600 111 222")
	email, _ := domain.NewEmail("ventas@textiles.local")
	contact, _ := domain.NewContactDetails(contactID, "Ventas", phone, email, nil)
	if err := orgProfile.AddContact(contact); err != nil {
		t.Fatalf("Failed to add contact: %v", err)
	}

	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, testUserID, personProfile, orgProfile)
	role, _ := domain.NewPartyRole(domain.PartyRoleClient)
	if err := party.AddRole(role); err != nil {
		t.Fatalf("Failed to add role: %v", err)
	}

	if err := repo.Save(ctx, party); err != nil {
		t.Fatalf("Save should not error, got: %v", err)
	}

	fetched, err := repo.FindByID(ctx, partyID)
	if err != nil {
		t.Fatalf("FindByID should not error, got: %v", err)
	}

	if fetched.OrganizationProfile() == nil || fetched.OrganizationProfile().Name() != "Textiles Pérez" {
		t.Fatalf("Organization profile mismatch")
	}
	if fetched.PersonProfile() == nil || fetched.PersonProfile().FirstName() != "Ana" {
		t.Fatalf("Person profile mismatch")
	}
	if len(fetched.Roles()) != 1 {
		t.Fatalf("Expected 1 role, got %d", len(fetched.Roles()))
	}
}

func TestPostgreSQLPartyRelationshipRepository_Save_And_FindByPartyID_Integration(t *testing.T) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUpParty(); err != nil {
		t.Fatalf("Failed to set up party schema: %v", err)
	}
	defer func() {
		if err := tdb.TearDownParty(); err != nil {
			t.Logf("Failed to tear down party schema: %v", err)
		}
	}()

	partyRepo := NewPostgreSQLPartyRepository(tdb.DB)
	relRepo := NewPostgreSQLPartyRelationshipRepository(tdb.DB)
	ctx := context.Background()

	partyID1, _ := domain.NewPartyID("party-a")
	partyID2, _ := domain.NewPartyID("party-b")
	personProfile, _ := domain.NewPersonProfile("A", "One")
	orgProfile, _ := domain.NewOrganizationProfile("Org B", nil, "")

	partyA, _ := domain.NewParty(partyID1, domain.PartyStatusActive, testUserID, personProfile, nil)
	partyB, _ := domain.NewParty(partyID2, domain.PartyStatusActive, testUserID, nil, orgProfile)

	if err := partyRepo.Save(ctx, partyA); err != nil {
		t.Fatalf("Failed to save party A: %v", err)
	}
	if err := partyRepo.Save(ctx, partyB); err != nil {
		t.Fatalf("Failed to save party B: %v", err)
	}

	relID, _ := domain.NewPartyRelationshipID("rel-001")
	rel, _ := domain.NewPartyRelationship(relID, partyID1, partyID2, domain.RelationshipIsEmployeeOf)

	if err := relRepo.Save(ctx, rel); err != nil {
		t.Fatalf("Failed to save relationship: %v", err)
	}

	rels, err := relRepo.FindByPartyID(ctx, partyID1)
	if err != nil {
		t.Fatalf("FindByPartyID should not error, got: %v", err)
	}
	if len(rels) != 1 {
		t.Fatalf("Expected 1 relationship, got %d", len(rels))
	}
}

func TestPostgreSQLPartyAddressRepository_Save_And_ListByParty_Integration(t *testing.T) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUpParty(); err != nil {
		t.Fatalf("Failed to set up party schema: %v", err)
	}
	defer func() {
		if err := tdb.TearDownParty(); err != nil {
			t.Logf("Failed to tear down party schema: %v", err)
		}
	}()

	partyRepo := NewPostgreSQLPartyRepository(tdb.DB)
	addressRepo := NewPostgreSQLPartyAddressRepository(tdb.DB)
	ctx := context.Background()

	partyID, _ := domain.NewPartyID("party-addr")
	personProfile, _ := domain.NewPersonProfile("Addr", "Test")
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, testUserID, personProfile, nil)
	if err := partyRepo.Save(ctx, party); err != nil {
		t.Fatalf("Failed to save party: %v", err)
	}

	addressID, _ := domain.NewAddressID("addr-001")
	address, _ := domain.NewAddress("Calle 1", "Madrid", "Madrid", "28001", "Spain")
	if err := addressRepo.Save(ctx, address, addressID, partyID, testUserID, testUserID); err != nil {
		t.Fatalf("Failed to save address: %v", err)
	}

	addresses, err := addressRepo.FindByPartyID(ctx, partyID)
	if err != nil {
		t.Fatalf("FindByPartyID should not error, got: %v", err)
	}
	if len(addresses) != 1 {
		t.Fatalf("Expected 1 address, got %d", len(addresses))
	}
}

func TestPostgreSQLPartyRelationshipRepository_Delete_Integration(t *testing.T) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUpParty(); err != nil {
		t.Fatalf("Failed to set up party schema: %v", err)
	}
	defer func() {
		if err := tdb.TearDownParty(); err != nil {
			t.Logf("Failed to tear down party schema: %v", err)
		}
	}()

	partyRepo := NewPostgreSQLPartyRepository(tdb.DB)
	relRepo := NewPostgreSQLPartyRelationshipRepository(tdb.DB)
	ctx := context.Background()

	partyID1, _ := domain.NewPartyID("party-rel-del-a")
	partyID2, _ := domain.NewPartyID("party-rel-del-b")
	personProfile, _ := domain.NewPersonProfile("Rel", "One")
	orgProfile, _ := domain.NewOrganizationProfile("Rel Org", nil, "")

	partyA, _ := domain.NewParty(partyID1, domain.PartyStatusActive, testUserID, personProfile, nil)
	partyB, _ := domain.NewParty(partyID2, domain.PartyStatusActive, testUserID, nil, orgProfile)

	if err := partyRepo.Save(ctx, partyA); err != nil {
		t.Fatalf("Failed to save party A: %v", err)
	}
	if err := partyRepo.Save(ctx, partyB); err != nil {
		t.Fatalf("Failed to save party B: %v", err)
	}

	relID, _ := domain.NewPartyRelationshipID("rel-del-001")
	rel, _ := domain.NewPartyRelationship(relID, partyID1, partyID2, domain.RelationshipIsEmployeeOf)
	if err := relRepo.Save(ctx, rel); err != nil {
		t.Fatalf("Failed to save relationship: %v", err)
	}

	if err := relRepo.Delete(ctx, relID); err != nil {
		t.Fatalf("Delete should not error: %v", err)
	}

	rels, err := relRepo.FindByPartyID(ctx, partyID1)
	if err != nil {
		t.Fatalf("FindByPartyID should not error: %v", err)
	}
	if len(rels) != 0 {
		t.Fatalf("Expected 0 relationships after delete, got %d", len(rels))
	}
}

func TestPostgreSQLPartyAddressRepository_FindPrimary_And_Delete_Integration(t *testing.T) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUpParty(); err != nil {
		t.Fatalf("Failed to set up party schema: %v", err)
	}
	defer func() {
		if err := tdb.TearDownParty(); err != nil {
			t.Logf("Failed to tear down party schema: %v", err)
		}
	}()

	partyRepo := NewPostgreSQLPartyRepository(tdb.DB)
	addressRepo := NewPostgreSQLPartyAddressRepository(tdb.DB)
	ctx := context.Background()

	partyID, _ := domain.NewPartyID("party-primary")
	personProfile, _ := domain.NewPersonProfile("Primary", "Address")
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, testUserID, personProfile, nil)
	if err := partyRepo.Save(ctx, party); err != nil {
		t.Fatalf("Failed to save party: %v", err)
	}

	addressID, _ := domain.NewAddressID("addr-primary")
	address, _ := domain.NewAddress("Calle 9", "Madrid", "Madrid", "28009", "Spain")
	if err := addressRepo.Save(ctx, address, addressID, partyID, testUserID, testUserID); err != nil {
		t.Fatalf("Failed to save address: %v", err)
	}

	if _, err := tdb.DB.ExecContext(ctx, "UPDATE party_addresses SET is_primary = true WHERE id = $1", addressID.Value()); err != nil {
		t.Fatalf("Failed to mark primary address: %v", err)
	}

	primary, err := addressRepo.FindPrimary(ctx, partyID)
	if err != nil {
		t.Fatalf("FindPrimary should not error: %v", err)
	}
	if primary.City() != "Madrid" {
		t.Fatalf("Expected primary address city Madrid")
	}

	if err := addressRepo.Delete(ctx, addressID); err != nil {
		t.Fatalf("Delete should not error: %v", err)
	}

	if _, err := addressRepo.FindPrimary(ctx, partyID); err == nil {
		t.Fatalf("Expected error when primary address is missing")
	}
	addresses, err := addressRepo.FindByPartyID(ctx, partyID)
	if err != nil {
		t.Fatalf("FindByPartyID should not error: %v", err)
	}
	if len(addresses) != 0 {
		t.Fatalf("Expected 0 addresses after delete, got %d", len(addresses))
	}
}

func TestPostgreSQLPartyRepository_FindAll_Filters_Integration(t *testing.T) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUpParty(); err != nil {
		t.Fatalf("Failed to set up party schema: %v", err)
	}
	defer func() {
		if err := tdb.TearDownParty(); err != nil {
			t.Logf("Failed to tear down party schema: %v", err)
		}
	}()

	repo := NewPostgreSQLPartyRepository(tdb.DB)
	ctx := context.Background()

	personID, _ := domain.NewPartyID("party-person")
	personProfile, _ := domain.NewPersonProfile("Ana", "Persona")
	partyPerson, _ := domain.NewParty(personID, domain.PartyStatusActive, testUserID, personProfile, nil)

	orgID, _ := domain.NewPartyID("party-org")
	orgProfile, _ := domain.NewOrganizationProfile("Org Name", nil, "")
	partyOrg, _ := domain.NewParty(orgID, domain.PartyStatusActive, testUserID, nil, orgProfile)
	role, _ := domain.NewPartyRole(domain.PartyRoleClient)
	_ = partyOrg.AddRole(role)

	if err := repo.Save(ctx, partyPerson); err != nil {
		t.Fatalf("Failed to save person party: %v", err)
	}
	if err := repo.Save(ctx, partyOrg); err != nil {
		t.Fatalf("Failed to save org party: %v", err)
	}

	parties, err := repo.FindAll(ctx, &PartyFilters{Type: "ORGANIZATION", PageNumber: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("FindAll should not error: %v", err)
	}
	if len(parties) != 1 {
		t.Fatalf("Expected 1 organization party, got %d", len(parties))
	}

	roleFilter := domain.PartyRoleClient
	parties, err = repo.FindAll(ctx, &PartyFilters{Role: &roleFilter, PageNumber: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("FindAll by role should not error: %v", err)
	}
	if len(parties) != 1 {
		t.Fatalf("Expected 1 party with role CLIENT, got %d", len(parties))
	}
}

func TestPostgreSQLPartyRepository_Count_Exists_Delete_Integration(t *testing.T) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUpParty(); err != nil {
		t.Fatalf("Failed to set up party schema: %v", err)
	}
	defer func() {
		if err := tdb.TearDownParty(); err != nil {
			t.Logf("Failed to tear down party schema: %v", err)
		}
	}()

	repo := NewPostgreSQLPartyRepository(tdb.DB)
	ctx := context.Background()

	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count should not error: %v", err)
	}
	if count != 0 {
		t.Fatalf("Expected count 0, got %d", count)
	}

	partyID, _ := domain.NewPartyID("party-delete")
	personProfile, _ := domain.NewPersonProfile("Delete", "Test")
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, testUserID, personProfile, nil)
	if err := repo.Save(ctx, party); err != nil {
		t.Fatalf("Save should not error: %v", err)
	}

	exists, err := repo.Exists(ctx, partyID)
	if err != nil || !exists {
		t.Fatalf("Expected party to exist")
	}

	if err := repo.Delete(ctx, partyID); err != nil {
		t.Fatalf("Delete should not error: %v", err)
	}

	exists, err = repo.Exists(ctx, partyID)
	if err != nil {
		t.Fatalf("Exists should not error: %v", err)
	}
	if exists {
		t.Fatalf("Expected party to be deleted")
	}
}
