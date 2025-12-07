package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Dubjay18/scenee/internal/auth"
	"github.com/Dubjay18/scenee/internal/models"
	"github.com/Dubjay18/scenee/internal/services"
	"github.com/Dubjay18/scenee/internal/validate"
)

type WatchlistHandler struct {
	Service *services.WatchlistService
	DB      *gorm.DB
}

func NewWatchlistHandler(s *services.WatchlistService, db *gorm.DB) *WatchlistHandler {
	return &WatchlistHandler{Service: s, DB: db}
}

// Routes is mounted under /watchlists in main.
func (h *WatchlistHandler) Routes(r chi.Router) {
	r.Get("/{id}", h.get)
	r.Get("/", h.listByOwner)
	r.Post("/", h.create)
	r.Patch("/{id}", h.update)
	r.Delete("/{id}", h.delete)
	// items
	r.Post("/{id}/items", h.addItem)
	r.Delete("/{id}/items/{itemId}", h.removeItem)
	// likes
	r.Post("/{id}/like", h.like)
	r.Delete("/{id}/like", h.unlike)
	// save
	r.Post("/{id}/save", h.save)
}

func (h *WatchlistHandler) GetPublic(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	wl, err := h.Service.GetBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "watchlist not found"})
		return
	}
	// Increment view count
	h.DB.WithContext(r.Context()).Model(&models.Watchlist{}).Where("id = ?", wl.ID).Update("view_count", gorm.Expr("view_count + 1"))
	_ = json.NewEncoder(w).Encode(wl)
}

func (h *WatchlistHandler) save(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	wlID := chi.URLParam(r, "id")
	if err := h.Service.SaveWatchlist(r.Context(), uid, wlID); err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			w.WriteHeader(http.StatusUnauthorized)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	// Create notification
	wl, err := h.Service.GetByID(r.Context(), wlID)
	if err == nil && wl.OwnerID != uid {
		// Notify the owner
		actorUUID, _ := uuid.Parse(uid)
		ownerUUID, _ := uuid.Parse(wl.OwnerID)
		notification := models.Notification{
			UserID:   ownerUUID,
			Type:     "save",
			ActorID:  actorUUID,
			EntityID: wl.ID, // the watchlist
		}
		h.DB.WithContext(r.Context()).Create(&notification)
	}
	w.WriteHeader(http.StatusNoContent)
}

// Public: /v1/search/movies
func (h *WatchlistHandler) SearchMovies(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "q is required"})
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	res, err := h.Service.SearchMovies(r.Context(), q, page)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(res)
}

// Public: GET /v1/movies/{id}
// Fetch a single movie from TMDb by its numeric ID.
func (h *WatchlistHandler) Movie(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "id is required"})
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "id must be a positive integer"})
		return
	}

	mv, err := h.Service.GetMovie(r.Context(), int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(mv)
}

// Public (or semi-public): /v1/trending?window=week|month&limit=20
func (h *WatchlistHandler) Trending(w http.ResponseWriter, r *http.Request) {
	type qT struct {
		Window string `validate:"oneof= week month"`
		Limit  int    `validate:"gte=1,lte=100"`
	}
	q := qT{Window: r.URL.Query().Get("window"), Limit: 20}
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			q.Limit = n
		}
	}
	
	if errs := validate.Map(q); errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(errs)
		_ = json.NewEncoder(w).Encode(errs)
		return
	}
	lists, err := h.Service.TrendingWatchlists(r.Context(), q.Window, q.Limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(lists)
}

func (h *WatchlistHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid := auth.UserID(r.Context())
	wl, err := h.Service.GetWatchlist(r.Context(), id, uid)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, gorm.ErrRecordNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	_ = json.NewEncoder(w).Encode(wl)
}

func (h *WatchlistHandler) listByOwner(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	uid := auth.UserID(r.Context())
	if owner == "" {
		if uid == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "owner required"})
			return
		}
		owner = uid
	}
	lists, err := h.Service.ListByOwner(r.Context(), owner, uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(lists)
}

func (h *WatchlistHandler) create(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	type bodyT struct {
		Title       string `validate:"required,min=1,max=200"`
		Description string `validate:"max=1000"`
		IsPublic    bool
	}
	var b bodyT
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errs := validate.Map(b); errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errs)
		return
	}
	wl := &models.Watchlist{Title: b.Title, Description: b.Description, Visibility: models.PublicVisibility}
	if err := h.Service.CreateWatchlist(r.Context(), uid, wl); err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(wl)
}

