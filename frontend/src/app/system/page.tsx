import { fetchHealth } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Database, Server, ShieldAlert } from "lucide-react"

export default async function SystemHealthPage() {
  const health = await fetchHealth().catch(() => ({ status: "Unavailable", postgres: "Unavailable", redis: "Unavailable", kafka: "Unavailable" }))

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">System Health</h1>
          <p className="text-muted-foreground mt-1">LedgerX Infrastructure Status.</p>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader className="flex flex-row items-start justify-between pb-2">
            <div>
              <CardTitle className="text-lg flex items-center gap-2"><Server className="h-5 w-5" /> Backend API</CardTitle>
              <CardDescription>Main Go backend service</CardDescription>
            </div>
            <Badge variant={health.status === "OK" ? "success" : "destructive"}>{health.status || "Unknown"}</Badge>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground mt-4">
              Handles HTTP ingestion of transactions and exposes data for this dashboard.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-start justify-between pb-2">
            <div>
              <CardTitle className="text-lg flex items-center gap-2"><Database className="h-5 w-5" /> PostgreSQL</CardTitle>
              <CardDescription>Primary persistence</CardDescription>
            </div>
            <Badge variant={health.postgres === "OK" ? "success" : "destructive"}>{health.postgres || "Unknown"}</Badge>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground mt-4">
              Stores transactions, events, and reconciliation history persistently.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-start justify-between pb-2">
            <div>
              <CardTitle className="text-lg flex items-center gap-2"><Database className="h-5 w-5 text-red-500" /> Redis Cache</CardTitle>
              <CardDescription>In-memory cache</CardDescription>
            </div>
            <Badge variant={health.redis === "OK" ? "success" : "destructive"}>{health.redis || "Unknown"}</Badge>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground mt-4">
              Stores detailed reconciliation results for fast retrieval and handles duplicate checking.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-start justify-between pb-2">
            <div>
              <CardTitle className="text-lg flex items-center gap-2"><ShieldAlert className="h-5 w-5 text-orange-500" /> Kafka Event Bus</CardTitle>
              <CardDescription>Message broker</CardDescription>
            </div>
            <Badge variant={health.kafka === "OK" ? "success" : "destructive"}>{health.kafka || "Unknown"}</Badge>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground mt-4">
              Asynchronous event stream processing. Consumes events emitted by upstream services.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
