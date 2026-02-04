package domain

import (
	"testing"
	"time"
)

func TestNewPartyRequiresProfile(t *testing.T) {
	partyID, _ := NewPartyID("party-001")
	if _, err := NewParty(partyID, PartyStatusActive, "user-1", nil, nil); err == nil {
		t.Fatalf("expected error when no profiles provided")
	}
}

func TestNewPartyValidatesStatus(t *testing.T) {
	partyID, _ := NewPartyID("party-002")
	personProfile, _ := NewPersonProfile("Ana", "Perez")
	if _, err := NewParty(partyID, PartyStatus("INVALID"), "user-1", personProfile, nil); err == nil {
		t.Fatalf("expected error for invalid status")
	}
}

func TestNewPartySetsAuditFields(t *testing.T) {
	partyID, _ := NewPartyID("party-003")
	personProfile, _ := NewPersonProfile("Ana", "Perez")
	party, err := NewParty(partyID, PartyStatusActive, "user-1", personProfile, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if party.CreatedBy() != "user-1" || party.ModifiedBy() != "user-1" {
		t.Fatalf("audit fields not set correctly")
	}
	if party.CreatedAt().IsZero() || party.ModifiedAt().IsZero() {
		t.Fatalf("timestamps should be set")
	}
}

func TestPartyActivateDeactivate(t *testing.T) {
	partyID, _ := NewPartyID("party-004")
	personProfile, _ := NewPersonProfile("Ana", "Perez")
	party, _ := NewParty(partyID, PartyStatusActive, "user-1", personProfile, nil)

	if err := party.Activate("user-2"); err == nil {
		t.Fatalf("expected error when activating already active party")
	}
	if err := party.Deactivate("user-2"); err != nil {
		t.Fatalf("unexpected error deactivating party: %v", err)
	}
	if err := party.Deactivate("user-2"); err == nil {
		t.Fatalf("expected error when deactivating already inactive party")
	}
	if err := party.Activate("user-2"); err != nil {
		t.Fatalf("unexpected error activating party: %v", err)
	}
}

func TestPartyRolesAddRemove(t *testing.T) {
	partyID, _ := NewPartyID("party-005")
	personProfile, _ := NewPersonProfile("Ana", "Perez")
	party, _ := NewParty(partyID, PartyStatusActive, "user-1", personProfile, nil)

	role, _ := NewPartyRole(PartyRoleClient)
	if err := party.AddRole(role); err != nil {
		t.Fatalf("unexpected error adding role: %v", err)
	}
	if err := party.AddRole(role); err == nil {
		t.Fatalf("expected error on duplicate role")
	}
	if err := party.RemoveRole(PartyRoleClient); err != nil {
		t.Fatalf("unexpected error removing role: %v", err)
	}
	if err := party.RemoveRole(PartyRoleClient); err == nil {
		t.Fatalf("expected error removing missing role")
	}
}

func TestPartyRelationshipAdd(t *testing.T) {
	partyID, _ := NewPartyID("party-006")
	personProfile, _ := NewPersonProfile("Ana", "Perez")
	party, _ := NewParty(partyID, PartyStatusActive, "user-1", personProfile, nil)

	invalidRel := PartyRelationship{typeValue: RelationshipType("INVALID")}
	if err := party.AddRelationship(invalidRel); err == nil {
		t.Fatalf("expected error for invalid relationship type")
	}

	fromID, _ := NewPartyID("party-006")
	toID, _ := NewPartyID("party-007")
	relID, _ := NewPartyRelationshipID("rel-001")
	rel, _ := NewPartyRelationship(relID, fromID, toID, RelationshipIsEmployeeOf)
	if err := party.AddRelationship(rel); err != nil {
		t.Fatalf("unexpected error adding relationship: %v", err)
	}
}

func TestPersonProfileValidation(t *testing.T) {
	if _, err := NewPersonProfile("", "Perez"); err == nil {
		t.Fatalf("expected error when first name is empty")
	}
	if _, err := NewPersonProfile("Ana", ""); err == nil {
		t.Fatalf("expected error when last name is empty")
	}
}

func TestOrganizationProfileAndContacts(t *testing.T) {
	if _, err := NewOrganizationProfile("", nil, ""); err == nil {
		t.Fatalf("expected error when organization name is empty")
	}

	profile, err := NewOrganizationProfile("Org", nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	contactID, _ := NewContactDetailsID("contact-001")
	contact, err := NewContactDetails(contactID, "Ventas", nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error creating contact: %v", err)
	}
	if err := profile.AddContact(contact); err != nil {
		t.Fatalf("unexpected error adding contact: %v", err)
	}
	if err := profile.AddContact(contact); err == nil {
		t.Fatalf("expected error adding duplicate contact")
	}

	updatedContact, err := NewContactDetails(contactID, "Soporte", nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error creating updated contact: %v", err)
	}
	if err := profile.UpdateContact(updatedContact); err != nil {
		t.Fatalf("unexpected error updating contact: %v", err)
	}
	if profile.Contacts()[0].TypeDescription() != "Soporte" {
		t.Fatalf("expected contact to be updated")
	}

	if err := profile.RemoveContact(contactID); err != nil {
		t.Fatalf("unexpected error removing contact: %v", err)
	}
	if len(profile.Contacts()) != 0 {
		t.Fatalf("expected contact list to be empty after removal")
	}
}

func TestContactDetailsValidation(t *testing.T) {
	if _, err := NewContactDetails(ContactDetailsID(""), "Ventas", nil, nil, nil); err == nil {
		t.Fatalf("expected error when contact id is empty")
	}

	id, _ := NewContactDetailsID("contact-002")
	if _, err := NewContactDetails(id, "", nil, nil, nil); err == nil {
		t.Fatalf("expected error when type description is empty")
	}
}

func TestPartyTypesValidation(t *testing.T) {
	if !PartyStatusActive.IsValid() {
		t.Fatalf("expected PartyStatusActive to be valid")
	}
	if PartyStatus("INVALID").IsValid() {
		t.Fatalf("expected invalid party status")
	}

	if !PartyRoleClient.IsValid() || PartyRoleType("INVALID").IsValid() {
		t.Fatalf("expected role validation to behave correctly")
	}

	if !RelationshipIsEmployeeOf.IsValid() || RelationshipType("INVALID").IsValid() {
		t.Fatalf("expected relationship validation to behave correctly")
	}

	if _, err := NewPartyID(""); err == nil {
		t.Fatalf("expected error for empty party ID")
	}
}

func TestNewPartyFromPersistence(t *testing.T) {
	partyID, _ := NewPartyID("party-008")
	personProfile, _ := NewPersonProfile("Ana", "Perez")
	createdAt := time.Now().Add(-time.Hour)
	modifiedAt := time.Now()

	role, _ := NewPartyRole(PartyRoleClient)
	party, err := NewPartyFromPersistence(
		partyID,
		PartyStatusActive,
		"user-1",
		createdAt,
		"user-2",
		modifiedAt,
		personProfile,
		nil,
		[]PartyRole{role},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if party.CreatedAt() != createdAt || party.ModifiedAt() != modifiedAt {
		t.Fatalf("expected timestamps to match persisted values")
	}
}

func TestAddressIDValidation(t *testing.T) {
	if _, err := NewAddressID(""); err == nil {
		t.Fatalf("expected error for empty address ID")
	}

	addrID, err := NewAddressID("addr-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addrID.String() != "addr-001" {
		t.Fatalf("expected address ID to match input")
	}
}
