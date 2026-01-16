import * as React from "react";
import { cn } from "@/lib/utils";

interface FieldGroupProps extends React.HTMLAttributes<HTMLDivElement> {
  label?: string;
  description?: string;
}

const FieldGroup = React.forwardRef<HTMLDivElement, FieldGroupProps>(
  ({ className, label, description, children, ...props }, ref) => (
    <div ref={ref} className={cn("space-y-4", className)} {...props}>
      {(label || description) && (
        <div className="space-y-1">
          {label && (
            <h3 className="text-base font-semibold text-neutral-800">{label}</h3>
          )}
          {description && (
            <p className="text-sm text-neutral-500">{description}</p>
          )}
        </div>
      )}
      <div className="space-y-4">{children}</div>
    </div>
  )
);
FieldGroup.displayName = "FieldGroup";

interface FieldProps extends React.HTMLAttributes<HTMLDivElement> {}

const Field = React.forwardRef<HTMLDivElement, FieldProps>(
  ({ className, ...props }, ref) => (
    <div ref={ref} className={cn("space-y-2", className)} {...props} />
  )
);
Field.displayName = "Field";

interface FieldLabelProps extends React.LabelHTMLAttributes<HTMLLabelElement> {}

const FieldLabel = React.forwardRef<HTMLLabelElement, FieldLabelProps>(
  ({ className, ...props }, ref) => (
    <label
      ref={ref}
      className={cn("text-sm font-medium text-neutral-700", className)}
      {...props}
    />
  )
);
FieldLabel.displayName = "FieldLabel";

interface FieldHintProps extends React.HTMLAttributes<HTMLParagraphElement> {}

const FieldHint = React.forwardRef<HTMLParagraphElement, FieldHintProps>(
  ({ className, ...props }, ref) => (
    <p
      ref={ref}
      className={cn("text-xs text-neutral-500", className)}
      {...props}
    />
  )
);
FieldHint.displayName = "FieldHint";

interface FieldErrorProps extends React.HTMLAttributes<HTMLParagraphElement> {}

const FieldError = React.forwardRef<HTMLParagraphElement, FieldErrorProps>(
  ({ className, ...props }, ref) => (
    <p
      ref={ref}
      className={cn("text-xs text-error-600 font-medium", className)}
      {...props}
    />
  )
);
FieldError.displayName = "FieldError";

export { FieldGroup, Field, FieldLabel, FieldHint, FieldError };
