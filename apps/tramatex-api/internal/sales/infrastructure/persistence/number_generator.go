package persistence

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/sales/domain"
)

// DocumentSequence is the GORM model for the document_sequences table.
type DocumentSequence struct {
	Prefix       string `gorm:"column:prefix;primaryKey"`
	Year         int    `gorm:"column:year;primaryKey"`
	CurrentValue int    `gorm:"column:current_value;not null;default:0"`
}

func (DocumentSequence) TableName() string { return "document_sequences" }

// SequentialNumberGenerator produces sequential document numbers using the DB.
// Format: PREFIX-YEAR-NNNN (e.g., PRE-2026-0001).
type SequentialNumberGenerator struct {
	db *gorm.DB
}

func NewSequentialNumberGenerator(db *gorm.DB) *SequentialNumberGenerator {
	return &SequentialNumberGenerator{db: db}
}

// nextSequential atomically increments the counter for the given prefix+year
// and returns the formatted number.  Uses INSERT ... ON CONFLICT to initialise
// the row when it doesn't exist yet (first document of the year).
func (g *SequentialNumberGenerator) nextSequential(ctx context.Context, prefix string, year int) (string, error) {
	var seq DocumentSequence

	err := getDB(ctx, g.db).Raw(`
		INSERT INTO document_sequences (prefix, year, current_value)
		VALUES (?, ?, 1)
		ON CONFLICT (prefix, year)
		DO UPDATE SET current_value = document_sequences.current_value + 1
		RETURNING current_value
	`, prefix, year).Scan(&seq.CurrentValue).Error
	if err != nil {
		return "", fmt.Errorf("failed to generate sequential number for %s/%d: %w", prefix, year, err)
	}

	return fmt.Sprintf("%s-%d-%04d", prefix, year, seq.CurrentValue), nil
}

func (g *SequentialNumberGenerator) NextQuoteNumber(ctx context.Context) (domain.QuoteNumber, error) {
	num, err := g.nextSequential(ctx, "PRE", time.Now().Year())
	if err != nil {
		return domain.QuoteNumber{}, err
	}
	return domain.NewQuoteNumber(num)
}

func (g *SequentialNumberGenerator) NextOrderNumber(ctx context.Context) (domain.OrderNumber, error) {
	num, err := g.nextSequential(ctx, "PED", time.Now().Year())
	if err != nil {
		return domain.OrderNumber{}, err
	}
	return domain.NewOrderNumber(num)
}

func (g *SequentialNumberGenerator) NextDeliveryNoteNumber(ctx context.Context) (domain.DeliveryNoteNumber, error) {
	num, err := g.nextSequential(ctx, "ALB", time.Now().Year())
	if err != nil {
		return domain.DeliveryNoteNumber{}, err
	}
	return domain.NewDeliveryNoteNumber(num)
}

func (g *SequentialNumberGenerator) NextInvoiceNumber(ctx context.Context, series domain.InvoiceSeries) (domain.InvoiceNumber, error) {
	num, err := g.nextSequential(ctx, series.Prefix(), series.Year())
	if err != nil {
		return domain.InvoiceNumber{}, err
	}
	return domain.NewInvoiceNumber(num)
}
