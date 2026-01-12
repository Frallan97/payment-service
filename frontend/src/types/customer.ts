export interface Customer {
  id: string;
  user_id: string;
  email: string;
  name: string;
  provider_customer_ids?: Record<string, string>;
  created_at: string;
  updated_at: string;
}
