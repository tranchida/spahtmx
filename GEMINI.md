# SPA HTMX avec Go - GEMINI.md

This project is a modern Single Page Application (SPA) built with **Go**, **HTMX**, **Tailwind CSS**, and **Templ**. It demonstrates a high-performance, low-complexity approach to web development by leveraging server-side rendering with dynamic client-side updates.

## 🚀 Project Overview

- **Architecture:** Clean/Hexagonal Architecture.
  - `internal/domain`: Core entities (`User`, `Prize`) and repository interfaces.
  - `internal/app`: Business logic services (`UserService`, `PrizeService`).
  - `internal/adapter/database`: PostgreSQL implementation using **Bun ORM**.
  - `internal/adapter/web`: **Echo** handlers and **Templ** templates.
- **Frontend Strategy:** SPA experience using **HTMX** for partial page updates (`hx-get`, `hx-target="#main-content"`, `hx-push-url="true"`).
- **Styling:** **Tailwind CSS** (v3/v4 style via `input.css`).
- **Templates:** **Templ** for type-safe, compiled Go templates.

## 🛠 Building and Running

### Prerequisites
- Go 1.25+
- Docker & Docker Compose (for PostgreSQL)
- `templ` and `air` (installed via `go tool` as defined in `go.mod`)

### Essential Commands
- **Development (Hot Reload):** `make dev` (runs `templ watch` and `air` in parallel).
- **Build Templates:** `go tool templ generate` or `make templ`.
- **Build CSS:** `make tailwind`.
- **Run Tests:** `go test ./...`.
- **Start Database:** `docker compose up -d`.
- **Run App (Manual):** `SEED_DB=true go run cmd/server/main.go`.

## 📏 Development Conventions

### 1. Templ Templates
- **Location:** `internal/adapter/web/templates/*.templ`.
- **Rule:** Never modify generated `*_templ.go` files. Always edit the `.templ` source and regenerate.
- **HTMX Integration:** Use `hx-get` for navigation to ensure only the `#main-content` fragment is swapped when requested via HTMX.

### 2. Clean Architecture Boundaries
- **Domain:** Must remain pure (no infrastructure/web imports).
- **Persistence:** Use `ToDomain()` and `FromDomain()` methods in `internal/adapter/database` to map between database models and domain entities.
- **Handlers:** Should be thin, delegating business logic to the `app` services.

### 3. HTMX Navigation
- The server detects HTMX requests via the `HX-Request` header.
- For HTMX requests, return only the specific fragment/template.
- For direct browser hits, return the full `Base` template wrapping the content.

### 4. Database & Models
- Use **Bun ORM** for database operations.
- Tables are auto-created on startup in `cmd/server/main.go`.
- Seed data is loaded from `nobel-prize.json` if `SEED_DB=true` is set.

## ⚙️ Configuration
Environment variables:
- `PORT`: Server port (default: `8080`).
- `SEED_DB`: Set to `true` to populate the database on startup.
- DSN is currently configured in `cmd/server/main.go`.
