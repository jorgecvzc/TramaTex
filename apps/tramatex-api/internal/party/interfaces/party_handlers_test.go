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

func (r *fakePartyRepo) Save(ctx context.Context, party *domain.Party, createdBy string, modifiedBy string) error {
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

func (r *fakeRelationshipRepo) Save(ctx context.Context, relationship domain.PartyRelationship, createdBy string, modifiedBy string) error {
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

type fakeAddressRepo struct {
	addressesByParty map[string][]*domain.Address
}

func newFakeAddressRepo() *fakeAddressRepo {
	return &fakeAddressRepo{addressesByParty: make(map[string][]*domain.Address)}
}

func (r *fakeAddressRepo) Save(ctx context.Context, address *domain.Address, addressID domain.AddressID, partyID domain.PartyID, createdBy string, modifiedBy string) error {
	if address == nil {
		return nil
	}
	r.addressesByParty[partyID.String()] = append(r.addressesByParty[partyID.String()], address)
	return nil
}

func (r *fakeAddressRepo) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]*domain.Address, error) {
	return r.addressesByParty[partyID.String()], nil
}

func (r *fakeAddressRepo) FindPrimary(ctx context.Context, partyID domain.PartyID) (*domain.Address, error) {
	addresses := r.addressesByParty[partyID.String()]
	if len(addresses) == 0 {
		return nil, nil
	}
	return addresses[0], nil
}

func (r *fakeAddressRepo) Delete(ctx context.Context, id domain.AddressID) error {
	return nil
}

