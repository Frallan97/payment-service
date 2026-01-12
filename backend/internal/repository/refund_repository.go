package repository

import (
	"context"
	"database/sql"
	"fmt"
	"payment-service/internal/models"

	"github.com/google/uuid"
)

type RefundRepository struct {
	db *sql.DB
}

func NewRefundRepository(db *sql.DB) *RefundRepository {
	return &RefundRepository{db: db}
}

// Create creates a new refund
func (r *RefundRepository) Create(ctx context.Context, refund *models.Refund) error {
	query := `
		INSERT INTO refunds (
			payment_id, provider, provider_refund_id,
			amount, currency, status, reason, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		refund.PaymentID,
		refund.Provider,
		refund.ProviderRefundID,
		refund.Amount,
		refund.Currency,
		refund.Status,
		refund.Reason,
		refund.Metadata,
	).Scan(&refund.ID, &refund.CreatedAt, &refund.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create refund: %w", err)
	}

	return nil
}

// GetByID retrieves a refund by ID
func (r *RefundRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Refund, error) {
	query := `
		SELECT
			id, payment_id, provider, provider_refund_id,
			amount, currency, status, reason, metadata,
			created_at, updated_at
		FROM refunds
		WHERE id = $1 AND deleted_at IS NULL`

	refund := &models.Refund{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&refund.ID,
		&refund.PaymentID,
		&refund.Provider,
		&refund.ProviderRefundID,
		&refund.Amount,
		&refund.Currency,
		&refund.Status,
		&refund.Reason,
		&refund.Metadata,
		&refund.CreatedAt,
		&refund.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get refund: %w", err)
	}

	return refund, nil
}

// GetByProviderRefundID retrieves a refund by provider refund ID
func (r *RefundRepository) GetByProviderRefundID(ctx context.Context, providerRefundID string) (*models.Refund, error) {
	query := `
		SELECT
			id, payment_id, provider, provider_refund_id,
			amount, currency, status, reason, metadata,
			created_at, updated_at
		FROM refunds
		WHERE provider_refund_id = $1 AND deleted_at IS NULL`

	refund := &models.Refund{}
	err := r.db.QueryRowContext(ctx, query, providerRefundID).Scan(
		&refund.ID,
		&refund.PaymentID,
		&refund.Provider,
		&refund.ProviderRefundID,
		&refund.Amount,
		&refund.Currency,
		&refund.Status,
		&refund.Reason,
		&refund.Metadata,
		&refund.CreatedAt,
		&refund.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get refund by provider ID: %w", err)
	}

	return refund, nil
}

// ListByPayment retrieves all refunds for a payment
func (r *RefundRepository) ListByPayment(ctx context.Context, paymentID uuid.UUID) ([]models.Refund, error) {
	query := `
		SELECT
			id, payment_id, provider, provider_refund_id,
			amount, currency, status, reason, metadata,
			created_at, updated_at
		FROM refunds
		WHERE payment_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list refunds by payment: %w", err)
	}
	defer rows.Close()

	var refunds []models.Refund
	for rows.Next() {
		var refund models.Refund
		err := rows.Scan(
			&refund.ID,
			&refund.PaymentID,
			&refund.Provider,
			&refund.ProviderRefundID,
			&refund.Amount,
			&refund.Currency,
			&refund.Status,
			&refund.Reason,
			&refund.Metadata,
			&refund.CreatedAt,
			&refund.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan refund: %w", err)
		}
		refunds = append(refunds, refund)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating refunds: %w", err)
	}

	return refunds, nil
}

// ListByCustomer retrieves all refunds for a customer with pagination
func (r *RefundRepository) ListByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]models.Refund, int, error) {
	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM refunds rf
		JOIN payments p ON rf.payment_id = p.id
		WHERE p.customer_id = $1 AND rf.deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, countQuery, customerID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count refunds: %w", err)
	}

	// Get refunds
	query := `
		SELECT
			rf.id, rf.payment_id, rf.provider, rf.provider_refund_id,
			rf.amount, rf.currency, rf.status, rf.reason, rf.metadata,
			rf.created_at, rf.updated_at
		FROM refunds rf
		JOIN payments p ON rf.payment_id = p.id
		WHERE p.customer_id = $1 AND rf.deleted_at IS NULL
		ORDER BY rf.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, customerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list refunds: %w", err)
	}
	defer rows.Close()

	var refunds []models.Refund
	for rows.Next() {
		var refund models.Refund
		err := rows.Scan(
			&refund.ID,
			&refund.PaymentID,
			&refund.Provider,
			&refund.ProviderRefundID,
			&refund.Amount,
			&refund.Currency,
			&refund.Status,
			&refund.Reason,
			&refund.Metadata,
			&refund.CreatedAt,
			&refund.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan refund: %w", err)
		}
		refunds = append(refunds, refund)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating refunds: %w", err)
	}

	return refunds, total, nil
}

// Update updates a refund status
func (r *RefundRepository) Update(ctx context.Context, refund *models.Refund) error {
	query := `
		UPDATE refunds SET
			status = $1,
			metadata = $2,
			updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		refund.Status,
		refund.Metadata,
		refund.ID,
	).Scan(&refund.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("refund not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update refund: %w", err)
	}

	return nil
}
