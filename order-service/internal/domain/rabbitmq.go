package domain

type OrderEventPublisher interface {
	Publish(order *Order) error
}
