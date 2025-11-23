package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Dubjay18/scenee/internal/models"
	"gorm.io/gorm"
)

type StatsHandler struct {
	DB *gorm.DB
}

func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{DB: db}
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	var userCount, watchlistCount, reviewCount int64

	h.DB.WithContext(r.Context()).Model(&models.User{}).Count(&userCount)
	h.DB.WithContext(r.Context()).Model(&models.Watchlist{}).Count(&watchlistCount)
	h.DB.WithContext(r.Context()).Model(&models.Review{}).Count(&reviewCount)

	stats := map[string]interface{}{
		"total_users":      userCount,
		"total_watchlists": watchlistCount,
		"total_reviews":    reviewCount,
	}

	_ = json.NewEncoder(w).Encode(stats)
}
