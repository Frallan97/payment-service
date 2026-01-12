# Payment Service

A centralized payment service providing unified payment processing for Stripe and Swish providers, supporting one-time payments, subscriptions, webhooks, and refunds.

## Features

- **Multiple Providers**: Stripe + Swish integration with unified API
- **Payment Types**: One-time payments, recurring subscriptions, refunds
- **Webhooks**: Real-time payment status updates with signature verification
- **Authentication**: JWT validation integrated with auth-service
- **Security**: PCI compliance, idempotency keys, rate limiting
- **Observability**: Structured logging, health checks

## Architecture

```
payment-service/
├── backend/
│   ├── cmd/api/main.go              # Application entry point
│   ├── internal/
│   │   ├── config/                  # Configuration management
│   │   ├── database/                # DB connection & migrations
│   │   ├── models/                  # Domain models
│   │   ├── providers/               # Payment provider abstraction
│   │   ├── services/                # Business logic (Phase 2)
│   │   ├── repository/              # Database operations (Phase 2)
│   │   ├── handlers/                # HTTP handlers (Phase 2)
│   │   └── middleware/              # HTTP middleware
│   ├── pkg/auth/                    # Auth-service integration
│   └── migrations/                  # SQL migrations
├── charts/                          # Helm chart (Phase 6)
└── docker-compose.yml               # Local development
```

## Quick Start

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Stripe API keys (test mode)
- SSH access to k3s cluster (for deployment)

### Deploy to Production

For production deployment to k3s cluster:

```bash
# 1. Use the deployment helper script
./scripts/deploy.sh

# 2. Or follow the detailed deployment guide
cat DEPLOYMENT.md
```

### Local Development

1. **Clone and navigate to the project:**
   ```bash
   cd payment-service
   ```

2. **Set up environment variables:**
   ```bash
   cp backend/.env.example backend/.env
   # Edit backend/.env with your Stripe API keys
   ```

3. **Start services with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

   This starts:
   - PostgreSQL on port 5434
   - Redis on port 6379
   - Payment service on port 8082

4. **Verify the service is running:**
   ```bash
   curl http://localhost:8082/health
   # Should return: OK
   ```

### Running Locally (Without Docker)

1. **Start PostgreSQL and Redis:**
   ```bash
   docker-compose up postgres redis -d
   ```

2. **Run the service:**
   ```bash
   cd backend
   go run cmd/api/main.go
   ```

## API Endpoints

### Authentication
All API endpoints require a JWT token from auth-service:
```
Authorization: Bearer <JWT_TOKEN>
```

### Payments
- `POST /api/payments` - Create a payment
- `GET /api/payments/:id` - Get payment details
- `GET /api/payments` - List payments

### Subscriptions
- `POST /api/subscriptions` - Create subscription
- `GET /api/subscriptions/:id` - Get subscription
- `PATCH /api/subscriptions/:id` - Update subscription
- `DELETE /api/subscriptions/:id` - Cancel subscription
- `GET /api/subscriptions` - List subscriptions

### Refunds
- `POST /api/refunds` - Create refund
- `GET /api/refunds/:id` - Get refund details
- `GET /api/refunds` - List refunds

### Webhooks (No Auth)
- `POST /api/webhooks/stripe` - Stripe webhook handler
- `POST /api/webhooks/swish` - Swish webhook handler

### Customer
- `GET /api/customers/me` - Get current user's customer record

## Example: Creating a Payment

```bash
curl -X POST http://localhost:8082/api/payments \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $(uuidgen)" \
  -d '{
    "provider": "stripe",
    "amount": 10000,
    "currency": "SEK",
    "description": "Premium subscription"
  }'
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "provider": "stripe",
  "provider_payment_id": "pi_3ABC123",
  "amount": 10000,
  "currency": "SEK",
  "status": "pending",
  "client_secret": "pi_3ABC123_secret_xyz",
  "created_at": "2026-01-11T10:00:00Z"
}
```

## Database Migrations

Migrations run automatically on service startup. To manually manage migrations:

