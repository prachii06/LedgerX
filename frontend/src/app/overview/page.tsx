import { fetchHealth, fetchTransactions, getOverview } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { StatusBadge } from "@/components/ui/status-badge"
import { formatCurrency, formatDate } from "@/lib/utils"
import Link from "next/link"
import { ActivitySquare, Database, ListOrdered, Server, ShieldAlert, FileWarning, RefreshCcw } from "lucide-react"
import { SimulationButton } from "@/components/SimulationButton"

export default async function OverviewPage() {
  const transactions = await fetchTransactions(10).catch(() => [])
  const health = await fetchHealth().catch(() => ({ status: "Unavailable", postgres: "Unavailable", redis: "Unavailable", kafka: "Unavailable" }))

  const stats = await getOverview().catch(() => ({ total_transactions: 0, reconciled: 0, issues: 0, pending: 0 }))

  const total = stats.total_transactions
  const matched = stats.reconciled
  const issues = stats.issues
  const pending = stats.pending

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Reconciliation Overview</h1>
          <p className="text-muted-foreground mt-1">Monitor transaction processing and reconciliation health.</p>
        </div>
        <div className="flex items-center gap-2">
          <SimulationButton />
        </div>
      </div>


      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Transactions</CardTitle>
            <ListOrdered className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{total}</div>
            <p className="text-xs text-muted-foreground">Processed by LedgerX</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Reconciled</CardTitle>
            <ActivitySquare className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-500">{matched}</div>
            <p className="text-xs text-muted-foreground">Successfully matched events</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Reconciliation Issues</CardTitle>
            <FileWarning className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-destructive">{issues}</div>
            <p className="text-xs text-muted-foreground">Require manual investigation</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Pending</CardTitle>
            <RefreshCcw className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-500">{pending}</div>
            <p className="text-xs text-muted-foreground">Awaiting events</p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle>Recent Transactions</CardTitle>
            <CardDescription>Latest processed transactions across the platform.</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-8">
              {transactions.slice(0, 5).map(t => (
                <div key={t.id} className="flex items-center">
                  <div className="ml-4 space-y-1">
                    <Link href={`/transactions/${t.id}`} className="text-sm font-medium leading-none hover:underline">
                      {t.external_id || t.id.slice(0, 8)}
                    </Link>
                    <p className="text-sm text-muted-foreground">
                      {formatDate(t.created_at)}
                    </p>
                  </div>
                  <div className="ml-auto flex items-center gap-4">
                    <StatusBadge status={t.status} />
                    <div className="font-medium">{formatCurrency(t.amount, t.currency)}</div>
                  </div>
                </div>
              ))}
              {transactions.length === 0 && (
                <div className="text-center py-4 text-sm text-muted-foreground">No transactions found.</div>
              )}
            </div>
          </CardContent>
        </Card>

        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>System Health</CardTitle>
            <CardDescription>Infrastructure status.</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between border-b pb-2">
                <div className="flex items-center gap-2">
                  <Server className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-medium">Backend API</span>
                </div>
                <Badge variant={health.status === "OK" ? "success" : "destructive"}>{health.status || "Unknown"}</Badge>
              </div>
              <div className="flex items-center justify-between border-b pb-2">
                <div className="flex items-center gap-2">
                  <Database className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-medium">PostgreSQL</span>
                </div>
                <Badge variant={health.postgres === "OK" ? "success" : "destructive"}>{health.postgres || "Unknown"}</Badge>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
