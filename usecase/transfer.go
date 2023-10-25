package usecase

import (
	"context"
	"fmt"

	"github.com/terajari/bank-api/dto"
	"github.com/terajari/bank-api/model"
	"github.com/terajari/bank-api/repository"
	"github.com/terajari/bank-api/utils"
)

type TransferUsecase interface {
	MakeTransfer(ctx context.Context, request dto.MakeTransferRequest) (dto.MakeTransferResponse, error)
}

type transferUsecase struct {
	accountRepo  repository.AccountsRepository
	entriesRepo  repository.EntryRepository
	transferRepo repository.TransferRepository
}

func NewTransferUsecase(acc repository.AccountsRepository, ent repository.EntryRepository, tr repository.TransferRepository) TransferUsecase {
	return &transferUsecase{accountRepo: acc, entriesRepo: ent, transferRepo: tr}
}

func (t *transferUsecase) MakeTransfer(ctx context.Context, request dto.MakeTransferRequest) (dto.MakeTransferResponse, error) {
	sender, err := t.accountRepo.Get(ctx, request.SenderId)
	if err != nil {
		return dto.MakeTransferResponse{}, err
	}
	receiver, err := t.accountRepo.Get(ctx, request.ReceiverId)
	if err != nil {
		return dto.MakeTransferResponse{}, err
	}

	if sender.Currency != receiver.Currency {
		return dto.MakeTransferResponse{}, fmt.Errorf("currency mismatch: %s != %s", sender.Currency, receiver.Currency)
	}
	if sender.Balance < request.Amount {
		return dto.MakeTransferResponse{}, fmt.Errorf("insufficient balance: %d < %d", sender.Balance, request.Amount)
	}

	transferId := utils.GenerateUUID()
	response, err := t.transferRepo.TransferTx(ctx, model.Transfer{
		ID:         transferId,
		SenderId:   request.SenderId,
		ReceiverId: request.ReceiverId,
		Amount:     request.Amount,
	})
	if err != nil {
		fmt.Println("error usecase")
		return dto.MakeTransferResponse{}, err
	}
	return response, nil
}
