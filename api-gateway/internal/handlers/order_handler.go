package handlers

import (
	"context"
	"net/http"

	"api-gateway/internal/clients"
	"api-gateway/internal/proto/orders/proto" // Путь к сгенерированным протобуфам

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	client *clients.OrderClient
}

func NewOrderHandler(client *clients.OrderClient) *OrderHandler {
	return &OrderHandler{client: client}
}

type OrderItemInput struct {
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type CreateOrderRequest struct {
	UserID     string           `json:"user_id"`
	Items      []OrderItemInput `json:"items"`
	TotalPrice float64          `json:"total_price"`
	Status     string           `json:"status"`
}

// CreateOrder обрабатывает создание нового заказа
// func (h *OrderHandler) CreateOrder(c *gin.Context) {
// 	var orderRequest proto.Order
// 	if err := c.ShouldBindJSON(&orderRequest); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	// Взаимодействие с OrderClient для создания заказа
// 	orderResponse, err := h.client.CreateOrder(context.Background(), &orderRequest)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order", "details": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, orderResponse)
// }

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Преобразование в proto.Order
	order := &proto.Order{
		UserId:     req.UserID,
		TotalPrice: req.TotalPrice,
		Items:      []*proto.OrderItem{},
	}

	for _, item := range req.Items {
		order.Items = append(order.Items, &proto.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	// gRPC вызов
	orderResponse, err := h.client.CreateOrder(context.Background(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderResponse)
}

// GetOrderByID обрабатывает получение заказа по ID
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")

	// Взаимодействие с OrderClient для получения заказа по ID
	order, err := h.client.GetOrderByID(context.Background(), &proto.OrderId{Id: orderID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ListOrders обрабатывает получение списка всех заказов
func (h *OrderHandler) ListOrders(c *gin.Context) {
	// Взаимодействие с OrderClient для получения списка всех заказов
	orderList, err := h.client.ListOrders(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order list", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderList)
}

// UpdateOrder обрабатывает обновление данных заказа
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	var orderRequest proto.Order
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Взаимодействие с OrderClient для обновления заказа
	orderResponse, err := h.client.UpdateOrder(context.Background(), &orderRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderResponse)
}

// DeleteOrder обрабатывает удаление заказа по ID
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	orderID := c.Param("id")

	// Взаимодействие с OrderClient для удаления заказа
	err := h.client.DeleteOrder(context.Background(), &proto.OrderId{Id: orderID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
