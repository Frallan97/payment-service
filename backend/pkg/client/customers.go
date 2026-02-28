package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetCurrentCustomer retrieves the customer record for the authenticated user.
func (c *Client) GetCurrentCustomer(ctx context.Context) (*Customer, error) {
	data, err := c.do(ctx, "GET", "/api/customers/me", nil)
	if err != nil {
		return nil, err
	}
	var customer Customer
	if err := json.Unmarshal(data, &customer); err != nil {
		return nil, fmt.Errorf("decode customer: %w", err)
	}
	return &customer, nil
}
