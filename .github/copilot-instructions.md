# SpaHTMX AI Coding Guide

## Architecture Overview

This is a **Go web application** using HTMX for SPA-like behavior without heavy JavaScript frameworks. It follows **Clean/Hexagonal Architecture**:

- **Domain** ([internal/domain](internal/domain/)): Core entities (`User`, `Prize`) and repository interfaces.
- **App** ([internal/app](internal/app/)): Business logic services (`UserService`, `PrizeService`, `AuthService`).
- **Adapters**: Implementation-specific code.
  - **Database** ([internal/adapter/database](internal/adapter/database/)): SQLite implementation using **Bun** ORM. Contains `*Bun` structs and domain converters.
  - **Web** ([internal/adapter/web](internal/adapter/web/)): Echo handlers, Templ templates, and static assets.

## HTMX SPA Pattern

The app serves **full pages** on initial load and **fragments** for HTMX navigation:

```go
// internal/adapter/web/handlers.go
func (h *Handler) handlePage(...) {
    // ...
    isHTMXRequest := c.Request().Header.Get("HX-Request") == "true"
    if isHTMXRequest {
        component = fragment  // Nav + content only
    } else {
        component = templates.Base(page, user, contents)  // Full page
    }
    // ...
}
```

- Templates use `hx-get` with `hx-target="#main-content"` and `hx-push-url="true"`.
- Server renders either complete HTML or just the content fragment.
- Auth redirects via `HX-Redirect` header for HTMX requests (see `AuthMiddleware` in [cmd/server/main.go](cmd/server/main.go)).

## Key Technologies

- **Echo v4**: Web framework.
- **Templ**: Type-safe Go templates compiled to `_templ.go` files (never edit these directly).
- **SQLite (via Bun)**: Database. **Note**: `compose.yaml` contains MongoDB, and `config.go` has `MongoDBURL`, but the application currently hardcodes `file:test.db` in `cmd/server/main.go`.
- **Tailwind CSS**: Compiled from `input.css` to [internal/adapter/web/static/css/styles.css](internal/adapter/web/static/css/styles.css).

## Development Workflow

### Templ & Tailwind
```bash
# Generate templates
go tool templ generate        # One-time
go tool templ generate -watch # Watch mode

# Build CSS
tailwindcss -i input.css -o internal/adapter/web/static/css/styles.css --minify
```

### Run Dev Server
```bash
make dev  # Runs templ watch + Air hot reload in parallel
```
Or use `go tool air`.

### Database
The app uses a local SQLite file `test.db`.
To seed the database (users and Nobel prizes):
```bash
SEED_DB=true go run cmd/server/main.go
```
This will create `test.db` if it doesn't exist and populate it.

**Note**: `docker compose up` starts MongoDB/Mongo Express, but they are **currently unused** by the Go application.

## Critical Patterns

### Repository Pattern with Domain Conversion
Database adapters use internal `*Bun` structs with `bun` tags and conversion functions.
Always convert at adapter boundaries - domain layer stays database-agnostic.

```go
// internal/adapter/database/prize_repository.go
func ToPrizeDomain(p PrizeBun) domain.Prize { ... }
func FromPrizeDomain(prize domain.Prize) (*PrizeBun, error) { ... }
```

### Handler Architecture
1. Extract params/form data.
2. Call service methods ([internal/app](internal/app/)).
3. Call `h.handlePage(c, route, templComponent)` to render response.

### Error Handling
Translate domain errors to HTTP errors in handlers using `translateError`.
Define errors in [internal/domain/errors.go](internal/domain/errors.go).

## Configuration
Env vars in [internal/config/config.go](internal/config/config.go):
- `PORT` (default: `8080`)
- `SEED_DB` (default: `false`)

## Static Assets
Static files are **embedded** into the binary using `embed.FS` in [internal/adapter/web/static_fs.go](internal/adapter/web/static_fs.go).
Served at `/static`.

## Testing
```bash
go test ./...
```
Tests use `testify` assertions.
