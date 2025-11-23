package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Dubjay18/scenee/internal/auth"
	"github.com/Dubjay18/scenee/internal/models"
	"github.com/Dubjay18/scenee/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ReviewHandler struct {
	Service *services.ReviewService
}

func NewReviewHandler(s *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{Service: s}
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	movieIDStr := chi.URLParam(r, "id")
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid movie id"})
		return
	}

	var body struct {
		Rating int    `json:"rating" validate:"required,min=1,max=10"`
		Review string `json:"review"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
		return
	}

	userID, _ := uuid.Parse(uid)
	review := &models.Review{
		UserID:  userID,
		MovieID: movieID,
		Rating:  body.Rating,
		Review:  body.Review,
	}

	if err := h.Service.Create(r.Context(), review); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to create review"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) GetByMovie(w http.ResponseWriter, r *http.Request) {
	movieID := chi.URLParam(r, "id")
	reviews, err := h.Service.GetByMovieID(r.Context(), movieID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get reviews"})
		return
	}
	_ = json.NewEncoder(w).Encode(reviews)
}

func (h *ReviewHandler) Update(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var body struct {
		Rating int    `json:"rating" validate:"min=1,max=10"`
		Review string `json:"review"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
		return
	}

	// Get existing review
	existing, err := h.Service.GetByUserAndMovie(r.Context(), uid, chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "review not found"})
		return
	}

	if body.Rating != 0 {
		existing.Rating = body.Rating
	}
	existing.Review = body.Review

	if err := h.Service.Update(r.Context(), existing); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update review"})
		return
	}

	_ = json.NewEncoder(w).Encode(existing)
}

func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	reviewID := chi.URLParam(r, "reviewID")
	if err := h.Service.Delete(r.Context(), reviewID, uid); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete review"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
