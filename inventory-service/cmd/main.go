package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"inventory-service/config"
	"inventory-service/internal/delivery"
	"inventory-service/internal/infrastructure/messaging"
	"inventory-service/internal/infrastructure/redis"
	"inventory-service/internal/infrastructure/repository"
	"inventory-service/internal/proto"
	"inventory-service/internal/usecase"

	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatalf("failed to load configuration")
	}

	db := config.ConnectDB(cfg)
	if db == nil {
		log.Fatalf("failed to connect to the database")
	}

	err := repository.RunMigrations(ctx, db)
	if err != nil {
		log.Fatalf("Ошибка миграции БД: %v", err)
	}

	if err := config.InitRedis(ctx); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}

	redisClient := config.GetRedisClient()

	productRepo := repository.NewProductRepository(db)
	productCache := redis.NewProductCache(redisClient, 5*time.Minute) // создаём отдельный кэш-слой
	productUC := usecase.NewProductUseCase(productRepo, productCache)
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
