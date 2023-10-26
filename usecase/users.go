package usecase

import (
	"context"
	"errors"

	"github.com/terajari/bank-api/dto"
	"github.com/terajari/bank-api/model"
	"github.com/terajari/bank-api/repository"
	"github.com/terajari/bank-api/utils"
)

type UsersUsecase interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserReponse, error)
	Login(ctx context.Context, req dto.LoginUserRequest) (model.Users, error)
}

type usersUsecase struct {
	repo repository.UsersRepository
}

func NewUsersUsecase(repo repository.UsersRepository) UsersUsecase {
	return &usersUsecase{
		repo: repo,
	}
}

func (u *usersUsecase) CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserReponse, error) {
	hashedPwd, err := utils.HashPasswrod(req.Password)
	if err != nil {
		return dto.UserReponse{}, err
	}

	user, err := u.repo.Create(ctx, model.Users{
		Username:       req.Username,
		HashedPassword: hashedPwd,
		FullName:       req.FullName,
		Email:          req.Email,
	})
	if err != nil {
		return dto.UserReponse{}, err
	}
	return dto.UserReponse{
		Username:     user.Username,
		FullName:     user.FullName,
		Email:        user.Email,
		PwdChangedAt: user.PasswordChangedAt,
		CreatedAt:    user.CreatedAt,
	}, nil
}

func (u *usersUsecase) Login(ctx context.Context, req dto.LoginUserRequest) (model.Users, error) {
	user, err := u.repo.Get(ctx, req.Username)
	if err != nil {
		return model.Users{}, errors.New("username or password incorrect")
	}

	err = utils.CheckPwd(req.Password, user.HashedPassword)
	if err != nil {
		return model.Users{}, errors.New("username or password incorrect")
	}

	return user, nil
}
