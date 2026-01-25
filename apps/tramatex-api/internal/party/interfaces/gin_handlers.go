
package interfaces

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joran-cortez/tramatex/internal/party/application"
)

// HTTP Handlers for Party module API

// OrganizationHandler handles organization endpoints
type OrganizationHandler struct {
	createHandler       *application.CreateOrganizationHandler
	updateHandler       *application.UpdateOrganizationHandler
	changeStatusHandler *application.ChangeOrganizationStatusHandler
	getHandler          *application.GetOrganizationHandler
	listHandler         *application.ListOrganizationsHandler
	listByRoleHandler   *application.ListOrganizationsByRoleHandler
}

// NewOrganizationHandler creates a new handler
func NewOrganizationHandler(
	createHandler *application.CreateOrganizationHandler,
	updateHandler *application.UpdateOrganizationHandler,
	changeStatusHandler *application.ChangeOrganizationStatusHandler,
	getHandler *application.GetOrganizationHandler,
	listHandler *application.ListOrganizationsHandler,
	listByRoleHandler *application.ListOrganizationsByRoleHandler,
) *OrganizationHandler {
	return &OrganizationHandler{
		createHandler:       createHandler,
		updateHandler:       updateHandler,
		changeStatusHandler: changeStatusHandler,
		getHandler:          getHandler,
		listHandler:         listHandler,
		listByRoleHandler:   listByRoleHandler,
	}
}

// CreateOrganization handles POST /organizations
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	cmd := &application.CreateOrganizationCommand{
		ID:        req.ID,
		Name:      req.Name,
		Role:      req.Role,
		TaxID:     req.TaxID,
		TaxIDType: req.TaxIDType,
		Website:   req.Website,
		CreatedBy: getUserIDFromContext(c), // Get from auth context
	}

	org, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, MapOrganizationToDTO(org))
}

// GetOrganization handles GET /organizations/{id}
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	query := &application.GetOrganizationQuery{ID: id}
	org, err := h.getHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MapOrganizationToDTO(org))
}

// ListOrganizations handles GET /organizations
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
	// Parse query parameters
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

	query := &application.ListOrganizationsQuery{
		Status:     c.Query("status"),
		Role:       c.Query("role"),
		Name:       c.Query("name"),
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}

	orgs, err := h.listHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]*OrganizationDTO, len(orgs))
	for i, org := range orgs {
		dtos[i] = MapOrganizationToDTO(org)
	}

	response := ListResponse{
		Data:       dtos,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      len(orgs),
	}
	c.JSON(http.StatusOK, response)
}

// UpdateOrganization handles PUT /organizations/{id}
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	var req UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	cmd := &application.UpdateOrganizationCommand{
		ID:         id,
		Name:       req.Name,
		Website:    req.Website,
		Notes:      req.Notes,
		ModifiedBy: getUserIDFromContext(c),
	}

	org, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MapOrganizationToDTO(org))
}

// ChangeStatus handles PATCH /organizations/{id}/status
func (h *OrganizationHandler) ChangeStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	var req ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	cmd := &application.ChangeOrganizationStatusCommand{
		ID:         id,
		Status:     req.Status,
		ModifiedBy: getUserIDFromContext(c),
	}

	org, err := h.changeStatusHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MapOrganizationToDTO(org))
}

// PersonHandler handles person endpoints
type PersonHandler struct {
	addPersonHandler         *application.AddPersonHandler
	getPersonHandler         *application.GetPersonHandler
	listPersonsHandler       *application.ListPersonsByOrganizationHandler
	getByEmailHandler        *application.GetPersonByEmailHandler
	getPrimaryContactHandler *application.GetPrimaryContactHandler
}

