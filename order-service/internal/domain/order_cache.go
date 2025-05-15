package domain

import "context"

type OrderCache interface {
	SetOrder(ctx context.Context, order *Order) error
	GetOrder(ctx context.Context, id string) (*Order, error)
	DeleteOrder(ctx context.Context, id string) error
	SetOrders(ctx context.Context, key string, orders []Order) error
	GetOrders(ctx context.Context, key string) ([]Order, error)
}
