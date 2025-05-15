package dto

import "time"

type CreateUserInput struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password_hash"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
