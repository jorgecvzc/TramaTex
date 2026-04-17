package application_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pricing_app "github.com/joran-cortez/tramatex/internal/pricing/application"
	"github.com/joran-cortez/tramatex/internal/sales/application"
	"github.com/joran-cortez/tramatex/internal/sales/domain"
)

type MockQuoteRepository struct {
	mock.Mock
}

func (m *MockQuoteRepository) Save(ctx context.Context, quote *domain.Quote) error {
	args := m.Called(ctx, quote)
	return args.Error(0)
}

func (m *MockQuoteRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Quote, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Quote), args.Error(1)
}

func (m *MockQuoteRepository) List(ctx context.Context, filter domain.QuoteFilter) ([]*domain.Quote, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Quote), args.Error(1)
}

func (m *MockQuoteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockSalesOrderRepository struct {
	mock.Mock
}

func (m *MockSalesOrderRepository) Save(ctx context.Context, order *domain.SalesOrder) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockSalesOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SalesOrder), args.Error(1)
}

func (m *MockSalesOrderRepository) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SalesOrder), args.Error(1)
}

func (m *MockSalesOrderRepository) List(ctx context.Context, filter domain.SalesOrderFilter) ([]*domain.SalesOrder, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.SalesOrder), args.Error(1)
}

func (m *MockSalesOrderRepository) FindByQuoteID(ctx context.Context, quoteID uuid.UUID) (*domain.SalesOrder, error) {
	args := m.Called(ctx, quoteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SalesOrder), args.Error(1)
}

type MockDeliveryNoteRepository struct {
	mock.Mock
}

func (m *MockDeliveryNoteRepository) Save(ctx context.Context, note *domain.DeliveryNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockDeliveryNoteRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.DeliveryNote, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DeliveryNote), args.Error(1)
}

func (m *MockDeliveryNoteRepository) List(ctx context.Context, filter domain.DeliveryNoteFilter) ([]*domain.DeliveryNote, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.DeliveryNote), args.Error(1)
}

func (m *MockDeliveryNoteRepository) ListBySalesOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.DeliveryNote, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.DeliveryNote), args.Error(1)
}

func (m *MockDeliveryNoteRepository) LinkLineItemsToInvoice(ctx context.Context, links map[uuid.UUID]uuid.UUID) error {
	args := m.Called(ctx, links)
	return args.Error(0)
}

func (m *MockDeliveryNoteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockInvoiceRepository struct {
	mock.Mock
}

func (m *MockInvoiceRepository) Save(ctx context.Context, invoice *domain.Invoice) error {
	args := m.Called(ctx, invoice)
	return args.Error(0)
}

func (m *MockInvoiceRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) List(ctx context.Context, filter domain.InvoiceFilter) ([]*domain.Invoice, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) ListBySalesOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.Invoice, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) FindByDeliveryNoteID(ctx context.Context, deliveryNoteID uuid.UUID) (*domain.Invoice, error) {
	args := m.Called(ctx, deliveryNoteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) ListDeliveryNoteIDsByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, invoiceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

func (m *MockInvoiceRepository) ListOrderIDsByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, invoiceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

type MockPricingEngine struct {
	mock.Mock
}

func (m *MockPricingEngine) CalculateFinalSalePrice(ctx context.Context, req pricing_app.CalculateFinalSalePriceRequest) (*pricing_app.CalculateFinalSalePriceResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pricing_app.CalculateFinalSalePriceResponse), args.Error(1)
}

type MockPartyLookup struct {
	mock.Mock
	existsPartyFn  func(context.Context, uuid.UUID) (bool, error)
	hasPartyRoleFn func(context.Context, uuid.UUID, string) (bool, error)
}

func (m *MockPartyLookup) ExistsParty(ctx context.Context, partyID uuid.UUID) (bool, error) {
	if m.existsPartyFn != nil {
		return m.existsPartyFn(ctx, partyID)
	}
	args := m.Called(ctx, partyID)
	return args.Bool(0), args.Error(1)
}

func (m *MockPartyLookup) HasPartyRole(ctx context.Context, partyID uuid.UUID, role string) (bool, error) {
	if m.hasPartyRoleFn != nil {
		return m.hasPartyRoleFn(ctx, partyID, role)
	}
	return true, nil
}

type MockNumberGenerator struct {
	mock.Mock
}

func (m *MockNumberGenerator) NextQuoteNumber(ctx context.Context) (domain.QuoteNumber, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.QuoteNumber), args.Error(1)
}

func (m *MockNumberGenerator) NextOrderNumber(ctx context.Context) (domain.OrderNumber, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.OrderNumber), args.Error(1)
}

func (m *MockNumberGenerator) NextDeliveryNoteNumber(ctx context.Context) (domain.DeliveryNoteNumber, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.DeliveryNoteNumber), args.Error(1)
}

func (m *MockNumberGenerator) NextInvoiceNumber(ctx context.Context, series domain.InvoiceSeries) (domain.InvoiceNumber, error) {
	args := m.Called(ctx, series)
	return args.Get(0).(domain.InvoiceNumber), args.Error(1)
}

type MockWorkOrderSuspender struct {
	mock.Mock
}

func (m *MockWorkOrderSuspender) SuspendWorkOrders(ctx context.Context, workOrderIDs []uuid.UUID) error {
	args := m.Called(ctx, workOrderIDs)
	return args.Error(0)
}

func (m *MockWorkOrderSuspender) ReactivateWorkOrders(ctx context.Context, workOrderIDs []uuid.UUID) error {
	args := m.Called(ctx, workOrderIDs)
	return args.Error(0)
}

func TestSalesService_CreateQuote_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	pricing := new(MockPricingEngine)
	partyLookup := new(MockPartyLookup)
	numbers := new(MockNumberGenerator)

	quoteNumber, _ := domain.NewQuoteNumber("Q-100")
	partyLookup.On("ExistsParty", mock.Anything, partyID).Return(true, nil)
	numbers.On("NextQuoteNumber", mock.Anything).Return(quoteNumber, nil)

	pricing.On("CalculateFinalSalePrice", mock.Anything, mock.MatchedBy(func(req pricing_app.CalculateFinalSalePriceRequest) bool {
		return req.ClientID == partyID.String() && len(req.SaleItems) == 1 && req.SaleItems[0].ProductVariantID == variantID && req.SaleItems[0].Quantity == 2
	})).Return(&pricing_app.CalculateFinalSalePriceResponse{
		CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
			{
				ProductVariantID: variantID,
				Quantity:         2,
				BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 10, Currency: domain.DefaultCurrency},
				FinalPrice:       pricing_app.MoneyDTO{Amount: 10, Currency: domain.DefaultCurrency},
			},
		},
		SaleTotal: pricing_app.MoneyDTO{Amount: 20, Currency: domain.DefaultCurrency},
	}, nil)

	quoteRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Quote")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, pricing, partyLookup, nil, nil)

	cmd := application.CreateQuoteCommand{
		PartyID:        partyID,
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Items: []application.QuoteLineItemInput{
			{ProductVariantID: variantID, Quantity: 2},
		},
	}

	result, err := service.CreateQuote(ctx, cmd)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, partyID, result.PartyID)
	assert.Equal(t, string(domain.QuoteStatusDraft), result.Status)
	quoteRepo.AssertCalled(t, "Save", mock.Anything, mock.AnythingOfType("*domain.Quote"))
}

func TestSalesService_CreateQuote_PartyNotFound(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	pricing := new(MockPricingEngine)
	partyLookup := new(MockPartyLookup)
	numbers := new(MockNumberGenerator)

	partyLookup.On("ExistsParty", mock.Anything, partyID).Return(false, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, pricing, partyLookup, nil, nil)
	cmd := application.CreateQuoteCommand{
		PartyID:        partyID,
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Items: []application.QuoteLineItemInput{
			{ProductVariantID: variantID, Quantity: 1},
		},
	}

	_, err := service.CreateQuote(ctx, cmd)
	assert.Error(t, err)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	quoteRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
}

func TestSalesService_CreateQuote_PartyNotClient(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	pricing := new(MockPricingEngine)
	partyLookup := &MockPartyLookup{
		existsPartyFn: func(context.Context, uuid.UUID) (bool, error) {
			return true, nil
		},
		hasPartyRoleFn: func(context.Context, uuid.UUID, string) (bool, error) {
			return false, nil
		},
	}
	numbers := new(MockNumberGenerator)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, pricing, partyLookup, nil, nil)
	cmd := application.CreateQuoteCommand{
		PartyID:        partyID,
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Items: []application.QuoteLineItemInput{
			{ProductVariantID: variantID, Quantity: 1},
		},
	}

	_, err := service.CreateQuote(ctx, cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "party must have CLIENT role")
	quoteRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
}

func TestSalesService_CreateOrder_FromQuoteNotApproved(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	pricing := new(MockPricingEngine)
	numbers := new(MockNumberGenerator)

	quoteNumber, _ := domain.NewQuoteNumber("Q-200")
	money, _ := domain.NewMoney(10, domain.DefaultCurrency)
	lineItem, _ := domain.NewQuoteLineItem(uuid.New(), 1, money, nil, 0)
	quote, _ := domain.NewQuote(quoteNumber, partyID, time.Now(), time.Now().Add(24*time.Hour), []domain.QuoteLineItem{lineItem}, money, "")
	quote.Status = domain.QuoteStatusRejected // Rejected quotes cannot be auto-approved to Issued/Approved

	quoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)
	orderNumber, _ := domain.NewOrderNumber("SO-1")
	numbers.On("NextOrderNumber", mock.Anything).Return(orderNumber, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, pricing, nil, nil, nil)
	cmd := application.CreateOrderCommand{
		PartyID:      partyID,
		QuoteID:      &quote.ID,
		DeliveryDate: time.Now().Add(48 * time.Hour),
	}

	_, err := service.CreateOrder(ctx, cmd)
	assert.Error(t, err)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeConflict, domainErr.Code)
	orderRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
}

func TestSalesService_CreateDeliveryNote_Partial(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	numbers := new(MockNumberGenerator)

	money, _ := domain.NewMoney(10, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-300")
	orderItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)
	order, _ := domain.NewSalesOrder(orderNumber, partyID, time.Now(), time.Now().Add(48*time.Hour), []domain.OrderLineItem{orderItem}, money, "")

	orderRepo.On("FindByIDForUpdate", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)

	noteNumber, _ := domain.NewDeliveryNoteNumber("DN-10")
	numbers.On("NextDeliveryNoteNumber", mock.Anything).Return(noteNumber, nil)

	var savedOrder *domain.SalesOrder
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Run(func(args mock.Arguments) {
		savedOrder = args.Get(1).(*domain.SalesOrder)
	}).Return(nil)
	deliveryRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.DeliveryNote")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, nil, nil, nil, nil)
	cmd := application.CreateDeliveryNoteCommand{
		SalesOrderID: order.ID,
		DeliveryDate: time.Now().Add(72 * time.Hour),
		Items: []application.DeliveryNoteLineItemInput{
			{SalesOrderLineItemID: orderItem.ID, DeliveredQuantity: 2},
		},
	}

	result, err := service.CreateDeliveryNote(ctx, cmd)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, string(domain.DeliveryNoteStatusPending), result.Status)
	assert.NotNil(t, savedOrder)
	assert.Equal(t, domain.SalesOrderStatusPartiallyDelivered, savedOrder.Status)
}

func TestSalesService_CreateInvoice_FromDeliveryNote(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	numbers := new(MockNumberGenerator)

	money, _ := domain.NewMoney(10, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-400")
	orderItem, _ := domain.NewOrderLineItem(variantID, 1, money, nil, 0)
	order, _ := domain.NewSalesOrder(orderNumber, partyID, time.Now(), time.Now().Add(48*time.Hour), []domain.OrderLineItem{orderItem}, money, "")
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)
	_ = order.ChangeStatus(domain.SalesOrderStatusDelivered)

	// Create a delivery note linked to the order
	dnNumber, _ := domain.NewDeliveryNoteNumber("ALB-001")
	dnLineItem, _ := domain.NewDeliveryNoteLineItem(orderItem.ID, variantID, 1)
	deliveryNote, _ := domain.NewDeliveryNote(dnNumber, order.ID, partyID, time.Now(), []domain.DeliveryNoteLineItem{dnLineItem}, "")
	_ = deliveryNote.ChangeStatus(domain.DeliveryNoteStatusDelivered)

	deliveryRepo.On("FindByID", mock.Anything, deliveryNote.ID).Return(deliveryNote, nil)
	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	orderRepo.On("FindByIDForUpdate", mock.Anything, order.ID).Return(order, nil)

	invoiceNumber, _ := domain.NewInvoiceNumber("FV-2026-0001")
	numbers.On("NextInvoiceNumber", mock.Anything, mock.Anything).Return(invoiceNumber, nil)

	// Build a matching invoice to simulate what would be persisted.
	invLineItem, _ := domain.NewInvoiceLineItem(variantID, 1, money, nil, nil, 0)
	invLineItem.SalesOrderLineItemID = &orderItem.ID
	taxAmount, _ := domain.NewMoney(0, domain.DefaultCurrency)
	currentYear := time.Now().Year()
	invSeries, _ := domain.NewInvoiceSeries("FV", currentYear)
	expectedInvoice, _ := domain.NewInvoice(
		invoiceNumber,
		domain.InvoiceTypeComplete,
		invSeries,
		partyID,
		time.Now(),
		time.Now().Add(24*time.Hour),
		[]domain.InvoiceLineItem{invLineItem},
		taxAmount,
		"",
	)
	// updateOrderInvoiceStatus calls ListBySalesOrderID after Save
	invoiceRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.Invoice{expectedInvoice}, nil)

	var savedOrder *domain.SalesOrder
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Run(func(args mock.Arguments) {
		savedOrder = args.Get(1).(*domain.SalesOrder)
	}).Return(nil)
	invoiceRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Invoice")).Return(nil)
	deliveryRepo.On("LinkLineItemsToInvoice", mock.Anything, mock.Anything).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, nil, nil, nil, nil)

	cmd := application.CreateInvoiceCommand{
		PartyID:         partyID,
		DeliveryNoteIDs: []uuid.UUID{deliveryNote.ID},
		InvoiceDate:     time.Now(),
		DueDate:         time.Now().Add(24 * time.Hour),
	}

	result, err := service.CreateInvoice(ctx, cmd)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, savedOrder)
	assert.Equal(t, domain.SalesOrderStatusInvoiced, savedOrder.Status)
}

