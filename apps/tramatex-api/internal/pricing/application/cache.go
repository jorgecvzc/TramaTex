package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type BasePriceCache interface {
	GetBasePrice(ctx context.Context, productID uuid.UUID, variantID uuid.UUID) (*domain.Money, error)
	SetBasePrice(ctx context.Context, productID uuid.UUID, variantID uuid.UUID, price domain.Money) error
}
