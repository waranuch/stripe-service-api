package server

import (
	"log"
	"net/http"
	"time"

	"stripe-service/internal/handlers"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer(stripeHandler *handlers.StripeHandler) *Server {
	s := &Server{}
	s.setupRouter(stripeHandler)
	return s
}

func (s *Server) Handler() http.Handler {
	return s.router
}

func (s *Server) setupRouter(stripeHandler *handlers.StripeHandler) {
	router := mux.NewRouter()

	// Add middleware
	router.Use(s.loggingMiddleware)
	router.Use(s.corsMiddleware)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check
	api.HandleFunc("/health", stripeHandler.HealthCheck).Methods("GET", "OPTIONS")

	// Customer routes
	api.HandleFunc("/customers", stripeHandler.CreateCustomer).Methods("POST")
	api.HandleFunc("/customers", stripeHandler.ListCustomers).Methods("GET")
	api.HandleFunc("/customers/{id}", stripeHandler.GetCustomer).Methods("GET")
	// Add OPTIONS support for all customer routes
	api.HandleFunc("/customers", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")

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

	s.router = router
}

// loggingMiddleware logs each HTTP request with structured information
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
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
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
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