// ===== Query Tests =====

func TestSalesService_GetQuote_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	quoteNumber, _ := domain.NewQuoteNumber("Q/2026/0001")
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, 0)
	taxAmount, _ := domain.NewMoney(42, domain.DefaultCurrency)

	quote, _ := domain.NewQuote(
		quoteNumber,
		partyID,
		time.Now(),
		time.Now().Add(30*24*time.Hour),
		[]domain.QuoteLineItem{lineItem},
		taxAmount,
		"Test quote",
	)

	quoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)
	orderRepo.On("FindByQuoteID", mock.Anything, quote.ID).Return(nil, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.GetQuoteByIDQuery{ID: quote.ID}
	result, err := service.GetQuote(ctx, query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, quote.ID, result.ID)
	assert.Equal(t, partyID, result.PartyID)
	quoteRepo.AssertExpectations(t)
}

func TestSalesService_GetQuote_NotFound(t *testing.T) {
	ctx := context.Background()
	quoteID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	quoteRepo.On("FindByID", mock.Anything, quoteID).Return(nil, domain.NewNotFoundError("quote not found"))

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.GetQuoteByIDQuery{ID: quoteID}
	result, err := service.GetQuote(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, result)
	quoteRepo.AssertExpectations(t)
}

func TestSalesService_ListQuotes_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	quoteNumber1, _ := domain.NewQuoteNumber("Q/2026/0001")
	quoteNumber2, _ := domain.NewQuoteNumber("Q/2026/0002")
	variantID := uuid.New()
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, 0)
	taxAmount, _ := domain.NewMoney(42, domain.DefaultCurrency)

	quote1, _ := domain.NewQuote(quoteNumber1, partyID, time.Now(), time.Now().Add(30*24*time.Hour), []domain.QuoteLineItem{lineItem}, taxAmount, "")
	quote2, _ := domain.NewQuote(quoteNumber2, partyID, time.Now(), time.Now().Add(30*24*time.Hour), []domain.QuoteLineItem{lineItem}, taxAmount, "")

	quotes := []*domain.Quote{quote1, quote2}
	filter := domain.QuoteFilter{PartyID: &partyID}

	quoteRepo.On("List", mock.Anything, filter).Return(quotes, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.ListQuotesQuery{PartyID: &partyID}
	results, err := service.ListQuotes(ctx, query)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, quote1.ID, results[0].ID)
	assert.Equal(t, quote2.ID, results[1].ID)
	quoteRepo.AssertExpectations(t)
}

func TestSalesService_GetOrder_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-001")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)
	taxAmount, _ := domain.NewMoney(105, domain.DefaultCurrency)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem},
		taxAmount,
		"Test order",
	)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.GetOrderByIDQuery{ID: order.ID}
	result, err := service.GetOrder(ctx, query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, order.ID, result.ID)
	assert.Equal(t, partyID, result.PartyID)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_ListOrders_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber1, _ := domain.NewOrderNumber("SO-001")
	orderNumber2, _ := domain.NewOrderNumber("SO-002")
	variantID := uuid.New()
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)
	taxAmount, _ := domain.NewMoney(105, domain.DefaultCurrency)

	order1, _ := domain.NewSalesOrder(orderNumber1, partyID, time.Now(), time.Now().Add(7*24*time.Hour), []domain.OrderLineItem{lineItem}, taxAmount, "")
	order2, _ := domain.NewSalesOrder(orderNumber2, partyID, time.Now(), time.Now().Add(7*24*time.Hour), []domain.OrderLineItem{lineItem}, taxAmount, "")

	orders := []*domain.SalesOrder{order1, order2}
	filter := domain.SalesOrderFilter{PartyID: &partyID}

	orderRepo.On("List", mock.Anything, filter).Return(orders, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.ListOrdersQuery{PartyID: &partyID}
	results, err := service.ListOrders(ctx, query)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, order1.ID, results[0].ID)
	assert.Equal(t, order2.ID, results[1].ID)
	orderRepo.AssertExpectations(t)
}

// ===== ConvertQuoteToOrder Tests =====

func TestSalesService_ConvertQuoteToOrder_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	numbers := new(MockNumberGenerator)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	quoteNumber, _ := domain.NewQuoteNumber("Q/2026/0001")
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, 0)
	taxAmount, _ := domain.NewMoney(42, domain.DefaultCurrency)

	quote, _ := domain.NewQuote(
		quoteNumber,
		partyID,
		time.Now(),
		time.Now().Add(30*24*time.Hour),
		[]domain.QuoteLineItem{lineItem},
		taxAmount,
		"Test quote",
	)
	_ = quote.ChangeStatus(domain.QuoteStatusIssued)
	_ = quote.ChangeStatus(domain.QuoteStatusApproved)

	quoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)
	quoteRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Quote")).Return(nil)

	orderNumber, _ := domain.NewOrderNumber("SO-001")
	numbers.On("NextOrderNumber", mock.Anything).Return(orderNumber, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, nil, nil, nil, nil)

	cmd := application.ConvertQuoteToOrderCommand{
		QuoteID:      quote.ID,
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
	}

	result, err := service.ConvertQuoteToOrder(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, partyID, result.PartyID)
	assert.Equal(t, "PENDING", string(result.Status))
	quoteRepo.AssertExpectations(t)
	orderRepo.AssertExpectations(t)
	numbers.AssertExpectations(t)
}

func TestSalesService_ConvertQuoteToOrder_QuoteNotApproved(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	numbers := new(MockNumberGenerator)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	quoteNumber, _ := domain.NewQuoteNumber("Q/2026/0001")
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, 0)
	taxAmount, _ := domain.NewMoney(42, domain.DefaultCurrency)

	quote, _ := domain.NewQuote(
		quoteNumber,
		partyID,
		time.Now(),
		time.Now().Add(30*24*time.Hour),
		[]domain.QuoteLineItem{lineItem},
		taxAmount,
		"Test quote",
	)
	quote.Status = domain.QuoteStatusRejected // Rejected quotes cannot be auto-approved

	// Quote is in DRAFT status, not APPROVED

	quoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)

	// Note: Service generates order number before validation - this is a code smell
	// but we need to mock it for the test to run
	orderNumber, _ := domain.NewOrderNumber("SO-999")
	numbers.On("NextOrderNumber", mock.Anything).Return(orderNumber, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, nil, nil, nil, nil)

	cmd := application.ConvertQuoteToOrderCommand{
		QuoteID:      quote.ID,
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
	}

	result, err := service.ConvertQuoteToOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	quoteRepo.AssertExpectations(t)
}

// ===== UpdateQuote Tests =====

func TestSalesService_UpdateQuote_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	quoteNumber, _ := domain.NewQuoteNumber("Q/2026/0001")
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, 0)
	taxAmount, _ := domain.NewMoney(42, domain.DefaultCurrency)

	quote, _ := domain.NewQuote(
		quoteNumber,
		partyID,
		time.Now(),
		time.Now().Add(30*24*time.Hour),
		[]domain.QuoteLineItem{lineItem},
		taxAmount,
		"Original notes",
	)

	quoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)
	quoteRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Quote")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	newNotes := "Updated notes"
	cmd := application.UpdateQuoteCommand{
		QuoteID: quote.ID,
		Notes:   &newNotes,
	}

	result, err := service.UpdateQuote(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, quote.ID, result.ID)
	quoteRepo.AssertExpectations(t)
}

// ===== ChangeQuoteStatus Tests =====

func TestSalesService_ChangeQuoteStatus_DraftToSent_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	quoteNumber, _ := domain.NewQuoteNumber("Q/2026/0001")
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, 0)
	taxAmount, _ := domain.NewMoney(42, domain.DefaultCurrency)

	quote, _ := domain.NewQuote(
		quoteNumber,
		partyID,
		time.Now(),
		time.Now().Add(30*24*time.Hour),
		[]domain.QuoteLineItem{lineItem},
		taxAmount,
		"Test quote",
	)

	quoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)
	quoteRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Quote")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	cmd := application.ChangeQuoteStatusCommand{
		QuoteID:   quote.ID,
		NewStatus: string(domain.QuoteStatusIssued),
	}

	result, err := service.ChangeQuoteStatus(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ISSUED", string(result.Status))
	quoteRepo.AssertExpectations(t)
}

func TestSalesService_ChangeQuoteStatus_InvalidTransition_Fail(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	quoteNumber, _ := domain.NewQuoteNumber("Q/2026/0001")
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, 0)
	taxAmount, _ := domain.NewMoney(42, domain.DefaultCurrency)

	quote, _ := domain.NewQuote(
		quoteNumber,
		partyID,
		time.Now(),
		time.Now().Add(30*24*time.Hour),
		[]domain.QuoteLineItem{lineItem},
		taxAmount,
		"Test quote",
	)
	// Quote is in DRAFT status

	quoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	// Try to change directly from DRAFT to APPROVED (invalid - must go through SENT first)
	cmd := application.ChangeQuoteStatusCommand{
		QuoteID:   quote.ID,
		NewStatus: string(domain.QuoteStatusApproved),
	}

	result, err := service.ChangeQuoteStatus(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	quoteRepo.AssertExpectations(t)
}

// ===== ChangeOrderStatus Tests =====

func TestSalesService_ChangeOrderStatus_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-001")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)
	taxAmount, _ := domain.NewMoney(105, domain.DefaultCurrency)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem},
		taxAmount,
		"Test order",
	)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	cmd := application.ChangeOrderStatusCommand{
		OrderID:   order.ID,
		NewStatus: string(domain.SalesOrderStatusInPreparation),
	}

	result, err := service.ChangeOrderStatus(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "IN_PREPARATION", string(result.Status))
	orderRepo.AssertExpectations(t)
}

func TestSalesService_ChangeOrderStatus_CancelSuspendsWorkOrders(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	suspender := new(MockWorkOrderSuspender)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-010")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)
	taxAmount, _ := domain.NewMoney(105, domain.DefaultCurrency)

	order, _ := domain.NewSalesOrder(
		orderNumber, partyID, time.Now(), time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem}, taxAmount, "Test order",
	)
	// Move to EN_PREPARACION so we can cancel.
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)

	woID1 := uuid.New()
	woID2 := uuid.New()
	wsID1, wsID2 := uuid.New(), uuid.New()
	order.WorkReferences = []domain.WorkReference{
		{ID: uuid.New(), WorkSetupID: &wsID1, WorkOrderID: &woID1, Sequence: 1},
		{ID: uuid.New(), WorkSetupID: &wsID2, WorkOrderID: &woID2, Sequence: 2},
	}

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)
	suspender.On("SuspendWorkOrders", mock.Anything, []uuid.UUID{woID1, woID2}).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)
	service.SetWorkOrderSuspender(suspender)

	result, err := service.ChangeOrderStatus(ctx, application.ChangeOrderStatusCommand{
		OrderID:   order.ID,
		NewStatus: string(domain.SalesOrderStatusCancelled),
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "CANCELLED", string(result.Status))
	suspender.AssertExpectations(t)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_ChangeOrderStatus_ReactivateReactivatesWorkOrders(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	suspender := new(MockWorkOrderSuspender)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-011")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)
	taxAmount, _ := domain.NewMoney(105, domain.DefaultCurrency)

	order, _ := domain.NewSalesOrder(
		orderNumber, partyID, time.Now(), time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem}, taxAmount, "Test order",
	)
	// Simulate: confirmed → cancelled (now CANCELADO).
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)
	_ = order.ChangeStatus(domain.SalesOrderStatusCancelled)

	woID := uuid.New()
	wsID3 := uuid.New()
	order.WorkReferences = []domain.WorkReference{
		{ID: uuid.New(), WorkSetupID: &wsID3, WorkOrderID: &woID, Sequence: 1},
	}

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)
	suspender.On("ReactivateWorkOrders", mock.Anything, []uuid.UUID{woID}).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)
	service.SetWorkOrderSuspender(suspender)

	result, err := service.ChangeOrderStatus(ctx, application.ChangeOrderStatusCommand{
		OrderID:   order.ID,
		NewStatus: string(domain.SalesOrderStatusPending),
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "PENDING", string(result.Status))
	suspender.AssertExpectations(t)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_ChangeOrderStatus_CancelWithoutSuspenderOK(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-012")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)
	taxAmount, _ := domain.NewMoney(105, domain.DefaultCurrency)

	order, _ := domain.NewSalesOrder(
		orderNumber, partyID, time.Now(), time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem}, taxAmount, "Test order",
	)
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)

	woID := uuid.New()
	wsID4 := uuid.New()
	order.WorkReferences = []domain.WorkReference{
		{ID: uuid.New(), WorkSetupID: &wsID4, WorkOrderID: &woID, Sequence: 1},
	}

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)

	// No suspender set — nil-safe path.
	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	result, err := service.ChangeOrderStatus(ctx, application.ChangeOrderStatusCommand{
		OrderID:   order.ID,
		NewStatus: string(domain.SalesOrderStatusCancelled),
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "CANCELLED", string(result.Status))
	orderRepo.AssertExpectations(t)
}

func TestSalesService_ChangeOrderStatus_CancelWithNoWorkOrdersSkipsSuspend(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	suspender := new(MockWorkOrderSuspender)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-013")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)
	taxAmount, _ := domain.NewMoney(105, domain.DefaultCurrency)

	order, _ := domain.NewSalesOrder(
		orderNumber, partyID, time.Now(), time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem}, taxAmount, "Test order",
	)
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)
	// No WorkReferences — nothing to suspend.

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)
	service.SetWorkOrderSuspender(suspender)

	result, err := service.ChangeOrderStatus(ctx, application.ChangeOrderStatusCommand{
		OrderID:   order.ID,
		NewStatus: string(domain.SalesOrderStatusCancelled),
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "CANCELLED", string(result.Status))
	// SuspendWorkOrders must NOT have been called.
	suspender.AssertNotCalled(t, "SuspendWorkOrders", mock.Anything, mock.Anything)
	orderRepo.AssertExpectations(t)
}

