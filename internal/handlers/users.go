package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Dubjay18/scenee/internal/auth"
	"github.com/Dubjay18/scenee/internal/services"
)

type UserHandler struct{ Users *services.UserService }

func NewUserHandler(s *services.UserService) *UserHandler { return &UserHandler{Users: s} }

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	u, err := h.Users.GetByID(r.Context(), uid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}
	_ = json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
		return
	}

	// Only allow updating bio and avatar_url
	allowedFields := map[string]bool{"bio": true, "avatar_url": true}
	filteredUpdates := make(map[string]interface{})
	for k, v := range updates {
		if allowedFields[k] {
			filteredUpdates[k] = v
		}
	}

	if err := h.Users.Update(r.Context(), uid, filteredUpdates); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update user"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "profile updated"})
}
