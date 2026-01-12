package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port string
	Env  string

	// Database
	DatabaseURL string

	// Redis
	RedisURL string

	// Stripe
	StripeAPIKey        string
	StripeWebhookSecret string

	// Swish
	SwishAPIURL        string
	SwishCertPath      string
	SwishKeyPath       string
	SwishWebhookSecret string

	// Auth Service
	AuthServiceURL string

	// CORS
	AllowedOrigins []string
}

func Load() (*Config, error) {
	// Load .env file if it exists (for local development)
	godotenv.Load()

	cfg := &Config{
		Port:                getEnv("PORT", "8080"),
		Env:                 getEnv("ENV", "development"),
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		RedisURL:            getEnv("REDIS_URL", "redis://localhost:6379"),
		StripeAPIKey:        getEnv("STRIPE_API_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		SwishAPIURL:         getEnv("SWISH_API_URL", "https://mss.cpc.getswish.net"),
		SwishCertPath:       getEnv("SWISH_CERT_PATH", ""),
		SwishKeyPath:        getEnv("SWISH_KEY_PATH", ""),
		SwishWebhookSecret:  getEnv("SWISH_WEBHOOK_SECRET", ""),
		AuthServiceURL:      getEnv("AUTH_SERVICE_URL", "https://auth.vibeoholic.com"),
		AllowedOrigins:      parseCSV(getEnv("ALLOWED_ORIGINS", "http://localhost:3000")),
	}

	// Validate required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.StripeAPIKey == "" {
		return nil, fmt.Errorf("STRIPE_API_KEY is required")
	}
	if cfg.StripeWebhookSecret == "" {
		return nil, fmt.Errorf("STRIPE_WEBHOOK_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseCSV(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
