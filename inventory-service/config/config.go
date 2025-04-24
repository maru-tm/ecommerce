package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoURI string
	Database string
}

func LoadConfig() *Config {
	return &Config{
		MongoURI: "mongodb://localhost:27017",
		// MongoURI: "mongodb+srv://aruzhanduyssenova:pj87dxbb0dtgtohk@cluster0.1rdt5vd.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0",
		Database: "ecommerce_products",
	}
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
