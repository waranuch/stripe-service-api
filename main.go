package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"stripe-service/config"
	"stripe-service/internal/handlers"
	"stripe-service/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate required configuration
	if cfg.Stripe.SecretKey == "" {
		log.Fatal("STRIPE_SECRET_KEY environment variable is required")
	}

	// Initialize services
	stripeService := service.NewStripeService(cfg)

	// Initialize handlers
	stripeHandler := handlers.NewStripeHandler(stripeService)

	// Setup router
	router := setupRouter(stripeHandler)

	// Setup server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Starting Stripe Service on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}

func setupRouter(stripeHandler *handlers.StripeHandler) *mux.Router {
	router := mux.NewRouter()

	// Add middleware
	router.Use(loggingMiddleware)
	router.Use(corsMiddleware)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check
	api.HandleFunc("/health", stripeHandler.HealthCheck).Methods("GET")

	// Customer routes
	api.HandleFunc("/customers", stripeHandler.CreateCustomer).Methods("POST")
	api.HandleFunc("/customers", stripeHandler.ListCustomers).Methods("GET")
	api.HandleFunc("/customers/{id}", stripeHandler.GetCustomer).Methods("GET")

	// Payment intent routes
	api.HandleFunc("/payment-intents", stripeHandler.CreatePaymentIntent).Methods("POST")
	api.HandleFunc("/payment-intents/{id}/confirm", stripeHandler.ConfirmPaymentIntent).Methods("POST")

	// Product routes
	api.HandleFunc("/products", stripeHandler.CreateProduct).Methods("POST")

	// Price routes
	api.HandleFunc("/prices", stripeHandler.CreatePrice).Methods("POST")

	// Subscription routes
	api.HandleFunc("/subscriptions", stripeHandler.CreateSubscription).Methods("POST")
	api.HandleFunc("/subscriptions/{id}", stripeHandler.CancelSubscription).Methods("DELETE")

	return router
}

// loggingMiddleware logs each HTTP request with structured information
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)

		// Structured logging with additional context
		log.Printf("HTTP Request - Method: %s, Path: %s, Status: %d, Duration: %v, UserAgent: %s, RemoteAddr: %s",
			r.Method,
			r.URL.Path,
			wrapper.statusCode,
			duration,
			r.UserAgent(),
			r.RemoteAddr,
		)
	})
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
