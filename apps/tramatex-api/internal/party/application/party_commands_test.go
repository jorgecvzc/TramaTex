package application

import (
	"context"
	"fmt"
	"testing"

	"github.com/joran-cortez/tramatex/internal/party/domain"
	"github.com/joran-cortez/tramatex/internal/party/persistence"
)

type fakePartyRepo struct {
	parties map[string]*domain.Party
}

func newFakePartyRepo() *fakePartyRepo {
	return &fakePartyRepo{parties: make(map[string]*domain.Party)}
}

func (r *fakePartyRepo) Save(ctx context.Context, party *domain.Party) error {
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

type fakeRelationshipRepo struct {
	relationships map[string]domain.PartyRelationship
}

func newFakeRelationshipRepo() *fakeRelationshipRepo {
	return &fakeRelationshipRepo{relationships: make(map[string]domain.PartyRelationship)}
}

func (r *fakeRelationshipRepo) Save(ctx context.Context, relationship domain.PartyRelationship) error {
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

type fakePartyAddressRepo struct {
	addressesByParty map[string][]*domain.Address
	addressesByID    map[string]*domain.Address
}

func newFakePartyAddressRepo() *fakePartyAddressRepo {
	return &fakePartyAddressRepo{
		addressesByParty: make(map[string][]*domain.Address),
		addressesByID:    make(map[string]*domain.Address),
	}
}

func (r *fakePartyAddressRepo) Save(ctx context.Context, address *domain.Address, addressID domain.AddressID, partyID domain.PartyID, createdBy string, modifiedBy string) error {
	if address == nil {
		return fmt.Errorf("address cannot be nil")
	}
	r.addressesByID[addressID.String()] = address
	r.addressesByParty[partyID.String()] = append(r.addressesByParty[partyID.String()], address)
	return nil
}

func (r *fakePartyAddressRepo) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]*domain.Address, error) {
	return r.addressesByParty[partyID.String()], nil
}

func (r *fakePartyAddressRepo) FindPrimary(ctx context.Context, partyID domain.PartyID) (*domain.Address, error) {
	addresses := r.addressesByParty[partyID.String()]
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no addresses")
	}
	return addresses[0], nil
}

func (r *fakePartyAddressRepo) Delete(ctx context.Context, id domain.AddressID) error {
	delete(r.addressesByID, id.String())
	return nil
}

func seedPartyWithProfiles(t *testing.T, repo *fakePartyRepo, id string, withOrg bool) *domain.Party {
	t.Helper()
	partyID, _ := domain.NewPartyID(id)
	personProfile, _ := domain.NewPersonProfile("Ana", "Perez")
	var orgProfile *domain.OrganizationProfile
	if withOrg {
		orgProfile, _ = domain.NewOrganizationProfile("Org", nil, "")
	}
	party, err := domain.NewParty(partyID, domain.PartyStatusActive, "user-1", personProfile, orgProfile)
	if err != nil {
		t.Fatalf("Failed to create party: %v", err)
	}
	if err := repo.Save(context.Background(), party); err != nil {
		t.Fatalf("Failed to seed party: %v", err)
	}
	return party
}

