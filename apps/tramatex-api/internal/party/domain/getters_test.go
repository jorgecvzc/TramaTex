package domain

import "testing"

func TestPartyGettersAndSetters(t *testing.T) {
	partyID, _ := NewPartyID("party-g1")
	personProfile, _ := NewPersonProfile("Ana", "Perez", nil, nil)
	taxID, _ := NewTaxID("B12345", "CIF")
	orgProfile, _ := NewOrganizationProfile("Org", taxID, "https://org.local", nil, nil)

	party, err := NewParty(partyID, PartyStatusActive, personProfile, orgProfile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if party.ID().String() != "party-g1" {
		t.Fatalf("expected party ID to match input")
	}
	if party.Status() != PartyStatusActive {
		t.Fatalf("expected status ACTIVE")
	}
	if party.PersonProfile() == nil || party.OrganizationProfile() == nil {
		t.Fatalf("expected both profiles to be present")
	}
	if party.PersonProfile().LastName() != "Perez" {
		t.Fatalf("expected last name to match input")
	}
	if len(party.Roles()) != 0 {
		t.Fatalf("expected no roles by default")
	}
	if len(party.Relationships()) != 0 {
		t.Fatalf("expected no relationships by default")
	}

	newPerson, _ := NewPersonProfile("Luis", "Lopez", nil, nil)
	if err := party.SetPersonProfile(newPerson); err != nil {
		t.Fatalf("unexpected error setting person profile: %v", err)
	}
	if party.PersonProfile().FirstName() != "Luis" {
		t.Fatalf("expected updated person profile")
	}

	newOrg, _ := NewOrganizationProfile("Org 2", taxID, "https://org2.local", nil, nil)
	if err := party.SetOrganizationProfile(newOrg); err != nil {
		t.Fatalf("unexpected error setting organization profile: %v", err)
	}
	if party.OrganizationProfile().Name() != "Org 2" {
		t.Fatalf("expected updated organization profile")
	}
}

func TestPartyRelationshipGetters(t *testing.T) {
	fromID, _ := NewPartyID("party-g2")
	toID, _ := NewPartyID("party-g3")
	relID, _ := NewPartyRelationshipID("rel-g1")
	rel, err := NewPartyRelationship(relID, fromID, toID, RelationshipIsEmployeeOf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rel.ID().String() != "rel-g1" {
		t.Fatalf("expected relationship ID to match input")
	}
	if rel.FromID().String() != "party-g2" || rel.ToID().String() != "party-g3" {
		t.Fatalf("expected relationship party IDs to match input")
	}
	if rel.Type() != RelationshipIsEmployeeOf {
		t.Fatalf("expected relationship type IS_EMPLOYEE_OF")
	}
}

func TestProfileAndContactGetters(t *testing.T) {
	taxID, _ := NewTaxID("B12345", "CIF")
	orgProfile, _ := NewOrganizationProfile("Org", taxID, "https://org.local", nil, nil)
	if orgProfile.Name() != "Org" {
		t.Fatalf("expected organization name to match input")
	}
	if orgProfile.TaxID() == nil || orgProfile.TaxID().Value() != "B12345" {
		t.Fatalf("expected tax ID to match input")
	}
	if orgProfile.Website() != "https://org.local" {
		t.Fatalf("expected website to match input")
	}

	phone, _ := NewPhone("+34 666 123456")
	email, _ := NewEmail("sales@org.local")
	relatedParty, _ := NewPartyID("party-g4")
	contactID, _ := NewContactDetailsID("contact-g1")
	contact, err := NewContactDetails(contactID, "Sales", phone, email, &relatedParty)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contact.Phone() == nil || contact.Phone().Value() != "+34 666 123456" {
		t.Fatalf("expected phone to match input")
	}
	if contact.Email() == nil || contact.Email().Value() != "sales@org.local" {
		t.Fatalf("expected email to match input")
	}
	if contact.RelatedPartyID() == nil || contact.RelatedPartyID().String() != "party-g4" {
		t.Fatalf("expected related party ID to match input")
	}
}

func TestValueObjectsAndIDValues(t *testing.T) {
	partyID, _ := NewPartyID("party-g5")
	if partyID.Value() != "party-g5" {
		t.Fatalf("expected party ID value to match input")
	}
	relID, _ := NewPartyRelationshipID("rel-g2")
	if relID.Value() != "rel-g2" {
		t.Fatalf("expected relationship ID value to match input")
	}
	contactID, _ := NewContactDetailsID("contact-g2")
	if contactID.Value() != "contact-g2" {
		t.Fatalf("expected contact ID value to match input")
	}
	addressID, _ := NewAddressID("addr-g1")
	if addressID.Value() != "addr-g1" {
		t.Fatalf("expected address ID value to match input")
	}

	email, _ := NewEmail("billing@org.local")
	if email.Value() != "billing@org.local" {
		t.Fatalf("expected email value to match input")
	}
	phone, _ := NewPhone("666123456")
	if phone.String() != "666123456" || phone.Value() != "666123456" {
		t.Fatalf("expected phone value to match input")
	}

	taxID, _ := NewTaxID("A12345", "VAT")
	if taxID.Value() != "A12345" || taxID.Type() != "VAT" {
		t.Fatalf("expected tax ID value and type to match input")
	}

	address, _ := NewAddress("Calle 1", "Madrid", "Madrid", "28001", "Spain")
	if address.Province() != "Madrid" || address.PostalCode() != "28001" || address.Country() != "Spain" {
		t.Fatalf("expected address fields to match input")
	}
}

func TestPartyStatusTransitions(t *testing.T) {
	partyID, _ := NewPartyID("party-g6")
	personProfile, _ := NewPersonProfile("Ana", "Perez", nil, nil)
	party, _ := NewParty(partyID, PartyStatusActive, personProfile, nil)

	if err := party.Activate(); err == nil {
		t.Fatalf("expected conflict when activating active party")
	}
	if err := party.Deactivate(); err != nil {
		t.Fatalf("expected to deactivate, got %v", err)
	}
	if err := party.Deactivate(); err == nil {
		t.Fatalf("expected conflict when deactivating inactive party")
	}
}

func TestPartyRoleTypeGetter(t *testing.T) {
	role, err := NewPartyRole(PartyRoleSupplier, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role.Type() != PartyRoleSupplier {
		t.Fatalf("expected role SUPPLIER")
	}
}
