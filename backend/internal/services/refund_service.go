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

type RefundService struct {
	refundRepo       repository.RefundRepositoryInterface
	paymentRepo      repository.PaymentRepositoryInterface
	customerRepo     repository.CustomerRepositoryInterface
	providerFactory  ProviderFactoryInterface
}

func NewRefundService(
	refundRepo repository.RefundRepositoryInterface,
	paymentRepo repository.PaymentRepositoryInterface,
	customerRepo repository.CustomerRepositoryInterface,
	providerFactory ProviderFactoryInterface,
) *RefundService {
	return &RefundService{
		refundRepo:      refundRepo,
		paymentRepo:     paymentRepo,
		customerRepo:    customerRepo,
		providerFactory: providerFactory,
	}
}

// CreateRefund creates a new refund
func (s *RefundService) CreateRefund(
	ctx context.Context,
	userID uuid.UUID,
	req *models.CreateRefundRequest,
) (*models.Refund, error) {
	// Get payment and verify ownership
	payment, err := s.paymentRepo.GetByID(ctx, req.PaymentID)
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

	// Check payment is refundable
	if payment.Status != models.PaymentStatusSucceeded {
		return nil, models.NewAPIError(
			models.ErrCodePaymentFailed,
			"Only successful payments can be refunded",
			http.StatusBadRequest,
		)
	}

	// Validate refund amount
	if req.Amount <= 0 {
		return nil, models.NewAPIError(
			models.ErrCodePaymentFailed,
			"Refund amount must be greater than zero",
			http.StatusBadRequest,
		)
	}

	if req.Amount > payment.Amount {
		return nil, models.NewAPIError(
			models.ErrCodePaymentFailed,
			"Refund amount cannot exceed payment amount",
			http.StatusBadRequest,
		)
	}

	// Check if payment has already been fully refunded
	existingRefunds, err := s.refundRepo.ListByPayment(ctx, payment.ID)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to check existing refunds",
			http.StatusInternalServerError,
		)
	}

	var totalRefunded int64
	for _, refund := range existingRefunds {
		if refund.Status == models.RefundStatusSucceeded {
			totalRefunded += refund.Amount
		}
	}

	if totalRefunded+req.Amount > payment.Amount {
		return nil, models.NewAPIError(
			models.ErrCodePaymentFailed,
			fmt.Sprintf("Cannot refund more than remaining amount. Already refunded: %d, Attempting: %d, Total: %d",
				totalRefunded, req.Amount, payment.Amount),
			http.StatusBadRequest,
		)
	}

	// Get provider
	provider, err := s.providerFactory.GetProvider(payment.Provider)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Provider not available",
			http.StatusBadRequest,
		)
	}

	// Create refund with provider
	providerReq := &providers.CreateRefundRequest{
		PaymentID: payment.ProviderPaymentID,
		Amount:    req.Amount,
		Reason:    req.Reason,
		Metadata:  convertMetadataToStrings(req.Metadata),
	}

	providerRefund, err := provider.CreateRefund(ctx, providerReq)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to create refund with provider",
			http.StatusBadGateway,
		)
	}

	// Save refund to database
	providerRefund.PaymentID = payment.ID
	if req.Notes != "" {
		providerRefund.Notes = &req.Notes
	}
	providerRefund.Metadata = req.Metadata

	if err := s.refundRepo.Create(ctx, providerRefund); err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to save refund to database",
			http.StatusInternalServerError,
		)
	}

	return providerRefund, nil
}

// GetRefund retrieves a refund by ID
func (s *RefundService) GetRefund(ctx context.Context, refundID, userID uuid.UUID) (*models.Refund, error) {
	refund, err := s.refundRepo.GetByID(ctx, refundID)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to retrieve refund",
			http.StatusInternalServerError,
		)
	}

	if refund == nil {
		return nil, models.NewAPIError(
			models.ErrCodeNotFound,
			"Refund not found",
			http.StatusNotFound,
		)
	}

	// Verify ownership through payment
	payment, err := s.paymentRepo.GetByID(ctx, refund.PaymentID)
	if err != nil || payment == nil {
		return nil, models.NewAPIError(
			models.ErrCodeNotFound,
			"Refund not found",
			http.StatusNotFound,
		)
	}

	customer, err := s.customerRepo.GetByID(ctx, payment.CustomerID)
	if err != nil || customer == nil || customer.UserID != userID {
		return nil, models.NewAPIError(
			models.ErrCodeNotFound,
			"Refund not found",
			http.StatusNotFound,
		)
	}

	return refund, nil
}

// ListRefunds lists refunds for a user
func (s *RefundService) ListRefunds(ctx context.Context, userID uuid.UUID, limit, offset int) (*models.RefundListResponse, error) {
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
		// No customer means no refunds
		return &models.RefundListResponse{
			Data:   []models.Refund{},
			Total:  0,
			Limit:  limit,
			Offset: offset,
		}, nil
	}

	// Get refunds
	refunds, total, err := s.refundRepo.ListByCustomer(ctx, customer.ID, limit, offset)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to list refunds",
			http.StatusInternalServerError,
		)
	}

	return &models.RefundListResponse{
		Data:   refunds,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// ListRefundsByPayment lists all refunds for a specific payment
func (s *RefundService) ListRefundsByPayment(ctx context.Context, paymentID, userID uuid.UUID) ([]models.Refund, error) {
	// Get payment and verify ownership
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

	customer, err := s.customerRepo.GetByID(ctx, payment.CustomerID)
	if err != nil || customer == nil || customer.UserID != userID {
		return nil, models.NewAPIError(
			models.ErrCodeNotFound,
			"Payment not found",
			http.StatusNotFound,
		)
	}

	// Get refunds for this payment
	refunds, err := s.refundRepo.ListByPayment(ctx, paymentID)
	if err != nil {
		return nil, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to list refunds",
			http.StatusInternalServerError,
		)
	}

	return refunds, nil
}
