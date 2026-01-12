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

type RefundHandler struct {
	refundService *services.RefundService
}

func NewRefundHandler(refundService *services.RefundService) *RefundHandler {
	return &RefundHandler{
		refundService: refundService,
	}
}

// CreateRefund handles POST /api/refunds
func (h *RefundHandler) CreateRefund(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"User not authenticated",
			http.StatusUnauthorized,
		))
		return
	}

	// Decode request
	var req models.CreateRefundRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Invalid request body",
			http.StatusBadRequest,
		))
		return
	}

	// Validate request
	if req.PaymentID == uuid.Nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Payment ID is required",
			http.StatusBadRequest,
		))
		return
	}

	if req.Amount <= 0 {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Refund amount must be greater than zero",
			http.StatusBadRequest,
		))
		return
	}

	// Create refund
	refund, err := h.refundService.CreateRefund(r.Context(), userID, &req)
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			WriteError(w, apiErr)
			return
		}
		WriteError(w, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to create refund",
			http.StatusInternalServerError,
		))
		return
	}

	WriteJSON(w, http.StatusCreated, refund)
}

// GetRefund handles GET /api/refunds/:id
func (h *RefundHandler) GetRefund(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"User not authenticated",
			http.StatusUnauthorized,
		))
		return
	}

	// Parse refund ID
	refundIDStr := chi.URLParam(r, "id")
	refundID, err := uuid.Parse(refundIDStr)
	if err != nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Invalid refund ID",
			http.StatusBadRequest,
		))
		return
	}

	// Get refund
	refund, err := h.refundService.GetRefund(r.Context(), refundID, userID)
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			WriteError(w, apiErr)
			return
		}
		WriteError(w, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to retrieve refund",
			http.StatusInternalServerError,
		))
		return
	}

	WriteJSON(w, http.StatusOK, refund)
}

// ListRefunds handles GET /api/refunds
func (h *RefundHandler) ListRefunds(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"User not authenticated",
			http.StatusUnauthorized,
		))
		return
	}

	// Parse pagination params
	limit := 20
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// List refunds
	response, err := h.refundService.ListRefunds(r.Context(), userID, limit, offset)
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			WriteError(w, apiErr)
			return
		}
		WriteError(w, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to list refunds",
			http.StatusInternalServerError,
		))
		return
	}

	WriteJSON(w, http.StatusOK, response)
}

// ListRefundsByPayment handles GET /api/payments/:id/refunds
func (h *RefundHandler) ListRefundsByPayment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"User not authenticated",
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

	// List refunds for this payment
	refunds, err := h.refundService.ListRefundsByPayment(r.Context(), paymentID, userID)
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			WriteError(w, apiErr)
			return
		}
		WriteError(w, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to list refunds",
			http.StatusInternalServerError,
		))
		return
	}

	WriteJSON(w, http.StatusOK, refunds)
}
