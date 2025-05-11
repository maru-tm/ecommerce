// internal/infrastructure/messaging/rabbitmq_order_publisher.go
package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"order-service/config"
	"order-service/internal/domain"

	"github.com/rabbitmq/amqp091-go"
)

type OrderPublisher struct {
	channel *amqp091.Channel
	queue   amqp091.Queue
}

// NewOrderPublisher создает новый OrderPublisher
func NewOrderPublisher(cfg *config.Config) (*OrderPublisher, error) {
	// Подключаемся к RabbitMQ
	_, ch, err := config.ConnectRabbitMQ(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Декларируем очередь
	q, err := ch.QueueDeclare(
		"order_created", // имя очереди
		true,            // durable
		false,           // auto-delete
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	// Возвращаем объект OrderPublisher
	return &OrderPublisher{
		channel: ch,
		queue:   q,
	}, nil
}

// Publish публикует сообщение о создании заказа
func (p *OrderPublisher) Publish(order *domain.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		log.Println("Ошибка сериализации заказа:", err)
		return fmt.Errorf("error serializing order: %v", err)
	}

	// Публикуем сообщение в очередь RabbitMQ
	err = p.channel.Publish(
		"",           // exchange
		p.queue.Name, // routing key (имя очереди)
		false,        // mandatory
		false,        // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Println("Ошибка публикации сообщения:", err)
		return fmt.Errorf("error publishing message: %v", err)
	}

	log.Printf("Заказ опубликован: %+v", order)
	return nil
}
