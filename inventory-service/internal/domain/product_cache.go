package domain

import "context"

type ProductCache interface {
	SetProduct(ctx context.Context, product *Product) error
	GetProduct(ctx context.Context, id string) (*Product, error)
	DeleteProduct(ctx context.Context, id string) error

	SetProducts(ctx context.Context, key string, products []Product) error
	GetProducts(ctx context.Context, key string) ([]Product, error)
}
