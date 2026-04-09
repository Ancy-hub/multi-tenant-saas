package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ancy-shibu/multi-tenant-saas/internal/services"
	"github.com/ancy-shibu/multi-tenant-saas/internal/utils"
	"github.com/go-chi/chi/v5"
)

// CreateUserRequest represents the request payload for creating a user.
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	service *services.UserService
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser handles POST /users - creates a new user.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	if !utils.IsValidEmail(req.Email) {
		utils.WriteError(w, http.StatusBadRequest, "Invalid email")
		return
	}
	if !utils.IsStrongPassword(req.Password) {
		utils.WriteError(w, http.StatusBadRequest,
			"Password must contain uppercase, lowercase, number, special char")
		return
	}
	err = h.service.CreateUser(r.Context(), req.Name, req.Password, req.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteSuccess(w, http.StatusCreated, map[string]string{
		"message": "user created",
	})
}

// GetUserByID handles GET /users/{id} - retrieves a user by ID.
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	utils.WriteSuccess(w, http.StatusOK, user)
}

// GetAllUsers handles GET /users - retrieves all users.
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Failed to fetch users")
	}

	utils.WriteSuccess(w, http.StatusOK, users)
}

// Login handles POST /login - authenticates a user and returns a JWT token.
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	utils.WriteSuccess(w, http.StatusOK, map[string]string{
		"token": token,
	})
}
