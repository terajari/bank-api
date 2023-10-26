package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/terajari/bank-api/dto"
	"github.com/terajari/bank-api/model"
	"github.com/terajari/bank-api/repository"
)

type SessionsUsecase interface {
	AddSessions(ctx context.Context, req dto.AddSessionsRequest) (dto.SessionResponse, error)
	GetSessions(ctx context.Context, id uuid.UUID) (dto.SessionResponse, error)
	UpdateBlockStatus(ctx context.Context, req dto.UpdateSessionBlockRequest) error
	LastSession(ctx context.Context) (dto.SessionResponse, error)
}

type sessionsUsecase struct {
	sessionsRepo repository.SessionsRepository
}

func NewSessionsUsecase(sessionsRepo repository.SessionsRepository) SessionsUsecase {
	return &sessionsUsecase{
		sessionsRepo: sessionsRepo,
	}
}

func (s *sessionsUsecase) AddSessions(ctx context.Context, req dto.AddSessionsRequest) (dto.SessionResponse, error) {
	ss, err := s.sessionsRepo.Create(ctx, model.Sessions{
		Id:           req.Id,
		Username:     req.Username,
		RefreshToken: req.RefreshToken,
		UserAgent:    req.UserAgent,
		ClientIp:     req.ClientIp,
		IsBlocked:    req.IsBlocked,
		ExpiresAt:    req.ExpiresAt,
	})
	if err != nil {
		return dto.SessionResponse{}, err
	}

	return dto.SessionResponse{
		Id:           ss.Id,
		Username:     ss.Username,
		RefreshToken: ss.RefreshToken,
		UserAgent:    ss.UserAgent,
		ClientIp:     ss.ClientIp,
		IsBlocked:    ss.IsBlocked,
		ExpiresAt:    ss.ExpiresAt,
		CreatedAt:    ss.CreatedAt,
	}, nil
}

func (s *sessionsUsecase) GetSessions(ctx context.Context, uuid uuid.UUID) (dto.SessionResponse, error) {
	ss, err := s.sessionsRepo.Get(ctx, uuid)
	if err != nil {
		return dto.SessionResponse{}, err
	}

	return dto.SessionResponse{
		Id:           ss.Id,
		Username:     ss.Username,
		RefreshToken: ss.RefreshToken,
		UserAgent:    ss.UserAgent,
		ClientIp:     ss.ClientIp,
		IsBlocked:    ss.IsBlocked,
		ExpiresAt:    ss.ExpiresAt,
		CreatedAt:    ss.CreatedAt,
	}, nil
}

func (s *sessionsUsecase) UpdateBlockStatus(ctx context.Context, req dto.UpdateSessionBlockRequest) error {
	err := s.sessionsRepo.Update(ctx, req.Id, req.IsBlocked)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionsUsecase) LastSession(ctx context.Context) (dto.SessionResponse, error) {
	ss, err := s.sessionsRepo.Last(ctx)
	if err != nil {
		return dto.SessionResponse{}, err
	}

	return dto.SessionResponse{
		Id:           ss.Id,
		Username:     ss.Username,
		RefreshToken: ss.RefreshToken,
		UserAgent:    ss.UserAgent,
		ClientIp:     ss.ClientIp,
		IsBlocked:    ss.IsBlocked,
		ExpiresAt:    ss.ExpiresAt,
		CreatedAt:    ss.CreatedAt,
	}, nil
}
