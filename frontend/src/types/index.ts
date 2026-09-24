export type TransactionStatus = "PENDING" | "PROCESSING" | "MATCHED" | "MISSING" | "MISMATCH" | "CURRENCY_MISMATCH" | "OUT_OF_ORDER" | "DUPLICATE";

export interface Transaction {
  id: string;
  external_id: string;
  amount: number;
  currency: string;
  status: TransactionStatus;
  created_at: string;
  updated_at: string;
}

export interface Event {
  id: string;
  transaction_id: string;
  source: string;
  event_type: string;
  sequence: number;
  payload: Record<string, any>;
  received_at: string;
  created_at: string;
}

export interface ReconciliationResult {
  transaction_id: string;
  status: TransactionStatus;
  message: string;
  expected_events: string[];
  received_events: string[];
  missing_events?: string[];
  duplicate_events?: string[];
  transaction_amount: number;
  event_amounts?: Record<string, number>;
  transaction_currency: string;
  event_currencies?: Record<string, string>;
  expected_sequence?: number[];
  actual_sequence?: number[];
}

export interface ReconciliationRecord {
  id: string;
  transaction_id: string;
  status: TransactionStatus;
  message: string;
  reconciled_at: string;
}

export interface SystemHealth {
  status: string;
  postgres: string;
  redis: string;
  kafka: string;
}