func TestCreatePartyHandler_Success(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewCreatePartyHandler(repo)
	ctx := context.Background()

	cmd := &CreatePartyCommand{
		ID:        "party-001",
		Status:    "ACTIVE",
		Roles:     []string{"CLIENT"},
		CreatedBy: "user-1",
		PersonProfile: &PersonProfileInput{
			FirstName: "Ana",
			LastName:  "Perez",
		},
		OrganizationProfile: &OrganizationProfileInput{
			Name:      "Textiles Perez",
			TaxID:     "B12345678",
			TaxIDType: "CIF",
			Website:   "https://textiles.local",
			Contacts: []ContactDetailsInput{
				{
					ID:              "contact-001",
					TypeDescription: "Ventas",
					Phone:           "+34 600 111 222",
					Email:           "ventas@textiles.local",
				},
			},
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
	if len(party.OrganizationProfile().Contacts()) != 1 {
		t.Fatalf("Expected 1 contact, got %d", len(party.OrganizationProfile().Contacts()))
	}
}

func TestUpdatePartyHandler_UpdatesProfilesAndStatus(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-002", false)

	handler := NewUpdatePartyHandler(repo)
	cmd := &UpdatePartyCommand{
		ID:         "party-002",
		Status:     "INACTIVE",
		ModifiedBy: "user-2",
		OrganizationProfile: &OrganizationProfileInput{
			Name: "Org Updated",
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

func TestChangePartyStatusHandler_Success(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-003", false)

	handler := NewChangePartyStatusHandler(repo)
	cmd := &ChangePartyStatusCommand{
		ID:         "party-003",
		Status:     "INACTIVE",
		ModifiedBy: "user-2",
	}

	party, err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if party.Status() != domain.PartyStatusInactive {
		t.Fatalf("Expected status INACTIVE, got %s", party.Status())
	}
}

func TestAddAndRemovePartyRoleHandlers(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-004", false)

	addHandler := NewAddPartyRoleHandler(repo)
	party, err := addHandler.Handle(context.Background(), &AddPartyRoleCommand{
		ID:   "party-004",
		Role: "CLIENT",
	})
	if err != nil {
		t.Fatalf("Add role should not error: %v", err)
	}
	if len(party.Roles()) != 1 {
		t.Fatalf("Expected 1 role, got %d", len(party.Roles()))
	}

	removeHandler := NewRemovePartyRoleHandler(repo)
	party, err = removeHandler.Handle(context.Background(), &RemovePartyRoleCommand{
		ID:   "party-004",
		Role: "CLIENT",
	})
	if err != nil {
		t.Fatalf("Remove role should not error: %v", err)
	}
	if len(party.Roles()) != 0 {
		t.Fatalf("Expected 0 roles after removal, got %d", len(party.Roles()))
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
	})
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if rel.Type() != domain.RelationshipIsEmployeeOf {
		t.Fatalf("Expected relationship type IS_EMPLOYEE_OF")
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
		ModifiedBy:      "user-2",
	})
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if contact.TypeDescription() != "Ventas" {
		t.Fatalf("Expected contact type Ventas")
	}
}

func TestUpdateAndRemoveContactDetailsHandlers(t *testing.T) {
	repo := newFakePartyRepo()
	party := seedPartyWithProfiles(t, repo, "party-006", true)

	contactID, _ := domain.NewContactDetailsID("contact-200")
	contact, _ := domain.NewContactDetails(contactID, "Ventas", nil, nil, nil)
	_ = party.OrganizationProfile().AddContact(contact)
	_ = repo.Save(context.Background(), party)

	updateHandler := NewUpdateContactDetailsHandler(repo)
	updatedType := "Soporte"
	updatedEmail := "soporte@org.local"
	updated, err := updateHandler.Handle(context.Background(), &UpdateContactDetailsCommand{
		PartyID:         "party-006",
		ContactID:       "contact-200",
		TypeDescription: &updatedType,
		Email:           &updatedEmail,
		ModifiedBy:      "user-2",
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
	}); err != nil {
		t.Fatalf("Remove should not error: %v", err)
	}

	contacts := party.OrganizationProfile().Contacts()
	if len(contacts) != 0 {
		t.Fatalf("Expected contact to be removed")
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
		CreatedBy:  "user-1",
	})
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if address.City() != "Madrid" {
		t.Fatalf("Expected city Madrid, got %s", address.City())
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
	_ = partyRepo.Save(context.Background(), party)

	relID, _ := domain.NewPartyRelationshipID("rel-200")
	fromID, _ := domain.NewPartyID("party-008")
	toID, _ := domain.NewPartyID("party-009")
	relationship, _ := domain.NewPartyRelationship(relID, fromID, toID, domain.RelationshipIsEmployeeOf)
	_ = relRepo.Save(context.Background(), relationship)

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