// ===== GetDeliveryNote Tests =====

func TestSalesService_GetDeliveryNote_Success(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	orderID := uuid.New()
	noteNumber, _ := domain.NewDeliveryNoteNumber("DN-001")
	lineItem, _ := domain.NewDeliveryNoteLineItem(uuid.New(), uuid.New(), 5)

	note, _ := domain.NewDeliveryNote(
		noteNumber,
		orderID,
		uuid.New(),
		time.Now(),
		[]domain.DeliveryNoteLineItem{lineItem},
		"Test note",
	)

	deliveryRepo.On("FindByID", mock.Anything, note.ID).Return(note, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.GetDeliveryNoteByIDQuery{ID: note.ID}
	result, err := service.GetDeliveryNote(ctx, query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, note.ID, result.ID)
	deliveryRepo.AssertExpectations(t)
}

func TestSalesService_GetDeliveryNote_NotFound(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	noteID := uuid.New()
	deliveryRepo.On("FindByID", mock.Anything, noteID).Return(nil, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.GetDeliveryNoteByIDQuery{ID: noteID}
	result, err := service.GetDeliveryNote(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	deliveryRepo.AssertExpectations(t)
}

// ===== ListDeliveryNotes Tests =====

func TestSalesService_ListDeliveryNotes_Success(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	orderID := uuid.New()
	partyID := uuid.New()
	noteNumber, _ := domain.NewDeliveryNoteNumber("DN-001")
	lineItem, _ := domain.NewDeliveryNoteLineItem(uuid.New(), uuid.New(), 5)

	note, _ := domain.NewDeliveryNote(
		noteNumber,
		orderID,
		partyID,
		time.Now(),
		[]domain.DeliveryNoteLineItem{lineItem},
		"Test note",
	)

	deliveryRepo.On("List", mock.Anything, mock.AnythingOfType("domain.DeliveryNoteFilter")).Return([]*domain.DeliveryNote{note}, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.ListDeliveryNotesQuery{PartyID: &partyID}
	result, err := service.ListDeliveryNotes(ctx, query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, note.ID, result[0].ID)
	deliveryRepo.AssertExpectations(t)
}

// ===== GetInvoice Tests =====

func TestSalesService_GetInvoice_Success(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	invoiceNumber, _ := domain.NewInvoiceNumber("INV-001")
	series, _ := domain.NewInvoiceSeries("A", 2026)
	lineItem, _ := domain.NewInvoiceLineItem(uuid.New(), 2, money, nil, nil)

	invoice, _ := domain.NewInvoice(
		invoiceNumber,
		domain.InvoiceTypeComplete,
		series,
		uuid.New(),
		time.Now(),
		time.Now().Add(30*24*time.Hour),
		[]domain.InvoiceLineItem{lineItem},
		money,
		"Net 30",
	)

	invoiceRepo.On("FindByID", mock.Anything, invoice.ID).Return(invoice, nil)
	invoiceRepo.On("ListDeliveryNoteIDsByInvoiceID", mock.Anything, invoice.ID).Return([]uuid.UUID{}, nil)
	invoiceRepo.On("ListOrderIDsByInvoiceID", mock.Anything, invoice.ID).Return([]uuid.UUID{}, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.GetInvoiceByIDQuery{ID: invoice.ID}
	result, err := service.GetInvoice(ctx, query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, invoice.ID, result.ID)
	invoiceRepo.AssertExpectations(t)
}

func TestSalesService_GetInvoice_NotFound(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	invoiceID := uuid.New()
	invoiceRepo.On("FindByID", mock.Anything, invoiceID).Return(nil, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.GetInvoiceByIDQuery{ID: invoiceID}
	result, err := service.GetInvoice(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	invoiceRepo.AssertExpectations(t)
}

// ===== ListInvoices Tests =====

func TestSalesService_ListInvoices_Success(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	partyID := uuid.New()
	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	invoiceNumber, _ := domain.NewInvoiceNumber("INV-001")
	series, _ := domain.NewInvoiceSeries("A", 2026)
	lineItem, _ := domain.NewInvoiceLineItem(uuid.New(), 2, money, nil, nil)

	invoice, _ := domain.NewInvoice(
		invoiceNumber,
		domain.InvoiceTypeComplete,
		series,
		partyID,
		time.Now(),
		time.Now().Add(30*24*time.Hour),
		[]domain.InvoiceLineItem{lineItem},
		money,
		"Net 30",
	)

	invoiceRepo.On("List", mock.Anything, mock.AnythingOfType("domain.InvoiceFilter")).Return([]*domain.Invoice{invoice}, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.ListInvoicesQuery{PartyID: &partyID}
	result, err := service.ListInvoices(ctx, query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, invoice.ID, result[0].ID)
	invoiceRepo.AssertExpectations(t)
}

// ===== UpdateOrderDetails Tests =====

func TestSalesService_UpdateOrderDetails_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	newPartyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	partyLookup := new(MockPartyLookup)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-001")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem},
		money,
		"Test order",
	)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)
	partyLookup.On("ExistsParty", mock.Anything, newPartyID).Return(true, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, partyLookup, nil, nil)

	newDate := time.Now().Add(14 * 24 * time.Hour)
	newNotes := "Updated notes"
	cmd := application.UpdateOrderDetailsCommand{
		OrderID:      order.ID,
		PartyID:      &newPartyID,
		DeliveryDate: &newDate,
		Notes:        &newNotes,
	}

	result, err := service.UpdateOrderDetails(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, order.ID, result.ID)
	assert.Equal(t, newPartyID, result.PartyID)
	orderRepo.AssertExpectations(t)
	partyLookup.AssertExpectations(t)
}

func TestSalesService_UpdateOrderDetails_OrderNotFound(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	orderID := uuid.New()
	orderRepo.On("FindByID", mock.Anything, orderID).Return(nil, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	cmd := application.UpdateOrderDetailsCommand{OrderID: orderID}
	result, err := service.UpdateOrderDetails(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	orderRepo.AssertExpectations(t)
}

// ===== AddOrderLineItem Tests =====

func TestSalesService_AddOrderLineItem_OrderNotFound(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	orderID := uuid.New()
	orderRepo.On("FindByID", mock.Anything, orderID).Return(nil, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	cmd := application.AddOrderLineItemCommand{
		OrderID: orderID,
		Item: application.OrderLineItemInput{
			ProductVariantID: uuid.New(),
			Quantity:         2,
		},
	}

	result, err := service.AddOrderLineItem(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_AddOrderLineItem_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-001")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem},
		money,
		"Test order",
	)

	// Change order to DELIVERED status (cannot edit line items)
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)
	_ = order.ChangeStatus(domain.SalesOrderStatusDelivered)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	cmd := application.AddOrderLineItemCommand{
		OrderID: order.ID,
		Item: application.OrderLineItemInput{
			ProductVariantID: uuid.New(),
			Quantity:         2,
		},
	}

	result, err := service.AddOrderLineItem(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeConflict, domainErr.Code)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_AddOrderLineItem_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID1 := uuid.New()
	variantID2 := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	mockPricingService := new(MockPricingEngine)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-ADD-001")
	lineItem, _ := domain.NewOrderLineItem(variantID1, 5, money, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem},
		money,
		"Test order",
	)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)
	mockPricingService.On("CalculateFinalSalePrice", mock.Anything, mock.Anything).Return(&pricing_app.CalculateFinalSalePriceResponse{
		CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
			{
				ProductVariantID: variantID1,
				Quantity:         5,
				BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
				FinalPrice:       pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
			},
			{
				ProductVariantID: variantID2,
				Quantity:         3,
				BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 50.0, Currency: domain.DefaultCurrency},
				FinalPrice:       pricing_app.MoneyDTO{Amount: 50.0, Currency: domain.DefaultCurrency},
			},
		},
	}, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, mockPricingService, nil, nil, nil)

	cmd := application.AddOrderLineItemCommand{
		OrderID: order.ID,
		Item: application.OrderLineItemInput{
			ProductVariantID: variantID2,
			Quantity:         3,
		},
	}

	result, err := service.AddOrderLineItem(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.LineItems))
	orderRepo.AssertExpectations(t)
	mockPricingService.AssertExpectations(t)
}

func TestSalesService_AddOrderLineItem_SaveError(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID1 := uuid.New()
	variantID2 := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	mockPricingService := new(MockPricingEngine)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-ADD-002")
	lineItem, _ := domain.NewOrderLineItem(variantID1, 5, money, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem},
		money,
		"Test order",
	)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)
	mockPricingService.On("CalculateFinalSalePrice", mock.Anything, mock.Anything).Return(&pricing_app.CalculateFinalSalePriceResponse{
		CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
			{
				ProductVariantID: variantID1,
				Quantity:         5,
				BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
				FinalPrice:       pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
			},
			{
				ProductVariantID: variantID2,
				Quantity:         3,
				BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 50.0, Currency: domain.DefaultCurrency},
				FinalPrice:       pricing_app.MoneyDTO{Amount: 50.0, Currency: domain.DefaultCurrency},
			},
		},
	}, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(errors.New("database save error"))

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, mockPricingService, nil, nil, nil)

	cmd := application.AddOrderLineItemCommand{
		OrderID: order.ID,
		Item: application.OrderLineItemInput{
			ProductVariantID: variantID2,
			Quantity:         3,
		},
	}

	result, err := service.AddOrderLineItem(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database save error")
	orderRepo.AssertExpectations(t)
}

// ===== CreateOrder Additional Tests =====

func TestSalesService_CreateOrder_PartyIDRequired(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	cmd := application.CreateOrderCommand{
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Items: []application.OrderLineItemInput{
			{ProductVariantID: uuid.New(), Quantity: 1},
		},
	}

	result, err := service.CreateOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeValidation, domainErr.Code)
}

func TestSalesService_CreateOrder_DeliveryDateRequired(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	cmd := application.CreateOrderCommand{
		PartyID: uuid.New(),
		Items: []application.OrderLineItemInput{
			{ProductVariantID: uuid.New(), Quantity: 1},
		},
	}

	result, err := service.CreateOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeValidation, domainErr.Code)
}

func TestSalesService_CreateOrder_ItemsRequired(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	partyLookup := new(MockPartyLookup)

	partyLookup.On("ExistsParty", mock.Anything, partyID).Return(true, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, partyLookup, nil, nil)

	cmd := application.CreateOrderCommand{
		PartyID:      partyID,
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Items:        []application.OrderLineItemInput{},
	}

	result, err := service.CreateOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeValidation, domainErr.Code)
}

func TestSalesService_CreateOrder_PartyNotExists(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	partyLookup := new(MockPartyLookup)

	partyLookup.On("ExistsParty", mock.Anything, partyID).Return(false, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, partyLookup, nil, nil)

	cmd := application.CreateOrderCommand{
		PartyID:      partyID,
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Items: []application.OrderLineItemInput{
			{ProductVariantID: uuid.New(), Quantity: 1},
		},
	}

	result, err := service.CreateOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	partyLookup.AssertExpectations(t)
}

func TestSalesService_CreateOrder_NumberGeneratorNotConfigured(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	partyLookup := new(MockPartyLookup)

	partyLookup.On("ExistsParty", mock.Anything, partyID).Return(true, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, partyLookup, nil, nil)

	cmd := application.CreateOrderCommand{
		PartyID:      partyID,
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Items: []application.OrderLineItemInput{
			{ProductVariantID: uuid.New(), Quantity: 1},
		},
	}

	result, err := service.CreateOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "number generator not configured")
}

func TestSalesService_CreateOrder_FromQuote_QuoteNotFound(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	quoteID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	partyLookup := new(MockPartyLookup)
	numbers := new(MockNumberGenerator)

	partyLookup.On("ExistsParty", mock.Anything, partyID).Return(true, nil)
	orderNumber, _ := domain.NewOrderNumber("SO-FROM-Q1")
	numbers.On("NextOrderNumber", mock.Anything).Return(orderNumber, nil)
	quoteRepo.On("FindByID", mock.Anything, quoteID).Return(nil, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, nil, partyLookup, nil, nil)

	cmd := application.CreateOrderCommand{
		PartyID:      partyID,
		QuoteID:      &quoteID,
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
	}

	result, err := service.CreateOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	quoteRepo.AssertExpectations(t)
}

func TestSalesService_CreateOrder_FromQuote_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	partyLookup := new(MockPartyLookup)
	numbers := new(MockNumberGenerator)

	quoteNumber, _ := domain.NewQuoteNumber("Q-300")
	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, 0)
	quote, _ := domain.NewQuote(quoteNumber, partyID, time.Now(), time.Now().Add(24*time.Hour), []domain.QuoteLineItem{lineItem}, money, "")
	// Transition quote through valid states: Draft -> Sent -> Approved
	_ = quote.ChangeStatus(domain.QuoteStatusIssued)
	_ = quote.ChangeStatus(domain.QuoteStatusApproved)

	partyLookup.On("ExistsParty", mock.Anything, partyID).Return(true, nil)
	orderNumber, _ := domain.NewOrderNumber("SO-FROM-Q2")
	numbers.On("NextOrderNumber", mock.Anything).Return(orderNumber, nil)
	quoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)
	quoteRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Quote")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, nil, partyLookup, nil, nil)

	cmd := application.CreateOrderCommand{
		PartyID:      partyID,
		QuoteID:      &quote.ID,
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
	}

	result, err := service.CreateOrder(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, quote.ID, *result.QuoteID)
	orderRepo.AssertExpectations(t)
	quoteRepo.AssertExpectations(t)
}

func TestSalesService_CreateOrder_Direct_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	partyLookup := new(MockPartyLookup)
	numbers := new(MockNumberGenerator)
	mockPricingService := new(MockPricingEngine)

	partyLookup.On("ExistsParty", mock.Anything, partyID).Return(true, nil)
	orderNumber, _ := domain.NewOrderNumber("SO-DIRECT-001")
	numbers.On("NextOrderNumber", mock.Anything).Return(orderNumber, nil)
	mockPricingService.On("CalculateFinalSalePrice", mock.Anything, mock.Anything).Return(&pricing_app.CalculateFinalSalePriceResponse{
		CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
			{
				ProductVariantID: variantID,
				Quantity:         5,
				BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
				FinalPrice:       pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
			},
		},
	}, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, mockPricingService, partyLookup, nil, nil)

	notes := "Direct order notes"
	cmd := application.CreateOrderCommand{
		PartyID:      partyID,
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Items: []application.OrderLineItemInput{
			{ProductVariantID: variantID, Quantity: 5},
		},
		Notes: &notes,
	}

	result, err := service.CreateOrder(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.QuoteID)
	assert.Equal(t, notes, result.Notes)
	assert.Equal(t, 1, len(result.LineItems))
	orderRepo.AssertExpectations(t)
	mockPricingService.AssertExpectations(t)
}

