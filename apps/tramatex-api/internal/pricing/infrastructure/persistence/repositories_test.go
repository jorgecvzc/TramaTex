package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		sqlDB.Close()
		t.Fatalf("failed to open gorm db: %v", err)
	}
	cleanup := func() { _ = sqlDB.Close() }
	return gormDB, mock, cleanup
}

func TestGORMPricingRuleRepository(t *testing.T) {
	gormDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewGORMPricingRuleRepository(gormDB)

	percentage, _ := domain.NewPercentage(0.1)
	rule, _ := domain.NewPricingRule("Rule", nil, nil, percentage, 1, nil, time.Now(), nil)

	mock.ExpectExec("UPDATE .*\"pricing_rules\"").WillReturnResult(sqlmock.NewResult(1, 1))
	if err := repo.Save(context.Background(), rule); err != nil {
		t.Fatalf("expected save success, got %v", err)
	}

	id := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "name", "markup_percentage", "min_quantity", "effective_from", "is_active"}).
		AddRow(id, "Rule", 0.1, 1, time.Now(), true)
	mock.ExpectQuery("SELECT .* FROM \"pricing_rules\"").WillReturnRows(rows)
	found, err := repo.FindByID(context.Background(), id)
	if err != nil || found == nil {
		t.Fatalf("expected find success")
	}

	listRows := sqlmock.NewRows([]string{"id", "name", "markup_percentage", "min_quantity", "effective_from", "is_active"}).
		AddRow(uuid.New(), "Rule A", 0.1, 1, time.Now(), true).
		AddRow(uuid.New(), "Rule B", 0.2, 2, time.Now(), true)
	mock.ExpectQuery("SELECT .* FROM \"pricing_rules\"").WillReturnRows(listRows)
	list, err := repo.List(context.Background())
	if err != nil || len(list) != 2 {
		t.Fatalf("expected list success")
	}

	applicableRows := sqlmock.NewRows([]string{"id", "name", "markup_percentage", "min_quantity", "effective_from", "is_active"}).
		AddRow(uuid.New(), "Rule", 0.1, 1, time.Now(), true)
	mock.ExpectQuery("SELECT .* FROM \"pricing_rules\"").WillReturnRows(applicableRows)
	foundRules, err := repo.FindApplicable(context.Background(), uuid.New(), 1, time.Now())
	if err != nil || len(foundRules) != 1 {
		t.Fatalf("expected applicable rules")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGORMClientPricingRepository(t *testing.T) {
	gormDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewGORMClientPricingRepository(gormDB)

	money, _ := domain.NewMoney(10, "EUR")
	override, _ := domain.NewClientPricing(uuid.New(), uuid.New(), money, time.Now(), nil)

	mock.ExpectExec("UPDATE .*\"client_pricing_overrides\"").WillReturnResult(sqlmock.NewResult(1, 1))
	if err := repo.Save(context.Background(), override); err != nil {
		t.Fatalf("expected save success, got %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "client_id", "product_variant_id", "fixed_price", "currency", "effective_from", "is_active"}).
		AddRow(uuid.New(), override.ClientID, override.ProductVariantID, 10.0, "EUR", time.Now(), true)
	mock.ExpectQuery("SELECT .* FROM \"client_pricing_overrides\"").WillReturnRows(rows)
	found, err := repo.FindApplicable(context.Background(), override.ClientID, override.ProductVariantID, time.Now())
	if err != nil || found == nil {
		t.Fatalf("expected find success")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGORMBrandProfitMarginRepository(t *testing.T) {
	gormDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewGORMBrandProfitMarginRepository(gormDB)

	percentage, _ := domain.NewPercentage(0.1)
	margin, _ := domain.NewBrandProfitMargin(uuid.New(), &percentage, nil, time.Now(), nil)

	mock.ExpectExec("UPDATE .*\"brand_profit_margins\"").WillReturnResult(sqlmock.NewResult(1, 1))
	if err := repo.Save(context.Background(), margin); err != nil {
		t.Fatalf("expected save success, got %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "brand_id", "percentage_value", "currency", "effective_from", "is_active"}).
		AddRow(uuid.New(), margin.BrandID, 0.1, "EUR", time.Now(), true)
	mock.ExpectQuery("SELECT .* FROM \"brand_profit_margins\"").WillReturnRows(rows)
	found, err := repo.FindApplicable(context.Background(), margin.BrandID, time.Now())
	if err != nil || found == nil {
		t.Fatalf("expected find success")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGORMSalesDiscountRuleRepository(t *testing.T) {
	gormDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewGORMSalesDiscountRuleRepository(gormDB)

	percentage, _ := domain.NewPercentage(0.1)
	rule, _ := domain.NewSalesDiscountRule("Discount", nil, nil, nil, domain.DiscountTypePercentage, &percentage, nil, 1, time.Now(), nil)

	mock.ExpectExec("UPDATE .*\"sales_discount_rules\"").WillReturnResult(sqlmock.NewResult(1, 1))
	if err := repo.Save(context.Background(), rule); err != nil {
		t.Fatalf("expected save success, got %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "discount_type", "percentage_value", "currency", "priority", "effective_from", "is_active"}).
		AddRow(uuid.New(), "Discount", string(domain.DiscountTypePercentage), 0.1, "EUR", 1, time.Now(), true)
	mock.ExpectQuery("SELECT .* FROM \"sales_discount_rules\"").WillReturnRows(rows)
	found, err := repo.FindApplicable(context.Background(), uuid.New(), uuid.New(), 1, time.Now())
	if err != nil || len(found) != 1 {
		t.Fatalf("expected applicable discount")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGORMPriceCalculationRepository(t *testing.T) {
	gormDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewGORMPriceCalculationRepository(gormDB)

	money, _ := domain.NewMoney(10, "EUR")
	calc, _ := domain.NewPriceCalculation(uuid.New(), uuid.New(), 1, money, money, []string{"Rule"})

	mock.ExpectExec("UPDATE .*\"price_calculations\"").WillReturnResult(sqlmock.NewResult(1, 1))
	if err := repo.Save(context.Background(), calc); err != nil {
		t.Fatalf("expected save success, got %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "product_variant_id", "client_id", "quantity", "base_cost", "final_price", "currency", "applied_rules", "calculated_at"}).
		AddRow(uuid.New(), calc.ProductVariantID, calc.ClientID, 1, 10.0, 10.0, "EUR", "[\"Rule\"]", time.Now())
	mock.ExpectQuery("SELECT .* FROM \"price_calculations\"").WillReturnRows(rows)
	list, err := repo.ListByProductVariantID(context.Background(), calc.ProductVariantID)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected list success")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGORMBaseSalesPriceRuleRepository(t *testing.T) {
	gormDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewGORMBaseSalesPriceRuleRepository(gormDB)

	percentage, _ := domain.NewPercentage(0.1)
	value, _ := domain.NewRuleValue(domain.RuleValuePercentageMarkup, &percentage, nil)
	rule, _ := domain.NewBaseSalesPriceRule("Base", nil, nil, nil, nil, value)

	mock.ExpectExec("UPDATE .*\"base_sales_price_rules\"").WillReturnResult(sqlmock.NewResult(1, 1))
	if err := repo.Save(context.Background(), rule); err != nil {
		t.Fatalf("expected save success, got %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "value_type", "percentage_value", "money_value_currency", "is_active"}).
		AddRow(uuid.New(), "Base", string(domain.RuleValuePercentageMarkup), 0.1, "EUR", true)
	mock.ExpectQuery("SELECT .* FROM \"base_sales_price_rules\"").WillReturnRows(rows)
	found, err := repo.FindByID(context.Background(), uuid.New())
	if err != nil || found == nil {
		t.Fatalf("expected find success")
	}

	listRows := sqlmock.NewRows([]string{"id", "name", "value_type", "percentage_value", "money_value_currency", "is_active"}).
		AddRow(uuid.New(), "Base", string(domain.RuleValuePercentageMarkup), 0.1, "EUR", true)
	mock.ExpectQuery("SELECT .* FROM \"base_sales_price_rules\"").WillReturnRows(listRows)
	list, err := repo.List(context.Background())
	if err != nil || len(list) != 1 {
		t.Fatalf("expected list success")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGORMSaleModificationRuleRepository(t *testing.T) {
	gormDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewGORMSaleModificationRuleRepository(gormDB)

	percentage, _ := domain.NewPercentage(0.1)
	value, _ := domain.NewRuleValue(domain.RuleValuePercentageMarkup, &percentage, nil)
	rule, _ := domain.NewSaleModificationRule("Sale", nil, nil, nil, value, 1, time.Now(), nil)

	mock.ExpectExec("UPDATE .*\"sale_modification_rules\"").WillReturnResult(sqlmock.NewResult(1, 1))
	if err := repo.Save(context.Background(), rule); err != nil {
		t.Fatalf("expected save success, got %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "client_ids", "value_type", "percentage_value", "money_value_currency", "priority", "effective_from", "is_active", "min_order_total_currency"}).
		AddRow(uuid.New(), "Sale", pq.StringArray{}, string(domain.RuleValuePercentageMarkup), 0.1, "EUR", 1, time.Now(), true, "EUR")
	mock.ExpectQuery("SELECT .* FROM \"sale_modification_rules\"").WillReturnRows(rows)
	found, err := repo.FindByID(context.Background(), uuid.New())
	if err != nil || found == nil {
		t.Fatalf("expected find success")
	}

	clientID := uuid.New()
	applicableRows := sqlmock.NewRows([]string{"id", "name", "client_ids", "value_type", "percentage_value", "money_value_currency", "priority", "effective_from", "is_active", "min_order_total_currency"}).
		AddRow(uuid.New(), "Sale", pq.StringArray{clientID.String()}, string(domain.RuleValuePercentageMarkup), 0.1, "EUR", 1, time.Now(), true, "EUR")
	mock.ExpectQuery("SELECT .* FROM \"sale_modification_rules\"").WillReturnRows(applicableRows)
	orderTotal, _ := domain.NewMoney(100, "EUR")
	list, err := repo.ListApplicable(context.Background(), clientID, nil, orderTotal, time.Now())
	if err != nil || len(list) != 1 {
		t.Fatalf("expected applicable list")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
