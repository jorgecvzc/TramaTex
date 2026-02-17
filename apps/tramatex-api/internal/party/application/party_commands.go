package application

import (
	"context"
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
	ActorID             string
}

// CreatePartyHandler handles party creation

type CreatePartyHandler struct {
	partyRepo persistence.PartyRepository
}

func NewCreatePartyHandler(partyRepo persistence.PartyRepository) *CreatePartyHandler {
	return &CreatePartyHandler{partyRepo: partyRepo}
}

func (h *CreatePartyHandler) Handle(ctx context.Context, cmd *CreatePartyCommand) (*domain.Party, error) {
	actorID := strings.TrimSpace(cmd.ActorID)
	if actorID == "" {
		return nil, domain.NewValidationError("actor ID is required")
	}
	if cmd.ID == "" {
		return nil, domain.NewValidationError("party ID cannot be empty")
	}

	partyID, err := domain.NewPartyID(cmd.ID)
	if err != nil {
		return nil, domain.WrapValidation("invalid party ID", err)
	}

	status := domain.PartyStatusActive
	if cmd.Status != "" {
		status = domain.PartyStatus(strings.ToUpper(cmd.Status))
		if !status.IsValid() {
			return nil, domain.NewValidationErrorf("invalid party status: %s", cmd.Status)
		}
	}

	var personProfile *domain.PersonProfile
	if cmd.PersonProfile != nil {
		personProfile, err = domain.NewPersonProfile(cmd.PersonProfile.FirstName, cmd.PersonProfile.LastName)
		if err != nil {
			return nil, domain.WrapValidation("invalid person profile", err)
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
				return nil, domain.WrapValidation("invalid tax ID", err)
			}
		}

		organizationProfile, err = domain.NewOrganizationProfile(
			cmd.OrganizationProfile.Name,
			taxID,
			cmd.OrganizationProfile.Website,
		)
		if err != nil {
			return nil, domain.WrapValidation("invalid organization profile", err)
		}

		for _, input := range cmd.OrganizationProfile.Contacts {
			contactID, err := domain.NewContactDetailsID(input.ID)
			if err != nil {
				return nil, domain.WrapValidation("invalid contact details ID", err)
			}

			var phone *domain.Phone
			if input.Phone != "" {
				phone, err = domain.NewPhone(input.Phone)
				if err != nil {
					return nil, domain.WrapValidation("invalid phone", err)
				}
			}

			var email *domain.Email
			if input.Email != "" {
				email, err = domain.NewEmail(input.Email)
				if err != nil {
					return nil, domain.WrapValidation("invalid email", err)
				}
			}

			var relatedPartyID *domain.PartyID
			if input.RelatedPartyID != "" {
				pid, err := domain.NewPartyID(input.RelatedPartyID)
				if err != nil {
					return nil, domain.WrapValidation("invalid related party ID", err)
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
				return nil, domain.WrapValidation("invalid contact details", err)
			}

			if err := organizationProfile.AddContact(contact); err != nil {
				return nil, err
			}
		}
	}

	party, err := domain.NewParty(partyID, status, personProfile, organizationProfile)
	if err != nil {
		return nil, domain.WrapValidation("failed to create party", err)
	}

	for _, role := range cmd.Roles {
		roleType, err := parsePartyRoleType(role)
		if err != nil {
			return nil, err
		}
		partyRole, err := domain.NewPartyRole(roleType, nil)
		if err != nil {
			return nil, err
		}
		if err := party.AddRole(partyRole); err != nil {
			return nil, err
		}
	}

	if err := h.partyRepo.Save(ctx, party, actorID, actorID); err != nil {
		return nil, domain.WrapPersistence("failed to save party", err)
	}

	return party, nil
}

// UpdatePartyCommand represents a command to update Party profiles

type UpdatePartyCommand struct {
	ID                  string
	Status              string
	PersonProfile       *PersonProfileInput
	OrganizationProfile *OrganizationProfileInput
	ActorID             string
}

// UpdatePartyHandler handles party updates

type UpdatePartyHandler struct {
	partyRepo persistence.PartyRepository
}

