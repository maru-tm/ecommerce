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
	Database            string
	RabbitMQURL         string
	InventoryServerPort string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Не удалось загрузить .env файл, программа завершена!")
	}

	return &Config{
		MongoURI:            getEnv("MONGO_URI", "mongodb://mongo:27017"),
		Database:            getEnv("MONGO_DB", "ecommerce_products"),
		RabbitMQURL:         getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		InventoryServerPort: getEnv("INVENTORY_SERVER_PORT", "50053"), // Читаем порт из переменной окружения

	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func ConnectDB(cfg *Config) *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Ошибка создания клиента MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Ошибка проверки соединения с MongoDB: %v", err)
	}

	fmt.Println("Подключение к MongoDB установлено!")
	return client.Database(cfg.Database)
}

// ConnectRabbitMQ подключается к RabbitMQ
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
