package models

import "time"

// Product represents a product
type Product struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Active      bool              `json:"active"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description,omitempty"`
	Active      bool              `json:"active"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Price represents a price for a product
type Price struct {
	ID                string            `json:"id"`
	ProductID         string            `json:"product_id"`
	UnitAmount        int64             `json:"unit_amount"`
	Currency          string            `json:"currency"`
	Type              string            `json:"type"`
	RecurringInterval string            `json:"recurring_interval,omitempty"`
	Active            bool              `json:"active"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// CreatePriceRequest represents the request to create a price
type CreatePriceRequest struct {
	ProductID         string            `json:"product_id" validate:"required"`
	UnitAmount        int64             `json:"unit_amount" validate:"required,min=1"`
	Currency          string            `json:"currency" validate:"required,len=3"`
	Type              string            `json:"type" validate:"required,oneof=one_time recurring"`
	RecurringInterval string            `json:"recurring_interval,omitempty"`
	Active            bool              `json:"active"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// Subscription represents a subscription
type Subscription struct {
	ID                 string            `json:"id"`
	CustomerID         string            `json:"customer_id"`
	PriceID            string            `json:"price_id"`
	Status             string            `json:"status"`
	CurrentPeriodStart time.Time         `json:"current_period_start"`
	CurrentPeriodEnd   time.Time         `json:"current_period_end"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

// CreateSubscriptionRequest represents the request to create a subscription
type CreateSubscriptionRequest struct {
	CustomerID string            `json:"customer_id" validate:"required"`
	PriceID    string            `json:"price_id" validate:"required"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}
