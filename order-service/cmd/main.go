package main

import (
	"fmt"
	"log"
	"net"

	"order-service/config"
	"order-service/internal/delivery"
	"order-service/internal/proto"
	"order-service/internal/repository"
	"order-service/internal/usecase"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatalf("failed to load configuration")
	}

	db := config.ConnectDB(cfg)
	if db == nil {
		log.Fatalf("failed to connect to the database")
	}

	orderRepo := repository.NewOrderRepository(db)
	orderUC := usecase.NewOrderUseCase(orderRepo)
	orderHandler := delivery.NewOrderServiceServer(orderUC)

	port := ":50052"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	proto.RegisterOrderServiceServer(grpcServer, orderHandler)

	fmt.Printf("gRPC server started on %s\n", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
