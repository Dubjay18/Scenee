package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourname/moodle/internal/auth"
	"github.com/yourname/moodle/internal/services"
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
