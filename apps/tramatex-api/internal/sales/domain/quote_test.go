package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ===== NewQuote Tests =====

func TestNewQuote_Success(t *testing.T) {
	partyID := uuid.New()
	quoteDate := time.Now()
	expirationDate := quoteDate.Add(30 * 24 * time.Hour)
	money, _ := NewMoney(100.0, "EUR")
	taxAmount, _ := NewMoney(21.0, "EUR")
	number, _ := NewQuoteNumber("Q/2026/0001")

	lineItem, err := NewQuoteLineItem(
		uuid.New(),
		2,
		money,
		nil,
		nil,
		nil,
	)
	require.NoError(t, err)

	quote, err := NewQuote(
		number,
		partyID,
		quoteDate,
		expirationDate,
		[]QuoteLineItem{lineItem},
		taxAmount,
		"Test quote",
	)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, quote.ID)
	assert.Equal(t, number, quote.QuoteNumber)
	assert.Equal(t, partyID, quote.PartyID)
	assert.Equal(t, QuoteStatusDraft, quote.Status)
	assert.Equal(t, 1, len(quote.LineItems))
	assert.Equal(t, "Test quote", quote.Notes)

	// Verify totals calculation
	assert.Equal(t, 200.0, quote.Subtotal.Amount()) // 100 * 2
	assert.Equal(t, 21.0, quote.TaxAmount.Amount())
	assert.Equal(t, 221.0, quote.Total.Amount()) // 200 + 21
}

func TestNewQuote_EmptyPartyID(t *testing.T) {
	quoteDate := time.Now()
	expirationDate := quoteDate.Add(30 * 24 * time.Hour)
	taxAmount, _ := NewMoney(0, "EUR")
	number, _ := NewQuoteNumber("Q/2026/0001")

	quote, err := NewQuote(
		number,
		uuid.Nil,
		quoteDate,
		expirationDate,
		[]QuoteLineItem{},
		taxAmount,
		"",
	)

	assert.Error(t, err)
	assert.Nil(t, quote)
	assert.Contains(t, err.Error(), "party ID cannot be empty")
}

func TestNewQuote_ExpirationBeforeQuoteDate(t *testing.T) {
	partyID := uuid.New()
	quoteDate := time.Now()
	expirationDate := quoteDate.Add(-10 * 24 * time.Hour) // 10 days before
	taxAmount, _ := NewMoney(0, "EUR")
	number, _ := NewQuoteNumber("Q/2026/0001")

	quote, err := NewQuote(
		number,
		partyID,
		quoteDate,
		expirationDate,
		[]QuoteLineItem{},
		taxAmount,
		"",
	)

	assert.Error(t, err)
	assert.Nil(t, quote)
	assert.Contains(t, err.Error(), "expiration date cannot be before quote date")
}

func TestNewQuote_MultipleLineItems(t *testing.T) {
	partyID := uuid.New()
	quoteDate := time.Now()
	expirationDate := quoteDate.Add(30 * 24 * time.Hour)
	taxAmount, _ := NewMoney(10.0, "EUR")
	number, _ := NewQuoteNumber("Q/2026/0001")

	price1, _ := NewMoney(50.0, "EUR")
	price2, _ := NewMoney(30.0, "EUR")

	lineItem1, _ := NewQuoteLineItem(uuid.New(), 2, price1, nil, nil, nil)
	lineItem2, _ := NewQuoteLineItem(uuid.New(), 3, price2, nil, nil, nil)

	quote, err := NewQuote(
		number,
		partyID,
		quoteDate,
		expirationDate,
		[]QuoteLineItem{lineItem1, lineItem2},
		taxAmount,
		"",
	)

	require.NoError(t, err)
	assert.Equal(t, 2, len(quote.LineItems))
	// Subtotal = (50*2) + (30*3) = 100 + 90 = 190
	assert.Equal(t, 190.0, quote.Subtotal.Amount())
	// Total = 190 + 10 = 200
	assert.Equal(t, 200.0, quote.Total.Amount())
}

