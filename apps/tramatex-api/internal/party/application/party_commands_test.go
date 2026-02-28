package application

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/joran-cortez/tramatex/internal/party/domain"
	"github.com/joran-cortez/tramatex/internal/party/persistence"
)

func strPtr(value string) *string {
	return &value
}

type fakePartyRepo struct {
	parties map[string]*domain.Party
}

func newFakePartyRepo() *fakePartyRepo {
	return &fakePartyRepo{parties: make(map[string]*domain.Party)}
}

func (r *fakePartyRepo) Save(ctx context.Context, party *domain.Party, createdBy string, modifiedBy string) error {
	if party == nil {
		return fmt.Errorf("party cannot be nil")
	}
	r.parties[party.ID().String()] = party
	return nil
}

func (r *fakePartyRepo) FindByID(ctx context.Context, id domain.PartyID) (*domain.Party, error) {
	party, ok := r.parties[id.String()]
	if !ok {
		return nil, fmt.Errorf("party not found")
	}
	return party, nil
}

func (r *fakePartyRepo) FindAll(ctx context.Context, filters *persistence.PartyFilters) ([]*domain.Party, error) {
	result := make([]*domain.Party, 0, len(r.parties))
	for _, party := range r.parties {
		result = append(result, party)
	}
	return result, nil
}

func (r *fakePartyRepo) Delete(ctx context.Context, id domain.PartyID) error {
	delete(r.parties, id.String())
	return nil
}

func (r *fakePartyRepo) Exists(ctx context.Context, id domain.PartyID) (bool, error) {
	_, ok := r.parties[id.String()]
	return ok, nil
}

func (r *fakePartyRepo) Count(ctx context.Context) (int64, error) {
	return int64(len(r.parties)), nil
}

func (r *fakePartyRepo) HasContactDetailsReferences(ctx context.Context, partyID domain.PartyID) (bool, error) {
	return false, nil
}

func (r *fakePartyRepo) HasMESWorkReferences(ctx context.Context, partyID domain.PartyID) (bool, error) {
	return false, nil
}

func (r *fakePartyRepo) HasSalesReferences(ctx context.Context, partyID domain.PartyID) (bool, error) {
	return false, nil
}

type errorPartyRepo struct {
	inner     *fakePartyRepo
	saveErr   error
	findErr   error
	existsErr error
}

func newErrorPartyRepo(saveErr, findErr error) *errorPartyRepo {
	return &errorPartyRepo{inner: newFakePartyRepo(), saveErr: saveErr, findErr: findErr}
}

func (r *errorPartyRepo) Save(ctx context.Context, party *domain.Party, createdBy string, modifiedBy string) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	return r.inner.Save(ctx, party, createdBy, modifiedBy)
}

func (r *errorPartyRepo) FindByID(ctx context.Context, id domain.PartyID) (*domain.Party, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	return r.inner.FindByID(ctx, id)
}

func (r *errorPartyRepo) FindAll(ctx context.Context, filters *persistence.PartyFilters) ([]*domain.Party, error) {
	return r.inner.FindAll(ctx, filters)
}

func (r *errorPartyRepo) Delete(ctx context.Context, id domain.PartyID) error {
	return r.inner.Delete(ctx, id)
}

func (r *errorPartyRepo) Exists(ctx context.Context, id domain.PartyID) (bool, error) {
	if r.existsErr != nil {
		return false, r.existsErr
	}
	return r.inner.Exists(ctx, id)
}

func (r *errorPartyRepo) Count(ctx context.Context) (int64, error) {
	return r.inner.Count(ctx)
}

func (r *errorPartyRepo) HasContactDetailsReferences(ctx context.Context, partyID domain.PartyID) (bool, error) {
	return r.inner.HasContactDetailsReferences(ctx, partyID)
}

func (r *errorPartyRepo) HasMESWorkReferences(ctx context.Context, partyID domain.PartyID) (bool, error) {
	return r.inner.HasMESWorkReferences(ctx, partyID)
}

func (r *errorPartyRepo) HasSalesReferences(ctx context.Context, partyID domain.PartyID) (bool, error) {
	return r.inner.HasSalesReferences(ctx, partyID)
}

type fakeRelationshipRepo struct {
	relationships map[string]domain.PartyRelationship
}

func newFakeRelationshipRepo() *fakeRelationshipRepo {
	return &fakeRelationshipRepo{relationships: make(map[string]domain.PartyRelationship)}
}

func (r *fakeRelationshipRepo) Save(ctx context.Context, relationship domain.PartyRelationship, createdBy string, modifiedBy string) error {
	r.relationships[relationship.ID().String()] = relationship
	return nil
}

func (r *fakeRelationshipRepo) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]domain.PartyRelationship, error) {
	result := make([]domain.PartyRelationship, 0)
	for _, rel := range r.relationships {
		if rel.FromID() == partyID || rel.ToID() == partyID {
			result = append(result, rel)
		}
	}
	return result, nil
}

func (r *fakeRelationshipRepo) Delete(ctx context.Context, id domain.PartyRelationshipID) error {
	delete(r.relationships, id.String())
	return nil
}

type errorRelationshipRepo struct {
	inner     *fakeRelationshipRepo
	saveErr   error
	deleteErr error
	findErr   error
}

func newErrorRelationshipRepo(saveErr, deleteErr, findErr error) *errorRelationshipRepo {
	return &errorRelationshipRepo{inner: newFakeRelationshipRepo(), saveErr: saveErr, deleteErr: deleteErr, findErr: findErr}
}

func (r *errorRelationshipRepo) Save(ctx context.Context, relationship domain.PartyRelationship, createdBy string, modifiedBy string) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	return r.inner.Save(ctx, relationship, createdBy, modifiedBy)
}

func (r *errorRelationshipRepo) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]domain.PartyRelationship, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	return r.inner.FindByPartyID(ctx, partyID)
}

func (r *errorRelationshipRepo) Delete(ctx context.Context, id domain.PartyRelationshipID) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}
	return r.inner.Delete(ctx, id)
}

type fakePartyAddressRepo struct {
	addressesByParty map[string][]*persistence.AddressWithID
	addressesByID    map[string]*domain.Address
}

func newFakePartyAddressRepo() *fakePartyAddressRepo {
	return &fakePartyAddressRepo{
		addressesByParty: make(map[string][]*persistence.AddressWithID),
		addressesByID:    make(map[string]*domain.Address),
	}
}

func (r *fakePartyAddressRepo) Save(ctx context.Context, address *domain.Address, addressID domain.AddressID, partyID domain.PartyID, createdBy string, modifiedBy string) error {
	if address == nil {
		return fmt.Errorf("address cannot be nil")
	}
	r.addressesByID[addressID.String()] = address
	r.addressesByParty[partyID.String()] = append(r.addressesByParty[partyID.String()], &persistence.AddressWithID{
		ID:      addressID.String(),
		Address: address,
	})
	return nil
}

