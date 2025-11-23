package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Dubjay18/scenee/internal/auth"
	"github.com/Dubjay18/scenee/internal/services"
	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	UserService *services.UserService
}

func NewAdminHandler(us *services.UserService) *AdminHandler {
	return &AdminHandler{UserService: us}
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	role, err := h.UserService.GetRole(r.Context(), uid)
	if err != nil || role != "admin" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userID := chi.URLParam(r, "id")
	if err := h.UserService.Delete(r.Context(), userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete user"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
