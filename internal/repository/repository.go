package repository

import (
	"calendar/internal/dto"
	"calendar/internal/models"
	"context"

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

func (p *postgres) CreateEvent(ctx context.Context, event models.Events) error {
	sql := `
	INSERT INTO events(user_id, date, event) 
	VALUES($1,$2,$3)
	`
	_, err := p.pool.Exec(ctx, sql, event.UserID, event.Date, event.Event)
	return err
}
func (p *postgres) UpdateEvent(ctx context.Context, updEvent models.Events) error {
	sql := `
	UPDATE events
	SET date=$1,event=$2
	WHERE user_id=$3
	`
	_, err := p.pool.Exec(ctx, sql, updEvent.Date, updEvent.Event, updEvent.UserID)
	return err
}
func (p *postgres) DeleteEvent(ctx context.Context, delEvent dto.RequestDTO) error {
	sql := `
	DELETE FROM events
	WHERE user_id=$1 AND date=$2
	`
	_, err := p.pool.Exec(ctx, sql, delEvent.UserID, delEvent.Date)
	return err
}
func (p *postgres) EventsForDay(ctx context.Context, Events dto.RequestDTO) ([]models.Events, error) {
	sql := `
	SELECT user_id, date, event
	FROM events
	WHERE user_id=$1 AND date=$2
	`
	rows, err := p.pool.Query(ctx, sql, Events.UserID, Events.Date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]models.Events, 0)
	for rows.Next() {
		var model models.Events
		err := rows.Scan(&model.UserID, &model.Date, &model.Event)
		if err != nil {
			return nil, err
		}
		res = append(res, model)
	}
	return res, nil
}
func (p *postgres) EventsForWeek(ctx context.Context, Events dto.RequestDTO) ([]models.Events, error) {
	sql := `
	SELECT user_id, date, event
	FROM events
	WHERE user_id=$1 AND EXTRACT(WEEK FROM $2)=EXTRACT(WEEK FROM date) AND EXTRACT(YEAR FROM date)=EXTRACT(YEAR FROM $3)
	`
	rows, err := p.pool.Query(ctx, sql, Events.UserID, Events.Date, Events.Date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]models.Events, 0)
	for rows.Next() {
		var model models.Events
		err := rows.Scan(&model.UserID, &model.Date, &model.Event)
		if err != nil {
			return nil, err
		}
		res = append(res, model)
	}
	return res, nil
}
func (p *postgres) EventsForMonth(ctx context.Context, Events dto.RequestDTO) ([]models.Events, error) {
	sql := `
	SELECT user_id, date, event
	FROM events
	WHERE user_id=$1 AND EXTRACT(MONTH FROM $2)=EXTRACT(MONTH FROM date) AND EXTRACT(YEAR FROM date)=EXTRACT(YEAR FROM $3)
	`
	rows, err := p.pool.Query(ctx, sql, Events.UserID, Events.Date, Events.Date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]models.Events, 0)
	for rows.Next() {
		var model models.Events
		err := rows.Scan(&model.UserID, &model.Date, &model.Event)
		if err != nil {
			return nil, err
		}
		res = append(res, model)
	}
	return res, nil
}
