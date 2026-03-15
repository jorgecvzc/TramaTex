package application

import (
	"context"
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

type VariantInfo struct {
	ProductName         string
	VariantSKU          string
	OptionConfiguration map[string]string
}

type ProductVariantLookup interface {
	GetVariantInfo(ctx context.Context, variantID uuid.UUID) (*VariantInfo, error)
}

// WorkOrderProgress is a Sales-local DTO representing the execution state
// of a MES WorkOrder. Sales does not know MES internals — it receives
// a pre-computed progress snapshot.
type WorkOrderProgress struct {
	WorkOrderID    uuid.UUID
	OrderNumber    string
	OrderName      string
	Status         string // MES production status (opaque to Sales)
	TotalTasks     int
	CompletedTasks int
	Lines          []WorkOrderLineProgress
}

// WorkOrderLineProgress represents progress of one line (work type at a position)
// within a WorkOrder. Sales treats this as read-only information from MES.
type WorkOrderLineProgress struct {
	WorkTypeID     uuid.UUID
	PositionID     uuid.UUID
	TotalTasks     int
	CompletedTasks int
}

// MESWorkLookup provides read-only access to MES WorkOrder execution state.
// Implementation lives in infrastructure as an adapter that calls MES service.
type MESWorkLookup interface {
	GetWorkOrderProgress(ctx context.Context, workOrderID uuid.UUID) (*WorkOrderProgress, error)
	GetWorkOrdersProgress(ctx context.Context, workOrderIDs []uuid.UUID) ([]WorkOrderProgress, error)
}

type DocumentNumberGenerator interface {
	NextQuoteNumber(ctx context.Context) (domain.QuoteNumber, error)
	NextOrderNumber(ctx context.Context) (domain.OrderNumber, error)
	NextDeliveryNoteNumber(ctx context.Context) (domain.DeliveryNoteNumber, error)
	NextInvoiceNumber(ctx context.Context, series domain.InvoiceSeries) (domain.InvoiceNumber, error)
}

type SalesService struct {
	quoteRepo     domain.QuoteRepository
	orderRepo     domain.SalesOrderRepository
	deliveryRepo  domain.DeliveryNoteRepository
	invoiceRepo   domain.InvoiceRepository
	numberGen     DocumentNumberGenerator
	pricingEngine PricingEngine
	partyLookup   PartyLookup
	productLookup ProductVariantLookup
	mesLookup     MESWorkLookup
	txManager     TransactionManager
}

func NewSalesService(
	quoteRepo domain.QuoteRepository,
	orderRepo domain.SalesOrderRepository,
	deliveryRepo domain.DeliveryNoteRepository,
	invoiceRepo domain.InvoiceRepository,
	numberGen DocumentNumberGenerator,
	pricingEngine PricingEngine,
	partyLookup PartyLookup,
	productLookup ProductVariantLookup,
	mesLookup MESWorkLookup,
) *SalesService {
	return &SalesService{
		quoteRepo:     quoteRepo,
		orderRepo:     orderRepo,
		deliveryRepo:  deliveryRepo,
		invoiceRepo:   invoiceRepo,
		numberGen:     numberGen,
		pricingEngine: pricingEngine,
		partyLookup:   partyLookup,
		productLookup: productLookup,
		mesLookup:     mesLookup,
	}
}

// SetTransactionManager configures service-level transaction support.
func (s *SalesService) SetTransactionManager(txManager TransactionManager) {
	s.txManager = txManager
}

// runInTransaction wraps fn in a DB transaction if a TransactionManager is configured.
// If no TransactionManager is set (e.g., in tests), fn runs directly.
func (s *SalesService) runInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if s.txManager != nil {
		return s.txManager.RunInTransaction(ctx, fn)
	}
	return fn(ctx)
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
		return nil, domain.NewConfigurationError("quote number generator not configured")
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
	quote.SalesWorkSetups = mesWorkRefsToDomain(cmd.MesWorkRefs)

	if err := s.quoteRepo.Save(ctx, quote); err != nil {
		return nil, err
	}

	dto := NewQuoteDTO(quote)
	s.enrichQuoteLineItems(ctx, dto.LineItems)
	return dto, nil
}