func setupHandlers() (*gin.Engine, *fakePartyRepo) {
	gin.SetMode(gin.TestMode)

	partyRepo := newFakePartyRepo()
	relRepo := newFakeRelationshipRepo()
	addressRepo := newFakeAddressRepo()

	createPartyHandler := application.NewCreatePartyHandler(partyRepo)
	updatePartyHandler := application.NewUpdatePartyHandler(partyRepo)
	changePartyStatusHandler := application.NewChangePartyStatusHandler(partyRepo)
	getPartyHandler := application.NewGetPartyHandler(partyRepo)
	listPartiesHandler := application.NewListPartiesHandler(partyRepo)
	getBatchHandler := application.NewGetPartiesBatchHandler(partyRepo)
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

	partyHandler := NewPartyHandler(createPartyHandler, updatePartyHandler, changePartyStatusHandler, getPartyHandler, listPartiesHandler, getBatchHandler)
	partyRoleHandler := NewPartyRoleHandler(addPartyRoleHandler, removePartyRoleHandler)
	partyRelationshipHandler := NewPartyRelationshipHandler(addRelationshipHandler, listRelationshipsHandler, removeRelationshipHandler)
	contactDetailsHandler := NewContactDetailsHandler(addContactHandler, updateContactHandler, listContactsHandler, removeContactHandler)
	partyAddressHandler := NewPartyAddressHandler(addAddressHandler, listAddressesHandler)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
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

func setupHandlersWithoutUser() *gin.Engine {
	gin.SetMode(gin.TestMode)

	partyRepo := newFakePartyRepo()
	relRepo := newFakeRelationshipRepo()
	addressRepo := newFakeAddressRepo()

	createPartyHandler := application.NewCreatePartyHandler(partyRepo)
	updatePartyHandler := application.NewUpdatePartyHandler(partyRepo)
	changePartyStatusHandler := application.NewChangePartyStatusHandler(partyRepo)
	getPartyHandler := application.NewGetPartyHandler(partyRepo)
	listPartiesHandler := application.NewListPartiesHandler(partyRepo)
	getBatchHandler := application.NewGetPartiesBatchHandler(partyRepo)
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

	partyHandler := NewPartyHandler(createPartyHandler, updatePartyHandler, changePartyStatusHandler, getPartyHandler, listPartiesHandler, getBatchHandler)
	partyRoleHandler := NewPartyRoleHandler(addPartyRoleHandler, removePartyRoleHandler)
	partyRelationshipHandler := NewPartyRelationshipHandler(addRelationshipHandler, listRelationshipsHandler, removeRelationshipHandler)
	contactDetailsHandler := NewContactDetailsHandler(addContactHandler, updateContactHandler, listContactsHandler, removeContactHandler)
	partyAddressHandler := NewPartyAddressHandler(addAddressHandler, listAddressesHandler)

	router := gin.New()
	router.POST("/parties", partyHandler.CreateParty)
	router.PATCH("/parties/:id/status", partyHandler.ChangePartyStatus)
	router.POST("/parties/:id/roles", partyRoleHandler.AddRole)
	router.POST("/parties/:id/relationships", partyRelationshipHandler.AddRelationship)
	router.POST("/parties/:id/contact-details", contactDetailsHandler.AddContactDetails)
	router.PUT("/parties/:id/contact-details/:contact_id", contactDetailsHandler.UpdateContactDetails)
	router.DELETE("/parties/:id/contact-details/:contact_id", contactDetailsHandler.RemoveContactDetails)
	router.POST("/parties/:id/addresses", partyAddressHandler.AddAddress)

	return router
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

func TestPartyHandler_CreateParty_InvalidJSON(t *testing.T) {
	router, _ := setupHandlers()

	resp := performRequest(router, http.MethodPost, "/parties", "{")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestPartyHandler_CreateParty_WithOrganizationProfile(t *testing.T) {
	router, _ := setupHandlers()

	createPayload := map[string]interface{}{
		"id":     "party-org-1",
		"status": "ACTIVE",
		"roles":  []string{"SUPPLIER"},
		"organization_profile": map[string]interface{}{
			"name":        "Org",
			"tax_id":      "B12345678",
			"tax_id_type": "CIF",
			"website":     "https://org.local",
			"contacts": []map[string]interface{}{
				{
					"id":               "contact-1",
					"type_description": "Ventas",
					"phone":            "+34 600 111 222",
					"email":            "ventas@org.local",
					"related_party_id": "party-related",
				},
			},
		},
	}

	resp := performRequest(router, http.MethodPost, "/parties", createPayload)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}
}

func TestPartyHandler_CreateParty_Unauthorized(t *testing.T) {
	router := setupHandlersWithoutUser()

	createPayload := map[string]interface{}{
		"id":     "party-unauth",
		"status": "ACTIVE",
		"roles":  []string{"CLIENT"},
		"person_profile": map[string]interface{}{
			"first_name": "Ana",
			"last_name":  "Perez",
		},
	}

	resp := performRequest(router, http.MethodPost, "/parties", createPayload)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestPartyHandler_ChangeStatus_Unauthorized(t *testing.T) {
	router := setupHandlersWithoutUser()

	resp := performRequest(router, http.MethodPatch, "/parties/party-unauth/status", map[string]interface{}{"status": "ACTIVE"})
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestPartyHandler_UpdateParty_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	partyRepo := newFakePartyRepo()
	createPartyHandler := application.NewCreatePartyHandler(partyRepo)
	updatePartyHandler := application.NewUpdatePartyHandler(partyRepo)
	changePartyStatusHandler := application.NewChangePartyStatusHandler(partyRepo)
	getPartyHandler := application.NewGetPartyHandler(partyRepo)
	listPartiesHandler := application.NewListPartiesHandler(partyRepo)
	getBatchHandler := application.NewGetPartiesBatchHandler(partyRepo)

	partyHandler := NewPartyHandler(createPartyHandler, updatePartyHandler, changePartyStatusHandler, getPartyHandler, listPartiesHandler, getBatchHandler)

	router := gin.New()
	router.PUT("/parties/:id", partyHandler.UpdateParty)

	resp := performRequest(router, http.MethodPut, "/parties/party-unauth", map[string]interface{}{"status": "ACTIVE"})
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestPartyHandler_UpdateParty_InvalidJSON(t *testing.T) {
	router, _ := setupHandlers()

	createPayload := map[string]interface{}{
		"id":     "party-101",
		"status": "ACTIVE",
		"roles":  []string{},
		"person_profile": map[string]interface{}{
			"first_name": "Ana",
			"last_name":  "Perez",
		},
	}
	_ = performRequest(router, http.MethodPost, "/parties", createPayload)

	resp := performRequest(router, http.MethodPut, "/parties/party-101", "{")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestPartyHandler_ChangeStatus_InvalidJSON(t *testing.T) {
	router, _ := setupHandlers()

	createPayload := map[string]interface{}{
		"id":     "party-102",
		"status": "ACTIVE",
		"roles":  []string{},
		"person_profile": map[string]interface{}{
			"first_name": "Ana",
			"last_name":  "Perez",
		},
	}
	_ = performRequest(router, http.MethodPost, "/parties", createPayload)

	resp := performRequest(router, http.MethodPatch, "/parties/party-102/status", "{")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestPartyHandler_ChangeStatus_InvalidStatus(t *testing.T) {
	router, _ := setupHandlers()

	createPayload := map[string]interface{}{
		"id":     "party-103",
		"status": "ACTIVE",
		"roles":  []string{},
		"person_profile": map[string]interface{}{
			"first_name": "Ana",
			"last_name":  "Perez",
		},
	}
	_ = performRequest(router, http.MethodPost, "/parties", createPayload)

	resp := performRequest(router, http.MethodPatch, "/parties/party-103/status", map[string]interface{}{"status": "BAD"})
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
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

func TestPartyRoleHandler_AddRole_InvalidJSON(t *testing.T) {
	router, _ := setupHandlers()

	createPayload := map[string]interface{}{
		"id":     "party-201",
		"status": "ACTIVE",
		"roles":  []string{},
		"person_profile": map[string]interface{}{
			"first_name": "Ana",
			"last_name":  "Perez",
		},
	}
	_ = performRequest(router, http.MethodPost, "/parties", createPayload)

	resp := performRequest(router, http.MethodPost, "/parties/party-201/roles", "{")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestPartyRoleHandler_RemoveRole_InvalidRole(t *testing.T) {
	router, _ := setupHandlers()

	createPayload := map[string]interface{}{
		"id":     "party-202",
		"status": "ACTIVE",
		"roles":  []string{},
		"person_profile": map[string]interface{}{
			"first_name": "Ana",
			"last_name":  "Perez",
		},
	}
	_ = performRequest(router, http.MethodPost, "/parties", createPayload)

	resp := performRequest(router, http.MethodDelete, "/parties/party-202/roles/BAD", nil)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestPartyRelationshipHandler_AddRelationship_InvalidJSON(t *testing.T) {
	router, _ := setupHandlers()

	resp := performRequest(router, http.MethodPost, "/parties/party-300/relationships", "{")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestPartyRelationshipHandler_AddRelationship_InvalidType(t *testing.T) {
	router, _ := setupHandlers()

	addPayload := map[string]interface{}{
		"id":          "rel-002",
		"to_party_id": "party-301",
		"type":        "BAD",
	}
	resp := performRequest(router, http.MethodPost, "/parties/party-300/relationships", addPayload)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestPartyRelationshipHandler_RemoveRelationship_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	relRepo := newFakeRelationshipRepo()
	partyRelationshipHandler := NewPartyRelationshipHandler(
		application.NewAddPartyRelationshipHandler(relRepo),
		application.NewListPartyRelationshipsHandler(relRepo),
		application.NewRemovePartyRelationshipHandler(relRepo),
	)

	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/parties/party-1/relationships/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "relationship_id", Value: ""}}

	partyRelationshipHandler.RemoveRelationship(ctx)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestContactDetailsHandler_AddUpdateRemove(t *testing.T) {
	router, repo := setupHandlers()

	partyID, _ := domain.NewPartyID("party-300")
	orgProfile, _ := domain.NewOrganizationProfile("Org", nil, "")
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, nil, orgProfile)
	_ = repo.Save(context.Background(), party, "system", "system")

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

func TestContactDetailsHandler_AddContactDetails_InvalidJSON(t *testing.T) {
	router, _ := setupHandlers()

	resp := performRequest(router, http.MethodPost, "/parties/party-400/contact-details", "{")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestContactDetailsHandler_UpdateContactDetails_InvalidJSON(t *testing.T) {
	router, _ := setupHandlers()

	resp := performRequest(router, http.MethodPut, "/parties/party-400/contact-details/contact-1", "{")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestContactDetailsHandler_ListContactDetails_PartyNotFound(t *testing.T) {
	router, _ := setupHandlers()

	resp := performRequest(router, http.MethodGet, "/parties/missing/contact-details", nil)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestPartyAddressHandler_AddAddress_InvalidJSON(t *testing.T) {
	router, _ := setupHandlers()

	resp := performRequest(router, http.MethodPost, "/parties/party-500/addresses", "{")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	if getUserIDFromContext(nil) != "system" {
		t.Fatalf("expected system for nil context")
	}

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set("userID", "user-9")
	if getUserIDFromContext(ctx) != "user-9" {
		t.Fatalf("expected user-9 from context")
	}
}

func TestActorIDFromContext_InvalidValues(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set("userID", "")
	if _, ok := actorIDFromContext(ctx); ok {
		t.Fatalf("expected empty actor ID to be rejected")
	}

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set("userID", 123)
	if _, ok := actorIDFromContext(ctx); ok {
		t.Fatalf("expected non-string actor ID to be rejected")
	}
}

func TestPartyHandler_ListParties(t *testing.T) {
	router, _ := setupHandlers()

	createPayload := map[string]interface{}{
		"id":     "party-400",
		"status": "ACTIVE",
		"roles":  []string{"CLIENT"},
		"person_profile": map[string]interface{}{
			"first_name": "Mia",
			"last_name":  "Lopez",
		},
	}
	_ = performRequest(router, http.MethodPost, "/parties", createPayload)

	resp := performRequest(router, http.MethodGet, "/parties?page=2&page_size=5", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var payload struct {
		Total      int `json:"total"`
		PageNumber int `json:"page_number"`
		PageSize   int `json:"page_size"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if payload.Total != 1 {
		t.Fatalf("expected total 1, got %d", payload.Total)
	}
	if payload.PageNumber != 2 || payload.PageSize != 5 {
		t.Fatalf("expected page 2 size 5, got %d size %d", payload.PageNumber, payload.PageSize)
	}
}

func TestPartyHandler_ListParties_InvalidStatus(t *testing.T) {
	router, _ := setupHandlers()

	resp := performRequest(router, http.MethodGet, "/parties?status=UNKNOWN", nil)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestPartyRelationshipHandler_AddListRemove(t *testing.T) {
	router, _ := setupHandlers()

	addPayload := map[string]interface{}{
		"id":          "rel-001",
		"to_party_id": "party-501",
		"type":        "IS_EMPLOYEE_OF",
	}
	resp := performRequest(router, http.MethodPost, "/parties/party-500/relationships", addPayload)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/parties/party-500/relationships", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
	var listPayload struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &listPayload); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if listPayload.Total != 1 {
		t.Fatalf("expected total 1, got %d", listPayload.Total)
	}

	resp = performRequest(router, http.MethodDelete, "/parties/party-500/relationships/rel-001", nil)
	if resp.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.Code)
	}
}

func TestPartyRelationshipHandler_ListRelationships_InvalidPartyID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	relRepo := newFakeRelationshipRepo()
	partyRelationshipHandler := NewPartyRelationshipHandler(
		application.NewAddPartyRelationshipHandler(relRepo),
		application.NewListPartyRelationshipsHandler(relRepo),
		application.NewRemovePartyRelationshipHandler(relRepo),
	)

	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/parties//relationships", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

	partyRelationshipHandler.ListRelationships(ctx)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestContactDetailsHandler_List(t *testing.T) {
	router, repo := setupHandlers()

	partyID, _ := domain.NewPartyID("party-600")
	orgProfile, _ := domain.NewOrganizationProfile("Org", nil, "")
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, nil, orgProfile)
	contactID, _ := domain.NewContactDetailsID("contact-600")
	contact, _ := domain.NewContactDetails(contactID, "Ventas", nil, nil, nil)
	_ = party.OrganizationProfile().AddContact(contact)
	_ = repo.Save(context.Background(), party, "system", "system")

	resp := performRequest(router, http.MethodGet, "/parties/party-600/contact-details", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
	var payload struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if payload.Total != 1 {
		t.Fatalf("expected total 1, got %d", payload.Total)
	}
}

func TestPartyAddressHandler_AddList(t *testing.T) {
	router, _ := setupHandlers()

	addPayload := map[string]interface{}{
		"id":          "addr-900",
		"street":      "Calle 9",
		"city":        "Valencia",
		"province":    "Valencia",
		"postal_code": "46001",
		"country":     "Spain",
	}
	resp := performRequest(router, http.MethodPost, "/parties/party-900/addresses", addPayload)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/parties/party-900/addresses", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
	var payload struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if payload.Total != 1 {
		t.Fatalf("expected total 1, got %d", payload.Total)
	}
}

func TestPartyAddressHandler_ListAddresses_Empty(t *testing.T) {
	router, _ := setupHandlers()

	resp := performRequest(router, http.MethodGet, "/parties/party-empty/addresses", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
	var payload struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if payload.Total != 0 {
		t.Fatalf("expected total 0, got %d", payload.Total)
	}
}

type failingAddressRepo struct{}

func (r *failingAddressRepo) Save(ctx context.Context, address *domain.Address, addressID domain.AddressID, partyID domain.PartyID, createdBy string, modifiedBy string) error {
	return nil
}

func (r *failingAddressRepo) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]*domain.Address, error) {
	return nil, errors.New("db error")
}

func (r *failingAddressRepo) FindPrimary(ctx context.Context, partyID domain.PartyID) (*domain.Address, error) {
	return nil, nil
}

func (r *failingAddressRepo) Delete(ctx context.Context, id domain.AddressID) error {
	return nil
}

func TestPartyAddressHandler_ListAddresses_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	addressRepo := &failingAddressRepo{}
	partyAddressHandler := NewPartyAddressHandler(
		application.NewAddPartyAddressHandler(addressRepo),
		application.NewListPartyAddressesHandler(addressRepo),
	)

	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/parties/party-1/addresses", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "party-1"}}

	partyAddressHandler.ListAddresses(ctx)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestContactDetailsHandler_UpdateContactDetails_NotFound(t *testing.T) {
	router, repo := setupHandlers()

	partyID, _ := domain.NewPartyID("party-910")
	orgProfile, _ := domain.NewOrganizationProfile("Org", nil, "")
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, nil, orgProfile)
	_ = repo.Save(context.Background(), party, "system", "system")

	updatePayload := map[string]interface{}{
		"type_description": "Soporte",
	}
	resp := performRequest(router, http.MethodPut, "/parties/party-910/contact-details/contact-missing", updatePayload)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestContactDetailsHandler_RemoveContactDetails_NotFound(t *testing.T) {
	router, repo := setupHandlers()

	partyID, _ := domain.NewPartyID("party-911")
	orgProfile, _ := domain.NewOrganizationProfile("Org", nil, "")
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, nil, orgProfile)
	_ = repo.Save(context.Background(), party, "system", "system")

	resp := performRequest(router, http.MethodDelete, "/parties/party-911/contact-details/contact-missing", nil)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestPartyHandler_GetParty_NotFound(t *testing.T) {
	router, _ := setupHandlers()

	resp := performRequest(router, http.MethodGet, "/parties/missing", nil)
	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.Code)
	}
}

func TestPartyHandler_GetParty_MissingID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	partyRepo := newFakePartyRepo()
	getBatchHandler := application.NewGetPartiesBatchHandler(partyRepo)
	partyHandler := NewPartyHandler(
		application.NewCreatePartyHandler(partyRepo),
		application.NewUpdatePartyHandler(partyRepo),
		application.NewChangePartyStatusHandler(partyRepo),
		application.NewGetPartyHandler(partyRepo),
		application.NewListPartiesHandler(partyRepo),
		getBatchHandler,
	)

	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/parties", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

	partyHandler.GetParty(ctx)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestPartyHandler_ListParties_InvalidRole(t *testing.T) {
	router, _ := setupHandlers()

	resp := performRequest(router, http.MethodGet, "/parties?role=BAD", nil)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestPartyRoleHandler_AddRole_Unauthorized(t *testing.T) {
	router := setupHandlersWithoutUser()

	resp := performRequest(router, http.MethodPost, "/parties/party-700/roles", map[string]interface{}{"role": "CLIENT"})
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestPartyRelationshipHandler_AddRelationship_Unauthorized(t *testing.T) {
	router := setupHandlersWithoutUser()

	resp := performRequest(router, http.MethodPost, "/parties/party-800/relationships", map[string]interface{}{"id": "rel-1", "to_party_id": "party-801", "type": "IS_EMPLOYEE_OF"})
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestContactDetailsHandler_UpdateContactDetails_Unauthorized(t *testing.T) {
	router := setupHandlersWithoutUser()

	resp := performRequest(router, http.MethodPut, "/parties/party-900/contact-details/contact-9", map[string]interface{}{"type_description": "Soporte"})
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestContactDetailsHandler_RemoveContactDetails_Unauthorized(t *testing.T) {
	router := setupHandlersWithoutUser()

	resp := performRequest(router, http.MethodDelete, "/parties/party-900/contact-details/contact-9", nil)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestPartyAddressHandler_AddAddress_Unauthorized(t *testing.T) {
	router := setupHandlersWithoutUser()

	resp := performRequest(router, http.MethodPost, "/parties/party-900/addresses", map[string]interface{}{"id": "addr-1", "street": "Calle 1", "city": "Madrid", "postal_code": "28001", "country": "Spain"})
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestPartyHandler_ListParties_InvalidPaging(t *testing.T) {
	router, _ := setupHandlers()

	resp := performRequest(router, http.MethodGet, "/parties?page=bad&page_size=bad", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestPartyHandlers_Constructors(t *testing.T) {
	partyRepo := newFakePartyRepo()
	relRepo := newFakeRelationshipRepo()
	addressRepo := newFakeAddressRepo()

	if NewPartyRelationshipHandler(application.NewAddPartyRelationshipHandler(relRepo), application.NewListPartyRelationshipsHandler(relRepo), application.NewRemovePartyRelationshipHandler(relRepo)) == nil {
		t.Fatalf("expected party relationship handler")
	}
	if NewContactDetailsHandler(application.NewAddContactDetailsHandler(partyRepo), application.NewUpdateContactDetailsHandler(partyRepo), application.NewListContactDetailsHandler(partyRepo), application.NewRemoveContactDetailsHandler(partyRepo)) == nil {
		t.Fatalf("expected contact details handler")
	}
	if NewPartyAddressHandler(application.NewAddPartyAddressHandler(addressRepo), application.NewListPartyAddressesHandler(addressRepo)) == nil {
		t.Fatalf("expected party address handler")
	}
}

func TestGetUserIDFromContext_Defaults(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	if getUserIDFromContext(ctx) != "system" {
		t.Fatalf("expected system for missing userID")
	}

	ctx.Set("userID", 123)
	if getUserIDFromContext(ctx) != "system" {
		t.Fatalf("expected system for non-string userID")
	}
}
