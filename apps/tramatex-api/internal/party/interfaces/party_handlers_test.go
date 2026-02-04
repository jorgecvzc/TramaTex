package interfaces

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joran-cortez/tramatex/internal/party/application"
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
		return nil
	}
	r.parties[party.ID().String()] = party
	return nil
}

func (r *fakePartyRepo) FindByID(ctx context.Context, id domain.PartyID) (*domain.Party, error) {
	party, ok := r.parties[id.String()]
	if !ok {
		return nil, errors.New("party not found")
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
	rels map[string]domain.PartyRelationship
}

func newFakeRelationshipRepo() *fakeRelationshipRepo {
	return &fakeRelationshipRepo{rels: make(map[string]domain.PartyRelationship)}
}

func (r *fakeRelationshipRepo) Save(ctx context.Context, relationship domain.PartyRelationship) error {
	r.rels[relationship.ID().String()] = relationship
	return nil
}

func (r *fakeRelationshipRepo) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]domain.PartyRelationship, error) {
	result := make([]domain.PartyRelationship, 0)
	for _, rel := range r.rels {
		if rel.FromID() == partyID || rel.ToID() == partyID {
			result = append(result, rel)
		}
	}
	return result, nil
}

func (r *fakeRelationshipRepo) Delete(ctx context.Context, id domain.PartyRelationshipID) error {
	delete(r.rels, id.String())
	return nil
}

type fakeAddressRepo struct{}

func (r *fakeAddressRepo) Save(ctx context.Context, address *domain.Address, addressID domain.AddressID, partyID domain.PartyID, createdBy string, modifiedBy string) error {
	return nil
}

func (r *fakeAddressRepo) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]*domain.Address, error) {
	return []*domain.Address{}, nil
}

func (r *fakeAddressRepo) FindPrimary(ctx context.Context, partyID domain.PartyID) (*domain.Address, error) {
	return nil, nil
}

func (r *fakeAddressRepo) Delete(ctx context.Context, id domain.AddressID) error {
	return nil
}

func setupHandlers() (*gin.Engine, *fakePartyRepo) {
	gin.SetMode(gin.TestMode)

	partyRepo := newFakePartyRepo()
	relRepo := newFakeRelationshipRepo()
	addressRepo := &fakeAddressRepo{}

	createPartyHandler := application.NewCreatePartyHandler(partyRepo)
	updatePartyHandler := application.NewUpdatePartyHandler(partyRepo)
	changePartyStatusHandler := application.NewChangePartyStatusHandler(partyRepo)
	getPartyHandler := application.NewGetPartyHandler(partyRepo)
	listPartiesHandler := application.NewListPartiesHandler(partyRepo)
	addPartyRoleHandler := application.NewAddPartyRoleHandler(partyRepo)
	removePartyRoleHandler := application.NewRemovePartyRoleHandler(partyRepo)
	addRelationshipHandler := application.NewAddPartyRelationshipHandler(relRepo)
	listRelationshipsHandler := application.NewListPartyRelationshipsHandler(relRepo)
	removeRelationshipHandler := application.NewRemovePartyRelationshipHandler(relRepo)
	addContactHandler := application.NewAddContactDetailsHandler(partyRepo)
	updateContactHandler := application.NewUpdateContactDetailsHandler(partyRepo)
	listContactsHandler := application.NewListContactDetailsHandler(partyRepo)
	removeContactHandler := application.NewRemoveContactDetailsHandler(partyRepo)
	addAddressHandler := application.NewAddPartyAddressHandler(addressRepo)
	listAddressesHandler := application.NewListPartyAddressesHandler(addressRepo)

	partyHandler := NewPartyHandler(createPartyHandler, updatePartyHandler, changePartyStatusHandler, getPartyHandler, listPartiesHandler)
	partyRoleHandler := NewPartyRoleHandler(addPartyRoleHandler, removePartyRoleHandler)
	partyRelationshipHandler := NewPartyRelationshipHandler(addRelationshipHandler, listRelationshipsHandler, removeRelationshipHandler)
	contactDetailsHandler := NewContactDetailsHandler(addContactHandler, updateContactHandler, listContactsHandler, removeContactHandler)
	partyAddressHandler := NewPartyAddressHandler(addAddressHandler, listAddressesHandler)

	router := gin.New()
	router.POST("/parties", partyHandler.CreateParty)
	router.GET("/parties", partyHandler.ListParties)
	router.GET("/parties/:id", partyHandler.GetParty)
	router.PUT("/parties/:id", partyHandler.UpdateParty)
	router.PATCH("/parties/:id/status", partyHandler.ChangePartyStatus)
	router.POST("/parties/:id/roles", partyRoleHandler.AddRole)
	router.DELETE("/parties/:id/roles/:role", partyRoleHandler.RemoveRole)
	router.POST("/parties/:id/relationships", partyRelationshipHandler.AddRelationship)
	router.GET("/parties/:id/relationships", partyRelationshipHandler.ListRelationships)
	router.DELETE("/parties/:id/relationships/:relationship_id", partyRelationshipHandler.RemoveRelationship)
	router.POST("/parties/:id/contact-details", contactDetailsHandler.AddContactDetails)
	router.PUT("/parties/:id/contact-details/:contact_id", contactDetailsHandler.UpdateContactDetails)
	router.DELETE("/parties/:id/contact-details/:contact_id", contactDetailsHandler.RemoveContactDetails)
	router.GET("/parties/:id/contact-details", contactDetailsHandler.ListContactDetails)
	router.POST("/parties/:id/addresses", partyAddressHandler.AddAddress)
	router.GET("/parties/:id/addresses", partyAddressHandler.ListAddresses)

	return router, partyRepo
}

