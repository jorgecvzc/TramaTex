package productclient

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	pricingapp "github.com/joran-cortez/tramatex/internal/pricing/application"
	productapp "github.com/joran-cortez/tramatex/internal/product/application"
)

// ProductDataProvider abstracts the Product module's application layer.
// Satisfied by *productapp.ProductService.
type ProductDataProvider interface {
	GetVariantPricingData(ctx context.Context, variantID uuid.UUID) (*productapp.VariantPricingDataDTO, error)
	GetVariantsPricingData(ctx context.Context, variantIDs []uuid.UUID) ([]*productapp.VariantPricingDataDTO, error)
	ListVariantsPricingData(ctx context.Context, productID uuid.UUID) ([]*productapp.VariantPricingDataDTO, error)
}

type ProductPricingClient struct {
	provider ProductDataProvider
}

func NewProductPricingClient(provider ProductDataProvider) *ProductPricingClient {
	return &ProductPricingClient{provider: provider}
}

func (c *ProductPricingClient) GetVariantPricingInfo(ctx context.Context, variantID uuid.UUID) (*pricingapp.ProductPricingInfo, error) {
	data, err := c.provider.GetVariantPricingData(ctx, variantID)
	if err != nil || data == nil {
		return nil, err
	}
	return mapToProductPricingInfo(data), nil
}

func (c *ProductPricingClient) GetVariantsPricingInfo(ctx context.Context, variantIDs []uuid.UUID) ([]*pricingapp.ProductPricingInfo, error) {
	dataList, err := c.provider.GetVariantsPricingData(ctx, variantIDs)
	if err != nil || dataList == nil {
		return nil, err
	}
	result := make([]*pricingapp.ProductPricingInfo, 0, len(dataList))
	for _, data := range dataList {
		result = append(result, mapToProductPricingInfo(data))
	}
	return result, nil
}

func (c *ProductPricingClient) ListVariantsPricingInfo(ctx context.Context, productID uuid.UUID) ([]*pricingapp.ProductPricingInfo, error) {
	dataList, err := c.provider.ListVariantsPricingData(ctx, productID)
	if err != nil || dataList == nil {
		return nil, err
	}
	result := make([]*pricingapp.ProductPricingInfo, 0, len(dataList))
	for _, data := range dataList {
		result = append(result, mapToProductPricingInfo(data))
	}
	return result, nil
}

func mapToProductPricingInfo(data *productapp.VariantPricingDataDTO) *pricingapp.ProductPricingInfo {
	return &pricingapp.ProductPricingInfo{
		VariantID:             data.VariantID,
		ProductID:             data.ProductID,
		BaseCost:              decimal.NewFromFloat(data.BaseCost),
		Currency:              data.Currency,
		BrandID:               data.BrandID,
		BrandMarkupPercentage: decimal.NewFromFloat(data.BrandMarkupPercentage),
		GroupIDs:              data.GroupIDs,
		TaxRate:               decimal.NewFromFloat(data.TaxRate),
	}
}
