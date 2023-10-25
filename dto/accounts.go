package dto

import "time"

type RegisterNewAccountsRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

type RegisterNewAccountsResponse struct {
	Id        string    `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type GetAccountRequest struct {
	Id string `uri:"id"`
}

type GetAccountResponse struct {
	Id       string `json:"id"`
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type ListAccountsRequest struct {
	Owner string `json:"owner"`
	Page  int    `json:"limit"`
	Size  int    `json:"offset"`
}

type UpdateAccountRequest struct {
	Id      string `json:"id"`
	Balance int64  `json:"balance"`
}

type UpdateAccountResponse struct {
	Id       string `json:"id"`
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}
