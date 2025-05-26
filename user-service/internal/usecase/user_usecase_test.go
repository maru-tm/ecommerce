package usecase_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
	"user-service/internal/domain"
	"user-service/internal/usecase"

	"github.com/google/uuid"
)

// --- Мок UserRepository ---
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
	list := make([]domain.User, 0, len(m.users))
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

// --- Мок UserCache ---
type mockUserCache struct {
	mu    sync.RWMutex
	users map[string]*domain.User // ключ — ID пользователя
}

func newMockUserCache() *mockUserCache {
	return &mockUserCache{
		users: make(map[string]*domain.User),
	}
}

func (c *mockUserCache) SetUser(ctx context.Context, user *domain.User) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.users[user.ID] = user
	return nil
}

func (c *mockUserCache) GetUser(ctx context.Context, id string) (*domain.User, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if user, ok := c.users[id]; ok {
		return user, nil
	}
	return nil, nil
}

func (c *mockUserCache) DeleteUser(ctx context.Context, id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.users, id)
	return nil
}

func (m *mockUserCache) SetUsers(ctx context.Context, key string, users []domain.User) error {
	// Упрощено — не реализуем
	return nil
}

func (m *mockUserCache) GetUsers(ctx context.Context, key string) ([]domain.User, error) {
	// Упрощено — не реализуем
	return nil, nil
}

// --- Вспомогательная функция для времени ---
func timePtr(t time.Time) *time.Time {
	return &t
}

// --- Тесты ---

func TestCreateUser_Success(t *testing.T) {
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
		t.Error("expected ID to be set")
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

func TestCreateUser_ExistingUsername(t *testing.T) {
	repo := newMockUserRepo()
	cache := newMockUserCache()
	uc := usecase.NewUserUseCase(repo, cache)

	existingUser := &domain.User{
		Username:     "existuser",
		Email:        "exist@example.com",
		PasswordHash: "hash",
		FullName:     "Exist User",
		ID:           uuid.New().String(),
	}
	repo.users[existingUser.Username] = existingUser

	user := &domain.User{
		Username:     "existuser",
		Email:        "new@example.com",
		PasswordHash: "hash",
		FullName:     "New User",
	}

	_, err := uc.CreateUser(user)
	if err == nil {
		t.Fatal("expected error for existing username")
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

func TestGetUserByID_DBHit(t *testing.T) {
	repo := newMockUserRepo()
	cache := newMockUserCache()
	uc := usecase.NewUserUseCase(repo, cache)

	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     "dbuser",
		Email:        "db@example.com",
		PasswordHash: "hash",
		FullName:     "DB User",
		Status:       domain.StatusActive,
		CreatedAt:    timePtr(time.Now()),
		UpdatedAt:    timePtr(time.Now()),
	}

	repo.users[user.Username] = user

	got, err := uc.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != user.ID {
		t.Fatalf("expected user from DB, got nil or wrong user")
	}

	// Проверим, что пользователь кешируется
	cachedUser, _ := cache.GetUser(context.Background(), user.ID)
	if cachedUser == nil || cachedUser.ID != user.ID {
		t.Errorf("expected user to be cached")
	}
}

func TestUpdateUser_Success(t *testing.T) {
	repo := newMockUserRepo()
	cache := newMockUserCache()
	uc := usecase.NewUserUseCase(repo, cache)

	// Создадим и добавим пользователя в репозиторий
	existingUser := &domain.User{
		ID:           uuid.New().String(),
		Username:     "updateuser",
		Email:        "old@example.com",
		PasswordHash: "oldhash",
		FullName:     "Old Name",
		Status:       domain.StatusActive,
		CreatedAt:    timePtr(time.Now().Add(-time.Hour)),
		UpdatedAt:    timePtr(time.Now().Add(-time.Hour)),
	}
	repo.users[existingUser.Username] = existingUser

	// Обновим пользователя
	updatedData := &domain.User{
		Username:     "updateuser",
		Email:        "new@example.com",
		PasswordHash: "newhash",
		FullName:     "New Name",
	}

	updatedUser, err := uc.UpdateUser(existingUser.ID, updatedData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updatedUser.Email != "new@example.com" {
		t.Errorf("expected email to be updated, got %s", updatedUser.Email)
	}
	if updatedUser.FullName != "New Name" {
		t.Errorf("expected full name to be updated, got %s", updatedUser.FullName)
	}
	if updatedUser.ID != existingUser.ID {
		t.Errorf("expected ID to remain unchanged")
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	repo := newMockUserRepo()
	cache := newMockUserCache()
	uc := usecase.NewUserUseCase(repo, cache)

	nonExistID := uuid.New().String()

	user := &domain.User{
		Username:     "nonexist",
		Email:        "nonexist@example.com",
		PasswordHash: "hash",
		FullName:     "Non Exist",
	}

	_, err := uc.UpdateUser(nonExistID, user)
	if err == nil {
		t.Fatal("expected error for non-existing user")
	}
}

func TestDeleteUser_Success(t *testing.T) {
	repo := newMockUserRepo()
	cache := newMockUserCache()
	uc := usecase.NewUserUseCase(repo, cache)

	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     "deluser",
		Email:        "del@example.com",
		PasswordHash: "hash",
		FullName:     "Delete User",
		Status:       domain.StatusActive,
		CreatedAt:    timePtr(time.Now()),
		UpdatedAt:    timePtr(time.Now()),
	}

	repo.users[user.Username] = user
	cache.SetUser(context.Background(), user)

	err := uc.DeleteUser(user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Проверяем, что пользователь удалён из репозитория
	if _, exists := repo.users[user.Username]; exists {
		t.Errorf("user should be deleted from repo")
	}

	// Проверяем, что пользователь удалён из кеша
	if cached, _ := cache.GetUser(context.Background(), user.ID); cached != nil {
		t.Errorf("user should be deleted from cache")
	}
}

func TestListUsers(t *testing.T) {
	repo := newMockUserRepo()
	cache := newMockUserCache()
	uc := usecase.NewUserUseCase(repo, cache)

	user1 := &domain.User{
		ID:           uuid.New().String(),
		Username:     "user1",
		Email:        "u1@example.com",
		PasswordHash: "hash",
		FullName:     "User One",
		Status:       domain.StatusActive,
	}
	user2 := &domain.User{
		ID:           uuid.New().String(),
		Username:     "user2",
		Email:        "u2@example.com",
		PasswordHash: "hash",
		FullName:     "User Two",
		Status:       domain.StatusActive,
	}

	repo.users[user1.Username] = user1
	repo.users[user2.Username] = user2

	users, err := uc.ListUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}
