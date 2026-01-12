package models

import (
	"time"

	"github.com/google/uuid"
)

// Refund represents a payment refund
type Refund struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	PaymentID uuid.UUID    `json:"payment_id" db:"payment_id"`

	// Refund details
	Provider        Provider     `json:"provider" db:"provider"`
	ProviderRefundID string      `json:"provider_refund_id" db:"provider_refund_id"`
	Amount          int64        `json:"amount" db:"amount"`
	Currency        Currency     `json:"currency" db:"currency"`
	Status          RefundStatus `json:"status" db:"status"`

	// Reason
	Reason *string `json:"reason,omitempty" db:"reason"`
	Notes  *string `json:"notes,omitempty" db:"notes"`

	// Error handling
	FailureCode    *string `json:"failure_code,omitempty" db:"failure_code"`
	FailureMessage *string `json:"failure_message,omitempty" db:"failure_message"`

	// Metadata
	Metadata map[string]any `json:"metadata,omitempty" db:"metadata"`

	// Timestamps
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

// CreateRefundRequest represents a request to create a refund
type CreateRefundRequest struct {
	PaymentID uuid.UUID      `json:"payment_id"`
	Amount    int64          `json:"amount"`
	Reason    string         `json:"reason,omitempty"`
	Notes     string         `json:"notes,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// RefundListResponse represents a list of refunds
type RefundListResponse struct {
	Data   []Refund `json:"data"`
	Total  int      `json:"total"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
}
