package services

import (
	"context"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/ancy-shibu/multi-tenant-saas/internal/repository"
	"github.com/google/uuid"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

// CREATE PROJECT
func (s *ProjectService) CreateProject(ctx context.Context, name, description string, orgID, userID uuid.UUID) error {
	p := models.Project{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		OrgID:       orgID,
		CreatedBy:   userID,
	}

	return s.repo.Create(ctx, p)
}
// GetProjects retrieves projects for an organization.
func (s *ProjectService) GetProjects(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]models.Project, error) {
	return s.repo.GetByOrg(ctx, orgID, limit, offset)
}

// GetProjectByID retrieves a single project by its ID.
func (s *ProjectService) GetProjectByID(ctx context.Context, projectID uuid.UUID) (*models.Project, error) {
	return s.repo.GetByID(ctx, projectID)
}

// DELETE PROJECT
func (s *ProjectService) DeleteProject(ctx context.Context, projectID uuid.UUID) error {
	return s.repo.Delete(ctx, projectID)
}