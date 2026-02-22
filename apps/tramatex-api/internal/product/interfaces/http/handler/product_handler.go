package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/joran-cortez/tramatex/internal/product/application"
	"github.com/joran-cortez/tramatex/internal/product/domain"
)

type ProductHandler struct {
	service *application.ProductService
}

func NewProductHandler(service *application.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func actorIDFromRequest(c *gin.Context) (string, bool) {
	actorID, ok := c.Request.Context().Value("actorID").(string)
	if !ok || actorID == "" {
		return "", false
	}
	return actorID, true
}

// mapErrorToHTTP maps domain errors to appropriate HTTP status codes
func mapErrorToHTTP(err error) int {
	if productErr, ok := err.(domain.ProductError); ok {
		switch productErr.Code {
		case domain.ErrCodeValidation:
			return http.StatusBadRequest
		case domain.ErrCodeNotFound:
			return http.StatusNotFound
		case domain.ErrCodeConflict:
			return http.StatusConflict
		case domain.ErrCodePersistence:
			return http.StatusInternalServerError
		}
	}
	return http.StatusInternalServerError
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var cmd application.CreateProductCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}
	cmd.ActorID = actorID

	product, err := h.service.CreateProduct(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(mapErrorToHTTP(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) AddGroupToProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var cmd application.AddGroupCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}
	cmd.ProductID = productID
	cmd.ActorID = actorID

	product, err := h.service.AddGroupToProduct(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) AddDirectAttributeToProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var cmd application.AddDirectAttributeCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}
	cmd.ProductID = productID
	cmd.ActorID = actorID

	product, err := h.service.AddDirectAttributeToProduct(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProductSKU(c *gin.Context) {
	var cmd application.UpdateProductSKUCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}
	cmd.ProductID = productID
	cmd.ActorID = actorID

	product, err := h.service.UpdateProductSKU(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req struct {
		Name               *string             `json:"name"`
		LongName           *string             `json:"long_name"`
		SKU                *string             `json:"sku"`
		Barcode            *string             `json:"barcode"`
		BasePrice          *float64            `json:"base_price"`
		ProductType        *domain.ProductType `json:"product_type"`
		Description        *string             `json:"description"`
		BrandID            *uuid.UUID          `json:"brand_id"`
		GroupIDs           []uuid.UUID         `json:"group_ids"`
		DirectAttributeIDs []uuid.UUID         `json:"direct_attribute_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}

	cmd := application.UpdateProductCommand{
		ActorID:            actorID,
		ProductID:          productID,
		Name:               req.Name,
		LongName:           req.LongName,
		SKU:                req.SKU,
		Barcode:            req.Barcode,
		BasePrice:          req.BasePrice,
		ProductType:        req.ProductType,
		Description:        req.Description,
		BrandID:            req.BrandID,
		GroupIDs:           req.GroupIDs,
		DirectAttributeIDs: req.DirectAttributeIDs,
	}

	product, err := h.service.UpdateProduct(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) CreateAttribute(c *gin.Context) {
	var cmd application.CreateAttributeCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}
	cmd.ActorID = actorID

	attribute, err := h.service.CreateAttribute(c.Request.Context(), cmd)
	if err != nil {
		// Check for validation errors (including duplicate code)
		if strings.Contains(err.Error(), "Ya existe un atributo con el código") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		// Check for database-level duplicate key constraint violation (fallback)
		if strings.Contains(err.Error(), "duplicate key value") && strings.Contains(err.Error(), "idx_attributes_code") {
			c.JSON(http.StatusConflict, gin.H{"error": "Ya existe un atributo con ese código"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, attribute)
}

func (h *ProductHandler) GetAttributeByID(c *gin.Context) {
	attributeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attribute id"})
		return
	}

	query := application.GetAttributeByIDQuery{ID: attributeID}
	attribute, err := h.service.GetAttributeByID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attribute)
}

func (h *ProductHandler) ListAttributes(c *gin.Context) {
	// Note: Scope-based filtering removed for MVP simplicity
	// All attributes are generic and returned without filtering
	var query application.ListAttributesQuery

	attributes, err := h.service.ListAttributes(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  attributes,
		"total": len(attributes),
	})
}

func (h *ProductHandler) UpdateAttribute(c *gin.Context) {
	attributeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attribute id"})
		return
	}

	var cmd application.UpdateAttributeCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}
	cmd.ID = attributeID // Set ID from URL parameter
	cmd.ActorID = actorID

	attribute, err := h.service.UpdateAttribute(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attribute)
}

// DeleteAttribute handles DELETE /attributes/:id
func (h *ProductHandler) DeleteAttribute(c *gin.Context) {
	attributeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attribute id"})
		return
	}

	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}

	cmd := application.DeleteAttributeCommand{
		ActorID: actorID,
		ID:      attributeID,
	}

	err = h.service.DeleteAttribute(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	query := application.GetProductByIDQuery{ID: productID}
	product, err := h.service.GetProductByID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	var query application.ListProductsQuery
	if brandIDStr := strings.TrimSpace(c.Query("brandId")); brandIDStr != "" {
		brandID, err := uuid.Parse(brandIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand id"})
			return
		}
		query.BrandID = &brandID
	}
	if groupIDStr := strings.TrimSpace(c.Query("groupId")); groupIDStr != "" {
		groupID, err := uuid.Parse(groupIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
			return
		}
		query.GroupID = &groupID
	}
	if isActiveStr := strings.TrimSpace(c.Query("isActive")); isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid isActive value"})
			return
		}
		query.IsActive = &isActive
	}

	products, err := h.service.ListProducts(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetCalculatedOptionSetsForProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	attributes, err := h.service.GetApplicableAttributesForProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attributes)
}

func (h *ProductHandler) GenerateProductVariants(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	cmd := application.GenerateProductVariantsCommand{ProductID: productID}
	err = h.service.GenerateProductVariants(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusAccepted) // 202 Accepted for async operation
}

func (h *ProductHandler) FindOrCreateProductVariant(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req struct {
		OptionConfiguration map[string]string `json:"optionConfiguration"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}

	cmd := application.FindOrCreateProductVariantCommand{
		ActorID:             actorID,
		ProductID:           productID,
		OptionConfiguration: req.OptionConfiguration,
	}

	variant, err := h.service.FindOrCreateProductVariant(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, variant)
}

func (h *ProductHandler) ListProductVariantsByProductID(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	query := application.ListProductVariantsByProductIDQuery{ProductID: productID}
	variants, err := h.service.ListProductVariantsByProductID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, variants)
}

func (h *ProductHandler) GetProductVariantByID(c *gin.Context) {
	variantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid variant id"})
		return
	}

	query := application.GetProductVariantByIDQuery{ID: variantID}
	variant, err := h.service.GetProductVariantByID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, variant)
}

