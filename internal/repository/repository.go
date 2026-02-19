package repository

import (
	"calendar/internal/apperr"
	"calendar/internal/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *postgres {
	return &postgres{
		pool: pool,
	}
}

func (p *postgres) CreateEvent(ctx context.Context, event models.Events) (int, error) {
	sql := `
	INSERT INTO events(user_id, event_date, event) 
	VALUES($1,$2,$3)
	RETURNING event_id
	`
	var id int
	err := p.pool.QueryRow(ctx, sql, event.UserID, event.Date, event.Event).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert event: %w", err)
	}
	return id, nil
}
func (p *postgres) UpdateEvent(ctx context.Context, updEvent models.Events) error {
	sql := `
	UPDATE events
	SET event_date=$1, event=$2
	WHERE event_id=$3
	`
	_, err := p.pool.Exec(ctx, sql, updEvent.Date, updEvent.Event, updEvent.EventID)
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}
	return nil
}
func (p *postgres) DeleteEvent(ctx context.Context, delEventID int) error {
	sql := `
	DELETE FROM events
	WHERE event_id=$1
	`
	resDel, err := p.pool.Exec(ctx, sql, delEventID)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}
	if resDel.RowsAffected() != 1 {
		return fmt.Errorf("delete event: %w", apperr.EventNotFound)
	}
	return nil
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
			return nil, apperr.ErrCancel
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, apperr.ErrTimeout
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
		return nil, fmt.Errorf("events for week failed: %w", err)
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
			return nil, apperr.ErrCancel
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, apperr.ErrTimeout
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

func (p *postgres) FindEvents(ctx context.Context, eventID int) error {
	sql := `SELECT event_id FROM events WHERE event_id=$1`
	row := p.pool.QueryRow(ctx, sql, eventID)
	var id int
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("find events: %w", err)
	}
	return nil
}
