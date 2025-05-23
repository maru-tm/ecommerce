package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"order-service/internal/domain"

	"github.com/google/uuid"
)

type orderUseCase struct {
	orderRepo      domain.OrderRepository
	orderPublisher domain.OrderEventPublisher
	orderCache     domain.OrderCache
}

func NewOrderUseCase(orderRepo domain.OrderRepository, orderPublisher domain.OrderEventPublisher, orderCache domain.OrderCache) domain.OrderUseCase {
	return &orderUseCase{
		orderRepo:      orderRepo,
		orderPublisher: orderPublisher,
		orderCache:     orderCache,
	}
}

func (uc *orderUseCase) log(level, msg string, fields map[string]interface{}) {
	logEntry := map[string]interface{}{
		"level": level,
		"msg":   msg,
		"time":  time.Now().Format(time.RFC3339),
	}
	for k, v := range fields {
		logEntry[k] = v
	}
	b, _ := json.Marshal(logEntry)
	fmt.Println(string(b))
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
	return nil
}

func (uc *orderUseCase) CreateOrder(order *domain.Order) error {
	if err := uc.validateOrder(order); err != nil {
		uc.log("ERROR", "Order validation failed", map[string]interface{}{"error": err.Error()})
		return err
	}

	order.ID = uuid.New().String()
	now := time.Now()
	order.CreatedAt = &now
	order.UpdatedAt = &now

	if err := uc.orderRepo.CreateOrder(order); err != nil {
		uc.log("ERROR", "Failed to create order in repository", map[string]interface{}{"error": err.Error(), "order_id": order.ID})
		return fmt.Errorf("failed to create order: %w", err)
	}

	if err := uc.orderPublisher.Publish(order); err != nil {
		uc.log("ERROR", "Failed to publish order creation event", map[string]interface{}{"error": err.Error(), "order_id": order.ID})
		return fmt.Errorf("failed to publish order creation event: %w", err)
	}

	uc.log("INFO", "Order created successfully", map[string]interface{}{"order_id": order.ID, "user_id": order.UserID})
	return nil
}

func (uc *orderUseCase) GetOrderByID(id string) (*domain.Order, error) {
	ctx := context.Background()

	order, err := uc.orderCache.GetOrder(ctx, id)
	if err == nil && order != nil {
		uc.log("INFO", "Order found in cache", map[string]interface{}{"order_id": id})
		return order, nil
	}

	uc.log("INFO", "Order not found in cache. Fetching from DB", map[string]interface{}{"order_id": id})
	order, err = uc.orderRepo.GetOrderByID(id)
	if err != nil {
		uc.log("ERROR", "Failed to fetch order from DB", map[string]interface{}{"error": err.Error(), "order_id": id})
		return nil, err
	}

	if err := uc.orderCache.SetOrder(ctx, order); err != nil {
		uc.log("WARN", "Failed to set order in cache", map[string]interface{}{"error": err.Error(), "order_id": id})
	}

	return order, nil
}

func (uc *orderUseCase) ListOrders() ([]domain.Order, error) {
	ctx := context.Background()
	cacheKey := "orders:all"

	orders, err := uc.orderCache.GetOrders(ctx, cacheKey)
	if err == nil && orders != nil {
		uc.log("INFO", "Orders found in cache", map[string]interface{}{"cache_key": cacheKey})
		return orders, nil
	}

	uc.log("INFO", "Orders not found in cache. Fetching from DB", nil)
	orders, err = uc.orderRepo.ListOrders()
	if err != nil {
		uc.log("ERROR", "Failed to fetch orders from DB", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	if err := uc.orderCache.SetOrders(ctx, cacheKey, orders); err != nil {
		uc.log("WARN", "Failed to set orders in cache", map[string]interface{}{"error": err.Error(), "cache_key": cacheKey})
	}

	return orders, nil
}

func (uc *orderUseCase) UpdateOrder(order *domain.Order) error {
	if err := uc.validateOrder(order); err != nil {
		uc.log("ERROR", "Order validation failed during update", map[string]interface{}{"error": err.Error(), "order_id": order.ID})
		return err
	}

	if err := uc.orderRepo.UpdateOrder(order); err != nil {
		uc.log("ERROR", "Failed to update order", map[string]interface{}{"error": err.Error(), "order_id": order.ID})
		return err
	}

	// Очистка кэша после успешного обновления
	ctx := context.Background()
	if err := uc.orderCache.DeleteOrder(ctx, order.ID); err != nil {
		uc.log("WARN", "Failed to delete order from cache after update", map[string]interface{}{"error": err.Error(), "order_id": order.ID})
	}

	uc.log("INFO", "Order updated successfully", map[string]interface{}{"order_id": order.ID})
	return nil
}

func (uc *orderUseCase) DeleteOrder(id string) error {
	if err := uc.orderRepo.DeleteOrder(id); err != nil {
		uc.log("ERROR", "Failed to delete order", map[string]interface{}{"error": err.Error(), "order_id": id})
		return err
	}

	// Очистка кэша после успешного удаления
	ctx := context.Background()
	if err := uc.orderCache.DeleteOrder(ctx, id); err != nil {
		uc.log("WARN", "Failed to delete order from cache after delete", map[string]interface{}{"error": err.Error(), "order_id": id})
	}

	uc.log("INFO", "Order deleted successfully", map[string]interface{}{"order_id": id})
	return nil
}
