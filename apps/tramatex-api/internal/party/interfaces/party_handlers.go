package interfaces

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joran-cortez/tramatex/internal/party/application"
)

// PartyHandler handles party endpoints

type PartyHandler struct {
	createHandler       *application.CreatePartyHandler
	updateHandler       *application.UpdatePartyHandler
	changeStatusHandler *application.ChangePartyStatusHandler
	deleteHandler       *application.DeletePartyHandler
	getHandler          *application.GetPartyHandler
	listHandler         *application.ListPartiesHandler
	getBatchHandler     *application.GetPartiesBatchHandler
}

func actorIDFromContext(c *gin.Context) (string, bool) {
	value, ok := c.Get("userID")
	if !ok {
		return "", false
	}
	actorID, ok := value.(string)
	if !ok || actorID == "" {
		return "", false
	}
	return actorID, true
}

// SplitAndTrim splits a string by separator and trims whitespace from each part
func SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func stringPtr(value string) *string {
	return &value
}

func NewPartyHandler(
	createHandler *application.CreatePartyHandler,
	updateHandler *application.UpdatePartyHandler,
	changeStatusHandler *application.ChangePartyStatusHandler,
	deleteHandler *application.DeletePartyHandler,
	getHandler *application.GetPartyHandler,
	listHandler *application.ListPartiesHandler,
	getBatchHandler *application.GetPartiesBatchHandler,
) *PartyHandler {
	return &PartyHandler{
		createHandler:       createHandler,
		updateHandler:       updateHandler,
		changeStatusHandler: changeStatusHandler,
		deleteHandler:       deleteHandler,
		getHandler:          getHandler,
		listHandler:         listHandler,
		getBatchHandler:     getBatchHandler,
	}
}

// CreateParty handles POST /parties
func (h *PartyHandler) CreateParty(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}
	var req CreatePartyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	cmd := &application.CreatePartyCommand{
		ID:                        req.ID,
		Status:                    req.Status,
		Roles:                     req.Roles,
		DefaultDiscountPercentage: req.DefaultDiscountPercentage,
		ActorID:                   actorID,
	}

	if req.PersonProfile != nil {
		cmd.PersonProfile = &application.PersonProfileInput{
			FirstName: stringPtr(req.PersonProfile.FirstName),
			LastName:  stringPtr(req.PersonProfile.LastName),
			Phone:     stringPtr(req.PersonProfile.Phone),
			Email:     stringPtr(req.PersonProfile.Email),
		}
	}

	if req.OrganizationProfile != nil {
		cmd.OrganizationProfile = &application.OrganizationProfileInput{
			Name:      stringPtr(req.OrganizationProfile.Name),
			TaxID:     stringPtr(req.OrganizationProfile.TaxID),
			TaxIDType: stringPtr(req.OrganizationProfile.TaxIDType),
			Website:   stringPtr(req.OrganizationProfile.Website),
			Phone:     stringPtr(req.OrganizationProfile.Phone),
			Email:     stringPtr(req.OrganizationProfile.Email),
			Notes:     stringPtr(req.OrganizationProfile.Notes),
		}
	}

	party, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, MapPartyToDTO(party))
}

// GetParty handles GET /parties/{id}
func (h *PartyHandler) GetParty(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Party ID is required"})
		return
	}

	party, err := h.getHandler.Handle(c.Request.Context(), &application.GetPartyQuery{ID: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MapPartyToDTO(party))
}

