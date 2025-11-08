package services

import (
	"context"
	"errors"

	"github.com/Dubjay18/scenee/internal/domain"
	"github.com/Dubjay18/scenee/internal/repositories"
	"github.com/Dubjay18/scenee/internal/tmdb"
	"gorm.io/gorm"
)

type IMovieService interface {
	// Define movie-related service methods here
	GetMovieByTMDBID(tmdbID int) (*domain.Movie, error)
	SearchMovies(query string, page int) (*domain.SearchResult, error)
	AddMovie(mov *domain.Movie) (*domain.Movie, error)
	// Add more methods as needed
}

type MovieService struct {
	tmdbClient tmdb.Client
	mrepo      repositories.MovieRepository
}

func NewMovieService(tmdbClient tmdb.Client, mrepo repositories.MovieRepository) *MovieService {
	return &MovieService{
		tmdbClient: tmdbClient,
		mrepo:      mrepo,
	}
}

func (s *MovieService) GetMovieByTMDBID(ctx context.Context, tmdbID int) (*domain.Movie, error) {
	m, err := s.mrepo.GetByTMDBID(ctx, tmdbID)
	if err == gorm.ErrRecordNotFound {
		tm, err := s.tmdbClient.GetMovie(ctx, int64(tmdbID))
		if err != nil {
			return nil, err
		}
		// Add movie to the database
		newMovie := s.tmdbClient.ToDomainMovie(tm)
		// Convert domain.Movie to models.Movie and save
		// (Assuming a ToModel method exists)
		modelMovie := newMovie.ToModel()
		if err := s.mrepo.Create(ctx, modelMovie); err != nil {
			return nil, err
		}
		return newMovie, nil
	}
	if err != nil {
		return nil, err
	}
	domainMovie := (&domain.Movie{}).FromModel(m)
	return domainMovie, nil
}

func (s *MovieService) SearchMovies(ctx context.Context, query string, page int) (*domain.SearchResult, error) {
	sm, err := s.mrepo.Search(ctx, query, page)
	if err == gorm.ErrRecordNotFound {
		resp, err := s.tmdbClient.SearchMovies(ctx, query, page)
		if err != nil {
			return nil, err
		}
		// Convert []tmdb.Movie to []domain.Movie
		movies := make([]domain.Movie, 0, len(resp.Results))
		for i := range resp.Results {
			dm := s.tmdbClient.ToDomainMovie(&resp.Results[i])
			if dm != nil {
				movies = append(movies, *dm)
			}
		}
		return &domain.SearchResult{
			Movies:     movies,
			TotalCount: resp.TotalResults,
			TotalPages: resp.TotalPages,
			Page:       resp.Page,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	response := &domain.SearchResult{
		Movies:     make([]domain.Movie, 0, len(sm)),
		TotalCount: len(sm), // Adjust as needed
		TotalPages: 1,       // Adjust as needed
		Page:       page,
	}
	for _, m := range sm {
		dm := (&domain.Movie{}).FromModel(&m)
		if dm != nil {
			response.Movies = append(response.Movies, *dm)
		}
	}
	return response, nil
}

func (s *MovieService) AddMovie(ctx context.Context, mov *domain.Movie) (*domain.Movie, error) {
	if mov == nil {
		return nil, errors.New("movie cannot be nil")
	}
	// Convert domain.Movie to models.Movie and save
	modelMovie := mov.ToModel()
	if err := s.mrepo.Create(ctx, modelMovie); err != nil {
		return nil, err
	}
	return mov, nil
}
