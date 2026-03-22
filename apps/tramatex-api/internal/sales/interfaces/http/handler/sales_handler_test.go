package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	pricing_app "github.com/joran-cortez/tramatex/internal/pricing/application"
	"github.com/joran-cortez/tramatex/internal/sales/application"
	"github.com/joran-cortez/tramatex/internal/sales/domain"
	infra_middleware "github.com/joran-cortez/tramatex/internal/shared/infrastructure/middleware"
)

// ===== STUB REPOSITORIES =====

type stubQuoteRepo struct {
	saveFn     func(context.Context, *domain.Quote) error
	findByIDFn func(context.Context, uuid.UUID) (*domain.Quote, error)
	listFn     func(context.Context, domain.QuoteFilter) ([]*domain.Quote, error)
	deleteFn   func(context.Context, uuid.UUID) error
}

func (s *stubQuoteRepo) Save(ctx context.Context, quote *domain.Quote) error {
	if s.saveFn != nil {
		return s.saveFn(ctx, quote)
	}
	return nil
}

func (s *stubQuoteRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Quote, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, domain.NewNotFoundError("quote not found")
}

func (s *stubQuoteRepo) List(ctx context.Context, filter domain.QuoteFilter) ([]*domain.Quote, error) {
	if s.listFn != nil {
		return s.listFn(ctx, filter)
	}
	return []*domain.Quote{}, nil
}

func (s *stubQuoteRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	return nil
}

type stubOrderRepo struct {
	saveFn     func(context.Context, *domain.SalesOrder) error
	findByIDFn func(context.Context, uuid.UUID) (*domain.SalesOrder, error)
	listFn     func(context.Context, domain.SalesOrderFilter) ([]*domain.SalesOrder, error)
}

func (s *stubOrderRepo) Save(ctx context.Context, order *domain.SalesOrder) error {
	if s.saveFn != nil {
		return s.saveFn(ctx, order)
	}
	return nil
}

func (s *stubOrderRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, domain.NewNotFoundError("order not found")
}

func (s *stubOrderRepo) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	return s.FindByID(ctx, id)
}

func (s *stubOrderRepo) List(ctx context.Context, filter domain.SalesOrderFilter) ([]*domain.SalesOrder, error) {
	if s.listFn != nil {
		return s.listFn(ctx, filter)
	}
	return []*domain.SalesOrder{}, nil
}

func (s *stubOrderRepo) FindByQuoteID(ctx context.Context, quoteID uuid.UUID) (*domain.SalesOrder, error) {
	return nil, nil
}

type stubDeliveryNoteRepo struct {
	saveFn               func(context.Context, *domain.DeliveryNote) error
	findByIDFn           func(context.Context, uuid.UUID) (*domain.DeliveryNote, error)
	listFn               func(context.Context, domain.DeliveryNoteFilter) ([]*domain.DeliveryNote, error)
	listBySalesOrderIDFn func(context.Context, uuid.UUID) ([]*domain.DeliveryNote, error)
	linkLineItemsFn      func(context.Context, map[uuid.UUID]uuid.UUID) error
}

func (s *stubDeliveryNoteRepo) Save(ctx context.Context, note *domain.DeliveryNote) error {
	if s.saveFn != nil {
		return s.saveFn(ctx, note)
	}
	return nil
}

func (s *stubDeliveryNoteRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.DeliveryNote, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, domain.NewNotFoundError("delivery note not found")
}

func (s *stubDeliveryNoteRepo) List(ctx context.Context, filter domain.DeliveryNoteFilter) ([]*domain.DeliveryNote, error) {
	if s.listFn != nil {
		return s.listFn(ctx, filter)
	}
	return []*domain.DeliveryNote{}, nil
}

func (s *stubDeliveryNoteRepo) ListBySalesOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.DeliveryNote, error) {
	if s.listBySalesOrderIDFn != nil {
		return s.listBySalesOrderIDFn(ctx, orderID)
	}
	return []*domain.DeliveryNote{}, nil
}

func (s *stubDeliveryNoteRepo) LinkLineItemsToInvoice(ctx context.Context, links map[uuid.UUID]uuid.UUID) error {
	if s.linkLineItemsFn != nil {
		return s.linkLineItemsFn(ctx, links)
	}
	return nil
}

type stubInvoiceRepo struct {
	saveFn                 func(context.Context, *domain.Invoice) error
	findByIDFn             func(context.Context, uuid.UUID) (*domain.Invoice, error)
	listFn                 func(context.Context, domain.InvoiceFilter) ([]*domain.Invoice, error)
	listBySalesOrderIDFn   func(context.Context, uuid.UUID) ([]*domain.Invoice, error)
	findByDeliveryNoteIDFn func(context.Context, uuid.UUID) (*domain.Invoice, error)
	listDeliveryNoteIDsFn  func(context.Context, uuid.UUID) ([]uuid.UUID, error)
}

func (s *stubInvoiceRepo) Save(ctx context.Context, invoice *domain.Invoice) error {
	if s.saveFn != nil {
		return s.saveFn(ctx, invoice)
	}
	return nil
}

func (s *stubInvoiceRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, domain.NewNotFoundError("invoice not found")
}

func (s *stubInvoiceRepo) List(ctx context.Context, filter domain.InvoiceFilter) ([]*domain.Invoice, error) {
	if s.listFn != nil {
		return s.listFn(ctx, filter)
	}
	return []*domain.Invoice{}, nil
}

