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
			payload,
			received_at,
			created_at
		FROM transaction_events
		WHERE transaction_id = $1
		ORDER BY received_at ASC
	`

	rows, err := r.db.Query(ctx, query, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var event models.Event

		err := rows.Scan(
			&event.ID,
			&event.TransactionID,
			&event.Source,
			&event.EventType,
			&event.Payload,
			&event.ReceivedAt,
			&event.CreatedAt,
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