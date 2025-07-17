package service

import (
	"context"
	"stripe-service/internal/models"
)

// StripeServiceInterface defines the interface for Stripe operations
type StripeServiceInterface interface {
	CreateCustomer(ctx context.Context, req *models.CreateCustomerRequest) (*models.Customer, error)
	GetCustomer(ctx context.Context, customerID string) (*models.Customer, error)
	ListCustomers(ctx context.Context, req *models.ListCustomersRequest) (*models.ListCustomersResponse, error)
	CreatePaymentIntent(ctx context.Context, req *models.CreatePaymentIntentRequest) (*models.PaymentIntent, error)
	ConfirmPaymentIntent(ctx context.Context, paymentIntentID string, req *models.ConfirmPaymentIntentRequest) (*models.PaymentIntent, error)
	CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error)
	CreatePrice(ctx context.Context, req *models.CreatePriceRequest) (*models.Price, error)
	CreateSubscription(ctx context.Context, req *models.CreateSubscriptionRequest) (*models.Subscription, error)
	CancelSubscription(ctx context.Context, subscriptionID string) (*models.Subscription, error)
}
