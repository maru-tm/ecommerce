package messaging

import (
	"encoding/json"
	"log"
	"time"

	"inventory-service/config"
	"inventory-service/internal/domain"
)

type OrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderItemEvent struct {
	OrderID    string      `json:"order_id"`
	UserID     string      `json:"user_id"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

func StartConsumer(inventoryUC domain.ProductUseCase) {
	cfg := config.LoadConfig()

	// Используем централизованное подключение
	conn, ch, err := config.ConnectRabbitMQ(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"order_created", // имя очереди
		true,            // durable
		false,           // autoDelete
		false,           // exclusive
		false,           // noWait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf("Ошибка объявления очереди: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // autoAck
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Ошибка регистрации потребителя: %v", err)
	}

	go func() {
		for msg := range msgs {
			var order OrderItemEvent
			if err := json.Unmarshal(msg.Body, &order); err != nil {
				log.Println("Ошибка при разборе сообщения:", err)
				continue
			}

			for _, item := range order.Items {
				if err := inventoryUC.DecreaseStock(item.ProductID, item.Quantity); err != nil {
					log.Printf("Ошибка обновления остатка для товара %s: %v", item.ProductID, err)
				} else {
					log.Printf("Остаток обновлён для товара %s: уменьшено на %d", item.ProductID, item.Quantity)
				}
			}
		}
	}()

	log.Println("Консьюмер запущен и слушает очередь order_created...")
	select {}
}
