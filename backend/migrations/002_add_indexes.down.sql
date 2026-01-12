-- Drop idempotency indexes
DROP INDEX IF EXISTS idx_idempotency_expires;
DROP INDEX IF EXISTS idx_idempotency_key;

-- Drop webhook indexes
DROP INDEX IF EXISTS idx_webhooks_event_type;
DROP INDEX IF EXISTS idx_webhooks_subscription_id;
DROP INDEX IF EXISTS idx_webhooks_payment_id;
DROP INDEX IF EXISTS idx_webhooks_processed;
DROP INDEX IF EXISTS idx_webhooks_provider_event;

-- Drop refund indexes
DROP INDEX IF EXISTS idx_refunds_created_at;
DROP INDEX IF EXISTS idx_refunds_provider;
DROP INDEX IF EXISTS idx_refunds_status;
DROP INDEX IF EXISTS idx_refunds_payment_id;

-- Drop subscription indexes
DROP INDEX IF EXISTS idx_subscriptions_trial_end;
DROP INDEX IF EXISTS idx_subscriptions_next_billing;
DROP INDEX IF EXISTS idx_subscriptions_provider_sub_id;
DROP INDEX IF EXISTS idx_subscriptions_provider;
DROP INDEX IF EXISTS idx_subscriptions_status;
DROP INDEX IF EXISTS idx_subscriptions_customer_id;

-- Drop payment indexes
DROP INDEX IF EXISTS idx_payments_customer_status;
DROP INDEX IF EXISTS idx_payments_created_at;
DROP INDEX IF EXISTS idx_idempotency_key;
DROP INDEX IF EXISTS idx_payments_provider_payment_id;
DROP INDEX IF EXISTS idx_payments_provider;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_subscription_id;
DROP INDEX IF EXISTS idx_payments_customer_id;

-- Drop customer indexes
DROP INDEX IF EXISTS idx_customers_deleted_at;
DROP INDEX IF EXISTS idx_customers_swish_id;
DROP INDEX IF EXISTS idx_customers_stripe_id;
DROP INDEX IF EXISTS idx_customers_email;
DROP INDEX IF EXISTS idx_customers_user_id;