// ===== NewQuoteLineItem Tests =====

func TestNewQuoteLineItem_Success(t *testing.T) {
	variantID := uuid.New()
	price, _ := NewMoney(100.0, "EUR")

	lineItem, err := NewQuoteLineItem(
		variantID,
		5,
		price,
		nil,
		nil,
		nil,
	)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, lineItem.ID)
	assert.Equal(t, variantID, lineItem.ProductVariantID)
	assert.Equal(t, 5, lineItem.Quantity)
	assert.Equal(t, 100.0, lineItem.CalculatedUnitPrice.Amount())
	assert.Equal(t, 100.0, lineItem.FinalUnitPrice.Amount())
	assert.Equal(t, 0.0, lineItem.FinalDiscountPerUnit.Amount())
	assert.Equal(t, 500.0, lineItem.Subtotal.Amount()) // 100 * 5
}

func TestNewQuoteLineItem_WithManualPrice(t *testing.T) {
	variantID := uuid.New()
	calculatedPrice, _ := NewMoney(100.0, "EUR")
	manualPrice, _ := NewMoney(90.0, "EUR")

	lineItem, err := NewQuoteLineItem(
		variantID,
		2,
		calculatedPrice,
		&manualPrice,
		nil,
		nil,
	)

	require.NoError(t, err)
	assert.Equal(t, 100.0, lineItem.CalculatedUnitPrice.Amount())
	assert.Equal(t, 90.0, lineItem.FinalUnitPrice.Amount()) // Manual overrides calculated
	assert.Equal(t, 180.0, lineItem.Subtotal.Amount())      // 90 * 2
}

func TestNewQuoteLineItem_WithManualDiscount(t *testing.T) {
	variantID := uuid.New()
	price, _ := NewMoney(100.0, "EUR")
	discount, _ := NewMoney(10.0, "EUR")

	lineItem, err := NewQuoteLineItem(
		variantID,
		3,
		price,
		nil,
		nil,
		&discount,
	)

	require.NoError(t, err)
	assert.Equal(t, 100.0, lineItem.FinalUnitPrice.Amount())
	assert.Equal(t, 10.0, lineItem.FinalDiscountPerUnit.Amount())
	// Subtotal = (100 - 10) * 3 = 90 * 3 = 270
	assert.Equal(t, 270.0, lineItem.Subtotal.Amount())
}

func TestNewQuoteLineItem_WithCalculatedDiscount(t *testing.T) {
	variantID := uuid.New()
	price, _ := NewMoney(100.0, "EUR")
	calculatedDiscount, _ := NewMoney(5.0, "EUR")

	lineItem, err := NewQuoteLineItem(
		variantID,
		4,
		price,
		nil,
		&calculatedDiscount,
		nil,
	)

	require.NoError(t, err)
	assert.Equal(t, 5.0, lineItem.FinalDiscountPerUnit.Amount())
	// Subtotal = (100 - 5) * 4 = 95 * 4 = 380
	assert.Equal(t, 380.0, lineItem.Subtotal.Amount())
}

func TestNewQuoteLineItem_ManualOverridesCalculatedDiscount(t *testing.T) {
	variantID := uuid.New()
	price, _ := NewMoney(100.0, "EUR")
	calculatedDiscount, _ := NewMoney(5.0, "EUR")
	manualDiscount, _ := NewMoney(15.0, "EUR")

	lineItem, err := NewQuoteLineItem(
		variantID,
		2,
		price,
		nil,
		&calculatedDiscount,
		&manualDiscount,
	)

	require.NoError(t, err)
	assert.Equal(t, 15.0, lineItem.FinalDiscountPerUnit.Amount()) // Manual overrides calculated
	// Subtotal = (100 - 15) * 2 = 85 * 2 = 170
	assert.Equal(t, 170.0, lineItem.Subtotal.Amount())
}

