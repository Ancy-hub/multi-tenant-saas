package repository

import (
	"context"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	DB *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool)*TaskRepository{
	return &TaskRepository{DB:db}
}
func (r *TaskRepository) Create(ctx context.Context, t models.Task) error {
	query := `
	INSERT INTO tasks (id, project_id, title, description, status, assigned_to, created_by)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	`

	_, err := r.DB.Exec(ctx,
		query,
		t.ID,
		t.ProjectID,
		t.Title,
		t.Description,
		t.Status,
		t.AssignedTo,
		t.CreatedBy,
	)

	return err
}

func (r *TaskRepository) GetByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]models.Task, error) {
	query := `
	SELECT id, project_id, title, description, status, assigned_to, created_by, created_at
	FROM tasks
	WHERE project_id = $1 AND deleted_at IS NULL
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.Query(ctx, query, projectID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var t models.Task
		err := rows.Scan(
			&t.ID,
			&t.ProjectID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.AssignedTo,
			&t.CreatedBy,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, t models.Task) error {
	query := `
	UPDATE tasks
	SET title=$1, description=$2, status=$3, assigned_to=$4, updated_at=NOW()
	WHERE id=$5 AND deleted_at IS NULL
	`

	_, err := r.DB.Exec(ctx,
		query,
		t.Title,
		t.Description,
		t.Status,
		t.AssignedTo,
		t.ID,
	)

	return err
}
func (r *TaskRepository) Delete(ctx context.Context, taskID uuid.UUID) error {
	query := `
	UPDATE tasks SET deleted_at = NOW()
	WHERE id=$1 AND deleted_at IS NULL
	`

	_, err := r.DB.Exec(ctx, query, taskID)
	return err
}