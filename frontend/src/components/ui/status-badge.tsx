import { Badge } from "@/components/ui/badge"
import { TransactionStatus } from "@/types"

export function StatusBadge({ status }: { status: TransactionStatus }) {
  let variant: "default" | "secondary" | "destructive" | "outline" | "success" | "warning" | "pending" = "default"
  let label = status

  switch (status) {
    case "MATCHED":
      variant = "success"
      break
    case "MISSING":
    case "OUT_OF_ORDER":
    case "DUPLICATE":
      variant = "warning"
      break
    case "MISMATCH":
    case "CURRENCY_MISMATCH":
      variant = "destructive"
      break
    case "PENDING":
    case "PROCESSING":
      variant = "pending"
      break
    default:
      variant = "secondary"
  }

  return <Badge variant={variant}>{label}</Badge>
}
