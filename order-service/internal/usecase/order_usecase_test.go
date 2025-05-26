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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}
func (m *mockOrderRepo) ListOrders() ([]domain.Order, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Order), args.Error(1)
}

// ==== TESTS ====

func TestCreateOrder_Success(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockPublisher := new(mockOrderPublisher)
	mockCache := new(mockOrderCache)

	order := &domain.Order{
		ID:     "1",
		UserID: "user123",
		Items: []domain.OrderItem{
			{ProductID: "prod-1", Quantity: 2},
		},
	}

	mockRepo.On("CreateOrder", order).Return(nil)
	mockPublisher.On("Publish", order).Return(nil)

	uc := usecase.NewOrderUseCase(mockRepo, mockPublisher, mockCache)
	err := uc.CreateOrder(order)

	assert.NoError(t, err)
}

func TestCreateOrder_ValidationError(t *testing.T) {
	uc := usecase.NewOrderUseCase(nil, nil, nil)

	err := uc.CreateOrder(&domain.Order{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user id cannot be empty")
}

func TestCreateOrder_PublishError(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockPublisher := new(mockOrderPublisher)
	mockCache := new(mockOrderCache)

	order := &domain.Order{
		UserID:     "user1",
		Items:      []domain.OrderItem{{ProductID: "p1", Quantity: 1}},
		TotalPrice: 10,
	}

	mockRepo.On("CreateOrder", order).Return(nil)
	mockPublisher.On("Publish", order).Return(errors.New("pub error"))

	uc := usecase.NewOrderUseCase(mockRepo, mockPublisher, mockCache)
	err := uc.CreateOrder(order)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pub error")
}

func TestGetOrderByID_CacheHit(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockPublisher := new(mockOrderPublisher)
	mockCache := new(mockOrderCache)

	order := &domain.Order{ID: "123", UserID: "user1"}

	mockCache.On("GetOrder", mock.Anything, "123").Return(order, nil)

	uc := usecase.NewOrderUseCase(mockRepo, mockPublisher, mockCache)
	res, err := uc.GetOrderByID("123")

	assert.NoError(t, err)
	assert.Equal(t, "123", res.ID)
}

func TestGetOrderByID_DBHit(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockPublisher := new(mockOrderPublisher)
	mockCache := new(mockOrderCache)

	order := &domain.Order{ID: "123", UserID: "user1"}

	mockCache.On("GetOrder", mock.Anything, "123").Return(nil, errors.New("miss"))
	mockRepo.On("GetOrderByID", "123").Return(order, nil)
	mockCache.On("SetOrder", mock.Anything, order).Return(nil)

	uc := usecase.NewOrderUseCase(mockRepo, mockPublisher, mockCache)
	res, err := uc.GetOrderByID("123")

	assert.NoError(t, err)
	assert.Equal(t, "123", res.ID)
}

func TestGetOrderByID_Error(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockPublisher := new(mockOrderPublisher)
	mockCache := new(mockOrderCache)

	mockCache.On("GetOrder", mock.Anything, "123").Return(nil, errors.New("miss"))
	mockRepo.On("GetOrderByID", "123").Return(nil, errors.New("not found"))

	uc := usecase.NewOrderUseCase(mockRepo, mockPublisher, mockCache)
	_, err := uc.GetOrderByID("123")

	assert.Error(t, err)
}

func TestListOrders_CacheHit(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockPublisher := new(mockOrderPublisher)
	mockCache := new(mockOrderCache)

	orders := []domain.Order{{ID: "1"}, {ID: "2"}}

	mockCache.On("GetOrders", mock.Anything, "orders:all").Return(orders, nil)

	uc := usecase.NewOrderUseCase(mockRepo, mockPublisher, mockCache)
	res, err := uc.ListOrders()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
}

func TestListOrders_DBHit(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockPublisher := new(mockOrderPublisher)
	mockCache := new(mockOrderCache)

	orders := []domain.Order{{ID: "1"}, {ID: "2"}}

	mockCache.On("GetOrders", mock.Anything, "orders:all").Return(nil, errors.New("miss"))
	mockRepo.On("ListOrders").Return(orders, nil)
	mockCache.On("SetOrders", mock.Anything, "orders:all", orders).Return(nil)

	uc := usecase.NewOrderUseCase(mockRepo, mockPublisher, mockCache)
	res, err := uc.ListOrders()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
}

func TestUpdateOrder_Success(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockPublisher := new(mockOrderPublisher)
	mockCache := new(mockOrderCache)

	order := &domain.Order{
		ID:     "1",
		UserID: "user1",
		Items: []domain.OrderItem{
			{ProductID: "prod-1", Quantity: 1},
		},
	}

	mockRepo.On("UpdateOrder", order).Return(nil)
	mockPublisher.On("Publish", order).Return(nil)
	mockCache.On("DeleteOrder", mock.Anything, order.ID).Return(nil) // вот тут

	uc := usecase.NewOrderUseCase(mockRepo, mockPublisher, mockCache)
	err := uc.UpdateOrder(order)

	assert.NoError(t, err)
}

func TestDeleteOrder_Success(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockPublisher := new(mockOrderPublisher)
	mockCache := new(mockOrderCache)

	mockRepo.On("DeleteOrder", "1").Return(nil)
	mockCache.On("DeleteOrder", mock.Anything, "1").Return(nil)

	uc := usecase.NewOrderUseCase(mockRepo, mockPublisher, mockCache)
	err := uc.DeleteOrder("1")

	assert.NoError(t, err)
}
