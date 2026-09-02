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
