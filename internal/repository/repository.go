package repository

import (
	"calendar/internal/models"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrTimeout = errors.New("timeout")
	ErrCancel  = errors.New("context cancelled")
)

type postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *postgres {
	return &postgres{
		pool: pool,
	}
}

func (p *postgres) CreateEvent(ctx context.Context, event models.Events) error {
	sql := `
	INSERT INTO events(user_id, event_date, event) 
	VALUES($1,$2,$3)
	`
	_, err := p.pool.Exec(ctx, sql, event.UserID, event.Date, event.Event)
	if errors.Is(err, context.Canceled) {
		return ErrCancel
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeout
	}
	return err
}
func (p *postgres) UpdateEvent(ctx context.Context, updEvent models.Events) error {
	sql := `
	UPDATE events
	SET event_date=$1, event=$2
	WHERE event_id=$3
	`
	_, err := p.pool.Exec(ctx, sql, updEvent.Date, updEvent.Event, updEvent.EventID)
	if errors.Is(err, context.Canceled) {
		return ErrCancel
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeout
	}
	return err
}
func (p *postgres) DeleteEvent(ctx context.Context, delEventID int) error {
	sql := `
	DELETE FROM events
	WHERE event_id=$1
	`
	_, err := p.pool.Exec(ctx, sql, delEventID)
	if errors.Is(err, context.Canceled) {
		return ErrCancel
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeout
	}
	return err
}
func (p *postgres) EventsForDay(ctx context.Context, userID int, date time.Time) ([]models.Events, error) {
	sql := `
	SELECT event_id, user_id, event_date, event
	FROM events
	WHERE user_id=$1 AND event_date=$2
	`
	rows, err := p.pool.Query(ctx, sql, userID, date)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, ErrCancel
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrTimeout
		}
		return nil, err
	}
	defer rows.Close()
	res := make([]models.Events, 0)
	for rows.Next() {
		var model models.Events
		err := rows.Scan(&model.EventID, &model.UserID, &model.Date, &model.Event)
		if err != nil {
			return nil, err
		}
		res = append(res, model)
	}
	return res, nil
}
func (p *postgres) EventsForWeek(ctx context.Context, userID int, date time.Time) ([]models.Events, error) {
	sql := `
	SELECT event_id, user_id, event_date, event
	FROM events
	WHERE user_id = $1
  		AND event_date >= date_trunc('week', $2)
  		AND event_date <  date_trunc('week', $2) + interval '1 week'
	`
	rows, err := p.pool.Query(ctx, sql, userID, date)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, ErrCancel
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrTimeout
		}
		return nil, err
	}
	defer rows.Close()
	res := make([]models.Events, 0)
	for rows.Next() {
		var model models.Events
		err := rows.Scan(&model.EventID, &model.UserID, &model.Date, &model.Event)
		if err != nil {
			return nil, err
		}
		res = append(res, model)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
func (p *postgres) EventsForMonth(ctx context.Context, userID int, date time.Time) ([]models.Events, error) {
	sql := `
	SELECT event_id, user_id, event_date, event
	FROM events
	WHERE user_id=$1 
		AND event_date >= date_trunc('month', $2)
		AND event_date < date_trunc('month', $2) + interval '1 month'
	`
	rows, err := p.pool.Query(ctx, sql, userID, date)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, ErrCancel
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrTimeout
		}
		return nil, err
	}
	defer rows.Close()
	res := make([]models.Events, 0)
	for rows.Next() {
		var model models.Events
		err := rows.Scan(&model.EventID, &model.UserID, &model.Date, &model.Event)
		if err != nil {
			return nil, err
		}
		res = append(res, model)
	}
	return res, nil
}
