package models

import "time"

// Customer represents a customer in the system
type Customer struct {
	ID          string            `json:"id"`
	Email       string            `json:"email"`
	Name        string            `json:"name"`
	Phone       string            `json:"phone,omitempty"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// CreateCustomerRequest represents the request to create a customer
type CreateCustomerRequest struct {
	Email       string            `json:"email" validate:"required,email"`
	Name        string            `json:"name" validate:"required"`
	Phone       string            `json:"phone,omitempty"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// UpdateCustomerRequest represents the request to update a customer
type UpdateCustomerRequest struct {
	Email       string            `json:"email,omitempty" validate:"omitempty,email"`
	Name        string            `json:"name,omitempty"`
	Phone       string            `json:"phone,omitempty"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ListCustomersRequest represents the request to list customers
type ListCustomersRequest struct {
	Limit  int64  `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// ListCustomersResponse represents the response when listing customers
type ListCustomersResponse struct {
	Customers  []Customer `json:"customers"`
	HasMore    bool       `json:"has_more"`
	NextCursor string     `json:"next_cursor,omitempty"`
}
