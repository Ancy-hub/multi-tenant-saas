package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ancy-shibu/multi-tenant-saas/internal/services"
	"github.com/ancy-shibu/multi-tenant-saas/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MembershipHandler struct {
	service *services.MembershipService
}

func NewMembershipHandler(service *services.MembershipService) *MembershipHandler {
	return &MembershipHandler{service: service}
}

func (h *MembershipHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	orgIDParam := chi.URLParam(r, "id")

	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid org_id")
		return
	}

	var req struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.UserID == "" || req.Role == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing fields")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}

	err = h.service.AddUserToOrg(r.Context(), userID, orgID, req.Role)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, map[string]string{
		"message": "user added to org",
	})
}

func (h *MembershipHandler) GetMembersByOrg(w http.ResponseWriter, r *http.Request) {
	orgIDParam := chi.URLParam(r, "id")

	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid org_id")
		return
	}
	//Pagination
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit := 10
	offset := 0

	if limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}
	if offsetStr != "" {
		fmt.Sscanf(offsetStr, "%d", &offset)
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	members, err := h.service.GetMembersByOrg(r.Context(), orgID, limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, members)
}

func (h *MembershipHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	orgIDParam := chi.URLParam(r, "org_id")
	userIDParam := chi.URLParam(r, "user_id")

	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid org_id")
		return
	}

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user_id")
		return
	}

	err = h.service.RemoveMember(r.Context(), userID, orgID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, map[string]string{
		"message": "member removed",
	})
}

func (h *MembershipHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	orgIDParam := chi.URLParam(r, "org_id")
	userIDParam := chi.URLParam(r, "user_id")

	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid org_id")
		return
	}

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user_id")
		return
	}

	var req struct {
		Role string `json:"role"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Role == "" {
		utils.WriteError(w, http.StatusBadRequest, "Invalid role")
		return
	}

	err = h.service.UpdateRole(r.Context(), userID, orgID, req.Role)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, map[string]string{
		"message": "role updated",
	})
}

func (h *MembershipHandler) GetUserOrgs(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "id")

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user_id")
		return
	}

	orgs, err := h.service.GetUserOrgs(r.Context(), userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, orgs)
}
