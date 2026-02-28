package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	pricing_app "github.com/joran-cortez/tramatex/internal/pricing/application"
	"github.com/joran-cortez/tramatex/internal/sales/domain"
)

type PricingEngine interface {
	CalculateFinalSalePrice(ctx context.Context, req pricing_app.CalculateFinalSalePriceRequest) (*pricing_app.CalculateFinalSalePriceResponse, error)
}

type PartyLookup interface {
	ExistsParty(ctx context.Context, partyID uuid.UUID) (bool, error)
	HasPartyRole(ctx context.Context, partyID uuid.UUID, role string) (bool, error)
}

type DocumentNumberGenerator interface {
	NextQuoteNumber(ctx context.Context) (domain.QuoteNumber, error)
	NextOrderNumber(ctx context.Context) (domain.OrderNumber, error)
	NextDeliveryNoteNumber(ctx context.Context) (domain.DeliveryNoteNumber, error)
	NextInvoiceNumber(ctx context.Context) (domain.InvoiceNumber, error)
}

type SalesService struct {
	quoteRepo     domain.QuoteRepository
	orderRepo     domain.SalesOrderRepository
	deliveryRepo  domain.DeliveryNoteRepository
	invoiceRepo   domain.InvoiceRepository
	numberGen     DocumentNumberGenerator
	pricingEngine PricingEngine
	partyLookup   PartyLookup
}

func NewSalesService(
	quoteRepo domain.QuoteRepository,
	orderRepo domain.SalesOrderRepository,
	deliveryRepo domain.DeliveryNoteRepository,
	invoiceRepo domain.InvoiceRepository,
	numberGen DocumentNumberGenerator,
	pricingEngine PricingEngine,
	partyLookup PartyLookup,
) *SalesService {
	return &SalesService{
		quoteRepo:     quoteRepo,
		orderRepo:     orderRepo,
		deliveryRepo:  deliveryRepo,
		invoiceRepo:   invoiceRepo,
		numberGen:     numberGen,
		pricingEngine: pricingEngine,
		partyLookup:   partyLookup,
	}
}

func (s *SalesService) CreateQuote(ctx context.Context, cmd CreateQuoteCommand) (*QuoteDTO, error) {
	if cmd.PartyID == uuid.Nil {
		return nil, domain.NewValidationError("partyId is required")
	}
	if cmd.ExpirationDate.IsZero() {
		return nil, domain.NewValidationError("expirationDate is required")
	}
	if len(cmd.Items) == 0 {
		return nil, domain.NewValidationError("items cannot be empty")
	}
	if err := s.ensurePartyExists(ctx, cmd.PartyID); err != nil {
		return nil, err
	}
	if s.numberGen == nil {
		return nil, fmt.Errorf("quote number generator not configured")
	}

	quoteNumber, err := s.numberGen.NextQuoteNumber(ctx)
	if err != nil {
		return nil, err
	}

	lineItems, err := s.buildQuoteLineItems(ctx, cmd.PartyID, cmd.Items)
	if err != nil {
		return nil, err
	}

	taxAmount, err := zeroMoney()
	if err != nil {
		return nil, err
	}
	for _, lineItem := range lineItems {
		taxAmount, err = taxAmount.Add(lineItem.TaxAmount)
		if err != nil {
			return nil, err
		}
	}
	notes := ""
	if cmd.Notes != nil {
		notes = *cmd.Notes
	}

	quote, err := domain.NewQuote(
		quoteNumber,
		cmd.PartyID,
		time.Now(),
		cmd.ExpirationDate,
		lineItems,
		taxAmount,
		notes,
	)
	if err != nil {
		return nil, err
	}

	if err := s.quoteRepo.Save(ctx, quote); err != nil {
		return nil, err
	}

	return NewQuoteDTO(quote), nil
}

func (s *SalesService) UpdateQuote(ctx context.Context, cmd UpdateQuoteCommand) (*QuoteDTO, error) {
	quote, err := s.quoteRepo.FindByID(ctx, cmd.QuoteID)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, domain.NewNotFoundError("quote not found")
	}
	if quote.Status != domain.QuoteStatusDraft {
		return nil, domain.NewConflictError("only draft quotes can be updated")
	}

	if cmd.ExpirationDate != nil {
		if cmd.ExpirationDate.Before(quote.QuoteDate) {
			return nil, domain.NewValidationError("expirationDate cannot be before quoteDate")
		}
		quote.ExpirationDate = *cmd.ExpirationDate
	}
	if cmd.Notes != nil {
		quote.Notes = *cmd.Notes
	}

	if cmd.Items != nil {
		if len(cmd.Items) == 0 {
			return nil, domain.NewValidationError("items cannot be empty")
		}
		lineItems, err := s.buildQuoteLineItems(ctx, quote.PartyID, cmd.Items)
		if err != nil {
			return nil, err
		}
		quote.LineItems = lineItems
		quote.TaxAmount, err = zeroMoney()
		if err != nil {
			return nil, err
		}
		if err := quote.RecalculateTotals(); err != nil {
			return nil, err
		}
	}

	if err := s.quoteRepo.Save(ctx, quote); err != nil {
		return nil, err
	}

	return NewQuoteDTO(quote), nil
}

func (s *SalesService) ChangeQuoteStatus(ctx context.Context, cmd ChangeQuoteStatusCommand) (*QuoteDTO, error) {
	quote, err := s.quoteRepo.FindByID(ctx, cmd.QuoteID)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, domain.NewNotFoundError("quote not found")
	}

	status, err := parseQuoteStatus(cmd.NewStatus)
	if err != nil {
		return nil, err
	}
	if err := quote.ChangeStatus(status); err != nil {
		return nil, err
	}

	if err := s.quoteRepo.Save(ctx, quote); err != nil {
		return nil, err
	}

	return NewQuoteDTO(quote), nil
}

