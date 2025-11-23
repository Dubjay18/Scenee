package services

import (
	"context"

	"github.com/Dubjay18/scenee/internal/models"
	"github.com/Dubjay18/scenee/internal/repositories"
)

type UserService struct {
	users repositories.UserRepository
}

func NewUserService(users repositories.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) Upsert(ctx context.Context, user *models.User) error {
	return s.users.Upsert(ctx, user)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.users.GetByEmail(ctx, email)
}

func (s *UserService) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return s.users.Update(ctx, id, updates)
}

func (s *UserService) GetRole(ctx context.Context, id string) (string, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	return user.Role, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.users.Delete(ctx, id)
}
