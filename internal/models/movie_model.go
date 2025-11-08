package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Movie struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TMDBID      int            `gorm:"uniqueIndex;not null"`
	Title       string         `gorm:"type:text;not null"`
	Year        int            `gorm:"index"`
	PosterURL   string         `gorm:"type:text"`
	BackdropURL string         `gorm:"type:text"`
	Genres      datatypes.JSON `gorm:"type:jsonb"` // store []string
	Runtime     *int
	Metadata    datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt   time.Time      `gorm:"not null;default:now()"`
	UpdatedAt   time.Time      `gorm:"not null;default:now()"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
