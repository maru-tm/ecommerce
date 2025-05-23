package config

import (
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

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
