package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/Dubjay18/scenee/internal/auth"
	"github.com/Dubjay18/scenee/internal/services"
	"github.com/Dubjay18/scenee/internal/validate"
)

type FeedHandler struct {
	Service *services.WatchlistService
}

func NewFeedHandler(s *services.WatchlistService) *FeedHandler {
	return &FeedHandler{Service: s}
}

// Routes mounts the feed routes
func (h *FeedHandler) Routes(r chi.Router) {
	r.Get("/", h.getFeed)
}

// getFeed handles GET /v1/feed → following + recommended (paginated)
// This returns a personalized feed for the authenticated user
func (h *FeedHandler) getFeed(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	type queryT struct {
		Page  int `validate:"omitempty,gte=1,lte=100"`
		Limit int `validate:"omitempty,gte=1,lte=100"`
	}

	q := queryT{Page: 1, Limit: 20}
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			q.Page = page
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			q.Limit = limit
		}
	}

	if errs := validate.Map(q); errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errs)
		return
	}

	// TODO: Implement actual feed service method
	// For now, return a placeholder response
	response := map[string]interface{}{
		"following":   []interface{}{},
		"recommended": []interface{}{},
		"page":        q.Page,
		"limit":       q.Limit,
		"has_more":    false,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// Mount returns a function that adds the routes under the given router
func (h *FeedHandler) Mount() func(r chi.Router) {
	return func(r chi.Router) { h.Routes(r) }
}
