package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/prachii06/LedgerX/internal/cache"
	"github.com/prachii06/LedgerX/internal/metrics"
	"github.com/prachii06/LedgerX/internal/models"
	"github.com/prachii06/LedgerX/internal/repository"
)

type ReconciliationService struct {
	repository       *repository.ReconciliationRepository
	resultRepository *repository.ReconciliationResultRepository
	cache            *cache.ReconciliationCache
	broadcaster      EventBroadcaster
}

func NewReconciliationService(
	repository *repository.ReconciliationRepository,
	resultRepository *repository.ReconciliationResultRepository,
	reconciliationCache *cache.ReconciliationCache,
	broadcaster EventBroadcaster,
) *ReconciliationService {
	return &ReconciliationService{
		repository:       repository,
		resultRepository: resultRepository,
		cache:            reconciliationCache,
		broadcaster:      broadcaster,
	}
}

func (s *ReconciliationService) ReconcileTransaction(
	ctx context.Context,
	transactionID string,
) (*models.ReconciliationResult, error) {

	// Get original transaction amount.
	transactionAmount, err := s.repository.GetTransactionAmount(
		ctx,
		transactionID,
	)
	if err != nil {
		return nil, err
	}

	// Get original transaction currency.
	transactionCurrency, err := s.repository.GetTransactionCurrency(
		ctx,
		transactionID,
	)
	if err != nil {
		return nil, err
	}

	// Get all events belonging to the transaction.
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
	eventCurrencies := make(map[string]string)

	// Process events.
	for _, event := range events {

		eventType := event.EventType

		receivedEvents = append(
			receivedEvents,
			eventType,
		)

		eventCounts[eventType]++

		// Extract amount.
		if payloadAmount, ok := event.Payload["amount"]; ok {

			switch amount := payloadAmount.(type) {

			case float64:
				eventAmounts[eventType] = amount

			case float32:
				eventAmounts[eventType] = float64(amount)

			case int:
				eventAmounts[eventType] = float64(amount)

			case int32:
				eventAmounts[eventType] = float64(amount)

			case int64:
				eventAmounts[eventType] = float64(amount)

			case json.Number:
				value, err := amount.Float64()
				if err == nil {
					eventAmounts[eventType] = value
				}
			}
		}

		// Extract currency.
		if payloadCurrency, ok := event.Payload["currency"]; ok {

			if currency, ok := payloadCurrency.(string); ok {
				eventCurrencies[eventType] = currency
			}
		}
	}

	// Expected sequence.
	expectedSequence := []int{1, 2, 3}

	actualSequence := make([]int, 0, len(events))

	for _, event := range events {
		actualSequence = append(
			actualSequence,
			event.Sequence,
		)
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

	// --------------------------------------------------
	// DUPLICATE
	// --------------------------------------------------

	if len(duplicateEvents) > 0 {

		result := &models.ReconciliationResult{
			TransactionID: transactionID,
			Status:        models.ReconciliationDuplicate,
			Message: fmt.Sprintf(
				"Duplicate events detected: %v",
				duplicateEvents,
			),
			ExpectedEvents:      expectedEvents,
			ReceivedEvents:      receivedEvents,
			DuplicateEvents:     duplicateEvents,
			TransactionAmount:   transactionAmount,
			EventAmounts:        eventAmounts,
			TransactionCurrency: transactionCurrency,
			EventCurrencies:     eventCurrencies,
			ExpectedSequence:    expectedSequence,
			ActualSequence:      actualSequence,
		}

		return s.saveAndReturn(ctx, result)
	}

	// --------------------------------------------------
	// MISSING
	// --------------------------------------------------

	if len(missingEvents) > 0 {

		result := &models.ReconciliationResult{
			TransactionID: transactionID,
			Status:        models.ReconciliationMissing,
			Message: fmt.Sprintf(
				"Missing events: %v",
				missingEvents,
			),
			ExpectedEvents:      expectedEvents,
			ReceivedEvents:      receivedEvents,
			MissingEvents:       missingEvents,
			TransactionAmount:   transactionAmount,
			EventAmounts:        eventAmounts,
			TransactionCurrency: transactionCurrency,
			EventCurrencies:     eventCurrencies,
			ExpectedSequence:    expectedSequence,
			ActualSequence:      actualSequence,
		}

		return s.saveAndReturn(ctx, result)
	}

	// --------------------------------------------------
	// OUT OF ORDER
	// --------------------------------------------------

	if len(actualSequence) == len(expectedSequence) {

		for i := range expectedSequence {

			if actualSequence[i] != expectedSequence[i] {

				result := &models.ReconciliationResult{
					TransactionID: transactionID,
					Status:        models.ReconciliationOutOfOrder,
					Message: fmt.Sprintf(
						"Events received out of order: expected sequence %v, received sequence %v",
						expectedSequence,
						actualSequence,
					),
					ExpectedEvents:      expectedEvents,
					ReceivedEvents:      receivedEvents,
					TransactionAmount:   transactionAmount,
					EventAmounts:        eventAmounts,
					TransactionCurrency: transactionCurrency,
					EventCurrencies:     eventCurrencies,
					ExpectedSequence:    expectedSequence,
					ActualSequence:      actualSequence,
				}

				return s.saveAndReturn(ctx, result)
			}
		}
	}

	// --------------------------------------------------
	// CURRENCY MISMATCH
	// --------------------------------------------------

	for eventType, eventCurrency := range eventCurrencies {

		if eventCurrency != transactionCurrency {

			result := &models.ReconciliationResult{
				TransactionID: transactionID,
				Status:        models.ReconciliationCurrencyMismatch,
				Message: fmt.Sprintf(
					"Currency mismatch in %s: transaction=%s, event=%s",
					eventType,
					transactionCurrency,
					eventCurrency,
				),
				ExpectedEvents:      expectedEvents,
				ReceivedEvents:      receivedEvents,
				TransactionAmount:   transactionAmount,
				EventAmounts:        eventAmounts,
				TransactionCurrency: transactionCurrency,
				EventCurrencies:     eventCurrencies,
				ExpectedSequence:    expectedSequence,
				ActualSequence:      actualSequence,
			}

			return s.saveAndReturn(ctx, result)
		}
	}

	// --------------------------------------------------
	// AMOUNT MISMATCH
	// --------------------------------------------------

	for eventType, eventAmount := range eventAmounts {

		if eventAmount != transactionAmount {

			result := &models.ReconciliationResult{
				TransactionID: transactionID,
				Status:        models.ReconciliationMismatch,
				Message: fmt.Sprintf(
					"Amount mismatch in %s: transaction=%.2f, event=%.2f",
					eventType,
					transactionAmount,
					eventAmount,
				),
				ExpectedEvents:      expectedEvents,
				ReceivedEvents:      receivedEvents,
				TransactionAmount:   transactionAmount,
				EventAmounts:        eventAmounts,
				TransactionCurrency: transactionCurrency,
				EventCurrencies:     eventCurrencies,
				ExpectedSequence:    expectedSequence,
				ActualSequence:      actualSequence,
			}

			return s.saveAndReturn(ctx, result)
		}
	}

	// --------------------------------------------------
	// MATCHED
	// --------------------------------------------------

	result := &models.ReconciliationResult{
		TransactionID:       transactionID,
		Status:              models.ReconciliationMatched,
		Message:             "All events received exactly once and amounts and currencies match",
		ExpectedEvents:      expectedEvents,
		ReceivedEvents:      receivedEvents,
		TransactionAmount:   transactionAmount,
		EventAmounts:        eventAmounts,
		TransactionCurrency: transactionCurrency,
		EventCurrencies:     eventCurrencies,
		ExpectedSequence:    expectedSequence,
		ActualSequence:      actualSequence,
	}

	return s.saveAndReturn(ctx, result)
}

// saveAndReturn saves the reconciliation attempt
// and caches the detailed result in Redis.
func (s *ReconciliationService) saveAndReturn(
	ctx context.Context,
	result *models.ReconciliationResult,
) (*models.ReconciliationResult, error) {

	// Save reconciliation history to PostgreSQL.
	if err := s.saveResult(ctx, result); err != nil {
		return nil, err
	}

	metrics.ReconciliationsTotal.WithLabelValues(string(result.Status)).Inc()

	// Store the detailed reconciliation result in Redis.
	if err := s.cache.Set(ctx, result); err != nil {
		return nil, err
	}

	if s.broadcaster != nil {
		s.broadcaster.BroadcastMessage("RECONCILIATION_COMPLETED", result.TransactionID, string(result.Status), result)
	}

	return result, nil
}

// saveResult persists the reconciliation result.
func (s *ReconciliationService) saveResult(
	ctx context.Context,
	result *models.ReconciliationResult,
) error {

	record := &models.ReconciliationRecord{
		ID:            uuid.New().String(),
		TransactionID: result.TransactionID,
		Status:        string(result.Status),
		Message:       result.Message,
		ReconciledAt:  time.Now(),
	}

	return s.resultRepository.Create(ctx, record)
}

func (s *ReconciliationService) GetReconciliationHistory(
	ctx context.Context,
	transactionID string,
) ([]models.ReconciliationRecord, error) {

	return s.resultRepository.GetByTransactionID(
		ctx,
		transactionID,
	)
}

func (s *ReconciliationService) GetCachedReconciliation(
	ctx context.Context,
	transactionID string,
) (*models.ReconciliationResult, error) {

	// 1. Try Redis first.
	cachedResult, err := s.cache.Get(
		ctx,
		transactionID,
	)

	if err != nil {
		return nil, err
	}

	// Cache HIT.
	if cachedResult != nil {
		return cachedResult, nil
	}

	// Cache MISS.
	// ReconcileTransaction will:
	// 1. Read events from PostgreSQL
	// 2. Calculate the result
	// 3. Save the reconciliation history
	// 4. Store the detailed result in Redis
	return s.ReconcileTransaction(
		ctx,
		transactionID,
	)
}
