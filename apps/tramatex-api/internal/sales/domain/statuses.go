package domain

import "fmt"

type QuoteStatus string

type SalesOrderStatus string

type DeliveryNoteStatus string

type InvoiceStatus string

const (
	QuoteStatusDraft     QuoteStatus = "BORRADOR"
	QuoteStatusIssued    QuoteStatus = "EMITIDA"
	QuoteStatusApproved  QuoteStatus = "APROBADA"
	QuoteStatusAccepted  QuoteStatus = "ACEPTADA" // Alias for backward compatibility or UI consistency
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
	case QuoteStatusDraft, QuoteStatusIssued, QuoteStatusApproved, QuoteStatusAccepted, QuoteStatusRejected, QuoteStatusExpired, QuoteStatusConverted:
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
		return to == QuoteStatusIssued
	case QuoteStatusIssued:
		return to == QuoteStatusApproved || to == QuoteStatusAccepted || to == QuoteStatusRejected || to == QuoteStatusExpired || to == QuoteStatusDraft
	case QuoteStatusApproved, QuoteStatusAccepted:
		return to == QuoteStatusConverted
	case QuoteStatusRejected:
		return to == QuoteStatusDraft
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
		return to == SalesOrderStatusDelivered || to == SalesOrderStatusPartiallyInvoiced || to == SalesOrderStatusInvoiced || to == SalesOrderStatusCanceled
	case SalesOrderStatusDelivered:
		return to == SalesOrderStatusPartiallyInvoiced || to == SalesOrderStatusInvoiced
	case SalesOrderStatusPartiallyInvoiced:
		return to == SalesOrderStatusInvoiced
	case SalesOrderStatusCanceled:
		return to == SalesOrderStatusPending
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