func (s *stubInvoiceRepo) ListBySalesOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.Invoice, error) {
	if s.listBySalesOrderIDFn != nil {
		return s.listBySalesOrderIDFn(ctx, orderID)
	}
	return []*domain.Invoice{}, nil
}

func (s *stubInvoiceRepo) FindByDeliveryNoteID(ctx context.Context, deliveryNoteID uuid.UUID) (*domain.Invoice, error) {
	if s.findByDeliveryNoteIDFn != nil {
		return s.findByDeliveryNoteIDFn(ctx, deliveryNoteID)
	}
	return nil, nil
}

func (s *stubInvoiceRepo) ListDeliveryNoteIDsByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]uuid.UUID, error) {
	if s.listDeliveryNoteIDsFn != nil {
		return s.listDeliveryNoteIDsFn(ctx, invoiceID)
	}
	return nil, nil
}

func (s *stubInvoiceRepo) ListOrderIDsByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}

type stubDocumentNumberGenerator struct {
	nextQuoteNumberFn        func(context.Context) (domain.QuoteNumber, error)
	nextOrderNumberFn        func(context.Context) (domain.OrderNumber, error)
	nextDeliveryNoteNumberFn func(context.Context) (domain.DeliveryNoteNumber, error)
	nextInvoiceNumberFn      func(context.Context, domain.InvoiceSeries) (domain.InvoiceNumber, error)
}

func (s *stubDocumentNumberGenerator) NextQuoteNumber(ctx context.Context) (domain.QuoteNumber, error) {
	if s.nextQuoteNumberFn != nil {
		return s.nextQuoteNumberFn(ctx)
	}
	qn, _ := domain.NewQuoteNumber("PRE-2026-0001")
	return qn, nil
}

func (s *stubDocumentNumberGenerator) NextOrderNumber(ctx context.Context) (domain.OrderNumber, error) {
	if s.nextOrderNumberFn != nil {
		return s.nextOrderNumberFn(ctx)
	}
	on, _ := domain.NewOrderNumber("PED-2026-0001")
	return on, nil
}

func (s *stubDocumentNumberGenerator) NextDeliveryNoteNumber(ctx context.Context) (domain.DeliveryNoteNumber, error) {
	if s.nextDeliveryNoteNumberFn != nil {
		return s.nextDeliveryNoteNumberFn(ctx)
	}
	dn, _ := domain.NewDeliveryNoteNumber("ALB-2026-0001")
	return dn, nil
}

func (s *stubDocumentNumberGenerator) NextInvoiceNumber(ctx context.Context, series domain.InvoiceSeries) (domain.InvoiceNumber, error) {
	if s.nextInvoiceNumberFn != nil {
		return s.nextInvoiceNumberFn(ctx, series)
	}
	in, _ := domain.NewInvoiceNumber("FV-2026-0001")
	return in, nil
}

type stubPricingEngine struct {
	calculateFinalSalePriceFn func(context.Context, pricing_app.CalculateFinalSalePriceRequest) (*pricing_app.CalculateFinalSalePriceResponse, error)
}

func (s *stubPricingEngine) CalculateFinalSalePrice(ctx context.Context, req pricing_app.CalculateFinalSalePriceRequest) (*pricing_app.CalculateFinalSalePriceResponse, error) {
	if s.calculateFinalSalePriceFn != nil {
		return s.calculateFinalSalePriceFn(ctx, req)
	}
	// Return a default response
	return &pricing_app.CalculateFinalSalePriceResponse{
		CalculatedItems: []pricing_app.CalculatedSaleItemResponse{},
		SaleTotal: pricing_app.MoneyDTO{
			Amount:   100.0,
			Currency: "EUR",
		},
	}, nil
}

type stubPartyLookup struct {
	existsPartyFn  func(context.Context, uuid.UUID) (bool, error)
	hasPartyRoleFn func(context.Context, uuid.UUID, string) (bool, error)
}

func (s *stubPartyLookup) ExistsParty(ctx context.Context, id uuid.UUID) (bool, error) {
	if s.existsPartyFn != nil {
		return s.existsPartyFn(ctx, id)
	}
	return true, nil
}

func (s *stubPartyLookup) HasPartyRole(ctx context.Context, id uuid.UUID, role string) (bool, error) {
	if s.hasPartyRoleFn != nil {
		return s.hasPartyRoleFn(ctx, id, role)
	}
	return true, nil
}

// ===== HELPER FUNCTIONS =====

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(infra_middleware.ErrorHandlerMiddleware("development"))
	return r
}

// ===== TESTS: QUOTE HANDLERS =====

