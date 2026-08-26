package models

import "time"

type ReconciliationRecord struct {
	ID            string    `json:"id"`
	TransactionID string    `json:"transaction_id"`
	Status        string    `json:"status"`
	Message       string    `json:"message"`
	ReconciledAt  time.Time `json:"reconciled_at"`
}