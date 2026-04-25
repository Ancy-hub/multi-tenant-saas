package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ancy-shibu/multi-tenant-saas/internal/middleware"
	"github.com/ancy-shibu/multi-tenant-saas/internal/services"
	"github.com/ancy-shibu/multi-tenant-saas/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	service *services.ProjectService
}

func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}
// CREATE PROJECT
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	orgIDParam := chi.URLParam(r, "id")

	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid org_id")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = h.service.CreateProject(r.Context(), req.Name, req.Description, orgID, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, "project created")
}
// GET PROJECTS
func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	orgIDParam := chi.URLParam(r, "id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	// Assuming limit/offset might be passed, defaulting for now
	limit := 50
	offset := 0

	projects, err := h.service.GetProjects(r.Context(), orgID, limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, projects)
}

// GetProjectByID handles GET /projects/{project_id} - retrieves a single project.
func (h *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	projectIDParam := chi.URLParam(r, "project_id")
	projectID, err := uuid.Parse(projectIDParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid project id")
		return
	}

	project, err := h.service.GetProjectByID(r.Context(), projectID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "project not found")
		return
	}

	utils.WriteSuccess(w, http.StatusOK, project)
}
// DELETE PROJECT
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectIDParam := chi.URLParam(r, "project_id")

	projectID, err := uuid.Parse(projectIDParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid project_id")
		return
	}

	err = h.service.DeleteProject(r.Context(), projectID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "project deleted")
}