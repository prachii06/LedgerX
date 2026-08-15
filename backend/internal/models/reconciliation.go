package models

type ReconciliationStatus string

const (
	ReconciliationMatched   ReconciliationStatus = "MATCHED"
	ReconciliationMissing   ReconciliationStatus = "MISSING"
	ReconciliationDuplicate ReconciliationStatus = "DUPLICATE"
	ReconciliationMismatch  ReconciliationStatus = "MISMATCH"
)

type ReconciliationResult struct {
	TransactionID string
	Status        ReconciliationStatus
	Message       string
}