func (s *SalesService) ConvertQuoteToOrder(ctx context.Context, cmd ConvertQuoteToOrderCommand) (*SalesOrderDTO, error) {
	quote, err := s.quoteRepo.FindByID(ctx, cmd.QuoteID)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, domain.NewNotFoundError("quote not found")
	}
	if s.numberGen == nil {
		return nil, fmt.Errorf("order number generator not configured")
	}

	orderNumber, err := s.numberGen.NextOrderNumber(ctx)
	if err != nil {
		return nil, err
	}

	order, err := quote.ConvertToOrder(orderNumber, cmd.DeliveryDate)
	if err != nil {
		return nil, err
	}

	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}
	if err := s.quoteRepo.Save(ctx, quote); err != nil {
		return nil, err
	}

	return NewSalesOrderDTO(order), nil
}

func (s *SalesService) CreateOrder(ctx context.Context, cmd CreateOrderCommand) (*SalesOrderDTO, error) {
	if cmd.PartyID == uuid.Nil {
		return nil, domain.NewValidationError("partyId is required")
	}
	if cmd.DeliveryDate.IsZero() {
		return nil, domain.NewValidationError("deliveryDate is required")
	}
	if cmd.QuoteID == nil && len(cmd.Items) == 0 {
		return nil, domain.NewValidationError("items cannot be empty")
	}
	if err := s.ensurePartyExists(ctx, cmd.PartyID); err != nil {
		return nil, err
	}
	if s.numberGen == nil {
		return nil, fmt.Errorf("order number generator not configured")
	}

	orderNumber, err := s.numberGen.NextOrderNumber(ctx)
	if err != nil {
		return nil, err
	}

	if cmd.QuoteID != nil {
		quote, err := s.quoteRepo.FindByID(ctx, *cmd.QuoteID)
		if err != nil {
			return nil, err
		}
		if quote == nil {
			return nil, domain.NewNotFoundError("quote not found")
		}
		if quote.Status != domain.QuoteStatusApproved {
			return nil, domain.NewConflictError("quote must be approved before creating an order")
		}
		order, err := quote.ConvertToOrder(orderNumber, cmd.DeliveryDate)
		if err != nil {
			return nil, err
		}
		if err := s.orderRepo.Save(ctx, order); err != nil {
			return nil, err
		}
		if err := s.quoteRepo.Save(ctx, quote); err != nil {
			return nil, err
		}
		return NewSalesOrderDTO(order), nil
	}

	lineItems, err := s.buildOrderLineItems(ctx, cmd.PartyID, cmd.Items, nil)
	if err != nil {
		return nil, err
	}

	taxAmount, err := zeroMoney()
	if err != nil {
		return nil, err
	}
	for _, lineItem := range lineItems {
		taxAmount, err = taxAmount.Add(lineItem.TaxAmount)
		if err != nil {
			return nil, err
		}
	}
	notes := ""
	if cmd.Notes != nil {
		notes = *cmd.Notes
	}

	order, err := domain.NewSalesOrder(
		orderNumber,
		cmd.PartyID,
		time.Now(),
		cmd.DeliveryDate,
		lineItems,
		taxAmount,
		notes,
	)
	if err != nil {
		return nil, err
	}

	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	return NewSalesOrderDTO(order), nil
}

func (s *SalesService) UpdateOrderDetails(ctx context.Context, cmd UpdateOrderDetailsCommand) (*SalesOrderDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, cmd.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.NewNotFoundError("order not found")
	}
	if !canEditOrderDetails(order.Status) {
		return nil, domain.NewConflictError("order details cannot be updated in current status")
	}

	if cmd.PartyID != nil {
		if err := s.ensurePartyExists(ctx, *cmd.PartyID); err != nil {
			return nil, err
		}
		order.PartyID = *cmd.PartyID
	}
	if cmd.DeliveryDate != nil {
		if cmd.DeliveryDate.Before(order.OrderDate) {
			return nil, domain.NewValidationError("deliveryDate cannot be before orderDate")
		}
		order.DeliveryDate = *cmd.DeliveryDate
	}
	if cmd.Notes != nil {
		order.Notes = *cmd.Notes
	}

	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	return NewSalesOrderDTO(order), nil
}

func (s *SalesService) ChangeOrderStatus(ctx context.Context, cmd ChangeOrderStatusCommand) (*SalesOrderDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, cmd.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.NewNotFoundError("order not found")
	}

	status, err := parseOrderStatus(cmd.NewStatus)
	if err != nil {
		return nil, err
	}
	if err := order.ChangeStatus(status); err != nil {
		return nil, err
	}

	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	return NewSalesOrderDTO(order), nil
}

func (s *SalesService) AddOrderLineItem(ctx context.Context, cmd AddOrderLineItemCommand) (*SalesOrderDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, cmd.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.NewNotFoundError("order not found")
	}
	if !canEditOrderLineItems(order.Status) {
		return nil, domain.NewConflictError("order line items cannot be updated in current status")
	}

	seeds := orderLineItemSeedsFromOrder(order.LineItems)
	seeds = append(seeds, orderLineItemSeed{
		ProductVariantID:      cmd.Item.ProductVariantID,
		Quantity:              cmd.Item.Quantity,
		ManualUnitPrice:       cmd.Item.ManualUnitPrice,
		ManualDiscountPerUnit: cmd.Item.ManualDiscountPerUnit,
	})

	lineItems, err := s.buildOrderLineItemsFromSeeds(ctx, order.PartyID, seeds)
	if err != nil {
		return nil, err
	}

	taxAmount, err := zeroMoney()
	if err != nil {
		return nil, err
	}
	for _, lineItem := range lineItems {
		taxAmount, err = taxAmount.Add(lineItem.TaxAmount)
		if err != nil {
			return nil, err
		}
	}
	order.LineItems = lineItems
	order.TaxAmount = taxAmount
	if err := order.RecalculateTotals(); err != nil {
		return nil, err
	}

	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	return NewSalesOrderDTO(order), nil
}

