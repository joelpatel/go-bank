package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/joelpatel/go-bank/db/sqlc"
	"github.com/joelpatel/go-bank/util"
)

type transferRequest struct {
	FromAccountID int64   `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64   `json:"to_account_id" binding:"required,min=1"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	// Currency      string  `json:"currency" binding:"required,oneof=USD INR EUR"` // currency of the money we want to transfer

}

func (server *Server) createTranfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validateTransfer(ctx, req.FromAccountID, req.ToAccountID, req.Amount) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validateTransfer(ctx *gin.Context, fromAccountID int64, toAccountID int64, amount float64) bool {
	account1, err := server.store.GetAccount(ctx, fromAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	account2, err := server.store.GetAccount(ctx, toAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account1.Currency != account2.Currency {
		err = fmt.Errorf("currency mismatch: %v -> %v", account1.Currency, account2.Currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	// check -ve.
	if util.RoundFloat(account1.Balance-amount, db.PRECISION) < 0 {
		err = fmt.Errorf("account %d does not have enough funds to transfer money", fromAccountID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true

}
