package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/joran-cortez/tramatex/internal/sales/application"
	"github.com/joran-cortez/tramatex/internal/sales/domain"
)

func TestSalesService_ConvertQuoteToOrder_AutoApprovesFromIssued(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	numbers := new(MockNumberGenerator)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	quoteNumber, _ := domain.NewQuoteNumber("Q/2026/AUTO")
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, 0)
	
	quote, _ := domain.NewQuote(
		quoteNumber,
		partyID,
		time.Now(),
		time.Now().Add(30*24*time.Hour),
		[]domain.QuoteLineItem{lineItem},
		money,
		"Test auto-approve",
	)
	// Set to ISSUED status
	_ = quote.ChangeStatus(domain.QuoteStatusIssued)

	quoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)
	quoteRepo.On("Save", mock.Anything, mock.MatchedBy(func(q *domain.Quote) bool {
		return q.Status == domain.QuoteStatusConverted
	})).Return(nil)

	orderNumber, _ := domain.NewOrderNumber("SO-AUTO-1")
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
	assert.Equal(t, "CONVERTED_TO_ORDER", string(quote.Status))
	quoteRepo.AssertExpectations(t)
}

func TestSalesService_AcceptAndConvertQuote_Success(t *testing.T) {
	ctx := context.Background()
	partyID := uuid.New()
	variantID := uuid.New()

	quoteRepo := new(MockQuoteRepository)
	orderRepo := new(MockSalesOrderRepository)
	deliveryRepo := new(MockDeliveryNoteRepository)
	invoiceRepo := new(MockInvoiceRepository)
	numbers := new(MockNumberGenerator)

	money, _ := domain.NewMoney(100, domain.DefaultCurrency)
	quoteNumber, _ := domain.NewQuoteNumber("Q/2026/ACCEPT")
	lineItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, 0)
	
	quote, _ := domain.NewQuote(
		quoteNumber,
		partyID,
		time.Now(),
		time.Now().Add(30*24*time.Hour),
		[]domain.QuoteLineItem{lineItem},
		money,
		"Test accept and convert",
	)

	quoteRepo.On("FindByID", mock.Anything, quote.ID).Return(quote, nil)
	quoteRepo.On("Save", mock.Anything, mock.MatchedBy(func(q *domain.Quote) bool {
		return q.Status == domain.QuoteStatusConverted
	})).Return(nil)

	orderNumber, _ := domain.NewOrderNumber("SO-ACCEPT-1")
	numbers.On("NextOrderNumber", mock.Anything).Return(orderNumber, nil)
	orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.SalesOrder")).Return(nil)

	service := application.NewSalesService(quoteRepo, orderRepo, deliveryRepo, invoiceRepo, numbers, nil, nil, nil, nil)

	cmd := application.AcceptAndConvertQuoteCommand{
		QuoteID:      quote.ID,
		DeliveryDate: time.Now().Add(7 * 24 * time.Hour),
	}

	result, err := service.AcceptAndConvertQuote(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "CONVERTED_TO_ORDER", string(quote.Status))
	quoteRepo.AssertExpectations(t)
}
