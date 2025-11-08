package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Like struct {
	CreatedAt   time.Time `json:"created_at"`
	UserID      uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	WatchlistID uuid.UUID `gorm:"type:uuid;index" json:"watchlist_id"`
}

func (Like) TableName() string { return "likes" }
func (l *Like) BeforeCreate(tx *gorm.DB) error {
	// composite PK via unique index; nothing special here.
	return nil
}

type Save struct {
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	WatchlistID uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
}

func (Save) TableName() string { return "saves" }

type Follow struct {
	FollowerID uuid.UUID `gorm:"type:uuid;not null;index"`
	FolloweeID uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt  time.Time `gorm:"not null;default:now()"`
}

func (Follow) TableName() string { return "follows" }
