package services

import (
	"context"
	"github.com/prachii06/LedgerX/internal/models"
	"github.com/prachii06/LedgerX/internal/repository"
)

type EventService struct {
	repository *repository.EventRepository
}

func NewEventService(
	repository *repository.EventRepository,
) *EventService {
	return &EventService{
		repository: repository,
	}
}

func (s *EventService) CreateEvent(
	ctx context.Context,
	event *models.Event,
) error {
	return s.repository.Create(ctx, event)
}