package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/terajari/bank-api/manager"
	"github.com/terajari/bank-api/middleware"
	"github.com/terajari/bank-api/token"
	"github.com/terajari/bank-api/utils"
	validators "github.com/terajari/bank-api/utils/validator"
)

type Server struct {
	AccountsHandler *AccountsHandler
	TransferHandler *TransferHandler
	UsersHandler    *UsersHandler
	SessionsHandler *SessionsHandler
	UsecaseManager  *manager.UsecaseManager
	Router          *gin.Engine
	Config          utils.Config
	TokenMaker      token.Maker
}

func NewServer(config utils.Config, usecase manager.UsecaseManager) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmtricKey)
	if err != nil {
		return nil, err
	}

	sessionsHandler, err := NewSessionHandler(usecase.SessionsUsecase(), tokenMaker, config)
	if err != nil {
		return nil, err
	}

	accHandler, err := NewAccountsHandler(usecase.AccountsUsecase(), usecase.SessionsUsecase())
	if err != nil {
		return nil, err
	}

	trfHandler, err := NewTransferHandler(usecase.TransferUsecase(), usecase.AccountsUsecase(), usecase.SessionsUsecase())
	if err != nil {
		return nil, err
	}

	usersHandler, err := NewUsersHandler(usecase.UsersUsecase(), usecase.SessionsUsecase(), tokenMaker, &config)
	if err != nil {
		return nil, err
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validators.ValidCurrency)
	}

	server := &Server{
		AccountsHandler: accHandler,
		TransferHandler: trfHandler,
		UsersHandler:    usersHandler,
		SessionsHandler: sessionsHandler,
		UsecaseManager:  &usecase,
		Router:          gin.Default(),
		Config:          config,
		TokenMaker:      tokenMaker,
	}
	return server, nil
}

func (s *Server) SetupRouter() {
	router := gin.Default()
	router.POST("/user", s.UsersHandler.createHandler)
	router.POST("/user/login", s.UsersHandler.loginHandler)
	router.POST("/token/renew", s.SessionsHandler.renewHandler)

	authRoute := router.Group("/").Use(middleware.AuthMiddleware(s.TokenMaker))
	authRoute.POST("/account", s.AccountsHandler.createHandler)
	authRoute.GET("/account/:id", s.AccountsHandler.getHandler)
	authRoute.GET("/account/", s.AccountsHandler.listHandlers)
	authRoute.POST("/user/logout", s.UsersHandler.logoutHandler)

	authRoute.POST("/transfer", s.TransferHandler.performTransfer)
	s.Router = router
}

func (s *Server) Start(address string) error {
	return s.Router.Run(address)
}
