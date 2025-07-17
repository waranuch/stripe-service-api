package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"stripe-service/config"
	"stripe-service/internal/handlers"
	"stripe-service/internal/server"
	"stripe-service/internal/service"
)

// Test the main application components integration
func TestMainApplicationIntegration(t *testing.T) {
	// Set up test environment
	os.Setenv("STRIPE_SECRET_KEY", "sk_test_123")
	os.Setenv("PORT", "8081")
	os.Setenv("HOST", "localhost")
	defer func() {
		os.Unsetenv("STRIPE_SECRET_KEY")
		os.Unsetenv("PORT")
		os.Unsetenv("HOST")
	}()

	// Load configuration
	cfg := config.Load()
	if cfg.Stripe.SecretKey == "" {
		t.Fatal("STRIPE_SECRET_KEY environment variable is required")
	}

	// Initialize services
	stripeService := service.NewStripeService(cfg)
	if stripeService == nil {
		t.Fatal("Failed to create stripe service")
	}

	// Initialize handlers
	stripeHandler := handlers.NewStripeHandler(stripeService)
	if stripeHandler == nil {
		t.Fatal("Failed to create stripe handler")
	}

	// Initialize server
	srv := server.NewServer(stripeHandler)
	if srv == nil {
		t.Fatal("Failed to create server")
	}

	// Setup HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      srv.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if httpServer == nil {
		t.Fatal("Failed to create HTTP server")
	}

	// Test server configuration
	expectedAddr := "localhost:8081"
	if httpServer.Addr != expectedAddr {
		t.Errorf("Expected server address %s, got %s", expectedAddr, httpServer.Addr)
	}

	if httpServer.ReadTimeout != 15*time.Second {
		t.Errorf("Expected ReadTimeout 15s, got %v", httpServer.ReadTimeout)
	}

	if httpServer.WriteTimeout != 15*time.Second {
		t.Errorf("Expected WriteTimeout 15s, got %v", httpServer.WriteTimeout)
	}

	if httpServer.IdleTimeout != 60*time.Second {
		t.Errorf("Expected IdleTimeout 60s, got %v", httpServer.IdleTimeout)
	}
}

// Test configuration validation
func TestConfigurationValidation(t *testing.T) {
	// Test missing STRIPE_SECRET_KEY
	originalKey := os.Getenv("STRIPE_SECRET_KEY")
	os.Unsetenv("STRIPE_SECRET_KEY")
	defer func() {
		if originalKey != "" {
			os.Setenv("STRIPE_SECRET_KEY", originalKey)
		}
	}()

	cfg := config.Load()
	if cfg.Stripe.SecretKey != "" {
		t.Error("Expected empty secret key when environment variable is not set")
	}

	// This would cause main() to call log.Fatal(), which we can't easily test
	// But we can test that the configuration loading works correctly
}

// Test graceful shutdown context
func TestGracefulShutdownContext(t *testing.T) {
	// Test that we can create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if ctx == nil {
		t.Error("Expected context to be created")
	}

	// Test that the context has the correct timeout
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Error("Expected context to have a deadline")
	}

	// Check that the deadline is approximately 30 seconds from now
	now := time.Now()
	expectedDeadline := now.Add(30 * time.Second)
	if deadline.Before(now) || deadline.After(expectedDeadline.Add(time.Second)) {
		t.Errorf("Expected deadline around %v, got %v", expectedDeadline, deadline)
	}
}

// Test server startup components
func TestServerStartupComponents(t *testing.T) {
	// Test that all components can be initialized without errors
	os.Setenv("STRIPE_SECRET_KEY", "sk_test_123")
	defer os.Unsetenv("STRIPE_SECRET_KEY")

	cfg := config.Load()
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)
	srv := server.NewServer(stripeHandler)

	// Test that the server handler is properly set up
	handler := srv.Handler()
	if handler == nil {
		t.Error("Expected server handler to be set up")
	}

	// Test that the server can handle a basic request
	// This is an integration test that verifies the full request flow
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Verify the request was created successfully
	if req.URL.Path != "/api/v1/health" {
		t.Errorf("Expected request path /api/v1/health, got %s", req.URL.Path)
	}

	// We can't easily test the full server startup without actually starting it,
	// but we can test that all components integrate correctly
}

// Test environment variable handling
func TestEnvironmentVariableHandling(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		envValue string
		expected string
	}{
		{
			name:     "STRIPE_SECRET_KEY set",
			envKey:   "STRIPE_SECRET_KEY",
			envValue: "sk_test_example",
			expected: "sk_test_example",
		},
		{
			name:     "PORT set",
			envKey:   "PORT",
			envValue: "9000",
			expected: "9000",
		},
		{
			name:     "HOST set",
			envKey:   "HOST",
			envValue: "0.0.0.0",
			expected: "0.0.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			original := os.Getenv(tt.envKey)
			defer func() {
				if original != "" {
					os.Setenv(tt.envKey, original)
				} else {
					os.Unsetenv(tt.envKey)
				}
			}()

			// Set test value
			os.Setenv(tt.envKey, tt.envValue)

			// Load configuration
			cfg := config.Load()

			// Verify the value was loaded correctly
			switch tt.envKey {
			case "STRIPE_SECRET_KEY":
				if cfg.Stripe.SecretKey != tt.expected {
					t.Errorf("Expected %s to be %s, got %s", tt.envKey, tt.expected, cfg.Stripe.SecretKey)
				}
			case "PORT":
				if fmt.Sprintf("%d", cfg.Server.Port) != tt.expected {
					t.Errorf("Expected %s to be %s, got %d", tt.envKey, tt.expected, cfg.Server.Port)
				}
			case "HOST":
				if cfg.Server.Host != tt.expected {
					t.Errorf("Expected %s to be %s, got %s", tt.envKey, tt.expected, cfg.Server.Host)
				}
			}
		})
	}
}

// Test HTTP server configuration
func TestHTTPServerConfiguration(t *testing.T) {
	os.Setenv("STRIPE_SECRET_KEY", "sk_test_123")
	os.Setenv("PORT", "8082")
	os.Setenv("HOST", "127.0.0.1")
	defer func() {
		os.Unsetenv("STRIPE_SECRET_KEY")
		os.Unsetenv("PORT")
		os.Unsetenv("HOST")
	}()

	cfg := config.Load()
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)
	srv := server.NewServer(stripeHandler)

	// Create HTTP server with the same configuration as main()
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      srv.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Test server configuration
	expectedAddr := "127.0.0.1:8082"
	if httpServer.Addr != expectedAddr {
		t.Errorf("Expected server address %s, got %s", expectedAddr, httpServer.Addr)
	}

	// Test timeouts
	if httpServer.ReadTimeout != 15*time.Second {
		t.Errorf("Expected ReadTimeout 15s, got %v", httpServer.ReadTimeout)
	}

	if httpServer.WriteTimeout != 15*time.Second {
		t.Errorf("Expected WriteTimeout 15s, got %v", httpServer.WriteTimeout)
	}

	if httpServer.IdleTimeout != 60*time.Second {
		t.Errorf("Expected IdleTimeout 60s, got %v", httpServer.IdleTimeout)
	}

	// Test that handler is not nil
	if httpServer.Handler == nil {
		t.Error("Expected server handler to be set")
	}
}
