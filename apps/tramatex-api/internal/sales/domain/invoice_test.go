package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewInvoice_Success(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-001")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(210, DefaultCurrency)
	unitPrice, _ := NewMoney(100, DefaultCurrency)

	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil)

	invoice, err := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"30 days",
	)

	if err != nil {
		t.Fatalf("NewInvoice() error = %v, want nil", err)
	}
	if invoice.ID == uuid.Nil {
		t.Error("NewInvoice() ID not set")
	}
	if invoice.Type != InvoiceTypeComplete {
		t.Errorf("NewInvoice() Type = %v, want %v", invoice.Type, InvoiceTypeComplete)
	}
	if invoice.Series.Code() != "A" {
		t.Errorf("NewInvoice() Series.Code = %v, want A", invoice.Series.Code())
	}
	if invoice.Status != InvoiceStatusDraft {
		t.Errorf("NewInvoice() Status = %v, want %v", invoice.Status, InvoiceStatusDraft)
	}
}

func TestNewInvoice_InvalidType(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-001")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(210, DefaultCurrency)

	_, err := NewInvoice(
		number,
		InvoiceType("INVALID"),
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{},
		taxAmount,
		"30 days",
	)

	if err == nil {
		t.Error("NewInvoice() with invalid type should return error")
	}
}

func TestNewInvoice_EmptyPartyID(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-001")
	series, _ := NewInvoiceSeries("A", 2026)
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(210, DefaultCurrency)

	_, err := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		uuid.Nil,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{},
		taxAmount,
		"30 days",
	)

	if err == nil {
		t.Error("NewInvoice() with empty party ID should return error")
	}
}

func TestInvoice_ValidateLegalLimits_CompleteInvoice(t *testing.T) {
	// Complete invoices have no amount limit
	number, _ := NewInvoiceNumber("INV-001")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(1000, DefaultCurrency)
	unitPrice, _ := NewMoney(5000, DefaultCurrency) // 5,000 EUR per unit

	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil) // 50,000 EUR total

	invoice, err := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"30 days",
	)

	if err != nil {
		t.Fatalf("NewInvoice() error = %v, want nil (complete invoices have no limit)", err)
	}
	if invoice.Total.Amount() < 3000 {
		t.Error("test setup error: total should be > 3,000 EUR to validate no limit enforcement")
	}
}

func TestInvoice_ValidateLegalLimits_SimplifiedInvoice_UnderLimit(t *testing.T) {
	// Simplified invoices (tickets) must be < 3,000 EUR
	number, _ := NewInvoiceNumber("TKT-001")
	series, _ := NewInvoiceSeries("TKT", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate
	taxAmount, _ := NewMoney(200, DefaultCurrency)
	unitPrice, _ := NewMoney(100, DefaultCurrency)

	lineItem, _ := NewInvoiceLineItem(uuid.New(), 20, unitPrice, nil, nil) // 2,000 EUR

	invoice, err := NewInvoice(
		number,
		InvoiceTypeSimplified,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"Immediate",
	)

	if err != nil {
		t.Fatalf("NewInvoice() error = %v, want nil (under 3,000 EUR limit)", err)
	}
	// Total = 2,000 + 200 = 2,200 EUR (under limit)
	if invoice.Total.Amount() >= 3000 {
		t.Errorf("test setup error: total = %.2f, want < 3,000", invoice.Total.Amount())
	}
}

func TestInvoice_ValidateLegalLimits_SimplifiedInvoice_OverLimit(t *testing.T) {
	// Simplified invoices (tickets) >= 3,000 EUR should fail
	number, _ := NewInvoiceNumber("TKT-002")
	series, _ := NewInvoiceSeries("TKT", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate
	taxAmount, _ := NewMoney(500, DefaultCurrency)
	unitPrice, _ := NewMoney(300, DefaultCurrency)

	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil) // 3,000 EUR

	_, err := NewInvoice(
		number,
		InvoiceTypeSimplified,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"Immediate",
	)

	if err == nil {
		t.Error("NewInvoice() should fail for simplified invoice >= 3,000 EUR")
	}
	// Total = 3,000 + 500 = 3,500 EUR (over limit)
}

func TestInvoice_ValidateLegalLimits_SimplifiedInvoice_ExactLimit(t *testing.T) {
	// Exactly 3,000 EUR should fail (limit is strictly less than 3,000)
	number, _ := NewInvoiceNumber("TKT-003")
	series, _ := NewInvoiceSeries("TKT", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate
	taxAmount, _ := NewMoney(0, DefaultCurrency)
	unitPrice, _ := NewMoney(3000, DefaultCurrency)

	lineItem, _ := NewInvoiceLineItem(uuid.New(), 1, unitPrice, nil, nil) // 3,000 EUR

	_, err := NewInvoice(
		number,
		InvoiceTypeSimplified,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"Immediate",
	)

	if err == nil {
		t.Error("NewInvoice() should fail for simplified invoice at exactly 3,000 EUR (limit is < 3,000)")
	}
}

func TestInvoice_RecalculateTotals(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-001")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(0, DefaultCurrency)
	unitPrice, _ := NewMoney(50, DefaultCurrency)

	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil, 20.0) // 500 EUR + 20% = 100 tax

	invoice, _ := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"30 days",
	)

	originalTotal := invoice.Total.Amount()

	// Modify line items
	newLineItem, _ := NewInvoiceLineItem(uuid.New(), 5, unitPrice, nil, nil, 20.0) // 250 EUR + 20% = 50 tax
	invoice.LineItems = append(invoice.LineItems, newLineItem)

	err := invoice.RecalculateTotals()
	if err != nil {
		t.Fatalf("RecalculateTotals() error = %v", err)
	}

	// New total should be (500 + 250) subtotal + (100 + 50) tax = 900 EUR
	expectedTotal := 900.0
	if invoice.Total.Amount() != expectedTotal {
		t.Errorf("RecalculateTotals() Total = %.2f, want %.2f", invoice.Total.Amount(), expectedTotal)
	}
	if invoice.Total.Amount() <= originalTotal {
		t.Error("RecalculateTotals() should increase total")
	}
}

