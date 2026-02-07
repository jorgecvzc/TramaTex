package domain

import "strings"

type QuoteNumber struct {
	value string
}

type OrderNumber struct {
	value string
}

type DeliveryNoteNumber struct {
	value string
}

type InvoiceNumber struct {
	value string
}

func NewQuoteNumber(value string) (QuoteNumber, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return QuoteNumber{}, NewValidationError("quote number cannot be empty")
	}
	return QuoteNumber{value: value}, nil
}

func NewOrderNumber(value string) (OrderNumber, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return OrderNumber{}, NewValidationError("order number cannot be empty")
	}
	return OrderNumber{value: value}, nil
}

func NewDeliveryNoteNumber(value string) (DeliveryNoteNumber, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return DeliveryNoteNumber{}, NewValidationError("delivery note number cannot be empty")
	}
	return DeliveryNoteNumber{value: value}, nil
}

func NewInvoiceNumber(value string) (InvoiceNumber, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return InvoiceNumber{}, NewValidationError("invoice number cannot be empty")
	}
	return InvoiceNumber{value: value}, nil
}

func (n QuoteNumber) String() string {
	return n.value
}

func (n OrderNumber) String() string {
	return n.value
}

func (n DeliveryNoteNumber) String() string {
	return n.value
}

func (n InvoiceNumber) String() string {
	return n.value
}
