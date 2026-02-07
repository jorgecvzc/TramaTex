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
}

type SalesOrderFilter struct {
	PartyID  *uuid.UUID
	Status   *SalesOrderStatus
	FromDate *time.Time
	ToDate   *time.Time
}

type DeliveryNoteFilter struct {
	SalesOrderID *uuid.UUID
	PartyID      *uuid.UUID
	Status       *DeliveryNoteStatus
	FromDate     *time.Time
	ToDate       *time.Time
}

type InvoiceFilter struct {
	PartyID  *uuid.UUID
	Status   *InvoiceStatus
	FromDate *time.Time
	ToDate   *time.Time
}

type QuoteRepository interface {
	Save(ctx context.Context, quote *Quote) error
	FindByID(ctx context.Context, id uuid.UUID) (*Quote, error)
	List(ctx context.Context, filter QuoteFilter) ([]*Quote, error)
}

type SalesOrderRepository interface {
	Save(ctx context.Context, order *SalesOrder) error
	FindByID(ctx context.Context, id uuid.UUID) (*SalesOrder, error)
	List(ctx context.Context, filter SalesOrderFilter) ([]*SalesOrder, error)
}

type DeliveryNoteRepository interface {
	Save(ctx context.Context, note *DeliveryNote) error
	FindByID(ctx context.Context, id uuid.UUID) (*DeliveryNote, error)
	List(ctx context.Context, filter DeliveryNoteFilter) ([]*DeliveryNote, error)
	ListBySalesOrderID(ctx context.Context, orderID uuid.UUID) ([]*DeliveryNote, error)
}

type InvoiceRepository interface {
	Save(ctx context.Context, invoice *Invoice) error
	FindByID(ctx context.Context, id uuid.UUID) (*Invoice, error)
	List(ctx context.Context, filter InvoiceFilter) ([]*Invoice, error)
	ListBySalesOrderID(ctx context.Context, orderID uuid.UUID) ([]*Invoice, error)
}
