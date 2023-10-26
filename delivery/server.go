package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/terajari/bank-api/manager"
	"github.com/terajari/bank-api/token"
	"github.com/terajari/bank-api/utils"
	validators "github.com/terajari/bank-api/utils/validator"
)

type Server struct {
	AccountsHandler *AccountsHandler
	TransferHandler *TransferHandler
	UsersHandler    *UsersHandler
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

	accHandler, err := NewAccountsHandler(usecase.AccountsUsecase())
	if err != nil {
		return nil, err
	}

	trfHandler, err := NewTransferHandler(usecase.TransferUsecase())
	if err != nil {
		return nil, err
	}

	usersHandler, err := NewUsersHandler(usecase.UsersUsecase(), tokenMaker, &config)
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
		UsecaseManager:  &usecase,
		Router:          gin.Default(),
		Config:          config,
		TokenMaker:      tokenMaker,
	}
	return server, nil
}

func (s *Server) SetupRouter() {
	router := gin.Default()
	router.POST("/account", s.AccountsHandler.createHandler)
	router.GET("/account/:id", s.AccountsHandler.getHandler)

	router.POST("/user", s.UsersHandler.createHandler)
	router.POST("/user/login", s.UsersHandler.loginHandler)

	router.POST("/transfer", s.TransferHandler.performTransfer)
	s.Router = router
}

func (s *Server) Start(address string) error {
	return s.Router.Run(address)
}
