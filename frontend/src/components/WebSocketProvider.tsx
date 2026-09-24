"use client";

import { useEffect, useRef } from "react";
import { useSWRConfig } from "swr";
// import { toast } from "sonner"; // If they have sonner, else console.log

const configuredWsUrl = process.env.NEXT_PUBLIC_WS_URL;
const WS_URL = (configuredWsUrl && configuredWsUrl.trim() !== "")
  ? configuredWsUrl.trim()
  : (process.env.NODE_ENV === "production" ? "MISSING_WS_URL" : "ws://localhost:8080/ws");

export function WebSocketProvider({ children }: { children: React.ReactNode }) {
  const { mutate } = useSWRConfig();
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    let reconnectTimeout: NodeJS.Timeout;

    const connect = () => {
      if (WS_URL === "MISSING_WS_URL") {
        console.error("NEXT_PUBLIC_WS_URL is missing. Please configure it in Vercel to point to the Go backend WebSocket.");
        return;
      }
      const ws = new WebSocket(WS_URL);
      wsRef.current = ws;

      ws.onopen = () => {
        console.log("WebSocket connected to LedgerX backend");
      };

      ws.onmessage = (event) => {
        try {
          const msg = JSON.parse(event.data);
          
          if (msg.type === "TRANSACTION_CREATED") {
            const tx = msg.details;
            // Update transactions list
            mutate("transactions-list", (data: any = []) => {
              const exists = data.find((t: any) => t.id === tx.id);
              if (exists) return data;
              return [tx, ...data].slice(0, 200);
            }, false);

            // Update overview transactions
            mutate("overview-transactions", (data: any = []) => {
              const exists = data.find((t: any) => t.id === tx.id);
              if (exists) return data;
              return [tx, ...data].slice(0, 10);
            }, false);

            // Update reconciliation transactions
            mutate("transactions-recon", (data: any = []) => {
              const exists = data.find((t: any) => t.id === tx.id);
              if (exists) return data;
              return [tx, ...data].slice(0, 200);
            }, false);

            // Update stats
            mutate("overview-stats", (stats: any) => {
              if (!stats) return stats;
              return { ...stats, total_transactions: (stats.total_transactions || 0) + 1, pending: (stats.pending || 0) + 1 };
            }, false);
          }
          
          if (msg.type === "EVENT_PERSISTED") {
            const ev = msg.details;
            // Update events stream
            mutate("events-stream", (data: any = []) => {
              if (data.some((e: any) => e.id === ev.id)) {
                return data; // Prevent duplicate keys
              }
              return [ev, ...data].slice(0, 300);
            }, false);
          }

          if (msg.type === "RECONCILIATION_COMPLETED") {
            const result = msg.details;
            
            const updateTxStatus = (data: any = []) => {
              return data.map((t: any) => {
                if (t.id === result.transaction_id) {
                  return { ...t, status: result.status, updated_at: new Date().toISOString() };
                }
                return t;
              });
            };

            // Update status in all transaction lists
            mutate("transactions-list", updateTxStatus, false);
            mutate("overview-transactions", updateTxStatus, false);
            mutate("transactions-recon", updateTxStatus, false);

            // Fetch latest stats since reconciliation changes them dynamically 
            // OR calculate the delta. For exact accuracy we can trigger a re-fetch of just the stats:
            mutate("overview-stats");
          }

        } catch (e) {
          console.error("Failed to parse WS message", e);
        }
      };

      ws.onclose = () => {
        console.log("WebSocket disconnected, reconnecting in 5s...");
        reconnectTimeout = setTimeout(connect, 5000);
      };

      ws.onerror = (err) => {
        console.warn("WebSocket error (transient)", err);
        ws.close();
      };
    };

    connect();

    return () => {
      clearTimeout(reconnectTimeout);
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [mutate]);

  return <>{children}</>;
}
