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

func (r *TransactionRepository) FindByID(
	ctx context.Context,
	id string,
) (*models.Transaction, error) {

	query := `
	SELECT
		t.id,
		t.external_id,
		t.amount,
		t.currency,
		COALESCE(r.status, t.status) as status,
		t.created_at,
		t.updated_at
	FROM transactions t
	LEFT JOIN LATERAL (
		SELECT status FROM reconciliation_results r
		WHERE r.transaction_id = t.id
		ORDER BY reconciled_at DESC LIMIT 1
	) r ON true
	WHERE t.id=$1
	`

	transaction := &models.Transaction{}

	err := r.DB.QueryRow(
		ctx,
		query,
		id,
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

func (r *TransactionRepository) FindAll(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.Transaction, error) {

	query := `
	SELECT
		t.id,
		t.external_id,
		t.amount,
		t.currency,
		COALESCE(r.status, t.status) as status,
		t.created_at,
		t.updated_at
	FROM transactions t
	LEFT JOIN LATERAL (
		SELECT status FROM reconciliation_results r
		WHERE r.transaction_id = t.id
		ORDER BY reconciled_at DESC LIMIT 1
	) r ON true
	ORDER BY t.created_at DESC
	LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
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
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) GetOverview(ctx context.Context) (int, int, int, int, error) {
	query := `
	WITH latest_status AS (
		SELECT
			t.id,
			COALESCE(r.status, t.status) as status
		FROM transactions t
		LEFT JOIN LATERAL (
			SELECT status FROM reconciliation_results r
			WHERE r.transaction_id = t.id
			ORDER BY reconciled_at DESC LIMIT 1
		) r ON true
	)
	SELECT 
		COUNT(*) as total,
		COALESCE(SUM(CASE WHEN status = 'RECONCILED' OR status = 'MATCHED' THEN 1 ELSE 0 END), 0) as reconciled,
		COALESCE(SUM(CASE WHEN status IN ('MISSING', 'MISMATCH', 'CURRENCY_MISMATCH', 'DUPLICATE') THEN 1 ELSE 0 END), 0) as issues,
		COALESCE(SUM(CASE WHEN status IN ('PENDING', 'PROCESSING') THEN 1 ELSE 0 END), 0) as pending
	FROM latest_status
	`

	var total, reconciled, issues, pending int

	err := r.DB.QueryRow(ctx, query).Scan(&total, &reconciled, &issues, &pending)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return total, reconciled, issues, pending, nil
}
