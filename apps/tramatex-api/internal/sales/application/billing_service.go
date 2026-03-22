package application

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	pricing_app "github.com/joran-cortez/tramatex/internal/pricing/application"
	"github.com/joran-cortez/tramatex/internal/sales/domain"
)

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
			// Calculate already-invoiced quantities to avoid double-invoicing
			// when some delivery notes for this order have already been invoiced.
			existingInvoices, err := s.invoiceRepo.ListBySalesOrderID(txCtx, order.ID)
			if err != nil {
				return err
			}
			alreadyInvoiced := make(map[uuid.UUID]int)
			for _, inv := range existingInvoices {
				for _, invItem := range inv.LineItems {
					if invItem.SalesOrderLineItemID != nil {
						alreadyInvoiced[*invItem.SalesOrderLineItemID] += invItem.Quantity
					}
				}
			}

			for _, item := range order.LineItems {
				remaining := item.Quantity - alreadyInvoiced[item.ID]
				if remaining <= 0 {
					continue // This line item is already fully invoiced
				}
				lineItem, err := buildInvoiceLineItemFromOrder(item, remaining)
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

		// Save invoice first, then link DN lines and update order statuses â€” all within the same transaction
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

func (s *SalesService) GetInvoice(ctx context.Context, query GetInvoiceByIDQuery) (*InvoiceDTO, error) {
	invoice, err := s.invoiceRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, domain.NewNotFoundError("invoice not found")
	}
	dnIDs, _ := s.invoiceRepo.ListDeliveryNoteIDsByInvoiceID(ctx, invoice.ID)
	orderIDs, _ := s.invoiceRepo.ListOrderIDsByInvoiceID(ctx, invoice.ID)
	dto := NewInvoiceDTO(invoice, orderIDs, dnIDs)
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
		if query.Type != nil {
			invoiceType := domain.InvoiceType(*query.Type)
			if err := invoiceType.IsValid(); err != nil {
				return nil, err
			}
			filter.Type = &invoiceType
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

func parseInvoiceStatus(input string) (domain.InvoiceStatus, error) {
	value := domain.InvoiceStatus(strings.ToUpper(strings.TrimSpace(input)))
	return value, value.IsValid()
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
	// dnLineItemID â†’ invoiceLineItemID
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

func (s *SalesService) enrichInvoiceLineItems(ctx context.Context, items []InvoiceLineItemDTO) {
	for i := range items {
		items[i].ProductName, items[i].VariantSKU, items[i].OptionConfiguration = s.lookupVariant(ctx, items[i].ProductVariantID)
	}
}
