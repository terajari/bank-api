package manager

import "github.com/terajari/bank-api/repository"

type RepositoryManager interface {
	AccountsRepo() repository.AccountsRepository
	EntryRepo() repository.EntryRepository
	TransferRepo() repository.TransferRepository
	UsersRepo() repository.UsersRepository
}

type repositoryManager struct {
	infra InfrastuctureManager
}

func (r *repositoryManager) AccountsRepo() repository.AccountsRepository {
	return repository.NewAccountsRepository(r.infra.Conn())
}

func (r *repositoryManager) EntryRepo() repository.EntryRepository {
	return repository.NewEntryRepository(r.infra.Conn())
}

func (r *repositoryManager) TransferRepo() repository.TransferRepository {
	return repository.NewTransferRepository(r.infra.Conn())
}

func (r *repositoryManager) UsersRepo() repository.UsersRepository {
	return repository.NewUsersRepository(r.infra.Conn())
}

func NewRepositoryManager(infra InfrastuctureManager) (RepositoryManager, error) {
	return &repositoryManager{
		infra: infra,
	}, nil
}
