package domain

func (i *Invoice) MarkPaid() error {
	return i.ChangeStatus(InvoiceStatusPaid)
}
