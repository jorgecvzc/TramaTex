package persistence

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/sales/domain"
)

type QuoteDataModel struct {
	gorm.Model
	ID               uuid.UUID `gorm:"type:uuid;primary_key;"`
	QuoteNumber      string    `gorm:"column:quote_number;not null"`
	PartyID          uuid.UUID `gorm:"column:party_id;not null"`
	QuoteDate        time.Time `gorm:"column:quote_date;not null"`
	ExpirationDate   time.Time `gorm:"column:expiration_date;not null"`
	Status           string    `gorm:"type:quote_status;not null"`
	SubtotalAmount   float64   `gorm:"column:subtotal_amount;type:numeric(12,2);not null"`
	SubtotalCurrency string    `gorm:"column:subtotal_currency;type:varchar(3);not null"`
	TaxAmount        float64   `gorm:"column:tax_amount;type:numeric(12,2);not null"`
	TaxCurrency      string    `gorm:"column:tax_currency;type:varchar(3);not null"`
	TotalAmount      float64   `gorm:"column:total_amount;type:numeric(12,2);not null"`
	TotalCurrency    string    `gorm:"column:total_currency;type:varchar(3);not null"`
	Notes            *string   `gorm:"column:notes"`
}

func (QuoteDataModel) TableName() string {
	return "quotes"
}

type QuoteLineItemDataModel struct {
	gorm.Model
	ID                          uuid.UUID  `gorm:"type:uuid;primary_key;"`
	QuoteID                     uuid.UUID  `gorm:"column:quote_id;not null"`
	MESWorkID                   *uuid.UUID `gorm:"column:mes_work_id;type:uuid"`
	ProductVariantID            uuid.UUID  `gorm:"column:product_variant_id;not null"`
	Quantity                    int        `gorm:"column:quantity;not null"`
	CalculatedUnitPriceAmount   float64    `gorm:"column:calculated_unit_price_amount;type:numeric(12,2);not null"`
	CalculatedUnitPriceCurrency string     `gorm:"column:calculated_unit_price_currency;type:varchar(3);not null"`
	ManualUnitPriceAmount       *float64   `gorm:"column:manual_unit_price_amount;type:numeric(12,2)"`
	ManualUnitPriceCurrency     *string    `gorm:"column:manual_unit_price_currency;type:varchar(3)"`
	FinalUnitPriceAmount        float64    `gorm:"column:final_unit_price_amount;type:numeric(12,2);not null"`
	FinalUnitPriceCurrency      string     `gorm:"column:final_unit_price_currency;type:varchar(3);not null"`
	CalculatedDiscountAmount    *float64   `gorm:"column:calculated_discount_per_unit_amount;type:numeric(12,2)"`
	CalculatedDiscountCurrency  *string    `gorm:"column:calculated_discount_per_unit_currency;type:varchar(3)"`
	ManualDiscountAmount        *float64   `gorm:"column:manual_discount_per_unit_amount;type:numeric(12,2)"`
	ManualDiscountCurrency      *string    `gorm:"column:manual_discount_per_unit_currency;type:varchar(3)"`
	FinalDiscountAmount         float64    `gorm:"column:final_discount_per_unit_amount;type:numeric(12,2);not null"`
	FinalDiscountCurrency       string     `gorm:"column:final_discount_per_unit_currency;type:varchar(3);not null"`
	SubtotalAmount              float64    `gorm:"column:subtotal_amount;type:numeric(12,2);not null"`
	SubtotalCurrency            string     `gorm:"column:subtotal_currency;type:varchar(3);not null"`
}

func (QuoteLineItemDataModel) TableName() string {
	return "quote_line_items"
}

