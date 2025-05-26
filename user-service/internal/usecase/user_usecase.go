package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"user-service/internal/domain"

	"github.com/google/uuid"
)

type userUseCase struct {
	repo  domain.UserRepository
	cache domain.UserCache
}

func NewUserUseCase(repo domain.UserRepository, cache domain.UserCache) domain.UserUseCase {
	log.Printf("[INFO] NewUserUseCase: initialized")
	return &userUseCase{repo: repo, cache: cache}
}

func (uc *userUseCase) validateUser(user *domain.User) error {
	log.Printf("[DEBUG] validateUser: validating user: %+v", user)

	if user.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if user.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if user.PasswordHash == "" {
		return fmt.Errorf("password hash cannot be empty")
	}
	if user.FullName == "" {
		return fmt.Errorf("full name cannot be empty")
	}
	return nil
}

func (uc *userUseCase) CreateUser(user *domain.User) (*domain.User, error) {
	log.Printf("[INFO] CreateUser: start creating user username=%s", user.Username)

	if err := uc.validateUser(user); err != nil {
		log.Printf("[ERROR] CreateUser: validation failed: %v", err)
		return nil, err
	}

	existingUser, err := uc.repo.GetUserByUsername(user.Username)
	if err != nil {
		log.Printf("[ERROR] CreateUser: error checking username uniqueness: %v", err)
		return nil, err
	}
	if existingUser != nil {
		err := fmt.Errorf("user with username '%s' already exists", user.Username)
		log.Printf("[ERROR] CreateUser: %v", err)
		return nil, err
	}

	user.ID = uuid.New().String()
	now := time.Now()
	user.CreatedAt = &now
	user.UpdatedAt = &now
	user.Status = domain.StatusActive

	createdUser, err := uc.repo.CreateUser(user)
	if err != nil {
		log.Printf("[ERROR] CreateUser: failed to create user: %v", err)
		return nil, err
	}

	log.Printf("[INFO] CreateUser: user created successfully ID=%s username=%s", createdUser.ID, createdUser.Username)
	return createdUser, nil
}

func (uc *userUseCase) GetUserByID(id string) (*domain.User, error) {
	log.Printf("[INFO] GetUserByID: fetching user ID=%s", id)
	ctx := context.Background()

	user, _ := uc.cache.GetUser(ctx, id)
	if user != nil {
		log.Printf("[DEBUG] GetUserByID: user found in cache ID=%s", id)
		return user, nil
	}

	log.Printf("[DEBUG] GetUserByID: user not found in cache, querying DB ID=%s", id)
	user, err := uc.repo.GetUserByID(id)
	if err != nil {
		log.Printf("[ERROR] GetUserByID: error fetching user from DB ID=%s: %v", id, err)
		return nil, err
	}

	if user == nil {
		log.Printf("[INFO] GetUserByID: user not found ID=%s", id)
		return nil, nil
	}

	if err := uc.cache.SetUser(ctx, user); err != nil {
		log.Printf("[ERROR] GetUserByID: failed to cache user ID=%s: %v", id, err)
	}

	return user, nil
}

func (uc *userUseCase) ListUsers() ([]domain.User, error) {
	log.Printf("[INFO] ListUsers: listing all users")
	ctx := context.Background()
	cacheKey := "users:all"

	users, _ := uc.cache.GetUsers(ctx, cacheKey)
	if users != nil {
		log.Printf("[DEBUG] ListUsers: users found in cache count=%d", len(users))
		return users, nil
	}

	log.Printf("[DEBUG] ListUsers: users not found in cache, querying DB")
	users, err := uc.repo.ListUsers()
	if err != nil {
		log.Printf("[ERROR] ListUsers: error listing users: %v", err)
		return nil, err
	}

	if err := uc.cache.SetUsers(ctx, cacheKey, users); err != nil {
		log.Printf("[ERROR] ListUsers: failed to cache users: %v", err)
	}

	log.Printf("[INFO] ListUsers: successfully retrieved users count=%d", len(users))
	return users, nil
}

func (uc *userUseCase) UpdateUser(id string, user *domain.User) (*domain.User, error) {
	log.Printf("[INFO] UpdateUser: updating user ID=%s", id)

	existingUser, err := uc.repo.GetUserByID(id)
	if err != nil {
		log.Printf("[ERROR] UpdateUser: error checking user existence ID=%s: %v", id, err)
		return nil, err
	}
	if existingUser == nil {
		err := fmt.Errorf("user with ID '%s' not found", id)
		log.Printf("[ERROR] UpdateUser: %v", err)
		return nil, err
	}

	if err := uc.validateUser(user); err != nil {
		log.Printf("[ERROR] UpdateUser: validation failed: %v", err)
		return nil, err
	}

	user.ID = id
	now := time.Now()
	user.UpdatedAt = &now
	if user.CreatedAt == nil {
		user.CreatedAt = existingUser.CreatedAt
	}

	updatedUser, err := uc.repo.UpdateUser(user)
	if err != nil {
		log.Printf("[ERROR] UpdateUser: failed to update user ID=%s: %v", id, err)
		return nil, err
	}

	// Обновляем кеш
	if err := uc.cache.SetUser(context.Background(), updatedUser); err != nil {
		log.Printf("[WARN] UpdateUser: failed to update user in cache ID=%s: %v", id, err)
	}

	log.Printf("[INFO] UpdateUser: user updated successfully ID=%s", updatedUser.ID)
	return updatedUser, nil
}

func (uc *userUseCase) DeleteUser(id string) error {
	log.Printf("[INFO] DeleteUser: deleting user ID=%s", id)

	existingUser, err := uc.repo.GetUserByID(id)
	if err != nil {
		log.Printf("[ERROR] DeleteUser: error checking user existence ID=%s: %v", id, err)
		return err
	}
	if existingUser == nil {
		err := fmt.Errorf("user with ID '%s' not found", id)
		log.Printf("[ERROR] DeleteUser: %v", err)
		return err
	}

	if err := uc.repo.DeleteUser(id); err != nil {
		log.Printf("[ERROR] DeleteUser: failed to delete user ID=%s: %v", id, err)
		return err
	}

	// Удаляем пользователя из кеша
	if err := uc.cache.DeleteUser(context.Background(), id); err != nil {
		log.Printf("[WARN] DeleteUser: failed to delete user from cache ID=%s: %v", id, err)
	}

	log.Printf("[INFO] DeleteUser: user deleted successfully ID=%s", id)
	return nil
}
