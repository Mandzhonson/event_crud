package models

import "time"

type Events struct {
	EventID int       `json:"event_id"`
	UserID  int       `json:"user_id"`
	Date    time.Time `json:"date"`
	Event   string    `json:"event"`
}
