package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"inventory-service/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type productCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewProductCache(client *redis.Client, ttl time.Duration) domain.ProductCache {
	return &productCache{
		client: client,
		ttl:    ttl,
	}
}

func (r *productCache) SetProduct(ctx context.Context, product *domain.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("product:%s", product.ID)
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *productCache) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	key := fmt.Sprintf("product:%s", id)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // кэш-промах
	}
	if err != nil {
		return nil, err
	}

	var product domain.Product
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productCache) DeleteProduct(ctx context.Context, id string) error {
	key := fmt.Sprintf("product:%s", id)
	return r.client.Del(ctx, key).Err()
}

func (r *productCache) SetProducts(ctx context.Context, key string, products []domain.Product) error {
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *productCache) GetProducts(ctx context.Context, key string) ([]domain.Product, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // кэш-промах
	}
	if err != nil {
		return nil, err
	}

	var products []domain.Product
	if err := json.Unmarshal([]byte(val), &products); err != nil {
		return nil, err
	}
	return products, nil
}
