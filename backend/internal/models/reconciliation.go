package models

type ReconciliationStatus string

const (
	ReconciliationMatched   ReconciliationStatus = "MATCHED"
	ReconciliationMissing   ReconciliationStatus = "MISSING"
	ReconciliationDuplicate ReconciliationStatus = "DUPLICATE"
	ReconciliationMismatch  ReconciliationStatus = "MISMATCH"
	ReconciliationCurrencyMismatch ReconciliationStatus = "CURRENCY_MISMATCH"
)

type ReconciliationResult struct {
	TransactionID string                 `json:"transaction_id"`
	Status        ReconciliationStatus   `json:"status"`
	Message       string                 `json:"message"`
	ExpectedEvents []string              `json:"expected_events"`
	ReceivedEvents []string              `json:"received_events"`
	MissingEvents  []string              `json:"missing_events,omitempty"`
	DuplicateEvents []string             `json:"duplicate_events,omitempty"`
	TransactionAmount float64            `json:"transaction_amount"`
	EventAmounts map[string]float64      `json:"event_amounts,omitempty"`
	TransactionCurrency string            `json:"transaction_currency"`
	EventCurrencies     map[string]string `json:"event_currencies,omitempty"`

}