type SalesOrderDataModel struct {
	gorm.Model
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;"`
	OrderNumber      string     `gorm:"column:order_number;not null"`
	QuoteID          *uuid.UUID `gorm:"column:quote_id"`
	PartyID          uuid.UUID  `gorm:"column:party_id;not null"`
	OrderDate        time.Time  `gorm:"column:order_date;not null"`
	DeliveryDate     time.Time  `gorm:"column:delivery_date;not null"`
	Status           string     `gorm:"type:sales_order_status;not null"`
	SubtotalAmount   float64    `gorm:"column:subtotal_amount;type:numeric(12,2);not null"`
	SubtotalCurrency string     `gorm:"column:subtotal_currency;type:varchar(3);not null"`
	TaxAmount        float64    `gorm:"column:tax_amount;type:numeric(12,2);not null"`
	TaxCurrency      string     `gorm:"column:tax_currency;type:varchar(3);not null"`
	TotalAmount      float64    `gorm:"column:total_amount;type:numeric(12,2);not null"`
	TotalCurrency    string     `gorm:"column:total_currency;type:varchar(3);not null"`
	Notes            *string    `gorm:"column:notes"`
}

func (SalesOrderDataModel) TableName() string {
	return "sales_orders"
}

type OrderLineItemDataModel struct {
	gorm.Model
	ID                          uuid.UUID  `gorm:"type:uuid;primary_key;"`
	SalesOrderID                uuid.UUID  `gorm:"column:sales_order_id;not null"`
	MESWorkID                   *uuid.UUID `gorm:"column:mes_work_id;type:uuid"`
	ProductVariantID            uuid.UUID  `gorm:"column:product_variant_id;not null"`
	Quantity                    int        `gorm:"column:quantity;not null"`
	CalculatedUnitPriceAmount   float64    `gorm:"column:calculated_unit_price_amount;type:numeric(12,2);not null"`
	CalculatedUnitPriceCurrency string     `gorm:"column:calculated_unit_price_currency;type:varchar(3);not null"`
	ManualUnitPriceAmount       *float64   `gorm:"column:manual_unit_price_amount;type:numeric(12,2)"`
	ManualUnitPriceCurrency     *string    `gorm:"column:manual_unit_price_currency;type:varchar(3)"`
	FinalUnitPriceAmount        float64    `gorm:"column:final_unit_price_amount;type:numeric(12,2);not null"`
	FinalUnitPriceCurrency      string     `gorm:"column:final_unit_price_currency;type:varchar(3);not null"`
	CalculatedDiscountAmount    *float64   `gorm:"column:calculated_discount_per_unit_amount;type:numeric(12,2)"`
	CalculatedDiscountCurrency  *string    `gorm:"column:calculated_discount_per_unit_currency;type:varchar(3)"`
	ManualDiscountAmount        *float64   `gorm:"column:manual_discount_per_unit_amount;type:numeric(12,2)"`
	ManualDiscountCurrency      *string    `gorm:"column:manual_discount_per_unit_currency;type:varchar(3)"`
	FinalDiscountAmount         float64    `gorm:"column:final_discount_per_unit_amount;type:numeric(12,2);not null"`
	FinalDiscountCurrency       string     `gorm:"column:final_discount_per_unit_currency;type:varchar(3);not null"`
	SubtotalAmount              float64    `gorm:"column:subtotal_amount;type:numeric(12,2);not null"`
	SubtotalCurrency            string     `gorm:"column:subtotal_currency;type:varchar(3);not null"`
}

func (OrderLineItemDataModel) TableName() string {
	return "order_line_items"
}

type DeliveryNoteDataModel struct {
	gorm.Model
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;"`
	DeliveryNoteNumber string    `gorm:"column:delivery_note_number;not null"`
	SalesOrderID       uuid.UUID `gorm:"column:sales_order_id;not null"`
	PartyID            uuid.UUID `gorm:"column:party_id;not null"`
	DeliveryDate       time.Time `gorm:"column:delivery_date;not null"`
	Status             string    `gorm:"type:delivery_note_status;not null"`
	Notes              *string   `gorm:"column:notes"`
}

func (DeliveryNoteDataModel) TableName() string {
	return "delivery_notes"
}

type DeliveryNoteLineItemDataModel struct {
	gorm.Model
	ID                   uuid.UUID `gorm:"type:uuid;primary_key;"`
	DeliveryNoteID       uuid.UUID `gorm:"column:delivery_note_id;not null"`
	SalesOrderLineItemID uuid.UUID `gorm:"column:sales_order_line_item_id;not null"`
	ProductVariantID     uuid.UUID `gorm:"column:product_variant_id;not null"`
	DeliveredQuantity    int       `gorm:"column:delivered_quantity;not null"`
}

