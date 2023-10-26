package manager

import "github.com/terajari/bank-api/usecase"

type UsecaseManager interface {
	AccountsUsecase() usecase.AccountsUsecase
	TransferUsecase() usecase.TransferUsecase
	UsersUsecase() usecase.UsersUsecase
	SessionsUsecase() usecase.SessionsUsecase
}

type usecaseManager struct {
	Repository RepositoryManager
}

func (u *usecaseManager) AccountsUsecase() usecase.AccountsUsecase {
	return usecase.NewAccountsUsecase(u.Repository.AccountsRepo())
}

func (u *usecaseManager) TransferUsecase() usecase.TransferUsecase {
	return usecase.NewTransferUsecase(u.Repository.AccountsRepo(), u.Repository.EntryRepo(), u.Repository.TransferRepo())
}

func (u *usecaseManager) UsersUsecase() usecase.UsersUsecase {
	return usecase.NewUsersUsecase(u.Repository.UsersRepo())
}

func (u *usecaseManager) SessionsUsecase() usecase.SessionsUsecase {
	return usecase.NewSessionsUsecase(u.Repository.SessionsRepo())
}

func NewUsecaseManager(repositoryManager RepositoryManager) (UsecaseManager, error) {
	return &usecaseManager{
		Repository: repositoryManager,
	}, nil
}
