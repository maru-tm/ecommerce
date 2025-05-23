package usecase_test

import (
	"context"
	"errors"
	"inventory-service/internal/domain"
	"inventory-service/internal/usecase"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockProductRepo struct {
	mock.Mock
}

func (m *mockProductRepo) GetProductByName(name string) (*domain.Product, error) {
	args := m.Called(name)
	if p, ok := args.Get(0).(*domain.Product); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProductRepo) CreateProduct(product *domain.Product) error {
	return m.Called(product).Error(0)
}

func (m *mockProductRepo) GetProductByID(id string) (*domain.Product, error) {
	args := m.Called(id)
	if p, ok := args.Get(0).(*domain.Product); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProductRepo) ListProducts() ([]domain.Product, error) {
	args := m.Called()
	if products, ok := args.Get(0).([]domain.Product); ok {
		return products, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProductRepo) UpdateProduct(product *domain.Product) error {
	return m.Called(product).Error(0)
}

func (m *mockProductRepo) DeleteProduct(id string) error {
	return m.Called(id).Error(0)
}

func (m *mockProductRepo) CheckStock(productID string) (bool, error) {
	args := m.Called(productID)
	return args.Bool(0), args.Error(1)
}

func (m *mockProductRepo) SearchProducts(query, categoryID string) ([]domain.Product, error) {
	args := m.Called(query, categoryID)
	if products, ok := args.Get(0).([]domain.Product); ok {
		return products, args.Error(1)
	}
	return nil, args.Error(1)
}

type mockProductCache struct {
	mock.Mock
}

func (m *mockProductCache) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if p, ok := args.Get(0).(*domain.Product); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProductCache) SetProduct(ctx context.Context, product *domain.Product) error {
	return m.Called(ctx, product).Error(0)
}

func (m *mockProductCache) DeleteProduct(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockProductCache) GetProducts(ctx context.Context, key string) ([]domain.Product, error) {
	args := m.Called(ctx, key)
	if products, ok := args.Get(0).([]domain.Product); ok {
		return products, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProductCache) SetProducts(ctx context.Context, key string, products []domain.Product) error {
	return m.Called(ctx, key, products).Error(0)
}

// --- Тесты ---

func createValidProduct() *domain.Product {
	return &domain.Product{
		ID:   "",
		Name: "Test Product",
		Category: domain.Category{
			ID:   "cat1",
			Name: "Category 1",
		},
		Price: 10.5,
		Stock: 100,
	}
}

func TestCreateProduct_Success(t *testing.T) {
	repo := new(mockProductRepo)
	cache := new(mockProductCache)

	uc := usecase.NewProductUseCase(repo, cache)

	product := createValidProduct()

	// Не существует с таким именем
	repo.On("GetProductByName", product.Name).Return(nil, nil)
	// Создание успешное
	repo.On("CreateProduct", mock.AnythingOfType("*domain.Product")).Return(nil)

	err := uc.CreateProduct(product)
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)

	repo.AssertExpectations(t)
}

func TestCreateProduct_ValidationFail(t *testing.T) {
	repo := new(mockProductRepo)
	cache := new(mockProductCache)
	uc := usecase.NewProductUseCase(repo, cache)

	// Пустое имя — ошибка валидации
	product := createValidProduct()
	product.Name = ""

	err := uc.CreateProduct(product)
	assert.ErrorContains(t, err, "product name cannot be empty")
}

func TestCreateProduct_AlreadyExists(t *testing.T) {
	repo := new(mockProductRepo)
	cache := new(mockProductCache)
	uc := usecase.NewProductUseCase(repo, cache)

	product := createValidProduct()

	repo.On("GetProductByName", product.Name).Return(product, nil)

	err := uc.CreateProduct(product)
	assert.ErrorContains(t, err, "already exists")
}

func TestGetProductByID_CacheHit(t *testing.T) {
	repo := new(mockProductRepo)
	cache := new(mockProductCache)
	uc := usecase.NewProductUseCase(repo, cache)

	product := createValidProduct()
	product.ID = uuid.New().String()

	cache.On("GetProduct", mock.Anything, product.ID).Return(product, nil)

	p, err := uc.GetProductByID(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, product, p)
}

func TestGetProductByID_CacheMiss_DBHit(t *testing.T) {
	repo := new(mockProductRepo)
	cache := new(mockProductCache)
	uc := usecase.NewProductUseCase(repo, cache)

	product := createValidProduct()
	product.ID = uuid.New().String()

	cache.On("GetProduct", mock.Anything, product.ID).Return(nil, errors.New("cache miss"))
	repo.On("GetProductByID", product.ID).Return(product, nil)
	cache.On("SetProduct", mock.Anything, product).Return(nil)

	p, err := uc.GetProductByID(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, product, p)
}

func TestListProducts_CacheHit(t *testing.T) {
	repo := new(mockProductRepo)
	cache := new(mockProductCache)
	uc := usecase.NewProductUseCase(repo, cache)

	products := []domain.Product{*createValidProduct()}
	cache.On("GetProducts", mock.Anything, "all_products").Return(products, nil)

	got, err := uc.ListProducts()
	assert.NoError(t, err)
	assert.Equal(t, products, got)
}

func TestListProducts_CacheMiss(t *testing.T) {
	repo := new(mockProductRepo)
	cache := new(mockProductCache)
	uc := usecase.NewProductUseCase(repo, cache)

	products := []domain.Product{*createValidProduct()}

	cache.On("GetProducts", mock.Anything, "all_products").Return(nil, errors.New("cache miss"))
	repo.On("ListProducts").Return(products, nil)
	cache.On("SetProducts", mock.Anything, "all_products", products).Return(nil)

	got, err := uc.ListProducts()
	assert.NoError(t, err)
	assert.Equal(t, products, got)
}

func TestUpdateProduct_Success(t *testing.T) {
	repo := new(mockProductRepo)
	cache := new(mockProductCache)
	uc := usecase.NewProductUseCase(repo, cache)

	product := createValidProduct()
	product.ID = uuid.New().String()

	repo.On("GetProductByID", product.ID).Return(product, nil)
	repo.On("UpdateProduct", product).Return(nil)
	cache.On("DeleteProduct", mock.Anything, product.ID).Return(nil)
	cache.On("SetProduct", mock.Anything, product).Return(nil)
	cache.On("DeleteProduct", mock.Anything, "all_products").Return(nil)

	err := uc.UpdateProduct(product.ID, product)
	assert.NoError(t, err)
}

func TestDecreaseStock_Success(t *testing.T) {
	repo := new(mockProductRepo)
	cache := new(mockProductCache)
	uc := usecase.NewProductUseCase(repo, cache)

	product := createValidProduct()
	product.ID = uuid.New().String()
	product.Stock = 10

	repo.On("GetProductByID", product.ID).Return(product, nil)
	repo.On("UpdateProduct", mock.MatchedBy(func(p *domain.Product) bool {
		return p.Stock == 7
	})).Return(nil)
	cache.On("DeleteProduct", mock.Anything, product.ID).Return(nil)
	cache.On("DeleteProduct", mock.Anything, "all_products").Return(nil)

	err := uc.DecreaseStock(product.ID, 3)
	assert.NoError(t, err)
}

func TestDeleteProduct_Success(t *testing.T) {
	repo := new(mockProductRepo)
	cache := new(mockProductCache)
	uc := usecase.NewProductUseCase(repo, cache)

	productID := uuid.New().String()
	product := createValidProduct()
	product.ID = productID

	repo.On("GetProductByID", productID).Return(product, nil)
	repo.On("DeleteProduct", productID).Return(nil)
	cache.On("DeleteProduct", mock.Anything, productID).Return(nil)
	cache.On("DeleteProduct", mock.Anything, "all_products").Return(nil)

	err := uc.DeleteProduct(productID)
	assert.NoError(t, err)
}

func TestCheckStock(t *testing.T) {
	repo := new(mockProductRepo)
	cache := new(mockProductCache)
	uc := usecase.NewProductUseCase(repo, cache)

	productID := uuid.New().String()

	repo.On("CheckStock", productID).Return(true, nil)

	ok, err := uc.CheckStock(productID, 5)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestSearchProducts(t *testing.T) {
	repo := new(mockProductRepo)
	cache := new(mockProductCache)
	uc := usecase.NewProductUseCase(repo, cache)

	products := []domain.Product{*createValidProduct()}
	query := "test"
	categoryID := "cat1"

	repo.On("SearchProducts", query, categoryID).Return(products, nil)

	got, err := uc.SearchProducts(query, categoryID)
	assert.NoError(t, err)
	assert.Equal(t, products, got)
}

func TestUpdateProductStock_Success(t *testing.T) {
	repo := new(mockProductRepo)
	cache := new(mockProductCache)
	uc := usecase.NewProductUseCase(repo, cache)

	product := createValidProduct()
	product.ID = uuid.New().String()

	repo.On("GetProductByID", product.ID).Return(product, nil)
	repo.On("UpdateProduct", product).Return(nil)
	cache.On("DeleteProduct", mock.Anything, product.ID).Return(nil)
	cache.On("DeleteProduct", mock.Anything, "all_products").Return(nil)

	err := uc.UpdateProductStock(product.ID, 50)
	assert.NoError(t, err)
}
