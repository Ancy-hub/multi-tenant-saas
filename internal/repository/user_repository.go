package repository

import (
	"context"
	"time"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository handles database operations for users.
type UserRepository struct {
	// DB is the database connection pool.
	DB *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(ctx context.Context, user models.User) error {
	query := `
	INSERT INTO users(id,email,password_hash,name)
	VALUES ($1,$2,$3,$4)
	`
	_, err := r.DB.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, user.Name)
	return err
}

// GetById retrieves a user by their ID.
func (r *UserRepository) GetById(ctx context.Context, id string) (*models.User, error) {
	query := `
	SELECT id, name, email, created_at
	From users
	Where id=$1
	`
	var user models.User
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// GetAll retrieves all users from the database.
func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `
	SELECT id, name, email, created_at
	FROM users
	ORDER by created_at desc
	`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

// GetByEmail retrieves a user by their email address.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	query := `
	SELECT id, name, email, password_hash
	FROM users
	WHERE email = $1
	`

	var user models.User

	err := r.DB.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)

	return user, err
}

func (r *UserRepository) StoreRefreshToken(ctx context.Context, userID uuid.UUID , token string)error{
	query:=`
	INSERT INTO refresh_tokens(id,user_id,token,expires_at)
	VALUES($1,$2,$3,$4)
	`
	_,err:=r.DB.Exec(
		ctx,
		query,
		uuid.New(),
		userID,
		token,
		time.Now().Add(7*24*time.Hour),//7 days
	)
	return err
}

func(r *UserRepository) GetUserByRefreshToken(ctx context.Context, token string)(uuid.UUID, error){
	query:=`
	SELECT user_id
	FROM refresh_tokens
	WHERE token= $1 AND expires_At>NOW()
	`
	var userID uuid.UUID
	err:=r.DB.QueryRow(ctx,query,token).Scan(&userID)
	if err!=nil{
		return uuid.Nil, err
	}
	return userID, nil
}
