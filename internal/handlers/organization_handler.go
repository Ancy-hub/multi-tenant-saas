package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ancy-shibu/multi-tenant-saas/internal/services"
	"github.com/go-chi/chi/v5"
)
type OrganizationHandler struct{
	service *services.OrganizationService
}

func NewOrganizationHandler(service *services.OrganizationService) *OrganizationHandler{
	return &OrganizationHandler{service: service}
}

func (h *OrganizationHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = h.service.CreateOrganization(r.Context(), req.Name)
	if err != nil {
		fmt.Println("ERROR CREATING ORG:", err)
		http.Error(w, "Failed to create organization", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "organization created",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *OrganizationHandler) GetOrganizations(w http.ResponseWriter, r *http.Request){
	orgs, err:= h.service.GetOrganizations(r.Context())
	if err!=nil{
		http.Error(w,"Failed to fetch organizations",http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/json")
	err = json.NewEncoder(w).Encode(orgs)
	if err!=nil{
		http.Error(w,"Failed to encode response",http.StatusInternalServerError)
	}
}

func (h *OrganizationHandler) GetOrganizationByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	org, err := h.service.GetOrganizationByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Organization not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(org)
}

func (h *OrganizationHandler) UpdateOrganization(w http.ResponseWriter, r *http.Request){
	id:= chi.URLParam(r,"id")
	var req struct{
		Name string `json:"name"`
	}

	err:= json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateOrganization(r.Context(),id,req.Name)
	if err != nil {
		http.Error(w, "Failed to update organization", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "organization name edited",
	}
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}