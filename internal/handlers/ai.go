package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Dubjay18/scenee/internal/services"
	"github.com/Dubjay18/scenee/internal/validate"
)

type AIHandler struct{ Service *services.AIService }

func NewAIHandler(s *services.AIService) *AIHandler { return &AIHandler{Service: s} }

func (h *AIHandler) Ask(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Query string `json:"query" validate:"required,min=1,max=500"`
	}
	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
		return
	}
	if errs := validate.Map(body); errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errs)
		return
	}
	answer, err := h.Service.Ask(r.Context(), body.Query)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"answer": answer})
}
