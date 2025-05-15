package delivery

import (
	"encoding/json"
	"inventory-service/internal/domain"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	useCase domain.ProductUseCase
}

func NewProductHandler(useCase domain.ProductUseCase) *ProductHandler {
	return &ProductHandler{useCase: useCase}
}

// Функция для обработки ошибок
func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	log.Printf("Error: %s, StatusCode: %d", message, statusCode) // Логирование ошибки
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
		"code":  statusCode,
	})
}

// CREATE /products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create product") // Логирование начала запроса

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		writeErrorResponse(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.useCase.CreateProduct(&product); err != nil {
		writeErrorResponse(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	log.Printf("Product created successfully: %+v", product) // Логирование успешного создания продукта
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// GET /products/{id}
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	log.Printf("Received request to get product with ID: %s", id) // Логирование запроса

	product, err := h.useCase.GetProductByID(id)
	if err != nil {
		writeErrorResponse(w, "Product not found", http.StatusNotFound)
		return
	}

	log.Printf("Product found: %+v", product) // Логирование успешного поиска продукта
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// GET /products
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to list all products") // Логирование запроса

	products, err := h.useCase.ListProducts()
	if err != nil {
		writeErrorResponse(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	log.Printf("Products fetched successfully: %+v", products) // Логирование успешного получения списка продуктов
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// PATCH /products/{id}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	log.Printf("Received request to update product with ID: %s", id) // Логирование запроса

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		writeErrorResponse(w, "Invalid input", http.StatusBadRequest)
		return
	}

	product.ID = id

	if err := h.useCase.UpdateProduct(id, &product); err != nil {
		writeErrorResponse(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	log.Printf("Product updated successfully: %+v", product) // Логирование успешного обновления продукта
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// DELETE /products/{id}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	log.Printf("Received request to delete product with ID: %s", id) // Логирование запроса

	if err := h.useCase.DeleteProduct(id); err != nil {
		writeErrorResponse(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	log.Printf("Product deleted successfully with ID: %s", id) // Логирование успешного удаления продукта
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) // 204 No Content for successful deletion
}