func (r *fakePartyAddressRepo) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]*persistence.AddressWithID, error) {
	return r.addressesByParty[partyID.String()], nil
}

func (r *fakePartyAddressRepo) FindPrimary(ctx context.Context, partyID domain.PartyID) (*domain.Address, error) {
	addresses := r.addressesByParty[partyID.String()]
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no addresses")
	}
	return addresses[0].Address, nil
}

func (r *fakePartyAddressRepo) Delete(ctx context.Context, id domain.AddressID) error {
	delete(r.addressesByID, id.String())
	return nil
}

type errorPartyAddressRepo struct {
	saveErr error
}

func (r *errorPartyAddressRepo) Save(ctx context.Context, address *domain.Address, addressID domain.AddressID, partyID domain.PartyID, createdBy string, modifiedBy string) error {
	return r.saveErr
}

func (r *errorPartyAddressRepo) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]*persistence.AddressWithID, error) {
	return nil, nil
}

func (r *errorPartyAddressRepo) FindPrimary(ctx context.Context, partyID domain.PartyID) (*domain.Address, error) {
	return nil, nil
}

func (r *errorPartyAddressRepo) Delete(ctx context.Context, id domain.AddressID) error {
	return nil
}

func seedPartyWithProfiles(t *testing.T, repo *fakePartyRepo, id string, withOrg bool) *domain.Party {
	t.Helper()
	partyID, _ := domain.NewPartyID(id)
	personProfile, _ := domain.NewPersonProfile("Ana", "Perez", nil, nil)
	var orgProfile *domain.OrganizationProfile
	if withOrg {
		orgProfile, _ = domain.NewOrganizationProfile("Org", nil, "", nil, nil)
	}
	party, err := domain.NewParty(partyID, domain.PartyStatusActive, personProfile, orgProfile)
	if err != nil {
		t.Fatalf("Failed to create party: %v", err)
	}
	if err := repo.Save(context.Background(), party, "user-1", "user-1"); err != nil {
		t.Fatalf("Failed to seed party: %v", err)
	}
	return party
}

func TestCreatePartyHandler_Success(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewCreatePartyHandler(repo)
	ctx := context.Background()

	cmd := &CreatePartyCommand{
		ID:      "party-001",
		Status:  "ACTIVE",
		Roles:   []string{"CLIENT"},
		ActorID: "user-1",
		PersonProfile: &PersonProfileInput{
			FirstName: strPtr("Ana"),
			LastName:  strPtr("Perez"),
		},
		OrganizationProfile: &OrganizationProfileInput{
			Name:      strPtr("Textiles Perez"),
			TaxID:     strPtr("B12345678"),
			TaxIDType: strPtr("CIF"),
			Website:   strPtr("https://textiles.local"),
		},
	}

	party, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if party.OrganizationProfile() == nil || party.PersonProfile() == nil {
		t.Fatalf("Expected both profiles to be set")
	}
	if len(party.Roles()) != 1 {
		t.Fatalf("Expected 1 role, got %d", len(party.Roles()))
	}
	if len(party.OrganizationProfile().Contacts()) != 0 {
		t.Fatalf("Expected 0 contacts, got %d", len(party.OrganizationProfile().Contacts()))
	}
}

func TestCreatePartyHandler_InvalidInputs(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewCreatePartyHandler(repo)

	if _, err := handler.Handle(context.Background(), &CreatePartyCommand{ID: "", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for empty party ID")
	}
	if _, err := handler.Handle(context.Background(), &CreatePartyCommand{ID: "party-x", ActorID: ""}); err == nil {
		t.Fatalf("expected error for empty actor ID")
	}
	if _, err := handler.Handle(context.Background(), &CreatePartyCommand{
		ID:            "party-x",
		Status:        "UNKNOWN",
		ActorID:       "user-1",
		PersonProfile: &PersonProfileInput{FirstName: strPtr("Ana"), LastName: strPtr("Perez")},
	}); err == nil {
		t.Fatalf("expected error for invalid status")
	}
}

func TestCreatePartyHandler_InvalidPersonProfile(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewCreatePartyHandler(repo)

	cmd := &CreatePartyCommand{
		ID:            "party-person",
		ActorID:       "user-1",
		PersonProfile: &PersonProfileInput{FirstName: strPtr(""), LastName: strPtr("Perez")},
	}

	if _, err := handler.Handle(context.Background(), cmd); err == nil {
		t.Fatalf("expected error for invalid person profile")
	}
}

func TestCreatePartyHandler_RequiresProfile(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewCreatePartyHandler(repo)

	cmd := &CreatePartyCommand{
		ID:      "party-no-profile",
		ActorID: "user-1",
	}

	if _, err := handler.Handle(context.Background(), cmd); err == nil {
		t.Fatalf("expected error for missing profiles")
	}
}

func TestCreatePartyHandler_InvalidTaxID(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewCreatePartyHandler(repo)

	cmd := &CreatePartyCommand{
		ID:      "party-bad-tax",
		ActorID: "user-1",
		OrganizationProfile: &OrganizationProfileInput{
			Name:  strPtr("Org"),
			TaxID: strPtr("BAD!"),
		},
	}

	if _, err := handler.Handle(context.Background(), cmd); err == nil {
		t.Fatalf("expected error for invalid tax ID")
	}
}

func TestCreatePartyHandler_SaveError(t *testing.T) {
	repo := newErrorPartyRepo(fmt.Errorf("save failed"), nil)
	handler := NewCreatePartyHandler(repo)

	cmd := &CreatePartyCommand{
		ID:            "party-save",
		ActorID:       "user-1",
		PersonProfile: &PersonProfileInput{FirstName: strPtr("Ana"), LastName: strPtr("Perez")},
	}

	if _, err := handler.Handle(context.Background(), cmd); err == nil {
		t.Fatalf("expected save error")
	}
}

func TestUpdatePartyHandler_UpdatesProfilesAndStatus(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-002", false)

	handler := NewUpdatePartyHandler(repo)
	cmd := &UpdatePartyCommand{
		ID:      "party-002",
		Status:  "INACTIVE",
		ActorID: "user-2",
		OrganizationProfile: &OrganizationProfileInput{
			Name: strPtr("Org Updated"),
		},
	}

	party, err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if party.Status() != domain.PartyStatusInactive {
		t.Fatalf("Expected status INACTIVE, got %s", party.Status())
	}
	if party.OrganizationProfile() == nil || party.OrganizationProfile().Name() != "Org Updated" {
		t.Fatalf("Organization profile not updated")
	}
}

func TestUpdatePartyHandler_ActivateStatus(t *testing.T) {
	repo := newFakePartyRepo()
	partyID, _ := domain.NewPartyID("party-activate")
	personProfile, _ := domain.NewPersonProfile("Ana", "Perez", nil, nil)
	party, _ := domain.NewParty(partyID, domain.PartyStatusInactive, personProfile, nil)
	_ = repo.Save(context.Background(), party, "user-1", "user-1")

	handler := NewUpdatePartyHandler(repo)
	updated, err := handler.Handle(context.Background(), &UpdatePartyCommand{ID: "party-activate", Status: "ACTIVE", ActorID: "user-1"})
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if updated.Status() != domain.PartyStatusActive {
		t.Fatalf("expected status ACTIVE, got %s", updated.Status())
	}
}

func TestUpdatePartyHandler_InvalidInputs(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewUpdatePartyHandler(repo)

	if _, err := handler.Handle(context.Background(), &UpdatePartyCommand{ID: "", ActorID: "user"}); err == nil {
		t.Fatalf("expected error for empty party ID")
	}
	if _, err := handler.Handle(context.Background(), &UpdatePartyCommand{ID: "party-x", ActorID: ""}); err == nil {
		t.Fatalf("expected error for empty actor ID")
	}
}

func TestUpdatePartyHandler_InvalidStatus(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-invalid-status", false)

	handler := NewUpdatePartyHandler(repo)
	if _, err := handler.Handle(context.Background(), &UpdatePartyCommand{ID: "party-invalid-status", Status: "UNKNOWN", ActorID: "user"}); err == nil {
		t.Fatalf("expected error for invalid status")
	}
}

func TestUpdatePartyHandler_OrganizationNameRequired(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-org", false)

	handler := NewUpdatePartyHandler(repo)
	cmd := &UpdatePartyCommand{
		ID:                  "party-org",
		OrganizationProfile: &OrganizationProfileInput{},
		ActorID:             "user-1",
	}
	if _, err := handler.Handle(context.Background(), cmd); err == nil {
		t.Fatalf("expected error for missing organization name")
	}
}

func TestAddContactDetailsHandler_RequiresOrganizationProfile(t *testing.T) {
	repo := newFakePartyRepo()
	party := seedPartyWithProfiles(t, repo, "party-no-org", false)

	handler := NewAddContactDetailsHandler(repo)
	cmd := &AddContactDetailsCommand{
		PartyID:         party.ID().String(),
		ContactID:       "contact-01",
		TypeDescription: "Sales",
		ActorID:         "user-1",
	}

	if _, err := handler.Handle(context.Background(), cmd); err == nil {
		t.Fatalf("expected error when organization profile is missing")
	}
}

func TestUpdateContactDetailsHandler_RequiresFields(t *testing.T) {
	repo := newFakePartyRepo()
	partyID, _ := domain.NewPartyID("party-contact-fields")
	orgProfile, _ := domain.NewOrganizationProfile("Org", nil, "", nil, nil)
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, nil, orgProfile)
	_ = repo.Save(context.Background(), party, "user-1", "user-1")

	handler := NewUpdateContactDetailsHandler(repo)
	cmd := &UpdateContactDetailsCommand{
		PartyID:   partyID.String(),
		ContactID: "contact-01",
		ActorID:   "user-1",
	}

	if _, err := handler.Handle(context.Background(), cmd); err == nil {
		t.Fatalf("expected error when no fields provided")
	}
}

