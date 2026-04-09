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

// UserService provides business logic for user operations.
type UserService struct {
	// repo is the user repository instance.
	repo *repository.UserRepository
}

// NewUserService creates a new UserService instance.
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser creates a new user with hashed password.
func (s *UserService) CreateUser(ctx context.Context, name, password, email string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := models.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
	return s.repo.Create(ctx, user)
}

// GetUserByID retrieves a user by their ID.
func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetById(ctx, id)
}

// GetAllUsers retrieves all users.
func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}

// Login authenticates a user and returns a JWT token.
func (s *UserService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "","",err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "","", errors.New("invalid credentials")
	}
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return "","", err
	}
	refreshToken:= utils.GenerateRefreshToken()

	err = s.repo.StoreRefreshToken(ctx,user.ID,refreshToken)
	if err != nil {
		return "", "",err
	}
	return accessToken, refreshToken, nil
}

func(s *UserService)ValidateRefreshToken(ctx context.Context,token string)(uuid.UUID,error){
	userID, err:= s.repo.GetUserByRefreshToken(ctx,token)
	if err!=nil{
		return uuid.Nil, errors.New("invalid or expired refresh token")
	}
	return userID,nil
}
