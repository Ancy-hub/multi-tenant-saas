package handlers

import (
	"encoding/json"
	"log"
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
	user, token, refresh, err := h.service.CreateUser(r.Context(), req.Name, req.Password, req.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteSuccess(w, http.StatusCreated, map[string]interface{}{
		"message":       "user created",
		"access_token":  token,
		"refresh_token": refresh,
		"user": map[string]interface{}{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
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
// LoginRequest represents the login credentials.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login authenticates a user and returns a JWT token.
// @Summary User Login
// @Description Authenticates a user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("Login request received from %s", r.RemoteAddr)
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		log.Printf("Login validation failed: %v", err)
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	log.Printf("Login attempt for email: %s", req.Email)
	token, refresh, user, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("Login failed for %s: %v", req.Email, err)
		utils.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	log.Printf("Login successful for %s", req.Email)
	utils.WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"access_token":  token,
		"refresh_token": refresh,
		"user": map[string]interface{}{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	})
}

func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.RefreshToken == "" {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	userID, err := h.service.ValidateRefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}
	newAccess, err := utils.GenerateAccessToken(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	utils.WriteSuccess(w, http.StatusOK, map[string]string{
		"access_token": newAccess,
	})
}
