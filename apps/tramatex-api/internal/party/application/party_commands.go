package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/joran-cortez/tramatex/internal/party/domain"
	"github.com/joran-cortez/tramatex/internal/party/persistence"
)

// Party command inputs

type PersonProfileInput struct {
	FirstName string
	LastName  string
}

type ContactDetailsInput struct {
	ID              string
	TypeDescription string
	Phone           string
	Email           string
	RelatedPartyID  string
}

type OrganizationProfileInput struct {
	Name      string
	TaxID     string
	TaxIDType string
	Website   string
	Contacts  []ContactDetailsInput
}

// CreatePartyCommand represents a command to create a Party
// At least one profile must be provided.
type CreatePartyCommand struct {
	ID                  string
	Status              string
	Roles               []string
	PersonProfile       *PersonProfileInput
	OrganizationProfile *OrganizationProfileInput
	CreatedBy           string
}

// CreatePartyHandler handles party creation

type CreatePartyHandler struct {
	partyRepo persistence.PartyRepository
}

func NewCreatePartyHandler(partyRepo persistence.PartyRepository) *CreatePartyHandler {
	return &CreatePartyHandler{partyRepo: partyRepo}
}

func (h *CreatePartyHandler) Handle(ctx context.Context, cmd *CreatePartyCommand) (*domain.Party, error) {
	if cmd.ID == "" {
		return nil, fmt.Errorf("party ID cannot be empty")
	}
	if cmd.CreatedBy == "" {
		return nil, fmt.Errorf("createdBy user ID cannot be empty")
	}

	partyID, err := domain.NewPartyID(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid party ID: %w", err)
	}

	status := domain.PartyStatusActive
	if cmd.Status != "" {
		status = domain.PartyStatus(strings.ToUpper(cmd.Status))
		if !status.IsValid() {
			return nil, fmt.Errorf("invalid party status: %s", cmd.Status)
		}
	}

	var personProfile *domain.PersonProfile
	if cmd.PersonProfile != nil {
		personProfile, err = domain.NewPersonProfile(cmd.PersonProfile.FirstName, cmd.PersonProfile.LastName)
		if err != nil {
			return nil, fmt.Errorf("invalid person profile: %w", err)
		}
	}

	var organizationProfile *domain.OrganizationProfile
	if cmd.OrganizationProfile != nil {
		var taxID *domain.TaxID
		if cmd.OrganizationProfile.TaxID != "" {
			taxType := cmd.OrganizationProfile.TaxIDType
			if taxType == "" {
				taxType = "NIF"
			}
			taxID, err = domain.NewTaxID(cmd.OrganizationProfile.TaxID, taxType)
			if err != nil {
				return nil, fmt.Errorf("invalid tax ID: %w", err)
			}
		}

		organizationProfile, err = domain.NewOrganizationProfile(
			cmd.OrganizationProfile.Name,
			taxID,
			cmd.OrganizationProfile.Website,
		)
		if err != nil {
			return nil, fmt.Errorf("invalid organization profile: %w", err)
		}

		for _, input := range cmd.OrganizationProfile.Contacts {
			contactID, err := domain.NewContactDetailsID(input.ID)
			if err != nil {
				return nil, fmt.Errorf("invalid contact details ID: %w", err)
			}

			var phone *domain.Phone
			if input.Phone != "" {
				phone, err = domain.NewPhone(input.Phone)
				if err != nil {
					return nil, fmt.Errorf("invalid phone: %w", err)
				}
			}

			var email *domain.Email
			if input.Email != "" {
				email, err = domain.NewEmail(input.Email)
				if err != nil {
					return nil, fmt.Errorf("invalid email: %w", err)
				}
			}

			var relatedPartyID *domain.PartyID
			if input.RelatedPartyID != "" {
				pid, err := domain.NewPartyID(input.RelatedPartyID)
				if err != nil {
					return nil, fmt.Errorf("invalid related party ID: %w", err)
				}
				relatedPartyID = &pid
			}

			contact, err := domain.NewContactDetails(
				contactID,
				input.TypeDescription,
				phone,
				email,
				relatedPartyID,
			)
			if err != nil {
				return nil, fmt.Errorf("invalid contact details: %w", err)
			}

			if err := organizationProfile.AddContact(contact); err != nil {
				return nil, fmt.Errorf("failed to add contact: %w", err)
			}
		}
	}

	party, err := domain.NewParty(partyID, status, cmd.CreatedBy, personProfile, organizationProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to create party: %w", err)
	}

	for _, role := range cmd.Roles {
		roleType, err := parsePartyRoleType(role)
		if err != nil {
			return nil, err
		}
		partyRole, err := domain.NewPartyRole(roleType)
		if err != nil {
			return nil, err
		}
		if err := party.AddRole(partyRole); err != nil {
			return nil, err
		}
	}

	if err := h.partyRepo.Save(ctx, party); err != nil {
		return nil, fmt.Errorf("failed to save party: %w", err)
	}

	return party, nil
}

