package repository

import (
	"order-service/internal/domain"
)

type OrderRepository interface {
	CreateOrder(order *domain.Order) error
	GetOrderByID(id string) (*domain.Order, error)
	ListOrders() ([]domain.Order, error)
	UpdateOrder(order *domain.Order) error
	DeleteOrder(id string) error
}
