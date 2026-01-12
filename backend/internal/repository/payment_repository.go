package repository

import (
	"context"
	"database/sql"
	"fmt"
	"payment-service/internal/models"

	"github.com/google/uuid"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create inserts a new payment
func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	query := `
		INSERT INTO payments (
			customer_id, provider, provider_payment_id, amount, currency, status,
			payment_method_type, payment_method_details, description, statement_descriptor,
			subscription_id, invoice_id, client_secret, failure_code, failure_message,
			metadata, idempotency_key
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		payment.CustomerID,
		payment.Provider,
		payment.ProviderPaymentID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.PaymentMethodType,
		payment.PaymentMethodDetails,
		payment.Description,
		payment.StatementDescriptor,
		payment.SubscriptionID,
		payment.InvoiceID,
		payment.ClientSecret,
		payment.FailureCode,
		payment.FailureMessage,
		payment.Metadata,
		payment.IdempotencyKey,
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

// GetByID retrieves a payment by ID
func (r *PaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	query := `
		SELECT id, customer_id, provider, provider_payment_id, amount, currency, status,
		       payment_method_type, payment_method_details, description, statement_descriptor,
		       subscription_id, invoice_id, client_secret, failure_code, failure_message,
		       metadata, idempotency_key, created_at, updated_at, completed_at
		FROM payments
		WHERE id = $1
	`

	payment := &models.Payment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.CustomerID,
		&payment.Provider,
		&payment.ProviderPaymentID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.PaymentMethodType,
		&payment.PaymentMethodDetails,
		&payment.Description,
		&payment.StatementDescriptor,
		&payment.SubscriptionID,
		&payment.InvoiceID,
		&payment.ClientSecret,
		&payment.FailureCode,
		&payment.FailureMessage,
		&payment.Metadata,
		&payment.IdempotencyKey,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&payment.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return payment, nil
}

// GetByProviderPaymentID retrieves a payment by provider payment ID
func (r *PaymentRepository) GetByProviderPaymentID(ctx context.Context, provider models.Provider, providerPaymentID string) (*models.Payment, error) {
	query := `
		SELECT id, customer_id, provider, provider_payment_id, amount, currency, status,
		       payment_method_type, payment_method_details, description, statement_descriptor,
		       subscription_id, invoice_id, client_secret, failure_code, failure_message,
		       metadata, idempotency_key, created_at, updated_at, completed_at
		FROM payments
		WHERE provider = $1 AND provider_payment_id = $2
	`

	payment := &models.Payment{}
	err := r.db.QueryRowContext(ctx, query, provider, providerPaymentID).Scan(
		&payment.ID,
		&payment.CustomerID,
		&payment.Provider,
		&payment.ProviderPaymentID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.PaymentMethodType,
		&payment.PaymentMethodDetails,
		&payment.Description,
		&payment.StatementDescriptor,
		&payment.SubscriptionID,
		&payment.InvoiceID,
		&payment.ClientSecret,
		&payment.FailureCode,
		&payment.FailureMessage,
		&payment.Metadata,
		&payment.IdempotencyKey,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&payment.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return payment, nil
}

// ListByCustomer retrieves payments for a customer
func (r *PaymentRepository) ListByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]models.Payment, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM payments WHERE customer_id = $1`
	if err := r.db.QueryRowContext(ctx, countQuery, customerID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %w", err)
	}

	// Get payments
	query := `
		SELECT id, customer_id, provider, provider_payment_id, amount, currency, status,
		       payment_method_type, payment_method_details, description, statement_descriptor,
		       subscription_id, invoice_id, client_secret, failure_code, failure_message,
		       metadata, idempotency_key, created_at, updated_at, completed_at
		FROM payments
		WHERE customer_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, customerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	payments := []models.Payment{}
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.ID,
			&payment.CustomerID,
			&payment.Provider,
			&payment.ProviderPaymentID,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.PaymentMethodType,
			&payment.PaymentMethodDetails,
			&payment.Description,
			&payment.StatementDescriptor,
			&payment.SubscriptionID,
			&payment.InvoiceID,
			&payment.ClientSecret,
			&payment.FailureCode,
			&payment.FailureMessage,
			&payment.Metadata,
			&payment.IdempotencyKey,
			&payment.CreatedAt,
			&payment.UpdatedAt,
			&payment.CompletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating payments: %w", err)
	}

	return payments, total, nil
}

// Update updates a payment
func (r *PaymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	query := `
		UPDATE payments
		SET status = $2, payment_method_type = $3, payment_method_details = $4,
		    failure_code = $5, failure_message = $6, completed_at = $7
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		payment.ID,
		payment.Status,
		payment.PaymentMethodType,
		payment.PaymentMethodDetails,
		payment.FailureCode,
		payment.FailureMessage,
		payment.CompletedAt,
	).Scan(&payment.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}