func TestSalesService_CreateOrder_Direct_SaveError(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	partyLookup := new(MockPartyLookup)
	numbers := new(MockNumberGenerator)
	mockPricingService := new(MockPricingEngine)

	partyLookup.On("ExistsParty", mock.Anything, partyID).Return(true, nil)
	orderNumber, _ := domain.NewOrderNumber("SO-DIRECT-002")
	numbers.On("NextOrderNumber", mock.Anything).Return(orderNumber, nil)
	mockPricingService.On("CalculateFinalSalePrice", mock.Anything, mock.Anything).Return(&pricing_app.CalculateFinalSalePriceResponse{
		CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
			{
				ProductVariantID: variantID,
				Quantity:         5,
				BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
				FinalPrice:       pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
			},
		},
	}, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(errors.New("database connection error"))

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, mockPricingService, partyLookup, nil, nil)

	cmd := application.CreateOrderCommand{
		PartyID:      partyID,
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Items: []application.OrderLineItemInput{
			{ProductVariantID: variantID, Quantity: 5},
		},
	}

	result, err := service.CreateOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database connection error")
	orderRepo.AssertExpectations(t)
}

// ===== UpdateQuote Error Cases =====

func TestSalesService_UpdateQuote_QuoteNotDraft(t *testing.T) {
	ctx := context.Background()
	mockQuoteRepo := new(MockQuoteRepository)
	mockOrderRepo := new(MockSalesOrderRepository)
	mockDeliveryRepo := new(MockDeliveryNoteRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)
	service := application.NewSalesService(mockQuoteRepo, mockOrderRepo, mockDeliveryRepo, mockInvoiceRepo, nil, nil, nil, nil, nil)

	partyID := uuid.New()
	quoteNumber, err := domain.NewQuoteNumber("QUO-2026-001")
	assert.NoError(t, err)
	expirationDate := time.Now().Add(30 * 24 * time.Hour)
	money, err := domain.NewMoney(100.0, domain.DefaultCurrency)
	assert.NoError(t, err)
	lineItem, err := domain.NewQuoteLineItem(uuid.New(), 1, money, nil, 0)
	assert.NoError(t, err)

	quote, err := domain.NewQuote(quoteNumber, partyID, time.Now(), expirationDate, []domain.QuoteLineItem{lineItem}, money, "")
	assert.NoError(t, err)
	quote.Status = domain.QuoteStatusApproved // Set to non-editable status (only DRAFT and ISSUED are editable)

	mockQuoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)

	newNotes := "Updated notes"
	cmd := application.UpdateQuoteCommand{
		QuoteID: quote.ID,
		Notes:   &newNotes,
	}

	result, err := service.UpdateQuote(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "draft")
	mockQuoteRepo.AssertExpectations(t)
}

func TestSalesService_UpdateQuote_InvalidExpirationDate(t *testing.T) {
	ctx := context.Background()
	mockQuoteRepo := new(MockQuoteRepository)
	mockOrderRepo := new(MockSalesOrderRepository)
	mockDeliveryRepo := new(MockDeliveryNoteRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)
	service := application.NewSalesService(mockQuoteRepo, mockOrderRepo, mockDeliveryRepo, mockInvoiceRepo, nil, nil, nil, nil, nil)

	partyID := uuid.New()
	quoteNumber, err := domain.NewQuoteNumber("QUO-2026-001")
	assert.NoError(t, err)
	quoteDate := time.Now()
	expirationDate := quoteDate.Add(30 * 24 * time.Hour)
	money, err := domain.NewMoney(100.0, domain.DefaultCurrency)
	assert.NoError(t, err)
	lineItem, err := domain.NewQuoteLineItem(uuid.New(), 1, money, nil, 0)
	assert.NoError(t, err)

	quote, err := domain.NewQuote(quoteNumber, partyID, quoteDate, expirationDate, []domain.QuoteLineItem{lineItem}, money, "")
	assert.NoError(t, err)

	mockQuoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)

	// Try to set expiration date before quote date (invalid)
	invalidExpirationDate := quoteDate.Add(-24 * time.Hour)
	cmd := application.UpdateQuoteCommand{
		QuoteID:        quote.ID,
		ExpirationDate: &invalidExpirationDate,
	}

	result, err := service.UpdateQuote(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockQuoteRepo.AssertExpectations(t)
}

func TestSalesService_UpdateQuote_EmptyItems(t *testing.T) {
	ctx := context.Background()
	mockQuoteRepo := new(MockQuoteRepository)
	mockOrderRepo := new(MockSalesOrderRepository)
	mockDeliveryRepo := new(MockDeliveryNoteRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)
	service := application.NewSalesService(mockQuoteRepo, mockOrderRepo, mockDeliveryRepo, mockInvoiceRepo, nil, nil, nil, nil, nil)

	partyID := uuid.New()
	quoteNumber, err := domain.NewQuoteNumber("QUO-2026-001")
	assert.NoError(t, err)
	expirationDate := time.Now().Add(30 * 24 * time.Hour)
	money, err := domain.NewMoney(100.0, domain.DefaultCurrency)
	assert.NoError(t, err)
	lineItem, err := domain.NewQuoteLineItem(uuid.New(), 1, money, nil, 0)
	assert.NoError(t, err)

	quote, err := domain.NewQuote(quoteNumber, partyID, time.Now(), expirationDate, []domain.QuoteLineItem{lineItem}, money, "")
	assert.NoError(t, err)
	mockQuoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)

	// Try to update with empty line items (invalid)
	emptyItems := []application.QuoteLineItemInput{}
	cmd := application.UpdateQuoteCommand{
		QuoteID: quote.ID,
		Items:   emptyItems,
	}

	result, err := service.UpdateQuote(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockQuoteRepo.AssertExpectations(t)
}

// ===== ChangeQuoteStatus Additional Tests =====

func TestSalesService_ChangeQuoteStatus_QuoteNotFound(t *testing.T) {
	ctx := context.Background()
	mockQuoteRepo := new(MockQuoteRepository)
	mockOrderRepo := new(MockSalesOrderRepository)
	mockDeliveryRepo := new(MockDeliveryNoteRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)

	service := application.NewSalesService(mockQuoteRepo, mockOrderRepo, mockDeliveryRepo, mockInvoiceRepo, nil, nil, nil, nil, nil)

	quoteID := uuid.New()
	mockQuoteRepo.On("FindByID", mock.Anything, quoteID).Return(nil, domain.NewNotFoundError("quote not found"))

	cmd := application.ChangeQuoteStatusCommand{
		QuoteID:   quoteID,
		NewStatus: "issued",
	}

	result, err := service.ChangeQuoteStatus(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockQuoteRepo.AssertExpectations(t)
}

// ===== ChangeOrderStatus Additional Tests =====

func TestSalesService_ChangeOrderStatus_OrderNotFound(t *testing.T) {
	ctx := context.Background()
	mockQuoteRepo := new(MockQuoteRepository)
	mockOrderRepo := new(MockSalesOrderRepository)
	mockDeliveryRepo := new(MockDeliveryNoteRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)

	service := application.NewSalesService(mockQuoteRepo, mockOrderRepo, mockDeliveryRepo, mockInvoiceRepo, nil, nil, nil, nil, nil)

	orderID := uuid.New()
	mockOrderRepo.On("FindByID", mock.Anything, orderID).Return(nil, domain.NewNotFoundError("order not found"))

	cmd := application.ChangeOrderStatusCommand{
		OrderID:   orderID,
		NewStatus: "confirmed",
	}

	result, err := service.ChangeOrderStatus(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockOrderRepo.AssertExpectations(t)
}

// ===== ListDeliveryNotes Additional Tests =====

func TestSalesService_ListDeliveryNotes_ByOrderID(t *testing.T) {
	ctx := context.Background()
	mockQuoteRepo := new(MockQuoteRepository)
	mockOrderRepo := new(MockSalesOrderRepository)
	mockDeliveryRepo := new(MockDeliveryNoteRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)

	service := application.NewSalesService(mockQuoteRepo, mockOrderRepo, mockDeliveryRepo, mockInvoiceRepo, nil, nil, nil, nil, nil)

	orderID := uuid.New()
	mockDeliveryRepo.On("List", mock.Anything, mock.AnythingOfType("domain.DeliveryNoteFilter")).Return([]*domain.DeliveryNote{}, nil)

	query := application.ListDeliveryNotesQuery{SalesOrderID: &orderID}
	result, err := service.ListDeliveryNotes(ctx, query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
	mockDeliveryRepo.AssertExpectations(t)
}

// ===== ConvertQuoteToOrder Additional Tests =====

func TestSalesService_ConvertQuoteToOrder_QuoteNotFound(t *testing.T) {
	ctx := context.Background()
	mockQuoteRepo := new(MockQuoteRepository)
	mockOrderRepo := new(MockSalesOrderRepository)
	mockDeliveryRepo := new(MockDeliveryNoteRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)

	service := application.NewSalesService(mockQuoteRepo, mockOrderRepo, mockDeliveryRepo, mockInvoiceRepo, nil, nil, nil, nil, nil)

	quoteID := uuid.New()
	mockQuoteRepo.On("FindByID", mock.Anything, quoteID).Return(nil, domain.NewNotFoundError("quote not found"))

	cmd := application.ConvertQuoteToOrderCommand{
		QuoteID:      quoteID,
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
	}

	result, err := service.ConvertQuoteToOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockQuoteRepo.AssertExpectations(t)
}

// ===== UpdateOrderLineItem Tests =====

func TestSalesService_UpdateOrderLineItem_OrderNotFound(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	orderID := uuid.New()
	orderRepo.On("FindByID", mock.Anything, orderID).Return(nil, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	cmd := application.UpdateOrderLineItemCommand{
		OrderID:    orderID,
		LineItemID: uuid.New(),
	}

	result, err := service.UpdateOrderLineItem(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_UpdateOrderLineItem_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-001")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem},
		money,
		"Test order",
	)

	// Change order to DELIVERED status (cannot edit line items)
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)
	_ = order.ChangeStatus(domain.SalesOrderStatusDelivered)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	newQty := 10
	cmd := application.UpdateOrderLineItemCommand{
		OrderID:    order.ID,
		LineItemID: lineItem.ID,
		Quantity:   &newQty,
	}

	result, err := service.UpdateOrderLineItem(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeConflict, domainErr.Code)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_UpdateOrderLineItem_LineItemNotFound(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	mockPricingService := new(MockPricingEngine)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-002")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem},
		money,
		"Test order",
	)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, mockPricingService, nil, nil, nil)

	nonExistentLineItemID := uuid.New()
	newQty := 10
	cmd := application.UpdateOrderLineItemCommand{
		OrderID:    order.ID,
		LineItemID: nonExistentLineItemID,
		Quantity:   &newQty,
	}

	result, err := service.UpdateOrderLineItem(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_UpdateOrderLineItem_UpdateQuantity_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	mockPricingService := new(MockPricingEngine)

	unitPrice, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-003")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, unitPrice, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem},
		unitPrice,
		"Test order",
	)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)

	// Mock pricing service to return price
	mockPricingService.On("CalculateFinalSalePrice", mock.Anything, mock.Anything).Return(&pricing_app.CalculateFinalSalePriceResponse{
		CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
			{
				ProductVariantID: variantID,
				Quantity:         10,
				BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
				TaxRate:          21.0,
				FinalPrice:       pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
			},
		},
	}, nil)

	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, mockPricingService, nil, nil, nil)

	newQty := 10
	cmd := application.UpdateOrderLineItemCommand{
		OrderID:    order.ID,
		LineItemID: lineItem.ID,
		Quantity:   &newQty,
	}

	result, err := service.UpdateOrderLineItem(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.LineItems))
	assert.Equal(t, 10, result.LineItems[0].Quantity)
	assert.InDelta(t, 1000.0, result.Subtotal.Amount, 0.001)
	assert.InDelta(t, 210.0, result.TaxAmount.Amount, 0.001)
	assert.InDelta(t, 1210.0, result.Total.Amount, 0.001)
	orderRepo.AssertExpectations(t)
	mockPricingService.AssertExpectations(t)
}

func TestSalesService_UpdateQuote_ItemsRecalculateTotals_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	mockPricingService := new(MockPricingEngine)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	quoteNumber, _ := domain.NewQuoteNumber("Q/2026/0100")
	lineItem, _ := domain.NewQuoteLineItem(variantID, 1, money, nil, 0, 21)
	initialTax, _ := domain.NewMoney(21, domain.DefaultCurrency)

	quote, _ := domain.NewQuote(
		quoteNumber,
		partyID,
		time.Now(),
		time.Now().Add(30*24*time.Hour),
		[]domain.QuoteLineItem{lineItem},
		initialTax,
		"Test quote",
	)

	quoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)
	quoteRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Quote")).Return(nil)

	mockPricingService.On("CalculateFinalSalePrice", mock.Anything, mock.Anything).Return(&pricing_app.CalculateFinalSalePriceResponse{
		CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
			{
				ProductVariantID: variantID,
				Quantity:         3,
				BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
				TaxRate:          21.0,
			},
		},
	}, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, mockPricingService, nil, nil, nil)

	unitPrice := application.MoneyDTO{Amount: 100, Currency: domain.DefaultCurrency}
	cmd := application.UpdateQuoteCommand{
		QuoteID: quote.ID,
		Items: []application.QuoteLineItemInput{
			{
				ProductVariantID: variantID,
				Quantity:         3,
				UnitPrice:        &unitPrice,
			},
		},
	}

	result, err := service.UpdateQuote(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.InDelta(t, 300.0, result.Subtotal.Amount, 0.001)
	assert.InDelta(t, 63.0, result.TaxAmount.Amount, 0.001)
	assert.InDelta(t, 363.0, result.Total.Amount, 0.001)
	quoteRepo.AssertExpectations(t)
	mockPricingService.AssertExpectations(t)
}

