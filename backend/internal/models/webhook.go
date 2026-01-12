package models

import (
	"time"

	"github.com/google/uuid"
)

// WebhookEvent represents a webhook event from a payment provider
type WebhookEvent struct {
	ID              uuid.UUID `json:"id" db:"id"`

	// Webhook identification
	Provider        Provider  `json:"provider" db:"provider"`
	ProviderEventID string    `json:"provider_event_id" db:"provider_event_id"`
	EventType       string    `json:"event_type" db:"event_type"`

	// Processing status
	Processed            bool    `json:"processed" db:"processed"`
	ProcessingAttempts   int     `json:"processing_attempts" db:"processing_attempts"`
	LastProcessingError  *string `json:"last_processing_error,omitempty" db:"last_processing_error"`

	// Payload
	Payload map[string]any `json:"payload" db:"payload"`

	// Related entity
	PaymentID      *uuid.UUID `json:"payment_id,omitempty" db:"payment_id"`
	SubscriptionID *uuid.UUID `json:"subscription_id,omitempty" db:"subscription_id"`
	RefundID       *uuid.UUID `json:"refund_id,omitempty" db:"refund_id"`

	// Timestamps
	ReceivedAt  time.Time  `json:"received_at" db:"received_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty" db:"processed_at"`
}
