-- Customer indexes
CREATE INDEX idx_customers_user_id ON customers(user_id);
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_stripe_id ON customers(stripe_customer_id);
CREATE INDEX idx_customers_swish_id ON customers(swish_customer_id);
CREATE INDEX idx_customers_deleted_at ON customers(deleted_at) WHERE deleted_at IS NOT NULL;

-- Payment indexes
CREATE INDEX idx_payments_customer_id ON payments(customer_id);
CREATE INDEX idx_payments_subscription_id ON payments(subscription_id) WHERE subscription_id IS NOT NULL;
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_provider ON payments(provider);
CREATE INDEX idx_payments_provider_payment_id ON payments(provider_payment_id);
CREATE INDEX idx_payments_idempotency_key ON payments(idempotency_key) WHERE idempotency_key IS NOT NULL;
CREATE INDEX idx_payments_created_at ON payments(created_at DESC);
CREATE INDEX idx_payments_customer_status ON payments(customer_id, status);

-- Subscription indexes
CREATE INDEX idx_subscriptions_customer_id ON subscriptions(customer_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_subscriptions_provider ON subscriptions(provider);
CREATE INDEX idx_subscriptions_provider_sub_id ON subscriptions(provider_subscription_id);
CREATE INDEX idx_subscriptions_next_billing ON subscriptions(current_period_end) WHERE status = 'active';
CREATE INDEX idx_subscriptions_trial_end ON subscriptions(trial_end) WHERE status = 'trialing';

-- Refund indexes
CREATE INDEX idx_refunds_payment_id ON refunds(payment_id);
CREATE INDEX idx_refunds_status ON refunds(status);
CREATE INDEX idx_refunds_provider ON refunds(provider);
CREATE INDEX idx_refunds_created_at ON refunds(created_at DESC);

-- Webhook indexes
CREATE INDEX idx_webhooks_provider_event ON webhook_events(provider, provider_event_id);
CREATE INDEX idx_webhooks_processed ON webhook_events(processed, received_at);
CREATE INDEX idx_webhooks_payment_id ON webhook_events(payment_id) WHERE payment_id IS NOT NULL;
CREATE INDEX idx_webhooks_subscription_id ON webhook_events(subscription_id) WHERE subscription_id IS NOT NULL;
CREATE INDEX idx_webhooks_event_type ON webhook_events(event_type);

-- Idempotency indexes
CREATE INDEX idx_idempotency_key ON idempotency_keys(key);
CREATE INDEX idx_idempotency_expires ON idempotency_keys(expires_at);