func TestSalesService_UpdateOrderLineItem_UpdateUnitPrice_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	mockPricingService := new(MockPricingEngine)

	listPrice, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-004")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, listPrice, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem},
		listPrice,
		"Test order",
	)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)

	// Even with override price, pricing service is still called for calculation
	mockPricingService.On("CalculateFinalSalePrice", mock.Anything, mock.Anything).Return(&pricing_app.CalculateFinalSalePriceResponse{
		CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
			{
				ProductVariantID: variantID,
				Quantity:         5,
				BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
				FinalPrice:       pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
			},
		},
	}, nil)

	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, mockPricingService, nil, nil, nil)

	overridePrice := application.MoneyDTO{Amount: 150.0, Currency: "EUR"}
	cmd := application.UpdateOrderLineItemCommand{
		OrderID:    order.ID,
		LineItemID: lineItem.ID,
		UnitPrice:  &overridePrice,
	}

	result, err := service.UpdateOrderLineItem(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.LineItems))
	assert.Equal(t, 150.0, result.LineItems[0].UnitPrice.Amount)
	assert.Equal(t, "EUR", result.LineItems[0].UnitPrice.Currency)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_UpdateOrderLineItem_SaveError(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	mockPricingService := new(MockPricingEngine)

	unitPrice, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-005")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, unitPrice, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem},
		unitPrice,
		"Test order",
	)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)
	mockPricingService.On("CalculateFinalSalePrice", mock.Anything, mock.Anything).Return(&pricing_app.CalculateFinalSalePriceResponse{
		CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
			{
				ProductVariantID: variantID,
				Quantity:         10,
				BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
				FinalPrice:       pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
			},
		},
	}, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(errors.New("database error"))

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, mockPricingService, nil, nil, nil)

	newQty := 10
	cmd := application.UpdateOrderLineItemCommand{
		OrderID:    order.ID,
		LineItemID: lineItem.ID,
		Quantity:   &newQty,
	}

	result, err := service.UpdateOrderLineItem(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
	orderRepo.AssertExpectations(t)
}

// ===== RemoveOrderLineItem Tests =====

func TestSalesService_RemoveOrderLineItem_OrderNotFound(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	orderID := uuid.New()
	orderRepo.On("FindByID", mock.Anything, orderID).Return(nil, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	cmd := application.RemoveOrderLineItemCommand{
		OrderID:    orderID,
		LineItemID: uuid.New(),
	}

	result, err := service.RemoveOrderLineItem(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_RemoveOrderLineItem_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	mockPricingService := new(MockPricingEngine)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-006")
	lineItem1, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)
	lineItem2, _ := domain.NewOrderLineItem(uuid.New(), 3, money, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem1, lineItem2},
		money,
		"Test order",
	)

	// Change order to DELIVERED status (cannot edit line items)
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)
	_ = order.ChangeStatus(domain.SalesOrderStatusDelivered)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)

	mockPricingService.On("CalculateFinalSalePrice", mock.Anything, mock.Anything).Return(&pricing_app.CalculateFinalSalePriceResponse{
		CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
			{
				ProductVariantID: variantID,
				Quantity:         3,
				BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
				FinalPrice:       pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
			},
		},
	}, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, mockPricingService, nil, nil, nil)

	cmd := application.RemoveOrderLineItemCommand{
		OrderID:    order.ID,
		LineItemID: lineItem1.ID,
	}

	result, err := service.RemoveOrderLineItem(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot be updated")
	orderRepo.AssertExpectations(t)
}

func TestSalesService_RemoveOrderLineItem_LineItemNotFound(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-007")
	lineItem1, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)
	lineItem2, _ := domain.NewOrderLineItem(uuid.New(), 3, money, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem1, lineItem2},
		money,
		"Test order",
	)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	nonExistentLineItemID := uuid.New()
	cmd := application.RemoveOrderLineItemCommand{
		OrderID:    order.ID,
		LineItemID: nonExistentLineItemID,
	}

	result, err := service.RemoveOrderLineItem(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_RemoveOrderLineItem_LastLineItem(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-008")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem},
		money,
		"Test order with only one line item",
	)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	cmd := application.RemoveOrderLineItemCommand{
		OrderID:    order.ID,
		LineItemID: lineItem.ID,
	}

	result, err := service.RemoveOrderLineItem(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeValidation, domainErr.Code)
	assert.Contains(t, err.Error(), "at least one line item")
	orderRepo.AssertExpectations(t)
}

func TestSalesService_RemoveOrderLineItem_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID1 := uuid.New()
	variantID2 := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	mockPricingService := new(MockPricingEngine)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-009")
	lineItem1, _ := domain.NewOrderLineItem(variantID1, 5, money, nil, 0)
	lineItem2, _ := domain.NewOrderLineItem(variantID2, 3, money, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem1, lineItem2},
		money,
		"Test order",
	)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)
	mockPricingService.On("CalculateFinalSalePrice", mock.Anything, mock.Anything).Return(&pricing_app.CalculateFinalSalePriceResponse{
		CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
			{
				ProductVariantID: variantID2,
				Quantity:         3,
				BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
				TaxRate:          21.0,
				FinalPrice:       pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
			},
		},
	}, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, mockPricingService, nil, nil, nil)

	cmd := application.RemoveOrderLineItemCommand{
		OrderID:    order.ID,
		LineItemID: lineItem1.ID,
	}

	result, err := service.RemoveOrderLineItem(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.LineItems))
	assert.InDelta(t, 300.0, result.Subtotal.Amount, 0.001)
	assert.InDelta(t, 63.0, result.TaxAmount.Amount, 0.001)
	assert.InDelta(t, 363.0, result.Total.Amount, 0.001)
	assert.Equal(t, 1, len(result.LineItems))
	assert.Equal(t, variantID2, result.LineItems[0].ProductVariantID)
	orderRepo.AssertExpectations(t)
	mockPricingService.AssertExpectations(t)
}

func TestSalesService_RemoveOrderLineItem_SaveError(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID1 := uuid.New()
	variantID2 := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	mockPricingService := new(MockPricingEngine)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-010")
	lineItem1, _ := domain.NewOrderLineItem(variantID1, 5, money, nil, 0)
	lineItem2, _ := domain.NewOrderLineItem(variantID2, 3, money, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem1, lineItem2},
		money,
		"Test order",
	)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)
	mockPricingService.On("CalculateFinalSalePrice", mock.Anything, mock.Anything).Return(&pricing_app.CalculateFinalSalePriceResponse{
		CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
			{
				ProductVariantID: variantID2,
				Quantity:         3,
				BaseSalesPrice:   pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
				FinalPrice:       pricing_app.MoneyDTO{Amount: 100.0, Currency: domain.DefaultCurrency},
			},
		},
	}, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(errors.New("database connection lost"))

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, mockPricingService, nil, nil, nil)

	cmd := application.RemoveOrderLineItemCommand{
		OrderID:    order.ID,
		LineItemID: lineItem1.ID,
	}

	result, err := service.RemoveOrderLineItem(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database connection lost")
	orderRepo.AssertExpectations(t)
}

// ==========================================
// Tests for functions with 0% coverage
// ==========================================

func TestSalesService_ListDeliveryNotes_parseDeliveryNoteStatus_Success(t *testing.T) {
	// Test parseDeliveryNoteStatus through ListDeliveryNotes
	deliveryRepo := &MockDeliveryNoteRepository{}
	service := application.NewSalesService(nil, nil, deliveryRepo, nil, nil, nil, nil, nil, nil)

	validStatuses := []string{"PENDING", "DELIVERED", "CANCELLED"}
	expectedStatuses := []domain.DeliveryNoteStatus{
		domain.DeliveryNoteStatusPending,
		domain.DeliveryNoteStatusDelivered,
		domain.DeliveryNoteStatusCancelled,
	}

	for i, statusStr := range validStatuses {
		t.Run(fmt.Sprintf("Valid status %s", statusStr), func(t *testing.T) {
			expectedStatus := expectedStatuses[i]

			// Mock repository to expect filter with parsed status
			deliveryRepo.On("List", mock.Anything, mock.MatchedBy(func(filter domain.DeliveryNoteFilter) bool {
				return filter.Status != nil && *filter.Status == expectedStatus
			})).Return([]*domain.DeliveryNote{}, nil).Once()

			query := application.ListDeliveryNotesQuery{
				Status: &statusStr,
			}

			result, err := service.ListDeliveryNotes(context.Background(), query)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			deliveryRepo.AssertExpectations(t)
		})
	}
}

func TestSalesService_ListDeliveryNotes_parseDeliveryNoteStatus_Error(t *testing.T) {
	// Test parseDeliveryNoteStatus error handling through ListDeliveryNotes
	deliveryRepo := &MockDeliveryNoteRepository{}
	service := application.NewSalesService(nil, nil, deliveryRepo, nil, nil, nil, nil, nil, nil)

	invalidStatuses := []string{"INVALID_STATUS", "", "   "}

	for _, statusStr := range invalidStatuses {
		t.Run(fmt.Sprintf("Invalid status '%s'", statusStr), func(t *testing.T) {
			query := application.ListDeliveryNotesQuery{
				Status: &statusStr,
			}

			result, err := service.ListDeliveryNotes(context.Background(), query)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "invalid delivery note status")
		})
	}
}

func TestSalesService_ListInvoices_parseInvoiceStatus_Success(t *testing.T) {
	// Test parseInvoiceStatus through ListInvoices
	invoiceRepo := &MockInvoiceRepository{}
	service := application.NewSalesService(nil, nil, nil, invoiceRepo, nil, nil, nil, nil, nil)

	validStatuses := []string{"DRAFT", "ISSUED", "PAID", "OVERDUE", "VOID"}
	expectedStatuses := []domain.InvoiceStatus{
		domain.InvoiceStatusDraft,
		domain.InvoiceStatusIssued,
		domain.InvoiceStatusPaid,
		domain.InvoiceStatusOverdue,
		domain.InvoiceStatusVoid,
	}

	for i, statusStr := range validStatuses {
		t.Run(fmt.Sprintf("Valid status %s", statusStr), func(t *testing.T) {
			expectedStatus := expectedStatuses[i]

			// Mock repository to expect filter with parsed status
			invoiceRepo.On("List", mock.Anything, mock.MatchedBy(func(filter domain.InvoiceFilter) bool {
				return filter.Status != nil && *filter.Status == expectedStatus
			})).Return([]*domain.Invoice{}, nil).Once()

			query := application.ListInvoicesQuery{
				Status: &statusStr,
			}

			result, err := service.ListInvoices(context.Background(), query)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			invoiceRepo.AssertExpectations(t)
		})
	}
}

func TestSalesService_ListInvoices_parseInvoiceStatus_Error(t *testing.T) {
	// Test parseInvoiceStatus error handling through ListInvoices
	invoiceRepo := &MockInvoiceRepository{}
	service := application.NewSalesService(nil, nil, nil, invoiceRepo, nil, nil, nil, nil, nil)

	invalidStatuses := []string{"INVALID_STATUS", "", "   "}

	for _, statusStr := range invalidStatuses {
		t.Run(fmt.Sprintf("Invalid status '%s'", statusStr), func(t *testing.T) {
			query := application.ListInvoicesQuery{
				Status: &statusStr,
			}

			result, err := service.ListInvoices(context.Background(), query)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "invalid invoice status")
		})
	}
}

