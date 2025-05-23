package routes

import (
	"api-gateway/internal/email"
	"api-gateway/internal/handlers"
	"api-gateway/internal/proto/users/proto"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, userClient proto.UserServiceClient, mailer *email.Mailer) {
	userHandler := handlers.NewUserHandler(userClient, mailer)

	userGroup := router.Group("/api/users")
	{
		userGroup.POST("/", userHandler.CreateUser)
		userGroup.GET("/:id", userHandler.GetUserByID)
		userGroup.GET("/", userHandler.ListUsers)
		userGroup.PUT("/:id", userHandler.UpdateUser)
		userGroup.DELETE("/:id", userHandler.DeleteUser)
	}
}