func (s *SalesService) UpdateOrderLineItem(ctx context.Context, cmd UpdateOrderLineItemCommand) (*SalesOrderDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, cmd.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.NewNotFoundError("order not found")
	}
	if !canEditOrderLineItems(order.Status) {
		return nil, domain.NewConflictError("order line items cannot be updated in current status")
	}

	seeds := orderLineItemSeedsFromOrder(order.LineItems)
	updated := false
	for i := range seeds {
		if seeds[i].ID != nil && *seeds[i].ID == cmd.LineItemID {
			updated = true
			if cmd.Quantity != nil {
				seeds[i].Quantity = *cmd.Quantity
			}
			if cmd.ManualUnitPrice != nil {
				seeds[i].ManualUnitPrice = cmd.ManualUnitPrice
			}
			if cmd.ManualDiscountPerUnit != nil {
				seeds[i].ManualDiscountPerUnit = cmd.ManualDiscountPerUnit
			}
			break
		}
	}
	if !updated {
		return nil, domain.NewNotFoundError("order line item not found")
	}

	lineItems, err := s.buildOrderLineItemsFromSeeds(ctx, order.PartyID, seeds)
	if err != nil {
		return nil, err
	}

	taxAmount, err := zeroMoney()
	if err != nil {
		return nil, err
	}
	order.LineItems = lineItems
	order.TaxAmount = taxAmount
	if err := order.RecalculateTotals(); err != nil {
		return nil, err
	}

	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	return NewSalesOrderDTO(order), nil
}

func (s *SalesService) RemoveOrderLineItem(ctx context.Context, cmd RemoveOrderLineItemCommand) (*SalesOrderDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, cmd.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.NewNotFoundError("order not found")
	}
	if !canEditOrderLineItems(order.Status) {
		return nil, domain.NewConflictError("order line items cannot be updated in current status")
	}

	seeds := orderLineItemSeedsFromOrder(order.LineItems)
	filtered := make([]orderLineItemSeed, 0, len(seeds))
	for _, seed := range seeds {
		if seed.ID != nil && *seed.ID == cmd.LineItemID {
			continue
		}
		filtered = append(filtered, seed)
	}
	if len(filtered) == len(seeds) {
		return nil, domain.NewNotFoundError("order line item not found")
	}
	if len(filtered) == 0 {
		return nil, domain.NewValidationError("order must have at least one line item")
	}

	lineItems, err := s.buildOrderLineItemsFromSeeds(ctx, order.PartyID, filtered)
	if err != nil {
		return nil, err
	}

	taxAmount, err := zeroMoney()
	if err != nil {
		return nil, err
	}
	order.LineItems = lineItems
	order.TaxAmount = taxAmount
	if err := order.RecalculateTotals(); err != nil {
		return nil, err
	}

	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	return NewSalesOrderDTO(order), nil
}

