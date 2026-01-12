package providers

import (
	"context"
	"fmt"
	"payment-service/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"github.com/stripe/stripe-go/v78/price"
	"github.com/stripe/stripe-go/v78/refund"
	"github.com/stripe/stripe-go/v78/subscription"
	"github.com/stripe/stripe-go/v78/webhook"
)

type StripeProvider struct {
	apiKey        string
	webhookSecret string
}

// NewStripeProvider creates a new Stripe provider
func NewStripeProvider(apiKey, webhookSecret string) *StripeProvider {
	stripe.Key = apiKey
	return &StripeProvider{
		apiKey:        apiKey,
		webhookSecret: webhookSecret,
	}
}

// Name returns the provider name
func (p *StripeProvider) Name() string {
	return "stripe"
}

// CreateCustomer creates a customer in Stripe
func (p *StripeProvider) CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*models.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(req.Email),
		Name:  stripe.String(req.Name),
	}

	for k, v := range req.Metadata {
		params.AddMetadata(k, v)
	}
	params.AddMetadata("user_id", req.UserID.String())

	stripeCustomer, err := customer.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe: failed to create customer: %w", err)
	}

	return &models.Customer{
		UserID:           req.UserID,
		Email:            req.Email,
		Name:             req.Name,
		StripeCustomerID: &stripeCustomer.ID,
	}, nil
}

// GetCustomer retrieves a customer from Stripe
func (p *StripeProvider) GetCustomer(ctx context.Context, providerCustomerID string) (*models.Customer, error) {
	stripeCustomer, err := customer.Get(providerCustomerID, nil)
	if err != nil {
		return nil, fmt.Errorf("stripe: failed to get customer: %w", err)
	}

	userIDStr := stripeCustomer.Metadata["user_id"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("stripe: invalid user_id in metadata: %w", err)
	}

	return &models.Customer{
		UserID:           userID,
		Email:            stripeCustomer.Email,
		Name:             stripeCustomer.Name,
		StripeCustomerID: &stripeCustomer.ID,
	}, nil
}

// CreatePayment creates a payment intent in Stripe
func (p *StripeProvider) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*models.Payment, error) {
	params := &stripe.PaymentIntentParams{
		Amount:      stripe.Int64(req.Amount),
		Currency:    stripe.String(req.Currency),
		Customer:    stripe.String(req.CustomerID),
		Description: stripe.String(req.Description),
		Confirm:     stripe.Bool(false),
	}

	if req.StatementDescriptor != "" {
		params.StatementDescriptor = stripe.String(req.StatementDescriptor)
	}

	for k, v := range req.Metadata {
		params.AddMetadata(k, v)
	}

	if req.IdempotencyKey != "" {
		params.IdempotencyKey = stripe.String(req.IdempotencyKey)
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe: failed to create payment intent: %w", err)
	}

	return mapPaymentIntentToPayment(pi), nil
}

// GetPayment retrieves a payment intent from Stripe
func (p *StripeProvider) GetPayment(ctx context.Context, providerPaymentID string) (*models.Payment, error) {
	pi, err := paymentintent.Get(providerPaymentID, nil)
	if err != nil {
		return nil, fmt.Errorf("stripe: failed to get payment intent: %w", err)
	}

	return mapPaymentIntentToPayment(pi), nil
}

// CancelPayment cancels a payment intent in Stripe
func (p *StripeProvider) CancelPayment(ctx context.Context, providerPaymentID string) (*models.Payment, error) {
	pi, err := paymentintent.Cancel(providerPaymentID, nil)
	if err != nil {
		return nil, fmt.Errorf("stripe: failed to cancel payment intent: %w", err)
	}

	return mapPaymentIntentToPayment(pi), nil
}

// VerifyWebhookSignature verifies the signature of a Stripe webhook
func (p *StripeProvider) VerifyWebhookSignature(payload []byte, signature string) error {
	_, err := webhook.ConstructEvent(payload, signature, p.webhookSecret)
	if err != nil {
		return fmt.Errorf("stripe: webhook signature verification failed: %w", err)
	}
	return nil
}

// ParseWebhookEvent parses a Stripe webhook event
func (p *StripeProvider) ParseWebhookEvent(payload []byte) (*WebhookEvent, error) {
	event, err := webhook.ConstructEvent(payload, "", p.webhookSecret)
	if err != nil {
		return nil, fmt.Errorf("stripe: failed to parse webhook event: %w", err)
	}

	webhookEvent := &WebhookEvent{
		ID:       event.ID,
		Type:     string(event.Type),
		Provider: "stripe",
		Payload:  make(map[string]any),
	}

	// Determine resource type based on event type
	switch event.Type {
	case "payment_intent.succeeded", "payment_intent.payment_failed",
		"payment_intent.canceled", "payment_intent.processing":
		webhookEvent.ResourceType = "payment"
	case "customer.created", "customer.updated", "customer.deleted":
		webhookEvent.ResourceType = "customer"
	}

	return webhookEvent, nil
}