// UpdatePartyCommand represents a command to update Party profiles

type UpdatePartyCommand struct {
	ID                  string
	Status              string
	PersonProfile       *PersonProfileInput
	OrganizationProfile *OrganizationProfileInput
	ModifiedBy          string
}

// UpdatePartyHandler handles party updates

type UpdatePartyHandler struct {
	partyRepo persistence.PartyRepository
}

func NewUpdatePartyHandler(partyRepo persistence.PartyRepository) *UpdatePartyHandler {
	return &UpdatePartyHandler{partyRepo: partyRepo}
}

func (h *UpdatePartyHandler) Handle(ctx context.Context, cmd *UpdatePartyCommand) (*domain.Party, error) {
	if cmd.ID == "" {
		return nil, fmt.Errorf("party ID cannot be empty")
	}
	if cmd.ModifiedBy == "" {
		return nil, fmt.Errorf("modifiedBy user ID cannot be empty")
	}

	partyID, err := domain.NewPartyID(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid party ID: %w", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, fmt.Errorf("party not found: %w", err)
	}

	if cmd.Status != "" {
		status := domain.PartyStatus(strings.ToUpper(cmd.Status))
		if !status.IsValid() {
			return nil, fmt.Errorf("invalid party status: %s", cmd.Status)
		}
		if status == domain.PartyStatusActive {
			if err := party.Activate(cmd.ModifiedBy); err != nil {
				return nil, err
			}
		} else {
			if err := party.Deactivate(cmd.ModifiedBy); err != nil {
				return nil, err
			}
		}
	}

	if cmd.PersonProfile != nil {
		existing := party.PersonProfile()
		firstName := cmd.PersonProfile.FirstName
		lastName := cmd.PersonProfile.LastName
		if existing != nil {
			if firstName == "" {
				firstName = existing.FirstName()
			}
			if lastName == "" {
				lastName = existing.LastName()
			}
		}
		profile, err := domain.NewPersonProfile(firstName, lastName)
		if err != nil {
			return nil, fmt.Errorf("invalid person profile: %w", err)
		}
		if err := party.SetPersonProfile(profile, cmd.ModifiedBy); err != nil {
			return nil, err
		}
	}

	if cmd.OrganizationProfile != nil {
		existing := party.OrganizationProfile()
		name := cmd.OrganizationProfile.Name
		if existing != nil && name == "" {
			name = existing.Name()
		}
		if name == "" {
			return nil, fmt.Errorf("organization name is required")
		}

		var taxID *domain.TaxID
		if cmd.OrganizationProfile.TaxID != "" {
			taxType := cmd.OrganizationProfile.TaxIDType
			if taxType == "" {
				taxType = "NIF"
			}
			taxID, err = domain.NewTaxID(cmd.OrganizationProfile.TaxID, taxType)
			if err != nil {
				return nil, fmt.Errorf("invalid tax ID: %w", err)
			}
		} else if existing != nil {
			taxID = existing.TaxID()
		}

		website := cmd.OrganizationProfile.Website
		if existing != nil && website == "" {
			website = existing.Website()
		}

		profile, err := domain.NewOrganizationProfile(name, taxID, website)
		if err != nil {
			return nil, fmt.Errorf("invalid organization profile: %w", err)
		}

		if err := party.SetOrganizationProfile(profile, cmd.ModifiedBy); err != nil {
			return nil, err
		}
	}

	if err := h.partyRepo.Save(ctx, party); err != nil {
		return nil, fmt.Errorf("failed to save party: %w", err)
	}

	return party, nil
}

