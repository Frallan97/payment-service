package repository

import (
	"context"
	"database/sql"
	"fmt"
	"payment-service/internal/models"

	"github.com/google/uuid"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// Create inserts a new customer
func (r *CustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	query := `
		INSERT INTO customers (user_id, email, name, stripe_customer_id, swish_customer_id, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		customer.UserID,
		customer.Email,
		customer.Name,
		customer.StripeCustomerID,
		customer.SwishCustomerID,
		customer.Metadata,
	).Scan(&customer.ID, &customer.CreatedAt, &customer.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}

	return nil
}

// GetByUserID retrieves a customer by user ID
func (r *CustomerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Customer, error) {
	query := `
		SELECT id, user_id, email, name, stripe_customer_id, swish_customer_id,
		       metadata, created_at, updated_at, deleted_at
		FROM customers
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	customer := &models.Customer{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&customer.ID,
		&customer.UserID,
		&customer.Email,
		&customer.Name,
		&customer.StripeCustomerID,
		&customer.SwishCustomerID,
		&customer.Metadata,
		&customer.CreatedAt,
		&customer.UpdatedAt,
		&customer.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

// GetByID retrieves a customer by ID
func (r *CustomerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	query := `
		SELECT id, user_id, email, name, stripe_customer_id, swish_customer_id,
		       metadata, created_at, updated_at, deleted_at
		FROM customers
		WHERE id = $1 AND deleted_at IS NULL
	`

	customer := &models.Customer{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&customer.ID,
		&customer.UserID,
		&customer.Email,
		&customer.Name,
		&customer.StripeCustomerID,
		&customer.SwishCustomerID,
		&customer.Metadata,
		&customer.CreatedAt,
		&customer.UpdatedAt,
		&customer.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

// Update updates a customer's information
func (r *CustomerRepository) Update(ctx context.Context, customer *models.Customer) error {
	query := `
		UPDATE customers
		SET email = $2, name = $3, stripe_customer_id = $4,
		    swish_customer_id = $5, metadata = $6
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		customer.ID,
		customer.Email,
		customer.Name,
		customer.StripeCustomerID,
		customer.SwishCustomerID,
		customer.Metadata,
	).Scan(&customer.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	return nil
}
