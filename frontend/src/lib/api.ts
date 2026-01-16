import { AuthService } from './auth';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface Payment {
  id: string;
  amount: number;
  currency: string;
  status: string;
  customer_email?: string;
  created_at: string;
  provider: string;
}

export interface Subscription {
  id: string;
  customer_id: string;
  status: string;
  plan_name?: string;
  amount: number;
  currency: string;
  current_period_start: string;
  current_period_end: string;
}

export interface Refund {
  id: string;
  payment_id: string;
  amount: number;
  currency: string;
  status: string;
  reason?: string;
  created_at: string;
}

async function fetchAPI(endpoint: string, options: RequestInit = {}) {
  const headers = {
    'Content-Type': 'application/json',
    ...AuthService.getAuthHeader(),
    ...options.headers,
  };

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    throw new Error(`API error: ${response.statusText}`);
  }

  return response.json();
}

export const api = {
  // Payments
  async getPayments(): Promise<Payment[]> {
    return fetchAPI('/api/payments');
  },

  async getPayment(id: string): Promise<Payment> {
    return fetchAPI(`/api/payments/${id}`);
  },

  // Subscriptions
  async getSubscriptions(): Promise<Subscription[]> {
    return fetchAPI('/api/subscriptions');
  },

  async getSubscription(id: string): Promise<Subscription> {
    return fetchAPI(`/api/subscriptions/${id}`);
  },

  // Refunds
  async getRefunds(): Promise<Refund[]> {
    return fetchAPI('/api/refunds');
  },

  async createRefund(paymentId: string, amount: number, reason?: string): Promise<Refund> {
    return fetchAPI('/api/refunds', {
      method: 'POST',
      body: JSON.stringify({ payment_id: paymentId, amount, reason }),
    });
  },
};
