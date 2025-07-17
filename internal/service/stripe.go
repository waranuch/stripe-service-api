package service

import (
	"context"
	"fmt"
	"time"

	"stripe-service/config"
	"stripe-service/internal/models"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
)

// Constants for default values
const (
	DefaultCustomerLimit = 10
	MaxCustomerLimit     = 100
)

// StripeService handles all Stripe operations
type StripeService struct {
	config *config.Config
	client *client.API
}

// NewStripeService creates a new Stripe service with its own client instance
func NewStripeService(cfg *config.Config) *StripeService {
	// Create a new Stripe client instance instead of using global state
	stripeClient := &client.API{}
	stripeClient.Init(cfg.Stripe.SecretKey, nil)

	return &StripeService{
		config: cfg,
		client: stripeClient,
	}
}

// Customer operations

// CreateCustomer creates a new customer in Stripe
func (s *StripeService) CreateCustomer(ctx context.Context, req *models.CreateCustomerRequest) (*models.Customer, error) {
	params := &stripe.CustomerParams{
		Email:       stripe.String(req.Email),
		Name:        stripe.String(req.Name),
		Description: stripe.String(req.Description),
	}

	// Set context for cancellation support
	params.Context = ctx

	if req.Phone != "" {
		params.Phone = stripe.String(req.Phone)
	}

	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	stripeCustomer, err := s.client.Customers.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return s.convertStripeCustomer(stripeCustomer), nil
}

// GetCustomer retrieves a customer by ID
func (s *StripeService) GetCustomer(ctx context.Context, customerID string) (*models.Customer, error) {
	params := &stripe.CustomerParams{}
	params.Context = ctx

	stripeCustomer, err := s.client.Customers.Get(customerID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return s.convertStripeCustomer(stripeCustomer), nil
}

// ListCustomers lists customers with pagination
func (s *StripeService) ListCustomers(ctx context.Context, req *models.ListCustomersRequest) (*models.ListCustomersResponse, error) {
	params := &stripe.CustomerListParams{}
	params.Context = ctx

	if req.Limit > 0 {
		params.Limit = stripe.Int64(req.Limit)
	} else {
		params.Limit = stripe.Int64(DefaultCustomerLimit)
	}

	if req.Cursor != "" {
		params.StartingAfter = stripe.String(req.Cursor)
	}

	iter := s.client.Customers.List(params)
	var customers []models.Customer

	for iter.Next() {
		customers = append(customers, *s.convertStripeCustomer(iter.Customer()))
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}

	return &models.ListCustomersResponse{
		Customers: customers,
		HasMore:   iter.Meta().HasMore,
	}, nil
}

// Payment operations

// CreatePaymentIntent creates a new payment intent
func (s *StripeService) CreatePaymentIntent(ctx context.Context, req *models.CreatePaymentIntentRequest) (*models.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Amount),
		Currency: stripe.String(req.Currency),
	}
	params.Context = ctx

	if req.CustomerID != "" {
		params.Customer = stripe.String(req.CustomerID)
	}

	if req.Description != "" {
		params.Description = stripe.String(req.Description)
	}

	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	if req.PaymentMethodID != "" {
		params.PaymentMethod = stripe.String(req.PaymentMethodID)
	}

	if req.ConfirmationMethod != "" {
		params.ConfirmationMethod = stripe.String(req.ConfirmationMethod)
	}

	stripePI, err := s.client.PaymentIntents.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return s.convertStripePaymentIntent(stripePI), nil
}

// ConfirmPaymentIntent confirms a payment intent
func (s *StripeService) ConfirmPaymentIntent(ctx context.Context, paymentIntentID string, req *models.ConfirmPaymentIntentRequest) (*models.PaymentIntent, error) {
	params := &stripe.PaymentIntentConfirmParams{}
	params.Context = ctx

	if req.PaymentMethodID != "" {
		params.PaymentMethod = stripe.String(req.PaymentMethodID)
	}

	if req.ReturnURL != "" {
		params.ReturnURL = stripe.String(req.ReturnURL)
	}

	stripePI, err := s.client.PaymentIntents.Confirm(paymentIntentID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	return s.convertStripePaymentIntent(stripePI), nil
}

// Product operations

// CreateProduct creates a new product
func (s *StripeService) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	params := &stripe.ProductParams{
		Name:        stripe.String(req.Name),
		Description: stripe.String(req.Description),
		Active:      stripe.Bool(req.Active),
	}
	params.Context = ctx

	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	stripeProduct, err := s.client.Products.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return s.convertStripeProduct(stripeProduct), nil
}

// CreatePrice creates a new price
func (s *StripeService) CreatePrice(ctx context.Context, req *models.CreatePriceRequest) (*models.Price, error) {
	params := &stripe.PriceParams{
		Product:    stripe.String(req.ProductID),
		UnitAmount: stripe.Int64(req.UnitAmount),
		Currency:   stripe.String(req.Currency),
		Active:     stripe.Bool(req.Active),
	}
	params.Context = ctx

	if req.Type == "recurring" && req.RecurringInterval != "" {
		params.Recurring = &stripe.PriceRecurringParams{
			Interval: stripe.String(req.RecurringInterval),
		}
	}

	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	stripePrice, err := s.client.Prices.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create price: %w", err)
	}

	return s.convertStripePrice(stripePrice), nil
}

// Subscription operations

// CreateSubscription creates a new subscription
func (s *StripeService) CreateSubscription(ctx context.Context, req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(req.CustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(req.PriceID),
			},
		},
	}
	params.Context = ctx

	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	stripeSub, err := s.client.Subscriptions.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return s.convertStripeSubscription(stripeSub), nil
}

