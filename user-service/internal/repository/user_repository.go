package repository

import (
	"user-service/internal/domain"
)

type UserRepository interface {
	CreateUser(user *domain.User) (*domain.User, error)
	GetUserByID(id string) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	ListUsers() ([]domain.User, error)
	UpdateUser(user *domain.User) (*domain.User, error)
	DeleteUser(id string) error
}
