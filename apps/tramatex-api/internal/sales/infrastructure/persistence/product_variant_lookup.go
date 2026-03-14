package persistence

import (
	"context"

	"github.com/google/uuid"

	product_domain "github.com/joran-cortez/tramatex/internal/product/domain"
	sales_app "github.com/joran-cortez/tramatex/internal/sales/application"
)

type ProductVariantLookupAdapter struct {
	variantRepo   product_domain.ProductVariantRepository
	productRepo   product_domain.ProductRepository
	attributeRepo product_domain.AttributeRepository
}

func NewProductVariantLookupAdapter(
	variantRepo product_domain.ProductVariantRepository,
	productRepo product_domain.ProductRepository,
	attributeRepo product_domain.AttributeRepository,
) *ProductVariantLookupAdapter {
	return &ProductVariantLookupAdapter{
		variantRepo:   variantRepo,
		productRepo:   productRepo,
		attributeRepo: attributeRepo,
	}
}

func (a *ProductVariantLookupAdapter) GetVariantInfo(ctx context.Context, variantID uuid.UUID) (*sales_app.VariantInfo, error) {
	variant, err := a.variantRepo.FindByID(ctx, variantID)
	if err != nil || variant == nil {
		return nil, err
	}

	product, err := a.productRepo.FindByID(ctx, variant.ProductID)
	if err != nil || product == nil {
		return &sales_app.VariantInfo{VariantSKU: variant.SKU}, err
	}

	optionConfig := a.buildOptionConfiguration(ctx, variant)

	return &sales_app.VariantInfo{
		ProductName:         product.Name,
		VariantSKU:          variant.SKU,
		OptionConfiguration: optionConfig,
	}, nil
}

func (a *ProductVariantLookupAdapter) buildOptionConfiguration(ctx context.Context, variant *product_domain.ProductVariant) map[string]string {
	if len(variant.AttributeValues) == 0 || a.attributeRepo == nil {
		return nil
	}

	allAttributes, err := a.attributeRepo.FindByScope(ctx, nil, nil)
	if err != nil || len(allAttributes) == 0 {
		return nil
	}

	optionConfig := make(map[string]string)
	for _, avID := range variant.AttributeValues {
		for _, attr := range allAttributes {
			for _, val := range attr.Values {
				if val.ID == avID {
					optionConfig[attr.Name] = val.Value
				}
			}
		}
	}

	if len(optionConfig) == 0 {
		return nil
	}
	return optionConfig
}
