package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/Dubjay18/scenee/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MovieRepository interface {
	// Create inserts a new movie into the database
	Create(ctx context.Context, movie *models.Movie) error

	// Update updates an existing movie
	Update(ctx context.Context, movie *models.Movie) error

	// Delete soft deletes a movie by ID
	Delete(ctx context.Context, id string) error

	// GetByID retrieves a movie by its UUID
	GetByID(ctx context.Context, id string) (*models.Movie, error)

	// GetByTMDBID retrieves a movie by its TMDB ID
	GetByTMDBID(ctx context.Context, tmdbID int) (*models.Movie, error)

	// Upsert creates a new movie or updates if TMDB ID already exists
	Upsert(ctx context.Context, movie *models.Movie) error

	// List retrieves movies with pagination
	List(ctx context.Context, limit, offset int) ([]models.Movie, error)

	// Search searches for movies by title
	Search(ctx context.Context, query string, limit int) ([]models.Movie, error)

	// GetByIDs retrieves multiple movies by their UUIDs
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Movie, error)

	// ListByYear retrieves movies released in a specific year
	ListByYear(ctx context.Context, year int, limit, offset int) ([]models.Movie, error)
}

type GormMovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *GormMovieRepository {
	return &GormMovieRepository{db: db}
}

func (r *GormMovieRepository) Create(ctx context.Context, movie *models.Movie) error {
	return r.db.WithContext(ctx).Create(movie).Error
}

func (r *GormMovieRepository) Update(ctx context.Context, movie *models.Movie) error {
	return r.db.WithContext(ctx).Model(&models.Movie{}).
		Where("id = ?", movie.ID).
		Updates(movie).Error
}

func (r *GormMovieRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Movie{}).Error
}

func (r *GormMovieRepository) GetByID(ctx context.Context, id string) (*models.Movie, error) {
	var movie models.Movie
	if err := r.db.WithContext(ctx).
		First(&movie, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *GormMovieRepository) GetByTMDBID(ctx context.Context, tmdbID int) (*models.Movie, error) {
	var movie models.Movie
	if err := r.db.WithContext(ctx).
		Where("tmdb_id = ?", tmdbID).
		First(&movie).Error; err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *GormMovieRepository) Upsert(ctx context.Context, movie *models.Movie) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "tmdb_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"title", "year", "poster_url", "backdrop_url", "genres", "runtime", "metadata", "updated_at"}),
		}).
		Create(movie).Error
}

func (r *GormMovieRepository) List(ctx context.Context, limit, offset int) ([]models.Movie, error) {
	var movies []models.Movie
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *GormMovieRepository) Search(ctx context.Context, query string, limit int) ([]models.Movie, error) {
	var movies []models.Movie
	searchPattern := "%" + query + "%"
	if err := r.db.WithContext(ctx).
		Where("title ILIKE ?", searchPattern).
		Order("year DESC").
		Limit(limit).
		Find(&movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *GormMovieRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Movie, error) {
	if len(ids) == 0 {
		return []models.Movie{}, nil
	}

	var movies []models.Movie
	if err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *GormMovieRepository) ListByYear(ctx context.Context, year int, limit, offset int) ([]models.Movie, error) {
	var movies []models.Movie
	if err := r.db.WithContext(ctx).
		Where("year = ?", year).
		Order("title ASC").
		Limit(limit).
		Offset(offset).
		Find(&movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}
