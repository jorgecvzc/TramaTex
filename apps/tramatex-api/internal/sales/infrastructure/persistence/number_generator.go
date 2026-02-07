package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/joran-cortez/tramatex/internal/sales/domain"
)

type TimeBasedNumberGenerator struct{}

func NewTimeBasedNumberGenerator() *TimeBasedNumberGenerator {
	return &TimeBasedNumberGenerator{}
}

func (g *TimeBasedNumberGenerator) NextQuoteNumber(ctx context.Context) (domain.QuoteNumber, error) {
	return domain.NewQuoteNumber(buildNumber("Q"))
}

func (g *TimeBasedNumberGenerator) NextOrderNumber(ctx context.Context) (domain.OrderNumber, error) {
	return domain.NewOrderNumber(buildNumber("SO"))
}

func (g *TimeBasedNumberGenerator) NextDeliveryNoteNumber(ctx context.Context) (domain.DeliveryNoteNumber, error) {
	return domain.NewDeliveryNoteNumber(buildNumber("DN"))
}

func (g *TimeBasedNumberGenerator) NextInvoiceNumber(ctx context.Context) (domain.InvoiceNumber, error) {
	return domain.NewInvoiceNumber(buildNumber("INV"))
}

func buildNumber(prefix string) string {
	timestamp := time.Now().UTC().Format("20060102-150405")
	rand := uuid.New().String()[:8]
	return fmt.Sprintf("%s-%s-%s", prefix, timestamp, rand)
}
