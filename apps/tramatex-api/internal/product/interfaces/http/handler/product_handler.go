package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/joran-cortez/tramatex/internal/product/application"
)

type ProductHandler struct {
	service *application.ProductService
}

func NewProductHandler(service *application.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var cmd application.CreateProductCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.CreateProduct(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) AddGroupToProduct(c *gin.Context) {
	var cmd application.AddGroupCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.AddGroupToProduct(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) AddDirectAttributeToProduct(c *gin.Context) {
	var cmd application.AddDirectAttributeCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	cmd.ProductID = productID

	product, err := h.service.UpdateProductSKU(c.Request.Context(), cmd)
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

	attribute, err := h.service.CreateAttribute(c.Request.Context(), cmd)
	if err != nil {
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
	var query application.ListAttributesQuery
	// TODO: Parse optional query parameters for filtering (ScopeType, BrandID, ProductGroupID)
	// For example:
	// if brandIDStr := c.Query("brandId"); brandIDStr != "" {
	// 	brandID, err := uuid.Parse(brandIDStr)
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand id"})
	// 		return
	// 	}
	// 	query.BrandID = &brandID
	// }

	attributes, err := h.service.ListAttributes(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attributes)
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
	cmd.ID = attributeID // Set ID from URL parameter

	attribute, err := h.service.UpdateAttribute(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attribute)
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
	// TODO: Parse optional query parameters for filtering

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
	productID, err := uuid.Parse(c.Param("productId"))
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
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var cmd application.FindOrCreateProductVariantCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ProductID = productID

	variant, err := h.service.FindOrCreateProductVariant(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, variant)
}

func (h *ProductHandler) ListProductVariantsByProductID(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
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
	cmd.ID = variantID

	variant, err := h.service.UpdateProductVariant(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, variant)
}

func (h *ProductHandler) CreatePartyServiceConfiguration(c *gin.Context) {
	partyID, err := uuid.Parse(c.Param("partyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid party id"})
		return
	}

	var cmd application.CreatePartyServiceConfigurationCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.PartyID = partyID

	config, err := h.service.CreatePartyServiceConfiguration(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, config)
}

func (h *ProductHandler) ListPartyServiceConfigurationsByPartyID(c *gin.Context) {
	partyID, err := uuid.Parse(c.Param("partyId"))
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
	partyID, err := uuid.Parse(c.Param("partyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid party id"})
		return
	}
	configID, err := uuid.Parse(c.Param("id"))
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
	partyID, err := uuid.Parse(c.Param("partyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid party id"})
		return
	}
	configID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid configuration id"})
		return
	}

	var cmd application.UpdatePartyServiceConfigurationCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ID = configID
	cmd.PartyID = partyID

	config, err := h.service.UpdatePartyServiceConfiguration(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

func (h *ProductHandler) DeletePartyServiceConfiguration(c *gin.Context) {
	partyID, err := uuid.Parse(c.Param("partyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid party id"})
		return
	}
	configID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid configuration id"})
		return
	}

	cmd := application.DeletePartyServiceConfigurationCommand{ID: configID, PartyID: partyID}
	err = h.service.DeletePartyServiceConfiguration(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