func (DeliveryNoteLineItemDataModel) TableName() string {
	return "delivery_note_line_items"
}

type InvoiceDataModel struct {
	gorm.Model
	ID               uuid.UUID `gorm:"type:uuid;primary_key;"`
	InvoiceNumber    string    `gorm:"column:invoice_number;not null"`
	Type             string    `gorm:"type:invoice_type;not null"` // COMPLETA | SIMPLIFICADA
	SeriesCode       string    `gorm:"column:series_code;not null"`
	SeriesYear       int       `gorm:"column:series_year;not null"`
	SeriesPrefix     string    `gorm:"column:series_prefix;not null"`
	PartyID          uuid.UUID `gorm:"column:party_id;not null"`
	InvoiceDate      time.Time `gorm:"column:invoice_date;not null"`
	DueDate          time.Time `gorm:"column:due_date;not null"`
	Status           string    `gorm:"type:invoice_status;not null"`
	PaymentTerms     *string   `gorm:"column:payment_terms"`
	SubtotalAmount   float64   `gorm:"column:subtotal_amount;type:numeric(12,2);not null"`
	SubtotalCurrency string    `gorm:"column:subtotal_currency;type:varchar(3);not null"`
	TaxAmount        float64   `gorm:"column:tax_amount;type:numeric(12,2);not null"`
	TaxCurrency      string    `gorm:"column:tax_currency;type:varchar(3);not null"`
	TotalAmount      float64   `gorm:"column:total_amount;type:numeric(12,2);not null"`
	TotalCurrency    string    `gorm:"column:total_currency;type:varchar(3);not null"`
}

func (InvoiceDataModel) TableName() string {
	return "invoices"
}

type InvoiceLineItemDataModel struct {
	gorm.Model
	ID                   uuid.UUID  `gorm:"type:uuid;primary_key;"`
	InvoiceID            uuid.UUID  `gorm:"column:invoice_id;not null"`
	SalesOrderLineItemID *uuid.UUID `gorm:"column:sales_order_line_item_id"`
	ProductVariantID     uuid.UUID  `gorm:"column:product_variant_id;not null"`
	Quantity             int        `gorm:"column:quantity;not null"`
	UnitPriceAmount      float64    `gorm:"column:unit_price_amount;type:numeric(12,2);not null"`
	UnitPriceCurrency    string     `gorm:"column:unit_price_currency;type:varchar(3);not null"`
	DiscountAmount       *float64   `gorm:"column:discount_amount;type:numeric(12,2)"`
	DiscountCurrency     *string    `gorm:"column:discount_currency;type:varchar(3)"`
	SubtotalAmount       float64    `gorm:"column:subtotal_amount;type:numeric(12,2);not null"`
	SubtotalCurrency     string     `gorm:"column:subtotal_currency;type:varchar(3);not null"`
	TaxAmount            *float64   `gorm:"column:tax_amount;type:numeric(12,2)"`
	TaxCurrency          *string    `gorm:"column:tax_currency;type:varchar(3)"`
	TotalAmount          float64    `gorm:"column:total_amount;type:numeric(12,2);not null"`
	TotalCurrency        string     `gorm:"column:total_currency;type:varchar(3);not null"`
}

func (InvoiceLineItemDataModel) TableName() string {
	return "invoice_line_items"
}

func quoteFromDomain(quote *domain.Quote) (*QuoteDataModel, error) {
	notes := toOptionalString(quote.Notes)
	return &QuoteDataModel{
		ID:               quote.ID,
		QuoteNumber:      quote.QuoteNumber.String(),
		PartyID:          quote.PartyID,
		QuoteDate:        quote.QuoteDate,
		ExpirationDate:   quote.ExpirationDate,
		Status:           string(quote.Status),
		SubtotalAmount:   quote.Subtotal.Amount(),
		SubtotalCurrency: quote.Subtotal.Currency(),
		TaxAmount:        quote.TaxAmount.Amount(),
		TaxCurrency:      quote.TaxAmount.Currency(),
		TotalAmount:      quote.Total.Amount(),
		TotalCurrency:    quote.Total.Currency(),
		Notes:            notes,
	}, nil
}

