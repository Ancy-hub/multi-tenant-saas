package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ancy-shibu/multi-tenant-saas/internal/services"
	"github.com/ancy-shibu/multi-tenant-saas/internal/utils"
	"github.com/go-chi/chi/v5"
)

type OrganizationHandler struct {
	service *services.OrganizationService
}

func NewOrganizationHandler(service *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: service}
}

// POST /organizations
func (h *OrganizationHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	err = h.service.CreateOrganization(r.Context(), req.Name)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create organization")
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, map[string]string{
		"message": "organization created",
	})
}

// GET /organizations
func (h *OrganizationHandler) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	orgs, err := h.service.GetOrganizations(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch organizations")
		return
	}

	utils.WriteSuccess(w, http.StatusOK, orgs)
}

// GET /organizations/{id}
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

// PUT /organizations/{id}
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
