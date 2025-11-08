package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/Dubjay18/scenee/internal/domain"
)

type Client struct {
	APIKey  string
	BaseURL string
	HTTP    *http.Client
}

type Movie struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	Overview     string  `json:"overview"`
	PosterPath   string  `json:"poster_path"`
	BackdropPath string  `json:"backdrop_path"`
	ReleaseDate  string  `json:"release_date"`
	Genres       []Genre `json:"genres"`
	Runtime      int     `json:"runtime"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type SearchMoviesResponse struct {
	Page         int     `json:"page"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
	Results      []Movie `json:"results"`
}

type TrendingResponse struct {
	Page    int     `json:"page"`
	Results []Movie `json:"results"`
}

type DiscoverResponse struct {
	Page    int     `json:"page"`
	Results []Movie `json:"results"`
}

func New(apiKey, base string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: base,
		HTTP:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) SearchMovies(ctx context.Context, query string, page int) (*SearchMoviesResponse, error) {
	u, _ := url.Parse(c.BaseURL + "/search/movie")
	q := u.Query()
	q.Set("api_key", c.APIKey)
	q.Set("query", query)
	if page > 0 {
		q.Set("page", fmt.Sprint(page))
	}
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tmdb status %d", res.StatusCode)
	}
	var out SearchMoviesResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetMovie(ctx context.Context, id int64) (*Movie, error) {
	u := fmt.Sprintf("%s/movie/%d?api_key=%s", c.BaseURL, id, url.QueryEscape(c.APIKey))
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tmdb status %d", res.StatusCode)
	}
	var out Movie
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TrendingMovies gets trending movies for a given window (day|week) and page.
func (c *Client) TrendingMovies(ctx context.Context, window string, page int, region string) (*TrendingResponse, error) {
	if window == "" {
		window = "day"
	}
	u, _ := url.Parse(c.BaseURL + "/trending/movie/" + window)
	q := u.Query()
	q.Set("api_key", c.APIKey)
	if page > 0 {
		q.Set("page", fmt.Sprint(page))
	}
	if region != "" {
		q.Set("region", region)
	}
	u.RawQuery = q.Encode()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tmdb status %d", res.StatusCode)
	}
	var out TrendingResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DiscoverMovies provides a randomized-like feed using discover with sort_by.
// Filters: with_genres, primary_release_year, region, sort_by (popularity.desc|vote_average.desc|release_date.desc)
func (c *Client) DiscoverMovies(ctx context.Context, page int, genre, year, region, sortBy string) (*DiscoverResponse, error) {
	u, _ := url.Parse(c.BaseURL + "/discover/movie")
	q := u.Query()
	q.Set("api_key", c.APIKey)
	if page > 0 {
		q.Set("page", fmt.Sprint(page))
	}
	if genre != "" {
		q.Set("with_genres", genre)
	}
	if year != "" {
		q.Set("primary_release_year", year)
	}
	if region != "" {
		q.Set("region", region)
	}
	if sortBy == "" {
		sortBy = "popularity.desc"
	}
	q.Set("sort_by", sortBy)
	u.RawQuery = q.Encode()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tmdb status %d", res.StatusCode)
	}
	var out DiscoverResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ToDomainMovie(tm *Movie) *domain.Movie {
	if tm == nil {
		return nil
	}

	year := 0
	var releaseDate *time.Time
	if tm.ReleaseDate != "" {
		if t, err := time.Parse("2006-01-02", tm.ReleaseDate); err == nil {
			year = t.Year()
			releaseDate = &t
		} else if len(tm.ReleaseDate) >= 4 {
			if y, err := strconv.Atoi(tm.ReleaseDate[:4]); err == nil {
				year = y
			}
		}
	}

	return &domain.Movie{
		ID:          uuid.New(),
		TMDBID:      int(tm.ID),
		Title:       tm.Title,
		Year:        year,
		ReleaseDate: releaseDate,
		PosterURL:   tm.PosterPath,
		BackdropURL: tm.BackdropPath,
		Genres:      convertGenres(tm.Genres),
		Runtime:     &tm.Runtime,
		Metadata: map[string]interface{}{
			"overview": tm.Overview,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func convertGenres(gs []Genre) []string {
	if len(gs) == 0 {
		return nil
	}
	out := make([]string, 0, len(gs))
	for _, g := range gs {
		out = append(out, g.Name)
	}
	return out
}
