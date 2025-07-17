package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"stripe-service/internal/models"
	"stripe-service/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// StripeHandler handles HTTP requests for Stripe operations
type StripeHandler struct {
	stripeService service.StripeServiceInterface
	validator     *validator.Validate
}

// NewStripeHandler creates a new Stripe handler
func NewStripeHandler(stripeService service.StripeServiceInterface) *StripeHandler {
	return &StripeHandler{
		stripeService: stripeService,
		validator:     validator.New(),
	}
}

// Helper methods for common operations

// handleServiceError provides consistent error handling for service operations
func (h *StripeHandler) handleServiceError(w http.ResponseWriter, err error, operation string, details map[string]interface{}) {
	// Structured logging with context
	logFields := map[string]interface{}{
		"operation": operation,
		"error":     err.Error(),
	}

	// Add additional details if provided
	for key, value := range details {
		logFields[key] = value
	}

	log.Printf("Service error - Operation: %s, Error: %v, Details: %+v", operation, err, details)
	h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to %s", operation))
}

// parseAndValidateJSON handles JSON parsing and validation
func (h *StripeHandler) parseAndValidateJSON(w http.ResponseWriter, r *http.Request, req interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return false
	}

	if err := h.validator.Struct(req); err != nil {
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
		return false
	}

	return true
}

// extractPathParameter extracts and validates path parameters
func (h *StripeHandler) extractPathParameter(w http.ResponseWriter, r *http.Request, paramName string) (string, bool) {
	vars := mux.Vars(r)
	value := vars[paramName]

	if value == "" {
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("Missing or empty %s parameter", paramName))
		return "", false
	}

	return value, true
}

// HealthCheck handles health check requests
func (h *StripeHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "healthy",
		"service": "stripe-service",
	}

	h.writeJSON(w, http.StatusOK, response)
}

// Customer handlers

// CreateCustomer handles customer creation requests
func (h *StripeHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCustomerRequest

	if !h.parseAndValidateJSON(w, r, &req) {
		return
	}

	customer, err := h.stripeService.CreateCustomer(r.Context(), &req)
	if err != nil {
		h.handleServiceError(w, err, "create customer", map[string]interface{}{
			"email": req.Email,
			"name":  req.Name,
		})
		return
	}

	h.writeJSON(w, http.StatusCreated, customer)
}

// GetCustomer handles customer retrieval requests
func (h *StripeHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	customerID, ok := h.extractPathParameter(w, r, "id")
	if !ok {
		return
	}

	customer, err := h.stripeService.GetCustomer(r.Context(), customerID)
	if err != nil {
		h.handleServiceError(w, err, "get customer", map[string]interface{}{
			"customer_id": customerID,
		})
		return
	}

	h.writeJSON(w, http.StatusOK, customer)
}

// ListCustomers handles customer listing requests
func (h *StripeHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	req := &models.ListCustomersRequest{}

	// Parse optional query parameters
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			req.Limit = limit
		}
	}

	if cursor := r.URL.Query().Get("cursor"); cursor != "" {
		req.Cursor = cursor
	}

	customers, err := h.stripeService.ListCustomers(r.Context(), req)
	if err != nil {
		h.handleServiceError(w, err, "list customers", map[string]interface{}{
			"limit":  req.Limit,
			"cursor": req.Cursor,
		})
		return
	}

	h.writeJSON(w, http.StatusOK, customers)
}

// Payment handlers

// CreatePaymentIntent handles payment intent creation requests
func (h *StripeHandler) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePaymentIntentRequest

	if !h.parseAndValidateJSON(w, r, &req) {
		return
	}

	paymentIntent, err := h.stripeService.CreatePaymentIntent(r.Context(), &req)
	if err != nil {
		h.handleServiceError(w, err, "create payment intent", map[string]interface{}{
			"amount":      req.Amount,
			"currency":    req.Currency,
			"customer_id": req.CustomerID,
		})
		return
	}

	h.writeJSON(w, http.StatusCreated, paymentIntent)
}

// ConfirmPaymentIntent handles payment intent confirmation requests
func (h *StripeHandler) ConfirmPaymentIntent(w http.ResponseWriter, r *http.Request) {
	paymentIntentID, ok := h.extractPathParameter(w, r, "id")
	if !ok {
		return
	}

	var req models.ConfirmPaymentIntentRequest
	if !h.parseAndValidateJSON(w, r, &req) {
		return
	}

	paymentIntent, err := h.stripeService.ConfirmPaymentIntent(r.Context(), paymentIntentID, &req)
	if err != nil {
		h.handleServiceError(w, err, "confirm payment intent", map[string]interface{}{
			"payment_intent_id": paymentIntentID,
			"payment_method_id": req.PaymentMethodID,
		})
		return
	}

	h.writeJSON(w, http.StatusOK, paymentIntent)
}

// Product handlers

// CreateProduct handles product creation requests
func (h *StripeHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductRequest

	if !h.parseAndValidateJSON(w, r, &req) {
		return
	}

	product, err := h.stripeService.CreateProduct(r.Context(), &req)
	if err != nil {
		h.handleServiceError(w, err, "create product", map[string]interface{}{
			"name":   req.Name,
			"active": req.Active,
		})
		return
	}

	h.writeJSON(w, http.StatusCreated, product)
}

// CreatePrice handles price creation requests
func (h *StripeHandler) CreatePrice(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePriceRequest

	if !h.parseAndValidateJSON(w, r, &req) {
		return
	}

	price, err := h.stripeService.CreatePrice(r.Context(), &req)
	if err != nil {
		h.handleServiceError(w, err, "create price", map[string]interface{}{
			"product_id":  req.ProductID,
			"unit_amount": req.UnitAmount,
			"currency":    req.Currency,
		})
		return
	}

	h.writeJSON(w, http.StatusCreated, price)
}

// Subscription handlers

// CreateSubscription handles subscription creation requests
func (h *StripeHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSubscriptionRequest

	if !h.parseAndValidateJSON(w, r, &req) {
		return
	}

	subscription, err := h.stripeService.CreateSubscription(r.Context(), &req)
	if err != nil {
		h.handleServiceError(w, err, "create subscription", map[string]interface{}{
			"customer_id": req.CustomerID,
			"price_id":    req.PriceID,
		})
		return
	}

	h.writeJSON(w, http.StatusCreated, subscription)
}

// CancelSubscription handles subscription cancellation requests
func (h *StripeHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptionID, ok := h.extractPathParameter(w, r, "id")
	if !ok {
		return
	}

	subscription, err := h.stripeService.CancelSubscription(r.Context(), subscriptionID)
	if err != nil {
		h.handleServiceError(w, err, "cancel subscription", map[string]interface{}{
			"subscription_id": subscriptionID,
		})
		return
	}

	h.writeJSON(w, http.StatusOK, subscription)
}

// Helper methods for response handling

func (h *StripeHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func (h *StripeHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}
