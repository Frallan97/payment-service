import { useNavigate } from "react-router-dom";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { Loading } from "@/components/ui/loading";
import { Alert, AlertDescription } from "@/components/ui/alert";
import StatusBadge from "./StatusBadge";
import { formatCurrency, formatDate } from "@/lib/utils";
import type { SubscriptionListResponse } from "@/types";

interface SubscriptionListProps {
  data: SubscriptionListResponse;
  isLoading: boolean;
  error: string | null;
  onPageChange: (offset: number) => void;
}

export default function SubscriptionList({
  data,
  isLoading,
  error,
  onPageChange,
}: SubscriptionListProps) {
  const navigate = useNavigate();

  if (isLoading) {
    return (
      <div className="py-12">
        <Loading size="lg" />
      </div>
    );
  }

  if (error) {
    return (
      <Alert variant="destructive">
        <AlertDescription>{error}</AlertDescription>
      </Alert>
    );
  }

  if (!data || data.data.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500 mb-4">No subscriptions found</p>
        <p className="text-sm text-gray-400">
          Create your first subscription to get started
        </p>
      </div>
    );
  }

  const currentPage = Math.floor(data.offset / data.limit) + 1;
  const totalPages = Math.ceil(data.total / data.limit);

  const handleRowClick = (id: string) => {
    navigate(`/dashboard/subscriptions/${id}`);
  };

  return (
    <div className="space-y-4">
      <div className="rounded-lg border bg-white shadow-sm">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Product</TableHead>
              <TableHead>Amount</TableHead>
              <TableHead>Interval</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Current Period End</TableHead>
              <TableHead>Created</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.data.map((subscription) => (
              <TableRow
                key={subscription.id}
                className="cursor-pointer"
                onClick={() => handleRowClick(subscription.id)}
              >
                <TableCell className="font-mono text-xs">
                  {subscription.id.substring(0, 8)}...
                </TableCell>
                <TableCell className="font-medium">
                  {subscription.product_name}
                </TableCell>
                <TableCell className="font-semibold">
                  {formatCurrency(subscription.amount, subscription.currency)}
                </TableCell>
                <TableCell>
                  Every {subscription.interval_count > 1 ? subscription.interval_count : ""}{" "}
                  {subscription.interval}
                  {subscription.interval_count > 1 ? "s" : ""}
                </TableCell>
                <TableCell>
                  <StatusBadge status={subscription.status} />
                </TableCell>
                <TableCell className="text-sm text-gray-500">
                  {formatDate(subscription.current_period_end)}
                </TableCell>
                <TableCell className="text-sm text-gray-500">
                  {formatDate(subscription.created_at)}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-gray-500">
            Showing {data.offset + 1} to{" "}
            {Math.min(data.offset + data.limit, data.total)} of {data.total} subscriptions
          </p>
          <Pagination>
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  onClick={() => {
                    if (currentPage > 1) {
                      onPageChange((currentPage - 2) * data.limit);
                    }
                  }}
                  className={
                    currentPage === 1
                      ? "pointer-events-none opacity-50"
                      : "cursor-pointer"
                  }
                />
              </PaginationItem>

              {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                const page = i + 1;
                return (
                  <PaginationItem key={page}>
                    <PaginationLink
                      onClick={() => onPageChange((page - 1) * data.limit)}
                      isActive={page === currentPage}
                      className="cursor-pointer"
                    >
                      {page}
                    </PaginationLink>
                  </PaginationItem>
                );
              })}

              <PaginationItem>
                <PaginationNext
                  onClick={() => {
                    if (currentPage < totalPages) {
                      onPageChange(currentPage * data.limit);
                    }
                  }}
                  className={
                    currentPage === totalPages
                      ? "pointer-events-none opacity-50"
                      : "cursor-pointer"
                  }
                />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        </div>
      )}
    </div>
  );
}
