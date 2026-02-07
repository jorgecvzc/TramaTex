package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/pricing/application"
)

type PricingHandler struct {
	service *application.PricingService
}

func NewPricingHandler(service *application.PricingService) *PricingHandler {
	return &PricingHandler{service: service}
}

func (h *PricingHandler) CalculatePrice(c *gin.Context) {
	var cmd application.CalculatePriceCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if cmd.ProductVariantID == uuid.Nil || cmd.ClientID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_variant_id and client_id are required"})
		return
	}

	result, err := h.service.CalculatePrice(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *PricingHandler) ListPricingRules(c *gin.Context) {
	rules, err := h.service.ListPricingRules(c.Request.Context(), application.ListPricingRulesQuery{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rules)
}

func (h *PricingHandler) CreatePricingRule(c *gin.Context) {
	var cmd application.CreatePricingRuleCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	rule, err := h.service.CreatePricingRule(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

func (h *PricingHandler) CreateClientPricingOverride(c *gin.Context) {
	var cmd application.CreateClientPricingCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	override, err := h.service.CreateClientPricing(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, override)
}

func (h *PricingHandler) GetPricingHistory(c *gin.Context) {
	variantID, err := uuid.Parse(c.Param("variantId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid variant id"})
		return
	}

	history, err := h.service.GetPricingHistory(c.Request.Context(), application.GetPricingHistoryQuery{ProductVariantID: variantID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
