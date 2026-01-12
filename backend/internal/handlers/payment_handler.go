package handlers

import (
	"net/http"
	"payment-service/internal/middleware"
	"payment-service/internal/models"
	"payment-service/internal/services"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreatePayment handles POST /api/payments
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, models.NewAPIError(
			models.ErrCodeAuthenticationFailed,
			"User ID not found in context",
			http.StatusUnauthorized,
		))
		return
	}

	email, _ := middleware.GetEmailFromContext(r.Context())
	name, _ := middleware.GetNameFromContext(r.Context())

	// Parse request
	var req models.CreatePaymentRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, err)
		return
	}

	// Validate request
	if req.Amount <= 0 {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Amount must be greater than 0",
			http.StatusBadRequest,
		))
		return
	}

	if req.Currency == "" {
		req.Currency = models.CurrencySEK // Default
	}

	if req.Provider == "" {
		req.Provider = models.ProviderStripe // Default
	}

	// Create payment
	payment, err := h.paymentService.CreatePayment(r.Context(), userID, email, name, &req)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, payment)
}

// GetPayment handles GET /api/payments/:id
func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, models.NewAPIError(
			models.ErrCodeAuthenticationFailed,
			"User ID not found in context",
			http.StatusUnauthorized,
		))
		return
	}

	// Parse payment ID
	paymentIDStr := chi.URLParam(r, "id")
	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Invalid payment ID",
			http.StatusBadRequest,
		))
		return
	}

	// Get payment
	payment, err := h.paymentService.GetPayment(r.Context(), paymentID, userID)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, payment)
}

// ListPayments handles GET /api/payments
func (h *PaymentHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, models.NewAPIError(
			models.ErrCodeAuthenticationFailed,
			"User ID not found in context",
			http.StatusUnauthorized,
		))
		return
	}

	// Parse query parameters
	limit := 20
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 100 {
				limit = 100 // Max limit
			}
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// List payments
	response, err := h.paymentService.ListPayments(r.Context(), userID, limit, offset)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, response)
}
