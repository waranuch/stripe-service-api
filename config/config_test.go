package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
	}{
		{
			name: "default values",
			envVars: map[string]string{
				"PORT":                   "",
				"HOST":                   "",
				"STRIPE_SECRET_KEY":      "",
				"STRIPE_PUBLISHABLE_KEY": "",
				"STRIPE_WEBHOOK_SECRET":  "",
			},
			expected: &Config{
				Server: ServerConfig{
					Port: 8080,
					Host: "localhost",
				},
				Stripe: StripeConfig{
					SecretKey:      "",
					PublishableKey: "",
					WebhookSecret:  "",
				},
			},
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"PORT":                   "9000",
				"HOST":                   "0.0.0.0",
				"STRIPE_SECRET_KEY":      "sk_test_123",
				"STRIPE_PUBLISHABLE_KEY": "pk_test_123",
				"STRIPE_WEBHOOK_SECRET":  "whsec_test_123",
			},
			expected: &Config{
				Server: ServerConfig{
					Port: 9000,
					Host: "0.0.0.0",
				},
				Stripe: StripeConfig{
					SecretKey:      "sk_test_123",
					PublishableKey: "pk_test_123",
					WebhookSecret:  "whsec_test_123",
				},
			},
		},
		{
			name: "invalid port defaults to 8080",
			envVars: map[string]string{
				"PORT":              "invalid",
				"STRIPE_SECRET_KEY": "",
			},
			expected: &Config{
				Server: ServerConfig{
					Port: 8080,
					Host: "localhost",
				},
				Stripe: StripeConfig{
					SecretKey:      "",
					PublishableKey: "",
					WebhookSecret:  "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original environment variables
			originalEnvVars := make(map[string]string)
			for key := range tt.envVars {
				originalEnvVars[key] = os.Getenv(key)
			}

			// Set test environment variables
			for key, value := range tt.envVars {
				if value == "" {
					os.Unsetenv(key)
				} else {
					os.Setenv(key, value)
				}
			}

			// Test Load function
			config := Load()

			// Assertions
			assert.Equal(t, tt.expected, config)

			// Restore original environment variables
			for key, value := range originalEnvVars {
				if value == "" {
					os.Unsetenv(key)
				} else {
					os.Setenv(key, value)
				}
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "returns env value when set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "env_value",
			expected:     "env_value",
		},
		{
			name:         "returns default when env not set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original value
			originalValue := os.Getenv(tt.key)

			// Set test value
			if tt.envValue == "" {
				os.Unsetenv(tt.key)
			} else {
				os.Setenv(tt.key, tt.envValue)
			}

			// Test function
			result := getEnv(tt.key, tt.defaultValue)

			// Assertion
			assert.Equal(t, tt.expected, result)

			// Restore original value
			if originalValue == "" {
				os.Unsetenv(tt.key)
			} else {
				os.Setenv(tt.key, originalValue)
			}
		})
	}
}

func TestGetEnvAsInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int
		envValue     string
		expected     int
	}{
		{
			name:         "returns env value when set and valid",
			key:          "TEST_INT_KEY",
			defaultValue: 8080,
			envValue:     "9000",
			expected:     9000,
		},
		{
			name:         "returns default when env not set",
			key:          "TEST_INT_KEY",
			defaultValue: 8080,
			envValue:     "",
			expected:     8080,
		},
		{
			name:         "returns default when env value is invalid",
			key:          "TEST_INT_KEY",
			defaultValue: 8080,
			envValue:     "invalid",
			expected:     8080,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original value
			originalValue := os.Getenv(tt.key)

			// Set test value
			if tt.envValue == "" {
				os.Unsetenv(tt.key)
			} else {
				os.Setenv(tt.key, tt.envValue)
			}

			// Test function
			result := getEnvAsInt(tt.key, tt.defaultValue)

			// Assertion
			assert.Equal(t, tt.expected, result)

			// Restore original value
			if originalValue == "" {
				os.Unsetenv(tt.key)
			} else {
				os.Setenv(tt.key, originalValue)
			}
		})
	}
}
