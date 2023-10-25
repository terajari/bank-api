package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/terajari/bank-api/model"
)

type AccountsRepository interface {
	Create(ctx context.Context, account model.Accounts) (model.Accounts, error)
	Get(ctx context.Context, id string) (model.Accounts, error)
	List(ctx context.Context, owner string, limit, offset int) ([]model.Accounts, error)
	Update(ctx context.Context, account model.Accounts) (model.Accounts, error)
	Delete(ctx context.Context, id string) error
	GetForUpdate(ctx context.Context, id string) (model.Accounts, error)
}

type accountsRepository struct {
	db *sqlx.DB
}

func NewAccountsRepository(db *sqlx.DB) AccountsRepository {
	return &accountsRepository{db: db}
}

func (r *accountsRepository) Create(ctx context.Context, account model.Accounts) (model.Accounts, error) {
	query := "INSERT INTO accounts (id, owner, balance, currency) VALUES ($1, $2, $3, $4) RETURNING id, owner, balance, currency, created_at"

	row := r.db.QueryRowContext(ctx, query, account.ID, account.Owner, account.Balance, account.Currency)
	var a model.Accounts
	if err := row.Scan(&a.ID, &a.Owner, &a.Balance, &a.Currency, &a.CreatedAt); err != nil {
		return model.Accounts{}, err
	}

	return a, nil
}

func (r *accountsRepository) Get(ctx context.Context, id string) (model.Accounts, error) {
	query := "SELECT id, owner, balance, currency FROM accounts WHERE id = $1 LIMIT 1"

	row := r.db.QueryRowContext(ctx, query, id)
	var a model.Accounts
	if err := row.Scan(&a.ID, &a.Owner, &a.Balance, &a.Currency); err != nil {
		return model.Accounts{}, err
	}

	return a, nil
}

func (r *accountsRepository) List(ctx context.Context, id string, limit, offset int) ([]model.Accounts, error) {
	query := "SELECT id, owner, balance, currency FROM accounts WHERE owner = $1 LIMIT $2 OFFSET $3 ORDER BY id"

	rows, err := r.db.QueryContext(ctx, query, id, limit, offset)
	if err != nil {
		return []model.Accounts{}, err
	}
	defer rows.Close()
	var accounts []model.Accounts
	for rows.Next() {
		var a model.Accounts
		if err := rows.Scan(&a.ID, &a.Owner, &a.Balance, &a.Currency); err != nil {
			return []model.Accounts{}, err
		}
		accounts = append(accounts, a)
	}

	return accounts, nil
}

func (r *accountsRepository) Update(ctx context.Context, account model.Accounts) (model.Accounts, error) {
	query := "UPDATE accounts SET balance = $2 WHERE id = $1 RETURNING id, owner, balance, currency, created_at"
	row := r.db.QueryRowContext(ctx, query, account.ID, account.Balance)
	var a model.Accounts
	if err := row.Scan(&a.ID, &a.Owner, &a.Balance, &a.Currency, &a.CreatedAt); err != nil {
		return model.Accounts{}, err
	}
	return a, nil
}

func (r *accountsRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM accounts WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil

}

func (r *accountsRepository) GetForUpdate(ctx context.Context, id string) (model.Accounts, error) {
	query := "SELECT id, owner, balance, currency FROM accounts WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE"

	row := r.db.QueryRowContext(ctx, query, id)
	var a model.Accounts
	if err := row.Scan(&a.ID, &a.Owner, &a.Balance, &a.Currency); err != nil {
		return model.Accounts{}, err
	}

	return a, nil
}
