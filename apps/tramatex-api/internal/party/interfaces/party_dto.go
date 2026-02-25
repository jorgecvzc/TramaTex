package interfaces

import "github.com/joran-cortez/tramatex/internal/party/domain"

// PartyDTO represents a party in API responses

type PartyDTO struct {
	ID                  string                  `json:"id"`
	Status              string                  `json:"status"`
	Roles               []string                `json:"roles"`
	PersonProfile       *PersonProfileDTO       `json:"person_profile,omitempty"`
	OrganizationProfile *OrganizationProfileDTO `json:"organization_profile,omitempty"`
}

type PersonProfileDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type OrganizationProfileDTO struct {
	Name      string              `json:"name"`
	TaxID     string              `json:"tax_id,omitempty"`
	TaxIDType string              `json:"tax_id_type,omitempty"`
	Website   string              `json:"website,omitempty"`
	Contacts  []ContactDetailsDTO `json:"contacts,omitempty"`
}

type ContactDetailsDTO struct {
	ID              string `json:"id"`
	TypeDescription string `json:"type_description"`
	Phone           string `json:"phone,omitempty"`
	Email           string `json:"email,omitempty"`
	RelatedPartyID  string `json:"related_party_id,omitempty"`
}

type PartyRelationshipDTO struct {
	ID          string `json:"id"`
	FromPartyID string `json:"from_party_id"`
	ToPartyID   string `json:"to_party_id"`
	Type        string `json:"type"`
}

type AddressDTO struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	Province   string `json:"province,omitempty"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// Requests

type CreatePartyRequest struct {
	ID                  string                      `json:"id"`
	Status              string                      `json:"status"`
	Roles               []string                    `json:"roles"`
	PersonProfile       *PersonProfileRequest       `json:"person_profile,omitempty"`
	OrganizationProfile *OrganizationProfileRequest `json:"organization_profile,omitempty"`
}

type PersonProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type OrganizationProfileRequest struct {
	Name      string                  `json:"name"`
	TaxID     string                  `json:"tax_id,omitempty"`
	TaxIDType string                  `json:"tax_id_type,omitempty"`
	Website   string                  `json:"website,omitempty"`
	Contacts  []ContactDetailsRequest `json:"contacts,omitempty"`
}

type ContactDetailsRequest struct {
	ID              string `json:"id"`
	TypeDescription string `json:"type_description"`
	Phone           string `json:"phone,omitempty"`
	Email           string `json:"email,omitempty"`
	RelatedPartyID  string `json:"related_party_id,omitempty"`
}

type UpdatePartyRequest struct {
	Status              string                      `json:"status,omitempty"`
	PersonProfile       *PersonProfileRequest       `json:"person_profile,omitempty"`
	OrganizationProfile *OrganizationProfileRequest `json:"organization_profile,omitempty"`
}

type ChangePartyStatusRequest struct {
	Status string `json:"status"`
}

type AddPartyRoleRequest struct {
	Role string `json:"role"`
}

type CreateRelationshipRequest struct {
	ID        string `json:"id"`
	ToPartyID string `json:"to_party_id"`
	Type      string `json:"type"`
}

type CreateContactDetailsRequest struct {
	ID              string `json:"id"`
	TypeDescription string `json:"type_description"`
	Phone           string `json:"phone,omitempty"`
	Email           string `json:"email,omitempty"`
	RelatedPartyID  string `json:"related_party_id,omitempty"`
}

type UpdateContactDetailsRequest struct {
	TypeDescription *string `json:"type_description,omitempty"`
	Phone           *string `json:"phone,omitempty"`
	Email           *string `json:"email,omitempty"`
	RelatedPartyID  *string `json:"related_party_id,omitempty"`
}

type CreatePartyAddressRequest struct {
	ID         string `json:"id"`
	Street     string `json:"street"`
	City       string `json:"city"`
	Province   string `json:"province,omitempty"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country,omitempty"`
}

// Mappers

