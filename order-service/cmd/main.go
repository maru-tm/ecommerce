package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"order-service/config"
	"order-service/internal/delivery"
	"order-service/internal/infrastructure/messaging"
	"order-service/internal/infrastructure/redis"
	"order-service/internal/infrastructure/repository"
	"order-service/internal/proto"
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

	ctx := context.Background()
	if err := config.InitRedis(ctx); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
	redisClient := config.GetRedisClient()
	orderCache := redis.NewOrderCache(redisClient, 5*time.Minute) // создаём отдельный кэш-слой

	orderRepo := repository.NewOrderRepository(db)
	publisher, err := messaging.NewOrderPublisher(cfg)
	if err != nil {
		log.Fatalf("Ошибка создания OrderPublisher: %v", err)
	}

	orderUC := usecase.NewOrderUseCase(orderRepo, publisher, orderCache)

	orderHandler := delivery.NewOrderServiceServer(orderUC)

	port := ":" + cfg.OrderServicePort
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
