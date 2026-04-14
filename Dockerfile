# Stage 1: Build
FROM docker.io/golang:1.26.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/app cmd/server/main.go

# Stage 2: Runtime
FROM docker.io/alpine:3.21
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Europe/Zurich /etc/localtime && \
    echo "Europe/Zurich" > /etc/timezone && \
    apk del tzdata && \
    adduser -D -u 1001 appuser
COPY --from=builder /app/bin/app .
COPY --from=builder /app/nobel-prize.json .
USER appuser
EXPOSE 8080
CMD ["./app"]