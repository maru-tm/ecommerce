package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Category struct {
	ID   string `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name"`
}

type Product struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Category    Category  `bson:"category" json:"category"`
	Price       float64   `bson:"price" json:"price"`
	Stock       int       `bson:"stock" json:"stock"`
	Description string    `bson:"description" json:"description"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

func getAllProducts(url string) ([]Product, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var products []Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return products, nil
}

func patchProductStock(productID string, newStock int) error {
	// Получаем текущий продукт
	getURL := fmt.Sprintf("http://localhost:8080/products/%s", productID)
	resp, err := http.Get(getURL)
	if err != nil {
		return fmt.Errorf("failed to get product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch product: status %d", resp.StatusCode)
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return fmt.Errorf("failed to decode product: %v", err)
	}

	// Обновляем только stock
	product.Stock = newStock
	product.UpdatedAt = time.Now()

	// Отправляем PATCH с полной структурой
	url := fmt.Sprintf("http://localhost:8080/products/%s", productID)

	jsonData, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal updated product: %v", err)
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create patch request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("patch request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("patch request returned status: %d", resp.StatusCode)
	}

	return nil
}
