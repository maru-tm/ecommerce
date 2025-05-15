package domain

import "context"

type UserCache interface {
	SetUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	SetUsers(ctx context.Context, key string, users []User) error
	GetUsers(ctx context.Context, key string) ([]User, error)
}
