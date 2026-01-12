import { useState, useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Card, CardContent } from "@/components/ui/card";
import { Plus } from "lucide-react";
import { PaymentServiceAPI } from "@/lib/api";
import { formatCurrency, formatDate } from "@/lib/utils";
import type { CreateRefundRequest, Payment } from "@/types";

interface CreateRefundFormProps {
  paymentId?: string;
  onSuccess?: () => void;
}

export default function CreateRefundForm({ paymentId, onSuccess }: CreateRefundFormProps) {
  const [open, setOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLoadingPayment, setIsLoadingPayment] = useState(false);
  const [error, setError] = useState("");
  const [selectedPayment, setSelectedPayment] = useState<Payment | null>(null);

  const [formData, setFormData] = useState<CreateRefundRequest>({
    payment_id: paymentId || "",
    amount: 0,
    reason: "",
    notes: "",
  });

  // Load payment details when payment ID is provided or changed
  useEffect(() => {
    if (formData.payment_id && formData.payment_id !== selectedPayment?.id) {
      loadPaymentDetails(formData.payment_id);
    }
  }, [formData.payment_id]);

  const loadPaymentDetails = async (id: string) => {
    setIsLoadingPayment(true);
    setError("");
    try {
      const payment = await PaymentServiceAPI.getPayment(id);
      if (payment.status !== "succeeded") {
        setError("Can only refund succeeded payments");
        setSelectedPayment(null);
      } else {
        setSelectedPayment(payment);
        // Set max refund amount to payment amount
        if (formData.amount === 0) {
          setFormData((prev) => ({ ...prev, amount: payment.amount }));
        }
      }
    } catch (err) {
      setError("Failed to load payment details");
      setSelectedPayment(null);
    } finally {
      setIsLoadingPayment(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    // Validation
    if (!formData.payment_id) {
      setError("Payment ID is required");
      return;
    }

    if (!selectedPayment) {
      setError("Please select a valid payment");
      return;
    }

    if (formData.amount <= 0) {
      setError("Refund amount must be greater than 0");
      return;
    }

    if (formData.amount > selectedPayment.amount) {
      setError("Refund amount cannot exceed payment amount");
      return;
    }

    setIsSubmitting(true);

    try {
      await PaymentServiceAPI.createRefund(formData);
      setOpen(false);
      // Reset form
      setFormData({
        payment_id: "",
        amount: 0,
        reason: "",
        notes: "",
      });
      setSelectedPayment(null);
      onSuccess?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create refund");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus className="h-4 w-4 mr-2" />
          Create Refund
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[600px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Create Refund</DialogTitle>
            <DialogDescription>
              Process a refund for a payment
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <div className="grid gap-2">
              <Label htmlFor="payment_id">Payment ID *</Label>
              <Input
                id="payment_id"
                value={formData.payment_id}
                onChange={(e) =>
                  setFormData({ ...formData, payment_id: e.target.value })
                }
                required
                placeholder="Enter payment ID"
                disabled={!!paymentId}
              />
              <p className="text-xs text-gray-500">
                The ID of the payment to refund
              </p>
            </div>

            {/* Show payment details when loaded */}
            {isLoadingPayment && (
              <div className="text-sm text-gray-500">Loading payment details...</div>
            )}

            {selectedPayment && (
              <Card>
                <CardContent className="pt-4">
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-gray-500">Payment Amount:</span>
                      <span className="font-semibold">
                        {formatCurrency(selectedPayment.amount, selectedPayment.currency)}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-500">Status:</span>
                      <span className="capitalize">{selectedPayment.status}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-500">Created:</span>
                      <span>{formatDate(selectedPayment.created_at)}</span>
                    </div>
                    {selectedPayment.description && (
                      <div className="flex justify-between">
                        <span className="text-gray-500">Description:</span>
                        <span className="text-right">{selectedPayment.description}</span>
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>
            )}

            <div className="grid gap-2">
              <Label htmlFor="amount">
                Refund Amount (in smallest currency unit) *
              </Label>
              <Input
                id="amount"
                type="number"
                min="1"
                max={selectedPayment?.amount || undefined}
                value={formData.amount || ""}
                onChange={(e) =>
                  setFormData({ ...formData, amount: parseInt(e.target.value) || 0 })
                }
                required
                disabled={!selectedPayment}
                placeholder={selectedPayment ? `Max: ${selectedPayment.amount}` : ""}
              />
              {selectedPayment && formData.amount > 0 && (
                <p className="text-xs text-gray-500">
                  = {(formData.amount / 100).toFixed(2)} {selectedPayment.currency} (Full refund:{" "}
                  {(selectedPayment.amount / 100).toFixed(2)} {selectedPayment.currency})
                </p>
              )}
            </div>

            <div className="grid gap-2">
              <Label htmlFor="reason">Reason</Label>
              <Select
                value={formData.reason}
                onValueChange={(value) =>
                  setFormData({ ...formData, reason: value })
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select a reason" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="duplicate">Duplicate</SelectItem>
                  <SelectItem value="fraudulent">Fraudulent</SelectItem>
                  <SelectItem value="requested_by_customer">
                    Requested by customer
                  </SelectItem>
                  <SelectItem value="other">Other</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="notes">Notes (optional)</Label>
              <Textarea
                id="notes"
                value={formData.notes}
                onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                placeholder="Additional information about this refund"
                rows={3}
              />
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setOpen(false)}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={isSubmitting || !selectedPayment}
              variant="destructive"
            >
              {isSubmitting ? "Processing..." : "Process Refund"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
