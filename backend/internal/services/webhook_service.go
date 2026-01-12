package services

import (
	"context"
	"encoding/json"
	"fmt"
	"payment-service/internal/models"
	"payment-service/internal/providers"
	"payment-service/internal/repository"
)

type WebhookService struct {
	webhookRepo      repository.WebhookRepositoryInterface
	paymentRepo      repository.PaymentRepositoryInterface
	subscriptionRepo repository.SubscriptionRepositoryInterface
	refundRepo       repository.RefundRepositoryInterface
}

func NewWebhookService(
	webhookRepo repository.WebhookRepositoryInterface,
	paymentRepo repository.PaymentRepositoryInterface,
	subscriptionRepo repository.SubscriptionRepositoryInterface,
	refundRepo repository.RefundRepositoryInterface,
) *WebhookService {
	return &WebhookService{
		webhookRepo:      webhookRepo,
		paymentRepo:      paymentRepo,
		subscriptionRepo: subscriptionRepo,
		refundRepo:       refundRepo,
	}
}

// ProcessWebhookEvent processes a webhook event with deduplication
func (s *WebhookService) ProcessWebhookEvent(ctx context.Context, event *providers.WebhookEvent, providerType string) error {
	// Convert provider string to models.Provider
	var provider models.Provider
	switch providerType {
	case "stripe":
		provider = models.ProviderStripe
	case "swish":
		provider = models.ProviderSwish
	default:
		return fmt.Errorf("unknown provider: %s", providerType)
	}

	// Check if event already processed (deduplication)
	existing, err := s.webhookRepo.GetByProviderEventID(ctx, provider, event.ID)
	if err != nil {
		return fmt.Errorf("failed to check for duplicate event: %w", err)
	}

	if existing != nil && existing.Processed {
		// Event already processed, skip
		return nil
	}

	// Marshal payload to JSON
	payloadJSON, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// Create webhook event record
	webhookEvent := &repository.WebhookEvent{
		Provider:        provider,
		ProviderEventID: event.ID,
		EventType:       event.Type,
		ResourceType:    &event.ResourceType,
		ResourceID:      &event.ResourceID,
		Processed:       false,
		Payload:         payloadJSON,
	}

	if existing != nil {
		// Event exists but not processed, use existing ID
		webhookEvent.ID = existing.ID
	} else {
		// Create new event record
		if err := s.webhookRepo.Create(ctx, webhookEvent); err != nil {
			return fmt.Errorf("failed to create webhook event record: %w", err)
		}
	}

	// Process event based on resource type
	var processingErr *string
	var processErr error

	switch event.ResourceType {
	case "payment":
		processErr = s.processPaymentEvent(ctx, event, provider)
	case "subscription":
		processErr = s.processSubscriptionEvent(ctx, event)
	case "refund":
		processErr = s.processRefundEvent(ctx, event)
	default:
		// Unknown resource type, just mark as processed
		processErr = nil
	}

	if processErr != nil {
		errMsg := processErr.Error()
		processingErr = &errMsg
	}

	// Mark event as processed
	if err := s.webhookRepo.MarkProcessed(ctx, webhookEvent.ID, processingErr); err != nil {
		return fmt.Errorf("failed to mark webhook event as processed: %w", err)
	}

	return processErr
}

// processPaymentEvent handles payment-related webhook events
func (s *WebhookService) processPaymentEvent(ctx context.Context, event *providers.WebhookEvent, provider models.Provider) error {
	if event.ResourceID == "" {
		return fmt.Errorf("payment event missing resource ID")
	}

	// Get payment from database by provider payment ID
	payment, err := s.paymentRepo.GetByProviderPaymentID(ctx, provider, event.ResourceID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	if payment == nil {
		// Payment not found in our database, might be from external source
		// Log and skip for now
		return nil
	}

	// Update payment status based on event type
	switch event.Type {
	case "payment_intent.succeeded":
		payment.Status = models.PaymentStatusSucceeded
	case "payment_intent.payment_failed":
		payment.Status = models.PaymentStatusFailed
	case "payment_intent.canceled":
		payment.Status = models.PaymentStatusCanceled
	case "payment_intent.processing":
		payment.Status = models.PaymentStatusProcessing
	default:
		// Unknown event type, skip
		return nil
	}

	// Update payment in database
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}

// processSubscriptionEvent handles subscription-related webhook events
func (s *WebhookService) processSubscriptionEvent(ctx context.Context, event *providers.WebhookEvent) error {
	if event.ResourceID == "" {
		return fmt.Errorf("subscription event missing resource ID")
	}

	// Get subscription from database by provider subscription ID
	subscription, err := s.subscriptionRepo.GetByProviderSubscriptionID(ctx, event.ResourceID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	if subscription == nil {
		// Subscription not found in our database
		return nil
	}

	// Update subscription status based on event type
	switch event.Type {
	case "customer.subscription.created":
		subscription.Status = models.SubscriptionStatusActive
	case "customer.subscription.updated":
		// Parse status from event - for now just keep existing status
		// In production, you'd want to parse the full subscription object from the payload
		// and map the status properly
		subscription.Status = models.SubscriptionStatusActive
	case "customer.subscription.deleted":
		subscription.Status = models.SubscriptionStatusCanceled
	case "customer.subscription.trial_will_end":
		// Notification event, no status change needed
		return nil
	default:
		// Unknown event type, skip
		return nil
	}

	// Update subscription in database
	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

// processRefundEvent handles refund-related webhook events
func (s *WebhookService) processRefundEvent(ctx context.Context, event *providers.WebhookEvent) error {
	if event.ResourceID == "" {
		return fmt.Errorf("refund event missing resource ID")
	}

	// Get refund from database by provider refund ID
	refund, err := s.refundRepo.GetByProviderRefundID(ctx, event.ResourceID)
	if err != nil {
		return fmt.Errorf("failed to get refund: %w", err)
	}

	if refund == nil {
		// Refund not found in our database
		return nil
	}

	// Update refund status based on event type
	switch event.Type {
	case "charge.refund.updated":
		// For now, just mark as succeeded
		// In production, parse the actual status from the payload
		refund.Status = models.RefundStatusSucceeded
	default:
		// Unknown event type, skip
		return nil
	}

	// Update refund in database
	if err := s.refundRepo.Update(ctx, refund); err != nil {
		return fmt.Errorf("failed to update refund: %w", err)
	}

	return nil
}
