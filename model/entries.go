package model

type Entries struct {
	ID        string `json:"id"`
	AccountId string `json:"account_id"`
	Amount    int64  `json:"amount"`
	CreatedAt string `json:"created_at"`
}
