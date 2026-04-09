package services

import (
	"context"
	"errors"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/ancy-shibu/multi-tenant-saas/internal/repository"
	"github.com/google/uuid"
)

// MembershipService provides business logic for organization membership operations.
type MembershipService struct {
	// repo is the membership repository instance.
	repo *repository.MembershipRepository
}

// NewMembershipService creates a new MembershipService instance.
func NewMembershipService(repo *repository.MembershipRepository) *MembershipService {
	return &MembershipService{repo: repo}
}

// AddUserToOrg adds a user to an organization with a specified role.
func (s *MembershipService) AddUserToOrg(ctx context.Context, userID uuid.UUID, orgID uuid.UUID, role string) error {
	if role != "admin" && role != "member" {
		return errors.New("invalid role")
	}
	m := models.Membership{
		ID:     uuid.New(),
		UserID: userID,
		OrgID:  orgID,
		Role:   role,
	}
	return s.repo.Create(ctx, m)
}

// GetMembersByOrg retrieves members of an organization with pagination.
func (s *MembershipService) GetMembersByOrg(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]models.Member, error) {
	if orgID == uuid.Nil {
		return nil, errors.New("invalid org id")
	}

	return s.repo.GetMembersByOrg(ctx, orgID, limit, offset)
}

// RemoveMember removes a user from an organization.
func (s *MembershipService) RemoveMember(ctx context.Context, userID, orgID uuid.UUID) error {
	return s.repo.RemoveMember(ctx, userID, orgID)
}

// UpdateRole updates a user's role in an organization.
func (s *MembershipService) UpdateRole(ctx context.Context, userID, orgID uuid.UUID, role string) error {
	if role != "admin" && role != "member" {
		return errors.New("invalid role")
	}

	return s.repo.UpdateRole(ctx, userID, orgID, role)
}

// GetUserOrgs retrieves all organizations a user belongs to.
func (s *MembershipService) GetUserOrgs(ctx context.Context, userID uuid.UUID) ([]models.UserOrg, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetOrgsByUser(ctx, userID)
}

// GetUserRole retrieves the role of a user in a specific organization.
func (s *MembershipService) GetUserRole(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) (string, error) {
	return s.repo.GetUserRole(ctx, userID, orgID)
}
