package repositories

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Dubjay18/scenee/internal/models"
	"github.com/google/uuid"
)

type WatchlistRepository interface {
	Create(ctx context.Context, watchlist *models.Watchlist) error
	Update(ctx context.Context, watchlist *models.Watchlist) error
	Delete(ctx context.Context, id, owner string) error
	Save(ctx context.Context, userID, watchlistID string) error
	Unsave(ctx context.Context, userID, watchlistID string) error
	GetByID(ctx context.Context, id string) (*models.Watchlist, error)
	GetBySlug(ctx context.Context, slug string) (*models.Watchlist, error)
	ListByOwner(ctx context.Context, owner string) ([]models.Watchlist, error)
	ListPublicByOwner(ctx context.Context, owner string) ([]models.Watchlist, error)
	EnsureOwner(ctx context.Context, watchlistID, owner string) error
	AddItem(ctx context.Context, item *models.WatchlistItem, owner string) error
	RemoveItem(ctx context.Context, watchlistID, itemID, owner string) error
	Like(ctx context.Context, userID, watchlistID string) error
	Unlike(ctx context.Context, userID, watchlistID string) error
	Top(ctx context.Context, window string, limit int) ([]models.Watchlist, error)
}

type GormWatchlistRepository struct {
	db *gorm.DB
}

func NewWatchlistRepository(db *gorm.DB) *GormWatchlistRepository {
	return &GormWatchlistRepository{db: db}
}

func (r *GormWatchlistRepository) Save(ctx context.Context, userID, watchlistID string) error {
	// First, check if the user has already saved this watchlist
	var watchlist models.Watchlist
	if err := r.db.WithContext(ctx).Where("id = ?", watchlistID).First(&watchlist).Error; err != nil {
		return err
	}

	// Check if userID is already in SavedBy
	for _, id := range watchlist.SavedBy {
		if id == userID {
			// Already saved, return without error
			return nil
		}
	}

	// Add userID to SavedBy array and increment SaveCount
	return r.db.WithContext(ctx).Model(&models.Watchlist{}).
		Where("id = ?", watchlistID).
		Updates(map[string]interface{}{
			"saved_by":   gorm.Expr("saved_by || ?::jsonb", `"`+userID+`"`),
			"save_count": gorm.Expr("save_count + 1"),
		}).Error
}

func (r *GormWatchlistRepository) Unsave(ctx context.Context, userID, watchlistID string) error {
	// First, check if the watchlist exists and user has saved it
	var watchlist models.Watchlist
	if err := r.db.WithContext(ctx).Where("id = ?", watchlistID).First(&watchlist).Error; err != nil {
		return err
	}

	// Check if userID is in SavedBy
	found := false
	for _, id := range watchlist.SavedBy {
		if id == userID {
			found = true
			break
		}
	}

	if !found {
		// Not saved, return without error
		return nil
	}

	// Remove userID from SavedBy array and decrement SaveCount
	// Using jsonb_set to remove the element
	return r.db.WithContext(ctx).Model(&models.Watchlist{}).
		Where("id = ?", watchlistID).
		Updates(map[string]interface{}{
			"saved_by":   gorm.Expr("(SELECT jsonb_agg(elem) FROM jsonb_array_elements(saved_by) elem WHERE elem::text != ?)", `"`+userID+`"`),
			"save_count": gorm.Expr("GREATEST(save_count - 1, 0)"),
		}).Error
}

func (r *GormWatchlistRepository) Create(ctx context.Context, watchlist *models.Watchlist) error {
	return r.db.WithContext(ctx).Create(watchlist).Error
}

func (r *GormWatchlistRepository) Update(ctx context.Context, watchlist *models.Watchlist) error {
	return r.db.WithContext(ctx).Model(&models.Watchlist{}).Where("id = ? AND owner_id = ?", watchlist.ID, watchlist.OwnerID).Updates(map[string]any{
		"title":       watchlist.Title,
		"description": watchlist.Description,
	}).Error
}

