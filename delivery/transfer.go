package delivery

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/terajari/bank-api/dto"
	"github.com/terajari/bank-api/middleware"
	"github.com/terajari/bank-api/model"
	"github.com/terajari/bank-api/token"
	"github.com/terajari/bank-api/usecase"
)

type TransferHandler struct {
	transferUsecase usecase.TransferUsecase
	accountUsecase  usecase.AccountsUsecase
	sessionsUsecase usecase.SessionsUsecase
}

func NewTransferHandler(tu usecase.TransferUsecase, au usecase.AccountsUsecase, su usecase.SessionsUsecase) (*TransferHandler, error) {
	return &TransferHandler{
		transferUsecase: tu,
		accountUsecase:  au,
		sessionsUsecase: su,
	}, nil
}

func (t *TransferHandler) performTransfer(ctx *gin.Context) {
	var req dto.MakeTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ls, err := t.sessionsUsecase.LastSession(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ls.IsBlocked {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "you must log in first"})
		return
	}

	sender, ok := t.validAccount(ctx, req.SenderId, req.Currency)
	if !ok {
		return
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)
	if sender.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "sender is not authorized to transfer"})
		return
	}

	_, ok = t.validAccount(ctx, req.ReceiverId, req.Currency)
	if !ok {
		return
	}

	resp, err := t.transferUsecase.MakeTransfer(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (a *TransferHandler) validAccount(ctx *gin.Context, accId, currency string) (model.Accounts, bool) {
	acc, err := a.accountUsecase.GetAccount(ctx, accId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return model.Accounts{}, false
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return model.Accounts{}, false
	}

	if acc.Currency != currency {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid currency"})
		return model.Accounts{}, false
	}

	return model.Accounts{
		ID:    acc.Id,
		Owner: acc.Owner,
	}, true
}
