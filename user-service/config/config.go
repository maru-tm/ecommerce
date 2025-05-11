package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseConfig struct {
	MongoURI string
	Database string
}

// LoadConfig загружает переменные окружения из .env файла и возвращает конфигурацию для базы данных.
func LoadConfig() *DatabaseConfig {
	// Загружаем .env файл
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл, используются переменные окружения")
	}

	// Возвращаем конфигурацию с переменными окружения или значениями по умолчанию
	return &DatabaseConfig{
		MongoURI: getEnv("MONGO_URI", "mongodb://mongo:27017"), // Если переменная окружения не найдена, используется значение по умолчанию
		Database: getEnv("MONGO_DB", "ecommerce_users"),
	}
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// ConnectDB устанавливает подключение к MongoDB с использованием конфигурации.
func ConnectDB(cfg *DatabaseConfig) *mongo.Database {
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
		log.Fatalf("Не удалось выполнить ping MongoDB: %v", err)
	}

	fmt.Println("Подключение к MongoDB установлено!")
	return client.Database(cfg.Database)
}
