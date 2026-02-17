package interfaces

import (
	"testing"

	"github.com/joran-cortez/tramatex/internal/party/domain"
)

func TestMapPartyToDTO(t *testing.T) {
	partyID, _ := domain.NewPartyID("party-100")
	personProfile, _ := domain.NewPersonProfile("Ana", "Perez")
	taxID, _ := domain.NewTaxID("B12345678", "CIF")
	orgProfile, _ := domain.NewOrganizationProfile("Textiles Perez", taxID, "https://textiles.local")

	contactID, _ := domain.NewContactDetailsID("contact-100")
	phone, _ := domain.NewPhone("+34 600 111 222")
	email, _ := domain.NewEmail("ventas@textiles.local")
	relatedPartyID, _ := domain.NewPartyID("party-related")
	contact, _ := domain.NewContactDetails(contactID, "Ventas", phone, email, &relatedPartyID)
	_ = orgProfile.AddContact(contact)

	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, personProfile, orgProfile)
	role, _ := domain.NewPartyRole(domain.PartyRoleClient, nil)
	_ = party.AddRole(role)

	dto := MapPartyToDTO(party)
	if dto == nil {
		t.Fatalf("Expected DTO, got nil")
	}
	if dto.ID != "party-100" || dto.Status != "ACTIVE" {
		t.Fatalf("Unexpected party DTO values")
	}
	if dto.PersonProfile == nil || dto.PersonProfile.FirstName != "Ana" {
		t.Fatalf("Person profile not mapped")
	}
	if dto.OrganizationProfile == nil || dto.OrganizationProfile.Name != "Textiles Perez" {
		t.Fatalf("Organization profile not mapped")
	}
	if len(dto.OrganizationProfile.Contacts) != 1 {
		t.Fatalf("Expected 1 contact, got %d", len(dto.OrganizationProfile.Contacts))
	}
}

func TestMapContactDetailsToDTO(t *testing.T) {
	contactID, _ := domain.NewContactDetailsID("contact-200")
	contact, _ := domain.NewContactDetails(contactID, "Soporte", nil, nil, nil)
	dto := MapContactDetailsToDTO(contact)
	if dto == nil || dto.ID != "contact-200" || dto.TypeDescription != "Soporte" {
		t.Fatalf("Contact details mapping failed")
	}
}

func TestMapContactDetailsToDTO_WithFields(t *testing.T) {
	contactID, _ := domain.NewContactDetailsID("contact-201")
	phone, _ := domain.NewPhone("+34 600 111 222")
	email, _ := domain.NewEmail("ventas@org.local")
	relatedPartyID, _ := domain.NewPartyID("party-201")
	contact, _ := domain.NewContactDetails(contactID, "Ventas", phone, email, &relatedPartyID)

	dto := MapContactDetailsToDTO(contact)
	if dto == nil {
		t.Fatalf("expected dto")
	}
	if dto.Phone != "+34 600 111 222" || dto.Email != "ventas@org.local" || dto.RelatedPartyID != "party-201" {
		t.Fatalf("expected phone, email, and related party ID to be mapped")
	}
}

func TestMapPartyRelationshipToDTO(t *testing.T) {
	relID, _ := domain.NewPartyRelationshipID("rel-200")
	fromID, _ := domain.NewPartyID("party-a")
	toID, _ := domain.NewPartyID("party-b")
	relationship, _ := domain.NewPartyRelationship(relID, fromID, toID, domain.RelationshipIsEmployeeOf)

	dto := MapPartyRelationshipToDTO(&relationship)
	if dto == nil || dto.ID != "rel-200" || dto.Type != "IS_EMPLOYEE_OF" {
		t.Fatalf("Relationship mapping failed")
	}
}

func TestMapAddressToDTO(t *testing.T) {
	address, _ := domain.NewAddress("Calle 1", "Madrid", "Madrid", "28001", "Spain")
	dto := MapAddressToDTO(address)
	if dto == nil || dto.City != "Madrid" {
		t.Fatalf("Address mapping failed")
	}
}

func TestMapDTOs_NilInputs(t *testing.T) {
	if MapPartyToDTO(nil) != nil {
		t.Fatalf("expected nil party dto")
	}
	if MapContactDetailsToDTO(nil) != nil {
		t.Fatalf("expected nil contact dto")
	}
	if MapPartyRelationshipToDTO(nil) != nil {
		t.Fatalf("expected nil relationship dto")
	}
	if MapAddressToDTO(nil) != nil {
		t.Fatalf("expected nil address dto")
	}
}