func TestInvoice_RecalculateTotals_SimplifiedOverLimit(t *testing.T) {
	number, _ := NewInvoiceNumber("TKT-001")
	series, _ := NewInvoiceSeries("TKT", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate
	taxAmount, _ := NewMoney(0, DefaultCurrency)
	unitPrice, _ := NewMoney(100, DefaultCurrency)

	lineItem, _ := NewInvoiceLineItem(uuid.New(), 20, unitPrice, nil, nil, 10.0) // 2,000 EUR + 10% = 200 tax

	invoice, _ := NewInvoice(
		number,
		InvoiceTypeSimplified,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"Immediate",
	)

	// Add more items to exceed limit
	newLineItem, _ := NewInvoiceLineItem(uuid.New(), 15, unitPrice, nil, nil, 10.0) // +1,500 EUR + 10% = 150 tax
	invoice.LineItems = append(invoice.LineItems, newLineItem)

	err := invoice.RecalculateTotals()
	if err == nil {
		t.Error("RecalculateTotals() should fail when simplified invoice exceeds 3,000 EUR after recalculation")
	}
	// Total would be (2,000 + 1,500) subtotal + (200 + 150) tax = 3,850 EUR (over limit)
}

func TestNewInvoiceLineItem_Success(t *testing.T) {
	productVariantID := uuid.New()
	quantity := 10
	unitPrice, _ := NewMoney(100, DefaultCurrency)

	lineItem, err := NewInvoiceLineItem(productVariantID, quantity, unitPrice, nil, nil)

	if err != nil {
		t.Fatalf("NewInvoiceLineItem() error = %v, want nil", err)
	}
	if lineItem.ID == uuid.Nil {
		t.Error("NewInvoiceLineItem() ID not set")
	}
	if lineItem.Quantity != quantity {
		t.Errorf("NewInvoiceLineItem() Quantity = %v, want %v", lineItem.Quantity, quantity)
	}
	expectedSubtotal := 1000.0 // 10 * 100
	if lineItem.Subtotal.Amount() != expectedSubtotal {
		t.Errorf("NewInvoiceLineItem() Subtotal = %.2f, want %.2f", lineItem.Subtotal.Amount(), expectedSubtotal)
	}
}

func TestNewInvoiceLineItem_WithDiscount(t *testing.T) {
	productVariantID := uuid.New()
	quantity := 10
	unitPrice, _ := NewMoney(100, DefaultCurrency)
	discount, _ := NewMoney(10, DefaultCurrency)

	lineItem, err := NewInvoiceLineItem(productVariantID, quantity, unitPrice, &discount, nil)

	if err != nil {
		t.Fatalf("NewInvoiceLineItem() error = %v, want nil", err)
	}
	expectedSubtotal := 900.0 // 10 * (100 - 10)
	if lineItem.Subtotal.Amount() != expectedSubtotal {
		t.Errorf("NewInvoiceLineItem() Subtotal with discount = %.2f, want %.2f", lineItem.Subtotal.Amount(), expectedSubtotal)
	}
}

func TestNewInvoiceLineItem_InvalidQuantity(t *testing.T) {
	productVariantID := uuid.New()
	unitPrice, _ := NewMoney(100, DefaultCurrency)

	_, err := NewInvoiceLineItem(productVariantID, 0, unitPrice, nil, nil)

	if err == nil {
		t.Error("NewInvoiceLineItem() with zero quantity should return error")
	}
}

func TestNewInvoiceLineItem_EmptyProductVariantID(t *testing.T) {
	unitPrice, _ := NewMoney(100, DefaultCurrency)

	_, err := NewInvoiceLineItem(uuid.Nil, 10, unitPrice, nil, nil)

	if err == nil {
		t.Error("NewInvoiceLineItem() with empty product variant ID should return error")
	}
}

// ===== ChangeStatus Tests =====

func TestInvoice_ChangeStatus_DraftToIssued_Success(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-001")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(100, DefaultCurrency)
	unitPrice, _ := NewMoney(50, DefaultCurrency)
	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil)

	invoice, _ := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"30 days",
	)

	err := invoice.ChangeStatus(InvoiceStatusIssued)

	if err != nil {
		t.Errorf("ChangeStatus(ISSUED) error = %v, want nil", err)
	}
	if invoice.Status != InvoiceStatusIssued {
		t.Errorf("ChangeStatus(ISSUED) Status = %v, want %v", invoice.Status, InvoiceStatusIssued)
	}
}

