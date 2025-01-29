package routes

import (
	"fmt"
	"library-api-book/internal/factory"
	"library-api-book/internal/grpc/client"
	"library-api-book/internal/middleware"
	"net/http"
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
			auth := v1.Use(middleware.CheckAuth(authClient))
			auth.GET("/books", provider.BookProvider.GetAllBooks)
			auth.GET("/books/:id", provider.BookProvider.GetDetailBook)
			auth.GET("/books/recommendation", provider.BookProvider.GetRecommendationBook)

			admin := v1.Use(middleware.CheckAuthIsAdminOrAuthor(authClient))
			admin.POST("/books", provider.BookProvider.CreateBook)
			admin.PUT("/books/:id", provider.BookProvider.UpdateBook)
			admin.DELETE("/books/id", provider.BookProvider.DeleteBook)
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
