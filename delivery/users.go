package delivery

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/terajari/bank-api/dto"
	"github.com/terajari/bank-api/token"
	"github.com/terajari/bank-api/usecase"
	"github.com/terajari/bank-api/utils"
)

type UsersHandler struct {
	usecase    usecase.UsersUsecase
	sessions   usecase.SessionsUsecase
	tokenMaker token.Maker
	config     *utils.Config
}

func NewUsersHandler(uc usecase.UsersUsecase, ss usecase.SessionsUsecase, token token.Maker, cfg *utils.Config) (*UsersHandler, error) {
	return &UsersHandler{
		usecase:    uc,
		sessions:   ss,
		tokenMaker: token,
		config:     cfg,
	}, nil
}

func (u *UsersHandler) createHandler(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := u.usecase.CreateUser(ctx, req)
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, response)
}

func (u *UsersHandler) loginHandler(ctx *gin.Context) {
	var req dto.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.usecase.Login(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessToken, Accesspayload, err := u.tokenMaker.CreateToken(req.Username, u.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, refreshPayload, err := u.tokenMaker.CreateToken(user.Username, u.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ss, err := u.sessions.AddSessions(ctx, dto.AddSessionsRequest{
		Id:           refreshPayload.ID,
		Username:     refreshPayload.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userRes := dto.UserReponse{
		Username:     user.Username,
		FullName:     user.FullName,
		Email:        user.Email,
		PwdChangedAt: user.PasswordChangedAt,
		CreatedAt:    user.CreatedAt,
	}

	response := dto.LoginUserResponse{
		SessionsId:            ss.Id,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  Accesspayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  userRes,
	}

	ctx.JSON(http.StatusOK, response)
}

func (u *UsersHandler) logoutHandler(ctx *gin.Context) {
	ls, err := u.sessions.LastSession(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if ls.IsBlocked {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "you already logged out"})
		return
	}

	err = u.sessions.UpdateBlockStatus(ctx, dto.UpdateSessionBlockRequest{
		Id:        ls.Id,
		IsBlocked: true,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully logout",
	})
}
