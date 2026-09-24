package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/prachii06/LedgerX/internal/models"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) Create(
	ctx context.Context,
	event *models.Event,
) error {

	query := `
		INSERT INTO transaction_events (
			id,
			transaction_id,
			source,
			event_type,
			sequence,
			payload,
			received_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		event.ID,
		event.TransactionID,
		event.Source,
		event.EventType,
		event.Sequence,
		event.Payload,
		event.ReceivedAt,
	)

	return err
}

func (r *EventRepository) FindAll(
	ctx context.Context,
	transactionID string,
	limit int,
	offset int,
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
		WHERE (CAST($1 AS text) = '' OR transaction_id = $1)
		ORDER BY received_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, transactionID, limit, offset)
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