// GetPartiesBatch handles GET /parties/batch?ids=uuid1,uuid2,uuid3
func (h *PartyHandler) GetPartiesBatch(c *gin.Context) {
	idsParam := c.Query("ids")
	if idsParam == "" {
		c.JSON(http.StatusOK, []PartyBatchDTO{})
		return
	}

	// Split comma-separated IDs
	ids := []string{}
	for _, id := range SplitAndTrim(idsParam, ",") {
		if id != "" {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		c.JSON(http.StatusOK, []PartyBatchDTO{})
		return
	}

	parties, err := h.getBatchHandler.Handle(c.Request.Context(), &application.GetPartiesBatchQuery{IDs: ids})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map to minimal DTOs for batch operations
	dtos := make([]PartyBatchDTO, len(parties))
	for i, party := range parties {
		dtos[i] = MapPartyToBatchDTO(party)
	}

	c.JSON(http.StatusOK, dtos)
}

// ListParties handles GET /parties
func (h *PartyHandler) ListParties(c *gin.Context) {
	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if size, err := strconv.Atoi(ps); err == nil && size > 0 {
			pageSize = size
		}
	}

	pageNumber := 1
	if pn := c.Query("page"); pn != "" {
		if num, err := strconv.Atoi(pn); err == nil && num > 0 {
			pageNumber = num
		}
	}

	query := &application.ListPartiesQuery{
		Status:     c.Query("status"),
		Role:       c.Query("role"),
		Type:       c.Query("type"),
		Name:       c.Query("name"),
		TaxID:      c.Query("tax_id"),
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}

	parties, err := h.listHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]*PartyDTO, len(parties))
	for i, party := range parties {
		dto := MapPartyToDTO(party)
		canDelete, err := h.deleteHandler.CanDelete(c.Request.Context(), party.ID().String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to evaluate party deletability"})
			return
		}
		dto.CanDelete = canDelete
		dtos[i] = dto
	}

	response := ListResponse{
		Data:       dtos,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      len(parties),
	}
	c.JSON(http.StatusOK, response)
}

// UpdateParty handles PUT /parties/{id}
func (h *PartyHandler) UpdateParty(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Party ID is required"})
		return
	}

	var req UpdatePartyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	cmd := &application.UpdatePartyCommand{
		ID:                        id,
		Status:                    req.Status,
		DefaultDiscountPercentage: req.DefaultDiscountPercentage,
		ActorID:                   actorID,
	}

	if req.PersonProfile != nil {
		cmd.PersonProfile = &application.PersonProfileInput{
			FirstName: stringPtr(req.PersonProfile.FirstName),
			LastName:  stringPtr(req.PersonProfile.LastName),
			Phone:     stringPtr(req.PersonProfile.Phone),
			Email:     stringPtr(req.PersonProfile.Email),
		}
	}

	if req.OrganizationProfile != nil {
		cmd.OrganizationProfile = &application.OrganizationProfileInput{
			Name:      stringPtr(req.OrganizationProfile.Name),
			TaxID:     stringPtr(req.OrganizationProfile.TaxID),
			TaxIDType: stringPtr(req.OrganizationProfile.TaxIDType),
			Website:   stringPtr(req.OrganizationProfile.Website),
			Phone:     stringPtr(req.OrganizationProfile.Phone),
			Email:     stringPtr(req.OrganizationProfile.Email),
			Notes:     stringPtr(req.OrganizationProfile.Notes),
		}
	}

	party, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MapPartyToDTO(party))
}

// ChangePartyStatus handles PATCH /parties/{id}/status
func (h *PartyHandler) ChangePartyStatus(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Party ID is required"})
		return
	}

	var req ChangePartyStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	cmd := &application.ChangePartyStatusCommand{
		ID:      id,
		Status:  req.Status,
		ActorID: actorID,
	}

	party, err := h.changeStatusHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MapPartyToDTO(party))
}

