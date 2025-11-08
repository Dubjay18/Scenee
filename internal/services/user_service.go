package services

import (
	"context"

	"github.com/yourname/moodle/internal/models"
	"github.com/yourname/moodle/internal/repositories"
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
