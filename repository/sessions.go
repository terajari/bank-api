package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/terajari/bank-api/model"
)

type SessionsRepository interface {
	Create(ctx context.Context, sessions model.Sessions) (model.Sessions, error)
	Get(ctx context.Context, id uuid.UUID) (model.Sessions, error)
	Update(ctx context.Context, id uuid.UUID, isBlocked bool) error
	Last(ctx context.Context) (model.Sessions, error)
}

type sessionsRepository struct {
	db *sqlx.DB
}

func NewSessionsRepository(db *sqlx.DB) SessionsRepository {
	return &sessionsRepository{db}
}

func (s *sessionsRepository) Create(ctx context.Context, sessions model.Sessions) (model.Sessions, error) {
	query := "INSERT INTO sessions (id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at"
	row := s.db.QueryRowContext(ctx, query, sessions.Id, sessions.Username, sessions.RefreshToken, sessions.UserAgent, sessions.ClientIp, sessions.IsBlocked, sessions.ExpiresAt)
	var ss model.Sessions
	err := row.Scan(&ss.Id, &ss.Username, &ss.RefreshToken, &ss.UserAgent, &ss.ClientIp, &ss.IsBlocked, &ss.ExpiresAt, &ss.CreatedAt)
	if err != nil {
		return model.Sessions{}, err
	}
	return ss, nil
}

func (s *sessionsRepository) Get(ctx context.Context, id uuid.UUID) (model.Sessions, error) {
	query := "SELECT id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at FROM sessions WHERE id = $1 LIMIT 1"
	row := s.db.QueryRowContext(ctx, query, id)
	var ss model.Sessions
	err := row.Scan(&ss.Id, &ss.Username, &ss.RefreshToken, &ss.UserAgent, &ss.ClientIp, &ss.IsBlocked, &ss.ExpiresAt, &ss.CreatedAt)
	if err != nil {
		return model.Sessions{}, err
	}
	return ss, nil
}

func (s *sessionsRepository) Update(ctx context.Context, id uuid.UUID, isBlocked bool) error {
	query := "UPDATE sessions SET is_blocked = $2 WHERE id = $1"
	_, err := s.db.ExecContext(ctx, query, id, isBlocked)
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionsRepository) Last(ctx context.Context) (model.Sessions, error) {
	query := "SELECT id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at FROM sessions ORDER BY created_at DESC LIMIT 1"
	row := s.db.QueryRowContext(ctx, query)
	var ss model.Sessions
	err := row.Scan(&ss.Id, &ss.Username, &ss.RefreshToken, &ss.UserAgent, &ss.ClientIp, &ss.IsBlocked, &ss.ExpiresAt, &ss.CreatedAt)
	if err != nil {
		return model.Sessions{}, err
	}
	return ss, nil
}
