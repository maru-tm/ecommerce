package main

import (
	"log"

	"api-gateway/internal/clients"
	"api-gateway/internal/middleware"
	"api-gateway/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	userClient, err := clients.NewUserClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create user client: %v", err)
	}

	orderClient, err := clients.NewOrderClient("localhost:50052")
	if err != nil {
		log.Fatalf("Failed to create order client: %v", err)
	}

	inventoryClient, err := clients.NewProductClient("localhost:50053")
	if err != nil {
		log.Fatalf("Failed to create inventory client: %v", err)
	}

	r := gin.Default()

	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.AuthMiddleware())

	routes.RegisterUserRoutes(r, userClient)
	routes.RegisterOrderRoutes(r, orderClient)
	routes.RegisterProductRoutes(r, inventoryClient)

	if err := r.Run(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
