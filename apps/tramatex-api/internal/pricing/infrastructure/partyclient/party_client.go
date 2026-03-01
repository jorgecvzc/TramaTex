package partyclient

import (
	"context"

	"gorm.io/gorm"
)

// PartyPricingClient provides client/party pricing info to the pricing module.
// This is an anti-corruption layer: pricing never depends on party domain directly.
type PartyPricingClient struct {
	db *gorm.DB
}

// partyRow is a minimal read model used only by the pricing module
type partyRow struct {
	DefaultDiscountPercentage float64 `gorm:"column:default_discount_percentage"`
}

func NewPartyPricingClient(db *gorm.DB) *PartyPricingClient {
	return &PartyPricingClient{db: db}
}

// GetClientDefaultDiscount returns the default discount percentage for a client party.
// Returns 0 if the party is not found or has no discount.
func (c *PartyPricingClient) GetClientDefaultDiscount(ctx context.Context, clientID string) (float64, error) {
	if clientID == "" {
		return 0, nil
	}

	var row partyRow
	if err := c.db.WithContext(ctx).
		Table("parties").
		Select("default_discount_percentage").
		Where("id = ?", clientID).
		First(&row).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}

	return row.DefaultDiscountPercentage, nil
}
