package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	DebugSQL    bool
	SeedDB      bool
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/spahtmx?sslmode=disable"),
		DebugSQL:    getEnv("DEBUG_SQL", "false") == "true",
		SeedDB:      getEnv("SEED_DB", "false") == "true",
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-key-change-me"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		fmt.Printf("Environment variable %s = %s\n", key, value)
		return value
	}
	return fallback
}
