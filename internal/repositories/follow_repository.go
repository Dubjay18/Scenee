package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Dubjay18/scenee/internal/models"
)

type FollowRepository interface {
	Follow(ctx context.Context, followerID, followeeID string) error
	Unfollow(ctx context.Context, followerID, followeeID string) error
	IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error)
	GetFollowers(ctx context.Context, userID string) ([]models.User, error)
	GetFollowing(ctx context.Context, userID string) ([]models.User, error)
}

type GormFollowRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *GormFollowRepository {
	return &GormFollowRepository{db: db}
}

func (r *GormFollowRepository) Follow(ctx context.Context, followerID, followeeID string) error {
	follow := models.Follow{FollowerID: parseUUID(followerID), FolloweeID: parseUUID(followeeID)}
	return r.db.WithContext(ctx).Create(&follow).Error
}

func (r *GormFollowRepository) Unfollow(ctx context.Context, followerID, followeeID string) error {
	return r.db.WithContext(ctx).Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&models.Follow{}).Error
}

func (r *GormFollowRepository) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Follow{}).Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Count(&count).Error
	return count > 0, err
}

func (r *GormFollowRepository) GetFollowers(ctx context.Context, userID string) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Joins("JOIN follows ON follows.follower_id = users.id").Where("follows.followee_id = ?", userID).Find(&users).Error
	return users, err
}

func (r *GormFollowRepository) GetFollowing(ctx context.Context, userID string) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Joins("JOIN follows ON follows.followee_id = users.id").Where("follows.follower_id = ?", userID).Find(&users).Error
	return users, err
}

func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
