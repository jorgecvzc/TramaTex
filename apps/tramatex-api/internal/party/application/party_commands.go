package application

import (
	"context"
	"strings"

	"github.com/joran-cortez/tramatex/internal/party/domain"
	"github.com/joran-cortez/tramatex/internal/party/persistence"
)

// Party command inputs

type PersonProfileInput struct {
	FirstName *string
	LastName  *string
	Phone     *string
	Email     *string
}

type ContactDetailsInput struct {
	ID              string
	TypeDescription string
	Phone           string
	Email           string
	RelatedPartyID  string
}

type OrganizationProfileInput struct {
	Name      *string
	TaxID     *string
	TaxIDType *string
	Website   *string
	Phone     *string
	Email     *string
	Notes     *string
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// CreatePartyCommand represents a command to create a Party
// At least one profile must be provided.
type CreatePartyCommand struct {
	ID                        string
	Status                    string
	Roles                     []string
	DefaultDiscountPercentage *float64
	PersonProfile             *PersonProfileInput
	OrganizationProfile       *OrganizationProfileInput
	ActorID                   string
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
		var phone *domain.Phone
		phoneValue := strings.TrimSpace(stringValue(cmd.PersonProfile.Phone))
		if phoneValue != "" {
			phone, err = domain.NewPhone(phoneValue)
			if err != nil {
				return nil, domain.WrapValidation("invalid phone", err)
			}
		}

		var email *domain.Email
		emailValue := strings.TrimSpace(stringValue(cmd.PersonProfile.Email))
		if emailValue != "" {
			email, err = domain.NewEmail(emailValue)
			if err != nil {
				return nil, domain.WrapValidation("invalid email", err)
			}
		}

		personProfile, err = domain.NewPersonProfile(
			stringValue(cmd.PersonProfile.FirstName),
			stringValue(cmd.PersonProfile.LastName),
			phone,
			email,
		)
		if err != nil {
			return nil, domain.WrapValidation("invalid person profile", err)
		}
	}

	var organizationProfile *domain.OrganizationProfile
	if cmd.OrganizationProfile != nil {
		var taxID *domain.TaxID
		taxValue := strings.TrimSpace(stringValue(cmd.OrganizationProfile.TaxID))
		if taxValue != "" {
			taxType := strings.TrimSpace(stringValue(cmd.OrganizationProfile.TaxIDType))
			if taxType == "" {
				taxType = "NIF"
			}
			taxID, err = domain.NewTaxID(taxValue, taxType)
			if err != nil {
				return nil, domain.WrapValidation("invalid tax ID", err)
			}
		}

		var phone *domain.Phone
		phoneValue := strings.TrimSpace(stringValue(cmd.OrganizationProfile.Phone))
		if phoneValue != "" {
			phone, err = domain.NewPhone(phoneValue)
			if err != nil {
				return nil, domain.WrapValidation("invalid phone", err)
			}
		}

		var email *domain.Email
		emailValue := strings.TrimSpace(stringValue(cmd.OrganizationProfile.Email))
		if emailValue != "" {
			email, err = domain.NewEmail(emailValue)
			if err != nil {
				return nil, domain.WrapValidation("invalid email", err)
			}
		}

		organizationProfile, err = domain.NewOrganizationProfile(
			strings.TrimSpace(stringValue(cmd.OrganizationProfile.Name)),
			taxID,
			strings.TrimSpace(stringValue(cmd.OrganizationProfile.Website)),
			phone,
			email,
			strings.TrimSpace(stringValue(cmd.OrganizationProfile.Notes)),
		)
		if err != nil {
			return nil, domain.WrapValidation("invalid organization profile", err)
		}
	}

	party, err := domain.NewParty(partyID, status, personProfile, organizationProfile)
	if err != nil {
		return nil, domain.WrapValidation("failed to create party", err)
	}

	isCustomer := false
	for _, role := range cmd.Roles {
		roleType, err := parsePartyRoleType(role)
		if err != nil {
			return nil, err
		}
		if roleType == domain.PartyRoleClient {
			isCustomer = true
		}
		partyRole, err := domain.NewPartyRole(roleType, nil)
		if err != nil {
			return nil, err
		}
		if err := party.AddRole(partyRole); err != nil {
			return nil, err
		}
	}

	if cmd.DefaultDiscountPercentage != nil {
		if !isCustomer && *cmd.DefaultDiscountPercentage != 0 {
			return nil, domain.NewValidationError("default discount percentage can only be assigned to customers")
		}
		if err := party.SetDefaultDiscountPercentage(*cmd.DefaultDiscountPercentage); err != nil {
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
	ID                        string
	Status                    string
	DefaultDiscountPercentage *float64
	PersonProfile             *PersonProfileInput
	OrganizationProfile       *OrganizationProfileInput
	ActorID                   string
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
		var firstName, lastName string
		var phone *domain.Phone
		var email *domain.Email

		if existing != nil {
			firstName = existing.FirstName()
			lastName = existing.LastName()
			phone = existing.Phone()
			email = existing.Email()
		}

		if cmd.PersonProfile.FirstName != nil {
			firstName = *cmd.PersonProfile.FirstName
		}
		if cmd.PersonProfile.LastName != nil {
			lastName = *cmd.PersonProfile.LastName
		}

		if cmd.PersonProfile.Phone != nil {
			phoneValue := strings.TrimSpace(*cmd.PersonProfile.Phone)
			if phoneValue == "" {
				phone = nil
			} else {
				phone, err = domain.NewPhone(phoneValue)
				if err != nil {
					return nil, domain.WrapValidation("invalid phone", err)
				}
			}
		}

		if cmd.PersonProfile.Email != nil {
			emailValue := strings.TrimSpace(*cmd.PersonProfile.Email)
			if emailValue == "" {
				email = nil
			} else {
				email, err = domain.NewEmail(emailValue)
				if err != nil {
					return nil, domain.WrapValidation("invalid email", err)
				}
			}
		}

		// If both are empty and there was no existing profile, it's an error.
		if firstName == "" && lastName == "" && existing == nil {
			return nil, domain.NewValidationError("person profile cannot be empty")
		}

		profile, err := domain.NewPersonProfile(firstName, lastName, phone, email)
		if err != nil {
			return nil, domain.WrapValidation("invalid person profile", err)
		}
		if err := party.SetPersonProfile(profile); err != nil {
			return nil, err
		}
	}

	if cmd.OrganizationProfile != nil {
		existing := party.OrganizationProfile()
		name := ""
		var taxID *domain.TaxID
		website := ""
		var phone *domain.Phone
		var email *domain.Email
		notes := ""

		if existing != nil {
			name = existing.Name()
			taxID = existing.TaxID()
			website = existing.Website()
			phone = existing.Phone()
			email = existing.Email()
			notes = existing.Notes()
		}

		if cmd.OrganizationProfile.Name != nil {
			name = strings.TrimSpace(*cmd.OrganizationProfile.Name)
		}
		if name == "" {
			return nil, domain.NewValidationError("organization name is required")
		}

		if cmd.OrganizationProfile.TaxID != nil {
			taxValue := strings.TrimSpace(*cmd.OrganizationProfile.TaxID)
			if taxValue == "" {
				taxID = nil
			} else {
				taxType := strings.TrimSpace(stringValue(cmd.OrganizationProfile.TaxIDType))
				if taxType == "" {
					taxType = "NIF"
				}
				taxID, err = domain.NewTaxID(taxValue, taxType)
				if err != nil {
					return nil, domain.WrapValidation("invalid tax ID", err)
				}
			}
		}

		if cmd.OrganizationProfile.Website != nil {
			website = strings.TrimSpace(*cmd.OrganizationProfile.Website)
		}

		if cmd.OrganizationProfile.Phone != nil {
			phoneValue := strings.TrimSpace(*cmd.OrganizationProfile.Phone)
			if phoneValue == "" {
				phone = nil
			} else {
				phone, err = domain.NewPhone(phoneValue)
				if err != nil {
					return nil, domain.WrapValidation("invalid phone", err)
				}
			}
		}

		if cmd.OrganizationProfile.Email != nil {
			emailValue := strings.TrimSpace(*cmd.OrganizationProfile.Email)
			if emailValue == "" {
				email = nil
			} else {
				email, err = domain.NewEmail(emailValue)
				if err != nil {
					return nil, domain.WrapValidation("invalid email", err)
				}
			}
		}

		if cmd.OrganizationProfile.Notes != nil {
			notes = strings.TrimSpace(*cmd.OrganizationProfile.Notes)
		}

		profile, err := domain.NewOrganizationProfile(name, taxID, website, phone, email, notes)
		if err != nil {
			return nil, domain.WrapValidation("invalid organization profile", err)
		}

		if err := party.SetOrganizationProfile(profile); err != nil {
			return nil, err
		}

		// Switching to organization type: clear any stale person profile
		if cmd.PersonProfile == nil {
			party.SetPersonProfile(nil)
		}
	}

	if cmd.PersonProfile != nil && cmd.OrganizationProfile == nil {
		// Switching to person type: clear any stale organization profile
		party.SetOrganizationProfile(nil)
	}

	if cmd.DefaultDiscountPercentage != nil {
		isCustomer := false
		for _, role := range party.Roles() {
			if role.Type() == domain.PartyRoleClient {
				isCustomer = true
				break
			}
		}
		if !isCustomer && *cmd.DefaultDiscountPercentage != 0 {
			return nil, domain.NewValidationError("default discount percentage can only be assigned to customers")
		}
		if err := party.SetDefaultDiscountPercentage(*cmd.DefaultDiscountPercentage); err != nil {
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

// DeletePartyCommand represents a hard delete request for a party
type DeletePartyCommand struct {
	ID      string
	ActorID string
}

// DeletePartyHandler handles deletion of orphan contact parties
type DeletePartyHandler struct {
	partyRepo persistence.PartyRepository
	relRepo   persistence.PartyRelationshipRepository
}

func NewDeletePartyHandler(partyRepo persistence.PartyRepository, relRepo persistence.PartyRelationshipRepository) *DeletePartyHandler {
	return &DeletePartyHandler{partyRepo: partyRepo, relRepo: relRepo}
}

func (h *DeletePartyHandler) CanDelete(ctx context.Context, partyID string) (bool, error) {
	id, err := domain.NewPartyID(partyID)
	if err != nil {
		return false, domain.WrapValidation("invalid party ID", err)
	}
	canDelete, _, err := h.canDeleteWithReason(ctx, id)
	return canDelete, err
}

func (h *DeletePartyHandler) canDeleteWithReason(ctx context.Context, partyID domain.PartyID) (bool, string, error) {
	// Verify party exists
	_, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return false, "", domain.WrapNotFound("party not found", err)
	}

	// Check party relationships
	relationships, err := h.relRepo.FindByPartyID(ctx, partyID)
	if err != nil {
		return false, "", domain.WrapPersistence("failed to load party relationships", err)
	}
	if len(relationships) > 0 {
		return false, "relationships", nil
	}

	hasContactRefs, err := h.partyRepo.HasContactDetailsReferences(ctx, partyID)
	if err != nil {
		return false, "", domain.WrapPersistence("failed to check contact details references", err)
	}
	if hasContactRefs {
		return false, "contact_details", nil
	}

	hasMESRefs, err := h.partyRepo.HasMESWorkReferences(ctx, partyID)
	if err != nil {
		return false, "", domain.WrapPersistence("failed to check MES work references", err)
	}
	if hasMESRefs {
		return false, "mes", nil
	}

	hasSalesRefs, err := h.partyRepo.HasSalesReferences(ctx, partyID)
	if err != nil {
		return false, "", domain.WrapPersistence("failed to check sales references", err)
	}

	if hasSalesRefs {
		return false, "sales", nil
	}

	return true, "", nil
}

func (h *DeletePartyHandler) Handle(ctx context.Context, cmd *DeletePartyCommand) error {
	actorID := strings.TrimSpace(cmd.ActorID)
	if actorID == "" {
		return domain.NewValidationError("actor ID is required")
	}
	if cmd.ID == "" {
		return domain.NewValidationError("party ID cannot be empty")
	}

	partyID, err := domain.NewPartyID(cmd.ID)
	if err != nil {
		return domain.WrapValidation("invalid party ID", err)
	}

	canDelete, reason, err := h.canDeleteWithReason(ctx, partyID)
	if err != nil {
		return err
	}
	if !canDelete {
		switch reason {
		case "relationships":
			return domain.NewValidationError("party is linked to other entities and cannot be deleted")
		case "contact_details":
			return domain.NewValidationError("contact is referenced in organization contact details and cannot be deleted")
		case "mes":
			return domain.NewValidationError("party has MES work records and cannot be deleted")
		case "sales":
			return domain.NewValidationError("party has sales documents (quotes, orders, invoices, or delivery notes) and cannot be deleted")
		default:
			return domain.NewValidationError("party cannot be deleted")
		}
	}

	if err := h.partyRepo.Delete(ctx, partyID); err != nil {
		return domain.WrapPersistence("failed to delete party", err)
	}

	return nil
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
	PartyID              string
	ContactID            string
	ActorID              string
	DeleteIfNoReferences bool // If true, also delete the contact party if it has no other references
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

	// If requested, check if the contact party can be deleted
	if cmd.DeleteIfNoReferences {
		contactPartyID, err := domain.NewPartyID(cmd.ContactID)
		if err != nil {
			// Ignore error, just don't delete
			return nil
		}

		// Check if contact has any other references
		hasRefs, err := h.partyRepo.HasContactDetailsReferences(ctx, contactPartyID)
		if err != nil {
			// Ignore error, just don't delete
			return nil
		}

		// If no references, delete the contact party
		if !hasRefs {
			if err := h.partyRepo.Delete(ctx, contactPartyID); err != nil {
				// Ignore error, the main operation succeeded
				return nil
			}
		}
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

// UpdatePartyAddressCommand represents updating an address of a party
type UpdatePartyAddressCommand struct {
	PartyID    string
	AddressID  string
	Street     string
	City       string
	Province   string
	PostalCode string
	Country    string
	ActorID    string
}

// UpdatePartyAddressHandler handles updating addresses
type UpdatePartyAddressHandler struct {
	addressRepo persistence.PartyAddressRepository
}

func NewUpdatePartyAddressHandler(addressRepo persistence.PartyAddressRepository) *UpdatePartyAddressHandler {
	return &UpdatePartyAddressHandler{addressRepo: addressRepo}
}

func (h *UpdatePartyAddressHandler) Handle(ctx context.Context, cmd *UpdatePartyAddressCommand) (*domain.Address, error) {
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

	// Save method handles both create and update (upsert)
	if err := h.addressRepo.Save(ctx, address, addressID, partyID, actorID, actorID); err != nil {
		return nil, domain.WrapPersistence("failed to update address", err)
	}

	return address, nil
}

// RemovePartyAddressCommand represents removing an address from a party
type RemovePartyAddressCommand struct {
	AddressID string
	ActorID   string
}

// RemovePartyAddressHandler handles address removal
type RemovePartyAddressHandler struct {
	addressRepo persistence.PartyAddressRepository
}

func NewRemovePartyAddressHandler(addressRepo persistence.PartyAddressRepository) *RemovePartyAddressHandler {
	return &RemovePartyAddressHandler{addressRepo: addressRepo}
}

func (h *RemovePartyAddressHandler) Handle(ctx context.Context, cmd *RemovePartyAddressCommand) error {
	actorID := strings.TrimSpace(cmd.ActorID)
	if actorID == "" {
		return domain.NewValidationError("actor ID is required")
	}

	addressID, err := domain.NewAddressID(cmd.AddressID)
	if err != nil {
		return domain.WrapValidation("invalid address ID", err)
	}

	if err := h.addressRepo.Delete(ctx, addressID); err != nil {
		return domain.WrapPersistence("failed to delete address", err)
	}

	return nil
}
