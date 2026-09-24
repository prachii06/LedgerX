package services

import (
	"context"
	"github.com/prachii06/LedgerX/internal/models"
	"github.com/prachii06/LedgerX/internal/repository"
)

type EventService struct {
	repository  *repository.EventRepository
	broadcaster EventBroadcaster
}

func NewEventService(
	repository *repository.EventRepository,
	broadcaster EventBroadcaster,
) *EventService {
	return &EventService{
		repository:  repository,
		broadcaster: broadcaster,
	}
}

func (s *EventService) CreateEvent(
	ctx context.Context,
	event *models.Event,
) error {
	err := s.repository.Create(ctx, event)
	if err == nil && s.broadcaster != nil {
		s.broadcaster.BroadcastMessage("EVENT_PERSISTED", event.TransactionID, "", event)
	}
	return err
}

func (s *EventService) GetEvents(
	ctx context.Context,
	transactionID string,
	limit int,
	offset int,
) ([]models.Event, error) {
	return s.repository.FindAll(ctx, transactionID, limit, offset)
}
