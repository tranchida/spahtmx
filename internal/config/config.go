package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	DebugSQL    bool
	SeedDB      bool
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/spahtmx?sslmode=disable"),
		DebugSQL:    getEnv("DEBUG_SQL", "false") == "true",
		SeedDB:      getEnv("SEED_DB", "false") == "true",
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
