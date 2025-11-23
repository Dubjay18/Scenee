package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/Dubjay18/scenee/internal/models"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	GetByMovieID(ctx context.Context, movieID string) ([]models.Review, error)
	GetByUserAndMovie(ctx context.Context, userID, movieID string) (*models.Review, error)
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id, userID string) error
}

type GormReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *GormReviewRepository {
	return &GormReviewRepository{db: db}
}

func (r *GormReviewRepository) Create(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *GormReviewRepository) GetByMovieID(ctx context.Context, movieID string) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.WithContext(ctx).Preload("User").Where("movie_id = ?", movieID).Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

func (r *GormReviewRepository) GetByUserAndMovie(ctx context.Context, userID, movieID string) (*models.Review, error) {
	var review models.Review
	err := r.db.WithContext(ctx).Where("user_id = ? AND movie_id = ?", userID, movieID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *GormReviewRepository) Update(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Save(review).Error
}

func (r *GormReviewRepository) Delete(ctx context.Context, id, userID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.Review{}).Error
}
