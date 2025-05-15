package delivery

import (
	"encoding/json"
	"log"
	"net/http"

	"user-service/internal/domain"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	useCase domain.UserUseCase
}

func NewUserHandler(useCase domain.UserUseCase) *UserHandler {
	return &UserHandler{useCase: useCase}
}

func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	log.Printf("Error: %s, StatusCode: %d", message, statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
		"code":  statusCode,
	})
}

// POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create user")

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeErrorResponse(w, "Invalid input", http.StatusBadRequest)
		return
	}

	createdUser, err := h.useCase.CreateUser(&user)
	if err != nil {
		writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("User created successfully: %+v", createdUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// GET /users/{id}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	log.Printf("Received request to get user with ID: %s", id)

	user, err := h.useCase.GetUserByID(id)
	if err != nil {
		writeErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("User found: %+v", user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GET /users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to list all users")

	users, err := h.useCase.ListUsers()
	if err != nil {
		writeErrorResponse(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	log.Printf("Users fetched successfully: %+v", users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// PATCH /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	log.Printf("Received request to update user with ID: %s", id)

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeErrorResponse(w, "Invalid input", http.StatusBadRequest)
		return
	}

	updatedUser, err := h.useCase.UpdateUser(id, &user)
	if err != nil {
		writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("User updated successfully: %+v", updatedUser)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// DELETE /users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	log.Printf("Received request to delete user with ID: %s", id)

	if err := h.useCase.DeleteUser(id); err != nil {
		writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("User deleted successfully with ID: %s", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
