package delivery

import (
	"context"
	"log"

	"order-service/internal/domain"
	"order-service/internal/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderServiceServer struct {
	proto.UnimplementedOrderServiceServer
	useCase domain.OrderUseCase
}

func NewOrderServiceServer(useCase domain.OrderUseCase) *OrderServiceServer {
	log.Println("[INFO] Initializing OrderServiceServer")
	return &OrderServiceServer{useCase: useCase}
}

func (s *OrderServiceServer) CreateOrder(ctx context.Context, req *proto.Order) (*proto.Order, error) {
	log.Printf("[INFO] CreateOrder: received request: ID=%s, UserID=%s", req.GetId(), req.GetUserId())

	order := mapProtoToDomain(req)
	err := s.useCase.CreateOrder(order)
	if err != nil {
		log.Printf("[ERROR] CreateOrder: failed to create order ID=%s: %v", order.ID, err)
		return nil, err
	}

	log.Printf("[INFO] CreateOrder: successfully created order ID=%s", order.ID)
	return mapDomainToProto(order), nil
}

func (s *OrderServiceServer) GetOrderByID(ctx context.Context, req *proto.OrderId) (*proto.Order, error) {
	log.Printf("[INFO] GetOrderByID: looking for order ID=%s", req.GetId())

	order, err := s.useCase.GetOrderByID(req.GetId())
	if err != nil {
		log.Printf("[ERROR] GetOrderByID: failed to find order ID=%s: %v", req.GetId(), err)
		return nil, err
	}

	log.Printf("[INFO] GetOrderByID: found order ID=%s", order.ID)
	return mapDomainToProto(order), nil
}

func (s *OrderServiceServer) ListOrders(ctx context.Context, req *proto.Empty) (*proto.OrderList, error) {
	log.Println("[INFO] ListOrders: fetching all orders")

	orders, err := s.useCase.ListOrders()
	if err != nil {
		log.Printf("[ERROR] ListOrders: failed to fetch orders: %v", err)
		return nil, err
	}

	log.Printf("[INFO] ListOrders: retrieved %d orders", len(orders))

	var protoOrders []*proto.Order
	for _, order := range orders {
		protoOrders = append(protoOrders, mapDomainToProto(&order))
	}

	return &proto.OrderList{Orders: protoOrders}, nil
}

func (s *OrderServiceServer) UpdateOrder(ctx context.Context, req *proto.Order) (*proto.Order, error) {
	log.Printf("[INFO] UpdateOrder: updating order ID=%s", req.GetId())

	order := mapProtoToDomain(req)
	err := s.useCase.UpdateOrder(order)
	if err != nil {
		log.Printf("[ERROR] UpdateOrder: failed to update order ID=%s: %v", order.ID, err)
		return nil, err
	}

	log.Printf("[INFO] UpdateOrder: successfully updated order ID=%s", order.ID)
	return mapDomainToProto(order), nil
}

func (s *OrderServiceServer) DeleteOrder(ctx context.Context, req *proto.OrderId) (*proto.Empty, error) {
	log.Printf("[INFO] DeleteOrder: deleting order ID=%s", req.GetId())

	err := s.useCase.DeleteOrder(req.GetId())
	if err != nil {
		log.Printf("[ERROR] DeleteOrder: failed to delete order ID=%s: %v", req.GetId(), err)
		return nil, err
	}

	log.Printf("[INFO] DeleteOrder: successfully deleted order ID=%s", req.GetId())
	return &proto.Empty{}, nil
}

func mapOrderItemsToProto(orderItems []domain.OrderItem) []*proto.OrderItem {
	var protoItems []*proto.OrderItem
	for _, item := range orderItems {
		protoItems = append(protoItems, &proto.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
	}
	return protoItems
}

func mapProtoToDomain(req *proto.Order) *domain.Order {
	createdAt := req.GetCreatedAt().AsTime()
	updatedAt := req.GetUpdatedAt().AsTime()

	var items []domain.OrderItem
	for _, item := range req.GetItems() {
		items = append(items, domain.OrderItem{
			ProductID: item.GetProductId(),
			Quantity:  int(item.GetQuantity()),
		})
	}

	return &domain.Order{
		ID:         req.GetId(),
		UserID:     req.GetUserId(),
		TotalPrice: req.GetTotalPrice(),
		Status:     domain.OrderStatus(req.GetStatus().String()),
		CreatedAt:  &createdAt,
		UpdatedAt:  &updatedAt,
		Items:      items,
	}
}

func mapDomainToProto(order *domain.Order) *proto.Order {
	return &proto.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		TotalPrice: order.TotalPrice,
		Status:     proto.OrderStatus(proto.OrderStatus_value[string(order.Status)]),
		CreatedAt:  timestamppb.New(*order.CreatedAt),
		UpdatedAt:  timestamppb.New(*order.UpdatedAt),
		Items:      mapOrderItemsToProto(order.Items),
	}
}
