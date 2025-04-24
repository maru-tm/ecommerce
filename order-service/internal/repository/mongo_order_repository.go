package repository

import (
	"context"
	"fmt"

	"order-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderRepository struct {
	db *mongo.Database
}

func NewOrderRepository(db *mongo.Database) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(order *domain.Order) error {
	collection := r.db.Collection("orders")

	_, err := collection.InsertOne(context.Background(), order)
	if err != nil {
		return fmt.Errorf("failed to create order with ID '%s': %w", order.ID, err)
	}

	return nil
}

func (r *orderRepository) GetOrderByID(id string) (*domain.Order, error) {
	collection := r.db.Collection("orders")

	var order domain.Order
	// Преобразуем строковый ID в ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID format for order ID '%s': %w", id, err)
	}

	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch order with ID '%s': %w", id, err)
	}

	return &order, nil
}

func (r *orderRepository) ListOrders() ([]domain.Order, error) {
	collection := r.db.Collection("orders")
	var orders []domain.Order

	cursor, err := collection.Find(context.Background(), bson.M{}, options.Find())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var order domain.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, fmt.Errorf("failed to decode order: %w", err)
		}
		orders = append(orders, order)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return orders, nil
}

func (r *orderRepository) UpdateOrder(order *domain.Order) error {
	collection := r.db.Collection("orders")
	// Преобразуем строковый ID обратно в ObjectID
	objectID, err := primitive.ObjectIDFromHex(order.ID)
	if err != nil {
		return fmt.Errorf("invalid ObjectID format for order ID '%s': %w", order.ID, err)
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"user_id":     order.UserID,
			"items":       order.Items,
			"total_price": order.TotalPrice,
			"status":      order.Status,
			"created_at":  order.CreatedAt,
			"updated_at":  order.UpdatedAt,
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to update order with ID '%s': %w", order.ID, err)
	}

	return nil
}

func (r *orderRepository) DeleteOrder(id string) error {
	collection := r.db.Collection("orders")
	// Преобразуем строковый ID в ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ObjectID format for order ID '%s': %w", id, err)
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete order with ID '%s': %w", id, err)
	}

	return nil
}
