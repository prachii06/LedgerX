import { Transaction, Event, ReconciliationResult, ReconciliationRecord, SystemHealth } from "@/types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

async function apiFetch(endpoint: string, options: RequestInit = {}) {
  try {
    const res = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      cache: options.cache || "no-store",
    });

    if (!res.ok) {
      // Non-2xx HTTP response
      throw new Error(`API Error: ${res.status} ${res.statusText}`);
    }

    try {
      return await res.json();
    } catch (e) {
      // Invalid JSON response
      throw new Error("API Error: Invalid JSON response from server");
    }
  } catch (error) {
    if (error instanceof TypeError && error.message === "Failed to fetch") {
      // Network/backend unreachable
      throw new Error(`Network Error: Backend unreachable at ${API_BASE_URL}. Ensure the backend is running.`);
    }
    throw error;
  }
}

export async function fetchHealth(): Promise<SystemHealth> {
  try {
    const liveRes = await apiFetch("/live").catch(() => ({ status: "Unavailable" }));
    const readyRes = await apiFetch("/ready").catch(() => ({ status: "Unavailable" }));
    
    return {
      status: liveRes.status === "alive" ? "OK" : "Unavailable",
      postgres: readyRes.status === "ready" ? "OK" : "Unavailable",
      redis: "Unavailable", // Not exposed by backend
      kafka: "Unavailable", // Not exposed by backend
    };
  } catch (error) {
    return { status: "Unavailable", postgres: "Unavailable", redis: "Unavailable", kafka: "Unavailable" };
  }
}

export async function getOverview(): Promise<{ total_transactions: number, reconciled: number, issues: number, pending: number }> {
  return apiFetch("/dashboard/overview");
}

export async function fetchTransactions(limit = 50, offset = 0): Promise<Transaction[]> {
  return apiFetch(`/transactions?limit=${limit}&offset=${offset}`);
}

export async function fetchTransaction(id: string): Promise<Transaction> {
  return apiFetch(`/transactions/${id}`);
}

export async function fetchEvents(limit = 50, offset = 0, transactionId?: string): Promise<Event[]> {
  let url = `/events?limit=${limit}&offset=${offset}`;
  if (transactionId) {
    url += `&transaction_id=${transactionId}`;
  }
  return apiFetch(url);
}

export async function getReconciliation(transactionId: string): Promise<ReconciliationResult> {
  return apiFetch(`/reconcile/${transactionId}`);
}

export async function getReconciliationHistory(transactionId: string): Promise<ReconciliationRecord[]> {
  return apiFetch(`/reconcile/${transactionId}/history`);
}

export async function runSimulation(count: number = 1): Promise<{ generated_transactions: any[] }> {
  return apiFetch(`/simulate`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ count }),
  });
}
