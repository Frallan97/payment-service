import { useParams, Link } from "react-router-dom";
import { usePayment } from "@/lib/hooks";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Loading } from "@/components/ui/loading";
import { Alert, AlertDescription } from "@/components/ui/alert";
import StatusBadge from "@/components/admin/StatusBadge";
import { formatCurrency, formatDate } from "@/lib/utils";
import { ArrowLeft, ExternalLink } from "lucide-react";

export default function PaymentDetailPage() {
  const { id } = useParams<{ id: string }>();
  const { data: payment, isLoading, error } = usePayment(id!);

  if (isLoading) {
    return (
      <div className="py-12">
        <Loading size="lg" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="space-y-4">
        <Link to="/dashboard/payments">
          <Button variant="ghost" size="sm">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Payments
          </Button>
        </Link>
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      </div>
    );
  }

  if (!payment) {
    return (
      <div className="space-y-4">
        <Link to="/dashboard/payments">
          <Button variant="ghost" size="sm">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Payments
          </Button>
        </Link>
        <p className="text-gray-500">Payment not found</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link to="/dashboard/payments">
            <Button variant="ghost" size="sm">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back
            </Button>
          </Link>
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Payment Details</h1>
            <p className="text-gray-500 font-mono text-sm mt-1">{payment.id}</p>
          </div>
        </div>
        <StatusBadge status={payment.status} />
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Payment Information */}
        <Card>
          <CardHeader>
            <CardTitle>Payment Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div>
              <p className="text-sm text-gray-500">Amount</p>
              <p className="text-lg font-semibold">
                {formatCurrency(payment.amount, payment.currency)}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Provider</p>
              <p className="capitalize">{payment.provider}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Provider Payment ID</p>
              <div className="flex items-center gap-2">
                <p className="font-mono text-sm">{payment.provider_payment_id}</p>
                {payment.provider === "stripe" && (
                  <a
                    href={`https://dashboard.stripe.com/payments/${payment.provider_payment_id}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-blue-600 hover:text-blue-700"
                  >
                    <ExternalLink className="h-4 w-4" />
                  </a>
                )}
              </div>
            </div>
            {payment.description && (
              <div>
                <p className="text-sm text-gray-500">Description</p>
                <p>{payment.description}</p>
              </div>
            )}
            {payment.statement_descriptor && (
              <div>
                <p className="text-sm text-gray-500">Statement Descriptor</p>
                <p>{payment.statement_descriptor}</p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Payment Method */}
        <Card>
          <CardHeader>
            <CardTitle>Payment Method</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {payment.payment_method_type && (
              <div>
                <p className="text-sm text-gray-500">Type</p>
                <p className="capitalize">{payment.payment_method_type}</p>
              </div>
            )}
            {payment.payment_method_details && (
              <div>
                <p className="text-sm text-gray-500">Details</p>
                <pre className="text-xs bg-gray-50 p-2 rounded">
                  {JSON.stringify(payment.payment_method_details, null, 2)}
                </pre>
              </div>
            )}
            {!payment.payment_method_type && !payment.payment_method_details && (
              <p className="text-sm text-gray-500">No payment method details available</p>
            )}
          </CardContent>
        </Card>

        {/* Status Information */}
        <Card>
          <CardHeader>
            <CardTitle>Status Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div>
              <p className="text-sm text-gray-500">Status</p>
              <div className="mt-1">
                <StatusBadge status={payment.status} />
              </div>
            </div>
            {payment.failure_code && (
              <div>
                <p className="text-sm text-gray-500">Failure Code</p>
                <p className="text-red-600">{payment.failure_code}</p>
              </div>
            )}
            {payment.failure_message && (
              <div>
                <p className="text-sm text-gray-500">Failure Message</p>
                <p className="text-red-600">{payment.failure_message}</p>
              </div>
            )}
            {payment.client_secret && (
              <div>
                <p className="text-sm text-gray-500">Client Secret</p>
                <p className="font-mono text-xs truncate">{payment.client_secret}</p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Timestamps */}
        <Card>
          <CardHeader>
            <CardTitle>Timestamps</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div>
              <p className="text-sm text-gray-500">Created</p>
              <p>{formatDate(payment.created_at)}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Last Updated</p>
              <p>{formatDate(payment.updated_at)}</p>
            </div>
            {payment.completed_at && (
              <div>
                <p className="text-sm text-gray-500">Completed</p>
                <p>{formatDate(payment.completed_at)}</p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Related Information */}
        {(payment.subscription_id || payment.invoice_id) && (
          <Card className="md:col-span-2">
            <CardHeader>
              <CardTitle>Related Information</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              {payment.subscription_id && (
                <div>
                  <p className="text-sm text-gray-500">Subscription ID</p>
                  <Link
                    to={`/dashboard/subscriptions/${payment.subscription_id}`}
                    className="text-blue-600 hover:text-blue-700 font-mono text-sm"
                  >
                    {payment.subscription_id}
                  </Link>
                </div>
              )}
              {payment.invoice_id && (
                <div>
                  <p className="text-sm text-gray-500">Invoice ID</p>
                  <p className="font-mono text-sm">{payment.invoice_id}</p>
                </div>
              )}
            </CardContent>
          </Card>
        )}

        {/* Metadata */}
        {payment.metadata && Object.keys(payment.metadata).length > 0 && (
          <Card className="md:col-span-2">
            <CardHeader>
              <CardTitle>Metadata</CardTitle>
            </CardHeader>
            <CardContent>
              <pre className="text-xs bg-gray-50 p-3 rounded overflow-auto">
                {JSON.stringify(payment.metadata, null, 2)}
              </pre>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