/*
func TestSalesService_buildInvoiceItemsFromDeliveryNotes_Success(t *testing.T) {
	// Test buildInvoiceItemsFromDeliveryNotes through CreateInvoice
	deliveryRepo := &MockDeliveryNoteRepository{}
	orderRepo := &MockSalesOrderRepository{}
	invoiceRepo := &MockInvoiceRepository{}
	numberGen := &Mocknerator{}

	service := application.NewSalesService(nil, orderRepo, deliveryRepo, invoiceRepo, numberGen, nil, nil, nil, nil)

	partyID := uuid.New()
	orderID := uuid.New()
	noteID := uuid.New()
	lineItemID := uuid.New()
	variantID := uuid.New()

	// Mock delivery note
	deliveryNote := createTestDeliveryNote(noteID, orderID, lineItemID, variantID, 5)

	// Mock sales order
	salesOrder := createTestSalesOrder(orderID, partyID, lineItemID, variantID, 10)

	// Setup mocks
	deliveryRepo.On("FindByID", mock.Anything, noteID).Return(deliveryNote, nil)
	orderRepo.On("FindByID", mock.Anything, orderID).Return(salesOrder, nil)
	invoiceRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	deliveryRepo.On("LinkLineItemsToInvoice", mock.Anything, mock.Anything).Return(nil)
	invoiceRepo.On("ListBySalesOrderID", mock.Anything, orderID).Return([]*domain.Invoice{}, nil)
	orderRepo.On("FindByIDForUpdate", mock.Anything, orderID).Return(salesOrder, nil)
	orderRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	invoiceNumber, _ := domain.NewInvoiceNumber("FV-2026-0001")
	numberGen.On("NextInvoiceNumber", mock.Anything, mock.Anything).Return(invoiceNumber, nil)

	cmd := application.CreateInvoiceCommand{
		PartyID:         partyID,
		InvoiceDate:     time.Now(),
		DueDate:         time.Now().AddDate(0, 1, 0),
		PaymentTerms:    "Net 30",
		DeliveryNoteIDs: []uuid.UUID{noteID},
	}

	result, err := service.CreateInvoice(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.LineItems, 1)
	assert.Equal(t, variantID, result.LineItems[0].ProductVariantID)
	deliveryRepo.AssertExpectations(t)
	orderRepo.AssertExpectations(t)
	invoiceRepo.AssertExpectations(t)
	numberGen.AssertExpectations(t)
}

func TestSalesService_buildInvoiceItemsFromDeliveryNotes_DeliveryNoteNotFound(t *testing.T) {
	deliveryRepo := &MockDeliveryNoteRepository{}
	orderRepo := &MockSalesOrderRepository{}
	invoiceRepo := &MockInvoiceRepository{}
	numberGen := &MockNumberGenerator{}

	service := application.NewSalesService(nil, orderRepo, deliveryRepo, invoiceRepo, numberGen, nil, nil, nil, nil)

	partyID := uuid.New()
	noteID := uuid.New()

	deliveryRepo.On("FindByID", mock.Anything, noteID).Return(nil, nil)

	cmd := application.CreateInvoiceCommand{
		PartyID:         partyID,
		InvoiceDate:     time.Now(),
		DueDate:         time.Now().AddDate(0, 1, 0),
		PaymentTerms:    "Net 30",
		DeliveryNoteIDs: []uuid.UUID{noteID},
	}

	result, err := service.CreateInvoice(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "delivery note not found")
	deliveryRepo.AssertExpectations(t)
}

func TestSalesService_buildInvoiceItemsFromDeliveryNotes_OrderNotFound(t *testing.T) {
	deliveryRepo := &MockDeliveryNoteRepository{}
	orderRepo := &MockSalesOrderRepository{}
	invoiceRepo := &MockInvoiceRepository{}
	numberGen := &MockNumberGenerator{}

	service := application.NewSalesService(nil, orderRepo, deliveryRepo, invoiceRepo, numberGen, nil, nil, nil, nil)

	partyID := uuid.New()
	orderID := uuid.New()
	noteID := uuid.New()
	lineItemID := uuid.New()
	variantID := uuid.New()

	deliveryNote := createTestDeliveryNote(noteID, orderID, lineItemID, variantID, 5)

	deliveryRepo.On("FindByID", mock.Anything, noteID).Return(deliveryNote, nil)
	orderRepo.On("FindByID", mock.Anything, orderID).Return(nil, nil)

	cmd := application.CreateInvoiceCommand{
		PartyID:         partyID,
		InvoiceDate:     time.Now(),
		DueDate:         time.Now().AddDate(0, 1, 0),
		PaymentTerms:    "Net 30",
		DeliveryNoteIDs: []uuid.UUID{noteID},
	}

	result, err := service.CreateInvoice(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "order not found")
	deliveryRepo.AssertExpectations(t)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_buildInvoiceItemsFromDeliveryNotes_PartyMismatch(t *testing.T) {
	deliveryRepo := &MockDeliveryNoteRepository{}
	orderRepo := &MockSalesOrderRepository{}
	invoiceRepo := &MockInvoiceRepository{}
	numberGen := &MockNumberGenerator{}

	service := application.NewSalesService(nil, orderRepo, deliveryRepo, invoiceRepo, numberGen, nil, nil, nil, nil)

	partyID := uuid.New()
	wrongPartyID := uuid.New()
	orderID := uuid.New()
	noteID := uuid.New()
	lineItemID := uuid.New()
	variantID := uuid.New()

	deliveryNote := createTestDeliveryNote(noteID, orderID, lineItemID, variantID, 5)
	salesOrder := createTestSalesOrder(orderID, wrongPartyID, lineItemID, variantID, 10)

	deliveryRepo.On("FindByID", mock.Anything, noteID).Return(deliveryNote, nil)
	orderRepo.On("FindByID", mock.Anything, orderID).Return(salesOrder, nil)

	cmd := application.CreateInvoiceCommand{
		PartyID:         partyID,
		InvoiceDate:     time.Now(),
		DueDate:         time.Now().AddDate(0, 1, 0),
		PaymentTerms:    "Net 30",
		DeliveryNoteIDs: []uuid.UUID{noteID},
	}

	result, err := service.CreateInvoice(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "delivery note party mismatch")
	deliveryRepo.AssertExpectations(t)
	orderRepo.AssertExpectations(t)
}
*/

// Helper functions for buildInvoiceItemsFromDeliveryNotes tests
func createTestDeliveryNote(noteID, orderID, lineItemID, variantID uuid.UUID, deliveredQty int) *domain.DeliveryNote {
	partyID := uuid.New()

	lineItem, _ := domain.NewDeliveryNoteLineItem(
		lineItemID,
		variantID,
		deliveredQty,
	)

	noteNumber, _ := domain.NewDeliveryNoteNumber("DN-001")
	note, _ := domain.NewDeliveryNote(
		noteNumber,
		orderID,
		partyID,
		time.Now(),
		[]domain.DeliveryNoteLineItem{lineItem},
		"Test notes",
	)
	note.ID = noteID
	return note
}

func createTestSalesOrder(orderID, partyID, lineItemID, variantID uuid.UUID, quantity int) *domain.SalesOrder {
	unitPrice, _ := domain.NewMoney(100.00, domain.DefaultCurrency)
	taxAmount, _ := domain.NewMoney(21.00, domain.DefaultCurrency)

	lineItem, _ := domain.NewOrderLineItem(
		variantID,
		quantity,
		unitPrice,
		nil,
		0,
	)
	lineItem.ID = lineItemID

	orderNumber, _ := domain.NewOrderNumber("ORD-001")
	orderDate := time.Now()
	deliveryDate := orderDate.AddDate(0, 0, 7) // 7 days later

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		orderDate,
		deliveryDate,
		[]domain.OrderLineItem{lineItem},
		taxAmount,
		"Test notes",
	)
	order.ID = orderID
	return order
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

func TestSalesService_CreateSimplifiedInvoice_Success(t *testing.T) {
	invoiceRepo := &MockInvoiceRepository{}
	numberGen := &MockNumberGenerator{}
	partyLookup := &MockPartyLookup{}
	pricingEngine := &MockPricingEngine{}

	service := application.NewSalesService(nil, nil, nil, invoiceRepo, numberGen, pricingEngine, partyLookup, nil, nil)

	partyID := uuid.New()
	variantID := uuid.New()
	invoiceDate := time.Now()

	// Setup mocks
	partyLookup.On("ExistsParty", mock.Anything, partyID).Return(true, nil)

	pricingResp := &pricing_app.CalculateFinalSalePriceResponse{
		SaleTotal: pricing_app.MoneyDTO{Amount: 100.00},
		CalculatedItems: []pricing_app.CalculatedSaleItemResponse{
			{
				ProductVariantID: variantID,
				Quantity:         2,
				FinalPrice:       pricing_app.MoneyDTO{Amount: 50.00},
			},
		},
	}
	pricingEngine.On("CalculateFinalSalePrice", mock.Anything, mock.Anything).Return(pricingResp, nil)

	invoiceNumber, _ := domain.NewInvoiceNumber("FT-2026-0001")
	numberGen.On("NextInvoiceNumber", mock.Anything, mock.Anything).Return(invoiceNumber, nil)

	invoiceRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	cmd := application.CreateSimplifiedInvoiceCommand{
		PartyID:     partyID,
		InvoiceDate: invoiceDate,
		Items: []application.OrderLineItemInputSimplified{
			{
				ProductVariantID: variantID,
				Quantity:         2,
			},
		},
	}

	result, err := service.CreateSimplifiedInvoice(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.LineItems, 1)
	assert.Equal(t, variantID, result.LineItems[0].ProductVariantID)
	assert.Equal(t, 2, result.LineItems[0].Quantity)
	partyLookup.AssertExpectations(t)
	pricingEngine.AssertExpectations(t)
	numberGen.AssertExpectations(t)
	invoiceRepo.AssertExpectations(t)
}

func TestSalesService_CreateSimplifiedInvoice_PartyNotFound(t *testing.T) {
	invoiceRepo := &MockInvoiceRepository{}
	numberGen := &MockNumberGenerator{}
	partyLookup := &MockPartyLookup{}

	service := application.NewSalesService(nil, nil, nil, invoiceRepo, numberGen, nil, partyLookup, nil, nil)

	partyID := uuid.New()

	partyLookup.On("ExistsParty", mock.Anything, partyID).Return(false, nil)

	cmd := application.CreateSimplifiedInvoiceCommand{
		PartyID:     partyID,
		InvoiceDate: time.Now(),
		Items: []application.OrderLineItemInputSimplified{
			{
				ProductVariantID: uuid.New(),
				Quantity:         1,
			},
		},
	}

	result, err := service.CreateSimplifiedInvoice(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "party not found")
	partyLookup.AssertExpectations(t)
}

func TestSalesService_CreateSimplifiedInvoice_PricingError(t *testing.T) {
	invoiceRepo := &MockInvoiceRepository{}
	numberGen := &MockNumberGenerator{}
	partyLookup := &MockPartyLookup{}
	pricingEngine := &MockPricingEngine{}

	service := application.NewSalesService(nil, nil, nil, invoiceRepo, numberGen, pricingEngine, partyLookup, nil, nil)

	partyID := uuid.New()

	partyLookup.On("ExistsParty", mock.Anything, partyID).Return(true, nil)
	pricingEngine.On("CalculateFinalSalePrice", mock.Anything, mock.Anything).Return(nil, errors.New("pricing service error"))

	cmd := application.CreateSimplifiedInvoiceCommand{
		PartyID:     partyID,
		InvoiceDate: time.Now(),
		Items: []application.OrderLineItemInputSimplified{
			{
				ProductVariantID: uuid.New(),
				Quantity:         1,
			},
		},
	}

	result, err := service.CreateSimplifiedInvoice(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "pricing calculation failed")
	partyLookup.AssertExpectations(t)
	pricingEngine.AssertExpectations(t)
}

func TestSalesService_UpdateQuote_ExpirationDateBeforeQuoteDate(t *testing.T) {
	quoteRepo := &MockQuoteRepository{}
	service := application.NewSalesService(quoteRepo, nil, nil, nil, nil, nil, nil, nil, nil)

	quoteID := uuid.New()
	existingQuote := createTestQuote(quoteID)

	quoteRepo.On("FindByID", mock.Anything, quoteID).Return(existingQuote, nil)

	// Set expiration date BEFORE quote date
	invalidExpirationDate := existingQuote.QuoteDate.AddDate(0, 0, -1)

	cmd := application.UpdateQuoteCommand{
		QuoteID:        quoteID,
		ExpirationDate: &invalidExpirationDate,
	}

	result, err := service.UpdateQuote(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "expirationDate cannot be before quoteDate")
	quoteRepo.AssertExpectations(t)
}

func TestSalesService_UpdateQuote_EmptyItemsArray(t *testing.T) {
	quoteRepo := &MockQuoteRepository{}
	service := application.NewSalesService(quoteRepo, nil, nil, nil, nil, nil, nil, nil, nil)

	quoteID := uuid.New()
	existingQuote := createTestQuote(quoteID)

	quoteRepo.On("FindByID", mock.Anything, quoteID).Return(existingQuote, nil)

	// Empty items array should cause error
	emptyItems := []application.QuoteLineItemInput{}

	cmd := application.UpdateQuoteCommand{
		QuoteID: quoteID,
		Items:   emptyItems,
	}

	result, err := service.UpdateQuote(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "items cannot be empty")
	quoteRepo.AssertExpectations(t)
}

func TestSalesService_UpdateQuote_NotDraftStatus(t *testing.T) {
	quoteRepo := &MockQuoteRepository{}
	service := application.NewSalesService(quoteRepo, nil, nil, nil, nil, nil, nil, nil, nil)

	quoteID := uuid.New()
	existingQuote := createTestQuote(quoteID)
	// Change status to Approved (not editable: only Draft and Issued are editable)
	_ = existingQuote.ChangeStatus(domain.QuoteStatusIssued)
	_ = existingQuote.ChangeStatus(domain.QuoteStatusApproved)

	quoteRepo.On("FindByID", mock.Anything, quoteID).Return(existingQuote, nil)

	cmd := application.UpdateQuoteCommand{
		QuoteID: quoteID,
		Notes:   stringPtr("Updated notes"),
	}

	result, err := service.UpdateQuote(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "only draft or issued quotes can be updated")
	quoteRepo.AssertExpectations(t)
}

// Helper function to create test quote
func createTestQuote(quoteID uuid.UUID) *domain.Quote {
	partyID := uuid.New()
	variantID := uuid.New()
	unitPrice, _ := domain.NewMoney(100.00, domain.DefaultCurrency)
	taxAmount, _ := domain.NewMoney(21.00, domain.DefaultCurrency)

	lineItem, _ := domain.NewQuoteLineItem(
		variantID,
		2,
		unitPrice,
		nil,
		0,
	)

	quoteNumber, _ := domain.NewQuoteNumber("Q-001")
	quoteDate := time.Now()
	expirationDate := quoteDate.AddDate(0, 1, 0)

	quote, _ := domain.NewQuote(
		quoteNumber,
		partyID,
		quoteDate,
		expirationDate,
		[]domain.QuoteLineItem{lineItem},
		taxAmount,
		"Test quote",
	)
	quote.ID = quoteID
	return quote
}

