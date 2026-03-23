package service

import (
	"calendar/internal/dto"
	"calendar/internal/models"
	"context"
)

//go:generate mockgen -source=interface.go -destination=mocks/mock_service.go
type EventService interface {
	CreateEvent(ctx context.Context, event dto.RequestDTO) (int, error)
	UpdateEvent(ctx context.Context, eventDTO dto.RequestDTO) error
	DeleteEvent(ctx context.Context, delEventID int) error
	EventsForDay(ctx context.Context, userID int, date string) ([]models.Events, error)
	EventsForWeek(ctx context.Context, userID int, date string) ([]models.Events, error)
	EventsForMonth(ctx context.Context, userID int, date string) ([]models.Events, error)
}
