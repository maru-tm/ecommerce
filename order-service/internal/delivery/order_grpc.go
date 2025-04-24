package delivery

import (
	"context"
	"log"

	"order-service/internal/domain"
	"order-service/internal/proto"
	"order-service/internal/usecase"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderServiceServer struct {
	proto.UnimplementedOrderServiceServer
	useCase usecase.OrderUseCase
}

func NewOrderServiceServer(useCase usecase.OrderUseCase) *OrderServiceServer {
	log.Println("Initializing OrderServiceServer")
	return &OrderServiceServer{
		useCase: useCase,
	}
}

func (s *OrderServiceServer) CreateOrder(ctx context.Context, req *proto.Order) (*proto.Order, error) {
	log.Printf("CreateOrder called with request: %v", req)

	order := &domain.Order{
		ID:         req.GetId(),
		UserID:     req.GetUserId(),
		TotalPrice: req.GetTotalPrice(),
		Status:     domain.OrderStatus(req.GetStatus().String()),
		CreatedAt:  req.GetCreatedAt().AsTime(),
		UpdatedAt:  req.GetUpdatedAt().AsTime(),
	}

	var orderItems []domain.OrderItem
	for _, item := range req.GetItems() {
		orderItems = append(orderItems, domain.OrderItem{
			ProductID: item.GetProductId(),
			Quantity:  int(item.GetQuantity()),
		})
	}
	order.Items = orderItems

	err := s.useCase.CreateOrder(order)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		return nil, err
	}

	log.Printf("Order successfully created: %v", order)

	return &proto.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		TotalPrice: order.TotalPrice,
		Status:     proto.OrderStatus(proto.OrderStatus_value[string(order.Status)]),
		CreatedAt:  timestamppb.New(order.CreatedAt),
		UpdatedAt:  timestamppb.New(order.UpdatedAt),
		Items:      mapOrderItemsToProto(order.Items),
	}, nil
}

func (s *OrderServiceServer) GetOrderByID(ctx context.Context, req *proto.OrderId) (*proto.Order, error) {
	log.Printf("GetOrderByID called with order ID: %v", req.GetId())

	order, err := s.useCase.GetOrderByID(req.GetId())
	if err != nil {
		log.Printf("Error getting order by ID: %v", err)
		return nil, err
	}

	log.Printf("Order found: %v", order)

	return &proto.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		TotalPrice: order.TotalPrice,
		Status:     proto.OrderStatus(proto.OrderStatus_value[string(order.Status)]),
		CreatedAt:  timestamppb.New(order.CreatedAt),
		UpdatedAt:  timestamppb.New(order.UpdatedAt),
		Items:      mapOrderItemsToProto(order.Items),
	}, nil
}

func (s *OrderServiceServer) ListOrders(ctx context.Context, req *proto.Empty) (*proto.OrderList, error) {
	log.Println("ListOrders called")

	orders, err := s.useCase.ListOrders()
	if err != nil {
		log.Printf("Error listing orders: %v", err)
		return nil, err
	}

	log.Printf("Found %d orders", len(orders))

	var orderList []*proto.Order
	for _, o := range orders {
		orderList = append(orderList, &proto.Order{
			Id:         o.ID,
			UserId:     o.UserID,
			TotalPrice: o.TotalPrice,
			Status:     proto.OrderStatus(proto.OrderStatus_value[string(o.Status)]),
			CreatedAt:  timestamppb.New(o.CreatedAt),
			UpdatedAt:  timestamppb.New(o.UpdatedAt),
			Items:      mapOrderItemsToProto(o.Items),
		})
	}

	log.Printf("Returning order list with %d orders", len(orderList))

	return &proto.OrderList{Orders: orderList}, nil
}

func (s *OrderServiceServer) UpdateOrder(ctx context.Context, req *proto.Order) (*proto.Order, error) {
	log.Printf("UpdateOrder called with request: %v", req)

	order := &domain.Order{
		ID:         req.GetId(),
		UserID:     req.GetUserId(),
		TotalPrice: req.GetTotalPrice(),
		Status:     domain.OrderStatus(req.GetStatus().String()),
		CreatedAt:  req.GetCreatedAt().AsTime(),
		UpdatedAt:  req.GetUpdatedAt().AsTime(),
	}

	var orderItems []domain.OrderItem
	for _, item := range req.GetItems() {
		orderItems = append(orderItems, domain.OrderItem{
			ProductID: item.GetProductId(),
			Quantity:  int(item.GetQuantity()),
		})
	}
	order.Items = orderItems

	err := s.useCase.UpdateOrder(order)
	if err != nil {
		log.Printf("Error updating order: %v", err)
		return nil, err
	}

	log.Printf("Order successfully updated: %v", order)

	return &proto.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		TotalPrice: order.TotalPrice,
		Status:     proto.OrderStatus(proto.OrderStatus_value[string(order.Status)]),
		CreatedAt:  timestamppb.New(order.CreatedAt),
		UpdatedAt:  timestamppb.New(order.UpdatedAt),
		Items:      mapOrderItemsToProto(order.Items),
	}, nil
}

func (s *OrderServiceServer) DeleteOrder(ctx context.Context, req *proto.OrderId) (*proto.Empty, error) {
	log.Printf("DeleteOrder called with order ID: %v", req.GetId())

	err := s.useCase.DeleteOrder(req.GetId())
	if err != nil {
		log.Printf("Error deleting order: %v", err)
		return nil, err
	}

	log.Printf("Order successfully deleted with ID: %v", req.GetId())

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
