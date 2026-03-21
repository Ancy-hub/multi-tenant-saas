package services

import (
	"context"
	"errors"

	"github.com/ancy-shibu/multi-tenant-saas/internal/models"
	"github.com/ancy-shibu/multi-tenant-saas/internal/repository"
	"github.com/ancy-shibu/multi-tenant-saas/internal/utils"
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

func( s *UserService) Login(ctx context.Context, email, password string)(string,error){
	user,err:=s.repo.GetByEmail(ctx,email)
	if err!=nil{
		return "",err
	}
	err= bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),[]byte(password))
	if err !=nil{
		return "",errors.New("invalid credentials")
	}
	token,err:=utils.GenerateToken(user.ID)
	if err!=nil{
		return "",err
	}
	return token,nil
}