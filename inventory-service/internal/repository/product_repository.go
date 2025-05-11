package repository

import (
	"inventory-service/internal/domain"
)

type ProductRepository interface {
	CreateProduct(product *domain.Product) error
	GetProductByID(id string) (*domain.Product, error)
	GetProductByName(name string) (*domain.Product, error)
	ListProducts() ([]domain.Product, error)
	UpdateProduct(product *domain.Product) error
	DeleteProduct(id string) error
	CheckStock(productID string) (bool, error)
	SearchProducts(keyword string, categoryID string) ([]domain.Product, error)
}
