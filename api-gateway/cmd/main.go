package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"api-gateway/config"
	"api-gateway/internal/clients"
	"api-gateway/internal/email"
	"api-gateway/internal/middleware"
	"api-gateway/internal/routes"
)

func main() {
	cfg := config.LoadConfig()
	mailer := email.NewMailer(cfg)

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

	// Включаем CORS — разрешаем запросы с фронтенда (localhost:3000)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middleware.LoggerMiddleware())
	// r.Use(middleware.AuthMiddleware()) // Пока отключено

	r.GET("/test-cors", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "CORS работает"})
	})

	routes.RegisterUserRoutes(r, userClient, mailer)
	routes.RegisterOrderRoutes(r, orderClient)
	routes.RegisterProductRoutes(r, inventoryClient)

	if err := r.Run(":" + cfg.APIGatewayPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