func TestUpdateContactDetailsHandler_NotFound(t *testing.T) {
	repo := newFakePartyRepo()
	partyID, _ := domain.NewPartyID("party-contact-missing")
	orgProfile, _ := domain.NewOrganizationProfile("Org", nil, "", nil, nil)
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, nil, orgProfile)
	_ = repo.Save(context.Background(), party, "user-1", "user-1")

	handler := NewUpdateContactDetailsHandler(repo)
	desc := "Support"
	cmd := &UpdateContactDetailsCommand{
		PartyID:         partyID.String(),
		ContactID:       "missing-contact",
		TypeDescription: &desc,
		ActorID:         "user-1",
	}

	if _, err := handler.Handle(context.Background(), cmd); err == nil {
		t.Fatalf("expected error for missing contact details")
	}
}

func TestRemoveContactDetailsHandler_RequiresOrganizationProfile(t *testing.T) {
	repo := newFakePartyRepo()
	party := seedPartyWithProfiles(t, repo, "party-remove-contact", false)

	handler := NewRemoveContactDetailsHandler(repo)
	cmd := &RemoveContactDetailsCommand{
		PartyID:   party.ID().String(),
		ContactID: "contact-01",
		ActorID:   "user-1",
	}

	if err := handler.Handle(context.Background(), cmd); err == nil {
		t.Fatalf("expected error when organization profile is missing")
	}
}

func TestUpdatePartyHandler_SaveError(t *testing.T) {
	repo := newErrorPartyRepo(fmt.Errorf("save failed"), nil)
	seedPartyWithProfiles(t, repo.inner, "party-save-error", false)

	handler := NewUpdatePartyHandler(repo)
	cmd := &UpdatePartyCommand{
		ID:      "party-save-error",
		Status:  "INACTIVE",
		ActorID: "user-1",
	}
	if _, err := handler.Handle(context.Background(), cmd); err == nil {
		t.Fatalf("expected save error")
	}
}

func TestChangePartyStatusHandler_Success(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-003", false)

	handler := NewChangePartyStatusHandler(repo)
	cmd := &ChangePartyStatusCommand{
		ID:      "party-003",
		Status:  "INACTIVE",
		ActorID: "user-2",
	}

	party, err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if party.Status() != domain.PartyStatusInactive {
		t.Fatalf("Expected status INACTIVE, got %s", party.Status())
	}
}

func TestChangePartyStatusHandler_InvalidInputs(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewChangePartyStatusHandler(repo)

	if _, err := handler.Handle(context.Background(), &ChangePartyStatusCommand{ID: "", ActorID: "user"}); err == nil {
		t.Fatalf("expected error for empty party ID")
	}
	if _, err := handler.Handle(context.Background(), &ChangePartyStatusCommand{ID: "party-x", ActorID: ""}); err == nil {
		t.Fatalf("expected error for empty actor ID")
	}
}

