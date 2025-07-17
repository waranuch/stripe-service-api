package models

import "time"

// PaymentIntent represents a payment intent
type PaymentIntent struct {
	ID                 string            `json:"id"`
	Amount             int64             `json:"amount"`
	Currency           string            `json:"currency"`
	Status             string            `json:"status"`
	CustomerID         string            `json:"customer_id,omitempty"`
	Description        string            `json:"description,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	ClientSecret       string            `json:"client_secret,omitempty"`
	PaymentMethodID    string            `json:"payment_method_id,omitempty"`
	ConfirmationMethod string            `json:"confirmation_method,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

// CreatePaymentIntentRequest represents the request to create a payment intent
type CreatePaymentIntentRequest struct {
	Amount             int64             `json:"amount" validate:"required,min=1"`
	Currency           string            `json:"currency" validate:"required,len=3"`
	CustomerID         string            `json:"customer_id,omitempty"`
	Description        string            `json:"description,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	PaymentMethodID    string            `json:"payment_method_id,omitempty"`
	ConfirmationMethod string            `json:"confirmation_method,omitempty"`
}

// ConfirmPaymentIntentRequest represents the request to confirm a payment intent
type ConfirmPaymentIntentRequest struct {
	PaymentMethodID string `json:"payment_method_id,omitempty"`
	ReturnURL       string `json:"return_url,omitempty"`
}
