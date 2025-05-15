package usecase

import (
	"context"
	"fmt"
	"log"

	"order-service/internal/domain"
	"order-service/internal/infrastructure/messaging"

	"github.com/google/uuid"
)

type orderUseCase struct {
	orderRepo      domain.OrderRepository
	orderPublisher *messaging.OrderPublisher
	orderCache     domain.OrderCache
}

func NewOrderUseCase(orderRepo domain.OrderRepository, orderPublisher *messaging.OrderPublisher, orderCache domain.OrderCache) domain.OrderUseCase {
	return &orderUseCase{
		orderRepo:      orderRepo,
		orderPublisher: orderPublisher,
		orderCache:     orderCache,
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
	ctx := context.Background()

	// Попытка получить из кэша
	if order, _ := uc.orderCache.GetOrder(ctx, id); order != nil {
		log.Printf("[GET] Order found in cache: %s", id)
		return order, nil
	}

	// Из базы
	log.Printf("[GET] Order not found in cache. Fetching from DB: %s", id)
	order, err := uc.orderRepo.GetOrderByID(id)
	if err != nil {
		log.Printf("[GET] Error fetching order from DB: %v", err)
		return nil, err
	}

	// Кэшируем
	if err := uc.orderCache.SetOrder(ctx, order); err != nil {
		log.Printf("[CACHE] Failed to set order in cache: %v", err)
	}

	return order, nil
}

func (uc *orderUseCase) ListOrders() ([]domain.Order, error) {
	ctx := context.Background()
	cacheKey := "orders:all"

	// Попытка получить из кэша
	if orders, _ := uc.orderCache.GetOrders(ctx, cacheKey); orders != nil {
		log.Println("[GET] Orders found in cache")
		return orders, nil
	}

	// Из базы
	log.Println("[GET] Orders not found in cache. Fetching from DB")
	orders, err := uc.orderRepo.ListOrders()
	if err != nil {
		log.Printf("[GET] Error fetching orders from DB: %v", err)
		return nil, err
	}

	// Кэшируем
	if err := uc.orderCache.SetOrders(ctx, cacheKey, orders); err != nil {
		log.Printf("[CACHE] Failed to set orders in cache: %v", err)
	}

	return orders, nil
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