func quoteLineItemsFromDomain(quoteID uuid.UUID, items []domain.QuoteLineItem) ([]QuoteLineItemDataModel, error) {
	models := make([]QuoteLineItemDataModel, 0, len(items))
	for _, item := range items {
		model, err := quoteLineItemFromDomain(quoteID, item)
		if err != nil {
			return nil, err
		}
		models = append(models, *model)
	}
	return models, nil
}

func quoteLineItemFromDomain(quoteID uuid.UUID, item domain.QuoteLineItem) (*QuoteLineItemDataModel, error) {
	return &QuoteLineItemDataModel{
		ID:                          item.ID,
		QuoteID:                     quoteID,
		MESWorkID:                   item.MESWorkID,
		ProductVariantID:            item.ProductVariantID,
		Quantity:                    item.Quantity,
		CalculatedUnitPriceAmount:   item.CalculatedUnitPrice.Amount(),
		CalculatedUnitPriceCurrency: item.CalculatedUnitPrice.Currency(),
		ManualUnitPriceAmount:       optionalAmount(item.ManualUnitPrice),
		ManualUnitPriceCurrency:     optionalCurrency(item.ManualUnitPrice),
		FinalUnitPriceAmount:        item.FinalUnitPrice.Amount(),
		FinalUnitPriceCurrency:      item.FinalUnitPrice.Currency(),
		CalculatedDiscountAmount:    optionalAmount(item.CalculatedDiscountPerUnit),
		CalculatedDiscountCurrency:  optionalCurrency(item.CalculatedDiscountPerUnit),
		ManualDiscountAmount:        optionalAmount(item.ManualDiscountPerUnit),
		ManualDiscountCurrency:      optionalCurrency(item.ManualDiscountPerUnit),
		FinalDiscountAmount:         item.FinalDiscountPerUnit.Amount(),
		FinalDiscountCurrency:       item.FinalDiscountPerUnit.Currency(),
		SubtotalAmount:              item.Subtotal.Amount(),
		SubtotalCurrency:            item.Subtotal.Currency(),
	}, nil
}

func salesOrderFromDomain(order *domain.SalesOrder) (*SalesOrderDataModel, error) {
	notes := toOptionalString(order.Notes)
	return &SalesOrderDataModel{
		ID:               order.ID,
		OrderNumber:      order.OrderNumber.String(),
		QuoteID:          order.QuoteID,
		PartyID:          order.PartyID,
		OrderDate:        order.OrderDate,
		DeliveryDate:     order.DeliveryDate,
		Status:           string(order.Status),
		SubtotalAmount:   order.Subtotal.Amount(),
		SubtotalCurrency: order.Subtotal.Currency(),
		TaxAmount:        order.TaxAmount.Amount(),
		TaxCurrency:      order.TaxAmount.Currency(),
		TotalAmount:      order.Total.Amount(),
		TotalCurrency:    order.Total.Currency(),
		Notes:            notes,
	}, nil
}

func orderLineItemsFromDomain(orderID uuid.UUID, items []domain.OrderLineItem) ([]OrderLineItemDataModel, error) {
	models := make([]OrderLineItemDataModel, 0, len(items))
	for _, item := range items {
		model, err := orderLineItemFromDomain(orderID, item)
		if err != nil {
			return nil, err
		}
		models = append(models, *model)
	}
	return models, nil
}

func orderLineItemFromDomain(orderID uuid.UUID, item domain.OrderLineItem) (*OrderLineItemDataModel, error) {
	return &OrderLineItemDataModel{
		ID:                          item.ID,
		SalesOrderID:                orderID,
		MESWorkID:                   item.MESWorkID,
		ProductVariantID:            item.ProductVariantID,
		Quantity:                    item.Quantity,
		CalculatedUnitPriceAmount:   item.CalculatedUnitPrice.Amount(),
		CalculatedUnitPriceCurrency: item.CalculatedUnitPrice.Currency(),
		ManualUnitPriceAmount:       optionalAmount(item.ManualUnitPrice),
		ManualUnitPriceCurrency:     optionalCurrency(item.ManualUnitPrice),
		FinalUnitPriceAmount:        item.FinalUnitPrice.Amount(),
		FinalUnitPriceCurrency:      item.FinalUnitPrice.Currency(),
		CalculatedDiscountAmount:    optionalAmount(item.CalculatedDiscountPerUnit),
		CalculatedDiscountCurrency:  optionalCurrency(item.CalculatedDiscountPerUnit),
		ManualDiscountAmount:        optionalAmount(item.ManualDiscountPerUnit),
		ManualDiscountCurrency:      optionalCurrency(item.ManualDiscountPerUnit),
		FinalDiscountAmount:         item.FinalDiscountPerUnit.Amount(),
		FinalDiscountCurrency:       item.FinalDiscountPerUnit.Currency(),
		SubtotalAmount:              item.Subtotal.Amount(),
		SubtotalCurrency:            item.Subtotal.Currency(),
	}, nil
}

