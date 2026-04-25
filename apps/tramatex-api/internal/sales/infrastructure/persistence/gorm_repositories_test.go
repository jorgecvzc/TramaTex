package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/sales/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupSalesTestDB(t *testing.T) (*TestDB, func()) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	require.NoError(t, tdb.SetUpSales())

	cleanup := func() {
		_ = tdb.TearDownSales()
		sqlDB, _ := tdb.DB.DB()
		_ = sqlDB.Close()
	}

	return tdb, cleanup
}

func TestSalesDataModelConversions(t *testing.T) {
	partyID := uuid.New()
	variantID := uuid.New()
	now := time.Now().UTC()

	quoteNumber, _ := domain.NewQuoteNumber("Q-100")
	unitPrice, _ := domain.NewMoney(10, domain.DefaultCurrency)
	tax, _ := domain.NewMoney(0, domain.DefaultCurrency)

	quoteItem, _ := domain.NewQuoteLineItem(variantID, 2, unitPrice, nil, 0)
	quote, err := domain.NewQuote(quoteNumber, partyID, now, now.Add(24*time.Hour), []domain.QuoteLineItem{quoteItem}, tax, "notes")
	assert.NoError(t, err)

	quoteModel, err := quoteFromDomain(quote)
	assert.NoError(t, err)
	itemsModel, err := quoteLineItemsFromDomain(quote.ID, quote.LineItems)
	assert.NoError(t, err)
	mappedQuote, err := quoteToDomain(quoteModel, itemsModel)
	assert.NoError(t, err)
	assert.Equal(t, quote.ID, mappedQuote.ID)

	orderNumber, _ := domain.NewOrderNumber("SO-200")
	orderItem, _ := domain.NewOrderLineItem(variantID, 1, unitPrice, nil, 0)
	order, err := domain.NewSalesOrder(orderNumber, partyID, now, now.Add(48*time.Hour), []domain.OrderLineItem{orderItem}, tax, "order")
	assert.NoError(t, err)

	orderModel, err := salesOrderFromDomain(order)
	assert.NoError(t, err)
	orderItems, err := orderLineItemsFromDomain(order.ID, order.LineItems)
	assert.NoError(t, err)
	mappedOrder, err := salesOrderToDomain(orderModel, orderItems)
	assert.NoError(t, err)
	assert.Equal(t, order.ID, mappedOrder.ID)

	noteNumber, _ := domain.NewDeliveryNoteNumber("DN-300")
	noteItem, _ := domain.NewDeliveryNoteLineItem(orderItem.ID, variantID, 1)
	note, err := domain.NewDeliveryNote(noteNumber, order.ID, partyID, now.Add(72*time.Hour), []domain.DeliveryNoteLineItem{noteItem}, "note")
	assert.NoError(t, err)

	noteModel, err := deliveryNoteFromDomain(note)
	assert.NoError(t, err)
	noteItems, err := deliveryNoteLineItemsFromDomain(note.ID, note.LineItems)
	assert.NoError(t, err)
	mappedNote, err := deliveryNoteToDomain(noteModel, noteItems)
	assert.NoError(t, err)
	assert.Equal(t, note.ID, mappedNote.ID)

	invoiceNumber, _ := domain.NewInvoiceNumber("INV-400")
	invoiceType := domain.InvoiceTypeComplete
	series, _ := domain.NewInvoiceSeries("A", 2026)
	invoiceItem, _ := domain.NewInvoiceLineItem(variantID, 1, unitPrice, nil, nil)
	invoiceItem.SalesOrderLineItemID = &orderItem.ID
	invoice, err := domain.NewInvoice(invoiceNumber, invoiceType, series, partyID, now, now.Add(96*time.Hour), []domain.InvoiceLineItem{invoiceItem}, tax, "NET30")
	assert.NoError(t, err)

	invoiceModel, err := invoiceFromDomain(invoice)
	assert.NoError(t, err)
	invoiceItems, err := invoiceLineItemsFromDomain(invoice.ID, invoice.LineItems)
	assert.NoError(t, err)
	mappedInvoice, err := invoiceToDomain(invoiceModel, invoiceItems)
	assert.NoError(t, err)
	assert.Equal(t, invoice.ID, mappedInvoice.ID)
}

