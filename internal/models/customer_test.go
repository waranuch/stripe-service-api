package models

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
)

func TestCreateCustomerRequest_Validation(t *testing.T) {
	validator := validator.New()

	tests := []struct {
		name    string
		request CreateCustomerRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreateCustomerRequest{
				Email: "test@example.com",
				Name:  "John Doe",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			request: CreateCustomerRequest{
				Name: "John Doe",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			request: CreateCustomerRequest{
				Email: "invalid-email",
				Name:  "John Doe",
			},
			wantErr: true,
		},
		{
			name: "missing name",
			request: CreateCustomerRequest{
				Email: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "with optional fields",
			request: CreateCustomerRequest{
				Email:       "test@example.com",
				Name:        "John Doe",
				Phone:       "+1234567890",
				Description: "Test customer",
				Metadata:    map[string]string{"key": "value"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Struct(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCustomerRequest validation = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCustomer_Structure(t *testing.T) {
	now := time.Now()
	customer := Customer{
		ID:          "cus_123456789",
		Email:       "test@example.com",
		Name:        "John Doe",
		Phone:       "+1234567890",
		Description: "Test customer",
		CreatedAt:   now,
		UpdatedAt:   now,
		Metadata:    map[string]string{"key": "value"},
	}

	// Test that all fields are properly set
	if customer.ID != "cus_123456789" {
		t.Errorf("Expected ID to be 'cus_123456789', got %s", customer.ID)
	}
	if customer.Email != "test@example.com" {
		t.Errorf("Expected Email to be 'test@example.com', got %s", customer.Email)
	}
	if customer.Name != "John Doe" {
		t.Errorf("Expected Name to be 'John Doe', got %s", customer.Name)
	}
	if customer.Phone != "+1234567890" {
		t.Errorf("Expected Phone to be '+1234567890', got %s", customer.Phone)
	}
	if customer.Description != "Test customer" {
		t.Errorf("Expected Description to be 'Test customer', got %s", customer.Description)
	}
	if customer.CreatedAt != now {
		t.Errorf("Expected CreatedAt to be %v, got %v", now, customer.CreatedAt)
	}
	if customer.UpdatedAt != now {
		t.Errorf("Expected UpdatedAt to be %v, got %v", now, customer.UpdatedAt)
	}
	if customer.Metadata["key"] != "value" {
		t.Errorf("Expected Metadata['key'] to be 'value', got %s", customer.Metadata["key"])
	}
}

func TestListCustomersResponse_Structure(t *testing.T) {
	customers := []Customer{
		{
			ID:    "cus_1",
			Email: "test1@example.com",
			Name:  "John Doe",
		},
		{
			ID:    "cus_2",
			Email: "test2@example.com",
			Name:  "Jane Smith",
		},
	}

	response := ListCustomersResponse{
		Customers:  customers,
		HasMore:    true,
		NextCursor: "next_cursor_value",
	}

	if len(response.Customers) != 2 {
		t.Errorf("Expected 2 customers, got %d", len(response.Customers))
	}
	if !response.HasMore {
		t.Error("Expected HasMore to be true")
	}
	if response.NextCursor != "next_cursor_value" {
		t.Errorf("Expected NextCursor to be 'next_cursor_value', got %s", response.NextCursor)
	}
	if response.Customers[0].ID != "cus_1" {
		t.Errorf("Expected first customer ID to be 'cus_1', got %s", response.Customers[0].ID)
	}
	if response.Customers[1].ID != "cus_2" {
		t.Errorf("Expected second customer ID to be 'cus_2', got %s", response.Customers[1].ID)
	}
}
