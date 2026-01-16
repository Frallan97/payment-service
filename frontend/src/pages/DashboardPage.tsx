import { Link } from "react-router-dom";
import DashboardStats from "@/components/admin/DashboardStats";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ArrowRight, CreditCard, RefreshCw, RotateCcw } from "lucide-react";

export default function DashboardPage() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-500 mt-1">
          Overview of your payment system
        </p>
      </div>

      {/* Stats Cards */}
      <DashboardStats />

      {/* Quick Actions */}
      <div>
        <h2 className="text-xl font-semibold text-gray-900 mb-4">Quick Actions</h2>
        <div className="grid gap-4 md:grid-cols-3">
          <Card className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-lg">
                <CreditCard className="h-5 w-5" />
                Payments
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-500 mb-4">
                View and manage all payment transactions
              </p>
              <Link to="/dashboard/payments">
                <Button variant="outline">
                  View Payments
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </Link>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-lg">
                <RefreshCw className="h-5 w-5" />
                Subscriptions
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-500 mb-4">
                Manage recurring subscription billing
              </p>
              <Link to="/dashboard/subscriptions">
                <Button variant="outline">
                  View Subscriptions
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </Link>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-lg">
                <RotateCcw className="h-5 w-5" />
                Refunds
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-500 mb-4">
                Process and track payment refunds
              </p>
              <Link to="/dashboard/refunds">
                <Button variant="outline">
                  View Refunds
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </Link>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Getting Started */}
      <Card>
        <CardHeader>
          <CardTitle>Getting Started</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="flex items-start gap-3">
            <div className="flex h-6 w-6 items-center justify-center rounded-full bg-blue-100 text-sm font-medium text-blue-600">
              1
            </div>
            <div>
              <p className="font-medium">Create a Test Payment</p>
              <p className="text-sm text-gray-500">
                Navigate to Payments and click "Create Payment" to process your first transaction
              </p>
            </div>
          </div>
          <div className="flex items-start gap-3">
            <div className="flex h-6 w-6 items-center justify-center rounded-full bg-blue-100 text-sm font-medium text-blue-600">
              2
            </div>
            <div>
              <p className="font-medium">Set Up a Subscription</p>
              <p className="text-sm text-gray-500">
                Go to Subscriptions to create recurring payment plans
              </p>
            </div>
          </div>
          <div className="flex items-start gap-3">
            <div className="flex h-6 w-6 items-center justify-center rounded-full bg-blue-100 text-sm font-medium text-blue-600">
              3
            </div>
            <div>
              <p className="font-medium">Process Refunds</p>
              <p className="text-sm text-gray-500">
                Use the Refunds section to return payments to customers
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
