package services

import (
	"context"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/ancy-shibu/multi-tenant-saas/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)	

type UserService struct{
	repo *repository.UserRepository
}

func NewUserService (repo *repository.UserRepository) *UserService{
	return &UserService{repo:repo}
}

func (s *UserService) CreateUser(ctx context.Context, name,password, email string) error{
	hashedPassword,err:=bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err!=nil{
		return err
	}
	user := models.User{
		ID: uuid.New(),
		Name: name,
		Email: email,
		PasswordHash: string(hashedPassword),
	}
	return s.repo.Create(ctx,user)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User,error){
	return s.repo.GetById(ctx,id)
}

func (s *UserService) GetAllUsers(ctx context.Context)([]models.User,error){
	return s.repo.GetAll(ctx)
}