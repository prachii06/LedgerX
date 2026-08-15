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

	// Get the original transaction amount.
	transactionAmount, err := s.repository.GetTransactionAmount(
		ctx,
		transactionID,
	)

	if err != nil {
		return nil, err
	}

	// Get all events for this transaction.
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

	// Check event amounts against the original transaction amount.
	for _, event := range events {

		payloadAmount, ok := event.Payload["amount"]

		if !ok {
			return &models.ReconciliationResult{
				TransactionID: transactionID,
				Status:        models.ReconciliationMismatch,
				Message: fmt.Sprintf(
					"Amount missing from event: %s",
					event.EventType,
				),
			}, nil
		}

		amount, ok := payloadAmount.(float64)

		if !ok {
			return &models.ReconciliationResult{
				TransactionID: transactionID,
				Status:        models.ReconciliationMismatch,
				Message: fmt.Sprintf(
					"Invalid amount in event: %s",
					event.EventType,
				),
			}, nil
		}

		if amount != transactionAmount {
			return &models.ReconciliationResult{
				TransactionID: transactionID,
				Status:        models.ReconciliationMismatch,
				Message: fmt.Sprintf(
					"Amount mismatch in %s: transaction=%.2f, event=%.2f",
					event.EventType,
					transactionAmount,
					amount,
				),
			}, nil
		}
	}

	// Everything matches.
	return &models.ReconciliationResult{
		TransactionID: transactionID,
		Status:        models.ReconciliationMatched,
		Message:       "All events received exactly once and amounts match",
	}, nil
}