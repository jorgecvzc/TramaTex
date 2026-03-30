package productclient

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	pricingapp "github.com/joran-cortez/tramatex/internal/pricing/application"
	productdomain "github.com/joran-cortez/tramatex/internal/product/domain"
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

	// Calculate BaseCost dynamically from product BasePrice + attribute modifiers
	baseCost, err := c.calculateVariantBaseCost(ctx, product.BasePrice, variant.AttributeValues)
	if err != nil {
		// Log the error, but use BasePrice as fallback
		baseCost = product.BasePrice
	}

	// Get brand markup percentage
	var brand productpersistence.BrandDataModel
	var brandMarkup float64
	brandID := uuid.Nil
	if product.BrandID != nil {
		brandID = *product.BrandID
		if err := c.db.WithContext(ctx).First(&brand, "id = ?", brandID).Error; err == nil {
			brandMarkup = brand.DefaultMarkupPercentage
		}
	}

	return &pricingapp.ProductPricingInfo{
		VariantID:             variant.ID,
		ProductID:             variant.ProductID,
		BaseCost:              baseCost,
		Currency:              "EUR",
		BrandID:               brandID,
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
	brandID := uuid.Nil
	if product.BrandID != nil {
		brandID = *product.BrandID
		if err := c.db.WithContext(ctx).First(&brand, "id = ?", brandID).Error; err == nil {
			brandMarkup = brand.DefaultMarkupPercentage
		}
	}

	groupIDs := parseUUIDs(product.GroupIDs)
	infos := make([]*pricingapp.ProductPricingInfo, 0, len(variants))
	for _, variant := range variants {
		// Calculate BaseCost dynamically for each variant
		baseCost, err := c.calculateVariantBaseCost(ctx, product.BasePrice, variant.AttributeValues)
		if err != nil {
			// Log error, use product BasePrice as fallback
			baseCost = product.BasePrice
		}

		infos = append(infos, &pricingapp.ProductPricingInfo{
			VariantID:             variant.ID,
			ProductID:             variant.ProductID,
			BaseCost:              baseCost,
			Currency:              "EUR",
			BrandID:               brandID,
			BrandMarkupPercentage: brandMarkup,
			GroupIDs:              groupIDs,
			TaxRate:               product.TaxRate,
		})
	}

	return infos, nil
}

// calculateVariantBaseCost calculates the base cost for a variant by loading attribute values
// and applying their price modifiers to the product's base price.
func (c *ProductPricingClient) calculateVariantBaseCost(ctx context.Context, productBasePrice float64, attributeValueIDStrings []string) (float64, error) {
	if len(attributeValueIDStrings) == 0 {
		return productBasePrice, nil
	}

	// Parse attribute value IDs
	attributeValueIDs := make([]uuid.UUID, 0, len(attributeValueIDStrings))
	for _, idStr := range attributeValueIDStrings {
		id, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}
		attributeValueIDs = append(attributeValueIDs, id)
	}

	if len(attributeValueIDs) == 0 {
		return productBasePrice, nil
	}

	// Load attribute values from database
	var attrValueDataModels []productpersistence.AttributeValueDataModel
	if err := c.db.WithContext(ctx).Where("id IN ?", attributeValueIDs).Find(&attrValueDataModels).Error; err != nil {
		return 0, err
	}

	// Convert to domain models
	attrValues := make([]productdomain.AttributeValue, 0, len(attrValueDataModels))
	for _, dm := range attrValueDataModels {
		attrValues = append(attrValues, *dm.ToDomain())
	}

	// Calculate base cost using domain logic
	return productdomain.CalculateBaseCost(productBasePrice, attrValues), nil
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
