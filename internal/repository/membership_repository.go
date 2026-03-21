package repository

import (
	"context"
	"errors"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MembershipRepository struct {
	DB *pgxpool.Pool
}

func NewMembershipRepository(db *pgxpool.Pool) *MembershipRepository{
	return &MembershipRepository{DB: db}
}

func( r *MembershipRepository) Create(ctx context.Context, m models.Membership)error{
	query:=`
	INSERT INTO memberships (id,user_id, org_id, role)
	VALUES ($1,$2,$3,$4)
	`
	_,err:=r.DB.Exec(ctx,query,m.ID,m.UserID,m.OrgID,m.Role)
	return err
}

func (r *MembershipRepository) GetMembersByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Member, error) {
	query := `
	SELECT u.id, u.name, u.email, m.role
	FROM memberships m
	JOIN users u ON m.user_id = u.id
	WHERE m.org_id = $1
	`

	rows, err := r.DB.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.Member

	for rows.Next() {
		var m models.Member
		err := rows.Scan(
			&m.UserID,
			&m.Name,
			&m.Email,
			&m.Role,
		)
		if err != nil {
			return nil, err
		}

		members = append(members, m)
	}

	return members, nil
}

func (r *MembershipRepository) RemoveMember(ctx context.Context, userID, orgID uuid.UUID) error {
	query := `
	DELETE FROM memberships
	WHERE user_id = $1 AND org_id = $2
	`

	cmdTag, err := r.DB.Exec(ctx, query, userID, orgID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("membership not found")
	}

	return nil
}

func (r *MembershipRepository) UpdateRole(ctx context.Context, userID, orgID uuid.UUID, role string) error {
	query := `
	UPDATE memberships
	SET role = $1
	WHERE user_id = $2 AND org_id = $3
	`

	cmdTag, err := r.DB.Exec(ctx, query, role, userID, orgID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("membership not found")
	}

	return nil
}

func (r *MembershipRepository) GetOrgsByUser(ctx context.Context, userID uuid.UUID) ([]models.UserOrg, error) {
	query := `
	SELECT o.id, o.name, m.role
	FROM memberships m
	JOIN organizations o ON m.org_id = o.id
	WHERE m.user_id = $1
	`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []models.UserOrg

	for rows.Next() {
		var o models.UserOrg
		err := rows.Scan(&o.OrgID, &o.Name, &o.Role)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, o)
	}

	return orgs, nil
}