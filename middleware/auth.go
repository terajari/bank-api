package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/terajari/bank-api/token"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := "authorization header is not provided"
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := "invalid authorization header format"
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		authrizationType := strings.ToLower(fields[0])
		if authrizationType != AuthorizationTypeBearer {
			err := "incorect authorization type"
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
