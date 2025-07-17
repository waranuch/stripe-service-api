# Makefile for Stripe Service API - Clean & Simple

# Variables
BINARY_NAME=stripe-service
DOCKER_IMAGE=stripe-service
PORT=8080

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help build run test clean docker-build docker-run fmt lint

# Default target
help: ## Show this help message
	@echo "$(GREEN)Stripe Service API - Clean & Simple$(NC)"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build the application
build: ## Build the application
	@echo "$(YELLOW)Building $(BINARY_NAME)...$(NC)"
	@go build -o $(BINARY_NAME) .
	@echo "$(GREEN)Build completed successfully!$(NC)"

# Run the application
run: ## Run the application
	@echo "$(YELLOW)Starting $(BINARY_NAME)...$(NC)"
	@go run .

# Run tests
test: ## Run tests
	@echo "$(YELLOW)Running tests...$(NC)"
	@go test ./...
	@echo "$(GREEN)Tests completed!$(NC)"

# Run tests with coverage
test-coverage: ## Run tests with coverage
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	@go test ./... -cover
	@echo "$(GREEN)Coverage analysis completed!$(NC)"

# Run tests with detailed coverage
test-coverage-detailed: ## Run tests with detailed coverage report
	@echo "$(YELLOW)Running tests with detailed coverage...$(NC)"
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# Create test data
test-data: ## Create test data using the API
	@echo "$(YELLOW)Creating test data...$(NC)"
	@go run scripts/create_test_data.go

# Test API endpoints
test-api: ## Test API endpoints with curl
	@echo "$(YELLOW)Testing API endpoints...$(NC)"
	@./scripts/test_api.sh

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "$(GREEN)Cleanup completed!$(NC)"

# Format code
fmt: ## Format code
	@echo "$(YELLOW)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)Code formatted!$(NC)"

# Lint code (requires golangci-lint)
lint: ## Lint code
	@echo "$(YELLOW)Linting code...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		echo "$(GREEN)Linting completed!$(NC)"; \
	else \
		echo "$(RED)golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

# Tidy dependencies
tidy: ## Tidy dependencies
	@echo "$(YELLOW)Tidying dependencies...$(NC)"
	@go mod tidy
	@echo "$(GREEN)Dependencies tidied!$(NC)"

# Build Docker image
docker-build: ## Build Docker image
	@echo "$(YELLOW)Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE) .
	@echo "$(GREEN)Docker image built successfully!$(NC)"

# Run Docker container
docker-run: ## Run Docker container
	@echo "$(YELLOW)Running Docker container...$(NC)"
	@docker run -p $(PORT):$(PORT) -e STRIPE_SECRET_KEY=$(STRIPE_SECRET_KEY) $(DOCKER_IMAGE)

# Development setup
dev-setup: ## Set up development environment
	@echo "$(YELLOW)Setting up development environment...$(NC)"
	@go mod download
	@echo "$(GREEN)Development environment ready!$(NC)"

# Quick start
start: build run ## Build and run the application

# Start with test Stripe key (for development)
start-dev: ## Start with test Stripe key for development
	@echo "$(YELLOW)Starting service with test Stripe key...$(NC)"
	@STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key_here go run main.go

docs: ## Generate and serve API documentation
	@echo "$(CYAN)Generating API documentation...$(NC)"
	@mkdir -p docs
	@echo "$(GREEN)✅ OpenAPI specification: openapi.yaml$(NC)"
	@echo "$(GREEN)✅ HTML documentation: docs/api-documentation.html$(NC)"
	@echo "$(YELLOW)📖 Open docs/api-documentation.html in your browser to view the documentation$(NC)"

docs-serve: ## Serve API documentation locally
	@echo "$(CYAN)Starting documentation server...$(NC)"
	@echo "$(YELLOW)📖 Documentation will be available at: http://localhost:8000/docs/api-documentation.html$(NC)"
	@python3 -m http.server 8000 2>/dev/null || python -m SimpleHTTPServer 8000

docs-validate: ## Validate OpenAPI specification
	@echo "$(CYAN)Validating OpenAPI specification...$(NC)"
	@python3 scripts/validate_openapi.py

# All-in-one development command
dev: fmt test build ## Format, test, and build

# Production build
prod-build: ## Build for production
	@echo "$(YELLOW)Building for production...$(NC)"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o $(BINARY_NAME) .
	@echo "$(GREEN)Production build completed!$(NC)" 