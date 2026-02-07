package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type RedisBasePriceCache struct {
	client *redis.Client
	ttl    time.Duration
}

type basePricePayload struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func NewRedisBasePriceCache(client *redis.Client, ttl time.Duration) *RedisBasePriceCache {
	return &RedisBasePriceCache{client: client, ttl: ttl}
}

func (c *RedisBasePriceCache) GetBasePrice(ctx context.Context, productID uuid.UUID, variantID uuid.UUID) (*domain.Money, error) {
	key := basePriceKey(productID)
	payload, err := c.client.HGet(ctx, key, variantID.String()).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var parsed basePricePayload
	if err := json.Unmarshal([]byte(payload), &parsed); err != nil {
		return nil, err
	}

	money, err := domain.NewMoney(parsed.Amount, parsed.Currency)
	if err != nil {
		return nil, err
	}

	return &money, nil
}

func (c *RedisBasePriceCache) SetBasePrice(ctx context.Context, productID uuid.UUID, variantID uuid.UUID, price domain.Money) error {
	payload, err := json.Marshal(basePricePayload{Amount: price.Amount(), Currency: price.Currency()})
	if err != nil {
		return err
	}
	key := basePriceKey(productID)
	if err := c.client.HSet(ctx, key, variantID.String(), payload).Err(); err != nil {
		return err
	}
	return c.client.Expire(ctx, key, c.ttl).Err()
}

func basePriceKey(productID uuid.UUID) string {
	return fmt.Sprintf("pricing:base_price:product:%s", productID.String())
}