// mapPaymentIntentToPayment converts a Stripe PaymentIntent to our Payment model
func mapPaymentIntentToPayment(pi *stripe.PaymentIntent) *models.Payment {
	payment := &models.Payment{
		Provider:          models.ProviderStripe,
		ProviderPaymentID: pi.ID,
		Amount:            pi.Amount,
		Currency:          models.Currency(pi.Currency),
		Status:            mapStripeStatus(string(pi.Status)),
		ClientSecret:      &pi.ClientSecret,
	}

	if pi.Description != "" {
		payment.Description = &pi.Description
	}

	if pi.StatementDescriptor != "" {
		payment.StatementDescriptor = &pi.StatementDescriptor
	}

	return payment
}

// CreateSubscription creates a subscription in Stripe
func (p *StripeProvider) CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*models.Subscription, error) {
	// Create price for the subscription
	priceParams := &stripe.PriceParams{
		Currency:   stripe.String(req.Currency),
		UnitAmount: stripe.Int64(req.Amount),
		Recurring: &stripe.PriceRecurringParams{
			Interval:      stripe.String(req.Interval),
			IntervalCount: stripe.Int64(int64(req.IntervalCount)),
		},
		Product: stripe.String("prod_payment_service"), // Use a generic product or create dynamically
	}

	priceObj, err := price.New(priceParams)
	if err != nil {
		return nil, fmt.Errorf("stripe: failed to create price: %w", err)
	}

	// Create subscription
	subParams := &stripe.SubscriptionParams{
		Customer: stripe.String(req.CustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			{Price: stripe.String(priceObj.ID)},
		},
	}

	if req.TrialPeriodDays > 0 {
		subParams.TrialPeriodDays = stripe.Int64(int64(req.TrialPeriodDays))
	}

	for k, v := range req.Metadata {
		subParams.AddMetadata(k, v)
	}

	sub, err := subscription.New(subParams)
	if err != nil {
		return nil, fmt.Errorf("stripe: failed to create subscription: %w", err)
	}

	return mapStripeSubscription(sub), nil
}

// GetSubscription retrieves a subscription from Stripe
func (p *StripeProvider) GetSubscription(ctx context.Context, providerSubscriptionID string) (*models.Subscription, error) {
	sub, err := subscription.Get(providerSubscriptionID, nil)
	if err != nil {
		return nil, fmt.Errorf("stripe: failed to get subscription: %w", err)
	}

	return mapStripeSubscription(sub), nil
}

// UpdateSubscription updates a subscription in Stripe
func (p *StripeProvider) UpdateSubscription(ctx context.Context, providerSubscriptionID string, req *UpdateSubscriptionRequest) (*models.Subscription, error) {
	params := &stripe.SubscriptionParams{}

	if req.CancelAtPeriodEnd != nil {
		params.CancelAtPeriodEnd = stripe.Bool(*req.CancelAtPeriodEnd)
	}

	for k, v := range req.Metadata {
		params.AddMetadata(k, v)
	}

	sub, err := subscription.Update(providerSubscriptionID, params)
	if err != nil {
		return nil, fmt.Errorf("stripe: failed to update subscription: %w", err)
	}

	return mapStripeSubscription(sub), nil
}

// CancelSubscription cancels a subscription in Stripe
func (p *StripeProvider) CancelSubscription(ctx context.Context, providerSubscriptionID string, immediate bool) (*models.Subscription, error) {
	if immediate {
		// Cancel immediately
		sub, err := subscription.Cancel(providerSubscriptionID, nil)
		if err != nil {
			return nil, fmt.Errorf("stripe: failed to cancel subscription: %w", err)
		}
		return mapStripeSubscription(sub), nil
	}

	// Cancel at period end
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}
	sub, err := subscription.Update(providerSubscriptionID, params)
	if err != nil {
		return nil, fmt.Errorf("stripe: failed to schedule cancellation: %w", err)
	}

	return mapStripeSubscription(sub), nil
}

