package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"api-gateway/config"
	"api-gateway/internal/clients"
	"api-gateway/internal/middleware"
	"api-gateway/internal/routes"
)

func main() {
	cfg := config.LoadConfig()

	userClient, err := clients.NewUserClient(cfg.UserServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create user client: %v", err)
	}

	orderClient, err := clients.NewOrderClient(cfg.OrderServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create order client: %v", err)
	}

	inventoryClient, err := clients.NewProductClient(cfg.ProductServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create product client: %v", err)
	}

	r := gin.Default()

	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.AuthMiddleware())

	routes.RegisterUserRoutes(r, userClient)
	routes.RegisterOrderRoutes(r, orderClient)
	routes.RegisterProductRoutes(r, inventoryClient)

	if err := r.Run(":" + cfg.APIGatewayPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
