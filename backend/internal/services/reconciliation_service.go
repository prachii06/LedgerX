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

	eventCounts := make(map[string]int)

	for _, event := range events {
		eventCounts[event.EventType]++
	}

	expectedEvents := []string{
		models.EventTransactionCreated,
		models.EventPaymentReceived,
		models.EventAccountingBooked,
	}

	// Check for duplicate events.
	for _, eventType := range expectedEvents {
		if eventCounts[eventType] > 1 {
			return &models.ReconciliationResult{
				TransactionID: transactionID,
				Status:        models.ReconciliationDuplicate,
				Message: fmt.Sprintf(
					"Duplicate event detected: %s",
					eventType,
				),
			}, nil
		}
	}

	// Check for missing events.
	for _, eventType := range expectedEvents {
		if eventCounts[eventType] == 0 {
			return &models.ReconciliationResult{
				TransactionID: transactionID,
				Status:        models.ReconciliationMissing,
				Message: fmt.Sprintf(
					"Missing event: %s",
					eventType,
				),
			}, nil
		}
	}

	// All expected events exist exactly once.
	return &models.ReconciliationResult{
		TransactionID: transactionID,
		Status:        models.ReconciliationMatched,
		Message:       "All expected events received exactly once",
	}, nil
}