// CreateRefund creates a refund in Stripe
func (p *StripeProvider) CreateRefund(ctx context.Context, req *CreateRefundRequest) (*models.Refund, error) {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(req.PaymentID),
	}

	if req.Amount > 0 {
		params.Amount = stripe.Int64(req.Amount)
	}

	if req.Reason != "" {
		params.Reason = stripe.String(req.Reason)
	}

	for k, v := range req.Metadata {
		params.AddMetadata(k, v)
	}

	ref, err := refund.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe: failed to create refund: %w", err)
	}

	return mapStripeRefund(ref), nil
}

// GetRefund retrieves a refund from Stripe
func (p *StripeProvider) GetRefund(ctx context.Context, providerRefundID string) (*models.Refund, error) {
	ref, err := refund.Get(providerRefundID, nil)
	if err != nil {
		return nil, fmt.Errorf("stripe: failed to get refund: %w", err)
	}

	return mapStripeRefund(ref), nil
}

// mapStripeSubscription converts a Stripe Subscription to our Subscription model
func mapStripeSubscription(sub *stripe.Subscription) *models.Subscription {
	subscription := &models.Subscription{
		Provider:               models.ProviderStripe,
		ProviderSubscriptionID: sub.ID,
		Amount:                 sub.Items.Data[0].Price.UnitAmount,
		Currency:               models.Currency(sub.Items.Data[0].Price.Currency),
		Interval:               string(sub.Items.Data[0].Price.Recurring.Interval),
		IntervalCount:          int(sub.Items.Data[0].Price.Recurring.IntervalCount),
		Status:                 mapStripeSubscriptionStatus(string(sub.Status)),
		CurrentPeriodStart:     time.Unix(sub.CurrentPeriodStart, 0),
		CurrentPeriodEnd:       time.Unix(sub.CurrentPeriodEnd, 0),
		CancelAtPeriodEnd:      sub.CancelAtPeriodEnd,
	}

	if sub.TrialStart != 0 {
		trialStart := time.Unix(sub.TrialStart, 0)
		subscription.TrialStart = &trialStart
	}

	if sub.TrialEnd != 0 {
		trialEnd := time.Unix(sub.TrialEnd, 0)
		subscription.TrialEnd = &trialEnd
	}

	if sub.CanceledAt != 0 {
		canceledAt := time.Unix(sub.CanceledAt, 0)
		subscription.CanceledAt = &canceledAt
	}

	return subscription
}

// mapStripeSubscriptionStatus maps Stripe subscription status to our SubscriptionStatus
func mapStripeSubscriptionStatus(stripeStatus string) models.SubscriptionStatus {
	switch stripeStatus {
	case "active":
		return models.SubscriptionStatusActive
	case "past_due":
		return models.SubscriptionStatusPastDue
	case "unpaid":
		return models.SubscriptionStatusUnpaid
	case "canceled":
		return models.SubscriptionStatusCanceled
	case "incomplete":
		return models.SubscriptionStatusIncomplete
	case "incomplete_expired":
		return models.SubscriptionStatusIncompleteExpired
	case "trialing":
		return models.SubscriptionStatusTrialing
	case "paused":
		return models.SubscriptionStatusPaused
	default:
		return models.SubscriptionStatusIncomplete
	}
}

// mapStripeRefund converts a Stripe Refund to our Refund model
func mapStripeRefund(ref *stripe.Refund) *models.Refund {
	refund := &models.Refund{
		Provider:        models.ProviderStripe,
		ProviderRefundID: ref.ID,
		Amount:          ref.Amount,
		Currency:        models.Currency(ref.Currency),
		Status:          mapStripeRefundStatus(string(ref.Status)),
	}

	if ref.Reason != "" {
		reason := string(ref.Reason)
		refund.Reason = &reason
	}

	return refund
}

// mapStripeRefundStatus maps Stripe refund status to our RefundStatus
func mapStripeRefundStatus(stripeStatus string) models.RefundStatus {
	switch stripeStatus {
	case "pending":
		return models.RefundStatusPending
	case "succeeded":
		return models.RefundStatusSucceeded
	case "failed":
		return models.RefundStatusFailed
	case "canceled":
		return models.RefundStatusCanceled
	default:
		return models.RefundStatusPending
	}
}

// mapStripeStatus maps Stripe payment intent status to our PaymentStatus
func mapStripeStatus(stripeStatus string) models.PaymentStatus {
	switch stripeStatus {
	case "requires_payment_method", "requires_confirmation":
		return models.PaymentStatusPending
	case "requires_action":
		return models.PaymentStatusRequiresAction
	case "processing":
		return models.PaymentStatusProcessing
	case "succeeded":
		return models.PaymentStatusSucceeded
	case "canceled":
		return models.PaymentStatusCanceled
	default:
		return models.PaymentStatusFailed
	}
}