func deliveryNoteFromDomain(note *domain.DeliveryNote) (*DeliveryNoteDataModel, error) {
	notes := toOptionalString(note.Notes)
	return &DeliveryNoteDataModel{
		ID:                 note.ID,
		DeliveryNoteNumber: note.DeliveryNoteNumber.String(),
		SalesOrderID:       note.SalesOrderID,
		PartyID:            note.PartyID,
		DeliveryDate:       note.DeliveryDate,
		Status:             string(note.Status),
		Notes:              notes,
	}, nil
}

func deliveryNoteLineItemsFromDomain(noteID uuid.UUID, items []domain.DeliveryNoteLineItem) ([]DeliveryNoteLineItemDataModel, error) {
	models := make([]DeliveryNoteLineItemDataModel, 0, len(items))
	for _, item := range items {
		models = append(models, DeliveryNoteLineItemDataModel{
			ID:                   item.ID,
			DeliveryNoteID:       noteID,
			SalesOrderLineItemID: item.SalesOrderLineItemID,
			ProductVariantID:     item.ProductVariantID,
			DeliveredQuantity:    item.DeliveredQuantity,
		})
	}
	return models, nil
}

func invoiceFromDomain(invoice *domain.Invoice) (*InvoiceDataModel, error) {
	terms := toOptionalString(invoice.PaymentTerms)
	return &InvoiceDataModel{
		ID:               invoice.ID,
		InvoiceNumber:    invoice.InvoiceNumber.String(),
		Type:             string(invoice.Type),
		SeriesCode:       invoice.Series.Code(),
		SeriesYear:       invoice.Series.Year(),
		SeriesPrefix:     invoice.Series.Prefix(),
		PartyID:          invoice.PartyID,
		InvoiceDate:      invoice.InvoiceDate,
		DueDate:          invoice.DueDate,
		Status:           string(invoice.Status),
		PaymentTerms:     terms,
		SubtotalAmount:   invoice.Subtotal.Amount(),
		SubtotalCurrency: invoice.Subtotal.Currency(),
		TaxAmount:        invoice.TaxAmount.Amount(),
		TaxCurrency:      invoice.TaxAmount.Currency(),
		TotalAmount:      invoice.Total.Amount(),
		TotalCurrency:    invoice.Total.Currency(),
	}, nil
}

func invoiceLineItemsFromDomain(invoiceID uuid.UUID, items []domain.InvoiceLineItem) ([]InvoiceLineItemDataModel, error) {
	models := make([]InvoiceLineItemDataModel, 0, len(items))
	for _, item := range items {
		models = append(models, InvoiceLineItemDataModel{
			ID:                   item.ID,
			InvoiceID:            invoiceID,
			SalesOrderLineItemID: item.SalesOrderLineItemID,
			ProductVariantID:     item.ProductVariantID,
			Quantity:             item.Quantity,
			UnitPriceAmount:      item.UnitPrice.Amount(),
			UnitPriceCurrency:    item.UnitPrice.Currency(),
			DiscountAmount:       optionalAmount(item.DiscountAmount),
			DiscountCurrency:     optionalCurrency(item.DiscountAmount),
			SubtotalAmount:       item.Subtotal.Amount(),
			SubtotalCurrency:     item.Subtotal.Currency(),
			TaxAmount:            optionalAmount(item.TaxAmount),
			TaxCurrency:          optionalCurrency(item.TaxAmount),
			TotalAmount:          item.Total.Amount(),
			TotalCurrency:        item.Total.Currency(),
		})
	}
	return models, nil
}

