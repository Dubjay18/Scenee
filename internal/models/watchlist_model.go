package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Watchlist struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OwnerID string `gorm:"type:uuid;index" json:"owner_id"`
	Owner   User   `gorm:"foreignKey:OwnerID"`
	Slug    string `gorm:"type:citext;uniqueIndex;not null" json:"slug"`

	Title       string   `gorm:"not null" json:"title"`
	Description string   `json:"description"`
	CoverUrl    string   `json:"cover_url"`
	LikeCount   int      `gorm:"default:0" json:"like_count"`
	SaveCount   int      `gorm:"default:0" json:"save_count"`
	ItemCount   int      `gorm:"default:0" json:"item_count"`
	Visibility  string   `gorm:"type:text;not null;check:visibility IN ('public','private','unlisted');default:'private'" json:"visibility"`
	SavedBy     []string `gorm:"type:jsonb;default:'[]'" json:"-"`
}

type WatchlistItem struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WatchlistID uuid.UUID `gorm:"type:uuid;not null;index"`
	MovieID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Note        string    `gorm:"type:text"`
	Position    int       `gorm:"not null;index"`
	AddedAt     time.Time `gorm:"not null;default:now()"`
}
