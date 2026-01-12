-- Drop triggers
DROP TRIGGER IF EXISTS update_refunds_updated_at ON refunds;
DROP TRIGGER IF EXISTS update_subscriptions_updated_at ON subscriptions;
DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;
DROP TRIGGER IF EXISTS update_customers_updated_at ON customers;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (in reverse order of creation due to foreign keys)
DROP TABLE IF EXISTS idempotency_keys;
DROP TABLE IF EXISTS webhook_events;
DROP TABLE IF EXISTS refunds;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS customers;

-- Drop enums
DROP TYPE IF EXISTS refund_status;
DROP TYPE IF EXISTS subscription_status;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS payment_provider;
DROP TYPE IF EXISTS currency_code;

-- Drop extension (only if no other tables are using it)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
