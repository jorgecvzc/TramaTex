package application

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/sales/domain"
)

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
		if order.Status == domain.SalesOrderStatusCancelled {
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
	var result *DeliveryNoteDTO
	err := s.runInTransaction(ctx, func(txCtx context.Context) error {
		note, err := s.deliveryRepo.FindByID(txCtx, cmd.DeliveryNoteID)
		if err != nil {
			return err
		}
		if note == nil {
			return domain.NewNotFoundError("delivery note not found")
		}

		status, err := parseDeliveryNoteStatus(cmd.NewStatus)
		if err != nil {
			return err
		}

		if err := note.ChangeStatus(status); err != nil {
			return err
		}

		if err := s.deliveryRepo.Save(txCtx, note); err != nil {
			return err
		}

		// Update order status based on current notes
		order, err := s.orderRepo.FindByIDForUpdate(txCtx, note.SalesOrderID)
		if err != nil {
			return err
		}
		if order != nil {
			delivered, err := s.deliveredQuantities(txCtx, order.ID)
			if err != nil {
				return err
			}

			allDelivered := true
			anyDelivered := false
			for _, li := range order.LineItems {
				qty := delivered[li.ID]
				if qty > 0 {
					anyDelivered = true
				}
				if qty < li.Quantity {
					allDelivered = false
				}
			}

			var targetStatus domain.SalesOrderStatus
			if allDelivered {
				targetStatus = domain.SalesOrderStatusDelivered
			} else if anyDelivered {
				targetStatus = domain.SalesOrderStatusPartiallyDelivered
			} else {
				targetStatus = domain.SalesOrderStatusInPreparation
			}

			if order.Status != targetStatus && order.Status != domain.SalesOrderStatusInvoiced && order.Status != domain.SalesOrderStatusPartiallyInvoiced && order.Status != domain.SalesOrderStatusCancelled {
				if err := order.ChangeStatus(targetStatus); err != nil {
					// Ignore transition errors if current status is incompatible
				}
				if err := s.orderRepo.Save(txCtx, order); err != nil {
					return err
				}
			}
		}

		result = NewDeliveryNoteDTO(note)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SalesService) DeleteDeliveryNote(ctx context.Context, cmd DeleteDeliveryNoteCommand) error {
	return s.runInTransaction(ctx, func(txCtx context.Context) error {
		note, err := s.deliveryRepo.FindByID(txCtx, cmd.DeliveryNoteID)
		if err != nil {
			return err
		}
		if note == nil {
			return domain.NewNotFoundError("delivery note not found")
		}

		// Check if it's already invoiced
		for _, li := range note.LineItems {
			if li.InvoiceLineItemID != nil {
				return domain.NewConflictError("cannot delete an invoiced delivery note")
			}
		}

		// Delete the note
		if err := s.deliveryRepo.Delete(txCtx, note.ID); err != nil {
			return err
		}

		// Recalculate Sales Order Status
		order, err := s.orderRepo.FindByIDForUpdate(txCtx, note.SalesOrderID)
		if err != nil {
			return err
		}
		if order != nil {
			delivered, err := s.deliveredQuantities(txCtx, order.ID)
			if err != nil {
				return err
			}

			allDelivered := true
			anyDelivered := false
			for _, li := range order.LineItems {
				qty := delivered[li.ID]
				if qty > 0 {
					anyDelivered = true
				}
				if qty < li.Quantity {
					allDelivered = false
				}
			}

			var targetStatus domain.SalesOrderStatus
			if allDelivered {
				targetStatus = domain.SalesOrderStatusDelivered
			} else if anyDelivered {
				targetStatus = domain.SalesOrderStatusPartiallyDelivered
			} else {
				targetStatus = domain.SalesOrderStatusInPreparation
			}

			// Do not change status if it is cancelled or invoiced, but allow returning from Delivered/PartiallyDelivered to InPreparation
			if order.Status != targetStatus && order.Status != domain.SalesOrderStatusInvoiced && order.Status != domain.SalesOrderStatusPartiallyInvoiced && order.Status != domain.SalesOrderStatusCancelled {
				if err := order.ChangeStatus(targetStatus); err != nil {
					// Fallback: If domain logic forbids moving back, at least we tried
				}
				if err := s.orderRepo.Save(txCtx, order); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func parseDeliveryNoteStatus(input string) (domain.DeliveryNoteStatus, error) {
	u := strings.ToUpper(strings.TrimSpace(input))

	switch u {
	case "PENDING", "PENDIENTE":
		return domain.DeliveryNoteStatusPending, nil
	case "DELIVERED", "ENTREGADO", "ENTREGADA":
		return domain.DeliveryNoteStatusDelivered, nil
	case "CANCELLED", "CANCELADO", "ANULADO", "ANULADA", "CANCELADA":
		return domain.DeliveryNoteStatusCancelled, nil
	default:
		// Fallback to direct cast and validation
		val := domain.DeliveryNoteStatus(u)
		if err := val.IsValid(); err != nil {
			return "", err
		}
		return val, nil
	}
}

func (s *SalesService) deliveredQuantities(ctx context.Context, orderID uuid.UUID) (map[uuid.UUID]int, error) {
	results := make(map[uuid.UUID]int)
	notes, err := s.deliveryRepo.ListBySalesOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	for _, note := range notes {
		if note.Status == domain.DeliveryNoteStatusCancelled {
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

func (s *SalesService) enrichDeliveryNoteLineItems(ctx context.Context, items []DeliveryNoteLineItemDTO) {
	for i := range items {
		items[i].ProductName, items[i].VariantSKU, items[i].OptionConfiguration = s.lookupVariant(ctx, items[i].ProductVariantID)
	}
}