func TestInvoice_ChangeStatus_DraftToVoid_Success(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-002")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(100, DefaultCurrency)
	unitPrice, _ := NewMoney(50, DefaultCurrency)
	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil)

	invoice, _ := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"30 days",
	)

	err := invoice.ChangeStatus(InvoiceStatusVoid)

	if err != nil {
		t.Errorf("ChangeStatus(VOID) error = %v, want nil", err)
	}
	if invoice.Status != InvoiceStatusVoid {
		t.Errorf("ChangeStatus(VOID) Status = %v, want %v", invoice.Status, InvoiceStatusVoid)
	}
}

func TestInvoice_ChangeStatus_IssuedToPaid_Success(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-003")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(100, DefaultCurrency)
	unitPrice, _ := NewMoney(50, DefaultCurrency)
	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil)

	invoice, _ := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"30 days",
	)

	_ = invoice.ChangeStatus(InvoiceStatusIssued)
	err := invoice.ChangeStatus(InvoiceStatusPaid)

	if err != nil {
		t.Errorf("ChangeStatus(PAID) from ISSUED error = %v, want nil", err)
	}
	if invoice.Status != InvoiceStatusPaid {
		t.Errorf("ChangeStatus(PAID) Status = %v, want %v", invoice.Status, InvoiceStatusPaid)
	}
}

func TestInvoice_ChangeStatus_IssuedToOverdue_Success(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-004")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(100, DefaultCurrency)
	unitPrice, _ := NewMoney(50, DefaultCurrency)
	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil)

	invoice, _ := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"30 days",
	)

	_ = invoice.ChangeStatus(InvoiceStatusIssued)
	err := invoice.ChangeStatus(InvoiceStatusOverdue)

	if err != nil {
		t.Errorf("ChangeStatus(OVERDUE) from ISSUED error = %v, want nil", err)
	}
	if invoice.Status != InvoiceStatusOverdue {
		t.Errorf("ChangeStatus(OVERDUE) Status = %v, want %v", invoice.Status, InvoiceStatusOverdue)
	}
}

func TestInvoice_ChangeStatus_IssuedToVoid_Success(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-005")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(100, DefaultCurrency)
	unitPrice, _ := NewMoney(50, DefaultCurrency)
	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil)

	invoice, _ := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"30 days",
	)

	_ = invoice.ChangeStatus(InvoiceStatusIssued)
	err := invoice.ChangeStatus(InvoiceStatusVoid)

	if err != nil {
		t.Errorf("ChangeStatus(VOID) from ISSUED error = %v, want nil", err)
	}
	if invoice.Status != InvoiceStatusVoid {
		t.Errorf("ChangeStatus(VOID) Status = %v, want %v", invoice.Status, InvoiceStatusVoid)
	}
}

