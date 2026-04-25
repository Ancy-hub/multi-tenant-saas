package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ancy-shibu/multi-tenant-saas/internal/middleware"
	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/ancy-shibu/multi-tenant-saas/internal/services"
	"github.com/ancy-shibu/multi-tenant-saas/internal/utils"
	"github.com/ancy-shibu/multi-tenant-saas/internal/worker"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service    *services.TaskService
	dispatcher worker.Dispatcher
}

func NewTaskHandler(service *services.TaskService, dispatcher worker.Dispatcher) *TaskHandler {
	return &TaskHandler{service: service, dispatcher: dispatcher}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	projectIDParam := chi.URLParam(r, "project_id")

	projectID, _ := uuid.Parse(projectIDParam)

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	if req.Title == "" {
		utils.WriteError(w, http.StatusBadRequest, "Title required")
		return
	}

	userID, _ := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	err := h.service.CreateTask(r.Context(), req.Title, req.Description, projectID, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Dispatch asynchronous background job
	if h.dispatcher != nil {
		h.dispatcher.Dispatch(&worker.TaskCreatedJob{
			TaskTitle:  req.Title,
			AssignedTo: uuid.Nil, // For now, unassigned initially
		})
	}

	utils.WriteSuccess(w, http.StatusCreated, "task created")
}
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	projectIDParam := chi.URLParam(r, "project_id")
	projectID, _ := uuid.Parse(projectIDParam)

	limit := 10
	offset := 0

	tasks, err := h.service.GetTasks(r.Context(), projectID, limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTasksByOrg(w http.ResponseWriter, r *http.Request) {
	orgIDParam := chi.URLParam(r, "id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	tasks, err := h.service.GetTasksByOrganization(r.Context(), orgID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, tasks)
}
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskIDParam := chi.URLParam(r, "task_id")
	taskID, _ := uuid.Parse(taskIDParam)

	var req models.Task
	json.NewDecoder(r.Body).Decode(&req)

	req.ID = taskID

	err := h.service.UpdateTask(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "task updated")
}

func(h *TaskHandler)DeleteTask(w http.ResponseWriter, r *http.Request){
	taskIDParam:= chi.URLParam(r,"task_id")
	taskID,_:=uuid.Parse(taskIDParam)
	err:=h.service.DeleteTask(r.Context(),taskID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteSuccess(w, http.StatusOK, "task deleted")
}