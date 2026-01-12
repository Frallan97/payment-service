package handlers

import (
	"io"
	"net/http"
	"payment-service/internal/models"
	"payment-service/internal/providers"
)

type WebhookHandler struct {
	providerFactory *providers.Factory
}

func NewWebhookHandler(providerFactory *providers.Factory) *WebhookHandler {
	return &WebhookHandler{
		providerFactory: providerFactory,
	}
}

// HandleStripeWebhook handles POST /api/webhooks/stripe
func (h *WebhookHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	// Read body
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Failed to read webhook payload",
			http.StatusBadRequest,
		))
		return
	}

	// Get Stripe signature header
	signature := r.Header.Get("Stripe-Signature")
	if signature == "" {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Missing Stripe-Signature header",
			http.StatusBadRequest,
		))
		return
	}

	// Get Stripe provider
	provider, err := h.providerFactory.GetProvider(models.ProviderStripe)
	if err != nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeProviderError,
			"Stripe provider not available",
			http.StatusInternalServerError,
		))
		return
	}

	// Verify webhook signature
	if err := provider.VerifyWebhookSignature(payload, signature); err != nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Invalid webhook signature",
			http.StatusBadRequest,
		))
		return
	}

	// Parse webhook event
	webhookEvent, err := provider.ParseWebhookEvent(payload)
	if err != nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Failed to parse webhook event",
			http.StatusBadRequest,
		))
		return
	}

	// TODO: Process webhook event with webhook service
	// This would involve:
	// 1. Checking for duplicate events using webhook_events table
	// 2. Processing the event based on type (payment.succeeded, subscription.updated, etc.)
	// 3. Updating database records accordingly
	// 4. Storing the event in webhook_events table for deduplication

	// For now, just log that we received the webhook
	// In production, this should be processed asynchronously
	_ = webhookEvent

	// Return 200 OK to acknowledge receipt
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"received": true}`))
}

// HandleSwishWebhook handles POST /api/webhooks/swish
func (h *WebhookHandler) HandleSwishWebhook(w http.ResponseWriter, r *http.Request) {
	// Swish webhook implementation
	// TODO: Implement when Swish provider is ready
	WriteError(w, models.NewAPIError(
		models.ErrCodeProviderError,
		"Swish webhooks not yet implemented",
		http.StatusNotImplemented,
	))
}
