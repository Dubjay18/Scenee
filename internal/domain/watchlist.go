package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourname/moodle/internal/models"
)

// Watchlist represents a watchlist in the domain layer
type Watchlist struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	OwnerID     string
	Owner       *User
	Slug        string
	Title       string
	Description string
	CoverUrl    string
	LikeCount   int
	SaveCount   int
	ItemCount   int
}

// FromModel converts models.Watchlist to domain.Watchlist
func (w *Watchlist) FromModel(model *models.Watchlist) *Watchlist {
	if model == nil {
		return nil
	}

	var owner *User
	if model.Owner.ID != uuid.Nil {
		owner = UserFromModel(&model.Owner)
	}

	return &Watchlist{
		ID:          model.ID,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		OwnerID:     model.OwnerID,
		Owner:       owner,
		Slug:        model.Slug,
		Title:       model.Title,
		Description: model.Description,
		CoverUrl:    model.CoverUrl,
		LikeCount:   model.LikeCount,
		SaveCount:   model.SaveCount,
		ItemCount:   model.ItemCount,
	}
}

// ToModel converts domain.Watchlist to models.Watchlist
func (w *Watchlist) ToModel() *models.Watchlist {
	if w == nil {
		return nil
	}

	modelWatchlist := &models.Watchlist{
		ID:          w.ID,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
		OwnerID:     w.OwnerID,
		Slug:        w.Slug,
		Title:       w.Title,
		Description: w.Description,
		CoverUrl:    w.CoverUrl,
		LikeCount:   w.LikeCount,
		SaveCount:   w.SaveCount,
		ItemCount:   w.ItemCount,
	}

	if w.Owner != nil {
		modelWatchlist.Owner = *w.Owner.ToModel()
	}

	return modelWatchlist
}

// WatchlistFromModel is a helper function to convert models.Watchlist to domain.Watchlist
func WatchlistFromModel(model *models.Watchlist) *Watchlist {
	if model == nil {
		return nil
	}
	var w Watchlist
	return w.FromModel(model)
}

// WatchlistsFromModel converts a slice of models.Watchlist to domain.Watchlist
func WatchlistsFromModel(modelWatchlists []models.Watchlist) []Watchlist {
	if modelWatchlists == nil {
		return nil
	}
	watchlists := make([]Watchlist, 0, len(modelWatchlists))
	for _, m := range modelWatchlists {
		watchlists = append(watchlists, *WatchlistFromModel(&m))
	}
	return watchlists
}

// WatchlistItem represents a watchlist item in the domain layer
type WatchlistItem struct {
	ID          uuid.UUID
	WatchlistID uuid.UUID
	MovieID     uuid.UUID
	Note        string
	Position    int
	AddedAt     time.Time
}

// FromModel converts models.WatchlistItem to domain.WatchlistItem
func (wi *WatchlistItem) FromModel(model *models.WatchlistItem) *WatchlistItem {
	if model == nil {
		return nil
	}
	return &WatchlistItem{
		ID:          model.ID,
		WatchlistID: model.WatchlistID,
		MovieID:     model.MovieID,
		Note:        model.Note,
		Position:    model.Position,
		AddedAt:     model.AddedAt,
	}
}

// ToModel converts domain.WatchlistItem to models.WatchlistItem
func (wi *WatchlistItem) ToModel() *models.WatchlistItem {
	if wi == nil {
		return nil
	}
	return &models.WatchlistItem{
		ID:          wi.ID,
		WatchlistID: wi.WatchlistID,
		MovieID:     wi.MovieID,
		Note:        wi.Note,
		Position:    wi.Position,
		AddedAt:     wi.AddedAt,
	}
}

// WatchlistItemFromModel is a helper function to convert models.WatchlistItem to domain.WatchlistItem
func WatchlistItemFromModel(model *models.WatchlistItem) *WatchlistItem {
	if model == nil {
		return nil
	}
	var wi WatchlistItem
	return wi.FromModel(model)
}

// WatchlistItemsFromModel converts a slice of models.WatchlistItem to domain.WatchlistItem
func WatchlistItemsFromModel(modelItems []models.WatchlistItem) []WatchlistItem {
	if modelItems == nil {
		return nil
	}
	items := make([]WatchlistItem, 0, len(modelItems))
	for _, m := range modelItems {
		items = append(items, *WatchlistItemFromModel(&m))
	}
	return items
}
