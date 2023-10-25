package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/terajari/bank-api/dto"
	"github.com/terajari/bank-api/model"
	"github.com/terajari/bank-api/utils"
)

type TransferRepository interface {
	Create(ctx context.Context, transfer model.Transfer) (model.Transfer, error)
	Get(ctx context.Context, id string) (model.Transfer, error)
	List(ctx context.Context, transfer model.Transfer, limit, offset int) ([]model.Transfer, error)
	TransferTx(ctx context.Context, transfer model.Transfer) (dto.MakeTransferResponse, error)
}

type transferRepository struct {
	accRepo AccountsRepository
	entRepo EntryRepository
	db      *sqlx.DB
}

func NewTransferRepository(db *sqlx.DB) TransferRepository {
	return &transferRepository{
		accRepo: NewAccountsRepository(db),
		entRepo: NewEntryRepository(db),
		db:      db,
	}
}

func (t *transferRepository) Create(ctx context.Context, transfer model.Transfer) (model.Transfer, error) {
	query := "INSERT INTO transfer (id, sender_id, receiver_id, amount) VALUES ($1, $2, $3) RETURNING id, sender_id, receiver_id, amount, created_at"
	row := t.db.QueryRowContext(ctx, query, transfer.ID, transfer.SenderId, transfer.ReceiverId)
	var tr model.Transfer
	if err := row.Scan(&tr.ID, &tr.SenderId, &tr.ReceiverId, &tr.Amount, &tr.CreatedAt); err != nil {
		return model.Transfer{}, err
	}
	return tr, nil
}

func (t *transferRepository) Get(ctx context.Context, id string) (model.Transfer, error) {
	query := "SELECT id, sender_id, receiver_id, amount, created_at FROM transfer WHERE id = $1 LIMIT 1"
	row := t.db.QueryRowContext(ctx, query, id)
	var tr model.Transfer
	if err := row.Scan(&tr.ID, &tr.SenderId, &tr.ReceiverId, &tr.Amount, &tr.CreatedAt); err != nil {
		return model.Transfer{}, err
	}
	return tr, nil
}

func (tr *transferRepository) List(ctx context.Context, transfer model.Transfer, limit, offset int) ([]model.Transfer, error) {
	query := "SELECT  id, sender_id, receiver_id, amount, created_at FROM transfer WHERE sender_id = $1 OR receiver_id = $2 ORDER BY id LIMIT $3 OFFSET $4"
	rows, err := tr.db.QueryContext(ctx, query, transfer.SenderId, transfer.ReceiverId, limit, offset)
	if err != nil {
		return []model.Transfer{}, err
	}
	defer rows.Close()
	var transfers []model.Transfer
	for rows.Next() {
		var tr model.Transfer
		if err := rows.Scan(&tr.ID, &tr.SenderId, &tr.ReceiverId, &tr.Amount, &tr.CreatedAt); err != nil {
			return []model.Transfer{}, err
		}
		transfers = append(transfers, tr)
	}
	return transfers, nil
}

func (t *transferRepository) TransferTx(ctx context.Context, transfer model.Transfer) (dto.MakeTransferResponse, error) {
	var response dto.MakeTransferResponse
	tx, err := t.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		fmt.Println("error repo1")
		return dto.MakeTransferResponse{}, err
	}
	defer tx.Rollback()

	queryTranfer := "INSERT INTO transfers (id, sender_id, receiver_id, amount) VALUES ($1, $2, $3, $4) RETURNING id, sender_id, receiver_id, amount, created_at"
	row := tx.QueryRowContext(ctx, queryTranfer, transfer.ID, transfer.SenderId, transfer.ReceiverId, transfer.Amount)
	var tr model.Transfer
	if err := row.Scan(&tr.ID, &tr.SenderId, &tr.ReceiverId, &tr.Amount, &tr.CreatedAt); err != nil {
		fmt.Println("error repo2")
		return dto.MakeTransferResponse{}, err
	}
	response.Transfer = tr

	senderEntryId := utils.GenerateUUID()
	querySenderEntry := "INSERT INTO entries (id, account_id, amount) VALUES ($1, $2, $3) RETURNING id, account_id, amount, created_at"
	row = tx.QueryRowContext(ctx, querySenderEntry, senderEntryId, transfer.SenderId, -transfer.Amount)
	var senderEnt model.Entries
	if err := row.Scan(&senderEnt.ID, &senderEnt.AccountId, &senderEnt.Amount, &senderEnt.CreatedAt); err != nil {
		fmt.Println("error repo3")
		return dto.MakeTransferResponse{}, err
	}
	response.SenderEntry = senderEnt

	receiverEntryId := utils.GenerateUUID()
	queryReceiverEntry := "INSERT INTO entries (id, account_id, amount) VALUES ($1, $2, $3) RETURNING id, account_id, amount, created_at"
	row = tx.QueryRowContext(ctx, queryReceiverEntry, receiverEntryId, transfer.ReceiverId, transfer.Amount)
	var receiverEnt model.Entries
	if err := row.Scan(&receiverEnt.ID, &receiverEnt.AccountId, &receiverEnt.Amount, &receiverEnt.CreatedAt); err != nil {
		fmt.Println("error repo4")
		return dto.MakeTransferResponse{}, err
	}
	response.ReceiverEntry = receiverEnt

	querySenderUpdate := "SELECT id, owner, balance, currency FROM accounts WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE"
	_, err = tx.ExecContext(ctx, querySenderUpdate, transfer.SenderId)
	if err != nil {
		return dto.MakeTransferResponse{}, err
	}

	updateBalanceQuery := "UPDATE accounts SET balance = balance + $2 WHERE id = $1 RETURNING id, owner, balance, currency, created_at"

	row = tx.QueryRowContext(ctx, updateBalanceQuery, transfer.SenderId, -transfer.Amount)
	var senderAcc model.Accounts
	if err := row.Scan(&senderAcc.ID, &senderAcc.Owner, &senderAcc.Balance, &senderAcc.Currency, &senderAcc.CreatedAt); err != nil {
		return dto.MakeTransferResponse{}, err
	}
	response.Sender = senderAcc

	row = tx.QueryRowContext(ctx, updateBalanceQuery, transfer.ReceiverId, transfer.Amount)
	var receiverAcc model.Accounts
	if err := row.Scan(&receiverAcc.ID, &receiverAcc.Owner, &receiverAcc.Balance, &receiverAcc.Currency, &receiverAcc.CreatedAt); err != nil {
		return dto.MakeTransferResponse{}, err
	}
	response.Receiver = receiverAcc

	err = tx.Commit()
	if err != nil {
		fmt.Println("error repo5")
		return dto.MakeTransferResponse{}, err
	}

	return response, nil

}
