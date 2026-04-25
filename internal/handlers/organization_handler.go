package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ancy-shibu/multi-tenant-saas/internal/middleware"
	"github.com/ancy-shibu/multi-tenant-saas/internal/services"
	"github.com/ancy-shibu/multi-tenant-saas/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// OrganizationHandler handles HTTP requests related to organizations.
type OrganizationHandler struct {
	// service is the organization service instance.
	service *services.OrganizationService
}

// NewOrganizationHandler creates a new OrganizationHandler instance.
func NewOrganizationHandler(service *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: service}
}

// CreateOrganization handles POST /organizations - creates a new organization.
func (h *OrganizationHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	err = h.service.CreateOrganization(r.Context(), req.Name, req.Description)
	if err != nil {
		log.Printf("CreateOrganization error: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create organization")
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, map[string]string{
		"message": "organization created",
	})
}

// GetOrganizations handles GET /organizations - retrieves all organizations.
// @Summary Get User Organizations
// @Description Retrieves a list of organizations the authenticated user belongs to
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Organization
// @Failure 401 {object} map[string]string
// @Router /organizations [get]
func (h *OrganizationHandler) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orgs, err := h.service.GetOrganizations(r.Context(), userID)
	if err != nil {
		log.Printf("GetOrganizations error: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch organizations")
		return
	}

	utils.WriteSuccess(w, http.StatusOK, orgs)
}

// GetOrganizationByID handles GET /organizations/{id} - retrieves an organization by ID.
func (h *OrganizationHandler) GetOrganizationByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "Organization ID is required")
		return
	}

	org, err := h.service.GetOrganizationByID(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Organization not found")
		return
	}

	utils.WriteSuccess(w, http.StatusOK, org)
}

// UpdateOrganization handles PATCH /organizations/{id} - updates an organization.
func (h *OrganizationHandler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	err = h.service.UpdateOrganization(r.Context(), id, req.Name)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update organization")
		return
	}

	utils.WriteSuccess(w, http.StatusOK, map[string]string{
		"message": "organization updated",
	})
}
