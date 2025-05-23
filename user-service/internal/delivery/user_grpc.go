package delivery

import (
	"context"
	"log"

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
	log.Printf("[INFO] NewUserServiceServer: server initialized")
	return &UserServer{uc: uc}
}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.User) (*proto.User, error) {
	log.Printf("[INFO] CreateUser: called for username=%s", req.GetUsername())

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
		log.Printf("[ERROR] CreateUser: failed to create user: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid user data: %v", err)
	}

	log.Printf("[INFO] CreateUser: successfully created user ID=%s", createdUser.ID)
	return userToProto(createdUser), nil
}

func (s *UserServer) GetUserByID(ctx context.Context, req *proto.UserId) (*proto.User, error) {
	log.Printf("[INFO] GetUserByID: called for ID=%s", req.GetId())

	user, err := s.uc.GetUserByID(req.GetId())
	if err != nil {
		log.Printf("[ERROR] GetUserByID: failed to get user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	if user == nil {
		log.Printf("[INFO] GetUserByID: user not found ID=%s", req.GetId())
		return nil, status.Error(codes.NotFound, "user not found")
	}

	log.Printf("[DEBUG] GetUserByID: user found ID=%s", user.ID)
	return userToProto(user), nil
}

func (s *UserServer) ListUsers(ctx context.Context, req *proto.Empty) (*proto.UserList, error) {
	log.Printf("[INFO] ListUsers: called")

	users, err := s.uc.ListUsers()
	if err != nil {
		log.Printf("[ERROR] ListUsers: failed to list users: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	log.Printf("[INFO] ListUsers: %d users listed", len(users))

	var protoUsers []*proto.User
	for _, user := range users {
		protoUsers = append(protoUsers, userToProto(&user))
	}

	return &proto.UserList{Users: protoUsers}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *proto.User) (*proto.User, error) {
	log.Printf("[INFO] UpdateUser: called for ID=%s", req.GetId())

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
		log.Printf("[ERROR] UpdateUser: failed to update user: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid user data: %v", err)
	}

	log.Printf("[INFO] UpdateUser: successfully updated user ID=%s", updatedUser.ID)
	return userToProto(updatedUser), nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *proto.UserId) (*proto.Empty, error) {
	log.Printf("[INFO] DeleteUser: called for ID=%s", req.GetId())

	err := s.uc.DeleteUser(req.GetId())
	if err != nil {
		log.Printf("[ERROR] DeleteUser: failed to delete user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	log.Printf("[INFO] DeleteUser: successfully deleted user ID=%s", req.GetId())
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
