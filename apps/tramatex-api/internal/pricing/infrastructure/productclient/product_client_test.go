package productclient

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		sqlDB.Close()
		t.Fatalf("failed to open gorm db: %v", err)
	}
	cleanup := func() { _ = sqlDB.Close() }
	return gormDB, mock, cleanup
}

func TestProductPricingClient_GetVariantPricingInfo_NotFound(t *testing.T) {
	gormDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	mock.ExpectQuery("SELECT .* FROM \"product_variants\"").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	client := NewProductPricingClient(gormDB)
	info, err := client.GetVariantPricingInfo(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if info != nil {
		t.Fatalf("expected nil info")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductPricingClient_GetVariantPricingInfo_Success(t *testing.T) {
	gormDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	variantID := uuid.New()
	productID := uuid.New()
	brandID := uuid.New()
	groupID := uuid.New()

	mock.ExpectQuery("SELECT .* FROM \"product_variants\"").
		WithArgs(variantID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "base_cost"}).
			AddRow(variantID, productID, 25.0))

	mock.ExpectQuery("SELECT .* FROM \"products\"").
		WithArgs(productID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "brand_id", "group_ids"}).
			AddRow(productID, brandID, pq.StringArray{groupID.String()}))

	client := NewProductPricingClient(gormDB)
	info, err := client.GetVariantPricingInfo(context.Background(), variantID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if info == nil || info.ProductID != productID || info.BrandID != brandID {
		t.Fatalf("unexpected info result")
	}
	if len(info.GroupIDs) != 1 || info.GroupIDs[0] != groupID {
		t.Fatalf("expected group id parsed")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