func (h *ProductHandler) GetProductVariantBySKU(c *gin.Context) {
	sku := c.Query("sku")
	if sku == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sku query parameter is required"})
		return
	}

	query := application.GetProductVariantBySKUQuery{SKU: sku}
	variant, err := h.service.GetProductVariantBySKU(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, variant)
}

func (h *ProductHandler) UpdateProductVariant(c *gin.Context) {
	variantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid variant id"})
		return
	}

	var cmd application.UpdateProductVariantCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}
	cmd.ID = variantID
	cmd.ActorID = actorID

	variant, err := h.service.UpdateProductVariant(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, variant)
}

func (h *ProductHandler) CreatePartyServiceConfiguration(c *gin.Context) {
	partyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid party id"})
		return
	}

	var cmd application.CreatePartyServiceConfigurationCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}
	cmd.PartyID = partyID
	cmd.ActorID = actorID

	config, err := h.service.CreatePartyServiceConfiguration(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, config)
}

func (h *ProductHandler) ListPartyServiceConfigurationsByPartyID(c *gin.Context) {
	partyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid party id"})
		return
	}

	query := application.ListPartyServiceConfigurationsByPartyIDQuery{PartyID: partyID}
	configs, err := h.service.ListPartyServiceConfigurationsByPartyID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}

