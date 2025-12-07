## Scenee — Copilot Instructions

This file gives targeted, actionable guidance for AI coding agents working on the Scenee codebase (Go backend + Expo client). Keep suggestions focused, cite code examples, and avoid speculative changes to infra or secrets.

- Big picture: The backend is a Go monolith (Go 1.24) exposing a REST API. Routing is handled with chi; the app follows a handlers -> services -> repositories layering. Key wiring happens in `cmd/api/main.go` where repositories, services, and handlers are instantiated and routes are mounted.

- Key files to reference:
  - `cmd/api/main.go` — application bootstrap and routes mounting (source of truth for endpoints and middleware).
  - `internal/handlers/*.go` — HTTP handlers; most handlers expose a `Routes(r chi.Router)` or handler methods (e.g. `internal/handlers/auth.go`).
  - `internal/services/*.go` — business logic; services call repositories and return `domain` types.
  - `internal/repositories/*_repository.go` — GORM-backed data access.
  - `internal/auth` — JWT verifier and auth middleware used in `main.go` (cookie name `access_token`).
  - `internal/ai` & `internal/services/ai_service.go` — Gemini AI integration and wrapper service.
  - `go.mod` — Go module / required libraries (chi, gorm, jwt, etc.).

- Dev workflows (explicit):
  - Local env: copy `.env.example` -> `.env` and fill required variables: DATABASE_URL or MIGRATION_URL, JWT_SECRET, TMDB_API_KEY, GEMINI_API_KEY, ENSEND_PROJECT_ID/SECRET.
  - Migrations: the project uses `goose` and Makefile targets. Use `make migrate-up` after setting `MIGRATION_URL` when running migrations against the DB directly.
  - Run the server: `make run` (the Makefile calls `go run` with the correct env and flags). The README lists the same quick-start steps.
  - Build/test: prefer `go build ./...` and `make test` for quick validation.

- Project-specific conventions and patterns:
  - Route registration is centralized in `cmd/api/main.go`. Add routes by creating a handler with a `Routes` method and hook it into the mounter function in `main.go`.
  - Handlers accept services or DB handles in constructors, e.g. `handlers.NewAuthHandler(service)`; prefer dependency injection via constructors already used in `main.go`.
  - Services return domain types (`internal/domain`) instead of GORM models where possible — convert using helpers like `domain.UserFromModel`.
  - Repositories are thin wrappers around GORM; follow existing naming and signature patterns in `internal/repositories/*_repository.go`.
  - Auth: JWT tokens are generated in `internal/services/auth_service.go` and set as an `access_token` cookie in `internal/handlers/auth.go`. Middleware is applied via `auth.NewJWTVerifier(cfg.JWTSecret)` in `main.go`.
  - Sensitive data and keys must come from env vars; never hard-code secrets in code or examples.

- Integration points to be aware of:
  - TMDb: proxied by `internal/tmdb` client; used by watchlist/movie handlers.
  - Google Gemini: AI integration lives in `internal/ai` and used by `internal/services/ai_service.go` and the `/ai/ask` endpoint.
  - EnSend: used for transactional emails (notification service created in `internal/services/auth_service.go`).

- When suggesting code changes, follow this checklist:
  1. Update or add a handler in `internal/handlers` and add its route in `cmd/api/main.go`'s mounter.
 2. Add service logic in `internal/services` and keep it focused (no direct HTTP logic in services).
 3. Use or extend repositories in `internal/repositories` for DB access; return domain types.
 4. Add or update unit tests near the package (use existing test conventions; run `make test`).

- Examples to cite in PRs or suggestions:
  - Wiring example (see `cmd/api/main.go`): repositories.NewUserRepository -> services.NewUserService -> handlers.NewUserHandler
  - Auth flow: `AuthHandler.login` sets a cookie name `access_token`; JWT created in `AuthService.generateJWT`.

- Quick notes for reviewers/agents:
  - The codebase targets Go 1.24; keep syntax and module changes compatible.
  - Migrations use `goose` and the Makefile—prefer `make migrate-*` targets rather than custom scripts.
  - Be conservative with database schema changes; migrations are required for runtime to remain consistent.

If any of this is unclear or you need more examples (specific handler/service patterns, or test patterns), ask and I will add short, concrete examples from the repository.