func (h *WatchlistHandler) update(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	type bodyT struct {
		Title       *string   `validate:"omitempty,min=1,max=200"`
		Description *string   `validate:"omitempty,max=1000"`
		Visibility  *string   `validate:"omitempty,oneof=public private unlisted"`
		Tags        *[]string `validate:"omitempty"`
	}
	var b bodyT
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errs := validate.Map(b); errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errs)
		return
	}
	updated, err := h.Service.UpdateWatchlist(r.Context(), uid, id, func(existing *models.Watchlist) {
		if b.Title != nil {
			existing.Title = *b.Title
		}
		if b.Description != nil {
			existing.Description = *b.Description
		}
		if b.Visibility != nil {
			existing.Visibility = *b.Visibility
		}
		if b.Tags != nil {
			existing.Tags = *b.Tags
		}
	})
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			w.WriteHeader(http.StatusUnauthorized)
		case errors.Is(err, services.ErrForbidden):
			w.WriteHeader(http.StatusForbidden)
		case errors.Is(err, gorm.ErrRecordNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		if err != nil {
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		}
		return
	}
	_ = json.NewEncoder(w).Encode(updated)
}

func (h *WatchlistHandler) delete(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	if err := h.Service.DeleteWatchlist(r.Context(), uid, id); err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			w.WriteHeader(http.StatusUnauthorized)
		case errors.Is(err, gorm.ErrRecordNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		if err != nil {
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WatchlistHandler) addItem(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	wlID := chi.URLParam(r, "id")
	type bodyT struct {
		TMDBID int64  `json:"tmdb_id" validate:"required,gt=0"`
		Notes  string `json:"notes" validate:"max=1000"`
	}
	var b bodyT
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errs := validate.Map(b); errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errs)
		return
	}
	item, err := h.Service.AddItem(r.Context(), uid, wlID, int(b.TMDBID), b.Notes)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			w.WriteHeader(http.StatusUnauthorized)
		case errors.Is(err, services.ErrForbidden):
			w.WriteHeader(http.StatusForbidden)
		case errors.Is(err, gorm.ErrRecordNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(item)
}

func (h *WatchlistHandler) removeItem(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	wlID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")
	if err := h.Service.RemoveItem(r.Context(), uid, wlID, itemID); err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			w.WriteHeader(http.StatusUnauthorized)
		case errors.Is(err, gorm.ErrRecordNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		if err != nil {
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WatchlistHandler) like(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	wlID := chi.URLParam(r, "id")
	if err := h.Service.Like(r.Context(), uid, wlID); err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			w.WriteHeader(http.StatusUnauthorized)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	// Create notification
	wl, err := h.Service.GetByID(r.Context(), wlID)
	if err == nil && wl.OwnerID != uid {
		// Notify the owner
		actorUUID, _ := uuid.Parse(uid)
		ownerUUID, _ := uuid.Parse(wl.OwnerID)
		notification := models.Notification{
			UserID:   ownerUUID,
			Type:     "like",
			ActorID:  actorUUID,
			EntityID: wl.ID, // the watchlist
		}
		h.DB.WithContext(r.Context()).Create(&notification)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WatchlistHandler) unlike(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	wlID := chi.URLParam(r, "id")
	if err := h.Service.Unlike(r.Context(), uid, wlID); err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			w.WriteHeader(http.StatusUnauthorized)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Feed: GET /v1/feed?type=trending|discover&window=day|week&page=1&genre=&year=&region=&sort_by=
// type=trending uses TMDb trending; type=discover uses TMDb discover with filters.
func (h *WatchlistHandler) Feed(w http.ResponseWriter, r *http.Request) {
	type qT struct {
		Type   string `validate:"required,oneof=trending discover"`
		Window string `validate:"omitempty,oneof=day week"`
		Page   int    `validate:"omitempty,gte=1,lte=1000"`
		Genre  string `validate:"omitempty"`
		Year   string `validate:"omitempty"`
		Region string `validate:"omitempty,len=2"`
		SortBy string `validate:"omitempty,oneof=popularity.desc vote_average.desc release_date.desc"`
	}
	q := qT{
		Type:   r.URL.Query().Get("type"),
		Window: r.URL.Query().Get("window"),
		Genre:  r.URL.Query().Get("genre"),
		Year:   r.URL.Query().Get("year"),
		Region: r.URL.Query().Get("region"),
		SortBy: r.URL.Query().Get("sort_by"),
	}
	if p := r.URL.Query().Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			q.Page = n
		}
	}
	if errs := validate.Map(q); errs != nil {
			fmt.Println(errs)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errs)
		return
	}
	opts := services.FeedOptions{
		Type:   q.Type,
		Window: q.Window,
		Page:   q.Page,
		Genre:  q.Genre,
		Year:   q.Year,
		Region: q.Region,
		SortBy: q.SortBy,
	}
	body, _, err := h.Service.FetchFeed(r.Context(), opts)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=60")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

// Mount returns a function that adds the routes under the given router
func (h *WatchlistHandler) Mount() func(r chi.Router) {
	return func(r chi.Router) { h.Routes(r) }
}
