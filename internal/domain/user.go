package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/Dubjay18/scenee/internal/models"
)

// User represents a user in the domain layer
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Bio       string    `json:"bio"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	AvatarUrl string    `json:"avatar_url"`
}

// FromModel converts models.User to domain.User
func (u *User) FromModel(model *models.User) *User {
	if model == nil {
		return nil
	}
	return &User{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Bio:       model.Bio,
		Email:     model.Email,
		Username:  model.Username,
		Password:  model.Password,
		AvatarUrl: model.AvatarUrl,
	}
}

// ToModel converts domain.User to models.User
func (u *User) ToModel() *models.User {
	if u == nil {
		return nil
	}
	return &models.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Bio:       u.Bio,
		Email:     u.Email,
		Username:  u.Username,
		Password:  u.Password,
		AvatarUrl: u.AvatarUrl,
	}
}

// UserFromModel is a helper function to convert models.User to domain.User
func UserFromModel(model *models.User) *User {
	if model == nil {
		return nil
	}
	return &User{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Bio:       model.Bio,
		Email:     model.Email,
		Username:  model.Username,
		Password:  model.Password,
		AvatarUrl: model.AvatarUrl,
	}
}

// UsersFromModel converts a slice of models.User to domain.User
func UsersFromModel(modelUsers []models.User) []User {
	if modelUsers == nil {
		return nil
	}
	users := make([]User, 0, len(modelUsers))
	for _, m := range modelUsers {
		users = append(users, *UserFromModel(&m))
	}
	return users
}
