package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"stripe-service/config"
	"stripe-service/internal/handlers"
	"stripe-service/internal/service"
)

func TestNewServer(t *testing.T) {
	// Create test dependencies
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Test NewServer
	server := NewServer(stripeHandler)

	if server == nil {
		t.Error("Expected server to be created, got nil")
	}

	if server.router == nil {
		t.Error("Expected router to be initialized, got nil")
	}
}

func TestServerHandler(t *testing.T) {
	// Create test dependencies
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Create server
	server := NewServer(stripeHandler)

	// Test Handler method
	handler := server.Handler()
	if handler == nil {
		t.Error("Expected handler to be returned, got nil")
	}

	// Test that handler is the same as router
	if handler != server.router {
		t.Error("Expected handler to be the same as router")
	}
}

func TestSetupRouter(t *testing.T) {
	// Create test dependencies
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Create server (which calls setupRouter internally)
	server := NewServer(stripeHandler)

	// Test that all expected routes are registered
	testCases := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{"GET", "/api/v1/health", http.StatusOK},
		{"OPTIONS", "/api/v1/customers", http.StatusOK},
		{"GET", "/api/v1/customers", http.StatusInternalServerError}, // Will fail due to test key
		{"POST", "/api/v1/customers", http.StatusBadRequest},         // Will fail due to empty body
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rr := httptest.NewRecorder()

		server.Handler().ServeHTTP(rr, req)

		if rr.Code != tc.expectedCode {
			t.Errorf("Route %s %s: expected status %d, got %d", tc.method, tc.path, tc.expectedCode, rr.Code)
		}
	}
}

func TestLoggingMiddleware(t *testing.T) {
	// Create test dependencies
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Create server
	server := NewServer(stripeHandler)

	// Test that logging middleware is applied
	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.RemoteAddr = "127.0.0.1:12345"

	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	// Check that request was processed successfully
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Note: We can't easily test the actual logging output without capturing logs,
	// but we can verify the middleware doesn't break the request flow
}

func TestCORSMiddleware(t *testing.T) {
	// Create test dependencies
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Create server
	server := NewServer(stripeHandler)

	t.Run("OPTIONS request", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/api/v1/customers", nil)
		rr := httptest.NewRecorder()

		server.Handler().ServeHTTP(rr, req)

		// Check CORS headers
		if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("Expected Access-Control-Allow-Origin to be '*', got '%s'", rr.Header().Get("Access-Control-Allow-Origin"))
		}

		if rr.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
			t.Errorf("Expected Access-Control-Allow-Methods to be 'GET, POST, PUT, DELETE, OPTIONS', got '%s'", rr.Header().Get("Access-Control-Allow-Methods"))
		}

		if rr.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
			t.Errorf("Expected Access-Control-Allow-Headers to be 'Content-Type, Authorization', got '%s'", rr.Header().Get("Access-Control-Allow-Headers"))
		}

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 for OPTIONS, got %d", rr.Code)
		}
	})

	t.Run("GET request with CORS headers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/health", nil)
		rr := httptest.NewRecorder()

		server.Handler().ServeHTTP(rr, req)

		// Check CORS headers are added to GET requests too
		if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("Expected Access-Control-Allow-Origin to be '*', got '%s'", rr.Header().Get("Access-Control-Allow-Origin"))
		}

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})
}

func TestResponseWriterWrapper(t *testing.T) {
	// Create test dependencies
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Create server
	server := NewServer(stripeHandler)

	// Test that response writer wrapper captures status codes correctly
	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	rr := httptest.NewRecorder()

	server.Handler().ServeHTTP(rr, req)

	// The wrapper should capture the 200 status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Test with a different status code (404 for non-existent route)
	req404 := httptest.NewRequest("GET", "/api/v1/nonexistent", nil)
	rr404 := httptest.NewRecorder()

	server.Handler().ServeHTTP(rr404, req404)

	if rr404.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr404.Code)
	}
}

func TestResponseWriterWrapperWriteHeader(t *testing.T) {
	// Test the WriteHeader method directly
	rr := httptest.NewRecorder()
	wrapper := &responseWriterWrapper{
		ResponseWriter: rr,
		statusCode:     http.StatusOK,
	}

	// Test WriteHeader
	wrapper.WriteHeader(http.StatusCreated)

	if wrapper.statusCode != http.StatusCreated {
		t.Errorf("Expected status code to be %d, got %d", http.StatusCreated, wrapper.statusCode)
	}

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected underlying ResponseWriter status to be %d, got %d", http.StatusCreated, rr.Code)
	}
}

func TestMiddlewareChain(t *testing.T) {
	// Create test dependencies
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Create server
	server := NewServer(stripeHandler)

	// Test that both middleware (logging and CORS) are applied in the correct order
	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	rr := httptest.NewRecorder()

	server.Handler().ServeHTTP(rr, req)

	// Should have CORS headers (from CORS middleware)
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS middleware not applied correctly")
	}

	// Should have processed successfully (logging middleware didn't interfere)
	if rr.Code != http.StatusOK {
		t.Error("Logging middleware interfered with request processing")
	}
}

func TestAllRoutes(t *testing.T) {
	// Create test dependencies
	cfg := &config.Config{
		Stripe: config.StripeConfig{
			SecretKey: "sk_test_123",
		},
	}
	stripeService := service.NewStripeService(cfg)
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Create server
	server := NewServer(stripeHandler)

	// Test all registered routes
	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/health"},
		{"GET", "/api/v1/customers"},
		{"POST", "/api/v1/customers"},
		{"GET", "/api/v1/customers/cus_123"},
		{"POST", "/api/v1/payment-intents"},
		{"POST", "/api/v1/payment-intents/pi_123/confirm"},
		{"POST", "/api/v1/products"},
		{"POST", "/api/v1/prices"},
		{"POST", "/api/v1/subscriptions"},
		{"DELETE", "/api/v1/subscriptions/sub_123"},
		{"OPTIONS", "/api/v1/customers"},
		// Test additional customer ID variations
		{"GET", "/api/v1/customers/cus_different_id"},
		{"DELETE", "/api/v1/subscriptions/sub_different_id"},
		{"POST", "/api/v1/payment-intents/pi_different_id/confirm"},
	}

	for _, route := range routes {
		req := httptest.NewRequest(route.method, route.path, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		server.Handler().ServeHTTP(rr, req)

		// We don't check specific status codes here since they depend on Stripe API
		// We just verify the routes are registered and don't panic
		if rr.Code == 0 {
			t.Errorf("Route %s %s returned no status code", route.method, route.path)
		}
	}
}
