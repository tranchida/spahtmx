FROM golang:1.25 AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate templates using the tool directive (Go 1.24+)
RUN go tool templ generate

# Build the application
# CGO_ENABLED=1 is required for go-sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o app main.go

# Runtime stage
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./app"]
