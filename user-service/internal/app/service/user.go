package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AnnaKhairetdinova/user-service/internal/domain"
	"github.com/AnnaKhairetdinova/user-service/internal/infrastructure/repository"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, name string, email string) (domain.User, error)
	GetUserByUUID(ctx context.Context, uuid uuid.UUID) (domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
}

type userService struct {
	ur repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{ur: ur}
}

func (us *userService) CreateUser(ctx context.Context, name string, email string) (domain.User, error) {
	if name == "" {
		return domain.User{}, fmt.Errorf("name is required")
	}

	if !strings.Contains(email, "@") {
		return domain.User{}, fmt.Errorf("invalid email")
	}

	user := domain.User{
		UUID:      uuid.New(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}

	createdUser, err := us.ur.Create(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}

func (us *userService) GetUserByUUID(ctx context.Context, uuid uuid.UUID) (domain.User, error) {
	user, err := us.ur.GetByUUID(ctx, uuid)
	if err != nil {
		return domain.User{}, fmt.Errorf("user not found")
	}

	return user, nil
}

func (us *userService) ListUsers(ctx context.Context) ([]domain.User, error) {
	return us.ur.List(ctx)
}
