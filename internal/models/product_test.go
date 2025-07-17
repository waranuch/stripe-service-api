package models

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
)

func TestCreateProductRequest_Validation(t *testing.T) {
	validator := validator.New()

	tests := []struct {
		name    string
		request CreateProductRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreateProductRequest{
				Name: "Test Product",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			request: CreateProductRequest{
				Description: "Test description",
			},
			wantErr: true,
		},
		{
			name: "with optional fields",
			request: CreateProductRequest{
				Name:        "Test Product",
				Description: "Test description",
				Active:      true,
				Metadata:    map[string]string{"category": "test"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Struct(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProductRequest validation = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreatePriceRequest_Validation(t *testing.T) {
	validator := validator.New()

	tests := []struct {
		name    string
		request CreatePriceRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreatePriceRequest{
				ProductID:  "prod_123",
				UnitAmount: 1000,
				Currency:   "usd",
				Type:       "one_time",
			},
			wantErr: false,
		},
		{
			name: "missing product id",
			request: CreatePriceRequest{
				UnitAmount: 1000,
				Currency:   "usd",
			},
			wantErr: true,
		},
		{
			name: "missing unit amount",
			request: CreatePriceRequest{
				ProductID: "prod_123",
				Currency:  "usd",
			},
			wantErr: true,
		},
		{
			name: "zero unit amount",
			request: CreatePriceRequest{
				ProductID:  "prod_123",
				UnitAmount: 0,
				Currency:   "usd",
			},
			wantErr: true,
		},
		{
			name: "negative unit amount",
			request: CreatePriceRequest{
				ProductID:  "prod_123",
				UnitAmount: -100,
				Currency:   "usd",
			},
			wantErr: true,
		},
		{
			name: "missing currency",
			request: CreatePriceRequest{
				ProductID:  "prod_123",
				UnitAmount: 1000,
			},
			wantErr: true,
		},
		{
			name: "with recurring",
			request: CreatePriceRequest{
				ProductID:         "prod_123",
				UnitAmount:        1000,
				Currency:          "usd",
				Type:              "recurring",
				RecurringInterval: "month",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Struct(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePriceRequest validation = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateSubscriptionRequest_Validation(t *testing.T) {
	validator := validator.New()

	tests := []struct {
		name    string
		request CreateSubscriptionRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreateSubscriptionRequest{
				CustomerID: "cus_123",
				PriceID:    "price_123",
			},
			wantErr: false,
		},
		{
			name: "missing customer id",
			request: CreateSubscriptionRequest{
				PriceID: "price_123",
			},
			wantErr: true,
		},
		{
			name: "missing price id",
			request: CreateSubscriptionRequest{
				CustomerID: "cus_123",
			},
			wantErr: true,
		},
		{
			name: "with optional fields",
			request: CreateSubscriptionRequest{
				CustomerID: "cus_123",
				PriceID:    "price_123",
				Metadata:   map[string]string{"plan": "premium"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Struct(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSubscriptionRequest validation = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProduct_Structure(t *testing.T) {
	now := time.Now()
	product := Product{
		ID:          "prod_123456789",
		Name:        "Test Product",
		Description: "Test description",
		Active:      true,
		Metadata:    map[string]string{"category": "test"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if product.ID != "prod_123456789" {
		t.Errorf("Expected ID to be 'prod_123456789', got %s", product.ID)
	}
	if product.Name != "Test Product" {
		t.Errorf("Expected Name to be 'Test Product', got %s", product.Name)
	}
	if product.Description != "Test description" {
		t.Errorf("Expected Description to be 'Test description', got %s", product.Description)
	}
	if !product.Active {
		t.Error("Expected Active to be true")
	}
	if product.Metadata["category"] != "test" {
		t.Errorf("Expected Metadata['category'] to be 'test', got %s", product.Metadata["category"])
	}
	if product.CreatedAt != now {
		t.Errorf("Expected CreatedAt to be %v, got %v", now, product.CreatedAt)
	}
	if product.UpdatedAt != now {
		t.Errorf("Expected UpdatedAt to be %v, got %v", now, product.UpdatedAt)
	}
}

func TestPrice_Structure(t *testing.T) {
	now := time.Now()
	price := Price{
		ID:                "price_123456789",
		ProductID:         "prod_123",
		UnitAmount:        1000,
		Currency:          "usd",
		Type:              "recurring",
		RecurringInterval: "month",
		Active:            true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if price.ID != "price_123456789" {
		t.Errorf("Expected ID to be 'price_123456789', got %s", price.ID)
	}
	if price.ProductID != "prod_123" {
		t.Errorf("Expected ProductID to be 'prod_123', got %s", price.ProductID)
	}
	if price.UnitAmount != 1000 {
		t.Errorf("Expected UnitAmount to be 1000, got %d", price.UnitAmount)
	}
	if price.Currency != "usd" {
		t.Errorf("Expected Currency to be 'usd', got %s", price.Currency)
	}
	if price.Type != "recurring" {
		t.Errorf("Expected Type to be 'recurring', got %s", price.Type)
	}
	if price.RecurringInterval != "month" {
		t.Errorf("Expected RecurringInterval to be 'month', got %s", price.RecurringInterval)
	}
	if !price.Active {
		t.Error("Expected Active to be true")
	}
	if price.CreatedAt != now {
		t.Errorf("Expected CreatedAt to be %v, got %v", now, price.CreatedAt)
	}
	if price.UpdatedAt != now {
		t.Errorf("Expected UpdatedAt to be %v, got %v", now, price.UpdatedAt)
	}
}

func TestSubscription_Structure(t *testing.T) {
	now := time.Now()
	subscription := Subscription{
		ID:         "sub_123456789",
		CustomerID: "cus_123",
		PriceID:    "price_123",
		Status:     "active",
		Metadata:   map[string]string{"plan": "premium"},
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if subscription.ID != "sub_123456789" {
		t.Errorf("Expected ID to be 'sub_123456789', got %s", subscription.ID)
	}
	if subscription.CustomerID != "cus_123" {
		t.Errorf("Expected CustomerID to be 'cus_123', got %s", subscription.CustomerID)
	}
	if subscription.PriceID != "price_123" {
		t.Errorf("Expected PriceID to be 'price_123', got %s", subscription.PriceID)
	}
	if subscription.Status != "active" {
		t.Errorf("Expected Status to be 'active', got %s", subscription.Status)
	}
	if subscription.Metadata["plan"] != "premium" {
		t.Errorf("Expected Metadata['plan'] to be 'premium', got %s", subscription.Metadata["plan"])
	}
	if subscription.CreatedAt != now {
		t.Errorf("Expected CreatedAt to be %v, got %v", now, subscription.CreatedAt)
	}
	if subscription.UpdatedAt != now {
		t.Errorf("Expected UpdatedAt to be %v, got %v", now, subscription.UpdatedAt)
	}
}
