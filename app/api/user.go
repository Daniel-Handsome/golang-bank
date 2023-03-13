package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/daniel/master-golang/db/sqlc"
	"github.com/daniel/master-golang/utils"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)


type createUserRequest struct {
	Name    string `json:"name" binding:"required,alphanumunicode"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Name             string    `json:"name"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreateAt         time.Time `json:"create_at"`
}


func newUserResponse(db db.User) userResponse {
	return userResponse{
		Name:             db.Name,
        FullName:         db.FullName,
        Email:            db.Email,
        PasswordChangeAt: db.PasswordChangeAt,
        CreateAt:         db.CreateAt,
	}
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
    if err := ctx.ShouldBindJSON(&req); err!= nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	arg := db.CreateUserParams{
		Name : req.Name,    
		Password: hashPassword,
		FullName : req.FullName,
		Email : req.Email,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation" :
				ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
                return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
	}

	ctx.JSON(http.StatusCreated, newUserResponse(user))
	return
}

type loginUserRequest struct {
	Name    string `json:"name" binding:"required,alphanumunicode"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionId uuid.UUID `json:"session_id"`
	AccessToken string `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	RefreshToken string `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User userResponse `json:"user"`
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	user, err := s.store.GetUser(ctx, req.Name)
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

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"err": err.Error(),
		})
		return
	}

	// create token
	token, payload, err := s.tokenMaker.CreateToken(user.Name, s.config.Access_token_duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	// refresh token
	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.Name, s.config.Access_token_duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID : refreshPayload.ID,
		Username : user.Name,
		RefreshToken: refreshToken,
		UserAgent : ctx.Request.UserAgent(),
		ClientID : ctx.ClientIP(),
		Isblacked: false,
		ExpiresAt: refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return 
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		SessionId : session.ID,
		AccessToken: token,
		AccessTokenExpiresAt: payload.ExpiredAt,
		RefreshToken: refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: newUserResponse(user),
	})
}