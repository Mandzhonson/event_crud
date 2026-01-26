package models

import "time"

type Events struct {
	UserID int       `json:"user_id"`
	Date   time.Time `json:"date"`
	Event  string    `json:"event"`
}
