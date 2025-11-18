package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	time "time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	Router *chi.Mux
}

// NewServer builds the base router and allows callers to mount versioned routes.
func NewServer(mounters ...func(r chi.Router)) *Server {
	r := chi.NewRouter()
	// Core middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// basic health
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Versioned API
	r.Route("/v1", func(v1 chi.Router) {
		for _, m := range mounters {
			if m != nil {
				m(v1)
			}
		}
	})

	// Log all registered routes for visibility when server starts. We collect them,
	// sort them for deterministic logging, and then print.
	type routeEntry struct{ method, path string }
	// Use a map to collect methods per path, then convert to a slice so we can sort
	routeMethods := map[string]map[string]bool{}
	_ = chi.Walk(r, func(method string, routePath string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if _, ok := routeMethods[routePath]; !ok {
			routeMethods[routePath] = map[string]bool{}
		}
		routeMethods[routePath][method] = true
		return nil
	})
	routes := make([]routeEntry, 0)
	for path, methods := range routeMethods {
		// If GET exists, skip HEAD to avoid duplicates
		if methods["GET"] && methods["HEAD"] {
			delete(methods, "HEAD")
		}
		for m := range methods {
			routes = append(routes, routeEntry{method: m, path: path})
		}
	}
	if len(routes) > 0 {
		sort.Slice(routes, func(i, j int) bool {
			if routes[i].path != routes[j].path {
				return routes[i].path < routes[j].path
			}
			return routes[i].method < routes[j].method
		})
		log.Println("registered routes:")
		for _, r := range routes {
			log.Printf("%s %s", r.method, r.path)
		}
	}

	return &Server{Router: r}
}
