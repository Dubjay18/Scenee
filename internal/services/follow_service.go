package services

import (
	"context"

	"github.com/Dubjay18/scenee/internal/models"
	"github.com/Dubjay18/scenee/internal/repositories"
)

type FollowService struct {
	follows repositories.FollowRepository
}

func NewFollowService(follows repositories.FollowRepository) *FollowService {
	return &FollowService{follows: follows}
}

func (s *FollowService) Follow(ctx context.Context, followerID, followeeID string) error {
	return s.follows.Follow(ctx, followerID, followeeID)
}

func (s *FollowService) Unfollow(ctx context.Context, followerID, followeeID string) error {
	return s.follows.Unfollow(ctx, followerID, followeeID)
}

func (s *FollowService) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	return s.follows.IsFollowing(ctx, followerID, followeeID)
}

func (s *FollowService) GetFollowers(ctx context.Context, userID string) ([]models.User, error) {
	return s.follows.GetFollowers(ctx, userID)
}

func (s *FollowService) GetFollowing(ctx context.Context, userID string) ([]models.User, error) {
	return s.follows.GetFollowing(ctx, userID)
}
