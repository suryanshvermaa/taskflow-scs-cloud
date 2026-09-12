package config

import (
	"os"
)

// Config holds all application configuration values loaded from environment variables.
type Config struct {
	Port       string
	DBPath     string
	CORSOrigin string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBPath:     getEnv("DB_PATH", "./data/taskflow.db"),
		CORSOrigin: getEnv("CORS_ORIGIN", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

