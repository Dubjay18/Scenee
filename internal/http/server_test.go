package httpserver

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

// Basic test to ensure NewServer mounts mounters and registers routes (with /v1 prefix)
func TestNewServer_RegistersRoutes(t *testing.T) {
	s := NewServer(func(r chi.Router) {
		r.Get("/foo", func(w http.ResponseWriter, r *http.Request) {})
		r.Post("/bar", func(w http.ResponseWriter, r *http.Request) {})
	})

	found := map[string]struct{}{}
	_ = chi.Walk(s.Router, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		found[method+" "+route] = struct{}{}
		return nil
	})

	if _, ok := found["GET /v1/foo"]; !ok {
		t.Fatalf("missing GET /v1/foo; routes: %v", found)
	}
	if _, ok := found["POST /v1/bar"]; !ok {
		t.Fatalf("missing POST /v1/bar; routes: %v", found)
	}
}
