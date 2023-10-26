package dto

import (
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserReponse struct {
	Username     string    `json:"username"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PwdChangedAt time.Time `json:"pwd_changed_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	AccessToken          string      `json:"access_token"`
	AccessTokenExpiresAt time.Time   `json:"access_token_expires_at"`
	User                 UserReponse `json:"user"`
}
