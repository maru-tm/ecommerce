package repository

import (
	"context"
	"fmt"
	"inventory-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productRepository struct {
	db *mongo.Database
}

func NewProductRepository(db *mongo.Database) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product *domain.Product) error {
	collection := r.db.Collection("products")

	_, err := collection.InsertOne(context.Background(), product)
	if err != nil {
		return fmt.Errorf("failed to create product '%s': %w", product.Name, err)
	}

	return nil
}

func (r *productRepository) GetProductByID(id string) (*domain.Product, error) {
	collection := r.db.Collection("products")

	var product domain.Product
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch product with ID '%s': %w", id, err)
	}

	return &product, nil
}

func (r *productRepository) GetProductByName(name string) (*domain.Product, error) {
	collection := r.db.Collection("products")

	var product domain.Product

	err := collection.FindOne(context.Background(), bson.M{"name": name}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch product with name '%s': %w", name, err)
	}

	return &product, nil
}

func (r *productRepository) ListProducts() ([]domain.Product, error) {
	collection := r.db.Collection("products")
	var products []domain.Product

	cursor, err := collection.Find(context.Background(), bson.M{}, options.Find())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, fmt.Errorf("failed to decode product: %w", err)
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return products, nil
}

func (r *productRepository) UpdateProduct(product *domain.Product) error {
	collection := r.db.Collection("products")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": product.ID},
		bson.M{"$set": bson.M{
			"name":        product.Name,
			"category":    product.Category,
			"price":       product.Price,
			"stock":       product.Stock,
			"description": product.Description,
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to update product '%s': %w", product.ID, err)
	}

	return nil
}

func (r *productRepository) DeleteProduct(id string) error {
	collection := r.db.Collection("products")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete product with ID '%s': %w", id, err)
	}

	return nil
}