func TestInvoice_ChangeStatus_OverdueToPaid_Success(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-006")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(100, DefaultCurrency)
	unitPrice, _ := NewMoney(50, DefaultCurrency)
	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil)

	invoice, _ := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"30 days",
	)

	_ = invoice.ChangeStatus(InvoiceStatusIssued)
	_ = invoice.ChangeStatus(InvoiceStatusOverdue)
	err := invoice.ChangeStatus(InvoiceStatusPaid)

	if err != nil {
		t.Errorf("ChangeStatus(PAID) from OVERDUE error = %v, want nil", err)
	}
	if invoice.Status != InvoiceStatusPaid {
		t.Errorf("ChangeStatus(PAID) Status = %v, want %v", invoice.Status, InvoiceStatusPaid)
	}
}

func TestInvoice_ChangeStatus_OverdueToVoid_Success(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-007")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(100, DefaultCurrency)
	unitPrice, _ := NewMoney(50, DefaultCurrency)
	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil)

	invoice, _ := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"30 days",
	)

	_ = invoice.ChangeStatus(InvoiceStatusIssued)
	_ = invoice.ChangeStatus(InvoiceStatusOverdue)
	err := invoice.ChangeStatus(InvoiceStatusVoid)

	if err != nil {
		t.Errorf("ChangeStatus(VOID) from OVERDUE error = %v, want nil", err)
	}
	if invoice.Status != InvoiceStatusVoid {
		t.Errorf("ChangeStatus(VOID) Status = %v, want %v", invoice.Status, InvoiceStatusVoid)
	}
}

func TestInvoice_ChangeStatus_InvalidTransition_DraftToPaid_Fail(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-008")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(100, DefaultCurrency)
	unitPrice, _ := NewMoney(50, DefaultCurrency)
	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil)

	invoice, _ := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"30 days",
	)

	err := invoice.ChangeStatus(InvoiceStatusPaid)

	if err == nil {
		t.Error("ChangeStatus(PAID) from DRAFT should fail")
	}
	if invoice.Status != InvoiceStatusDraft {
		t.Errorf("Status should remain DRAFT, got %v", invoice.Status)
	}
}

func TestInvoice_ChangeStatus_InvalidTransition_PaidToAnything_Fail(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-009")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(100, DefaultCurrency)
	unitPrice, _ := NewMoney(50, DefaultCurrency)
	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil)

	invoice, _ := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"30 days",
	)

	_ = invoice.ChangeStatus(InvoiceStatusIssued)
	_ = invoice.ChangeStatus(InvoiceStatusPaid)
	err := invoice.ChangeStatus(InvoiceStatusVoid)

	if err == nil {
		t.Error("ChangeStatus(VOID) from PAID should fail (terminal state)")
	}
	if invoice.Status != InvoiceStatusPaid {
		t.Errorf("Status should remain PAID, got %v", invoice.Status)
	}
}

func TestInvoice_ChangeStatus_InvalidTransition_VoidToAnything_Fail(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-010")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(100, DefaultCurrency)
	unitPrice, _ := NewMoney(50, DefaultCurrency)
	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil)

	invoice, _ := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"30 days",
	)

	_ = invoice.ChangeStatus(InvoiceStatusVoid)
	err := invoice.ChangeStatus(InvoiceStatusIssued)

	if err == nil {
		t.Error("ChangeStatus(ISSUED) from VOID should fail (terminal state)")
	}
	if invoice.Status != InvoiceStatusVoid {
		t.Errorf("Status should remain VOID, got %v", invoice.Status)
	}
}

func TestInvoice_ChangeStatus_InvalidStatus_Fail(t *testing.T) {
	number, _ := NewInvoiceNumber("INV-011")
	series, _ := NewInvoiceSeries("A", 2026)
	partyID := uuid.New()
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, 30)
	taxAmount, _ := NewMoney(100, DefaultCurrency)
	unitPrice, _ := NewMoney(50, DefaultCurrency)
	lineItem, _ := NewInvoiceLineItem(uuid.New(), 10, unitPrice, nil, nil)

	invoice, _ := NewInvoice(
		number,
		InvoiceTypeComplete,
		series,
		partyID,
		invoiceDate,
		dueDate,
		[]InvoiceLineItem{lineItem},
		taxAmount,
		"30 days",
	)

	err := invoice.ChangeStatus(InvoiceStatus("INVALID_STATUS"))

	if err == nil {
		t.Error("ChangeStatus() with invalid status should fail")
	}
	if invoice.Status != InvoiceStatusDraft {
		t.Errorf("Status should remain DRAFT, got %v", invoice.Status)
	}
}
