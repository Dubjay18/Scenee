package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Dubjay18/scenee/internal/ai"
	"github.com/Dubjay18/scenee/internal/auth"
	"github.com/Dubjay18/scenee/internal/handlers"
	httpserver "github.com/Dubjay18/scenee/internal/http"
	"github.com/Dubjay18/scenee/internal/repositories"
	"github.com/Dubjay18/scenee/internal/services"
	"github.com/Dubjay18/scenee/internal/tmdb"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Port                string `envconfig:"PORT" default:"8080"`
	DatabaseURL         string `envconfig:"DATABASE_URL" required:"false"`
	MigrationURL        string `envconfig:"MIGRATION_URL" required:"false"`
	JWTSecret           string `envconfig:"JWT_SECRET" required:"true"`
	ClientURL           string `envconfig:"CLIENT_URL" default:"exp://192.168.0.5:8081/--/auth"`
	TMDBAPIKey          string `envconfig:"TMDB_API_KEY" required:"true"`
	TMDBBaseURL         string `envconfig:"TMDB_BASE_URL" default:"https://api.themoviedb.org/3"`
	GeminiAPIKey        string `envconfig:"GEMINI_API_KEY" required:"true"`
	GeminiModel         string `envconfig:"GEMINI_MODEL" default:"gemini-1.5-flash"`
	EnSendProjectID     string `envconfig:"ENSEND_PROJECT_ID" required:"true"`
	EnSendProjectSecret string `envconfig:"ENSEND_PROJECT_SECRET" required:"true"`
}

func mustLoadEnv() Config {
	_ = godotenv.Load() // Load .env from current working directory
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		log.Fatalf("env error: %v", err)
	}
	// Prefer MIGRATION_URL (direct DB connection) when present. This mirrors the Makefile
	// behavior which uses MIGRATION_URL for direct connections (bypassing poolers like PgBouncer).
	if c.MigrationURL != "" {
		c.DatabaseURL = c.MigrationURL
	}
	if c.DatabaseURL == "" {
		log.Fatalf("env error: DATABASE_URL or MIGRATION_URL must be set")
	}
	return c
}

func mustDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sqlDB, _ := db.DB()
	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("db ping error: %v", err)
	}
	return db
}

func main() {
	cfg := mustLoadEnv()
	db := mustDB(cfg.DatabaseURL)
	tmdbClient := tmdb.New(cfg.TMDBAPIKey, cfg.TMDBBaseURL)
	aiClient := ai.NewGemini(cfg.GeminiAPIKey, cfg.GeminiModel)

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	watchlistRepo := repositories.NewWatchlistRepository(db)
	followRepo := repositories.NewFollowRepository(db)
	reviewRepo := repositories.NewReviewRepository(db)

	// Services
	userService := services.NewUserService(userRepo)
	watchlistService := services.NewWatchlistService(watchlistRepo, tmdbClient)
	aiService := services.NewAIService(aiClient)
	authService := services.NewAuthService(userService, cfg.JWTSecret, cfg.EnSendProjectID, cfg.EnSendProjectSecret)
	followService := services.NewFollowService(followRepo)
	reviewService := services.NewReviewService(reviewRepo)

	// Handlers
	wlHandler := handlers.NewWatchlistHandler(watchlistService, db)
	aiHandler := handlers.NewAIHandler(aiService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	followHandler := handlers.NewFollowHandler(followService, db)
	notificationHandler := handlers.NewNotificationHandler(db)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	discoverHandler := handlers.NewDiscoverHandler(watchlistService)
	adminHandler := handlers.NewAdminHandler(userService)
	statsHandler := handlers.NewStatsHandler(db)

	// Auth middleware
	verifier := auth.NewJWTVerifier(cfg.JWTSecret)

	mounter := func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Get("/search/movies", wlHandler.SearchMovies)
			r.Get("/movies/{id}", wlHandler.Movie)
			r.Get("/feed", wlHandler.Feed)
			r.Get("/watchlists/public/{slug}", wlHandler.GetPublic)
			r.Route("/discover", discoverHandler.Routes)
			r.Post("/ai/ask", aiHandler.Ask)
			// Auth routes (public)
			r.Route("/auth", authHandler.Routes)
		})
		// Authed routes
		r.Group(func(r chi.Router) {
			r.Use(verifier.Middleware)
			r.Get("/me", userHandler.Me)
			r.Patch("/me", userHandler.UpdateMe)
			r.Route("/watchlists", wlHandler.Routes)
			// trending can be public but keep here for now or move above
			r.Get("/trending", wlHandler.Trending)
			r.Route("/users/{id}", func(r chi.Router) {
				r.Post("/follow", followHandler.Follow)
				r.Delete("/follow", followHandler.Unfollow)
				r.Get("/followers", followHandler.GetFollowers)
				r.Get("/following", followHandler.GetFollowing)
			})
			r.Route("/movies/{id}/reviews", func(r chi.Router) {
				r.Get("/", reviewHandler.GetByMovie)
				r.Post("/", reviewHandler.Create)
				r.Put("/", reviewHandler.Update)
				r.Delete("/{reviewID}", reviewHandler.Delete)
			})
			r.Delete("/admin/users/{id}", adminHandler.DeleteUser)
			r.Get("/stats", statsHandler.GetStats)
			r.Route("/notifications", notificationHandler.Routes)
		})
	}

	srv := httpserver.NewServer(mounter)

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, srv.Router); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
