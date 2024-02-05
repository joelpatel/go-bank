package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joelpatel/go-bank/currency"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var request createAccountRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !currency.IsSupportedCurrency(request.Currency) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is an unsupported currency.", request.Currency)})
		return
	}

	createdAccount, err := server.store.CreateAccount(ctx, request.Owner, 0, request.Currency)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createdAccount)
}

type getAccountByIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccountByID(ctx *gin.Context) {
	var request getAccountByIDRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccountByID(ctx, request.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Account with id %d not found.", request.ID)})
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsByOwnerRequestQuery struct {
	PageID   int64 `form:"page_id" binding:"required,min=1"`
	PageSize int64 `form:"page_size" binding:"required,min=1,max=100"`
}

type listAccountsByOwnerRequestJSON struct {
	Owner string `json:"owner" binding:"required"`
}

func (server *Server) listAccountsByOwner(ctx *gin.Context) {
	var requestQueryParam listAccountsByOwnerRequestQuery
	var requestJSON listAccountsByOwnerRequestJSON

	if err := ctx.ShouldBindJSON(&requestJSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&requestQueryParam); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	accounts, err := server.store.ListAccounts(ctx, requestJSON.Owner, requestQueryParam.PageSize, requestQueryParam.PageSize*(requestQueryParam.PageID-1))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if accounts == nil || len(*accounts) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Accounts from entry %d not found.", requestQueryParam.PageSize*(requestQueryParam.PageID-1))})
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type updateAccountOwnerRequest struct {
	ID       int64  `json:"id" binding:"required,min=1"`
	NewOwner string `json:"new_owner" binding:"required"`
}

func (server *Server) updateAccountOwner(ctx *gin.Context) {
	var request updateAccountOwnerRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	rowsAffected, err := server.store.UpdateAccountOwner(ctx, request.ID, request.NewOwner)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if rowsAffected != 1 {
		ctx.Status(http.StatusNotModified)
	} else {
		ctx.Status(http.StatusNoContent)
	}
}

type deleteAccountByIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccountByID(ctx *gin.Context) {
	var request deleteAccountByIDRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rowsAffected, err := server.store.DeleteAccountByID(ctx, request.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if rowsAffected != 1 {
		ctx.Status(http.StatusNotFound)
	} else {
		ctx.Status(http.StatusNoContent)
	}
}
