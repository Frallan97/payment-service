package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// CreateSubscription creates a new subscription.
func (c *Client) CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*Subscription, error) {
	data, err := c.do(ctx, "POST", "/api/subscriptions", req)
	if err != nil {
		return nil, err
	}
	var sub Subscription
	if err := json.Unmarshal(data, &sub); err != nil {
		return nil, fmt.Errorf("decode subscription: %w", err)
	}
	return &sub, nil
}

// GetSubscription retrieves a subscription by ID.
func (c *Client) GetSubscription(ctx context.Context, id uuid.UUID) (*Subscription, error) {
	data, err := c.do(ctx, "GET", "/api/subscriptions/"+id.String(), nil)
	if err != nil {
		return nil, err
	}
	var sub Subscription
	if err := json.Unmarshal(data, &sub); err != nil {
		return nil, fmt.Errorf("decode subscription: %w", err)
	}
	return &sub, nil
}

// ListSubscriptions lists subscriptions with pagination.
func (c *Client) ListSubscriptions(ctx context.Context, limit, offset int) (*SubscriptionListResponse, error) {
	path := fmt.Sprintf("/api/subscriptions?limit=%d&offset=%d", limit, offset)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp SubscriptionListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode subscription list: %w", err)
	}
	return &resp, nil
}

// UpdateSubscription updates a subscription.
func (c *Client) UpdateSubscription(ctx context.Context, id uuid.UUID, req *UpdateSubscriptionRequest) (*Subscription, error) {
	data, err := c.do(ctx, "PATCH", "/api/subscriptions/"+id.String(), req)
	if err != nil {
		return nil, err
	}
	var sub Subscription
	if err := json.Unmarshal(data, &sub); err != nil {
		return nil, fmt.Errorf("decode subscription: %w", err)
	}
	return &sub, nil
}

// CancelSubscription cancels a subscription. If immediate is true, it cancels
// right away; otherwise it cancels at the end of the current billing period.
func (c *Client) CancelSubscription(ctx context.Context, id uuid.UUID, immediate bool) (*Subscription, error) {
	path := fmt.Sprintf("/api/subscriptions/%s?immediate=%t", id.String(), immediate)
	data, err := c.do(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	var sub Subscription
	if err := json.Unmarshal(data, &sub); err != nil {
		return nil, fmt.Errorf("decode subscription: %w", err)
	}
	return &sub, nil
}
