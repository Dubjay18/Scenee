package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/Dubjay18/scenee/internal/services"
	"github.com/Dubjay18/scenee/internal/validate"
)

type DiscoverHandler struct {
	Service *services.WatchlistService
}

func NewDiscoverHandler(s *services.WatchlistService) *DiscoverHandler {
	return &DiscoverHandler{Service: s}
}

// Routes mounts the discover routes
func (h *DiscoverHandler) Routes(r chi.Router) {
	r.Get("/trending", h.getTrending)
	r.Get("/new", h.getNew)
}

// getTrending handles GET /v1/discover/trending?window=7d|30d
// Returns trending movies based on time window
func (h *DiscoverHandler) getTrending(w http.ResponseWriter, r *http.Request) {
	type queryT struct {
		Window string `validate:"omitempty,oneof=7d 30d day week month"`
		Page   int    `validate:"omitempty,gte=1,lte=1000"`
	}

	window := r.URL.Query().Get("window")
	if window == "" {
		window = "week" // default
	}
	// Convert 7d to week, 30d to month for consistency
	if window == "7d" {
		window = "week"
	} else if window == "30d" {
		window = "month"
	}

	q := queryT{Window: window, Page: 1}
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

	// Use the existing Feed functionality with type=trending
	opts := services.FeedOptions{
		Type:   "trending",
		Window: q.Window,
		Page:   q.Page,
	}

	body, _, err := h.Service.FetchFeed(r.Context(), opts)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

// getNew handles GET /v1/discover/new
// Returns newly released movies
func (h *DiscoverHandler) getNew(w http.ResponseWriter, r *http.Request) {
	type queryT struct {
		Page   int    `validate:"omitempty,gte=1,lte=1000"`
		Genre  string `validate:"omitempty"`
		Region string `validate:"omitempty,len=2"`
	}

	q := queryT{Page: 1}
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			q.Page = page
		}
	}
	q.Genre = r.URL.Query().Get("genre")
	q.Region = r.URL.Query().Get("region")

	if errs := validate.Map(q); errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errs)
		return
	}

	// Use discover endpoint with sort by release date
	opts := services.FeedOptions{
		Type:   "discover",
		Page:   q.Page,
		Genre:  q.Genre,
		Region: q.Region,
		SortBy: "release_date.desc",
	}

	body, _, err := h.Service.FetchFeed(r.Context(), opts)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

// Mount returns a function that adds the routes under the given router
func (h *DiscoverHandler) Mount() func(r chi.Router) {
	return func(r chi.Router) { h.Routes(r) }
}