func TestChangePartyStatusHandler_InvalidStatus(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-invalid-status", false)
	handler := NewChangePartyStatusHandler(repo)

	if _, err := handler.Handle(context.Background(), &ChangePartyStatusCommand{ID: "party-invalid-status", Status: "BAD", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for invalid status")
	}
}

func TestChangePartyStatusHandler_SaveError(t *testing.T) {
	repo := newErrorPartyRepo(fmt.Errorf("save failed"), nil)
	seedPartyWithProfiles(t, repo.inner, "party-status-save", false)
	handler := NewChangePartyStatusHandler(repo)
	if _, err := handler.Handle(context.Background(), &ChangePartyStatusCommand{ID: "party-status-save", Status: "INACTIVE", ActorID: "user"}); err == nil {
		t.Fatalf("expected save error")
	}
}

func TestAddAndRemovePartyRoleHandlers(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-004", false)

	addHandler := NewAddPartyRoleHandler(repo)
	party, err := addHandler.Handle(context.Background(), &AddPartyRoleCommand{
		ID:      "party-004",
		Role:    "CLIENT",
		ActorID: "user-1",
	})
	if err != nil {
		t.Fatalf("Add role should not error: %v", err)
	}
	if len(party.Roles()) != 1 {
		t.Fatalf("Expected 1 role, got %d", len(party.Roles()))
	}

	removeHandler := NewRemovePartyRoleHandler(repo)
	party, err = removeHandler.Handle(context.Background(), &RemovePartyRoleCommand{
		ID:      "party-004",
		Role:    "CLIENT",
		ActorID: "user-1",
	})
	if err != nil {
		t.Fatalf("Remove role should not error: %v", err)
	}
	if len(party.Roles()) != 0 {
		t.Fatalf("Expected 0 roles after removal, got %d", len(party.Roles()))
	}
}

func TestAddPartyRoleHandler_InvalidRole(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-role", false)

	handler := NewAddPartyRoleHandler(repo)
	if _, err := handler.Handle(context.Background(), &AddPartyRoleCommand{ID: "party-role", Role: "BAD", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for invalid role")
	}
}

func TestAddPartyRoleHandler_InvalidPartyID(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewAddPartyRoleHandler(repo)

	if _, err := handler.Handle(context.Background(), &AddPartyRoleCommand{ID: "", Role: "CLIENT", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for invalid party ID")
	}
}

func TestAddPartyRoleHandler_PartyNotFound(t *testing.T) {
	repo := newErrorPartyRepo(nil, fmt.Errorf("not found"))
	handler := NewAddPartyRoleHandler(repo)
	if _, err := handler.Handle(context.Background(), &AddPartyRoleCommand{ID: "missing", Role: "CLIENT", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected party not found error")
	}
}

func TestAddPartyRoleHandler_SaveError(t *testing.T) {
	repo := newErrorPartyRepo(fmt.Errorf("save failed"), nil)
	seedPartyWithProfiles(t, repo.inner, "party-role-save", false)
	handler := NewAddPartyRoleHandler(repo)
	if _, err := handler.Handle(context.Background(), &AddPartyRoleCommand{ID: "party-role-save", Role: "CLIENT", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected save error")
	}
}

func TestRemovePartyRoleHandler_InvalidRole(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-role-remove", false)

	handler := NewRemovePartyRoleHandler(repo)
	if _, err := handler.Handle(context.Background(), &RemovePartyRoleCommand{ID: "party-role-remove", Role: "BAD", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for invalid role")
	}
}

func TestRemovePartyRoleHandler_InvalidPartyID(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewRemovePartyRoleHandler(repo)

	if _, err := handler.Handle(context.Background(), &RemovePartyRoleCommand{ID: "", Role: "CLIENT", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for invalid party ID")
	}
}

func TestRemovePartyRoleHandler_SaveError(t *testing.T) {
	repo := newErrorPartyRepo(fmt.Errorf("save failed"), nil)
	seedPartyWithProfiles(t, repo.inner, "party-role-remove-save", false)

	removeHandler := NewRemovePartyRoleHandler(repo)
	if _, err := removeHandler.Handle(context.Background(), &RemovePartyRoleCommand{ID: "party-role-remove-save", Role: "CLIENT", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected save error")
	}
}

func TestAddPartyRelationshipHandler_Success(t *testing.T) {
	relRepo := newFakeRelationshipRepo()
	handler := NewAddPartyRelationshipHandler(relRepo)

	rel, err := handler.Handle(context.Background(), &AddPartyRelationshipCommand{
		ID:             "party-005",
		RelationshipID: "rel-001",
		ToPartyID:      "party-006",
		Type:           "IS_EMPLOYEE_OF",
		ActorID:        "user-1",
	})
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if rel.Type() != domain.RelationshipIsEmployeeOf {
		t.Fatalf("Expected relationship type IS_EMPLOYEE_OF")
	}
}

func TestAddPartyRelationshipHandler_InvalidType(t *testing.T) {
	relRepo := newFakeRelationshipRepo()
	handler := NewAddPartyRelationshipHandler(relRepo)

	if _, err := handler.Handle(context.Background(), &AddPartyRelationshipCommand{
		ID:             "party-011",
		RelationshipID: "rel-bad",
		ToPartyID:      "party-012",
		Type:           "BAD",
		ActorID:        "user-1",
	}); err == nil {
		t.Fatalf("expected error for invalid relationship type")
	}
}

func TestAddPartyRelationshipHandler_InvalidRelationshipID(t *testing.T) {
	relRepo := newFakeRelationshipRepo()
	handler := NewAddPartyRelationshipHandler(relRepo)

	if _, err := handler.Handle(context.Background(), &AddPartyRelationshipCommand{
		ID:             "party-011",
		RelationshipID: "",
		ToPartyID:      "party-012",
		Type:           "IS_EMPLOYEE_OF",
		ActorID:        "user-1",
	}); err == nil {
		t.Fatalf("expected error for invalid relationship ID")
	}
}

func TestAddPartyRelationshipHandler_SaveError(t *testing.T) {
	relRepo := newErrorRelationshipRepo(fmt.Errorf("save failed"), nil, nil)
	handler := NewAddPartyRelationshipHandler(relRepo)

	if _, err := handler.Handle(context.Background(), &AddPartyRelationshipCommand{
		ID:             "party-020",
		RelationshipID: "rel-020",
		ToPartyID:      "party-021",
		Type:           "IS_EMPLOYEE_OF",
		ActorID:        "user-1",
	}); err == nil {
		t.Fatalf("expected save error")
	}
}

func TestRemovePartyRelationshipHandler_Success(t *testing.T) {
	relRepo := newFakeRelationshipRepo()
	relID, _ := domain.NewPartyRelationshipID("rel-010")
	fromID, _ := domain.NewPartyID("party-010")
	toID, _ := domain.NewPartyID("party-011")
	relationship, _ := domain.NewPartyRelationship(relID, fromID, toID, domain.RelationshipIsEmployeeOf)
	_ = relRepo.Save(context.Background(), relationship, "user-1", "user-1")

	handler := NewRemovePartyRelationshipHandler(relRepo)
	if err := handler.Handle(context.Background(), &RemovePartyRelationshipCommand{RelationshipID: "rel-010"}); err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if _, ok := relRepo.relationships["rel-010"]; ok {
		t.Fatalf("Expected relationship to be removed")
	}
}

func TestRemovePartyRelationshipHandler_InvalidID(t *testing.T) {
	relRepo := newFakeRelationshipRepo()
	handler := NewRemovePartyRelationshipHandler(relRepo)
	if err := handler.Handle(context.Background(), &RemovePartyRelationshipCommand{RelationshipID: ""}); err == nil {
		t.Fatalf("expected error for empty relationship ID")
	}
}

func TestRemovePartyRelationshipHandler_DeleteError(t *testing.T) {
	relRepo := newErrorRelationshipRepo(nil, fmt.Errorf("delete failed"), nil)
	handler := NewRemovePartyRelationshipHandler(relRepo)
	if err := handler.Handle(context.Background(), &RemovePartyRelationshipCommand{RelationshipID: "rel-030"}); err == nil {
		t.Fatalf("expected delete error")
	}
}

func TestAddContactDetailsHandler_Success(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-006", true)

	handler := NewAddContactDetailsHandler(repo)
	contact, err := handler.Handle(context.Background(), &AddContactDetailsCommand{
		PartyID:         "party-006",
		ContactID:       "contact-100",
		TypeDescription: "Ventas",
		Email:           "ventas@org.local",
		ActorID:         "user-2",
	})
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if contact.TypeDescription() != "Ventas" {
		t.Fatalf("Expected contact type Ventas")
	}
}

func TestAddContactDetailsHandler_InvalidPartyID(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewAddContactDetailsHandler(repo)

	if _, err := handler.Handle(context.Background(), &AddContactDetailsCommand{PartyID: "", ContactID: "c-1", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for empty party ID")
	}
}

func TestAddContactDetailsHandler_InvalidContactID(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-contact-id", true)

	handler := NewAddContactDetailsHandler(repo)
	if _, err := handler.Handle(context.Background(), &AddContactDetailsCommand{PartyID: "party-contact-id", ContactID: "", TypeDescription: "Sales", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for empty contact ID")
	}
}

func TestAddContactDetailsHandler_NoOrganizationProfile(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-no-org", false)

	handler := NewAddContactDetailsHandler(repo)
	if _, err := handler.Handle(context.Background(), &AddContactDetailsCommand{PartyID: "party-no-org", ContactID: "c-1", TypeDescription: "Ventas", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error when organization profile missing")
	}
}

func TestAddContactDetailsHandler_InvalidPhone(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-bad-phone", true)

	handler := NewAddContactDetailsHandler(repo)
	if _, err := handler.Handle(context.Background(), &AddContactDetailsCommand{PartyID: "party-bad-phone", ContactID: "c-1", TypeDescription: "Ventas", Phone: "bad", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for invalid phone")
	}
}

func TestAddContactDetailsHandler_SaveError(t *testing.T) {
	repo := newErrorPartyRepo(fmt.Errorf("save failed"), nil)
	seedPartyWithProfiles(t, repo.inner, "party-save-contact", true)

	handler := NewAddContactDetailsHandler(repo)
	if _, err := handler.Handle(context.Background(), &AddContactDetailsCommand{PartyID: "party-save-contact", ContactID: "c-1", TypeDescription: "Ventas", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected save error")
	}
}

func TestUpdateAndRemoveContactDetailsHandlers(t *testing.T) {
	repo := newFakePartyRepo()
	party := seedPartyWithProfiles(t, repo, "party-006", true)

	contactID, _ := domain.NewContactDetailsID("contact-200")
	contact, _ := domain.NewContactDetails(contactID, "Ventas", nil, nil, nil)
	_ = party.OrganizationProfile().AddContact(contact)
	_ = repo.Save(context.Background(), party, "user-1", "user-1")

	updateHandler := NewUpdateContactDetailsHandler(repo)
	updatedType := "Soporte"
	updatedEmail := "soporte@org.local"
	updated, err := updateHandler.Handle(context.Background(), &UpdateContactDetailsCommand{
		PartyID:         "party-006",
		ContactID:       "contact-200",
		TypeDescription: &updatedType,
		Email:           &updatedEmail,
		ActorID:         "user-2",
	})
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if updated.TypeDescription() != "Soporte" {
		t.Fatalf("Expected contact type Soporte")
	}
	if updated.Email() == nil || updated.Email().Value() != "soporte@org.local" {
		t.Fatalf("Expected updated email")
	}

	removeHandler := NewRemoveContactDetailsHandler(repo)
	if err := removeHandler.Handle(context.Background(), &RemoveContactDetailsCommand{
		PartyID:   "party-006",
		ContactID: "contact-200",
		ActorID:   "user-2",
	}); err != nil {
		t.Fatalf("Remove should not error: %v", err)
	}

	contacts := party.OrganizationProfile().Contacts()
	if len(contacts) != 0 {
		t.Fatalf("Expected contact to be removed")
	}
}

func TestUpdateContactDetailsHandler_InvalidInputs(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewUpdateContactDetailsHandler(repo)

	if _, err := handler.Handle(context.Background(), &UpdateContactDetailsCommand{PartyID: "", ContactID: "c-1", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for empty party ID")
	}
	if _, err := handler.Handle(context.Background(), &UpdateContactDetailsCommand{PartyID: "p-1", ContactID: "", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for empty contact ID")
	}
	if _, err := handler.Handle(context.Background(), &UpdateContactDetailsCommand{PartyID: "p-1", ContactID: "c-1", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error when no fields provided")
	}
}

func TestUpdateContactDetailsHandler_ContactNotFound(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-no-contact", true)

	handler := NewUpdateContactDetailsHandler(repo)
	updateType := "Soporte"
	if _, err := handler.Handle(context.Background(), &UpdateContactDetailsCommand{PartyID: "party-no-contact", ContactID: "c-missing", TypeDescription: &updateType, ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error when contact not found")
	}
}

func TestUpdateContactDetailsHandler_InvalidEmail(t *testing.T) {
	repo := newFakePartyRepo()
	party := seedPartyWithProfiles(t, repo, "party-bad-email", true)
	contactID, _ := domain.NewContactDetailsID("contact-301")
	contact, _ := domain.NewContactDetails(contactID, "Ventas", nil, nil, nil)
	_ = party.OrganizationProfile().AddContact(contact)
	_ = repo.Save(context.Background(), party, "user-1", "user-1")

	handler := NewUpdateContactDetailsHandler(repo)
	badEmail := "bad"
	if _, err := handler.Handle(context.Background(), &UpdateContactDetailsCommand{PartyID: "party-bad-email", ContactID: "contact-301", Email: &badEmail, ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for invalid email")
	}
}

func TestUpdateContactDetailsHandler_InvalidPhone(t *testing.T) {
	repo := newFakePartyRepo()
	party := seedPartyWithProfiles(t, repo, "party-bad-phone-update", true)
	contactID, _ := domain.NewContactDetailsID("contact-302")
	contact, _ := domain.NewContactDetails(contactID, "Ventas", nil, nil, nil)
	_ = party.OrganizationProfile().AddContact(contact)
	_ = repo.Save(context.Background(), party, "user-1", "user-1")

	handler := NewUpdateContactDetailsHandler(repo)
	badPhone := "bad"
	if _, err := handler.Handle(context.Background(), &UpdateContactDetailsCommand{PartyID: "party-bad-phone-update", ContactID: "contact-302", Phone: &badPhone, ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for invalid phone")
	}
}

func TestUpdateContactDetailsHandler_InvalidTypeDescription(t *testing.T) {
	repo := newFakePartyRepo()
	party := seedPartyWithProfiles(t, repo, "party-bad-type", true)
	contactID, _ := domain.NewContactDetailsID("contact-303")
	contact, _ := domain.NewContactDetails(contactID, "Ventas", nil, nil, nil)
	_ = party.OrganizationProfile().AddContact(contact)
	_ = repo.Save(context.Background(), party, "user-1", "user-1")

	handler := NewUpdateContactDetailsHandler(repo)
	badType := ""
	if _, err := handler.Handle(context.Background(), &UpdateContactDetailsCommand{PartyID: "party-bad-type", ContactID: "contact-303", TypeDescription: &badType, ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for invalid type description")
	}
}

func TestUpdateContactDetailsHandler_ClearFields(t *testing.T) {
	repo := newFakePartyRepo()
	party := seedPartyWithProfiles(t, repo, "party-clear-fields", true)
	contactID, _ := domain.NewContactDetailsID("contact-401")
	phone, _ := domain.NewPhone("+34 600 111 222")
	email, _ := domain.NewEmail("ventas@org.local")
	contact, _ := domain.NewContactDetails(contactID, "Ventas", phone, email, nil)
	_ = party.OrganizationProfile().AddContact(contact)
	_ = repo.Save(context.Background(), party, "user-1", "user-1")

	handler := NewUpdateContactDetailsHandler(repo)
	empty := ""
	updated, err := handler.Handle(context.Background(), &UpdateContactDetailsCommand{PartyID: "party-clear-fields", ContactID: "contact-401", Phone: &empty, Email: &empty, ActorID: "user-1"})
	if err != nil {
		t.Fatalf("expected update without error: %v", err)
	}
	if updated.Phone() != nil || updated.Email() != nil {
		t.Fatalf("expected phone and email to be cleared")
	}
}

func TestUpdateContactDetailsHandler_UpdateRelatedParty(t *testing.T) {
	repo := newFakePartyRepo()
	party := seedPartyWithProfiles(t, repo, "party-related-update", true)
	contactID, _ := domain.NewContactDetailsID("contact-related")
	oldRelated, _ := domain.NewPartyID("party-old")
	contact, _ := domain.NewContactDetails(contactID, "Ventas", nil, nil, &oldRelated)
	_ = party.OrganizationProfile().AddContact(contact)
	_ = repo.Save(context.Background(), party, "user-1", "user-1")

	newRelated := "party-new"
	handler := NewUpdateContactDetailsHandler(repo)
	updated, err := handler.Handle(context.Background(), &UpdateContactDetailsCommand{
		PartyID:        "party-related-update",
		ContactID:      "contact-related",
		RelatedPartyID: &newRelated,
		ActorID:        "user-1",
	})
	if err != nil {
		t.Fatalf("expected update without error: %v", err)
	}
	if updated.RelatedPartyID() == nil || updated.RelatedPartyID().String() != "party-new" {
		t.Fatalf("expected related party ID to be updated")
	}
}

func TestUpdateContactDetailsHandler_ClearRelatedParty(t *testing.T) {
	repo := newFakePartyRepo()
	party := seedPartyWithProfiles(t, repo, "party-related-clear", true)
	contactID, _ := domain.NewContactDetailsID("contact-related-clear")
	oldRelated, _ := domain.NewPartyID("party-old")
	contact, _ := domain.NewContactDetails(contactID, "Ventas", nil, nil, &oldRelated)
	_ = party.OrganizationProfile().AddContact(contact)
	_ = repo.Save(context.Background(), party, "user-1", "user-1")

	empty := ""
	handler := NewUpdateContactDetailsHandler(repo)
	updated, err := handler.Handle(context.Background(), &UpdateContactDetailsCommand{
		PartyID:        "party-related-clear",
		ContactID:      "contact-related-clear",
		RelatedPartyID: &empty,
		ActorID:        "user-1",
	})
	if err != nil {
		t.Fatalf("expected update without error: %v", err)
	}
	if updated.RelatedPartyID() != nil {
		t.Fatalf("expected related party ID to be cleared")
	}
}

func TestRemoveContactDetailsHandler_InvalidInputs(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewRemoveContactDetailsHandler(repo)

	if err := handler.Handle(context.Background(), &RemoveContactDetailsCommand{PartyID: "", ContactID: "c-1", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for empty party ID")
	}
	if err := handler.Handle(context.Background(), &RemoveContactDetailsCommand{PartyID: "p-1", ContactID: "", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for empty contact ID")
	}
}

func TestRemoveContactDetailsHandler_NoOrganizationProfile(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-no-org-remove", false)

	handler := NewRemoveContactDetailsHandler(repo)
	if err := handler.Handle(context.Background(), &RemoveContactDetailsCommand{PartyID: "party-no-org-remove", ContactID: "c-1", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error when organization profile missing")
	}
}

func TestRemoveContactDetailsHandler_SaveError(t *testing.T) {
	repo := newErrorPartyRepo(fmt.Errorf("save failed"), nil)
	party := seedPartyWithProfiles(t, repo.inner, "party-remove-save", true)
	contactID, _ := domain.NewContactDetailsID("contact-901")
	contact, _ := domain.NewContactDetails(contactID, "Ventas", nil, nil, nil)
	_ = party.OrganizationProfile().AddContact(contact)
	_ = repo.inner.Save(context.Background(), party, "user-1", "user-1")

	handler := NewRemoveContactDetailsHandler(repo)
	if err := handler.Handle(context.Background(), &RemoveContactDetailsCommand{PartyID: "party-remove-save", ContactID: "contact-901", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected save error")
	}
}

func TestAddPartyAddressHandler_Success(t *testing.T) {
	addressRepo := newFakePartyAddressRepo()
	handler := NewAddPartyAddressHandler(addressRepo)

	address, err := handler.Handle(context.Background(), &AddPartyAddressCommand{
		PartyID:    "party-007",
		AddressID:  "addr-001",
		Street:     "Calle 1",
		City:       "Madrid",
		Province:   "Madrid",
		PostalCode: "28001",
		Country:    "Spain",
		ActorID:    "user-1",
	})
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if address.City() != "Madrid" {
		t.Fatalf("Expected city Madrid, got %s", address.City())
	}
}

func TestAddPartyAddressHandler_InvalidInputs(t *testing.T) {
	addressRepo := newFakePartyAddressRepo()
	handler := NewAddPartyAddressHandler(addressRepo)

	if _, err := handler.Handle(context.Background(), &AddPartyAddressCommand{PartyID: "", AddressID: "addr-1", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for empty party ID")
	}
	if _, err := handler.Handle(context.Background(), &AddPartyAddressCommand{PartyID: "party-1", AddressID: "", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for empty address ID")
	}
}

func TestAddPartyAddressHandler_SaveError(t *testing.T) {
	addressRepo := &errorPartyAddressRepo{saveErr: fmt.Errorf("save failed")}
	handler := NewAddPartyAddressHandler(addressRepo)

	cmd := &AddPartyAddressCommand{
		PartyID:    "party-addr",
		AddressID:  "addr-err",
		Street:     "Calle 1",
		City:       "Madrid",
		PostalCode: "28001",
		Country:    "Spain",
		ActorID:    "user-1",
	}
	if _, err := handler.Handle(context.Background(), cmd); err == nil {
		t.Fatalf("expected save error")
	}
}

func TestAddPartyAddressHandler_InvalidAddress(t *testing.T) {
	addressRepo := newFakePartyAddressRepo()
	handler := NewAddPartyAddressHandler(addressRepo)

	cmd := &AddPartyAddressCommand{
		PartyID:    "party-addr-bad",
		AddressID:  "addr-bad",
		Street:     "Calle 1",
		City:       "",
		PostalCode: "28001",
		Country:    "Spain",
		ActorID:    "user-1",
	}
	if _, err := handler.Handle(context.Background(), cmd); err == nil {
		t.Fatalf("expected error for invalid address")
	}
}

func TestListHandlersForParty(t *testing.T) {
	partyRepo := newFakePartyRepo()
	relRepo := newFakeRelationshipRepo()
	addressRepo := newFakePartyAddressRepo()

	party := seedPartyWithProfiles(t, partyRepo, "party-008", true)
	contactID, _ := domain.NewContactDetailsID("contact-200")
	contact, _ := domain.NewContactDetails(contactID, "Soporte", nil, nil, nil)
	_ = party.OrganizationProfile().AddContact(contact)
	_ = partyRepo.Save(context.Background(), party, "user-1", "user-1")

	relID, _ := domain.NewPartyRelationshipID("rel-200")
	fromID, _ := domain.NewPartyID("party-008")
	toID, _ := domain.NewPartyID("party-009")
	relationship, _ := domain.NewPartyRelationship(relID, fromID, toID, domain.RelationshipIsEmployeeOf)
	_ = relRepo.Save(context.Background(), relationship, "user-1", "user-1")

	addressID, _ := domain.NewAddressID("addr-200")
	address, _ := domain.NewAddress("Calle 2", "Madrid", "Madrid", "28002", "Spain")
	_ = addressRepo.Save(context.Background(), address, addressID, fromID, "user-1", "user-1")

	listPartiesHandler := NewListPartiesHandler(partyRepo)
	parties, err := listPartiesHandler.Handle(context.Background(), &ListPartiesQuery{PageNumber: 1, PageSize: 10})
	if err != nil || len(parties) != 1 {
		t.Fatalf("Expected 1 party, got %d (err=%v)", len(parties), err)
	}

	listRelsHandler := NewListPartyRelationshipsHandler(relRepo)
	rels, err := listRelsHandler.Handle(context.Background(), &ListPartyRelationshipsQuery{PartyID: "party-008"})
	if err != nil || len(rels) != 1 {
		t.Fatalf("Expected 1 relationship, got %d (err=%v)", len(rels), err)
	}

	listContactsHandler := NewListContactDetailsHandler(partyRepo)
	contacts, err := listContactsHandler.Handle(context.Background(), &ListContactDetailsQuery{PartyID: "party-008"})
	if err != nil || len(contacts) != 1 {
		t.Fatalf("Expected 1 contact, got %d (err=%v)", len(contacts), err)
	}

	listAddressesHandler := NewListPartyAddressesHandler(addressRepo)
	addresses, err := listAddressesHandler.Handle(context.Background(), &ListPartyAddressesQuery{PartyID: "party-008"})
	if err != nil || len(addresses) != 1 {
		t.Fatalf("Expected 1 address, got %d (err=%v)", len(addresses), err)
	}
}

func TestParsePartyRoleType(t *testing.T) {
	role, err := parsePartyRoleType("client")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if role != domain.PartyRoleClient {
		t.Fatalf("expected CLIENT, got %s", role)
	}

	if _, err := parsePartyRoleType("bad"); err == nil {
		t.Fatalf("expected error for invalid role")
	}
}

func TestCreatePartyHandler_DefaultStatusAndTaxIDType(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewCreatePartyHandler(repo)

	cmd := &CreatePartyCommand{
		ID:      "party-defaults",
		ActorID: "user-1",
		OrganizationProfile: &OrganizationProfileInput{
			Name:  strPtr("Org"),
			TaxID: strPtr("B12345678"),
		},
	}

	party, err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if party.Status() != domain.PartyStatusActive {
		t.Fatalf("expected status ACTIVE")
	}
	if party.OrganizationProfile().TaxID() == nil || party.OrganizationProfile().TaxID().Type() != "NIF" {
		t.Fatalf("expected default tax ID type NIF")
	}
}

func TestCreatePartyHandler_DuplicateRole(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewCreatePartyHandler(repo)

	cmd := &CreatePartyCommand{
		ID:      "party-dup-role",
		ActorID: "user-1",
		Roles:   []string{"CLIENT", "CLIENT"},
		PersonProfile: &PersonProfileInput{
			FirstName: strPtr("Ana"),
			LastName:  strPtr("Perez"),
		},
	}

	_, err := handler.Handle(context.Background(), cmd)
	if err == nil {
		t.Fatalf("expected error for duplicate role")
	}
	var partyErr domain.PartyError
	if !errors.As(err, &partyErr) || partyErr.Code != domain.ErrCodeConflict {
		t.Fatalf("expected conflict error, got %v", err)
	}
}

func TestUpdatePartyHandler_UsesExistingPersonProfile(t *testing.T) {
	repo := newFakePartyRepo()
	party := seedPartyWithProfiles(t, repo, "party-existing-person", false)

	handler := NewUpdatePartyHandler(repo)
	cmd := &UpdatePartyCommand{
		ID:            party.ID().String(),
		ActorID:       "user-1",
		PersonProfile: &PersonProfileInput{},
	}

	updated, err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if updated.PersonProfile() == nil || updated.PersonProfile().FirstName() == "" || updated.PersonProfile().LastName() == "" {
		t.Fatalf("expected existing person profile to be preserved")
	}
}

func TestUpdatePartyHandler_InvalidTaxID(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-bad-tax", true)

	handler := NewUpdatePartyHandler(repo)
	cmd := &UpdatePartyCommand{
		ID:      "party-bad-tax",
		ActorID: "user-1",
		OrganizationProfile: &OrganizationProfileInput{
			Name:  strPtr("Org"),
			TaxID: strPtr("BAD!"),
		},
	}

	if _, err := handler.Handle(context.Background(), cmd); err == nil {
		t.Fatalf("expected error for invalid tax ID")
	}
}

func TestUpdatePartyHandler_NotFound(t *testing.T) {
	repo := newErrorPartyRepo(nil, fmt.Errorf("missing"))
	handler := NewUpdatePartyHandler(repo)

	_, err := handler.Handle(context.Background(), &UpdatePartyCommand{ID: "missing", ActorID: "user-1"})
	if err == nil {
		t.Fatalf("expected not found error")
	}
	var partyErr domain.PartyError
	if !errors.As(err, &partyErr) || partyErr.Code != domain.ErrCodeNotFound {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestChangePartyStatusHandler_PartyNotFound(t *testing.T) {
	repo := newErrorPartyRepo(nil, fmt.Errorf("missing"))
	handler := NewChangePartyStatusHandler(repo)

	_, err := handler.Handle(context.Background(), &ChangePartyStatusCommand{ID: "missing", Status: "ACTIVE", ActorID: "user-1"})
	if err == nil {
		t.Fatalf("expected not found error")
	}
}

func TestAddPartyRoleHandler_EmptyActorID(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-role-actor", false)

	handler := NewAddPartyRoleHandler(repo)
	if _, err := handler.Handle(context.Background(), &AddPartyRoleCommand{ID: "party-role-actor", Role: "CLIENT", ActorID: ""}); err == nil {
		t.Fatalf("expected error for empty actor ID")
	}
}

func TestRemovePartyRoleHandler_EmptyActorID(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-role-actor-remove", false)

	handler := NewRemovePartyRoleHandler(repo)
	if _, err := handler.Handle(context.Background(), &RemovePartyRoleCommand{ID: "party-role-actor-remove", Role: "CLIENT", ActorID: ""}); err == nil {
		t.Fatalf("expected error for empty actor ID")
	}
}

func TestAddPartyRelationshipHandler_EmptyActorID(t *testing.T) {
	relRepo := newFakeRelationshipRepo()
	handler := NewAddPartyRelationshipHandler(relRepo)

	if _, err := handler.Handle(context.Background(), &AddPartyRelationshipCommand{ID: "party-1", RelationshipID: "rel-1", ToPartyID: "party-2", Type: "IS_EMPLOYEE_OF", ActorID: ""}); err == nil {
		t.Fatalf("expected error for empty actor ID")
	}
}

func TestAddPartyRelationshipHandler_InvalidToPartyID(t *testing.T) {
	relRepo := newFakeRelationshipRepo()
	handler := NewAddPartyRelationshipHandler(relRepo)

	if _, err := handler.Handle(context.Background(), &AddPartyRelationshipCommand{ID: "party-1", RelationshipID: "rel-1", ToPartyID: "", Type: "IS_EMPLOYEE_OF", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for invalid to party ID")
	}
}

func TestAddContactDetailsHandler_InvalidEmail(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-bad-email-add", true)

	handler := NewAddContactDetailsHandler(repo)
	if _, err := handler.Handle(context.Background(), &AddContactDetailsCommand{PartyID: "party-bad-email-add", ContactID: "c-1", TypeDescription: "Sales", Email: "bad", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for invalid email")
	}
}

func TestAddContactDetailsHandler_DuplicateContact(t *testing.T) {
	repo := newFakePartyRepo()
	party := seedPartyWithProfiles(t, repo, "party-dup-contact", true)
	contactID, _ := domain.NewContactDetailsID("contact-dup")
	contact, _ := domain.NewContactDetails(contactID, "Sales", nil, nil, nil)
	_ = party.OrganizationProfile().AddContact(contact)
	_ = repo.Save(context.Background(), party, "user-1", "user-1")

	handler := NewAddContactDetailsHandler(repo)
	_, err := handler.Handle(context.Background(), &AddContactDetailsCommand{PartyID: "party-dup-contact", ContactID: "contact-dup", TypeDescription: "Sales", ActorID: "user-1"})
	if err == nil {
		t.Fatalf("expected error for duplicate contact")
	}
}

func TestUpdateContactDetailsHandler_EmptyActorID(t *testing.T) {
	repo := newFakePartyRepo()
	party := seedPartyWithProfiles(t, repo, "party-update-actor", true)
	contactID, _ := domain.NewContactDetailsID("contact-actor")
	contact, _ := domain.NewContactDetails(contactID, "Sales", nil, nil, nil)
	_ = party.OrganizationProfile().AddContact(contact)
	_ = repo.Save(context.Background(), party, "user-1", "user-1")

	handler := NewUpdateContactDetailsHandler(repo)
	updatedType := "Support"
	if _, err := handler.Handle(context.Background(), &UpdateContactDetailsCommand{PartyID: "party-update-actor", ContactID: "contact-actor", TypeDescription: &updatedType, ActorID: ""}); err == nil {
		t.Fatalf("expected error for empty actor ID")
	}
}

func TestUpdateContactDetailsHandler_PartyNotFound(t *testing.T) {
	repo := newErrorPartyRepo(nil, fmt.Errorf("missing"))
	handler := NewUpdateContactDetailsHandler(repo)
	updatedType := "Support"

	if _, err := handler.Handle(context.Background(), &UpdateContactDetailsCommand{PartyID: "missing", ContactID: "contact-1", TypeDescription: &updatedType, ActorID: "user-1"}); err == nil {
		t.Fatalf("expected not found error")
	}
}

func TestRemoveContactDetailsHandler_EmptyActorID(t *testing.T) {
	repo := newFakePartyRepo()
	party := seedPartyWithProfiles(t, repo, "party-remove-actor", true)
	contactID, _ := domain.NewContactDetailsID("contact-remove")
	contact, _ := domain.NewContactDetails(contactID, "Sales", nil, nil, nil)
	_ = party.OrganizationProfile().AddContact(contact)
	_ = repo.Save(context.Background(), party, "user-1", "user-1")

	handler := NewRemoveContactDetailsHandler(repo)
	if err := handler.Handle(context.Background(), &RemoveContactDetailsCommand{PartyID: "party-remove-actor", ContactID: "contact-remove", ActorID: ""}); err == nil {
		t.Fatalf("expected error for empty actor ID")
	}
}

func TestRemoveContactDetailsHandler_PartyNotFound(t *testing.T) {
	repo := newErrorPartyRepo(nil, fmt.Errorf("missing"))
	handler := NewRemoveContactDetailsHandler(repo)

	if err := handler.Handle(context.Background(), &RemoveContactDetailsCommand{PartyID: "missing", ContactID: "contact-1", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected not found error")
	}
}

func TestRemoveContactDetailsHandler_ContactNotFound(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-contact-missing", true)

	handler := NewRemoveContactDetailsHandler(repo)
	if err := handler.Handle(context.Background(), &RemoveContactDetailsCommand{PartyID: "party-contact-missing", ContactID: "contact-missing", ActorID: "user-1"}); err == nil {
		t.Fatalf("expected error for missing contact")
	}
}

func TestAddPartyAddressHandler_EmptyActorID(t *testing.T) {
	addressRepo := newFakePartyAddressRepo()
	handler := NewAddPartyAddressHandler(addressRepo)

	if _, err := handler.Handle(context.Background(), &AddPartyAddressCommand{PartyID: "party-addr", AddressID: "addr-1", Street: "Calle 1", City: "Madrid", PostalCode: "28001", Country: "Spain", ActorID: ""}); err == nil {
		t.Fatalf("expected error for empty actor ID")
	}
}
