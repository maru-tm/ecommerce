package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"user-service/config"
	"user-service/internal/delivery"
	"user-service/internal/infrastructure/redis"
	"user-service/internal/infrastructure/repository"
	"user-service/internal/proto"
	"user-service/internal/usecase"

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

	config.InitMetrics("9101", "user-service")

	ctx := context.Background()
	if err := config.InitRedis(ctx); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}

	redisClient := config.GetRedisClient()
	userCache := redis.NewUserCache(redisClient, 5*time.Minute) // создаём отдельный кэш-слой

	userRepo := repository.NewUserRepository(db)
	userUC := usecase.NewUserUseCase(userRepo, userCache)
	userHandler := delivery.NewUserServiceServer(userUC)

	port := ":50051"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	proto.RegisterUserServiceServer(grpcServer, userHandler)

	fmt.Printf("gRPC server started on %s\n", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
