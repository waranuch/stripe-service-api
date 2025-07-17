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
	"stripe-service/internal/server"
	"stripe-service/internal/service"
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

	// Initialize server
	srv := server.NewServer(stripeHandler)

	// Setup HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      srv.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Starting Stripe Service on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
