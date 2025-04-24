package clients

import (
	"context"
	"fmt"

	proto "api-gateway/internal/proto/products/proto" // Путь к сгенерированным протобуфам

	"google.golang.org/grpc"
)

type ProductClient struct {
	client proto.ProductServiceClient
}

// NewProductClient создает новый gRPC-клиент для сервиса продуктов
func NewProductClient(address string) (*ProductClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}
	client := proto.NewProductServiceClient(conn)
	return &ProductClient{client: client}, nil
}

// CreateProduct создает новый продукт
func (pc *ProductClient) CreateProduct(ctx context.Context, product *proto.Product) (*proto.Product, error) {
	return pc.client.CreateProduct(ctx, product)
}

// GetProductByID возвращает продукт по ID
func (pc *ProductClient) GetProductByID(ctx context.Context, productId *proto.ProductId) (*proto.Product, error) {
	return pc.client.GetProductByID(ctx, productId)
}

// ListProducts возвращает список всех продуктов
func (pc *ProductClient) ListProducts(ctx context.Context) (*proto.ProductList, error) {
	return pc.client.ListProducts(ctx, &proto.Empty{})
}

// UpdateProduct обновляет данные продукта
func (pc *ProductClient) UpdateProduct(ctx context.Context, product *proto.Product) (*proto.Product, error) {
	return pc.client.UpdateProduct(ctx, product)
}

// DeleteProduct удаляет продукт
func (pc *ProductClient) DeleteProduct(ctx context.Context, productId *proto.ProductId) error {
	_, err := pc.client.DeleteProduct(ctx, productId)
	return err
}
