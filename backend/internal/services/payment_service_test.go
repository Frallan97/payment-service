package services

import (
	"context"
	"errors"
	"payment-service/internal/models"
	"payment-service/internal/providers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPaymentRepository is a mock for PaymentRepository
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetByProviderPaymentID(ctx context.Context, provider models.Provider, providerPaymentID string) (*models.Payment, error) {
	args := m.Called(ctx, provider, providerPaymentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) ListByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]models.Payment, int, error) {
	args := m.Called(ctx, customerID, limit, offset)
	return args.Get(0).([]models.Payment), args.Int(1), args.Error(2)
}

func (m *MockPaymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

// MockCustomerRepository is a mock for CustomerRepository
type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Customer, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockCustomerRepository) Update(ctx context.Context, customer *models.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

// MockPaymentProvider is a mock for PaymentProvider
type MockPaymentProvider struct {
	mock.Mock
}

func (m *MockPaymentProvider) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockPaymentProvider) CreateCustomer(ctx context.Context, req *providers.CreateCustomerRequest) (*models.Customer, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockPaymentProvider) GetCustomer(ctx context.Context, providerCustomerID string) (*models.Customer, error) {
	args := m.Called(ctx, providerCustomerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockPaymentProvider) CreatePayment(ctx context.Context, req *providers.CreatePaymentRequest) (*models.Payment, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentProvider) GetPayment(ctx context.Context, providerPaymentID string) (*models.Payment, error) {
	args := m.Called(ctx, providerPaymentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentProvider) CancelPayment(ctx context.Context, providerPaymentID string) (*models.Payment, error) {
	args := m.Called(ctx, providerPaymentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentProvider) CreateSubscription(ctx context.Context, req *providers.CreateSubscriptionRequest) (*models.Subscription, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockPaymentProvider) GetSubscription(ctx context.Context, providerSubscriptionID string) (*models.Subscription, error) {
	args := m.Called(ctx, providerSubscriptionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockPaymentProvider) UpdateSubscription(ctx context.Context, providerSubscriptionID string, req *providers.UpdateSubscriptionRequest) (*models.Subscription, error) {
	args := m.Called(ctx, providerSubscriptionID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockPaymentProvider) CancelSubscription(ctx context.Context, providerSubscriptionID string, immediate bool) (*models.Subscription, error) {
	args := m.Called(ctx, providerSubscriptionID, immediate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockPaymentProvider) CreateRefund(ctx context.Context, req *providers.CreateRefundRequest) (*models.Refund, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Refund), args.Error(1)
}

func (m *MockPaymentProvider) GetRefund(ctx context.Context, providerRefundID string) (*models.Refund, error) {
	args := m.Called(ctx, providerRefundID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Refund), args.Error(1)
}

func (m *MockPaymentProvider) VerifyWebhookSignature(payload []byte, signature string) error {
	args := m.Called(payload, signature)
	return args.Error(0)
}

func (m *MockPaymentProvider) ParseWebhookEvent(payload []byte) (*providers.WebhookEvent, error) {
	args := m.Called(payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*providers.WebhookEvent), args.Error(1)
}

// MockProviderFactory is a mock for ProviderFactoryInterface
type MockProviderFactory struct {
	mock.Mock
}

func (m *MockProviderFactory) GetProvider(provider models.Provider) (providers.PaymentProvider, error) {
	args := m.Called(provider)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(providers.PaymentProvider), args.Error(1)
}

// Ensure Factory implements the interface
var _ ProviderFactoryInterface = (*MockProviderFactory)(nil)

func TestPaymentService_CreatePayment_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	userID := uuid.New()
	customerID := uuid.New()
	stripeCustomerID := "cus_test123"

	mockPaymentRepo := new(MockPaymentRepository)
	mockCustomerRepo := new(MockCustomerRepository)
	mockProvider := new(MockPaymentProvider)
	mockFactory := new(MockProviderFactory)

	service := NewPaymentService(mockPaymentRepo, mockCustomerRepo, mockFactory)

	// Test data
	req := &models.CreatePaymentRequest{
		Provider:    models.ProviderStripe,
		Amount:      10000,
		Currency:    models.CurrencySEK,
		Description: "Test payment",
	}

	existingCustomer := &models.Customer{
		ID:               customerID,
		UserID:           userID,
		Email:            "test@example.com",
		Name:             "Test User",
		StripeCustomerID: &stripeCustomerID,
	}

	providerPayment := &models.Payment{
		ID:                uuid.New(),
		Provider:          models.ProviderStripe,
		ProviderPaymentID: "pi_test123",
		Amount:            10000,
		Currency:          models.CurrencySEK,
		Status:            models.PaymentStatusPending,
	}

	// Mock expectations
	mockCustomerRepo.On("GetByUserID", ctx, userID).Return(existingCustomer, nil)
	mockFactory.On("GetProvider", models.ProviderStripe).Return(mockProvider, nil)
	mockProvider.On("CreatePayment", ctx, mock.AnythingOfType("*providers.CreatePaymentRequest")).Return(providerPayment, nil)
	mockPaymentRepo.On("Create", ctx, mock.AnythingOfType("*models.Payment")).Return(nil)

	// Execute
	result, err := service.CreatePayment(ctx, userID, "test@example.com", "Test User", req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, providerPayment.ProviderPaymentID, result.ProviderPaymentID)
	assert.Equal(t, customerID, result.CustomerID)

	mockCustomerRepo.AssertExpectations(t)
	mockFactory.AssertExpectations(t)
	mockProvider.AssertExpectations(t)
	mockPaymentRepo.AssertExpectations(t)
}

func TestPaymentService_CreatePayment_NewCustomer(t *testing.T) {
	// Setup
	ctx := context.Background()
	userID := uuid.New()
	stripeCustomerID := "cus_test123"

	mockPaymentRepo := new(MockPaymentRepository)
	mockCustomerRepo := new(MockCustomerRepository)
	mockProvider := new(MockPaymentProvider)
	mockFactory := new(MockProviderFactory)

	service := NewPaymentService(mockPaymentRepo, mockCustomerRepo, mockFactory)

	req := &models.CreatePaymentRequest{
		Provider:    models.ProviderStripe,
		Amount:      10000,
		Currency:    models.CurrencySEK,
		Description: "Test payment",
	}

	newCustomer := &models.Customer{
		ID:               uuid.New(),
		UserID:           userID,
		Email:            "test@example.com",
		Name:             "Test User",
		StripeCustomerID: &stripeCustomerID,
	}

	providerPayment := &models.Payment{
		ID:                uuid.New(),
		Provider:          models.ProviderStripe,
		ProviderPaymentID: "pi_test123",
		Amount:            10000,
		Currency:          models.CurrencySEK,
		Status:            models.PaymentStatusPending,
	}

	// Mock expectations - customer doesn't exist
	mockCustomerRepo.On("GetByUserID", ctx, userID).Return(nil, nil)
	mockFactory.On("GetProvider", models.ProviderStripe).Return(mockProvider, nil).Times(2)
	mockProvider.On("CreateCustomer", ctx, mock.AnythingOfType("*providers.CreateCustomerRequest")).Return(newCustomer, nil)
	mockCustomerRepo.On("Create", ctx, mock.AnythingOfType("*models.Customer")).Return(nil)
	mockProvider.On("CreatePayment", ctx, mock.AnythingOfType("*providers.CreatePaymentRequest")).Return(providerPayment, nil)
	mockPaymentRepo.On("Create", ctx, mock.AnythingOfType("*models.Payment")).Return(nil)

	// Execute
	result, err := service.CreatePayment(ctx, userID, "test@example.com", "Test User", req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockCustomerRepo.AssertExpectations(t)
	mockFactory.AssertExpectations(t)
	mockProvider.AssertExpectations(t)
	mockPaymentRepo.AssertExpectations(t)
}

func TestPaymentService_CreatePayment_ProviderError(t *testing.T) {
	// Setup
	ctx := context.Background()
	userID := uuid.New()
	customerID := uuid.New()
	stripeCustomerID := "cus_test123"

	mockPaymentRepo := new(MockPaymentRepository)
	mockCustomerRepo := new(MockCustomerRepository)
	mockProvider := new(MockPaymentProvider)
	mockFactory := new(MockProviderFactory)

	service := NewPaymentService(mockPaymentRepo, mockCustomerRepo, mockFactory)

	req := &models.CreatePaymentRequest{
		Provider:    models.ProviderStripe,
		Amount:      10000,
		Currency:    models.CurrencySEK,
		Description: "Test payment",
	}

	existingCustomer := &models.Customer{
		ID:               customerID,
		UserID:           userID,
		Email:            "test@example.com",
		Name:             "Test User",
		StripeCustomerID: &stripeCustomerID,
	}

	// Mock expectations
	mockCustomerRepo.On("GetByUserID", ctx, userID).Return(existingCustomer, nil)
	mockFactory.On("GetProvider", models.ProviderStripe).Return(mockProvider, nil)
	mockProvider.On("CreatePayment", ctx, mock.AnythingOfType("*providers.CreatePaymentRequest")).Return(nil, errors.New("provider error"))

	// Execute
	result, err := service.CreatePayment(ctx, userID, "test@example.com", "Test User", req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	apiErr, ok := err.(*models.APIError)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodePaymentFailed, apiErr.Code)
}

func TestPaymentService_GetPayment_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	userID := uuid.New()
	customerID := uuid.New()
	paymentID := uuid.New()

	mockPaymentRepo := new(MockPaymentRepository)
	mockCustomerRepo := new(MockCustomerRepository)
	mockFactory := new(MockProviderFactory)

	service := NewPaymentService(mockPaymentRepo, mockCustomerRepo, mockFactory)

	payment := &models.Payment{
		ID:                paymentID,
		CustomerID:        customerID,
		Provider:          models.ProviderStripe,
		ProviderPaymentID: "pi_test123",
		Amount:            10000,
		Currency:          models.CurrencySEK,
		Status:            models.PaymentStatusSucceeded,
	}

	customer := &models.Customer{
		ID:     customerID,
		UserID: userID,
		Email:  "test@example.com",
		Name:   "Test User",
	}

	// Mock expectations
	mockPaymentRepo.On("GetByID", ctx, paymentID).Return(payment, nil)
	mockCustomerRepo.On("GetByID", ctx, customerID).Return(customer, nil)

	// Execute
	result, err := service.GetPayment(ctx, paymentID, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, paymentID, result.ID)

	mockPaymentRepo.AssertExpectations(t)
	mockCustomerRepo.AssertExpectations(t)
}

func TestPaymentService_GetPayment_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	userID := uuid.New()
	paymentID := uuid.New()

	mockPaymentRepo := new(MockPaymentRepository)
	mockCustomerRepo := new(MockCustomerRepository)
	mockFactory := new(MockProviderFactory)

	service := NewPaymentService(mockPaymentRepo, mockCustomerRepo, mockFactory)

	// Mock expectations
	mockPaymentRepo.On("GetByID", ctx, paymentID).Return(nil, nil)

	// Execute
	result, err := service.GetPayment(ctx, paymentID, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	apiErr, ok := err.(*models.APIError)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeNotFound, apiErr.Code)
}

func TestPaymentService_GetPayment_UnauthorizedAccess(t *testing.T) {
	// Setup
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	customerID := uuid.New()
	paymentID := uuid.New()

	mockPaymentRepo := new(MockPaymentRepository)
	mockCustomerRepo := new(MockCustomerRepository)
	mockFactory := new(MockProviderFactory)

	service := NewPaymentService(mockPaymentRepo, mockCustomerRepo, mockFactory)

	payment := &models.Payment{
		ID:                paymentID,
		CustomerID:        customerID,
		Provider:          models.ProviderStripe,
		ProviderPaymentID: "pi_test123",
		Amount:            10000,
		Currency:          models.CurrencySEK,
		Status:            models.PaymentStatusSucceeded,
	}

	customer := &models.Customer{
		ID:     customerID,
		UserID: otherUserID, // Different user
		Email:  "other@example.com",
		Name:   "Other User",
	}

	// Mock expectations
	mockPaymentRepo.On("GetByID", ctx, paymentID).Return(payment, nil)
	mockCustomerRepo.On("GetByID", ctx, customerID).Return(customer, nil)

	// Execute
	result, err := service.GetPayment(ctx, paymentID, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	apiErr, ok := err.(*models.APIError)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeNotFound, apiErr.Code)
}

func TestPaymentService_ListPayments_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	userID := uuid.New()
	customerID := uuid.New()

	mockPaymentRepo := new(MockPaymentRepository)
	mockCustomerRepo := new(MockCustomerRepository)
	mockFactory := new(MockProviderFactory)

	service := NewPaymentService(mockPaymentRepo, mockCustomerRepo, mockFactory)

	customer := &models.Customer{
		ID:     customerID,
		UserID: userID,
		Email:  "test@example.com",
		Name:   "Test User",
	}

	payments := []models.Payment{
		{
			ID:                uuid.New(),
			CustomerID:        customerID,
			Provider:          models.ProviderStripe,
			ProviderPaymentID: "pi_test1",
			Amount:            10000,
			Currency:          models.CurrencySEK,
			Status:            models.PaymentStatusSucceeded,
		},
		{
			ID:                uuid.New(),
			CustomerID:        customerID,
			Provider:          models.ProviderStripe,
			ProviderPaymentID: "pi_test2",
			Amount:            20000,
			Currency:          models.CurrencySEK,
			Status:            models.PaymentStatusSucceeded,
		},
	}

	// Mock expectations
	mockCustomerRepo.On("GetByUserID", ctx, userID).Return(customer, nil)
	mockPaymentRepo.On("ListByCustomer", ctx, customerID, 20, 0).Return(payments, 2, nil)

	// Execute
	result, err := service.ListPayments(ctx, userID, 20, 0)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Data))
	assert.Equal(t, 2, result.Total)

	mockCustomerRepo.AssertExpectations(t)
	mockPaymentRepo.AssertExpectations(t)
}

func TestPaymentService_ListPayments_NoCustomer(t *testing.T) {
	// Setup
	ctx := context.Background()
	userID := uuid.New()

	mockPaymentRepo := new(MockPaymentRepository)
	mockCustomerRepo := new(MockCustomerRepository)
	mockFactory := new(MockProviderFactory)

	service := NewPaymentService(mockPaymentRepo, mockCustomerRepo, mockFactory)

	// Mock expectations
	mockCustomerRepo.On("GetByUserID", ctx, userID).Return(nil, nil)

	// Execute
	result, err := service.ListPayments(ctx, userID, 20, 0)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Data))
	assert.Equal(t, 0, result.Total)

	mockCustomerRepo.AssertExpectations(t)
}
