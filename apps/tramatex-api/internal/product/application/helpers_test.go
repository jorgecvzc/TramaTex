package application

import (
	"testing"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// COMMENTED OUT: attributeMatchesScopeType function was removed during scope refactoring
// This test is obsolete since the system now uses DirectAttributeIDs instead of scope-based inheritance
/*
func TestAttributeMatchesScopeType(t *testing.T) {
	// brandID := uuid.New()
	// groupID := uuid.New()

	genericAttr, _ := domain.NewAttribute("Generic", "G")
	brandAttr, _ := domain.NewAttribute("Brand", "B")
	groupBrandAttr, _ := domain.NewAttribute("GroupBrand", "GB")

	if !attributeMatchesScopeType(genericAttr, nil) {
		t.Fatalf("expected nil scope type to match")
	}

	scopeGeneric := "GENERIC"
	if !attributeMatchesScopeType(genericAttr, &scopeGeneric) {
		t.Fatalf("expected generic attr to match GENERIC scope")
	}
	if attributeMatchesScopeType(brandAttr, &scopeGeneric) {
		t.Fatalf("expected brand attr to not match GENERIC scope")
	}

	scopeBrand := "BRAND"
	if !attributeMatchesScopeType(brandAttr, &scopeBrand) {
		t.Fatalf("expected brand attr to match BRAND scope")
	}
	if attributeMatchesScopeType(groupBrandAttr, &scopeBrand) {
		t.Fatalf("expected brand+group attr to not match BRAND scope")
	}

	scopeBrandGroup := "BRAND_GROUP"
	if !attributeMatchesScopeType(groupBrandAttr, &scopeBrandGroup) {
		t.Fatalf("expected brand+group attr to match BRAND_GROUP scope")
	}

	scopeInvalid := "INVALID"
	if attributeMatchesScopeType(genericAttr, &scopeInvalid) {
		t.Fatalf("expected invalid scope to not match")
	}
}
*/

func TestProductMatchesQueryAndHasGroup(t *testing.T) {
	brandID := uuid.New()
	groupID := uuid.New()
	otherGroupID := uuid.New()
	product := &domain.Product{
		ID:       uuid.New(),
		SKU:      "P-1",
		Name:     "Product",
		BrandID:  brandID,
		GroupIDs: []uuid.UUID{groupID},
		IsActive: true,
	}

	query := ListProductsQuery{BrandID: &brandID, GroupID: &groupID}
	if !productMatchesQuery(product, query) {
		t.Fatalf("expected product to match brand and group")
	}

	query = ListProductsQuery{GroupID: &otherGroupID}
	if productMatchesQuery(product, query) {
		t.Fatalf("expected product to not match missing group")
	}

	inactive := false
	query = ListProductsQuery{IsActive: &inactive}
	if productMatchesQuery(product, query) {
		t.Fatalf("expected product to not match inactive filter")
	}

	if !productHasGroup(product, groupID) {
		t.Fatalf("expected productHasGroup to find group")
	}
	if productHasGroup(product, otherGroupID) {
		t.Fatalf("expected productHasGroup to return false for missing group")
	}
}
