package clients

import (
	"context"
	"fmt"

	proto "api-gateway/internal/proto/users/proto" // Путь к сгенерированным протобуфам

	"google.golang.org/grpc"
)

type UserClient struct {
	client proto.UserServiceClient
}

// NewUserClient создает новый gRPC-клиент для сервиса пользователей
func NewUserClient(address string) (*UserClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}
	client := proto.NewUserServiceClient(conn)
	return &UserClient{client: client}, nil
}

// CreateUser создает нового пользователя
func (uc *UserClient) CreateUser(ctx context.Context, user *proto.User) (*proto.User, error) {
	return uc.client.CreateUser(ctx, user)
}

// GetUserByID возвращает пользователя по ID
func (uc *UserClient) GetUserByID(ctx context.Context, userId *proto.UserId) (*proto.User, error) {
	return uc.client.GetUserByID(ctx, userId)
}

// ListUsers возвращает список всех пользователей
func (uc *UserClient) ListUsers(ctx context.Context) (*proto.UserList, error) {
	return uc.client.ListUsers(ctx, &proto.Empty{})
}

// UpdateUser обновляет данные пользователя
func (uc *UserClient) UpdateUser(ctx context.Context, user *proto.User) (*proto.User, error) {
	return uc.client.UpdateUser(ctx, user)
}

// DeleteUser удаляет пользователя
func (uc *UserClient) DeleteUser(ctx context.Context, userId *proto.UserId) error {
	_, err := uc.client.DeleteUser(ctx, userId)
	return err
}
