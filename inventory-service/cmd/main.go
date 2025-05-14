package main

import (
	"fmt"
	"log"
	"net"

	"inventory-service/config"
	"inventory-service/internal/delivery"
	"inventory-service/internal/messaging"
	"inventory-service/internal/proto"
	"inventory-service/internal/repository"
	"inventory-service/internal/usecase"

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
	if err := config.InitRedis(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}

	productRepo := repository.NewProductRepository(db)
	productUC := usecase.NewProductUseCase(productRepo)
	productHandler := delivery.NewProductServiceServer(productUC)

	port := ":" + cfg.InventoryServerPort
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	proto.RegisterProductServiceServer(grpcServer, productHandler)

	go messaging.StartConsumer(productUC)

	fmt.Printf("gRPC server started on %s\n", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
