package persistence

import (
	"context"
	"testing"
)

func TestActorIDFromContext(t *testing.T) {
	if got := actorIDFromContext(context.Background()); got != "system" {
		t.Fatalf("expected system, got %s", got)
	}

	ctxWithActor := context.WithValue(context.Background(), "actorID", "actor-1")
	if got := actorIDFromContext(ctxWithActor); got != "actor-1" {
		t.Fatalf("expected actor-1, got %s", got)
	}

	ctxWithUser := context.WithValue(context.Background(), "userID", "user-1")
	if got := actorIDFromContext(ctxWithUser); got != "user-1" {
		t.Fatalf("expected user-1, got %s", got)
	}

	ctxWithBoth := context.WithValue(ctxWithActor, "userID", "user-2")
	if got := actorIDFromContext(ctxWithBoth); got != "user-2" {
		t.Fatalf("expected user-2, got %s", got)
	}
}

func TestTableNames(t *testing.T) {
	if (AttributeDataModel{}).TableName() != "attributes" {
		t.Fatalf("unexpected attributes table name")
	}
	if (AttributeValueDataModel{}).TableName() != "attribute_values" {
		t.Fatalf("unexpected attribute values table name")
	}
	if (ProductDataModel{}).TableName() != "products" {
		t.Fatalf("unexpected products table name")
	}
	if (VariantDataModel{}).TableName() != "product_variants" {
		t.Fatalf("unexpected variants table name")
	}
	if (BrandDataModel{}).TableName() != "brands" {
		t.Fatalf("unexpected brands table name")
	}
	if (ProductGroupDataModel{}).TableName() != "product_groups" {
		t.Fatalf("unexpected product groups table name")
	}
	if (PartyServiceConfigurationModel{}).TableName() != "party_service_configurations" {
		t.Fatalf("unexpected party service configurations table name")
	}
}
