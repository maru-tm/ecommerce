package routes

import (
	"api-gateway/internal/clients"
	"api-gateway/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(router *gin.Engine, orderClient *clients.OrderClient) {
	orderHandler := handlers.NewOrderHandler(orderClient)

	orderGroup := router.Group("/api/orders")
	{
		orderGroup.POST("/", orderHandler.CreateOrder)
		orderGroup.GET("/", orderHandler.ListOrders)
		orderGroup.GET("/:id", orderHandler.GetOrderByID)
		orderGroup.PUT("/:id", orderHandler.UpdateOrder)
		orderGroup.DELETE("/:id", orderHandler.DeleteOrder)
	}
}
