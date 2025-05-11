package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoURI            string
	MongoDatabase       string
	RabbitMQURL         string
	OrderServicePort    string
	InventoryServiceURL string
}

// LoadConfig загружает конфигурацию из .env файла или переменных окружения
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env файл не найден, используются переменные окружения")
	}

	return &Config{
		MongoURI:            getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:       getEnv("MONGO_DB", "ecommerce_orders"),
		RabbitMQURL:         getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		OrderServicePort:    getEnv("ORDER_SERVICE_PORT", "order-service:50052"),
		InventoryServiceURL: getEnv("INVENTORY_SERVICE_URL", "inventory-service:50053"),
	}
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// ConnectDB устанавливает соединение с MongoDB
func ConnectDB(cfg *Config) *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("❌ Не удалось создать клиент MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("❌ Ошибка подключения к MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ MongoDB не отвечает: %v", err)
	}

	fmt.Println("✅ Подключение к MongoDB успешно установлено")
	return client.Database(cfg.MongoDatabase)
}

// ConnectRabbitMQ устанавливает соединение и канал с RabbitMQ
func ConnectRabbitMQ(cfg *Config) (*amqp091.Connection, *amqp091.Channel, error) {
	var conn *amqp091.Connection
	var ch *amqp091.Channel
	var err error

	retries := 10

	for i := 0; i < retries; i++ {
		conn, err = amqp091.Dial(cfg.RabbitMQURL)
		if err == nil {
			ch, err = conn.Channel()
			if err == nil {
				log.Println("✅ Успешное подключение к RabbitMQ")
				return conn, ch, nil
			}
			conn.Close()
		}

		log.Printf("⏳ Попытка подключения к RabbitMQ (%d/%d): %v", i+1, retries, err)
		time.Sleep(10 * time.Second)
	}

	return nil, nil, fmt.Errorf("❌ Ошибка подключения к RabbitMQ после %d попыток: %v", retries, err)
}
