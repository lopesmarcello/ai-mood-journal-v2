// Package config contains configurations variables
package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	AppEnv      string
	RouteLLMKey string
	GroqAPIKey  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		Port:        getEnv("PORT", "8080"),
		AppEnv:      getEnv("APP_ENV", "development"),
		RouteLLMKey: getEnv("ROUTELLM_API_KEY", ""),
		GroqAPIKey:  getEnv("GROQ_API_KEY", ""),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
