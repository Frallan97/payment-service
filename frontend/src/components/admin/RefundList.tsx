import { Link } from "react-router-dom";
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
import type { RefundListResponse } from "@/types";

interface RefundListProps {
  data: RefundListResponse;
  isLoading: boolean;
  error: string | null;
  onPageChange: (offset: number) => void;
}

export default function RefundList({
  data,
  isLoading,
  error,
  onPageChange,
}: RefundListProps) {
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
        <p className="text-gray-500 mb-4">No refunds found</p>
        <p className="text-sm text-gray-400">
          Refunds will appear here when you process them
        </p>
      </div>
    );
  }

  const currentPage = Math.floor(data.offset / data.limit) + 1;
  const totalPages = Math.ceil(data.total / data.limit);

  return (
    <div className="space-y-4">
      <div className="rounded-lg border bg-white shadow-sm">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Payment</TableHead>
              <TableHead>Amount</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Reason</TableHead>
              <TableHead>Provider</TableHead>
              <TableHead>Created</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.data.map((refund) => (
              <TableRow key={refund.id}>
                <TableCell className="font-mono text-xs">
                  {refund.id.substring(0, 8)}...
                </TableCell>
                <TableCell>
                  <Link
                    to={`/dashboard/payments/${refund.payment_id}`}
                    className="text-blue-600 hover:text-blue-700 font-mono text-xs"
                  >
                    {refund.payment_id.substring(0, 8)}...
                  </Link>
                </TableCell>
                <TableCell className="font-semibold">
                  {formatCurrency(refund.amount, refund.currency)}
                </TableCell>
                <TableCell>
                  <StatusBadge status={refund.status} />
                </TableCell>
                <TableCell className="max-w-xs truncate">
                  {refund.reason || "-"}
                </TableCell>
                <TableCell className="capitalize">{refund.provider}</TableCell>
                <TableCell className="text-sm text-gray-500">
                  {formatDate(refund.created_at)}
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
            {Math.min(data.offset + data.limit, data.total)} of {data.total} refunds
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
