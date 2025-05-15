package delivery

import (
	"context"

	"user-service/internal/domain"
	"user-service/internal/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServer struct {
	proto.UnimplementedUserServiceServer
	uc domain.UserUseCase
}

func NewUserServiceServer(uc domain.UserUseCase) *UserServer {
	return &UserServer{uc: uc}
}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.User) (*proto.User, error) {
	createdAt := req.GetCreatedAt().AsTime()
	updatedAt := req.GetUpdatedAt().AsTime()

	user := &domain.User{
		Username:     req.GetUsername(),
		PasswordHash: req.GetPasswordHash(),
		Email:        req.GetEmail(),
		FullName:     req.GetFullName(),
		Status:       domain.UserStatus(req.GetStatus().String()),
		CreatedAt:    &createdAt,
		UpdatedAt:    &updatedAt,
	}

	createdUser, err := s.uc.CreateUser(user)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user data: %v", err)
	}

	return userToProto(createdUser), nil
}

func (s *UserServer) GetUserByID(ctx context.Context, req *proto.UserId) (*proto.User, error) {
	user, err := s.uc.GetUserByID(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return userToProto(user), nil
}

func (s *UserServer) ListUsers(ctx context.Context, req *proto.Empty) (*proto.UserList, error) {
	users, err := s.uc.ListUsers()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	var protoUsers []*proto.User
	for _, user := range users {
		protoUsers = append(protoUsers, userToProto(&user))
	}

	return &proto.UserList{Users: protoUsers}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *proto.User) (*proto.User, error) {
	updatedAt := req.GetUpdatedAt().AsTime()

	user := &domain.User{
		ID:           req.GetId(),
		Username:     req.GetUsername(),
		PasswordHash: req.GetPasswordHash(),
		Email:        req.GetEmail(),
		FullName:     req.GetFullName(),
		Status:       domain.UserStatus(req.GetStatus().String()),
		CreatedAt:    nil,
		UpdatedAt:    &updatedAt,
	}

	updatedUser, err := s.uc.UpdateUser(req.GetId(), user)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user data: %v", err)
	}

	return userToProto(updatedUser), nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *proto.UserId) (*proto.Empty, error) {
	err := s.uc.DeleteUser(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &proto.Empty{}, nil
}

func userToProto(user *domain.User) *proto.User {
	return &proto.User{
		Id:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Email:        user.Email,
		FullName:     user.FullName,
		Status:       proto.UserStatus(proto.UserStatus_value[string(user.Status)]),
		CreatedAt:    timestamppb.New(*user.CreatedAt),
		UpdatedAt:    timestamppb.New(*user.UpdatedAt),
	}
}
