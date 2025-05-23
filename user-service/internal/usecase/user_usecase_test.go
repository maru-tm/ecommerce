package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"user-service/internal/domain"
	"user-service/internal/usecase"

	"github.com/google/uuid"
)

// Мок UserRepository
type mockUserRepo struct {
	users map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*domain.User)}
}

func (m *mockUserRepo) CreateUser(user *domain.User) (*domain.User, error) {
	if _, exists := m.users[user.Username]; exists {
		return nil, errors.New("user exists")
	}
	m.users[user.Username] = user
	return user, nil
}

func (m *mockUserRepo) GetUserByID(id string) (*domain.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) GetUserByUsername(username string) (*domain.User, error) {
	if u, ok := m.users[username]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *mockUserRepo) ListUsers() ([]domain.User, error) {
	var list []domain.User
	for _, u := range m.users {
		list = append(list, *u)
	}
	return list, nil
}

func (m *mockUserRepo) UpdateUser(user *domain.User) (*domain.User, error) {
	if _, exists := m.users[user.Username]; !exists {
		return nil, errors.New("user not found")
	}
	m.users[user.Username] = user
	return user, nil
}

func (m *mockUserRepo) DeleteUser(id string) error {
	for k, u := range m.users {
		if u.ID == id {
			delete(m.users, k)
			return nil
		}
	}
	return errors.New("user not found")
}

// Мок UserCache
type mockUserCache struct {
	cache map[string]*domain.User
}

func newMockUserCache() *mockUserCache {
	return &mockUserCache{cache: make(map[string]*domain.User)}
}

func (m *mockUserCache) SetUser(ctx context.Context, user *domain.User) error {
	m.cache[user.ID] = user
	return nil
}

func (m *mockUserCache) GetUser(ctx context.Context, id string) (*domain.User, error) {
	if u, ok := m.cache[id]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *mockUserCache) DeleteUser(ctx context.Context, id string) error {
	delete(m.cache, id)
	return nil
}

func (m *mockUserCache) SetUsers(ctx context.Context, key string, users []domain.User) error {
	// для упрощения — не реализуем
	return nil
}

func (m *mockUserCache) GetUsers(ctx context.Context, key string) ([]domain.User, error) {
	// для упрощения — не реализуем
	return nil, nil
}

func TestCreateUser(t *testing.T) {
	repo := newMockUserRepo()
	cache := newMockUserCache()
	uc := usecase.NewUserUseCase(repo, cache)

	user := &domain.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash",
		FullName:     "Test User",
	}

	createdUser, err := uc.CreateUser(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if createdUser.ID == "" {
		t.Errorf("expected ID to be set")
	}
	if createdUser.Status != domain.StatusActive {
		t.Errorf("expected status active, got %s", createdUser.Status)
	}
}

func TestCreateUser_EmptyUsername(t *testing.T) {
	repo := newMockUserRepo()
	cache := newMockUserCache()
	uc := usecase.NewUserUseCase(repo, cache)

	user := &domain.User{
		Username:     "",
		Email:        "test@example.com",
		PasswordHash: "hash",
		FullName:     "Test User",
	}

	_, err := uc.CreateUser(user)
	if err == nil {
		t.Fatal("expected error for empty username")
	}
}

func TestGetUserByID_CacheHit(t *testing.T) {
	repo := newMockUserRepo()
	cache := newMockUserCache()
	uc := usecase.NewUserUseCase(repo, cache)

	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     "cacheuser",
		Email:        "cache@example.com",
		PasswordHash: "hash",
		FullName:     "Cache User",
		Status:       domain.StatusActive,
		CreatedAt:    timePtr(time.Now()),
		UpdatedAt:    timePtr(time.Now()),
	}

	cache.SetUser(context.Background(), user)

	got, err := uc.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != user.ID {
		t.Fatalf("expected user from cache, got nil or wrong user")
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
