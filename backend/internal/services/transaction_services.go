package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/prachii06/LedgerX/internal/models"
	"github.com/prachii06/LedgerX/internal/repository"
	"time"
)

type TransactionService struct {
	repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{
		repo: repo,
	}
}

func (s *TransactionService) CreateTransaction(
	ctx context.Context,
	transaction *models.Transaction,
) error {

	if transaction.ExternalID == "" {
		return errors.New("external_id is required")
	}

	if transaction.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if transaction.Currency == "" {
		return errors.New("currency is required")
	}

	existing, err := s.repo.FindByExternalID(ctx, transaction.ExternalID)

	if err == nil && existing != nil {
		return errors.New("transaction already exists")
	}

	transaction.ID = uuid.New().String()
	transaction.Status = models.StatusPending
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	return s.repo.Create(ctx, transaction)
}

func (s *TransactionService) GetTransactions(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.Transaction, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *TransactionService) GetTransaction(
	ctx context.Context,
	id string,
) (*models.Transaction, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *TransactionService) GetOverview(ctx context.Context) (int, int, int, int, error) {
	return s.repo.GetOverview(ctx)
}