// NewPersonHandler creates a new handler
func NewPersonHandler(
	addPersonHandler *application.AddPersonHandler,
	getPersonHandler *application.GetPersonHandler,
	listPersonsHandler *application.ListPersonsByOrganizationHandler,
	getByEmailHandler *application.GetPersonByEmailHandler,
	getPrimaryContactHandler *application.GetPrimaryContactHandler,
) *PersonHandler {
	return &PersonHandler{
		addPersonHandler:         addPersonHandler,
		getPersonHandler:         getPersonHandler,
		listPersonsHandler:       listPersonsHandler,
		getByEmailHandler:        getByEmailHandler,
		getPrimaryContactHandler: getPrimaryContactHandler,
	}
}

// AddPerson handles POST /organizations/{org_id}/persons
func (h *PersonHandler) AddPerson(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	var req CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	cmd := &application.AddPersonCommand{
		ID:             req.ID,
		OrganizationID: orgID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Phone:          req.Phone,
		JobTitle:       req.JobTitle,
		IsPrimary:      req.IsPrimary,
		CreatedBy:      getUserIDFromContext(c),
	}

	person, err := h.addPersonHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, MapPersonToDTO(person))
}

// GetPerson handles GET /persons/{id}
func (h *PersonHandler) GetPerson(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Person ID is required"})
		return
	}

	query := &application.GetPersonQuery{ID: id}
	person, err := h.getPersonHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MapPersonToDTO(person))
}

// ListPersons handles GET /organizations/{org_id}/persons
func (h *PersonHandler) ListPersons(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	query := &application.ListPersonsByOrganizationQuery{OrganizationID: orgID}
	persons, err := h.listPersonsHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]*PersonDTO, len(persons))
	for i, p := range persons {
		dtos[i] = MapPersonToDTO(p)
	}

	response := ListResponse{
		Data:  dtos,
		Total: len(persons),
	}
	c.JSON(http.StatusOK, response)
}

// GetPrimaryContact handles GET /organizations/{org_id}/primary-contact
func (h *PersonHandler) GetPrimaryContact(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	query := &application.GetPrimaryContactQuery{OrganizationID: orgID}
	person, err := h.getPrimaryContactHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MapPersonToDTO(person))
}

// AddressHandler handles address endpoints
type AddressHandler struct {
	addAddressHandler        *application.AddAddressHandler
	listAddressesHandler     *application.ListAddressesByOrganizationHandler
	getPrimaryAddressHandler *application.GetPrimaryAddressHandler
}

// NewAddressHandler creates a new handler
func NewAddressHandler(
	addAddressHandler *application.AddAddressHandler,
	listAddressesHandler *application.ListAddressesByOrganizationHandler,
	getPrimaryAddressHandler *application.GetPrimaryAddressHandler,
) *AddressHandler {
	return &AddressHandler{
		addAddressHandler:        addAddressHandler,
		listAddressesHandler:     listAddressesHandler,
		getPrimaryAddressHandler: getPrimaryAddressHandler,
	}
}

// AddAddress handles POST /organizations/{org_id}/addresses
func (h *AddressHandler) AddAddress(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Country == "" {
		req.Country = "Spain"
	}

	cmd := &application.AddAddressCommand{
		ID:             req.ID,
		OrganizationID: orgID,
		Street:         req.Street,
		City:           req.City,
		Province:       req.Province,
		PostalCode:     req.PostalCode,
		Country:        req.Country,
		IsPrimary:      false,
		CreatedBy:      getUserIDFromContext(c),
	}

	address, err := h.addAddressHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, MapAddressToDTO(address))
}

// ListAddresses handles GET /organizations/{org_id}/addresses
func (h *AddressHandler) ListAddresses(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	query := &application.ListAddressesByOrganizationQuery{OrganizationID: orgID}
	addresses, err := h.listAddressesHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]*AddressDTO, len(addresses))
	for i, a := range addresses {
		dtos[i] = MapAddressToDTO(a)
	}

	response := ListResponse{
		Data:  dtos,
		Total: len(addresses),
	}
	c.JSON(http.StatusOK, response)
}

func getUserIDFromContext(c *gin.Context) string {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		return "anonymous"
	}
	return userID
}
