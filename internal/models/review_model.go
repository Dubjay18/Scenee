package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	MovieID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"movie_id"`
	Rating    int            `gorm:"not null;check:rating >= 1 AND rating <= 10" json:"rating"`
	Review    string         `gorm:"type:text" json:"review"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
}

func (Review) TableName() string { return "reviews" }
