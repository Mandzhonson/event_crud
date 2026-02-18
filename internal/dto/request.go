package dto

type GetDTO struct {
	Period string `json:"period"`
	UserID int    `json:"user_id"`
	Date   string `json:"event_date"`
}

type RequestDTO struct {
	EventID int    `json:"event_id"`
	UserID  int    `json:"user_id"`
	Date    string `json:"event_date"`
	Event   string `json:"event"`
}
