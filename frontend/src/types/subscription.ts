import type { Provider, Currency, SubscriptionStatus, PaginatedResponse } from "./common";

export interface Subscription {
  id: string;
  customer_id: string;
  provider: Provider;
  provider_subscription_id: string;
  status: SubscriptionStatus;
  amount: number;
  currency: Currency;
  interval: string;
  interval_count: number;
  current_period_start: string;
  current_period_end: string;
  trial_start?: string;
  trial_end?: string;
  cancel_at?: string;
  canceled_at?: string;
  cancel_at_period_end: boolean;
  latest_payment_id?: string;
  product_name: string;
  product_description?: string;
  metadata?: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface CreateSubscriptionRequest {
  provider: Provider;
  amount: number;
  currency: Currency;
  interval: string;
  interval_count: number;
  product_name: string;
  product_description?: string;
  trial_period_days?: number;
  metadata?: Record<string, any>;
}

export interface UpdateSubscriptionRequest {
  cancel_at_period_end?: boolean;
  metadata?: Record<string, any>;
}

export type SubscriptionListResponse = PaginatedResponse<Subscription>;
