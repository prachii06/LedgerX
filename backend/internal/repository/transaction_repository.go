package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prachii06/LedgerX/internal/models"
)

type TransactionRepository struct {
	DB *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		DB: db,
	}
}

func (r *TransactionRepository) Create(
	ctx context.Context,
	transaction *models.Transaction,
) error {

	query := `
	INSERT INTO transactions
	(
		id,
		external_id,
		amount,
		currency,
		status,
		created_at,
		updated_at
	)
	VALUES
	($1,$2,$3,$4,$5,$6,$7)
	`

	_, err := r.DB.Exec(
		ctx,
		query,
		transaction.ID,
		transaction.ExternalID,
		transaction.Amount,
		transaction.Currency,
		transaction.Status,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)

	return err
}

func (r *TransactionRepository) FindByExternalID(
	ctx context.Context,
	externalID string,
) (*models.Transaction, error) {

	query := `
	SELECT
		id,
		external_id,
		amount,
		currency,
		status,
		created_at,
		updated_at
	FROM transactions
	WHERE external_id=$1
	`

	transaction := &models.Transaction{}

	err := r.DB.QueryRow(
		ctx,
		query,
		externalID,
	).Scan(
		&transaction.ID,
		&transaction.ExternalID,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}
