package repository

import (
	"context"
	"fmt"
	"log"

	"order-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type orderRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) domain.OrderRepository {
	return &orderRepository{
		db:         db,
		collection: db.Collection("orders"),
	}
}

func (r *orderRepository) CreateOrder(order *domain.Order) error {
	ctx := context.Background()
	log.Printf("[CreateOrder] Start creating order ID=%s", order.ID)

	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		err = fmt.Errorf("failed to create order with ID '%s': %w", order.ID, err)
		log.Printf("[CreateOrder] Error: %v", err)
		return err
	}

	log.Printf("[CreateOrder] Successfully created order ID=%s", order.ID)
	return nil
}

func (r *orderRepository) GetOrderByID(id string) (*domain.Order, error) {
	ctx := context.Background()
	log.Printf("[GetOrderByID] Fetching order ID=%s", id)

	var order domain.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("[GetOrderByID] Order ID=%s not found", id)
			return nil, nil
		}
		err = fmt.Errorf("failed to fetch order with ID '%s': %w", id, err)
		log.Printf("[GetOrderByID] Error: %v", err)
		return nil, err
	}

	log.Printf("[GetOrderByID] Order ID=%s found", id)
	return &order, nil
}

func (r *orderRepository) ListOrders() ([]domain.Order, error) {
	ctx := context.Background()
	log.Println("[ListOrders] Fetching all orders")

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		err = fmt.Errorf("failed to fetch orders: %w", err)
		log.Printf("[ListOrders] Error: %v", err)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("[ListOrders] Cursor close error: %v", err)
		}
	}()

	var orders []domain.Order
	for cursor.Next(ctx) {
		var order domain.Order
		if err := cursor.Decode(&order); err != nil {
			err = fmt.Errorf("failed to decode order: %w", err)
			log.Printf("[ListOrders] Error: %v", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := cursor.Err(); err != nil {
		err = fmt.Errorf("cursor iteration error: %w", err)
		log.Printf("[ListOrders] Error: %v", err)
		return nil, err
	}

	log.Printf("[ListOrders] Fetched %d orders", len(orders))
	return orders, nil
}

func (r *orderRepository) UpdateOrder(order *domain.Order) error {
	ctx := context.Background()
	log.Printf("[UpdateOrder] Updating order ID=%s", order.ID)

	objectID, err := primitive.ObjectIDFromHex(order.ID)
	if err != nil {
		err = fmt.Errorf("invalid ObjectID format for order ID '%s': %w", order.ID, err)
		log.Printf("[UpdateOrder] Error: %v", err)
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"user_id":     order.UserID,
			"items":       order.Items,
			"total_price": order.TotalPrice,
			"status":      order.Status,
			"created_at":  order.CreatedAt,
			"updated_at":  order.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		err = fmt.Errorf("failed to update order with ID '%s': %w", order.ID, err)
		log.Printf("[UpdateOrder] Error: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		log.Printf("[UpdateOrder] No order found with ID=%s", order.ID)
		return fmt.Errorf("no order found with ID %s", order.ID)
	}

	log.Printf("[UpdateOrder] Successfully updated order ID=%s", order.ID)
	return nil
}

func (r *orderRepository) DeleteOrder(id string) error {
	ctx := context.Background()
	log.Printf("[DeleteOrder] Deleting order ID=%s", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		err = fmt.Errorf("invalid ObjectID format for order ID '%s': %w", id, err)
		log.Printf("[DeleteOrder] Error: %v", err)
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		err = fmt.Errorf("failed to delete order with ID '%s': %w", id, err)
		log.Printf("[DeleteOrder] Error: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		log.Printf("[DeleteOrder] No order found to delete with ID=%s", id)
		return fmt.Errorf("no order found to delete with ID %s", id)
	}

	log.Printf("[DeleteOrder] Successfully deleted order ID=%s", id)
	return nil
}
