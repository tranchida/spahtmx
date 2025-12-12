# Stage 1: Build
FROM docker.io/golang:1.25-alpine AS builder

WORKDIR /app

# Installer les dépendances système
RUN apk add --no-cache git make

# Copier les fichiers du module Go
COPY go.mod go.sum* ./

# Télécharger les dépendances
RUN go mod download

# Copier le code source
COPY . .

# Générer les templates templ
RUN go install github.com/a-h/templ/cmd/templ@latest && templ generate

# Compiler l'application
RUN go build -o bin/app main.go

# Stage 2: Runtime
FROM docker.io/alpine:latest

WORKDIR /app

# Installer les dépendances runtime (si nécessaire)
RUN apk add --no-cache ca-certificates

# Copier l'application compilée depuis le builder
COPY --from=builder /app/bin/app .

# Copier les ressources statiques
COPY --from=builder /app/static ./static

# Exposer le port (adapter selon vos besoins)
EXPOSE 8080

# Lancer l'application
CMD ["./app"]
