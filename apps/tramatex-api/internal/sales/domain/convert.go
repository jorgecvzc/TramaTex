package domain

import (
	"time"
)

func (q *Quote) ConvertToOrder(orderNumber OrderNumber, deliveryDate time.Time) (*SalesOrder, error) {
	orderItems := make([]OrderLineItem, len(q.LineItems))
	for i, item := range q.LineItems {
		orderItems[i] = OrderLineItem{
			ID:               item.ID,
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			ListUnitPrice:    item.ListUnitPrice,
			UnitPrice:        item.UnitPrice,
			TaxRate:          item.TaxRate,
			DiscountPercent:  item.DiscountPercent,
			DiscountPerUnit:  item.DiscountPerUnit,
			Subtotal:         item.Subtotal,
			TaxAmount:        item.TaxAmount,
		}
	}
	now := time.Now()
	// Truncate to date-only (start of day UTC) so the comparison is date-level,
	// avoiding timezone mismatches between frontend date-picker and server clock.
	orderDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	deliveryDay := time.Date(deliveryDate.Year(), deliveryDate.Month(), deliveryDate.Day(), 0, 0, 0, 0, time.UTC)
	order, err := NewSalesOrder(orderNumber, q.PartyID, orderDay, deliveryDay, orderItems, q.TaxAmount, q.Notes)
	if err != nil {
		return nil, err
	}
	order.QuoteID = &q.ID

	// Copy work references from quote to order with new IDs
	orderWorkRefs := make([]WorkReference, len(q.WorkReferences))
	for i, ws := range q.WorkReferences {
		orderWorkRefs[i] = WorkReference{
			ID:          ws.ID,
			WorkSetupID: ws.WorkSetupID,
			Sequence:    ws.Sequence,
			Description: ws.Description,
		}
	}
	order.WorkReferences = orderWorkRefs

	// Force transition to CONVERTED_TO_ORDER. 
	// The application layer should ensure the quote was APPROVED first.
	// If the state was not valid for conversion, ChangeStatus will return the error.
	if err := q.ChangeStatus(QuoteStatusConverted); err != nil {
		return nil, err
	}
	return order, nil
}