// DeleteParty handles DELETE /parties/{id}
func (h *PartyHandler) DeleteParty(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Party ID is required"})
		return
	}

	if err := h.deleteHandler.Handle(c.Request.Context(), &application.DeletePartyCommand{
		ID:      id,
		ActorID: actorID,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// PartyRoleHandler handles role endpoints

type PartyRoleHandler struct {
	addHandler    *application.AddPartyRoleHandler
	removeHandler *application.RemovePartyRoleHandler
}

func NewPartyRoleHandler(addHandler *application.AddPartyRoleHandler, removeHandler *application.RemovePartyRoleHandler) *PartyRoleHandler {
	return &PartyRoleHandler{addHandler: addHandler, removeHandler: removeHandler}
}

// AddRole handles POST /parties/{id}/roles
func (h *PartyRoleHandler) AddRole(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}
	id := c.Param("id")
	var req AddPartyRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	party, err := h.addHandler.Handle(c.Request.Context(), &application.AddPartyRoleCommand{
		ID:      id,
		Role:    req.Role,
		ActorID: actorID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MapPartyToDTO(party))
}

// RemoveRole handles DELETE /parties/{id}/roles/{role}
func (h *PartyRoleHandler) RemoveRole(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}
	id := c.Param("id")
	role := c.Param("role")

	party, err := h.removeHandler.Handle(c.Request.Context(), &application.RemovePartyRoleCommand{
		ID:      id,
		Role:    role,
		ActorID: actorID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MapPartyToDTO(party))
}

// PartyRelationshipHandler handles relationship endpoints

type PartyRelationshipHandler struct {
	addHandler    *application.AddPartyRelationshipHandler
	listHandler   *application.ListPartyRelationshipsHandler
	removeHandler *application.RemovePartyRelationshipHandler
}

func NewPartyRelationshipHandler(
	addHandler *application.AddPartyRelationshipHandler,
	listHandler *application.ListPartyRelationshipsHandler,
	removeHandler *application.RemovePartyRelationshipHandler,
) *PartyRelationshipHandler {
	return &PartyRelationshipHandler{
		addHandler:    addHandler,
		listHandler:   listHandler,
		removeHandler: removeHandler,
	}
}

// AddRelationship handles POST /parties/{id}/relationships
func (h *PartyRelationshipHandler) AddRelationship(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}
	id := c.Param("id")
	var req CreateRelationshipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	relationship, err := h.addHandler.Handle(c.Request.Context(), &application.AddPartyRelationshipCommand{
		ID:             id,
		RelationshipID: req.ID,
		ToPartyID:      req.ToPartyID,
		Type:           req.Type,
		ActorID:        actorID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, MapPartyRelationshipToDTO(relationship))
}

// ListRelationships handles GET /parties/{id}/relationships
func (h *PartyRelationshipHandler) ListRelationships(c *gin.Context) {
	id := c.Param("id")
	relationships, err := h.listHandler.Handle(c.Request.Context(), &application.ListPartyRelationshipsQuery{PartyID: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]*PartyRelationshipDTO, len(relationships))
	for i := range relationships {
		rel := relationships[i]
		dtos[i] = MapPartyRelationshipToDTO(&rel)
	}

	response := ListResponse{Data: dtos, Total: len(dtos)}
	c.JSON(http.StatusOK, response)
}

// RemoveRelationship handles DELETE /parties/{id}/relationships/{relationship_id}
func (h *PartyRelationshipHandler) RemoveRelationship(c *gin.Context) {
	relID := c.Param("relationship_id")
	if err := h.removeHandler.Handle(c.Request.Context(), &application.RemovePartyRelationshipCommand{RelationshipID: relID}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ContactDetailsHandler handles contact detail endpoints

type ContactDetailsHandler struct {
	addHandler    *application.AddContactDetailsHandler
	updateHandler *application.UpdateContactDetailsHandler
	listHandler   *application.ListContactDetailsHandler
	removeHandler *application.RemoveContactDetailsHandler
}

func NewContactDetailsHandler(
	addHandler *application.AddContactDetailsHandler,
	updateHandler *application.UpdateContactDetailsHandler,
	listHandler *application.ListContactDetailsHandler,
	removeHandler *application.RemoveContactDetailsHandler,
) *ContactDetailsHandler {
	return &ContactDetailsHandler{
		addHandler:    addHandler,
		updateHandler: updateHandler,
		listHandler:   listHandler,
		removeHandler: removeHandler,
	}
}

// AddContactDetails handles POST /parties/{id}/contact-details
func (h *ContactDetailsHandler) AddContactDetails(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}
	id := c.Param("id")
	var req CreateContactDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	contact, err := h.addHandler.Handle(c.Request.Context(), &application.AddContactDetailsCommand{
		PartyID:         id,
		ContactID:       req.ID,
		TypeDescription: req.TypeDescription,
		Phone:           req.Phone,
		Email:           req.Email,
		RelatedPartyID:  req.RelatedPartyID,
		ActorID:         actorID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, MapContactDetailsToDTO(contact))
}

// ListContactDetails handles GET /parties/{id}/contact-details
func (h *ContactDetailsHandler) ListContactDetails(c *gin.Context) {
	id := c.Param("id")
	contacts, err := h.listHandler.Handle(c.Request.Context(), &application.ListContactDetailsQuery{PartyID: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]*ContactDetailsDTO, len(contacts))
	for i, contact := range contacts {
		dtos[i] = MapContactDetailsToDTO(contact)
	}

	response := ListResponse{Data: dtos, Total: len(dtos)}
	c.JSON(http.StatusOK, response)
}

// UpdateContactDetails handles PUT /parties/{id}/contact-details/{contact_id}
func (h *ContactDetailsHandler) UpdateContactDetails(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}
	partyID := c.Param("id")
	contactID := c.Param("contact_id")

	var req UpdateContactDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	contact, err := h.updateHandler.Handle(c.Request.Context(), &application.UpdateContactDetailsCommand{
		PartyID:         partyID,
		ContactID:       contactID,
		TypeDescription: req.TypeDescription,
		Phone:           req.Phone,
		Email:           req.Email,
		RelatedPartyID:  req.RelatedPartyID,
		ActorID:         actorID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MapContactDetailsToDTO(contact))
}

// RemoveContactDetails handles DELETE /parties/{id}/contact-details/{contact_id}
func (h *ContactDetailsHandler) RemoveContactDetails(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}
	partyID := c.Param("id")
	contactID := c.Param("contact_id")

	// Check query parameter to decide if we should delete the party if no other references
	deleteIfNoRefs := c.Query("deleteIfNoReferences") == "true"

	if err := h.removeHandler.Handle(c.Request.Context(), &application.RemoveContactDetailsCommand{
		PartyID:              partyID,
		ContactID:            contactID,
		ActorID:              actorID,
		DeleteIfNoReferences: deleteIfNoRefs,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// PartyAddressHandler handles party address endpoints

type PartyAddressHandler struct {
	addHandler    *application.AddPartyAddressHandler
	updateHandler *application.UpdatePartyAddressHandler
	removeHandler *application.RemovePartyAddressHandler
	listHandler   *application.ListPartyAddressesHandler
}

func NewPartyAddressHandler(
	addHandler *application.AddPartyAddressHandler,
	updateHandler *application.UpdatePartyAddressHandler,
	removeHandler *application.RemovePartyAddressHandler,
	listHandler *application.ListPartyAddressesHandler,
) *PartyAddressHandler {
	return &PartyAddressHandler{
		addHandler:    addHandler,
		updateHandler: updateHandler,
		removeHandler: removeHandler,
		listHandler:   listHandler,
	}
}

// AddAddress handles POST /parties/{id}/addresses
func (h *PartyAddressHandler) AddAddress(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}
	id := c.Param("id")
	var req CreatePartyAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	address, err := h.addHandler.Handle(c.Request.Context(), &application.AddPartyAddressCommand{
		PartyID:    id,
		AddressID:  req.ID,
		Street:     req.Street,
		City:       req.City,
		Province:   req.Province,
		PostalCode: req.PostalCode,
		Country:    req.Country,
		ActorID:    actorID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, MapAddressToDTO(address, req.ID))
}

// ListAddresses handles GET /parties/{id}/addresses
func (h *PartyAddressHandler) ListAddresses(c *gin.Context) {
	id := c.Param("id")
	addresses, err := h.listHandler.Handle(c.Request.Context(), &application.ListPartyAddressesQuery{PartyID: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]*AddressDTO, len(addresses))
	for i, addressWithID := range addresses {
		dtos[i] = MapAddressToDTO(addressWithID.Address, addressWithID.ID)
	}

	response := ListResponse{Data: dtos, Total: len(dtos)}
	c.JSON(http.StatusOK, response)
}

// UpdateAddress handles PUT /parties/{id}/addresses/{addressId}
func (h *PartyAddressHandler) UpdateAddress(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}
	partyID := c.Param("id")
	addressID := c.Param("addressId")

	var req CreatePartyAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	address, err := h.updateHandler.Handle(c.Request.Context(), &application.UpdatePartyAddressCommand{
		PartyID:    partyID,
		AddressID:  addressID,
		Street:     req.Street,
		City:       req.City,
		Province:   req.Province,
		PostalCode: req.PostalCode,
		Country:    req.Country,
		ActorID:    actorID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MapAddressToDTO(address, addressID))
}

// DeleteAddress handles DELETE /parties/{id}/addresses/{addressId}
func (h *PartyAddressHandler) DeleteAddress(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}
	addressID := c.Param("addressId")

	if err := h.removeHandler.Handle(c.Request.Context(), &application.RemovePartyAddressCommand{
		AddressID: addressID,
		ActorID:   actorID,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
