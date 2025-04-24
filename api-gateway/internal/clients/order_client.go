package clients

import (
	"context"
	"fmt"

	proto "api-gateway/internal/proto/orders/proto" // Путь к сгенерированным протобуфам

	"google.golang.org/grpc"
)

type OrderClient struct {
	client proto.OrderServiceClient
}

// NewOrderClient создает новый gRPC-клиент для сервиса заказов
func NewOrderClient(address string) (*OrderClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}
	client := proto.NewOrderServiceClient(conn)
	return &OrderClient{client: client}, nil
}

// CreateOrder создает новый заказ
func (oc *OrderClient) CreateOrder(ctx context.Context, order *proto.Order) (*proto.Order, error) {
	return oc.client.CreateOrder(ctx, order)
}

// GetOrderByID возвращает заказ по ID
func (oc *OrderClient) GetOrderByID(ctx context.Context, orderId *proto.OrderId) (*proto.Order, error) {
	return oc.client.GetOrderByID(ctx, orderId)
}

// ListOrders возвращает список всех заказов
func (oc *OrderClient) ListOrders(ctx context.Context) (*proto.OrderList, error) {
	return oc.client.ListOrders(ctx, &proto.Empty{})
}

// UpdateOrder обновляет данные заказа
func (oc *OrderClient) UpdateOrder(ctx context.Context, order *proto.Order) (*proto.Order, error) {
	return oc.client.UpdateOrder(ctx, order)
}

// DeleteOrder удаляет заказ
func (oc *OrderClient) DeleteOrder(ctx context.Context, orderId *proto.OrderId) error {
	_, err := oc.client.DeleteOrder(ctx, orderId)
	return err
}
