import { useState } from "react";
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
import { Plus } from "lucide-react";
import { PaymentServiceAPI } from "@/lib/api";
import type { CreateSubscriptionRequest, Provider, Currency } from "@/types";

interface CreateSubscriptionFormProps {
  onSuccess?: () => void;
}

export default function CreateSubscriptionForm({ onSuccess }: CreateSubscriptionFormProps) {
  const [open, setOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");

  const [formData, setFormData] = useState<CreateSubscriptionRequest>({
    provider: "stripe",
    amount: 0,
    currency: "SEK",
    interval: "month",
    interval_count: 1,
    product_name: "",
    product_description: "",
    trial_period_days: 0,
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    // Validation
    if (formData.amount <= 0) {
      setError("Amount must be greater than 0");
      return;
    }

    if (!formData.product_name.trim()) {
      setError("Product name is required");
      return;
    }

    if (formData.interval_count < 1) {
      setError("Interval count must be at least 1");
      return;
    }

    setIsSubmitting(true);

    try {
      await PaymentServiceAPI.createSubscription(formData);
      setOpen(false);
      // Reset form
      setFormData({
        provider: "stripe",
        amount: 0,
        currency: "SEK",
        interval: "month",
        interval_count: 1,
        product_name: "",
        product_description: "",
        trial_period_days: 0,
      });
      onSuccess?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create subscription");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus className="h-4 w-4 mr-2" />
          Create Subscription
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Create Subscription</DialogTitle>
            <DialogDescription>
              Create a new recurring subscription
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <div className="grid gap-2">
              <Label htmlFor="product_name">Product Name *</Label>
              <Input
                id="product_name"
                value={formData.product_name}
                onChange={(e) =>
                  setFormData({ ...formData, product_name: e.target.value })
                }
                required
                placeholder="Premium Plan"
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="product_description">Product Description</Label>
              <Textarea
                id="product_description"
                value={formData.product_description}
                onChange={(e) =>
                  setFormData({ ...formData, product_description: e.target.value })
                }
                placeholder="Access to all premium features"
                rows={2}
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="provider">Provider</Label>
                <Select
                  value={formData.provider}
                  onValueChange={(value: Provider) =>
                    setFormData({ ...formData, provider: value })
                  }
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="stripe">Stripe</SelectItem>
                    <SelectItem value="swish">Swish</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="grid gap-2">
                <Label htmlFor="currency">Currency</Label>
                <Select
                  value={formData.currency}
                  onValueChange={(value: Currency) =>
                    setFormData({ ...formData, currency: value })
                  }
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="SEK">SEK</SelectItem>
                    <SelectItem value="USD">USD</SelectItem>
                    <SelectItem value="EUR">EUR</SelectItem>
                    <SelectItem value="GBP">GBP</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="amount">
                Amount (in smallest currency unit, e.g., öre/cents) *
              </Label>
              <Input
                id="amount"
                type="number"
                min="1"
                value={formData.amount || ""}
                onChange={(e) =>
                  setFormData({ ...formData, amount: parseInt(e.target.value) || 0 })
                }
                required
                placeholder="9900 = 99.00 SEK"
              />
              <p className="text-xs text-gray-500">
                {formData.amount > 0 && `= ${(formData.amount / 100).toFixed(2)} ${formData.currency}`}
              </p>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="interval">Interval *</Label>
                <Select
                  value={formData.interval}
                  onValueChange={(value) =>
                    setFormData({ ...formData, interval: value })
                  }
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="day">Day</SelectItem>
                    <SelectItem value="week">Week</SelectItem>
                    <SelectItem value="month">Month</SelectItem>
                    <SelectItem value="year">Year</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="grid gap-2">
                <Label htmlFor="interval_count">Interval Count *</Label>
                <Input
                  id="interval_count"
                  type="number"
                  min="1"
                  value={formData.interval_count}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      interval_count: parseInt(e.target.value) || 1,
                    })
                  }
                  required
                />
              </div>
            </div>

            <p className="text-sm text-gray-500">
              Billing frequency: Every{" "}
              {formData.interval_count > 1 ? formData.interval_count : ""}{" "}
              {formData.interval}
              {formData.interval_count > 1 ? "s" : ""}
            </p>

            <div className="grid gap-2">
              <Label htmlFor="trial_period_days">Trial Period (days)</Label>
              <Input
                id="trial_period_days"
                type="number"
                min="0"
                value={formData.trial_period_days || ""}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    trial_period_days: parseInt(e.target.value) || 0,
                  })
                }
                placeholder="0 = No trial"
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
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Creating..." : "Create Subscription"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
