package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_service_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "payment_service_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// Payment-specific metrics
	paymentsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_service_payments_total",
			Help: "Total number of payments created",
		},
		[]string{"provider", "status"},
	)

	paymentAmount = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "payment_service_payment_amount",
			Help:    "Payment amounts",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 500000, 1000000},
		},
		[]string{"provider", "currency"},
	)

	// Subscription metrics
	subscriptionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_service_subscriptions_total",
			Help: "Total number of subscriptions",
		},
		[]string{"provider", "status"},
	)

	subscriptionsActive = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "payment_service_subscriptions_active",
			Help: "Number of active subscriptions",
		},
		[]string{"provider"},
	)

	// Refund metrics
	refundsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_service_refunds_total",
			Help: "Total number of refunds",
		},
		[]string{"provider", "status"},
	)

	refundAmount = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "payment_service_refund_amount",
			Help:    "Refund amounts",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 500000, 1000000},
		},
		[]string{"provider", "currency"},
	)

	// Webhook metrics
	webhooksReceived = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_service_webhooks_received_total",
			Help: "Total number of webhooks received",
		},
		[]string{"provider", "event_type"},
	)

	webhookProcessingErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_service_webhook_processing_errors_total",
			Help: "Total number of webhook processing errors",
		},
		[]string{"provider", "event_type"},
	)

	webhookProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "payment_service_webhook_processing_duration_seconds",
			Help:    "Webhook processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"provider", "event_type"},
	)
)

// MetricsMiddleware records HTTP request metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		ww := &metricsResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call next handler
		next.ServeHTTP(ww, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(ww.statusCode)

		// Get route pattern from chi context
		routePattern := chi.RouteContext(r.Context()).RoutePattern()
		if routePattern == "" {
			routePattern = r.URL.Path
		}

		httpRequestsTotal.WithLabelValues(r.Method, routePattern, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, routePattern, status).Observe(duration)
	})
}

// metricsResponseWriter wraps http.ResponseWriter to capture status code
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *metricsResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// RecordPayment records a payment metric
func RecordPayment(provider, status string, amount int64, currency string) {
	paymentsTotal.WithLabelValues(provider, status).Inc()
	paymentAmount.WithLabelValues(provider, currency).Observe(float64(amount))
}

// RecordSubscription records a subscription metric
func RecordSubscription(provider, status string) {
	subscriptionsTotal.WithLabelValues(provider, status).Inc()
}

// UpdateActiveSubscriptions updates the active subscriptions gauge
func UpdateActiveSubscriptions(provider string, count int) {
	subscriptionsActive.WithLabelValues(provider).Set(float64(count))
}

// RecordRefund records a refund metric
func RecordRefund(provider, status string, amount int64, currency string) {
	refundsTotal.WithLabelValues(provider, status).Inc()
	refundAmount.WithLabelValues(provider, currency).Observe(float64(amount))
}

// RecordWebhook records a webhook metric
func RecordWebhook(provider, eventType string) {
	webhooksReceived.WithLabelValues(provider, eventType).Inc()
}

// RecordWebhookError records a webhook processing error
func RecordWebhookError(provider, eventType string) {
	webhookProcessingErrors.WithLabelValues(provider, eventType).Inc()
}

// RecordWebhookDuration records webhook processing duration
func RecordWebhookDuration(provider, eventType string, duration time.Duration) {
	webhookProcessingDuration.WithLabelValues(provider, eventType).Observe(duration.Seconds())
}
