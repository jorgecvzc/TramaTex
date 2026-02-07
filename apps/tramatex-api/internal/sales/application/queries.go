package application

import (
	"time"

	"github.com/google/uuid"
)

type GetQuoteByIDQuery struct {
	ID uuid.UUID
}

type ListQuotesQuery struct {
	PartyID    *uuid.UUID
	Status     *string
	FromDate   *time.Time
	ToDate     *time.Time
	PageSize   int
	PageNumber int
}

type GetOrderByIDQuery struct {
	ID uuid.UUID
}

type ListOrdersQuery struct {
	PartyID    *uuid.UUID
	Status     *string
	FromDate   *time.Time
	ToDate     *time.Time
	PageSize   int
	PageNumber int
}

type GetDeliveryNoteByIDQuery struct {
	ID uuid.UUID
}

type ListDeliveryNotesQuery struct {
	SalesOrderID *uuid.UUID
	PartyID      *uuid.UUID
	Status       *string
	FromDate     *time.Time
	ToDate       *time.Time
	PageSize     int
	PageNumber   int
}

type GetInvoiceByIDQuery struct {
	ID uuid.UUID
}

type ListInvoicesQuery struct {
	PartyID    *uuid.UUID
	Status     *string
	FromDate   *time.Time
	ToDate     *time.Time
	PageSize   int
	PageNumber int
}
