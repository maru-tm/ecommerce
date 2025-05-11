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
func (r *productRepository) CheckStock(productID string) (bool, error) {
	collection := r.db.Collection("products")

	// Ищем товар по ID
	var product domain.Product
	err := collection.FindOne(context.Background(), bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil // Если товар не найден, возвращаем false
		}
		return false, fmt.Errorf("failed to fetch product with ID '%s': %w", productID, err)
	}

	// Проверяем наличие товара на складе
	if product.Stock > 0 {
		return true, nil // Если в наличии, возвращаем true
	}

	return false, nil // Если товар отсутствует, возвращаем false
}

func (r *productRepository) SearchProducts(query string, categoryID string) ([]domain.Product, error) {
	collection := r.db.Collection("products")

	filter := bson.M{
		"$and": []bson.M{},
	}

	// Добавляем фильтр по ключевым словам, если есть
	if query != "" {
		filter["$and"] = append(filter["$and"].([]bson.M), bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": query, "$options": "i"}},
				{"description": bson.M{"$regex": query, "$options": "i"}},
			},
		})
	}

	// Добавляем фильтр по categoryID, если передан
	if categoryID != "" {
		filter["$and"] = append(filter["$and"].([]bson.M), bson.M{
			"category.id": categoryID,
		})
	}

	// Если фильтров нет, ищем все
	if len(filter["$and"].([]bson.M)) == 0 {
		filter = bson.M{}
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer cursor.Close(context.Background())

	var results []domain.Product
	for cursor.Next(context.Background()) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, fmt.Errorf("failed to decode product: %w", err)
		}
		results = append(results, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return results, nil
}
