import * as React from "react"
import { cn } from "@/lib/utils"

export interface BadgeProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "secondary" | "destructive" | "outline" | "success" | "warning" | "pending"
}

function Badge({ className, variant = "default", ...props }: BadgeProps) {
  return (
    <div
      className={cn(
        "inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
        {
          "border-border bg-muted/10 text-foreground hover:bg-muted/20": variant === "default",
          "border-muted bg-secondary/20 text-muted-foreground hover:bg-secondary/40": variant === "secondary",
          "border-red-500/30 bg-red-500/10 text-red-500 hover:bg-red-500/20": variant === "destructive",
          "border-green-500/30 bg-green-500/10 text-green-500 hover:bg-green-500/20": variant === "success",
          "border-yellow-500/30 bg-yellow-500/10 text-yellow-500 hover:bg-yellow-500/20": variant === "warning",
          "border-blue-500/30 bg-blue-500/10 text-blue-500 hover:bg-blue-500/20": variant === "pending",
          "text-foreground border-border": variant === "outline",
        },
        className
      )}
      {...props}
    />
  )
}

export { Badge }
