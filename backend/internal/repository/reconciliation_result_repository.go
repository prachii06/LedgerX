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

func (r *ReconciliationResultRepository) GetByTransactionID(
	ctx context.Context,
	transactionID string,
) ([]models.ReconciliationRecord, error) {

	query := `
		SELECT
			id,
			transaction_id,
			status,
			message,
			reconciled_at
		FROM reconciliation_results
		WHERE transaction_id = $1
		ORDER BY reconciled_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		transactionID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	records := make([]models.ReconciliationRecord, 0)

	for rows.Next() {

		var record models.ReconciliationRecord

		err := rows.Scan(
			&record.ID,
			&record.TransactionID,
			&record.Status,
			&record.Message,
			&record.ReconciledAt,
		)

		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
