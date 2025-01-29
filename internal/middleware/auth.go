package middleware

import (
	"context"
	"library-api-book/internal/commons/response"
	"library-api-book/internal/grpc/client"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckAuth(authClient *client.AuthClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		bearerToken := strings.Split(header, "Bearer ")

		if len(bearerToken) != 2 {
			resp := response.UnauthorizedErrorWithAdditionalInfo("len token must be 2")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		valid, payload := authClient.ValidateToken(context.Background(), bearerToken[1])
		if !valid {
			resp := response.UnauthorizedErrorWithAdditionalInfo("Invalid token")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		ctx.Set("authId", payload.AuthId)
		ctx.Set("role", payload.Role)
		ctx.Next()
	}
}

func CheckAuthIsAdmin(authClient *client.AuthClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		bearerToken := strings.Split(header, "Bearer ")

		if len(bearerToken) != 2 {
			resp := response.UnauthorizedErrorWithAdditionalInfo("len token must be 2")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		valid, payload := authClient.ValidateToken(context.Background(), bearerToken[1])
		if !valid {
			resp := response.UnauthorizedErrorWithAdditionalInfo("Invalid token")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		if payload.Role != "admin" {
			resp := response.UnauthorizedErrorWithAdditionalInfo("user doesn't have permission to access")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		ctx.Set("authId", payload.AuthId)
		ctx.Set("role", payload.Role)
		ctx.Next()
	}
}

func CheckAuthIsAdminOrAuthor(authClient *client.AuthClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		bearerToken := strings.Split(header, "Bearer ")

		if len(bearerToken) != 2 {
			resp := response.UnauthorizedErrorWithAdditionalInfo("len token must be 2")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		valid, payload := authClient.ValidateToken(context.Background(), bearerToken[1])
		if !valid {
			resp := response.UnauthorizedErrorWithAdditionalInfo("Invalid token")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		if payload.Role == "user" {
			resp := response.UnauthorizedErrorWithAdditionalInfo("user doesn't have permission to access")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}
		ctx.Set("authId", payload.AuthId)
		ctx.Set("role", payload.Role)
		ctx.Next()
	}
}
