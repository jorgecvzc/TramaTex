package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type QuoteFilter struct {
	PartyID  *uuid.UUID
	Status   *QuoteStatus
	FromDate *time.Time
	ToDate   *time.Time
	Search   *string
	Limit    int
}

type SalesOrderFilter struct {
	PartyID  *uuid.UUID
	Status   *SalesOrderStatus
	FromDate *time.Time
	ToDate   *time.Time
	Search   *string
	Limit    int
}

type DeliveryNoteFilter struct {
	SalesOrderID *uuid.UUID
	PartyID      *uuid.UUID
	Status       *DeliveryNoteStatus
	FromDate     *time.Time
	ToDate       *time.Time
	Search       *string
	Limit        int
}

type InvoiceFilter struct {
	PartyID  *uuid.UUID
	Status   *InvoiceStatus
	Type     *InvoiceType // Filter by COMPLETA or SIMPLIFICADA
	FromDate *time.Time
	ToDate   *time.Time
	Search   *string
	Limit    int
}

type QuoteRepository interface {
	Save(ctx context.Context, quote *Quote) error
	FindByID(ctx context.Context, id uuid.UUID) (*Quote, error)
	List(ctx context.Context, filter QuoteFilter) ([]*Quote, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type SalesOrderRepository interface {
	Save(ctx context.Context, order *SalesOrder) error
	FindByID(ctx context.Context, id uuid.UUID) (*SalesOrder, error)
	FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*SalesOrder, error)
	FindByQuoteID(ctx context.Context, quoteID uuid.UUID) (*SalesOrder, error)
	List(ctx context.Context, filter SalesOrderFilter) ([]*SalesOrder, error)
}

type DeliveryNoteRepository interface {
	Save(ctx context.Context, note *DeliveryNote) error
	FindByID(ctx context.Context, id uuid.UUID) (*DeliveryNote, error)
	List(ctx context.Context, filter DeliveryNoteFilter) ([]*DeliveryNote, error)
	ListBySalesOrderID(ctx context.Context, orderID uuid.UUID) ([]*DeliveryNote, error)
	LinkLineItemsToInvoice(ctx context.Context, links map[uuid.UUID]uuid.UUID) error
}

type InvoiceRepository interface {
	Save(ctx context.Context, invoice *Invoice) error
	FindByID(ctx context.Context, id uuid.UUID) (*Invoice, error)
	List(ctx context.Context, filter InvoiceFilter) ([]*Invoice, error)
	ListBySalesOrderID(ctx context.Context, orderID uuid.UUID) ([]*Invoice, error)
	FindByDeliveryNoteID(ctx context.Context, deliveryNoteID uuid.UUID) (*Invoice, error)
	ListDeliveryNoteIDsByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]uuid.UUID, error)
}
