package routes

import (
	"fmt"
	"library-api-book/internal/commons/response"
	"library-api-book/internal/factory"
	"library-api-book/internal/grpc/client"
	"library-api-book/internal/middleware"
	"library-api-book/pkg/token"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(provider *factory.Provider, authClient *client.AuthClient) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger(), CORS())

	router.GET("/", func(ctx *gin.Context) {
		currentYear := time.Now().Year()
		message := fmt.Sprintf("Library API Book %d", currentYear)

		ctx.JSON(http.StatusOK, message)
	})

	api := router.Group("/api")
	{
		v1 := api.Group("v1")
		{
			v1.GET("/books", provider.BookProvider.GetAllBooks)
			v1.GET("/books/:id", provider.BookProvider.GetDetailBook)
			auth := v1.Use(middleware.CheckAuth(authClient))
			auth.POST("/books", provider.BookProvider.CreateBook)
			auth.PUT("/books/:id", provider.BookProvider.UpdateBook)
			auth.DELETE("/books/id", provider.BookProvider.DeleteBook)
			auth.GET("/books/recommendation", provider.BookProvider.GetRecommendationBook)
		}
	}

	return router
}

func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, accept, access-control-allow-origin, access-control-allow-headers")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		ctx.Next()
	}
}

func CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		bearerToken := strings.Split(header, "Bearer ")

		if len(bearerToken) != 2 {
			resp := response.UnauthorizedErrorWithAdditionalInfo("len token must be 2")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		payload, err := token.ValidateToken(bearerToken[1])
		if err != nil {
			resp := response.UnauthorizedErrorWithAdditionalInfo(err.Error())
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}
		ctx.Set("authId", payload.AuthId)
		ctx.Next()
	}
}
