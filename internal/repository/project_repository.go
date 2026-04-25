package repository

import (
	"context"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	DB *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{DB: db}
}

// CREATE PROJECT
func (r *ProjectRepository) Create(ctx context.Context, p models.Project) error {
	query := `
	INSERT INTO projects (id, name, description, org_id, created_by)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.DB.Exec(ctx, query,
		p.ID,
		p.Name,
		p.Description,
		p.OrgID,
		p.CreatedBy,
	)
	return err
}


// GET PROJECTS BY ORG (with pagination)
func (r *ProjectRepository) GetByOrg(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]models.Project, error) {
	query := `
	SELECT id, name, COALESCE(description, ''), org_id, created_by, created_at
	FROM projects
	WHERE org_id = $1 AND deleted_at IS NULL
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var p models.Project
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.OrgID,
			&p.CreatedBy,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

// GetByID retrieves a single project by its ID.
func (r *ProjectRepository) GetByID(ctx context.Context, projectID uuid.UUID) (*models.Project, error) {
	query := `
	SELECT id, name, COALESCE(description, ''), org_id, created_by, created_at
	FROM projects
	WHERE id = $1 AND deleted_at IS NULL
	`

	var p models.Project
	err := r.DB.QueryRow(ctx, query, projectID).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.OrgID,
		&p.CreatedBy,
		&p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
// DELETE PROJECT
func (r *ProjectRepository) Delete(ctx context.Context, projectID uuid.UUID) error {
	query := `
	UPDATE projects
	SET deleted_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := r.DB.Exec(ctx, query, projectID)
	return err
}