func NewUpdatePartyHandler(partyRepo persistence.PartyRepository) *UpdatePartyHandler {
	return &UpdatePartyHandler{partyRepo: partyRepo}
}

func (h *UpdatePartyHandler) Handle(ctx context.Context, cmd *UpdatePartyCommand) (*domain.Party, error) {
	actorID := strings.TrimSpace(cmd.ActorID)
	if actorID == "" {
		return nil, domain.NewValidationError("actor ID is required")
	}
	if cmd.ID == "" {
		return nil, domain.NewValidationError("party ID cannot be empty")
	}

	partyID, err := domain.NewPartyID(cmd.ID)
	if err != nil {
		return nil, domain.WrapValidation("invalid party ID", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, domain.WrapNotFound("party not found", err)
	}

	if cmd.Status != "" {
		status := domain.PartyStatus(strings.ToUpper(cmd.Status))
		if !status.IsValid() {
			return nil, domain.NewValidationErrorf("invalid party status: %s", cmd.Status)
		}
		if status == domain.PartyStatusActive {
			if err := party.Activate(); err != nil {
				return nil, err
			}
		} else {
			if err := party.Deactivate(); err != nil {
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
			return nil, domain.WrapValidation("invalid person profile", err)
		}
		if err := party.SetPersonProfile(profile); err != nil {
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
			return nil, domain.NewValidationError("organization name is required")
		}

		var taxID *domain.TaxID
		if cmd.OrganizationProfile.TaxID != "" {
			taxType := cmd.OrganizationProfile.TaxIDType
			if taxType == "" {
				taxType = "NIF"
			}
			taxID, err = domain.NewTaxID(cmd.OrganizationProfile.TaxID, taxType)
			if err != nil {
				return nil, domain.WrapValidation("invalid tax ID", err)
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
			return nil, domain.WrapValidation("invalid organization profile", err)
		}

		if err := party.SetOrganizationProfile(profile); err != nil {
			return nil, err
		}
	}

	if err := h.partyRepo.Save(ctx, party, actorID, actorID); err != nil {
		return nil, domain.WrapPersistence("failed to save party", err)
	}

	return party, nil
}

// ChangePartyStatusCommand represents a status change

type ChangePartyStatusCommand struct {
	ID      string
	Status  string
	ActorID string
}

// ChangePartyStatusHandler handles status changes

type ChangePartyStatusHandler struct {
	partyRepo persistence.PartyRepository
}

func NewChangePartyStatusHandler(partyRepo persistence.PartyRepository) *ChangePartyStatusHandler {
	return &ChangePartyStatusHandler{partyRepo: partyRepo}
}

func (h *ChangePartyStatusHandler) Handle(ctx context.Context, cmd *ChangePartyStatusCommand) (*domain.Party, error) {
	actorID := strings.TrimSpace(cmd.ActorID)
	if actorID == "" {
		return nil, domain.NewValidationError("actor ID is required")
	}
	if cmd.ID == "" {
		return nil, domain.NewValidationError("party ID cannot be empty")
	}

	partyID, err := domain.NewPartyID(cmd.ID)
	if err != nil {
		return nil, domain.WrapValidation("invalid party ID", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, domain.WrapNotFound("party not found", err)
	}

	status := domain.PartyStatus(strings.ToUpper(cmd.Status))
	if !status.IsValid() {
		return nil, domain.NewValidationErrorf("invalid party status: %s", cmd.Status)
	}

	if status == domain.PartyStatusActive {
		if err := party.Activate(); err != nil {
			return nil, err
		}
	} else {
		if err := party.Deactivate(); err != nil {
			return nil, err
		}
	}

	if err := h.partyRepo.Save(ctx, party, actorID, actorID); err != nil {
		return nil, domain.WrapPersistence("failed to save party", err)
	}

	return party, nil
}

// AddPartyRoleCommand represents adding a role to a party

type AddPartyRoleCommand struct {
	ID      string
	Role    string
	ActorID string
}

// AddPartyRoleHandler handles adding a role

type AddPartyRoleHandler struct {
	partyRepo persistence.PartyRepository
}

func NewAddPartyRoleHandler(partyRepo persistence.PartyRepository) *AddPartyRoleHandler {
	return &AddPartyRoleHandler{partyRepo: partyRepo}
}

func (h *AddPartyRoleHandler) Handle(ctx context.Context, cmd *AddPartyRoleCommand) (*domain.Party, error) {
	actorID := strings.TrimSpace(cmd.ActorID)
	if actorID == "" {
		return nil, domain.NewValidationError("actor ID is required")
	}
	partyID, err := domain.NewPartyID(cmd.ID)
	if err != nil {
		return nil, domain.WrapValidation("invalid party ID", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, domain.WrapNotFound("party not found", err)
	}

	roleType, err := parsePartyRoleType(cmd.Role)
	if err != nil {
		return nil, err
	}

	role, err := domain.NewPartyRole(roleType, nil)
	if err != nil {
		return nil, err
	}
	if err := party.AddRole(role); err != nil {
		return nil, err
	}

	if err := h.partyRepo.Save(ctx, party, actorID, actorID); err != nil {
		return nil, domain.WrapPersistence("failed to save party", err)
	}

	return party, nil
}

// RemovePartyRoleCommand represents removing a role from a party

type RemovePartyRoleCommand struct {
	ID      string
	Role    string
	ActorID string
}

// RemovePartyRoleHandler handles removing a role

type RemovePartyRoleHandler struct {
	partyRepo persistence.PartyRepository
}

func NewRemovePartyRoleHandler(partyRepo persistence.PartyRepository) *RemovePartyRoleHandler {
	return &RemovePartyRoleHandler{partyRepo: partyRepo}
}

func (h *RemovePartyRoleHandler) Handle(ctx context.Context, cmd *RemovePartyRoleCommand) (*domain.Party, error) {
	actorID := strings.TrimSpace(cmd.ActorID)
	if actorID == "" {
		return nil, domain.NewValidationError("actor ID is required")
	}
	partyID, err := domain.NewPartyID(cmd.ID)
	if err != nil {
		return nil, domain.WrapValidation("invalid party ID", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, domain.WrapNotFound("party not found", err)
	}

	roleType, err := parsePartyRoleType(cmd.Role)
	if err != nil {
		return nil, err
	}

	if err := party.RemoveRole(roleType); err != nil {
		return nil, err
	}

	if err := h.partyRepo.Save(ctx, party, actorID, actorID); err != nil {
		return nil, domain.WrapPersistence("failed to save party", err)
	}

	return party, nil
}

// AddPartyRelationshipCommand represents adding a relationship

type AddPartyRelationshipCommand struct {
	ID             string
	RelationshipID string
	ToPartyID      string
	Type           string
	ActorID        string
}

// AddPartyRelationshipHandler handles adding relationships

type AddPartyRelationshipHandler struct {
	relRepo persistence.PartyRelationshipRepository
}

func NewAddPartyRelationshipHandler(relRepo persistence.PartyRelationshipRepository) *AddPartyRelationshipHandler {
	return &AddPartyRelationshipHandler{relRepo: relRepo}
}

func (h *AddPartyRelationshipHandler) Handle(ctx context.Context, cmd *AddPartyRelationshipCommand) (*domain.PartyRelationship, error) {
	actorID := strings.TrimSpace(cmd.ActorID)
	if actorID == "" {
		return nil, domain.NewValidationError("actor ID is required")
	}
	fromID, err := domain.NewPartyID(cmd.ID)
	if err != nil {
		return nil, domain.WrapValidation("invalid party ID", err)
	}

	toID, err := domain.NewPartyID(cmd.ToPartyID)
	if err != nil {
		return nil, domain.WrapValidation("invalid to party ID", err)
	}

	relID, err := domain.NewPartyRelationshipID(cmd.RelationshipID)
	if err != nil {
		return nil, domain.WrapValidation("invalid relationship ID", err)
	}

	relType := domain.RelationshipType(strings.ToUpper(cmd.Type))
	if !relType.IsValid() {
		return nil, domain.NewValidationErrorf("invalid relationship type: %s", cmd.Type)
	}

	relationship, err := domain.NewPartyRelationship(relID, fromID, toID, relType)
	if err != nil {
		return nil, err
	}

	if err := h.relRepo.Save(ctx, relationship, actorID, actorID); err != nil {
		return nil, domain.WrapPersistence("failed to save relationship", err)
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
		return domain.WrapValidation("invalid relationship ID", err)
	}

	if err := h.relRepo.Delete(ctx, relID); err != nil {
		return domain.WrapPersistence("failed to delete relationship", err)
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
	ActorID         string
}

// AddContactDetailsHandler handles contact details

type AddContactDetailsHandler struct {
	partyRepo persistence.PartyRepository
}

func NewAddContactDetailsHandler(partyRepo persistence.PartyRepository) *AddContactDetailsHandler {
	return &AddContactDetailsHandler{partyRepo: partyRepo}
}

func (h *AddContactDetailsHandler) Handle(ctx context.Context, cmd *AddContactDetailsCommand) (*domain.ContactDetails, error) {
	actorID := strings.TrimSpace(cmd.ActorID)
	if actorID == "" {
		return nil, domain.NewValidationError("actor ID is required")
	}
	partyID, err := domain.NewPartyID(cmd.PartyID)
	if err != nil {
		return nil, domain.WrapValidation("invalid party ID", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, domain.WrapNotFound("party not found", err)
	}

	orgProfile := party.OrganizationProfile()
	if orgProfile == nil {
		return nil, domain.NewValidationError("organization profile is required to add contact details")
	}

	contactID, err := domain.NewContactDetailsID(cmd.ContactID)
	if err != nil {
		return nil, domain.WrapValidation("invalid contact ID", err)
	}

	var phone *domain.Phone
	if cmd.Phone != "" {
		phone, err = domain.NewPhone(cmd.Phone)
		if err != nil {
			return nil, domain.WrapValidation("invalid phone", err)
		}
	}

	var email *domain.Email
	if cmd.Email != "" {
		email, err = domain.NewEmail(cmd.Email)
		if err != nil {
			return nil, domain.WrapValidation("invalid email", err)
		}
	}

	var relatedPartyID *domain.PartyID
	if cmd.RelatedPartyID != "" {
		pid, err := domain.NewPartyID(cmd.RelatedPartyID)
		if err != nil {
			return nil, domain.WrapValidation("invalid related party ID", err)
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

	if err := h.partyRepo.Save(ctx, party, actorID, actorID); err != nil {
		return nil, domain.WrapPersistence("failed to save party", err)
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
	ActorID         string
}

// UpdateContactDetailsHandler handles updating contact details

type UpdateContactDetailsHandler struct {
	partyRepo persistence.PartyRepository
}

func NewUpdateContactDetailsHandler(partyRepo persistence.PartyRepository) *UpdateContactDetailsHandler {
	return &UpdateContactDetailsHandler{partyRepo: partyRepo}
}

func (h *UpdateContactDetailsHandler) Handle(ctx context.Context, cmd *UpdateContactDetailsCommand) (*domain.ContactDetails, error) {
	actorID := strings.TrimSpace(cmd.ActorID)
	if actorID == "" {
		return nil, domain.NewValidationError("actor ID is required")
	}
	if cmd.PartyID == "" {
		return nil, domain.NewValidationError("party ID cannot be empty")
	}
	if cmd.ContactID == "" {
		return nil, domain.NewValidationError("contact ID cannot be empty")
	}
	if cmd.TypeDescription == nil && cmd.Phone == nil && cmd.Email == nil && cmd.RelatedPartyID == nil {
		return nil, domain.NewValidationError("at least one field must be provided")
	}

	partyID, err := domain.NewPartyID(cmd.PartyID)
	if err != nil {
		return nil, domain.WrapValidation("invalid party ID", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, domain.WrapNotFound("party not found", err)
	}

	orgProfile := party.OrganizationProfile()
	if orgProfile == nil {
		return nil, domain.NewValidationError("organization profile is required to update contact details")
	}

	contactID, err := domain.NewContactDetailsID(cmd.ContactID)
	if err != nil {
		return nil, domain.WrapValidation("invalid contact ID", err)
	}

	var existing *domain.ContactDetails
	for _, contact := range orgProfile.Contacts() {
		if contact.ID() == contactID {
			existing = contact
			break
		}
	}
	if existing == nil {
		return nil, domain.NewNotFoundError("contact details not found")
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
				return nil, domain.WrapValidation("invalid phone", err)
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
				return nil, domain.WrapValidation("invalid email", err)
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
				return nil, domain.WrapValidation("invalid related party ID", err)
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

	if err := h.partyRepo.Save(ctx, party, actorID, actorID); err != nil {
		return nil, domain.WrapPersistence("failed to save party", err)
	}

	return updated, nil
}

// RemoveContactDetailsCommand represents removing contact details from an organization profile

type RemoveContactDetailsCommand struct {
	PartyID   string
	ContactID string
	ActorID   string
}

// RemoveContactDetailsHandler handles removing contact details

type RemoveContactDetailsHandler struct {
	partyRepo persistence.PartyRepository
}

func NewRemoveContactDetailsHandler(partyRepo persistence.PartyRepository) *RemoveContactDetailsHandler {
	return &RemoveContactDetailsHandler{partyRepo: partyRepo}
}

func (h *RemoveContactDetailsHandler) Handle(ctx context.Context, cmd *RemoveContactDetailsCommand) error {
	actorID := strings.TrimSpace(cmd.ActorID)
	if actorID == "" {
		return domain.NewValidationError("actor ID is required")
	}
	if cmd.PartyID == "" {
		return domain.NewValidationError("party ID cannot be empty")
	}
	if cmd.ContactID == "" {
		return domain.NewValidationError("contact ID cannot be empty")
	}

	partyID, err := domain.NewPartyID(cmd.PartyID)
	if err != nil {
		return domain.WrapValidation("invalid party ID", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return domain.WrapNotFound("party not found", err)
	}

	orgProfile := party.OrganizationProfile()
	if orgProfile == nil {
		return domain.NewValidationError("organization profile is required to remove contact details")
	}

	contactID, err := domain.NewContactDetailsID(cmd.ContactID)
	if err != nil {
		return domain.WrapValidation("invalid contact ID", err)
	}

	if err := orgProfile.RemoveContact(contactID); err != nil {
		return err
	}

	if err := h.partyRepo.Save(ctx, party, actorID, actorID); err != nil {
		return domain.WrapPersistence("failed to save party", err)
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
	ActorID    string
}

// AddPartyAddressHandler handles adding addresses

type AddPartyAddressHandler struct {
	addressRepo persistence.PartyAddressRepository
}

func NewAddPartyAddressHandler(addressRepo persistence.PartyAddressRepository) *AddPartyAddressHandler {
	return &AddPartyAddressHandler{addressRepo: addressRepo}
}

func (h *AddPartyAddressHandler) Handle(ctx context.Context, cmd *AddPartyAddressCommand) (*domain.Address, error) {
	actorID := strings.TrimSpace(cmd.ActorID)
	if actorID == "" {
		return nil, domain.NewValidationError("actor ID is required")
	}
	partyID, err := domain.NewPartyID(cmd.PartyID)
	if err != nil {
		return nil, domain.WrapValidation("invalid party ID", err)
	}

	addressID, err := domain.NewAddressID(cmd.AddressID)
	if err != nil {
		return nil, domain.WrapValidation("invalid address ID", err)
	}

	address, err := domain.NewAddress(cmd.Street, cmd.City, cmd.Province, cmd.PostalCode, cmd.Country)
	if err != nil {
		return nil, err
	}

	if err := h.addressRepo.Save(ctx, address, addressID, partyID, actorID, actorID); err != nil {
		return nil, domain.WrapPersistence("failed to save address", err)
	}

	return address, nil
}

func parsePartyRoleType(role string) (domain.PartyRoleType, error) {
	roleType := domain.PartyRoleType(strings.ToUpper(role))
	if !roleType.IsValid() {
		return "", domain.NewValidationErrorf("invalid party role: %s", role)
	}
	return roleType, nil
}
