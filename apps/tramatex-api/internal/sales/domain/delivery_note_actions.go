package domain

func (d *DeliveryNote) RegisterDelivery() error {
	return d.ChangeStatus(DeliveryNoteStatusDelivered)
}
