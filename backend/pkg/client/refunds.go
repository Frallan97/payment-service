package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// CreateRefund creates a new refund for a payment.
func (c *Client) CreateRefund(ctx context.Context, req *CreateRefundRequest) (*Refund, error) {
	data, err := c.do(ctx, "POST", "/api/refunds", req)
	if err != nil {
		return nil, err
	}
	var refund Refund
	if err := json.Unmarshal(data, &refund); err != nil {
		return nil, fmt.Errorf("decode refund: %w", err)
	}
	return &refund, nil
}

// GetRefund retrieves a refund by ID.
func (c *Client) GetRefund(ctx context.Context, id uuid.UUID) (*Refund, error) {
	data, err := c.do(ctx, "GET", "/api/refunds/"+id.String(), nil)
	if err != nil {
		return nil, err
	}
	var refund Refund
	if err := json.Unmarshal(data, &refund); err != nil {
		return nil, fmt.Errorf("decode refund: %w", err)
	}
	return &refund, nil
}

// ListRefunds lists refunds with pagination.
func (c *Client) ListRefunds(ctx context.Context, limit, offset int) (*RefundListResponse, error) {
	path := fmt.Sprintf("/api/refunds?limit=%d&offset=%d", limit, offset)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp RefundListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode refund list: %w", err)
	}
	return &resp, nil
}

// ListRefundsByPayment lists refunds for a specific payment.
func (c *Client) ListRefundsByPayment(ctx context.Context, paymentID uuid.UUID) (*RefundListResponse, error) {
	path := fmt.Sprintf("/api/payments/%s/refunds", paymentID.String())
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp RefundListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode refund list: %w", err)
	}
	return &resp, nil
}
