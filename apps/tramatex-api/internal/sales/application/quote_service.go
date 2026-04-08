package application

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	pricing_app "github.com/joran-cortez/tramatex/internal/pricing/application"
	"github.com/joran-cortez/tramatex/internal/sales/domain"
)

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

	now := time.Now()
	quoteDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	quote, err := domain.NewQuote(
		quoteNumber,
		cmd.PartyID,
		quoteDay,
		cmd.ExpirationDate,
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
	quote.WorkReferences = workSetups

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
		workSetups, err := s.processMesWorkRefs(ctx, quote.PartyID, cmd.MesWorkRefs)
		if err != nil {
			return nil, err
		}
		quote.WorkReferences = workSetups
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
		quote.TaxAmount, err = sumQuoteLineItemTaxes(lineItems)
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

	// Step 1: Accept the quote (EMITIDA â†’ APROBADA)
	if err := quote.ChangeStatus(domain.QuoteStatusApproved); err != nil {
		return nil, err
	}

	// Step 2: Convert to order (APROBADA â†’ CONVERTIDA_A_PEDIDO + new order)
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

		// Tariff/list price: base cost plus variant attribute modifiers, without brand markup.
		listUnitPrice, err := toDomainMoney(calcItem.BaseCost)
		if err != nil {
			return nil, err
		}

		// Sale price: pricing engine base sales price (includes brand margin/rules), unless user overrides it.
		effectiveUnitPrice, err := toDomainMoney(calcItem.BaseSalesPrice)
		if err != nil {
			return nil, err
		}
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

func parseQuoteStatus(input string) (domain.QuoteStatus, error) {
	value := domain.QuoteStatus(strings.ToUpper(strings.TrimSpace(input)))
	return value, value.IsValid()
}

func (s *SalesService) enrichQuoteLineItems(ctx context.Context, items []QuoteLineItemDTO) {
	for i := range items {
		items[i].ProductName, items[i].VariantSKU, items[i].OptionConfiguration = s.lookupVariant(ctx, items[i].ProductVariantID)
	}
}
