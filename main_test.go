package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"stripe-service/config"
	"stripe-service/internal/handlers"
	"stripe-service/internal/service"
)

func TestSetupRouter(t *testing.T) {
	// Create a test service and handler
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Setup router
	router := setupRouter(stripeHandler)

	if router == nil {
		t.Error("Expected router to be created, got nil")
	}

	// Test that router is not nil
	if router == nil {
		t.Error("Expected router to be non-nil")
	}
}

func TestHealthEndpoint(t *testing.T) {
	// Create a test service and handler
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Setup router
	router := setupRouter(stripeHandler)

	// Create a test request
	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	// Check content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	// Create a simple test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Wrap with logging middleware
	wrappedHandler := loggingMiddleware(testHandler)

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Serve the request
	wrappedHandler.ServeHTTP(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	if body := rr.Body.String(); body != "test response" {
		t.Errorf("Expected body 'test response', got '%s'", body)
	}
}

func TestCorsMiddleware(t *testing.T) {
	// Create a simple test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Wrap with CORS middleware
	wrappedHandler := corsMiddleware(testHandler)

	t.Run("GET request", func(t *testing.T) {
		// Create a test request
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		// Serve the request
		wrappedHandler.ServeHTTP(rr, req)

		// Check CORS headers
		if header := rr.Header().Get("Access-Control-Allow-Origin"); header != "*" {
			t.Errorf("Expected Access-Control-Allow-Origin '*', got '%s'", header)
		}

		if header := rr.Header().Get("Access-Control-Allow-Methods"); header != "GET, POST, PUT, DELETE, OPTIONS" {
			t.Errorf("Expected Access-Control-Allow-Methods 'GET, POST, PUT, DELETE, OPTIONS', got '%s'", header)
		}

		if header := rr.Header().Get("Access-Control-Allow-Headers"); header != "Content-Type, Authorization" {
			t.Errorf("Expected Access-Control-Allow-Headers 'Content-Type, Authorization', got '%s'", header)
		}
	})

	t.Run("OPTIONS request", func(t *testing.T) {
		// Create an OPTIONS request
		req := httptest.NewRequest("OPTIONS", "/test", nil)
		rr := httptest.NewRecorder()

		// Serve the request
		wrappedHandler.ServeHTTP(rr, req)

		// Check that OPTIONS returns 200
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status code %d for OPTIONS, got %d", http.StatusOK, status)
		}

		// Check CORS headers
		if header := rr.Header().Get("Access-Control-Allow-Origin"); header != "*" {
			t.Errorf("Expected Access-Control-Allow-Origin '*', got '%s'", header)
		}
	})
}

func TestResponseWriterWrapper(t *testing.T) {
	rr := httptest.NewRecorder()
	wrapper := &responseWriterWrapper{
		ResponseWriter: rr,
		statusCode:     http.StatusOK,
	}

	// Test WriteHeader
	wrapper.WriteHeader(http.StatusCreated)
	if wrapper.statusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, wrapper.statusCode)
	}

	// Test Write
	testData := []byte("test data")
	n, err := wrapper.Write(testData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
	}
}

func TestServerConfiguration(t *testing.T) {
	// Test that we can create a server configuration
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}

	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)
	router := setupRouter(stripeHandler)

	server := &http.Server{
		Addr:         "localhost:8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if server.Addr != "localhost:8080" {
		t.Errorf("Expected server address 'localhost:8080', got '%s'", server.Addr)
	}

	if server.ReadTimeout != 15*time.Second {
		t.Errorf("Expected ReadTimeout 15s, got %v", server.ReadTimeout)
	}

	if server.WriteTimeout != 15*time.Second {
		t.Errorf("Expected WriteTimeout 15s, got %v", server.WriteTimeout)
	}

	if server.IdleTimeout != 60*time.Second {
		t.Errorf("Expected IdleTimeout 60s, got %v", server.IdleTimeout)
	}
}
