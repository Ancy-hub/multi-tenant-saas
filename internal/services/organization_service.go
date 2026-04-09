package services

import (
	"context"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/ancy-shibu/multi-tenant-saas/internal/repository"
	"github.com/google/uuid"
)

// OrganizationService provides business logic for organization operations.
type OrganizationService struct {
	// repo is the organization repository instance.
	repo *repository.OrganizationRepository
}

// NewOrganizationService creates a new OrganizationService instance.
func NewOrganizationService(repo *repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repo: repo}
}

// CreateOrganization creates a new organization.
func (s *OrganizationService) CreateOrganization(ctx context.Context, name string) error {
	org := models.Organization{
		ID:   uuid.New(),
		Name: name,
	}

	return s.repo.Create(ctx, org)
}

// GetOrganizations retrieves all organizations.
func (s *OrganizationService) GetOrganizations(ctx context.Context) ([]models.Organization, error) {
	return s.repo.GetAll(ctx)
}

// GetOrganizationByID retrieves an organization by its ID.
func (s *OrganizationService) GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateOrganization updates an organization's name.
func (s *OrganizationService) UpdateOrganization(ctx context.Context, id string, name string) error {
	return s.repo.Update(ctx, id, name)
}
