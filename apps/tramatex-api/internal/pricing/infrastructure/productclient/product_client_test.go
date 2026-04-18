package productclient

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	productapp "github.com/joran-cortez/tramatex/internal/product/application"
)

type fakeProductDataProvider struct {
	variantData  *productapp.VariantPricingDataDTO
	variantErr   error
	variantsList []*productapp.VariantPricingDataDTO
	listErr      error
}

func (f *fakeProductDataProvider) GetVariantPricingData(_ context.Context, _ uuid.UUID) (*productapp.VariantPricingDataDTO, error) {
	return f.variantData, f.variantErr
}

func (f *fakeProductDataProvider) GetVariantsPricingData(_ context.Context, _ []uuid.UUID) ([]*productapp.VariantPricingDataDTO, error) {
	return f.variantsList, f.listErr
}

func (f *fakeProductDataProvider) ListVariantsPricingData(_ context.Context, _ uuid.UUID) ([]*productapp.VariantPricingDataDTO, error) {
	return f.variantsList, f.listErr
}

func TestProductPricingClient_GetVariantPricingInfo_NotFound(t *testing.T) {
	fake := &fakeProductDataProvider{variantData: nil, variantErr: nil}
	client := NewProductPricingClient(fake)

	info, err := client.GetVariantPricingInfo(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if info != nil {
		t.Fatalf("expected nil info")
	}
}

func TestProductPricingClient_GetVariantPricingInfo_Error(t *testing.T) {
	fake := &fakeProductDataProvider{variantErr: errors.New("db error")}
	client := NewProductPricingClient(fake)

	info, err := client.GetVariantPricingInfo(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error")
	}
	if info != nil {
		t.Fatal("expected nil info on error")
	}
}

func TestProductPricingClient_GetVariantPricingInfo_Success(t *testing.T) {
	variantID := uuid.New()
	productID := uuid.New()
	brandID := uuid.New()
	groupID := uuid.New()

	fake := &fakeProductDataProvider{
		variantData: &productapp.VariantPricingDataDTO{
			VariantID:             variantID,
			ProductID:             productID,
			BaseCost:              25.0,
			Currency:              "EUR",
			BrandID:               brandID,
			BrandMarkupPercentage: 10.0,
			GroupIDs:              []uuid.UUID{groupID},
			TaxRate:               21.0,
		},
	}
	client := NewProductPricingClient(fake)

	info, err := client.GetVariantPricingInfo(context.Background(), variantID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if info == nil {
		t.Fatal("expected non-nil info")
	}
	if info.VariantID != variantID || info.ProductID != productID {
		t.Fatalf("unexpected IDs: variant=%s product=%s", info.VariantID, info.ProductID)
	}
	if info.BrandID != brandID || !info.BrandMarkupPercentage.Equal(decimal.NewFromFloat(10.0)) {
		t.Fatal("unexpected brand data")
	}
	if len(info.GroupIDs) != 1 || info.GroupIDs[0] != groupID {
		t.Fatal("expected group id parsed")
	}
	if !info.BaseCost.Equal(decimal.NewFromFloat(25.0)) || !info.TaxRate.Equal(decimal.NewFromFloat(21.0)) {
		t.Fatal("unexpected pricing data")
	}
}

func TestProductPricingClient_ListVariantsPricingInfo_Success(t *testing.T) {
	productID := uuid.New()
	v1 := uuid.New()
	v2 := uuid.New()

	fake := &fakeProductDataProvider{
		variantsList: []*productapp.VariantPricingDataDTO{
			{VariantID: v1, ProductID: productID, BaseCost: 10.0, Currency: "EUR"},
			{VariantID: v2, ProductID: productID, BaseCost: 20.0, Currency: "EUR"},
		},
	}
	client := NewProductPricingClient(fake)

	infos, err := client.ListVariantsPricingInfo(context.Background(), productID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(infos) != 2 {
		t.Fatalf("expected 2 infos, got %d", len(infos))
	}
	if infos[0].VariantID != v1 || infos[1].VariantID != v2 {
		t.Fatal("unexpected variant order")
	}
}
