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
	tokenMaker token.Maker
	config     *utils.Config
}

func NewUsersHandler(uc usecase.UsersUsecase, token token.Maker, cfg *utils.Config) (*UsersHandler, error) {
	return &UsersHandler{
		usecase:    uc,
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

	accessToken, payload, err := u.tokenMaker.CreateToken(req.Username, u.config.AccessTokenDuration)
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
		AccessToken:          accessToken,
		AccessTokenExpiresAt: payload.ExpiredAt,
		User:                 userRes,
	}

	ctx.JSON(http.StatusOK, response)
}
