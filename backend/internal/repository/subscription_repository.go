package repository

import (
	"context"
	"database/sql"
	"fmt"
	"payment-service/internal/models"

	"github.com/google/uuid"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create creates a new subscription
func (r *SubscriptionRepository) Create(ctx context.Context, subscription *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (
			customer_id, provider, provider_subscription_id,
			amount, currency, interval, interval_count,
			status, current_period_start, current_period_end,
			trial_start, trial_end, cancel_at_period_end,
			canceled_at, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		) RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		subscription.CustomerID,
		subscription.Provider,
		subscription.ProviderSubscriptionID,
		subscription.Amount,
		subscription.Currency,
		subscription.Interval,
		subscription.IntervalCount,
		subscription.Status,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.TrialStart,
		subscription.TrialEnd,
		subscription.CancelAtPeriodEnd,
		subscription.CanceledAt,
		subscription.Metadata,
	).Scan(&subscription.ID, &subscription.CreatedAt, &subscription.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

// GetByID retrieves a subscription by ID
func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	query := `
		SELECT
			id, customer_id, provider, provider_subscription_id,
			amount, currency, interval, interval_count,
			status, current_period_start, current_period_end,
			trial_start, trial_end, cancel_at_period_end,
			canceled_at, metadata, created_at, updated_at
		FROM subscriptions
		WHERE id = $1 AND deleted_at IS NULL`

	subscription := &models.Subscription{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.CustomerID,
		&subscription.Provider,
		&subscription.ProviderSubscriptionID,
		&subscription.Amount,
		&subscription.Currency,
		&subscription.Interval,
		&subscription.IntervalCount,
		&subscription.Status,
		&subscription.CurrentPeriodStart,
		&subscription.CurrentPeriodEnd,
		&subscription.TrialStart,
		&subscription.TrialEnd,
		&subscription.CancelAtPeriodEnd,
		&subscription.CanceledAt,
		&subscription.Metadata,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return subscription, nil
}

// GetByProviderSubscriptionID retrieves a subscription by provider subscription ID
func (r *SubscriptionRepository) GetByProviderSubscriptionID(ctx context.Context, providerSubscriptionID string) (*models.Subscription, error) {
	query := `
		SELECT
			id, customer_id, provider, provider_subscription_id,
			amount, currency, interval, interval_count,
			status, current_period_start, current_period_end,
			trial_start, trial_end, cancel_at_period_end,
			canceled_at, metadata, created_at, updated_at
		FROM subscriptions
		WHERE provider_subscription_id = $1 AND deleted_at IS NULL`

	subscription := &models.Subscription{}
	err := r.db.QueryRowContext(ctx, query, providerSubscriptionID).Scan(
		&subscription.ID,
		&subscription.CustomerID,
		&subscription.Provider,
		&subscription.ProviderSubscriptionID,
		&subscription.Amount,
		&subscription.Currency,
		&subscription.Interval,
		&subscription.IntervalCount,
		&subscription.Status,
		&subscription.CurrentPeriodStart,
		&subscription.CurrentPeriodEnd,
		&subscription.TrialStart,
		&subscription.TrialEnd,
		&subscription.CancelAtPeriodEnd,
		&subscription.CanceledAt,
		&subscription.Metadata,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription by provider ID: %w", err)
	}

	return subscription, nil
}

// ListByCustomer retrieves all subscriptions for a customer with pagination
func (r *SubscriptionRepository) ListByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]models.Subscription, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM subscriptions WHERE customer_id = $1 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, countQuery, customerID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count subscriptions: %w", err)
	}

	// Get subscriptions
	query := `
		SELECT
			id, customer_id, provider, provider_subscription_id,
			amount, currency, interval, interval_count,
			status, current_period_start, current_period_end,
			trial_start, trial_end, cancel_at_period_end,
			canceled_at, metadata, created_at, updated_at
		FROM subscriptions
		WHERE customer_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, customerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []models.Subscription
	for rows.Next() {
		var subscription models.Subscription
		err := rows.Scan(
			&subscription.ID,
			&subscription.CustomerID,
			&subscription.Provider,
			&subscription.ProviderSubscriptionID,
			&subscription.Amount,
			&subscription.Currency,
			&subscription.Interval,
			&subscription.IntervalCount,
			&subscription.Status,
			&subscription.CurrentPeriodStart,
			&subscription.CurrentPeriodEnd,
			&subscription.TrialStart,
			&subscription.TrialEnd,
			&subscription.CancelAtPeriodEnd,
			&subscription.CanceledAt,
			&subscription.Metadata,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating subscriptions: %w", err)
	}

	return subscriptions, total, nil
}

// Update updates a subscription
func (r *SubscriptionRepository) Update(ctx context.Context, subscription *models.Subscription) error {
	query := `
		UPDATE subscriptions SET
			status = $1,
			current_period_start = $2,
			current_period_end = $3,
			cancel_at_period_end = $4,
			canceled_at = $5,
			metadata = $6,
			updated_at = NOW()
		WHERE id = $7 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		subscription.Status,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.CancelAtPeriodEnd,
		subscription.CanceledAt,
		subscription.Metadata,
		subscription.ID,
	).Scan(&subscription.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("subscription not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}
