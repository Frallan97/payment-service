import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { CreditCard } from "lucide-react";
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
import { SearchInput } from "@/components/ui/search-input";
import {
  TableToolbar,
  TableToolbarLeft,
  TableToolbarRight,
} from "@/components/ui/table-toolbar";
import {
  TooltipProvider,
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from "@/components/ui/tooltip";
import { EmptyState } from "@/components/ui/empty-state";
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
  const [searchQuery, setSearchQuery] = useState("");

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
      <EmptyState
        icon={CreditCard}
        title="No payments found"
        description="Create your first payment to get started"
      />
    );
  }

  // Filter payments based on search query
  const filteredPayments = data.data.filter((payment) => {
    if (!searchQuery) return true;
    const query = searchQuery.toLowerCase();
    return (
      payment.id.toLowerCase().includes(query) ||
      payment.description?.toLowerCase().includes(query) ||
      payment.status.toLowerCase().includes(query) ||
      payment.provider.toLowerCase().includes(query)
    );
  });

  const currentPage = Math.floor(data.offset / data.limit) + 1;
  const totalPages = Math.ceil(data.total / data.limit);

  const handleRowClick = (id: string) => {
    navigate(`/dashboard/payments/${id}`);
  };

  return (
    <div className="space-y-4">
      {/* Toolbar with search */}
      <TableToolbar>
        <TableToolbarLeft>
          <SearchInput
            placeholder="Search by ID, description, status, or provider..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onClear={() => setSearchQuery("")}
            className="max-w-md"
          />
        </TableToolbarLeft>
        <TableToolbarRight>
          <span className="text-sm text-neutral-500">
            {filteredPayments.length} of {data.total} payments
          </span>
        </TableToolbarRight>
      </TableToolbar>

      {/* Table wrapper with mobile responsive */}
      <div className="rounded-xl border border-neutral-200 bg-white shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
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
              {filteredPayments.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6}>
                    <EmptyState
                      icon={CreditCard}
                      title="No matching payments"
                      description="Try adjusting your search query"
                    />
                  </TableCell>
                </TableRow>
              ) : (
                filteredPayments.map((payment) => (
                  <TableRow
                    key={payment.id}
                    className="cursor-pointer"
                    onClick={() => handleRowClick(payment.id)}
                  >
                    <TableCell className="font-mono text-xs">
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <span className="text-neutral-600">
                              {payment.id.substring(0, 8)}...
                            </span>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p className="font-mono text-xs">{payment.id}</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                    </TableCell>
                    <TableCell className="font-semibold text-neutral-800">
                      {formatCurrency(payment.amount, payment.currency)}
                    </TableCell>
                    <TableCell>
                      <StatusBadge status={payment.status} />
                    </TableCell>
                    <TableCell className="capitalize text-neutral-700">
                      {payment.provider}
                    </TableCell>
                    <TableCell className="max-w-xs truncate text-neutral-700">
                      {payment.description || "-"}
                    </TableCell>
                    <TableCell className="text-sm text-neutral-500">
                      {formatDate(payment.created_at)}
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-neutral-500">
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
