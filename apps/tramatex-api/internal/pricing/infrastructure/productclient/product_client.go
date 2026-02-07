package productclient

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	pricingapp "github.com/joran-cortez/tramatex/internal/pricing/application"
	productpersistence "github.com/joran-cortez/tramatex/internal/product/infrastructure/persistence"
)

type ProductPricingClient struct {
	db *gorm.DB
}

func NewProductPricingClient(db *gorm.DB) *ProductPricingClient {
	return &ProductPricingClient{db: db}
}

func (c *ProductPricingClient) GetVariantPricingInfo(ctx context.Context, variantID uuid.UUID) (*pricingapp.ProductPricingInfo, error) {
	var variant productpersistence.VariantDataModel
	if err := c.db.WithContext(ctx).First(&variant, "id = ?", variantID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var product productpersistence.ProductDataModel
	if err := c.db.WithContext(ctx).First(&product, "id = ?", variant.ProductID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &pricingapp.ProductPricingInfo{
		VariantID: variant.ID,
		ProductID: variant.ProductID,
		BaseCost:  variant.BaseCost,
		Currency:  "EUR",
		BrandID:   product.BrandID,
		GroupIDs:  parseUUIDs(product.GroupIDs),
	}, nil
}

func (c *ProductPricingClient) ListVariantsPricingInfo(ctx context.Context, productID uuid.UUID) ([]*pricingapp.ProductPricingInfo, error) {
	var product productpersistence.ProductDataModel
	if err := c.db.WithContext(ctx).First(&product, "id = ?", productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var variants []productpersistence.VariantDataModel
	if err := c.db.WithContext(ctx).Where("product_id = ?", productID).Find(&variants).Error; err != nil {
		return nil, err
	}

	groupIDs := parseUUIDs(product.GroupIDs)
	infos := make([]*pricingapp.ProductPricingInfo, 0, len(variants))
	for _, variant := range variants {
		infos = append(infos, &pricingapp.ProductPricingInfo{
			VariantID: variant.ID,
			ProductID: variant.ProductID,
			BaseCost:  variant.BaseCost,
			Currency:  "EUR",
			BrandID:   product.BrandID,
			GroupIDs:  groupIDs,
		})
	}

	return infos, nil
}

func parseUUIDs(values []string) []uuid.UUID {
	result := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		parsed, err := uuid.Parse(value)
		if err != nil {
			continue
		}
		result = append(result, parsed)
	}
	return result
}
