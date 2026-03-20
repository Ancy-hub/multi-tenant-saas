package handlers

import (
	"encoding/json"
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

func (h *MembershipHandler) AddUser(w http.ResponseWriter, r *http.Request){
	orgID:= chi.URLParam(r,"id")
	var req struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}

	err:=json.NewDecoder(r.Body).Decode(&req)
	if err!=nil{
		utils.WriteError(w,http.StatusBadRequest,"Invalid request")
	}

	if req.UserID==""|| req.Role==""{
		utils.WriteError(w,http.StatusBadRequest,"Missing fields")
		return
	}

	_, err = uuid.Parse(req.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}
	
	err=h.service.AddUserToOrg(r.Context(),req.UserID,orgID,req.Role)
	if err!=nil{
		utils.WriteError(w,http.StatusBadRequest,err.Error())
		return
	}

	utils.WriteJSON(w,http.StatusCreated,map[string]string{
		"message":"user added to org",
	})
}