func (s *SalesService) UpdateQuote(ctx context.Context, cmd UpdateQuoteCommand) (*QuoteDTO, error) {
	quote, err := s.quoteRepo.FindByID(ctx, cmd.QuoteID)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, domain.NewNotFoundError("quote not found")
	}
	if quote.Status != domain.QuoteStatusDraft && quote.Status != domain.QuoteStatusIssued {
		return nil, domain.NewConflictError("only draft or issued quotes can be updated")
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
	if cmd.MesWorkRefs != nil {
		quote.SalesWorkSetups = mesWorkRefsToDomain(cmd.MesWorkRefs)
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

func (s *SalesService) PreviewQuoteCalculation(ctx context.Context, cmd PreviewQuoteCommand) (*QuotePreviewDTO, error) {
	if cmd.PartyID == uuid.Nil {
		return nil, domain.NewValidationError("partyId is required")
	}
	if len(cmd.Items) == 0 {
		return nil, domain.NewValidationError("items cannot be empty")
	}

	lineItems, err := s.buildQuoteLineItems(ctx, cmd.PartyID, cmd.Items)
	if err != nil {
		return nil, err
	}

	taxAmount, err := zeroMoney()
	if err != nil {
		return nil, err
	}
	subtotal, err := zeroMoney()
	if err != nil {
		return nil, err
	}
	for _, li := range lineItems {
		taxAmount, err = taxAmount.Add(li.TaxAmount)
		if err != nil {
			return nil, err
		}
		subtotal, err = subtotal.Add(li.Subtotal)
		if err != nil {
			return nil, err
		}
	}

	total, err := subtotal.Add(taxAmount)
	if err != nil {
		return nil, err
	}

	dtoItems := make([]QuoteLineItemDTO, 0, len(lineItems))
	for _, li := range lineItems {
		dto := NewQuoteLineItemDTO(li)
		dto.ProductName, dto.VariantSKU, dto.OptionConfiguration = s.lookupVariant(ctx, li.ProductVariantID)
		dtoItems = append(dtoItems, dto)
	}

	return &QuotePreviewDTO{
		LineItems: dtoItems,
		Subtotal:  NewMoneyDTO(subtotal),
		TaxAmount: NewMoneyDTO(taxAmount),
		Total:     NewMoneyDTO(total),
	}, nil
}

func (s *SalesService) PreviewOrderCalculation(ctx context.Context, cmd PreviewOrderCommand) (*OrderPreviewDTO, error) {
	if cmd.PartyID == uuid.Nil {
		return nil, domain.NewValidationError("partyId is required")
	}
	if len(cmd.Items) == 0 {
		return nil, domain.NewValidationError("items cannot be empty")
	}

	lineItems, err := s.buildOrderLineItems(ctx, cmd.PartyID, cmd.Items, nil)
	if err != nil {
		return nil, err
	}

	taxAmount, err := zeroMoney()
	if err != nil {
		return nil, err
	}
	subtotal, err := zeroMoney()
	if err != nil {
		return nil, err
	}
	for _, li := range lineItems {
		taxAmount, err = taxAmount.Add(li.TaxAmount)
		if err != nil {
			return nil, err
		}
		subtotal, err = subtotal.Add(li.Subtotal)
		if err != nil {
			return nil, err
		}
	}

	total, err := subtotal.Add(taxAmount)
	if err != nil {
		return nil, err
	}

	dtoItems := make([]OrderLineItemDTO, 0, len(lineItems))
	for _, li := range lineItems {
		dto := NewOrderLineItemDTO(li)
		dto.ProductName, dto.VariantSKU, dto.OptionConfiguration = s.lookupVariant(ctx, li.ProductVariantID)
		dtoItems = append(dtoItems, dto)
	}

	return &OrderPreviewDTO{
		LineItems: dtoItems,
		Subtotal:  NewMoneyDTO(subtotal),
		TaxAmount: NewMoneyDTO(taxAmount),
		Total:     NewMoneyDTO(total),
	}, nil
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

func (s *SalesService) DeleteQuote(ctx context.Context, cmd DeleteQuoteCommand) error {
	quote, err := s.quoteRepo.FindByID(ctx, cmd.QuoteID)
	if err != nil {
		return err
	}
	if quote == nil {
		return domain.NewNotFoundError("quote not found")
	}
	if quote.Status != domain.QuoteStatusDraft {
		return domain.NewConflictError("only draft quotes can be deleted")
	}

	return s.quoteRepo.Delete(ctx, cmd.QuoteID)
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
		return nil, domain.NewConfigurationError("order number generator not configured")
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

func (s *SalesService) AcceptAndConvertQuote(ctx context.Context, cmd AcceptAndConvertQuoteCommand) (*SalesOrderDTO, error) {
	quote, err := s.quoteRepo.FindByID(ctx, cmd.QuoteID)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, domain.NewNotFoundError("quote not found")
	}

	// Step 1: Accept the quote (EMITIDA → APROBADA)
	if err := quote.ChangeStatus(domain.QuoteStatusApproved); err != nil {
		return nil, err
	}

	// Step 2: Convert to order (APROBADA → CONVERTIDA_A_PEDIDO + new order)
	if s.numberGen == nil {
		return nil, domain.NewConfigurationError("order number generator not configured")
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
		return nil, domain.NewConfigurationError("order number generator not configured")
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

	now := time.Now()
	orderDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	deliveryDay := time.Date(cmd.DeliveryDate.Year(), cmd.DeliveryDate.Month(), cmd.DeliveryDate.Day(), 0, 0, 0, 0, time.UTC)

	order, err := domain.NewSalesOrder(
		orderNumber,
		cmd.PartyID,
		orderDay,
		deliveryDay,
		lineItems,
		taxAmount,
		notes,
	)
	if err != nil {
		return nil, err
	}
	order.SalesWorkSetups = mesWorkRefsToDomain(cmd.MesWorkRefs)

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
	if cmd.MesWorkRefs != nil {
		order.SalesWorkSetups = mesWorkRefsToDomain(cmd.MesWorkRefs)
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
		ProductVariantID: cmd.Item.ProductVariantID,
		Quantity:         cmd.Item.Quantity,
		UnitPrice:        cmd.Item.UnitPrice,
		DiscountPercent:  cmd.Item.DiscountPercent,
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
			if cmd.UnitPrice != nil {
				seeds[i].UnitPrice = cmd.UnitPrice
			}
			if cmd.DiscountPercent != nil {
				seeds[i].DiscountPercent = cmd.DiscountPercent
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
		return nil, domain.NewConfigurationError("delivery note number generator not configured")
	}

	var result *DeliveryNoteDTO
	err := s.runInTransaction(ctx, func(txCtx context.Context) error {
		// Lock the order row to prevent concurrent delivery note creation
		order, err := s.orderRepo.FindByIDForUpdate(txCtx, cmd.SalesOrderID)
		if err != nil {
			return err
		}
		if order == nil {
			return domain.NewNotFoundError("order not found")
		}
		if order.Status == domain.SalesOrderStatusCanceled {
			return domain.NewConflictError("cannot create delivery note for canceled order")
		}
		if order.Status == domain.SalesOrderStatusInvoiced || order.Status == domain.SalesOrderStatusPartiallyInvoiced {
			return domain.NewConflictError("cannot create delivery note for invoiced order")
		}

		alreadyDelivered, err := s.deliveredQuantities(txCtx, order.ID)
		if err != nil {
			return err
		}

		lineItems := make([]domain.DeliveryNoteLineItem, 0, len(cmd.Items))
		for _, item := range cmd.Items {
			orderLine := findOrderLineItem(order.LineItems, item.SalesOrderLineItemID)
			if orderLine == nil {
				return domain.NewValidationError("sales order line item not found")
			}
			previous := alreadyDelivered[item.SalesOrderLineItemID]
			if previous+item.DeliveredQuantity > orderLine.Quantity {
				return domain.NewValidationError("delivered quantity exceeds ordered quantity")
			}
			created, err := domain.NewDeliveryNoteLineItem(item.SalesOrderLineItemID, orderLine.ProductVariantID, item.DeliveredQuantity)
			if err != nil {
				return err
			}
			lineItems = append(lineItems, created)
		}

		noteNumber, err := s.numberGen.NextDeliveryNoteNumber(txCtx)
		if err != nil {
			return err
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
			return err
		}

		deliveredAll := isOrderFullyDelivered(order.LineItems, alreadyDelivered, lineItems)
		if order.Status == domain.SalesOrderStatusPending {
			if err := order.ChangeStatus(domain.SalesOrderStatusInPreparation); err != nil {
				return err
			}
		}
		if deliveredAll {
			if order.Status != domain.SalesOrderStatusDelivered {
				if err := order.ChangeStatus(domain.SalesOrderStatusDelivered); err != nil {
					return err
				}
			}
			if err := note.ChangeStatus(domain.DeliveryNoteStatusDelivered); err != nil {
				return err
			}
		} else {
			if order.Status != domain.SalesOrderStatusPartiallyDelivered {
				if err := order.ChangeStatus(domain.SalesOrderStatusPartiallyDelivered); err != nil {
					return err
				}
			}
		}

		if err := s.deliveryRepo.Save(txCtx, note); err != nil {
			return err
		}
		if err := s.orderRepo.Save(txCtx, order); err != nil {
			return err
		}

		result = NewDeliveryNoteDTO(note)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
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
		return nil, domain.NewConfigurationError("invoice number generator not configured")
	}

	var result *InvoiceDTO
	err := s.runInTransaction(ctx, func(txCtx context.Context) error {
		lineItems := make([]domain.InvoiceLineItem, 0)
		relatedOrders := make(map[uuid.UUID]struct{})
		dnToInvoiceLinks := make(map[uuid.UUID]uuid.UUID)

		orders, err := s.fetchOrdersForInvoice(txCtx, cmd.PartyID, cmd.SalesOrderIDs)
		if err != nil {
			return err
		}
		for _, order := range orders {
			for _, item := range order.LineItems {
				lineItem, err := buildInvoiceLineItemFromOrder(item, item.Quantity)
				if err != nil {
					return err
				}
				lineItems = append(lineItems, lineItem)
			}
			relatedOrders[order.ID] = struct{}{}
		}

		if len(cmd.DeliveryNoteIDs) > 0 {
			noteItems, noteOrders, noteLinks, err := s.buildInvoiceItemsFromDeliveryNotes(txCtx, cmd.PartyID, cmd.DeliveryNoteIDs)
			if err != nil {
				return err
			}
			lineItems = append(lineItems, noteItems...)
			for _, orderID := range noteOrders {
				relatedOrders[orderID] = struct{}{}
			}
			for k, v := range noteLinks {
				dnToInvoiceLinks[k] = v
			}
		}
		if len(lineItems) == 0 {
			return domain.NewValidationError("invoice must have at least one line item")
		}

		taxAmount, err := zeroMoney()
		if err != nil {
			return err
		}
		paymentTerms := ""
		if cmd.PaymentTerms != nil {
			paymentTerms = *cmd.PaymentTerms
		}

		// Invoices from orders: COMPLETA type with series "FV" (Factura de Venta)
		invoiceType := domain.InvoiceTypeComplete
		currentYear := time.Now().Year()
		series, err := domain.NewInvoiceSeries("FV", currentYear)
		if err != nil {
			return err
		}

		invoiceNumber, err := s.numberGen.NextInvoiceNumber(txCtx, series)
		if err != nil {
			return err
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
			return err
		}

		// Save invoice first, then link DN lines and update order statuses — all within the same transaction
		if err := s.invoiceRepo.Save(txCtx, invoice); err != nil {
			return err
		}

		// Link delivery note line items to their corresponding invoice line items
		if len(dnToInvoiceLinks) > 0 {
			if err := s.deliveryRepo.LinkLineItemsToInvoice(txCtx, dnToInvoiceLinks); err != nil {
				return err
			}
		}

		for orderID := range relatedOrders {
			order, err := s.orderRepo.FindByIDForUpdate(txCtx, orderID)
			if err != nil {
				return err
			}
			if order == nil {
				return domain.NewNotFoundError("order not found")
			}

			if err := s.updateOrderInvoiceStatus(txCtx, order, lineItems); err != nil {
				return err
			}
			if err := s.orderRepo.Save(txCtx, order); err != nil {
				return err
			}
		}

		relatedIDs := make([]uuid.UUID, 0, len(relatedOrders))
		for id := range relatedOrders {
			relatedIDs = append(relatedIDs, id)
		}

		result = NewInvoiceDTO(invoice, relatedIDs, cmd.DeliveryNoteIDs)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SalesService) GetQuote(ctx context.Context, query GetQuoteByIDQuery) (*QuoteDTO, error) {
	quote, err := s.quoteRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, domain.NewNotFoundError("quote not found")
	}
	dto := NewQuoteDTO(quote)
	s.enrichQuoteLineItems(ctx, dto.LineItems)
	s.enrichQuoteWithOrderInfo(ctx, dto)
	return dto, nil
}

func (s *SalesService) enrichQuoteWithOrderInfo(ctx context.Context, dto *QuoteDTO) {
	order, err := s.orderRepo.FindByQuoteID(ctx, dto.ID)
	if err != nil || order == nil {
		return
	}
	dto.GeneratedOrderID = &order.ID
	dto.GeneratedOrderNumber = order.OrderNumber.String()
}

func (s *SalesService) ListQuotes(ctx context.Context, query ListQuotesQuery) ([]*QuoteDTO, error) {
	filter := domain.QuoteFilter{PartyID: query.PartyID, FromDate: query.FromDate, ToDate: query.ToDate, Search: query.Search, Limit: query.PageSize}
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
		dto := NewQuoteDTO(quote)
		s.enrichQuoteLineItems(ctx, dto.LineItems)
		result = append(result, dto)
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
	dto := NewSalesOrderDTO(order)
	s.enrichOrderLineItems(ctx, dto.LineItems)
	return dto, nil
}

func (s *SalesService) ListOrders(ctx context.Context, query ListOrdersQuery) ([]*SalesOrderDTO, error) {
	filter := domain.SalesOrderFilter{PartyID: query.PartyID, FromDate: query.FromDate, ToDate: query.ToDate, Search: query.Search, Limit: query.PageSize}
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
		dto := NewSalesOrderDTO(order)
		s.enrichOrderLineItems(ctx, dto.LineItems)
		result = append(result, dto)
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
	dto := NewDeliveryNoteDTO(note)
	s.enrichDeliveryNoteLineItems(ctx, dto.LineItems)
	// Derive invoiceId from line items (in MVP all lines point to the same invoice)
	for _, li := range note.LineItems {
		if li.InvoiceLineItemID != nil {
			inv, findErr := s.invoiceRepo.FindByDeliveryNoteID(ctx, note.ID)
			if findErr == nil && inv != nil {
				dto.InvoiceID = &inv.ID
			}
			break
		}
	}
	return dto, nil
}

func (s *SalesService) ListDeliveryNotes(ctx context.Context, query ListDeliveryNotesQuery) ([]*DeliveryNoteDTO, error) {
	filter := domain.DeliveryNoteFilter{
		SalesOrderID: query.SalesOrderID,
		PartyID:      query.PartyID,
		FromDate:     query.FromDate,
		ToDate:       query.ToDate,
		Search:       query.Search,
		Limit:        query.PageSize,
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
		dto := NewDeliveryNoteDTO(note)
		s.enrichDeliveryNoteLineItems(ctx, dto.LineItems)
		result = append(result, dto)
	}
	return result, nil
}

func (s *SalesService) ChangeDeliveryNoteStatus(ctx context.Context, cmd ChangeDeliveryNoteStatusCommand) (*DeliveryNoteDTO, error) {
	note, err := s.deliveryRepo.FindByID(ctx, cmd.DeliveryNoteID)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, domain.NewNotFoundError("delivery note not found")
	}

	status, err := parseDeliveryNoteStatus(cmd.NewStatus)
	if err != nil {
		return nil, err
	}
	if err := note.ChangeStatus(status); err != nil {
		return nil, err
	}

	if err := s.deliveryRepo.Save(ctx, note); err != nil {
		return nil, err
	}

	dto := NewDeliveryNoteDTO(note)
	return dto, nil
}

func (s *SalesService) GetInvoice(ctx context.Context, query GetInvoiceByIDQuery) (*InvoiceDTO, error) {
	invoice, err := s.invoiceRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, domain.NewNotFoundError("invoice not found")
	}
	dnIDs, _ := s.invoiceRepo.ListDeliveryNoteIDsByInvoiceID(ctx, invoice.ID)
	dto := NewInvoiceDTO(invoice, nil, dnIDs)
	s.enrichInvoiceLineItems(ctx, dto.LineItems)
	return dto, nil
}

func (s *SalesService) ChangeInvoiceStatus(ctx context.Context, cmd ChangeInvoiceStatusCommand) (*InvoiceDTO, error) {
	invoice, err := s.invoiceRepo.FindByID(ctx, cmd.InvoiceID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, domain.NewNotFoundError("invoice not found")
	}

	status, err := parseInvoiceStatus(cmd.NewStatus)
	if err != nil {
		return nil, err
	}
	if err := invoice.ChangeStatus(status); err != nil {
		return nil, err
	}

	if err := s.invoiceRepo.Save(ctx, invoice); err != nil {
		return nil, err
	}

	return NewInvoiceDTO(invoice, nil, nil), nil
}

func (s *SalesService) ListInvoices(ctx context.Context, query ListInvoicesQuery) ([]*InvoiceDTO, error) {
	var invoices []*domain.Invoice
	var err error

	if query.DeliveryNoteID != nil {
		inv, findErr := s.invoiceRepo.FindByDeliveryNoteID(ctx, *query.DeliveryNoteID)
		if findErr != nil {
			return nil, findErr
		}
		if inv != nil {
			invoices = []*domain.Invoice{inv}
		}
	} else if query.SalesOrderID != nil {
		invoices, err = s.invoiceRepo.ListBySalesOrderID(ctx, *query.SalesOrderID)
	} else {
		filter := domain.InvoiceFilter{PartyID: query.PartyID, FromDate: query.FromDate, ToDate: query.ToDate, Search: query.Search, Limit: query.PageSize}
		if query.Status != nil {
			status, err := parseInvoiceStatus(*query.Status)
			if err != nil {
				return nil, err
			}
			filter.Status = &status
		}
		invoices, err = s.invoiceRepo.List(ctx, filter)
	}
	if err != nil {
		return nil, err
	}
	result := make([]*InvoiceDTO, 0, len(invoices))
	for _, invoice := range invoices {
		dto := NewInvoiceDTO(invoice, nil, nil)
		s.enrichInvoiceLineItems(ctx, dto.LineItems)
		result = append(result, dto)
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
	ID               *uuid.UUID
	ProductVariantID uuid.UUID
	Quantity         int
	UnitPrice        *MoneyDTO
	DiscountPercent  *float64
}

func orderLineItemSeedsFromOrder(items []domain.OrderLineItem) []orderLineItemSeed {
	seeds := make([]orderLineItemSeed, 0, len(items))
	for _, item := range items {
		id := item.ID
		var unitPriceDTO *MoneyDTO
		if item.UnitPrice.Amount() != item.ListUnitPrice.Amount() {
			dto := NewMoneyDTO(item.UnitPrice)
			unitPriceDTO = &dto
		}
		discountPct := &item.DiscountPercent
		seeds = append(seeds, orderLineItemSeed{
			ID:               &id,
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			UnitPrice:        unitPriceDTO,
			DiscountPercent:  discountPct,
		})
	}
	return seeds
}

func (s *SalesService) buildQuoteLineItems(ctx context.Context, partyID uuid.UUID, items []QuoteLineItemInput) ([]domain.QuoteLineItem, error) {
	if s.pricingEngine == nil {
		return nil, domain.NewConfigurationError("pricing engine not configured")
	}
	request := pricing_app.CalculateFinalSalePriceRequest{
		ClientID:  partyID.String(),
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
		return nil, domain.NewValidationError("pricing response mismatch")
	}

	lineItems := make([]domain.QuoteLineItem, 0, len(items))
	for i, item := range items {
		calcItem := pricing.CalculatedItems[i]

		listUnitPrice, err := toDomainMoney(calcItem.BaseSalesPrice)
		if err != nil {
			return nil, err
		}

		// Effective unit price: user override or pricing engine's baseSalesPrice
		effectiveUnitPrice := listUnitPrice
		if item.UnitPrice != nil {
			effectiveUnitPrice, err = domain.NewMoney(item.UnitPrice.Amount, item.UnitPrice.Currency)
			if err != nil {
				return nil, err
			}
		}

		// Discount: user-specified takes precedence, then pricing engine
		discountPercent := calcItem.DiscountPercent
		if item.DiscountPercent != nil {
			discountPercent = *item.DiscountPercent
		}

		// When user overrides price or discount, recalculate line values from the override.
		// Otherwise, use Pricing engine's pre-calculated values (single source of truth).
		if item.UnitPrice != nil || item.DiscountPercent != nil {
			lineItem, err := domain.NewQuoteLineItem(
				item.ProductVariantID,
				item.Quantity,
				listUnitPrice,
				&effectiveUnitPrice,
				discountPercent,
				calcItem.TaxRate,
			)
			if err != nil {
				return nil, err
			}
			lineItems = append(lineItems, lineItem)
		} else {
			// Use pre-calculated values from Pricing — single source of truth
			discountPerUnit, err := toDomainMoney(calcItem.DiscountAmount)
			if err != nil {
				return nil, err
			}
			subtotal, err := toDomainMoney(calcItem.LineSubtotal)
			if err != nil {
				return nil, err
			}
			taxAmount, err := toDomainMoney(calcItem.LineTaxAmount)
			if err != nil {
				return nil, err
			}

			lineItem, err := domain.NewQuoteLineItemFromCalculated(
				item.ProductVariantID,
				item.Quantity,
				listUnitPrice,
				effectiveUnitPrice,
				discountPercent,
				discountPerUnit,
				subtotal,
				calcItem.TaxRate,
				taxAmount,
			)
			if err != nil {
				return nil, err
			}
			lineItems = append(lineItems, lineItem)
		}
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
			ID:               id,
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			UnitPrice:        item.UnitPrice,
			DiscountPercent:  item.DiscountPercent,
		})
	}
	return s.buildOrderLineItemsFromSeeds(ctx, partyID, seeds)
}

func (s *SalesService) buildOrderLineItemsFromSeeds(ctx context.Context, partyID uuid.UUID, seeds []orderLineItemSeed) ([]domain.OrderLineItem, error) {
	if s.pricingEngine == nil {
		return nil, domain.NewConfigurationError("pricing engine not configured")
	}
	request := pricing_app.CalculateFinalSalePriceRequest{
		ClientID:  partyID.String(),
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
		return nil, domain.NewValidationError("pricing response mismatch")
	}

	lineItems := make([]domain.OrderLineItem, 0, len(seeds))
	for i, seed := range seeds {
		calcItem := pricing.CalculatedItems[i]

		listUnitPrice, err := toDomainMoney(calcItem.BaseSalesPrice)
		if err != nil {
			return nil, err
		}

		// Effective unit price: user override or pricing engine's baseSalesPrice
		effectiveUnitPrice := listUnitPrice
		if seed.UnitPrice != nil {
			effectiveUnitPrice, err = domain.NewMoney(seed.UnitPrice.Amount, seed.UnitPrice.Currency)
			if err != nil {
				return nil, err
			}
		}

		// Discount: user-specified takes precedence, then pricing engine
		discountPercent := calcItem.DiscountPercent
		if seed.DiscountPercent != nil {
			discountPercent = *seed.DiscountPercent
		}

		var lineItem domain.OrderLineItem
		// When user overrides price or discount, recalculate line values from the override.
		// Otherwise, use Pricing engine's pre-calculated values (single source of truth).
		if seed.UnitPrice != nil || seed.DiscountPercent != nil {
			lineItem, err = domain.NewOrderLineItem(
				seed.ProductVariantID,
				seed.Quantity,
				listUnitPrice,
				&effectiveUnitPrice,
				discountPercent,
				calcItem.TaxRate,
			)
			if err != nil {
				return nil, err
			}
		} else {
			// Use pre-calculated values from Pricing — single source of truth
			discountPerUnit, err := toDomainMoney(calcItem.DiscountAmount)
			if err != nil {
				return nil, err
			}
			subtotal, err := toDomainMoney(calcItem.LineSubtotal)
			if err != nil {
				return nil, err
			}
			taxAmount, err := toDomainMoney(calcItem.LineTaxAmount)
			if err != nil {
				return nil, err
			}

			lineItem, err = domain.NewOrderLineItemFromCalculated(
				seed.ProductVariantID,
				seed.Quantity,
				listUnitPrice,
				effectiveUnitPrice,
				discountPercent,
				discountPerUnit,
				subtotal,
				calcItem.TaxRate,
				taxAmount,
			)
			if err != nil {
				return nil, err
			}
		}

		if seed.ID != nil {
			lineItem.ID = *seed.ID
		}
		lineItems = append(lineItems, lineItem)
	}

	return lineItems, nil
}

func toDomainMoney(dto pricing_app.MoneyDTO) (domain.Money, error) {
	return domain.NewMoney(dto.Amount, dto.Currency)
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
	switch status {
	case domain.SalesOrderStatusPending, domain.SalesOrderStatusInPreparation:
		return true
	default:
		return false
	}
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
		if note.Status == domain.DeliveryNoteStatusCanceled {
			continue
		}
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
	if item.DiscountPerUnit.Amount() > 0 {
		value := item.DiscountPerUnit
		discount = &value
	}
	lineItem, err := domain.NewInvoiceLineItem(
		item.ProductVariantID,
		quantity,
		item.UnitPrice,
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
		order, err := s.orderRepo.FindByIDForUpdate(ctx, orderID)
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
		if order.Status != domain.SalesOrderStatusDelivered && order.Status != domain.SalesOrderStatusPartiallyDelivered && order.Status != domain.SalesOrderStatusPartiallyInvoiced && order.Status != domain.SalesOrderStatusInvoiced {
			return nil, domain.NewConflictError("order must be delivered before invoicing")
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (s *SalesService) buildInvoiceItemsFromDeliveryNotes(ctx context.Context, partyID uuid.UUID, noteIDs []uuid.UUID) ([]domain.InvoiceLineItem, []uuid.UUID, map[uuid.UUID]uuid.UUID, error) {
	lineItems := make([]domain.InvoiceLineItem, 0)
	orderIDs := make([]uuid.UUID, 0)
	seen := make(map[uuid.UUID]struct{})
	// dnLineItemID → invoiceLineItemID
	dnToInvoiceLinks := make(map[uuid.UUID]uuid.UUID)

	for _, noteID := range noteIDs {
		note, err := s.deliveryRepo.FindByID(ctx, noteID)
		if err != nil {
			return nil, nil, nil, err
		}
		if note == nil {
			return nil, nil, nil, domain.NewNotFoundError("delivery note not found")
		}
		order, err := s.orderRepo.FindByID(ctx, note.SalesOrderID)
		if err != nil {
			return nil, nil, nil, err
		}
		if order == nil {
			return nil, nil, nil, domain.NewNotFoundError("order not found")
		}
		if order.PartyID != partyID {
			return nil, nil, nil, domain.NewValidationError("delivery note party mismatch")
		}
		for _, item := range note.LineItems {
			if item.InvoiceLineItemID != nil {
				return nil, nil, nil, domain.NewValidationError("delivery note line item already invoiced")
			}
			orderLine := findOrderLineItem(order.LineItems, item.SalesOrderLineItemID)
			if orderLine == nil {
				return nil, nil, nil, domain.NewValidationError("sales order line item not found")
			}
			lineItem, err := buildInvoiceLineItemFromOrder(*orderLine, item.DeliveredQuantity)
			if err != nil {
				return nil, nil, nil, err
			}
			dnToInvoiceLinks[item.ID] = lineItem.ID
			lineItems = append(lineItems, lineItem)
		}
		if _, ok := seen[order.ID]; !ok {
			orderIDs = append(orderIDs, order.ID)
			seen[order.ID] = struct{}{}
		}
	}

	return lineItems, orderIDs, dnToInvoiceLinks, nil
}

func (s *SalesService) updateOrderInvoiceStatus(ctx context.Context, order *domain.SalesOrder, newInvoiceItems []domain.InvoiceLineItem) error {
	if order.Status == domain.SalesOrderStatusCanceled {
		return domain.NewConflictError("cannot update invoice status for canceled order")
	}
	if order.Status != domain.SalesOrderStatusDelivered && order.Status != domain.SalesOrderStatusPartiallyDelivered && order.Status != domain.SalesOrderStatusPartiallyInvoiced && order.Status != domain.SalesOrderStatusInvoiced {
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
// Optimized for TPV/POS workflow: validates < 3,000 EUR limit, uses series "FT", allows any CLIENT party
func (s *SalesService) CreateSimplifiedInvoice(ctx context.Context, cmd CreateSimplifiedInvoiceCommand) (*InvoiceDTO, error) {
	if err := s.ensurePartyExists(ctx, cmd.PartyID); err != nil {
		return nil, err
	}

	// Build sale items for pricing calculation
	saleItems := make([]pricing_app.SaleItemRequest, 0, len(cmd.Items))
	discountByVariant := make(map[uuid.UUID]float64)
	for _, itemInput := range cmd.Items {
		saleItems = append(saleItems, pricing_app.SaleItemRequest{
			ProductVariantID: itemInput.ProductVariantID,
			Quantity:         itemInput.Quantity,
		})
		discountByVariant[itemInput.ProductVariantID] = itemInput.DiscountPercent
	}

	// Calculate prices using pricing engine
	priceReq := pricing_app.CalculateFinalSalePriceRequest{
		SaleItems: saleItems,
		ClientID:  cmd.PartyID.String(),
		SaleDate:  cmd.InvoiceDate,
	}
	priceResp, err := s.pricingEngine.CalculateFinalSalePrice(ctx, priceReq)
	if err != nil {
		return nil, domain.NewConfigurationError("pricing calculation failed: " + err.Error())
	}

	// Build invoice line items from calculated prices
	lineItems := make([]domain.InvoiceLineItem, 0, len(priceResp.CalculatedItems))
	for _, calculatedItem := range priceResp.CalculatedItems {
		// Use BaseSalesPrice (catalogue price before client-specific discounts)
		// so the manual discountPercent from the frontend maps to the client discount
		// shown in the UI, avoiding double-discount when the pricing engine also applies it.
		unitPrice, err := domain.NewMoney(calculatedItem.BaseSalesPrice.Amount, domain.DefaultCurrency)
		if err != nil {
			return nil, err
		}

		// Apply manual discount if provided; otherwise use the pricing engine's discount
		var discountAmount *domain.Money
		if dp, ok := discountByVariant[calculatedItem.ProductVariantID]; ok && dp > 0 {
			da := unitPrice.Amount() * dp / 100
			dm, err := domain.NewMoney(da, domain.DefaultCurrency)
			if err != nil {
				return nil, err
			}
			discountAmount = &dm
		} else if calculatedItem.DiscountPercent > 0 {
			da := unitPrice.Amount() * calculatedItem.DiscountPercent / 100
			dm, err := domain.NewMoney(da, domain.DefaultCurrency)
			if err != nil {
				return nil, err
			}
			discountAmount = &dm
		}

		// For tickets, we use direct pricing with optional manual discount
		lineItem, err := domain.NewInvoiceLineItem(
			calculatedItem.ProductVariantID,
			calculatedItem.Quantity,
			unitPrice,
			discountAmount,
			nil, // No tax breakdown per line
			calculatedItem.TaxRate,
		)
		if err != nil {
			return nil, err
		}

		lineItems = append(lineItems, lineItem)
	}

	// Simplified invoice: immediate payment, no payment terms
	// Ticket series: "FT" (Factura de Ticket) for current year
	invoiceType := domain.InvoiceTypeSimplified
	currentYear := time.Now().Year()
	series, err := domain.NewInvoiceSeries("FT", currentYear)
	if err != nil {
		return nil, err
	}

	invoiceNumber, err := s.numberGen.NextInvoiceNumber(ctx, series)
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
	return NewInvoiceDTO(invoice, []uuid.UUID{}, nil), nil
}

// --- Product variant enrichment helpers ---

func (s *SalesService) lookupVariant(ctx context.Context, variantID uuid.UUID) (string, string, map[string]string) {
	if s.productLookup == nil {
		return "", "", nil
	}
	info, err := s.productLookup.GetVariantInfo(ctx, variantID)
	if err != nil || info == nil {
		return "", "", nil
	}
	return info.ProductName, info.VariantSKU, info.OptionConfiguration
}

func (s *SalesService) enrichQuoteLineItems(ctx context.Context, items []QuoteLineItemDTO) {
	for i := range items {
		items[i].ProductName, items[i].VariantSKU, items[i].OptionConfiguration = s.lookupVariant(ctx, items[i].ProductVariantID)
	}
}

func (s *SalesService) enrichOrderLineItems(ctx context.Context, items []OrderLineItemDTO) {
	for i := range items {
		items[i].ProductName, items[i].VariantSKU, items[i].OptionConfiguration = s.lookupVariant(ctx, items[i].ProductVariantID)
	}
}

func (s *SalesService) enrichDeliveryNoteLineItems(ctx context.Context, items []DeliveryNoteLineItemDTO) {
	for i := range items {
		items[i].ProductName, items[i].VariantSKU, items[i].OptionConfiguration = s.lookupVariant(ctx, items[i].ProductVariantID)
	}
}

func (s *SalesService) enrichInvoiceLineItems(ctx context.Context, items []InvoiceLineItemDTO) {
	for i := range items {
		items[i].ProductName, items[i].VariantSKU, items[i].OptionConfiguration = s.lookupVariant(ctx, items[i].ProductVariantID)
	}
}
