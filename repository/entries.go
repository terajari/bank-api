package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/terajari/bank-api/model"
)

type EntryRepository interface {
	Create(ctx context.Context, entry model.Entries) (model.Entries, error)
	Get(ctx context.Context, id string) (model.Entries, error)
	List(ctx context.Context, accountId string, limit, offset int) ([]model.Entries, error)
}

type entryRepository struct {
	db *sqlx.DB
}

func NewEntryRepository(db *sqlx.DB) EntryRepository {
	return &entryRepository{db: db}
}

func (r *entryRepository) Create(ctx context.Context, entry model.Entries) (model.Entries, error) {
	query := "INSERT INTO entries (id, account_id, amount, created_at) VALUES ($1, $2, $3, $4) RETURNING id, account_id, amount, created_at"

	row := r.db.QueryRowContext(ctx, query, entry.ID, entry.AccountId, entry.Amount, entry.CreatedAt)
	var e model.Entries
	if err := row.Scan(&e.ID, &e.AccountId, &e.Amount, &e.CreatedAt); err != nil {
		return model.Entries{}, err
	}

	return e, nil
}

func (r *entryRepository) Get(ctx context.Context, id string) (model.Entries, error) {
	query := "SELECT id, account_id, amount, created_at FROM entries WHERE id = $1 LIMIT 1"

	row := r.db.QueryRowContext(ctx, query, id)
	var e model.Entries
	if err := row.Scan(&e.ID, &e.AccountId, &e.Amount, &e.CreatedAt); err != nil {
		return model.Entries{}, err
	}

	return e, nil
}

func (r *entryRepository) List(ctx context.Context, accountId string, limit, offset int) ([]model.Entries, error) {
	query := "SELECT id, account_id, amount, created_at FROM entries WHERE account_id = $1 ORDER BY id limit = $2 offset = $3"
	rows, err := r.db.QueryContext(ctx, query, accountId, limit, offset)
	if err != nil {
		return []model.Entries{}, err
	}
	defer rows.Close()
	var entries []model.Entries
	for rows.Next() {
		var e model.Entries
		if err := rows.Scan(&e.ID, &e.AccountId, &e.Amount, &e.CreatedAt); err != nil {
			return []model.Entries{}, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
