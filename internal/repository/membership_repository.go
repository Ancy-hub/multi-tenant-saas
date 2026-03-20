package repository

import (
	"context"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
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
	_,err:=r.DB.Exec(ctx,query,m.ID,m.UserID,m.OrganizationID,m.Role)
	return err
}