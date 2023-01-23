package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/daniel/master-golang/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationTypeBearer = "Bearer"
	authorizationPayloadKey = "authorization_payload_key"
)

func authMiddleware(tokenMaker token.Marker) gin.HandlerFunc {
	return func (ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		fields := strings.Fields(authorizationHeader)

		authorizationType := fields[0]
		if (len(fields) != 2) || (authorizationType != authorizationTypeBearer) {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		accessToken := fields[1]
		payload, err :=tokenMaker.VerifyToken(accessToken)
		if err!= nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
                "error": err.Error(),
            })
            return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
