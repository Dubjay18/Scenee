package services

import (
	"context"

	"github.com/Dubjay18/scenee/internal/models"
	"github.com/Dubjay18/scenee/internal/repositories"
)

type ReviewService struct {
	reviews repositories.ReviewRepository
}

func NewReviewService(reviews repositories.ReviewRepository) *ReviewService {
	return &ReviewService{reviews: reviews}
}

func (s *ReviewService) Create(ctx context.Context, review *models.Review) error {
	return s.reviews.Create(ctx, review)
}

func (s *ReviewService) GetByMovieID(ctx context.Context, movieID string) ([]models.Review, error) {
	return s.reviews.GetByMovieID(ctx, movieID)
}

func (s *ReviewService) GetByUserAndMovie(ctx context.Context, userID, movieID string) (*models.Review, error) {
	return s.reviews.GetByUserAndMovie(ctx, userID, movieID)
}

func (s *ReviewService) Update(ctx context.Context, review *models.Review) error {
	return s.reviews.Update(ctx, review)
}

func (s *ReviewService) Delete(ctx context.Context, id, userID string) error {
	return s.reviews.Delete(ctx, id, userID)
}
