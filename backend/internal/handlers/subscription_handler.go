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

type SubscriptionHandler struct {
	subscriptionService *services.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// CreateSubscription handles POST /api/subscriptions
func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	// Get user info from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"User not authenticated",
			http.StatusUnauthorized,
		))
		return
	}
	email, _ := middleware.GetEmailFromContext(r.Context())
	name, _ := middleware.GetNameFromContext(r.Context())

	// Decode request
	var req models.CreateSubscriptionRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Invalid request body",
			http.StatusBadRequest,
		))
		return
	}

	// Validate request
	if req.Amount <= 0 {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Amount must be greater than zero",
			http.StatusBadRequest,
		))
		return
	}

	if req.Currency == "" {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Currency is required",
			http.StatusBadRequest,
		))
		return
	}

	if req.Interval == "" {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Interval is required",
			http.StatusBadRequest,
		))
		return
	}

	if req.IntervalCount <= 0 {
		req.IntervalCount = 1
	}

	if req.ProductName == "" {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Product name is required",
			http.StatusBadRequest,
		))
		return
	}

	// Create subscription
	subscription, err := h.subscriptionService.CreateSubscription(
		r.Context(),
		userID,
		email,
		name,
		&req,
	)
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			WriteError(w, apiErr)
			return
		}
		WriteError(w, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to create subscription",
			http.StatusInternalServerError,
		))
		return
	}

	WriteJSON(w, http.StatusCreated, subscription)
}

// GetSubscription handles GET /api/subscriptions/:id
func (h *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"User not authenticated",
			http.StatusUnauthorized,
		))
		return
	}

	// Parse subscription ID
	subscriptionIDStr := chi.URLParam(r, "id")
	subscriptionID, err := uuid.Parse(subscriptionIDStr)
	if err != nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Invalid subscription ID",
			http.StatusBadRequest,
		))
		return
	}

	// Get subscription
	subscription, err := h.subscriptionService.GetSubscription(r.Context(), subscriptionID, userID)
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			WriteError(w, apiErr)
			return
		}
		WriteError(w, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to retrieve subscription",
			http.StatusInternalServerError,
		))
		return
	}

	WriteJSON(w, http.StatusOK, subscription)
}

// ListSubscriptions handles GET /api/subscriptions
func (h *SubscriptionHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
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

	// List subscriptions
	response, err := h.subscriptionService.ListSubscriptions(r.Context(), userID, limit, offset)
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			WriteError(w, apiErr)
			return
		}
		WriteError(w, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to list subscriptions",
			http.StatusInternalServerError,
		))
		return
	}

	WriteJSON(w, http.StatusOK, response)
}

// UpdateSubscription handles PATCH /api/subscriptions/:id
func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"User not authenticated",
			http.StatusUnauthorized,
		))
		return
	}

	// Parse subscription ID
	subscriptionIDStr := chi.URLParam(r, "id")
	subscriptionID, err := uuid.Parse(subscriptionIDStr)
	if err != nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Invalid subscription ID",
			http.StatusBadRequest,
		))
		return
	}

	// Decode request
	var req models.UpdateSubscriptionRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Invalid request body",
			http.StatusBadRequest,
		))
		return
	}

	// Update subscription
	subscription, err := h.subscriptionService.UpdateSubscription(
		r.Context(),
		subscriptionID,
		userID,
		&req,
	)
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			WriteError(w, apiErr)
			return
		}
		WriteError(w, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to update subscription",
			http.StatusInternalServerError,
		))
		return
	}

	WriteJSON(w, http.StatusOK, subscription)
}

// CancelSubscription handles DELETE /api/subscriptions/:id
func (h *SubscriptionHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"User not authenticated",
			http.StatusUnauthorized,
		))
		return
	}

	// Parse subscription ID
	subscriptionIDStr := chi.URLParam(r, "id")
	subscriptionID, err := uuid.Parse(subscriptionIDStr)
	if err != nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Invalid subscription ID",
			http.StatusBadRequest,
		))
		return
	}

	// Parse query param for immediate cancellation
	immediate := r.URL.Query().Get("immediate") == "true"

	// Cancel subscription
	subscription, err := h.subscriptionService.CancelSubscription(
		r.Context(),
		subscriptionID,
		userID,
		immediate,
	)
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			WriteError(w, apiErr)
			return
		}
		WriteError(w, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to cancel subscription",
			http.StatusInternalServerError,
		))
		return
	}

	WriteJSON(w, http.StatusOK, subscription)
}
