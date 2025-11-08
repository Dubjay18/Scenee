package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/Dubjay18/scenee/internal/services"
	"github.com/Dubjay18/scenee/internal/validate"
)

type SearchHandler struct {
	WatchlistService *services.WatchlistService
	UserService      *services.UserService
}

func NewSearchHandler(ws *services.WatchlistService, us *services.UserService) *SearchHandler {
	return &SearchHandler{
		WatchlistService: ws,
		UserService:      us,
	}
}

// Routes mounts the search routes
func (h *SearchHandler) Routes(r chi.Router) {
	r.Get("/", h.search)
}

// search handles GET /v1/search?q=...&type=movie|user|watchlist
// Searches across movies, users, or watchlists based on type parameter
func (h *SearchHandler) search(w http.ResponseWriter, r *http.Request) {
	type queryT struct {
		Q    string `validate:"required,min=1"`
		Type string `validate:"required,oneof=movie user watchlist"`
		Page int    `validate:"omitempty,gte=1,lte=100"`
	}

	q := queryT{
		Q:    r.URL.Query().Get("q"),
		Type: r.URL.Query().Get("type"),
		Page: 1,
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			q.Page = page
		}
	}

	if errs := validate.Map(q); errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errs)
		return
	}

	var result interface{}
	var err error

	switch q.Type {
	case "movie":
		// Search movies using TMDb
		result, err = h.WatchlistService.SearchMovies(r.Context(), q.Q, q.Page)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

	case "user":
		// TODO: Implement SearchUsers method in UserService
		// users, searchErr := h.UserService.SearchUsers(r.Context(), q.Q, 20, (q.Page-1)*20)
		result = map[string]interface{}{
			"results": []interface{}{},
			"page":    q.Page,
			"message": "User search not yet implemented",
		}

	case "watchlist":
		// TODO: Implement SearchWatchlists method in WatchlistService
		// watchlists, searchErr := h.WatchlistService.SearchWatchlists(r.Context(), q.Q, 20, (q.Page-1)*20)
		result = map[string]interface{}{
			"results": []interface{}{},
			"page":    q.Page,
			"message": "Watchlist search not yet implemented",
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid type"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	_ = json.NewEncoder(w).Encode(result)
}

// Mount returns a function that adds the routes under the given router
func (h *SearchHandler) Mount() func(r chi.Router) {
	return func(r chi.Router) { h.Routes(r) }
}
