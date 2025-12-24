# Stage 1: Build
FROM docker.io/golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o bin/app cmd/server/main.go

# Stage 2: Runtime
FROM docker.io/alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/bin/app .
EXPOSE 8080
CMD ["./app"]