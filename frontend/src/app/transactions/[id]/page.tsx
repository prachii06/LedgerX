import { fetchTransaction, fetchEvents, getReconciliation, getReconciliationHistory } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { StatusBadge } from "@/components/ui/status-badge"
import { formatCurrency, formatDate } from "@/lib/utils"
import { Badge } from "@/components/ui/badge"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { AlertCircle, CheckCircle2, ChevronRight, Copy, RefreshCcw, XCircle } from "lucide-react"

export default async function TransactionDetailPage({ params }: { params: Promise<{ id: string }> | { id: string } }) {
  const resolvedParams = await params
  const transactionId = resolvedParams.id

  const [transaction, events, reconciliation, history] = await Promise.all([
    fetchTransaction(transactionId).catch(() => null),
    fetchEvents(100, 0, transactionId).catch(() => []),
    getReconciliation(transactionId).catch(() => null),
    getReconciliationHistory(transactionId).catch(() => []),
  ])

  if (!transaction) {
    return (
      <div className="flex flex-col items-center justify-center h-[60vh] text-center space-y-4">
        <XCircle className="h-12 w-12 text-muted-foreground" />
        <h2 className="text-xl font-semibold">Transaction Not Found</h2>
        <p className="text-muted-foreground">The transaction ID {transactionId} does not exist or could not be loaded.</p>
      </div>
    )
  }

  return (
    <div className="space-y-6 pb-20">
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-3">
            <h1 className="text-3xl font-bold tracking-tight font-mono text-sm uppercase">Txn: {transaction.external_id || transaction.id.split("-")[0]}</h1>
            <StatusBadge status={transaction.status} />
          </div>
          <p className="text-muted-foreground mt-1 text-sm font-mono">{transaction.id}</p>
        </div>
        <div className="flex items-center gap-4">
          <div className="text-right">
            <div className="text-2xl font-bold">{formatCurrency(transaction.amount, transaction.currency)}</div>
            <div className="text-sm text-muted-foreground">Amount</div>
          </div>
        </div>
      </div>

      {reconciliation && (
        <Card className="border-primary/20 bg-primary/5">
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-medium flex items-center gap-2">
              <RefreshCcw className="h-4 w-4" />
              Reconciliation Summary
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-4 text-sm">
              <div>
                <span className="text-muted-foreground block mb-1">Status</span>
                <StatusBadge status={reconciliation.status} />
              </div>
              <div>
                <span className="text-muted-foreground block mb-1">Expected Events</span>
                <span className="font-semibold text-base">{reconciliation.expected_events.length}</span>
              </div>
              <div>
                <span className="text-muted-foreground block mb-1">Received Events</span>
                <span className="font-semibold text-base">{reconciliation.received_events.length}</span>
              </div>
              <div>
                <span className="text-muted-foreground block mb-1">Missing Events</span>
                <span className="font-semibold text-base text-destructive">{reconciliation.missing_events?.length || 0}</span>
              </div>
            </div>
            {reconciliation.message && (
              <div className="mt-4 p-3 bg-background rounded-md text-sm border font-medium">
                {reconciliation.message}
              </div>
            )}
          </CardContent>
        </Card>
      )}

      <Card>
        <CardHeader>
          <CardTitle>Event Flow</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-6 pl-4 border-l-2 ml-4 pb-4">
            {reconciliation?.expected_events.map((eventType, index) => {
              const matchedEvent = events.find(e => e.event_type === eventType)
              return (
                <div key={eventType} className="relative">
                  <div className={`absolute -left-[25px] top-1 rounded-full p-0.5 bg-background border ${matchedEvent ? 'text-green-500 border-green-500' : 'text-muted-foreground border-muted'}`}>
                    {matchedEvent ? <CheckCircle2 className="h-5 w-5 bg-background rounded-full" /> : <AlertCircle className="h-5 w-5 bg-background rounded-full" />}
                  </div>
                  <div className="pl-6 space-y-2">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <h4 className={`font-semibold ${matchedEvent ? 'text-foreground' : 'text-muted-foreground'}`}>{eventType}</h4>
                        {matchedEvent && <Badge variant="outline" className="text-[10px]">Seq {matchedEvent.sequence}</Badge>}
                      </div>
                      <span className="text-xs text-muted-foreground">
                        {matchedEvent ? formatDate(matchedEvent.received_at) : 'Waiting...'}
                      </span>
                    </div>
                    {matchedEvent ? (
                      <div className="text-sm text-muted-foreground">
                        Source: {matchedEvent.source}
                        {matchedEvent.payload?.amount && ` • Amount: ${matchedEvent.payload.amount}`}
                      </div>
                    ) : (
                      <div className="text-sm text-destructive font-medium">
                        Missing Event
                      </div>
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        </CardContent>
      </Card>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Transaction Details</CardTitle>
          </CardHeader>
          <CardContent>
            <dl className="space-y-4 text-sm">
              <div className="grid grid-cols-3 gap-4">
                <dt className="text-muted-foreground">Internal ID</dt>
                <dd className="col-span-2 font-mono break-all">{transaction.id}</dd>
              </div>
              <div className="grid grid-cols-3 gap-4">
                <dt className="text-muted-foreground">External ID</dt>
                <dd className="col-span-2 font-mono break-all">{transaction.external_id}</dd>
              </div>
              <div className="grid grid-cols-3 gap-4">
                <dt className="text-muted-foreground">Created At</dt>
                <dd className="col-span-2">{formatDate(transaction.created_at)}</dd>
              </div>
              <div className="grid grid-cols-3 gap-4">
                <dt className="text-muted-foreground">Updated At</dt>
                <dd className="col-span-2">{formatDate(transaction.updated_at)}</dd>
              </div>
            </dl>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader>
            <CardTitle>Reconciliation History</CardTitle>
          </CardHeader>
          <CardContent>
            {history.length === 0 ? (
              <div className="text-sm text-muted-foreground text-center py-4">No history available</div>
            ) : (
              <div className="space-y-4">
                {history.map(record => (
                  <div key={record.id} className="flex justify-between items-start text-sm border-b pb-4 last:border-0 last:pb-0">
                    <div className="space-y-1">
                      <StatusBadge status={record.status} />
                      <p className="text-muted-foreground">{record.message}</p>
                    </div>
                    <span className="text-xs text-muted-foreground whitespace-nowrap ml-4">
                      {formatDate(record.reconciled_at)}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Raw Event Data</CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="px-6 py-3">Sequence</TableHead>
                <TableHead className="px-6 py-3">Type</TableHead>
                <TableHead className="px-6 py-3">Source</TableHead>
                <TableHead className="px-6 py-3">Payload</TableHead>
                <TableHead className="px-6 py-3">Received</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {events.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} className="h-24 text-center text-muted-foreground">
                    No events received yet.
                  </TableCell>
                </TableRow>
              ) : (
                events.map(event => (
                  <TableRow key={event.id}>
                    <TableCell className="px-6 py-4">{event.sequence}</TableCell>
                    <TableCell className="px-6 py-4 font-medium">{event.event_type}</TableCell>
                    <TableCell className="px-6 py-4">{event.source}</TableCell>
                    <TableCell className="px-6 py-4 font-mono text-xs text-muted-foreground max-w-xs break-all">
                      {JSON.stringify(event.payload)}
                    </TableCell>
                    <TableCell className="px-6 py-4 text-muted-foreground">{formatDate(event.received_at)}</TableCell>
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
