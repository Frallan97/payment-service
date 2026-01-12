export type Provider = "stripe" | "swish";

export type Currency = "SEK" | "USD" | "EUR" | "GBP";

export type PaymentStatus =
  | "pending"
  | "processing"
  | "requires_action"
  | "succeeded"
  | "failed"
  | "canceled";

export type SubscriptionStatus =
  | "active"
  | "past_due"
  | "unpaid"
  | "canceled"
  | "incomplete"
  | "incomplete_expired"
  | "trialing"
  | "paused";

export type RefundStatus =
  | "pending"
  | "processing"
  | "succeeded"
  | "failed"
  | "canceled";

export interface APIError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  limit: number;
  offset: number;
}