// ChangePartyStatusCommand represents a status change

type ChangePartyStatusCommand struct {
	ID         string
	Status     string
	ModifiedBy string
}

// ChangePartyStatusHandler handles status changes

type ChangePartyStatusHandler struct {
	partyRepo persistence.PartyRepository
}

func NewChangePartyStatusHandler(partyRepo persistence.PartyRepository) *ChangePartyStatusHandler {
	return &ChangePartyStatusHandler{partyRepo: partyRepo}
}

func (h *ChangePartyStatusHandler) Handle(ctx context.Context, cmd *ChangePartyStatusCommand) (*domain.Party, error) {
	if cmd.ID == "" {
		return nil, fmt.Errorf("party ID cannot be empty")
	}
	if cmd.ModifiedBy == "" {
		return nil, fmt.Errorf("modifiedBy user ID cannot be empty")
	}

	partyID, err := domain.NewPartyID(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid party ID: %w", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, fmt.Errorf("party not found: %w", err)
	}

	status := domain.PartyStatus(strings.ToUpper(cmd.Status))
	if !status.IsValid() {
		return nil, fmt.Errorf("invalid party status: %s", cmd.Status)
	}

	if status == domain.PartyStatusActive {
		if err := party.Activate(cmd.ModifiedBy); err != nil {
			return nil, err
		}
	} else {
		if err := party.Deactivate(cmd.ModifiedBy); err != nil {
			return nil, err
		}
	}

	if err := h.partyRepo.Save(ctx, party); err != nil {
		return nil, fmt.Errorf("failed to save party: %w", err)
	}

	return party, nil
}

// AddPartyRoleCommand represents adding a role to a party

type AddPartyRoleCommand struct {
	ID         string
	Role       string
	ModifiedBy string
}

// AddPartyRoleHandler handles adding a role

type AddPartyRoleHandler struct {
	partyRepo persistence.PartyRepository
}

func NewAddPartyRoleHandler(partyRepo persistence.PartyRepository) *AddPartyRoleHandler {
	return &AddPartyRoleHandler{partyRepo: partyRepo}
}

func (h *AddPartyRoleHandler) Handle(ctx context.Context, cmd *AddPartyRoleCommand) (*domain.Party, error) {
	partyID, err := domain.NewPartyID(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid party ID: %w", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, fmt.Errorf("party not found: %w", err)
	}

	roleType, err := parsePartyRoleType(cmd.Role)
	if err != nil {
		return nil, err
	}

	role, err := domain.NewPartyRole(roleType)
	if err != nil {
		return nil, err
	}
	if err := party.AddRole(role); err != nil {
		return nil, err
	}

	if err := h.partyRepo.Save(ctx, party); err != nil {
		return nil, fmt.Errorf("failed to save party: %w", err)
	}

	return party, nil
}

// RemovePartyRoleCommand represents removing a role from a party

type RemovePartyRoleCommand struct {
	ID         string
	Role       string
	ModifiedBy string
}

// RemovePartyRoleHandler handles removing a role

type RemovePartyRoleHandler struct {
	partyRepo persistence.PartyRepository
}

func NewRemovePartyRoleHandler(partyRepo persistence.PartyRepository) *RemovePartyRoleHandler {
	return &RemovePartyRoleHandler{partyRepo: partyRepo}
}

func (h *RemovePartyRoleHandler) Handle(ctx context.Context, cmd *RemovePartyRoleCommand) (*domain.Party, error) {
	partyID, err := domain.NewPartyID(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid party ID: %w", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, fmt.Errorf("party not found: %w", err)
	}

	roleType, err := parsePartyRoleType(cmd.Role)
	if err != nil {
		return nil, err
	}

	if err := party.RemoveRole(roleType); err != nil {
		return nil, err
	}

	if err := h.partyRepo.Save(ctx, party); err != nil {
		return nil, fmt.Errorf("failed to save party: %w", err)
	}

	return party, nil
}

