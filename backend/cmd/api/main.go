package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"payment-service/internal/config"
	"payment-service/internal/database"
	"payment-service/internal/handlers"
	"payment-service/internal/middleware"
	"payment-service/internal/providers"
	"payment-service/internal/repository"
	"payment-service/internal/services"
	"payment-service/pkg/auth"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting payment service in %s mode on port %s", cfg.Env, cfg.Port)

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Connect to Redis
	redisOpt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	redisClient := redis.NewClient(redisOpt)
	defer redisClient.Close()

	// Test Redis connection
	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(pingCtx).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Println("Rate limiting will be disabled")
		redisClient = nil
	} else {
		log.Println("Successfully connected to Redis")
	}

	// Run migrations
	migrationsPath := "./migrations"
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		migrationsPath = "/app/migrations"
	}

	if err := database.RunMigrations(db, migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Fetch auth-service public key
	publicKey, err := auth.FetchPublicKey(cfg.AuthServiceURL)
	if err != nil {
		log.Fatalf("Failed to fetch auth-service public key: %v", err)
	}
	log.Println("Successfully fetched auth-service public key")

	// Initialize provider factory
	providerFactory := providers.NewFactory(cfg.StripeAPIKey, cfg.StripeWebhookSecret)

	// Initialize repositories
	customerRepo := repository.NewCustomerRepository(db.DB)
	paymentRepo := repository.NewPaymentRepository(db.DB)
	subscriptionRepo := repository.NewSubscriptionRepository(db.DB)
	refundRepo := repository.NewRefundRepository(db.DB)
	webhookRepo := repository.NewWebhookRepository(db.DB)

	// Initialize services
	paymentService := services.NewPaymentService(paymentRepo, customerRepo, providerFactory)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, customerRepo, providerFactory)
	refundService := services.NewRefundService(refundRepo, paymentRepo, customerRepo, providerFactory)
	// TODO: Wire webhookService into webhookHandler for async event processing
	_ = services.NewWebhookService(webhookRepo, paymentRepo, subscriptionRepo, refundRepo)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	customerHandler := handlers.NewCustomerHandler(customerRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)
	refundHandler := handlers.NewRefundHandler(refundService)
	webhookHandler := handlers.NewWebhookHandler(providerFactory)

	// Initialize router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.MetricsMiddleware)
	r.Use(middleware.CORSMiddleware(cfg.AllowedOrigins))

	// Rate limiting (if Redis is available)
	if redisClient != nil {
		rateLimiter := middleware.NewRateLimiter(redisClient, 100, time.Minute)
		r.Use(rateLimiter.RateLimitMiddleware)
	}

	// Health check endpoint (no auth required)
	r.Get("/health", healthHandler.Health)

	// Metrics endpoint (no auth required)
	r.Handle("/metrics", promhttp.Handler())

	// API routes (auth required)
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(publicKey))
		r.Use(middleware.IdempotencyMiddleware)

		// Customer endpoints
		r.Get("/customers/me", customerHandler.GetMe)

		// Payment endpoints
		r.Post("/payments", paymentHandler.CreatePayment)
		r.Get("/payments/{id}", paymentHandler.GetPayment)
		r.Get("/payments/{id}/refunds", refundHandler.ListRefundsByPayment)
		r.Get("/payments", paymentHandler.ListPayments)

		// Subscription endpoints
		r.Post("/subscriptions", subscriptionHandler.CreateSubscription)
		r.Get("/subscriptions/{id}", subscriptionHandler.GetSubscription)
		r.Patch("/subscriptions/{id}", subscriptionHandler.UpdateSubscription)
		r.Delete("/subscriptions/{id}", subscriptionHandler.CancelSubscription)
		r.Get("/subscriptions", subscriptionHandler.ListSubscriptions)

		// Refund endpoints
		r.Post("/refunds", refundHandler.CreateRefund)
		r.Get("/refunds/{id}", refundHandler.GetRefund)
		r.Get("/refunds", refundHandler.ListRefunds)
	})

	// Webhook endpoints (no auth, verified by signature)
	r.Post("/api/webhooks/stripe", webhookHandler.HandleStripeWebhook)
	r.Post("/api/webhooks/swish", webhookHandler.HandleSwishWebhook)

	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	if err := server.Close(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
