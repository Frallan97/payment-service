package models

import (
	"time"

	"github.com/google/uuid"
)

// Customer represents a customer in the payment system
type Customer struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`

	// Provider-specific customer IDs
	StripeCustomerID *string `json:"stripe_customer_id,omitempty" db:"stripe_customer_id"`
	SwishCustomerID  *string `json:"swish_customer_id,omitempty" db:"swish_customer_id"`

	// Metadata
	Metadata map[string]any `json:"metadata,omitempty" db:"metadata"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// CreateCustomerRequest represents a request to create a customer
type CreateCustomerRequest struct {
	UserID   uuid.UUID      `json:"user_id"`
	Email    string         `json:"email"`
	Name     string         `json:"name"`
	Provider Provider       `json:"provider"`
	Metadata map[string]any `json:"metadata,omitempty"`
}
