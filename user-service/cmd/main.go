package main

import (
	"fmt"
	"log"
	"net"

	"user-service/config"
	"user-service/internal/delivery"
	"user-service/internal/proto"
	"user-service/internal/repository"
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

	if err := config.InitRedis(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userUC := usecase.NewUserUseCase(userRepo)
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
