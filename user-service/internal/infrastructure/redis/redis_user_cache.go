package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"user-service/internal/domain"

	"github.com/redis/go-redis/v9"
)

type userCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewUserCache(client *redis.Client, ttl time.Duration) domain.UserCache {
	return &userCache{
		client: client,
		ttl:    ttl,
	}
}

func (r *userCache) SetUser(ctx context.Context, user *domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("user:%s", user.ID)
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *userCache) GetUser(ctx context.Context, id string) (*domain.User, error) {
	key := fmt.Sprintf("user:%s", id)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // cache miss
	}
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userCache) DeleteUser(ctx context.Context, id string) error {
	key := fmt.Sprintf("user:%s", id)
	return r.client.Del(ctx, key).Err()
}

func (r *userCache) SetUsers(ctx context.Context, key string, users []domain.User) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *userCache) GetUsers(ctx context.Context, key string) ([]domain.User, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // cache miss
	}
	if err != nil {
		return nil, err
	}

	var users []domain.User
	if err := json.Unmarshal([]byte(val), &users); err != nil {
		return nil, err
	}
	return users, nil
}
