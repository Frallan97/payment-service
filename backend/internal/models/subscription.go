package models

import (
	"time"

	"github.com/google/uuid"
)

// Subscription represents a recurring payment subscription
type Subscription struct {
	ID         uuid.UUID          `json:"id" db:"id"`
	CustomerID uuid.UUID          `json:"customer_id" db:"customer_id"`

	// Subscription details
	Provider               Provider           `json:"provider" db:"provider"`
	ProviderSubscriptionID string             `json:"provider_subscription_id" db:"provider_subscription_id"`
	Status                 SubscriptionStatus `json:"status" db:"status"`

	// Billing
	Amount        int64    `json:"amount" db:"amount"`
	Currency      Currency `json:"currency" db:"currency"`
	Interval      string   `json:"interval" db:"interval"`
	IntervalCount int      `json:"interval_count" db:"interval_count"`

	// Billing dates
	CurrentPeriodStart time.Time  `json:"current_period_start" db:"current_period_start"`
	CurrentPeriodEnd   time.Time  `json:"current_period_end" db:"current_period_end"`
	TrialStart         *time.Time `json:"trial_start,omitempty" db:"trial_start"`
	TrialEnd           *time.Time `json:"trial_end,omitempty" db:"trial_end"`

	// Cancellation
	CancelAt           *time.Time `json:"cancel_at,omitempty" db:"cancel_at"`
	CanceledAt         *time.Time `json:"canceled_at,omitempty" db:"canceled_at"`
	CancelAtPeriodEnd  bool       `json:"cancel_at_period_end" db:"cancel_at_period_end"`

	// Latest payment
	LatestPaymentID *uuid.UUID `json:"latest_payment_id,omitempty" db:"latest_payment_id"`

	// Product info
	ProductName        string  `json:"product_name" db:"product_name"`
	ProductDescription *string `json:"product_description,omitempty" db:"product_description"`

	// Metadata
	Metadata map[string]any `json:"metadata,omitempty" db:"metadata"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateSubscriptionRequest represents a request to create a subscription
type CreateSubscriptionRequest struct {
	Provider           Provider       `json:"provider"`
	Amount             int64          `json:"amount"`
	Currency           Currency       `json:"currency"`
	Interval           string         `json:"interval"`
	IntervalCount      int            `json:"interval_count"`
	ProductName        string         `json:"product_name"`
	ProductDescription string         `json:"product_description,omitempty"`
	TrialPeriodDays    int            `json:"trial_period_days,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
}

// UpdateSubscriptionRequest represents a request to update a subscription
type UpdateSubscriptionRequest struct {
	CancelAtPeriodEnd *bool          `json:"cancel_at_period_end,omitempty"`
	Metadata          map[string]any `json:"metadata,omitempty"`
}

// SubscriptionListResponse represents a list of subscriptions
type SubscriptionListResponse struct {
	Data   []Subscription `json:"data"`
	Total  int            `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}
