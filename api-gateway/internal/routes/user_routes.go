package routes

import (
	"api-gateway/internal/clients"
	"api-gateway/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, userClient *clients.UserClient) {
	userHandler := handlers.NewUserHandler(userClient)

	userGroup := router.Group("/api/users")
	{
		userGroup.POST("/", userHandler.CreateUser)
		userGroup.GET("/:id", userHandler.GetUserByID)
		userGroup.GET("/", userHandler.ListUsers)
		userGroup.PUT("/:id", userHandler.UpdateUser)
		userGroup.DELETE("/:id", userHandler.DeleteUser)
	}
}
