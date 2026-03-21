package repository

import (
	"context"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository{
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(ctx context.Context, user models.User)error{
	query:=`
	INSERT INTO users(id,email,password_hash,name)
	VALUES ($1,$2,$3,$4)
	`
	 _,err:=r.DB.Exec(ctx,query,user.ID,user.Email, user.PasswordHash, user.Name)
	 return err
}

func (r *UserRepository) GetById(ctx context.Context, id string)(*models.User,error){
	query:=`
	SELECT id, name, email, created_at
	From users
	Where id=$1
	`
	var user models.User
	err:=r.DB.QueryRow(ctx,query,id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if err!=nil{
		return nil,err
	}
	return &user,err
}

func ( r *UserRepository) GetAll(ctx context.Context)([] models.User,error){
	query:=`
	SELECT id, name, email, created_at
	FROM users
	ORDER by created_at desc
	`
	rows,err:=r.DB.Query(ctx,query)
	if err!= nil{
		return nil,err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next(){
		var user models.User
		err:=rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
		)
		if err!=nil{
			return nil,err
		}

		users = append(users, user)
	}
	return users,nil
}

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