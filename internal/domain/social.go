package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/Dubjay18/scenee/internal/models"
)

// Like represents a like in the domain layer
type Like struct {
	CreatedAt   time.Time
	UserID      uuid.UUID
	WatchlistID uuid.UUID
}

// FromModel converts models.Like to domain.Like
func (l *Like) FromModel(model *models.Like) *Like {
	if model == nil {
		return nil
	}
	return &Like{
		CreatedAt:   model.CreatedAt,
		UserID:      model.UserID,
		WatchlistID: model.WatchlistID,
	}
}

// ToModel converts domain.Like to models.Like
func (l *Like) ToModel() *models.Like {
	if l == nil {
		return nil
	}
	return &models.Like{
		CreatedAt:   l.CreatedAt,
		UserID:      l.UserID,
		WatchlistID: l.WatchlistID,
	}
}

// LikeFromModel is a helper function to convert models.Like to domain.Like
func LikeFromModel(model *models.Like) *Like {
	if model == nil {
		return nil
	}
	var l Like
	return l.FromModel(model)
}

// LikesFromModel converts a slice of models.Like to domain.Like
func LikesFromModel(modelLikes []models.Like) []Like {
	if modelLikes == nil {
		return nil
	}
	likes := make([]Like, 0, len(modelLikes))
	for _, m := range modelLikes {
		likes = append(likes, *LikeFromModel(&m))
	}
	return likes
}

// Save represents a save in the domain layer
type Save struct {
	UserID      uuid.UUID
	WatchlistID uuid.UUID
	CreatedAt   time.Time
}

// FromModel converts models.Save to domain.Save
func (s *Save) FromModel(model *models.Save) *Save {
	if model == nil {
		return nil
	}
	return &Save{
		UserID:      model.UserID,
		WatchlistID: model.WatchlistID,
		CreatedAt:   model.CreatedAt,
	}
}

// ToModel converts domain.Save to models.Save
func (s *Save) ToModel() *models.Save {
	if s == nil {
		return nil
	}
	return &models.Save{
		UserID:      s.UserID,
		WatchlistID: s.WatchlistID,
		CreatedAt:   s.CreatedAt,
	}
}

// SaveFromModel is a helper function to convert models.Save to domain.Save
func SaveFromModel(model *models.Save) *Save {
	if model == nil {
		return nil
	}
	var s Save
	return s.FromModel(model)
}

// SavesFromModel converts a slice of models.Save to domain.Save
func SavesFromModel(modelSaves []models.Save) []Save {
	if modelSaves == nil {
		return nil
	}
	saves := make([]Save, 0, len(modelSaves))
	for _, m := range modelSaves {
		saves = append(saves, *SaveFromModel(&m))
	}
	return saves
}

// Follow represents a follow relationship in the domain layer
type Follow struct {
	FollowerID uuid.UUID
	FolloweeID uuid.UUID
	CreatedAt  time.Time
}

// FromModel converts models.Follow to domain.Follow
func (f *Follow) FromModel(model *models.Follow) *Follow {
	if model == nil {
		return nil
	}
	return &Follow{
		FollowerID: model.FollowerID,
		FolloweeID: model.FolloweeID,
		CreatedAt:  model.CreatedAt,
	}
}

// ToModel converts domain.Follow to models.Follow
func (f *Follow) ToModel() *models.Follow {
	if f == nil {
		return nil
	}
	return &models.Follow{
		FollowerID: f.FollowerID,
		FolloweeID: f.FolloweeID,
		CreatedAt:  f.CreatedAt,
	}
}

// FollowFromModel is a helper function to convert models.Follow to domain.Follow
func FollowFromModel(model *models.Follow) *Follow {
	if model == nil {
		return nil
	}
	var f Follow
	return f.FromModel(model)
}

// FollowsFromModel converts a slice of models.Follow to domain.Follow
func FollowsFromModel(modelFollows []models.Follow) []Follow {
	if modelFollows == nil {
		return nil
	}
	follows := make([]Follow, 0, len(modelFollows))
	for _, m := range modelFollows {
		follows = append(follows, *FollowFromModel(&m))
	}
	return follows
}
