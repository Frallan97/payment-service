import { useState } from "react";
import { usePayments } from "@/lib/hooks";
import PaymentList from "@/components/admin/PaymentList";
import CreatePaymentForm from "@/components/admin/CreatePaymentForm";

export default function PaymentsPage() {
  const [offset, setOffset] = useState(0);
  const limit = 20;

  const { data, isLoading, error, refetch } = usePayments({ limit, offset });

  const handlePageChange = (newOffset: number) => {
    setOffset(newOffset);
  };

  const handlePaymentCreated = () => {
    // Refresh the list and go back to first page
    setOffset(0);
    refetch();
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Payments</h1>
          <p className="text-gray-500 mt-1">
            View and manage all payment transactions
          </p>
        </div>
        <CreatePaymentForm onSuccess={handlePaymentCreated} />
      </div>

      <PaymentList
        data={data!}
        isLoading={isLoading}
        error={error}
        onPageChange={handlePageChange}
      />
    </div>
  );
}
