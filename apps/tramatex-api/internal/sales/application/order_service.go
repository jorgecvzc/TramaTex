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
		// Full conversion path: quote is APPROVED â†’ copy line items, change quote status
		if quote.Status != domain.QuoteStatusApproved {
			return nil, domain.NewConflictError("quote must be approved before conversion")
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
	workSetups, err := s.processMesWorkRefs(ctx, cmd.PartyID, cmd.MesWorkRefs)
	if err != nil {
		return nil, err
	}
	order.WorkReferences = workSetups

	// Preserve source quote reference when order is created from a non-approved quote
	if cmd.QuoteID != nil {
		order.QuoteID = cmd.QuoteID
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
	if cmd.MesWorkRefs != nil {
		workSetups, err := s.processMesWorkRefs(ctx, order.PartyID, cmd.MesWorkRefs)
		if err != nil {
			return nil, err
		}
		order.WorkReferences = workSetups
	}

	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	// For confirmed orders, create MES WorkOrders for any newly added WorkReferences
	// that don't yet have a linked WorkOrder (e.g. refs added after initial confirmation).
	if order.Status == domain.SalesOrderStatusInPreparation {
		if err := s.createWorkOrdersForOrder(ctx, order); err != nil {
			return nil, err
		}
		// Re-save to persist the WorkOrderIDs returned by MES.
		if err := s.orderRepo.Save(ctx, order); err != nil {
			return nil, err
		}
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

	previousStatus := order.Status

	status, err := parseOrderStatus(cmd.NewStatus)
	if err != nil {
		return nil, err
	}
	if err := order.ChangeStatus(status); err != nil {
		return nil, err
	}

	// On confirmation (PENDIENTE â†’ EN_PREPARACION), create MES WorkOrders
	// for every WorkReference that doesn't yet have a WorkOrderID.
	if previousStatus == domain.SalesOrderStatusPending && status == domain.SalesOrderStatusInPreparation {
		if err := s.createWorkOrdersForOrder(ctx, order); err != nil {
			return nil, err
		}
	}

	// On cancellation of a confirmed order, suspend its MES WorkOrders.
	if status == domain.SalesOrderStatusCanceled {
		if err := s.suspendWorkOrdersForOrder(ctx, order); err != nil {
			return nil, err
		}
	}

	// On reactivation (CANCELADO â†’ PENDIENTE) followed by re-confirmation
	// (PENDIENTE â†’ EN_PREPARACION), reactivate the MES WorkOrders. We also
	// handle CANCELADO â†’ PENDIENTE directly so WorkOrders are ready even
	// before the second status change.
	if previousStatus == domain.SalesOrderStatusCanceled && status == domain.SalesOrderStatusPending {
		if err := s.reactivateWorkOrdersForOrder(ctx, order); err != nil {
			return nil, err
		}
	}

	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	return NewSalesOrderDTO(order), nil
}

// createWorkOrdersForOrder calls MES to create a WorkOrder for each
// WorkReference that has no WorkOrderID yet. Notes from the order are
// copied to the WorkOrder.

func (s *SalesService) createWorkOrdersForOrder(ctx context.Context, order *domain.SalesOrder) error {
	if s.workOrderCreator == nil {
		return nil
	}
	for i := range order.WorkReferences {
		wr := &order.WorkReferences[i]
		if wr.WorkOrderID != nil {
			continue // already has a WorkOrder
		}
		wsID := wr.WorkSetupID // already *uuid.UUID, nil-safe
		workName := wr.Description
		if workName == "" {
			workName = order.OrderNumber.String()
		}
		workOrderID, err := s.workOrderCreator.CreateWorkOrder(
			ctx,
			workName,               // workName (falls back to order number)
			order.PartyID.String(), // partyID
			order.Notes,            // notes (observaciones)
			wsID,                   // workSetupID (optional)
			wr.ID,                  // orderWorkSetupID for linking back
		)
		if err != nil {
			return fmt.Errorf("create work order for work reference %s: %w", wr.ID, err)
		}
		wr.WorkOrderID = &workOrderID
	}
	return nil
}

// collectWorkOrderIDs returns the non-nil WorkOrderIDs from the order's WorkReferences.

func collectWorkOrderIDs(order *domain.SalesOrder) []uuid.UUID {
	var ids []uuid.UUID
	for _, wr := range order.WorkReferences {
		if wr.WorkOrderID != nil {
			ids = append(ids, *wr.WorkOrderID)
		}
	}
	return ids
}

// suspendWorkOrdersForOrder tells MES to suspend all linked WorkOrders.

func (s *SalesService) suspendWorkOrdersForOrder(ctx context.Context, order *domain.SalesOrder) error {
	if s.workOrderSuspender == nil {
		return nil
	}
	ids := collectWorkOrderIDs(order)
	if len(ids) == 0 {
		return nil
	}
	return s.workOrderSuspender.SuspendWorkOrders(ctx, ids)
}

// reactivateWorkOrdersForOrder tells MES to reactivate all linked WorkOrders.

func (s *SalesService) reactivateWorkOrdersForOrder(ctx context.Context, order *domain.SalesOrder) error {
	if s.workOrderSuspender == nil {
		return nil
	}
	ids := collectWorkOrderIDs(order)
	if len(ids) == 0 {
		return nil
	}
	return s.workOrderSuspender.ReactivateWorkOrders(ctx, ids)
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

	taxAmount, err := sumOrderLineItemTaxes(lineItems)
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

	taxAmount, err := sumOrderLineItemTaxes(lineItems)
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
	s.enrichOrderWithQuoteInfo(ctx, dto)
	s.enrichOrderMESStatus(ctx, dto)
	return dto, nil
}


func (s *SalesService) enrichOrderWithQuoteInfo(ctx context.Context, dto *SalesOrderDTO) {
	if dto.QuoteID == nil {
		return
	}
	quote, err := s.quoteRepo.FindByID(ctx, *dto.QuoteID)
	if err != nil || quote == nil {
		return
	}
	dto.SourceQuoteNumber = quote.QuoteNumber.String()
}


func (s *SalesService) enrichOrderMESStatus(ctx context.Context, dto *SalesOrderDTO) {
	if s.workOrderStatusProvider == nil || len(dto.MesWorkRefs) == 0 {
		return
	}
	var ids []uuid.UUID
	for _, ref := range dto.MesWorkRefs {
		if ref.WorkOrderID != nil {
			ids = append(ids, *ref.WorkOrderID)
		}
	}
	if len(ids) == 0 {
		return
	}
	statuses, err := s.workOrderStatusProvider.GetWorkOrderStatuses(ctx, ids)
	if err != nil {
		return
	}
	allComplete := true
	for i, ref := range dto.MesWorkRefs {
		if ref.WorkOrderID == nil {
			allComplete = false
			continue
		}
		if st, ok := statuses[*ref.WorkOrderID]; ok {
			s := st
			dto.MesWorkRefs[i].WorkOrderStatus = &s
			if st != "COMPLETED" {
				allComplete = false
			}
		} else {
			allComplete = false
		}
	}
	if len(ids) == len(dto.MesWorkRefs) {
		dto.ProductionReady = allComplete
	}
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


func (s *SalesService) ListPendingWorkSetups(ctx context.Context) ([]PendingWorkSetupDTO, error) {
	status := domain.SalesOrderStatusInPreparation
	filter := domain.SalesOrderFilter{Status: &status}
	orders, err := s.orderRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	var result []PendingWorkSetupDTO
	for _, order := range orders {
		for _, ws := range order.WorkReferences {
			if ws.WorkOrderID != nil {
				continue
			}
			result = append(result, PendingWorkSetupDTO{
				ID:           ws.ID,
				WorkSetupID:  ws.WorkSetupID,
				Description:  ws.Description,
				OrderID:      order.ID,
				OrderNumber:  order.OrderNumber.String(),
				DeliveryDate: order.DeliveryDate,
				PartyID:      order.PartyID,
			})
		}
	}
	return result, nil
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
			// Use pre-calculated values from Pricing â€” single source of truth
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


func parseOrderStatus(input string) (domain.SalesOrderStatus, error) {
	value := domain.SalesOrderStatus(strings.ToUpper(strings.TrimSpace(input)))
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


func (s *SalesService) enrichOrderLineItems(ctx context.Context, items []OrderLineItemDTO) {
	for i := range items {
		items[i].ProductName, items[i].VariantSKU, items[i].OptionConfiguration = s.lookupVariant(ctx, items[i].ProductVariantID)
	}
}


