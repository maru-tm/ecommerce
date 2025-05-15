package domain

import (
	"time"

	"user-service/internal/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusBanned   UserStatus = "banned"
)

type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"password"`
	Email        string     `json:"email"`
	FullName     string     `json:"full_name"`
	Status       UserStatus `json:"status"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

type UserProfile struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	FullName  string     `json:"full_name"`
	Status    UserStatus `json:"status"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (u *User) ToProto() *proto.User {
	var createdAt *timestamppb.Timestamp
	if u.CreatedAt != nil {
		createdAt = timestamppb.New(*u.CreatedAt)
	}

	var updatedAt *timestamppb.Timestamp
	if u.UpdatedAt != nil {
		updatedAt = timestamppb.New(*u.UpdatedAt)
	}

	return &proto.User{
		Id:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		Email:        u.Email,
		FullName:     u.FullName,
		Status:       proto.UserStatus(proto.UserStatus_value[string(u.Status)]),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

func UserFromProto(p *proto.User) (*User, error) {
	var createdAt *time.Time
	if p.GetCreatedAt() != nil {
		t := p.GetCreatedAt().AsTime()
		createdAt = &t
	}

	var updatedAt *time.Time
	if p.GetUpdatedAt() != nil {
		t := p.GetUpdatedAt().AsTime()
		updatedAt = &t
	}

	return &User{
		ID:           p.GetId(),
		Username:     p.GetUsername(),
		PasswordHash: p.GetPasswordHash(),
		Email:        p.GetEmail(),
		FullName:     p.GetFullName(),
		Status:       UserStatus(p.GetStatus().String()),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}
