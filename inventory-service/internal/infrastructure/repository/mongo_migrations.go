package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RunMigrations(ctx context.Context, db *mongo.Database) error {
	// Пример: создаем коллекцию products, если ее нет
	collectionNames, err := db.ListCollectionNames(ctx, struct{}{})
	if err != nil {
		return err
	}

	// Проверяем, есть ли коллекция "products"
	found := false
	for _, name := range collectionNames {
		if name == "products" {
			found = true
			break
		}
	}

	if !found {
		// Создаем коллекцию
		err = db.CreateCollection(ctx, "products")
		if err != nil {
			return err
		}
		log.Println("Коллекция 'products' создана")
	}

	// Создаем индекс, например, уникальный по полю "sku"
	productsColl := db.Collection("products")

	indexModel := mongo.IndexModel{
		Keys:    map[string]int{"sku": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err = productsColl.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}
	log.Println("Индекс для 'sku' создан или уже существует")

	// Добавь сюда другие миграции, если нужно

	return nil
}
