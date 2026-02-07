package domain

import "fmt"

type QuoteStatus string

type SalesOrderStatus string

type DeliveryNoteStatus string

type InvoiceStatus string

const (
	QuoteStatusDraft     QuoteStatus = "BORRADOR"
	QuoteStatusSent      QuoteStatus = "ENVIADA"
	QuoteStatusApproved  QuoteStatus = "APROBADA"
	QuoteStatusRejected  QuoteStatus = "RECHAZADA"
	QuoteStatusExpired   QuoteStatus = "EXPIRADA"
	QuoteStatusConverted QuoteStatus = "CONVERTIDA_A_PEDIDO"

	SalesOrderStatusPending            SalesOrderStatus = "PENDIENTE"
	SalesOrderStatusInPreparation      SalesOrderStatus = "EN_PREPARACION"
	SalesOrderStatusPartiallyDelivered SalesOrderStatus = "ENTREGADO_PARCIALMENTE"
	SalesOrderStatusDelivered          SalesOrderStatus = "ENTREGADO"
	SalesOrderStatusCanceled           SalesOrderStatus = "CANCELADO"
	SalesOrderStatusPartiallyInvoiced  SalesOrderStatus = "FACTURADO_PARCIALMENTE"
	SalesOrderStatusInvoiced           SalesOrderStatus = "FACTURADO_COMPLETAMENTE"

	DeliveryNoteStatusPending   DeliveryNoteStatus = "PENDIENTE"
	DeliveryNoteStatusDelivered DeliveryNoteStatus = "ENTREGADO"
	DeliveryNoteStatusCanceled  DeliveryNoteStatus = "CANCELADO"

	InvoiceStatusDraft   InvoiceStatus = "BORRADOR"
	InvoiceStatusIssued  InvoiceStatus = "EMITIDA"
	InvoiceStatusPaid    InvoiceStatus = "PAGADA"
	InvoiceStatusOverdue InvoiceStatus = "VENCIDA"
	InvoiceStatusVoid    InvoiceStatus = "ANULADA"
)

func (s QuoteStatus) IsValid() error {
	switch s {
	case QuoteStatusDraft, QuoteStatusSent, QuoteStatusApproved, QuoteStatusRejected, QuoteStatusExpired, QuoteStatusConverted:
		return nil
	default:
		return NewValidationError(fmt.Sprintf("invalid quote status: %s", s))
	}
}

func (s SalesOrderStatus) IsValid() error {
	switch s {
	case SalesOrderStatusPending, SalesOrderStatusInPreparation, SalesOrderStatusPartiallyDelivered, SalesOrderStatusDelivered, SalesOrderStatusCanceled,
		SalesOrderStatusPartiallyInvoiced, SalesOrderStatusInvoiced:
		return nil
	default:
		return NewValidationError(fmt.Sprintf("invalid sales order status: %s", s))
	}
}

func (s DeliveryNoteStatus) IsValid() error {
	switch s {
	case DeliveryNoteStatusPending, DeliveryNoteStatusDelivered, DeliveryNoteStatusCanceled:
		return nil
	default:
		return NewValidationError(fmt.Sprintf("invalid delivery note status: %s", s))
	}
}

func (s InvoiceStatus) IsValid() error {
	switch s {
	case InvoiceStatusDraft, InvoiceStatusIssued, InvoiceStatusPaid, InvoiceStatusOverdue, InvoiceStatusVoid:
		return nil
	default:
		return NewValidationError(fmt.Sprintf("invalid invoice status: %s", s))
	}
}

func canTransitionQuote(from QuoteStatus, to QuoteStatus) bool {
	switch from {
	case QuoteStatusDraft:
		return to == QuoteStatusSent
	case QuoteStatusSent:
		return to == QuoteStatusApproved || to == QuoteStatusRejected || to == QuoteStatusExpired
	case QuoteStatusApproved:
		return to == QuoteStatusConverted
	default:
		return false
	}
}

func canTransitionOrder(from SalesOrderStatus, to SalesOrderStatus) bool {
	switch from {
	case SalesOrderStatusPending:
		return to == SalesOrderStatusInPreparation || to == SalesOrderStatusCanceled
	case SalesOrderStatusInPreparation:
		return to == SalesOrderStatusPartiallyDelivered || to == SalesOrderStatusDelivered || to == SalesOrderStatusCanceled
	case SalesOrderStatusPartiallyDelivered:
		return to == SalesOrderStatusDelivered || to == SalesOrderStatusCanceled
	case SalesOrderStatusDelivered:
		return to == SalesOrderStatusPartiallyInvoiced || to == SalesOrderStatusInvoiced
	case SalesOrderStatusPartiallyInvoiced:
		return to == SalesOrderStatusInvoiced
	default:
		return false
	}
}

func canTransitionDeliveryNote(from DeliveryNoteStatus, to DeliveryNoteStatus) bool {
	switch from {
	case DeliveryNoteStatusPending:
		return to == DeliveryNoteStatusDelivered || to == DeliveryNoteStatusCanceled
	default:
		return false
	}
}

func canTransitionInvoice(from InvoiceStatus, to InvoiceStatus) bool {
	switch from {
	case InvoiceStatusDraft:
		return to == InvoiceStatusIssued || to == InvoiceStatusVoid
	case InvoiceStatusIssued:
		return to == InvoiceStatusPaid || to == InvoiceStatusOverdue || to == InvoiceStatusVoid
	case InvoiceStatusOverdue:
		return to == InvoiceStatusPaid || to == InvoiceStatusVoid
	default:
		return false
	}
}
