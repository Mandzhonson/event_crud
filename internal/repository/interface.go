package repository

import (
	"calendar/internal/models"
	"context"
	"time"
)

type Repo interface {
	CreateEvent(ctx context.Context, event models.Events) error
	UpdateEvent(ctx context.Context, updEvent models.Events) error
	DeleteEvent(ctx context.Context, delEventID int) error
	EventsForDay(ctx context.Context, userID int, date time.Time) ([]models.Events, error)
	EventsForWeek(ctx context.Context, userID int, date time.Time) ([]models.Events, error)
	EventsForMonth(ctx context.Context, userID int, date time.Time) ([]models.Events, error)
}