func (s *SalesService) CreateDeliveryNote(ctx context.Context, cmd CreateDeliveryNoteCommand) (*DeliveryNoteDTO, error) {
	if cmd.SalesOrderID == uuid.Nil {
		return nil, domain.NewValidationError("salesOrderId is required")
	}
	if cmd.DeliveryDate.IsZero() {
		return nil, domain.NewValidationError("deliveryDate is required")
	}
	if len(cmd.Items) == 0 {
		return nil, domain.NewValidationError("items cannot be empty")
	}
	if s.numberGen == nil {
		return nil, fmt.Errorf("delivery note number generator not configured")
	}

	order, err := s.orderRepo.FindByID(ctx, cmd.SalesOrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.NewNotFoundError("order not found")
	}
	if order.Status == domain.SalesOrderStatusCanceled {
		return nil, domain.NewConflictError("cannot create delivery note for canceled order")
	}
	if order.Status == domain.SalesOrderStatusInvoiced || order.Status == domain.SalesOrderStatusPartiallyInvoiced {
		return nil, domain.NewConflictError("cannot create delivery note for invoiced order")
	}

	alreadyDelivered, err := s.deliveredQuantities(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	lineItems := make([]domain.DeliveryNoteLineItem, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		orderLine := findOrderLineItem(order.LineItems, item.SalesOrderLineItemID)
		if orderLine == nil {
			return nil, domain.NewValidationError("sales order line item not found")
		}
		previous := alreadyDelivered[item.SalesOrderLineItemID]
		if previous+item.DeliveredQuantity > orderLine.Quantity {
			return nil, domain.NewValidationError("delivered quantity exceeds ordered quantity")
		}
		created, err := domain.NewDeliveryNoteLineItem(item.SalesOrderLineItemID, orderLine.ProductVariantID, item.DeliveredQuantity)
		if err != nil {
			return nil, err
		}
		lineItems = append(lineItems, created)
	}

	noteNumber, err := s.numberGen.NextDeliveryNoteNumber(ctx)
	if err != nil {
		return nil, err
	}

	notes := ""
	if cmd.Notes != nil {
		notes = *cmd.Notes
	}

	note, err := domain.NewDeliveryNote(
		noteNumber,
		order.ID,
		order.PartyID,
		cmd.DeliveryDate,
		lineItems,
		notes,
	)
	if err != nil {
		return nil, err
	}

	deliveredAll := isOrderFullyDelivered(order.LineItems, alreadyDelivered, lineItems)
	if order.Status == domain.SalesOrderStatusPending {
		if err := order.ChangeStatus(domain.SalesOrderStatusInPreparation); err != nil {
			return nil, err
		}
	}
	if deliveredAll {
		if order.Status != domain.SalesOrderStatusDelivered {
			if err := order.ChangeStatus(domain.SalesOrderStatusDelivered); err != nil {
				return nil, err
			}
		}
		if err := note.ChangeStatus(domain.DeliveryNoteStatusDelivered); err != nil {
			return nil, err
		}
	} else {
		if order.Status != domain.SalesOrderStatusPartiallyDelivered {
			if err := order.ChangeStatus(domain.SalesOrderStatusPartiallyDelivered); err != nil {
				return nil, err
			}
		}
	}

	if err := s.deliveryRepo.Save(ctx, note); err != nil {
		return nil, err
	}
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	return NewDeliveryNoteDTO(note), nil
}

func (s *SalesService) CreateInvoice(ctx context.Context, cmd CreateInvoiceCommand) (*InvoiceDTO, error) {
	if cmd.PartyID == uuid.Nil {
		return nil, domain.NewValidationError("partyId is required")
	}
	if cmd.InvoiceDate.IsZero() {
		return nil, domain.NewValidationError("invoiceDate is required")
	}
	if cmd.DueDate.IsZero() {
		return nil, domain.NewValidationError("dueDate is required")
	}
	if len(cmd.SalesOrderIDs) == 0 && len(cmd.DeliveryNoteIDs) == 0 {
		return nil, domain.NewValidationError("salesOrderIds or deliveryNoteIds must be provided")
	}
	if len(cmd.SalesOrderIDs) > 0 && len(cmd.DeliveryNoteIDs) > 0 {
		return nil, domain.NewValidationError("provide either salesOrderIds or deliveryNoteIds, not both")
	}
	if s.numberGen == nil {
		return nil, fmt.Errorf("invoice number generator not configured")
	}

	lineItems := make([]domain.InvoiceLineItem, 0)
	relatedOrders := make(map[uuid.UUID]struct{})

	orders, err := s.fetchOrdersForInvoice(ctx, cmd.PartyID, cmd.SalesOrderIDs)
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		for _, item := range order.LineItems {
			lineItem, err := buildInvoiceLineItemFromOrder(item, item.Quantity)
			if err != nil {
				return nil, err
			}
			lineItems = append(lineItems, lineItem)
		}
		relatedOrders[order.ID] = struct{}{}
	}

	if len(cmd.DeliveryNoteIDs) > 0 {
		noteItems, noteOrders, err := s.buildInvoiceItemsFromDeliveryNotes(ctx, cmd.PartyID, cmd.DeliveryNoteIDs)
		if err != nil {
			return nil, err
		}
		lineItems = append(lineItems, noteItems...)
		for _, orderID := range noteOrders {
			relatedOrders[orderID] = struct{}{}
		}
	}
	if len(lineItems) == 0 {
		return nil, domain.NewValidationError("invoice must have at least one line item")
	}

	taxAmount, err := zeroMoney()
	if err != nil {
		return nil, err
	}
	paymentTerms := ""
	if cmd.PaymentTerms != nil {
		paymentTerms = *cmd.PaymentTerms
	}

	invoiceNumber, err := s.numberGen.NextInvoiceNumber(ctx)
	if err != nil {
		return nil, err
	}

	// Default for invoices from orders: COMPLETA type with series "A"
	invoiceType := domain.InvoiceTypeComplete
	currentYear := time.Now().Year()
	series, err := domain.NewInvoiceSeries("A", currentYear)
	if err != nil {
		return nil, err
	}

	invoice, err := domain.NewInvoice(
		invoiceNumber,
		invoiceType,
		series,
		cmd.PartyID,
		cmd.InvoiceDate,
		cmd.DueDate,
		lineItems,
		taxAmount,
		paymentTerms,
	)
	if err != nil {
		return nil, err
	}

	for orderID := range relatedOrders {
		order, err := s.orderRepo.FindByID(ctx, orderID)
		if err != nil {
			return nil, err
		}
		if order == nil {
			return nil, domain.NewNotFoundError("order not found")
		}

		if err := s.updateOrderInvoiceStatus(ctx, order, lineItems); err != nil {
			return nil, err
		}
		if err := s.orderRepo.Save(ctx, order); err != nil {
			return nil, err
		}
	}

	if err := s.invoiceRepo.Save(ctx, invoice); err != nil {
		return nil, err
	}

	relatedIDs := make([]uuid.UUID, 0, len(relatedOrders))
	for id := range relatedOrders {
		relatedIDs = append(relatedIDs, id)
	}

	return NewInvoiceDTO(invoice, relatedIDs), nil
}

func (s *SalesService) GetQuote(ctx context.Context, query GetQuoteByIDQuery) (*QuoteDTO, error) {
	quote, err := s.quoteRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, domain.NewNotFoundError("quote not found")
	}
	return NewQuoteDTO(quote), nil
}

func (s *SalesService) ListQuotes(ctx context.Context, query ListQuotesQuery) ([]*QuoteDTO, error) {
	filter := domain.QuoteFilter{PartyID: query.PartyID, FromDate: query.FromDate, ToDate: query.ToDate, Search: query.Search}
	if query.Status != nil {
		status, err := parseQuoteStatus(*query.Status)
		if err != nil {
			return nil, err
		}
		filter.Status = &status
	}
	quotes, err := s.quoteRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	result := make([]*QuoteDTO, 0, len(quotes))
	for _, quote := range quotes {
		result = append(result, NewQuoteDTO(quote))
	}
	return result, nil
}

func (s *SalesService) GetOrder(ctx context.Context, query GetOrderByIDQuery) (*SalesOrderDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.NewNotFoundError("order not found")
	}
	return NewSalesOrderDTO(order), nil
}