func MapPartyToDTO(party *domain.Party) *PartyDTO {
	if party == nil {
		return nil
	}

	dto := &PartyDTO{
		ID:     party.ID().String(),
		Status: string(party.Status()),
		Roles:  make([]string, 0),
	}

	for _, role := range party.Roles() {
		dto.Roles = append(dto.Roles, string(role.Type()))
	}

	if profile := party.PersonProfile(); profile != nil {
		dto.PersonProfile = &PersonProfileDTO{
			FirstName: profile.FirstName(),
			LastName:  profile.LastName(),
		}
	}

	if profile := party.OrganizationProfile(); profile != nil {
		orgDTO := &OrganizationProfileDTO{
			Name:     profile.Name(),
			Website:  profile.Website(),
			Contacts: make([]ContactDetailsDTO, 0),
		}
		if profile.TaxID() != nil {
			orgDTO.TaxID = profile.TaxID().Value()
			orgDTO.TaxIDType = profile.TaxID().Type()
		}

		for _, contact := range profile.Contacts() {
			contactDTO := ContactDetailsDTO{
				ID:              contact.ID().String(),
				TypeDescription: contact.TypeDescription(),
			}
			if contact.Phone() != nil {
				contactDTO.Phone = contact.Phone().Value()
			}
			if contact.Email() != nil {
				contactDTO.Email = contact.Email().Value()
			}
			if contact.RelatedPartyID() != nil {
				contactDTO.RelatedPartyID = contact.RelatedPartyID().String()
			}
			orgDTO.Contacts = append(orgDTO.Contacts, contactDTO)
		}

		dto.OrganizationProfile = orgDTO
	}

	return dto
}

func MapContactDetailsToDTO(contact *domain.ContactDetails) *ContactDetailsDTO {
	if contact == nil {
		return nil
	}

	dto := &ContactDetailsDTO{
		ID:              contact.ID().String(),
		TypeDescription: contact.TypeDescription(),
	}
	if contact.Phone() != nil {
		dto.Phone = contact.Phone().Value()
	}
	if contact.Email() != nil {
		dto.Email = contact.Email().Value()
	}
	if contact.RelatedPartyID() != nil {
		dto.RelatedPartyID = contact.RelatedPartyID().String()
	}

	return dto
}

func MapPartyRelationshipToDTO(relationship *domain.PartyRelationship) *PartyRelationshipDTO {
	if relationship == nil {
		return nil
	}

	return &PartyRelationshipDTO{
		ID:          relationship.ID().String(),
		FromPartyID: relationship.FromID().String(),
		ToPartyID:   relationship.ToID().String(),
		Type:        string(relationship.Type()),
	}
}

func MapAddressToDTO(address *domain.Address) *AddressDTO {
	if address == nil {
		return nil
	}

	return &AddressDTO{
		Street:     address.Street(),
		City:       address.City(),
		Province:   address.Province(),
		PostalCode: address.PostalCode(),
		Country:    address.Country(),
	}
}

// PartyBatchDTO is a minimal DTO for batch operations
type PartyBatchDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Reference string `json:"reference,omitempty"`
	TaxID     string `json:"tax_id,omitempty"`
	TaxIDType string `json:"tax_id_type,omitempty"`
}

// MapPartyToBatchDTO maps a Party domain model to a minimal batch DTO
func MapPartyToBatchDTO(party *domain.Party) PartyBatchDTO {
	if party == nil {
		return PartyBatchDTO{}
	}

	dto := PartyBatchDTO{
		ID: party.ID().String(),
	}

	// Extract name from organization or person profile
	if orgProfile := party.OrganizationProfile(); orgProfile != nil {
		dto.Name = orgProfile.Name()
		if taxID := orgProfile.TaxID(); taxID != nil {
			dto.TaxID = taxID.Value()
			dto.TaxIDType = taxID.Type()
		}
	} else if personProfile := party.PersonProfile(); personProfile != nil {
		dto.Name = personProfile.FirstName() + " " + personProfile.LastName()
	}

	for _, role := range party.Roles() {
		if role.CreationIdentifier() == nil {
			continue
		}

		identifier := *role.CreationIdentifier()
		if identifier == "" {
			continue
		}

		dto.Reference = identifier
		break
	}

	return dto
}
