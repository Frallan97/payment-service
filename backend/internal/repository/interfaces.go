package repository

import (
	"context"
	"payment-service/internal/models"

	"github.com/google/uuid"
)

// PaymentRepositoryInterface defines the interface for payment repository operations
type PaymentRepositoryInterface interface {
	Create(ctx context.Context, payment *models.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error)
	GetByProviderPaymentID(ctx context.Context, provider models.Provider, providerPaymentID string) (*models.Payment, error)
	ListByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]models.Payment, int, error)
	Update(ctx context.Context, payment *models.Payment) error
}

// CustomerRepositoryInterface defines the interface for customer repository operations
type CustomerRepositoryInterface interface {
	Create(ctx context.Context, customer *models.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Customer, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Customer, error)
	Update(ctx context.Context, customer *models.Customer) error
}

// SubscriptionRepositoryInterface defines the interface for subscription repository operations
type SubscriptionRepositoryInterface interface {
	Create(ctx context.Context, subscription *models.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	GetByProviderSubscriptionID(ctx context.Context, providerSubscriptionID string) (*models.Subscription, error)
	ListByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]models.Subscription, int, error)
	Update(ctx context.Context, subscription *models.Subscription) error
}

// RefundRepositoryInterface defines the interface for refund repository operations
type RefundRepositoryInterface interface {
	Create(ctx context.Context, refund *models.Refund) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Refund, error)
	GetByProviderRefundID(ctx context.Context, providerRefundID string) (*models.Refund, error)
	ListByPayment(ctx context.Context, paymentID uuid.UUID) ([]models.Refund, error)
	ListByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]models.Refund, int, error)
	Update(ctx context.Context, refund *models.Refund) error
}

// WebhookRepositoryInterface defines the interface for webhook repository operations
type WebhookRepositoryInterface interface {
	Create(ctx context.Context, event *WebhookEvent) error
	GetByProviderEventID(ctx context.Context, provider models.Provider, providerEventID string) (*WebhookEvent, error)
	MarkProcessed(ctx context.Context, id uuid.UUID, processingError *string) error
}
