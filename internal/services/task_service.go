package services

import (
	"context"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/ancy-shibu/multi-tenant-saas/internal/repository"
	"github.com/google/uuid"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, title, desc string, projectID, userID uuid.UUID) error {
	t := models.Task{
		ID:        uuid.New(),
		ProjectID: projectID,
		Title:     title,
		Description: desc,
		Status:    "todo",
		CreatedBy: userID,
	}

	return s.repo.Create(ctx, t)
}

func (s *TaskService) GetTasks(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]models.Task, error) {
	return s.repo.GetByProject(ctx, projectID, limit, offset)
}

func (s *TaskService) GetTasksByOrganization(ctx context.Context, orgID uuid.UUID) ([]models.Task, error) {
	return s.repo.GetByOrganization(ctx, orgID)
}

func (s *TaskService) UpdateTask(ctx context.Context, t models.Task) error {
	return s.repo.Update(ctx, t)
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	return s.repo.Delete(ctx, taskID)
}

