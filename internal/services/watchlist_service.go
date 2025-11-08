package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/Dubjay18/scenee/internal/cache"
	"github.com/Dubjay18/scenee/internal/domain"
	"github.com/Dubjay18/scenee/internal/models"
	"github.com/Dubjay18/scenee/internal/repositories"
	"github.com/Dubjay18/scenee/internal/tmdb"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type WatchlistService struct {
	watchlists repositories.WatchlistRepository
	msvc       *MovieService
	feedCache  *cache.TTLCache[string, []byte]
}

func NewWatchlistService(repo repositories.WatchlistRepository, tmdbClient *tmdb.Client) *WatchlistService {
	return &WatchlistService{
		watchlists: repo,
		msvc:       NewMovieService(*tmdbClient, nil),
		feedCache:  cache.NewTTL[string, []byte](60 * time.Second),
	}
}

func (s *WatchlistService) SaveWatchlist(ctx context.Context, userID, watchlistID string) error {
	if userID == "" {
		return ErrUnauthorized
	}
	watchlist, err := s.watchlists.GetByID(ctx, watchlistID)
	if err != nil {
		return err
	}
	if watchlist.OwnerID == userID {
		return errors.New("cannot save your own watchlist")
	}
	return s.watchlists.Save(ctx, userID, watchlistID)
}

func (s *WatchlistService) SearchMovies(ctx context.Context, query string, page int) (*domain.SearchResult, error) {
	return s.msvc.SearchMovies(ctx, query, page)
}

func (s *WatchlistService) GetMovie(ctx context.Context, id int) (*domain.Movie, error) {
	return s.msvc.GetMovieByTMDBID(ctx, id)
}

func (s *WatchlistService) TrendingWatchlists(ctx context.Context, window string, limit int) ([]models.Watchlist, error) {
	return s.watchlists.Top(ctx, window, limit)
}

func (s *WatchlistService) GetWatchlist(ctx context.Context, id, requester string) (*models.Watchlist, error) {
	wl, err := s.watchlists.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if wl.Visibility == models.PrivateVisibility && wl.OwnerID != requester {
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

func (s *WatchlistService) AddItem(ctx context.Context, owner, watchlistID string, tmdbID int, notes string) (*models.WatchlistItem, error) {
	if owner == "" {
		return nil, ErrUnauthorized
	}
	movie, err := s.msvc.GetMovieByTMDBID(ctx, tmdbID)
	if err != nil {
		return nil, err
	}
	item := &models.WatchlistItem{
		WatchlistID: uuid.MustParse(watchlistID),
		MovieID:     movie.ID,
		Note:        notes,
		ID:          uuid.New(),
		Position:    0, // Will be set in repository
		AddedAt:     time.Now(),
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
		res, err = s.msvc.tmdbClient.TrendingMovies(ctx, opts.Window, opts.Page, opts.Region)
		if err != nil {
			return nil, false, err
		}
		body, err = json.Marshal(res)
		if err != nil {
			return nil, false, err
		}
	default:
		var res *tmdb.DiscoverResponse
		res, err = s.msvc.tmdbClient.DiscoverMovies(ctx, opts.Page, opts.Genre, opts.Year, opts.Region, opts.SortBy)
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
