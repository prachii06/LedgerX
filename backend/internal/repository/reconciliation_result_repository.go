package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prachii06/LedgerX/internal/models"
)

type ReconciliationResultRepository struct {
	db *pgxpool.Pool
}

func NewReconciliationResultRepository(
	db *pgxpool.Pool,
) *ReconciliationResultRepository {
	return &ReconciliationResultRepository{
		db: db,
	}
}

func (r *ReconciliationResultRepository) Create(
	ctx context.Context,
	record *models.ReconciliationRecord,
) error {

	query := `
		INSERT INTO reconciliation_results (
			id,
			transaction_id,
			status,
			message,
			reconciled_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		record.ID,
		record.TransactionID,
		record.Status,
		record.Message,
		record.ReconciledAt,
	)

	return err
}