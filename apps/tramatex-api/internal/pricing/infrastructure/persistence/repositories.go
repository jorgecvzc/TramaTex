package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type GORMPricingRuleRepository struct {
	db *gorm.DB
}

func NewGORMPricingRuleRepository(db *gorm.DB) *GORMPricingRuleRepository {
	return &GORMPricingRuleRepository{db: db}
}

func (r *GORMPricingRuleRepository) Save(ctx context.Context, rule *domain.PricingRule) error {
	data := PricingRuleFromDomain(rule)
	return r.db.WithContext(ctx).Save(data).Error
}

func (r *GORMPricingRuleRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.PricingRule, error) {
	var data PricingRuleDataModel
	err := r.db.WithContext(ctx).First(&data, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return data.ToDomain()
}

func (r *GORMPricingRuleRepository) List(ctx context.Context) ([]*domain.PricingRule, error) {
	var data []PricingRuleDataModel
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&data).Error
	if err != nil {
		return nil, err
	}

	rules := make([]*domain.PricingRule, 0, len(data))
	for i := range data {
		rule, err := data[i].ToDomain()
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (r *GORMPricingRuleRepository) FindApplicable(ctx context.Context, variantID uuid.UUID, quantity int, at time.Time) ([]*domain.PricingRule, error) {
	var data []PricingRuleDataModel
	query := r.db.WithContext(ctx).
		Where("is_active = true").
		Where("effective_from <= ?", at).
		Where("effective_to IS NULL OR effective_to >= ?", at).
		Where("(product_variant_id IS NULL OR product_variant_id = ?)", variantID).
		Where("min_quantity <= ?", quantity).
		Where("max_quantity IS NULL OR max_quantity >= ?", quantity)

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	rules := make([]*domain.PricingRule, 0, len(data))
	for i := range data {
		rule, err := data[i].ToDomain()
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

type GORMClientPricingRepository struct {
	db *gorm.DB
}

func NewGORMClientPricingRepository(db *gorm.DB) *GORMClientPricingRepository {
	return &GORMClientPricingRepository{db: db}
}

func (r *GORMClientPricingRepository) Save(ctx context.Context, override *domain.ClientPricing) error {
	data := ClientPricingFromDomain(override)
	return r.db.WithContext(ctx).Save(data).Error
}

func (r *GORMClientPricingRepository) FindApplicable(ctx context.Context, clientID uuid.UUID, variantID uuid.UUID, at time.Time) (*domain.ClientPricing, error) {
	var data ClientPricingDataModel
	err := r.db.WithContext(ctx).
		Where("is_active = true").
		Where("client_id = ? AND product_variant_id = ?", clientID, variantID).
		Where("effective_from <= ?", at).
		Where("effective_to IS NULL OR effective_to >= ?", at).
		Order("created_at desc").
		First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return data.ToDomain()
}

type GORMBrandProfitMarginRepository struct {
	db *gorm.DB
}

func NewGORMBrandProfitMarginRepository(db *gorm.DB) *GORMBrandProfitMarginRepository {
	return &GORMBrandProfitMarginRepository{db: db}
}

func (r *GORMBrandProfitMarginRepository) Save(ctx context.Context, margin *domain.BrandProfitMargin) error {
	data := BrandProfitMarginFromDomain(margin)
	return r.db.WithContext(ctx).Save(data).Error
}

func (r *GORMBrandProfitMarginRepository) FindApplicable(ctx context.Context, brandID uuid.UUID, at time.Time) (*domain.BrandProfitMargin, error) {
	var data BrandProfitMarginDataModel
	err := r.db.WithContext(ctx).
		Where("is_active = true").
		Where("brand_id = ?", brandID).
		Where("effective_from <= ?", at).
		Where("effective_to IS NULL OR effective_to >= ?", at).
		Order("created_at desc").
		First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return data.ToDomain()
}

type GORMSalesDiscountRuleRepository struct {
	db *gorm.DB
}

func NewGORMSalesDiscountRuleRepository(db *gorm.DB) *GORMSalesDiscountRuleRepository {
	return &GORMSalesDiscountRuleRepository{db: db}
}

func (r *GORMSalesDiscountRuleRepository) Save(ctx context.Context, rule *domain.SalesDiscountRule) error {
	data := SalesDiscountRuleFromDomain(rule)
	return r.db.WithContext(ctx).Save(data).Error
}

func (r *GORMSalesDiscountRuleRepository) FindApplicable(ctx context.Context, clientID uuid.UUID, variantID uuid.UUID, quantity int, at time.Time) ([]*domain.SalesDiscountRule, error) {
	var data []SalesDiscountRuleDataModel
	query := r.db.WithContext(ctx).
		Where("is_active = true").
		Where("effective_from <= ?", at).
		Where("effective_to IS NULL OR effective_to >= ?", at).
		Where("(client_id IS NULL OR client_id = ?)", clientID).
		Where("(product_variant_id IS NULL OR product_variant_id = ?)", variantID).
		Where("min_quantity IS NULL OR min_quantity <= ?", quantity)

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	rules := make([]*domain.SalesDiscountRule, 0, len(data))
	for i := range data {
		rule, err := data[i].ToDomain()
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

type GORMPriceCalculationRepository struct {
	db *gorm.DB
}

func NewGORMPriceCalculationRepository(db *gorm.DB) *GORMPriceCalculationRepository {
	return &GORMPriceCalculationRepository{db: db}
}

func (r *GORMPriceCalculationRepository) Save(ctx context.Context, calc *domain.PriceCalculation) error {
	data, err := PriceCalculationFromDomain(calc)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(data).Error
}

func (r *GORMPriceCalculationRepository) ListByProductVariantID(ctx context.Context, variantID uuid.UUID) ([]*domain.PriceCalculation, error) {
	var data []PriceCalculationDataModel
	if err := r.db.WithContext(ctx).
		Where("product_variant_id = ?", variantID).
		Order("calculated_at desc").
		Find(&data).Error; err != nil {
		return nil, err
	}

	results := make([]*domain.PriceCalculation, 0, len(data))
	for i := range data {
		calc, err := data[i].ToDomain()
		if err != nil {
			return nil, err
		}
		results = append(results, calc)
	}
	return results, nil
}

type GORMBaseSalesPriceRuleRepository struct {
	db *gorm.DB
}

func NewGORMBaseSalesPriceRuleRepository(db *gorm.DB) *GORMBaseSalesPriceRuleRepository {
	return &GORMBaseSalesPriceRuleRepository{db: db}
}

func (r *GORMBaseSalesPriceRuleRepository) Save(ctx context.Context, rule *domain.BaseSalesPriceRule) error {
	data := BaseSalesPriceRuleFromDomain(rule)
	return r.db.WithContext(ctx).Save(data).Error
}

func (r *GORMBaseSalesPriceRuleRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.BaseSalesPriceRule, error) {
	var data BaseSalesPriceRuleDataModel
	err := r.db.WithContext(ctx).First(&data, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return data.ToDomain()
}

func (r *GORMBaseSalesPriceRuleRepository) List(ctx context.Context) ([]*domain.BaseSalesPriceRule, error) {
	var data []BaseSalesPriceRuleDataModel
	if err := r.db.WithContext(ctx).Order("created_at desc").Find(&data).Error; err != nil {
		return nil, err
	}

	rules := make([]*domain.BaseSalesPriceRule, 0, len(data))
	for i := range data {
		rule, err := data[i].ToDomain()
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

type GORMSaleModificationRuleRepository struct {
	db *gorm.DB
}

func NewGORMSaleModificationRuleRepository(db *gorm.DB) *GORMSaleModificationRuleRepository {
	return &GORMSaleModificationRuleRepository{db: db}
}

func (r *GORMSaleModificationRuleRepository) Save(ctx context.Context, rule *domain.SaleModificationRule) error {
	data := SaleModificationRuleFromDomain(rule)
	return r.db.WithContext(ctx).Save(data).Error
}

func (r *GORMSaleModificationRuleRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.SaleModificationRule, error) {
	var data SaleModificationRuleDataModel
	err := r.db.WithContext(ctx).First(&data, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return data.ToDomain()
}

func (r *GORMSaleModificationRuleRepository) ListApplicable(ctx context.Context, clientID string, productGroupID *uuid.UUID, orderTotal domain.Money, at time.Time) ([]*domain.SaleModificationRule, error) {
	var data []SaleModificationRuleDataModel
	query := r.db.WithContext(ctx).
		Where("is_active = true").
		Where("effective_from <= ?", at).
		Where("effective_to IS NULL OR effective_to >= ?", at)

	if productGroupID != nil {
		query = query.Where("product_group_id IS NULL OR product_group_id = ?", *productGroupID)
	}

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	rules := make([]*domain.SaleModificationRule, 0, len(data))
	for i := range data {
		rule, err := data[i].ToDomain()
		if err != nil {
			return nil, err
		}
		if rule.AppliesTo(clientID, productGroupID, orderTotal, at) {
			rules = append(rules, rule)
		}
	}
	return rules, nil
}
