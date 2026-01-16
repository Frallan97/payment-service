import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

const badgeVariants = cva(
  "inline-flex items-center rounded-md border px-2.5 py-1 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
  {
    variants: {
      variant: {
        default:
          "border-primary-500/20 bg-primary-100 text-primary-600 shadow-sm",
        secondary:
          "border-neutral-300 bg-neutral-100 text-neutral-700 shadow-sm",
        destructive:
          "border-error-500/20 bg-error-50 text-error-600 shadow-sm",
        error:
          "border-error-500/20 bg-error-50 text-error-600 shadow-sm",
        outline: "border-neutral-300 text-foreground",
        success: "border-success-500/20 bg-success-50 text-success-600 shadow-sm",
        warning: "border-warning-500/20 bg-warning-50 text-warning-600 shadow-sm",
        info: "border-info-500/20 bg-info-50 text-info-600 shadow-sm",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
);

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <div className={cn(badgeVariants({ variant }), className)} {...props} />
  );
}

export { Badge, badgeVariants };
