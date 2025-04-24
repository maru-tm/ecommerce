package usecase

import (
	"fmt"

	"order-service/internal/domain"
	"order-service/internal/repository"

	"github.com/google/uuid"
)

type OrderUseCase interface {
	CreateOrder(order *domain.Order) error
	GetOrderByID(id string) (*domain.Order, error)
	ListOrders() ([]domain.Order, error)
	UpdateOrder(order *domain.Order) error
	DeleteOrder(id string) error
}

type orderUseCase struct {
	orderRepo repository.OrderRepository
}

func NewOrderUseCase(orderRepo repository.OrderRepository) OrderUseCase {
	return &orderUseCase{orderRepo: orderRepo}
}

func (uc *orderUseCase) validateOrder(order *domain.Order) error {
	if order.UserID == "" {
		return fmt.Errorf("user id cannot be empty")
	}
	if len(order.Items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}
	for _, item := range order.Items {
		if item.ProductID == "" {
			return fmt.Errorf("order item must have a valid product ID")
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("order item quantity must be positive")
		}
	}

	products, err := getAllProducts("http://localhost:8080/products")
	if err != nil {
		return fmt.Errorf("failed to fetch products: %v", err)
	}

	for _, orderItem := range order.Items {
		var foundProduct *Product
		for _, product := range products {
			if product.ID == orderItem.ProductID {
				foundProduct = &product
				break
			}
		}

		if foundProduct == nil {
			return fmt.Errorf("product with ID %s not found in inventory", orderItem.ProductID)
		}

		if foundProduct.Stock < orderItem.Quantity {
			return fmt.Errorf("not enough stock for product %s", foundProduct.Name)
		}

		newStock := foundProduct.Stock - orderItem.Quantity
		if err := patchProductStock(orderItem.ProductID, newStock); err != nil {
			return fmt.Errorf("failed to update stock for product %s: %v", foundProduct.Name, err)
		}
	}

	return nil
}

func (uc *orderUseCase) CreateOrder(order *domain.Order) error {
	// if err := uc.validateOrder(order); err != nil {
	// 	return err
	// }

	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	return uc.orderRepo.CreateOrder(order)
}

func (uc *orderUseCase) GetOrderByID(id string) (*domain.Order, error) {
	return uc.orderRepo.GetOrderByID(id)
}

func (uc *orderUseCase) ListOrders() ([]domain.Order, error) {
	return uc.orderRepo.ListOrders()
}

func (uc *orderUseCase) UpdateOrder(order *domain.Order) error {
	if err := uc.validateOrder(order); err != nil {
		return err
	}

	return uc.orderRepo.UpdateOrder(order)
}

func (uc *orderUseCase) DeleteOrder(id string) error {
	return uc.orderRepo.DeleteOrder(id)
}
