package handlers

import (
	"context"
	"net/http"

	"api-gateway/internal/clients"
	"api-gateway/internal/proto/products/proto" // Путь к сгенерированным протобуфам

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	client *clients.ProductClient
}

func NewProductHandler(client *clients.ProductClient) *ProductHandler {
	return &ProductHandler{client: client}
}

// CreateProduct обрабатывает создание нового продукта
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var productRequest proto.Product
	if err := c.ShouldBindJSON(&productRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Взаимодействие с ProductClient для создания продукта
	productResponse, err := h.client.CreateProduct(context.Background(), &productRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productResponse)
}

// GetProductByID обрабатывает получение продукта по ID
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productID := c.Param("id")

	// Взаимодействие с ProductClient для получения продукта по ID
	product, err := h.client.GetProductByID(context.Background(), &proto.ProductId{Id: productID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// ListProducts обрабатывает получение списка всех продуктов
func (h *ProductHandler) ListProducts(c *gin.Context) {
	// Взаимодействие с ProductClient для получения списка всех продуктов
	productList, err := h.client.ListProducts(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product list", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productList)
}

// UpdateProduct обрабатывает обновление данных продукта
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var productRequest proto.Product
	if err := c.ShouldBindJSON(&productRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Взаимодействие с ProductClient для обновления продукта
	productResponse, err := h.client.UpdateProduct(context.Background(), &productRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productResponse)
}

// DeleteProduct обрабатывает удаление продукта по ID
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	// Взаимодействие с ProductClient для удаления продукта
	err := h.client.DeleteProduct(context.Background(), &proto.ProductId{Id: productID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
