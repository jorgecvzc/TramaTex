package application_test

import (
	"context"
	"errors"
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

func (m *MockSalesOrderRepository) List(ctx context.Context, filter domain.SalesOrderFilter) ([]*domain.SalesOrder, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.SalesOrder), args.Error(1)
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
}

func (m *MockPartyLookup) ExistsParty(ctx context.Context, partyID uuid.UUID) (bool, error) {
	args := m.Called(ctx, partyID)
	return args.Bool(0), args.Error(1)
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

func (m *MockNumberGenerator) NextInvoiceNumber(ctx context.Context) (domain.InvoiceNumber, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.InvoiceNumber), args.Error(1)
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
		return req.ClientID == partyID && len(req.SaleItems) == 1 && req.SaleItems[0].ProductVariantID == variantID && req.SaleItems[0].Quantity == 2
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

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, pricing, partyLookup)

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

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, pricing, partyLookup)
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
	lineItem, _ := domain.NewQuoteLineItem(uuid.New(), 1, money, nil, nil, nil)
	quote, _ := domain.NewQuote(quoteNumber, partyID, time.Now(), time.Now().Add(24*time.Hour), []domain.QuoteLineItem{lineItem}, money, "")

	quoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)
	orderNumber, _ := domain.NewOrderNumber("SO-1")
	numbers.On("NextOrderNumber", mock.Anything).Return(orderNumber, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, pricing, nil)
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
	orderItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, nil, nil)
	order, _ := domain.NewSalesOrder(orderNumber, partyID, time.Now(), time.Now().Add(48*time.Hour), []domain.OrderLineItem{orderItem}, money, "")

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	deliveryRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.DeliveryNote{}, nil)

	noteNumber, _ := domain.NewDeliveryNoteNumber("DN-10")
	numbers.On("NextDeliveryNoteNumber", mock.Anything).Return(noteNumber, nil)

	var savedOrder *domain.SalesOrder
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Run(func(args mock.Arguments) {
		savedOrder = args.Get(1).(*domain.SalesOrder)
	}).Return(nil)
	deliveryRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.DeliveryNote")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, nil, nil)
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

func TestSalesService_CreateInvoice_FromDeliveredOrder(t *testing.T) {
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
	orderItem, _ := domain.NewOrderLineItem(variantID, 1, money, nil, nil, nil)
	order, _ := domain.NewSalesOrder(orderNumber, partyID, time.Now(), time.Now().Add(48*time.Hour), []domain.OrderLineItem{orderItem}, money, "")
	_ = order.ChangeStatus(domain.SalesOrderStatusInPreparation)
	_ = order.ChangeStatus(domain.SalesOrderStatusDelivered)

	orderRepo.On("FindByID", mock.Anything, order.ID).Return(order, nil)
	invoiceRepo.On("ListBySalesOrderID", mock.Anything, order.ID).Return([]*domain.Invoice{}, nil)

	invoiceNumber, _ := domain.NewInvoiceNumber("INV-20")
	numbers.On("NextInvoiceNumber", mock.Anything).Return(invoiceNumber, nil)

	var savedOrder *domain.SalesOrder
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Run(func(args mock.Arguments) {
		savedOrder = args.Get(1).(*domain.SalesOrder)
	}).Return(nil)
	invoiceRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Invoice")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, nil, nil)

	cmd := application.CreateInvoiceCommand{
		PartyID:       partyID,
		SalesOrderIDs: []uuid.UUID{order.ID},
		InvoiceDate:   time.Now(),
		DueDate:       time.Now().Add(24 * time.Hour),
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
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, nil, nil)
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

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil)

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

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil)

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
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, nil, nil)
	taxAmount, _ := domain.NewMoney(42, domain.DefaultCurrency)

	quote1, _ := domain.NewQuote(quoteNumber1, partyID, time.Now(), time.Now().Add(30*24*time.Hour), []domain.QuoteLineItem{lineItem}, taxAmount, "")
	quote2, _ := domain.NewQuote(quoteNumber2, partyID, time.Now(), time.Now().Add(30*24*time.Hour), []domain.QuoteLineItem{lineItem}, taxAmount, "")

	quotes := []*domain.Quote{quote1, quote2}
	filter := domain.QuoteFilter{PartyID: &partyID}

	quoteRepo.On("List", mock.Anything, filter).Return(quotes, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil)

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
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, nil, nil)
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

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil)

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
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, nil, nil)
	taxAmount, _ := domain.NewMoney(105, domain.DefaultCurrency)

	order1, _ := domain.NewSalesOrder(orderNumber1, partyID, time.Now(), time.Now().Add(7*24*time.Hour), []domain.OrderLineItem{lineItem}, taxAmount, "")
	order2, _ := domain.NewSalesOrder(orderNumber2, partyID, time.Now(), time.Now().Add(7*24*time.Hour), []domain.OrderLineItem{lineItem}, taxAmount, "")

	orders := []*domain.SalesOrder{order1, order2}
	filter := domain.SalesOrderFilter{PartyID: &partyID}

	orderRepo.On("List", mock.Anything, filter).Return(orders, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil)

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
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, nil, nil)
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
	_ = quote.ChangeStatus(domain.QuoteStatusSent)
	_ = quote.ChangeStatus(domain.QuoteStatusApproved)

	quoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)
	quoteRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Quote")).Return(nil)

	orderNumber, _ := domain.NewOrderNumber("SO-001")
	numbers.On("NextOrderNumber", mock.Anything).Return(orderNumber, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, nil, nil)

	cmd := application.ConvertQuoteToOrderCommand{
		QuoteID:      quote.ID,
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
	}

	result, err := service.ConvertQuoteToOrder(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, partyID, result.PartyID)
	assert.Equal(t, "PENDIENTE", string(result.Status))
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
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, nil, nil)
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
	// Quote is in DRAFT status, not APPROVED

	quoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)

	// Note: Service generates order number before validation - this is a code smell
	// but we need to mock it for the test to run
	orderNumber, _ := domain.NewOrderNumber("SO-999")
	numbers.On("NextOrderNumber", mock.Anything).Return(orderNumber, nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, nil, nil)

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
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, nil, nil)
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

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil)

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
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, nil, nil)
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

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil)

	cmd := application.ChangeQuoteStatusCommand{
		QuoteID:   quote.ID,
		NewStatus: string(domain.QuoteStatusSent),
	}

	result, err := service.ChangeQuoteStatus(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ENVIADA", string(result.Status))
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
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, nil, nil)
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

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil)

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
	lineItem, _ := domain.NewOrderLineItem(variantID, 5, money, nil, nil, nil)
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

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, nil, nil, nil)

	cmd := application.ChangeOrderStatusCommand{
		OrderID:   order.ID,
		NewStatus: string(domain.SalesOrderStatusInPreparation),
	}

	result, err := service.ChangeOrderStatus(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "EN_PREPARACION", string(result.Status))
	orderRepo.AssertExpectations(t)
}
