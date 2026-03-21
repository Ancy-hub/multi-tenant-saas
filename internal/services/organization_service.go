package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/ancy-shibu/multi-tenant-saas/internal/repository"
)

type OrganizationService struct {
	repo *repository.OrganizationRepository
}

func NewOrganizationService(repo *repository.OrganizationRepository) *OrganizationService{
	return &OrganizationService{repo: repo}
}

func (s *OrganizationService) CreateOrganization(ctx context.Context, name string)error{
	org:=models.Organization{
		ID: uuid.New(),
		Name: name,
	}

	return s.repo.Create(ctx,org)
}

func (s *OrganizationService) GetOrganizations(ctx context.Context)([]models.Organization,error){
	return s.repo.GetAll(ctx)
}

func (s *OrganizationService) GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error){
	return s.repo.GetByID(ctx,id)
}

func (s *OrganizationService) UpdateOrganization(ctx context.Context, id string, name string) error{
	return s.repo.Update(ctx,id,name)
}