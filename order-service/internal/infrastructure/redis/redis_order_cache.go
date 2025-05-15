package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type orderCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewOrderCache(client *redis.Client, ttl time.Duration) domain.OrderCache {
	return &orderCache{
		client: client,
		ttl:    ttl,
	}
}

func (r *orderCache) SetOrder(ctx context.Context, order *domain.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("order:%s", order.ID)
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *orderCache) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	key := fmt.Sprintf("order:%s", id)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // кэш-промах
	}
	if err != nil {
		return nil, err
	}

	var order domain.Order
	if err := json.Unmarshal([]byte(val), &order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderCache) DeleteOrder(ctx context.Context, id string) error {
	key := fmt.Sprintf("order:%s", id)
	return r.client.Del(ctx, key).Err()
}

func (r *orderCache) SetOrders(ctx context.Context, key string, orders []domain.Order) error {
	data, err := json.Marshal(orders)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *orderCache) GetOrders(ctx context.Context, key string) ([]domain.Order, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // кэш-промах
	}
	if err != nil {
		return nil, err
	}

	var orders []domain.Order
	if err := json.Unmarshal([]byte(val), &orders); err != nil {
		return nil, err
	}
	return orders, nil
}
