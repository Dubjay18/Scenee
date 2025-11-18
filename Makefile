# Use fish-compatible commands

APP_NAME=scenee
PKG=github.com/Dubjay18/scenee
DB_DSN?=$(DATABASE_URL)
MIGRATIONS_DIR=./migrations
MIGRATIONS_TABLE?=goose_db_version

.PHONY: build run test tidy migrate-up migrate-down migrate-status migrate-create migrate-reset migrate-drop-table goose

build:
	GO111MODULE=on go build -o bin/$(APP_NAME) ./cmd/api

run: build
	./bin/$(APP_NAME)

test:
	go test ./...

watch:
		@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi



tidy:
	go mod tidy

# Helper to compute DSN for goose: append flags to avoid pgx stmtcache issues
# Prioritizes MIGRATION_URL over DATABASE_URL for direct DB connections (bypassing poolers like PgBouncer)
define RESOLVE_DSN
DB_DSN="$${MIGRATION_URL:-$${DB_DSN:-$${DATABASE_URL}}}"; \
if [ -z "$$DB_DSN" ] && [ -f .env ]; then \
	DB_DSN="$$(grep -E '^MIGRATION_URL=' .env | head -n1 | cut -d= -f2- | sed -e 's/^"//' -e 's/"$$//')"; \
	if [ -z "$$DB_DSN" ]; then \
		DB_DSN="$$(grep -E '^DATABASE_URL=' .env | head -n1 | cut -d= -f2- | sed -e 's/^"//' -e 's/"$$//')"; \
	fi; \
fi; \
if [ -z "$$DB_DSN" ]; then echo "DB_DSN/MIGRATION_URL/DATABASE_URL not set. Put MIGRATION_URL in .env (direct connection, not pooler) or pass DB_DSN=postgres://user:pass@host:5432/dbname" >&2; exit 2; fi; \
DSN_FOR_GOOSE="$$(echo "$$DB_DSN" | sed -e 's/[&?]pgbouncer=true//g')"; \
case "$$DSN_FOR_GOOSE" in \
*\?*) DSN_FOR_GOOSE="$$DSN_FOR_GOOSE&prefer_simple_protocol=true&statement_cache_mode=none" ;; \
*) DSN_FOR_GOOSE="$$DSN_FOR_GOOSE?prefer_simple_protocol=true&statement_cache_mode=none" ;; \
esac; \
command -v goose >/dev/null 2>&1 || go install github.com/pressly/goose/v3/cmd/goose@latest;
endef

migrate-up:
	@$(RESOLVE_DSN) \
	goose -table $(MIGRATIONS_TABLE) -dir $(MIGRATIONS_DIR) postgres "$$DSN_FOR_GOOSE" up

migrate-down:
	@$(RESOLVE_DSN) \
	goose -table $(MIGRATIONS_TABLE) -dir $(MIGRATIONS_DIR) postgres "$$DSN_FOR_GOOSE" down

migrate-status:
	@$(RESOLVE_DSN) \
	goose -table $(MIGRATIONS_TABLE) -dir $(MIGRATIONS_DIR) postgres "$$DSN_FOR_GOOSE" status

migrate-create:
	@if test -z "$(name)"; then echo "Usage: make migrate-create name=<name>"; exit 1; fi
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

migrate-reset:
	@$(RESOLVE_DSN) \
	goose -table $(MIGRATIONS_TABLE) -dir $(MIGRATIONS_DIR) postgres "$$DSN_FOR_GOOSE" reset

# Drop the goose migrations table (requires psql). Handy when the table is left in a bad state.
migrate-drop-table:
	@DB_DSN="$${DB_DSN:-$${DATABASE_URL}}"; \
if [ -z "$$DB_DSN" ] && [ -f .env ]; then DB_DSN="$$(grep -E '^DATABASE_URL=' .env | head -n1 | cut -d= -f2- | sed -e 's/^"//' -e 's/"$$//')"; fi; \
if ! command -v psql >/dev/null 2>&1; then echo "psql not found. Install PostgreSQL client tools or drop table $(MIGRATIONS_TABLE) manually." >&2; exit 2; fi; \
psql "$$DB_DSN" -c 'DROP TABLE IF EXISTS $(MIGRATIONS_TABLE);'

goose:
	@command -v goose >/dev/null 2>&1 || go install github.com/pressly/goose/v3/cmd/goose@latest
