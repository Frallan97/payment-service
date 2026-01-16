import { AuthService, PAYMENT_SERVICE_URL } from "./auth";
import type {
  Payment,
  CreatePaymentRequest,
  PaymentListResponse,
  Subscription,
  CreateSubscriptionRequest,
  UpdateSubscriptionRequest,
  SubscriptionListResponse,
  Refund,
  CreateRefundRequest,
  RefundListResponse,
  Customer,
  APIError,
} from "@/types";

export class APIClientError extends Error {
  code: string;
  statusCode: number;

  constructor(
    message: string,
    code: string,
    statusCode: number
  ) {
    super(message);
    this.name = "APIClientError";
    this.code = code;
    this.statusCode = statusCode;
  }
}

interface ListParams {
  limit?: number;
  offset?: number;
}

export class PaymentServiceAPI {
  private static baseURL = PAYMENT_SERVICE_URL;

  private static async fetch<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    try {
      const response = await AuthService.fetchWithAuth(
        `${this.baseURL}${endpoint}`,
        options
      );

      if (!response.ok) {
        let error: APIError;
        try {
          error = await response.json();
        } catch {
          // Handle non-JSON error responses (e.g., HTML error pages)
          error = {
            code: 'unknown_error',
            message: `HTTP ${response.status} error`,
          };
        }
        throw new APIClientError(error.message, error.code, response.status);
      }

      return await response.json();
    } catch (error) {
      if (error instanceof APIClientError) {
        throw error;
      }
      // Handle network errors and other unexpected errors
      throw new APIClientError(
        error instanceof Error ? error.message : "An unexpected error occurred",
        "network_error",
        0
      );
    }
  }

  // Payment methods
  static async listPayments(params?: ListParams): Promise<PaymentListResponse> {
    const query = new URLSearchParams({
      limit: params?.limit?.toString() || "20",
      offset: params?.offset?.toString() || "0",
    });
    return this.fetch<PaymentListResponse>(`/api/payments?${query}`);
  }

  static async getPayment(id: string): Promise<Payment> {
    return this.fetch<Payment>(`/api/payments/${id}`);
  }

  static async createPayment(data: CreatePaymentRequest): Promise<Payment> {
    return this.fetch<Payment>(`/api/payments`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
  }

  // Subscription methods
  static async listSubscriptions(params?: ListParams): Promise<SubscriptionListResponse> {
    const query = new URLSearchParams({
      limit: params?.limit?.toString() || "20",
      offset: params?.offset?.toString() || "0",
    });
    return this.fetch<SubscriptionListResponse>(`/api/subscriptions?${query}`);
  }

  static async getSubscription(id: string): Promise<Subscription> {
    return this.fetch<Subscription>(`/api/subscriptions/${id}`);
  }

  static async createSubscription(data: CreateSubscriptionRequest): Promise<Subscription> {
    return this.fetch<Subscription>(`/api/subscriptions`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
  }

  static async updateSubscription(
    id: string,
    data: UpdateSubscriptionRequest
  ): Promise<Subscription> {
    return this.fetch<Subscription>(`/api/subscriptions/${id}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
  }

  static async cancelSubscription(
    id: string,
    immediate: boolean = false
  ): Promise<Subscription> {
    const query = immediate ? "?immediate=true" : "";
    return this.fetch<Subscription>(`/api/subscriptions/${id}${query}`, {
      method: "DELETE",
    });
  }

  // Refund methods
  static async listRefunds(params?: ListParams): Promise<RefundListResponse> {
    const query = new URLSearchParams({
      limit: params?.limit?.toString() || "20",
      offset: params?.offset?.toString() || "0",
    });
    return this.fetch<RefundListResponse>(`/api/refunds?${query}`);
  }

  static async getRefund(id: string): Promise<Refund> {
    return this.fetch<Refund>(`/api/refunds/${id}`);
  }

  static async createRefund(data: CreateRefundRequest): Promise<Refund> {
    return this.fetch<Refund>(`/api/refunds`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
  }

  static async listRefundsByPayment(paymentId: string): Promise<RefundListResponse> {
    return this.fetch<RefundListResponse>(`/api/payments/${paymentId}/refunds`);
  }

  // Customer methods
  static async getMe(): Promise<Customer> {
    return this.fetch<Customer>(`/api/customers/me`);
  }
}
