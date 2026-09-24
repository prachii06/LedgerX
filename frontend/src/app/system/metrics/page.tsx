import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Server } from "lucide-react"

// A simple function to fetch and parse the prometheus text format
async function fetchMetrics() {
  try {
    const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"
    const res = await fetch(`${API_BASE_URL}/metrics`, { cache: "no-store" })
    if (!res.ok) return null
    const text = await res.text()
    
    // Basic parser for demonstration
    const metrics: Record<string, string> = {}
    text.split('\n').forEach(line => {
      if (line && !line.startsWith('#')) {
        const [key, value] = line.split(' ')
        if (key && value) metrics[key] = value
      }
    })
    return metrics
  } catch (error) {
    return null
  }
}

export default async function MetricsPage() {
  const metrics = await fetchMetrics()

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">System Metrics</h1>
          <p className="text-muted-foreground mt-1">Application telemetry from Prometheus endpoint.</p>
        </div>
      </div>

      {!metrics ? (
        <Card>
          <CardContent className="pt-6 text-center text-muted-foreground">
            Unable to fetch metrics from the backend.
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">Reconciliations (MATCHED)</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {metrics['ledgerx_reconciliations_total{status="MATCHED"}'] || "0"}
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">Reconciliations (MISSING)</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {metrics['ledgerx_reconciliations_total{status="MISSING"}'] || "0"}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">Reconciliations (MISMATCH)</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {metrics['ledgerx_reconciliations_total{status="MISMATCH"}'] || "0"}
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">Total Processed Events</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {metrics['ledgerx_events_processed_total'] || "0"}
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">Go Goroutines</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {metrics['go_goroutines'] || "N/A"}
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">Go Allocated Memory</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {metrics['go_memstats_alloc_bytes'] ? `${(parseInt(metrics['go_memstats_alloc_bytes']) / 1024 / 1024).toFixed(2)} MB` : "N/A"}
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  )
}
