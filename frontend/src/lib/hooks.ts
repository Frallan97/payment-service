import { useState, useEffect, useCallback } from "react";
import { PaymentServiceAPI } from "./api";
import type {
  PaymentListResponse,
  Payment,
  SubscriptionListResponse,
  Subscription,
  RefundListResponse,
  Refund,
} from "@/types";

interface ListParams {
  limit?: number;
  offset?: number;
}

interface UseListResult<T> {
  data: T | null;
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

interface UseDetailResult<T> {
  data: T | null;
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

// Payment hooks
export function usePayments(params?: ListParams): UseListResult<PaymentListResponse> {
  const [data, setData] = useState<PaymentListResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchPayments = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const result = await PaymentServiceAPI.listPayments(params);
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load payments");
    } finally {
      setIsLoading(false);
    }
  }, [params?.limit, params?.offset]);

  useEffect(() => {
    fetchPayments();
  }, [fetchPayments]);

  return { data, isLoading, error, refetch: fetchPayments };
}

export function usePayment(id: string): UseDetailResult<Payment> {
  const [data, setData] = useState<Payment | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchPayment = useCallback(async () => {
    if (!id) return;

    setIsLoading(true);
    setError(null);
    try {
      const result = await PaymentServiceAPI.getPayment(id);
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load payment");
    } finally {
      setIsLoading(false);
    }
  }, [id]);

  useEffect(() => {
    fetchPayment();
  }, [fetchPayment]);

  return { data, isLoading, error, refetch: fetchPayment };
}

// Subscription hooks
export function useSubscriptions(params?: ListParams): UseListResult<SubscriptionListResponse> {
  const [data, setData] = useState<SubscriptionListResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchSubscriptions = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const result = await PaymentServiceAPI.listSubscriptions(params);
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load subscriptions");
    } finally {
      setIsLoading(false);
    }
  }, [params?.limit, params?.offset]);

  useEffect(() => {
    fetchSubscriptions();
  }, [fetchSubscriptions]);

  return { data, isLoading, error, refetch: fetchSubscriptions };
}

export function useSubscription(id: string): UseDetailResult<Subscription> {
  const [data, setData] = useState<Subscription | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchSubscription = useCallback(async () => {
    if (!id) return;

    setIsLoading(true);
    setError(null);
    try {
      const result = await PaymentServiceAPI.getSubscription(id);
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load subscription");
    } finally {
      setIsLoading(false);
    }
  }, [id]);

  useEffect(() => {
    fetchSubscription();
  }, [fetchSubscription]);

  return { data, isLoading, error, refetch: fetchSubscription };
}

// Refund hooks
export function useRefunds(params?: ListParams): UseListResult<RefundListResponse> {
  const [data, setData] = useState<RefundListResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchRefunds = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const result = await PaymentServiceAPI.listRefunds(params);
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load refunds");
    } finally {
      setIsLoading(false);
    }
  }, [params?.limit, params?.offset]);

  useEffect(() => {
    fetchRefunds();
  }, [fetchRefunds]);

  return { data, isLoading, error, refetch: fetchRefunds };
}

export function useRefund(id: string): UseDetailResult<Refund> {
  const [data, setData] = useState<Refund | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchRefund = useCallback(async () => {
    if (!id) return;

    setIsLoading(true);
    setError(null);
    try {
      const result = await PaymentServiceAPI.getRefund(id);
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load refund");
    } finally {
      setIsLoading(false);
    }
  }, [id]);

  useEffect(() => {
    fetchRefund();
  }, [fetchRefund]);

  return { data, isLoading, error, refetch: fetchRefund };
}
