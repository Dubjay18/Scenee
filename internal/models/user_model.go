package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Bio       string         `json:"bio"`
	Email     string         `gorm:"uniqueIndex" json:"email"`
	Username  string         `gorm:"uniqueIndex" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	AvatarUrl string         `json:"avatar_url"`
	Role      string         `gorm:"type:text;not null;default:'user';check:role IN ('user','admin')" json:"role"`
}
