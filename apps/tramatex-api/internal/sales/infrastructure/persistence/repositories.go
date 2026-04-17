package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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

	return getDB(ctx, r.db).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(data).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("quote_id = ?", quote.ID).Delete(&QuoteLineItemDataModel{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("quote_id = ?", quote.ID).Delete(&QuoteWorkRefModel{}).Error; err != nil {
			return err
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		if len(quote.WorkReferences) > 0 {
			wsModels := quoteWorkRefsFromDomain(quote.ID, quote.WorkReferences)
			if err := tx.Create(&wsModels).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GORMQuoteRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Quote, error) {
	var data QuoteDataModel
	if err := getDB(ctx, r.db).First(&data, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	items, err := r.loadQuoteLineItems(ctx, id)
	if err != nil {
		return nil, err
	}

	quote, qErr := quoteToDomain(&data, items)
	if qErr != nil {
		return nil, qErr
	}

	workSetups, err := r.loadQuoteWorkSetups(ctx, id)
	if err != nil {
		return nil, err
	}
	quote.WorkReferences = workSetups

	return quote, nil
}

func (r *GORMQuoteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return getDB(ctx, r.db).Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("quote_id = ?", id).Delete(&QuoteWorkRefModel{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("quote_id = ?", id).Delete(&QuoteLineItemDataModel{}).Error; err != nil {
			return err
		}
		return tx.Unscoped().Where("id = ?", id).Delete(&QuoteDataModel{}).Error
	})
}

func (r *GORMQuoteRepository) List(ctx context.Context, filter domain.QuoteFilter) ([]*domain.Quote, error) {
	query := getDB(ctx, r.db).Model(&QuoteDataModel{})
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
	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(`
			quote_number ILIKE ?
			OR EXISTS (
				SELECT 1
				FROM organization_profiles op
				WHERE op.party_id = CAST("quotes".party_id AS TEXT)
				  AND op.name ILIKE ?
			)
			OR EXISTS (
				SELECT 1
				FROM person_profiles pp
				WHERE pp.party_id = CAST(quotes.party_id AS TEXT)
				  AND (
					pp.first_name ILIKE ?
					OR pp.last_name ILIKE ?
					OR (pp.first_name || ' ' || pp.last_name) ILIKE ?
				  )
			)
			OR EXISTS (
				SELECT 1
				FROM party_roles pr
				WHERE pr.party_id = CAST(quotes.party_id AS TEXT)
				  AND pr.creation_identifier ILIKE ?
			)
		`, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
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
		workSetups, err := r.loadQuoteWorkSetups(ctx, data[i].ID)
		if err != nil {
			return nil, err
		}
		mapped.WorkReferences = workSetups
		result = append(result, mapped)
	}
	return result, nil
}

func (r *GORMQuoteRepository) loadQuoteLineItems(ctx context.Context, quoteID uuid.UUID) ([]QuoteLineItemDataModel, error) {
	var items []QuoteLineItemDataModel
	if err := getDB(ctx, r.db).Where("quote_id = ?", quoteID).Order("created_at asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GORMQuoteRepository) loadQuoteWorkSetups(ctx context.Context, quoteID uuid.UUID) ([]domain.WorkReference, error) {
	var rows []quoteWorkRefRow
	err := getDB(ctx, r.db).Raw(`
		SELECT qs.id, qs.work_setup_id, qs.sequence, qs.description
		FROM quote_work_setups qs
		WHERE qs.quote_id = ? AND qs.deleted_at IS NULL
		ORDER BY qs.sequence ASC
	`, quoteID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return quoteWorkRefRowsToDomain(rows), nil
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

	return getDB(ctx, r.db).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(data).Error; err != nil {
			return err
		}

		// Build a set of incoming line-item IDs so we only delete the removed ones.
		// This avoids FK violations when delivery_note_line_items reference existing items.
		incomingIDs := make(map[uuid.UUID]struct{}, len(items))
		for _, item := range items {
			incomingIDs[item.ID] = struct{}{}
		}

		var existingIDs []uuid.UUID
		if err := tx.Model(&OrderLineItemDataModel{}).
			Where("sales_order_id = ?", order.ID).
			Pluck("id", &existingIDs).Error; err != nil {
			return err
		}

		var toDelete []uuid.UUID
		for _, eid := range existingIDs {
			if _, keep := incomingIDs[eid]; !keep {
				toDelete = append(toDelete, eid)
			}
		}
		if len(toDelete) > 0 {
			if err := tx.Unscoped().Where("id IN ?", toDelete).Delete(&OrderLineItemDataModel{}).Error; err != nil {
				return err
			}
		}

		// Upsert remaining items
		for i := range items {
			if err := tx.Save(&items[i]).Error; err != nil {
				return err
			}
		}

		// Replace work setups
		if err := tx.Unscoped().Where("order_id = ?", order.ID).Delete(&OrderWorkRefModel{}).Error; err != nil {
			return err
		}
		if len(order.WorkReferences) > 0 {
			wsModels := orderWorkRefsFromDomain(order.ID, order.WorkReferences)
			if err := tx.Create(&wsModels).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GORMSalesOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	var data SalesOrderDataModel
	if err := getDB(ctx, r.db).First(&data, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	items, err := r.loadOrderLineItems(ctx, id)
	if err != nil {
		return nil, err
	}

	order, oErr := salesOrderToDomain(&data, items)
	if oErr != nil {
		return nil, oErr
	}

	workSetups, err := r.loadOrderWorkSetups(ctx, id)
	if err != nil {
		return nil, err
	}
	order.WorkReferences = workSetups

	return order, nil
}

func (r *GORMSalesOrderRepository) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	var data SalesOrderDataModel
	if err := getDB(ctx, r.db).Clauses(clause.Locking{Strength: "UPDATE"}).First(&data, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	items, err := r.loadOrderLineItems(ctx, id)
	if err != nil {
		return nil, err
	}

	order, oErr := salesOrderToDomain(&data, items)
	if oErr != nil {
		return nil, oErr
	}

	workSetups, err := r.loadOrderWorkSetups(ctx, id)
	if err != nil {
		return nil, err
	}
	order.WorkReferences = workSetups

	return order, nil
}

func (r *GORMSalesOrderRepository) List(ctx context.Context, filter domain.SalesOrderFilter) ([]*domain.SalesOrder, error) {
	query := getDB(ctx, r.db).Model(&SalesOrderDataModel{})
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
	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(`
			order_number ILIKE ?
			OR EXISTS (
				SELECT 1
				FROM organization_profiles op
				WHERE op.party_id = CAST(sales_orders.party_id AS TEXT)
				  AND op.name ILIKE ?
			)
			OR EXISTS (
				SELECT 1
				FROM person_profiles pp
				WHERE pp.party_id = CAST(sales_orders.party_id AS TEXT)
				  AND (
					pp.first_name ILIKE ?
					OR pp.last_name ILIKE ?
					OR (pp.first_name || ' ' || pp.last_name) ILIKE ?
				  )
			)
			OR EXISTS (
				SELECT 1
				FROM party_roles pr
				WHERE pr.party_id = CAST(sales_orders.party_id AS TEXT)
				  AND pr.creation_identifier ILIKE ?
			)
		`, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
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
		workSetups, err := r.loadOrderWorkSetups(ctx, data[i].ID)
		if err != nil {
			return nil, err
		}
		mapped.WorkReferences = workSetups
		result = append(result, mapped)
	}
	return result, nil
}

func (r *GORMSalesOrderRepository) FindByQuoteID(ctx context.Context, quoteID uuid.UUID) (*domain.SalesOrder, error) {
	var data SalesOrderDataModel
	if err := getDB(ctx, r.db).First(&data, "quote_id = ?", quoteID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	items, err := r.loadOrderLineItems(ctx, data.ID)
	if err != nil {
		return nil, err
	}

	order, oErr := salesOrderToDomain(&data, items)
	if oErr != nil {
		return nil, oErr
	}

	workSetups, err := r.loadOrderWorkSetups(ctx, data.ID)
	if err != nil {
		return nil, err
	}
	order.WorkReferences = workSetups

	return order, nil
}

func (r *GORMSalesOrderRepository) loadOrderLineItems(ctx context.Context, orderID uuid.UUID) ([]OrderLineItemDataModel, error) {
	var items []OrderLineItemDataModel
	if err := getDB(ctx, r.db).Where("sales_order_id = ?", orderID).Order("created_at asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GORMSalesOrderRepository) loadOrderWorkSetups(ctx context.Context, orderID uuid.UUID) ([]domain.WorkReference, error) {
	var rows []orderWorkRefRow
	err := getDB(ctx, r.db).Raw(`
		SELECT os.id, os.work_setup_id, os.work_order_id, os.sequence, os.description
		FROM order_work_setups os
		WHERE os.order_id = ? AND os.deleted_at IS NULL
		ORDER BY os.sequence ASC
	`, orderID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return orderWorkRefRowsToDomain(rows), nil
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

	return getDB(ctx, r.db).Transaction(func(tx *gorm.DB) error {
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
	if err := getDB(ctx, r.db).First(&data, "id = ?", id).Error; err != nil {
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
	query := getDB(ctx, r.db).Model(&DeliveryNoteDataModel{})
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
	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(`
			delivery_note_number ILIKE ?
			OR EXISTS (
				SELECT 1
				FROM organization_profiles op
				WHERE op.party_id = CAST(delivery_notes.party_id AS TEXT)
				  AND op.name ILIKE ?
			)
			OR EXISTS (
				SELECT 1
				FROM person_profiles pp
				WHERE pp.party_id = CAST(delivery_notes.party_id AS TEXT)
				  AND (
					pp.first_name ILIKE ?
					OR pp.last_name ILIKE ?
					OR (pp.first_name || ' ' || pp.last_name) ILIKE ?
				  )
			)
			OR EXISTS (
				SELECT 1
				FROM party_roles pr
				WHERE pr.party_id = CAST(delivery_notes.party_id AS TEXT)
				  AND pr.creation_identifier ILIKE ?
			)
		`, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
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

func (r *GORMDeliveryNoteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return getDB(ctx, r.db).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("delivery_note_id = ?", id).Delete(&DeliveryNoteLineItemDataModel{}).Error; err != nil {
			return err
		}
		return tx.Delete(&DeliveryNoteDataModel{}, "id = ?", id).Error
	})
}

func (r *GORMDeliveryNoteRepository) loadDeliveryNoteLineItems(ctx context.Context, noteID uuid.UUID) ([]DeliveryNoteLineItemDataModel, error) {
	var items []DeliveryNoteLineItemDataModel
	if err := getDB(ctx, r.db).Where("delivery_note_id = ?", noteID).Order("created_at asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GORMDeliveryNoteRepository) LinkLineItemsToInvoice(ctx context.Context, links map[uuid.UUID]uuid.UUID) error {
	if len(links) == 0 {
		return nil
	}
	db := getDB(ctx, r.db)
	for dnLineItemID, invoiceLineItemID := range links {
		ilID := invoiceLineItemID
		if err := db.Model(&DeliveryNoteLineItemDataModel{}).
			Where("id = ?", dnLineItemID).
			Update("invoice_line_item_id", ilID).Error; err != nil {
			return err
		}
	}
	return nil
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

	return getDB(ctx, r.db).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(data).Error; err != nil {
			return err
		}

		// Save existing DN-to-invoice line item links before deleting
		type dnInvoiceLink struct {
			DNLineItemID      uuid.UUID `gorm:"column:dn_line_item_id"`
			InvoiceLineItemID uuid.UUID `gorm:"column:invoice_line_item_id"`
		}
		var existingLinks []dnInvoiceLink
		if err := tx.Raw(
			"SELECT id AS dn_line_item_id, invoice_line_item_id FROM delivery_note_line_items WHERE invoice_line_item_id IN (SELECT id FROM invoice_line_items WHERE invoice_id = ?)",
			invoice.ID,
		).Scan(&existingLinks).Error; err != nil {
			return err
		}

		// Clear FK references before deleting invoice line items
		if len(existingLinks) > 0 {
			linkedIDs := make([]uuid.UUID, len(existingLinks))
			for i, l := range existingLinks {
				linkedIDs[i] = l.DNLineItemID
			}
			if err := tx.Exec(
				"UPDATE delivery_note_line_items SET invoice_line_item_id = NULL WHERE id = ANY(?)",
				pq.Array(linkedIDs),
			).Error; err != nil {
				return err
			}
		}

		if err := tx.Unscoped().Where("invoice_id = ?", invoice.ID).Delete(&InvoiceLineItemDataModel{}).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		// Restore FK references for line items that still exist
		if len(existingLinks) > 0 {
			newItemIDSet := make(map[uuid.UUID]struct{}, len(items))
			for _, item := range items {
				newItemIDSet[item.ID] = struct{}{}
			}
			for _, link := range existingLinks {
				if _, exists := newItemIDSet[link.InvoiceLineItemID]; exists {
					if err := tx.Exec(
						"UPDATE delivery_note_line_items SET invoice_line_item_id = ? WHERE id = ?",
						link.InvoiceLineItemID, link.DNLineItemID,
					).Error; err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}

func (r *GORMInvoiceRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	var data InvoiceDataModel
	if err := getDB(ctx, r.db).First(&data, "id = ?", id).Error; err != nil {
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
	query := getDB(ctx, r.db).Model(&InvoiceDataModel{})
	if filter.PartyID != nil {
		query = query.Where("party_id = ?", *filter.PartyID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.Type != nil {
		query = query.Where("type = ?", string(*filter.Type))
	}
	if filter.FromDate != nil {
		query = query.Where("invoice_date >= ?", *filter.FromDate)
	}
	if filter.ToDate != nil {
		query = query.Where("invoice_date <= ?", *filter.ToDate)
	}
	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(`
			invoice_number ILIKE ?
			OR EXISTS (
				SELECT 1
				FROM organization_profiles op
				WHERE op.party_id = CAST(invoices.party_id AS TEXT)
				  AND op.name ILIKE ?
			)
			OR EXISTS (
				SELECT 1
				FROM person_profiles pp
				WHERE pp.party_id = CAST(invoices.party_id AS TEXT)
				  AND (
					pp.first_name ILIKE ?
					OR pp.last_name ILIKE ?
					OR (pp.first_name || ' ' || pp.last_name) ILIKE ?
				  )
			)
			OR EXISTS (
				SELECT 1
				FROM party_roles pr
				WHERE pr.party_id = CAST(invoices.party_id AS TEXT)
				  AND pr.creation_identifier ILIKE ?
			)
		`, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
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
	query := getDB(ctx, r.db).
		Model(&InvoiceDataModel{}).
		Joins("JOIN invoice_line_items ON invoice_line_items.invoice_id = invoices.id").
		Joins("JOIN order_line_items ON order_line_items.id = invoice_line_items.sales_order_line_item_id").
		Where("order_line_items.sales_order_id = ?", orderID).
		Distinct("invoices.*")

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
	if err := getDB(ctx, r.db).Where("invoice_id = ?", invoiceID).Order("created_at asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GORMInvoiceRepository) FindByDeliveryNoteID(ctx context.Context, deliveryNoteID uuid.UUID) (*domain.Invoice, error) {
	var data InvoiceDataModel
	err := getDB(ctx, r.db).
		Model(&InvoiceDataModel{}).
		Joins("JOIN invoice_line_items ON invoice_line_items.invoice_id = invoices.id").
		Joins("JOIN delivery_note_line_items ON delivery_note_line_items.invoice_line_item_id = invoice_line_items.id").
		Where("delivery_note_line_items.delivery_note_id = ?", deliveryNoteID).
		First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	items, err := r.loadInvoiceLineItems(ctx, data.ID)
	if err != nil {
		return nil, err
	}
	return invoiceToDomain(&data, items)
}

func (r *GORMInvoiceRepository) ListDeliveryNoteIDsByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := getDB(ctx, r.db).
		Model(&DeliveryNoteLineItemDataModel{}).
		Joins("JOIN invoice_line_items ON invoice_line_items.id = delivery_note_line_items.invoice_line_item_id").
		Where("invoice_line_items.invoice_id = ?", invoiceID).
		Distinct("delivery_note_line_items.delivery_note_id").
		Pluck("delivery_note_line_items.delivery_note_id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *GORMInvoiceRepository) ListOrderIDsByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := getDB(ctx, r.db).
		Model(&DeliveryNoteDataModel{}).
		Select("DISTINCT delivery_notes.sales_order_id").
		Joins("JOIN delivery_note_line_items ON delivery_note_line_items.delivery_note_id = delivery_notes.id").
		Joins("JOIN invoice_line_items ON invoice_line_items.id = delivery_note_line_items.invoice_line_item_id").
		Where("invoice_line_items.invoice_id = ?", invoiceID).
		Pluck("delivery_notes.sales_order_id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}
