package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/daniel/master-golang/db/sqlc"
	"github.com/daniel/master-golang/token"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var request createTransferRequest

	if err := ctx.BindJSON(&request); err!= nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}

	// validated request
	account, ok := s.IsAccountsCurrentMatch(ctx, request.FromAccountID, request.Currency)
	if  !ok {
		return
	}


	if  payload :=ctx.MustGet(authorizationPayloadKey).(*token.Payload); account.Owner != payload.UserName {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":  errors.New("the account is not belong to you "),
		})
	}

	account, ok = s.IsAccountsCurrentMatch(ctx, request.ToAccountID, request.Currency)
	if  !ok{
		return
	}

	transfer, err := s.store.TransferTx(ctx,  db.TransferTxParams{
		FromAccountID: request.FromAccountID,
        ToAccountID:   request.ToAccountID,
        Amount:        request.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
	}

	ctx.JSON(http.StatusOK, transfer)
}

func (s *Server) IsAccountsCurrentMatch(ctx *gin.Context, accountID int64, currentName string) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, accountID)

	if err != nil {
		if 	err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return account, false
	}

	// match
	if account.Currency != currentName {
		err := fmt.Errorf("account [%d] misMatch : %s vs %s", account.ID, account.Currency, currentName)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return account, false
	}

	return account, true
}