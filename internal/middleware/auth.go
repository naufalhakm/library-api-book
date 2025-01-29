package middleware

import (
	"context"
	"library-api-book/internal/grpc/client"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckAuth(authClient *client.AuthClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		bearerToken := strings.Split(header, "Bearer ")

		if len(bearerToken) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		valid, payload := authClient.ValidateToken(context.Background(), bearerToken[1])
		if !valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		ctx.Set("authId", payload.AuthId)
		ctx.Set("role", payload.Role)
		ctx.Next()
	}
}
