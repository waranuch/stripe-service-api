package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig
	Stripe StripeConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port int
	Host string
}

// StripeConfig holds Stripe-related configuration
type StripeConfig struct {
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
}

// Load loads configuration from environment variables
func Load() *Config {
	config := &Config{
		Server: ServerConfig{
			Port: getEnvAsInt("PORT", 8080),
			Host: getEnv("HOST", "localhost"),
		},
		Stripe: StripeConfig{
			SecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
			PublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
			WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
		},
	}

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
} 