func TestNewQuoteLineItem_EmptyVariantID(t *testing.T) {
	price, _ := NewMoney(100.0, "EUR")

	lineItem, err := NewQuoteLineItem(
		uuid.Nil,
		1,
		price,
		nil,
		nil,
		nil,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product variant ID cannot be empty")
	assert.Equal(t, uuid.Nil, lineItem.ID)
}

func TestNewQuoteLineItem_ZeroQuantity(t *testing.T) {
	price, _ := NewMoney(100.0, "EUR")

	_, err := NewQuoteLineItem(
		uuid.New(),
		0,
		price,
		nil,
		nil,
		nil,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quantity must be greater than zero")
}

func TestNewQuoteLineItem_NegativeQuantity(t *testing.T) {
	price, _ := NewMoney(100.0, "EUR")

	_, err := NewQuoteLineItem(
		uuid.New(),
		-5,
		price,
		nil,
		nil,
		nil,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quantity must be greater than zero")
}

func TestNewQuoteLineItem_CurrencyMismatch_ManualPrice(t *testing.T) {
	variantID := uuid.New()
	calculatedPrice, _ := NewMoney(100.0, "EUR")
	manualPrice, _ := NewMoney(90.0, "USD") // Different currency

	_, err := NewQuoteLineItem(
		variantID,
		1,
		calculatedPrice,
		&manualPrice,
		nil,
		nil,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "currency mismatch")
}

func TestNewQuoteLineItem_CurrencyMismatch_ManualDiscount(t *testing.T) {
	variantID := uuid.New()
	price, _ := NewMoney(100.0, "EUR")
	discount, _ := NewMoney(10.0, "GBP") // Different currency

	_, err := NewQuoteLineItem(
		variantID,
		1,
		price,
		nil,
		nil,
		&discount,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "currency mismatch")
}

// ===== ChangeStatus Tests =====

func TestQuote_ChangeStatus_DraftToSent(t *testing.T) {
	quote := createValidQuote(t)

	err := quote.ChangeStatus(QuoteStatusSent)

	assert.NoError(t, err)
	assert.Equal(t, QuoteStatusSent, quote.Status)
}

func TestQuote_ChangeStatus_SentToApproved(t *testing.T) {
	quote := createValidQuote(t)
	quote.Status = QuoteStatusSent

	err := quote.ChangeStatus(QuoteStatusApproved)

	assert.NoError(t, err)
	assert.Equal(t, QuoteStatusApproved, quote.Status)
}

func TestQuote_ChangeStatus_SentToRejected(t *testing.T) {
	quote := createValidQuote(t)
	quote.Status = QuoteStatusSent

	err := quote.ChangeStatus(QuoteStatusRejected)

	assert.NoError(t, err)
	assert.Equal(t, QuoteStatusRejected, quote.Status)
}

func TestQuote_ChangeStatus_SentToExpired(t *testing.T) {
	quote := createValidQuote(t)
	quote.Status = QuoteStatusSent

	err := quote.ChangeStatus(QuoteStatusExpired)

	assert.NoError(t, err)
	assert.Equal(t, QuoteStatusExpired, quote.Status)
}

func TestQuote_ChangeStatus_ApprovedToConverted(t *testing.T) {
	quote := createValidQuote(t)
	quote.Status = QuoteStatusApproved

	err := quote.ChangeStatus(QuoteStatusConverted)

	assert.NoError(t, err)
	assert.Equal(t, QuoteStatusConverted, quote.Status)
}

func TestQuote_ChangeStatus_InvalidTransition_DraftToApproved(t *testing.T) {
	quote := createValidQuote(t)

	err := quote.ChangeStatus(QuoteStatusApproved)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid quote status transition")
	assert.Equal(t, QuoteStatusDraft, quote.Status) // Status unchanged
}

func TestQuote_ChangeStatus_InvalidTransition_DraftToConverted(t *testing.T) {
	quote := createValidQuote(t)

	err := quote.ChangeStatus(QuoteStatusConverted)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid quote status transition")
}

func TestQuote_ChangeStatus_InvalidTransition_SentToDraft(t *testing.T) {
	quote := createValidQuote(t)
	quote.Status = QuoteStatusSent

	err := quote.ChangeStatus(QuoteStatusDraft)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid quote status transition")
}

func TestQuote_ChangeStatus_InvalidTransition_RejectedToAnything(t *testing.T) {
	quote := createValidQuote(t)
	quote.Status = QuoteStatusRejected

	err := quote.ChangeStatus(QuoteStatusSent)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid quote status transition")
}

func TestQuote_ChangeStatus_InvalidStatus(t *testing.T) {
	quote := createValidQuote(t)

	err := quote.ChangeStatus(QuoteStatus("INVALID_STATUS"))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid quote status")
}

// ===== RecalculateTotals Tests =====

func TestQuote_RecalculateTotals_SingleItem(t *testing.T) {
	quote := createValidQuote(t)

	// Modify line item subtotal
	newPrice, _ := NewMoney(200.0, "EUR")
	quote.LineItems[0].Subtotal = newPrice

	err := quote.RecalculateTotals()

	require.NoError(t, err)
	assert.Equal(t, 200.0, quote.Subtotal.Amount())
	// Total = 200 (subtotal) + 10 (tax)
	assert.Equal(t, 210.0, quote.Total.Amount())
}

func TestQuote_RecalculateTotals_MultipleItems(t *testing.T) {
	partyID := uuid.New()
	quoteDate := time.Now()
	expirationDate := quoteDate.Add(30 * 24 * time.Hour)
	taxAmount, _ := NewMoney(20.0, "EUR")
	number, _ := NewQuoteNumber("Q/2026/0001")

	price1, _ := NewMoney(50.0, "EUR")
	price2, _ := NewMoney(30.0, "EUR")

	lineItem1, _ := NewQuoteLineItem(uuid.New(), 2, price1, nil, nil, nil)
	lineItem2, _ := NewQuoteLineItem(uuid.New(), 3, price2, nil, nil, nil)

	quote, _ := NewQuote(
		number,
		partyID,
		quoteDate,
		expirationDate,
		[]QuoteLineItem{lineItem1, lineItem2},
		taxAmount,
		"",
	)

	// Modify line items
	newSubtotal1, _ := NewMoney(120.0, "EUR")
	newSubtotal2, _ := NewMoney(80.0, "EUR")
	quote.LineItems[0].Subtotal = newSubtotal1
	quote.LineItems[1].Subtotal = newSubtotal2

	err := quote.RecalculateTotals()

	require.NoError(t, err)
	// Subtotal = 120 + 80 = 200
	assert.Equal(t, 200.0, quote.Subtotal.Amount())
	// Total = 200 + 20 = 220
	assert.Equal(t, 220.0, quote.Total.Amount())
}

func TestQuote_RecalculateTotals_WithZeroTax(t *testing.T) {
	quote := createValidQuote(t)
	zeroTax, _ := NewMoney(0, "EUR")
	quote.TaxAmount = zeroTax

	// Modify line item subtotal
	newPrice, _ := NewMoney(150.0, "EUR")
	quote.LineItems[0].Subtotal = newPrice

	err := quote.RecalculateTotals()

	require.NoError(t, err)
	assert.Equal(t, 150.0, quote.Subtotal.Amount())
	assert.Equal(t, 150.0, quote.Total.Amount()) // No tax
}

// ===== Helper Functions =====

func createValidQuote(t *testing.T) *Quote {
	partyID := uuid.New()
	quoteDate := time.Now()
	expirationDate := quoteDate.Add(30 * 24 * time.Hour)
	price, _ := NewMoney(100.0, "EUR")
	taxAmount, _ := NewMoney(10.0, "EUR")
	number, _ := NewQuoteNumber("Q/2026/0001")

	lineItem, _ := NewQuoteLineItem(
		uuid.New(),
		1,
		price,
		nil,
		nil,
		nil,
	)

	quote, err := NewQuote(
		number,
		partyID,
		quoteDate,
		expirationDate,
		[]QuoteLineItem{lineItem},
		taxAmount,
		"",
	)
	require.NoError(t, err)
	return quote
}
