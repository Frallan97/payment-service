package services

import (
	"context"
	"fmt"
	"net/http"
	"payment-service/internal/models"
	"payment-service/internal/providers"
	"payment-service/internal/repository"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	subscriptionRepo repository.SubscriptionRepositoryInterface
	customerRepo     repository.CustomerRepositoryInterface
	providerFactory  ProviderFactoryInterface
}

func NewSubscriptionService(
	subscriptionRepo repository.SubscriptionRepositoryInterface,
	customerRepo repository.CustomerRepositoryInterface,
	providerFactory ProviderFactoryInterface,
) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		customerRepo:     customerRepo,
		providerFactory:  providerFactory,
	}
}

// CreateSubscription creates a new subscription
func (s *SubscriptionService) CreateSubscription(
	ctx context.Context,
	userID uuid.UUID,
	email, name string,
	req *models.CreateSubscriptionRequest,
) (*models.Subscription, error) {
	// Get or create customer
	customer, err := s.getOrCreateCustomer(ctx, userID, email, name, req.Provider)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to get or create customer",
			http.StatusInternalServerError,
		)
	}

	// Get provider
	provider, err := s.providerFactory.GetProvider(req.Provider)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			fmt.Sprintf("Provider %s not available", req.Provider),
			http.StatusBadRequest,
		)
	}

	// Get provider customer ID
	var providerCustomerID string
	if req.Provider == models.ProviderStripe && customer.StripeCustomerID != nil {
		providerCustomerID = *customer.StripeCustomerID
	} else if req.Provider == models.ProviderSwish && customer.SwishCustomerID != nil {
		providerCustomerID = *customer.SwishCustomerID
	} else {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Customer not configured for this provider",
			http.StatusBadRequest,
		)
	}

	// Create subscription with provider
	providerReq := &providers.CreateSubscriptionRequest{
		CustomerID:         providerCustomerID,
		Amount:             req.Amount,
		Currency:           string(req.Currency),
		Interval:           req.Interval,
		IntervalCount:      req.IntervalCount,
		ProductName:        req.ProductName,
		ProductDescription: req.ProductDescription,
		TrialPeriodDays:    req.TrialPeriodDays,
		Metadata:           convertMetadataToStrings(req.Metadata),
	}

	providerSubscription, err := provider.CreateSubscription(ctx, providerReq)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to create subscription with provider",
			http.StatusBadGateway,
		)
	}

	// Save subscription to database
	providerSubscription.CustomerID = customer.ID
	providerSubscription.ProductName = req.ProductName
	if req.ProductDescription != "" {
		providerSubscription.ProductDescription = &req.ProductDescription
	}
	providerSubscription.Metadata = req.Metadata

	if err := s.subscriptionRepo.Create(ctx, providerSubscription); err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to save subscription to database",
			http.StatusInternalServerError,
		)
	}

	return providerSubscription, nil
}

// GetSubscription retrieves a subscription by ID
func (s *SubscriptionService) GetSubscription(ctx context.Context, subscriptionID, userID uuid.UUID) (*models.Subscription, error) {
	subscription, err := s.subscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to retrieve subscription",
			http.StatusInternalServerError,
		)
	}

	if subscription == nil {
		return nil, models.NewAPIError(
			models.ErrCodeNotFound,
			"Subscription not found",
			http.StatusNotFound,
		)
	}

	// Verify customer owns this subscription
	customer, err := s.customerRepo.GetByID(ctx, subscription.CustomerID)
	if err != nil || customer == nil || customer.UserID != userID {
		return nil, models.NewAPIError(
			models.ErrCodeNotFound,
			"Subscription not found",
			http.StatusNotFound,
		)
	}

	return subscription, nil
}

