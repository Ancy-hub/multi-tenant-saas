package repository

import (
	"context"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
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
	INSERT INTO organizations (id, name)
	Values($1,$2)
	`
	_, err := r.DB.Exec(ctx, query, org.ID, org.Name)
	return err
}

// GetAll retrieves all organizations from the database.
func (r *OrganizationRepository) GetAll(ctx context.Context) ([]models.Organization, error) {
	query := `
	SELECT id, name, created_at
	From organizations
	Order BY created_at desc
	`
	rows, err := r.DB.Query(ctx, query)
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
	SELECT id, name, created_at
	FROM organizations
	WHERE id=$1
	`
	var org models.Organization
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&org.ID,
		&org.Name,
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
