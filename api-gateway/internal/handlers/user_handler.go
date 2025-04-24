package handlers

import (
	"context"
	"net/http"

	"api-gateway/internal/clients"
	"api-gateway/internal/proto/users/proto" // Updated path for the generated proto files

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	client *clients.UserClient
}

func NewUserHandler(client *clients.UserClient) *UserHandler {
	return &UserHandler{client: client}
}

// CreateUser handles the creation of a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var userRequest proto.User
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Call the UserClient to create the user
	userResponse, err := h.client.CreateUser(context.Background(), &userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

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
	userList, err := h.client.ListUsers(context.Background())
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
	err := h.client.DeleteUser(context.Background(), &proto.UserId{Id: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
