package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Dubjay18/scenee/internal/auth"
	"github.com/Dubjay18/scenee/internal/models"
	"github.com/Dubjay18/scenee/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FollowHandler struct {
	Follows *services.FollowService
	DB      *gorm.DB
}

func NewFollowHandler(s *services.FollowService, db *gorm.DB) *FollowHandler {
	return &FollowHandler{Follows: s, DB: db}
}

func (h *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	followerID := auth.UserID(r.Context())
	followeeID := chi.URLParam(r, "id")
	if followerID == followeeID {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "cannot follow yourself"})
		return
	}
	if err := h.Follows.Follow(r.Context(), followerID, followeeID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to follow"})
		return
	}
	// Create notification
	followerUUID, _ := uuid.Parse(followerID)
	followeeUUID, _ := uuid.Parse(followeeID)
	notification := models.Notification{
		UserID:   followeeUUID,
		Type:     "follow",
		ActorID:  followerUUID,
		EntityID: followerUUID, // or followee? for follow, entity is the follower
	}
	h.DB.WithContext(r.Context()).Create(&notification)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "followed"})
}

func (h *FollowHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	followerID := auth.UserID(r.Context())
	followeeID := chi.URLParam(r, "id")
	if err := h.Follows.Unfollow(r.Context(), followerID, followeeID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to unfollow"})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "unfollowed"})
}

func (h *FollowHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	followers, err := h.Follows.GetFollowers(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get followers"})
		return
	}
	_ = json.NewEncoder(w).Encode(followers)
}

func (h *FollowHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	following, err := h.Follows.GetFollowing(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get following"})
		return
	}
	_ = json.NewEncoder(w).Encode(following)
}