func (s *SalesService) ListOrders(ctx context.Context, query ListOrdersQuery) ([]*SalesOrderDTO, error) {
	filter := domain.SalesOrderFilter{PartyID: query.PartyID, FromDate: query.FromDate, ToDate: query.ToDate, Search: query.Search}
	if query.Status != nil {
		status, err := parseOrderStatus(*query.Status)
		if err != nil {
			return nil, err
		}
		filter.Status = &status
	}
	orders, err := s.orderRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	result := make([]*SalesOrderDTO, 0, len(orders))
	for _, order := range orders {
		result = append(result, NewSalesOrderDTO(order))
	}
	return result, nil
}

func (s *SalesService) GetDeliveryNote(ctx context.Context, query GetDeliveryNoteByIDQuery) (*DeliveryNoteDTO, error) {
	note, err := s.deliveryRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, domain.NewNotFoundError("delivery note not found")
	}
	return NewDeliveryNoteDTO(note), nil
}

func (s *SalesService) ListDeliveryNotes(ctx context.Context, query ListDeliveryNotesQuery) ([]*DeliveryNoteDTO, error) {
	filter := domain.DeliveryNoteFilter{
		SalesOrderID: query.SalesOrderID,
		PartyID:      query.PartyID,
		FromDate:     query.FromDate,
		ToDate:       query.ToDate,
		Search:       query.Search,
	}
	if query.Status != nil {
		status, err := parseDeliveryNoteStatus(*query.Status)
		if err != nil {
			return nil, err
		}
		filter.Status = &status
	}
	notes, err := s.deliveryRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	result := make([]*DeliveryNoteDTO, 0, len(notes))
	for _, note := range notes {
		result = append(result, NewDeliveryNoteDTO(note))
	}
	return result, nil
}

func (s *SalesService) GetInvoice(ctx context.Context, query GetInvoiceByIDQuery) (*InvoiceDTO, error) {
	invoice, err := s.invoiceRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, domain.NewNotFoundError("invoice not found")
	}
	return NewInvoiceDTO(invoice, nil), nil
}

func (s *SalesService) ListInvoices(ctx context.Context, query ListInvoicesQuery) ([]*InvoiceDTO, error) {
	filter := domain.InvoiceFilter{PartyID: query.PartyID, FromDate: query.FromDate, ToDate: query.ToDate, Search: query.Search}
	if query.Status != nil {
		status, err := parseInvoiceStatus(*query.Status)
		if err != nil {
			return nil, err
		}
		filter.Status = &status
	}
	invoices, err := s.invoiceRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	result := make([]*InvoiceDTO, 0, len(invoices))
	for _, invoice := range invoices {
		result = append(result, NewInvoiceDTO(invoice, nil))
	}
	return result, nil
}

func (s *SalesService) ensurePartyExists(ctx context.Context, partyID uuid.UUID) error {
	if s.partyLookup == nil {
		return nil
	}
	exists, err := s.partyLookup.ExistsParty(ctx, partyID)
	if err != nil {
		return err
	}
	if !exists {
		return domain.NewNotFoundError("party not found")
	}
	isClient, err := s.partyLookup.HasPartyRole(ctx, partyID, "CLIENT")
	if err != nil {
		return err
	}
	if !isClient {
		return domain.NewValidationError("party must have CLIENT role to sell")
	}
	return nil
}

type orderLineItemSeed struct {
	ID                    *uuid.UUID
	MesWorkID             *uuid.UUID
	ProductVariantID      uuid.UUID
	Quantity              int
	ManualUnitPrice       *MoneyDTO
	ManualDiscountPerUnit *MoneyDTO
}

func orderLineItemSeedsFromOrder(items []domain.OrderLineItem) []orderLineItemSeed {
	seeds := make([]orderLineItemSeed, 0, len(items))
	for _, item := range items {
		id := item.ID
		seeds = append(seeds, orderLineItemSeed{
			ID:                    &id,
			MesWorkID:             item.MESWorkID,
			ProductVariantID:      item.ProductVariantID,
			Quantity:              item.Quantity,
			ManualUnitPrice:       toMoneyDTOPtr(item.ManualUnitPrice),
			ManualDiscountPerUnit: toMoneyDTOPtr(item.ManualDiscountPerUnit),
		})
	}
	return seeds
}

func (s *SalesService) buildQuoteLineItems(ctx context.Context, partyID uuid.UUID, items []QuoteLineItemInput) ([]domain.QuoteLineItem, error) {
	if s.pricingEngine == nil {
		return nil, fmt.Errorf("pricing engine not configured")
	}
	request := pricing_app.CalculateFinalSalePriceRequest{
		ClientID:  partyID,
		SaleItems: make([]pricing_app.SaleItemRequest, len(items)),
	}
	for i, item := range items {
		if item.ProductVariantID == uuid.Nil {
			return nil, domain.NewValidationError("productVariantId is required")
		}
		if item.Quantity <= 0 {
			return nil, domain.NewValidationError("quantity must be greater than zero")
		}
		request.SaleItems[i] = pricing_app.SaleItemRequest{
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
		}
	}

	pricing, err := s.pricingEngine.CalculateFinalSalePrice(ctx, request)
	if err != nil {
		return nil, err
	}
	if len(pricing.CalculatedItems) != len(items) {
		return nil, fmt.Errorf("pricing response mismatch")
	}

	lineItems := make([]domain.QuoteLineItem, 0, len(items))
	for i, item := range items {
		calculatedUnit, calculatedDiscount, err := deriveCalculatedPrices(pricing.CalculatedItems[i])
		if err != nil {
			return nil, err
		}
		manualUnit, err := toDomainMoneyPtr(item.ManualUnitPrice)
		if err != nil {
			return nil, err
		}
		manualDiscount, err := toDomainMoneyPtr(item.ManualDiscountPerUnit)
		if err != nil {
			return nil, err
		}
		lineItem, err := domain.NewQuoteLineItem(
			item.ProductVariantID,
			item.Quantity,
			calculatedUnit,
			manualUnit,
			calculatedDiscount,
			manualDiscount,
			pricing.CalculatedItems[i].TaxRate,
		)
		if err != nil {
			return nil, err
		}
		lineItem.MESWorkID = item.MesWorkID
		lineItems = append(lineItems, lineItem)
	}
	return lineItems, nil
}

