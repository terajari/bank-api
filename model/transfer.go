package model

type Transfer struct {
	ID         string `json:"id"`
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
	Amount     int64  `json:"amount"`
	CreatedAt  string `json:"created_at"`
}
