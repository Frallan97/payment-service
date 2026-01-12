-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Currency enum type
CREATE TYPE currency_code AS ENUM ('SEK', 'USD', 'EUR', 'GBP');

-- Provider enum
CREATE TYPE payment_provider AS ENUM ('stripe', 'swish');

-- Payment status enum
CREATE TYPE payment_status AS ENUM (
    'pending',           -- Payment initiated
    'processing',        -- Being processed
    'requires_action',   -- Needs user action (3DS, etc)
    'succeeded',         -- Completed successfully
    'failed',            -- Failed
    'canceled'           -- Canceled by user/system
);

-- Subscription status enum
CREATE TYPE subscription_status AS ENUM (
    'active',
    'past_due',
    'unpaid',
    'canceled',
    'incomplete',
    'incomplete_expired',
    'trialing',
    'paused'
);

-- Refund status enum
CREATE TYPE refund_status AS ENUM (
    'pending',
    'processing',
    'succeeded',
    'failed',
    'canceled'
);

-- Customers table - maps auth-service users to provider customer IDs
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,                           -- From auth-service JWT
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255),

    -- Provider-specific customer IDs
    stripe_customer_id VARCHAR(255) UNIQUE,
    swish_customer_id VARCHAR(255) UNIQUE,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT unique_user_id UNIQUE (user_id)
);

-- Payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id),

    -- Payment details
    provider payment_provider NOT NULL,
    provider_payment_id VARCHAR(255) NOT NULL,       -- Stripe PaymentIntent ID or Swish payment ref
    amount BIGINT NOT NULL,                          -- Amount in smallest currency unit (öre, cents)
    currency currency_code NOT NULL DEFAULT 'SEK',
    status payment_status NOT NULL DEFAULT 'pending',

    -- Payment method
    payment_method_type VARCHAR(50),                 -- card, swish, bank_transfer, etc
    payment_method_details JSONB,                    -- Last4, brand, etc

    -- Destination/purpose
    description TEXT,
    statement_descriptor VARCHAR(255),               -- Appears on card statement

    -- Related entities
    subscription_id UUID,                            -- NULL for one-time payments
    invoice_id VARCHAR(255),                         -- Provider invoice ID (if applicable)

    -- Client secret for frontend confirmation (Stripe)
    client_secret VARCHAR(255),

    -- Error handling
    failure_code VARCHAR(100),
    failure_message TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',                     -- Custom app-specific data

    -- Idempotency
    idempotency_key VARCHAR(255),

    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,

    CONSTRAINT unique_provider_payment UNIQUE (provider, provider_payment_id)
);

-- Subscriptions table
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id),

    -- Subscription details
    provider payment_provider NOT NULL,
    provider_subscription_id VARCHAR(255) NOT NULL,
    status subscription_status NOT NULL DEFAULT 'incomplete',

    -- Billing
    amount BIGINT NOT NULL,
    currency currency_code NOT NULL DEFAULT 'SEK',
    interval VARCHAR(20) NOT NULL,                   -- day, week, month, year
    interval_count INTEGER NOT NULL DEFAULT 1,       -- Every X intervals

    -- Billing dates
    current_period_start TIMESTAMP NOT NULL,
    current_period_end TIMESTAMP NOT NULL,
    trial_start TIMESTAMP,
    trial_end TIMESTAMP,

    -- Cancellation
    cancel_at TIMESTAMP,
    canceled_at TIMESTAMP,
    cancel_at_period_end BOOLEAN DEFAULT false,

    -- Latest payment
    latest_payment_id UUID REFERENCES payments(id),

    -- Product info
    product_name VARCHAR(255),
    product_description TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT unique_provider_subscription UNIQUE (provider, provider_subscription_id)
);

-- Add foreign key from payments to subscriptions
ALTER TABLE payments
    ADD CONSTRAINT fk_payments_subscription
    FOREIGN KEY (subscription_id)
    REFERENCES subscriptions(id)
    ON DELETE SET NULL;

-- Refunds table
CREATE TABLE refunds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL REFERENCES payments(id),

    -- Refund details
    provider payment_provider NOT NULL,
    provider_refund_id VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL,                          -- Amount to refund
    currency currency_code NOT NULL,
    status refund_status NOT NULL DEFAULT 'pending',

    -- Reason
    reason VARCHAR(100),                             -- duplicate, fraudulent, requested_by_customer
    notes TEXT,

    -- Error handling
    failure_code VARCHAR(100),
    failure_message TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,

    CONSTRAINT unique_provider_refund UNIQUE (provider, provider_refund_id)
);

-- Webhook events table (for deduplication and audit)
CREATE TABLE webhook_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Webhook identification
    provider payment_provider NOT NULL,
    provider_event_id VARCHAR(255) NOT NULL,
    event_type VARCHAR(100) NOT NULL,

    -- Processing status
    processed BOOLEAN DEFAULT false,
    processing_attempts INTEGER DEFAULT 0,
    last_processing_error TEXT,

    -- Payload
    payload JSONB NOT NULL,

    -- Related entity (if identified)
    payment_id UUID REFERENCES payments(id),
    subscription_id UUID REFERENCES subscriptions(id),
    refund_id UUID REFERENCES refunds(id),

    -- Timestamps
    received_at TIMESTAMP DEFAULT NOW(),
    processed_at TIMESTAMP,

    CONSTRAINT unique_provider_event UNIQUE (provider, provider_event_id)
);

-- Idempotency keys table (for preventing duplicate operations)
CREATE TABLE idempotency_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    key VARCHAR(255) NOT NULL,
    request_method VARCHAR(10) NOT NULL,
    request_path VARCHAR(500) NOT NULL,

    -- Response cache
    response_status_code INTEGER,
    response_body JSONB,

    -- Associated resource
    resource_type VARCHAR(50),                       -- payment, subscription, refund
    resource_id UUID,

    -- Expiry (keys valid for 24 hours)
    expires_at TIMESTAMP NOT NULL,

    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT unique_idempotency_key UNIQUE (key)
);

-- Update triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_customers_updated_at
    BEFORE UPDATE ON customers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payments_updated_at
    BEFORE UPDATE ON payments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_subscriptions_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_refunds_updated_at
    BEFORE UPDATE ON refunds
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
