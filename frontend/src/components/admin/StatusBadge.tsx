import { Badge } from "@/components/ui/badge";
import type { PaymentStatus, SubscriptionStatus, RefundStatus } from "@/types";

interface StatusBadgeProps {
  status: PaymentStatus | SubscriptionStatus | RefundStatus | string;
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  const getVariant = (status: string): "success" | "destructive" | "warning" | "info" | "default" | "secondary" => {
    const lowerStatus = status.toLowerCase();

    // Success statuses (green)
    if (["succeeded", "active", "completed"].includes(lowerStatus)) {
      return "success";
    }

    // Error statuses (red)
    if (["failed", "canceled", "cancelled", "unpaid"].includes(lowerStatus)) {
      return "destructive";
    }

    // Warning statuses (yellow/orange)
    if (["pending", "processing", "past_due", "requires_action"].includes(lowerStatus)) {
      return "warning";
    }

    // Info statuses (blue)
    if (["trialing", "incomplete"].includes(lowerStatus)) {
      return "info";
    }

    // Paused/inactive (gray)
    if (["paused", "incomplete_expired"].includes(lowerStatus)) {
      return "secondary";
    }

    return "default";
  };

  return (
    <Badge variant={getVariant(status)}>
      {status.replace(/_/g, " ")}
    </Badge>
  );
}
