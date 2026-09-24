import { Transaction, Event, ReconciliationResult, ReconciliationRecord, SystemHealth } from "@/types";

const configuredApiUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
const API_BASE_URL = (configuredApiUrl && configuredApiUrl.trim() !== "") 
  ? configuredApiUrl.trim() 
  : (process.env.NODE_ENV === "production" ? "MISSING_BACKEND_URL" : "http://localhost:8080");

async function apiFetch(endpoint: string, options: RequestInit = {}) {
  if (API_BASE_URL === "MISSING_BACKEND_URL") {
    console.error("NEXT_PUBLIC_API_BASE_URL is missing or empty. Please configure it in Vercel.");
    throw new Error("NEXT_PUBLIC_API_BASE_URL is missing. Please configure it in Vercel to point to the Go backend.");
  }
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
    const readyRes = await apiFetch("/ready").catch(() => ({
      status: "Unavailable",
      services: { postgresql: "Unavailable", redis: "Unavailable", kafka: "Unavailable" }
    }));

    const services = readyRes.services || {};

    return {
      status: liveRes.status === "alive" ? "OK" : "Unavailable",
      postgres: services.postgresql === "ok" ? "OK" : "Unavailable",
      redis: services.redis === "ok" ? "OK" : "Unavailable",
      kafka: services.kafka === "ok" ? "OK" : "Unavailable",
    };
  } catch (error) {
    return { status: "Unavailable", postgres: "Unavailable", redis: "Unavailable", kafka: "Unavailable" };
  }
}

export async function getOverview(): Promise<{ total_transactions: number, reconciled: number, issues: number, pending: number }> {
  return apiFetch("/overview");
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
