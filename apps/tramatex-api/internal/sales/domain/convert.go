package domain

import "time"

func (q *Quote) ConvertToOrder(orderNumber OrderNumber, deliveryDate time.Time) (*SalesOrder, error) {
	if q.Status != QuoteStatusApproved {
		return nil, NewConflictError("quote must be approved before conversion")
	}
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
	order.MESWorkRefs = q.MESWorkRefs // Carry document-level MES references
	if err := q.ChangeStatus(QuoteStatusConverted); err != nil {
		return nil, err
	}
	return order, nil
}