func (r *GormWatchlistRepository) Delete(ctx context.Context, id, owner string) error {
	return r.db.WithContext(ctx).Where("id = ? AND owner_id = ?", id, owner).Delete(&models.Watchlist{}).Error
}

func (r *GormWatchlistRepository) GetByID(ctx context.Context, id string) (*models.Watchlist, error) {
	var watchlist models.Watchlist
	if err := r.db.WithContext(ctx).Preload("Items", func(tx *gorm.DB) *gorm.DB { return tx.Order("position ASC") }).First(&watchlist, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &watchlist, nil
}

func (r *GormWatchlistRepository) GetBySlug(ctx context.Context, slug string) (*models.Watchlist, error) {
	var watchlist models.Watchlist
	if err := r.db.WithContext(ctx).Preload("Items", func(tx *gorm.DB) *gorm.DB { return tx.Order("position ASC") }).Where("slug = ? AND visibility = 'public'", slug).First(&watchlist).Error; err != nil {
		return nil, err
	}
	return &watchlist, nil
}

func (r *GormWatchlistRepository) ListByOwner(ctx context.Context, owner string) ([]models.Watchlist, error) {
	var out []models.Watchlist
	if err := r.db.WithContext(ctx).Where("owner_id = ?", owner).Order("updated_at DESC").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *GormWatchlistRepository) ListPublicByOwner(ctx context.Context, owner string) ([]models.Watchlist, error) {
	var out []models.Watchlist
	if err := r.db.WithContext(ctx).Where("owner_id = ? AND is_public = TRUE", owner).Order("updated_at DESC").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *GormWatchlistRepository) EnsureOwner(ctx context.Context, watchlistID, owner string) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Watchlist{}).Where("id = ? AND owner_id = ?", watchlistID, owner).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *GormWatchlistRepository) AddItem(ctx context.Context, item *models.WatchlistItem, owner string) error {
	if err := r.EnsureOwner(ctx, item.WatchlistID.String(), owner); err != nil {
		return err
	}
	var pos int
	if err := r.db.WithContext(ctx).Model(&models.WatchlistItem{}).Where("watchlist_id = ?", item.WatchlistID).Select("COALESCE(MAX(position), -1)+1").Scan(&pos).Error; err != nil {
		return err
	}
	item.Position = pos
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GormWatchlistRepository) RemoveItem(ctx context.Context, watchlistID, itemID, owner string) error {
	if err := r.EnsureOwner(ctx, watchlistID, owner); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ? AND watchlist_id = ?", itemID, watchlistID).Delete(&models.WatchlistItem{}).Error
}

func (r *GormWatchlistRepository) Like(ctx context.Context, userID, watchlistID string) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "user_id"}, {Name: "watchlist_id"}}, DoNothing: true}).Create(&models.Like{UserID: uuid.MustParse(userID), WatchlistID: uuid.MustParse(watchlistID)}).Error
}

func (r *GormWatchlistRepository) Unlike(ctx context.Context, userID, watchlistID string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND watchlist_id = ?", uuid.MustParse(userID), uuid.MustParse(watchlistID)).Delete(&models.Like{}).Error
}

func (r *GormWatchlistRepository) Top(ctx context.Context, window string, limit int) ([]models.Watchlist, error) {
	var out []models.Watchlist
	q := r.db.WithContext(ctx).Table("watchlists w").Select("w.*").Joins("LEFT JOIN likes l ON l.watchlist_id = w.id").Where("w.is_public = TRUE")
	switch window {
	case "week":
		q = q.Where("l.created_at >= NOW() - interval '7 days'")
	case "month":
		q = q.Where("l.created_at >= NOW() - interval '30 days'")
	}
	if err := q.Group("w.id").Order("COUNT(l.id) DESC, w.updated_at DESC").Limit(limit).Scan(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}
