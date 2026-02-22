package domain

import "time"

func (q *Quote) ConvertToOrder(orderNumber OrderNumber, deliveryDate time.Time) (*SalesOrder, error) {
	if q.Status != QuoteStatusApproved {
		return nil, NewConflictError("quote must be approved before conversion")
	}
	orderItems := make([]OrderLineItem, len(q.LineItems))
	for i, item := range q.LineItems {
		orderItems[i] = OrderLineItem{
			ID:                        item.ID,
			MESWorkID:                 item.MESWorkID,
			ProductVariantID:          item.ProductVariantID,
			Quantity:                  item.Quantity,
			CalculatedUnitPrice:       item.CalculatedUnitPrice,
			ManualUnitPrice:           item.ManualUnitPrice,
			FinalUnitPrice:            item.FinalUnitPrice,
			CalculatedDiscountPerUnit: item.CalculatedDiscountPerUnit,
			ManualDiscountPerUnit:     item.ManualDiscountPerUnit,
			FinalDiscountPerUnit:      item.FinalDiscountPerUnit,
			Subtotal:                  item.Subtotal,
		}
	}
	order, err := NewSalesOrder(orderNumber, q.PartyID, time.Now(), deliveryDate, orderItems, q.TaxAmount, q.Notes)
	if err != nil {
		return nil, err
	}
	order.QuoteID = &q.ID
	if err := q.ChangeStatus(QuoteStatusConverted); err != nil {
		return nil, err
	}
	return order, nil
}
