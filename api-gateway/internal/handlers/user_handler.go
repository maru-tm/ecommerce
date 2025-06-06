package handlers

import (
	"context"
	"net/http"

	"api-gateway/internal/email"
	"api-gateway/internal/handlers/dto"
	"api-gateway/internal/proto/users/proto" // Updated path for the generated proto files

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserHandler struct {
	client proto.UserServiceClient
	mailer *email.Mailer
}

func NewUserHandler(client proto.UserServiceClient, mailer *email.Mailer) *UserHandler {
	return &UserHandler{
		client: client,
		mailer: mailer,
	}
}

// CreateUser handles the creation of a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var input dto.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "details": err.Error()})
		return
	}

	// // Преобразуем статус в enum
	// statusVal, ok := proto.UserStatus_value[strings.ToUpper(input.Status)]
	// if !ok {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user status"})
	// 	return
	// }

	// Преобразуем в proto.User
	userProto := &proto.User{
		Id:           input.ID,
		Username:     input.Username,
		PasswordHash: input.Password,
		Email:        input.Email,
		FullName:     input.FullName,
		CreatedAt:    timestamppb.New(input.CreatedAt),
		UpdatedAt:    timestamppb.New(input.UpdatedAt),
	}

	// Отправляем в gRPC
	userResponse, err := h.client.CreateUser(context.Background(), userProto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	go func(email, name string) {
		if err := h.mailer.SendWelcomeEmail(email, name); err != nil {
			// Логируем ошибку, но не мешаем API
			// log.Printf("Failed to send welcome email to %s: %v", email, err)
		}
	}(input.Email, input.FullName)

	c.JSON(http.StatusOK, userResponse)
}

// GetUserByID handles the retrieval of a user by their ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	// Call the UserClient to get the user by ID
	userProfile, err := h.client.GetUserByID(context.Background(), &proto.UserId{Id: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userProfile)
}

// ListUsers handles the retrieval of all users
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Call the UserClient to list all users
	userList, err := h.client.ListUsers(context.Background(), &proto.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user list", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userList)
}

// UpdateUser handles updating the user details
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var userRequest proto.User
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Call the UserClient to update the user
	userResponse, err := h.client.UpdateUser(context.Background(), &userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userResponse)
}

// DeleteUser handles deleting a user by their ID
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	// Call the UserClient to delete the user
	_, err := h.client.DeleteUser(context.Background(), &proto.UserId{Id: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
