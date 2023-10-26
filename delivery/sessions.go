package delivery

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/terajari/bank-api/dto"
	"github.com/terajari/bank-api/token"
	"github.com/terajari/bank-api/usecase"
	"github.com/terajari/bank-api/utils"
)

type SessionsHandler struct {
	sessionUsecase usecase.SessionsUsecase
	tokenMaker     token.Maker
	config         *utils.Config
}

func NewSessionHandler(su usecase.SessionsUsecase, tm token.Maker, cfg utils.Config) (*SessionsHandler, error) {
	return &SessionsHandler{
		sessionUsecase: su,
		tokenMaker:     tm,
		config:         &cfg,
	}, nil
}

func (sessionHandler *SessionsHandler) renewHandler(
	ctx *gin.Context) {
	var req dto.RenewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	refreshPayload, err := sessionHandler.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	}

	session, err := sessionHandler.sessionUsecase.GetSessions(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(
			http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(
			http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatch session token")
		ctx.JSON(
			http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(
			http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		return
	}

	accessToken, accessPayload, err := sessionHandler.tokenMaker.CreateToken(
		refreshPayload.Username, sessionHandler.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	rsp := dto.RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
