package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/Dubjay18/scenee/internal/auth"
	"github.com/Dubjay18/scenee/internal/models"
	"gorm.io/gorm"
)

type NotificationHandler struct {
	DB *gorm.DB
}

func NewNotificationHandler(db *gorm.DB) *NotificationHandler {
	return &NotificationHandler{DB: db}
}

// Routes mounts the notification routes
func (h *NotificationHandler) Routes(r chi.Router) {
	r.Get("/", h.getNotifications)
	r.Post("/{id}/mark-read", h.markAsRead)
}

// getNotifications handles GET /v1/notifications?unread=true
// Returns notifications for the authenticated user
func (h *NotificationHandler) getNotifications(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid user id"})
		return
	}

	unreadOnly := r.URL.Query().Get("unread") == "true"

	var notifications []models.Notification
	query := h.DB.WithContext(r.Context()).Where("user_id = ?", userUUID)

	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	if err := query.Order("created_at DESC").Limit(100).Find(&notifications).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"notifications": notifications,
		"count":         len(notifications),
	})
}

// markAsRead handles POST /v1/notifications/:id/mark-read
// Marks a notification as read
func (h *NotificationHandler) markAsRead(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid user id"})
		return
	}

	notificationID := chi.URLParam(r, "id")
	notificationUUID, err := uuid.Parse(notificationID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid notification id"})
		return
	}

	// Update the notification, ensuring it belongs to the user
	result := h.DB.WithContext(r.Context()).
		Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationUUID, userUUID).
		Update("is_read", true)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "notification not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Mount returns a function that adds the routes under the given router
func (h *NotificationHandler) Mount() func(r chi.Router) {
	return func(r chi.Router) { h.Routes(r) }
}
