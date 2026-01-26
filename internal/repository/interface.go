package repository

import (
	"calendar/internal/dto"
	"calendar/internal/models"
	"context"
	"time"
)

type Repo interface {
	CreateEvent(ctx context.Context, event models.Events) error
	UpdateEvent(ctx context.Context, updEvent models.Events) error
	DeleteEvent(ctx context.Context, delEvent dto.RequestDTO) error
	EventsForDay(ctx context.Context, date time.Time) ([]models.Events, error)
	EventsForWeek(ctx context.Context, week int) ([]models.Events, error)
	EventsForMonth(ctx context.Context, month int) ([]models.Events, error)
}
