package routes

import (
	"api-gateway/internal/clients"
	"api-gateway/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(router *gin.Engine, productClient *clients.ProductClient) {
	productHandler := handlers.NewProductHandler(productClient)

	// Роуты для работы с продуктами
	productGroup := router.Group("/api/products")
	{
		productGroup.POST("/", productHandler.CreateProduct)      // Создание продукта
		productGroup.GET("/", productHandler.ListProducts)        // Список продуктов
		productGroup.GET("/:id", productHandler.GetProductByID)   // Получить продукт по ID
		productGroup.PUT("/:id", productHandler.UpdateProduct)    // Обновить продукт
		productGroup.DELETE("/:id", productHandler.DeleteProduct) // Удалить продукт
	}
}
