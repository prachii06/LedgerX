import { fetchTransactions } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { StatusBadge } from "@/components/ui/status-badge"
import { formatCurrency, formatDate } from "@/lib/utils"
import Link from "next/link"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

export default async function ReconciliationPage() {
  // Fetch transactions and filter to those that have been reconciled or are in issue states
  const allTransactions = await fetchTransactions(200).catch(() => [])
  
  // For the reconciliation queue, we focus on anomalies and matches.
  // PENDING/PROCESSING might not have been checked yet if expected time hasn't passed,
  // but let's show all with their statuses.
  const transactions = allTransactions

  const matched = transactions.filter(t => t.status === "MATCHED").length
  const missing = transactions.filter(t => t.status === "MISSING").length
  const mismatch = transactions.filter(t => t.status === "MISMATCH" || t.status === "CURRENCY_MISMATCH").length
  const pending = transactions.filter(t => t.status === "PENDING" || t.status === "PROCESSING").length

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Reconciliation Operations</h1>
          <p className="text-muted-foreground mt-1">Review and resolve transaction discrepancies.</p>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-4">
        <Card className="bg-success/10 border-success/20">
          <CardHeader className="py-4">
            <CardTitle className="text-sm font-medium">MATCHED</CardTitle>
            <div className="text-2xl font-bold text-green-500">{matched}</div>
          </CardHeader>
        </Card>
        <Card className="bg-warning/10 border-warning/20">
          <CardHeader className="py-4">
            <CardTitle className="text-sm font-medium">MISSING</CardTitle>
            <div className="text-2xl font-bold text-yellow-500">{missing}</div>
          </CardHeader>
        </Card>
        <Card className="bg-destructive/10 border-destructive/20">
          <CardHeader className="py-4">
            <CardTitle className="text-sm font-medium">MISMATCH</CardTitle>
            <div className="text-2xl font-bold text-destructive">{mismatch}</div>
          </CardHeader>
        </Card>
        <Card className="bg-secondary/30">
          <CardHeader className="py-4">
            <CardTitle className="text-sm font-medium">PENDING</CardTitle>
            <div className="text-2xl font-bold text-blue-500">{pending}</div>
          </CardHeader>
        </Card>
      </div>

      <Card>
        <CardHeader className="py-4 px-6 border-b flex flex-row items-center justify-between">
          <CardTitle className="text-base font-medium">Reconciliation Queue</CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="px-6 py-3">Transaction ID</TableHead>
                <TableHead className="px-6 py-3">Amount</TableHead>
                <TableHead className="px-6 py-3">Status</TableHead>
                <TableHead className="px-6 py-3">Last Checked</TableHead>
                <TableHead className="px-6 py-3 text-right">Action</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {transactions.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} className="h-24 text-center text-muted-foreground">
                    No transactions in queue.
                  </TableCell>
                </TableRow>
              ) : (
                transactions.map((t) => (
                  <TableRow key={t.id} className="group cursor-pointer">
                    <TableCell className="px-6 py-4 font-medium font-mono text-xs">
                      {t.external_id || t.id.slice(0, 12) + "..."}
                    </TableCell>
                    <TableCell className="px-6 py-4">
                      {formatCurrency(t.amount, t.currency)}
                    </TableCell>
                    <TableCell className="px-6 py-4">
                      <StatusBadge status={t.status} />
                    </TableCell>
                    <TableCell className="px-6 py-4 text-muted-foreground">
                      {formatDate(t.updated_at)}
                    </TableCell>
                    <TableCell className="px-6 py-4 text-right">
                      <Link
                        href={`/transactions/${t.id}`}
                        className="text-sm font-medium text-primary hover:underline transition-opacity"
                      >
                        Investigate
                      </Link>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  )
}
