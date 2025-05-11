package usecase

import (
	"fmt"

	"order-service/internal/domain"
	"order-service/internal/infrastructure/messaging"

	"github.com/google/uuid"
)

type orderUseCase struct {
	orderRepo      domain.OrderRepository
	orderPublisher *messaging.OrderPublisher
}

func NewOrderUseCase(orderRepo domain.OrderRepository, orderPublisher *messaging.OrderPublisher) domain.OrderUseCase {
	return &orderUseCase{
		orderRepo:      orderRepo,
		orderPublisher: orderPublisher,
	}
}

func (uc *orderUseCase) validateOrder(order *domain.Order) error {
	// Проверяем обязательные поля в заказе
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

	return nil
}

func (uc *orderUseCase) CreateOrder(order *domain.Order) error {
	// Валидация заказа
	if err := uc.validateOrder(order); err != nil {
		return err
	}

	// Генерируем ID для заказа, если его нет
	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	// Сохраняем заказ в репозитории
	if err := uc.orderRepo.CreateOrder(order); err != nil {
		return fmt.Errorf("failed to create order: %v", err)
	}

	// Публикуем сообщение о создании заказа в RabbitMQ
	if err := uc.orderPublisher.Publish(order); err != nil {
		return fmt.Errorf("failed to publish order creation event: %v", err)
	}

	return nil
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
