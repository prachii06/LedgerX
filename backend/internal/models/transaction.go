package models

import "time"

type Transaction struct {
	ID         string    `json:"id"`          //internal ledgerx id
	ExternalID string    `json:"external_id"` //id from external systems
	Amount     float64   `json:"amount"`
	Currency   string    `json:"currency"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
