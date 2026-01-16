import * as React from "react";
import { cn } from "@/lib/utils";

interface TableToolbarProps extends React.HTMLAttributes<HTMLDivElement> {}

const TableToolbar = React.forwardRef<HTMLDivElement, TableToolbarProps>(
  ({ className, ...props }, ref) => (
    <div
      ref={ref}
      className={cn("flex items-center justify-between py-4 gap-4", className)}
      {...props}
    />
  )
);
TableToolbar.displayName = "TableToolbar";

interface TableToolbarLeftProps extends React.HTMLAttributes<HTMLDivElement> {}

const TableToolbarLeft = React.forwardRef<HTMLDivElement, TableToolbarLeftProps>(
  ({ className, ...props }, ref) => (
    <div
      ref={ref}
      className={cn("flex items-center gap-2 flex-1", className)}
      {...props}
    />
  )
);
TableToolbarLeft.displayName = "TableToolbarLeft";

interface TableToolbarRightProps extends React.HTMLAttributes<HTMLDivElement> {}

const TableToolbarRight = React.forwardRef<HTMLDivElement, TableToolbarRightProps>(
  ({ className, ...props }, ref) => (
    <div
      ref={ref}
      className={cn("flex items-center gap-2", className)}
      {...props}
    />
  )
);
TableToolbarRight.displayName = "TableToolbarRight";

export { TableToolbar, TableToolbarLeft, TableToolbarRight };
