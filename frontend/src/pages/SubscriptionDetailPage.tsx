import { useState } from "react";
import { useParams, Link } from "react-router-dom";
import { useSubscription } from "@/lib/hooks";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Loading } from "@/components/ui/loading";
import { Alert, AlertDescription } from "@/components/ui/alert";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import StatusBadge from "@/components/admin/StatusBadge";
import { formatCurrency, formatDate } from "@/lib/utils";
import { ArrowLeft, XCircle } from "lucide-react";
import { PaymentServiceAPI } from "@/lib/api";

export default function SubscriptionDetailPage() {
  const { id } = useParams<{ id: string }>();
  const { data: subscription, isLoading, error, refetch } = useSubscription(id!);

  const [showCancelDialog, setShowCancelDialog] = useState(false);
  const [cancelImmediate, setCancelImmediate] = useState(false);
  const [isCancelling, setIsCancelling] = useState(false);
  const [cancelError, setCancelError] = useState("");

  const handleCancelSubscription = async () => {
    setCancelError("");
    setIsCancelling(true);

    try {
      await PaymentServiceAPI.cancelSubscription(id!, cancelImmediate);
      setShowCancelDialog(false);
      refetch();
    } catch (err) {
      setCancelError(err instanceof Error ? err.message : "Failed to cancel subscription");
    } finally {
      setIsCancelling(false);
    }
  };

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
        <Link to="/dashboard/subscriptions">
          <Button variant="ghost" size="sm">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Subscriptions
          </Button>
        </Link>
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      </div>
    );
  }

  if (!subscription) {
    return (
      <div className="space-y-4">
        <Link to="/dashboard/subscriptions">
          <Button variant="ghost" size="sm">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Subscriptions
          </Button>
        </Link>
        <p className="text-gray-500">Subscription not found</p>
      </div>
    );
  }

  const canCancel = ["active", "trialing", "past_due"].includes(subscription.status);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link to="/dashboard/subscriptions">
            <Button variant="ghost" size="sm">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back
            </Button>
          </Link>
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Subscription Details</h1>
            <p className="text-gray-500 font-mono text-sm mt-1">{subscription.id}</p>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <StatusBadge status={subscription.status} />
          {canCancel && (
            <Button
              variant="destructive"
              size="sm"
              onClick={() => setShowCancelDialog(true)}
            >
              <XCircle className="h-4 w-4 mr-2" />
              Cancel Subscription
            </Button>
          )}
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Product Information */}
        <Card>
          <CardHeader>
            <CardTitle>Product Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div>
              <p className="text-sm text-gray-500">Product Name</p>
              <p className="text-lg font-semibold">{subscription.product_name}</p>
            </div>
            {subscription.product_description && (
              <div>
                <p className="text-sm text-gray-500">Description</p>
                <p>{subscription.product_description}</p>
              </div>
            )}
            <div>
              <p className="text-sm text-gray-500">Provider</p>
              <p className="capitalize">{subscription.provider}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Provider Subscription ID</p>
              <p className="font-mono text-sm">{subscription.provider_subscription_id}</p>
            </div>
          </CardContent>
        </Card>

        {/* Billing Information */}
        <Card>
          <CardHeader>
            <CardTitle>Billing Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div>
              <p className="text-sm text-gray-500">Amount</p>
              <p className="text-lg font-semibold">
                {formatCurrency(subscription.amount, subscription.currency)}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Billing Frequency</p>
              <p>
                Every {subscription.interval_count > 1 ? subscription.interval_count : ""}{" "}
                {subscription.interval}
                {subscription.interval_count > 1 ? "s" : ""}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Current Period</p>
              <p className="text-sm">
                {formatDate(subscription.current_period_start)} -{" "}
                {formatDate(subscription.current_period_end)}
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Trial Information */}
        {(subscription.trial_start || subscription.trial_end) && (
          <Card>
            <CardHeader>
              <CardTitle>Trial Information</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              {subscription.trial_start && (
                <div>
                  <p className="text-sm text-gray-500">Trial Start</p>
                  <p>{formatDate(subscription.trial_start)}</p>
                </div>
              )}
              {subscription.trial_end && (
                <div>
                  <p className="text-sm text-gray-500">Trial End</p>
                  <p>{formatDate(subscription.trial_end)}</p>
                </div>
              )}
            </CardContent>
          </Card>
        )}

        {/* Cancellation Information */}
        {(subscription.cancel_at || subscription.canceled_at || subscription.cancel_at_period_end) && (
          <Card>
            <CardHeader>
              <CardTitle>Cancellation Information</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              {subscription.cancel_at_period_end && (
                <div>
                  <p className="text-sm text-gray-500">Status</p>
                  <p className="text-orange-600">Scheduled to cancel at period end</p>
                </div>
              )}
              {subscription.cancel_at && (
                <div>
                  <p className="text-sm text-gray-500">Cancel At</p>
                  <p>{formatDate(subscription.cancel_at)}</p>
                </div>
              )}
              {subscription.canceled_at && (
                <div>
                  <p className="text-sm text-gray-500">Canceled At</p>
                  <p>{formatDate(subscription.canceled_at)}</p>
                </div>
              )}
            </CardContent>
          </Card>
        )}

        {/* Timestamps */}
        <Card>
          <CardHeader>
            <CardTitle>Timestamps</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div>
              <p className="text-sm text-gray-500">Created</p>
              <p>{formatDate(subscription.created_at)}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Last Updated</p>
              <p>{formatDate(subscription.updated_at)}</p>
            </div>
          </CardContent>
        </Card>

        {/* Metadata */}
        {subscription.metadata && Object.keys(subscription.metadata).length > 0 && (
          <Card className="md:col-span-2">
            <CardHeader>
              <CardTitle>Metadata</CardTitle>
            </CardHeader>
            <CardContent>
              <pre className="text-xs bg-gray-50 p-3 rounded overflow-auto">
                {JSON.stringify(subscription.metadata, null, 2)}
              </pre>
            </CardContent>
          </Card>
        )}
      </div>

      {/* Cancel Subscription Dialog */}
      <Dialog open={showCancelDialog} onOpenChange={setShowCancelDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Cancel Subscription</DialogTitle>
            <DialogDescription>
              Choose how you want to cancel this subscription
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            {cancelError && (
              <Alert variant="destructive">
                <AlertDescription>{cancelError}</AlertDescription>
              </Alert>
            )}

            <div className="space-y-3">
              <div
                className={`border rounded-lg p-4 cursor-pointer transition-colors ${
                  !cancelImmediate ? "border-blue-500 bg-blue-50" : "border-gray-200 hover:border-gray-300"
                }`}
                onClick={() => setCancelImmediate(false)}
              >
                <div className="flex items-start gap-3">
                  <input
                    type="radio"
                    checked={!cancelImmediate}
                    onChange={() => setCancelImmediate(false)}
                    className="mt-1"
                  />
                  <div>
                    <Label className="font-medium">Cancel at period end</Label>
                    <p className="text-sm text-gray-500 mt-1">
                      Subscription will remain active until {formatDate(subscription.current_period_end)}
                    </p>
                  </div>
                </div>
              </div>

              <div
                className={`border rounded-lg p-4 cursor-pointer transition-colors ${
                  cancelImmediate ? "border-red-500 bg-red-50" : "border-gray-200 hover:border-gray-300"
                }`}
                onClick={() => setCancelImmediate(true)}
              >
                <div className="flex items-start gap-3">
                  <input
                    type="radio"
                    checked={cancelImmediate}
                    onChange={() => setCancelImmediate(true)}
                    className="mt-1"
                  />
                  <div>
                    <Label className="font-medium">Cancel immediately</Label>
                    <p className="text-sm text-gray-500 mt-1">
                      Subscription will be canceled right away
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setShowCancelDialog(false)}
              disabled={isCancelling}
            >
              Go Back
            </Button>
            <Button
              variant="destructive"
              onClick={handleCancelSubscription}
              disabled={isCancelling}
            >
              {isCancelling ? "Cancelling..." : "Confirm Cancellation"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