func (s *SalesService) buildOrderLineItems(ctx context.Context, partyID uuid.UUID, items []OrderLineItemInput, existingIDs []uuid.UUID) ([]domain.OrderLineItem, error) {
	seeds := make([]orderLineItemSeed, 0, len(items))
	for i, item := range items {
		var id *uuid.UUID
		if i < len(existingIDs) && existingIDs[i] != uuid.Nil {
			value := existingIDs[i]
			id = &value
		}
		seeds = append(seeds, orderLineItemSeed{
			ID:                    id,
			MesWorkID:             item.MesWorkID,
			ProductVariantID:      item.ProductVariantID,
			Quantity:              item.Quantity,
			ManualUnitPrice:       item.ManualUnitPrice,
			ManualDiscountPerUnit: item.ManualDiscountPerUnit,
		})
	}
	return s.buildOrderLineItemsFromSeeds(ctx, partyID, seeds)
}

func (s *SalesService) buildOrderLineItemsFromSeeds(ctx context.Context, partyID uuid.UUID, seeds []orderLineItemSeed) ([]domain.OrderLineItem, error) {
	if s.pricingEngine == nil {
		return nil, fmt.Errorf("pricing engine not configured")
	}
	request := pricing_app.CalculateFinalSalePriceRequest{
		ClientID:  partyID,
		SaleItems: make([]pricing_app.SaleItemRequest, len(seeds)),
	}
	for i, seed := range seeds {
		if seed.ProductVariantID == uuid.Nil {
			return nil, domain.NewValidationError("productVariantId is required")
		}
		if seed.Quantity <= 0 {
			return nil, domain.NewValidationError("quantity must be greater than zero")
		}
		request.SaleItems[i] = pricing_app.SaleItemRequest{
			ProductVariantID: seed.ProductVariantID,
			Quantity:         seed.Quantity,
		}
	}

	pricing, err := s.pricingEngine.CalculateFinalSalePrice(ctx, request)
	if err != nil {
		return nil, err
	}
	if len(pricing.CalculatedItems) != len(seeds) {
		return nil, fmt.Errorf("pricing response mismatch")
	}

	lineItems := make([]domain.OrderLineItem, 0, len(seeds))
	for i, seed := range seeds {
		calculatedUnit, calculatedDiscount, err := deriveCalculatedPrices(pricing.CalculatedItems[i])
		if err != nil {
			return nil, err
		}
		manualUnit, err := toDomainMoneyPtr(seed.ManualUnitPrice)
		if err != nil {
			return nil, err
		}
		manualDiscount, err := toDomainMoneyPtr(seed.ManualDiscountPerUnit)
		if err != nil {
			return nil, err
		}

		lineItem, err := domain.NewOrderLineItem(
			seed.ProductVariantID,
			seed.Quantity,
			calculatedUnit,
			manualUnit,
			calculatedDiscount,
			manualDiscount,
			pricing.CalculatedItems[i].TaxRate,
		)
		if err != nil {
			return nil, err
		}
		if seed.ID != nil {
			lineItem.ID = *seed.ID
		}
		lineItem.MESWorkID = seed.MesWorkID
		lineItems = append(lineItems, lineItem)
	}

	return lineItems, nil
}

func deriveCalculatedPrices(item pricing_app.CalculatedSaleItemResponse) (domain.Money, *domain.Money, error) {
	base, err := domain.NewMoney(item.BaseSalesPrice.Amount, item.BaseSalesPrice.Currency)
	if err != nil {
		return domain.Money{}, nil, err
	}
	final, err := domain.NewMoney(item.FinalPrice.Amount, item.FinalPrice.Currency)
	if err != nil {
		return domain.Money{}, nil, err
	}
	if final.Amount() <= base.Amount() {
		discount, err := base.Subtract(final)
		if err != nil {
			return domain.Money{}, nil, err
		}
		if discount.Amount() > 0 {
			return base, &discount, nil
		}
		return base, nil, nil
	}
	return final, nil, nil
}

func toDomainMoneyPtr(dto *MoneyDTO) (*domain.Money, error) {
	if dto == nil {
		return nil, nil
	}
	money, err := domain.NewMoney(dto.Amount, dto.Currency)
	if err != nil {
		return nil, err
	}
	return &money, nil
}

func zeroMoney() (domain.Money, error) {
	return domain.NewMoney(0, domain.DefaultCurrency)
}

func sumQuoteLineItemTaxes(items []domain.QuoteLineItem) (domain.Money, error) {
	total, err := zeroMoney()
	if err != nil {
		return domain.Money{}, err
	}
	for _, item := range items {
		total, err = total.Add(item.TaxAmount)
		if err != nil {
			return domain.Money{}, err
		}
	}
	return total, nil
}

func sumOrderLineItemTaxes(items []domain.OrderLineItem) (domain.Money, error) {
	total, err := zeroMoney()
	if err != nil {
		return domain.Money{}, err
	}
	for _, item := range items {
		total, err = total.Add(item.TaxAmount)
		if err != nil {
			return domain.Money{}, err
		}
	}
	return total, nil
}

