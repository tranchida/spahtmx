# SPA HTMX avec Go

Une application Single Page Application (SPA) moderne utilisant HTMX, Tailwind CSS et Templ.

## 🚀 Fonctionnalités

- **4 pages** : Accueil, Admin, Prix Nobel, À propos
- **Navigation SPA** : Pas de rechargement de page grâce à HTMX
- **Templates Templ** : Rendu côté serveur avec Templ (type-safe Go templates)
- **Design moderne** : Interface responsive avec Tailwind CSS et animations fluides
- **Gestion d'utilisateurs** : Page admin avec liste d'utilisateurs et statistiques
- **Prix Nobel** : Consultation des prix Nobel (données SQLite)
- **API interactive** : Toggle du statut utilisateur avec HTMX
- **Base de données** : Persistance avec SQLite (via Bun)
- **Fichiers statiques embarqués** : Déploiement simplifié avec embed.FS
- **CI/CD** : Pipeline GitHub Actions pour Docker et déploiement automatique

## 📁 Structure du projet

```
.
├── cmd/
│   └── server/
│       └── main.go      # Point d'entrée de l'application
├── internal/
│   ├── adapter/
│   │   ├── database/    # Implémentation des dépôts SQLite
│   │   └── web/         # Handlers Echo, templates et assets statiques
│   │       ├── static/  # Fichiers JS (htmx, tailwind)
│   │       └── templates/ # Templates Templ
│   ├── app/             # Logique métier (Services)
│   ├── domain/          # Modèles et interfaces (Dépôts)
│   └── config/          # Configuration via variables d'environnement
├── .github/workflows/   # CI/CD (Docker publish & Deploy)
├── compose.yaml         # Configuration Docker Compose (Prometheus)
├── Dockerfile           # Build multi-stage pour la production
├── Makefile             # Raccourcis pour le développement
├── go.mod               # Dépendances Go
└── nobel-prize.json     # Données initiales pour le seed
```

## 🛠️ Installation et démarrage

1. Assurez-vous d'avoir Go (1.23+) et Docker installés.

2. Clonez le projet et accédez au répertoire :
```bash
git clone <url-du-repo>
cd spahtmx
```

3. Installez les dépendances Go :
```bash
go mod download
```

4. Générez les templates Templ :
```bash
go tool templ generate
```

5. Lancez le serveur :
```bash
# Avec les variables d'environnement par défaut
go run cmd/server/main.go
```
Ou avec le peuplement de la base de données (Seed) :
```bash
SEED_DB=true go run cmd/server/main.go
```
Ou utilisez Air pour le développement (nécessite l'installation de air) :
```bash
air
```

6. Ouvrez votre navigateur à l'adresse : **http://localhost:8080**

## 🎯 Comment ça fonctionne

### Architecture
L'application suit les principes de la **Clean Architecture** (ou Hexagonale) :
- **Domain** : Entités et interfaces fondamentales.
- **App** : Services orchestrant la logique métier.
- **Adapters** : Implémentations spécifiques (SQLite pour le stockage, Web/Echo pour l'interface).

### Architecture SPA avec HTMX
L'application utilise HTMX pour créer une expérience SPA sans framework JavaScript lourd :
- **Navigation** : Les clics sur les liens déclenchent des requêtes AJAX (`hx-get`) qui ne chargent que le contenu de la page cible (`#main-content`).
- **Historique** : Géré avec `hx-push-url="true"`.
- **Rendu** : Le serveur Echo retourne soit la page complète (premier chargement), soit uniquement le fragment de contenu (navigation HTMX) grâce à une détection des headers HTMX.

### Routes
- `/` : Accueil
- `/admin` : Administration des utilisateurs
- `/prizes` : Liste des prix Nobel (données SQLite)
- `/about` : À propos
- `/api/switch/{id}` : Toggle du statut utilisateur

## 🎨 Développement

Utilisez le Makefile pour les tâches courantes :
```bash
make build   # Compile l'application
make dev     # Lance air pour le rechargement automatique
```

### Configuration
L'application se configure via des variables d'environnement :
- `PORT` : Port d'écoute (défaut : 8080)
- `SEED_DB` : Si "true", remplit la base de données au démarrage

## 📝 Technologies

- **Go 1.23** - Backend robuste
- **Echo** - Framework web performant
- **SQLite** - Base de données relationnelle (via Bun)
- **Templ** - Templates type-safe pour Go
- **HTMX** - Frontend dynamique sans JS complexe
- **Tailwind CSS** - Styling rapide
- **Docker & Docker Compose** - Conteneurisation
- **GitHub Actions** - CI/CD et déploiement continu

## 🚀 Déploiement

Le projet inclut une configuration CI/CD via GitHub Actions (`.github/workflows/docker-publish.yml`) :
1. **Build** : À chaque push sur `master`, une image Docker est construite et poussée sur GitHub Container Registry (GHCR).
2. **Deploy** : L'image est automatiquement déployée sur le serveur cible via SSH.

### Compilation manuelle
```bash
docker build -t spahtmx .
docker run -p 8080:8080 spahtmx
```

## 📄 License

Projet d'exemple - Libre d'utilisation
