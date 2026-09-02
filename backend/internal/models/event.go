package models

import "time"

type Event struct {
	ID            string                 `json:"id"`
	TransactionID string                 `json:"transaction_id"`
	Source        string                 `json:"source"`
	EventType     string                 `json:"event_type"`
	Sequence      int                    `json:"sequence"`
	Payload       map[string]interface{} `json:"payload"`
	ReceivedAt    time.Time              `json:"received_at"`
	CreatedAt     time.Time              `json:"created_at"`
}

const (
	SourceOrderService      = "ORDER_SERVICE"
	SourcePaymentGateway    = "PAYMENT_GATEWAY"
	SourceAccountingService = "ACCOUNTING_SERVICE"
)

const (
	EventTransactionCreated = "TRANSACTION_CREATED"
	EventPaymentReceived    = "PAYMENT_RECEIVED"
	EventAccountingBooked   = "ACCOUNTING_BOOKED"
)
