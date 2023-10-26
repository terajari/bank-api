package dto

import (
	"github.com/terajari/bank-api/model"
)

type MakeTransferRequest struct {
	SenderId   string `json:"sender_id" binding:"required"`
	ReceiverId string `json:"receiver_id" binding:"required"`
	Amount     int64  `json:"amount" binding:"required,gt=0"`
	Currency   string `json:"currency" binding:"required,currency"`
}

type MakeTransferResponse struct {
	Transfer      model.Transfer `json:"transfer"`
	Sender        model.Accounts `json:"sender"`
	Receiver      model.Accounts `json:"receiver"`
	SenderEntry   model.Entries  `json:"sender_entry"`
	ReceiverEntry model.Entries  `json:"receiver_entry"`
}