// CancelSubscription cancels a subscription
func (s *StripeService) CancelSubscription(ctx context.Context, subscriptionID string) (*models.Subscription, error) {
	params := &stripe.SubscriptionCancelParams{}
	params.Context = ctx

	stripeSub, err := s.client.Subscriptions.Cancel(subscriptionID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel subscription: %w", err)
	}

	return s.convertStripeSubscription(stripeSub), nil
}

// Helper methods to convert Stripe objects to internal models

// StripeCustomer interface for testing
type StripeCustomer interface {
	GetID() string
	GetEmail() string
	GetName() string
	GetPhone() string
	GetDescription() string
	GetMetadata() map[string]string
	GetCreated() int64
}

// Adapter for real Stripe customer
type stripeCustomerAdapter struct {
	customer *stripe.Customer
}

func (a *stripeCustomerAdapter) GetID() string {
	if a.customer == nil {
		return ""
	}
	return a.customer.ID
}

func (a *stripeCustomerAdapter) GetEmail() string {
	if a.customer == nil {
		return ""
	}
	return a.customer.Email
}

func (a *stripeCustomerAdapter) GetName() string {
	if a.customer == nil {
		return ""
	}
	return a.customer.Name
}

func (a *stripeCustomerAdapter) GetPhone() string {
	if a.customer == nil {
		return ""
	}
	return a.customer.Phone
}

func (a *stripeCustomerAdapter) GetDescription() string {
	if a.customer == nil {
		return ""
	}
	return a.customer.Description
}

func (a *stripeCustomerAdapter) GetMetadata() map[string]string {
	if a.customer == nil {
		return nil
	}
	return a.customer.Metadata
}

func (a *stripeCustomerAdapter) GetCreated() int64 {
	if a.customer == nil {
		return 0
	}
	return a.customer.Created
}

func (s *StripeService) convertStripeCustomer(stripeCustomer *stripe.Customer) *models.Customer {
	if stripeCustomer == nil {
		return nil
	}
	adapter := &stripeCustomerAdapter{customer: stripeCustomer}
	return s.convertStripeCustomerInterface(adapter)
}

func (s *StripeService) convertStripeCustomerInterface(stripeCustomer StripeCustomer) *models.Customer {
	// Use Stripe's created timestamp for both created and updated
	createdAt := time.Unix(stripeCustomer.GetCreated(), 0)

	return &models.Customer{
		ID:          stripeCustomer.GetID(),
		Email:       stripeCustomer.GetEmail(),
		Name:        stripeCustomer.GetName(),
		Phone:       stripeCustomer.GetPhone(),
		Description: stripeCustomer.GetDescription(),
		Metadata:    stripeCustomer.GetMetadata(),
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt, // Stripe doesn't provide separate updated_at
	}
}

func (s *StripeService) convertStripePaymentIntent(stripePI *stripe.PaymentIntent) *models.PaymentIntent {
	if stripePI == nil {
		return nil
	}
	customerID := ""
	if stripePI.Customer != nil {
		customerID = stripePI.Customer.ID
	}

	createdAt := time.Unix(stripePI.Created, 0)

	return &models.PaymentIntent{
		ID:           stripePI.ID,
		Amount:       stripePI.Amount,
		Currency:     string(stripePI.Currency),
		Status:       string(stripePI.Status),
		CustomerID:   customerID,
		Description:  stripePI.Description,
		Metadata:     stripePI.Metadata,
		ClientSecret: stripePI.ClientSecret,
		CreatedAt:    createdAt,
		UpdatedAt:    createdAt,
	}
}

func (s *StripeService) convertStripeProduct(stripeProduct *stripe.Product) *models.Product {
	if stripeProduct == nil {
		return nil
	}
	createdAt := time.Unix(stripeProduct.Created, 0)
	updatedAt := createdAt
	if stripeProduct.Updated > 0 {
		updatedAt = time.Unix(stripeProduct.Updated, 0)
	}

	return &models.Product{
		ID:          stripeProduct.ID,
		Name:        stripeProduct.Name,
		Description: stripeProduct.Description,
		Active:      stripeProduct.Active,
		Metadata:    stripeProduct.Metadata,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (s *StripeService) convertStripePrice(stripePrice *stripe.Price) *models.Price {
	if stripePrice == nil {
		return nil
	}
	createdAt := time.Unix(stripePrice.Created, 0)

	priceType := "one_time"
	recurringInterval := ""

	if stripePrice.Recurring != nil {
		priceType = "recurring"
		recurringInterval = string(stripePrice.Recurring.Interval)
	}

	return &models.Price{
		ID:                stripePrice.ID,
		ProductID:         stripePrice.Product.ID,
		UnitAmount:        stripePrice.UnitAmount,
		Currency:          string(stripePrice.Currency),
		Type:              priceType,
		RecurringInterval: recurringInterval,
		Active:            stripePrice.Active,
		Metadata:          stripePrice.Metadata,
		CreatedAt:         createdAt,
		UpdatedAt:         createdAt,
	}
}

func (s *StripeService) convertStripeSubscription(stripeSub *stripe.Subscription) *models.Subscription {
	if stripeSub == nil {
		return nil
	}
	createdAt := time.Unix(stripeSub.Created, 0)

	return &models.Subscription{
		ID:                 stripeSub.ID,
		CustomerID:         stripeSub.Customer.ID,
		PriceID:            stripeSub.Items.Data[0].Price.ID,
		Status:             string(stripeSub.Status),
		CurrentPeriodStart: time.Unix(stripeSub.CurrentPeriodStart, 0),
		CurrentPeriodEnd:   time.Unix(stripeSub.CurrentPeriodEnd, 0),
		Metadata:           stripeSub.Metadata,
		CreatedAt:          createdAt,
		UpdatedAt:          createdAt,
	}
}
