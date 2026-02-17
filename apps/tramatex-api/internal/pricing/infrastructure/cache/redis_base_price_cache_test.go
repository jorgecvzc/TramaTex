package cache

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

func TestRedisBasePriceCache_ErrorsWithoutServer(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"})
	cache := NewRedisBasePriceCache(client, time.Second)

	productID := uuid.New()
	variantID := uuid.New()

	_, err := cache.GetBasePrice(context.Background(), productID, variantID)
	if err == nil {
		t.Fatalf("expected error without redis server")
	}

	price, _ := domain.NewMoney(10, "EUR")
	if err := cache.SetBasePrice(context.Background(), productID, variantID, price); err == nil {
		t.Fatalf("expected error without redis server")
	}
}
