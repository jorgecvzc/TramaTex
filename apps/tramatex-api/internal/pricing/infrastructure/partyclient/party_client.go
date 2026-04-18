package partyclient

import (
	"context"

	partyapp "github.com/joran-cortez/tramatex/internal/party/application"
)

// PartyPricingClient provides client/party pricing info to the pricing module.
// This is an anti-corruption layer: pricing consumes Party's application layer,
// never its domain or persistence directly.
type PartyPricingClient struct {
	handler *partyapp.GetClientDefaultDiscountHandler
}

func NewPartyPricingClient(handler *partyapp.GetClientDefaultDiscountHandler) *PartyPricingClient {
	return &PartyPricingClient{handler: handler}
}

// GetClientDefaultDiscount returns the default discount percentage for a client party.
// Returns 0 if the party is not found or has no discount.
func (c *PartyPricingClient) GetClientDefaultDiscount(ctx context.Context, clientID string) (float64, error) {
	return c.handler.Handle(ctx, &partyapp.GetClientDefaultDiscountQuery{PartyID: clientID})
}
