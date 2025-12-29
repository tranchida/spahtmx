# SpaHTMX AI Coding Guide

## Architecture Overview

This is a **Go web application** using HTMX for SPA-like behavior without heavy JavaScript frameworks. It follows **Clean/Hexagonal Architecture**:

- **Domain** ([internal/domain](internal/domain/)): Core entities (`User`, `Prize`) and repository interfaces
- **App** ([internal/app](internal/app/)): Business logic services (`UserService`, `PrizeService`, `AuthService`)
- **Adapters**: Implementation-specific code
  - **MongoDB** ([internal/adapter/mongodb](internal/adapter/mongodb/)): Repository implementations with `*Mongo` structs and domain converters (`ToPrizeDomain`, `FromPrizeDomain`)
  - **Web** ([internal/adapter/web](internal/adapter/web/)): Echo handlers, Templ templates, and static assets

## HTMX SPA Pattern

The app serves **full pages** on initial load and **fragments** for HTMX navigation:

```go
// handlers.go - handlePage method detects HTMX requests
isHTMXRequest := c.Request().Header.Get("HX-Request") == "true"
if isHTMXRequest {
    component = fragment  // Nav + content only
} else {
    component = templates.Base(page, user, contents)  // Full page
}
```

- Templates use `hx-get` with `hx-target="#main-content"` and `hx-push-url="true"`
- Server renders either complete HTML or just the content fragment
- Auth redirects via `HX-Redirect` header for HTMX requests (see `AuthMiddleware` in [cmd/server/main.go](cmd/server/main.go))

## Key Technologies

- **Echo v4**: Web framework with middleware for logging, gzip, auth
- **Templ**: Type-safe Go templates compiled to `_templ.go` files (never edit these directly)
- **MongoDB v2**: Database with manual type conversions between domain and Mongo models
- **Tailwind CSS**: Compiled from `input.css` to [internal/adapter/web/static/css/styles.css](internal/adapter/web/static/css/styles.css)

## Development Workflow

### Generate Templ Templates
```bash
go tool templ generate        # One-time generation
go tool templ generate -watch # Auto-regenerate (use in dev)
```

Edit `.templ` files in [internal/adapter/web/templates](internal/adapter/web/templates/). Templ generates `_templ.go` files automatically.

### Build Tailwind CSS
```bash
tailwindcss -i input.css -o internal/adapter/web/static/css/styles.css --minify
```

Run this when modifying Tailwind classes in templates.

### Development Mode
```bash
make dev  # Runs templ watch + Air hot reload in parallel
```

Or manually:
```bash
go tool air  # Hot reload on Go file changes (see .air.toml for config)
```

### Database Setup
```bash
docker compose up -d  # Start MongoDB + Mongo Express
```

- MongoDB: `mongodb://root:example@localhost:27017`
- Mongo Express UI: `http://localhost:8081` (user: `mongoexpressuser`, pass: `mongoexpresspass`)

Seed database on startup:
```bash
SEED_DB=true go run cmd/server/main.go
```

Seeds users (username/password: `alice`, `bob`, `charlie` / `password`) and Nobel Prize data from [nobel-prize.json](nobel-prize.json).

## Critical Patterns

### Repository Pattern with Domain Conversion

MongoDB adapters use internal `*Mongo` structs with `bson` tags and conversion functions:

```go
// In prize_repository.go
func ToPrizeDomain(p PrizeMongo) domain.Prize { ... }
func FromPrizeDomain(prize domain.Prize) (*PrizeMongo, error) { ... }
```

Always convert at adapter boundaries - domain layer stays database-agnostic.

### Handler Architecture

All handlers follow this pattern:
1. Extract params/form data from `echo.Context`
2. Call service methods (business logic in `internal/app`)
3. Call `h.handlePage(c, route, templComponent)` to render response

The `handlePage` method:
- Detects HTMX requests via `HX-Request` header
- Manages user context from cookies/session
- Renders full page or fragment based on request type

### Authentication

Simple cookie-based auth (not production-ready):
- Login sets `session` cookie with username
- `AuthMiddleware` checks cookie, redirects unauthenticated users
- For HTMX requests, uses `HX-Redirect` header instead of server-side redirect

### Error Handling

Translate domain errors to HTTP errors in handlers:
```go
func translateError(err error) error {
    if errors.Is(err, domain.ErrUserNotFound) {
        return echo.NewHTTPError(http.StatusNotFound, "User not found")
    }
    // ...
}
```

Define errors in [internal/domain/errors.go](internal/domain/errors.go).

## Configuration

Env vars in [internal/config/config.go](internal/config/config.go):
- `PORT` (default: `8080`)
- `MONGODB_URL` (default: `mongodb://root:example@localhost:27017`)
- `SEED_DB` (default: `false`) - set to `true` to seed database on startup

## Static Assets

Static files are **embedded** into the binary using `embed.FS` in [internal/adapter/web/static_fs.go](internal/adapter/web/static_fs.go):
```go
//go:embed static/*
var StaticFS embed.FS
```

Served via Echo's `StaticFS` at `/static` route. Update files in [internal/adapter/web/static](internal/adapter/web/static/).

## Testing

Run tests:
```bash
go test ./...
```

Tests use `testify` assertions (see `go.mod`). Add tests alongside code files (`*_test.go`).

## Building for Production

```bash
make build  # Outputs binary to bin/app
```

Or use Docker:
```bash
docker build -t spahtmx .
```

The Dockerfile uses multi-stage builds with Templ generation and static embedding.

## Common Tasks

**Add a new page:**
1. Create `.templ` file in [internal/adapter/web/templates](internal/adapter/web/templates/)
2. Add handler method in [internal/adapter/web/handlers.go](internal/adapter/web/handlers.go)
3. Register route in `initWeb()` in [cmd/server/main.go](cmd/server/main.go)
4. Add nav link in [nav.templ](internal/adapter/web/templates/nav.templ)

**Add MongoDB query:**
1. Add method to repository interface in [internal/domain/repositories.go](internal/domain/repositories.go)
2. Implement in `*MongoRepository` ([internal/adapter/mongodb](internal/adapter/mongodb/))
3. Call from service in [internal/app](internal/app/)

**Modify domain model:**
1. Update struct in [internal/domain/models.go](internal/domain/models.go)
2. Update corresponding `*Mongo` struct in adapter
3. Update conversion functions (`To*Domain`, `From*Domain`)
