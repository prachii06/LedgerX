package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/prachii06/LedgerX/internal/models"
)

type ReconciliationRepository struct {
	db *pgxpool.Pool
}

func NewReconciliationRepository(
	db *pgxpool.Pool,
) *ReconciliationRepository {
	return &ReconciliationRepository{
		db: db,
	}
}

// GetEventsByTransactionID returns all events belonging to a transaction.
func (r *ReconciliationRepository) GetEventsByTransactionID(
	ctx context.Context,
	transactionID string,
) ([]models.Event, error) {

	query := `
		SELECT
			id,
			transaction_id,
			source,
			event_type,
			sequence,
			payload,
			received_at
		FROM transaction_events
		WHERE transaction_id = $1
		ORDER BY received_at ASC
	`

	rows, err := r.db.Query(ctx, query, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]models.Event, 0)

	for rows.Next() {

		var event models.Event

		err := rows.Scan(
			&event.ID,
			&event.TransactionID,
			&event.Source,
			&event.EventType,
			&event.Sequence,
			&event.Payload,
			&event.ReceivedAt,
		)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// GetTransactionAmount returns the original transaction amount.
func (r *ReconciliationRepository) GetTransactionAmount(
	ctx context.Context,
	transactionID string,
) (float64, error) {

	query := `
		SELECT amount
		FROM transactions
		WHERE id = $1
	`

	var amount float64

	err := r.db.QueryRow(
		ctx,
		query,
		transactionID,
	).Scan(&amount)

	if err != nil {
		return 0, err
	}

	return amount, nil
}

// GetTransactionCurrency returns the original transaction currency.
func (r *ReconciliationRepository) GetTransactionCurrency(
	ctx context.Context,
	transactionID string,
) (string, error) {

	query := `
		SELECT currency
		FROM transactions
		WHERE id = $1
	`

	var currency string

	err := r.db.QueryRow(
		ctx,
		query,
		transactionID,
	).Scan(&currency)

	if err != nil {
		return "", err
	}

	return currency, nil
}

// GetTransactionsNeedingReconciliation returns transactions that
// have not been reconciled yet.
func (r *ReconciliationRepository) GetTransactionsNeedingReconciliation(
	ctx context.Context,
) ([]string, error) {

	query := `
		SELECT t.id
		FROM transactions t
		LEFT JOIN reconciliation_results rr
			ON rr.transaction_id = t.id
		WHERE rr.transaction_id IS NULL
		AND (
			SELECT MAX(te.received_at)
			FROM transaction_events te
			WHERE te.transaction_id = t.id
		) < NOW() - INTERVAL '10 seconds'
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactionIDs := make([]string, 0)

	for rows.Next() {

		var transactionID string

		if err := rows.Scan(&transactionID); err != nil {
			return nil, err
		}

		transactionIDs = append(
			transactionIDs,
			transactionID,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactionIDs, nil
}
