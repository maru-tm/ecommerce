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
	return &userUseCase{repo: repo}
}

func (uc *userUseCase) validateUser(user *domain.User) error {
	log.Printf("Validating user: %+v", user)

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
	log.Println("Creating user...")

	if err := uc.validateUser(user); err != nil {
		log.Printf("Validation failed: %v", err)
		return nil, err
	}

	existingUser, err := uc.repo.GetUserByUsername(user.Username)
	if err != nil {
		log.Printf("Error checking username uniqueness: %v", err)
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with username '%s' already exists", user.Username)
	}

	// Generate ID and timestamps
	user.ID = uuid.New().String()
	now := time.Now()
	user.CreatedAt = &now
	user.UpdatedAt = &now
	user.Status = domain.StatusActive

	createdUser, err := uc.repo.CreateUser(user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil, err
	}
	log.Printf("User created: %+v", createdUser)
	return createdUser, nil
}

func (uc *userUseCase) GetUserByID(id string) (*domain.User, error) {
	ctx := context.Background()

	// Сначала пробуем кэш
	if user, _ := uc.cache.GetUser(ctx, id); user != nil {
		log.Printf("[GET] User found in cache: %s", id)
		return user, nil
	}

	// Иначе — база
	log.Printf("[GET] User not found in cache. Fetching from DB: %s", id)
	user, err := uc.repo.GetUserByID(id)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return nil, err
	}
	if user == nil {
		log.Printf("User with ID '%s' not found", id)
		return nil, fmt.Errorf("user not found")
	}
	if err := uc.cache.SetUser(ctx, user); err != nil {
		log.Printf("[CACHE] Failed to set user in cache: %v", err)
	}
	return user, nil
}

func (uc *userUseCase) ListUsers() ([]domain.User, error) {
	log.Println("Listing users...")

	ctx := context.Background()
	cacheKey := "users:all"

	if users, _ := uc.cache.GetUsers(ctx, cacheKey); users != nil {
		log.Println("[GET] Users found in cache")
		return users, nil
	}

	log.Println("[GET] Users not found in cache. Fetching from DB")
	users, err := uc.repo.ListUsers()
	if err != nil {
		log.Printf("Error listing users: %v", err)
		return nil, err
	}
	if err := uc.cache.SetUsers(ctx, cacheKey, users); err != nil {
		log.Printf("[CACHE] Failed to set users in cache: %v", err)
	}
	return users, nil
}

func (uc *userUseCase) UpdateUser(id string, user *domain.User) (*domain.User, error) {
	log.Printf("Updating user ID: %s", id)

	existingUser, err := uc.repo.GetUserByID(id)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		return nil, err
	}
	if existingUser == nil {
		return nil, fmt.Errorf("user with ID '%s' not found", id)
	}

	if err := uc.validateUser(user); err != nil {
		log.Printf("Validation failed: %v", err)
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
		log.Printf("Failed to update user: %v", err)
		return nil, err
	}
	log.Printf("User updated: %+v", updatedUser)
	return updatedUser, nil
}

func (uc *userUseCase) DeleteUser(id string) error {
	log.Printf("Deleting user ID: %s", id)

	existingUser, err := uc.repo.GetUserByID(id)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		return err
	}
	if existingUser == nil {
		return fmt.Errorf("user with ID '%s' not found", id)
	}

	if err := uc.repo.DeleteUser(id); err != nil {
		log.Printf("Failed to delete user: %v", err)
		return err
	}
	log.Printf("User with ID '%s' deleted", id)
	return nil
}
