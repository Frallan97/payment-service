package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Provider represents a payment provider
type Provider string

const (
	ProviderStripe Provider = "stripe"
	ProviderSwish  Provider = "swish"
)

// Currency represents a currency code
type Currency string

const (
	CurrencySEK Currency = "SEK"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending        PaymentStatus = "pending"
	PaymentStatusProcessing     PaymentStatus = "processing"
	PaymentStatusRequiresAction PaymentStatus = "requires_action"
	PaymentStatusSucceeded      PaymentStatus = "succeeded"
	PaymentStatusFailed         PaymentStatus = "failed"
	PaymentStatusCanceled       PaymentStatus = "canceled"
)

// SubscriptionStatus represents the status of a subscription
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

// RefundStatus represents the status of a refund
type RefundStatus string

const (
	RefundStatusPending    RefundStatus = "pending"
	RefundStatusProcessing RefundStatus = "processing"
	RefundStatusSucceeded  RefundStatus = "succeeded"
	RefundStatusFailed     RefundStatus = "failed"
	RefundStatusCanceled   RefundStatus = "canceled"
)

// JSONBMap is a map[string]any that implements sql.Scanner and driver.Valuer
// for PostgreSQL JSONB columns.
type JSONBMap map[string]any

// Scan implements sql.Scanner for reading JSONB from PostgreSQL.
func (m *JSONBMap) Scan(value any) error {
	if value == nil {
		*m = JSONBMap{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("JSONBMap.Scan: expected []byte, got %T", value)
	}
	return json.Unmarshal(b, m)
}

// Value implements driver.Valuer for writing JSONB to PostgreSQL.
func (m JSONBMap) Value() (driver.Value, error) {
	if m == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(m)
}

// ErrorCode represents API error codes
type ErrorCode string

const (
	ErrCodeInvalidRequest       ErrorCode = "invalid_request"
	ErrCodeAuthenticationFailed ErrorCode = "authentication_failed"
	ErrCodePaymentFailed        ErrorCode = "payment_failed"
	ErrCodeInsufficientFunds    ErrorCode = "insufficient_funds"
	ErrCodeProviderError        ErrorCode = "provider_error"
	ErrCodeNotFound             ErrorCode = "not_found"
	ErrCodeDuplicate            ErrorCode = "duplicate"
	ErrCodeRateLimitExceeded    ErrorCode = "rate_limit_exceeded"
)

// APIError represents an API error
type APIError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	StatusCode int                    `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAPIError creates a new API error
func NewAPIError(code ErrorCode, message string, statusCode int) *APIError {
	return &APIError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    make(map[string]interface{}),
	}
}
