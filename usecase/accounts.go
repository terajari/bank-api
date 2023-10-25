package usecase

import (
	"context"

	"github.com/terajari/bank-api/dto"
	"github.com/terajari/bank-api/model"
	"github.com/terajari/bank-api/repository"
	"github.com/terajari/bank-api/utils"
)

type AccountsUsecase interface {
	RegisterNewAccounts(ctx context.Context, req dto.RegisterNewAccountsRequest) (dto.RegisterNewAccountsResponse, error)
	GetAccount(ctx context.Context, id string) (dto.GetAccountResponse, error)
	ListAccounts(ctx context.Context, req dto.ListAccountsRequest) ([]dto.GetAccountResponse, error)
	UpdateAccount(ctx context.Context, req dto.UpdateAccountRequest) (dto.UpdateAccountResponse, error)
	DeleteAccount(ctx context.Context, id string) error
}

type accountsUsecase struct {
	repo repository.AccountsRepository
}

func NewAccountsUsecase(repo repository.AccountsRepository) AccountsUsecase {
	return &accountsUsecase{repo: repo}
}

func (a *accountsUsecase) RegisterNewAccounts(ctx context.Context, req dto.RegisterNewAccountsRequest) (dto.RegisterNewAccountsResponse, error) {
	id := utils.GenerateUUID()
	account, err := a.repo.Create(ctx, model.Accounts{
		ID:       id,
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	})
	if err != nil {
		return dto.RegisterNewAccountsResponse{}, err
	}
	return dto.RegisterNewAccountsResponse{
		Id:        account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}, nil
}

func (a *accountsUsecase) GetAccount(ctx context.Context, id string) (dto.GetAccountResponse, error) {
	account, err := a.repo.Get(ctx, id)
	if err != nil {
		return dto.GetAccountResponse{}, err
	}
	return dto.GetAccountResponse{
		Id:       account.ID,
		Owner:    account.Owner,
		Balance:  account.Balance,
		Currency: account.Currency,
	}, nil
}

func (a *accountsUsecase) ListAccounts(ctx context.Context, req dto.ListAccountsRequest) ([]dto.GetAccountResponse, error) {
	size := req.Size
	if size == 0 {
		size = 5
	}
	page := (req.Page - 1) * size
	accounts, err := a.repo.List(ctx, req.Owner, size, page)
	if err != nil {
		return []dto.GetAccountResponse{}, err
	}
	var accountsDto []dto.GetAccountResponse
	for _, account := range accounts {
		accountsDto = append(accountsDto, dto.GetAccountResponse{
			Id:       account.ID,
			Owner:    account.Owner,
			Balance:  account.Balance,
			Currency: account.Currency,
		})
	}
	return accountsDto, nil
}

func (a *accountsUsecase) UpdateAccount(ctx context.Context, req dto.UpdateAccountRequest) (dto.UpdateAccountResponse, error) {
	acc, err := a.repo.Get(ctx, req.Id)
	if err != nil {
		return dto.UpdateAccountResponse{}, err
	}
	req.Id = acc.ID

	updatedAccount, err := a.repo.Update(ctx, model.Accounts{
		ID:      req.Id,
		Balance: req.Balance,
	})
	if err != nil {
		return dto.UpdateAccountResponse{}, err
	}
	return dto.UpdateAccountResponse{
		Id:       updatedAccount.ID,
		Owner:    updatedAccount.Owner,
		Balance:  updatedAccount.Balance,
		Currency: updatedAccount.Currency,
	}, nil
}

func (a *accountsUsecase) DeleteAccount(ctx context.Context, id string) error {
	acc, err := a.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	return a.repo.Delete(ctx, acc.ID)
}
