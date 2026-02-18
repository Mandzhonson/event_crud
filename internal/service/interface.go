package service

import (
	"calendar/internal/dto"
	"calendar/internal/models"
	"context"
	"time"
)

type EventService interface {
	CreateEvent(ctx context.Context, event dto.RequestDTO) error
	UpdateEvent(ctx context.Context, eventDTO dto.RequestDTO) error
	DeleteEvent(ctx context.Context, delEventID int) error
	EventsForDay(ctx context.Context, userID int, date time.Time) ([]models.Events, error)
	EventsForWeek(ctx context.Context, userID int, date time.Time) ([]models.Events, error)
	EventsForMonth(ctx context.Context, userID int, date time.Time) ([]models.Events, error)
}
