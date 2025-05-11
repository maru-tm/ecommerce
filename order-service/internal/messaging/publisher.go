package messaging

// import (
// 	"encoding/json"
// 	"log"
// 	"order-service/config" // импортируем config для доступа к конфигурации
// 	"order-service/internal/domain"

// 	"github.com/rabbitmq/amqp091-go"
// )

// type OrderPublisher struct {
// 	channel *amqp091.Channel
// 	queue   amqp091.Queue
// }

// // NewOrderPublisher создаёт нового OrderPublisher, подключается к RabbitMQ и настраивает очередь
// func NewOrderPublisher(cfg *config.Config) (*OrderPublisher, error) {
// 	// Используем функцию из config.go для подключения к RabbitMQ
// 	_, ch, err := config.ConnectRabbitMQ(cfg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Объявляем очередь для сообщений
// 	q, err := ch.QueueDeclare(
// 		"order_created", // имя очереди
// 		true,            // durable
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Возвращаем объект OrderPublisher с каналом и очередью
// 	return &OrderPublisher{channel: ch, queue: q}, nil
// }

// // Publish публикует заказ в очередь
// func (p *OrderPublisher) Publish(order *domain.Order) {
// 	body, err := json.Marshal(order)
// 	if err != nil {
// 		log.Println("Ошибка сериализации заказа:", err)
// 		return
// 	}

// 	err = p.channel.Publish(
// 		"", p.queue.Name, false, false,
// 		amqp091.Publishing{
// 			ContentType: "application/json",
// 			Body:        body,
// 		},
// 	)
// 	if err != nil {
// 		log.Println("Ошибка публикации сообщения:", err)
// 	} else {
// 		log.Printf("Заказ опубликован: %+v", order)
// 	}
// }
