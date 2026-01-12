package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"payment-service/internal/models"
	"time"

	"github.com/google/uuid"
)

type WebhookEvent struct {
	ID               uuid.UUID      `db:"id"`
	Provider         models.Provider `db:"provider"`
	ProviderEventID  string         `db:"provider_event_id"`
	EventType        string         `db:"event_type"`
	ResourceType     *string        `db:"resource_type"`
	ResourceID       *string        `db:"resource_id"`
	Processed        bool           `db:"processed"`
	ProcessedAt      *time.Time     `db:"processed_at"`
	Payload          json.RawMessage `db:"payload"`
	ProcessingError  *string        `db:"processing_error"`
	CreatedAt        time.Time      `db:"created_at"`
}

type WebhookRepository struct {
	db *sql.DB
}

func NewWebhookRepository(db *sql.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

// Create creates a new webhook event record
func (r *WebhookRepository) Create(ctx context.Context, event *WebhookEvent) error {
	query := `
		INSERT INTO webhook_events (
			provider, provider_event_id, event_type,
			resource_type, resource_id, processed,
			payload
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING id, created_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		event.Provider,
		event.ProviderEventID,
		event.EventType,
		event.ResourceType,
		event.ResourceID,
		event.Processed,
		event.Payload,
	).Scan(&event.ID, &event.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create webhook event: %w", err)
	}

	return nil
}

// GetByProviderEventID retrieves a webhook event by provider event ID
func (r *WebhookRepository) GetByProviderEventID(ctx context.Context, provider models.Provider, providerEventID string) (*WebhookEvent, error) {
	query := `
		SELECT
			id, provider, provider_event_id, event_type,
			resource_type, resource_id, processed, processed_at,
			payload, processing_error, created_at
		FROM webhook_events
		WHERE provider = $1 AND provider_event_id = $2`

	event := &WebhookEvent{}
	err := r.db.QueryRowContext(ctx, query, provider, providerEventID).Scan(
		&event.ID,
		&event.Provider,
		&event.ProviderEventID,
		&event.EventType,
		&event.ResourceType,
		&event.ResourceID,
		&event.Processed,
		&event.ProcessedAt,
		&event.Payload,
		&event.ProcessingError,
		&event.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook event: %w", err)
	}

	return event, nil
}

// MarkProcessed marks a webhook event as processed
func (r *WebhookRepository) MarkProcessed(ctx context.Context, id uuid.UUID, processingError *string) error {
	query := `
		UPDATE webhook_events SET
			processed = true,
			processed_at = NOW(),
			processing_error = $1
		WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, processingError, id)
	if err != nil {
		return fmt.Errorf("failed to mark webhook event as processed: %w", err)
	}

	return nil
}

// CleanupOldEvents deletes processed webhook events older than the specified duration
func (r *WebhookRepository) CleanupOldEvents(ctx context.Context, olderThan time.Duration) (int64, error) {
	query := `
		DELETE FROM webhook_events
		WHERE processed = true
		AND processed_at < $1`

	cutoffTime := time.Now().Add(-olderThan)
	result, err := r.db.ExecContext(ctx, query, cutoffTime)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old webhook events: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