func TestGORMRepositories_Sales(t *testing.T) {
	tdb, cleanup := setupSalesTestDB(t)
	defer cleanup()

	ctx := context.Background()
	variantID := uuid.New()
	partyID := uuid.New()

	assert.NoError(t, tdb.DB.WithContext(ctx).Exec("INSERT INTO product_variants (id) VALUES (?)", variantID).Error)

	quoteRepo := NewGORMQuoteRepository(tdb.DB)
	orderRepo := NewGORMSalesOrderRepository(tdb.DB)
	deliveryRepo := NewGORMDeliveryNoteRepository(tdb.DB)
	invoiceRepo := NewGORMInvoiceRepository(tdb.DB)

	now := time.Now().UTC()
	money, _ := domain.NewMoney(10, domain.DefaultCurrency)
	tax, _ := domain.NewMoney(0, domain.DefaultCurrency)

	quoteNumber, _ := domain.NewQuoteNumber("Q-500")
	quoteItem, _ := domain.NewQuoteLineItem(variantID, 2, money, nil, 0)
	quote, err := domain.NewQuote(quoteNumber, partyID, now, now.Add(24*time.Hour), []domain.QuoteLineItem{quoteItem}, tax, "quote")
	assert.NoError(t, err)
	assert.NoError(t, quoteRepo.Save(ctx, quote))

	foundQuote, err := quoteRepo.FindByID(ctx, quote.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundQuote)
	assert.Equal(t, quote.ID, foundQuote.ID)
	quotes, err := quoteRepo.List(ctx, domain.QuoteFilter{PartyID: &partyID})
	assert.NoError(t, err)
	assert.Len(t, quotes, 1)

	orderNumber, _ := domain.NewOrderNumber("SO-600")
	orderItem, _ := domain.NewOrderLineItem(variantID, 1, money, nil, 0)
	order, err := domain.NewSalesOrder(orderNumber, partyID, now, now.Add(48*time.Hour), []domain.OrderLineItem{orderItem}, tax, "order")
	assert.NoError(t, err)
	assert.NoError(t, orderRepo.Save(ctx, order))

	foundOrder, err := orderRepo.FindByID(ctx, order.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundOrder)
	orders, err := orderRepo.List(ctx, domain.SalesOrderFilter{PartyID: &partyID})
	assert.NoError(t, err)
	assert.Len(t, orders, 1)

	orderLineID := order.LineItems[0].ID
	noteNumber, _ := domain.NewDeliveryNoteNumber("DN-700")
	noteItem, _ := domain.NewDeliveryNoteLineItem(orderLineID, variantID, 1)
	note, err := domain.NewDeliveryNote(noteNumber, order.ID, partyID, now.Add(72*time.Hour), []domain.DeliveryNoteLineItem{noteItem}, "note")
	assert.NoError(t, err)
	assert.NoError(t, deliveryRepo.Save(ctx, note))

	foundNote, err := deliveryRepo.FindByID(ctx, note.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundNote)
	notes, err := deliveryRepo.ListBySalesOrderID(ctx, order.ID)
	assert.NoError(t, err)
	assert.Len(t, notes, 1)

	invoiceNumber, _ := domain.NewInvoiceNumber("INV-800")
	invoiceType := domain.InvoiceTypeComplete
	series, _ := domain.NewInvoiceSeries("A", 2026)
	invoiceItem, _ := domain.NewInvoiceLineItem(variantID, 1, order.LineItems[0].UnitPrice, nil, nil)
	invoiceItem.SalesOrderLineItemID = &orderLineID
	invoice, err := domain.NewInvoice(invoiceNumber, invoiceType, series, partyID, now, now.Add(96*time.Hour), []domain.InvoiceLineItem{invoiceItem}, tax, "NET30")
	assert.NoError(t, err)
	assert.NoError(t, invoiceRepo.Save(ctx, invoice))

	foundInvoice, err := invoiceRepo.FindByID(ctx, invoice.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundInvoice)
	invoices, err := invoiceRepo.ListBySalesOrderID(ctx, order.ID)
	assert.NoError(t, err)
	assert.Len(t, invoices, 1)
}