func quoteToDomain(quote *QuoteDataModel, items []QuoteLineItemDataModel) (*domain.Quote, error) {
	subtotal, err := moneyFromParts(quote.SubtotalAmount, quote.SubtotalCurrency)
	if err != nil {
		return nil, err
	}
	tax, err := moneyFromParts(quote.TaxAmount, quote.TaxCurrency)
	if err != nil {
		return nil, err
	}
	total, err := moneyFromParts(quote.TotalAmount, quote.TotalCurrency)
	if err != nil {
		return nil, err
	}

	lineItems := make([]domain.QuoteLineItem, 0, len(items))
	for _, item := range items {
		line, err := quoteLineItemToDomain(item)
		if err != nil {
			return nil, err
		}
		lineItems = append(lineItems, line)
	}

	quoteNumber, err := domain.NewQuoteNumber(quote.QuoteNumber)
	if err != nil {
		return nil, err
	}

	return &domain.Quote{
		ID:             quote.ID,
		QuoteNumber:    quoteNumber,
		PartyID:        quote.PartyID,
		QuoteDate:      quote.QuoteDate,
		ExpirationDate: quote.ExpirationDate,
		Status:         domain.QuoteStatus(quote.Status),
		LineItems:      lineItems,
		Subtotal:       subtotal,
		TaxAmount:      tax,
		Total:          total,
		Notes:          optionalStringValue(quote.Notes),
	}, nil
}

func quoteLineItemToDomain(item QuoteLineItemDataModel) (domain.QuoteLineItem, error) {
	calculatedUnit, err := moneyFromParts(item.CalculatedUnitPriceAmount, item.CalculatedUnitPriceCurrency)
	if err != nil {
		return domain.QuoteLineItem{}, err
	}
	finalUnit, err := moneyFromParts(item.FinalUnitPriceAmount, item.FinalUnitPriceCurrency)
	if err != nil {
		return domain.QuoteLineItem{}, err
	}
	finalDiscount, err := moneyFromParts(item.FinalDiscountAmount, item.FinalDiscountCurrency)
	if err != nil {
		return domain.QuoteLineItem{}, err
	}
	subtotal, err := moneyFromParts(item.SubtotalAmount, item.SubtotalCurrency)
	if err != nil {
		return domain.QuoteLineItem{}, err
	}

	manualUnit, err := optionalMoneyFromParts(item.ManualUnitPriceAmount, item.ManualUnitPriceCurrency)
	if err != nil {
		return domain.QuoteLineItem{}, err
	}
	calculatedDiscount, err := optionalMoneyFromParts(item.CalculatedDiscountAmount, item.CalculatedDiscountCurrency)
	if err != nil {
		return domain.QuoteLineItem{}, err
	}
	manualDiscount, err := optionalMoneyFromParts(item.ManualDiscountAmount, item.ManualDiscountCurrency)
	if err != nil {
		return domain.QuoteLineItem{}, err
	}

	return domain.QuoteLineItem{
		ID:                        item.ID,
		MESWorkID:                 item.MESWorkID,
		ProductVariantID:          item.ProductVariantID,
		Quantity:                  item.Quantity,
		CalculatedUnitPrice:       calculatedUnit,
		ManualUnitPrice:           manualUnit,
		FinalUnitPrice:            finalUnit,
		CalculatedDiscountPerUnit: calculatedDiscount,
		ManualDiscountPerUnit:     manualDiscount,
		FinalDiscountPerUnit:      finalDiscount,
		Subtotal:                  subtotal,
	}, nil
}

