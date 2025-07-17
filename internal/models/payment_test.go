package models

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
)

func TestCreatePaymentIntentRequest_Validation(t *testing.T) {
	validator := validator.New()

	tests := []struct {
		name    string
		request CreatePaymentIntentRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreatePaymentIntentRequest{
				Amount:   1000,
				Currency: "usd",
			},
			wantErr: false,
		},
		{
			name: "missing amount",
			request: CreatePaymentIntentRequest{
				Currency: "usd",
			},
			wantErr: true,
		},
		{
			name: "zero amount",
			request: CreatePaymentIntentRequest{
				Amount:   0,
				Currency: "usd",
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			request: CreatePaymentIntentRequest{
				Amount:   -100,
				Currency: "usd",
			},
			wantErr: true,
		},
		{
			name: "missing currency",
			request: CreatePaymentIntentRequest{
				Amount: 1000,
			},
			wantErr: true,
		},
		{
			name: "with optional fields",
			request: CreatePaymentIntentRequest{
				Amount:      1000,
				Currency:    "usd",
				CustomerID:  "cus_123",
				Description: "Test payment",
				Metadata:    map[string]string{"order_id": "123"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Struct(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePaymentIntentRequest validation = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfirmPaymentIntentRequest_Validation(t *testing.T) {
	validator := validator.New()

	tests := []struct {
		name    string
		request ConfirmPaymentIntentRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: ConfirmPaymentIntentRequest{
				PaymentMethodID: "pm_card_visa",
			},
			wantErr: false,
		},
		{
			name:    "empty request",
			request: ConfirmPaymentIntentRequest{},
			wantErr: false, // No required fields for confirm request
		},
		{
			name: "with return url",
			request: ConfirmPaymentIntentRequest{
				PaymentMethodID: "pm_card_visa",
				ReturnURL:       "https://example.com/return",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Struct(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfirmPaymentIntentRequest validation = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPaymentIntent_Structure(t *testing.T) {
	now := time.Now()
	payment := PaymentIntent{
		ID:                 "pi_123456789",
		Amount:             1000,
		Currency:           "usd",
		Status:             "requires_payment_method",
		ClientSecret:       "pi_123_secret_456",
		CustomerID:         "cus_123",
		Description:        "Test payment",
		PaymentMethodID:    "pm_card_visa",
		ConfirmationMethod: "automatic",
		Metadata:           map[string]string{"order_id": "123"},
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	// Test that all fields are properly set
	if payment.ID != "pi_123456789" {
		t.Errorf("Expected ID to be 'pi_123456789', got %s", payment.ID)
	}
	if payment.Amount != 1000 {
		t.Errorf("Expected Amount to be 1000, got %d", payment.Amount)
	}
	if payment.Currency != "usd" {
		t.Errorf("Expected Currency to be 'usd', got %s", payment.Currency)
	}
	if payment.Status != "requires_payment_method" {
		t.Errorf("Expected Status to be 'requires_payment_method', got %s", payment.Status)
	}
	if payment.ClientSecret != "pi_123_secret_456" {
		t.Errorf("Expected ClientSecret to be 'pi_123_secret_456', got %s", payment.ClientSecret)
	}
	if payment.CustomerID != "cus_123" {
		t.Errorf("Expected CustomerID to be 'cus_123', got %s", payment.CustomerID)
	}
	if payment.Description != "Test payment" {
		t.Errorf("Expected Description to be 'Test payment', got %s", payment.Description)
	}
	if payment.PaymentMethodID != "pm_card_visa" {
		t.Errorf("Expected PaymentMethodID to be 'pm_card_visa', got %s", payment.PaymentMethodID)
	}
	if payment.ConfirmationMethod != "automatic" {
		t.Errorf("Expected ConfirmationMethod to be 'automatic', got %s", payment.ConfirmationMethod)
	}
	if payment.Metadata["order_id"] != "123" {
		t.Errorf("Expected Metadata['order_id'] to be '123', got %s", payment.Metadata["order_id"])
	}
	if payment.CreatedAt != now {
		t.Errorf("Expected CreatedAt to be %v, got %v", now, payment.CreatedAt)
	}
	if payment.UpdatedAt != now {
		t.Errorf("Expected UpdatedAt to be %v, got %v", now, payment.UpdatedAt)
	}
}
