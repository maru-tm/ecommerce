package usecase

import (
	"fmt"
	"inventory-service/internal/domain" // импортируем пакет с протобуфами
	"inventory-service/internal/repository"
	"log"

	"github.com/google/uuid"
)

type ProductUseCase interface {
	CreateProduct(product *domain.Product) error
	GetProductByID(id string) (*domain.Product, error)
	ListProducts() ([]domain.Product, error)
	UpdateProduct(id string, product *domain.Product) error
	DecreaseStock(productID string, quantity int) error
	DeleteProduct(id string) error
	CheckStock(productID string, quantity int) (bool, error)
	SearchProducts(query string, categoryID string) ([]domain.Product, error) // Новый метод поиска продуктов
	UpdateProductStock(productID string, quantity int) error                  // Новый метод для обновления количества товара
}

type productUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return &productUseCase{repo: repo}
}

func (uc *productUseCase) validateProduct(product *domain.Product) error {
	log.Printf("Validating product: %+v", product)

	if product.Name == "" {
		return fmt.Errorf("product name cannot be empty")
	}
	if product.Category.ID == "" {
		return fmt.Errorf("product category ID cannot be empty")
	}
	if product.Category.Name == "" {
		return fmt.Errorf("product category name cannot be empty")
	}
	if product.Price <= 0 {
		return fmt.Errorf("product price must be positive")
	}
	if product.Stock < 0 {
		return fmt.Errorf("product stock cannot be negative")
	}
	return nil
}

func (uc *productUseCase) CreateProduct(product *domain.Product) error {
	log.Println("Creating product...")
	if err := uc.validateProduct(product); err != nil {
		log.Printf("Validation failed: %v", err)
		return err
	}

	existingProduct, err := uc.repo.GetProductByName(product.Name)
	if err != nil {
		log.Printf("Error checking product uniqueness: %v", err)
		return fmt.Errorf("error checking product uniqueness: %w", err)
	}
	if existingProduct != nil {
		log.Printf("Product with name '%s' already exists", product.Name)
		return fmt.Errorf("product with name '%s' already exists", product.Name)
	}

	if product.ID == "" {
		product.ID = uuid.New().String()
		log.Printf("Generated new UUID for product: %s", product.ID)
	}

	err = uc.repo.CreateProduct(product)
	if err != nil {
		log.Printf("Failed to create product: %v", err)
	} else {
		log.Printf("Product created successfully: %+v", product)
	}
	return err
}

func (uc *productUseCase) GetProductByID(id string) (*domain.Product, error) {
	log.Printf("Getting product by ID: %s", id)
	product, err := uc.repo.GetProductByID(id)
	if err != nil {
		log.Printf("Error getting product: %v", err)
	}
	return product, err
}

func (uc *productUseCase) ListProducts() ([]domain.Product, error) {
	log.Println("Listing all products...")
	products, err := uc.repo.ListProducts()
	if err != nil {
		log.Printf("Error listing products: %v", err)
	}
	return products, err
}

func (uc *productUseCase) UpdateProduct(id string, product *domain.Product) error {
	log.Printf("Updating product ID: %s", id)
	if id == "" {
		return fmt.Errorf("product ID cannot be empty")
	}
	product.ID = id

	existingProduct, err := uc.repo.GetProductByID(id)
	if err != nil {
		log.Printf("Error checking product existence: %v", err)
		return fmt.Errorf("error checking product existence: %w", err)
	}
	if existingProduct == nil {
		log.Printf("Product with ID '%s' not found", id)
		return fmt.Errorf("product with ID '%s' not found", id)
	}

	if err := uc.validateProduct(product); err != nil {
		log.Printf("Validation failed: %v", err)
		return err
	}

	err = uc.repo.UpdateProduct(product)
	if err != nil {
		log.Printf("Failed to update product: %v", err)
	} else {
		log.Printf("Product updated successfully: %+v", product)
	}
	return err
}

func (uc *productUseCase) DecreaseStock(productID string, quantity int) error {
	log.Printf("Decreasing stock for product ID: %s by %d", productID, quantity)

	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	product, err := uc.repo.GetProductByID(productID)
	if err != nil {
		log.Printf("Error fetching product: %v", err)
		return fmt.Errorf("failed to fetch product: %w", err)
	}
	if product == nil {
		log.Printf("Product with ID '%s' not found", productID)
		return fmt.Errorf("product with ID '%s' not found", productID)
	}

	if product.Stock < quantity {
		log.Printf("Not enough stock: available %d, requested %d", product.Stock, quantity)
		return fmt.Errorf("not enough stock for product '%s'", productID)
	}

	product.Stock -= quantity
	err = uc.repo.UpdateProduct(product)
	if err != nil {
		log.Printf("Failed to update product stock: %v", err)
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	log.Printf("Stock decreased successfully. New stock: %d", product.Stock)
	return nil
}

func (uc *productUseCase) DeleteProduct(id string) error {
	log.Printf("Deleting product ID: %s", id)
	existingProduct, err := uc.repo.GetProductByID(id)
	if err != nil {
		log.Printf("Error checking product existence: %v", err)
		return fmt.Errorf("error checking product existence: %w", err)
	}
	if existingProduct == nil {
		log.Printf("Product with ID '%s' not found", id)
		return fmt.Errorf("product with ID '%s' not found", id)
	}

	err = uc.repo.DeleteProduct(id)
	if err != nil {
		log.Printf("Failed to delete product: %v", err)
	} else {
		log.Printf("Product with ID '%s' deleted successfully", id)
	}
	return err
}

func (uc *productUseCase) CheckStock(productID string, quantity int) (bool, error) {
	log.Printf("Checking stock for product ID: %s with quantity: %d", productID, quantity)

	available, err := uc.repo.CheckStock(productID)
	if err != nil {
		log.Printf("Error checking stock: %v", err)
		return false, fmt.Errorf("failed to check stock: %w", err)
	}

	return available, nil
}

func (uc *productUseCase) SearchProducts(query string, categoryID string) ([]domain.Product, error) {
	log.Printf("Searching products with query: %s and category ID: %s", query, categoryID)

	// Выполним поиск в репозитории
	products, err := uc.repo.SearchProducts(query, categoryID)
	if err != nil {
		log.Printf("Error searching products: %v", err)
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	log.Printf("Found %d products matching the search", len(products))
	return products, nil
}

func (uc *productUseCase) UpdateProductStock(productID string, quantity int) error {
	log.Printf("Updating stock for product ID: %s by %d", productID, quantity)

	// Получаем продукт из репозитория
	product, err := uc.repo.GetProductByID(productID)
	if err != nil {
		log.Printf("Error fetching product: %v", err)
		return fmt.Errorf("failed to fetch product: %w", err)
	}
	if product == nil {
		log.Printf("Product with ID '%s' not found", productID)
		return fmt.Errorf("product with ID '%s' not found", productID)
	}

	// Обновляем количество товара
	product.Stock = quantity
	err = uc.repo.UpdateProduct(product)
	if err != nil {
		log.Printf("Failed to update product stock: %v", err)
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	log.Printf("Product stock updated successfully to %d", quantity)
	return nil
}