func salesOrderToDomain(order *SalesOrderDataModel, items []OrderLineItemDataModel) (*domain.SalesOrder, error) {
	subtotal, err := moneyFromParts(order.SubtotalAmount, order.SubtotalCurrency)
	if err != nil {
		return nil, err
	}
	tax, err := moneyFromParts(order.TaxAmount, order.TaxCurrency)
	if err != nil {
		return nil, err
	}
	total, err := moneyFromParts(order.TotalAmount, order.TotalCurrency)
	if err != nil {
		return nil, err
	}

	lineItems := make([]domain.OrderLineItem, 0, len(items))
	for _, item := range items {
		line, err := orderLineItemToDomain(item)
		if err != nil {
			return nil, err
		}
		lineItems = append(lineItems, line)
	}

	orderNumber, err := domain.NewOrderNumber(order.OrderNumber)
	if err != nil {
		return nil, err
	}

	return &domain.SalesOrder{
		ID:           order.ID,
		OrderNumber:  orderNumber,
		QuoteID:      order.QuoteID,
		PartyID:      order.PartyID,
		OrderDate:    order.OrderDate,
		DeliveryDate: order.DeliveryDate,
		Status:       domain.SalesOrderStatus(order.Status),
		LineItems:    lineItems,
		Subtotal:     subtotal,
		TaxAmount:    tax,
		Total:        total,
		Notes:        optionalStringValue(order.Notes),
	}, nil
}

func orderLineItemToDomain(item OrderLineItemDataModel) (domain.OrderLineItem, error) {
	calculatedUnit, err := moneyFromParts(item.CalculatedUnitPriceAmount, item.CalculatedUnitPriceCurrency)
	if err != nil {
		return domain.OrderLineItem{}, err
	}
	finalUnit, err := moneyFromParts(item.FinalUnitPriceAmount, item.FinalUnitPriceCurrency)
	if err != nil {
		return domain.OrderLineItem{}, err
	}
	finalDiscount, err := moneyFromParts(item.FinalDiscountAmount, item.FinalDiscountCurrency)
	if err != nil {
		return domain.OrderLineItem{}, err
	}
	subtotal, err := moneyFromParts(item.SubtotalAmount, item.SubtotalCurrency)
	if err != nil {
		return domain.OrderLineItem{}, err
	}

	manualUnit, err := optionalMoneyFromParts(item.ManualUnitPriceAmount, item.ManualUnitPriceCurrency)
	if err != nil {
		return domain.OrderLineItem{}, err
	}
	calculatedDiscount, err := optionalMoneyFromParts(item.CalculatedDiscountAmount, item.CalculatedDiscountCurrency)
	if err != nil {
		return domain.OrderLineItem{}, err
	}
	manualDiscount, err := optionalMoneyFromParts(item.ManualDiscountAmount, item.ManualDiscountCurrency)
	if err != nil {
		return domain.OrderLineItem{}, err
	}

	return domain.OrderLineItem{
		ID:                        item.ID,
		MESWorkID:                 item.MESWorkID,
		ProductVariantID:          item.ProductVariantID,
		Quantity:                  item.Quantity,
		CalculatedUnitPrice:       calculatedUnit,
		ManualUnitPrice:           manualUnit,
		FinalUnitPrice:            finalUnit,
		CalculatedDiscountPerUnit: calculatedDiscount,
		ManualDiscountPerUnit:     manualDiscount,
		FinalDiscountPerUnit:      finalDiscount,
		Subtotal:                  subtotal,
	}, nil
}

func deliveryNoteToDomain(note *DeliveryNoteDataModel, items []DeliveryNoteLineItemDataModel) (*domain.DeliveryNote, error) {
	lineItems := make([]domain.DeliveryNoteLineItem, 0, len(items))
	for _, item := range items {
		lineItems = append(lineItems, domain.DeliveryNoteLineItem{
			ID:                   item.ID,
			SalesOrderLineItemID: item.SalesOrderLineItemID,
			ProductVariantID:     item.ProductVariantID,
			DeliveredQuantity:    item.DeliveredQuantity,
		})
	}

	noteNumber, err := domain.NewDeliveryNoteNumber(note.DeliveryNoteNumber)
	if err != nil {
		return nil, err
	}

	return &domain.DeliveryNote{
		ID:                 note.ID,
		DeliveryNoteNumber: noteNumber,
		SalesOrderID:       note.SalesOrderID,
		PartyID:            note.PartyID,
		DeliveryDate:       note.DeliveryDate,
		Status:             domain.DeliveryNoteStatus(note.Status),
		LineItems:          lineItems,
		Notes:              optionalStringValue(note.Notes),
	}, nil
}

