package clients

import (
	"context"

	proto "api-gateway/internal/proto/users/proto" // Путь к сгенерированным протобуфам

	"google.golang.org/grpc"
)

type UserClient struct {
	client proto.UserServiceClient
}

// NewUserClient создает новый gRPC-клиент для сервиса пользователей
func NewUserClient(addr string) (proto.UserServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := proto.NewUserServiceClient(conn)
	return client, nil
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