// AddPartyRelationshipCommand represents adding a relationship

type AddPartyRelationshipCommand struct {
	ID             string
	RelationshipID string
	ToPartyID      string
	Type           string
}

// AddPartyRelationshipHandler handles adding relationships

type AddPartyRelationshipHandler struct {
	relRepo persistence.PartyRelationshipRepository
}

func NewAddPartyRelationshipHandler(relRepo persistence.PartyRelationshipRepository) *AddPartyRelationshipHandler {
	return &AddPartyRelationshipHandler{relRepo: relRepo}
}

func (h *AddPartyRelationshipHandler) Handle(ctx context.Context, cmd *AddPartyRelationshipCommand) (*domain.PartyRelationship, error) {
	fromID, err := domain.NewPartyID(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid party ID: %w", err)
	}

	toID, err := domain.NewPartyID(cmd.ToPartyID)
	if err != nil {
		return nil, fmt.Errorf("invalid to party ID: %w", err)
	}

	relID, err := domain.NewPartyRelationshipID(cmd.RelationshipID)
	if err != nil {
		return nil, fmt.Errorf("invalid relationship ID: %w", err)
	}

	relType := domain.RelationshipType(strings.ToUpper(cmd.Type))
	if !relType.IsValid() {
		return nil, fmt.Errorf("invalid relationship type: %s", cmd.Type)
	}

	relationship, err := domain.NewPartyRelationship(relID, fromID, toID, relType)
	if err != nil {
		return nil, err
	}

	if err := h.relRepo.Save(ctx, relationship); err != nil {
		return nil, fmt.Errorf("failed to save relationship: %w", err)
	}

	return &relationship, nil
}

// RemovePartyRelationshipCommand represents removing a relationship

type RemovePartyRelationshipCommand struct {
	RelationshipID string
}

// RemovePartyRelationshipHandler handles removing relationships

type RemovePartyRelationshipHandler struct {
	relRepo persistence.PartyRelationshipRepository
}

func NewRemovePartyRelationshipHandler(relRepo persistence.PartyRelationshipRepository) *RemovePartyRelationshipHandler {
	return &RemovePartyRelationshipHandler{relRepo: relRepo}
}

func (h *RemovePartyRelationshipHandler) Handle(ctx context.Context, cmd *RemovePartyRelationshipCommand) error {
	relID, err := domain.NewPartyRelationshipID(cmd.RelationshipID)
	if err != nil {
		return fmt.Errorf("invalid relationship ID: %w", err)
	}

	if err := h.relRepo.Delete(ctx, relID); err != nil {
		return fmt.Errorf("failed to delete relationship: %w", err)
	}

	return nil
}

// AddContactDetailsCommand represents adding contact details to an organization profile

type AddContactDetailsCommand struct {
	PartyID         string
	ContactID       string
	TypeDescription string
	Phone           string
	Email           string
	RelatedPartyID  string
	ModifiedBy      string
}

// AddContactDetailsHandler handles contact details

type AddContactDetailsHandler struct {
	partyRepo persistence.PartyRepository
}

func NewAddContactDetailsHandler(partyRepo persistence.PartyRepository) *AddContactDetailsHandler {
	return &AddContactDetailsHandler{partyRepo: partyRepo}
}