// ListSubscriptions lists subscriptions for a user
func (s *SubscriptionService) ListSubscriptions(ctx context.Context, userID uuid.UUID, limit, offset int) (*models.SubscriptionListResponse, error) {
	// Get customer
	customer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to retrieve customer",
			http.StatusInternalServerError,
		)
	}

	if customer == nil {
		// No customer means no subscriptions
		return &models.SubscriptionListResponse{
			Data:   []models.Subscription{},
			Total:  0,
			Limit:  limit,
			Offset: offset,
		}, nil
	}

	// Get subscriptions
	subscriptions, total, err := s.subscriptionRepo.ListByCustomer(ctx, customer.ID, limit, offset)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to list subscriptions",
			http.StatusInternalServerError,
		)
	}

	return &models.SubscriptionListResponse{
		Data:   subscriptions,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// UpdateSubscription updates a subscription
func (s *SubscriptionService) UpdateSubscription(
	ctx context.Context,
	subscriptionID, userID uuid.UUID,
	req *models.UpdateSubscriptionRequest,
) (*models.Subscription, error) {
	// Get and verify ownership
	subscription, err := s.GetSubscription(ctx, subscriptionID, userID)
	if err != nil {
		return nil, err
	}

	// Get provider
	provider, err := s.providerFactory.GetProvider(subscription.Provider)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Provider not available",
			http.StatusBadRequest,
		)
	}

	// Update with provider
	providerReq := &providers.UpdateSubscriptionRequest{
		CancelAtPeriodEnd: req.CancelAtPeriodEnd,
		Metadata:          convertMetadataToStrings(req.Metadata),
	}

	updatedSubscription, err := provider.UpdateSubscription(
		ctx,
		subscription.ProviderSubscriptionID,
		providerReq,
	)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to update subscription with provider",
			http.StatusBadGateway,
		)
	}

	// Update in database
	subscription.Status = updatedSubscription.Status
	subscription.CancelAtPeriodEnd = updatedSubscription.CancelAtPeriodEnd
	subscription.CanceledAt = updatedSubscription.CanceledAt
	if req.Metadata != nil {
		subscription.Metadata = req.Metadata
	}

	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to update subscription in database",
			http.StatusInternalServerError,
		)
	}

	return subscription, nil
}

// CancelSubscription cancels a subscription
func (s *SubscriptionService) CancelSubscription(
	ctx context.Context,
	subscriptionID, userID uuid.UUID,
	immediate bool,
) (*models.Subscription, error) {
	// Get and verify ownership
	subscription, err := s.GetSubscription(ctx, subscriptionID, userID)
	if err != nil {
		return nil, err
	}

	// Get provider
	provider, err := s.providerFactory.GetProvider(subscription.Provider)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Provider not available",
			http.StatusBadRequest,
		)
	}

	// Cancel with provider
	canceledSubscription, err := provider.CancelSubscription(
		ctx,
		subscription.ProviderSubscriptionID,
		immediate,
	)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to cancel subscription with provider",
			http.StatusBadGateway,
		)
	}

	// Update in database
	subscription.Status = canceledSubscription.Status
	subscription.CancelAtPeriodEnd = canceledSubscription.CancelAtPeriodEnd
	subscription.CanceledAt = canceledSubscription.CanceledAt

	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to update subscription in database",
			http.StatusInternalServerError,
		)
	}

	return subscription, nil
}

// getOrCreateCustomer gets an existing customer or creates a new one
func (s *SubscriptionService) getOrCreateCustomer(
	ctx context.Context,
	userID uuid.UUID,
	email, name string,
	provider models.Provider,
) (*models.Customer, error) {
	// Check if customer exists
	customer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Customer exists, check if they have provider customer ID
	if customer != nil {
		hasProviderID := (provider == models.ProviderStripe && customer.StripeCustomerID != nil) ||
			(provider == models.ProviderSwish && customer.SwishCustomerID != nil)

		if hasProviderID {
			return customer, nil
		}

		// Need to create provider customer
		return s.addProviderToCustomer(ctx, customer, provider)
	}

	// Create new customer
	providerInstance, err := s.providerFactory.GetProvider(provider)
	if err != nil {
		return nil, err
	}

	providerCustomer, err := providerInstance.CreateCustomer(ctx, &providers.CreateCustomerRequest{
		UserID: userID,
		Email:  email,
		Name:   name,
		Metadata: map[string]string{
			"user_id": userID.String(),
		},
	})
	if err != nil {
		return nil, err
	}

	// Save to database
	if err := s.customerRepo.Create(ctx, providerCustomer); err != nil {
		return nil, err
	}

	return providerCustomer, nil
}

// addProviderToCustomer adds a provider customer ID to an existing customer
func (s *SubscriptionService) addProviderToCustomer(
	ctx context.Context,
	customer *models.Customer,
	provider models.Provider,
) (*models.Customer, error) {
	providerInstance, err := s.providerFactory.GetProvider(provider)
	if err != nil {
		return nil, err
	}

	providerCustomer, err := providerInstance.CreateCustomer(ctx, &providers.CreateCustomerRequest{
		UserID: customer.UserID,
		Email:  customer.Email,
		Name:   customer.Name,
		Metadata: map[string]string{
			"user_id": customer.UserID.String(),
		},
	})
	if err != nil {
		return nil, err
	}

	// Update customer with new provider ID
	if provider == models.ProviderStripe {
		customer.StripeCustomerID = providerCustomer.StripeCustomerID
	} else if provider == models.ProviderSwish {
		customer.SwishCustomerID = providerCustomer.SwishCustomerID
	}

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}
