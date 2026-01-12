import type { Provider, Currency, RefundStatus, PaginatedResponse } from "./common";

export interface Refund {
  id: string;
  payment_id: string;
  provider: Provider;
  provider_refund_id: string;
  amount: number;
  currency: Currency;
  status: RefundStatus;
  reason?: string;
  notes?: string;
  failure_code?: string;
  failure_message?: string;
  metadata?: Record<string, any>;
  created_at: string;
  updated_at: string;
  completed_at?: string;
}

export interface CreateRefundRequest {
  payment_id: string;
  amount: number;
  reason?: string;
  notes?: string;
  metadata?: Record<string, any>;
}

export type RefundListResponse = PaginatedResponse<Refund>;
