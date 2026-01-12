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

type PaymentService struct {
	paymentRepo  repository.PaymentRepositoryInterface
	customerRepo repository.CustomerRepositoryInterface
	providerFactory ProviderFactoryInterface
}

// ProviderFactoryInterface defines the interface for provider factory
type ProviderFactoryInterface interface {
	GetProvider(provider models.Provider) (providers.PaymentProvider, error)
}

func NewPaymentService(
	paymentRepo repository.PaymentRepositoryInterface,
	customerRepo repository.CustomerRepositoryInterface,
	providerFactory ProviderFactoryInterface,
) *PaymentService {
	return &PaymentService{
		paymentRepo:     paymentRepo,
		customerRepo:    customerRepo,
		providerFactory: providerFactory,
	}
}

// CreatePayment creates a new payment
func (s *PaymentService) CreatePayment(
	ctx context.Context,
	userID uuid.UUID,
	email, name string,
	req *models.CreatePaymentRequest,
) (*models.Payment, error) {
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

	// Create payment with provider
	providerReq := &providers.CreatePaymentRequest{
		CustomerID:          providerCustomerID,
		Amount:              req.Amount,
		Currency:            string(req.Currency),
		Description:         req.Description,
		StatementDescriptor: req.StatementDescriptor,
		Metadata:            convertMetadataToStrings(req.Metadata),
	}

	providerPayment, err := provider.CreatePayment(ctx, providerReq)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodePaymentFailed,
			"Failed to create payment with provider",
			http.StatusBadGateway,
		)
	}

	// Save payment to database
	providerPayment.CustomerID = customer.ID
	providerPayment.Metadata = req.Metadata

	if err := s.paymentRepo.Create(ctx, providerPayment); err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to save payment to database",
			http.StatusInternalServerError,
		)
	}

	return providerPayment, nil
}

// GetPayment retrieves a payment by ID
func (s *PaymentService) GetPayment(ctx context.Context, paymentID, userID uuid.UUID) (*models.Payment, error) {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to retrieve payment",
			http.StatusInternalServerError,
		)
	}

	if payment == nil {
		return nil, models.NewAPIError(
			models.ErrCodeNotFound,
			"Payment not found",
			http.StatusNotFound,
		)
	}

	// Verify customer owns this payment
	customer, err := s.customerRepo.GetByID(ctx, payment.CustomerID)
	if err != nil || customer == nil || customer.UserID != userID {
		return nil, models.NewAPIError(
			models.ErrCodeNotFound,
			"Payment not found",
			http.StatusNotFound,
		)
	}

	return payment, nil
}

// ListPayments lists payments for a user
func (s *PaymentService) ListPayments(ctx context.Context, userID uuid.UUID, limit, offset int) (*models.PaymentListResponse, error) {
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
		// No customer means no payments
		return &models.PaymentListResponse{
			Data:   []models.Payment{},
			Total:  0,
			Limit:  limit,
			Offset: offset,
		}, nil
	}

	// Get payments
	payments, total, err := s.paymentRepo.ListByCustomer(ctx, customer.ID, limit, offset)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to list payments",
			http.StatusInternalServerError,
		)
	}

	return &models.PaymentListResponse{
		Data:   payments,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// getOrCreateCustomer gets an existing customer or creates a new one
func (s *PaymentService) getOrCreateCustomer(
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
func (s *PaymentService) addProviderToCustomer(
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

// convertMetadataToStrings converts map[string]any to map[string]string for provider
func convertMetadataToStrings(metadata map[string]any) map[string]string {
	result := make(map[string]string)
	for k, v := range metadata {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}
