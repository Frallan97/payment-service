import type { Provider, Currency, PaymentStatus, PaginatedResponse } from "./common";

export interface Payment {
  id: string;
  customer_id: string;
  provider: Provider;
  provider_payment_id: string;
  amount: number;
  currency: Currency;
  status: PaymentStatus;
  payment_method_type?: string;
  payment_method_details?: Record<string, any>;
  description?: string;
  statement_descriptor?: string;
  subscription_id?: string;
  invoice_id?: string;
  client_secret?: string;
  failure_code?: string;
  failure_message?: string;
  metadata?: Record<string, any>;
  idempotency_key?: string;
  created_at: string;
  updated_at: string;
  completed_at?: string;
}

export interface CreatePaymentRequest {
  provider: Provider;
  amount: number;
  currency: Currency;
  description?: string;
  statement_descriptor?: string;
  metadata?: Record<string, any>;
}

export type PaymentListResponse = PaginatedResponse<Payment>;
