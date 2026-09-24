"use client";

import { fetchEvents } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { formatDate } from "@/lib/utils"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import useSWR from "swr"

export default function EventsPage() {
  const { data: eventsData } = useSWR(
    'events-stream',
    () => fetchEvents(300).catch(() => []),
    { revalidateOnFocus: false }
  )
  const events = eventsData || [];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Events Stream</h1>
          <p className="text-muted-foreground mt-1">Raw event ingestion from Kafka.</p>
        </div>
      </div>

      <Card>
        <CardHeader className="py-4 px-6 border-b">
          <CardTitle className="text-base font-medium">All Events</CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="px-6 py-3">Event ID</TableHead>
                <TableHead className="px-6 py-3">Transaction ID</TableHead>
                <TableHead className="px-6 py-3">Type</TableHead>
                <TableHead className="px-6 py-3">Source</TableHead>
                <TableHead className="px-6 py-3">Received At</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {events.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} className="h-24 text-center text-muted-foreground">
                    No events found.
                  </TableCell>
                </TableRow>
              ) : (
                events.map((e) => (
                  <TableRow key={e.id} className="group">
                    <TableCell className="px-6 py-4 font-medium font-mono text-xs text-muted-foreground max-w-[120px] truncate">
                      {e.id}
                    </TableCell>
                    <TableCell className="px-6 py-4 font-mono text-xs">
                      {e.transaction_id}
                    </TableCell>
                    <TableCell className="px-6 py-4 font-medium">
                      {e.event_type}
                    </TableCell>
                    <TableCell className="px-6 py-4">
                      {e.source}
                    </TableCell>
                    <TableCell className="px-6 py-4 text-muted-foreground">
                      {formatDate(e.received_at)}
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
