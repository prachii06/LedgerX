package services

import (
	"context"
	"fmt"

	"github.com/prachii06/LedgerX/internal/models"
	"github.com/prachii06/LedgerX/internal/repository"
)

type ReconciliationService struct {
	repository *repository.ReconciliationRepository
}

func NewReconciliationService(
	repository *repository.ReconciliationRepository,
) *ReconciliationService {
	return &ReconciliationService{
		repository: repository,
	}
}

func (s *ReconciliationService) ReconcileTransaction(
	ctx context.Context,
	transactionID string,
) (*models.ReconciliationResult, error) {

	events, err := s.repository.GetEventsByTransactionID(
		ctx,
		transactionID,
	)

	if err != nil {
		return nil, err
	}

	eventCount := len(events)

	switch {
	case eventCount == 3:
		return &models.ReconciliationResult{
			TransactionID: transactionID,
			Status:        models.ReconciliationMatched,
			Message:       "All expected events received",
		}, nil

	case eventCount < 3:
		return &models.ReconciliationResult{
			TransactionID: transactionID,
			Status:        models.ReconciliationMissing,
			Message: fmt.Sprintf(
				"Expected 3 events but received %d",
				eventCount,
			),
		}, nil

	default:
		return &models.ReconciliationResult{
			TransactionID: transactionID,
			Status:        models.ReconciliationDuplicate,
			Message: fmt.Sprintf(
				"Expected 3 events but received %d",
				eventCount,
			),
		}, nil
	}
}