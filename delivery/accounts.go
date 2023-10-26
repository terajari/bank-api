package delivery

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/terajari/bank-api/dto"
	"github.com/terajari/bank-api/middleware"
	"github.com/terajari/bank-api/token"
	"github.com/terajari/bank-api/usecase"
)

type AccountsHandler struct {
	sessionsUsecase usecase.SessionsUsecase
	usecase         usecase.AccountsUsecase
}

func NewAccountsHandler(uc usecase.AccountsUsecase, ss usecase.SessionsUsecase) (*AccountsHandler, error) {
	return &AccountsHandler{
		sessionsUsecase: ss,
		usecase:         uc,
	}, nil
}

type reqCreate struct {
	Currency string `json:"currency" binding:"required"`
}

func (a *AccountsHandler) createHandler(ctx *gin.Context) {
	var req reqCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ls, err := a.sessionsUsecase.LastSession(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ls.IsBlocked {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "you must login first"})
		return
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

	resp, err := a.usecase.RegisterNewAccounts(ctx, dto.RegisterNewAccountsRequest{
		Owner:    authPayload.Username,
		Currency: req.Currency,
	})
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

	ls, err := a.sessionsUsecase.LastSession(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ls.IsBlocked {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "you must login first"})
		return
	}

	resp, err := a.usecase.GetAccount(ctx, req.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)
	if resp.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user is not authorized to access this"})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (a *AccountsHandler) listHandlers(ctx *gin.Context) {
	var req dto.ListAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ls, err := a.sessionsUsecase.LastSession(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ls.IsBlocked {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "you must login first"})
		return
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)
	req.Owner = authPayload.Username

	resp, err := a.usecase.ListAccounts(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
