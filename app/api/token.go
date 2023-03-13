package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}


func (s *Server) RenewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
    if err := ctx.ShouldBindJSON(&req); err!= nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// get session
	payload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	session, err := s.store.GetSession(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"err": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	if session.Isblacked {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"err": fmt.Errorf("blocked sesssion"),
		})
	}

	// check information
	if session.Username != payload.UserName {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"err": fmt.Errorf("incorrect username"),
		})
	}

	if session.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"err": fmt.Errorf("mismatch sesssion"),
		})
	}

	// create token
	token, payload, err := s.tokenMaker.CreateToken(payload.UserName, s.config.Access_token_duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}


	ctx.JSON(http.StatusOK, renewAccessTokenResponse{
		AccessToken: token,
		AccessTokenExpiresAt: payload.ExpiredAt, 
	})
}
