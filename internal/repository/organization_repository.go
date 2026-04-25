package repository

import (
	"context"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrganizationRepository handles database operations for organizations.
type OrganizationRepository struct {
	// DB is the database connection pool.
	DB *pgxpool.Pool
}

// NewOrganizationRepository creates a new OrganizationRepository instance.
func NewOrganizationRepository(db *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{DB: db}
}

// Create inserts a new organization into the database.
func (r *OrganizationRepository) Create(ctx context.Context, org models.Organization) error {
	query := `
	INSERT INTO organizations (id, name, description)
	Values($1,$2,$3)
	`
	_, err := r.DB.Exec(ctx, query, org.ID, org.Name, org.Description)
	return err
}

// GetByUserID retrieves all organizations a user has access to.
func (r *OrganizationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Organization, error) {
	query := `
	SELECT DISTINCT o.id, o.name, COALESCE(o.description, ''), o.created_at
	FROM organizations o
	LEFT JOIN memberships m ON o.id = m.org_id
	WHERE o.created_by = $1 OR m.user_id = $1
	ORDER BY o.created_at DESC
	`
	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []models.Organization
	for rows.Next() {
		var org models.Organization
		err := rows.Scan(
			&org.ID,
			&org.Name,
			&org.Description,
			&org.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}
	return orgs, nil
}

// GetByID retrieves an organization by its ID.
func (r *OrganizationRepository) GetByID(ctx context.Context, id string) (*models.Organization, error) {
	query := `
	SELECT id, name, COALESCE(description, ''), created_at
	FROM organizations
	WHERE id=$1
	`
	var org models.Organization
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&org.ID,
		&org.Name,
		&org.Description,
		&org.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &org, err
}

// Update modifies an organization's name.
func (r *OrganizationRepository) Update(ctx context.Context, id string, name string) error {
	query := `
	UPDATE organizations
	SET name =$1
	WHERE id = $2
	`
	_, err := r.DB.Exec(ctx, query, name, id)
	return err
}
