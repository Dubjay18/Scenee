# Scenee API (Go)

Backend for a mobile app to create, share, like, and discover movie watchlists. Includes TMDb integration for movie metadata and Gemini-powered AI assistant ("Scenee").

## Stack
- Go 1.24.0
- chi router
- GORM (with PostgreSQL)
- goose for migrations
- TMDb API for movie data
- Google Gemini for AI

## Features
- Auth (JWT-based with custom user registration/login)
- Users & Profiles
- Watchlists (create/update/delete/save)
- Watchlist items (movies from TMDb)
- Likes & Saves
- Trending/top watchlists (weekly/monthly)
- Feed (trending/discover with filters)
- Search via TMDb proxy endpoints
- AI endpoint `/ai/ask` powered by Gemini

## Local setup

1. Prereqs
- Go 1.24.0
- PostgreSQL
- goose installed (`go install github.com/pressly/goose/v3/cmd/goose@latest`)

2. Copy env
```
cp .env.example .env
```

3. Fill env
- DATABASE_URL: PostgreSQL connection string
- JWT_SECRET: secret key for JWT signing
- TMDB_API_KEY: your TMDb API key
- GEMINI_API_KEY: your Google AI API key

4. Run migrations
```
make migrate-up
```

5. Run server
```
make run
```

## Makefile targets
- build, run
- migrate-up, migrate-down, migrate-status, migrate-create name=<name>
- test

## API sketch
- GET /healthz (health check)
- POST /v1/auth/register
- POST /v1/auth/login
- POST /v1/auth/logout
- GET /v1/auth/user
- GET /v1/me
- GET /v1/users/{id}
- POST /v1/watchlists
- GET /v1/watchlists?owner=<id>
- GET /v1/watchlists/{id}
- PATCH /v1/watchlists/{id}
- DELETE /v1/watchlists/{id}
- POST /v1/watchlists/{id}/items
- DELETE /v1/watchlists/{id}/items/{itemId}
- POST /v1/watchlists/{id}/like
- DELETE /v1/watchlists/{id}/like
- POST /v1/watchlists/{id}/save
- GET /v1/trending?window=week|month&limit=20
- GET /v1/feed?type=trending|discover&window=day|week&page=1&genre=&year=&region=&sort_by=
- GET /v1/search/movies?q=...
- GET /v1/movies/{id}
- POST /v1/ai/ask {"query":"..."}