func (h *ProductHandler) GetPartyServiceConfigurationByID(c *gin.Context) {
	partyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid party id"})
		return
	}
	configID, err := uuid.Parse(c.Param("configId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid configuration id"})
		return
	}

	query := application.GetPartyServiceConfigurationByIDQuery{PartyID: partyID, ID: configID}
	config, err := h.service.GetPartyServiceConfigurationByID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

func (h *ProductHandler) UpdatePartyServiceConfiguration(c *gin.Context) {
	partyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid party id"})
		return
	}
	configID, err := uuid.Parse(c.Param("configId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid configuration id"})
		return
	}

	var cmd application.UpdatePartyServiceConfigurationCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}
	cmd.ID = configID
	cmd.PartyID = partyID
	cmd.ActorID = actorID

	config, err := h.service.UpdatePartyServiceConfiguration(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

func (h *ProductHandler) DeletePartyServiceConfiguration(c *gin.Context) {
	partyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid party id"})
		return
	}
	configID, err := uuid.Parse(c.Param("configId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid configuration id"})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}

	cmd := application.DeletePartyServiceConfigurationCommand{ID: configID, PartyID: partyID, ActorID: actorID}
	err = h.service.DeletePartyServiceConfiguration(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ============================================================================
// Brand Handlers
// ============================================================================

// ListBrands handles GET /brands
func (h *ProductHandler) ListBrands(c *gin.Context) {
	brands, err := h.service.ListBrands(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  brands,
		"total": len(brands),
	})
}

// GetBrandByID handles GET /brands/:id
func (h *ProductHandler) GetBrandByID(c *gin.Context) {
	brandID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand id"})
		return
	}

	brand, err := h.service.GetBrandByID(c.Request.Context(), brandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, brand)
}

// CreateBrand handles POST /brands
func (h *ProductHandler) CreateBrand(c *gin.Context) {
	var cmd application.CreateBrandCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}
	cmd.ActorID = actorID

	brand, err := h.service.CreateBrand(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, brand)
}

// UpdateBrand handles PUT /brands/:id
func (h *ProductHandler) UpdateBrand(c *gin.Context) {
	brandID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand id"})
		return
	}

	var cmd application.UpdateBrandCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}
	cmd.ID = brandID
	cmd.ActorID = actorID

	brand, err := h.service.UpdateBrand(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, brand)
}

// DeleteBrand handles DELETE /brands/:id
func (h *ProductHandler) DeleteBrand(c *gin.Context) {
	brandID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand id"})
		return
	}

	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}

	cmd := application.DeleteBrandCommand{
		ActorID: actorID,
		ID:      brandID,
	}

	err = h.service.DeleteBrand(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ============================================================================
// Product Group Handlers
// ============================================================================

// ListProductGroups handles GET /product-groups
func (h *ProductHandler) ListProductGroups(c *gin.Context) {
	groups, err := h.service.ListProductGroups(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  groups,
		"total": len(groups),
	})
}

// GetProductGroupByID handles GET /product-groups/:id
func (h *ProductHandler) GetProductGroupByID(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product group id"})
		return
	}

	group, err := h.service.GetProductGroupByID(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// CreateProductGroup handles POST /product-groups
func (h *ProductHandler) CreateProductGroup(c *gin.Context) {
	var cmd application.CreateProductGroupCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}
	cmd.ActorID = actorID

	group, err := h.service.CreateProductGroup(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// UpdateProductGroup handles PUT /product-groups/:id
func (h *ProductHandler) UpdateProductGroup(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product group id"})
		return
	}

	var cmd application.UpdateProductGroupCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}
	cmd.ID = groupID
	cmd.ActorID = actorID

	group, err := h.service.UpdateProductGroup(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// DeleteProductGroup handles DELETE /product-groups/:id
func (h *ProductHandler) DeleteProductGroup(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product group id"})
		return
	}

	actorID, ok := actorIDFromRequest(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing actor ID"})
		return
	}

	cmd := application.DeleteProductGroupCommand{
		ActorID: actorID,
		ID:      groupID,
	}

	err = h.service.DeleteProductGroup(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