func (h *AddContactDetailsHandler) Handle(ctx context.Context, cmd *AddContactDetailsCommand) (*domain.ContactDetails, error) {
	partyID, err := domain.NewPartyID(cmd.PartyID)
	if err != nil {
		return nil, fmt.Errorf("invalid party ID: %w", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, fmt.Errorf("party not found: %w", err)
	}

	orgProfile := party.OrganizationProfile()
	if orgProfile == nil {
		return nil, fmt.Errorf("organization profile is required to add contact details")
	}

	contactID, err := domain.NewContactDetailsID(cmd.ContactID)
	if err != nil {
		return nil, fmt.Errorf("invalid contact ID: %w", err)
	}

	var phone *domain.Phone
	if cmd.Phone != "" {
		phone, err = domain.NewPhone(cmd.Phone)
		if err != nil {
			return nil, fmt.Errorf("invalid phone: %w", err)
		}
	}

	var email *domain.Email
	if cmd.Email != "" {
		email, err = domain.NewEmail(cmd.Email)
		if err != nil {
			return nil, fmt.Errorf("invalid email: %w", err)
		}
	}

	var relatedPartyID *domain.PartyID
	if cmd.RelatedPartyID != "" {
		pid, err := domain.NewPartyID(cmd.RelatedPartyID)
		if err != nil {
			return nil, fmt.Errorf("invalid related party ID: %w", err)
		}
		relatedPartyID = &pid
	}

	contact, err := domain.NewContactDetails(contactID, cmd.TypeDescription, phone, email, relatedPartyID)
	if err != nil {
		return nil, err
	}

	if err := orgProfile.AddContact(contact); err != nil {
		return nil, err
	}

	if err := h.partyRepo.Save(ctx, party); err != nil {
		return nil, fmt.Errorf("failed to save party: %w", err)
	}

	return contact, nil
}

// UpdateContactDetailsCommand represents updating contact details on an organization profile

type UpdateContactDetailsCommand struct {
	PartyID         string
	ContactID       string
	TypeDescription *string
	Phone           *string
	Email           *string
	RelatedPartyID  *string
	ModifiedBy      string
}

// UpdateContactDetailsHandler handles updating contact details

type UpdateContactDetailsHandler struct {
	partyRepo persistence.PartyRepository
}

func NewUpdateContactDetailsHandler(partyRepo persistence.PartyRepository) *UpdateContactDetailsHandler {
	return &UpdateContactDetailsHandler{partyRepo: partyRepo}
}

func (h *UpdateContactDetailsHandler) Handle(ctx context.Context, cmd *UpdateContactDetailsCommand) (*domain.ContactDetails, error) {
	if cmd.PartyID == "" {
		return nil, fmt.Errorf("party ID cannot be empty")
	}
	if cmd.ContactID == "" {
		return nil, fmt.Errorf("contact ID cannot be empty")
	}
	if cmd.TypeDescription == nil && cmd.Phone == nil && cmd.Email == nil && cmd.RelatedPartyID == nil {
		return nil, fmt.Errorf("at least one field must be provided")
	}

	partyID, err := domain.NewPartyID(cmd.PartyID)
	if err != nil {
		return nil, fmt.Errorf("invalid party ID: %w", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, fmt.Errorf("party not found: %w", err)
	}

	orgProfile := party.OrganizationProfile()
	if orgProfile == nil {
		return nil, fmt.Errorf("organization profile is required to update contact details")
	}

	contactID, err := domain.NewContactDetailsID(cmd.ContactID)
	if err != nil {
		return nil, fmt.Errorf("invalid contact ID: %w", err)
	}

	var existing *domain.ContactDetails
	for _, contact := range orgProfile.Contacts() {
		if contact.ID() == contactID {
			existing = contact
			break
		}
	}
	if existing == nil {
		return nil, fmt.Errorf("contact details not found")
	}

	typeDescription := existing.TypeDescription()
	if cmd.TypeDescription != nil {
		typeDescription = *cmd.TypeDescription
	}

	phone := existing.Phone()
	if cmd.Phone != nil {
		if *cmd.Phone == "" {
			phone = nil
		} else {
			phone, err = domain.NewPhone(*cmd.Phone)
			if err != nil {
				return nil, fmt.Errorf("invalid phone: %w", err)
			}
		}
	}

	email := existing.Email()
	if cmd.Email != nil {
		if *cmd.Email == "" {
			email = nil
		} else {
			email, err = domain.NewEmail(*cmd.Email)
			if err != nil {
				return nil, fmt.Errorf("invalid email: %w", err)
			}
		}
	}

	relatedPartyID := existing.RelatedPartyID()
	if cmd.RelatedPartyID != nil {
		if *cmd.RelatedPartyID == "" {
			relatedPartyID = nil
		} else {
			pid, err := domain.NewPartyID(*cmd.RelatedPartyID)
			if err != nil {
				return nil, fmt.Errorf("invalid related party ID: %w", err)
			}
			relatedPartyID = &pid
		}
	}

	updated, err := domain.NewContactDetails(contactID, typeDescription, phone, email, relatedPartyID)
	if err != nil {
		return nil, err
	}

	if err := orgProfile.UpdateContact(updated); err != nil {
		return nil, err
	}

	if err := h.partyRepo.Save(ctx, party); err != nil {
		return nil, fmt.Errorf("failed to save party: %w", err)
	}

	return updated, nil
}

// RemoveContactDetailsCommand represents removing contact details from an organization profile

type RemoveContactDetailsCommand struct {
	PartyID    string
	ContactID  string
	ModifiedBy string
}

// RemoveContactDetailsHandler handles removing contact details

type RemoveContactDetailsHandler struct {
	partyRepo persistence.PartyRepository
}

func NewRemoveContactDetailsHandler(partyRepo persistence.PartyRepository) *RemoveContactDetailsHandler {
	return &RemoveContactDetailsHandler{partyRepo: partyRepo}
}

func (h *RemoveContactDetailsHandler) Handle(ctx context.Context, cmd *RemoveContactDetailsCommand) error {
	if cmd.PartyID == "" {
		return fmt.Errorf("party ID cannot be empty")
	}
	if cmd.ContactID == "" {
		return fmt.Errorf("contact ID cannot be empty")
	}

	partyID, err := domain.NewPartyID(cmd.PartyID)
	if err != nil {
		return fmt.Errorf("invalid party ID: %w", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return fmt.Errorf("party not found: %w", err)
	}

	orgProfile := party.OrganizationProfile()
	if orgProfile == nil {
		return fmt.Errorf("organization profile is required to remove contact details")
	}

	contactID, err := domain.NewContactDetailsID(cmd.ContactID)
	if err != nil {
		return fmt.Errorf("invalid contact ID: %w", err)
	}

	if err := orgProfile.RemoveContact(contactID); err != nil {
		return err
	}

	if err := h.partyRepo.Save(ctx, party); err != nil {
		return fmt.Errorf("failed to save party: %w", err)
	}

	return nil
}

// AddPartyAddressCommand represents adding an address to a party

type AddPartyAddressCommand struct {
	PartyID    string
	AddressID  string
	Street     string
	City       string
	Province   string
	PostalCode string
	Country    string
	CreatedBy  string
}

// AddPartyAddressHandler handles adding addresses

type AddPartyAddressHandler struct {
	addressRepo persistence.PartyAddressRepository
}

func NewAddPartyAddressHandler(addressRepo persistence.PartyAddressRepository) *AddPartyAddressHandler {
	return &AddPartyAddressHandler{addressRepo: addressRepo}
}

func (h *AddPartyAddressHandler) Handle(ctx context.Context, cmd *AddPartyAddressCommand) (*domain.Address, error) {
	partyID, err := domain.NewPartyID(cmd.PartyID)
	if err != nil {
		return nil, fmt.Errorf("invalid party ID: %w", err)
	}

	addressID, err := domain.NewAddressID(cmd.AddressID)
	if err != nil {
		return nil, fmt.Errorf("invalid address ID: %w", err)
	}

	address, err := domain.NewAddress(cmd.Street, cmd.City, cmd.Province, cmd.PostalCode, cmd.Country)
	if err != nil {
		return nil, err
	}

	if err := h.addressRepo.Save(ctx, address, addressID, partyID, cmd.CreatedBy, cmd.CreatedBy); err != nil {
		return nil, fmt.Errorf("failed to save address: %w", err)
	}

	return address, nil
}

func parsePartyRoleType(role string) (domain.PartyRoleType, error) {
	roleType := domain.PartyRoleType(strings.ToUpper(role))
	if !roleType.IsValid() {
		return "", fmt.Errorf("invalid party role: %s", role)
	}
	return roleType, nil
}
