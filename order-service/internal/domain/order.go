package domain

import (
	"time"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusCompleted OrderStatus = "completed"
	StatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	ProductID string `bson:"product_id" json:"product_id"` // изменено на string
	Quantity  int    `bson:"quantity" json:"quantity"`
}

type Order struct {
	ID         string      `bson:"_id,omitempty" json:"id"` // изменено на string
	UserID     string      `bson:"user_id" json:"user_id"`
	Items      []OrderItem `bson:"items" json:"items"`
	TotalPrice float64     `bson:"total_price" json:"total_price"`
	Status     OrderStatus `bson:"status" json:"status"`
	CreatedAt  time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time   `bson:"updated_at" json:"updated_at"`
}
