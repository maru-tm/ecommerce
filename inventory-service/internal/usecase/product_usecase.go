package usecase

import (
	"context"
	"fmt"
	"inventory-service/internal/domain"
	"log"

	"github.com/google/uuid"
)

type productUseCase struct {
	repo  domain.ProductRepository
	cache domain.ProductCache
}

func NewProductUseCase(repo domain.ProductRepository, cache domain.ProductCache) domain.ProductUseCase {
	return &productUseCase{
		repo:  repo,
		cache: cache,
	}
}

func (uc *productUseCase) validateProduct(product *domain.Product) error {
	log.Printf("[VALIDATE] Validating product: %+v", product)

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
	log.Println("[CREATE] Creating product...")

	if err := uc.validateProduct(product); err != nil {
		log.Printf("[CREATE] Validation failed: %v", err)
		return err
	}

	existingProduct, err := uc.repo.GetProductByName(product.Name)
	if err != nil {
		log.Printf("[CREATE] Error checking product uniqueness: %v", err)
		return fmt.Errorf("error checking product uniqueness: %w", err)
	}
	if existingProduct != nil {
		log.Printf("[CREATE] Product with name '%s' already exists", product.Name)
		return fmt.Errorf("product with name '%s' already exists", product.Name)
	}

	if product.ID == "" {
		product.ID = uuid.New().String()
		log.Printf("[CREATE] Generated new UUID for product: %s", product.ID)
	}

	if err := uc.repo.CreateProduct(product); err != nil {
		log.Printf("[CREATE] Failed to create product: %v", err)
		return err
	}

	log.Printf("[CREATE] Product created successfully: %+v", product)
	return nil
}

func (uc *productUseCase) GetProductByID(id string) (*domain.Product, error) {
	ctx := context.Background()

	if p, _ := uc.cache.GetProduct(ctx, id); p != nil {
		log.Printf("[GET] Product found in cache: %s", id)
		return p, nil
	}

	log.Printf("[GET] Product not found in cache. Fetching from DB: %s", id)
	product, err := uc.repo.GetProductByID(id)
	if err != nil {
		log.Printf("[GET] Error fetching product from DB: %v", err)
		return nil, err
	}
	if product == nil {
		log.Printf("[GET] Product with ID '%s' not found in DB", id)
		return nil, nil
	}

	if err := uc.cache.SetProduct(ctx, product); err != nil {
		log.Printf("[CACHE] Failed to set product in cache: %v", err)
	} else {
		log.Printf("[CACHE] Product cached successfully: %s", id)
	}

	return product, nil
}

func (uc *productUseCase) ListProducts() ([]domain.Product, error) {
	log.Println("[LIST] Fetching all products...")

	ctx := context.Background()
	cacheKey := "all_products"

	products, err := uc.cache.GetProducts(ctx, cacheKey)
	if err != nil {
		log.Printf("[LIST] Error retrieving products from cache: %v", err)
	}
	if products != nil {
		log.Println("[LIST] Products loaded from cache")
		return products, nil
	}

	log.Println("[LIST] Cache is empty. Fetching from repository...")
	products, err = uc.repo.ListProducts()
	if err != nil {
		log.Printf("[LIST] Error fetching products from repository: %v", err)
		return nil, err
	}

	if err := uc.cache.SetProducts(ctx, cacheKey, products); err != nil {
		log.Printf("[LIST] Failed to cache products: %v", err)
	} else {
		log.Printf("[LIST] Products cached successfully. Key: %s, Count: %d", cacheKey, len(products))
	}

	return products, nil
}