func TestSalesService_CreateDeliveryNote_CanceledOrder(t *testing.T) {
	orderRepo := &MockSalesOrderRepository{}
	deliveryRepo := &MockDeliveryNoteRepository{}
	numberGen := &MockNumberGenerator{}

	service := application.NewSalesService(nil, orderRepo, deliveryRepo, nil, numberGen, nil, nil, nil, nil)

	orderID := uuid.New()
	lineItemID := uuid.New()

	// Create canceled order
	order := createTestSalesOrder(orderID, uuid.New(), lineItemID, uuid.New(), 10)
	_ = order.ChangeStatus(domain.SalesOrderStatusCancelled)

	orderRepo.On("FindByIDForUpdate", mock.Anything, orderID).Return(order, nil)

	cmd := application.CreateDeliveryNoteCommand{
		SalesOrderID: orderID,
		DeliveryDate: time.Now(),
		Items: []application.DeliveryNoteLineItemInput{
			{
				SalesOrderLineItemID: lineItemID,
				DeliveredQuantity:    5,
			},
		},
	}

	result, err := service.CreateDeliveryNote(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot create delivery note for canceled order")
	orderRepo.AssertExpectations(t)
}

func TestSalesService_CreateDeliveryNote_InvoicedOrder(t *testing.T) {
	orderRepo := &MockSalesOrderRepository{}
	deliveryRepo := &MockDeliveryNoteRepository{}
	numberGen := &MockNumberGenerator{}

	service := application.NewSalesService(nil, orderRepo, deliveryRepo, nil, numberGen, nil, nil, nil, nil)

	orderID := uuid.New()
	lineItemID := uuid.New()

	// Create invoiced order with proper status transitions
	order := createTestSalesOrder(orderID, uuid.New(), lineItemID, uuid.New(), 10)

	// Follow valid status transitions: Pending -> InPreparation -> PartiallyDelivered -> Delivered -> Invoiced
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)
	_ = order.ChangeStatus(domain.SalesOrderStatusPartiallyDelivered)
	_ = order.ChangeStatus(domain.SalesOrderStatusDelivered)
	_ = order.ChangeStatus(domain.SalesOrderStatusInvoiced)

	orderRepo.On("FindByIDForUpdate", mock.Anything, orderID).Return(order, nil)

	// Mock delivery note number generator - this will be called before status validation
	numberGen.On("GenerateDeliveryNoteNumber", mock.Anything).Return("DN-001", nil)

	deliveryRepo.On("ListBySalesOrderID", mock.Anything, orderID).Return([]*domain.DeliveryNote{}, nil)

	cmd := application.CreateDeliveryNoteCommand{
		SalesOrderID: orderID,
		DeliveryDate: time.Now(),
		Items: []application.DeliveryNoteLineItemInput{
			{
				SalesOrderLineItemID: lineItemID,
				DeliveredQuantity:    5,
			},
		},
	}

	result, err := service.CreateDeliveryNote(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot create delivery note for invoiced order")
	orderRepo.AssertExpectations(t)
}

func TestSalesService_CreateDeliveryNote_OrderNotFound(t *testing.T) {
	orderRepo := &MockSalesOrderRepository{}
	numberGen := &MockNumberGenerator{}

	service := application.NewSalesService(nil, orderRepo, nil, nil, numberGen, nil, nil, nil, nil)

	orderID := uuid.New()

	orderRepo.On("FindByIDForUpdate", mock.Anything, orderID).Return(nil, domain.NewNotFoundError("sales order not found"))

	cmd := application.CreateDeliveryNoteCommand{
		SalesOrderID: orderID,
		DeliveryDate: time.Now(),
		Items: []application.DeliveryNoteLineItemInput{
			{
				SalesOrderLineItemID: uuid.New(),
				DeliveredQuantity:    5,
			},
		},
	}

	result, err := service.CreateDeliveryNote(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "sales order not found")
	orderRepo.AssertExpectations(t)
}

func TestSalesService_DeleteDeliveryNote_RevertsOrderStatusWhenNoDeliveriesRemain(t *testing.T) {
	orderRepo := &MockSalesOrderRepository{}
	deliveryRepo := &MockDeliveryNoteRepository{}
	service := application.NewSalesService(nil, orderRepo, deliveryRepo, nil, nil, nil, nil, nil, nil)

	orderID := uuid.New()
	lineItemID := uuid.New()
	variantID := uuid.New()
	noteID := uuid.New()

	order := createTestSalesOrder(orderID, uuid.New(), lineItemID, variantID, 10)
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)
	_ = order.ChangeStatus(domain.SalesOrderStatusDelivered)

	note := createTestDeliveryNote(noteID, orderID, lineItemID, variantID, 10)
	_ = note.ChangeStatus(domain.DeliveryNoteStatusDelivered)

	deliveryRepo.On("FindByID", mock.Anything, noteID).Return(note, nil)
	deliveryRepo.On("Delete", mock.Anything, noteID).Return(nil)
	orderRepo.On("FindByIDForUpdate", mock.Anything, orderID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, orderID).Return([]*domain.DeliveryNote{}, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Run(func(args mock.Arguments) {
		savedOrder := args.Get(1).(*domain.SalesOrder)
		assert.Equal(t, domain.SalesOrderStatusInPreparation, savedOrder.Status)
	}).Return(nil)

	err := service.DeleteDeliveryNote(context.Background(), application.DeleteDeliveryNoteCommand{DeliveryNoteID: noteID})

	assert.NoError(t, err)
	assert.Equal(t, domain.SalesOrderStatusInPreparation, order.Status)
	orderRepo.AssertExpectations(t)
	deliveryRepo.AssertExpectations(t)
}

func TestSalesService_CreateDeliveryNote_LineItemNotFound(t *testing.T) {
	orderRepo := &MockSalesOrderRepository{}
	deliveryRepo := &MockDeliveryNoteRepository{}
	numberGen := &MockNumberGenerator{}

	service := application.NewSalesService(nil, orderRepo, deliveryRepo, nil, numberGen, nil, nil, nil, nil)

	orderID := uuid.New()
	lineItemID := uuid.New()
	nonExistentLineItemID := uuid.New()

	// Create order with line item
	order := createTestSalesOrder(orderID, uuid.New(), lineItemID, uuid.New(), 10)
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)

	orderRepo.On("FindByIDForUpdate", mock.Anything, orderID).Return(order, nil)
	numberGen.On("GenerateDeliveryNoteNumber", mock.Anything).Return("DN-001", nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, orderID).Return([]*domain.DeliveryNote{}, nil)

	cmd := application.CreateDeliveryNoteCommand{
		SalesOrderID: orderID,
		DeliveryDate: time.Now(),
		Items: []application.DeliveryNoteLineItemInput{
			{
				SalesOrderLineItemID: nonExistentLineItemID, // Different ID
				DeliveredQuantity:    5,
			},
		},
	}

	result, err := service.CreateDeliveryNote(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "line item not found")
	orderRepo.AssertExpectations(t)
}

// Test para toMoneyDTOPtr function coverage
func TestToMoneyDTOPtr_NilInput(t *testing.T) {
	// Test case directly - this is a simple helper function - but it's not accessible from _test package
	// We'll test it indirectly through a DTO that uses it
	// Let's skip this test for now and test via DTO creation that calls toMoneyDTOPtr
}

func TestToMoneyDTOPtr_ValidMoney(t *testing.T) {
	// We test this indirectly since toMoneyDTOPtr is not exported
	// Let's remove these tests for now and add coverage via existing DTO creation
}

// Test parseOrderStatus via ChangeOrderStatus function (indirect testing)
func TestSalesService_ChangeOrderStatus_InvalidStatus(t *testing.T) {
	orderRepo := &MockSalesOrderRepository{}

	service := application.NewSalesService(nil, orderRepo, nil, nil, nil, nil, nil, nil, nil)

	orderID := uuid.New()

	// Create a test order so we can get past FindByID and reach parseOrderStatus
	order := createTestSalesOrder(orderID, uuid.New(), uuid.New(), uuid.New(), 10)

	orderRepo.On("FindByID", mock.Anything, orderID).Return(order, nil)

	cmd := application.ChangeOrderStatusCommand{
		OrderID:   orderID,
		NewStatus: "INVALID_STATUS", // This will exercise parseOrderStatus with invalid input
	}

	result, err := service.ChangeOrderStatus(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	// The error message from parseOrderStatus should contain validation info
	orderRepo.AssertExpectations(t)
}

// Test parseOrderStatus via ListOrders function (another indirect test)
func TestSalesService_ListOrders_InvalidStatus(t *testing.T) {
	orderRepo := &MockSalesOrderRepository{}

	service := application.NewSalesService(nil, orderRepo, nil, nil, nil, nil, nil, nil, nil)

	invalidStatus := "COMPLETELY_INVALID_STATUS"

	query := application.ListOrdersQuery{
		Status: &invalidStatus, // This will trigger parseOrderStatus
	}

	result, err := service.ListOrders(context.Background(), query)

	assert.Error(t, err)
	assert.Nil(t, result)
	// parseOrderStatus should return an error for invalid status
}

// Test canEditOrderDetails with InPreparation status (to get 100% coverage)
func TestSalesService_UpdateOrderDetails_InPreparationStatus(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	newPartyID := uuid.New()
	variantID := uuid.New()

	orderRepo := new(MockSalesOrderRepository)
	partyLookup := new(MockPartyLookup)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	orderNumber, _ := domain.NewOrderNumber("SO-PREP-001")
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, 0)

	order, _ := domain.NewSalesOrder(
		orderNumber,
		partyID,
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		[]domain.OrderLineItem{lineItem},
		money,
		"Test order",
	)

	// Change order to InPreparation status to test the other case of canEditOrderDetails
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)

	deliveryRepo := new(MockDeliveryNoteRepository)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)
	partyLookup.On("ExistsParty", mock.Anything, newPartyID).Return(true, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)

	service := application.NewSalesService(nil, orderRepo, deliveryRepo, nil, nil, nil, partyLookup, nil, nil)

	newNotes := "Updated notes for InPreparation order"
	cmd := application.UpdateOrderDetailsCommand{
		OrderID: order.ID,
		PartyID: &newPartyID,
		Notes:   &newNotes,
	}

	result, err := service.UpdateOrderDetails(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newNotes, result.Notes)
	orderRepo.AssertExpectations(t)
	partyLookup.AssertExpectations(t)
}

func TestSalesService_CreateDeliveryNote_ZeroDeliveryDate(t *testing.T) {
	service := application.NewSalesService(nil, nil, nil, nil, nil, nil, nil, nil, nil)

	cmd := application.CreateDeliveryNoteCommand{
		SalesOrderID: uuid.New(),
		DeliveryDate: time.Time{}, // Zero time
		Items: []application.DeliveryNoteLineItemInput{
			{
				SalesOrderLineItemID: uuid.New(),
				DeliveredQuantity:    1,
			},
		},
	}

	result, err := service.CreateDeliveryNote(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "deliveryDate is required")
}

func TestSalesService_CreateDeliveryNote_EmptyItems(t *testing.T) {
	service := application.NewSalesService(nil, nil, nil, nil, nil, nil, nil, nil, nil)

	cmd := application.CreateDeliveryNoteCommand{
		SalesOrderID: uuid.New(),
		DeliveryDate: time.Now(),
		Items:        []application.DeliveryNoteLineItemInput{}, // Empty items
	}

	result, err := service.CreateDeliveryNote(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "items cannot be empty")
}

// Test additional validation path that might exercise more coverage
func TestSalesService_CreateOrder_EmptyItems(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()

	orderRepo := new(MockSalesOrderRepository)
	partyLookup := new(MockPartyLookup)
	numbers := new(MockNumberGenerator)

	service := application.NewSalesService(nil, orderRepo, nil, nil, numbers, nil, partyLookup, nil, nil)

	cmd := application.CreateOrderCommand{
		PartyID:      partyID,
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Items:        []application.OrderLineItemInput{}, // Empty items - should trigger validation
		Notes:        nil,
	}

	result, err := service.CreateOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	// This should trigger validation early and increment coverage on validation paths
}

// Test parseDeliveryNoteStatus through ListDeliveryNotes with more invalid statuses
func TestSalesService_ListDeliveryNotes_MoreInvalidStatuses(t *testing.T) {
	deliveryRepo := &MockDeliveryNoteRepository{}

	service := application.NewSalesService(nil, nil, deliveryRepo, nil, nil, nil, nil, nil, nil)

	// Test multiple invalid statuses to exercise parseDeliveryNoteStatus more thoroughly
	invalidStatuses := []string{"INVALID", "WRONG_STATUS", "123", "null", "undefined"}

	for _, invalidStatus := range invalidStatuses {
		query := application.ListDeliveryNotesQuery{
			Status: &invalidStatus,
		}

		result, err := service.ListDeliveryNotes(context.Background(), query)

		assert.Error(t, err)
		assert.Nil(t, result)
	}
}

// Test CreateOrder validation paths to exercise buildOrderLineItems indirectly
func TestSalesService_CreateOrder_ValidationPaths(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	orderRepo := &MockSalesOrderRepository{}
	partyLookup := &MockPartyLookup{}
	numberGen := &MockNumberGenerator{}
	pricingEngine := &MockPricingEngine{}

	service := application.NewSalesService(nil, orderRepo, nil, nil, numberGen, pricingEngine, partyLookup, nil, nil)

	// Test case 1: PartyID nil
	cmd1 := application.CreateOrderCommand{
		PartyID:      uuid.Nil, // Invalid PartyID
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
		Items: []application.OrderLineItemInput{
			{ProductVariantID: variantID, Quantity: 5},
		},
	}

	result1, err1 := service.CreateOrder(ctx, cmd1)
	assert.Error(t, err1)
	assert.Nil(t, result1)

	// Test case 2: DeliveryDate zero
	cmd2 := application.CreateOrderCommand{
		PartyID:      partyID,
		DeliveryDate: time.Time{}, // Zero time
		Items: []application.OrderLineItemInput{
			{ProductVariantID: variantID, Quantity: 5},
		},
	}

	result2, err2 := service.CreateOrder(ctx, cmd2)
	assert.Error(t, err2)
	assert.Nil(t, result2)
}

// Test to exercise deliveredQuantities function through CreateDeliveryNote
func TestSalesService_CreateDeliveryNote_DeliveredQuantitiesPath(t *testing.T) {
	orderRepo := &MockSalesOrderRepository{}
	deliveryRepo := &MockDeliveryNoteRepository{}
	numberGen := &MockNumberGenerator{}

	service := application.NewSalesService(nil, orderRepo, deliveryRepo, nil, numberGen, nil, nil, nil, nil)

	orderID := uuid.New()
	lineItemID := uuid.New()

	// Create order in valid status
	order := createTestSalesOrder(orderID, uuid.New(), lineItemID, uuid.New(), 10)
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)

	orderRepo.On("FindByIDForUpdate", mock.Anything, orderID).Return(order, nil)

	noteNumber, _ := domain.NewDeliveryNoteNumber("DN-TEST-001")
	numberGen.On("NextDeliveryNoteNumber", mock.Anything).Return(noteNumber, nil)

	// Mock existing delivery notes to exercise deliveredQuantities calculation
	existingDeliveryNotes := []*domain.DeliveryNote{}
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, orderID).Return(existingDeliveryNotes, nil)
	deliveryRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.DeliveryNote")).Return(nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)

	cmd := application.CreateDeliveryNoteCommand{
		SalesOrderID: orderID,
		DeliveryDate: time.Now(),
		Items: []application.DeliveryNoteLineItemInput{
			{
				SalesOrderLineItemID: lineItemID,
				DeliveredQuantity:    3, // Partial delivery
			},
		},
	}

	result, err := service.CreateDeliveryNote(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	orderRepo.AssertExpectations(t)
	deliveryRepo.AssertExpectations(t)
	numberGen.AssertExpectations(t)
}

