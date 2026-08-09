package generator

import (
	"context"
	"time"
	"github.com/google/uuid"
	"github.com/prachii06/LedgerX/internal/models"
	"github.com/prachii06/LedgerX/internal/services"
)

type EventGenerator struct {
	service *services.EventService
}

func NewEventGenerator(
	service *services.EventService,
) *EventGenerator {
	return &EventGenerator{
		service: service,
	}
}

func (g *EventGenerator) GenerateForTransaction(
	ctx context.Context,
	transactionID string,
) error {
	events := []models.Event{
		{
			ID:            uuid.New().String(),
			TransactionID: transactionID,
			Source:        models.SourceOrderService,
			EventType:     models.EventTransactionCreated,
			Payload:       map[string]interface{}{},
			ReceivedAt:    time.Now(),
		},
		{
			ID:            uuid.New().String(),
			TransactionID: transactionID,
			Source:        models.SourcePaymentGateway,
			EventType:     models.EventPaymentReceived,
			Payload:       map[string]interface{}{},
			ReceivedAt:    time.Now(),
		},
		{
			ID:            uuid.New().String(),
			TransactionID: transactionID,
			Source:        models.SourceAccountingService,
			EventType:     models.EventAccountingBooked,
			Payload:       map[string]interface{}{},
			ReceivedAt:    time.Now(),
		},
	}

	for _, event := range events {
		if err := g.service.CreateEvent(ctx, &event); err != nil {
			return err
		}
	}

	return nil
}