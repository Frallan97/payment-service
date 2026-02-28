package client

import (
	"time"

	"github.com/google/uuid"
)

// Provider represents a payment provider.
type Provider string

const (
	ProviderStripe Provider = "stripe"
	ProviderSwish  Provider = "swish"
)

// Currency represents a supported currency.
type Currency string

const (
	CurrencySEK Currency = "SEK"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
)

// --- Payment types ---

// PaymentStatus represents the status of a payment.
type PaymentStatus string

const (
	StatusPending        PaymentStatus = "pending"
	StatusProcessing     PaymentStatus = "processing"
	StatusRequiresAction PaymentStatus = "requires_action"
	StatusSucceeded      PaymentStatus = "succeeded"
	StatusFailed         PaymentStatus = "failed"
	StatusCanceled       PaymentStatus = "canceled"
)

// Payment represents a payment returned by the API.
type Payment struct {
	ID                   uuid.UUID      `json:"id"`
	CustomerID           uuid.UUID      `json:"customer_id"`
	Provider             Provider       `json:"provider"`
	ProviderPaymentID    string         `json:"provider_payment_id"`
	Amount               int64          `json:"amount"`
	Currency             Currency       `json:"currency"`
	Status               PaymentStatus  `json:"status"`
	PaymentMethodType    *string        `json:"payment_method_type,omitempty"`
	PaymentMethodDetails map[string]any `json:"payment_method_details,omitempty"`
	Description          *string        `json:"description,omitempty"`
	ClientSecret         *string        `json:"client_secret,omitempty"`
	FailureCode          *string        `json:"failure_code,omitempty"`
	FailureMessage       *string        `json:"failure_message,omitempty"`
	Metadata             map[string]any `json:"metadata,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	CompletedAt          *time.Time     `json:"completed_at,omitempty"`
}

// CreatePaymentRequest is the request body for creating a payment.
type CreatePaymentRequest struct {
	Provider    Provider       `json:"provider"`
	Amount      int64          `json:"amount"`
	Currency    Currency       `json:"currency"`
	Description string         `json:"description,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// PaymentListResponse is the response for listing payments.
type PaymentListResponse struct {
	Data   []Payment `json:"data"`
	Total  int       `json:"total"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}

// --- Subscription types ---

// SubscriptionStatus represents the status of a subscription.
type SubscriptionStatus string

const (
	SubscriptionStatusActive            SubscriptionStatus = "active"
	SubscriptionStatusPastDue           SubscriptionStatus = "past_due"
	SubscriptionStatusUnpaid            SubscriptionStatus = "unpaid"
	SubscriptionStatusCanceled          SubscriptionStatus = "canceled"
	SubscriptionStatusIncomplete        SubscriptionStatus = "incomplete"
	SubscriptionStatusIncompleteExpired SubscriptionStatus = "incomplete_expired"
	SubscriptionStatusTrialing          SubscriptionStatus = "trialing"
	SubscriptionStatusPaused            SubscriptionStatus = "paused"
)

// Subscription represents a subscription returned by the API.
type Subscription struct {
	ID                     uuid.UUID          `json:"id"`
	CustomerID             uuid.UUID          `json:"customer_id"`
	Provider               Provider           `json:"provider"`
	ProviderSubscriptionID string             `json:"provider_subscription_id"`
	Status                 SubscriptionStatus `json:"status"`
	Amount                 int64              `json:"amount"`
	Currency               Currency           `json:"currency"`
	Interval               string             `json:"interval"`
	IntervalCount          int                `json:"interval_count"`
	CurrentPeriodStart     time.Time          `json:"current_period_start"`
	CurrentPeriodEnd       time.Time          `json:"current_period_end"`
	TrialStart             *time.Time         `json:"trial_start,omitempty"`
	TrialEnd               *time.Time         `json:"trial_end,omitempty"`
	CancelAt               *time.Time         `json:"cancel_at,omitempty"`
	CanceledAt             *time.Time         `json:"canceled_at,omitempty"`
	CancelAtPeriodEnd      bool               `json:"cancel_at_period_end"`
	LatestPaymentID        *uuid.UUID         `json:"latest_payment_id,omitempty"`
	ProductName            string             `json:"product_name"`
	ProductDescription     *string            `json:"product_description,omitempty"`
	Metadata               map[string]any     `json:"metadata,omitempty"`
	CreatedAt              time.Time          `json:"created_at"`
	UpdatedAt              time.Time          `json:"updated_at"`
}

// CreateSubscriptionRequest is the request body for creating a subscription.
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

// UpdateSubscriptionRequest is the request body for updating a subscription.
type UpdateSubscriptionRequest struct {
	CancelAtPeriodEnd *bool          `json:"cancel_at_period_end,omitempty"`
	Metadata          map[string]any `json:"metadata,omitempty"`
}

// SubscriptionListResponse is the response for listing subscriptions.
type SubscriptionListResponse struct {
	Data   []Subscription `json:"data"`
	Total  int            `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

// --- Refund types ---

// RefundStatus represents the status of a refund.
type RefundStatus string

const (
	RefundStatusPending    RefundStatus = "pending"
	RefundStatusProcessing RefundStatus = "processing"
	RefundStatusSucceeded  RefundStatus = "succeeded"
	RefundStatusFailed     RefundStatus = "failed"
	RefundStatusCanceled   RefundStatus = "canceled"
)

// Refund represents a refund returned by the API.
type Refund struct {
	ID               uuid.UUID      `json:"id"`
	PaymentID        uuid.UUID      `json:"payment_id"`
	Provider         Provider       `json:"provider"`
	ProviderRefundID string         `json:"provider_refund_id"`
	Amount           int64          `json:"amount"`
	Currency         Currency       `json:"currency"`
	Status           RefundStatus   `json:"status"`
	Reason           *string        `json:"reason,omitempty"`
	Notes            *string        `json:"notes,omitempty"`
	FailureCode      *string        `json:"failure_code,omitempty"`
	FailureMessage   *string        `json:"failure_message,omitempty"`
	Metadata         map[string]any `json:"metadata,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	CompletedAt      *time.Time     `json:"completed_at,omitempty"`
}

// CreateRefundRequest is the request body for creating a refund.
type CreateRefundRequest struct {
	PaymentID uuid.UUID      `json:"payment_id"`
	Amount    int64          `json:"amount"`
	Reason    string         `json:"reason,omitempty"`
	Notes     string         `json:"notes,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// RefundListResponse is the response for listing refunds.
type RefundListResponse struct {
	Data   []Refund `json:"data"`
	Total  int      `json:"total"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
}

// --- Customer types ---

// Customer represents a customer returned by the API.
type Customer struct {
	ID               uuid.UUID      `json:"id"`
	UserID           uuid.UUID      `json:"user_id"`
	Email            string         `json:"email"`
	Name             string         `json:"name"`
	StripeCustomerID *string        `json:"stripe_customer_id,omitempty"`
	SwishCustomerID  *string        `json:"swish_customer_id,omitempty"`
	Metadata         map[string]any `json:"metadata,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        *time.Time     `json:"deleted_at,omitempty"`
}

// --- Error types ---

// APIError represents an error response from the payment API.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type apiErrorWrapper struct {
	Error APIError `json:"error"`
}
