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
