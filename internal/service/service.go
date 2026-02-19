package service

import (
	"calendar/internal/apperr"
	"calendar/internal/dto"
	"calendar/internal/models"
	"calendar/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
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

func (evSer *eventService) CreateEvent(ctx context.Context, eventDTO dto.RequestDTO) (int, error) {
	if eventDTO.Date == "" || eventDTO.Event == "" || eventDTO.UserID <= 0 {
		return 0, apperr.InvalidReqParams
	}
	date, err := time.Parse("2006-01-02", eventDTO.Date)
	if err != nil {
		return 0, apperr.InvalidReqParams
	}
	event := models.Events{UserID: eventDTO.UserID, Event: eventDTO.Event, Date: date}

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	id, err := evSer.Repo.CreateEvent(dbCtx, event)
	if err != nil {
		slog.Error("failed to insert event", slog.Any("error", err))
		if errors.Is(err, context.Canceled) {
			return 0, apperr.ErrCancel
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, apperr.ErrTimeout
		}
		slog.Error("failed to update event", slog.Any("error", err))
		return 0, apperr.InternalServErr
	}
	slog.Debug("create event is successfull", slog.Int("event_id", event.EventID))
	return id, nil
}

func (evSer *eventService) UpdateEvent(ctx context.Context, eventDTO dto.RequestDTO) error {
	if eventDTO.Date == "" && eventDTO.Event == "" {
		return apperr.InternalServErr
	}
	err := evSer.Repo.FindEvents(ctx, eventDTO.EventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperr.EventNotFound
		}
		return apperr.InternalServErr
	}
	date, err := time.Parse("2006-01-02", eventDTO.Date)
	if err != nil {
		return apperr.InvalidReqParams
	}
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	updEvent := models.Events{EventID: eventDTO.EventID, Event: eventDTO.Event, Date: date}
	if err := evSer.Repo.UpdateEvent(dbCtx, updEvent); err != nil {
		if errors.Is(err, context.Canceled) {
			return apperr.ErrCancel
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return apperr.ErrTimeout
		}
		slog.Error("failed to update event", slog.Any("error", err))
		return apperr.InternalServErr
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
