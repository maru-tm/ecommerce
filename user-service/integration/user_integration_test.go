//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"user-service/internal/domain"
	"user-service/internal/usecase"

	"github.com/stretchr/testify/assert"
)

// DummyRepo implements domain.UserRepository for integration test purposes.
type DummyRepo struct {
	users map[string]*domain.User
}

func NewDummyRepo() *DummyRepo {
	return &DummyRepo{users: make(map[string]*domain.User)}
}

func (r *DummyRepo) CreateUser(user *domain.User) (*domain.User, error) {
	r.users[user.ID] = user
	return user, nil
}

func (r *DummyRepo) GetUserByID(id string) (*domain.User, error) {
	if user, ok := r.users[id]; ok {
		return user, nil
	}
	return nil, nil
}

func (r *DummyRepo) GetUserByUsername(username string) (*domain.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, nil
}

func (r *DummyRepo) ListUsers() ([]domain.User, error) {
	var result []domain.User
	for _, u := range r.users {
		result = append(result, *u)
	}
	return result, nil
}

func (r *DummyRepo) UpdateUser(user *domain.User) (*domain.User, error) {
	r.users[user.ID] = user
	return user, nil
}

func (r *DummyRepo) DeleteUser(id string) error {
	delete(r.users, id)
	return nil
}

// DummyCache implements domain.UserCache for testing purposes.
type DummyCache struct {
	cache map[string]*domain.User
}

func NewDummyCache() *DummyCache {
	return &DummyCache{cache: make(map[string]*domain.User)}
}

func (c *DummyCache) GetUser(_ context.Context, id string) (*domain.User, error) {
	user, ok := c.cache[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (c *DummyCache) SetUser(_ context.Context, user *domain.User) error {
	c.cache[user.ID] = user
	return nil
}

func (c *DummyCache) DeleteUser(_ context.Context, id string) error {
	delete(c.cache, id)
	return nil
}

func (c *DummyCache) GetUsers(_ context.Context, _ string) ([]domain.User, error) {
	// Возвращаем пустой срез, т.к. тесты не требуют реализации
	return nil, nil
}

func (c *DummyCache) SetUsers(_ context.Context, _ string, _ []domain.User) error {
	return nil
}

func TestCreateAndGetUserIntegration(t *testing.T) {
	repo := NewDummyRepo()
	cache := NewDummyCache()
	uc := usecase.NewUserUseCase(repo, cache)

	user := &domain.User{
		Username:     "testuser",
		PasswordHash: "hashedpass",
		Email:        "test@example.com",
		FullName:     "Test User",
	}

	createdUser, err := uc.CreateUser(user)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.NotEmpty(t, createdUser.ID)

	// Should be cached now
	cachedUser, err := uc.GetUserByID(createdUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, cachedUser)
	assert.Equal(t, createdUser.ID, cachedUser.ID)
}
