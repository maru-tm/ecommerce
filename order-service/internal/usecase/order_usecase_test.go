package usecase_test

import (
	"context"
	"errors"
	"testing"

	"order-service/internal/domain"
	"order-service/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ==== MOCKS ====

type mockOrderRepo struct{ mock.Mock }

func (m *mockOrderRepo) CreateOrder(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}
func (m *mockOrderRepo) GetOrderByID(id string) (*domain.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Order), args.Error(1)
}
func (m *mockOrderRepo) ListOrders() ([]domain.Order, error) {
	args := m.Called()
	return args.Get(0).([]domain.Order), args.Error(1)
}
func (m *mockOrderRepo) UpdateOrder(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}
func (m *mockOrderRepo) DeleteOrder(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type mockOrderPublisher struct{ mock.Mock }

func (m *mockOrderPublisher) Publish(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

type mockOrderCache struct{ mock.Mock }

func (m *mockOrderCache) SetOrder(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}
func (m *mockOrderCache) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Order), args.Error(1)
}
func (m *mockOrderCache) DeleteOrder(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockOrderCache) SetOrders(ctx context.Context, key string, orders []domain.Order) error {
	args := m.Called(ctx, key, orders)
	return args.Error(0)
}
func (m *mockOrderCache) GetOrders(ctx context.Context, key string) ([]domain.Order, error) {
	args := m.Called(ctx, key)
	return args.Get(0).([]domain.Order), args.Error(1)
}

// ==== TEST ====

func TestCreateOrder_Success(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockPublisher := new(mockOrderPublisher)
	mockCache := new(mockOrderCache)

	order := &domain.Order{
		UserID: "user123",
		Items: []domain.OrderItem{
			{ProductID: "prod123", Quantity: 2},
		},
		TotalPrice: 100,
	}

	mockRepo.On("CreateOrder", mock.AnythingOfType("*domain.Order")).Return(nil)
	mockPublisher.On("Publish", mock.AnythingOfType("*domain.Order")).Return(nil)

	uc := usecase.NewOrderUseCase(mockRepo, mockPublisher, mockCache)

	err := uc.CreateOrder(order)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestCreateOrder_ValidationError(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockPublisher := new(mockOrderPublisher)
	mockCache := new(mockOrderCache)

	order := &domain.Order{
		UserID: "", // ошибка
		Items:  []domain.OrderItem{},
	}

	uc := usecase.NewOrderUseCase(mockRepo, mockPublisher, mockCache)

	err := uc.CreateOrder(order)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user id cannot be empty")
}

func TestGetOrderByID_CacheHit(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockPublisher := new(mockOrderPublisher)
	mockCache := new(mockOrderCache)

	order := &domain.Order{ID: "order123", UserID: "user123"}

	mockCache.On("GetOrder", mock.Anything, "order123").Return(order, nil)

	uc := usecase.NewOrderUseCase(mockRepo, mockPublisher, mockCache)

	res, err := uc.GetOrderByID("order123")

	assert.NoError(t, err)
	assert.Equal(t, "order123", res.ID)
	mockCache.AssertExpectations(t)
}

func TestGetOrderByID_DBHit(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockPublisher := new(mockOrderPublisher)
	mockCache := new(mockOrderCache)

	order := &domain.Order{ID: "order123", UserID: "user123"}

	mockCache.On("GetOrder", mock.Anything, "order123").Return(nil, errors.New("not found"))
	mockRepo.On("GetOrderByID", "order123").Return(order, nil)
	mockCache.On("SetOrder", mock.Anything, order).Return(nil)

	uc := usecase.NewOrderUseCase(mockRepo, mockPublisher, mockCache)

	res, err := uc.GetOrderByID("order123")

	assert.NoError(t, err)
	assert.Equal(t, "order123", res.ID)
}

func TestDeleteOrder_Success(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockPublisher := new(mockOrderPublisher)
	mockCache := new(mockOrderCache)

	mockRepo.On("DeleteOrder", "order123").Return(nil)
	mockCache.On("DeleteOrder", mock.Anything, "order123").Return(nil)

	uc := usecase.NewOrderUseCase(mockRepo, mockPublisher, mockCache)

	err := uc.DeleteOrder("order123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
