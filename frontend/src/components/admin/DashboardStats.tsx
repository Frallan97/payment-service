import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Loading } from "@/components/ui/loading";
import { CreditCard, RefreshCw, RotateCcw, TrendingUp } from "lucide-react";
import { PaymentServiceAPI } from "@/lib/api";
import { formatCurrency } from "@/lib/utils";

interface Stats {
  totalPayments: number;
  succeededPayments: number;
  totalPaymentAmount: number;
  activeSubscriptions: number;
  totalRefunds: number;
}

export default function DashboardStats() {
  const [stats, setStats] = useState<Stats | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      // Fetch data from all endpoints
      const [payments, subscriptions, refunds] = await Promise.all([
        PaymentServiceAPI.listPayments({ limit: 100, offset: 0 }),
        PaymentServiceAPI.listSubscriptions({ limit: 100, offset: 0 }),
        PaymentServiceAPI.listRefunds({ limit: 100, offset: 0 }),
      ]);

      // Calculate stats
      const succeededPayments = payments.data.filter((p) => p.status === "succeeded");
      const totalPaymentAmount = succeededPayments.reduce((sum, p) => sum + p.amount, 0);
      const activeSubscriptions = subscriptions.data.filter(
        (s) => s.status === "active" || s.status === "trialing"
      );

      setStats({
        totalPayments: payments.total,
        succeededPayments: succeededPayments.length,
        totalPaymentAmount,
        activeSubscriptions: activeSubscriptions.length,
        totalRefunds: refunds.total,
      });
    } catch (error) {
      console.error("Failed to load stats:", error);
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) {
    return (
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        {[...Array(4)].map((_, i) => (
          <Card key={i}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Loading...</CardTitle>
            </CardHeader>
            <CardContent>
              <Loading size="sm" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (!stats) {
    return null;
  }

  const successRate = stats.totalPayments > 0
    ? ((stats.succeededPayments / stats.totalPayments) * 100).toFixed(1)
    : "0";

  return (
    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
      {/* Total Payments */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Total Payments</CardTitle>
          <CreditCard className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{stats.totalPayments}</div>
          <p className="text-xs text-muted-foreground">
            {stats.succeededPayments} succeeded
          </p>
        </CardContent>
      </Card>

      {/* Payment Volume */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Payment Volume</CardTitle>
          <TrendingUp className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {formatCurrency(stats.totalPaymentAmount, "SEK")}
          </div>
          <p className="text-xs text-muted-foreground">
            From succeeded payments
          </p>
        </CardContent>
      </Card>

      {/* Active Subscriptions */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Active Subscriptions</CardTitle>
          <RefreshCw className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{stats.activeSubscriptions}</div>
          <p className="text-xs text-muted-foreground">
            Recurring revenue streams
          </p>
        </CardContent>
      </Card>

      {/* Total Refunds */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Refunds</CardTitle>
          <RotateCcw className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{stats.totalRefunds}</div>
          <p className="text-xs text-muted-foreground">
            Success rate: {successRate}%
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
