import { useState } from "react";
import { useSubscriptions } from "@/lib/hooks";
import SubscriptionList from "@/components/admin/SubscriptionList";
import CreateSubscriptionForm from "@/components/admin/CreateSubscriptionForm";

export default function SubscriptionsPage() {
  const [offset, setOffset] = useState(0);
  const limit = 20;

  const { data, isLoading, error, refetch } = useSubscriptions({ limit, offset });

  const handlePageChange = (newOffset: number) => {
    setOffset(newOffset);
  };

  const handleSubscriptionCreated = () => {
    // Refresh the list and go back to first page
    setOffset(0);
    refetch();
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Subscriptions</h1>
          <p className="text-gray-500 mt-1">
            View and manage all recurring subscriptions
          </p>
        </div>
        <CreateSubscriptionForm onSuccess={handleSubscriptionCreated} />
      </div>

      <SubscriptionList
        data={data!}
        isLoading={isLoading}
        error={error}
        onPageChange={handlePageChange}
      />
    </div>
  );
}
