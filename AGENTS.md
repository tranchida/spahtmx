# AGENTS.md

## Mission rapide
- Ce repo est une SPA server-side Go: Echo + HTMX + Templ + Bun/PostgreSQL.
- Objectif des agents: livrer des changements cohérents avec l'architecture hexagonale (`internal/domain`, `internal/app`, `internal/adapter`).

## Architecture essentielle
- `internal/domain/`: entités + interfaces de dépôts (`UserRepository`, `PrizeRepository`) + erreurs métier.
- `internal/app/`: services métier orchestrant les interfaces domaine (ex: `UserService`, `PrizeService`, `AuthService`).
- `internal/adapter/database/`: implémentations Bun; conversion explicite `To*Domain` / `From*Domain` aux frontières.
- `internal/adapter/web/`: handlers Echo, routes, rendu Templ, assets embarqués via `embed.FS`.
- Composition racine dans `cmd/server/main.go` (config, DB, schema, seed, middleware, routes).

## Flux HTTP/SPA (pattern critique)
- Navigation HTMX via `hx-get`, `hx-target="#content"`, `hx-push-url` dans `internal/adapter/web/templates/nav.templ`.
- `handlePage` (`internal/adapter/web/handlers.go`) détecte `HX-Request`:
  - HTMX: renvoie fragment `Nav + contents`.
  - Non-HTMX: renvoie `Base(...)` complet.
- Auth/redirect:
  - `AuthMiddleware` (`cmd/server/main.go`) utilise `HX-Redirect` pour HTMX, sinon HTTP 303 vers `/login`.
  - Login/logout suivent la même logique dans `HandleLoginPost` / `HandleLogout`.

## Données et persistence
- PostgreSQL via Bun, DSN configurable (`DATABASE_URL`) dans `internal/config/config.go`.
- Schéma créé au démarrage (`createSchema` dans `cmd/server/main.go`).
- Seed optionnel via `SEED_DB=true`:
  - users par défaut (`alice`, `bob`, `charlie`)
  - import Nobel depuis `nobel-prize.json`.

## Workflows dev indispensables
- Commandes utiles (`Makefile`): `make dev`, `make build`, `make templ`, `make air`, `make tailwind`.
- Génération Templ one-shot: `go tool templ generate`.
- Tests: `go test ./...`.
- Prérequis local: PostgreSQL actif (voir `compose.yaml`, ex: `docker compose up -d`).

## Conventions de changement (spécifiques projet)
- Ne pas éditer `*_templ.go`; modifier `*.templ` puis régénérer Templ.
- Ajouter une route via constantes de `handlers.go` (`Route*`) puis enregistrement dans `initWeb`.
- Préserver la séparation hexagonale: pas d'accès DB direct depuis `internal/app`.
- Côté UI HTMX, garder les cibles cohérentes (`#content` pour pages, `#userlist` pour liste admin).
- Pour assets statiques, passer par `internal/adapter/web/static/` (servis via `StaticFS`).

## Points d'intégration externes
- Web: Echo v4 + middleware gzip/logger.
- Templates: `github.com/a-h/templ`.
- ORM: Bun + `pgdriver` PostgreSQL.
- Auth: cookie `session` (24h), mot de passe hashé bcrypt.
