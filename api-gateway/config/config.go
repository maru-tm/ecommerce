package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	UserServiceAddr    string
	OrderServiceAddr   string
	ProductServiceAddr string
	APIGatewayPort     string
	RabbitMQURL        string
}

func LoadConfig() *Config {
	// Загружаем .env файл (опционально)
	if err := godotenv.Load(); err != nil {
		log.Println(".env файл не найден — используются переменные окружения или значения по умолчанию")
	}

	return &Config{
		UserServiceAddr:    getEnv("USER_SERVICE_ADDR", "user-service:50051"),
		OrderServiceAddr:   getEnv("ORDER_SERVICE_ADDR", "order-service:50054"),
		ProductServiceAddr: getEnv("PRODUCT_SERVICE_ADDR", "inventory-service:50053"),
		APIGatewayPort:     getEnv("API_GATEWAY_PORT", "8000"),
		RabbitMQURL:        getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
