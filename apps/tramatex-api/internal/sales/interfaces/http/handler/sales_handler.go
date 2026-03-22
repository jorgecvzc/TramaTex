package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/joran-cortez/tramatex/internal/sales/application"
)

type SalesHandler struct {
	service *application.SalesService
}

func NewSalesHandler(service *application.SalesService) *SalesHandler {
	return &SalesHandler{service: service}
}

func (h *SalesHandler) CreateQuote(c *gin.Context) {
	var cmd application.CreateQuoteCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.service.CreateQuote(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *SalesHandler) GetQuote(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	result, err := h.service.GetQuote(c.Request.Context(), application.GetQuoteByIDQuery{ID: id})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) ListQuotes(c *gin.Context) {
	query := application.ListQuotesQuery{}
	if partyID, ok := parseUUIDQuery(c, "partyId"); ok {
		query.PartyID = partyID
	}
	if search := strings.TrimSpace(c.Query("search")); search != "" {
		query.Search = &search
	}
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query.Status = &status
	}
	if fromDate, ok := parseTimeQuery(c, "fromDate"); ok {
		query.FromDate = fromDate
	} else {
		return
	}
	if toDate, ok := parseEndDateQuery(c, "toDate"); ok {
		query.ToDate = toDate
	} else {
		return
	}
	if limit := strings.TrimSpace(c.Query("limit")); limit != "" {
		if n, err := strconv.Atoi(limit); err == nil && n > 0 {
			query.PageSize = n
		}
	}

	result, err := h.service.ListQuotes(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) UpdateQuote(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var cmd application.UpdateQuoteCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cmd.QuoteID = id

	result, err := h.service.UpdateQuote(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) ChangeQuoteStatus(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var cmd application.ChangeQuoteStatusCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cmd.QuoteID = id

	result, err := h.service.ChangeQuoteStatus(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) DeleteQuote(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	cmd := application.DeleteQuoteCommand{QuoteID: id}
	if err := h.service.DeleteQuote(c.Request.Context(), cmd); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SalesHandler) PreviewQuoteCalculation(c *gin.Context) {
	var cmd application.PreviewQuoteCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.service.PreviewQuoteCalculation(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) PreviewOrderCalculation(c *gin.Context) {
	var cmd application.PreviewOrderCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.service.PreviewOrderCalculation(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) ConvertQuoteToOrder(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var cmd application.ConvertQuoteToOrderCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cmd.QuoteID = id

	result, err := h.service.ConvertQuoteToOrder(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *SalesHandler) AcceptAndConvertQuote(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var cmd application.AcceptAndConvertQuoteCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cmd.QuoteID = id

	result, err := h.service.AcceptAndConvertQuote(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *SalesHandler) CreateOrder(c *gin.Context) {
	var cmd application.CreateOrderCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.service.CreateOrder(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *SalesHandler) GetOrder(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	result, err := h.service.GetOrder(c.Request.Context(), application.GetOrderByIDQuery{ID: id})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) ListOrders(c *gin.Context) {
	query := application.ListOrdersQuery{}
	if partyID, ok := parseUUIDQuery(c, "partyId"); ok {
		query.PartyID = partyID
	}
	if search := strings.TrimSpace(c.Query("search")); search != "" {
		query.Search = &search
	}
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query.Status = &status
	}
	if fromDate, ok := parseTimeQuery(c, "fromDate"); ok {
		query.FromDate = fromDate
	} else {
		return
	}
	if toDate, ok := parseEndDateQuery(c, "toDate"); ok {
		query.ToDate = toDate
	} else {
		return
	}
	if limit := strings.TrimSpace(c.Query("limit")); limit != "" {
		if n, err := strconv.Atoi(limit); err == nil && n > 0 {
			query.PageSize = n
		}
	}

	result, err := h.service.ListOrders(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) ListPendingWorkSetups(c *gin.Context) {
	result, err := h.service.ListPendingWorkSetups(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) UpdateOrderDetails(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var cmd application.UpdateOrderDetailsCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cmd.OrderID = id

	result, err := h.service.UpdateOrderDetails(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) ChangeOrderStatus(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var cmd application.ChangeOrderStatusCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cmd.OrderID = id

	result, err := h.service.ChangeOrderStatus(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) AddOrderLineItem(c *gin.Context) {
	orderID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var cmd application.AddOrderLineItemCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cmd.OrderID = orderID

	result, err := h.service.AddOrderLineItem(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) UpdateOrderLineItem(c *gin.Context) {
	orderID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}
	lineItemID, ok := parseUUIDParam(c, "lineItemId")
	if !ok {
		return
	}

	var cmd application.UpdateOrderLineItemCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cmd.OrderID = orderID
	cmd.LineItemID = lineItemID

	result, err := h.service.UpdateOrderLineItem(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) RemoveOrderLineItem(c *gin.Context) {
	orderID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}
	lineItemID, ok := parseUUIDParam(c, "lineItemId")
	if !ok {
		return
	}

	cmd := application.RemoveOrderLineItemCommand{
		OrderID:    orderID,
		LineItemID: lineItemID,
	}

	result, err := h.service.RemoveOrderLineItem(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) CreateDeliveryNote(c *gin.Context) {
	var cmd application.CreateDeliveryNoteCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.service.CreateDeliveryNote(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *SalesHandler) GetDeliveryNote(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	result, err := h.service.GetDeliveryNote(c.Request.Context(), application.GetDeliveryNoteByIDQuery{ID: id})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) ListDeliveryNotes(c *gin.Context) {
	query := application.ListDeliveryNotesQuery{}
	if orderID, ok := parseUUIDQuery(c, "salesOrderId"); ok {
		query.SalesOrderID = orderID
	}
	if partyID, ok := parseUUIDQuery(c, "partyId"); ok {
		query.PartyID = partyID
	}
	if search := strings.TrimSpace(c.Query("search")); search != "" {
		query.Search = &search
	}
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query.Status = &status
	}
	if fromDate, ok := parseTimeQuery(c, "fromDate"); ok {
		query.FromDate = fromDate
	} else {
		return
	}
	if toDate, ok := parseEndDateQuery(c, "toDate"); ok {
		query.ToDate = toDate
	} else {
		return
	}
	if limit := strings.TrimSpace(c.Query("limit")); limit != "" {
		if n, err := strconv.Atoi(limit); err == nil && n > 0 {
			query.PageSize = n
		}
	}

	result, err := h.service.ListDeliveryNotes(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) CreateInvoice(c *gin.Context) {
	var cmd application.CreateInvoiceCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.service.CreateInvoice(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

// CreateSimplifiedInvoice creates a ticket (factura simplificada) for retail sales < 3,000 EUR
// POST /api/sales/invoices/simplified
func (h *SalesHandler) CreateSimplifiedInvoice(c *gin.Context) {
	var cmd application.CreateSimplifiedInvoiceCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.service.CreateSimplifiedInvoice(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *SalesHandler) ChangeDeliveryNoteStatus(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var cmd application.ChangeDeliveryNoteStatusCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cmd.DeliveryNoteID = id

	result, err := h.service.ChangeDeliveryNoteStatus(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) GetInvoice(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	result, err := h.service.GetInvoice(c.Request.Context(), application.GetInvoiceByIDQuery{ID: id})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) ChangeInvoiceStatus(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var cmd application.ChangeInvoiceStatusCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cmd.InvoiceID = id

	result, err := h.service.ChangeInvoiceStatus(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SalesHandler) ListInvoices(c *gin.Context) {
	query := application.ListInvoicesQuery{}
	if partyID, ok := parseUUIDQuery(c, "partyId"); ok {
		query.PartyID = partyID
	}
	if orderID, ok := parseUUIDQuery(c, "orderId"); ok {
		query.SalesOrderID = orderID
	}
	if deliveryNoteID, ok := parseUUIDQuery(c, "deliveryNoteId"); ok {
		query.DeliveryNoteID = deliveryNoteID
	}
	if search := strings.TrimSpace(c.Query("search")); search != "" {
		query.Search = &search
	}
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query.Status = &status
	}
	if invoiceType := strings.TrimSpace(c.Query("type")); invoiceType != "" {
		query.Type = &invoiceType
	}
	if fromDate, ok := parseTimeQuery(c, "fromDate"); ok {
		query.FromDate = fromDate
	} else {
		return
	}
	if toDate, ok := parseEndDateQuery(c, "toDate"); ok {
		query.ToDate = toDate
	} else {
		return
	}
	if limit := strings.TrimSpace(c.Query("limit")); limit != "" {
		if n, err := strconv.Atoi(limit); err == nil && n > 0 {
			query.PageSize = n
		}
	}

	result, err := h.service.ListInvoices(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func parseUUIDParam(c *gin.Context, name string) (uuid.UUID, bool) {
	value, err := uuid.Parse(c.Param(name))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return uuid.Nil, false
	}
	return value, true
}

func parseUUIDQuery(c *gin.Context, name string) (*uuid.UUID, bool) {
	value := strings.TrimSpace(c.Query(name))
	if value == "" {
		return nil, true
	}
	parsed, err := uuid.Parse(value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return nil, false
	}
	return &parsed, true
}

func parseTimeQuery(c *gin.Context, name string) (*time.Time, bool) {
	value := strings.TrimSpace(c.Query(name))
	if value == "" {
		return nil, true
	}
	// Try RFC3339 format first (2006-01-02T15:04:05Z)
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		// Try date-only format (2006-01-02)
		parsed, err = time.Parse("2006-01-02", value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name + ". Expected format: YYYY-MM-DD or RFC3339"})
			return nil, false
		}
	}
	return &parsed, true
}

// parseEndDateQuery parses a "toDate" query parameter. When the value is a
// date-only string (YYYY-MM-DD) it is adjusted to the end of that day
// (23:59:59.999999999 UTC) so that records created at any point during that
// day are included in the results.
func parseEndDateQuery(c *gin.Context, name string) (*time.Time, bool) {
	value := strings.TrimSpace(c.Query(name))
	if value == "" {
		return nil, true
	}
	// Try RFC3339 first — caller provided a precise timestamp, use as-is.
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		// Date-only format: extend to end of day.
		parsed, err = time.Parse("2006-01-02", value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name + ". Expected format: YYYY-MM-DD or RFC3339"})
			return nil, false
		}
		endOfDay := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 999999999, time.UTC)
		return &endOfDay, true
	}
	return &parsed, true
}
