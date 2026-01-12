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
import type { PaymentListResponse } from "@/types";

interface PaymentListProps {
  data: PaymentListResponse;
  isLoading: boolean;
  error: string | null;
  onPageChange: (offset: number) => void;
}

export default function PaymentList({
  data,
  isLoading,
  error,
  onPageChange,
}: PaymentListProps) {
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
        <p className="text-gray-500 mb-4">No payments found</p>
        <p className="text-sm text-gray-400">
          Create your first payment to get started
        </p>
      </div>
    );
  }

  const currentPage = Math.floor(data.offset / data.limit) + 1;
  const totalPages = Math.ceil(data.total / data.limit);

  const handleRowClick = (id: string) => {
    navigate(`/dashboard/payments/${id}`);
  };

  return (
    <div className="space-y-4">
      <div className="rounded-lg border bg-white shadow-sm">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Amount</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Provider</TableHead>
              <TableHead>Description</TableHead>
              <TableHead>Created</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.data.map((payment) => (
              <TableRow
                key={payment.id}
                className="cursor-pointer"
                onClick={() => handleRowClick(payment.id)}
              >
                <TableCell className="font-mono text-xs">
                  {payment.id.substring(0, 8)}...
                </TableCell>
                <TableCell className="font-semibold">
                  {formatCurrency(payment.amount, payment.currency)}
                </TableCell>
                <TableCell>
                  <StatusBadge status={payment.status} />
                </TableCell>
                <TableCell className="capitalize">{payment.provider}</TableCell>
                <TableCell className="max-w-xs truncate">
                  {payment.description || "-"}
                </TableCell>
                <TableCell className="text-sm text-gray-500">
                  {formatDate(payment.created_at)}
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
            {Math.min(data.offset + data.limit, data.total)} of {data.total} payments
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
