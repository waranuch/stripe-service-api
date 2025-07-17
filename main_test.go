package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"stripe-service/config"
	"stripe-service/internal/handlers"
	"stripe-service/internal/server"
	"stripe-service/internal/service"
)

func TestServerCreation(t *testing.T) {
	// Create a test service and handler
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Create server
	srv := server.NewServer(stripeHandler)

	if srv == nil {
		t.Error("Expected server to be created, got nil")
	}

	// Test that handler is not nil
	if srv.Handler() == nil {
		t.Error("Expected handler to be non-nil")
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

	// Create server
	srv := server.NewServer(stripeHandler)

	// Create a test request
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	srv.Handler().ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body contains expected content
	expected := `"status":"healthy"`
	if !contains(rr.Body.String(), expected) {
		t.Errorf("Handler returned unexpected body: got %v want to contain %v", rr.Body.String(), expected)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	// Create a test service and handler
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Create server (which includes logging middleware)
	srv := server.NewServer(stripeHandler)

	// Create a test request
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler with middleware
	srv.Handler().ServeHTTP(rr, req)

	// Check that the request was processed (middleware didn't block it)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Middleware blocked request: got %v want %v", status, http.StatusOK)
	}
}

func TestCORSMiddleware(t *testing.T) {
	// Create a test service and handler
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Create server (which includes CORS middleware)
	srv := server.NewServer(stripeHandler)

	// Test OPTIONS request
	req, err := http.NewRequest("OPTIONS", "/api/v1/customers", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	// Check CORS headers
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected CORS header Access-Control-Allow-Origin to be '*', got %v", rr.Header().Get("Access-Control-Allow-Origin"))
	}

	if rr.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("Expected CORS header Access-Control-Allow-Methods to be set")
	}

	if rr.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Error("Expected CORS header Access-Control-Allow-Headers to be set")
	}

	// OPTIONS request should return 200
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("OPTIONS request returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestResponseWriterWrapper(t *testing.T) {
	// Create a test service and handler
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Create server
	srv := server.NewServer(stripeHandler)

	// Create a test request
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler
	srv.Handler().ServeHTTP(rr, req)

	// The response writer wrapper should capture the status code correctly
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Response writer wrapper didn't capture status correctly: got %v want %v", status, http.StatusOK)
	}
}

func TestServerIntegration(t *testing.T) {
	// Create a test service and handler
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Create server
	srv := server.NewServer(stripeHandler)

	// Test various endpoints
	endpoints := []struct {
		method   string
		path     string
		expected int
	}{
		{"GET", "/api/v1/health", http.StatusOK},
		{"GET", "/api/v1/customers", http.StatusInternalServerError}, // Should fail without proper Stripe key
		{"OPTIONS", "/api/v1/customers", http.StatusOK},              // Should succeed with CORS
	}

	for _, endpoint := range endpoints {
		req, err := http.NewRequest(endpoint.method, endpoint.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rr, req)

		if status := rr.Code; status != endpoint.expected {
			t.Errorf("Endpoint %s %s returned wrong status code: got %v want %v",
				endpoint.method, endpoint.path, status, endpoint.expected)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		(len(s) > len(substr) && contains(s[1:], substr))
}
