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

	// Get brand markup percentage
	var brand productpersistence.BrandDataModel
	var brandMarkup float64
	if err := c.db.WithContext(ctx).First(&brand, "id = ?", product.BrandID).Error; err == nil {
		brandMarkup = brand.DefaultMarkupPercentage
	}

	return &pricingapp.ProductPricingInfo{
		VariantID:             variant.ID,
		ProductID:             variant.ProductID,
		BaseCost:              variant.BaseCost,
		Currency:              "EUR",
		BrandID:               product.BrandID,
		BrandMarkupPercentage: brandMarkup,
		GroupIDs:              parseUUIDs(product.GroupIDs),
		TaxRate:               product.TaxRate,
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

	// Get brand markup percentage
	var brand productpersistence.BrandDataModel
	var brandMarkup float64
	if err := c.db.WithContext(ctx).First(&brand, "id = ?", product.BrandID).Error; err == nil {
		brandMarkup = brand.DefaultMarkupPercentage
	}

	groupIDs := parseUUIDs(product.GroupIDs)
	infos := make([]*pricingapp.ProductPricingInfo, 0, len(variants))
	for _, variant := range variants {
		infos = append(infos, &pricingapp.ProductPricingInfo{
			VariantID:             variant.ID,
			ProductID:             variant.ProductID,
			BaseCost:              variant.BaseCost,
			Currency:              "EUR",
			BrandID:               product.BrandID,
			BrandMarkupPercentage: brandMarkup,
			GroupIDs:              groupIDs,
			TaxRate:               product.TaxRate,
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
