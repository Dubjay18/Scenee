package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/Dubjay18/scenee/internal/models"
	"gorm.io/datatypes"
)

// Movie represents a movie in the domain layer
type Movie struct {
	ID          uuid.UUID
	TMDBID      int
	Title       string
	Year        int
	ReleaseDate *time.Time
	PosterURL   string
	BackdropURL string
	Genres      []string
	Runtime     *int
	Metadata    map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SearchResult struct {
	Movies     []Movie
	TotalCount int
	TotalPages int
	Page       int
}

// FromModel converts models.Movie to domain.Movie
func (m *Movie) FromModel(model *models.Movie) *Movie {
	if model == nil {
		return nil
	}

	// Decode genres
	var genres []string
	if len(model.Genres) > 0 {
		_ = json.Unmarshal(model.Genres, &genres)
	}

	// Decode metadata
	var metadata map[string]interface{}
	if len(model.Metadata) > 0 {
		_ = json.Unmarshal(model.Metadata, &metadata)
	}

	return &Movie{
		ID:          model.ID,
		TMDBID:      model.TMDBID,
		Title:       model.Title,
		Year:        model.Year,
		ReleaseDate: model.ReleaseDate,
		PosterURL:   model.PosterURL,
		BackdropURL: model.BackdropURL,
		Genres:      genres,
		Runtime:     model.Runtime,
		Metadata:    metadata,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

// ToModel converts domain.Movie to models.Movie
func (m *Movie) ToModel() *models.Movie {
	if m == nil {
		return nil
	}

	// Encode genres
	var genresJSON datatypes.JSON
	if m.Genres != nil {
		genresBytes, _ := json.Marshal(m.Genres)
		genresJSON = datatypes.JSON(genresBytes)
	}

	// Encode metadata
	var metadataJSON datatypes.JSON
	if m.Metadata != nil {
		metadataBytes, _ := json.Marshal(m.Metadata)
		metadataJSON = datatypes.JSON(metadataBytes)
	}

	return &models.Movie{
		ID:          m.ID,
		TMDBID:      m.TMDBID,
		Title:       m.Title,
		Year:        m.Year,
		ReleaseDate: m.ReleaseDate,
		PosterURL:   m.PosterURL,
		BackdropURL: m.BackdropURL,
		Genres:      genresJSON,
		Runtime:     m.Runtime,
		Metadata:    metadataJSON,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// MovieFromModel is a helper function to convert models.Movie to domain.Movie
func MovieFromModel(model *models.Movie) *Movie {
	if model == nil {
		return nil
	}
	var m Movie
	return m.FromModel(model)
}

// MoviesFromModel converts a slice of models.Movie to domain.Movie
func MoviesFromModel(modelMovies []models.Movie) []Movie {
	if modelMovies == nil {
		return nil
	}
	movies := make([]Movie, 0, len(modelMovies))
	for _, m := range modelMovies {
		movies = append(movies, *MovieFromModel(&m))
	}
	return movies
}
