package models

import (
	"time"

	"github.com/google/uuid"
)

// Payment represents a payment transaction
type Payment struct {
	ID                uuid.UUID     `json:"id" db:"id"`
	CustomerID        uuid.UUID     `json:"customer_id" db:"customer_id"`

	// Payment details
	Provider          Provider      `json:"provider" db:"provider"`
	ProviderPaymentID string        `json:"provider_payment_id" db:"provider_payment_id"`
	Amount            int64         `json:"amount" db:"amount"`
	Currency          Currency      `json:"currency" db:"currency"`
	Status            PaymentStatus `json:"status" db:"status"`

	// Payment method
	PaymentMethodType    *string        `json:"payment_method_type,omitempty" db:"payment_method_type"`
	PaymentMethodDetails map[string]any `json:"payment_method_details,omitempty" db:"payment_method_details"`

	// Description
	Description         *string `json:"description,omitempty" db:"description"`
	StatementDescriptor *string `json:"statement_descriptor,omitempty" db:"statement_descriptor"`

	// Related entities
	SubscriptionID *uuid.UUID `json:"subscription_id,omitempty" db:"subscription_id"`
	InvoiceID      *string    `json:"invoice_id,omitempty" db:"invoice_id"`

	// Client secret for frontend confirmation
	ClientSecret *string `json:"client_secret,omitempty" db:"client_secret"`

	// Error handling
	FailureCode    *string `json:"failure_code,omitempty" db:"failure_code"`
	FailureMessage *string `json:"failure_message,omitempty" db:"failure_message"`

	// Metadata
	Metadata map[string]any `json:"metadata,omitempty" db:"metadata"`

	// Idempotency
	IdempotencyKey *string `json:"idempotency_key,omitempty" db:"idempotency_key"`

	// Timestamps
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

// CreatePaymentRequest represents a request to create a payment
type CreatePaymentRequest struct {
	Provider            Provider       `json:"provider"`
	Amount              int64          `json:"amount"`
	Currency            Currency       `json:"currency"`
	Description         string         `json:"description,omitempty"`
	StatementDescriptor string         `json:"statement_descriptor,omitempty"`
	Metadata            map[string]any `json:"metadata,omitempty"`
}

// PaymentListResponse represents a list of payments
type PaymentListResponse struct {
	Data   []Payment `json:"data"`
	Total  int       `json:"total"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}
