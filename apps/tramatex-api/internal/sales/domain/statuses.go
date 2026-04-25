package domain

import "fmt"

type QuoteStatus string
type SalesOrderStatus string
type DeliveryNoteStatus string
type InvoiceStatus string

const (
	QuoteStatusDraft     QuoteStatus = "DRAFT"
	QuoteStatusIssued    QuoteStatus = "ISSUED"
	QuoteStatusApproved  QuoteStatus = "APPROVED"
	QuoteStatusRejected  QuoteStatus = "REJECTED"
	QuoteStatusExpired   QuoteStatus = "EXPIRED"
	QuoteStatusConverted QuoteStatus = "CONVERTED_TO_ORDER"

	SalesOrderStatusPending            SalesOrderStatus = "PENDING"
	SalesOrderStatusInPreparation      SalesOrderStatus = "IN_PREPARATION"
	SalesOrderStatusReadyForProduction SalesOrderStatus = "READY_FOR_PRODUCTION"
	SalesOrderStatusPartiallyDelivered SalesOrderStatus = "PARTIALLY_DELIVERED"
	SalesOrderStatusDelivered          SalesOrderStatus = "DELIVERED"
	SalesOrderStatusCancelled          SalesOrderStatus = "CANCELLED"
	SalesOrderStatusPartiallyInvoiced  SalesOrderStatus = "PARTIALLY_INVOICED"
	SalesOrderStatusInvoiced           SalesOrderStatus = "INVOICED"

	DeliveryNoteStatusPending   DeliveryNoteStatus = "PENDING"
	DeliveryNoteStatusDelivered DeliveryNoteStatus = "DELIVERED"
	DeliveryNoteStatusCancelled DeliveryNoteStatus = "CANCELLED"

	InvoiceStatusDraft   InvoiceStatus = "DRAFT"
	InvoiceStatusIssued  InvoiceStatus = "ISSUED"
	InvoiceStatusPaid    InvoiceStatus = "PAID"
	InvoiceStatusOverdue InvoiceStatus = "OVERDUE"
	InvoiceStatusVoid    InvoiceStatus = "VOID"
)

func (s QuoteStatus) IsValid() error {
	switch s {
	case QuoteStatusDraft,
		QuoteStatusIssued,
		QuoteStatusApproved,
		QuoteStatusRejected,
		QuoteStatusExpired,
		QuoteStatusConverted:
		return nil
	default:
		return NewValidationError(fmt.Sprintf("invalid quote status: %s", s))
	}
}

func (s SalesOrderStatus) IsValid() error {
	switch s {
	case SalesOrderStatusPending,
		SalesOrderStatusInPreparation,
		SalesOrderStatusReadyForProduction,
		SalesOrderStatusPartiallyDelivered,
		SalesOrderStatusDelivered,
		SalesOrderStatusCancelled,
		SalesOrderStatusPartiallyInvoiced,
		SalesOrderStatusInvoiced:
		return nil
	default:
		return NewValidationError(fmt.Sprintf("invalid sales order status: %s", s))
	}
}

func (s DeliveryNoteStatus) IsValid() error {
	switch s {
	case DeliveryNoteStatusPending,
		DeliveryNoteStatusDelivered,
		DeliveryNoteStatusCancelled:
		return nil
	default:
		return NewValidationError(fmt.Sprintf("invalid delivery note status: %s", s))
	}
}

func (s InvoiceStatus) IsValid() error {
	switch s {
	case InvoiceStatusDraft,
		InvoiceStatusIssued,
		InvoiceStatusPaid,
		InvoiceStatusOverdue,
		InvoiceStatusVoid:
		return nil
	default:
		return NewValidationError(fmt.Sprintf("invalid invoice status: %s", s))
	}
}

func canTransitionQuote(from QuoteStatus, to QuoteStatus) bool {
	if from == to {
		return true
	}
	switch from {
	case QuoteStatusDraft:
		return to == QuoteStatusIssued
	case QuoteStatusIssued:
		return to == QuoteStatusApproved || to == QuoteStatusRejected || to == QuoteStatusExpired || to == QuoteStatusDraft
	case QuoteStatusApproved:
		return to == QuoteStatusConverted || to == QuoteStatusIssued || to == QuoteStatusDraft
	case QuoteStatusRejected:
		return to == QuoteStatusDraft
	default:
		return false
	}
}

func canTransitionOrder(from SalesOrderStatus, to SalesOrderStatus) bool {
	if from == to {
		return true
	}
	switch from {
	case SalesOrderStatusPending:
		return to == SalesOrderStatusInPreparation || to == SalesOrderStatusReadyForProduction || to == SalesOrderStatusCancelled
	case SalesOrderStatusInPreparation:
		return to == SalesOrderStatusPartiallyDelivered || to == SalesOrderStatusDelivered || to == SalesOrderStatusReadyForProduction || to == SalesOrderStatusCancelled
	case SalesOrderStatusReadyForProduction:
		return to == SalesOrderStatusInPreparation || to == SalesOrderStatusPartiallyDelivered || to == SalesOrderStatusDelivered || to == SalesOrderStatusCancelled
	case SalesOrderStatusPartiallyDelivered:
		return to == SalesOrderStatusInPreparation || to == SalesOrderStatusDelivered || to == SalesOrderStatusPartiallyInvoiced || to == SalesOrderStatusInvoiced || to == SalesOrderStatusCancelled
	case SalesOrderStatusDelivered:
		return to == SalesOrderStatusInPreparation || to == SalesOrderStatusPartiallyDelivered || to == SalesOrderStatusPartiallyInvoiced || to == SalesOrderStatusInvoiced
	case SalesOrderStatusPartiallyInvoiced:
		return to == SalesOrderStatusInvoiced
	case SalesOrderStatusCancelled:
		return to == SalesOrderStatusPending
	default:
		return false
	}
}

func canTransitionDeliveryNote(from DeliveryNoteStatus, to DeliveryNoteStatus) bool {
	if from == to {
		return true
	}
	switch from {
	case DeliveryNoteStatusPending:
		return to == DeliveryNoteStatusDelivered || to == DeliveryNoteStatusCancelled
	default:
		// Delivered and Cancelled are terminal states
		return false
	}
}

func canTransitionInvoice(from InvoiceStatus, to InvoiceStatus) bool {
	if from == to {
		return true
	}
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
