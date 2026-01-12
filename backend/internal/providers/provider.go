package providers

import (
	"context"
	"payment-service/internal/models"

	"github.com/google/uuid"
)

// PaymentProvider defines the interface all payment providers must implement
type PaymentProvider interface {
	// Provider identification
	Name() string

	// Customer management
	CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*models.Customer, error)
	GetCustomer(ctx context.Context, providerCustomerID string) (*models.Customer, error)

	// One-time payments
	CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*models.Payment, error)
	GetPayment(ctx context.Context, providerPaymentID string) (*models.Payment, error)
	CancelPayment(ctx context.Context, providerPaymentID string) (*models.Payment, error)

	// Subscriptions
	CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*models.Subscription, error)
	GetSubscription(ctx context.Context, providerSubscriptionID string) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, providerSubscriptionID string, req *UpdateSubscriptionRequest) (*models.Subscription, error)
	CancelSubscription(ctx context.Context, providerSubscriptionID string, immediate bool) (*models.Subscription, error)

	// Refunds
	CreateRefund(ctx context.Context, req *CreateRefundRequest) (*models.Refund, error)
	GetRefund(ctx context.Context, providerRefundID string) (*models.Refund, error)

	// Webhooks
	VerifyWebhookSignature(payload []byte, signature string) error
	ParseWebhookEvent(payload []byte) (*WebhookEvent, error)
}

// CreateCustomerRequest represents a request to create a customer
type CreateCustomerRequest struct {
	UserID   uuid.UUID
	Email    string
	Name     string
	Metadata map[string]string
}

// CreatePaymentRequest represents a request to create a payment
type CreatePaymentRequest struct {
	CustomerID          string
	Amount              int64
	Currency            string
	Description         string
	StatementDescriptor string
	Metadata            map[string]string
	IdempotencyKey      string
}

// CreateSubscriptionRequest represents a request to create a subscription
type CreateSubscriptionRequest struct {
	CustomerID         string
	Amount             int64
	Currency           string
	Interval           string
	IntervalCount      int
	ProductName        string
	ProductDescription string
	TrialPeriodDays    int
	Metadata           map[string]string
}

// UpdateSubscriptionRequest represents a request to update a subscription
type UpdateSubscriptionRequest struct {
	CancelAtPeriodEnd *bool
	Metadata          map[string]string
}

// CreateRefundRequest represents a request to create a refund
type CreateRefundRequest struct {
	PaymentID string
	Amount    int64
	Reason    string
	Metadata  map[string]string
}

// WebhookEvent represents a parsed webhook event
type WebhookEvent struct {
	ID           string
	Type         string
	Provider     string
	ResourceType string // payment, subscription, refund
	ResourceID   string
	Status       string
	Payload      map[string]any
}