func invoiceToDomain(invoice *InvoiceDataModel, items []InvoiceLineItemDataModel) (*domain.Invoice, error) {
	subtotal, err := moneyFromParts(invoice.SubtotalAmount, invoice.SubtotalCurrency)
	if err != nil {
		return nil, err
	}
	tax, err := moneyFromParts(invoice.TaxAmount, invoice.TaxCurrency)
	if err != nil {
		return nil, err
	}
	total, err := moneyFromParts(invoice.TotalAmount, invoice.TotalCurrency)
	if err != nil {
		return nil, err
	}

	lineItems := make([]domain.InvoiceLineItem, 0, len(items))
	for _, item := range items {
		line, err := invoiceLineItemToDomain(item)
		if err != nil {
			return nil, err
		}
		lineItems = append(lineItems, line)
	}

	invoiceNumber, err := domain.NewInvoiceNumber(invoice.InvoiceNumber)
	if err != nil {
		return nil, err
	}

	invoiceType := domain.InvoiceType(invoice.Type)
	if err := invoiceType.IsValid(); err != nil {
		return nil, err
	}

	series, err := domain.NewInvoiceSeriesWithPrefix(invoice.SeriesCode, invoice.SeriesYear, invoice.SeriesPrefix)
	if err != nil {
		return nil, err
	}

	return &domain.Invoice{
		ID:            invoice.ID,
		InvoiceNumber: invoiceNumber,
		Type:          invoiceType,
		Series:        series,
		PartyID:       invoice.PartyID,
		InvoiceDate:   invoice.InvoiceDate,
		DueDate:       invoice.DueDate,
		Status:        domain.InvoiceStatus(invoice.Status),
		LineItems:     lineItems,
		Subtotal:      subtotal,
		TaxAmount:     tax,
		Total:         total,
		PaymentTerms:  optionalStringValue(invoice.PaymentTerms),
	}, nil
}

func invoiceLineItemToDomain(item InvoiceLineItemDataModel) (domain.InvoiceLineItem, error) {
	unitPrice, err := moneyFromParts(item.UnitPriceAmount, item.UnitPriceCurrency)
	if err != nil {
		return domain.InvoiceLineItem{}, err
	}
	subtotal, err := moneyFromParts(item.SubtotalAmount, item.SubtotalCurrency)
	if err != nil {
		return domain.InvoiceLineItem{}, err
	}
	total, err := moneyFromParts(item.TotalAmount, item.TotalCurrency)
	if err != nil {
		return domain.InvoiceLineItem{}, err
	}

	discount, err := optionalMoneyFromParts(item.DiscountAmount, item.DiscountCurrency)
	if err != nil {
		return domain.InvoiceLineItem{}, err
	}
	tax, err := optionalMoneyFromParts(item.TaxAmount, item.TaxCurrency)
	if err != nil {
		return domain.InvoiceLineItem{}, err
	}

	return domain.InvoiceLineItem{
		ID:                   item.ID,
		SalesOrderLineItemID: item.SalesOrderLineItemID,
		ProductVariantID:     item.ProductVariantID,
		Quantity:             item.Quantity,
		UnitPrice:            unitPrice,
		DiscountAmount:       discount,
		Subtotal:             subtotal,
		TaxAmount:            tax,
		Total:                total,
	}, nil
}

func moneyFromParts(amount float64, currency string) (domain.Money, error) {
	value := strings.TrimSpace(currency)
	if value == "" {
		value = domain.DefaultCurrency
	}
	return domain.NewMoney(amount, value)
}

func optionalMoneyFromParts(amount *float64, currency *string) (*domain.Money, error) {
	if amount == nil {
		return nil, nil
	}
	value := domain.DefaultCurrency
	if currency != nil && strings.TrimSpace(*currency) != "" {
		value = *currency
	}
	money, err := domain.NewMoney(*amount, value)
	if err != nil {
		return nil, err
	}
	return &money, nil
}

func toOptionalString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func optionalStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func optionalAmount(money *domain.Money) *float64 {
	if money == nil {
		return nil
	}
	value := money.Amount()
	return &value
}

func optionalCurrency(money *domain.Money) *string {
	if money == nil {
		return nil
	}
	value := money.Currency()
	return &value
}