func TestCreateQuote_Success(t *testing.T) {
	quoteRepo := &stubQuoteRepo{
		saveFn: func(ctx context.Context, quote *domain.Quote) error {
			return nil
		},
	}
	partyLookup := &stubPartyLookup{
		existsPartyFn: func(ctx context.Context, id uuid.UUID) (bool, error) {
			return true, nil
		},
	}
	numberGen := &stubDocumentNumberGenerator{}

	variantID := uuid.New()
	pricingEngine := &stubPricingEngine{
		calculateFinalSalePriceFn: func(ctx context.Context, req pricing_app.CalculateFinalSalePriceRequest) (*pricing_app.CalculateFinalSalePriceResponse, error) {
			return &pricing_app.CalculateFinalSalePriceResponse{
				CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
					{
						ProductVariantID: variantID,
						Quantity:         10,
						BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 10.0, Currency: "EUR"},
						FinalPrice:       pricing_app.MoneyDTO{Amount: 10.0, Currency: "EUR"},
					},
				},
				SaleTotal: pricing_app.MoneyDTO{Amount: 100.0, Currency: "EUR"},
			}, nil
		},
	}

	service := application.NewSalesService(
		quoteRepo, &stubOrderRepo{}, &stubDeliveryNoteRepo{}, &stubInvoiceRepo{},
		numberGen, pricingEngine, partyLookup, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.POST("/quotes", handler.CreateQuote)

	partyID := uuid.New()
	reqBody := map[string]interface{}{
		"partyId":        partyID.String(),
		"expirationDate": time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		"items": []map[string]interface{}{
			{
				"productVariantId": variantID.String(),
				"quantity":         10,
			},
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Logf("❌ Expected 201, got %d. Response: %s", rec.Code, rec.Body.String())
	}
	assert.Equal(t, http.StatusCreated, rec.Code)

	var result map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.NotEmpty(t, result["id"])
	assert.NotEmpty(t, result["quoteNumber"])
}

func TestGetQuote_NotFound(t *testing.T) {
	quoteRepo := &stubQuoteRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Quote, error) {
			return nil, domain.NewNotFoundError("quote not found")
		},
	}

	service := application.NewSalesService(
		quoteRepo, &stubOrderRepo{}, &stubDeliveryNoteRepo{}, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.GET("/quotes/:id", handler.GetQuote)

	quoteID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/quotes/"+quoteID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetQuote_InvalidID(t *testing.T) {
	service := application.NewSalesService(
		&stubQuoteRepo{}, &stubOrderRepo{}, &stubDeliveryNoteRepo{}, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.GET("/quotes/:id", handler.GetQuote)

	req := httptest.NewRequest(http.MethodGet, "/quotes/invalid-uuid", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ===== TESTS: ORDER HANDLERS =====

func TestCreateOrder_Success(t *testing.T) {
	orderRepo := &stubOrderRepo{
		saveFn: func(ctx context.Context, order *domain.SalesOrder) error {
			return nil
		},
	}
	partyLookup := &stubPartyLookup{
		existsPartyFn: func(ctx context.Context, id uuid.UUID) (bool, error) {
			return true, nil
		},
	}
	numberGen := &stubDocumentNumberGenerator{}

	variantID := uuid.New()
	pricingEngine := &stubPricingEngine{
		calculateFinalSalePriceFn: func(ctx context.Context, req pricing_app.CalculateFinalSalePriceRequest) (*pricing_app.CalculateFinalSalePriceResponse, error) {
			return &pricing_app.CalculateFinalSalePriceResponse{
				CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
					{
						ProductVariantID: variantID,
						Quantity:         5,
						BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 10.0, Currency: "EUR"},
						FinalPrice:       pricing_app.MoneyDTO{Amount: 10.0, Currency: "EUR"},
					},
				},
				SaleTotal: pricing_app.MoneyDTO{Amount: 50.0, Currency: "EUR"},
			}, nil
		},
	}

	service := application.NewSalesService(
		&stubQuoteRepo{}, orderRepo, &stubDeliveryNoteRepo{}, &stubInvoiceRepo{},
		numberGen, pricingEngine, partyLookup, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.POST("/orders", handler.CreateOrder)

	partyID := uuid.New()
	reqBody := map[string]interface{}{
		"partyId":      partyID.String(),
		"deliveryDate": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
		"items": []map[string]interface{}{
			{
				"productVariantId": variantID.String(),
				"quantity":         5,
			},
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Logf("❌ Expected 201, got %d. Response: %s", rec.Code, rec.Body.String())
	}
	assert.Equal(t, http.StatusCreated, rec.Code)

	var result map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.NotEmpty(t, result["id"])
	assert.NotEmpty(t, result["orderNumber"])
}

func TestGetOrder_NotFound(t *testing.T) {
	orderRepo := &stubOrderRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
			return nil, domain.NewNotFoundError("order not found")
		},
	}

	service := application.NewSalesService(
		&stubQuoteRepo{}, orderRepo, &stubDeliveryNoteRepo{}, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.GET("/orders/:id", handler.GetOrder)

	orderID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ===== TESTS: DELIVERY NOTE HANDLERS =====

func TestCreateDeliveryNote_Success(t *testing.T) {
	orderID := uuid.New()
	lineItemID := uuid.New()
	variantID := uuid.New()

	// Create a valid order with line items using struct literals
	money, _ := domain.NewMoney(100.0, "EUR")
	order := &domain.SalesOrder{
		ID:           orderID,
		PartyID:      uuid.New(),
		OrderDate:    time.Now(),
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Status:       domain.SalesOrderStatusPending,
		LineItems: []domain.OrderLineItem{
			{
				ID:               lineItemID,
				ProductVariantID: variantID,
				Quantity:         10,
				ListUnitPrice:    money,
				UnitPrice:        money,
				DiscountPerUnit:  money,
				Subtotal:         money,
			},
		},
		Subtotal:  money,
		TaxAmount: money,
		Total:     money,
	}

	orderRepo := &stubOrderRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
			if id == orderID {
				return order, nil
			}
			return nil, domain.NewNotFoundError("order not found")
		},
	}

	noteRepo := &stubDeliveryNoteRepo{
		saveFn: func(ctx context.Context, note *domain.DeliveryNote) error {
			return nil
		},
	}

	service := application.NewSalesService(
		&stubQuoteRepo{}, orderRepo, noteRepo, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.POST("/delivery-notes", handler.CreateDeliveryNote)

	reqBody := map[string]interface{}{
		"salesOrderId": orderID.String(),
		"deliveryDate": time.Now().Format(time.RFC3339),
		"items": []map[string]interface{}{
			{
				"salesOrderLineItemId": lineItemID.String(),
				"deliveredQuantity":    5,
			},
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/delivery-notes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Logf("❌ Expected 201, got %d. Response: %s", rec.Code, rec.Body.String())
	}
	assert.Equal(t, http.StatusCreated, rec.Code)

	var result map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.NotEmpty(t, result["id"])
}

// ===== TESTS: INVOICE HANDLERS =====

func TestCreateInvoice_Success(t *testing.T) {
	orderID := uuid.New()
	partyID := uuid.New()
	variantID := uuid.New()

	// Create a valid order with line items using struct literals
	money, _ := domain.NewMoney(100.0, "EUR")
	order := &domain.SalesOrder{
		ID:           orderID,
		PartyID:      partyID,
		OrderDate:    time.Now(),
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Status:       domain.SalesOrderStatusDelivered, // Status must be Delivered for invoicing
		LineItems: []domain.OrderLineItem{
			{
				ID:               uuid.New(),
				ProductVariantID: variantID,
				Quantity:         10,
				ListUnitPrice:    money,
				UnitPrice:        money,
				DiscountPerUnit:  money,
				Subtotal:         money,
			},
		},
		Subtotal:  money,
		TaxAmount: money,
		Total:     money,
	}

	orderRepo := &stubOrderRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
			if id == orderID {
				return order, nil
			}
			return nil, domain.NewNotFoundError("order not found")
		},
	}

	invoiceRepo := &stubInvoiceRepo{
		saveFn: func(ctx context.Context, invoice *domain.Invoice) error {
			return nil
		},
	}
	numberGen := &stubDocumentNumberGenerator{}

	service := application.NewSalesService(
		&stubQuoteRepo{}, orderRepo, &stubDeliveryNoteRepo{}, invoiceRepo,
		numberGen, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.POST("/invoices", handler.CreateInvoice)

	paymentTerms := "Net 30"
	reqBody := map[string]interface{}{
		"partyId":         partyID.String(),
		"salesOrderIds":   []string{orderID.String()},
		"deliveryNoteIds": []string{},
		"invoiceDate":     time.Now().Format(time.RFC3339),
		"dueDate":         time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		"paymentTerms":    &paymentTerms,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/invoices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Logf("❌ Expected 201, got %d. Response: %s", rec.Code, rec.Body.String())
	}
	assert.Equal(t, http.StatusCreated, rec.Code)

	var result map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.NotEmpty(t, result["id"])
	assert.NotEmpty(t, result["invoiceNumber"])
}

func TestGetInvoice_NotFound(t *testing.T) {
	invoiceRepo := &stubInvoiceRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
			return nil, domain.NewNotFoundError("invoice not found")
		},
	}

	service := application.NewSalesService(
		&stubQuoteRepo{}, &stubOrderRepo{}, &stubDeliveryNoteRepo{}, invoiceRepo,
		&stubDocumentNumberGenerator{}, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.GET("/invoices/:id", handler.GetInvoice)

	invoiceID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/invoices/"+invoiceID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ===== ADDITIONAL QUOTE HANDLER TESTS =====

func TestListQuotes_Success(t *testing.T) {
	partyID := uuid.New()
	quoteID1 := uuid.New()
	quoteID2 := uuid.New()

	money, _ := domain.NewMoney(100.0, "EUR")
	quote1 := &domain.Quote{
		ID:             quoteID1,
		PartyID:        partyID,
		QuoteDate:      time.Now(),
		ExpirationDate: time.Now().Add(30 * 24 * time.Hour),
		Status:         domain.QuoteStatusDraft,
		Subtotal:       money,
		TaxAmount:      money,
		Total:          money,
	}
	quote2 := &domain.Quote{
		ID:             quoteID2,
		PartyID:        partyID,
		QuoteDate:      time.Now(),
		ExpirationDate: time.Now().Add(30 * 24 * time.Hour),
		Status:         domain.QuoteStatusIssued,
		Subtotal:       money,
		TaxAmount:      money,
		Total:          money,
	}

	quoteRepo := &stubQuoteRepo{
		listFn: func(ctx context.Context, filter domain.QuoteFilter) ([]*domain.Quote, error) {
			return []*domain.Quote{quote1, quote2}, nil
		},
	}

	service := application.NewSalesService(
		quoteRepo, &stubOrderRepo{}, &stubDeliveryNoteRepo{}, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.GET("/quotes", handler.ListQuotes)

	req := httptest.NewRequest(http.MethodGet, "/quotes?partyId="+partyID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUpdateQuote_Success(t *testing.T) {
	quoteID := uuid.New()
	money, _ := domain.NewMoney(100.0, "EUR")

	existingQuote := &domain.Quote{
		ID:             quoteID,
		PartyID:        uuid.New(),
		QuoteDate:      time.Now(),
		ExpirationDate: time.Now().Add(30 * 24 * time.Hour),
		Status:         domain.QuoteStatusDraft,
		Subtotal:       money,
		TaxAmount:      money,
		Total:          money,
	}

	quoteRepo := &stubQuoteRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Quote, error) {
			if id == quoteID {
				return existingQuote, nil
			}
			return nil, domain.NewNotFoundError("quote not found")
		},
		saveFn: func(ctx context.Context, quote *domain.Quote) error {
			return nil
		},
	}

	variantID := uuid.New()
	pricingEngine := &stubPricingEngine{
		calculateFinalSalePriceFn: func(ctx context.Context, req pricing_app.CalculateFinalSalePriceRequest) (*pricing_app.CalculateFinalSalePriceResponse, error) {
			return &pricing_app.CalculateFinalSalePriceResponse{
				CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
					{
						ProductVariantID: variantID,
						Quantity:         5,
						BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 20.0, Currency: "EUR"},
						FinalPrice:       pricing_app.MoneyDTO{Amount: 20.0, Currency: "EUR"},
					},
				},
				SaleTotal: pricing_app.MoneyDTO{Amount: 100.0, Currency: "EUR"},
			}, nil
		},
	}

	service := application.NewSalesService(
		quoteRepo, &stubOrderRepo{}, &stubDeliveryNoteRepo{}, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, pricingEngine, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.PUT("/quotes/:id", handler.UpdateQuote)

	reqBody := map[string]interface{}{
		"expirationDate": time.Now().Add(60 * 24 * time.Hour).Format(time.RFC3339),
		"items": []map[string]interface{}{
			{
				"productVariantId": variantID.String(),
				"quantity":         5,
			},
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/quotes/"+quoteID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Logf("❌ Expected 200, got %d. Response: %s", rec.Code, rec.Body.String())
	}
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestChangeQuoteStatus_Success(t *testing.T) {
	quoteID := uuid.New()
	money, _ := domain.NewMoney(100.0, "EUR")

	existingQuote := &domain.Quote{
		ID:             quoteID,
		PartyID:        uuid.New(),
		QuoteDate:      time.Now(),
		ExpirationDate: time.Now().Add(30 * 24 * time.Hour),
		Status:         domain.QuoteStatusDraft,
		Subtotal:       money,
		TaxAmount:      money,
		Total:          money,
	}

	quoteRepo := &stubQuoteRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Quote, error) {
			if id == quoteID {
				return existingQuote, nil
			}
			return nil, domain.NewNotFoundError("quote not found")
		},
		saveFn: func(ctx context.Context, quote *domain.Quote) error {
			return nil
		},
	}

	service := application.NewSalesService(
		quoteRepo, &stubOrderRepo{}, &stubDeliveryNoteRepo{}, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.POST("/quotes/:id/status", handler.ChangeQuoteStatus)

	reqBody := map[string]interface{}{
		"newStatus": "EMITIDA",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/quotes/"+quoteID.String()+"/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Logf("❌ Expected 200, got %d. Response: %s", rec.Code, rec.Body.String())
	}
	assert.Equal(t, http.StatusOK, rec.Code)
}

// ===== ADDITIONAL ORDER HANDLER TESTS =====

func TestListOrders_Success(t *testing.T) {
	partyID := uuid.New()
	orderID1 := uuid.New()
	orderID2 := uuid.New()

	money, _ := domain.NewMoney(100.0, "EUR")
	order1 := &domain.SalesOrder{
		ID:           orderID1,
		PartyID:      partyID,
		OrderDate:    time.Now(),
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Status:       domain.SalesOrderStatusPending,
		Subtotal:     money,
		TaxAmount:    money,
		Total:        money,
	}
	order2 := &domain.SalesOrder{
		ID:           orderID2,
		PartyID:      partyID,
		OrderDate:    time.Now(),
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Status:       domain.SalesOrderStatusDelivered,
		Subtotal:     money,
		TaxAmount:    money,
		Total:        money,
	}

	orderRepo := &stubOrderRepo{
		listFn: func(ctx context.Context, filter domain.SalesOrderFilter) ([]*domain.SalesOrder, error) {
			return []*domain.SalesOrder{order1, order2}, nil
		},
	}

	service := application.NewSalesService(
		&stubQuoteRepo{}, orderRepo, &stubDeliveryNoteRepo{}, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.GET("/orders", handler.ListOrders)

	req := httptest.NewRequest(http.MethodGet, "/orders?partyId="+partyID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUpdateOrderDetails_Success(t *testing.T) {
	orderID := uuid.New()
	money, _ := domain.NewMoney(100.0, "EUR")

	existingOrder := &domain.SalesOrder{
		ID:           orderID,
		PartyID:      uuid.New(),
		OrderDate:    time.Now(),
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Status:       domain.SalesOrderStatusPending,
		Subtotal:     money,
		TaxAmount:    money,
		Total:        money,
	}

	orderRepo := &stubOrderRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
			if id == orderID {
				return existingOrder, nil
			}
			return nil, domain.NewNotFoundError("order not found")
		},
		saveFn: func(ctx context.Context, order *domain.SalesOrder) error {
			return nil
		},
	}

	service := application.NewSalesService(
		&stubQuoteRepo{}, orderRepo, &stubDeliveryNoteRepo{}, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.PUT("/orders/:id", handler.UpdateOrderDetails)

	newDeliveryDate := time.Now().Add(14 * 24 * time.Hour)
	reqBody := map[string]interface{}{
		"deliveryDate": newDeliveryDate.Format(time.RFC3339),
		"notes":        "Updated notes",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/orders/"+orderID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Logf("❌ Expected 200, got %d. Response: %s", rec.Code, rec.Body.String())
	}
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestChangeOrderStatus_Success(t *testing.T) {
	orderID := uuid.New()
	money, _ := domain.NewMoney(100.0, "EUR")

	existingOrder := &domain.SalesOrder{
		ID:           orderID,
		PartyID:      uuid.New(),
		OrderDate:    time.Now(),
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Status:       domain.SalesOrderStatusPending,
		Subtotal:     money,
		TaxAmount:    money,
		Total:        money,
	}

	orderRepo := &stubOrderRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
			if id == orderID {
				return existingOrder, nil
			}
			return nil, domain.NewNotFoundError("order not found")
		},
		saveFn: func(ctx context.Context, order *domain.SalesOrder) error {
			return nil
		},
	}

	service := application.NewSalesService(
		&stubQuoteRepo{}, orderRepo, &stubDeliveryNoteRepo{}, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.POST("/orders/:id/status", handler.ChangeOrderStatus)

	reqBody := map[string]interface{}{
		"newStatus": "EN_PREPARACION",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/orders/"+orderID.String()+"/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Logf("❌ Expected 200, got %d. Response: %s", rec.Code, rec.Body.String())
	}
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAddOrderLineItem_Success(t *testing.T) {
	orderID := uuid.New()
	variantID := uuid.New()
	money, _ := domain.NewMoney(100.0, "EUR")

	existingOrder := &domain.SalesOrder{
		ID:           orderID,
		PartyID:      uuid.New(),
		OrderDate:    time.Now(),
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Status:       domain.SalesOrderStatusPending,
		LineItems:    []domain.OrderLineItem{},
		Subtotal:     money,
		TaxAmount:    money,
		Total:        money,
	}

	orderRepo := &stubOrderRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
			if id == orderID {
				return existingOrder, nil
			}
			return nil, domain.NewNotFoundError("order not found")
		},
		saveFn: func(ctx context.Context, order *domain.SalesOrder) error {
			return nil
		},
	}

	pricingEngine := &stubPricingEngine{
		calculateFinalSalePriceFn: func(ctx context.Context, req pricing_app.CalculateFinalSalePriceRequest) (*pricing_app.CalculateFinalSalePriceResponse, error) {
			return &pricing_app.CalculateFinalSalePriceResponse{
				CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
					{
						ProductVariantID: variantID,
						Quantity:         3,
						BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 15.0, Currency: "EUR"},
						FinalPrice:       pricing_app.MoneyDTO{Amount: 15.0, Currency: "EUR"},
					},
				},
				SaleTotal: pricing_app.MoneyDTO{Amount: 45.0, Currency: "EUR"},
			}, nil
		},
	}

	service := application.NewSalesService(
		&stubQuoteRepo{}, orderRepo, &stubDeliveryNoteRepo{}, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, pricingEngine, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.POST("/orders/:id/line-items", handler.AddOrderLineItem)

	reqBody := map[string]interface{}{
		"item": map[string]interface{}{
			"productVariantId": variantID.String(),
			"quantity":         3,
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/orders/"+orderID.String()+"/line-items", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Logf("❌ Expected 200, got %d. Response: %s", rec.Code, rec.Body.String())
	}
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateOrderLineItem_Success(t *testing.T) {
	orderID := uuid.New()
	lineItemID := uuid.New()
	variantID := uuid.New()
	money, _ := domain.NewMoney(100.0, "EUR")
	zero, _ := domain.NewMoney(0.0, "EUR")

	existingOrder := &domain.SalesOrder{
		ID:           orderID,
		PartyID:      uuid.New(),
		OrderDate:    time.Now(),
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Status:       domain.SalesOrderStatusPending,
		LineItems: []domain.OrderLineItem{
			{
				ID:               lineItemID,
				ProductVariantID: variantID,
				Quantity:         2,
				ListUnitPrice:    money,
				UnitPrice:        money,
				DiscountPerUnit:  zero,
				Subtotal:         money,
			},
		},
		Subtotal:  money,
		TaxAmount: money,
		Total:     money,
	}

	orderRepo := &stubOrderRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
			if id == orderID {
				return existingOrder, nil
			}
			return nil, domain.NewNotFoundError("order not found")
		},
		saveFn: func(ctx context.Context, order *domain.SalesOrder) error {
			return nil
		},
	}

	pricingEngine := &stubPricingEngine{
		calculateFinalSalePriceFn: func(ctx context.Context, req pricing_app.CalculateFinalSalePriceRequest) (*pricing_app.CalculateFinalSalePriceResponse, error) {
			return &pricing_app.CalculateFinalSalePriceResponse{
				CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
					{
						ProductVariantID: variantID,
						Quantity:         5,
						BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 15.0, Currency: "EUR"},
						FinalPrice:       pricing_app.MoneyDTO{Amount: 15.0, Currency: "EUR"},
					},
				},
				SaleTotal: pricing_app.MoneyDTO{Amount: 75.0, Currency: "EUR"},
			}, nil
		},
	}

	service := application.NewSalesService(
		&stubQuoteRepo{}, orderRepo, &stubDeliveryNoteRepo{}, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, pricingEngine, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.PUT("/orders/:id/line-items/:lineItemId", handler.UpdateOrderLineItem)

	reqBody := map[string]interface{}{
		"quantity": 5,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/orders/"+orderID.String()+"/line-items/"+lineItemID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Logf("❌ Expected 200, got %d. Response: %s", rec.Code, rec.Body.String())
	}
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRemoveOrderLineItem_Success(t *testing.T) {
	orderID := uuid.New()
	lineItemID := uuid.New()
	variantID := uuid.New()
	variantID2 := uuid.New()
	money, _ := domain.NewMoney(100.0, "EUR")

	existingOrder := &domain.SalesOrder{
		ID:           orderID,
		PartyID:      uuid.New(),
		OrderDate:    time.Now(),
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Status:       domain.SalesOrderStatusPending,
		LineItems: []domain.OrderLineItem{
			{
				ID:               lineItemID,
				ProductVariantID: variantID,
				Quantity:         2,
				ListUnitPrice:    money,
				UnitPrice:        money,
				DiscountPerUnit:  money,
				Subtotal:         money,
			},
			{
				ID:               uuid.New(),
				ProductVariantID: variantID2,
				Quantity:         1,
				ListUnitPrice:    money,
				UnitPrice:        money,
				DiscountPerUnit:  money,
				Subtotal:         money,
			},
		},
		Subtotal:  money,
		TaxAmount: money,
		Total:     money,
	}

	orderRepo := &stubOrderRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
			if id == orderID {
				return existingOrder, nil
			}
			return nil, domain.NewNotFoundError("order not found")
		},
		saveFn: func(ctx context.Context, order *domain.SalesOrder) error {
			return nil
		},
	}

	pricingEngine := &stubPricingEngine{
		calculateFinalSalePriceFn: func(ctx context.Context, req pricing_app.CalculateFinalSalePriceRequest) (*pricing_app.CalculateFinalSalePriceResponse, error) {
			// After removing, should recalc the remaining item
			items := make([]pricing_app.CalculatedSaleItemResponse, len(req.SaleItems))
			var total float64
			for i, item := range req.SaleItems {
				price := 100.0
				subtotal := price * float64(item.Quantity)
				total += subtotal
				items[i] = pricing_app.CalculatedSaleItemResponse{
					ProductVariantID: item.ProductVariantID,
					Quantity:         item.Quantity,
					BaseSalesPrice:   pricing_app.MoneyDTO{Amount: price, Currency: "EUR"},
					FinalPrice:       pricing_app.MoneyDTO{Amount: price, Currency: "EUR"},
				}
			}
			return &pricing_app.CalculateFinalSalePriceResponse{
				CalculatedItems: items,
				SaleTotal:       pricing_app.MoneyDTO{Amount: total, Currency: "EUR"},
			}, nil
		},
	}

	service := application.NewSalesService(
		&stubQuoteRepo{}, orderRepo, &stubDeliveryNoteRepo{}, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, pricingEngine, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.DELETE("/orders/:id/line-items/:lineItemId", handler.RemoveOrderLineItem)

	req := httptest.NewRequest(http.MethodDelete, "/orders/"+orderID.String()+"/line-items/"+lineItemID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Logf("❌ Expected 200, got %d. Response: %s", rec.Code, rec.Body.String())
	}
	assert.Equal(t, http.StatusOK, rec.Code)
}

// ===== ADDITIONAL DELIVERY NOTE HANDLER TESTS =====

func TestGetDeliveryNote_Success(t *testing.T) {
	noteID := uuid.New()

	noteNumber, _ := domain.NewDeliveryNoteNumber("DN/2026/0001")
	note := &domain.DeliveryNote{
		ID:                 noteID,
		DeliveryNoteNumber: noteNumber,
		SalesOrderID:       uuid.New(),
		PartyID:            uuid.New(),
		DeliveryDate:       time.Now(),
		Status:             domain.DeliveryNoteStatusPending,
		LineItems:          []domain.DeliveryNoteLineItem{},
	}

	noteRepo := &stubDeliveryNoteRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.DeliveryNote, error) {
			if id == noteID {
				return note, nil
			}
			return nil, domain.NewNotFoundError("delivery note not found")
		},
	}

	service := application.NewSalesService(
		&stubQuoteRepo{}, &stubOrderRepo{}, noteRepo, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.GET("/delivery-notes/:id", handler.GetDeliveryNote)

	req := httptest.NewRequest(http.MethodGet, "/delivery-notes/"+noteID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.NotEmpty(t, result["id"])
}

func TestListDeliveryNotes_Success(t *testing.T) {
	orderID := uuid.New()
	noteID1 := uuid.New()
	noteID2 := uuid.New()

	noteNumber1, _ := domain.NewDeliveryNoteNumber("DN/2026/0001")
	noteNumber2, _ := domain.NewDeliveryNoteNumber("DN/2026/0002")
	note1 := &domain.DeliveryNote{
		ID:                 noteID1,
		DeliveryNoteNumber: noteNumber1,
		SalesOrderID:       orderID,
		PartyID:            uuid.New(),
		DeliveryDate:       time.Now(),
		Status:             domain.DeliveryNoteStatusPending,
		LineItems:          []domain.DeliveryNoteLineItem{},
	}
	note2 := &domain.DeliveryNote{
		ID:                 noteID2,
		DeliveryNoteNumber: noteNumber2,
		SalesOrderID:       orderID,
		PartyID:            uuid.New(),
		DeliveryDate:       time.Now(),
		Status:             domain.DeliveryNoteStatusDelivered,
		LineItems:          []domain.DeliveryNoteLineItem{},
	}

	noteRepo := &stubDeliveryNoteRepo{
		listFn: func(ctx context.Context, filter domain.DeliveryNoteFilter) ([]*domain.DeliveryNote, error) {
			return []*domain.DeliveryNote{note1, note2}, nil
		},
	}

	service := application.NewSalesService(
		&stubQuoteRepo{}, &stubOrderRepo{}, noteRepo, &stubInvoiceRepo{},
		&stubDocumentNumberGenerator{}, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.GET("/delivery-notes", handler.ListDeliveryNotes)

	req := httptest.NewRequest(http.MethodGet, "/delivery-notes?salesOrderId="+orderID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

// ===== ADDITIONAL INVOICE HANDLER TESTS =====

func TestCreateSimplifiedInvoice_Success(t *testing.T) {
	orderID := uuid.New()
	partyID := uuid.New()
	variantID := uuid.New()

	money, _ := domain.NewMoney(100.0, "EUR")
	order := &domain.SalesOrder{
		ID:           orderID,
		PartyID:      partyID,
		OrderDate:    time.Now(),
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Status:       domain.SalesOrderStatusDelivered,
		LineItems: []domain.OrderLineItem{
			{
				ID:               uuid.New(),
				ProductVariantID: variantID,
				Quantity:         10,
				ListUnitPrice:    money,
				UnitPrice:        money,
				DiscountPerUnit:  money,
				Subtotal:         money,
			},
		},
		Subtotal:  money,
		TaxAmount: money,
		Total:     money,
	}

	orderRepo := &stubOrderRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
			if id == orderID {
				return order, nil
			}
			return nil, domain.NewNotFoundError("order not found")
		},
	}

	invoiceRepo := &stubInvoiceRepo{
		saveFn: func(ctx context.Context, invoice *domain.Invoice) error {
			return nil
		},
	}
	numberGen := &stubDocumentNumberGenerator{}

	service := application.NewSalesService(
		&stubQuoteRepo{}, orderRepo, &stubDeliveryNoteRepo{}, invoiceRepo,
		numberGen, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.POST("/invoices/simplified", handler.CreateSimplifiedInvoice)

	reqBody := map[string]interface{}{
		"partyId":         partyID.String(),
		"salesOrderIds":   []string{orderID.String()},
		"deliveryNoteIds": []string{},
		"invoiceDate":     time.Now().Format(time.RFC3339),
		"dueDate":         time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/invoices/simplified", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Logf("❌ Expected 201, got %d. Response: %s", rec.Code, rec.Body.String())
	}
	assert.Equal(t, http.StatusCreated, rec.Code)

	var result map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.NotEmpty(t, result["id"])
	assert.NotEmpty(t, result["invoiceNumber"])
}

func TestListInvoices_Success(t *testing.T) {
	partyID := uuid.New()
	invoiceID1 := uuid.New()
	invoiceID2 := uuid.New()

	money, _ := domain.NewMoney(100.0, "EUR")
	invoiceNumber1, _ := domain.NewInvoiceNumber("A/2026/0001")
	invoiceNumber2, _ := domain.NewInvoiceNumber("A/2026/0002")

	invoice1 := &domain.Invoice{
		ID:            invoiceID1,
		InvoiceNumber: invoiceNumber1,
		PartyID:       partyID,
		InvoiceDate:   time.Now(),
		DueDate:       time.Now().Add(30 * 24 * time.Hour),
		Status:        domain.InvoiceStatusIssued,
		Subtotal:      money,
		TaxAmount:     money,
		Total:         money,
	}
	invoice2 := &domain.Invoice{
		ID:            invoiceID2,
		InvoiceNumber: invoiceNumber2,
		PartyID:       partyID,
		InvoiceDate:   time.Now(),
		DueDate:       time.Now().Add(30 * 24 * time.Hour),
		Status:        domain.InvoiceStatusPaid,
		Subtotal:      money,
		TaxAmount:     money,
		Total:         money,
	}

	invoiceRepo := &stubInvoiceRepo{
		listFn: func(ctx context.Context, filter domain.InvoiceFilter) ([]*domain.Invoice, error) {
			return []*domain.Invoice{invoice1, invoice2}, nil
		},
	}

	service := application.NewSalesService(
		&stubQuoteRepo{}, &stubOrderRepo{}, &stubDeliveryNoteRepo{}, invoiceRepo,
		&stubDocumentNumberGenerator{}, &stubPricingEngine{}, &stubPartyLookup{}, nil, nil,
	)
	handler := NewSalesHandler(service)

	router := setupTestRouter()
	router.GET("/invoices", handler.ListInvoices)

	req := httptest.NewRequest(http.MethodGet, "/invoices?partyId="+partyID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}
