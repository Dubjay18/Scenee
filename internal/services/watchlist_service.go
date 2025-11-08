package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/yourname/moodle/internal/cache"
	"github.com/yourname/moodle/internal/models"
	"github.com/yourname/moodle/internal/repositories"
	"github.com/yourname/moodle/internal/tmdb"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type WatchlistService struct {
	watchlists repositories.WatchlistRepository
	tmdb       *tmdb.Client
	feedCache  *cache.TTLCache[string, []byte]
}

func NewWatchlistService(repo repositories.WatchlistRepository, tmdbClient *tmdb.Client) *WatchlistService {
	return &WatchlistService{
		watchlists: repo,
		tmdb:       tmdbClient,
		feedCache:  cache.NewTTL[string, []byte](60 * time.Second),
	}
}

func (s *WatchlistService) SearchMovies(ctx context.Context, query string, page int) (*tmdb.SearchMoviesResponse, error) {
	return s.tmdb.SearchMovies(ctx, query, page)
}

func (s *WatchlistService) GetMovie(ctx context.Context, id int64) (*tmdb.Movie, error) {
	return s.tmdb.GetMovie(ctx, id)
}

func (s *WatchlistService) TrendingWatchlists(ctx context.Context, window string, limit int) ([]models.Watchlist, error) {
	return s.watchlists.Top(ctx, window, limit)
}

func (s *WatchlistService) GetWatchlist(ctx context.Context, id, requester string) (*models.Watchlist, error) {
	wl, err := s.watchlists.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !wl.IsPublic && wl.OwnerID != requester {
		return nil, ErrForbidden
	}
	return wl, nil
}

func (s *WatchlistService) ListByOwner(ctx context.Context, owner, requester string) ([]models.Watchlist, error) {
	if requester != "" && owner == requester {
		return s.watchlists.ListByOwner(ctx, owner)
	}
	return s.watchlists.ListPublicByOwner(ctx, owner)
}

func (s *WatchlistService) CreateWatchlist(ctx context.Context, owner string, watchlist *models.Watchlist) error {
	if owner == "" {
		return ErrUnauthorized
	}
	watchlist.OwnerID = owner
	return s.watchlists.Create(ctx, watchlist)
}

func (s *WatchlistService) UpdateWatchlist(ctx context.Context, owner, id string, updater func(existing *models.Watchlist)) (*models.Watchlist, error) {
	if owner == "" {
		return nil, ErrUnauthorized
	}
	existing, err := s.watchlists.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing.OwnerID != owner {
		return nil, ErrForbidden
	}
	if updater != nil {
		updater(existing)
	}
	if err := s.watchlists.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *WatchlistService) DeleteWatchlist(ctx context.Context, owner, id string) error {
	if owner == "" {
		return ErrUnauthorized
	}
	return s.watchlists.Delete(ctx, id, owner)
}

func (s *WatchlistService) AddItem(ctx context.Context, owner, watchlistID string, tmdbID int64, notes string) (*models.WatchlistItem, error) {
	if owner == "" {
		return nil, ErrUnauthorized
	}
	movie, err := s.tmdb.GetMovie(ctx, tmdbID)
	if err != nil {
		return nil, err
	}
	item := &models.WatchlistItem{
		WatchlistID: watchlistID,
		TMDBID:      tmdbID,
		Title:       movie.Title,
		PosterPath:  movie.PosterPath,
		ReleaseDate: movie.ReleaseDate,
		Notes:       notes,
	}
	if err := s.watchlists.AddItem(ctx, item, owner); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *WatchlistService) RemoveItem(ctx context.Context, owner, watchlistID, itemID string) error {
	if owner == "" {
		return ErrUnauthorized
	}
	return s.watchlists.RemoveItem(ctx, watchlistID, itemID, owner)
}

func (s *WatchlistService) Like(ctx context.Context, owner, watchlistID string) error {
	if owner == "" {
		return ErrUnauthorized
	}
	return s.watchlists.Like(ctx, owner, watchlistID)
}

func (s *WatchlistService) Unlike(ctx context.Context, owner, watchlistID string) error {
	if owner == "" {
		return ErrUnauthorized
	}
	return s.watchlists.Unlike(ctx, owner, watchlistID)
}

type FeedOptions struct {
	Type   string
	Window string
	Page   int
	Genre  string
	Year   string
	Region string
	SortBy string
}

func (o FeedOptions) cacheKey() string {
	return fmt.Sprintf("type=%s|window=%s|page=%d|genre=%s|year=%s|region=%s|sort=%s", o.Type, o.Window, o.Page, o.Genre, o.Year, o.Region, o.SortBy)
}

func (s *WatchlistService) FetchFeed(ctx context.Context, opts FeedOptions) ([]byte, bool, error) {
	key := opts.cacheKey()
	if key == "" {
		key = "_empty"
	}
	if body, ok := s.feedCache.Get(key); ok {
		return body, true, nil
	}

	var body []byte
	var err error
	switch opts.Type {
	case "trending":
		var res *tmdb.TrendingResponse
		res, err = s.tmdb.TrendingMovies(ctx, opts.Window, opts.Page, opts.Region)
		if err != nil {
			return nil, false, err
		}
		body, err = json.Marshal(res)
		if err != nil {
			return nil, false, err
		}
	default:
		var res *tmdb.DiscoverResponse
		res, err = s.tmdb.DiscoverMovies(ctx, opts.Page, opts.Genre, opts.Year, opts.Region, opts.SortBy)
		if err != nil {
			return nil, false, err
		}
		body, err = json.Marshal(res)
		if err != nil {
			return nil, false, err
		}
	}

	s.feedCache.Set(key, body)
	return body, false, nil
}
