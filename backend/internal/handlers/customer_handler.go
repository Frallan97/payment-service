package handlers

import (
	"net/http"
	"payment-service/internal/middleware"
	"payment-service/internal/models"
	"payment-service/internal/repository"
)

type CustomerHandler struct {
	customerRepo *repository.CustomerRepository
}

func NewCustomerHandler(customerRepo *repository.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{
		customerRepo: customerRepo,
	}
}

// GetMe handles GET /api/customers/me
func (h *CustomerHandler) GetMe(w http.ResponseWriter, r *http.Request) {
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

	// Get customer
	customer, err := h.customerRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeProviderError,
			"Failed to retrieve customer",
			http.StatusInternalServerError,
		))
		return
	}

	if customer == nil {
		WriteError(w, models.NewAPIError(
			models.ErrCodeNotFound,
			"Customer not found",
			http.StatusNotFound,
		))
		return
	}

	WriteJSON(w, http.StatusOK, customer)
}