func sumInvoiceLineItemTaxes(items []domain.InvoiceLineItem) (domain.Money, error) {
	total, err := zeroMoney()
	if err != nil {
		return domain.Money{}, err
	}
	for _, item := range items {
		if item.TaxAmount == nil {
			continue
		}
		total, err = total.Add(*item.TaxAmount)
		if err != nil {
			return domain.Money{}, err
		}
	}
	return total, nil
}

func parseQuoteStatus(input string) (domain.QuoteStatus, error) {
	value := domain.QuoteStatus(strings.ToUpper(strings.TrimSpace(input)))
	return value, value.IsValid()
}

func parseOrderStatus(input string) (domain.SalesOrderStatus, error) {
	value := domain.SalesOrderStatus(strings.ToUpper(strings.TrimSpace(input)))
	return value, value.IsValid()
}

func parseDeliveryNoteStatus(input string) (domain.DeliveryNoteStatus, error) {
	value := domain.DeliveryNoteStatus(strings.ToUpper(strings.TrimSpace(input)))
	return value, value.IsValid()
}

func parseInvoiceStatus(input string) (domain.InvoiceStatus, error) {
	value := domain.InvoiceStatus(strings.ToUpper(strings.TrimSpace(input)))
	return value, value.IsValid()
}

func canEditOrderDetails(status domain.SalesOrderStatus) bool {
	switch status {
	case domain.SalesOrderStatusPending, domain.SalesOrderStatusInPreparation:
		return true
	default:
		return false
	}
}

func canEditOrderLineItems(status domain.SalesOrderStatus) bool {
	return status == domain.SalesOrderStatusPending
}

func findOrderLineItem(items []domain.OrderLineItem, lineItemID uuid.UUID) *domain.OrderLineItem {
	for i := range items {
		if items[i].ID == lineItemID {
			return &items[i]
		}
	}
	return nil
}

func (s *SalesService) deliveredQuantities(ctx context.Context, orderID uuid.UUID) (map[uuid.UUID]int, error) {
	results := make(map[uuid.UUID]int)
	notes, err := s.deliveryRepo.ListBySalesOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	for _, note := range notes {
		for _, item := range note.LineItems {
			results[item.SalesOrderLineItemID] += item.DeliveredQuantity
		}
	}
	return results, nil
}

func isOrderFullyDelivered(items []domain.OrderLineItem, previous map[uuid.UUID]int, newItems []domain.DeliveryNoteLineItem) bool {
	current := make(map[uuid.UUID]int)
	for id, qty := range previous {
		current[id] = qty
	}
	for _, item := range newItems {
		current[item.SalesOrderLineItemID] += item.DeliveredQuantity
	}
	for _, item := range items {
		if current[item.ID] < item.Quantity {
			return false
		}
	}
	return true
}

func buildInvoiceLineItemFromOrder(item domain.OrderLineItem, quantity int) (domain.InvoiceLineItem, error) {
	var discount *domain.Money
	if item.FinalDiscountPerUnit.Amount() > 0 {
		value := item.FinalDiscountPerUnit
		discount = &value
	}
	lineItem, err := domain.NewInvoiceLineItem(
		item.ProductVariantID,
		quantity,
		item.FinalUnitPrice,
		discount,
		nil,
		item.TaxRate,
	)
	if err != nil {
		return domain.InvoiceLineItem{}, err
	}
	lineItem.SalesOrderLineItemID = &item.ID
	return lineItem, nil
}

