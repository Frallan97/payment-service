import { useState } from "react";
import { useRefunds } from "@/lib/hooks";
import RefundList from "@/components/admin/RefundList";
import CreateRefundForm from "@/components/admin/CreateRefundForm";

export default function RefundsPage() {
  const [offset, setOffset] = useState(0);
  const limit = 20;

  const { data, isLoading, error, refetch } = useRefunds({ limit, offset });

  const handlePageChange = (newOffset: number) => {
    setOffset(newOffset);
  };

  const handleRefundCreated = () => {
    // Refresh the list and go back to first page
    setOffset(0);
    refetch();
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Refunds</h1>
          <p className="text-gray-500 mt-1">
            View and manage all payment refunds
          </p>
        </div>
        <CreateRefundForm onSuccess={handleRefundCreated} />
      </div>

      <RefundList
        data={data!}
        isLoading={isLoading}
        error={error}
        onPageChange={handlePageChange}
      />
    </div>
  );
}
