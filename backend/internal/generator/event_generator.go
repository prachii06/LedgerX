package generator

import (
	"context"
	"time"
	"math/rand"
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

	events := []models.Event{
		{
			ID:            uuid.New().String(),
			TransactionID: transactionID,
			Source:        models.SourceOrderService,
			EventType:     models.EventTransactionCreated,
			Payload: map[string]interface{}{
				"transaction_id": transactionID,
				"amount":         amount,
				"currency":       currency,
				"sequence":       1,
			},
			ReceivedAt: time.Now(),
		},
		{
			ID:            uuid.New().String(),
			TransactionID: transactionID,
			Source:        models.SourcePaymentGateway,
			EventType:     models.EventPaymentReceived,
			Payload: map[string]interface{}{
				"transaction_id": transactionID,
				"amount":         amount,
				"currency":       currency,
				"payment_status": "SUCCESS",
				"sequence":       2,
			},
			ReceivedAt: time.Now(),
		},
		{
			ID:            uuid.New().String(),
			TransactionID: transactionID,
			Source:        models.SourceAccountingService,
			EventType:     models.EventAccountingBooked,
			Payload: map[string]interface{}{
				"transaction_id": transactionID,
				"amount":         amount,
				"currency":       currency,
				"ledger_status":  "BOOKED",
				"sequence":       3,
			},
			ReceivedAt: time.Now(),
		},
	}

	for i, event := range events {
	if err := g.service.CreateEvent(ctx, &event); err != nil {
		return err
	}

	// Simulate a duplicate event for EventPaymentReceived with a 30% chance
	if event.EventType == models.EventPaymentReceived {
	if rand.Float64() < 0.3 {
		if err := g.generateDuplicateEvent(ctx, event); err != nil {
			return err
		}
	}
}

	if i < len(events)-1 {
		time.Sleep(1 * time.Second)
	}
}


	

	return nil
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