```bash
# Apply migrations
migrate -path backend/migrations -database "postgresql://paymentuser:paymentpass@localhost:5434/paymentdb?sslmode=disable" up

# Rollback last migration
migrate -path backend/migrations -database "postgresql://paymentuser:paymentpass@localhost:5434/paymentdb?sslmode=disable" down 1
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| ENV | Environment (development/production) | development |
| DATABASE_URL | PostgreSQL connection string | Required |
| REDIS_URL | Redis connection string | redis://localhost:6379 |
| STRIPE_API_KEY | Stripe secret key | Required |
| STRIPE_WEBHOOK_SECRET | Stripe webhook secret | Required |
| SWISH_API_URL | Swish API URL | https://mss.cpc.getswish.net |
| SWISH_CERT_PATH | Path to Swish TLS certificate | - |
| SWISH_KEY_PATH | Path to Swish TLS key | - |
| AUTH_SERVICE_URL | Auth service URL | https://auth.vibeoholic.com |
| ALLOWED_ORIGINS | CORS allowed origins (comma-separated) | http://localhost:3000 |

## Development Roadmap

### ✅ Phase 1: Foundation (Completed)
- [x] Project structure and Go modules
- [x] Database schema and migrations
- [x] Configuration management
- [x] Auth-service JWT integration
- [x] Basic HTTP server with routes
- [x] Docker & Docker Compose setup

### ✅ Phase 2: Core Payments (Completed)
- [x] PaymentProvider interface
- [x] Stripe provider implementation
- [x] Customer & Payment repositories
- [x] Payment service with business logic
- [x] Payment handlers (create, get, list)
- [x] Customer endpoint (GET /api/customers/me)
- [x] Idempotency middleware
- [x] Complete payment flow integration

### ✅ Phase 3: Subscriptions (Completed)
- [x] Subscription models
- [x] Stripe subscription implementation
- [x] Subscription service & handlers
- [x] Subscription repository with CRUD operations
- [x] Create, get, update, cancel, list subscriptions
- [x] Cancel at period end or immediate cancellation support

### ✅ Phase 4: Refunds & Webhooks (Completed)
- [x] Refund implementation in Stripe provider
- [x] Refund repository and service
- [x] Refund handlers (create, get, list)
- [x] Webhook handler for Stripe with signature verification
- [x] Webhook service for processing events
- [x] Webhook repository with event deduplication
- [x] Payment/subscription/refund status updates via webhooks

### Phase 5: Swish Integration
- [ ] Swish provider with TLS
- [ ] Swish payment flow
- [ ] Swish webhooks

### ✅ Phase 6: Production Readiness (Completed)
- [x] Unit tests for payment service with testify/mock
- [x] Repository interface abstraction for better testability
- [x] Prometheus metrics middleware
- [x] Custom payment/subscription/refund/webhook metrics
- [x] HTTP request duration and count metrics
- [x] Redis-based rate limiting (sliding window algorithm)
- [x] Rate limit headers (X-RateLimit-*)
- [x] Helm chart with deployment, service, ingress
- [x] GitHub Actions CI/CD pipeline
- [x] Automated testing, building, and pushing to GHCR
- [x] Health and readiness probes
- [x] /metrics endpoint for Prometheus scraping

### 🚀 Phase 7: Deployment (Ready)
- [x] Deployment guide created (`DEPLOYMENT.md`)
- [x] Kubernetes manifests (secrets, ArgoCD application)
- [x] Deployment helper script (`scripts/deploy.sh`)
- [x] Grafana dashboard configuration
- [ ] Deploy to k3s cluster at payments.vibeoholic.com
- [ ] Configure Stripe webhooks
- [ ] Production testing and validation

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v -run TestPaymentService ./internal/services
```

## Project Status

**Current Phase**: Phase 6 Complete ✅

The payment service is now production-ready with comprehensive functionality:
- ✅ Complete Stripe integration (payments, subscriptions, refunds)
- ✅ Customer management with provider mapping
- ✅ Payment, subscription, and refund full lifecycle
- ✅ Webhook handlers with signature verification and event processing
- ✅ **Unit tests** with comprehensive mocking
- ✅ **Prometheus metrics** for monitoring
- ✅ **Redis-based rate limiting** (100 req/min per user)
- ✅ **Helm chart** for Kubernetes deployment
- ✅ **GitHub Actions CI/CD** pipeline
- ✅ Database repositories with interface abstraction
- ✅ Service layer with business logic
- ✅ HTTP handlers with JSON responses
- ✅ Authentication via auth-service JWT
- ✅ Idempotency support
- ✅ Health and metrics endpoints
- ✅ Production-grade error handling and logging

**Next Steps**: Phase 7 - Deployment to k3s cluster via ArgoCD

## Contributing

1. Follow the existing code structure and patterns
2. Write tests for new features
3. Run `go fmt` before committing
4. Ensure migrations have both up and down files

## Architecture Decisions

- **Provider Abstraction**: Unified interface allows easy addition of new payment providers
- **JWT Validation**: Integrates with centralized auth-service for consistent authentication
- **Database Enums**: PostgreSQL enums ensure data consistency
- **Idempotency**: Prevents duplicate charges from network retries
- **Webhooks**: Asynchronous processing with retry logic for reliability

## License

Internal project - All rights reserved
