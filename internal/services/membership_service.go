package services

import (
	"context"
	"errors"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/ancy-shibu/multi-tenant-saas/internal/repository"
	"github.com/google/uuid"
)

type MembershipService struct {
	repo *repository.MembershipRepository
}

func NewMembershipService(repo *repository.MembershipRepository) *MembershipService {
	return &MembershipService{repo: repo}
}

func (s *MembershipService) AddUserToOrg(ctx context.Context,userID, orgID, role string)error{
	if role !="admin" && role !="member"{
		return errors.New("invalid role")
	}
	m:=models.Membership{
		ID: uuid.New().String(),
		UserID: userID,
		OrganizationID: orgID,
		Role: role,
	}
	return s.repo.Create(ctx,m)
}