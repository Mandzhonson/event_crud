package service

import (
	"calendar/internal/models"
	"calendar/internal/repository"
	"context"
	"fmt"
	"time"
)

type eventService struct {
	Repo repository.Repo
}

func NewEventService(r repository.Repo) *eventService {
	return &eventService{
		Repo: r,
	}
}

func (evSer *eventService) CreateEvent(ctx context.Context, event models.Events) error {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := evSer.Repo.CreateEvent(dbCtx, event); err != nil {
		return fmt.Errorf("error creating event: %w", err)
	}
	return nil
}
func (evSer *eventService) UpdateEvent(ctx context.Context, updEvent models.Events) error {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := evSer.Repo.UpdateEvent(dbCtx, updEvent); err != nil {
		return fmt.Errorf("error updating event: %w", err)
	}
	return nil
}
func (evSer *eventService) DeleteEvent(ctx context.Context, delEventID int) error {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := evSer.Repo.DeleteEvent(dbCtx, delEventID); err != nil {
		return fmt.Errorf("error deleting event: %w", err)
	}
	return nil
}
func (evSer *eventService) EventsForDay(ctx context.Context, userID int, date time.Time) ([]models.Events, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	events, err := evSer.Repo.EventsForDay(dbCtx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("error EvenstForDay: %w", err)
	}
	return events, nil
}
func (evSer *eventService) EventsForWeek(ctx context.Context, userID int, date time.Time) ([]models.Events, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	events, err := evSer.Repo.EventsForWeek(dbCtx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("error EvenstForWeek: %w", err)
	}
	return events, nil

}
func (evSer *eventService) EventsForMonth(ctx context.Context, userID int, date time.Time) ([]models.Events, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	events, err := evSer.Repo.EventsForMonth(dbCtx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("error EvenstForMonth: %w", err)
	}
	return events, nil
}
