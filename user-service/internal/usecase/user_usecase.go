package usecase

import (
	"fmt"
	"log"
	"time"

	"user-service/internal/domain"
	"user-service/internal/repository"

	"github.com/google/uuid"
)

type UserUseCase interface {
	CreateUser(user *domain.User) (*domain.User, error)
	GetUserByID(id string) (*domain.User, error)
	ListUsers() ([]domain.User, error)
	UpdateUser(id string, user *domain.User) (*domain.User, error)
	DeleteUser(id string) error
}

type userUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
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
	log.Printf("Getting user by ID: %s", id)

	user, err := uc.repo.GetUserByID(id)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return nil, err
	}
	if user == nil {
		log.Printf("User with ID '%s' not found", id)
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (uc *userUseCase) ListUsers() ([]domain.User, error) {
	log.Println("Listing users...")
	users, err := uc.repo.ListUsers()
	if err != nil {
		log.Printf("Error listing users: %v", err)
		return nil, err
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
