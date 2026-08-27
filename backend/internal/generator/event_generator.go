package generator

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/prachii06/LedgerX/internal/models"
)

type EventCreator interface {
	CreateEvent(ctx context.Context, event *models.Event) error
}

type EventGenerator struct {
	service EventCreator
}

func NewEventGenerator(
	service EventCreator,
) *EventGenerator {
	return &EventGenerator{
		service: service,
	}
}

func (g *EventGenerator) GenerateForTransaction(
	ctx context.Context,
	transactionID string,
	amount float64,
	currency string,
) error {

	accountingAmount := amount

	// Simulate an occasional accounting amount mismatch.
	if rand.Float64() < 0.2 {
		accountingAmount = amount + 500
	}

	events := []models.Event{
		{
			ID:            uuid.New().String(),
			TransactionID: transactionID,
			Source:        models.SourceOrderService,
			EventType:     models.EventTransactionCreated,
			Sequence:      1,
			Payload: map[string]interface{}{
				"transaction_id": transactionID,
				"amount":         amount,
				"currency":       currency,
			},
			ReceivedAt: time.Now(),
		},
		{
			ID:            uuid.New().String(),
			TransactionID: transactionID,
			Source:        models.SourcePaymentGateway,
			EventType:     models.EventPaymentReceived,
			Sequence:      2,
			Payload: map[string]interface{}{
				"transaction_id": transactionID,
				"amount":         amount,
				"currency":       currency,
				"payment_status": "SUCCESS",
			},
			ReceivedAt: time.Now(),
		},
		{
			ID:            uuid.New().String(),
			TransactionID: transactionID,
			Source:        models.SourceAccountingService,
			EventType:     models.EventAccountingBooked,
			Sequence:      3,
			Payload: map[string]interface{}{
				"transaction_id": transactionID,
				"amount":         accountingAmount,
				"currency":       currency,
				"ledger_status":  "BOOKED",
			},
			ReceivedAt: time.Now(),
		},
	}

	// Save the first event immediately.
	if err := g.service.CreateEvent(ctx, &events[0]); err != nil {
		return err
	}

	// Generate remaining events asynchronously.
	go g.generateRemainingEvents(
		transactionID,
		events[1:],
	)

	return nil
}


func (g *EventGenerator) generateRemainingEvents(
	transactionID string,
	events []models.Event,
) {
	ctx := context.Background()

	for _, event := range events {

		// Simulate delayed arrival.
		time.Sleep(4 * time.Second)

		// Simulate a missing accounting event.
		if event.EventType == models.EventAccountingBooked {
			if rand.Float64() < 0.3 {
				continue
			}
		}

		if err := g.service.CreateEvent(ctx, &event); err != nil {
			return
		}

		// Simulate duplicate payment event.
		if event.EventType == models.EventPaymentReceived {
			if rand.Float64() < 0.3 {
				_ = g.generateDuplicateEvent(ctx, event)
			}
		}
	}
}

func (g *EventGenerator) generateDuplicateEvent(
	ctx context.Context,
	event models.Event,
) error {

	event.ID = uuid.New().String()

	if err := g.service.CreateEvent(ctx, &event); err != nil {
		return err
	}

	return nil
}