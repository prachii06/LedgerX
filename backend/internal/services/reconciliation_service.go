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

	expectedEvents := []string{
		models.EventTransactionCreated,
		models.EventPaymentReceived,
		models.EventAccountingBooked,
	}

	receivedEvents := make([]string, 0)

	eventCounts := make(map[string]int)
	eventAmounts := make(map[string]float64)

	for _, event := range events {

		eventType := event.EventType

		receivedEvents = append(
			receivedEvents,
			eventType,
		)

		eventCounts[eventType]++

		// Extract amount from event payload.
		if payloadAmount, ok := event.Payload["amount"]; ok {

			if amount, ok := payloadAmount.(float64); ok {
				eventAmounts[eventType] = amount
			}
		}
	}

	// Find missing events.
	missingEvents := make([]string, 0)

	for _, eventType := range expectedEvents {
		if eventCounts[eventType] == 0 {
			missingEvents = append(
				missingEvents,
				eventType,
			)
		}
	}

	// Find duplicate events.
	duplicateEvents := make([]string, 0)

	for _, eventType := range expectedEvents {
		if eventCounts[eventType] > 1 {
			duplicateEvents = append(
				duplicateEvents,
				eventType,
			)
		}
	}

	// Check duplicates.
	if len(duplicateEvents) > 0 {
		return &models.ReconciliationResult{
			TransactionID:     transactionID,
			Status:            models.ReconciliationDuplicate,
			Message:           fmt.Sprintf(
				"Duplicate events detected: %v",
				duplicateEvents,
			),
			ExpectedEvents:    expectedEvents,
			ReceivedEvents:    receivedEvents,
			DuplicateEvents:   duplicateEvents,
			TransactionAmount: transactionAmount,
			EventAmounts:      eventAmounts,
		}, nil
	}

	// Check missing events.
	if len(missingEvents) > 0 {
		return &models.ReconciliationResult{
			TransactionID:     transactionID,
			Status:            models.ReconciliationMissing,
			Message:           fmt.Sprintf(
				"Missing events: %v",
				missingEvents,
			),
			ExpectedEvents:    expectedEvents,
			ReceivedEvents:    receivedEvents,
			MissingEvents:     missingEvents,
			TransactionAmount: transactionAmount,
			EventAmounts:      eventAmounts,
		}, nil
	}

	// Check amount mismatches.
	for eventType, eventAmount := range eventAmounts {

		if eventAmount != transactionAmount {

			return &models.ReconciliationResult{
				TransactionID: transactionID,
				Status:        models.ReconciliationMismatch,
				Message: fmt.Sprintf(
					"Amount mismatch in %s: transaction=%.2f, event=%.2f",
					eventType,
					transactionAmount,
					eventAmount,
				),
				ExpectedEvents:    expectedEvents,
				ReceivedEvents:    receivedEvents,
				TransactionAmount: transactionAmount,
				EventAmounts:      eventAmounts,
			}, nil
		}
	}

	// Everything matches.
	return &models.ReconciliationResult{
		TransactionID:     transactionID,
		Status:            models.ReconciliationMatched,
		Message:           "All events received exactly once and amounts match",
		ExpectedEvents:    expectedEvents,
		ReceivedEvents:    receivedEvents,
		TransactionAmount: transactionAmount,
		EventAmounts:      eventAmounts,
	}, nil
}