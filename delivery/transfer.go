package delivery

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/terajari/bank-api/dto"
	"github.com/terajari/bank-api/usecase"
)

type TransferHandler struct {
	usecase usecase.TransferUsecase
}

func NewTransferHandler(uc usecase.TransferUsecase) (*TransferHandler, error) {
	return &TransferHandler{
		usecase: uc,
	}, nil
}

func (t *TransferHandler) performTransfer(ctx *gin.Context) {
	var req dto.MakeTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("error deliver1")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := t.usecase.MakeTransfer(ctx, req)
	if err != nil {
		fmt.Println("error deliver2")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
