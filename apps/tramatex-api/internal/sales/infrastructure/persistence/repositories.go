package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/sales/domain"
)

type GORMQuoteRepository struct {
	db *gorm.DB
}

func NewGORMQuoteRepository(db *gorm.DB) *GORMQuoteRepository {
	return &GORMQuoteRepository{db: db}
}

func (r *GORMQuoteRepository) Save(ctx context.Context, quote *domain.Quote) error {
	data, err := quoteFromDomain(quote)
	if err != nil {
		return err
	}
	items, err := quoteLineItemsFromDomain(quote.ID, quote.LineItems)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(data).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("quote_id = ?", quote.ID).Delete(&QuoteLineItemDataModel{}).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		return tx.Create(&items).Error
	})
}

func (r *GORMQuoteRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Quote, error) {
	var data QuoteDataModel
	if err := r.db.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	items, err := r.loadQuoteLineItems(ctx, id)
	if err != nil {
		return nil, err
	}

	return quoteToDomain(&data, items)
}

func (r *GORMQuoteRepository) List(ctx context.Context, filter domain.QuoteFilter) ([]*domain.Quote, error) {
	query := r.db.WithContext(ctx).Model(&QuoteDataModel{})
	if filter.PartyID != nil {
		query = query.Where("party_id = ?", *filter.PartyID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.FromDate != nil {
		query = query.Where("quote_date >= ?", *filter.FromDate)
	}
	if filter.ToDate != nil {
		query = query.Where("quote_date <= ?", *filter.ToDate)
	}

	var data []QuoteDataModel
	if err := query.Order("created_at desc").Find(&data).Error; err != nil {
		return nil, err
	}

	result := make([]*domain.Quote, 0, len(data))
	for i := range data {
		items, err := r.loadQuoteLineItems(ctx, data[i].ID)
		if err != nil {
			return nil, err
		}
		mapped, err := quoteToDomain(&data[i], items)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}
	return result, nil
}

func (r *GORMQuoteRepository) loadQuoteLineItems(ctx context.Context, quoteID uuid.UUID) ([]QuoteLineItemDataModel, error) {
	var items []QuoteLineItemDataModel
	if err := r.db.WithContext(ctx).Where("quote_id = ?", quoteID).Order("created_at asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type GORMSalesOrderRepository struct {
	db *gorm.DB
}

func NewGORMSalesOrderRepository(db *gorm.DB) *GORMSalesOrderRepository {
	return &GORMSalesOrderRepository{db: db}
}

func (r *GORMSalesOrderRepository) Save(ctx context.Context, order *domain.SalesOrder) error {
	data, err := salesOrderFromDomain(order)
	if err != nil {
		return err
	}
	items, err := orderLineItemsFromDomain(order.ID, order.LineItems)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(data).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("sales_order_id = ?", order.ID).Delete(&OrderLineItemDataModel{}).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		return tx.Create(&items).Error
	})
}

func (r *GORMSalesOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	var data SalesOrderDataModel
	if err := r.db.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	items, err := r.loadOrderLineItems(ctx, id)
	if err != nil {
		return nil, err
	}

	return salesOrderToDomain(&data, items)
}

func (r *GORMSalesOrderRepository) List(ctx context.Context, filter domain.SalesOrderFilter) ([]*domain.SalesOrder, error) {
	query := r.db.WithContext(ctx).Model(&SalesOrderDataModel{})
	if filter.PartyID != nil {
		query = query.Where("party_id = ?", *filter.PartyID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.FromDate != nil {
		query = query.Where("order_date >= ?", *filter.FromDate)
	}
	if filter.ToDate != nil {
		query = query.Where("order_date <= ?", *filter.ToDate)
	}

	var data []SalesOrderDataModel
	if err := query.Order("created_at desc").Find(&data).Error; err != nil {
		return nil, err
	}

	result := make([]*domain.SalesOrder, 0, len(data))
	for i := range data {
		items, err := r.loadOrderLineItems(ctx, data[i].ID)
		if err != nil {
			return nil, err
		}
		mapped, err := salesOrderToDomain(&data[i], items)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}
	return result, nil
}

func (r *GORMSalesOrderRepository) loadOrderLineItems(ctx context.Context, orderID uuid.UUID) ([]OrderLineItemDataModel, error) {
	var items []OrderLineItemDataModel
	if err := r.db.WithContext(ctx).Where("sales_order_id = ?", orderID).Order("created_at asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type GORMDeliveryNoteRepository struct {
	db *gorm.DB
}

func NewGORMDeliveryNoteRepository(db *gorm.DB) *GORMDeliveryNoteRepository {
	return &GORMDeliveryNoteRepository{db: db}
}

func (r *GORMDeliveryNoteRepository) Save(ctx context.Context, note *domain.DeliveryNote) error {
	data, err := deliveryNoteFromDomain(note)
	if err != nil {
		return err
	}
	items, err := deliveryNoteLineItemsFromDomain(note.ID, note.LineItems)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(data).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("delivery_note_id = ?", note.ID).Delete(&DeliveryNoteLineItemDataModel{}).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		return tx.Create(&items).Error
	})
}

func (r *GORMDeliveryNoteRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.DeliveryNote, error) {
	var data DeliveryNoteDataModel
	if err := r.db.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	items, err := r.loadDeliveryNoteLineItems(ctx, id)
	if err != nil {
		return nil, err
	}

	return deliveryNoteToDomain(&data, items)
}

func (r *GORMDeliveryNoteRepository) List(ctx context.Context, filter domain.DeliveryNoteFilter) ([]*domain.DeliveryNote, error) {
	query := r.db.WithContext(ctx).Model(&DeliveryNoteDataModel{})
	if filter.SalesOrderID != nil {
		query = query.Where("sales_order_id = ?", *filter.SalesOrderID)
	}
	if filter.PartyID != nil {
		query = query.Where("party_id = ?", *filter.PartyID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.FromDate != nil {
		query = query.Where("delivery_date >= ?", *filter.FromDate)
	}
	if filter.ToDate != nil {
		query = query.Where("delivery_date <= ?", *filter.ToDate)
	}

	var data []DeliveryNoteDataModel
	if err := query.Order("created_at desc").Find(&data).Error; err != nil {
		return nil, err
	}

	result := make([]*domain.DeliveryNote, 0, len(data))
	for i := range data {
		items, err := r.loadDeliveryNoteLineItems(ctx, data[i].ID)
		if err != nil {
			return nil, err
		}
		mapped, err := deliveryNoteToDomain(&data[i], items)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}
	return result, nil
}

func (r *GORMDeliveryNoteRepository) ListBySalesOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.DeliveryNote, error) {
	filter := domain.DeliveryNoteFilter{SalesOrderID: &orderID}
	return r.List(ctx, filter)
}

func (r *GORMDeliveryNoteRepository) loadDeliveryNoteLineItems(ctx context.Context, noteID uuid.UUID) ([]DeliveryNoteLineItemDataModel, error) {
	var items []DeliveryNoteLineItemDataModel
	if err := r.db.WithContext(ctx).Where("delivery_note_id = ?", noteID).Order("created_at asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type GORMInvoiceRepository struct {
	db *gorm.DB
}

func NewGORMInvoiceRepository(db *gorm.DB) *GORMInvoiceRepository {
	return &GORMInvoiceRepository{db: db}
}

func (r *GORMInvoiceRepository) Save(ctx context.Context, invoice *domain.Invoice) error {
	data, err := invoiceFromDomain(invoice)
	if err != nil {
		return err
	}
	items, err := invoiceLineItemsFromDomain(invoice.ID, invoice.LineItems)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(data).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("invoice_id = ?", invoice.ID).Delete(&InvoiceLineItemDataModel{}).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		return tx.Create(&items).Error
	})
}

func (r *GORMInvoiceRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	var data InvoiceDataModel
	if err := r.db.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	items, err := r.loadInvoiceLineItems(ctx, id)
	if err != nil {
		return nil, err
	}

	return invoiceToDomain(&data, items)
}

func (r *GORMInvoiceRepository) List(ctx context.Context, filter domain.InvoiceFilter) ([]*domain.Invoice, error) {
	query := r.db.WithContext(ctx).Model(&InvoiceDataModel{})
	if filter.PartyID != nil {
		query = query.Where("party_id = ?", *filter.PartyID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.FromDate != nil {
		query = query.Where("invoice_date >= ?", *filter.FromDate)
	}
	if filter.ToDate != nil {
		query = query.Where("invoice_date <= ?", *filter.ToDate)
	}

	var data []InvoiceDataModel
	if err := query.Order("created_at desc").Find(&data).Error; err != nil {
		return nil, err
	}

	result := make([]*domain.Invoice, 0, len(data))
	for i := range data {
		items, err := r.loadInvoiceLineItems(ctx, data[i].ID)
		if err != nil {
			return nil, err
		}
		mapped, err := invoiceToDomain(&data[i], items)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}
	return result, nil
}

func (r *GORMInvoiceRepository) ListBySalesOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.Invoice, error) {
	var data []InvoiceDataModel
	query := r.db.WithContext(ctx).
		Model(&InvoiceDataModel{}).
		Joins("JOIN invoice_line_items ON invoice_line_items.invoice_id = invoices.id").
		Joins("JOIN order_line_items ON order_line_items.id = invoice_line_items.sales_order_line_item_id").
		Where("order_line_items.sales_order_id = ?", orderID).
		Distinct("invoices.id")

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	result := make([]*domain.Invoice, 0, len(data))
	for i := range data {
		items, err := r.loadInvoiceLineItems(ctx, data[i].ID)
		if err != nil {
			return nil, err
		}
		mapped, err := invoiceToDomain(&data[i], items)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}
	return result, nil
}

func (r *GORMInvoiceRepository) loadInvoiceLineItems(ctx context.Context, invoiceID uuid.UUID) ([]InvoiceLineItemDataModel, error) {
	var items []InvoiceLineItemDataModel
	if err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceID).Order("created_at asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
