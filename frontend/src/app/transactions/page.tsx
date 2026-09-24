import { fetchTransactions } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { StatusBadge } from "@/components/ui/status-badge"
import { formatCurrency, formatDate } from "@/lib/utils"
import Link from "next/link"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

export default async function TransactionsPage() {
  const transactions = await fetchTransactions(200).catch(() => [])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Transactions</h1>
          <p className="text-muted-foreground mt-1">View and filter all system transactions.</p>
        </div>
      </div>

      <Card>
        <CardHeader className="py-4 px-6 border-b">
          <CardTitle className="text-base font-medium">All Transactions</CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="px-6 py-3">Transaction ID</TableHead>
                <TableHead className="px-6 py-3">Amount</TableHead>
                <TableHead className="px-6 py-3">Status</TableHead>
                <TableHead className="px-6 py-3">Created At</TableHead>
                <TableHead className="px-6 py-3">Updated At</TableHead>
                <TableHead className="px-6 py-3 text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {transactions.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="h-24 text-center text-muted-foreground">
                    No transactions found.
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
                      {formatDate(t.created_at)}
                    </TableCell>
                    <TableCell className="px-6 py-4 text-muted-foreground">
                      {formatDate(t.updated_at)}
                    </TableCell>
                    <TableCell className="px-6 py-4 text-right">
                      <Link
                        href={`/transactions/${t.id}`}
                        className="text-sm font-medium text-primary hover:underline opacity-0 group-hover:opacity-100 transition-opacity"
                      >
                        View Details
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
