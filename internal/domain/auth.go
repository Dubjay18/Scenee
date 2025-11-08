package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourname/moodle/internal/models"
)

// RefreshToken represents a refresh token in the domain layer
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	UserAgent string
	IP        string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// ToDomain converts models.RefreshToken to domain.RefreshToken
func (rt *RefreshToken) FromModel(model *models.RefreshToken) *RefreshToken {
	if model == nil {
		return nil
	}
	return &RefreshToken{
		ID:        model.ID,
		UserID:    model.UserID,
		TokenHash: model.TokenHash,
		UserAgent: model.UserAgent,
		IP:        model.IP,
		ExpiresAt: model.ExpiresAt,
		RevokedAt: model.RevokedAt,
		CreatedAt: model.CreatedAt,
	}
}

// ToModel converts domain.RefreshToken to models.RefreshToken
func (rt *RefreshToken) ToModel() *models.RefreshToken {
	if rt == nil {
		return nil
	}
	return &models.RefreshToken{
		ID:        rt.ID,
		UserID:    rt.UserID,
		TokenHash: rt.TokenHash,
		UserAgent: rt.UserAgent,
		IP:        rt.IP,
		ExpiresAt: rt.ExpiresAt,
		RevokedAt: rt.RevokedAt,
		CreatedAt: rt.CreatedAt,
	}
}

// EmailVerification represents an email verification in the domain layer
type EmailVerification struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	CodeHash   string
	SentAt     time.Time
	ConsumedAt *time.Time
}

// FromModel converts models.EmailVerification to domain.EmailVerification
func (ev *EmailVerification) FromModel(model *models.EmailVerification) *EmailVerification {
	if model == nil {
		return nil
	}
	return &EmailVerification{
		ID:         model.ID,
		UserID:     model.UserID,
		CodeHash:   model.CodeHash,
		SentAt:     model.SentAt,
		ConsumedAt: model.ConsumedAt,
	}
}

// ToModel converts domain.EmailVerification to models.EmailVerification
func (ev *EmailVerification) ToModel() *models.EmailVerification {
	if ev == nil {
		return nil
	}
	return &models.EmailVerification{
		ID:         ev.ID,
		UserID:     ev.UserID,
		CodeHash:   ev.CodeHash,
		SentAt:     ev.SentAt,
		ConsumedAt: ev.ConsumedAt,
	}
}

// AuthProvider represents an authentication provider in the domain layer
type AuthProvider struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Provider       string
	ProviderUserID string
}

// FromModel converts models.AuthProvider to domain.AuthProvider
func (ap *AuthProvider) FromModel(model *models.AuthProvider) *AuthProvider {
	if model == nil {
		return nil
	}
	return &AuthProvider{
		ID:             model.ID,
		UserID:         model.UserID,
		Provider:       model.Provider,
		ProviderUserID: model.ProviderUserID,
	}
}

// ToModel converts domain.AuthProvider to models.AuthProvider
func (ap *AuthProvider) ToModel() *models.AuthProvider {
	if ap == nil {
		return nil
	}
	return &models.AuthProvider{
		ID:             ap.ID,
		UserID:         ap.UserID,
		Provider:       ap.Provider,
		ProviderUserID: ap.ProviderUserID,
	}
}
