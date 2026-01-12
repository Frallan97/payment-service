package providers

import (
	"fmt"
	"payment-service/internal/models"
)

// Factory creates payment providers based on configuration
type Factory struct {
	stripeProvider PaymentProvider
}

// NewFactory creates a new provider factory
func NewFactory(stripeAPIKey, stripeWebhookSecret string) *Factory {
	return &Factory{
		stripeProvider: NewStripeProvider(stripeAPIKey, stripeWebhookSecret),
	}
}

// GetProvider returns a provider by name
func (f *Factory) GetProvider(provider models.Provider) (PaymentProvider, error) {
	switch provider {
	case models.ProviderStripe:
		return f.stripeProvider, nil
	case models.ProviderSwish:
		return nil, fmt.Errorf("swish provider not yet implemented")
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}
