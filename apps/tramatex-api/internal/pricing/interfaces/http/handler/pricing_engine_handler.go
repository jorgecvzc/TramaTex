package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/pricing/application"
)

type PricingEngineHandler struct {
	service *application.PricingEngineService
}

func NewPricingEngineHandler(service *application.PricingEngineService) *PricingEngineHandler {
	return &PricingEngineHandler{service: service}
}

func (h *PricingEngineHandler) CreateBaseSalesPriceRule(c *gin.Context) {
	var cmd application.CreateBaseSalesPriceRuleCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	rule, err := h.service.CreateBaseSalesPriceRule(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

func (h *PricingEngineHandler) UpdateBaseSalesPriceRule(c *gin.Context) {
	ruleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	var cmd application.UpdateBaseSalesPriceRuleCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cmd.ID = ruleID

	rule, err := h.service.UpdateBaseSalesPriceRule(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rule)
}

func (h *PricingEngineHandler) CreateSaleModificationRule(c *gin.Context) {
	var cmd application.CreateSaleModificationRuleCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	rule, err := h.service.CreateSaleModificationRule(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

func (h *PricingEngineHandler) UpdateSaleModificationRule(c *gin.Context) {
	ruleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	var cmd application.UpdateSaleModificationRuleCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cmd.ID = ruleID

	rule, err := h.service.UpdateSaleModificationRule(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rule)
}

func (h *PricingEngineHandler) CalculateBaseSalesPrice(c *gin.Context) {
	var req application.CalculateBaseSalesPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.VariantID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "variantId is required"})
		return
	}

	result, err := h.service.CalculateBaseSalesPrice(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *PricingEngineHandler) CalculateFinalSalePrice(c *gin.Context) {
	var req application.CalculateFinalSalePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.ClientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clientId is required"})
		return
	}

	result, err := h.service.CalculateFinalSalePrice(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
