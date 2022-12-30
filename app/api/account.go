package api

import (
	"database/sql"
	"net/http"

	db "github.com/daniel/master-golang/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)


type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
    if err := ctx.ShouldBindJSON(&req); err!= nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	//用db抓 不要s.store 這是給query用
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
        Currency: req.Currency,
		Balance: 0,
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); !ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation" :
				ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
                return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
        return
	}

	ctx.JSON(http.StatusCreated, account)
	return
}

// uri 是代表資源的意思
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`	
}

func (s *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}
	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		if 	err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

//form  給query string
type getAccountsRequest struct {
	Page int32 `form:"page" binding:"required,min=1"` 
	PerPage int32 `form:"per_page" binding:"required"`
}

func (s *Server) getAccounts(ctx *gin.Context) {
	var req getAccountsRequest
	
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}

	accounts, err := s.store.GetAccounts(ctx, db.GetAccountsParams{
		Limit: req.PerPage,
		Offset: (req.Page - 1) * req.PerPage,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
	}

	ctx.JSON(http.StatusOK, accounts)
}