func TestGORMRepositories_SalesListsAreOrderedByDocumentNumberDescending(t *testing.T) {
	tdb, cleanup := setupSalesTestDB(t)
	defer cleanup()

	ctx := context.Background()
	variantID := uuid.New()
	partyID := uuid.New()

	assert.NoError(t, tdb.DB.WithContext(ctx).Exec("INSERT INTO product_variants (id) VALUES (?)", variantID).Error)

	quoteRepo := NewGORMQuoteRepository(tdb.DB)
	orderRepo := NewGORMSalesOrderRepository(tdb.DB)
	deliveryRepo := NewGORMDeliveryNoteRepository(tdb.DB)
	invoiceRepo := NewGORMInvoiceRepository(tdb.DB)

	now := time.Now().UTC()
	money, _ := domain.NewMoney(10, domain.DefaultCurrency)
	tax, _ := domain.NewMoney(0, domain.DefaultCurrency)

	highQuoteNumber, _ := domain.NewQuoteNumber("PRE-2026-0010")
	lowQuoteNumber, _ := domain.NewQuoteNumber("PRE-2026-0002")
	quoteItemA, _ := domain.NewQuoteLineItem(variantID, 1, money, nil, 0)
	quoteItemB, _ := domain.NewQuoteLineItem(variantID, 1, money, nil, 0)
	highQuote, _ := domain.NewQuote(highQuoteNumber, partyID, now, now.Add(24*time.Hour), []domain.QuoteLineItem{quoteItemA}, tax, "high quote")
	lowQuote, _ := domain.NewQuote(lowQuoteNumber, partyID, now.Add(time.Minute), now.Add(25*time.Hour), []domain.QuoteLineItem{quoteItemB}, tax, "low quote")
	assert.NoError(t, quoteRepo.Save(ctx, highQuote))
	assert.NoError(t, quoteRepo.Save(ctx, lowQuote))

	quotes, err := quoteRepo.List(ctx, domain.QuoteFilter{PartyID: &partyID})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(quotes), 2)
	assert.Equal(t, highQuote.ID, quotes[0].ID)
	assert.Equal(t, lowQuote.ID, quotes[1].ID)

	highOrderNumber, _ := domain.NewOrderNumber("PED-2026-0010")
	lowOrderNumber, _ := domain.NewOrderNumber("PED-2026-0002")
	orderItemA, _ := domain.NewOrderLineItem(variantID, 1, money, nil, 0)
	orderItemB, _ := domain.NewOrderLineItem(variantID, 1, money, nil, 0)
	highOrder, _ := domain.NewSalesOrder(highOrderNumber, partyID, now, now.Add(48*time.Hour), []domain.OrderLineItem{orderItemA}, tax, "high order")
	lowOrder, _ := domain.NewSalesOrder(lowOrderNumber, partyID, now.Add(time.Minute), now.Add(49*time.Hour), []domain.OrderLineItem{orderItemB}, tax, "low order")
	assert.NoError(t, orderRepo.Save(ctx, highOrder))
	assert.NoError(t, orderRepo.Save(ctx, lowOrder))

	orders, err := orderRepo.List(ctx, domain.SalesOrderFilter{PartyID: &partyID})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(orders), 2)
	assert.Equal(t, highOrder.ID, orders[0].ID)
	assert.Equal(t, lowOrder.ID, orders[1].ID)

	highDeliveryNumber, _ := domain.NewDeliveryNoteNumber("ALB-2026-0010")
	lowDeliveryNumber, _ := domain.NewDeliveryNoteNumber("ALB-2026-0002")
	highDeliveryItem, _ := domain.NewDeliveryNoteLineItem(highOrder.LineItems[0].ID, variantID, 1)
	lowDeliveryItem, _ := domain.NewDeliveryNoteLineItem(lowOrder.LineItems[0].ID, variantID, 1)
	highDelivery, _ := domain.NewDeliveryNote(highDeliveryNumber, highOrder.ID, partyID, now.Add(72*time.Hour), []domain.DeliveryNoteLineItem{highDeliveryItem}, "high delivery")
	lowDelivery, _ := domain.NewDeliveryNote(lowDeliveryNumber, lowOrder.ID, partyID, now.Add(73*time.Hour), []domain.DeliveryNoteLineItem{lowDeliveryItem}, "low delivery")
	assert.NoError(t, deliveryRepo.Save(ctx, highDelivery))
	assert.NoError(t, deliveryRepo.Save(ctx, lowDelivery))

	deliveryNotes, err := deliveryRepo.List(ctx, domain.DeliveryNoteFilter{PartyID: &partyID})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(deliveryNotes), 2)
	assert.Equal(t, highDelivery.ID, deliveryNotes[0].ID)
	assert.Equal(t, lowDelivery.ID, deliveryNotes[1].ID)

	series, _ := domain.NewInvoiceSeries("FV", 2026)
	highInvoiceNumber, _ := domain.NewInvoiceNumber("FV-2026-0010")
	lowInvoiceNumber, _ := domain.NewInvoiceNumber("FV-2026-0002")
	highInvoiceItem, _ := domain.NewInvoiceLineItem(variantID, 1, money, nil, nil)
	lowInvoiceItem, _ := domain.NewInvoiceLineItem(variantID, 1, money, nil, nil)
	highInvoice := mustNewInvoiceForSortTest(t, highInvoiceNumber, series, partyID, now, highInvoiceItem, tax)
	lowInvoice := mustNewInvoiceForSortTest(t, lowInvoiceNumber, series, partyID, now.Add(time.Minute), lowInvoiceItem, tax)
	assert.NoError(t, invoiceRepo.Save(ctx, highInvoice))
	assert.NoError(t, invoiceRepo.Save(ctx, lowInvoice))

	invoices, err := invoiceRepo.List(ctx, domain.InvoiceFilter{PartyID: &partyID})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(invoices), 2)
	assert.Equal(t, highInvoice.ID, invoices[0].ID)
	assert.Equal(t, lowInvoice.ID, invoices[1].ID)
}

func mustNewInvoiceForSortTest(t *testing.T, number domain.InvoiceNumber, series domain.InvoiceSeries, partyID uuid.UUID, invoiceDate time.Time, item domain.InvoiceLineItem, tax domain.Money) *domain.Invoice {
	t.Helper()
	invoice, err := domain.NewInvoice(number, domain.InvoiceTypeComplete, series, partyID, invoiceDate, invoiceDate.Add(24*time.Hour), []domain.InvoiceLineItem{item}, tax, "NET30")
	assert.NoError(t, err)
	return invoice
}
