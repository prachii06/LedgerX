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
          "border-transparent bg-primary/20 text-primary hover:bg-primary/30": variant === "default",
          "border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80": variant === "secondary",
          "border-transparent bg-red-500/20 text-red-400 hover:bg-red-500/30": variant === "destructive",
          "border-transparent bg-green-500/20 text-green-400 hover:bg-green-500/30": variant === "success",
          "border-transparent bg-yellow-500/20 text-yellow-400 hover:bg-yellow-500/30": variant === "warning",
          "border-transparent bg-blue-500/20 text-blue-400 hover:bg-blue-500/30": variant === "pending",
          "text-foreground border-border": variant === "outline",
        },
        className
      )}
      {...props}
    />
  )
}

export { Badge }