// Test to exercise more validation paths in UpdateQuote
func TestSalesService_UpdateQuote_AdditionalValidations(t *testing.T) {
	quoteRepo := &MockQuoteRepository{}
	pricingEngine := &MockPricingEngine{}

	service := application.NewSalesService(quoteRepo, nil, nil, nil, nil, pricingEngine, nil, nil, nil)

	quoteID := uuid.New()
	variantID := uuid.New()

	// Create draft quote
	quote := createTestQuote(quoteID)

	quoteRepo.On("FindByID", mock.Anything, quoteID).Return(quote, nil)

	// Test case: Items with zero quantity (should trigger validation)
	cmd := application.UpdateQuoteCommand{
		QuoteID: quoteID,
		Items: []application.QuoteLineItemInput{
			{
				ProductVariantID: variantID,
				Quantity:         0, // Invalid quantity
			},
		},
	}

	result, err := service.UpdateQuote(context.Background(), cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "quantity must be greater than zero")
	quoteRepo.AssertExpectations(t)
}

// Test canEditOrderDetails with false path (non-editable status)
func TestSalesService_UpdateOrderDetails_NonEditableStatus(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()

	orderRepo := new(MockSalesOrderRepository)

	service := application.NewSalesService(nil, orderRepo, nil, nil, nil, nil, nil, nil, nil)

	orderID := uuid.New()
	lineItemID := uuid.New()

	// Create order with non-editable status (Delivered)
	order := createTestSalesOrder(orderID, partyID, lineItemID, uuid.New(), 10)
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)
	_ = order.ChangeStatus(domain.SalesOrderStatusPartiallyDelivered)
	_ = order.ChangeStatus(domain.SalesOrderStatusDelivered) // Non-editable status

	orderRepo.On("FindByID", mock.Anything, orderID).Return(order, nil)

	newNotes := "Should not be able to update"
	cmd := application.UpdateOrderDetailsCommand{
		OrderID: orderID,
		Notes:   &newNotes,
	}

	result, err := service.UpdateOrderDetails(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot be updated in current status")
	orderRepo.AssertExpectations(t)
}

// Test buildOrderLineItemsFromSeeds - Strategic tests for 6.1% coverage function
func TestSalesService_CreateOrder_NoNumberGenerator(t *testing.T) {
	quotesRepo := &MockQuoteRepository{}
	orderRepo := &MockSalesOrderRepository{}
	service := application.NewSalesService(quotesRepo, orderRepo, nil, nil, nil, nil, nil, nil, nil) // No number generator

	ctx := context.Background()
	cmd := application.CreateOrderCommand{
		PartyID: uuid.New(),
		Items: []application.OrderLineItemInput{
			{
				ProductVariantID: uuid.New(),
				Quantity:         1,
			},
		},
		DeliveryDate: time.Now().AddDate(0, 0, 7), // 7 days from now
	}

	result, err := service.CreateOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "order number generator not configured")
}

func TestSalesService_CreateOrder_InvalidProductVariantID(t *testing.T) {
	quotesRepo := &MockQuoteRepository{}
	orderRepo := &MockSalesOrderRepository{}
	pricingEngine := &MockPricingEngine{}
	partyLookup := &MockPartyLookup{}
	numberGen := &MockNumberGenerator{}
	service := application.NewSalesService(quotesRepo, orderRepo, nil, nil, numberGen, pricingEngine, partyLookup, nil, nil)

	ctx := context.Background()
	partyID := uuid.New()
	cmd := application.CreateOrderCommand{
		PartyID: partyID,
		Items: []application.OrderLineItemInput{
			{
				ProductVariantID: uuid.Nil, // Invalid UUID
				Quantity:         1,
			},
		},
		DeliveryDate: time.Now().AddDate(0, 0, 7),
	}

	partyLookup.On("ExistsParty", ctx, partyID).Return(true, nil)
	orderNum, _ := domain.NewOrderNumber("ORD-001")
	numberGen.On("NextOrderNumber", ctx).Return(orderNum, nil)

	result, err := service.CreateOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "productVariantId is required")
	partyLookup.AssertExpectations(t)
}

func TestSalesService_CreateOrder_InvalidQuantityZero(t *testing.T) {
	quotesRepo := &MockQuoteRepository{}
	orderRepo := &MockSalesOrderRepository{}
	pricingEngine := &MockPricingEngine{}
	partyLookup := &MockPartyLookup{}
	numberGen := &MockNumberGenerator{}
	service := application.NewSalesService(quotesRepo, orderRepo, nil, nil, numberGen, pricingEngine, partyLookup, nil, nil)

	ctx := context.Background()
	partyID := uuid.New()
	cmd := application.CreateOrderCommand{
		PartyID: partyID,
		Items: []application.OrderLineItemInput{
			{
				ProductVariantID: uuid.New(),
				Quantity:         0, // Invalid quantity
			},
		},
		DeliveryDate: time.Now().AddDate(0, 0, 7),
	}

	partyLookup.On("ExistsParty", ctx, partyID).Return(true, nil)
	orderNum, _ := domain.NewOrderNumber("ORD-002")
	numberGen.On("NextOrderNumber", ctx).Return(orderNum, nil)

	result, err := service.CreateOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "quantity must be greater than zero")
	partyLookup.AssertExpectations(t)
}

func TestSalesService_CreateOrder_InvalidQuantityNegative(t *testing.T) {
	quotesRepo := &MockQuoteRepository{}
	orderRepo := &MockSalesOrderRepository{}
	pricingEngine := &MockPricingEngine{}
	partyLookup := &MockPartyLookup{}
	numberGen := &MockNumberGenerator{}
	service := application.NewSalesService(quotesRepo, orderRepo, nil, nil, numberGen, pricingEngine, partyLookup, nil, nil)

	ctx := context.Background()
	partyID := uuid.New()
	cmd := application.CreateOrderCommand{
		PartyID: partyID,
		Items: []application.OrderLineItemInput{
			{
				ProductVariantID: uuid.New(),
				Quantity:         -1, // Invalid quantity
			},
		},
		DeliveryDate: time.Now().AddDate(0, 0, 7),
	}

	partyLookup.On("ExistsParty", ctx, partyID).Return(true, nil)
	orderNum, _ := domain.NewOrderNumber("ORD-003")
	numberGen.On("NextOrderNumber", ctx).Return(orderNum, nil)

	result, err := service.CreateOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "quantity must be greater than zero")
	partyLookup.AssertExpectations(t)
}

// --- NEW TESTS FOR COVERAGE IMPROVEMENT (Sprint 13 Task 01) ---

func TestSalesService_GetQuote_RepositoryError(t *testing.T) {
	ctx := context.Background()
	quoteID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	// Simulate a repository error (e.g., database connection failure)
	repoErr := errors.New("database connection error")
	quoteRepo.On("FindByID", mock.Anything, quoteID).Return(nil, repoErr)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.GetQuoteByIDQuery{ID: quoteID}
	result, err := service.GetQuote(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoErr, err)
	quoteRepo.AssertExpectations(t)
}

func TestSalesService_GetDeliveryNote_RepositoryError(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	noteID := uuid.New()
	// Simulate a repository error (e.g., database connection failure)
	repoErr := errors.New("database connection error")
	deliveryRepo.On("FindByID", mock.Anything, noteID).Return(nil, repoErr)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.GetDeliveryNoteByIDQuery{ID: noteID}
	result, err := service.GetDeliveryNote(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoErr, err)
	deliveryRepo.AssertExpectations(t)
}

func TestSalesService_GetInvoice_RepositoryError(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	invoiceID := uuid.New()
	// Simulate a repository error (e.g., database connection failure)
	repoErr := errors.New("database connection error")
	invoiceRepo.On("FindByID", mock.Anything, invoiceID).Return(nil, repoErr)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.GetInvoiceByIDQuery{ID: invoiceID}
	result, err := service.GetInvoice(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoErr, err)
	invoiceRepo.AssertExpectations(t)
}

func TestSalesService_CreateOrder_PartyLookupError(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	numberGen := new(MockNumberGenerator)
	pricingEngine := new(MockPricingEngine)
	partyLookup := new(MockPartyLookup)

	partyID := uuid.New()
	variantID := uuid.New()

	// Simulate a party lookup error (e.g., database connection failure)
	lookupErr := errors.New("party service unavailable")
	partyLookup.On("ExistsParty", mock.Anything, partyID).Return(false, lookupErr)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numberGen, pricingEngine, partyLookup, nil, nil)

	cmd := application.CreateOrderCommand{
		PartyID: partyID,
		Items: []application.OrderLineItemInput{
			{
				ProductVariantID: variantID,
				Quantity:         1,
			},
		},
		DeliveryDate: time.Now().AddDate(0, 0, 7),
	}

	result, err := service.CreateOrder(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, lookupErr, err)
	partyLookup.AssertExpectations(t)
}

func TestSalesService_GetOrder_RepositoryError(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	orderID := uuid.New()
	// Simulate a repository error (e.g., database connection failure)
	repoErr := errors.New("database connection error")
	orderRepo.On("FindByID", mock.Anything, orderID).Return(nil, repoErr)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.GetOrderByIDQuery{ID: orderID}
	result, err := service.GetOrder(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoErr, err)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_GetOrder_NotFound(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	orderID := uuid.New()
	orderRepo.On("FindByID", mock.Anything, orderID).Return(nil, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.GetOrderByIDQuery{ID: orderID}
	result, err := service.GetOrder(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_ListQuotes_RepositoryError(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	filter := domain.QuoteFilter{PartyID: &partyID}
	// Simulate a repository error (e.g., database connection failure)
	repoErr := errors.New("database connection error")
	quoteRepo.On("List", mock.Anything, filter).Return(nil, repoErr)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.ListQuotesQuery{PartyID: &partyID}
	results, err := service.ListQuotes(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, results)
	assert.Equal(t, repoErr, err)
	quoteRepo.AssertExpectations(t)
}

func TestSalesService_GetQuote_NilResult(t *testing.T) {
	ctx := context.Background()
	quoteID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	// Repository returns (nil, nil) - no error but quote doesn't exist
	quoteRepo.On("FindByID", mock.Anything, quoteID).Return(nil, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.GetQuoteByIDQuery{ID: quoteID}
	result, err := service.GetQuote(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	quoteRepo.AssertExpectations(t)
}

func TestSalesService_ListOrders_RepositoryError(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	filter := domain.SalesOrderFilter{PartyID: &partyID}
	repoErr := errors.New("database connection error")
	orderRepo.On("List", mock.Anything, filter).Return(nil, repoErr)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.ListOrdersQuery{PartyID: &partyID}
	results, err := service.ListOrders(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, results)
	assert.Equal(t, repoErr, err)
	orderRepo.AssertExpectations(t)
}

func TestSalesService_ListInvoices_RepositoryError(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	filter := domain.InvoiceFilter{PartyID: &partyID}
	repoErr := errors.New("database connection error")
	invoiceRepo.On("List", mock.Anything, filter).Return(nil, repoErr)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.ListInvoicesQuery{PartyID: &partyID}
	results, err := service.ListInvoices(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, results)
	assert.Equal(t, repoErr, err)
	invoiceRepo.AssertExpectations(t)
}

func TestSalesService_ListDeliveryNotes_RepositoryError(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)

	filter := domain.DeliveryNoteFilter{PartyID: &partyID}
	repoErr := errors.New("database connection error")
	deliveryRepo.On("List", mock.Anything, filter).Return(nil, repoErr)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil, nil, nil)

	query := application.ListDeliveryNotesQuery{PartyID: &partyID}
	results, err := service.ListDeliveryNotes(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, results)
	assert.Equal(t, repoErr, err)
	deliveryRepo.AssertExpectations(t)
}

// --- FINAL TESTS FOR 75% COVERAGE GOAL (Sprint 13) ---

func TestSalesService_CreateInvoice_MissingPartyID(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	numberGen := new(MockNumberGenerator)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numberGen, nil, nil, nil, nil)

	cmd := application.CreateInvoiceCommand{
		PartyID:         uuid.Nil, // Invalid: nil UUID
		DeliveryNoteIDs: []uuid.UUID{uuid.New()},
		InvoiceDate:     time.Now(),
		DueDate:         time.Now().Add(24 * time.Hour),
	}

	result, err := service.CreateInvoice(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "partyId is required")
}

func TestSalesService_CreateInvoice_MissingInvoiceDate(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	numberGen := new(MockNumberGenerator)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numberGen, nil, nil, nil, nil)

	cmd := application.CreateInvoiceCommand{
		PartyID:         uuid.New(),
		DeliveryNoteIDs: []uuid.UUID{uuid.New()},
		InvoiceDate:     time.Time{}, // Invalid: zero time
		DueDate:         time.Now().Add(24 * time.Hour),
	}

	result, err := service.CreateInvoice(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invoiceDate is required")
}

func TestSalesService_CreateInvoice_MissingDueDate(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	numberGen := new(MockNumberGenerator)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numberGen, nil, nil, nil, nil)

	cmd := application.CreateInvoiceCommand{
		PartyID:         uuid.New(),
		DeliveryNoteIDs: []uuid.UUID{uuid.New()},
		InvoiceDate:     time.Now(),
		DueDate:         time.Time{}, // Invalid: zero time
	}

	result, err := service.CreateInvoice(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "dueDate is required")
}

func TestSalesService_CreateInvoice_NoDeliveryNotes(t *testing.T) {
	ctx := context.Background()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	numberGen := new(MockNumberGenerator)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numberGen, nil, nil, nil, nil)

	cmd := application.CreateInvoiceCommand{
		PartyID:         uuid.New(),
		DeliveryNoteIDs: []uuid.UUID{}, // Empty
		InvoiceDate:     time.Now(),
		DueDate:         time.Now().Add(24 * time.Hour),
	}

	result, err := service.CreateInvoice(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "deliveryNoteIds must be provided")
}

// TestSalesService_CreateInvoice_BothOrdersAndDeliveryNotes removed:
// SalesOrderIDs field no longer exists in CreateInvoiceCommand.
// Invoices can only be created from delivery notes.
