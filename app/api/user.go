package api

import (
	"net/http"
	"time"

	db "github.com/daniel/master-golang/db/sqlc"
	"github.com/daniel/master-golang/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)


type createUserRequest struct {
	Name    string `json:"name" binding:"required,alphanumunicode"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Name             string    `json:"name"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreateAt         time.Time `json:"create_at"`
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

	ctx.JSON(http.StatusCreated, createUserResponse{
		Name:             user.Name,
        FullName:         user.FullName,
        Email:            user.Email,
        PasswordChangeAt: user.PasswordChangeAt,
        CreateAt:         user.CreateAt,
	})
	return
}