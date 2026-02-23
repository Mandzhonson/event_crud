package dto

type GetDTO struct {
	Period string `form:"period"`
	UserID int    `form:"user_id"`
	Date   string `form:"event_date"`
}

type RequestDTO struct {
	EventID int    `json:"event_id"`
	UserID  int    `json:"user_id"`
	Date    string `json:"event_date"`
	Event   string `json:"event"`
}
