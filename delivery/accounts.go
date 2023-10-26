package delivery

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/terajari/bank-api/dto"
	"github.com/terajari/bank-api/usecase"
)

type AccountsHandler struct {
	usecase usecase.AccountsUsecase
}

func NewAccountsHandler(uc usecase.AccountsUsecase) (*AccountsHandler, error) {
	return &AccountsHandler{
		usecase: uc,
	}, nil
}

func (a *AccountsHandler) createHandler(ctx *gin.Context) {
	var req dto.RegisterNewAccountsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := a.usecase.RegisterNewAccounts(ctx, req)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, gin.H{
					"error": errors.New(pqErr.Message),
				})
				return
			}
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (a *AccountsHandler) getHandler(ctx *gin.Context) {
	var req dto.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := a.usecase.GetAccount(ctx, req.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