func performRequest(router *gin.Engine, method, path string, payload interface{}) *httptest.ResponseRecorder {
	var body *bytes.Buffer
	if payload != nil {
		data, _ := json.Marshal(payload)
		body = bytes.NewBuffer(data)
	} else {
		body = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func TestPartyHandler_CreateGetUpdateStatus(t *testing.T) {
	router, _ := setupHandlers()

	createPayload := map[string]interface{}{
		"id":     "party-100",
		"status": "ACTIVE",
		"roles":  []string{"CLIENT"},
		"person_profile": map[string]interface{}{
			"first_name": "Ana",
			"last_name":  "Perez",
		},
	}

	resp := performRequest(router, http.MethodPost, "/parties", createPayload)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/parties/party-100", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	updatePayload := map[string]interface{}{
		"status": "INACTIVE",
		"person_profile": map[string]interface{}{
			"first_name": "Ana",
			"last_name":  "Diaz",
		},
	}
	resp = performRequest(router, http.MethodPut, "/parties/party-100", updatePayload)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodPatch, "/parties/party-100/status", map[string]interface{}{"status": "ACTIVE"})
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestPartyRoleHandler_AddRemove(t *testing.T) {
	router, _ := setupHandlers()

	createPayload := map[string]interface{}{
		"id":     "party-200",
		"status": "ACTIVE",
		"roles":  []string{},
		"person_profile": map[string]interface{}{
			"first_name": "Luis",
			"last_name":  "Lopez",
		},
	}
	_ = performRequest(router, http.MethodPost, "/parties", createPayload)

	resp := performRequest(router, http.MethodPost, "/parties/party-200/roles", map[string]interface{}{"role": "CLIENT"})
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodDelete, "/parties/party-200/roles/CLIENT", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestContactDetailsHandler_AddUpdateRemove(t *testing.T) {
	router, repo := setupHandlers()

	partyID, _ := domain.NewPartyID("party-300")
	orgProfile, _ := domain.NewOrganizationProfile("Org", nil, "")
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, "system", nil, orgProfile)
	_ = repo.Save(context.Background(), party)

	addPayload := map[string]interface{}{
		"id":               "contact-1",
		"type_description": "Ventas",
		"email":            "ventas@org.local",
	}
	resp := performRequest(router, http.MethodPost, "/parties/party-300/contact-details", addPayload)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}

	updatePayload := map[string]interface{}{
		"type_description": "Soporte",
	}
	resp = performRequest(router, http.MethodPut, "/parties/party-300/contact-details/contact-1", updatePayload)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodDelete, "/parties/party-300/contact-details/contact-1", nil)
	if resp.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.Code)
	}
}
