package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string    `gorm:"type:text;not null;index"`
	UserAgent string    `gorm:"type:text"`
	IP        string    `gorm:"type:inet"`
	ExpiresAt time.Time `gorm:"not null;index"`
	RevokedAt *time.Time
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

func (RefreshToken) TableName() string { return "auth_refresh_tokens" }

type EmailVerification struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	CodeHash   string    `gorm:"type:text;not null;index"`
	SentAt     time.Time `gorm:"not null;default:now()"`
	ConsumedAt *time.Time
}

func (EmailVerification) TableName() string { return "auth_email_verifications" }

type AuthProvider struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index"`
	Provider       string    `gorm:"type:text;not null;index"` // 'google','apple'
	ProviderUserID string    `gorm:"type:text;not null"`
}

func (AuthProvider) TableName() string { return "auth_providers" }