func (s *SalesService) fetchOrdersForInvoice(ctx context.Context, partyID uuid.UUID, orderIDs []uuid.UUID) ([]*domain.SalesOrder, error) {
	orders := make([]*domain.SalesOrder, 0, len(orderIDs))
	for _, orderID := range orderIDs {
		order, err := s.orderRepo.FindByID(ctx, orderID)
		if err != nil {
			return nil, err
		}
		if order == nil {
			return nil, domain.NewNotFoundError("order not found")
		}
		if order.PartyID != partyID {
			return nil, domain.NewValidationError("order party mismatch")
		}
		if order.Status == domain.SalesOrderStatusCanceled {
			return nil, domain.NewConflictError("cannot invoice canceled order")
		}
		if order.Status != domain.SalesOrderStatusDelivered && order.Status != domain.SalesOrderStatusPartiallyInvoiced && order.Status != domain.SalesOrderStatusInvoiced {
			return nil, domain.NewConflictError("order must be delivered before invoicing")
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (s *SalesService) buildInvoiceItemsFromDeliveryNotes(ctx context.Context, partyID uuid.UUID, noteIDs []uuid.UUID) ([]domain.InvoiceLineItem, []uuid.UUID, error) {
	lineItems := make([]domain.InvoiceLineItem, 0)
	orderIDs := make([]uuid.UUID, 0)
	seen := make(map[uuid.UUID]struct{})

	for _, noteID := range noteIDs {
		note, err := s.deliveryRepo.FindByID(ctx, noteID)
		if err != nil {
			return nil, nil, err
		}
		if note == nil {
			return nil, nil, domain.NewNotFoundError("delivery note not found")
		}
		order, err := s.orderRepo.FindByID(ctx, note.SalesOrderID)
		if err != nil {
			return nil, nil, err
		}
		if order == nil {
			return nil, nil, domain.NewNotFoundError("order not found")
		}
		if order.PartyID != partyID {
			return nil, nil, domain.NewValidationError("delivery note party mismatch")
		}
		for _, item := range note.LineItems {
			orderLine := findOrderLineItem(order.LineItems, item.SalesOrderLineItemID)
			if orderLine == nil {
				return nil, nil, domain.NewValidationError("sales order line item not found")
			}
			lineItem, err := buildInvoiceLineItemFromOrder(*orderLine, item.DeliveredQuantity)
			if err != nil {
				return nil, nil, err
			}
			lineItems = append(lineItems, lineItem)
		}
		if _, ok := seen[order.ID]; !ok {
			orderIDs = append(orderIDs, order.ID)
			seen[order.ID] = struct{}{}
		}
	}

	return lineItems, orderIDs, nil
}

func (s *SalesService) updateOrderInvoiceStatus(ctx context.Context, order *domain.SalesOrder, newInvoiceItems []domain.InvoiceLineItem) error {
	if order.Status == domain.SalesOrderStatusCanceled {
		return domain.NewConflictError("cannot update invoice status for canceled order")
	}
	if order.Status != domain.SalesOrderStatusDelivered && order.Status != domain.SalesOrderStatusPartiallyInvoiced && order.Status != domain.SalesOrderStatusInvoiced {
		return domain.NewConflictError("order must be delivered before invoicing")
	}

	invoiced := make(map[uuid.UUID]int)
	existing, err := s.invoiceRepo.ListBySalesOrderID(ctx, order.ID)
	if err != nil {
		return err
	}
	for _, invoice := range existing {
		for _, item := range invoice.LineItems {
			if item.SalesOrderLineItemID != nil {
				invoiced[*item.SalesOrderLineItemID] += item.Quantity
			}
		}
	}
	for _, item := range newInvoiceItems {
		if item.SalesOrderLineItemID != nil {
			invoiced[*item.SalesOrderLineItemID] += item.Quantity
		}
	}

	fullyInvoiced := true
	for _, item := range order.LineItems {
		if invoiced[item.ID] < item.Quantity {
			fullyInvoiced = false
			break
		}
	}

	if order.Status == domain.SalesOrderStatusInvoiced {
		return nil
	}

	var target domain.SalesOrderStatus
	if fullyInvoiced {
		target = domain.SalesOrderStatusInvoiced
	} else {
		target = domain.SalesOrderStatusPartiallyInvoiced
	}

	if order.Status == target {
		return nil
	}
	if err := order.ChangeStatus(target); err != nil {
		return err
	}
	return nil
}

// CreateSimplifiedInvoice (CU-S-019) creates a ticket (factura simplificada) for retail sales
// Optimized for TPV/POS workflow: validates < 3,000 EUR limit, uses series "TKT", allows CONSUMIDOR_FINAL
func (s *SalesService) CreateSimplifiedInvoice(ctx context.Context, cmd CreateSimplifiedInvoiceCommand) (*InvoiceDTO, error) {
	if err := s.ensurePartyExists(ctx, cmd.PartyID); err != nil {
		return nil, err
	}

	// Build sale items for pricing calculation
	saleItems := make([]pricing_app.SaleItemRequest, 0, len(cmd.Items))
	for _, itemInput := range cmd.Items {
		saleItems = append(saleItems, pricing_app.SaleItemRequest{
			ProductVariantID: itemInput.ProductVariantID,
			Quantity:         itemInput.Quantity,
		})
	}

	// Calculate prices using pricing engine
	priceReq := pricing_app.CalculateFinalSalePriceRequest{
		SaleItems: saleItems,
		ClientID:  cmd.PartyID,
		SaleDate:  cmd.InvoiceDate,
	}
	priceResp, err := s.pricingEngine.CalculateFinalSalePrice(ctx, priceReq)
	if err != nil {
		return nil, fmt.Errorf("pricing calculation failed: %w", err)
	}

	// Build invoice line items from calculated prices
	lineItems := make([]domain.InvoiceLineItem, 0, len(priceResp.CalculatedItems))
	for _, calculatedItem := range priceResp.CalculatedItems {
		unitPrice, err := domain.NewMoney(calculatedItem.FinalPrice.Amount, domain.DefaultCurrency)
		if err != nil {
			return nil, err
		}

		// For tickets, we use direct pricing without manual overrides
		lineItem, err := domain.NewInvoiceLineItem(
			calculatedItem.ProductVariantID,
			calculatedItem.Quantity,
			unitPrice,
			nil, // No discount
			nil, // No tax breakdown per line
			calculatedItem.TaxRate,
		)
		if err != nil {
			return nil, err
		}

		lineItems = append(lineItems, lineItem)
	}

	// Simplified invoice: immediate payment, no payment terms
	invoiceNumber, err := s.numberGen.NextInvoiceNumber(ctx)
	if err != nil {
		return nil, err
	}

	// Ticket series: "TKT" for current year
	invoiceType := domain.InvoiceTypeSimplified
	currentYear := time.Now().Year()
	series, err := domain.NewInvoiceSeries("TKT", currentYear)
	if err != nil {
		return nil, err
	}

	taxAmount, err := zeroMoney()
	if err != nil {
		return nil, err
	}
	for _, lineItem := range lineItems {
		if lineItem.TaxAmount == nil {
			continue
		}
		taxAmount, err = taxAmount.Add(*lineItem.TaxAmount)
		if err != nil {
			return nil, err
		}
	}

	// Create invoice - will validate < 3,000 EUR automatically via ValidateLegalLimits()
	invoice, err := domain.NewInvoice(
		invoiceNumber,
		invoiceType,
		series,
		cmd.PartyID,
		cmd.InvoiceDate,
		cmd.InvoiceDate, // Due date = invoice date for immediate payment
		lineItems,
		taxAmount,
		"Immediate Payment", // Payment terms for tickets
	)
	if err != nil {
		return nil, err
	}

	if err := s.invoiceRepo.Save(ctx, invoice); err != nil {
		return nil, err
	}

	// No related orders for simplified invoices (direct sale)
	return NewInvoiceDTO(invoice, []uuid.UUID{}), nil
}
