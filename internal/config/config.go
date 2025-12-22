package config

import (
	"os"
)

type Config struct {
	Port       string
	MongoDBURL string
	SeedDB     bool
}

func Load() *Config {
	return &Config{
		Port:       getEnv("PORT", "8765"),
		MongoDBURL: getEnv("MONGODB_URL", "mongodb://root:example@localhost:27017"),
		SeedDB:     getEnv("SEED_DB", "false") == "true",
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
