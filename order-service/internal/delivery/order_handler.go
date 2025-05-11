package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"order-service/internal/domain"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	useCase domain.OrderUseCase
}

func NewOrderHandler(useCase domain.OrderUseCase) *OrderHandler {
	return &OrderHandler{useCase: useCase}
}

// writeErrorResponse sends a JSON error response with a given status code.
func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	log.Printf("Error: %s, StatusCode: %d\n", message, statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
		"code":  statusCode,
	})
}

// CREATE /orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create order")

	w.Header().Set("Content-Type", "application/json")

	var order domain.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		writeErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.useCase.CreateOrder(&order); err != nil {
		writeErrorResponse(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(order)
	log.Printf("Successfully created order with ID: %s\n", order.ID)
}

// GET /orders/{id}
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get order by ID")

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("Fetching order with ID: %s\n", id)

	order, err := h.useCase.GetOrderByID(id)
	if err != nil {
		writeErrorResponse(w, "Failed to get order", http.StatusInternalServerError)
		return
	}
	if order == nil {
		writeErrorResponse(w, "Order not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(order)
	log.Printf("Successfully fetched order with ID: %s\n", id)
}

// GET /orders
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to list all orders")

	w.Header().Set("Content-Type", "application/json")

	orders, err := h.useCase.ListOrders()
	if err != nil {
		writeErrorResponse(w, "Failed to list orders", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(orders)
	log.Printf("Successfully fetched list of orders\n")
}

// PATCH /orders/{id}
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to update order")

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("Fetching order with ID: %s\n", id)

	var order domain.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		writeErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}
	order.ID = id

	if err := h.useCase.UpdateOrder(&order); err != nil {
		writeErrorResponse(w, "Failed to update order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(order)
	log.Printf("Successfully updated order with ID: %s\n", id)
}

// DELETE /orders/{id}
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to delete order")

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("Deleting order with ID: %s\n", id)

	if err := h.useCase.DeleteOrder(id); err != nil {
		writeErrorResponse(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Successfully deleted order with ID: %s\n", id)
}
