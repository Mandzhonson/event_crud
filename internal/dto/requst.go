package dto

import "time"

type RequestDTO struct {
	Period string    `json:"period"`
	UserID int       `json:"user_id"`
	Date   time.Time `json:"event_date"`
}