func (uc *productUseCase) UpdateProduct(id string, product *domain.Product) error {
	log.Printf("[UPDATE] Updating product: %s", id)

	if id == "" {
		return fmt.Errorf("product ID cannot be empty")
	}
	product.ID = id

	existingProduct, err := uc.repo.GetProductByID(id)
	if err != nil {
		log.Printf("[UPDATE] Error checking product existence: %v", err)
		return fmt.Errorf("error checking product existence: %w", err)
	}
	if existingProduct == nil {
		log.Printf("[UPDATE] Product not found: %s", id)
		return fmt.Errorf("product with ID '%s' not found", id)
	}

	if err := uc.validateProduct(product); err != nil {
		log.Printf("[UPDATE] Validation failed: %v", err)
		return err
	}

	if err := uc.repo.UpdateProduct(product); err != nil {
		log.Printf("[UPDATE] Failed to update product: %v", err)
		return err
	}

	log.Printf("[UPDATE] Product updated successfully: %+v", product)
	return nil
}

func (uc *productUseCase) DecreaseStock(productID string, quantity int) error {
	log.Printf("[STOCK] Decreasing stock for product: %s by %d", productID, quantity)

	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	product, err := uc.repo.GetProductByID(productID)
	if err != nil {
		log.Printf("[STOCK] Error fetching product: %v", err)
		return fmt.Errorf("failed to fetch product: %w", err)
	}
	if product == nil {
		log.Printf("[STOCK] Product not found: %s", productID)
		return fmt.Errorf("product with ID '%s' not found", productID)
	}

	if product.Stock < quantity {
		log.Printf("[STOCK] Not enough stock. Available: %d, Requested: %d", product.Stock, quantity)
		return fmt.Errorf("not enough stock for product '%s'", productID)
	}

	product.Stock -= quantity
	if err := uc.repo.UpdateProduct(product); err != nil {
		log.Printf("[STOCK] Failed to update stock: %v", err)
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	log.Printf("[STOCK] Stock updated successfully. New stock: %d", product.Stock)
	return nil
}

func (uc *productUseCase) DeleteProduct(id string) error {
	log.Printf("[DELETE] Deleting product: %s", id)

	existingProduct, err := uc.repo.GetProductByID(id)
	if err != nil {
		log.Printf("[DELETE] Error checking product existence: %v", err)
		return fmt.Errorf("error checking product existence: %w", err)
	}
	if existingProduct == nil {
		log.Printf("[DELETE] Product not found: %s", id)
		return fmt.Errorf("product with ID '%s' not found", id)
	}

	if err := uc.repo.DeleteProduct(id); err != nil {
		log.Printf("[DELETE] Failed to delete product: %v", err)
		return err
	}

	log.Printf("[DELETE] Product deleted successfully: %s", id)
	return nil
}

func (uc *productUseCase) CheckStock(productID string, quantity int) (bool, error) {
	log.Printf("[CHECK] Checking stock for product: %s, Quantity: %d", productID, quantity)

	available, err := uc.repo.CheckStock(productID)
	if err != nil {
		log.Printf("[CHECK] Error checking stock: %v", err)
		return false, fmt.Errorf("failed to check stock: %w", err)
	}

	return available, nil
}

func (uc *productUseCase) SearchProducts(query string, categoryID string) ([]domain.Product, error) {
	log.Printf("[SEARCH] Searching products. Query: '%s', Category ID: '%s'", query, categoryID)

	products, err := uc.repo.SearchProducts(query, categoryID)
	if err != nil {
		log.Printf("[SEARCH] Error searching products: %v", err)
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	log.Printf("[SEARCH] Found %d products matching criteria", len(products))
	return products, nil
}

func (uc *productUseCase) UpdateProductStock(productID string, quantity int) error {
	log.Printf("[STOCK] Updating stock for product: %s to %d", productID, quantity)

	product, err := uc.repo.GetProductByID(productID)
	if err != nil {
		log.Printf("[STOCK] Error fetching product: %v", err)
		return fmt.Errorf("failed to fetch product: %w", err)
	}
	if product == nil {
		log.Printf("[STOCK] Product not found: %s", productID)
		return fmt.Errorf("product with ID '%s' not found", productID)
	}

	product.Stock = quantity
	if err := uc.repo.UpdateProduct(product); err != nil {
		log.Printf("[STOCK] Failed to update product stock: %v", err)
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	log.Printf("[STOCK] Product stock updated successfully to %d", quantity)
	return nil
}
